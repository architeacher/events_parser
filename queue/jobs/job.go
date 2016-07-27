package jobs

import (
	"time"
)

// A buffered channel that we can send work requests on.
var JobsQueue chan Job

type Job struct {
	id      int
	delay   time.Duration
	created time.Time
	// Making the payload it self separated from the job.
	Payload *Payload
}

func NewJob(id int, delay time.Duration, created time.Time, payload *Payload) *Job {
	return &Job{
		id:id,
		delay: delay,
		created: created,
		Payload: payload,
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

func PushToChanel (jobCollection *Collection) {

	for _, work := range jobCollection.jobs {

		// Push the work to the queue.
		JobsQueue <- work
	}
}