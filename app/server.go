package app

import (
	"splash/processing"
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
)

type Server struct {
	config map[string]string
}

func NewServer(config map[string]string) *Server {
	return &Server{
		config:config,
	}
}

func (self *Server) Start() {

	listenAddr := self.config["host"] + ":" + self.config["port"]

	http.HandleFunc(self.config["path"], handler(self))

	aggregator := processing.NewAggregator()

	go aggregator.MonitorNewData()
	go aggregator.Aggregate(processing.NewOperator())

	// Should be last line as it is a blocking.
	log.Fatal(http.ListenAndServe(listenAddr, nil))
}

func handler(self *Server) func(http.ResponseWriter, *http.Request) {

	return func(w http.ResponseWriter, r *http.Request) {

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
		fmt.Println(len(jobs))
		jobsLib.PushToChanel(jobsLib.NewCollection(jobs))

		self.respond(w)

		return
	}
}

func (*Server) buildJob(id int, eventPayload *protobuf.Event_Payload) *jobsLib.Job {

	payload := jobsLib.NewPayloadFromEventPayload(eventPayload)
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

		job := self.buildJob(index, item)
		jobs = append(jobs, *job)
	}

	return jobs, nil
}

func (*Server) respond(w http.ResponseWriter) {

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(map[string]interface{}{"created": true})
}
