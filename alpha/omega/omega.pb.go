// Code generated by protoc-gen-go. DO NOT EDIT.
// versions:
// 	protoc-gen-go v1.28.1
// 	protoc        v3.21.5
// source: omega.proto

package omega

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

type DelayType int32

const (
	DelayType_DelayTypeNone         DelayType = 0
	DelayType_DelayTypeRelateTime   DelayType = 1 // delay_time is relative time
	DelayType_DelayTypeAbsoluteTime DelayType = 2 // delay_time is absolute time
)

// Enum value maps for DelayType.
var (
	DelayType_name = map[int32]string{
		0: "DelayTypeNone",
		1: "DelayTypeRelateTime",
		2: "DelayTypeAbsoluteTime",
	}
	DelayType_value = map[string]int32{
		"DelayTypeNone":         0,
		"DelayTypeRelateTime":   1,
		"DelayTypeAbsoluteTime": 2,
	}
)

func (x DelayType) Enum() *DelayType {
	p := new(DelayType)
	*p = x
	return p
}

func (x DelayType) String() string {
	return protoimpl.X.EnumStringOf(x.Descriptor(), protoreflect.EnumNumber(x))
}

func (DelayType) Descriptor() protoreflect.EnumDescriptor {
	return file_omega_proto_enumTypes[0].Descriptor()
}

func (DelayType) Type() protoreflect.EnumType {
	return &file_omega_proto_enumTypes[0]
}

func (x DelayType) Number() protoreflect.EnumNumber {
	return protoreflect.EnumNumber(x)
}

// Deprecated: Use DelayType.Descriptor instead.
func (DelayType) EnumDescriptor() ([]byte, []int) {
	return file_omega_proto_rawDescGZIP(), []int{0}
}

type Response struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Status  int32  `protobuf:"zigzag32,1,opt,name=status,proto3" json:"status,omitempty"`
	Message string `protobuf:"bytes,2,opt,name=message,proto3" json:"message,omitempty"`
}

func (x *Response) Reset() {
	*x = Response{}
	if protoimpl.UnsafeEnabled {
		mi := &file_omega_proto_msgTypes[0]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *Response) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*Response) ProtoMessage() {}

func (x *Response) ProtoReflect() protoreflect.Message {
	mi := &file_omega_proto_msgTypes[0]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use Response.ProtoReflect.Descriptor instead.
func (*Response) Descriptor() ([]byte, []int) {
	return file_omega_proto_rawDescGZIP(), []int{0}
}

func (x *Response) GetStatus() int32 {
	if x != nil {
		return x.Status
	}
	return 0
}

func (x *Response) GetMessage() string {
	if x != nil {
		return x.Message
	}
	return ""
}

type PubReq struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Topic    string `protobuf:"bytes,1,opt,name=topic,proto3" json:"topic,omitempty"`
	HashCode uint64 `protobuf:"varint,2,opt,name=hash_code,json=hashCode,proto3" json:"hash_code,omitempty"`
	Data     []byte `protobuf:"bytes,3,opt,name=data,proto3" json:"data,omitempty"`
	// @desc: priority of the message, the larger the value, the higher the priority
	Priority uint32 `protobuf:"varint,4,opt,name=priority,proto3" json:"priority,omitempty"`
	// @ref_type: DelayType
	DelayType uint32 `protobuf:"varint,5,opt,name=delay_type,json=delayType,proto3" json:"delay_type,omitempty"`
	DelayTime uint32 `protobuf:"varint,6,opt,name=delay_time,json=delayTime,proto3" json:"delay_time,omitempty"`
}

func (x *PubReq) Reset() {
	*x = PubReq{}
	if protoimpl.UnsafeEnabled {
		mi := &file_omega_proto_msgTypes[1]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *PubReq) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*PubReq) ProtoMessage() {}

