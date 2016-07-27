package workers

import (
	"time"
	"runtime"
	"splash/processing"
	"splash/services"
	"splash/queue/jobs"
	"splash/communication/protocols/protobuf"
)

// Possible worker stats
const (
	PAUSED = 0
	RUNNING = 1
	STOPPED = 2
)

type Worker struct {
	id                int
	state             chan int
	jobsPool 	*jobs.Pool
}

func New(id int, jobsPool *jobs.Pool) *Worker {
	return &Worker{
		id: id,
		state: make(chan int),
		jobsPool: jobsPool,
	}
}

func (self *Worker) Start() {

	serviceLocator := services.NewLocator()
	logger := serviceLocator.Logger()

	go self.setState(RUNNING)

	for {
		// The current worker will register it self the the jobs pool,
		// by passing in a work request chanel to receive work on it later.
		self.jobsPool.JobsRequestsPool <- self.jobsPool.JobRequest

		select {
		// A job is received, on the worker's channel, and picked up by the worker.
		case job := <-self.jobsPool.JobRequest:
			time.Sleep(job.GetDelay())
			err := self.process(&job)

			if err != nil {
				logger.Error(err.Error())
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

func (self *Worker) process(job *jobs.Job) error {

	serviceLocator := services.NewLocator()
	//logger := serviceLocator.Logger()

	//logger.Info("Worker ", self.id, "is processing Job", job.Id(), " - Creeated at:", job.GetCreated())

	eventType := job.Payload.GetType()

	time := serviceLocator.GetAsTimestamp(job.Payload.GetTime())
	day := time.Format("2006-01-02")

	switch eventType {
	case protobuf.Event_SIGNUP:

		// Pushing numbers to merge channel.
		processing.DailyActiveUsers <- day
		break
	}

	return nil
}