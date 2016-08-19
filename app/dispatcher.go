package app

import (
	"splash/processing"
	"splash/processing/aggregation"
	"splash/queue/jobs"
	"strconv"
	"splash/processing/map_reduce"
	workersLib "splash/queue/workers"
	"splash/processing/map_reduce/reducers"
	"splash/services"
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

	// Initializing the queue
	jobs.JobsQueue = make(chan interface{}, self.maxQueuedItems)

	go self.dispatch()

	serviceLocator := services.NewLocator()
	logger := serviceLocator.Logger()

	// Launching the server
	server := NewServer(self.config[BASE_SERVER], processing.NewOperator(), logger)
	analyticsServer := NewAnalyticsServer(self.config[ANALYTICS_SERVER], aggregation.NewAggregator())

	go analyticsServer.Start()
	server.Start()
}

func (self *Dispatcher) dispatch() interface{} {
	return <-self.mapReduce(reducers.Reducer, jobs.JobsQueue)
}


func (self *Dispatcher) mapReduce(reducer map_reduce.ReducerFunc, input chan interface{}) chan interface{} {

	workers := make([]*workersLib.Worker, self.maxWorkers)
	workersPool := workersLib.NewPool(make(chan chan jobs.Job, self.maxWorkers), make(chan interface{}))

	for index := range workers {
		workers[index] = workersLib.New(strconv.Itoa(index), make(chan jobs.Job), workersPool, make(chan int))
	}

	workersCollection := workersLib.NewCollection(workers, workersPool, false)

	reducerInput := make(chan interface{})
	reducerOutput := make(chan interface{})
	mapperCollector := make(map_reduce.MapperCollector, self.maxWorkers)

	go reducer(reducerInput, reducerOutput)
	go map_reduce.ReducerDispatcher(mapperCollector, reducerInput)

	// Starting workers mapper
	go workersCollection.DispatchMappers(input, mapperCollector)

	return reducerOutput
}