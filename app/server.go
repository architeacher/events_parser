package app

import (
	//"splash/processing"
	jobsLib "splash/queue/jobs"
	"splash/services"
	"log"
	"net/http"
	"io/ioutil"
	"splash/communication/protocols/protobuf"
	"github.com/golang/protobuf/proto"
	"encoding/json"
	"time"
	"fmt"
	"splash/processing"
)

type Server struct {
	config map[string]string
	aggregator *processing.Aggregator
	operator *processing.Operator
}

func NewServer(config map[string]string, aggregator *processing.Aggregator, operator *processing.Operator) *Server {
	return &Server{
		config: config,
		aggregator: aggregator,
		operator: operator,
	}
}
// Todo: Remove this tracking
var signups, follows, creations, impressions, total int

func (self *Server) Start() {

	s := &http.Server{
		Addr:           self.config["port"],
		Handler:        self,
	}

	go self.aggregator.MonitorNewData()
	go self.aggregator.Aggregate(self.operator)

	// Todo: Use the logger.
	log.Fatal(s.ListenAndServe())
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

	jobs, err := self.getJobsFromBodyData(bodyData)

	if err != nil {
		logger.Error(err.Error())
		return
	}

	jobsLib.PushToChanel(jobsLib.NewCollection(jobs))

	self.respond(w)

	return
}

func (self *Server) handler() func(http.ResponseWriter, *http.Request) {

	return func (w http.ResponseWriter, r *http.Request) {

		self.ServeHTTP(w, r)
	}
}

func (*Server) buildJob(id int, eventPayload *protobuf.Event_Payload) *jobsLib.Job {

	payload := jobsLib.NewPayloadFromEventPayload(eventPayload)
	// Todo: duration should be configurable.
	return jobsLib.NewJob(id, 1000, time.Now(), payload)
}

func (self *Server) getJobsFromBodyData(bodyData []byte) ([]jobsLib.Job, error) {

	protoData := new(protobuf.Event)
	err := proto.Unmarshal(bodyData, protoData)

	if err != nil {
		return nil, err
	}

	jobs := []jobsLib.Job{}

	for index, item := range protoData.GetPayloadCollection() {

		switch item.GetEventType() {
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

		job := self.buildJob(index, item)
		jobs = append(jobs, *job)
	}

	total += len(jobs)

	return jobs, nil
}

func (*Server) respond(w http.ResponseWriter) {

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	if err := json.NewEncoder(w).Encode(map[string]interface{}{"created": true, "total": total, "signups": signups, "follows": follows, "creations": creations, "impressions": impressions}); err != nil {
		fmt.Println(err.Error())
	}
}
