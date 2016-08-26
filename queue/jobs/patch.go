package jobs

import (
	"sync/atomic"
	"runtime"
)

type Patch struct {
	id                         string
	jobs                       *[]*Job
	length, finishedJobsLength uint64
}

func NewPatch(id string, jobs *[]*Job) *Patch {

	patch := &Patch{
		id:         id,
		jobs:       jobs,
		length:     uint64(len(*jobs)),
	}

	for _, job := range *jobs {
		job.setPatch(patch)
	}

	return patch
}

func (self *Patch) GetId() string {
	return self.id
}

func (self *Patch) GetLength() uint64 {
	return self.length
}

func (self *Patch) SetFinished(job *Job) *Patch {

	if job == nil {
		return nil
	}

	atomic.AddUint64(&self.finishedJobsLength, 1)
	runtime.Gosched()

	return self
}

func (self *Patch) GetFinishedLength() uint64 {
	return self.finishedJobsLength
}

func (self *Patch) IsFinished() bool {
	return self.GetLength() == self.finishedJobsLength
}
