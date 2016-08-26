package workers

import (
	"splash/processing/map_reduce"
	"splash/queue/jobs"
	"sync"
)

type Collection struct {
	workers     []*Worker
	length      int
	workersPool *Pool
	isGrouped   bool
}

var (
	WorkersCollectionHandler *Collection
)

func NewCollection(workers []*Worker, workersPool *Pool, isGrouped bool) *Collection {
	return &Collection{
		workers: workers,
		length:  len(workers),
		// Should workers work as a group.
		workersPool: workersPool,
		isGrouped:   isGrouped,
	}
}

func (self *Collection) DispatchMappers(tasks jobs.JobsQueue, collector map_reduce.MapperCollector) {

	defer close(collector)

	// A new task is received.
	for task := range tasks {
		job := task.(*jobs.Job)

		go func(job jobs.Job) {
			// Blocking till an idle worker is available, then trying to obtain this available worker's job channel,
			// to send a job request on.
			jobRequest := <-self.workersPool.GetWorkersPoolChannel()
			// Dispatch the job to the worker job channel.
			jobRequest <- job
		}(*job)

		collector <- self.workersPool.GetOutputChannel()
	}
}

func (self *Collection) Start() {

	var wg *sync.WaitGroup

	if self.isGrouped {
		wg = new(sync.WaitGroup)
	}

	for _, worker := range self.workers {

		if self.isGrouped {
			wg.Add(1)
		}

		if self.isGrouped {
			go worker.Start(wg)
		} else {
			go worker.Start(nil)
		}
	}

	if self.isGrouped {
		wg.Wait()
	}
}

func (self *Collection) Pause() {
	for _, worker := range self.workers {
		worker.Pause()
	}
}

func (self *Collection) Stop() {
	for _, worker := range self.workers {
		worker.Stop()
	}
}

func (self *Collection) Restart() {
	self.Stop()
	self.Start()
}

func (self *Collection) GetLength() int {
	return self.length
}
