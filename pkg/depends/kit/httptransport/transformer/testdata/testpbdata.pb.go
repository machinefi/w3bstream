// Code generated by protoc-gen-go. DO NOT EDIT.
// versions:
// 	protoc-gen-go v1.28.1
// 	protoc        v3.21.12
// source: testpbdata.proto

package testdatapb

import (
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

type Header struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	EventType string `protobuf:"bytes,1,opt,name=event_type,json=eventType,proto3" json:"event_type,omitempty"` // event type
	PubId     string `protobuf:"bytes,2,opt,name=pub_id,json=pubId,proto3" json:"pub_id,omitempty"`             // the unique identifier for publisher
	Token     string `protobuf:"bytes,3,opt,name=token,proto3" json:"token,omitempty"`                          // for validation message
	PubTime   int64  `protobuf:"varint,4,opt,name=pub_time,json=pubTime,proto3" json:"pub_time,omitempty"`      // event pub timestamp
	EventId   string `protobuf:"bytes,5,opt,name=event_id,json=eventId,proto3" json:"event_id,omitempty"`       // event id for tracing
}

func (x *Header) Reset() {
	*x = Header{}
	if protoimpl.UnsafeEnabled {
		mi := &file_testpbdata_proto_msgTypes[0]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *Header) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*Header) ProtoMessage() {}

func (x *Header) ProtoReflect() protoreflect.Message {
	mi := &file_testpbdata_proto_msgTypes[0]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use Header.ProtoReflect.Descriptor instead.
func (*Header) Descriptor() ([]byte, []int) {
	return file_testpbdata_proto_rawDescGZIP(), []int{0}
}

func (x *Header) GetEventType() string {
	if x != nil {
		return x.EventType
	}
	return ""
}

func (x *Header) GetPubId() string {
	if x != nil {
		return x.PubId
	}
	return ""
}

func (x *Header) GetToken() string {
	if x != nil {
		return x.Token
	}
	return ""
}

func (x *Header) GetPubTime() int64 {
	if x != nil {
		return x.PubTime
	}
	return 0
}

func (x *Header) GetEventId() string {
	if x != nil {
		return x.EventId
	}
	return ""
}

type Event struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Header  *Header `protobuf:"bytes,1,opt,name=header,proto3" json:"header,omitempty"`
	Payload []byte  `protobuf:"bytes,2,opt,name=payload,proto3" json:"payload,omitempty"`
}

func (x *Event) Reset() {
	*x = Event{}
	if protoimpl.UnsafeEnabled {
		mi := &file_testpbdata_proto_msgTypes[1]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *Event) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*Event) ProtoMessage() {}

func (x *Event) ProtoReflect() protoreflect.Message {
	mi := &file_testpbdata_proto_msgTypes[1]
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
	return file_testpbdata_proto_rawDescGZIP(), []int{1}
}

func (x *Event) GetHeader() *Header {
	if x != nil {
		return x.Header
	}
	return nil
}

func (x *Event) GetPayload() []byte {
	if x != nil {
		return x.Payload
	}
	return nil
}

type Events struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Data []*Event `protobuf:"bytes,1,rep,name=data,proto3" json:"data,omitempty"`
}

func (x *Events) Reset() {
	*x = Events{}
	if protoimpl.UnsafeEnabled {
		mi := &file_testpbdata_proto_msgTypes[2]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *Events) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*Events) ProtoMessage() {}