func (x *PubReq) ProtoReflect() protoreflect.Message {
	mi := &file_omega_proto_msgTypes[1]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use PubReq.ProtoReflect.Descriptor instead.
func (*PubReq) Descriptor() ([]byte, []int) {
	return file_omega_proto_rawDescGZIP(), []int{1}
}

func (x *PubReq) GetTopic() string {
	if x != nil {
		return x.Topic
	}
	return ""
}

func (x *PubReq) GetHashCode() uint64 {
	if x != nil {
		return x.HashCode
	}
	return 0
}

func (x *PubReq) GetData() []byte {
	if x != nil {
		return x.Data
	}
	return nil
}

func (x *PubReq) GetPriority() uint32 {
	if x != nil {
		return x.Priority
	}
	return 0
}

func (x *PubReq) GetDelayType() uint32 {
	if x != nil {
		return x.DelayType
	}
	return 0
}

func (x *PubReq) GetDelayTime() uint32 {
	if x != nil {
		return x.DelayTime
	}
	return 0
}

type SubReq struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	ClientId string `protobuf:"bytes,1,opt,name=client_id,json=clientId,proto3" json:"client_id,omitempty"`
	Topic    string `protobuf:"bytes,3,opt,name=topic,proto3" json:"topic,omitempty"`
	Group    string `protobuf:"bytes,4,opt,name=group,proto3" json:"group,omitempty"`
}

func (x *SubReq) Reset() {
	*x = SubReq{}
	if protoimpl.UnsafeEnabled {
		mi := &file_omega_proto_msgTypes[2]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *SubReq) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*SubReq) ProtoMessage() {}

func (x *SubReq) ProtoReflect() protoreflect.Message {
	mi := &file_omega_proto_msgTypes[2]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use SubReq.ProtoReflect.Descriptor instead.
func (*SubReq) Descriptor() ([]byte, []int) {
	return file_omega_proto_rawDescGZIP(), []int{2}
}

func (x *SubReq) GetClientId() string {
	if x != nil {
		return x.ClientId
	}
	return ""
}

func (x *SubReq) GetTopic() string {
	if x != nil {
		return x.Topic
	}
	return ""
}

func (x *SubReq) GetGroup() string {
	if x != nil {
		return x.Group
	}
	return ""
}

type SubRsp struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Status   int32  `protobuf:"zigzag32,1,opt,name=status,proto3" json:"status,omitempty"`
	Token    string `protobuf:"bytes,2,opt,name=token,proto3" json:"token,omitempty"`
	Sequence uint64 `protobuf:"varint,3,opt,name=sequence,proto3" json:"sequence,omitempty"`
}

func (x *SubRsp) Reset() {
	*x = SubRsp{}
	if protoimpl.UnsafeEnabled {
		mi := &file_omega_proto_msgTypes[3]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *SubRsp) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*SubRsp) ProtoMessage() {}

func (x *SubRsp) ProtoReflect() protoreflect.Message {
	mi := &file_omega_proto_msgTypes[3]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use SubRsp.ProtoReflect.Descriptor instead.
func (*SubRsp) Descriptor() ([]byte, []int) {
	return file_omega_proto_rawDescGZIP(), []int{3}
}

func (x *SubRsp) GetStatus() int32 {
	if x != nil {
		return x.Status
	}
	return 0
}

func (x *SubRsp) GetToken() string {
	if x != nil {
		return x.Token
	}
	return ""
}

func (x *SubRsp) GetSequence() uint64 {
	if x != nil {
		return x.Sequence
	}
	return 0
}

type SubCfg struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Topic          string `protobuf:"bytes,1,opt,name=topic,proto3" json:"topic,omitempty"`
	Group          string `protobuf:"bytes,2,opt,name=group,proto3" json:"group,omitempty"`
	Version        uint32 `protobuf:"varint,3,opt,name=version,proto3" json:"version,omitempty"`
	IsHashDispatch bool   `protobuf:"varint,4,opt,name=is_hash_dispatch,json=isHashDispatch,proto3" json:"is_hash_dispatch,omitempty"`
}

