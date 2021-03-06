// Code generated by protoc-gen-go. DO NOT EDIT.
// versions:
// 	protoc-gen-go v1.27.1
// 	protoc        v3.13.0
// source: event.proto

package event

import (
	context "context"
	timestamp "github.com/golang/protobuf/ptypes/timestamp"
	grpc "google.golang.org/grpc"
	codes "google.golang.org/grpc/codes"
	status "google.golang.org/grpc/status"
	protoreflect "google.golang.org/protobuf/reflect/protoreflect"
	protoimpl "google.golang.org/protobuf/runtime/protoimpl"
	reflect "reflect"
	sync "sync"
)

const (
	// Verify that this generated code is sufficiently up-to-date.
	_ = protoimpl.EnforceVersion(20 - protoimpl.MinVersion)
	// Verify that runtime/protoimpl is sufficiently up-to-date.
	_ = protoimpl.EnforceVersion(protoimpl.MaxVersion - 20)
)

type Event struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	EventId         string               `protobuf:"bytes,1,opt,name=eventId,proto3" json:"eventId,omitempty"`
	EventType       string               `protobuf:"bytes,2,opt,name=eventType,proto3" json:"eventType,omitempty"`
	EntityType      string               `protobuf:"bytes,3,opt,name=entityType,proto3" json:"entityType,omitempty"`
	EntityId        string               `protobuf:"bytes,4,opt,name=entityId,proto3" json:"entityId,omitempty"`
	Recorded        *timestamp.Timestamp `protobuf:"bytes,5,opt,name=recorded,proto3" json:"recorded,omitempty"`
	Payload         string               `protobuf:"bytes,6,opt,name=payload,proto3" json:"payload,omitempty"`
	PreviousEventId string               `protobuf:"bytes,7,opt,name=previousEventId,proto3" json:"previousEventId,omitempty"`
}

func (x *Event) Reset() {
	*x = Event{}
	if protoimpl.UnsafeEnabled {
		mi := &file_event_proto_msgTypes[0]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *Event) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*Event) ProtoMessage() {}

func (x *Event) ProtoReflect() protoreflect.Message {
	mi := &file_event_proto_msgTypes[0]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use Event.ProtoReflect.Descriptor instead.
func (*Event) Descriptor() ([]byte, []int) {
	return file_event_proto_rawDescGZIP(), []int{0}
}

func (x *Event) GetEventId() string {
	if x != nil {
		return x.EventId
	}
	return ""
}

func (x *Event) GetEventType() string {
	if x != nil {
		return x.EventType
	}
	return ""
}

func (x *Event) GetEntityType() string {
	if x != nil {
		return x.EntityType
	}
	return ""
}

func (x *Event) GetEntityId() string {
	if x != nil {
		return x.EntityId
	}
	return ""
}

func (x *Event) GetRecorded() *timestamp.Timestamp {
	if x != nil {
		return x.Recorded
	}
	return nil
}

func (x *Event) GetPayload() string {
	if x != nil {
		return x.Payload
	}
	return ""
}

func (x *Event) GetPreviousEventId() string {
	if x != nil {
		return x.PreviousEventId
	}
	return ""
}

type PushRequest struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	IdempotencyKey string `protobuf:"bytes,1,opt,name=idempotencyKey,proto3" json:"idempotencyKey,omitempty"`
	EventType      string `protobuf:"bytes,2,opt,name=eventType,proto3" json:"eventType,omitempty"`
	EntityType     string `protobuf:"bytes,3,opt,name=entityType,proto3" json:"entityType,omitempty"`
	EntityId       string `protobuf:"bytes,4,opt,name=entityId,proto3" json:"entityId,omitempty"`
	Payload        string `protobuf:"bytes,5,opt,name=payload,proto3" json:"payload,omitempty"`
}

func (x *PushRequest) Reset() {
	*x = PushRequest{}
	if protoimpl.UnsafeEnabled {
		mi := &file_event_proto_msgTypes[1]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *PushRequest) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*PushRequest) ProtoMessage() {}

func (x *PushRequest) ProtoReflect() protoreflect.Message {
	mi := &file_event_proto_msgTypes[1]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use PushRequest.ProtoReflect.Descriptor instead.
func (*PushRequest) Descriptor() ([]byte, []int) {
	return file_event_proto_rawDescGZIP(), []int{1}
}

func (x *PushRequest) GetIdempotencyKey() string {
	if x != nil {
		return x.IdempotencyKey
	}
	return ""
}

func (x *PushRequest) GetEventType() string {
	if x != nil {
		return x.EventType
	}
	return ""
}

func (x *PushRequest) GetEntityType() string {
	if x != nil {
		return x.EntityType
	}
	return ""
}

func (x *PushRequest) GetEntityId() string {
	if x != nil {
		return x.EntityId
	}
	return ""
}

func (x *PushRequest) GetPayload() string {
	if x != nil {
		return x.Payload
	}
	return ""
}

type GetRequest struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	EventId string `protobuf:"bytes,1,opt,name=eventId,proto3" json:"eventId,omitempty"`
}

func (x *GetRequest) Reset() {
	*x = GetRequest{}
	if protoimpl.UnsafeEnabled {
		mi := &file_event_proto_msgTypes[2]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *GetRequest) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*GetRequest) ProtoMessage() {}

func (x *GetRequest) ProtoReflect() protoreflect.Message {
	mi := &file_event_proto_msgTypes[2]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use GetRequest.ProtoReflect.Descriptor instead.
func (*GetRequest) Descriptor() ([]byte, []int) {
	return file_event_proto_rawDescGZIP(), []int{2}
}

func (x *GetRequest) GetEventId() string {
	if x != nil {
		return x.EventId
	}
	return ""
}

type EntityStreamRequest struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	EntityType string                      `protobuf:"bytes,1,opt,name=entityType,proto3" json:"entityType,omitempty"`
	EntityId   string                      `protobuf:"bytes,2,opt,name=entityId,proto3" json:"entityId,omitempty"`
	Filter     *EntityStreamRequest_Filter `protobuf:"bytes,3,opt,name=filter,proto3" json:"filter,omitempty"`
}

func (x *EntityStreamRequest) Reset() {
	*x = EntityStreamRequest{}
	if protoimpl.UnsafeEnabled {
		mi := &file_event_proto_msgTypes[3]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *EntityStreamRequest) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*EntityStreamRequest) ProtoMessage() {}

