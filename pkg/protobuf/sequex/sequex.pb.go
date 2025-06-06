// Code generated by protoc-gen-go. DO NOT EDIT.
// versions:
// 	protoc-gen-go v1.36.5
// 	protoc        (unknown)
// source: protobuf/sequex/sequex.proto

package sequex

import (
	protoreflect "google.golang.org/protobuf/reflect/protoreflect"
	protoimpl "google.golang.org/protobuf/runtime/protoimpl"
	timestamppb "google.golang.org/protobuf/types/known/timestamppb"
	reflect "reflect"
	sync "sync"
	unsafe "unsafe"
)

const (
	// Verify that this generated code is sufficiently up-to-date.
	_ = protoimpl.EnforceVersion(20 - protoimpl.MinVersion)
	// Verify that runtime/protoimpl is sufficiently up-to-date.
	_ = protoimpl.EnforceVersion(protoimpl.MaxVersion - 20)
)

// Enum for Event Type
type EventType int32

const (
	// Unknown
	EventType_UNKNOWN_EVENT EventType = 0
	// Kline Events
	EventType_KLINE_UPDATE   EventType = 1
	EventType_KLINE_ACK      EventType = 2
	EventType_KLINE_FAILED   EventType = 3
	EventType_KLINE_FINISHED EventType = 4
	// Order Events
	EventType_ORDER_UPDATE   EventType = 5
	EventType_ORDER_ACK      EventType = 6
	EventType_ORDER_FAILED   EventType = 7
	EventType_ORDER_FINISHED EventType = 8
	// Execution Events
	EventType_EXECUTION_UPDATE   EventType = 9
	EventType_EXECUTION_ACK      EventType = 10
	EventType_EXECUTION_FAILED   EventType = 11
	EventType_EXECUTION_FINISHED EventType = 12
)

// Enum value maps for EventType.
var (
	EventType_name = map[int32]string{
		0:  "UNKNOWN_EVENT",
		1:  "KLINE_UPDATE",
		2:  "KLINE_ACK",
		3:  "KLINE_FAILED",
		4:  "KLINE_FINISHED",
		5:  "ORDER_UPDATE",
		6:  "ORDER_ACK",
		7:  "ORDER_FAILED",
		8:  "ORDER_FINISHED",
		9:  "EXECUTION_UPDATE",
		10: "EXECUTION_ACK",
		11: "EXECUTION_FAILED",
		12: "EXECUTION_FINISHED",
	}
	EventType_value = map[string]int32{
		"UNKNOWN_EVENT":      0,
		"KLINE_UPDATE":       1,
		"KLINE_ACK":          2,
		"KLINE_FAILED":       3,
		"KLINE_FINISHED":     4,
		"ORDER_UPDATE":       5,
		"ORDER_ACK":          6,
		"ORDER_FAILED":       7,
		"ORDER_FINISHED":     8,
		"EXECUTION_UPDATE":   9,
		"EXECUTION_ACK":      10,
		"EXECUTION_FAILED":   11,
		"EXECUTION_FINISHED": 12,
	}
)

func (x EventType) Enum() *EventType {
	p := new(EventType)
	*p = x
	return p
}

func (x EventType) String() string {
	return protoimpl.X.EnumStringOf(x.Descriptor(), protoreflect.EnumNumber(x))
}

func (EventType) Descriptor() protoreflect.EnumDescriptor {
	return file_protobuf_sequex_sequex_proto_enumTypes[0].Descriptor()
}

func (EventType) Type() protoreflect.EnumType {
	return &file_protobuf_sequex_sequex_proto_enumTypes[0]
}

func (x EventType) Number() protoreflect.EnumNumber {
	return protoreflect.EnumNumber(x)
}

// Deprecated: Use EventType.Descriptor instead.
func (EventType) EnumDescriptor() ([]byte, []int) {
	return file_protobuf_sequex_sequex_proto_rawDescGZIP(), []int{0}
}