func (x *Events) ProtoReflect() protoreflect.Message {
	mi := &file_testpbdata_proto_msgTypes[2]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use Events.ProtoReflect.Descriptor instead.
func (*Events) Descriptor() ([]byte, []int) {
	return file_testpbdata_proto_rawDescGZIP(), []int{2}
}

func (x *Events) GetData() []*Event {
	if x != nil {
		return x.Data
	}
	return nil
}

var File_testpbdata_proto protoreflect.FileDescriptor

var file_testpbdata_proto_rawDesc = []byte{
	0x0a, 0x10, 0x74, 0x65, 0x73, 0x74, 0x70, 0x62, 0x64, 0x61, 0x74, 0x61, 0x2e, 0x70, 0x72, 0x6f,
	0x74, 0x6f, 0x12, 0x0a, 0x74, 0x65, 0x73, 0x74, 0x64, 0x61, 0x74, 0x61, 0x70, 0x62, 0x22, 0x8a,
	0x01, 0x0a, 0x06, 0x48, 0x65, 0x61, 0x64, 0x65, 0x72, 0x12, 0x1d, 0x0a, 0x0a, 0x65, 0x76, 0x65,
	0x6e, 0x74, 0x5f, 0x74, 0x79, 0x70, 0x65, 0x18, 0x01, 0x20, 0x01, 0x28, 0x09, 0x52, 0x09, 0x65,
	0x76, 0x65, 0x6e, 0x74, 0x54, 0x79, 0x70, 0x65, 0x12, 0x15, 0x0a, 0x06, 0x70, 0x75, 0x62, 0x5f,
	0x69, 0x64, 0x18, 0x02, 0x20, 0x01, 0x28, 0x09, 0x52, 0x05, 0x70, 0x75, 0x62, 0x49, 0x64, 0x12,
	0x14, 0x0a, 0x05, 0x74, 0x6f, 0x6b, 0x65, 0x6e, 0x18, 0x03, 0x20, 0x01, 0x28, 0x09, 0x52, 0x05,
	0x74, 0x6f, 0x6b, 0x65, 0x6e, 0x12, 0x19, 0x0a, 0x08, 0x70, 0x75, 0x62, 0x5f, 0x74, 0x69, 0x6d,
	0x65, 0x18, 0x04, 0x20, 0x01, 0x28, 0x03, 0x52, 0x07, 0x70, 0x75, 0x62, 0x54, 0x69, 0x6d, 0x65,
	0x12, 0x19, 0x0a, 0x08, 0x65, 0x76, 0x65, 0x6e, 0x74, 0x5f, 0x69, 0x64, 0x18, 0x05, 0x20, 0x01,
	0x28, 0x09, 0x52, 0x07, 0x65, 0x76, 0x65, 0x6e, 0x74, 0x49, 0x64, 0x22, 0x4d, 0x0a, 0x05, 0x45,
	0x76, 0x65, 0x6e, 0x74, 0x12, 0x2a, 0x0a, 0x06, 0x68, 0x65, 0x61, 0x64, 0x65, 0x72, 0x18, 0x01,
	0x20, 0x01, 0x28, 0x0b, 0x32, 0x12, 0x2e, 0x74, 0x65, 0x73, 0x74, 0x64, 0x61, 0x74, 0x61, 0x70,
	0x62, 0x2e, 0x48, 0x65, 0x61, 0x64, 0x65, 0x72, 0x52, 0x06, 0x68, 0x65, 0x61, 0x64, 0x65, 0x72,
	0x12, 0x18, 0x0a, 0x07, 0x70, 0x61, 0x79, 0x6c, 0x6f, 0x61, 0x64, 0x18, 0x02, 0x20, 0x01, 0x28,
	0x0c, 0x52, 0x07, 0x70, 0x61, 0x79, 0x6c, 0x6f, 0x61, 0x64, 0x22, 0x2f, 0x0a, 0x06, 0x45, 0x76,
	0x65, 0x6e, 0x74, 0x73, 0x12, 0x25, 0x0a, 0x04, 0x64, 0x61, 0x74, 0x61, 0x18, 0x01, 0x20, 0x03,
	0x28, 0x0b, 0x32, 0x11, 0x2e, 0x74, 0x65, 0x73, 0x74, 0x64, 0x61, 0x74, 0x61, 0x70, 0x62, 0x2e,
	0x45, 0x76, 0x65, 0x6e, 0x74, 0x52, 0x04, 0x64, 0x61, 0x74, 0x61, 0x42, 0x0f, 0x5a, 0x0d, 0x2e,
	0x2f, 0x3b, 0x74, 0x65, 0x73, 0x74, 0x64, 0x61, 0x74, 0x61, 0x70, 0x62, 0x62, 0x06, 0x70, 0x72,
	0x6f, 0x74, 0x6f, 0x33,
}

var (
	file_testpbdata_proto_rawDescOnce sync.Once
	file_testpbdata_proto_rawDescData = file_testpbdata_proto_rawDesc
)

func file_testpbdata_proto_rawDescGZIP() []byte {
	file_testpbdata_proto_rawDescOnce.Do(func() {
		file_testpbdata_proto_rawDescData = protoimpl.X.CompressGZIP(file_testpbdata_proto_rawDescData)
	})
	return file_testpbdata_proto_rawDescData
}

var file_testpbdata_proto_msgTypes = make([]protoimpl.MessageInfo, 3)
var file_testpbdata_proto_goTypes = []interface{}{
	(*Header)(nil), // 0: testdatapb.Header
	(*Event)(nil),  // 1: testdatapb.Event
	(*Events)(nil), // 2: testdatapb.Events
}
var file_testpbdata_proto_depIdxs = []int32{
	0, // 0: testdatapb.Event.header:type_name -> testdatapb.Header
	1, // 1: testdatapb.Events.data:type_name -> testdatapb.Event
	2, // [2:2] is the sub-list for method output_type
	2, // [2:2] is the sub-list for method input_type
	2, // [2:2] is the sub-list for extension type_name
	2, // [2:2] is the sub-list for extension extendee
	0, // [0:2] is the sub-list for field type_name
}

func init() { file_testpbdata_proto_init() }
func file_testpbdata_proto_init() {
	if File_testpbdata_proto != nil {
		return
	}
	if !protoimpl.UnsafeEnabled {
		file_testpbdata_proto_msgTypes[0].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*Header); i {
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
		file_testpbdata_proto_msgTypes[1].Exporter = func(v interface{}, i int) interface{} {
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
		file_testpbdata_proto_msgTypes[2].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*Events); i {
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
			RawDescriptor: file_testpbdata_proto_rawDesc,
			NumEnums:      0,
			NumMessages:   3,
			NumExtensions: 0,
			NumServices:   0,
		},
		GoTypes:           file_testpbdata_proto_goTypes,
		DependencyIndexes: file_testpbdata_proto_depIdxs,
		MessageInfos:      file_testpbdata_proto_msgTypes,
	}.Build()
	File_testpbdata_proto = out.File
	file_testpbdata_proto_rawDesc = nil
	file_testpbdata_proto_goTypes = nil
	file_testpbdata_proto_depIdxs = nil
}
