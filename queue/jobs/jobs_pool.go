package jobs

type Pool struct {
	// A pool of job requests.
	JobsRequestsPool chan chan Job
	JobRequest       chan Job
}

func NewPool(jobsRequestsPool chan chan Job, jobRequest chan Job) *Pool {
	return &Pool{
		JobsRequestsPool: jobsRequestsPool,
		JobRequest: jobRequest,
	}
}

func (self *Pool)GetWorkerChannel() (chan chan Job){

	return self.JobsRequestsPool
}
