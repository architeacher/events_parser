package workers

import (
	"splash/queue/jobs"
)

type Pool struct {
	// A pool of job requests.
	JobsRequestsPool chan chan jobs.Job
}

func NewPool(jobsRequestsPool chan chan jobs.Job) *Pool {
	return &Pool{
		JobsRequestsPool: jobsRequestsPool,
	}
}

func (self *Pool)GetWorkersPoolChannel() (chan chan jobs.Job){

	return self.JobsRequestsPool
}
