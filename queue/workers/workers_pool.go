package workers

import (
	"splash/queue/jobs"
)

type Pool struct {
	// A pool of job requests.
	jobsRequestsPool chan chan jobs.Job
	// The output of the mappers, after it is being processed.
	output           chan interface{}
}

func NewPool(jobsRequestsPool chan chan jobs.Job, output chan interface{}) *Pool {
	return &Pool{
		jobsRequestsPool: jobsRequestsPool,
		output:           output,
	}
}

func (self *Pool) AssignIdleWorkerChannel(jobRequest chan jobs.Job) {
	self.jobsRequestsPool <- jobRequest
}

func (self *Pool) GetWorkersPoolChannel() chan chan jobs.Job {
	return self.jobsRequestsPool
}

func (self *Pool) GetOutputChannel() chan interface{} {
	return self.output
}
