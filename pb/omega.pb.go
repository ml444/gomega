// Code generated by protoc-gen-go. DO NOT EDIT.
// versions:
// 	protoc-gen-go v1.28.0
// 	protoc        v3.15.1
// source: omega.proto

package pb

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

type SubReq_Policy int32

const (
	SubReq_PolicyConcurrence SubReq_Policy = 0
	SubReq_PolicySerial      SubReq_Policy = 1
)

// Enum value maps for SubReq_Policy.
var (
	SubReq_Policy_name = map[int32]string{
		0: "PolicyConcurrence",
		1: "PolicySerial",
	}
	SubReq_Policy_value = map[string]int32{
		"PolicyConcurrence": 0,
		"PolicySerial":      1,
	}
)

func (x SubReq_Policy) Enum() *SubReq_Policy {
	p := new(SubReq_Policy)
	*p = x
	return p
}

func (x SubReq_Policy) String() string {
	return protoimpl.X.EnumStringOf(x.Descriptor(), protoreflect.EnumNumber(x))
}

func (SubReq_Policy) Descriptor() protoreflect.EnumDescriptor {
	return file_omega_proto_enumTypes[0].Descriptor()
}

func (SubReq_Policy) Type() protoreflect.EnumType {
	return &file_omega_proto_enumTypes[0]
}

func (x SubReq_Policy) Number() protoreflect.EnumNumber {
	return protoreflect.EnumNumber(x)
}

// Deprecated: Use SubReq_Policy.Descriptor instead.
func (SubReq_Policy) EnumDescriptor() ([]byte, []int) {
	return file_omega_proto_rawDescGZIP(), []int{2, 0}
}

type PubReq struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Namespace string `protobuf:"bytes,1,opt,name=namespace,proto3" json:"namespace,omitempty"`
	TopicName string `protobuf:"bytes,2,opt,name=topic_name,json=topicName,proto3" json:"topic_name,omitempty"`
	Partition uint32 `protobuf:"varint,3,opt,name=partition,proto3" json:"partition,omitempty"`
	HashCode  uint64 `protobuf:"varint,4,opt,name=hash_code,json=hashCode,proto3" json:"hash_code,omitempty"`
	Data      []byte `protobuf:"bytes,5,opt,name=data,proto3" json:"data,omitempty"`
}

func (x *PubReq) Reset() {
	*x = PubReq{}
	if protoimpl.UnsafeEnabled {
		mi := &file_omega_proto_msgTypes[0]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *PubReq) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*PubReq) ProtoMessage() {}

