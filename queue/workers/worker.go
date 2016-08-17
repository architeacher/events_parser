package workers

import (
	"time"
	"runtime"
	"sync"
	"splash/services"
	"splash/queue/jobs"
	"splash/processing/map_reduce"
)

// Possible worker stats
const (
	PAUSED = 0
	RUNNING = 1
	STOPPED = 2
)

type Worker struct {
	id          int
	state       chan int
	// A chanel that the worker can receive a work on.
	jobRequest  chan jobs.Job
	workersPool *Pool
}

func New(id int, jobRequest chan jobs.Job, workersPool *Pool, state chan int) *Worker {
	return &Worker{
		id: id,
		state: state,
		jobRequest: jobRequest,
		workersPool: workersPool,
	}
}

func (self *Worker) Start(mapper map_reduce.MapperFunc, waitGroup interface{}) {

	serviceLocator := services.NewLocator()
	logger := serviceLocator.Logger()

	go self.setState(RUNNING)

	for {
		// The current worker will register it self as an idle, by adding its own job chanel to the jobs pool,
		// so that it will receive work on this chanel later.
		self.workersPool.JobsRequestsPool <- self.jobRequest

		select {
		// A job is received, on the worker's channel, and picked up by the worker.
		case job := <-self.jobRequest:

			time.Sleep(job.GetDelay())

			err := mapper(&job, self.workersPool.GetOutputChannel())
			if err != nil {
				logger.Error("Error processing job:", job.Id(), err.Error())
			}

			job.SetFinished(time.Now())

			if waitGroup != nil {
				wg := waitGroup.(sync.WaitGroup)
				wg.Done()
			}

		// Workers will stop working after 24 hours, taking a nap :P
		case <-time.After(time.Hour * 24):
			self.Stop()

		case state := <-self.state:

			switch state {
			case PAUSED:
				logger.Info("Worker", self.id, "is paused.")

			case RUNNING:
				logger.Info("Worker", self.id, "is started.")

			case STOPPED:
				logger.Info("Worker", self.id, "is stopped.")
				return

			default:
				// We use runtime.Gosched() to prevent a deadlock in this case.
				// It will not be needed of work is performed here which yields
				// to the scheduler.
				runtime.Gosched()

				if PAUSED == state {
					break
				}
			}
		}
	}

}

func (self *Worker) Pause() {

	go self.setState(PAUSED)
}

func (self *Worker) Stop() {

	go self.setState(STOPPED)
}

func (self *Worker) setState(status int) {

	self.state <- status
}