func (x *EntityStreamRequest) ProtoReflect() protoreflect.Message {
	mi := &file_event_proto_msgTypes[3]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use EntityStreamRequest.ProtoReflect.Descriptor instead.
func (*EntityStreamRequest) Descriptor() ([]byte, []int) {
	return file_event_proto_rawDescGZIP(), []int{3}
}

func (x *EntityStreamRequest) GetEntityType() string {
	if x != nil {
		return x.EntityType
	}
	return ""
}

func (x *EntityStreamRequest) GetEntityId() string {
	if x != nil {
		return x.EntityId
	}
	return ""
}

func (x *EntityStreamRequest) GetFilter() *EntityStreamRequest_Filter {
	if x != nil {
		return x.Filter
	}
	return nil
}

type GlobalStreamRequest struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	EventIdFrom string   `protobuf:"bytes,1,opt,name=eventIdFrom,proto3" json:"eventIdFrom,omitempty"`
	EntityType  []string `protobuf:"bytes,2,rep,name=entityType,proto3" json:"entityType,omitempty"`
	EventType   []string `protobuf:"bytes,3,rep,name=eventType,proto3" json:"eventType,omitempty"`
}

func (x *GlobalStreamRequest) Reset() {
	*x = GlobalStreamRequest{}
	if protoimpl.UnsafeEnabled {
		mi := &file_event_proto_msgTypes[4]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *GlobalStreamRequest) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*GlobalStreamRequest) ProtoMessage() {}

func (x *GlobalStreamRequest) ProtoReflect() protoreflect.Message {
	mi := &file_event_proto_msgTypes[4]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use GlobalStreamRequest.ProtoReflect.Descriptor instead.
func (*GlobalStreamRequest) Descriptor() ([]byte, []int) {
	return file_event_proto_rawDescGZIP(), []int{4}
}

func (x *GlobalStreamRequest) GetEventIdFrom() string {
	if x != nil {
		return x.EventIdFrom
	}
	return ""
}

func (x *GlobalStreamRequest) GetEntityType() []string {
	if x != nil {
		return x.EntityType
	}
	return nil
}

func (x *GlobalStreamRequest) GetEventType() []string {
	if x != nil {
		return x.EventType
	}
	return nil
}

type SubscribeGlobalStreamRequest struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	EventIdFrom string   `protobuf:"bytes,1,opt,name=eventIdFrom,proto3" json:"eventIdFrom,omitempty"`
	EntityType  []string `protobuf:"bytes,2,rep,name=entityType,proto3" json:"entityType,omitempty"`
	EventType   []string `protobuf:"bytes,3,rep,name=eventType,proto3" json:"eventType,omitempty"`
}

func (x *SubscribeGlobalStreamRequest) Reset() {
	*x = SubscribeGlobalStreamRequest{}
	if protoimpl.UnsafeEnabled {
		mi := &file_event_proto_msgTypes[5]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *SubscribeGlobalStreamRequest) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*SubscribeGlobalStreamRequest) ProtoMessage() {}

func (x *SubscribeGlobalStreamRequest) ProtoReflect() protoreflect.Message {
	mi := &file_event_proto_msgTypes[5]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use SubscribeGlobalStreamRequest.ProtoReflect.Descriptor instead.
func (*SubscribeGlobalStreamRequest) Descriptor() ([]byte, []int) {
	return file_event_proto_rawDescGZIP(), []int{5}
}

func (x *SubscribeGlobalStreamRequest) GetEventIdFrom() string {
	if x != nil {
		return x.EventIdFrom
	}
	return ""
}

func (x *SubscribeGlobalStreamRequest) GetEntityType() []string {
	if x != nil {
		return x.EntityType
	}
	return nil
}

func (x *SubscribeGlobalStreamRequest) GetEventType() []string {
	if x != nil {
		return x.EventType
	}
	return nil
}

type EventStreamQuantum struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Event *Event            `protobuf:"bytes,1,opt,name=event,proto3" json:"event,omitempty"`
	Meta  map[string]string `protobuf:"bytes,2,rep,name=meta,proto3" json:"meta,omitempty" protobuf_key:"bytes,1,opt,name=key,proto3" protobuf_val:"bytes,2,opt,name=value,proto3"`
}

func (x *EventStreamQuantum) Reset() {
	*x = EventStreamQuantum{}
	if protoimpl.UnsafeEnabled {
		mi := &file_event_proto_msgTypes[6]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *EventStreamQuantum) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*EventStreamQuantum) ProtoMessage() {}

func (x *EventStreamQuantum) ProtoReflect() protoreflect.Message {
	mi := &file_event_proto_msgTypes[6]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use EventStreamQuantum.ProtoReflect.Descriptor instead.
func (*EventStreamQuantum) Descriptor() ([]byte, []int) {
	return file_event_proto_rawDescGZIP(), []int{6}
}

func (x *EventStreamQuantum) GetEvent() *Event {
	if x != nil {
		return x.Event
	}
	return nil
}

func (x *EventStreamQuantum) GetMeta() map[string]string {
	if x != nil {
		return x.Meta
	}
	return nil
}

type EntityStreamRequest_Filter struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	EventIdFrom string `protobuf:"bytes,1,opt,name=eventIdFrom,proto3" json:"eventIdFrom,omitempty"`
}

func (x *EntityStreamRequest_Filter) Reset() {
	*x = EntityStreamRequest_Filter{}
	if protoimpl.UnsafeEnabled {
		mi := &file_event_proto_msgTypes[7]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *EntityStreamRequest_Filter) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*EntityStreamRequest_Filter) ProtoMessage() {}

