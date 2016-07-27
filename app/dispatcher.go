package app

import (
	workersLib "splash/queue/workers"
	"splash/queue/jobs"
	"splash/services"
)

type Dispatcher struct {
	config                     map[string]map[string]string
	maxWorkers, maxQueuedItems int
	workersCollection          *workersLib.Collection
	jobsPool		   *jobs.Pool
}

func NewDispatcher(maxWorkers, maxQueuedItems int, config map[string]map[string]string) *Dispatcher {

	workers := make([]*workersLib.Worker, maxWorkers)
	jobsRequestsPool := make(chan chan jobs.Job, maxWorkers)
	jobRequest := make(chan jobs.Job)

	jobsPool := jobs.NewPool(jobsRequestsPool, jobRequest)

	for index := range workers {
		workers[index] = workersLib.New(index, jobsPool)
	}

	workersCollection := workersLib.NewCollection(workers)

	return &Dispatcher{
		config : config,
		maxWorkers: maxWorkers,
		maxQueuedItems: maxQueuedItems,
		workersCollection: workersCollection,
		jobsPool: jobsPool,
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
	server := NewServer(self.config["base"])
	analyticsServer := NewAnalyticsServer(self.config["analytics"])

	go analyticsServer.Start()
	// Should be last line as it is a blocking.
	server.Start()
}

func (self *Dispatcher) dispatch() {

	// A new job is received.
	for job := range jobs.JobsQueue {

		go func() {

			// Blocking till an idle worker is available, then trying to obtain this available worker's job channel.
			workerChannel := <-self.jobsPool.GetWorkerChannel()

			// Dispatch the job to the worker job channel.
			workerChannel <- job
		}()

	}

	defer close(jobs.JobsQueue)
}