func (x *SubCfg) Reset() {
	*x = SubCfg{}
	if protoimpl.UnsafeEnabled {
		mi := &file_omega_proto_msgTypes[4]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *SubCfg) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*SubCfg) ProtoMessage() {}

func (x *SubCfg) ProtoReflect() protoreflect.Message {
	mi := &file_omega_proto_msgTypes[4]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use SubCfg.ProtoReflect.Descriptor instead.
func (*SubCfg) Descriptor() ([]byte, []int) {
	return file_omega_proto_rawDescGZIP(), []int{4}
}

func (x *SubCfg) GetTopic() string {
	if x != nil {
		return x.Topic
	}
	return ""
}

func (x *SubCfg) GetGroup() string {
	if x != nil {
		return x.Group
	}
	return ""
}

func (x *SubCfg) GetVersion() uint32 {
	if x != nil {
		return x.Version
	}
	return 0
}

func (x *SubCfg) GetIsHashDispatch() bool {
	if x != nil {
		return x.IsHashDispatch
	}
	return false
}

type ConsumeReq struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Token    string `protobuf:"bytes,4,opt,name=token,proto3" json:"token,omitempty"`
	Sequence uint64 `protobuf:"varint,5,opt,name=sequence,proto3" json:"sequence,omitempty"`
}

func (x *ConsumeReq) Reset() {
	*x = ConsumeReq{}
	if protoimpl.UnsafeEnabled {
		mi := &file_omega_proto_msgTypes[5]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *ConsumeReq) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*ConsumeReq) ProtoMessage() {}

func (x *ConsumeReq) ProtoReflect() protoreflect.Message {
	mi := &file_omega_proto_msgTypes[5]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use ConsumeReq.ProtoReflect.Descriptor instead.
func (*ConsumeReq) Descriptor() ([]byte, []int) {
	return file_omega_proto_rawDescGZIP(), []int{5}
}

func (x *ConsumeReq) GetToken() string {
	if x != nil {
		return x.Token
	}
	return ""
}

func (x *ConsumeReq) GetSequence() uint64 {
	if x != nil {
		return x.Sequence
	}
	return 0
}

type ConsumeRsp struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Sequence uint64 `protobuf:"varint,3,opt,name=sequence,proto3" json:"sequence,omitempty"`
	Data     []byte `protobuf:"bytes,4,opt,name=data,proto3" json:"data,omitempty"`
}

func (x *ConsumeRsp) Reset() {
	*x = ConsumeRsp{}
	if protoimpl.UnsafeEnabled {
		mi := &file_omega_proto_msgTypes[6]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *ConsumeRsp) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*ConsumeRsp) ProtoMessage() {}

