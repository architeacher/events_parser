// Code generated by protoc-gen-go.
// source: events.proto
// DO NOT EDIT!

/*
Package protobuf is a generated protocol buffer package.

It is generated from these files:
	events.proto

It has these top-level messages:
	Event
*/
package protobuf

import proto "github.com/golang/protobuf/proto"
import fmt "fmt"
import math "math"

// Reference imports to suppress errors if they are not otherwise used.
var _ = proto.Marshal
var _ = fmt.Errorf
var _ = math.Inf

// This is a compile-time assertion to ensure that this generated file
// is compatible with the proto package it is being compiled against.
// A compilation error at this line likely means your copy of the
// proto package needs to be updated.
const _ = proto.ProtoPackageIsVersion2 // please upgrade the proto package

type Event_EventType int32

const (
	Event_SIGNUP          Event_EventType = 0
	Event_FOLLOW          Event_EventType = 1
	Event_SPLASH_CREATION Event_EventType = 2
	Event_IMPRESSION      Event_EventType = 3
)

var Event_EventType_name = map[int32]string{
	0: "SIGNUP",
	1: "FOLLOW",
	2: "SPLASH_CREATION",
	3: "IMPRESSION",
}
var Event_EventType_value = map[string]int32{
	"SIGNUP":          0,
	"FOLLOW":          1,
	"SPLASH_CREATION": 2,
	"IMPRESSION":      3,
}

func (x Event_EventType) Enum() *Event_EventType {
	p := new(Event_EventType)
	*p = x
	return p
}
func (x Event_EventType) String() string {
	return proto.EnumName(Event_EventType_name, int32(x))
}
func (x *Event_EventType) UnmarshalJSON(data []byte) error {
	value, err := proto.UnmarshalJSONEnum(Event_EventType_value, data, "Event_EventType")
	if err != nil {
		return err
	}
	*x = Event_EventType(value)
	return nil
}
func (Event_EventType) EnumDescriptor() ([]byte, []int) { return fileDescriptor0, []int{0, 0} }

type Event struct {
	PayloadCollection []*Event_Payload `protobuf:"bytes,4,rep,name=payloadCollection" json:"payloadCollection,omitempty"`
	XXX_unrecognized  []byte           `json:"-"`
}

func (m *Event) Reset()                    { *m = Event{} }
func (m *Event) String() string            { return proto.CompactTextString(m) }
func (*Event) ProtoMessage()               {}
func (*Event) Descriptor() ([]byte, []int) { return fileDescriptor0, []int{0} }

func (m *Event) GetPayloadCollection() []*Event_Payload {
	if m != nil {
		return m.PayloadCollection
	}
	return nil
}

type Event_Payload struct {
	Time      *int64           `protobuf:"varint,1,req,name=time" json:"time,omitempty"`
	EventType *Event_EventType `protobuf:"varint,2,req,name=eventType,enum=protobuf.Event_EventType" json:"eventType,omitempty"`
	// Id of the user who took the action.
	ActorId *int64 `protobuf:"varint,3,opt,name=actorId" json:"actorId,omitempty"`
	// Id of the user who is affected by the action...
	SubjectId        *int64 `protobuf:"varint,4,opt,name=subjectId" json:"subjectId,omitempty"`
	XXX_unrecognized []byte `json:"-"`
}

func (m *Event_Payload) Reset()                    { *m = Event_Payload{} }
func (m *Event_Payload) String() string            { return proto.CompactTextString(m) }
func (*Event_Payload) ProtoMessage()               {}
func (*Event_Payload) Descriptor() ([]byte, []int) { return fileDescriptor0, []int{0, 0} }

func (m *Event_Payload) GetTime() int64 {
	if m != nil && m.Time != nil {
		return *m.Time
	}
	return 0
}

func (m *Event_Payload) GetEventType() Event_EventType {
	if m != nil && m.EventType != nil {
		return *m.EventType
	}
	return Event_SIGNUP
}

