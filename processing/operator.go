package processing

import (
	"splash/communication/protocols/protobuf"
	"github.com/golang/protobuf/proto"
	"strconv"
	"splash/processing/map_reduce"
	"splash/processing/map_reduce/mappers"
	jobsLib "splash/queue/jobs"
)
// Todo: Remove this tracking
var Signups, Follows, Creations, Impressions, Total int

type Operator struct {
}

func NewOperator() *Operator {
	return &Operator{}
}

func (self *Operator) EnumerateData(bodyData []byte) (chan interface{}, error) {

	output := make(chan interface{})

	go func() {

		protoData := new(protobuf.Event)
		proto.Unmarshal(bodyData, protoData)

		output <- protoData.GetPayloadCollection()

		close(output)
	}()

	return output, nil
}

func (self *Operator) EnumerateJobs(input chan interface{}) (*jobsLib.Collection, error){

	jobs := []jobsLib.Job{}

	for item := range input {
		eventPayload := item.([]*protobuf.Event_Payload)

		for index, event := range eventPayload {

			// Todo: Remove debugging code.
			switch event.EventType {
			case protobuf.Event_SIGNUP:
				Signups++
				break
			case protobuf.Event_FOLLOW:
				Follows++
				break
			case protobuf.Event_SPLASH_CREATION:
				Creations++
				break
			case protobuf.Event_IMPRESSION:
				Impressions++
				break
			}

			job := jobsLib.NewJobFromEventPayload(strconv.Itoa(index), event,[]map_reduce.MapperFunc{mappers.Mapper})
			jobs = append(jobs, *job)
		}
	}

	Total += len(jobs)

	return jobsLib.NewCollection(jobs), nil
}