func (x *PubReq) ProtoReflect() protoreflect.Message {
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

// Deprecated: Use PubReq.ProtoReflect.Descriptor instead.
func (*PubReq) Descriptor() ([]byte, []int) {
	return file_omega_proto_rawDescGZIP(), []int{0}
}

func (x *PubReq) GetNamespace() string {
	if x != nil {
		return x.Namespace
	}
	return ""
}

func (x *PubReq) GetTopicName() string {
	if x != nil {
		return x.TopicName
	}
	return ""
}

func (x *PubReq) GetPartition() uint32 {
	if x != nil {
		return x.Partition
	}
	return 0
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

type PubRsp struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Status  int32  `protobuf:"zigzag32,1,opt,name=status,proto3" json:"status,omitempty"`
	Message string `protobuf:"bytes,2,opt,name=message,proto3" json:"message,omitempty"`
}

func (x *PubRsp) Reset() {
	*x = PubRsp{}
	if protoimpl.UnsafeEnabled {
		mi := &file_omega_proto_msgTypes[1]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *PubRsp) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*PubRsp) ProtoMessage() {}

func (x *PubRsp) ProtoReflect() protoreflect.Message {
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

// Deprecated: Use PubRsp.ProtoReflect.Descriptor instead.
func (*PubRsp) Descriptor() ([]byte, []int) {
	return file_omega_proto_rawDescGZIP(), []int{1}
}

func (x *PubRsp) GetStatus() int32 {
	if x != nil {
		return x.Status
	}
	return 0
}

func (x *PubRsp) GetMessage() string {
	if x != nil {
		return x.Message
	}
	return ""
}

// strategy
type SubReq struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	ClientId  string `protobuf:"bytes,1,opt,name=client_id,json=clientId,proto3" json:"client_id,omitempty"`
	Namespace string `protobuf:"bytes,2,opt,name=namespace,proto3" json:"namespace,omitempty"`
	Topic     string `protobuf:"bytes,3,opt,name=topic,proto3" json:"topic,omitempty"`
	Partition uint32 `protobuf:"varint,4,opt,name=partition,proto3" json:"partition,omitempty"`
	Sequence  uint64 `protobuf:"varint,5,opt,name=sequence,proto3" json:"sequence,omitempty"`
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

func (x *SubReq) GetNamespace() string {
	if x != nil {
		return x.Namespace
	}
	return ""
}

func (x *SubReq) GetTopic() string {
	if x != nil {
		return x.Topic
	}
	return ""
}

func (x *SubReq) GetPartition() uint32 {
	if x != nil {
		return x.Partition
	}
	return 0
}

func (x *SubReq) GetSequence() uint64 {
	if x != nil {
		return x.Sequence
	}
	return 0
}

type SubRsp struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Status int32  `protobuf:"zigzag32,1,opt,name=status,proto3" json:"status,omitempty"`
	Token  string `protobuf:"bytes,2,opt,name=token,proto3" json:"token,omitempty"`
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

type ConsumeReq struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Token            string `protobuf:"bytes,1,opt,name=token,proto3" json:"token,omitempty"`
	Sequence         uint64 `protobuf:"varint,2,opt,name=sequence,proto3" json:"sequence,omitempty"`
	IsRetry          bool   `protobuf:"varint,3,opt,name=is_retry,json=isRetry,proto3" json:"is_retry,omitempty"`
	IgnoreRetryCount bool   `protobuf:"varint,4,opt,name=ignore_retry_count,json=ignoreRetryCount,proto3" json:"ignore_retry_count,omitempty"`
	RetryIntervalMs  int64  `protobuf:"zigzag64,5,opt,name=retry_interval_ms,json=retryIntervalMs,proto3" json:"retry_interval_ms,omitempty"`
}

func (x *ConsumeReq) Reset() {
	*x = ConsumeReq{}
	if protoimpl.UnsafeEnabled {
		mi := &file_omega_proto_msgTypes[4]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *ConsumeReq) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*ConsumeReq) ProtoMessage() {}

func (x *ConsumeReq) ProtoReflect() protoreflect.Message {
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

// Deprecated: Use ConsumeReq.ProtoReflect.Descriptor instead.
func (*ConsumeReq) Descriptor() ([]byte, []int) {
	return file_omega_proto_rawDescGZIP(), []int{4}
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

func (x *ConsumeReq) GetIsRetry() bool {
	if x != nil {
		return x.IsRetry
	}
	return false
}

func (x *ConsumeReq) GetIgnoreRetryCount() bool {
	if x != nil {
		return x.IgnoreRetryCount
	}
	return false
}

func (x *ConsumeReq) GetRetryIntervalMs() int64 {
	if x != nil {
		return x.RetryIntervalMs
	}
	return 0
}

type ConsumeRsp struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Partition  uint32 `protobuf:"varint,2,opt,name=partition,proto3" json:"partition,omitempty"`
	Sequence   uint64 `protobuf:"varint,3,opt,name=sequence,proto3" json:"sequence,omitempty"`
	Data       []byte `protobuf:"bytes,4,opt,name=data,proto3" json:"data,omitempty"`
	RetryCount uint32 `protobuf:"varint,5,opt,name=retry_count,json=retryCount,proto3" json:"retry_count,omitempty"`
}

func (x *ConsumeRsp) Reset() {
	*x = ConsumeRsp{}
	if protoimpl.UnsafeEnabled {
		mi := &file_omega_proto_msgTypes[5]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *ConsumeRsp) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*ConsumeRsp) ProtoMessage() {}

func (x *ConsumeRsp) ProtoReflect() protoreflect.Message {
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

// Deprecated: Use ConsumeRsp.ProtoReflect.Descriptor instead.
func (*ConsumeRsp) Descriptor() ([]byte, []int) {
	return file_omega_proto_rawDescGZIP(), []int{5}
}

func (x *ConsumeRsp) GetPartition() uint32 {
	if x != nil {
		return x.Partition
	}
	return 0
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

func (x *ConsumeRsp) GetRetryCount() uint32 {
	if x != nil {
		return x.RetryCount
	}
	return 0
}

var File_omega_proto protoreflect.FileDescriptor

var file_omega_proto_rawDesc = []byte{
	0x0a, 0x0b, 0x6f, 0x6d, 0x65, 0x67, 0x61, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x12, 0x05, 0x6f,
	0x6d, 0x65, 0x67, 0x61, 0x22, 0x94, 0x01, 0x0a, 0x06, 0x50, 0x75, 0x62, 0x52, 0x65, 0x71, 0x12,
	0x1c, 0x0a, 0x09, 0x6e, 0x61, 0x6d, 0x65, 0x73, 0x70, 0x61, 0x63, 0x65, 0x18, 0x01, 0x20, 0x01,
	0x28, 0x09, 0x52, 0x09, 0x6e, 0x61, 0x6d, 0x65, 0x73, 0x70, 0x61, 0x63, 0x65, 0x12, 0x1d, 0x0a,
	0x0a, 0x74, 0x6f, 0x70, 0x69, 0x63, 0x5f, 0x6e, 0x61, 0x6d, 0x65, 0x18, 0x02, 0x20, 0x01, 0x28,
	0x09, 0x52, 0x09, 0x74, 0x6f, 0x70, 0x69, 0x63, 0x4e, 0x61, 0x6d, 0x65, 0x12, 0x1c, 0x0a, 0x09,
	0x70, 0x61, 0x72, 0x74, 0x69, 0x74, 0x69, 0x6f, 0x6e, 0x18, 0x03, 0x20, 0x01, 0x28, 0x0d, 0x52,
	0x09, 0x70, 0x61, 0x72, 0x74, 0x69, 0x74, 0x69, 0x6f, 0x6e, 0x12, 0x1b, 0x0a, 0x09, 0x68, 0x61,
	0x73, 0x68, 0x5f, 0x63, 0x6f, 0x64, 0x65, 0x18, 0x04, 0x20, 0x01, 0x28, 0x04, 0x52, 0x08, 0x68,
	0x61, 0x73, 0x68, 0x43, 0x6f, 0x64, 0x65, 0x12, 0x12, 0x0a, 0x04, 0x64, 0x61, 0x74, 0x61, 0x18,
	0x05, 0x20, 0x01, 0x28, 0x0c, 0x52, 0x04, 0x64, 0x61, 0x74, 0x61, 0x22, 0x3a, 0x0a, 0x06, 0x50,
	0x75, 0x62, 0x52, 0x73, 0x70, 0x12, 0x16, 0x0a, 0x06, 0x73, 0x74, 0x61, 0x74, 0x75, 0x73, 0x18,
	0x01, 0x20, 0x01, 0x28, 0x11, 0x52, 0x06, 0x73, 0x74, 0x61, 0x74, 0x75, 0x73, 0x12, 0x18, 0x0a,
	0x07, 0x6d, 0x65, 0x73, 0x73, 0x61, 0x67, 0x65, 0x18, 0x02, 0x20, 0x01, 0x28, 0x09, 0x52, 0x07,
	0x6d, 0x65, 0x73, 0x73, 0x61, 0x67, 0x65, 0x22, 0xc6, 0x01, 0x0a, 0x06, 0x53, 0x75, 0x62, 0x52,
	0x65, 0x71, 0x12, 0x1b, 0x0a, 0x09, 0x63, 0x6c, 0x69, 0x65, 0x6e, 0x74, 0x5f, 0x69, 0x64, 0x18,
	0x01, 0x20, 0x01, 0x28, 0x09, 0x52, 0x08, 0x63, 0x6c, 0x69, 0x65, 0x6e, 0x74, 0x49, 0x64, 0x12,
	0x1c, 0x0a, 0x09, 0x6e, 0x61, 0x6d, 0x65, 0x73, 0x70, 0x61, 0x63, 0x65, 0x18, 0x02, 0x20, 0x01,
	0x28, 0x09, 0x52, 0x09, 0x6e, 0x61, 0x6d, 0x65, 0x73, 0x70, 0x61, 0x63, 0x65, 0x12, 0x14, 0x0a,
	0x05, 0x74, 0x6f, 0x70, 0x69, 0x63, 0x18, 0x03, 0x20, 0x01, 0x28, 0x09, 0x52, 0x05, 0x74, 0x6f,
	0x70, 0x69, 0x63, 0x12, 0x1c, 0x0a, 0x09, 0x70, 0x61, 0x72, 0x74, 0x69, 0x74, 0x69, 0x6f, 0x6e,
	0x18, 0x04, 0x20, 0x01, 0x28, 0x0d, 0x52, 0x09, 0x70, 0x61, 0x72, 0x74, 0x69, 0x74, 0x69, 0x6f,
	0x6e, 0x12, 0x1a, 0x0a, 0x08, 0x73, 0x65, 0x71, 0x75, 0x65, 0x6e, 0x63, 0x65, 0x18, 0x05, 0x20,
	0x01, 0x28, 0x04, 0x52, 0x08, 0x73, 0x65, 0x71, 0x75, 0x65, 0x6e, 0x63, 0x65, 0x22, 0x31, 0x0a,
	0x06, 0x50, 0x6f, 0x6c, 0x69, 0x63, 0x79, 0x12, 0x15, 0x0a, 0x11, 0x50, 0x6f, 0x6c, 0x69, 0x63,
	0x79, 0x43, 0x6f, 0x6e, 0x63, 0x75, 0x72, 0x72, 0x65, 0x6e, 0x63, 0x65, 0x10, 0x00, 0x12, 0x10,
	0x0a, 0x0c, 0x50, 0x6f, 0x6c, 0x69, 0x63, 0x79, 0x53, 0x65, 0x72, 0x69, 0x61, 0x6c, 0x10, 0x01,
	0x22, 0x36, 0x0a, 0x06, 0x53, 0x75, 0x62, 0x52, 0x73, 0x70, 0x12, 0x16, 0x0a, 0x06, 0x73, 0x74,
	0x61, 0x74, 0x75, 0x73, 0x18, 0x01, 0x20, 0x01, 0x28, 0x11, 0x52, 0x06, 0x73, 0x74, 0x61, 0x74,
	0x75, 0x73, 0x12, 0x14, 0x0a, 0x05, 0x74, 0x6f, 0x6b, 0x65, 0x6e, 0x18, 0x02, 0x20, 0x01, 0x28,
	0x09, 0x52, 0x05, 0x74, 0x6f, 0x6b, 0x65, 0x6e, 0x22, 0xb3, 0x01, 0x0a, 0x0a, 0x43, 0x6f, 0x6e,
	0x73, 0x75, 0x6d, 0x65, 0x52, 0x65, 0x71, 0x12, 0x14, 0x0a, 0x05, 0x74, 0x6f, 0x6b, 0x65, 0x6e,
	0x18, 0x01, 0x20, 0x01, 0x28, 0x09, 0x52, 0x05, 0x74, 0x6f, 0x6b, 0x65, 0x6e, 0x12, 0x1a, 0x0a,
	0x08, 0x73, 0x65, 0x71, 0x75, 0x65, 0x6e, 0x63, 0x65, 0x18, 0x02, 0x20, 0x01, 0x28, 0x04, 0x52,
	0x08, 0x73, 0x65, 0x71, 0x75, 0x65, 0x6e, 0x63, 0x65, 0x12, 0x19, 0x0a, 0x08, 0x69, 0x73, 0x5f,
	0x72, 0x65, 0x74, 0x72, 0x79, 0x18, 0x03, 0x20, 0x01, 0x28, 0x08, 0x52, 0x07, 0x69, 0x73, 0x52,
	0x65, 0x74, 0x72, 0x79, 0x12, 0x2c, 0x0a, 0x12, 0x69, 0x67, 0x6e, 0x6f, 0x72, 0x65, 0x5f, 0x72,
	0x65, 0x74, 0x72, 0x79, 0x5f, 0x63, 0x6f, 0x75, 0x6e, 0x74, 0x18, 0x04, 0x20, 0x01, 0x28, 0x08,
	0x52, 0x10, 0x69, 0x67, 0x6e, 0x6f, 0x72, 0x65, 0x52, 0x65, 0x74, 0x72, 0x79, 0x43, 0x6f, 0x75,
	0x6e, 0x74, 0x12, 0x2a, 0x0a, 0x11, 0x72, 0x65, 0x74, 0x72, 0x79, 0x5f, 0x69, 0x6e, 0x74, 0x65,
	0x72, 0x76, 0x61, 0x6c, 0x5f, 0x6d, 0x73, 0x18, 0x05, 0x20, 0x01, 0x28, 0x12, 0x52, 0x0f, 0x72,
	0x65, 0x74, 0x72, 0x79, 0x49, 0x6e, 0x74, 0x65, 0x72, 0x76, 0x61, 0x6c, 0x4d, 0x73, 0x22, 0x7b,
	0x0a, 0x0a, 0x43, 0x6f, 0x6e, 0x73, 0x75, 0x6d, 0x65, 0x52, 0x73, 0x70, 0x12, 0x1c, 0x0a, 0x09,
	0x70, 0x61, 0x72, 0x74, 0x69, 0x74, 0x69, 0x6f, 0x6e, 0x18, 0x02, 0x20, 0x01, 0x28, 0x0d, 0x52,
	0x09, 0x70, 0x61, 0x72, 0x74, 0x69, 0x74, 0x69, 0x6f, 0x6e, 0x12, 0x1a, 0x0a, 0x08, 0x73, 0x65,
	0x71, 0x75, 0x65, 0x6e, 0x63, 0x65, 0x18, 0x03, 0x20, 0x01, 0x28, 0x04, 0x52, 0x08, 0x73, 0x65,
	0x71, 0x75, 0x65, 0x6e, 0x63, 0x65, 0x12, 0x12, 0x0a, 0x04, 0x64, 0x61, 0x74, 0x61, 0x18, 0x04,
	0x20, 0x01, 0x28, 0x0c, 0x52, 0x04, 0x64, 0x61, 0x74, 0x61, 0x12, 0x1f, 0x0a, 0x0b, 0x72, 0x65,
	0x74, 0x72, 0x79, 0x5f, 0x63, 0x6f, 0x75, 0x6e, 0x74, 0x18, 0x05, 0x20, 0x01, 0x28, 0x0d, 0x52,
	0x0a, 0x72, 0x65, 0x74, 0x72, 0x79, 0x43, 0x6f, 0x75, 0x6e, 0x74, 0x32, 0x93, 0x01, 0x0a, 0x0c,
	0x4f, 0x6d, 0x65, 0x67, 0x61, 0x53, 0x65, 0x72, 0x76, 0x69, 0x63, 0x65, 0x12, 0x25, 0x0a, 0x03,
	0x50, 0x75, 0x62, 0x12, 0x0d, 0x2e, 0x6f, 0x6d, 0x65, 0x67, 0x61, 0x2e, 0x50, 0x75, 0x62, 0x52,
	0x65, 0x71, 0x1a, 0x0d, 0x2e, 0x6f, 0x6d, 0x65, 0x67, 0x61, 0x2e, 0x50, 0x75, 0x62, 0x52, 0x73,
	0x70, 0x22, 0x00, 0x12, 0x25, 0x0a, 0x03, 0x53, 0x75, 0x62, 0x12, 0x0d, 0x2e, 0x6f, 0x6d, 0x65,
	0x67, 0x61, 0x2e, 0x53, 0x75, 0x62, 0x52, 0x65, 0x71, 0x1a, 0x0d, 0x2e, 0x6f, 0x6d, 0x65, 0x67,
	0x61, 0x2e, 0x53, 0x75, 0x62, 0x52, 0x73, 0x70, 0x22, 0x00, 0x12, 0x35, 0x0a, 0x07, 0x43, 0x6f,
	0x6e, 0x73, 0x75, 0x6d, 0x65, 0x12, 0x11, 0x2e, 0x6f, 0x6d, 0x65, 0x67, 0x61, 0x2e, 0x43, 0x6f,
	0x6e, 0x73, 0x75, 0x6d, 0x65, 0x52, 0x65, 0x71, 0x1a, 0x11, 0x2e, 0x6f, 0x6d, 0x65, 0x67, 0x61,
	0x2e, 0x43, 0x6f, 0x6e, 0x73, 0x75, 0x6d, 0x65, 0x52, 0x73, 0x70, 0x22, 0x00, 0x28, 0x01, 0x30,
	0x01, 0x42, 0x22, 0x5a, 0x20, 0x67, 0x69, 0x74, 0x68, 0x75, 0x62, 0x2e, 0x63, 0x6f, 0x6d, 0x2f,
	0x6d, 0x6c, 0x34, 0x34, 0x34, 0x2f, 0x73, 0x63, 0x68, 0x65, 0x64, 0x75, 0x6c, 0x65, 0x72, 0x2f,
	0x70, 0x62, 0x2f, 0x70, 0x62, 0x62, 0x06, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x33,
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
var file_omega_proto_msgTypes = make([]protoimpl.MessageInfo, 6)
var file_omega_proto_goTypes = []interface{}{
	(SubReq_Policy)(0), // 0: omega.SubReq.Policy
	(*PubReq)(nil),     // 1: omega.PubReq
	(*PubRsp)(nil),     // 2: omega.PubRsp
	(*SubReq)(nil),     // 3: omega.SubReq
	(*SubRsp)(nil),     // 4: omega.SubRsp
	(*ConsumeReq)(nil), // 5: omega.ConsumeReq
	(*ConsumeRsp)(nil), // 6: omega.ConsumeRsp
}
var file_omega_proto_depIdxs = []int32{
	1, // 0: omega.OmegaService.Pub:input_type -> omega.PubReq
	3, // 1: omega.OmegaService.Sub:input_type -> omega.SubReq
	5, // 2: omega.OmegaService.Consume:input_type -> omega.ConsumeReq
	2, // 3: omega.OmegaService.Pub:output_type -> omega.PubRsp
	4, // 4: omega.OmegaService.Sub:output_type -> omega.SubRsp
	6, // 5: omega.OmegaService.Consume:output_type -> omega.ConsumeRsp
	3, // [3:6] is the sub-list for method output_type
	0, // [0:3] is the sub-list for method input_type
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
		file_omega_proto_msgTypes[1].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*PubRsp); i {
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
		file_omega_proto_msgTypes[5].Exporter = func(v interface{}, i int) interface{} {
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
			NumMessages:   6,
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
