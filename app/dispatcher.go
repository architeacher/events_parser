package app

import (
	"strconv"
	"splash/processing"
	"splash/processing/aggregation"
	"splash/queue/jobs"
	workersLib "splash/queue/workers"
	"splash/services"
	"splash/logger"
)

const (
	BASE_SERVER      = "base"
	ANALYTICS_SERVER = "analytics"
	MAX_WORKERS      = "maxWorkers"
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
		config:         config,
		maxWorkers:     maxWorkers,
		maxQueuedItems: maxQueuedItems,
	}
}

func (self *Dispatcher) Run() {

	aggregation.AggregationQueue = make(aggregation.AggregationTypeChannel)

	serviceLocator := services.NewLocator()
	logger := serviceLocator.Logger()

	go self.dispatch(logger)

	// Launching the server
	server := NewServer(make(map[string]*RequestsHandler), make(chan *RequestsHandler), self.config[BASE_SERVER], processing.NewOperator(), logger)
	analyticsServer := NewAnalyticsServer(self.config[ANALYTICS_SERVER], aggregation.NewAggregator(), logger)

	go analyticsServer.Start()
	server.Start()
}

func (self *Dispatcher) dispatch(logger *logger.Logger) {

	workers := make([]*workersLib.Worker, self.maxWorkers)
	workersPool := workersLib.NewPool(make(chan chan jobs.Job, self.maxWorkers), make(chan interface{}))

	for index := range workers {
		workers[index] = workersLib.New(strconv.Itoa(index), make(chan int), workersPool, make(chan jobs.Job), logger)
	}

	workersLib.WorkersCollectionHandler = workersLib.NewCollection(workers, workersPool, false)
	workersLib.WorkersCollectionHandler.Start()
}

