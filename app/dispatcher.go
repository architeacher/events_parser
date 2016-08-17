package app

import (
	"splash/processing"
	"splash/queue/jobs"
	"strconv"
	"splash/processing/map_reduce"
	workersLib "splash/queue/workers"
)

const (
	BASE_SERVER = "base"
	ANALYTICS_SERVER = "analytics"
	MAX_WORKERS = "maxWorkers"
	MAX_QUEUED_ITEMS = "maxQueuedItems"
)

type Dispatcher struct {
	config                     map[string]map[string]string
	maxWorkers, maxQueuedItems int
}

func NewDispatcher(config map[string]map[string]string) *Dispatcher {

	maxWorkers, _ := strconv.Atoi(config[BASE_SERVER][MAX_WORKERS])
	maxQueuedItems, _ := strconv.Atoi(config[BASE_SERVER][MAX_QUEUED_ITEMS])

	return &Dispatcher{
		config : config,
		maxWorkers: maxWorkers,
		maxQueuedItems: maxQueuedItems,
	}
}

func (self *Dispatcher) Run() {

	//serviceLocator := services.NewLocator()
	//go serviceLocator.Stats()

	// Initializing the queue
	jobs.JobsQueue = make(chan interface{}, self.maxQueuedItems)

	go self.dispatch()

	// Launching the server
	server := NewServer(self.config[BASE_SERVER])
	analyticsServer := NewAnalyticsServer(self.config[ANALYTICS_SERVER], processing.NewAggregator(), processing.NewOperator())

	go analyticsServer.Start()
	server.Start()
}

func (self *Dispatcher) dispatch() interface{} {

	results := self.mapReduce(map_reduce.Mapper, map_reduce.Reducer, jobs.JobsQueue)
	return results
}


func (self *Dispatcher) mapReduce(mapper map_reduce.MapperFunc, reducer map_reduce.ReducerFunc, input chan interface{}) interface{} {

	workers := make([]*workersLib.Worker, self.maxWorkers)
	workersPool := workersLib.NewPool(make(chan chan jobs.Job, self.maxWorkers), make(chan interface{}))

	for index := range workers {
		workers[index] = workersLib.New(index, make(chan jobs.Job), workersPool, make(chan int))
	}

	workersCollection := workersLib.NewCollection(workers, workersPool, false)

	reducerInput := make(chan interface{})
	reducerOutput := make(chan interface{})
	mapperCollector := make(map_reduce.MapperCollector, self.maxWorkers)

	go reducer(reducerInput, reducerOutput)
	go map_reduce.ReducerDispatcher(mapperCollector, reducerInput)

	// Starting workers mapper
	go workersCollection.DispatchMappers(mapper, input, mapperCollector)

	return <-reducerOutput
}