func (x *ConsumeRsp) ProtoReflect() protoreflect.Message {
	mi := &file_omega_proto_msgTypes[6]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use ConsumeRsp.ProtoReflect.Descriptor instead.
func (*ConsumeRsp) Descriptor() ([]byte, []int) {
	return file_omega_proto_rawDescGZIP(), []int{6}
}

func (x *ConsumeRsp) GetSequence() uint64 {
	if x != nil {
		return x.Sequence
	}
	return 0
}

func (x *ConsumeRsp) GetData() []byte {
	if x != nil {
		return x.Data
	}
	return nil
}

var File_omega_proto protoreflect.FileDescriptor

var file_omega_proto_rawDesc = []byte{
	0x0a, 0x0b, 0x6f, 0x6d, 0x65, 0x67, 0x61, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x12, 0x05, 0x6f,
	0x6d, 0x65, 0x67, 0x61, 0x22, 0x3c, 0x0a, 0x08, 0x52, 0x65, 0x73, 0x70, 0x6f, 0x6e, 0x73, 0x65,
	0x12, 0x16, 0x0a, 0x06, 0x73, 0x74, 0x61, 0x74, 0x75, 0x73, 0x18, 0x01, 0x20, 0x01, 0x28, 0x11,
	0x52, 0x06, 0x73, 0x74, 0x61, 0x74, 0x75, 0x73, 0x12, 0x18, 0x0a, 0x07, 0x6d, 0x65, 0x73, 0x73,
	0x61, 0x67, 0x65, 0x18, 0x02, 0x20, 0x01, 0x28, 0x09, 0x52, 0x07, 0x6d, 0x65, 0x73, 0x73, 0x61,
	0x67, 0x65, 0x22, 0xa9, 0x01, 0x0a, 0x06, 0x50, 0x75, 0x62, 0x52, 0x65, 0x71, 0x12, 0x14, 0x0a,
	0x05, 0x74, 0x6f, 0x70, 0x69, 0x63, 0x18, 0x01, 0x20, 0x01, 0x28, 0x09, 0x52, 0x05, 0x74, 0x6f,
	0x70, 0x69, 0x63, 0x12, 0x1b, 0x0a, 0x09, 0x68, 0x61, 0x73, 0x68, 0x5f, 0x63, 0x6f, 0x64, 0x65,
	0x18, 0x02, 0x20, 0x01, 0x28, 0x04, 0x52, 0x08, 0x68, 0x61, 0x73, 0x68, 0x43, 0x6f, 0x64, 0x65,
	0x12, 0x12, 0x0a, 0x04, 0x64, 0x61, 0x74, 0x61, 0x18, 0x03, 0x20, 0x01, 0x28, 0x0c, 0x52, 0x04,
	0x64, 0x61, 0x74, 0x61, 0x12, 0x1a, 0x0a, 0x08, 0x70, 0x72, 0x69, 0x6f, 0x72, 0x69, 0x74, 0x79,
	0x18, 0x04, 0x20, 0x01, 0x28, 0x0d, 0x52, 0x08, 0x70, 0x72, 0x69, 0x6f, 0x72, 0x69, 0x74, 0x79,
	0x12, 0x1d, 0x0a, 0x0a, 0x64, 0x65, 0x6c, 0x61, 0x79, 0x5f, 0x74, 0x79, 0x70, 0x65, 0x18, 0x05,
	0x20, 0x01, 0x28, 0x0d, 0x52, 0x09, 0x64, 0x65, 0x6c, 0x61, 0x79, 0x54, 0x79, 0x70, 0x65, 0x12,
	0x1d, 0x0a, 0x0a, 0x64, 0x65, 0x6c, 0x61, 0x79, 0x5f, 0x74, 0x69, 0x6d, 0x65, 0x18, 0x06, 0x20,
	0x01, 0x28, 0x0d, 0x52, 0x09, 0x64, 0x65, 0x6c, 0x61, 0x79, 0x54, 0x69, 0x6d, 0x65, 0x22, 0x51,
	0x0a, 0x06, 0x53, 0x75, 0x62, 0x52, 0x65, 0x71, 0x12, 0x1b, 0x0a, 0x09, 0x63, 0x6c, 0x69, 0x65,
	0x6e, 0x74, 0x5f, 0x69, 0x64, 0x18, 0x01, 0x20, 0x01, 0x28, 0x09, 0x52, 0x08, 0x63, 0x6c, 0x69,
	0x65, 0x6e, 0x74, 0x49, 0x64, 0x12, 0x14, 0x0a, 0x05, 0x74, 0x6f, 0x70, 0x69, 0x63, 0x18, 0x03,
	0x20, 0x01, 0x28, 0x09, 0x52, 0x05, 0x74, 0x6f, 0x70, 0x69, 0x63, 0x12, 0x14, 0x0a, 0x05, 0x67,
	0x72, 0x6f, 0x75, 0x70, 0x18, 0x04, 0x20, 0x01, 0x28, 0x09, 0x52, 0x05, 0x67, 0x72, 0x6f, 0x75,
	0x70, 0x22, 0x52, 0x0a, 0x06, 0x53, 0x75, 0x62, 0x52, 0x73, 0x70, 0x12, 0x16, 0x0a, 0x06, 0x73,
	0x74, 0x61, 0x74, 0x75, 0x73, 0x18, 0x01, 0x20, 0x01, 0x28, 0x11, 0x52, 0x06, 0x73, 0x74, 0x61,
	0x74, 0x75, 0x73, 0x12, 0x14, 0x0a, 0x05, 0x74, 0x6f, 0x6b, 0x65, 0x6e, 0x18, 0x02, 0x20, 0x01,
	0x28, 0x09, 0x52, 0x05, 0x74, 0x6f, 0x6b, 0x65, 0x6e, 0x12, 0x1a, 0x0a, 0x08, 0x73, 0x65, 0x71,
	0x75, 0x65, 0x6e, 0x63, 0x65, 0x18, 0x03, 0x20, 0x01, 0x28, 0x04, 0x52, 0x08, 0x73, 0x65, 0x71,
	0x75, 0x65, 0x6e, 0x63, 0x65, 0x22, 0x78, 0x0a, 0x06, 0x53, 0x75, 0x62, 0x43, 0x66, 0x67, 0x12,
	0x14, 0x0a, 0x05, 0x74, 0x6f, 0x70, 0x69, 0x63, 0x18, 0x01, 0x20, 0x01, 0x28, 0x09, 0x52, 0x05,
	0x74, 0x6f, 0x70, 0x69, 0x63, 0x12, 0x14, 0x0a, 0x05, 0x67, 0x72, 0x6f, 0x75, 0x70, 0x18, 0x02,
	0x20, 0x01, 0x28, 0x09, 0x52, 0x05, 0x67, 0x72, 0x6f, 0x75, 0x70, 0x12, 0x18, 0x0a, 0x07, 0x76,
	0x65, 0x72, 0x73, 0x69, 0x6f, 0x6e, 0x18, 0x03, 0x20, 0x01, 0x28, 0x0d, 0x52, 0x07, 0x76, 0x65,
	0x72, 0x73, 0x69, 0x6f, 0x6e, 0x12, 0x28, 0x0a, 0x10, 0x69, 0x73, 0x5f, 0x68, 0x61, 0x73, 0x68,
	0x5f, 0x64, 0x69, 0x73, 0x70, 0x61, 0x74, 0x63, 0x68, 0x18, 0x04, 0x20, 0x01, 0x28, 0x08, 0x52,
	0x0e, 0x69, 0x73, 0x48, 0x61, 0x73, 0x68, 0x44, 0x69, 0x73, 0x70, 0x61, 0x74, 0x63, 0x68, 0x22,
	0x3e, 0x0a, 0x0a, 0x43, 0x6f, 0x6e, 0x73, 0x75, 0x6d, 0x65, 0x52, 0x65, 0x71, 0x12, 0x14, 0x0a,
	0x05, 0x74, 0x6f, 0x6b, 0x65, 0x6e, 0x18, 0x04, 0x20, 0x01, 0x28, 0x09, 0x52, 0x05, 0x74, 0x6f,
	0x6b, 0x65, 0x6e, 0x12, 0x1a, 0x0a, 0x08, 0x73, 0x65, 0x71, 0x75, 0x65, 0x6e, 0x63, 0x65, 0x18,
	0x05, 0x20, 0x01, 0x28, 0x04, 0x52, 0x08, 0x73, 0x65, 0x71, 0x75, 0x65, 0x6e, 0x63, 0x65, 0x22,
	0x3c, 0x0a, 0x0a, 0x43, 0x6f, 0x6e, 0x73, 0x75, 0x6d, 0x65, 0x52, 0x73, 0x70, 0x12, 0x1a, 0x0a,
	0x08, 0x73, 0x65, 0x71, 0x75, 0x65, 0x6e, 0x63, 0x65, 0x18, 0x03, 0x20, 0x01, 0x28, 0x04, 0x52,
	0x08, 0x73, 0x65, 0x71, 0x75, 0x65, 0x6e, 0x63, 0x65, 0x12, 0x12, 0x0a, 0x04, 0x64, 0x61, 0x74,
	0x61, 0x18, 0x04, 0x20, 0x01, 0x28, 0x0c, 0x52, 0x04, 0x64, 0x61, 0x74, 0x61, 0x2a, 0x52, 0x0a,
	0x09, 0x44, 0x65, 0x6c, 0x61, 0x79, 0x54, 0x79, 0x70, 0x65, 0x12, 0x11, 0x0a, 0x0d, 0x44, 0x65,
	0x6c, 0x61, 0x79, 0x54, 0x79, 0x70, 0x65, 0x4e, 0x6f, 0x6e, 0x65, 0x10, 0x00, 0x12, 0x17, 0x0a,
	0x13, 0x44, 0x65, 0x6c, 0x61, 0x79, 0x54, 0x79, 0x70, 0x65, 0x52, 0x65, 0x6c, 0x61, 0x74, 0x65,
	0x54, 0x69, 0x6d, 0x65, 0x10, 0x01, 0x12, 0x19, 0x0a, 0x15, 0x44, 0x65, 0x6c, 0x61, 0x79, 0x54,
	0x79, 0x70, 0x65, 0x41, 0x62, 0x73, 0x6f, 0x6c, 0x75, 0x74, 0x65, 0x54, 0x69, 0x6d, 0x65, 0x10,
	0x02, 0x32, 0xc5, 0x01, 0x0a, 0x0c, 0x4f, 0x6d, 0x65, 0x67, 0x61, 0x53, 0x65, 0x72, 0x76, 0x69,
	0x63, 0x65, 0x12, 0x27, 0x0a, 0x03, 0x50, 0x75, 0x62, 0x12, 0x0d, 0x2e, 0x6f, 0x6d, 0x65, 0x67,
	0x61, 0x2e, 0x50, 0x75, 0x62, 0x52, 0x65, 0x71, 0x1a, 0x0f, 0x2e, 0x6f, 0x6d, 0x65, 0x67, 0x61,
	0x2e, 0x52, 0x65, 0x73, 0x70, 0x6f, 0x6e, 0x73, 0x65, 0x22, 0x00, 0x12, 0x25, 0x0a, 0x03, 0x53,
	0x75, 0x62, 0x12, 0x0d, 0x2e, 0x6f, 0x6d, 0x65, 0x67, 0x61, 0x2e, 0x53, 0x75, 0x62, 0x52, 0x65,
	0x71, 0x1a, 0x0d, 0x2e, 0x6f, 0x6d, 0x65, 0x67, 0x61, 0x2e, 0x53, 0x75, 0x62, 0x52, 0x73, 0x70,
	0x22, 0x00, 0x12, 0x30, 0x0a, 0x0c, 0x55, 0x70, 0x64, 0x61, 0x74, 0x65, 0x53, 0x75, 0x62, 0x43,
	0x66, 0x67, 0x12, 0x0d, 0x2e, 0x6f, 0x6d, 0x65, 0x67, 0x61, 0x2e, 0x53, 0x75, 0x62, 0x43, 0x66,
	0x67, 0x1a, 0x0f, 0x2e, 0x6f, 0x6d, 0x65, 0x67, 0x61, 0x2e, 0x52, 0x65, 0x73, 0x70, 0x6f, 0x6e,
	0x73, 0x65, 0x22, 0x00, 0x12, 0x33, 0x0a, 0x07, 0x43, 0x6f, 0x6e, 0x73, 0x75, 0x6d, 0x65, 0x12,
	0x11, 0x2e, 0x6f, 0x6d, 0x65, 0x67, 0x61, 0x2e, 0x43, 0x6f, 0x6e, 0x73, 0x75, 0x6d, 0x65, 0x52,
	0x65, 0x71, 0x1a, 0x11, 0x2e, 0x6f, 0x6d, 0x65, 0x67, 0x61, 0x2e, 0x43, 0x6f, 0x6e, 0x73, 0x75,
	0x6d, 0x65, 0x52, 0x73, 0x70, 0x22, 0x00, 0x30, 0x01, 0x42, 0x25, 0x5a, 0x23, 0x67, 0x69, 0x74,
	0x68, 0x75, 0x62, 0x2e, 0x63, 0x6f, 0x6d, 0x2f, 0x6d, 0x6c, 0x34, 0x34, 0x34, 0x2f, 0x67, 0x6f,
	0x6d, 0x65, 0x67, 0x61, 0x2f, 0x61, 0x6c, 0x70, 0x68, 0x61, 0x2f, 0x6f, 0x6d, 0x65, 0x67, 0x61,
	0x62, 0x06, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x33,
}

var (
	file_omega_proto_rawDescOnce sync.Once
	file_omega_proto_rawDescData = file_omega_proto_rawDesc
)

func file_omega_proto_rawDescGZIP() []byte {
	file_omega_proto_rawDescOnce.Do(func() {
		file_omega_proto_rawDescData = protoimpl.X.CompressGZIP(file_omega_proto_rawDescData)
	})
	return file_omega_proto_rawDescData
}

var file_omega_proto_enumTypes = make([]protoimpl.EnumInfo, 1)
var file_omega_proto_msgTypes = make([]protoimpl.MessageInfo, 7)
var file_omega_proto_goTypes = []interface{}{
	(DelayType)(0),     // 0: omega.DelayType
	(*Response)(nil),   // 1: omega.Response
	(*PubReq)(nil),     // 2: omega.PubReq
	(*SubReq)(nil),     // 3: omega.SubReq
	(*SubRsp)(nil),     // 4: omega.SubRsp
	(*SubCfg)(nil),     // 5: omega.SubCfg
	(*ConsumeReq)(nil), // 6: omega.ConsumeReq
	(*ConsumeRsp)(nil), // 7: omega.ConsumeRsp
}
var file_omega_proto_depIdxs = []int32{
	2, // 0: omega.OmegaService.Pub:input_type -> omega.PubReq
	3, // 1: omega.OmegaService.Sub:input_type -> omega.SubReq
	5, // 2: omega.OmegaService.UpdateSubCfg:input_type -> omega.SubCfg
	6, // 3: omega.OmegaService.Consume:input_type -> omega.ConsumeReq
	1, // 4: omega.OmegaService.Pub:output_type -> omega.Response
	4, // 5: omega.OmegaService.Sub:output_type -> omega.SubRsp
	1, // 6: omega.OmegaService.UpdateSubCfg:output_type -> omega.Response
	7, // 7: omega.OmegaService.Consume:output_type -> omega.ConsumeRsp
	4, // [4:8] is the sub-list for method output_type
	0, // [0:4] is the sub-list for method input_type
	0, // [0:0] is the sub-list for extension type_name
	0, // [0:0] is the sub-list for extension extendee
	0, // [0:0] is the sub-list for field type_name
}

func init() { file_omega_proto_init() }
func file_omega_proto_init() {
	if File_omega_proto != nil {
		return
	}
	if !protoimpl.UnsafeEnabled {
		file_omega_proto_msgTypes[0].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*Response); i {
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
		file_omega_proto_msgTypes[1].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*PubReq); i {
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
		file_omega_proto_msgTypes[2].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*SubReq); i {
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
		file_omega_proto_msgTypes[3].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*SubRsp); i {
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
		file_omega_proto_msgTypes[4].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*SubCfg); i {
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
		file_omega_proto_msgTypes[5].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*ConsumeReq); i {
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
		file_omega_proto_msgTypes[6].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*ConsumeRsp); i {
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
			RawDescriptor: file_omega_proto_rawDesc,
			NumEnums:      1,
			NumMessages:   7,
			NumExtensions: 0,
			NumServices:   1,
		},
		GoTypes:           file_omega_proto_goTypes,
		DependencyIndexes: file_omega_proto_depIdxs,
		EnumInfos:         file_omega_proto_enumTypes,
		MessageInfos:      file_omega_proto_msgTypes,
	}.Build()
	File_omega_proto = out.File
	file_omega_proto_rawDesc = nil
	file_omega_proto_goTypes = nil
	file_omega_proto_depIdxs = nil
}