package runner

import (
	"splash/queue/workers"
	"splash/processing/map_reduce"
	"splash/queue/jobs"
)

func MapReduce(workersCollection *workers.Collection, mapperCollector map_reduce.MapperCollector, reducer map_reduce.ReducerFunc, input jobs.JobsQueue) chan interface{} {

	reducerInput := make(chan interface{})
	reducerOutput := make(chan interface{})

	go reducer(reducerInput, reducerOutput)

	// Starting workers mapper
	go workersCollection.DispatchMappers(input, mapperCollector)
	go map_reduce.ReducerDispatcher(mapperCollector, reducerInput)

	return reducerOutput
}