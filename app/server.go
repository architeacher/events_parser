package app

import (
	"log"
	"net/http"
	"io/ioutil"
	"encoding/json"
	"time"
	"fmt"
	jobsLib "splash/queue/jobs"
	"splash/services"
	"github.com/golang/protobuf/proto"
	"splash/communication/protocols/protobuf"
)

type Server struct {
	config map[string]string
}

func NewServer(config map[string]string) *Server {
	return &Server{
		config: config,
	}
}
// Todo: Remove this tracking
var signups, follows, creations, impressions, total int

func (self *Server) Start() {

	listenAddr := self.config["host"] + self.config["port"]

	http.HandleFunc(self.config["path"], self.handler())

	// Todo: Use the logger.
	log.Fatal("HTTP server:", http.ListenAndServe(listenAddr, nil))
}

func (self *Server) ServeHTTP(w http.ResponseWriter, r *http.Request) {

	// Todo: Logger should be injected.
	serviceLocator := services.NewLocator()
	logger := serviceLocator.Logger()

	bodyData, err := ioutil.ReadAll(r.Body)

	if err != nil {
		logger.Error(err.Error())
		return
	}

	//jobs, err := self.getJobsFromBodyData(bodyData)
	data, err := self.enumerateData(bodyData)

	if err != nil {
		logger.Error(err.Error())
		return
	}

	jobsCollection, err := self.enumerateJobs(data)

	if err != nil {
		logger.Error(err.Error())
		return
	}

	jobsLib.PushToChanel(jobsCollection, jobsLib.JobsQueue)

	self.respond(w)
}

func (self *Server) handler() func(http.ResponseWriter, *http.Request) {

	return func (w http.ResponseWriter, r *http.Request) {

		self.ServeHTTP(w, r)
	}
}

func (self *Server) enumerateData(bodyData []byte) (chan interface{}, error) {

	output := make(chan interface{})

	go func() {

		protoData := new(protobuf.Event)
		proto.Unmarshal(bodyData, protoData)

		output <- protoData.GetPayloadCollection()

		close(output)
	}()

	return output, nil
}

func (self *Server) enumerateJobs(input chan interface{}) (*jobsLib.Collection, error){

	jobs := []jobsLib.Job{}

	for item := range input {
		eventPayload := item.([]*protobuf.Event_Payload)

		for index, event := range eventPayload {

			// Todo: Remove debugging code.
			switch event.GetEventType() {
			case protobuf.Event_SIGNUP:
				signups++
				break
			case protobuf.Event_FOLLOW:
				follows++
				break
			case protobuf.Event_SPLASH_CREATION:
				creations++
				break
			case protobuf.Event_IMPRESSION:
				impressions++
				break
			}

			job := self.buildJob(index, event)
			jobs = append(jobs, *job)
		}
	}

	total += len(jobs)

	return jobsLib.NewCollection(jobs), nil
}

func (*Server) buildJob(id int, eventPayload *protobuf.Event_Payload) *jobsLib.Job {

	payload := jobsLib.NewPayloadFromEventPayload(eventPayload)
	// Todo: duration should be configurable.
	return jobsLib.NewJob(id, 1000, time.Now(), payload)
}

func (*Server) respond(w http.ResponseWriter) {

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)

	if err := json.NewEncoder(w).Encode(
		map[string]interface{}{
			"created": true, "total": total, "signups": signups, "follows": follows, "creations": creations, "impressions": impressions,
		}); err != nil {

		fmt.Println(err.Error())
	}
}