// Enum for Event Source
type EventSource int32

const (
	EventSource_UNKNOWN_SOURCE EventSource = 0
	EventSource_SEQUEX         EventSource = 1
	EventSource_STRATEGIST     EventSource = 2 // Add more event sources as needed
)

// Enum value maps for EventSource.
var (
	EventSource_name = map[int32]string{
		0: "UNKNOWN_SOURCE",
		1: "SEQUEX",
		2: "STRATEGIST",
	}
	EventSource_value = map[string]int32{
		"UNKNOWN_SOURCE": 0,
		"SEQUEX":         1,
		"STRATEGIST":     2,
	}
)

func (x EventSource) Enum() *EventSource {
	p := new(EventSource)
	*p = x
	return p
}

func (x EventSource) String() string {
	return protoimpl.X.EnumStringOf(x.Descriptor(), protoreflect.EnumNumber(x))
}

func (EventSource) Descriptor() protoreflect.EnumDescriptor {
	return file_protobuf_sequex_sequex_proto_enumTypes[1].Descriptor()
}

func (EventSource) Type() protoreflect.EnumType {
	return &file_protobuf_sequex_sequex_proto_enumTypes[1]
}

func (x EventSource) Number() protoreflect.EnumNumber {
	return protoreflect.EnumNumber(x)
}

// Deprecated: Use EventSource.Descriptor instead.
func (EventSource) EnumDescriptor() ([]byte, []int) {
	return file_protobuf_sequex_sequex_proto_rawDescGZIP(), []int{1}
}

// Message representing an Event
type Event struct {
	state         protoimpl.MessageState `protogen:"open.v1"`
	Id            string                 `protobuf:"bytes,1,opt,name=id,proto3" json:"id,omitempty"` // UUID
	Type          EventType              `protobuf:"varint,2,opt,name=type,proto3,enum=sequex.EventType" json:"type,omitempty"`
	Source        EventSource            `protobuf:"varint,3,opt,name=source,proto3,enum=sequex.EventSource" json:"source,omitempty"`
	CreatedAt     *timestamppb.Timestamp `protobuf:"bytes,4,opt,name=created_at,json=createdAt,proto3" json:"created_at,omitempty"`
	Payload       []byte                 `protobuf:"bytes,6,opt,name=payload,proto3" json:"payload,omitempty"`
	unknownFields protoimpl.UnknownFields
	sizeCache     protoimpl.SizeCache
}

func (x *Event) Reset() {
	*x = Event{}
	mi := &file_protobuf_sequex_sequex_proto_msgTypes[0]
	ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
	ms.StoreMessageInfo(mi)
}

func (x *Event) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*Event) ProtoMessage() {}