func (x *EntityStreamRequest_Filter) ProtoReflect() protoreflect.Message {
	mi := &file_event_proto_msgTypes[7]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use EntityStreamRequest_Filter.ProtoReflect.Descriptor instead.
func (*EntityStreamRequest_Filter) Descriptor() ([]byte, []int) {
	return file_event_proto_rawDescGZIP(), []int{3, 0}
}

func (x *EntityStreamRequest_Filter) GetEventIdFrom() string {
	if x != nil {
		return x.EventIdFrom
	}
	return ""
}

var File_event_proto protoreflect.FileDescriptor

var file_event_proto_rawDesc = []byte{
	0x0a, 0x0b, 0x65, 0x76, 0x65, 0x6e, 0x74, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x12, 0x15, 0x65,
	0x76, 0x65, 0x6e, 0x74, 0x68, 0x6f, 0x75, 0x73, 0x65, 0x2e, 0x67, 0x72, 0x70, 0x63, 0x2e, 0x65,
	0x76, 0x65, 0x6e, 0x74, 0x1a, 0x1f, 0x67, 0x6f, 0x6f, 0x67, 0x6c, 0x65, 0x2f, 0x70, 0x72, 0x6f,
	0x74, 0x6f, 0x62, 0x75, 0x66, 0x2f, 0x74, 0x69, 0x6d, 0x65, 0x73, 0x74, 0x61, 0x6d, 0x70, 0x2e,
	0x70, 0x72, 0x6f, 0x74, 0x6f, 0x22, 0xf7, 0x01, 0x0a, 0x05, 0x45, 0x76, 0x65, 0x6e, 0x74, 0x12,
	0x18, 0x0a, 0x07, 0x65, 0x76, 0x65, 0x6e, 0x74, 0x49, 0x64, 0x18, 0x01, 0x20, 0x01, 0x28, 0x09,
	0x52, 0x07, 0x65, 0x76, 0x65, 0x6e, 0x74, 0x49, 0x64, 0x12, 0x1c, 0x0a, 0x09, 0x65, 0x76, 0x65,
	0x6e, 0x74, 0x54, 0x79, 0x70, 0x65, 0x18, 0x02, 0x20, 0x01, 0x28, 0x09, 0x52, 0x09, 0x65, 0x76,
	0x65, 0x6e, 0x74, 0x54, 0x79, 0x70, 0x65, 0x12, 0x1e, 0x0a, 0x0a, 0x65, 0x6e, 0x74, 0x69, 0x74,
	0x79, 0x54, 0x79, 0x70, 0x65, 0x18, 0x03, 0x20, 0x01, 0x28, 0x09, 0x52, 0x0a, 0x65, 0x6e, 0x74,
	0x69, 0x74, 0x79, 0x54, 0x79, 0x70, 0x65, 0x12, 0x1a, 0x0a, 0x08, 0x65, 0x6e, 0x74, 0x69, 0x74,
	0x79, 0x49, 0x64, 0x18, 0x04, 0x20, 0x01, 0x28, 0x09, 0x52, 0x08, 0x65, 0x6e, 0x74, 0x69, 0x74,
	0x79, 0x49, 0x64, 0x12, 0x36, 0x0a, 0x08, 0x72, 0x65, 0x63, 0x6f, 0x72, 0x64, 0x65, 0x64, 0x18,
	0x05, 0x20, 0x01, 0x28, 0x0b, 0x32, 0x1a, 0x2e, 0x67, 0x6f, 0x6f, 0x67, 0x6c, 0x65, 0x2e, 0x70,
	0x72, 0x6f, 0x74, 0x6f, 0x62, 0x75, 0x66, 0x2e, 0x54, 0x69, 0x6d, 0x65, 0x73, 0x74, 0x61, 0x6d,
	0x70, 0x52, 0x08, 0x72, 0x65, 0x63, 0x6f, 0x72, 0x64, 0x65, 0x64, 0x12, 0x18, 0x0a, 0x07, 0x70,
	0x61, 0x79, 0x6c, 0x6f, 0x61, 0x64, 0x18, 0x06, 0x20, 0x01, 0x28, 0x09, 0x52, 0x07, 0x70, 0x61,
	0x79, 0x6c, 0x6f, 0x61, 0x64, 0x12, 0x28, 0x0a, 0x0f, 0x70, 0x72, 0x65, 0x76, 0x69, 0x6f, 0x75,
	0x73, 0x45, 0x76, 0x65, 0x6e, 0x74, 0x49, 0x64, 0x18, 0x07, 0x20, 0x01, 0x28, 0x09, 0x52, 0x0f,
	0x70, 0x72, 0x65, 0x76, 0x69, 0x6f, 0x75, 0x73, 0x45, 0x76, 0x65, 0x6e, 0x74, 0x49, 0x64, 0x22,
	0xa9, 0x01, 0x0a, 0x0b, 0x50, 0x75, 0x73, 0x68, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x12,
	0x26, 0x0a, 0x0e, 0x69, 0x64, 0x65, 0x6d, 0x70, 0x6f, 0x74, 0x65, 0x6e, 0x63, 0x79, 0x4b, 0x65,
	0x79, 0x18, 0x01, 0x20, 0x01, 0x28, 0x09, 0x52, 0x0e, 0x69, 0x64, 0x65, 0x6d, 0x70, 0x6f, 0x74,
	0x65, 0x6e, 0x63, 0x79, 0x4b, 0x65, 0x79, 0x12, 0x1c, 0x0a, 0x09, 0x65, 0x76, 0x65, 0x6e, 0x74,
	0x54, 0x79, 0x70, 0x65, 0x18, 0x02, 0x20, 0x01, 0x28, 0x09, 0x52, 0x09, 0x65, 0x76, 0x65, 0x6e,
	0x74, 0x54, 0x79, 0x70, 0x65, 0x12, 0x1e, 0x0a, 0x0a, 0x65, 0x6e, 0x74, 0x69, 0x74, 0x79, 0x54,
	0x79, 0x70, 0x65, 0x18, 0x03, 0x20, 0x01, 0x28, 0x09, 0x52, 0x0a, 0x65, 0x6e, 0x74, 0x69, 0x74,
	0x79, 0x54, 0x79, 0x70, 0x65, 0x12, 0x1a, 0x0a, 0x08, 0x65, 0x6e, 0x74, 0x69, 0x74, 0x79, 0x49,
	0x64, 0x18, 0x04, 0x20, 0x01, 0x28, 0x09, 0x52, 0x08, 0x65, 0x6e, 0x74, 0x69, 0x74, 0x79, 0x49,
	0x64, 0x12, 0x18, 0x0a, 0x07, 0x70, 0x61, 0x79, 0x6c, 0x6f, 0x61, 0x64, 0x18, 0x05, 0x20, 0x01,
	0x28, 0x09, 0x52, 0x07, 0x70, 0x61, 0x79, 0x6c, 0x6f, 0x61, 0x64, 0x22, 0x26, 0x0a, 0x0a, 0x47,
	0x65, 0x74, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x12, 0x18, 0x0a, 0x07, 0x65, 0x76, 0x65,
	0x6e, 0x74, 0x49, 0x64, 0x18, 0x01, 0x20, 0x01, 0x28, 0x09, 0x52, 0x07, 0x65, 0x76, 0x65, 0x6e,
	0x74, 0x49, 0x64, 0x22, 0xc8, 0x01, 0x0a, 0x13, 0x45, 0x6e, 0x74, 0x69, 0x74, 0x79, 0x53, 0x74,
	0x72, 0x65, 0x61, 0x6d, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x12, 0x1e, 0x0a, 0x0a, 0x65,
	0x6e, 0x74, 0x69, 0x74, 0x79, 0x54, 0x79, 0x70, 0x65, 0x18, 0x01, 0x20, 0x01, 0x28, 0x09, 0x52,
	0x0a, 0x65, 0x6e, 0x74, 0x69, 0x74, 0x79, 0x54, 0x79, 0x70, 0x65, 0x12, 0x1a, 0x0a, 0x08, 0x65,
	0x6e, 0x74, 0x69, 0x74, 0x79, 0x49, 0x64, 0x18, 0x02, 0x20, 0x01, 0x28, 0x09, 0x52, 0x08, 0x65,
	0x6e, 0x74, 0x69, 0x74, 0x79, 0x49, 0x64, 0x12, 0x49, 0x0a, 0x06, 0x66, 0x69, 0x6c, 0x74, 0x65,
	0x72, 0x18, 0x03, 0x20, 0x01, 0x28, 0x0b, 0x32, 0x31, 0x2e, 0x65, 0x76, 0x65, 0x6e, 0x74, 0x68,
	0x6f, 0x75, 0x73, 0x65, 0x2e, 0x67, 0x72, 0x70, 0x63, 0x2e, 0x65, 0x76, 0x65, 0x6e, 0x74, 0x2e,
	0x45, 0x6e, 0x74, 0x69, 0x74, 0x79, 0x53, 0x74, 0x72, 0x65, 0x61, 0x6d, 0x52, 0x65, 0x71, 0x75,
	0x65, 0x73, 0x74, 0x2e, 0x46, 0x69, 0x6c, 0x74, 0x65, 0x72, 0x52, 0x06, 0x66, 0x69, 0x6c, 0x74,
	0x65, 0x72, 0x1a, 0x2a, 0x0a, 0x06, 0x46, 0x69, 0x6c, 0x74, 0x65, 0x72, 0x12, 0x20, 0x0a, 0x0b,
	0x65, 0x76, 0x65, 0x6e, 0x74, 0x49, 0x64, 0x46, 0x72, 0x6f, 0x6d, 0x18, 0x01, 0x20, 0x01, 0x28,
	0x09, 0x52, 0x0b, 0x65, 0x76, 0x65, 0x6e, 0x74, 0x49, 0x64, 0x46, 0x72, 0x6f, 0x6d, 0x22, 0x75,
	0x0a, 0x13, 0x47, 0x6c, 0x6f, 0x62, 0x61, 0x6c, 0x53, 0x74, 0x72, 0x65, 0x61, 0x6d, 0x52, 0x65,
	0x71, 0x75, 0x65, 0x73, 0x74, 0x12, 0x20, 0x0a, 0x0b, 0x65, 0x76, 0x65, 0x6e, 0x74, 0x49, 0x64,
	0x46, 0x72, 0x6f, 0x6d, 0x18, 0x01, 0x20, 0x01, 0x28, 0x09, 0x52, 0x0b, 0x65, 0x76, 0x65, 0x6e,
	0x74, 0x49, 0x64, 0x46, 0x72, 0x6f, 0x6d, 0x12, 0x1e, 0x0a, 0x0a, 0x65, 0x6e, 0x74, 0x69, 0x74,
	0x79, 0x54, 0x79, 0x70, 0x65, 0x18, 0x02, 0x20, 0x03, 0x28, 0x09, 0x52, 0x0a, 0x65, 0x6e, 0x74,
	0x69, 0x74, 0x79, 0x54, 0x79, 0x70, 0x65, 0x12, 0x1c, 0x0a, 0x09, 0x65, 0x76, 0x65, 0x6e, 0x74,
	0x54, 0x79, 0x70, 0x65, 0x18, 0x03, 0x20, 0x03, 0x28, 0x09, 0x52, 0x09, 0x65, 0x76, 0x65, 0x6e,
	0x74, 0x54, 0x79, 0x70, 0x65, 0x22, 0x7e, 0x0a, 0x1c, 0x53, 0x75, 0x62, 0x73, 0x63, 0x72, 0x69,
	0x62, 0x65, 0x47, 0x6c, 0x6f, 0x62, 0x61, 0x6c, 0x53, 0x74, 0x72, 0x65, 0x61, 0x6d, 0x52, 0x65,
	0x71, 0x75, 0x65, 0x73, 0x74, 0x12, 0x20, 0x0a, 0x0b, 0x65, 0x76, 0x65, 0x6e, 0x74, 0x49, 0x64,
	0x46, 0x72, 0x6f, 0x6d, 0x18, 0x01, 0x20, 0x01, 0x28, 0x09, 0x52, 0x0b, 0x65, 0x76, 0x65, 0x6e,
	0x74, 0x49, 0x64, 0x46, 0x72, 0x6f, 0x6d, 0x12, 0x1e, 0x0a, 0x0a, 0x65, 0x6e, 0x74, 0x69, 0x74,
	0x79, 0x54, 0x79, 0x70, 0x65, 0x18, 0x02, 0x20, 0x03, 0x28, 0x09, 0x52, 0x0a, 0x65, 0x6e, 0x74,
	0x69, 0x74, 0x79, 0x54, 0x79, 0x70, 0x65, 0x12, 0x1c, 0x0a, 0x09, 0x65, 0x76, 0x65, 0x6e, 0x74,
	0x54, 0x79, 0x70, 0x65, 0x18, 0x03, 0x20, 0x03, 0x28, 0x09, 0x52, 0x09, 0x65, 0x76, 0x65, 0x6e,
	0x74, 0x54, 0x79, 0x70, 0x65, 0x22, 0xca, 0x01, 0x0a, 0x12, 0x45, 0x76, 0x65, 0x6e, 0x74, 0x53,
	0x74, 0x72, 0x65, 0x61, 0x6d, 0x51, 0x75, 0x61, 0x6e, 0x74, 0x75, 0x6d, 0x12, 0x32, 0x0a, 0x05,
	0x65, 0x76, 0x65, 0x6e, 0x74, 0x18, 0x01, 0x20, 0x01, 0x28, 0x0b, 0x32, 0x1c, 0x2e, 0x65, 0x76,
	0x65, 0x6e, 0x74, 0x68, 0x6f, 0x75, 0x73, 0x65, 0x2e, 0x67, 0x72, 0x70, 0x63, 0x2e, 0x65, 0x76,
	0x65, 0x6e, 0x74, 0x2e, 0x45, 0x76, 0x65, 0x6e, 0x74, 0x52, 0x05, 0x65, 0x76, 0x65, 0x6e, 0x74,
	0x12, 0x47, 0x0a, 0x04, 0x6d, 0x65, 0x74, 0x61, 0x18, 0x02, 0x20, 0x03, 0x28, 0x0b, 0x32, 0x33,
	0x2e, 0x65, 0x76, 0x65, 0x6e, 0x74, 0x68, 0x6f, 0x75, 0x73, 0x65, 0x2e, 0x67, 0x72, 0x70, 0x63,
	0x2e, 0x65, 0x76, 0x65, 0x6e, 0x74, 0x2e, 0x45, 0x76, 0x65, 0x6e, 0x74, 0x53, 0x74, 0x72, 0x65,
	0x61, 0x6d, 0x51, 0x75, 0x61, 0x6e, 0x74, 0x75, 0x6d, 0x2e, 0x4d, 0x65, 0x74, 0x61, 0x45, 0x6e,
	0x74, 0x72, 0x79, 0x52, 0x04, 0x6d, 0x65, 0x74, 0x61, 0x1a, 0x37, 0x0a, 0x09, 0x4d, 0x65, 0x74,
	0x61, 0x45, 0x6e, 0x74, 0x72, 0x79, 0x12, 0x10, 0x0a, 0x03, 0x6b, 0x65, 0x79, 0x18, 0x01, 0x20,
	0x01, 0x28, 0x09, 0x52, 0x03, 0x6b, 0x65, 0x79, 0x12, 0x14, 0x0a, 0x05, 0x76, 0x61, 0x6c, 0x75,
	0x65, 0x18, 0x02, 0x20, 0x01, 0x28, 0x09, 0x52, 0x05, 0x76, 0x61, 0x6c, 0x75, 0x65, 0x3a, 0x02,
	0x38, 0x01, 0x32, 0xee, 0x03, 0x0a, 0x03, 0x41, 0x50, 0x49, 0x12, 0x4a, 0x0a, 0x04, 0x50, 0x75,
	0x73, 0x68, 0x12, 0x22, 0x2e, 0x65, 0x76, 0x65, 0x6e, 0x74, 0x68, 0x6f, 0x75, 0x73, 0x65, 0x2e,
	0x67, 0x72, 0x70, 0x63, 0x2e, 0x65, 0x76, 0x65, 0x6e, 0x74, 0x2e, 0x50, 0x75, 0x73, 0x68, 0x52,
	0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x1a, 0x1c, 0x2e, 0x65, 0x76, 0x65, 0x6e, 0x74, 0x68, 0x6f,
	0x75, 0x73, 0x65, 0x2e, 0x67, 0x72, 0x70, 0x63, 0x2e, 0x65, 0x76, 0x65, 0x6e, 0x74, 0x2e, 0x45,
	0x76, 0x65, 0x6e, 0x74, 0x22, 0x00, 0x12, 0x48, 0x0a, 0x03, 0x47, 0x65, 0x74, 0x12, 0x21, 0x2e,
	0x65, 0x76, 0x65, 0x6e, 0x74, 0x68, 0x6f, 0x75, 0x73, 0x65, 0x2e, 0x67, 0x72, 0x70, 0x63, 0x2e,
	0x65, 0x76, 0x65, 0x6e, 0x74, 0x2e, 0x47, 0x65, 0x74, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74,
	0x1a, 0x1c, 0x2e, 0x65, 0x76, 0x65, 0x6e, 0x74, 0x68, 0x6f, 0x75, 0x73, 0x65, 0x2e, 0x67, 0x72,
	0x70, 0x63, 0x2e, 0x65, 0x76, 0x65, 0x6e, 0x74, 0x2e, 0x45, 0x76, 0x65, 0x6e, 0x74, 0x22, 0x00,
	0x12, 0x69, 0x0a, 0x0c, 0x45, 0x6e, 0x74, 0x69, 0x74, 0x79, 0x53, 0x74, 0x72, 0x65, 0x61, 0x6d,
	0x12, 0x2a, 0x2e, 0x65, 0x76, 0x65, 0x6e, 0x74, 0x68, 0x6f, 0x75, 0x73, 0x65, 0x2e, 0x67, 0x72,
	0x70, 0x63, 0x2e, 0x65, 0x76, 0x65, 0x6e, 0x74, 0x2e, 0x45, 0x6e, 0x74, 0x69, 0x74, 0x79, 0x53,
	0x74, 0x72, 0x65, 0x61, 0x6d, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x1a, 0x29, 0x2e, 0x65,
	0x76, 0x65, 0x6e, 0x74, 0x68, 0x6f, 0x75, 0x73, 0x65, 0x2e, 0x67, 0x72, 0x70, 0x63, 0x2e, 0x65,
	0x76, 0x65, 0x6e, 0x74, 0x2e, 0x45, 0x76, 0x65, 0x6e, 0x74, 0x53, 0x74, 0x72, 0x65, 0x61, 0x6d,
	0x51, 0x75, 0x61, 0x6e, 0x74, 0x75, 0x6d, 0x22, 0x00, 0x30, 0x01, 0x12, 0x69, 0x0a, 0x0c, 0x47,
	0x6c, 0x6f, 0x62, 0x61, 0x6c, 0x53, 0x74, 0x72, 0x65, 0x61, 0x6d, 0x12, 0x2a, 0x2e, 0x65, 0x76,
	0x65, 0x6e, 0x74, 0x68, 0x6f, 0x75, 0x73, 0x65, 0x2e, 0x67, 0x72, 0x70, 0x63, 0x2e, 0x65, 0x76,
	0x65, 0x6e, 0x74, 0x2e, 0x47, 0x6c, 0x6f, 0x62, 0x61, 0x6c, 0x53, 0x74, 0x72, 0x65, 0x61, 0x6d,
	0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x1a, 0x29, 0x2e, 0x65, 0x76, 0x65, 0x6e, 0x74, 0x68,
	0x6f, 0x75, 0x73, 0x65, 0x2e, 0x67, 0x72, 0x70, 0x63, 0x2e, 0x65, 0x76, 0x65, 0x6e, 0x74, 0x2e,
	0x45, 0x76, 0x65, 0x6e, 0x74, 0x53, 0x74, 0x72, 0x65, 0x61, 0x6d, 0x51, 0x75, 0x61, 0x6e, 0x74,
	0x75, 0x6d, 0x22, 0x00, 0x30, 0x01, 0x12, 0x7b, 0x0a, 0x15, 0x53, 0x75, 0x62, 0x73, 0x63, 0x72,
	0x69, 0x62, 0x65, 0x47, 0x6c, 0x6f, 0x62, 0x61, 0x6c, 0x53, 0x74, 0x72, 0x65, 0x61, 0x6d, 0x12,
	0x33, 0x2e, 0x65, 0x76, 0x65, 0x6e, 0x74, 0x68, 0x6f, 0x75, 0x73, 0x65, 0x2e, 0x67, 0x72, 0x70,
	0x63, 0x2e, 0x65, 0x76, 0x65, 0x6e, 0x74, 0x2e, 0x53, 0x75, 0x62, 0x73, 0x63, 0x72, 0x69, 0x62,
	0x65, 0x47, 0x6c, 0x6f, 0x62, 0x61, 0x6c, 0x53, 0x74, 0x72, 0x65, 0x61, 0x6d, 0x52, 0x65, 0x71,
	0x75, 0x65, 0x73, 0x74, 0x1a, 0x29, 0x2e, 0x65, 0x76, 0x65, 0x6e, 0x74, 0x68, 0x6f, 0x75, 0x73,
	0x65, 0x2e, 0x67, 0x72, 0x70, 0x63, 0x2e, 0x65, 0x76, 0x65, 0x6e, 0x74, 0x2e, 0x45, 0x76, 0x65,
	0x6e, 0x74, 0x53, 0x74, 0x72, 0x65, 0x61, 0x6d, 0x51, 0x75, 0x61, 0x6e, 0x74, 0x75, 0x6d, 0x22,
	0x00, 0x30, 0x01, 0x42, 0x1f, 0x5a, 0x1d, 0x65, 0x76, 0x65, 0x6e, 0x74, 0x68, 0x6f, 0x75, 0x73,
	0x65, 0x2f, 0x61, 0x70, 0x69, 0x2f, 0x63, 0x6f, 0x6d, 0x70, 0x69, 0x6c, 0x65, 0x64, 0x2f, 0x65,
	0x76, 0x65, 0x6e, 0x74, 0x62, 0x06, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x33,
}

var (
	file_event_proto_rawDescOnce sync.Once
	file_event_proto_rawDescData = file_event_proto_rawDesc
)

func file_event_proto_rawDescGZIP() []byte {
	file_event_proto_rawDescOnce.Do(func() {
		file_event_proto_rawDescData = protoimpl.X.CompressGZIP(file_event_proto_rawDescData)
	})
	return file_event_proto_rawDescData
}

var file_event_proto_msgTypes = make([]protoimpl.MessageInfo, 9)
var file_event_proto_goTypes = []interface{}{
	(*Event)(nil),                        // 0: eventhouse.grpc.event.Event
	(*PushRequest)(nil),                  // 1: eventhouse.grpc.event.PushRequest
	(*GetRequest)(nil),                   // 2: eventhouse.grpc.event.GetRequest
	(*EntityStreamRequest)(nil),          // 3: eventhouse.grpc.event.EntityStreamRequest
	(*GlobalStreamRequest)(nil),          // 4: eventhouse.grpc.event.GlobalStreamRequest
	(*SubscribeGlobalStreamRequest)(nil), // 5: eventhouse.grpc.event.SubscribeGlobalStreamRequest
	(*EventStreamQuantum)(nil),           // 6: eventhouse.grpc.event.EventStreamQuantum
	(*EntityStreamRequest_Filter)(nil),   // 7: eventhouse.grpc.event.EntityStreamRequest.Filter
	nil,                                  // 8: eventhouse.grpc.event.EventStreamQuantum.MetaEntry
	(*timestamp.Timestamp)(nil),          // 9: google.protobuf.Timestamp
}
var file_event_proto_depIdxs = []int32{
	9, // 0: eventhouse.grpc.event.Event.recorded:type_name -> google.protobuf.Timestamp
	7, // 1: eventhouse.grpc.event.EntityStreamRequest.filter:type_name -> eventhouse.grpc.event.EntityStreamRequest.Filter
	0, // 2: eventhouse.grpc.event.EventStreamQuantum.event:type_name -> eventhouse.grpc.event.Event
	8, // 3: eventhouse.grpc.event.EventStreamQuantum.meta:type_name -> eventhouse.grpc.event.EventStreamQuantum.MetaEntry
	1, // 4: eventhouse.grpc.event.API.Push:input_type -> eventhouse.grpc.event.PushRequest
	2, // 5: eventhouse.grpc.event.API.Get:input_type -> eventhouse.grpc.event.GetRequest
	3, // 6: eventhouse.grpc.event.API.EntityStream:input_type -> eventhouse.grpc.event.EntityStreamRequest
	4, // 7: eventhouse.grpc.event.API.GlobalStream:input_type -> eventhouse.grpc.event.GlobalStreamRequest
	5, // 8: eventhouse.grpc.event.API.SubscribeGlobalStream:input_type -> eventhouse.grpc.event.SubscribeGlobalStreamRequest
	0, // 9: eventhouse.grpc.event.API.Push:output_type -> eventhouse.grpc.event.Event
	0, // 10: eventhouse.grpc.event.API.Get:output_type -> eventhouse.grpc.event.Event
	6, // 11: eventhouse.grpc.event.API.EntityStream:output_type -> eventhouse.grpc.event.EventStreamQuantum
	6, // 12: eventhouse.grpc.event.API.GlobalStream:output_type -> eventhouse.grpc.event.EventStreamQuantum
	6, // 13: eventhouse.grpc.event.API.SubscribeGlobalStream:output_type -> eventhouse.grpc.event.EventStreamQuantum
	9, // [9:14] is the sub-list for method output_type
	4, // [4:9] is the sub-list for method input_type
	4, // [4:4] is the sub-list for extension type_name
	4, // [4:4] is the sub-list for extension extendee
	0, // [0:4] is the sub-list for field type_name
}

func init() { file_event_proto_init() }
func file_event_proto_init() {
	if File_event_proto != nil {
		return
	}
	if !protoimpl.UnsafeEnabled {
		file_event_proto_msgTypes[0].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*Event); i {
			case 0:
				return &v.state
			case 1:
				return &v.sizeCache
			case 2:
				return &v.unknownFields
			default:
				return nil
			}
		}
		file_event_proto_msgTypes[1].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*PushRequest); i {
			case 0:
				return &v.state
			case 1:
				return &v.sizeCache
			case 2:
				return &v.unknownFields
			default:
				return nil
			}
		}
		file_event_proto_msgTypes[2].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*GetRequest); i {
			case 0:
				return &v.state
			case 1:
				return &v.sizeCache
			case 2:
				return &v.unknownFields
			default:
				return nil
			}
		}
		file_event_proto_msgTypes[3].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*EntityStreamRequest); i {
			case 0:
				return &v.state
			case 1:
				return &v.sizeCache
			case 2:
				return &v.unknownFields
			default:
				return nil
			}
		}
		file_event_proto_msgTypes[4].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*GlobalStreamRequest); i {
			case 0:
				return &v.state
			case 1:
				return &v.sizeCache
			case 2:
				return &v.unknownFields
			default:
				return nil
			}
		}
		file_event_proto_msgTypes[5].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*SubscribeGlobalStreamRequest); i {
			case 0:
				return &v.state
			case 1:
				return &v.sizeCache
			case 2:
				return &v.unknownFields
			default:
				return nil
			}
		}
		file_event_proto_msgTypes[6].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*EventStreamQuantum); i {
			case 0:
				return &v.state
			case 1:
				return &v.sizeCache
			case 2:
				return &v.unknownFields
			default:
				return nil
			}
		}
		file_event_proto_msgTypes[7].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*EntityStreamRequest_Filter); i {
			case 0:
				return &v.state
			case 1:
				return &v.sizeCache
			case 2:
				return &v.unknownFields
			default:
				return nil
			}
		}
	}
	type x struct{}
	out := protoimpl.TypeBuilder{
		File: protoimpl.DescBuilder{
			GoPackagePath: reflect.TypeOf(x{}).PkgPath(),
			RawDescriptor: file_event_proto_rawDesc,
			NumEnums:      0,
			NumMessages:   9,
			NumExtensions: 0,
			NumServices:   1,
		},
		GoTypes:           file_event_proto_goTypes,
		DependencyIndexes: file_event_proto_depIdxs,
		MessageInfos:      file_event_proto_msgTypes,
	}.Build()
	File_event_proto = out.File
	file_event_proto_rawDesc = nil
	file_event_proto_goTypes = nil
	file_event_proto_depIdxs = nil
}

