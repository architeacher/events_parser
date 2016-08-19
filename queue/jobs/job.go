package jobs

import (
	"time"
	"splash/communication/protocols/protobuf"
	"splash/processing/map_reduce"
)

// A buffered channel that we can send work requests on.
var JobsQueue chan interface{}

type Job struct {
	id      string
	// Delay the current job for a certain duration.
	delay   time.Duration
	// Tracking time
	created time.Time
	finished time.Time

	// Making the payload it self separated from the job.
	payload interface{}
	// List of mappers that should be applied to this job.
	mappers []map_reduce.MapperFunc
}

func NewJob(id string, delay time.Duration, created time.Time, payload interface{}, mappers []map_reduce.MapperFunc) *Job {
	return &Job{
		id:id,
		delay: delay,
		created: created,
		payload: payload,
		mappers: mappers,
	}
}

func NewJobFromEventPayload(id string, eventPayload *protobuf.Event_Payload, mappers []map_reduce.MapperFunc) *Job {

	payload := NewPayloadFromEventPayload(eventPayload)
	// Todo: duration should be configurable.
	return NewJob(id, 100, time.Now(), payload, mappers)
}

func (self *Job) Id() string {
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

func (self *Job) AddMapper(mapper map_reduce.MapperFunc) *Job {
	self.mappers = append(self.mappers, mapper)
	return self
}

func (self *Job) GetMappers() []map_reduce.MapperFunc {
	return self.mappers
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