func (m *Event_Payload) GetActorId() int64 {
	if m != nil && m.ActorId != nil {
		return *m.ActorId
	}
	return 0
}

func (m *Event_Payload) GetSubjectId() int64 {
	if m != nil && m.SubjectId != nil {
		return *m.SubjectId
	}
	return 0
}

func init() {
	proto.RegisterType((*Event)(nil), "protobuf.Event")
	proto.RegisterType((*Event_Payload)(nil), "protobuf.Event.Payload")
	proto.RegisterEnum("protobuf.Event_EventType", Event_EventType_name, Event_EventType_value)
}

func init() { proto.RegisterFile("events.proto", fileDescriptor0) }

var fileDescriptor0 = []byte{
	// 226 bytes of a gzipped FileDescriptorProto
	0x1f, 0x8b, 0x08, 0x00, 0x00, 0x09, 0x6e, 0x88, 0x02, 0xff, 0xe2, 0xe2, 0x49, 0x2d, 0x4b, 0xcd,
	0x2b, 0x29, 0xd6, 0x2b, 0x28, 0xca, 0x2f, 0xc9, 0x17, 0xe2, 0x00, 0x53, 0x49, 0xa5, 0x69, 0x4a,
	0x1f, 0x18, 0xb9, 0x58, 0x5d, 0x41, 0x52, 0x42, 0x46, 0x5c, 0x82, 0x05, 0x89, 0x95, 0x39, 0xf9,
	0x89, 0x29, 0xce, 0xf9, 0x39, 0x39, 0xa9, 0xc9, 0x25, 0x99, 0xf9, 0x79, 0x12, 0x2c, 0x0a, 0xcc,
	0x1a, 0xdc, 0x46, 0xe2, 0x7a, 0x30, 0xf5, 0x7a, 0x60, 0xb5, 0x7a, 0x01, 0x10, 0x85, 0x52, 0x99,
	0x5c, 0xec, 0x50, 0xa6, 0x10, 0x0f, 0x17, 0x4b, 0x49, 0x66, 0x6e, 0xaa, 0x04, 0xa3, 0x02, 0x93,
	0x06, 0xb3, 0x90, 0x0e, 0x17, 0x27, 0xd8, 0xc2, 0x90, 0xca, 0x82, 0x54, 0x09, 0x26, 0xa0, 0x10,
	0x9f, 0x91, 0x24, 0xba, 0x21, 0xae, 0x30, 0x05, 0x42, 0xfc, 0x5c, 0xec, 0x89, 0xc9, 0x25, 0xf9,
	0x45, 0x9e, 0x29, 0x12, 0xcc, 0x0a, 0x8c, 0x40, 0xed, 0x82, 0x5c, 0x9c, 0xc5, 0xa5, 0x49, 0x59,
	0x40, 0x47, 0x00, 0x85, 0x58, 0x40, 0x42, 0x4a, 0x1e, 0x5c, 0x9c, 0x08, 0x0d, 0x5c, 0x5c, 0x6c,
	0xc1, 0x9e, 0xee, 0x7e, 0xa1, 0x01, 0x02, 0x0c, 0x20, 0xb6, 0x9b, 0xbf, 0x8f, 0x8f, 0x7f, 0xb8,
	0x00, 0xa3, 0x90, 0x30, 0x17, 0x7f, 0x70, 0x80, 0x8f, 0x63, 0xb0, 0x47, 0xbc, 0x73, 0x90, 0xab,
	0x63, 0x88, 0xa7, 0xbf, 0x9f, 0x00, 0x93, 0x10, 0x1f, 0x17, 0x97, 0xa7, 0x6f, 0x40, 0x90, 0x6b,
	0x70, 0x30, 0x88, 0xcf, 0x0c, 0x08, 0x00, 0x00, 0xff, 0xff, 0x2f, 0xad, 0x71, 0x34, 0x0b, 0x01,
	0x00, 0x00,
}