package jobs

import (
	"sync/atomic"
	"runtime"
)

type Patch struct {
	id                         string
	collection                 *PatchCollection
	jobs                       *[]*Job
	length, finishedJobsLength uint64
}

func NewPatch(id string, collection *PatchCollection, jobs *[]*Job) *Patch {

	patch := &Patch{
		id:         id,
		collection: collection,
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

func (self *Patch) GetCollection() *PatchCollection {
	return self.collection
}

func (self *Patch) SetFinished(jobId string) {

	if jobId != "" {
		atomic.AddUint64(&self.finishedJobsLength, 1)
		runtime.Gosched()
	}

	if self.IsFinished() {
		self.GetCollection().SetFinished(self.GetId())
	}
}

func (self *Patch) GetFinishedLength() uint64 {
	return self.finishedJobsLength
}

func (self *Patch) IsFinished() bool {
	return self.length == self.finishedJobsLength
}
