package processing

import (
	"github.com/golang/protobuf/proto"
	"splash/communication/protocols/protobuf"
	"splash/processing/map_reduce"
	"splash/processing/map_reduce/mappers"
	jobsLib "splash/queue/jobs"
	"splash/services"
	"strconv"
)

// Todo: Remove this tracking
var Signups, Follows, Creations, Impressions, Total uint64

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

func (self *Operator) EnumeratePatch(collection *jobsLib.PatchCollection, input chan interface{}) (*jobsLib.Patch, error) {

	jobs := []*jobsLib.Job{}

	serviceLocator := services.NewLocator()

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

			job := jobsLib.NewJobFromEventPayload(serviceLocator.RandString("job-"+strconv.Itoa(index)+"-", 55), event, []map_reduce.MapperFunc{mappers.Mapper})
			jobs = append(jobs, job)
		}
	}

	patch := jobsLib.NewPatch(serviceLocator.RandString("patch-", 55), collection, &jobs)

	collection.Append(patch)

	return patch, nil
}
