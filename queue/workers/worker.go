package workers

import (
	"runtime"
	"splash/queue/jobs"
	"sync"
	"time"
	"splash/logger"
)

// Possible worker stats
const (
	PAUSED  = 0
	RUNNING = 1
	STOPPED = 2
)

var (
	workingTimes uint64 = 0
)

type Worker struct {
	id    string
	state chan int
	workersPool *Pool
	// A chanel that the worker can receive a work on.
	jobRequest  chan jobs.Job
	logger *logger.Logger
}

func New(id string, state chan int, workersPool *Pool, jobRequest chan jobs.Job, logger *logger.Logger) *Worker {
	return &Worker{
		id:          id,
		state:       state,
		workersPool: workersPool,
		jobRequest:  jobRequest,
		logger: logger,
	}
}

func (self *Worker) Start(waitGroup interface{}) {
	go self.setState(RUNNING)

	for {
		// The current worker will register it self as an idle, by adding its own job chanel to the jobs pool,
		// so that it will receive work on this chanel later.
		self.workersPool.AssignIdleWorkerChannel(self.jobRequest)

		select {
		// A job is received, on the worker's channel, and picked up by the worker.
		case job := <-self.jobRequest:
			self.process(&job, waitGroup)

		// Workers will stop working after 24 hours, taking a nap :P
		//case <-time.After(time.Hour * 24):
		//	self.Stop()

		case state := <-self.state:
			self.checkState(state)
		}
	}

}

func (self *Worker) Pause() {
	go self.setState(PAUSED)
}

func (self *Worker) Stop() {
	go self.setState(STOPPED)
}

func (self *Worker) Id() string {
	return self.id
}

func (self *Worker) setState(status int) {
	self.state <- status
}

func (self *Worker) process(job *jobs.Job, waitGroup interface{}) {
	time.Sleep(job.GetDelay())

	for _, mapper := range job.GetMappers() {
		err := mapper(job, self.workersPool.GetOutputChannel())
		if err != nil {
			self.logger.Error("Error processing job:", job.GetId(), err.Error())
		}
	}

	job.SetFinished(true, time.Now())

	if waitGroup != nil {
		wg := waitGroup.(sync.WaitGroup)
		wg.Done()
	}
}

func (self *Worker) checkState(state int) {
	switch state {
	case PAUSED:
		self.logger.Info("Worker", self.id, "is paused.")

	case RUNNING:
		self.logger.Info("Worker", self.id, "is started.")

	case STOPPED:
		self.logger.Info("Worker", self.id, "is stopped.")
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