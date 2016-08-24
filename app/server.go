package app

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"splash/logger"
	"splash/processing"
	jobsLib "splash/queue/jobs"
	"splash/services"
)

type Server struct {
	config   map[string]string
	operator *processing.Operator
	logger   *logger.Logger
}

func NewServer(config map[string]string, operator *processing.Operator, logger *logger.Logger) *Server {
	return &Server{
		config:   config,
		operator: operator,
		logger:   logger,
	}
}

func (self *Server) Start() {
	listenAddr := self.config["host"] + self.config["port"]

	http.HandleFunc(self.config["path"], self.handler())

	self.logger.Alert("HTTP server:", http.ListenAndServe(listenAddr, nil))
}

func (self *Server) ServeHTTP(w http.ResponseWriter, r *http.Request) {

	authorizationToken := r.Header["Authorization"][0]

	serviceLocator := services.NewLocator()

	collection := PatchesCollection[authorizationToken]

	if collection == nil {
		collection = jobsLib.NewPatchCollection(serviceLocator.RandString("Collection[patch]-", 55))
		PatchesCollection[authorizationToken] = collection
	}

	bodyData, err := ioutil.ReadAll(r.Body)

	if err != nil {
		self.logger.Error(err.Error())
		return
	}

	data, err := self.operator.EnumerateData(bodyData)

	if err != nil {
		self.logger.Error(err.Error())
		return
	}

	patch, err := self.operator.EnumeratePatch(collection, data)

	fmt.Println(PatchesCollection[authorizationToken].GetLength())
	if err != nil {
		self.logger.Error(err.Error())
		return
	}

	processing.Total += patch.GetLength()

	jobsLib.PushToChanel(patch, jobsLib.JobsQueue)

	self.respond(w)
}

func (self *Server) handler() func(http.ResponseWriter, *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		self.ServeHTTP(w, r)
	}
}

func (*Server) respond(w http.ResponseWriter) {

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	// Todo: remove debugging code.
	if err := json.NewEncoder(w).Encode(
		map[string]interface{}{
			"created": true, "total": processing.Total, "signups": processing.Signups, "follows": processing.Follows, "creations": processing.Creations, "impressions": processing.Impressions,
		}); err != nil {

		fmt.Println(err.Error())
	}
}
