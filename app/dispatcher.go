package app

import (
	workersLib "splash/queue/workers"
	"splash/processing"
	"splash/queue/jobs"
	"strconv"
	"splash/processing/map_reduce"
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
	workersCollection          *workersLib.Collection
}

func NewDispatcher(config map[string]map[string]string) *Dispatcher {

	maxWorkers, _ := strconv.Atoi(config[BASE_SERVER][MAX_WORKERS])
	maxQueuedItems, _ := strconv.Atoi(config[BASE_SERVER][MAX_QUEUED_ITEMS])

	workers := make([]*workersLib.Worker, maxWorkers)
	jobsRequestsPool := make(chan chan jobs.Job, maxWorkers)

	workersPool := workersLib.NewPool(jobsRequestsPool, make(chan interface{}))

	for index := range workers {
		workers[index] = workersLib.New(index, make(chan jobs.Job), workersPool, make(chan int))
	}

	workersCollection := workersLib.NewCollection(workers, workersPool, map_reduce.Mapper, false)

	return &Dispatcher{
		config : config,
		maxWorkers: maxWorkers,
		maxQueuedItems: maxQueuedItems,
		workersCollection: workersCollection,
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

func (self *Dispatcher) dispatch() {

	// Starting workers
	self.workersCollection.Dispatch(jobs.JobsQueue, make(map_reduce.MapperCollector, self.maxWorkers))

	//reducerInput := make(chan interface{})
	//reducerOutput := make(chan interface{})
	//
	//go map_reduce.Reducer(reducerInput, reducerOutput)
}