// Reference imports to suppress errors if they are not otherwise used.
var _ context.Context
var _ grpc.ClientConnInterface

// This is a compile-time assertion to ensure that this generated file
// is compatible with the grpc package it is being compiled against.
const _ = grpc.SupportPackageIsVersion6

// APIClient is the client API for API service.
//
// For semantics around ctx use and closing/ending streaming RPCs, please refer to https://godoc.org/google.golang.org/grpc#ClientConn.NewStream.
type APIClient interface {
	Push(ctx context.Context, in *PushRequest, opts ...grpc.CallOption) (*Event, error)
	Get(ctx context.Context, in *GetRequest, opts ...grpc.CallOption) (*Event, error)
	EntityStream(ctx context.Context, in *EntityStreamRequest, opts ...grpc.CallOption) (API_EntityStreamClient, error)
	GlobalStream(ctx context.Context, in *GlobalStreamRequest, opts ...grpc.CallOption) (API_GlobalStreamClient, error)
	SubscribeGlobalStream(ctx context.Context, in *SubscribeGlobalStreamRequest, opts ...grpc.CallOption) (API_SubscribeGlobalStreamClient, error)
}

type aPIClient struct {
	cc grpc.ClientConnInterface
}

func NewAPIClient(cc grpc.ClientConnInterface) APIClient {
	return &aPIClient{cc}
}

