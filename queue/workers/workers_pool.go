package workers

import (
	"splash/queue/jobs"
)

type Pool struct {
	// A pool of job requests.
	JobsRequestsPool chan chan jobs.Job
	// The output of the mappers, after it is being processed.
	output chan interface{}
}

func NewPool(jobsRequestsPool chan chan jobs.Job, output chan interface{}) *Pool {
	return &Pool{
		JobsRequestsPool: jobsRequestsPool,
		output: output,
	}
}

func (self *Pool)GetWorkersPoolChannel() (chan chan jobs.Job){
	return self.JobsRequestsPool
}


func (self *Pool) GetOutputChannel() chan interface{} {
	return self.output
}