package app

import (
	"encoding/json"
	"strconv"
	"io/ioutil"
	"net/http"
	"fmt"
	"splash/logger"
	"splash/processing"
	jobsLib "splash/queue/jobs"
	"splash/processing/map_reduce"
	"splash/processing/map_reduce/mappers"
	"splash/processing/map_reduce/runner"
	"splash/queue/workers"
	"splash/processing/map_reduce/reducers"
	"splash/processing/aggregation"
	"io"
)

// Todo: Move this to another file.
type RequestsChannel chan *RequestsHandler

type RequestsHandler struct {
	AuthorizationToken                            *string
	Jobs                                          jobsLib.JobsQueue
	Collector                                     map_reduce.MapperCollector
	TotalJobs, PatchesCount, PatchesReceivedCount int
	IsNew, IsDone	bool
	Result	aggregation.AggregationType
}

type Server struct {
	// Todo: Handle the part when to stop assigning jobs to the same request, the client needs to send the length of the jobs.
	requests        map[string]*RequestsHandler
	requestsChannel RequestsChannel
	config          map[string]string
	operator        *processing.Operator
	logger          *logger.Logger
}

func NewServer(requests map[string]*RequestsHandler, requestsChannel RequestsChannel, config map[string]string, operator *processing.Operator, logger *logger.Logger) *Server {

	return &Server{
		requests: requests,
		requestsChannel: requestsChannel,
		config:   config,
		operator: operator,
		logger:   logger,
	}
}

func (self *Server) Start() {
	listenAddr := self.config["host"] + self.config["port"]

	go self.MonitorNewRequests()

	http.HandleFunc(self.config["path"], self.handler())

	self.logger.Alert("HTTP server:", http.ListenAndServe(listenAddr, nil))
}

func (self *Server) ServeHTTP(w http.ResponseWriter, r *http.Request) {

	header := r.Header

	authorizationToken := header["Authorization"][0]
	patchesCount, err := strconv.Atoi(header["Patches-Count"][0])

	if err != nil {
		self.logger.Alert(err.Error())
	}

	patch, err := self.getPatch(r.Body)

	if err != nil {
		self.logger.Error(err.Error())
	}

	request := self.getRequest(authorizationToken, patchesCount)
	request.PatchesReceivedCount++

	shouldClose := false

	if request.PatchesReceivedCount == request.PatchesCount {
		fmt.Println("Closed")
		shouldClose = true
	}

	if request.IsNew {
		self.requestsChannel <- request
	}

	jobsLib.PopulateChannel(patch, request.Jobs, shouldClose)

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
			"created": true, "total": processing.Total, "patches": processing.Patches, "largest patch": processing.LargestPatch, "signups": processing.Signups, "follows": processing.Follows, "creations": processing.Creations, "impressions": processing.Impressions,
		}); err != nil {

		fmt.Println(err.Error())
	}
}

func (self *Server) getPatch(body io.ReadCloser) (*jobsLib.Patch, error) {
	bodyData, err := ioutil.ReadAll(body)

	if err != nil {
		self.logger.Error(err.Error())
	}

	data, err := self.operator.EnumerateData(bodyData)

	if err != nil {
		self.logger.Error(err.Error())
	}

	patch, err := self.operator.EnumeratePatch(data, []map_reduce.MapperFunc{mappers.Mapper})

	if err != nil {
		return nil, err
	}

	return patch, nil
}

func (self *Server) getRequest(authorizationToken string, patchesCount int) *RequestsHandler {

	request := self.requests[authorizationToken]

	if request == nil {
		request = &RequestsHandler{
			AuthorizationToken: &authorizationToken,
			Jobs: make(jobsLib.JobsQueue, 9 * patchesCount),
			// Todo: Trying to get the number of workers.
			// Buffered channel size should be relative to the current workers number.
			Collector: make(map_reduce.MapperCollector, 3 * patchesCount),
			PatchesCount: patchesCount,
			PatchesReceivedCount: 0,
			IsNew: true,
			IsDone: false,
		}

		self.requests[authorizationToken] = request
	} else {
		request.IsNew = false
	}

	return request
}

func (self *Server) MonitorNewRequests() {
	go func() {
		for {
			select {
			case request := <- self.requestsChannel:
				fmt.Println("Received request")
				result := runner.MapReduce(workers.WorkersCollectionHandler, request.Collector, reducers.Reducer, request.Jobs)

				for item := range result {
					fmt.Println("We are here.")
					output := item.(aggregation.AggregationType)
					request.IsDone = true
					request.Result = output
					// Todo: delete the result from the map
					delete(self.requests, *request.AuthorizationToken)
					aggregation.AggregationQueue <- output
				}
			}
		}
	}()
}