func (c *aPIClient) Push(ctx context.Context, in *PushRequest, opts ...grpc.CallOption) (*Event, error) {
	out := new(Event)
	err := c.cc.Invoke(ctx, "/eventhouse.grpc.event.API/Push", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *aPIClient) Get(ctx context.Context, in *GetRequest, opts ...grpc.CallOption) (*Event, error) {
	out := new(Event)
	err := c.cc.Invoke(ctx, "/eventhouse.grpc.event.API/Get", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *aPIClient) EntityStream(ctx context.Context, in *EntityStreamRequest, opts ...grpc.CallOption) (API_EntityStreamClient, error) {
	stream, err := c.cc.NewStream(ctx, &_API_serviceDesc.Streams[0], "/eventhouse.grpc.event.API/EntityStream", opts...)
	if err != nil {
		return nil, err
	}
	x := &aPIEntityStreamClient{stream}
	if err := x.ClientStream.SendMsg(in); err != nil {
		return nil, err
	}
	if err := x.ClientStream.CloseSend(); err != nil {
		return nil, err
	}
	return x, nil
}

type API_EntityStreamClient interface {
	Recv() (*EventStreamQuantum, error)
	grpc.ClientStream
}

type aPIEntityStreamClient struct {
	grpc.ClientStream
}

func (x *aPIEntityStreamClient) Recv() (*EventStreamQuantum, error) {
	m := new(EventStreamQuantum)
	if err := x.ClientStream.RecvMsg(m); err != nil {
		return nil, err
	}
	return m, nil
}

func (c *aPIClient) GlobalStream(ctx context.Context, in *GlobalStreamRequest, opts ...grpc.CallOption) (API_GlobalStreamClient, error) {
	stream, err := c.cc.NewStream(ctx, &_API_serviceDesc.Streams[1], "/eventhouse.grpc.event.API/GlobalStream", opts...)
	if err != nil {
		return nil, err
	}
	x := &aPIGlobalStreamClient{stream}
	if err := x.ClientStream.SendMsg(in); err != nil {
		return nil, err
	}
	if err := x.ClientStream.CloseSend(); err != nil {
		return nil, err
	}
	return x, nil
}

type API_GlobalStreamClient interface {
	Recv() (*EventStreamQuantum, error)
	grpc.ClientStream
}

type aPIGlobalStreamClient struct {
	grpc.ClientStream
}

func (x *aPIGlobalStreamClient) Recv() (*EventStreamQuantum, error) {
	m := new(EventStreamQuantum)
	if err := x.ClientStream.RecvMsg(m); err != nil {
		return nil, err
	}
	return m, nil
}

func (c *aPIClient) SubscribeGlobalStream(ctx context.Context, in *SubscribeGlobalStreamRequest, opts ...grpc.CallOption) (API_SubscribeGlobalStreamClient, error) {
	stream, err := c.cc.NewStream(ctx, &_API_serviceDesc.Streams[2], "/eventhouse.grpc.event.API/SubscribeGlobalStream", opts...)
	if err != nil {
		return nil, err
	}
	x := &aPISubscribeGlobalStreamClient{stream}
	if err := x.ClientStream.SendMsg(in); err != nil {
		return nil, err
	}
	if err := x.ClientStream.CloseSend(); err != nil {
		return nil, err
	}
	return x, nil
}

type API_SubscribeGlobalStreamClient interface {
	Recv() (*EventStreamQuantum, error)
	grpc.ClientStream
}

type aPISubscribeGlobalStreamClient struct {
	grpc.ClientStream
}

func (x *aPISubscribeGlobalStreamClient) Recv() (*EventStreamQuantum, error) {
	m := new(EventStreamQuantum)
	if err := x.ClientStream.RecvMsg(m); err != nil {
		return nil, err
	}
	return m, nil
}

// APIServer is the server API for API service.
type APIServer interface {
	Push(context.Context, *PushRequest) (*Event, error)
	Get(context.Context, *GetRequest) (*Event, error)
	EntityStream(*EntityStreamRequest, API_EntityStreamServer) error
	GlobalStream(*GlobalStreamRequest, API_GlobalStreamServer) error
	SubscribeGlobalStream(*SubscribeGlobalStreamRequest, API_SubscribeGlobalStreamServer) error
}

// UnimplementedAPIServer can be embedded to have forward compatible implementations.
type UnimplementedAPIServer struct {
}

func (*UnimplementedAPIServer) Push(context.Context, *PushRequest) (*Event, error) {
	return nil, status.Errorf(codes.Unimplemented, "method Push not implemented")
}
func (*UnimplementedAPIServer) Get(context.Context, *GetRequest) (*Event, error) {
	return nil, status.Errorf(codes.Unimplemented, "method Get not implemented")
}
func (*UnimplementedAPIServer) EntityStream(*EntityStreamRequest, API_EntityStreamServer) error {
	return status.Errorf(codes.Unimplemented, "method EntityStream not implemented")
}
func (*UnimplementedAPIServer) GlobalStream(*GlobalStreamRequest, API_GlobalStreamServer) error {
	return status.Errorf(codes.Unimplemented, "method GlobalStream not implemented")
}
func (*UnimplementedAPIServer) SubscribeGlobalStream(*SubscribeGlobalStreamRequest, API_SubscribeGlobalStreamServer) error {
	return status.Errorf(codes.Unimplemented, "method SubscribeGlobalStream not implemented")
}

func RegisterAPIServer(s *grpc.Server, srv APIServer) {
	s.RegisterService(&_API_serviceDesc, srv)
}

func _API_Push_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(PushRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(APIServer).Push(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/eventhouse.grpc.event.API/Push",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(APIServer).Push(ctx, req.(*PushRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _API_Get_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(GetRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(APIServer).Get(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/eventhouse.grpc.event.API/Get",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(APIServer).Get(ctx, req.(*GetRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _API_EntityStream_Handler(srv interface{}, stream grpc.ServerStream) error {
	m := new(EntityStreamRequest)
	if err := stream.RecvMsg(m); err != nil {
		return err
	}
	return srv.(APIServer).EntityStream(m, &aPIEntityStreamServer{stream})
}

type API_EntityStreamServer interface {
	Send(*EventStreamQuantum) error
	grpc.ServerStream
}

type aPIEntityStreamServer struct {
	grpc.ServerStream
}

func (x *aPIEntityStreamServer) Send(m *EventStreamQuantum) error {
	return x.ServerStream.SendMsg(m)
}

func _API_GlobalStream_Handler(srv interface{}, stream grpc.ServerStream) error {
	m := new(GlobalStreamRequest)
	if err := stream.RecvMsg(m); err != nil {
		return err
	}
	return srv.(APIServer).GlobalStream(m, &aPIGlobalStreamServer{stream})
}

type API_GlobalStreamServer interface {
	Send(*EventStreamQuantum) error
	grpc.ServerStream
}

type aPIGlobalStreamServer struct {
	grpc.ServerStream
}

func (x *aPIGlobalStreamServer) Send(m *EventStreamQuantum) error {
	return x.ServerStream.SendMsg(m)
}

func _API_SubscribeGlobalStream_Handler(srv interface{}, stream grpc.ServerStream) error {
	m := new(SubscribeGlobalStreamRequest)
	if err := stream.RecvMsg(m); err != nil {
		return err
	}
	return srv.(APIServer).SubscribeGlobalStream(m, &aPISubscribeGlobalStreamServer{stream})
}

type API_SubscribeGlobalStreamServer interface {
	Send(*EventStreamQuantum) error
	grpc.ServerStream
}

type aPISubscribeGlobalStreamServer struct {
	grpc.ServerStream
}

func (x *aPISubscribeGlobalStreamServer) Send(m *EventStreamQuantum) error {
	return x.ServerStream.SendMsg(m)
}

var _API_serviceDesc = grpc.ServiceDesc{
	ServiceName: "eventhouse.grpc.event.API",
	HandlerType: (*APIServer)(nil),
	Methods: []grpc.MethodDesc{
		{
			MethodName: "Push",
			Handler:    _API_Push_Handler,
		},
		{
			MethodName: "Get",
			Handler:    _API_Get_Handler,
		},
	},
	Streams: []grpc.StreamDesc{
		{
			StreamName:    "EntityStream",
			Handler:       _API_EntityStream_Handler,
			ServerStreams: true,
		},
		{
			StreamName:    "GlobalStream",
			Handler:       _API_GlobalStream_Handler,
			ServerStreams: true,
		},
		{
			StreamName:    "SubscribeGlobalStream",
			Handler:       _API_SubscribeGlobalStream_Handler,
			ServerStreams: true,
		},
	},
	Metadata: "event.proto",
}
