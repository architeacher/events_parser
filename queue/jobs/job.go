package jobs

import (
	"time"
)

// A buffered channel that we can send work requests on.
var JobsQueue chan interface{}

type Job struct {
	id      int
	delay   time.Duration
	created time.Time
	finished time.Time

	// Making the payload it self separated from the job.
	payload interface{}
}

func NewJob(id int, delay time.Duration, created time.Time, payload interface{}) *Job {
	return &Job{
		id:id,
		delay: delay,
		created: created,
		payload: payload,
	}
}

func (self *Job) Id() int {
	return self.id
}

func (self *Job) GetDelay() time.Duration {
	return self.delay
}

func (self *Job) GetCreated() time.Time {
	return self.created
}

func (self *Job) SetFinished(finished time.Time) *Job {
	self.finished = finished
	return self
}

func (self *Job) GetFinished() time.Time {
	return self.finished
}

func (self *Job) GetPayload() interface{} {
	return self.payload
}

func PushToChanel (jobCollection *Collection, JobsQueue chan interface{}) chan interface{} {

	go func() {
		for _, work := range jobCollection.jobs {

			// Push the work to the queue.
			JobsQueue <- work
		}
	}()

	return JobsQueue
}