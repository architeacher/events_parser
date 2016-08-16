package map_reduce

import (
	"splash/services"
	"splash/queue/jobs"
	"splash/communication/protocols/protobuf"
)

// MapperCollector is a channel that collects the output from mapper tasks
type MapperCollector chan chan interface{}

// MapperFunc is a function that performs the mapping part of the MapReduce job
type MapperFunc func(interface{}, chan interface{}) error

// ReducerFunc is a function that performs the reduce part of the MapReduce job
type ReducerFunc func(chan interface{}, chan interface{})

type MapReduce struct {
}


func NewMapReduce() *MapReduce {
	return &MapReduce{}
}

func Mapper(input interface{}, output chan interface{}) error {

	results := map[string]int{}

	serviceLocator := services.NewLocator()

	//logger := serviceLocator.Logger()
	//logger.Info("Worker ", self.id, "is processing Job", job.Id(), " - Created at:", job.GetCreated())

	job := input.(*jobs.Job)
	payload := job.GetPayload()
	eventType := payload.GetType()

	time := serviceLocator.GetAsTimestamp(payload.GetTime())
	day := time.Format("2006-01-02")

	switch eventType {
	case protobuf.Event_SIGNUP:
		// Pushing numbers to merge channel.
		results[day]++
		break
	}

	output <- results

	return nil
}

func reducerDispatcher(collector MapperCollector, reducerInput chan interface{}) {
	for output := range collector {
		reducerInput <- <-output
	}
	close(reducerInput)
}

//func Reducer(input chan interface{}, output chan interface{}) {
//	results := map[Telemetry]int{}
//	for matches := range input {
//		for key, value := range matches.(map[Telemetry]int) {
//			_, exists := results[key]
//			if !exists {
//				results[key] = value
//			} else {
//				results[key] = results[key] + value
//			}
//		}
//	}
//	output <- results
//}
//
//func mapReduce(mapper MapperFunc, reducer ReducerFunc, input chan interface{}) interface{} {
//
//	reducerInput := make(chan interface{})
//	reducerOutput := make(chan interface{})
//	mapperCollector := make(MapperCollector, MaxWorkers)
//
//	go reducer(reducerInput, reducerOutput)
//	go reducerDispatcher(mapperCollector, reducerInput)
//	go mapperDispatcher(mapper, input, mapperCollector)
//
//	return <-reducerOutput
//}