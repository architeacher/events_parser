package jobs

import (
	"splash/communication/protocols/protobuf"
)

type Payload struct {
	time               int64
	actorId, subjectId string
	eventType          protobuf.Event_EventType
}

func NewPayload(time int64, eventType protobuf.Event_EventType, actorId, subjectId string) *Payload {
	return &Payload{
		time:      time,
		eventType: eventType,
		actorId:   actorId,
		subjectId: subjectId,
	}
}

func NewPayloadFromEventPayload(eventPayload *protobuf.Event_Payload) *Payload {
	return NewPayload(
		eventPayload.Time,
		eventPayload.EventType,
		eventPayload.ActorId,
		eventPayload.SubjectId,
	)
}

func (self *Payload) GetTime() int64 {
	return self.time
}

func (self *Payload) GetType() protobuf.Event_EventType {
	return self.eventType
}

func (self *Payload) GetActorId() string {
	return self.actorId
}

func (self *Payload) GetSubjectId() string {
	return self.subjectId
}
