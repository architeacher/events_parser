package app

import (
	"strconv"
	"splash/processing"
	"splash/processing/aggregation"
	"splash/processing/map_reduce"
	"splash/processing/map_reduce/reducers"
	"splash/queue/jobs"
	workersLib "splash/queue/workers"
	"splash/services"
	"time"
)

const (
	BASE_SERVER      = "base"
	ANALYTICS_SERVER = "analytics"
	MAX_WORKERS      = "maxWorkers"
	MAX_QUEUED_ITEMS = "maxQueuedItems"
)

var (
	// Todo: Handle the part when to stop assigning jobs to the same request, the client needs to send the length of the jobs.
	PatchesCollection map[string]*jobs.PatchCollection
)

type Dispatcher struct {
	config                     map[string]map[string]string
	maxWorkers, maxQueuedItems int
}

func NewDispatcher(config map[string]map[string]string) *Dispatcher {
	PatchesCollection = make(map[string]*jobs.PatchCollection)

	maxWorkers, _ := strconv.Atoi(config[BASE_SERVER][MAX_WORKERS])
	maxQueuedItems, _ := strconv.Atoi(config[BASE_SERVER][MAX_QUEUED_ITEMS])

	return &Dispatcher{
		config:         config,
		maxWorkers:     maxWorkers,
		maxQueuedItems: maxQueuedItems,
	}
}

func (self *Dispatcher) Run() {

	// Initializing the queue
	jobs.JobsQueue = make(chan interface{}, self.maxQueuedItems)
	map_reduce.IsPatchFinished = make(chan interface{})

	go self.dispatch()

	serviceLocator := services.NewLocator()
	logger := serviceLocator.Logger()

	// Launching the server
	server := NewServer(self.config[BASE_SERVER], processing.NewOperator(), logger)
	analyticsServer := NewAnalyticsServer(self.config[ANALYTICS_SERVER], aggregation.NewAggregator(), logger)

	go analyticsServer.Start()
	server.Start()
}

func (self *Dispatcher) dispatch() chan interface{} {
	return self.mapReduce(reducers.Reducer, jobs.JobsQueue)
}

func (self *Dispatcher) mapReduce(reducer map_reduce.ReducerFunc, input chan interface{}) chan interface{} {

	workers := make([]*workersLib.Worker, self.maxWorkers)
	workersPool := workersLib.NewPool(make(chan chan jobs.Job, self.maxWorkers), make(chan interface{}))

	serviceLocator := services.NewLocator()
	logger := serviceLocator.Logger()

	for index := range workers {
		workers[index] = workersLib.New(strconv.Itoa(index), make(chan int), workersPool, make(chan jobs.Job), logger)
	}

	workersCollection := workersLib.NewCollection(workers, workersPool, false)

	reducerInput := make(chan interface{})
	reducerOutput := make(chan interface{})
	// Buffered channel size should be relative to the current workers number.
	mapperCollector := make(map_reduce.MapperCollector, 2 * self.maxWorkers)

	go reducer(reducerInput, reducerOutput)
	go map_reduce.ReducerDispatcher(mapperCollector, reducerInput)

	// Starting workers mapper
	go workersCollection.DispatchMappers(input, mapperCollector)

	self.monitorPatches()

	return reducerOutput
}

func (self *Dispatcher) monitorPatches() {
	go func() {
		for {
			time.Sleep(time.Duration(1000))
			//for _, patchCollection := range PatchesCollection {
			//
			//	//if patchCollection.IsFinished() {
			//	//	map_reduce.IsPatchFinished <- true
			//	//}
			//
			//	fmt.Println("Length", patchCollection.GetId(), patchCollection.GetLength(), patchCollection.GetFinishedLength())
			//}

		}
	}()
}
