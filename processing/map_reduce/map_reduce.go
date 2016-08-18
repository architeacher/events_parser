package map_reduce

import (
	"splash/queue/jobs"
	"splash/communication/protocols/protobuf"
	"fmt"
	"splash/processing"
	"splash/services"
)

// MapperCollector is a channel that collects the output from mapper tasks
type MapperCollector chan chan interface{}

// MapperFunc is a function that performs the mapping part of the MapReduce job
type MapperFunc func(interface{}, chan interface{}) error

// ReducerFunc is a function that performs the reduce part of the MapReduce job
type ReducerFunc func(chan interface{}, chan interface{})

func Mapper(input interface{}, output chan interface{}) error {

	serviceLocator := services.NewLocator()

	results := map[processing.GroupingType]interface{}{}

	job := input.(*jobs.Job)

	payload := job.GetPayload().(*jobs.Payload)
	eventType := payload.GetType()

	time := serviceLocator.GetAsTimestamp(payload.GetTime())

	year, week := time.ISOWeek()

	results[processing.TYPE_GROUPING_BY_TIME] = map[string]interface{}{
		processing.KEY_TIME : time,
		processing.KEY_ACTOR: payload.GetActorId(),
	}

	switch eventType {
	case protobuf.Event_IMPRESSION:
		results[processing.TYPE_GROUPING_BY_IMPRESSION] = map[string]interface{}{
			processing.KEY_WEEK: week,
			processing.KEY_YEAR: year,
			processing.KEY_SUBJECT: payload.GetSubjectId(),
		}
		break
	}

	output <- results

	return nil
}

func Reducer(input chan interface{}, output chan interface{}) {
	results := map[protobuf.Event_EventType]int{}
	for matches := range input {
		for key, value := range matches.(map[protobuf.Event_EventType][]interface{}) {
		//	_, exists := results[key]
		//	if !exists {
		//		results[key] = value
		//	} else {
		//		results[key] = results[key] + value
		//	}
			for index, event := range value {

				fmt.Println(key, index, event)
			}
		}
	}
	output <- results
}

func ReducerDispatcher(collector MapperCollector, reducerInput chan interface{}) {
	for output := range collector {
		reducerInput <- <-output
	}
	close(reducerInput)
}
