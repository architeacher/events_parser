package jobs

import (
	"splash/communication/protocols/protobuf"
	"splash/processing/map_reduce"
	"time"
)

// A buffered channel that we can send work channel requests on.
type JobsQueue chan interface{}

type Job struct {
	id    string
	patch *Patch

	// Delay the current job for a certain duration.
	delay time.Duration
	// Tracking time
	created  time.Time
	finished time.Time

	isFinished bool

	// Making the payload it self separated from the job.
	payload interface{}
	// List of mappers that should be applied to this job.
	mappers []map_reduce.MapperFunc
}

func NewJob(id string, delay time.Duration, created time.Time, payload interface{}, mappers []map_reduce.MapperFunc) *Job {
	return &Job{
		id:      id,
		delay:   delay,
		created: created,
		payload: payload,
		mappers: mappers,
	}
}

func NewJobFromEventPayload(id string, eventPayload *protobuf.Event_Payload, delay time.Duration, mappers []map_reduce.MapperFunc) *Job {
	payload := NewPayloadFromEventPayload(eventPayload)
	return NewJob(id, delay, time.Now(), payload, mappers)
}

func (self *Job) GetId() string {
	return self.id
}

func (self *Job) setPatch(patch *Patch) *Job {
	self.patch = patch
	return self
}

func (self *Job) GetPatch() *Patch {
	return self.patch
}

func (self *Job) GetDelay() time.Duration {
	return self.delay
}

func (self *Job) GetCreated() time.Time {
	return self.created
}

func (self *Job) SetFinished(isFinished bool, finished time.Time) *Job {
	self.finished = finished
	self.isFinished = isFinished
	self.GetPatch().SetFinished(self)
	return self
}

func (self *Job) GetFinished() time.Time {
	return self.finished
}

func (self *Job) IsFinished() bool {
	return self.isFinished
}

func (self *Job) GetPayload() interface{} {
	return self.payload
}

func (self *Job) AddMapper(mapper map_reduce.MapperFunc) *Job {
	self.mappers = append(self.mappers, mapper)
	return self
}

func (self *Job) GetMappers() []map_reduce.MapperFunc {
	return self.mappers
}

func PopulateChannel(patch *Patch, jobsQueue JobsQueue, shouldClose bool) chan interface{} {

	for _, work := range *patch.jobs {
		// Push the work to the queue.
		jobsQueue <- work
	}

	if shouldClose {
		defer close(jobsQueue)
	}

	return jobsQueue
}