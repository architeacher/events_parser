package jobs

import (
	"splash/structure"
	"sync/atomic"
	"runtime"
)

type PatchCollection struct {
	collection            *structure.Collection
	finishedPatchesLength uint64
	structure.CollectionInterface
}

func NewPatchCollection(id string) *PatchCollection {
	return &PatchCollection{
		collection: structure.NewCollection(id),
	}
}

func (self *PatchCollection) GetId() string {
	return self.collection.GetId()
}

func (self *PatchCollection) GetLength() uint64 {
	return self.collection.GetLength()
}

func (self *PatchCollection) Append(item interface{}) *PatchCollection {
	self.collection.Append(item)
	return self
}

func (self *PatchCollection) AppendItems(items []interface{}) *PatchCollection {
	self.collection.Append(items)
	return self
}

func (self *PatchCollection) GetFinishedLength() uint64 {
	return self.finishedPatchesLength
}

func (self *PatchCollection) SetFinished(patchId string) *PatchCollection{

	if patchId != "" {
		atomic.AddUint64(&self.finishedPatchesLength, 1)
		runtime.Gosched()
	}
	return self
}

func (self *PatchCollection) IsFinished() bool {
	return self.collection.GetLength() == self.finishedPatchesLength
}