func (x *Event) ProtoReflect() protoreflect.Message {
	mi := &file_protobuf_sequex_sequex_proto_msgTypes[0]
	if x != nil {
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
	return file_protobuf_sequex_sequex_proto_rawDescGZIP(), []int{0}
}

func (x *Event) GetId() string {
	if x != nil {
		return x.Id
	}
	return ""
}

func (x *Event) GetType() EventType {
	if x != nil {
		return x.Type
	}
	return EventType_UNKNOWN_EVENT
}

func (x *Event) GetSource() EventSource {
	if x != nil {
		return x.Source
	}
	return EventSource_UNKNOWN_SOURCE
}

func (x *Event) GetCreatedAt() *timestamppb.Timestamp {
	if x != nil {
		return x.CreatedAt
	}
	return nil
}

func (x *Event) GetPayload() []byte {
	if x != nil {
		return x.Payload
	}
	return nil
}

var File_protobuf_sequex_sequex_proto protoreflect.FileDescriptor

var file_protobuf_sequex_sequex_proto_rawDesc = string([]byte{
	0x0a, 0x1c, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x62, 0x75, 0x66, 0x2f, 0x73, 0x65, 0x71, 0x75, 0x65,
	0x78, 0x2f, 0x73, 0x65, 0x71, 0x75, 0x65, 0x78, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x12, 0x06,
	0x73, 0x65, 0x71, 0x75, 0x65, 0x78, 0x1a, 0x1f, 0x67, 0x6f, 0x6f, 0x67, 0x6c, 0x65, 0x2f, 0x70,
	0x72, 0x6f, 0x74, 0x6f, 0x62, 0x75, 0x66, 0x2f, 0x74, 0x69, 0x6d, 0x65, 0x73, 0x74, 0x61, 0x6d,
	0x70, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x22, 0xc0, 0x01, 0x0a, 0x05, 0x45, 0x76, 0x65, 0x6e,
	0x74, 0x12, 0x0e, 0x0a, 0x02, 0x69, 0x64, 0x18, 0x01, 0x20, 0x01, 0x28, 0x09, 0x52, 0x02, 0x69,
	0x64, 0x12, 0x25, 0x0a, 0x04, 0x74, 0x79, 0x70, 0x65, 0x18, 0x02, 0x20, 0x01, 0x28, 0x0e, 0x32,
	0x11, 0x2e, 0x73, 0x65, 0x71, 0x75, 0x65, 0x78, 0x2e, 0x45, 0x76, 0x65, 0x6e, 0x74, 0x54, 0x79,
	0x70, 0x65, 0x52, 0x04, 0x74, 0x79, 0x70, 0x65, 0x12, 0x2b, 0x0a, 0x06, 0x73, 0x6f, 0x75, 0x72,
	0x63, 0x65, 0x18, 0x03, 0x20, 0x01, 0x28, 0x0e, 0x32, 0x13, 0x2e, 0x73, 0x65, 0x71, 0x75, 0x65,
	0x78, 0x2e, 0x45, 0x76, 0x65, 0x6e, 0x74, 0x53, 0x6f, 0x75, 0x72, 0x63, 0x65, 0x52, 0x06, 0x73,
	0x6f, 0x75, 0x72, 0x63, 0x65, 0x12, 0x39, 0x0a, 0x0a, 0x63, 0x72, 0x65, 0x61, 0x74, 0x65, 0x64,
	0x5f, 0x61, 0x74, 0x18, 0x04, 0x20, 0x01, 0x28, 0x0b, 0x32, 0x1a, 0x2e, 0x67, 0x6f, 0x6f, 0x67,
	0x6c, 0x65, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x62, 0x75, 0x66, 0x2e, 0x54, 0x69, 0x6d, 0x65,
	0x73, 0x74, 0x61, 0x6d, 0x70, 0x52, 0x09, 0x63, 0x72, 0x65, 0x61, 0x74, 0x65, 0x64, 0x41, 0x74,
	0x12, 0x18, 0x0a, 0x07, 0x70, 0x61, 0x79, 0x6c, 0x6f, 0x61, 0x64, 0x18, 0x06, 0x20, 0x01, 0x28,
	0x0c, 0x52, 0x07, 0x70, 0x61, 0x79, 0x6c, 0x6f, 0x61, 0x64, 0x2a, 0x83, 0x02, 0x0a, 0x09, 0x45,
	0x76, 0x65, 0x6e, 0x74, 0x54, 0x79, 0x70, 0x65, 0x12, 0x11, 0x0a, 0x0d, 0x55, 0x4e, 0x4b, 0x4e,
	0x4f, 0x57, 0x4e, 0x5f, 0x45, 0x56, 0x45, 0x4e, 0x54, 0x10, 0x00, 0x12, 0x10, 0x0a, 0x0c, 0x4b,
	0x4c, 0x49, 0x4e, 0x45, 0x5f, 0x55, 0x50, 0x44, 0x41, 0x54, 0x45, 0x10, 0x01, 0x12, 0x0d, 0x0a,
	0x09, 0x4b, 0x4c, 0x49, 0x4e, 0x45, 0x5f, 0x41, 0x43, 0x4b, 0x10, 0x02, 0x12, 0x10, 0x0a, 0x0c,
	0x4b, 0x4c, 0x49, 0x4e, 0x45, 0x5f, 0x46, 0x41, 0x49, 0x4c, 0x45, 0x44, 0x10, 0x03, 0x12, 0x12,
	0x0a, 0x0e, 0x4b, 0x4c, 0x49, 0x4e, 0x45, 0x5f, 0x46, 0x49, 0x4e, 0x49, 0x53, 0x48, 0x45, 0x44,
	0x10, 0x04, 0x12, 0x10, 0x0a, 0x0c, 0x4f, 0x52, 0x44, 0x45, 0x52, 0x5f, 0x55, 0x50, 0x44, 0x41,
	0x54, 0x45, 0x10, 0x05, 0x12, 0x0d, 0x0a, 0x09, 0x4f, 0x52, 0x44, 0x45, 0x52, 0x5f, 0x41, 0x43,
	0x4b, 0x10, 0x06, 0x12, 0x10, 0x0a, 0x0c, 0x4f, 0x52, 0x44, 0x45, 0x52, 0x5f, 0x46, 0x41, 0x49,
	0x4c, 0x45, 0x44, 0x10, 0x07, 0x12, 0x12, 0x0a, 0x0e, 0x4f, 0x52, 0x44, 0x45, 0x52, 0x5f, 0x46,
	0x49, 0x4e, 0x49, 0x53, 0x48, 0x45, 0x44, 0x10, 0x08, 0x12, 0x14, 0x0a, 0x10, 0x45, 0x58, 0x45,
	0x43, 0x55, 0x54, 0x49, 0x4f, 0x4e, 0x5f, 0x55, 0x50, 0x44, 0x41, 0x54, 0x45, 0x10, 0x09, 0x12,
	0x11, 0x0a, 0x0d, 0x45, 0x58, 0x45, 0x43, 0x55, 0x54, 0x49, 0x4f, 0x4e, 0x5f, 0x41, 0x43, 0x4b,
	0x10, 0x0a, 0x12, 0x14, 0x0a, 0x10, 0x45, 0x58, 0x45, 0x43, 0x55, 0x54, 0x49, 0x4f, 0x4e, 0x5f,
	0x46, 0x41, 0x49, 0x4c, 0x45, 0x44, 0x10, 0x0b, 0x12, 0x16, 0x0a, 0x12, 0x45, 0x58, 0x45, 0x43,
	0x55, 0x54, 0x49, 0x4f, 0x4e, 0x5f, 0x46, 0x49, 0x4e, 0x49, 0x53, 0x48, 0x45, 0x44, 0x10, 0x0c,
	0x2a, 0x3d, 0x0a, 0x0b, 0x45, 0x76, 0x65, 0x6e, 0x74, 0x53, 0x6f, 0x75, 0x72, 0x63, 0x65, 0x12,
	0x12, 0x0a, 0x0e, 0x55, 0x4e, 0x4b, 0x4e, 0x4f, 0x57, 0x4e, 0x5f, 0x53, 0x4f, 0x55, 0x52, 0x43,
	0x45, 0x10, 0x00, 0x12, 0x0a, 0x0a, 0x06, 0x53, 0x45, 0x51, 0x55, 0x45, 0x58, 0x10, 0x01, 0x12,
	0x0e, 0x0a, 0x0a, 0x53, 0x54, 0x52, 0x41, 0x54, 0x45, 0x47, 0x49, 0x53, 0x54, 0x10, 0x02, 0x32,
	0x3c, 0x0a, 0x0d, 0x53, 0x65, 0x71, 0x75, 0x65, 0x78, 0x53, 0x65, 0x72, 0x76, 0x69, 0x63, 0x65,
	0x12, 0x2b, 0x0a, 0x07, 0x4f, 0x6e, 0x45, 0x76, 0x65, 0x6e, 0x74, 0x12, 0x0d, 0x2e, 0x73, 0x65,
	0x71, 0x75, 0x65, 0x78, 0x2e, 0x45, 0x76, 0x65, 0x6e, 0x74, 0x1a, 0x0d, 0x2e, 0x73, 0x65, 0x71,
	0x75, 0x65, 0x78, 0x2e, 0x45, 0x76, 0x65, 0x6e, 0x74, 0x28, 0x01, 0x30, 0x01, 0x42, 0x13, 0x5a,
	0x11, 0x2e, 0x2f, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x62, 0x75, 0x66, 0x2f, 0x73, 0x65, 0x71, 0x75,
	0x65, 0x78, 0x62, 0x06, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x33,
})

var (
	file_protobuf_sequex_sequex_proto_rawDescOnce sync.Once
	file_protobuf_sequex_sequex_proto_rawDescData []byte
)

func file_protobuf_sequex_sequex_proto_rawDescGZIP() []byte {
	file_protobuf_sequex_sequex_proto_rawDescOnce.Do(func() {
		file_protobuf_sequex_sequex_proto_rawDescData = protoimpl.X.CompressGZIP(unsafe.Slice(unsafe.StringData(file_protobuf_sequex_sequex_proto_rawDesc), len(file_protobuf_sequex_sequex_proto_rawDesc)))
	})
	return file_protobuf_sequex_sequex_proto_rawDescData
}

var file_protobuf_sequex_sequex_proto_enumTypes = make([]protoimpl.EnumInfo, 2)
var file_protobuf_sequex_sequex_proto_msgTypes = make([]protoimpl.MessageInfo, 1)
var file_protobuf_sequex_sequex_proto_goTypes = []any{
	(EventType)(0),                // 0: sequex.EventType
	(EventSource)(0),              // 1: sequex.EventSource
	(*Event)(nil),                 // 2: sequex.Event
	(*timestamppb.Timestamp)(nil), // 3: google.protobuf.Timestamp
}
var file_protobuf_sequex_sequex_proto_depIdxs = []int32{
	0, // 0: sequex.Event.type:type_name -> sequex.EventType
	1, // 1: sequex.Event.source:type_name -> sequex.EventSource
	3, // 2: sequex.Event.created_at:type_name -> google.protobuf.Timestamp
	2, // 3: sequex.SequexService.OnEvent:input_type -> sequex.Event
	2, // 4: sequex.SequexService.OnEvent:output_type -> sequex.Event
	4, // [4:5] is the sub-list for method output_type
	3, // [3:4] is the sub-list for method input_type
	3, // [3:3] is the sub-list for extension type_name
	3, // [3:3] is the sub-list for extension extendee
	0, // [0:3] is the sub-list for field type_name
}

func init() { file_protobuf_sequex_sequex_proto_init() }
func file_protobuf_sequex_sequex_proto_init() {
	if File_protobuf_sequex_sequex_proto != nil {
		return
	}
	type x struct{}
	out := protoimpl.TypeBuilder{
		File: protoimpl.DescBuilder{
			GoPackagePath: reflect.TypeOf(x{}).PkgPath(),
			RawDescriptor: unsafe.Slice(unsafe.StringData(file_protobuf_sequex_sequex_proto_rawDesc), len(file_protobuf_sequex_sequex_proto_rawDesc)),
			NumEnums:      2,
			NumMessages:   1,
			NumExtensions: 0,
			NumServices:   1,
		},
		GoTypes:           file_protobuf_sequex_sequex_proto_goTypes,
		DependencyIndexes: file_protobuf_sequex_sequex_proto_depIdxs,
		EnumInfos:         file_protobuf_sequex_sequex_proto_enumTypes,
		MessageInfos:      file_protobuf_sequex_sequex_proto_msgTypes,
	}.Build()
	File_protobuf_sequex_sequex_proto = out.File
	file_protobuf_sequex_sequex_proto_goTypes = nil
	file_protobuf_sequex_sequex_proto_depIdxs = nil
}
