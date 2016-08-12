package app

import (
	workersLib "splash/queue/workers"
	"splash/services"
	"splash/processing"
	"splash/queue/jobs"
)

type Dispatcher struct {
	config                     map[string]map[string]string
	maxWorkers, maxQueuedItems int
	workersCollection          *workersLib.Collection
	workersPool                *workersLib.Pool
}

func NewDispatcher(maxWorkers, maxQueuedItems int, config map[string]map[string]string) *Dispatcher {

	workers := make([]*workersLib.Worker, maxWorkers)
	jobsRequestsPool := make(chan chan jobs.Job, maxWorkers)

	workersPool := workersLib.NewPool(jobsRequestsPool)

	for index := range workers {
		workers[index] = workersLib.New(index, make(chan jobs.Job), workersPool, make(chan int))
	}

	workersCollection := workersLib.NewCollection(workers, true)

	return &Dispatcher{
		config : config,
		maxWorkers: maxWorkers,
		maxQueuedItems: maxQueuedItems,
		workersCollection: workersCollection,
		workersPool: workersPool,
	}
}

func (self *Dispatcher) Run() {

	serviceLocator := services.NewLocator()
	go serviceLocator.Stats()

	// Initializing the queue
	jobs.JobsQueue = make(chan jobs.Job, self.maxQueuedItems)

	// Starting workers
	self.workersCollection.Start()

	go self.dispatch()

	// Launching the server
	server := NewServer(self.config["base"], processing.NewAggregator(), processing.NewOperator())
	analyticsServer := NewAnalyticsServer(self.config["analytics"])

	go server.Start()
	go analyticsServer.Start()
}

func (self *Dispatcher) dispatch() {

	// A new job is received.
	for job := range jobs.JobsQueue {

		go func(job jobs.Job) {

			// Blocking till an idle worker is available, then trying to obtain this available worker's job channel,
			// to send a job request on.
			jobRequest := <-self.workersPool.GetWorkersPoolChannel()

			// Dispatch the job to the worker job channel.
			jobRequest <- job
		}(job)

	}

	defer close(jobs.JobsQueue)
}