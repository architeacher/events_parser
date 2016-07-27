package jobs

import (
	"splash/communication/protocols/protobuf"
)

type Payload struct {
	time, actorId, subjectId   	int64
	eventType protobuf.Event_EventType
}

func NewPayload(time int64, eventType protobuf.Event_EventType, actorId, subjectId int64) *Payload {
	return &Payload{
		time: time,
		eventType: eventType,
		actorId: actorId,
		subjectId: subjectId,
	}
}

func NewPayloadFromEventPayload(eventPayload *protobuf.Event_Payload) *Payload {
	return NewPayload(
		eventPayload.GetTime(),
		eventPayload.GetEventType(),
		eventPayload.GetActorId(),
		eventPayload.GetSubjectId(),
	)
}

func (self * Payload) GetTime() int64 {
	return self.time
}

func (self * Payload) GetType() protobuf.Event_EventType {
	return self.eventType
}

func (self * Payload) GetActorId() int64 {
	return self.actorId
}

func (self * Payload) GetSubjectId() int64 {
	return self.subjectId
}
