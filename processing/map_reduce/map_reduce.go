package map_reduce

import (
	"splash/queue/jobs"
	"splash/communication/protocols/protobuf"
	"fmt"
	"reflect"
)

// MapperCollector is a channel that collects the output from mapper tasks
type MapperCollector chan chan interface{}

// MapperFunc is a function that performs the mapping part of the MapReduce job
type MapperFunc func(interface{}, chan interface{}) error

// ReducerFunc is a function that performs the reduce part of the MapReduce job
type ReducerFunc func(chan interface{}, chan interface{})

func Mapper(input interface{}, output chan interface{}) error {

	results := map[protobuf.Event_EventType][]interface{}{}

	job := input.(*jobs.Job)

	payload := job.GetPayload().(*jobs.Payload)
	eventType := payload.GetType()

	//time := serviceLocator.GetAsTimestamp(payload.GetTime())
	//day := time.Format("2006-01-02")

	results[eventType] = append(results[eventType], payload)

	switch eventType {
	case protobuf.Event_SIGNUP:
		// Pushing numbers to merge channel.
		break
	}

	output <- results

	return nil
}

func Reducer(input chan interface{}, output chan interface{}) {
	results := map[protobuf.Event_EventType]int{}
	for matches := range input {
		//for key, value := range matches.(map[protobuf.Event_EventType]int) {
		//	_, exists := results[key]
		//	if !exists {
		//		results[key] = value
		//	} else {
		//		results[key] = results[key] + value
		//	}
		//}
		fmt.Println(reflect.TypeOf(matches))
	}
	output <- results
}

func ReducerDispatcher(collector MapperCollector, reducerInput chan interface{}) {
	for output := range collector {
		reducerInput <- <-output
	}
	close(reducerInput)
}
