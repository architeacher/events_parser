package workers

import (
	"sync"
	"splash/processing/map_reduce"
	"splash/queue/jobs"
)

type Collection struct {
	workers      []*Worker
	length       int
	workersPool  *Pool
	mapper	     map_reduce.MapperFunc
	isGrouped    bool
}

func NewCollection(workers []*Worker, workersPool *Pool, isGrouped bool) *Collection {
	return &Collection{
		workers: workers,
		length: len(workers),
		// Should workers work as a group.
		workersPool: workersPool,
		isGrouped: isGrouped,
	}
}

func (self *Collection) DispatchMappers(mapper map_reduce.MapperFunc, tasks chan interface{}, collector map_reduce.MapperCollector) {

	self.SetMapper(mapper)
	self.Start()

	// A new job is received.
	for task := range tasks {

		job := task.(jobs.Job)

		go func(job jobs.Job) {

			// Blocking till an idle worker is available, then trying to obtain this available worker's job channel,
			// to send a job request on.
			jobRequest := <-self.workersPool.GetWorkersPoolChannel()

			// Dispatch the job to the worker job channel.
			jobRequest <- job
		}(job)

		collector <- self.workersPool.GetOutputChannel()
	}

	defer close(collector)
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

		if self.isGrouped{
			go worker.Start(self.mapper, wg)
		} else {
			go worker.Start(self.mapper, nil)
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

func (self *Collection) SetMapper(mapper map_reduce.MapperFunc) *Collection {
	self.mapper = mapper
	return self
}