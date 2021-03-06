// Code generated by protoc-gen-go. DO NOT EDIT.
// versions:
// 	protoc-gen-go v1.25.0-devel
// 	protoc        v3.6.1
// source: flow_manager.proto

package flow

import (
	proto "github.com/golang/protobuf/proto"
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

// This is a compile-time assertion that a sufficiently up-to-date version
// of the legacy proto package is being used.
const _ = proto.ProtoPackageIsVersion4

type BreakBridgeRequest struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	ConversationId string `protobuf:"bytes,1,opt,name=conversation_id,json=conversationId,proto3" json:"conversation_id,omitempty"`
	Cause          string `protobuf:"bytes,2,opt,name=cause,proto3" json:"cause,omitempty"`
}

func (x *BreakBridgeRequest) Reset() {
	*x = BreakBridgeRequest{}
	if protoimpl.UnsafeEnabled {
		mi := &file_flow_manager_proto_msgTypes[0]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *BreakBridgeRequest) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*BreakBridgeRequest) ProtoMessage() {}

func (x *BreakBridgeRequest) ProtoReflect() protoreflect.Message {
	mi := &file_flow_manager_proto_msgTypes[0]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use BreakBridgeRequest.ProtoReflect.Descriptor instead.
func (*BreakBridgeRequest) Descriptor() ([]byte, []int) {
	return file_flow_manager_proto_rawDescGZIP(), []int{0}
}

func (x *BreakBridgeRequest) GetConversationId() string {
	if x != nil {
		return x.ConversationId
	}
	return ""
}

func (x *BreakBridgeRequest) GetCause() string {
	if x != nil {
		return x.Cause
	}
	return ""
}

type BreakBridgeResponse struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Error *Error `protobuf:"bytes,1,opt,name=error,proto3" json:"error,omitempty"`
}

func (x *BreakBridgeResponse) Reset() {
	*x = BreakBridgeResponse{}
	if protoimpl.UnsafeEnabled {
		mi := &file_flow_manager_proto_msgTypes[1]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *BreakBridgeResponse) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*BreakBridgeResponse) ProtoMessage() {}

func (x *BreakBridgeResponse) ProtoReflect() protoreflect.Message {
	mi := &file_flow_manager_proto_msgTypes[1]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use BreakBridgeResponse.ProtoReflect.Descriptor instead.
func (*BreakBridgeResponse) Descriptor() ([]byte, []int) {
	return file_flow_manager_proto_rawDescGZIP(), []int{1}
}

func (x *BreakBridgeResponse) GetError() *Error {
	if x != nil {
		return x.Error
	}
	return nil
}

type Message struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Id   int64  `protobuf:"varint,1,opt,name=id,proto3" json:"id,omitempty"`
	Type string `protobuf:"bytes,2,opt,name=type,proto3" json:"type,omitempty"`
	// Types that are assignable to Value:
	//	*Message_Text
	//	*Message_File_
	Value isMessage_Value `protobuf_oneof:"value"`
}

func (x *Message) Reset() {
	*x = Message{}
	if protoimpl.UnsafeEnabled {
		mi := &file_flow_manager_proto_msgTypes[2]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *Message) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*Message) ProtoMessage() {}

func (x *Message) ProtoReflect() protoreflect.Message {
	mi := &file_flow_manager_proto_msgTypes[2]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use Message.ProtoReflect.Descriptor instead.
func (*Message) Descriptor() ([]byte, []int) {
	return file_flow_manager_proto_rawDescGZIP(), []int{2}
}

func (x *Message) GetId() int64 {
	if x != nil {
		return x.Id
	}
	return 0
}

func (x *Message) GetType() string {
	if x != nil {
		return x.Type
	}
	return ""
}

func (m *Message) GetValue() isMessage_Value {
	if m != nil {
		return m.Value
	}
	return nil
}

func (x *Message) GetText() string {
	if x, ok := x.GetValue().(*Message_Text); ok {
		return x.Text
	}
	return ""
}

func (x *Message) GetFile() *Message_File {
	if x, ok := x.GetValue().(*Message_File_); ok {
		return x.File
	}
	return nil
}

type isMessage_Value interface {
	isMessage_Value()
}

type Message_Text struct {
	Text string `protobuf:"bytes,3,opt,name=text,proto3,oneof"`
}

type Message_File_ struct {
	File *Message_File `protobuf:"bytes,4,opt,name=file,proto3,oneof"`
}

func (*Message_Text) isMessage_Value() {}

func (*Message_File_) isMessage_Value() {}

type StartRequest struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	ConversationId string            `protobuf:"bytes,1,opt,name=conversation_id,json=conversationId,proto3" json:"conversation_id,omitempty"`
	ProfileId      int64             `protobuf:"varint,2,opt,name=profile_id,json=profileId,proto3" json:"profile_id,omitempty"`
	DomainId       int64             `protobuf:"varint,3,opt,name=domain_id,json=domainId,proto3" json:"domain_id,omitempty"`
	Message        *Message          `protobuf:"bytes,4,opt,name=message,proto3" json:"message,omitempty"`
	Variables      map[string]string `protobuf:"bytes,5,rep,name=variables,proto3" json:"variables,omitempty" protobuf_key:"bytes,1,opt,name=key,proto3" protobuf_val:"bytes,2,opt,name=value,proto3"`
}

func (x *StartRequest) Reset() {
	*x = StartRequest{}
	if protoimpl.UnsafeEnabled {
		mi := &file_flow_manager_proto_msgTypes[3]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *StartRequest) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*StartRequest) ProtoMessage() {}

func (x *StartRequest) ProtoReflect() protoreflect.Message {
	mi := &file_flow_manager_proto_msgTypes[3]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use StartRequest.ProtoReflect.Descriptor instead.
func (*StartRequest) Descriptor() ([]byte, []int) {
	return file_flow_manager_proto_rawDescGZIP(), []int{3}
}

func (x *StartRequest) GetConversationId() string {
	if x != nil {
		return x.ConversationId
	}
	return ""
}

func (x *StartRequest) GetProfileId() int64 {
	if x != nil {
		return x.ProfileId
	}
	return 0
}

func (x *StartRequest) GetDomainId() int64 {
	if x != nil {
		return x.DomainId
	}
	return 0
}

func (x *StartRequest) GetMessage() *Message {
	if x != nil {
		return x.Message
	}
	return nil
}

func (x *StartRequest) GetVariables() map[string]string {
	if x != nil {
		return x.Variables
	}
	return nil
}

type Error struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Id      string `protobuf:"bytes,1,opt,name=id,proto3" json:"id,omitempty"`
	Message string `protobuf:"bytes,2,opt,name=message,proto3" json:"message,omitempty"`
}

func (x *Error) Reset() {
	*x = Error{}
	if protoimpl.UnsafeEnabled {
		mi := &file_flow_manager_proto_msgTypes[4]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *Error) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*Error) ProtoMessage() {}

func (x *Error) ProtoReflect() protoreflect.Message {
	mi := &file_flow_manager_proto_msgTypes[4]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use Error.ProtoReflect.Descriptor instead.
func (*Error) Descriptor() ([]byte, []int) {
	return file_flow_manager_proto_rawDescGZIP(), []int{4}
}

func (x *Error) GetId() string {
	if x != nil {
		return x.Id
	}
	return ""
}

func (x *Error) GetMessage() string {
	if x != nil {
		return x.Message
	}
	return ""
}

type StartResponse struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Error *Error `protobuf:"bytes,1,opt,name=error,proto3" json:"error,omitempty"`
}

func (x *StartResponse) Reset() {
	*x = StartResponse{}
	if protoimpl.UnsafeEnabled {
		mi := &file_flow_manager_proto_msgTypes[5]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *StartResponse) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*StartResponse) ProtoMessage() {}

func (x *StartResponse) ProtoReflect() protoreflect.Message {
	mi := &file_flow_manager_proto_msgTypes[5]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use StartResponse.ProtoReflect.Descriptor instead.
func (*StartResponse) Descriptor() ([]byte, []int) {
	return file_flow_manager_proto_rawDescGZIP(), []int{5}
}

func (x *StartResponse) GetError() *Error {
	if x != nil {
		return x.Error
	}
	return nil
}

type BreakRequest struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	ConversationId string `protobuf:"bytes,1,opt,name=conversation_id,json=conversationId,proto3" json:"conversation_id,omitempty"`
}

func (x *BreakRequest) Reset() {
	*x = BreakRequest{}
	if protoimpl.UnsafeEnabled {
		mi := &file_flow_manager_proto_msgTypes[6]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *BreakRequest) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*BreakRequest) ProtoMessage() {}

func (x *BreakRequest) ProtoReflect() protoreflect.Message {
	mi := &file_flow_manager_proto_msgTypes[6]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use BreakRequest.ProtoReflect.Descriptor instead.
func (*BreakRequest) Descriptor() ([]byte, []int) {
	return file_flow_manager_proto_rawDescGZIP(), []int{6}
}

func (x *BreakRequest) GetConversationId() string {
	if x != nil {
		return x.ConversationId
	}
	return ""
}

type BreakResponse struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Error *Error `protobuf:"bytes,1,opt,name=error,proto3" json:"error,omitempty"`
}

func (x *BreakResponse) Reset() {
	*x = BreakResponse{}
	if protoimpl.UnsafeEnabled {
		mi := &file_flow_manager_proto_msgTypes[7]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *BreakResponse) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*BreakResponse) ProtoMessage() {}

func (x *BreakResponse) ProtoReflect() protoreflect.Message {
	mi := &file_flow_manager_proto_msgTypes[7]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use BreakResponse.ProtoReflect.Descriptor instead.
func (*BreakResponse) Descriptor() ([]byte, []int) {
	return file_flow_manager_proto_rawDescGZIP(), []int{7}
}

func (x *BreakResponse) GetError() *Error {
	if x != nil {
		return x.Error
	}
	return nil
}

type ConfirmationMessageRequest struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	ConversationId string     `protobuf:"bytes,1,opt,name=conversation_id,json=conversationId,proto3" json:"conversation_id,omitempty"`
	ConfirmationId string     `protobuf:"bytes,2,opt,name=confirmation_id,json=confirmationId,proto3" json:"confirmation_id,omitempty"`
	Messages       []*Message `protobuf:"bytes,3,rep,name=messages,proto3" json:"messages,omitempty"`
}

func (x *ConfirmationMessageRequest) Reset() {
	*x = ConfirmationMessageRequest{}
	if protoimpl.UnsafeEnabled {
		mi := &file_flow_manager_proto_msgTypes[8]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *ConfirmationMessageRequest) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*ConfirmationMessageRequest) ProtoMessage() {}

func (x *ConfirmationMessageRequest) ProtoReflect() protoreflect.Message {
	mi := &file_flow_manager_proto_msgTypes[8]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use ConfirmationMessageRequest.ProtoReflect.Descriptor instead.
func (*ConfirmationMessageRequest) Descriptor() ([]byte, []int) {
	return file_flow_manager_proto_rawDescGZIP(), []int{8}
}

func (x *ConfirmationMessageRequest) GetConversationId() string {
	if x != nil {
		return x.ConversationId
	}
	return ""
}

func (x *ConfirmationMessageRequest) GetConfirmationId() string {
	if x != nil {
		return x.ConfirmationId
	}
	return ""
}

func (x *ConfirmationMessageRequest) GetMessages() []*Message {
	if x != nil {
		return x.Messages
	}
	return nil
}

type ConfirmationMessageResponse struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Error *Error `protobuf:"bytes,1,opt,name=error,proto3" json:"error,omitempty"`
}

func (x *ConfirmationMessageResponse) Reset() {
	*x = ConfirmationMessageResponse{}
	if protoimpl.UnsafeEnabled {
		mi := &file_flow_manager_proto_msgTypes[9]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *ConfirmationMessageResponse) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*ConfirmationMessageResponse) ProtoMessage() {}

func (x *ConfirmationMessageResponse) ProtoReflect() protoreflect.Message {
	mi := &file_flow_manager_proto_msgTypes[9]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use ConfirmationMessageResponse.ProtoReflect.Descriptor instead.
func (*ConfirmationMessageResponse) Descriptor() ([]byte, []int) {
	return file_flow_manager_proto_rawDescGZIP(), []int{9}
}

func (x *ConfirmationMessageResponse) GetError() *Error {
	if x != nil {
		return x.Error
	}
	return nil
}

type Message_File struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Id       int64  `protobuf:"varint,1,opt,name=id,proto3" json:"id,omitempty"`
	Url      string `protobuf:"bytes,2,opt,name=url,proto3" json:"url,omitempty"`
	MimeType string `protobuf:"bytes,3,opt,name=mime_type,json=mimeType,proto3" json:"mime_type,omitempty"`
}

func (x *Message_File) Reset() {
	*x = Message_File{}
	if protoimpl.UnsafeEnabled {
		mi := &file_flow_manager_proto_msgTypes[10]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *Message_File) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*Message_File) ProtoMessage() {}

func (x *Message_File) ProtoReflect() protoreflect.Message {
	mi := &file_flow_manager_proto_msgTypes[10]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use Message_File.ProtoReflect.Descriptor instead.
func (*Message_File) Descriptor() ([]byte, []int) {
	return file_flow_manager_proto_rawDescGZIP(), []int{2, 0}
}

func (x *Message_File) GetId() int64 {
	if x != nil {
		return x.Id
	}
	return 0
}

func (x *Message_File) GetUrl() string {
	if x != nil {
		return x.Url
	}
	return ""
}

func (x *Message_File) GetMimeType() string {
	if x != nil {
		return x.MimeType
	}
	return ""
}

var File_flow_manager_proto protoreflect.FileDescriptor

var file_flow_manager_proto_rawDesc = []byte{
	0x0a, 0x12, 0x66, 0x6c, 0x6f, 0x77, 0x5f, 0x6d, 0x61, 0x6e, 0x61, 0x67, 0x65, 0x72, 0x2e, 0x70,
	0x72, 0x6f, 0x74, 0x6f, 0x12, 0x04, 0x66, 0x6c, 0x6f, 0x77, 0x22, 0x53, 0x0a, 0x12, 0x42, 0x72,
	0x65, 0x61, 0x6b, 0x42, 0x72, 0x69, 0x64, 0x67, 0x65, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74,
	0x12, 0x27, 0x0a, 0x0f, 0x63, 0x6f, 0x6e, 0x76, 0x65, 0x72, 0x73, 0x61, 0x74, 0x69, 0x6f, 0x6e,
	0x5f, 0x69, 0x64, 0x18, 0x01, 0x20, 0x01, 0x28, 0x09, 0x52, 0x0e, 0x63, 0x6f, 0x6e, 0x76, 0x65,
	0x72, 0x73, 0x61, 0x74, 0x69, 0x6f, 0x6e, 0x49, 0x64, 0x12, 0x14, 0x0a, 0x05, 0x63, 0x61, 0x75,
	0x73, 0x65, 0x18, 0x02, 0x20, 0x01, 0x28, 0x09, 0x52, 0x05, 0x63, 0x61, 0x75, 0x73, 0x65, 0x22,
	0x38, 0x0a, 0x13, 0x42, 0x72, 0x65, 0x61, 0x6b, 0x42, 0x72, 0x69, 0x64, 0x67, 0x65, 0x52, 0x65,
	0x73, 0x70, 0x6f, 0x6e, 0x73, 0x65, 0x12, 0x21, 0x0a, 0x05, 0x65, 0x72, 0x72, 0x6f, 0x72, 0x18,
	0x01, 0x20, 0x01, 0x28, 0x0b, 0x32, 0x0b, 0x2e, 0x66, 0x6c, 0x6f, 0x77, 0x2e, 0x45, 0x72, 0x72,
	0x6f, 0x72, 0x52, 0x05, 0x65, 0x72, 0x72, 0x6f, 0x72, 0x22, 0xbd, 0x01, 0x0a, 0x07, 0x4d, 0x65,
	0x73, 0x73, 0x61, 0x67, 0x65, 0x12, 0x0e, 0x0a, 0x02, 0x69, 0x64, 0x18, 0x01, 0x20, 0x01, 0x28,
	0x03, 0x52, 0x02, 0x69, 0x64, 0x12, 0x12, 0x0a, 0x04, 0x74, 0x79, 0x70, 0x65, 0x18, 0x02, 0x20,
	0x01, 0x28, 0x09, 0x52, 0x04, 0x74, 0x79, 0x70, 0x65, 0x12, 0x14, 0x0a, 0x04, 0x74, 0x65, 0x78,
	0x74, 0x18, 0x03, 0x20, 0x01, 0x28, 0x09, 0x48, 0x00, 0x52, 0x04, 0x74, 0x65, 0x78, 0x74, 0x12,
	0x28, 0x0a, 0x04, 0x66, 0x69, 0x6c, 0x65, 0x18, 0x04, 0x20, 0x01, 0x28, 0x0b, 0x32, 0x12, 0x2e,
	0x66, 0x6c, 0x6f, 0x77, 0x2e, 0x4d, 0x65, 0x73, 0x73, 0x61, 0x67, 0x65, 0x2e, 0x46, 0x69, 0x6c,
	0x65, 0x48, 0x00, 0x52, 0x04, 0x66, 0x69, 0x6c, 0x65, 0x1a, 0x45, 0x0a, 0x04, 0x46, 0x69, 0x6c,
	0x65, 0x12, 0x0e, 0x0a, 0x02, 0x69, 0x64, 0x18, 0x01, 0x20, 0x01, 0x28, 0x03, 0x52, 0x02, 0x69,
	0x64, 0x12, 0x10, 0x0a, 0x03, 0x75, 0x72, 0x6c, 0x18, 0x02, 0x20, 0x01, 0x28, 0x09, 0x52, 0x03,
	0x75, 0x72, 0x6c, 0x12, 0x1b, 0x0a, 0x09, 0x6d, 0x69, 0x6d, 0x65, 0x5f, 0x74, 0x79, 0x70, 0x65,
	0x18, 0x03, 0x20, 0x01, 0x28, 0x09, 0x52, 0x08, 0x6d, 0x69, 0x6d, 0x65, 0x54, 0x79, 0x70, 0x65,
	0x42, 0x07, 0x0a, 0x05, 0x76, 0x61, 0x6c, 0x75, 0x65, 0x22, 0x9b, 0x02, 0x0a, 0x0c, 0x53, 0x74,
	0x61, 0x72, 0x74, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x12, 0x27, 0x0a, 0x0f, 0x63, 0x6f,
	0x6e, 0x76, 0x65, 0x72, 0x73, 0x61, 0x74, 0x69, 0x6f, 0x6e, 0x5f, 0x69, 0x64, 0x18, 0x01, 0x20,
	0x01, 0x28, 0x09, 0x52, 0x0e, 0x63, 0x6f, 0x6e, 0x76, 0x65, 0x72, 0x73, 0x61, 0x74, 0x69, 0x6f,
	0x6e, 0x49, 0x64, 0x12, 0x1d, 0x0a, 0x0a, 0x70, 0x72, 0x6f, 0x66, 0x69, 0x6c, 0x65, 0x5f, 0x69,
	0x64, 0x18, 0x02, 0x20, 0x01, 0x28, 0x03, 0x52, 0x09, 0x70, 0x72, 0x6f, 0x66, 0x69, 0x6c, 0x65,
	0x49, 0x64, 0x12, 0x1b, 0x0a, 0x09, 0x64, 0x6f, 0x6d, 0x61, 0x69, 0x6e, 0x5f, 0x69, 0x64, 0x18,
	0x03, 0x20, 0x01, 0x28, 0x03, 0x52, 0x08, 0x64, 0x6f, 0x6d, 0x61, 0x69, 0x6e, 0x49, 0x64, 0x12,
	0x27, 0x0a, 0x07, 0x6d, 0x65, 0x73, 0x73, 0x61, 0x67, 0x65, 0x18, 0x04, 0x20, 0x01, 0x28, 0x0b,
	0x32, 0x0d, 0x2e, 0x66, 0x6c, 0x6f, 0x77, 0x2e, 0x4d, 0x65, 0x73, 0x73, 0x61, 0x67, 0x65, 0x52,
	0x07, 0x6d, 0x65, 0x73, 0x73, 0x61, 0x67, 0x65, 0x12, 0x3f, 0x0a, 0x09, 0x76, 0x61, 0x72, 0x69,
	0x61, 0x62, 0x6c, 0x65, 0x73, 0x18, 0x05, 0x20, 0x03, 0x28, 0x0b, 0x32, 0x21, 0x2e, 0x66, 0x6c,
	0x6f, 0x77, 0x2e, 0x53, 0x74, 0x61, 0x72, 0x74, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x2e,
	0x56, 0x61, 0x72, 0x69, 0x61, 0x62, 0x6c, 0x65, 0x73, 0x45, 0x6e, 0x74, 0x72, 0x79, 0x52, 0x09,
	0x76, 0x61, 0x72, 0x69, 0x61, 0x62, 0x6c, 0x65, 0x73, 0x1a, 0x3c, 0x0a, 0x0e, 0x56, 0x61, 0x72,
	0x69, 0x61, 0x62, 0x6c, 0x65, 0x73, 0x45, 0x6e, 0x74, 0x72, 0x79, 0x12, 0x10, 0x0a, 0x03, 0x6b,
	0x65, 0x79, 0x18, 0x01, 0x20, 0x01, 0x28, 0x09, 0x52, 0x03, 0x6b, 0x65, 0x79, 0x12, 0x14, 0x0a,
	0x05, 0x76, 0x61, 0x6c, 0x75, 0x65, 0x18, 0x02, 0x20, 0x01, 0x28, 0x09, 0x52, 0x05, 0x76, 0x61,
	0x6c, 0x75, 0x65, 0x3a, 0x02, 0x38, 0x01, 0x22, 0x31, 0x0a, 0x05, 0x45, 0x72, 0x72, 0x6f, 0x72,
	0x12, 0x0e, 0x0a, 0x02, 0x69, 0x64, 0x18, 0x01, 0x20, 0x01, 0x28, 0x09, 0x52, 0x02, 0x69, 0x64,
	0x12, 0x18, 0x0a, 0x07, 0x6d, 0x65, 0x73, 0x73, 0x61, 0x67, 0x65, 0x18, 0x02, 0x20, 0x01, 0x28,
	0x09, 0x52, 0x07, 0x6d, 0x65, 0x73, 0x73, 0x61, 0x67, 0x65, 0x22, 0x32, 0x0a, 0x0d, 0x53, 0x74,
	0x61, 0x72, 0x74, 0x52, 0x65, 0x73, 0x70, 0x6f, 0x6e, 0x73, 0x65, 0x12, 0x21, 0x0a, 0x05, 0x65,
	0x72, 0x72, 0x6f, 0x72, 0x18, 0x01, 0x20, 0x01, 0x28, 0x0b, 0x32, 0x0b, 0x2e, 0x66, 0x6c, 0x6f,
	0x77, 0x2e, 0x45, 0x72, 0x72, 0x6f, 0x72, 0x52, 0x05, 0x65, 0x72, 0x72, 0x6f, 0x72, 0x22, 0x37,
	0x0a, 0x0c, 0x42, 0x72, 0x65, 0x61, 0x6b, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x12, 0x27,
	0x0a, 0x0f, 0x63, 0x6f, 0x6e, 0x76, 0x65, 0x72, 0x73, 0x61, 0x74, 0x69, 0x6f, 0x6e, 0x5f, 0x69,
	0x64, 0x18, 0x01, 0x20, 0x01, 0x28, 0x09, 0x52, 0x0e, 0x63, 0x6f, 0x6e, 0x76, 0x65, 0x72, 0x73,
	0x61, 0x74, 0x69, 0x6f, 0x6e, 0x49, 0x64, 0x22, 0x32, 0x0a, 0x0d, 0x42, 0x72, 0x65, 0x61, 0x6b,
	0x52, 0x65, 0x73, 0x70, 0x6f, 0x6e, 0x73, 0x65, 0x12, 0x21, 0x0a, 0x05, 0x65, 0x72, 0x72, 0x6f,
	0x72, 0x18, 0x01, 0x20, 0x01, 0x28, 0x0b, 0x32, 0x0b, 0x2e, 0x66, 0x6c, 0x6f, 0x77, 0x2e, 0x45,
	0x72, 0x72, 0x6f, 0x72, 0x52, 0x05, 0x65, 0x72, 0x72, 0x6f, 0x72, 0x22, 0x99, 0x01, 0x0a, 0x1a,
	0x43, 0x6f, 0x6e, 0x66, 0x69, 0x72, 0x6d, 0x61, 0x74, 0x69, 0x6f, 0x6e, 0x4d, 0x65, 0x73, 0x73,
	0x61, 0x67, 0x65, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x12, 0x27, 0x0a, 0x0f, 0x63, 0x6f,
	0x6e, 0x76, 0x65, 0x72, 0x73, 0x61, 0x74, 0x69, 0x6f, 0x6e, 0x5f, 0x69, 0x64, 0x18, 0x01, 0x20,
	0x01, 0x28, 0x09, 0x52, 0x0e, 0x63, 0x6f, 0x6e, 0x76, 0x65, 0x72, 0x73, 0x61, 0x74, 0x69, 0x6f,
	0x6e, 0x49, 0x64, 0x12, 0x27, 0x0a, 0x0f, 0x63, 0x6f, 0x6e, 0x66, 0x69, 0x72, 0x6d, 0x61, 0x74,
	0x69, 0x6f, 0x6e, 0x5f, 0x69, 0x64, 0x18, 0x02, 0x20, 0x01, 0x28, 0x09, 0x52, 0x0e, 0x63, 0x6f,
	0x6e, 0x66, 0x69, 0x72, 0x6d, 0x61, 0x74, 0x69, 0x6f, 0x6e, 0x49, 0x64, 0x12, 0x29, 0x0a, 0x08,
	0x6d, 0x65, 0x73, 0x73, 0x61, 0x67, 0x65, 0x73, 0x18, 0x03, 0x20, 0x03, 0x28, 0x0b, 0x32, 0x0d,
	0x2e, 0x66, 0x6c, 0x6f, 0x77, 0x2e, 0x4d, 0x65, 0x73, 0x73, 0x61, 0x67, 0x65, 0x52, 0x08, 0x6d,
	0x65, 0x73, 0x73, 0x61, 0x67, 0x65, 0x73, 0x22, 0x40, 0x0a, 0x1b, 0x43, 0x6f, 0x6e, 0x66, 0x69,
	0x72, 0x6d, 0x61, 0x74, 0x69, 0x6f, 0x6e, 0x4d, 0x65, 0x73, 0x73, 0x61, 0x67, 0x65, 0x52, 0x65,
	0x73, 0x70, 0x6f, 0x6e, 0x73, 0x65, 0x12, 0x21, 0x0a, 0x05, 0x65, 0x72, 0x72, 0x6f, 0x72, 0x18,
	0x01, 0x20, 0x01, 0x28, 0x0b, 0x32, 0x0b, 0x2e, 0x66, 0x6c, 0x6f, 0x77, 0x2e, 0x45, 0x72, 0x72,
	0x6f, 0x72, 0x52, 0x05, 0x65, 0x72, 0x72, 0x6f, 0x72, 0x32, 0xa3, 0x02, 0x0a, 0x15, 0x46, 0x6c,
	0x6f, 0x77, 0x43, 0x68, 0x61, 0x74, 0x53, 0x65, 0x72, 0x76, 0x65, 0x72, 0x53, 0x65, 0x72, 0x76,
	0x69, 0x63, 0x65, 0x12, 0x32, 0x0a, 0x05, 0x53, 0x74, 0x61, 0x72, 0x74, 0x12, 0x12, 0x2e, 0x66,
	0x6c, 0x6f, 0x77, 0x2e, 0x53, 0x74, 0x61, 0x72, 0x74, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74,
	0x1a, 0x13, 0x2e, 0x66, 0x6c, 0x6f, 0x77, 0x2e, 0x53, 0x74, 0x61, 0x72, 0x74, 0x52, 0x65, 0x73,
	0x70, 0x6f, 0x6e, 0x73, 0x65, 0x22, 0x00, 0x12, 0x32, 0x0a, 0x05, 0x42, 0x72, 0x65, 0x61, 0x6b,
	0x12, 0x12, 0x2e, 0x66, 0x6c, 0x6f, 0x77, 0x2e, 0x42, 0x72, 0x65, 0x61, 0x6b, 0x52, 0x65, 0x71,
	0x75, 0x65, 0x73, 0x74, 0x1a, 0x13, 0x2e, 0x66, 0x6c, 0x6f, 0x77, 0x2e, 0x42, 0x72, 0x65, 0x61,
	0x6b, 0x52, 0x65, 0x73, 0x70, 0x6f, 0x6e, 0x73, 0x65, 0x22, 0x00, 0x12, 0x44, 0x0a, 0x0b, 0x42,
	0x72, 0x65, 0x61, 0x6b, 0x42, 0x72, 0x69, 0x64, 0x67, 0x65, 0x12, 0x18, 0x2e, 0x66, 0x6c, 0x6f,
	0x77, 0x2e, 0x42, 0x72, 0x65, 0x61, 0x6b, 0x42, 0x72, 0x69, 0x64, 0x67, 0x65, 0x52, 0x65, 0x71,
	0x75, 0x65, 0x73, 0x74, 0x1a, 0x19, 0x2e, 0x66, 0x6c, 0x6f, 0x77, 0x2e, 0x42, 0x72, 0x65, 0x61,
	0x6b, 0x42, 0x72, 0x69, 0x64, 0x67, 0x65, 0x52, 0x65, 0x73, 0x70, 0x6f, 0x6e, 0x73, 0x65, 0x22,
	0x00, 0x12, 0x5c, 0x0a, 0x13, 0x43, 0x6f, 0x6e, 0x66, 0x69, 0x72, 0x6d, 0x61, 0x74, 0x69, 0x6f,
	0x6e, 0x4d, 0x65, 0x73, 0x73, 0x61, 0x67, 0x65, 0x12, 0x20, 0x2e, 0x66, 0x6c, 0x6f, 0x77, 0x2e,
	0x43, 0x6f, 0x6e, 0x66, 0x69, 0x72, 0x6d, 0x61, 0x74, 0x69, 0x6f, 0x6e, 0x4d, 0x65, 0x73, 0x73,
	0x61, 0x67, 0x65, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x1a, 0x21, 0x2e, 0x66, 0x6c, 0x6f,
	0x77, 0x2e, 0x43, 0x6f, 0x6e, 0x66, 0x69, 0x72, 0x6d, 0x61, 0x74, 0x69, 0x6f, 0x6e, 0x4d, 0x65,
	0x73, 0x73, 0x61, 0x67, 0x65, 0x52, 0x65, 0x73, 0x70, 0x6f, 0x6e, 0x73, 0x65, 0x22, 0x00, 0x62,
	0x06, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x33,
}

var (
	file_flow_manager_proto_rawDescOnce sync.Once
	file_flow_manager_proto_rawDescData = file_flow_manager_proto_rawDesc
)

func file_flow_manager_proto_rawDescGZIP() []byte {
	file_flow_manager_proto_rawDescOnce.Do(func() {
		file_flow_manager_proto_rawDescData = protoimpl.X.CompressGZIP(file_flow_manager_proto_rawDescData)
	})
	return file_flow_manager_proto_rawDescData
}

var file_flow_manager_proto_msgTypes = make([]protoimpl.MessageInfo, 12)
var file_flow_manager_proto_goTypes = []interface{}{
	(*BreakBridgeRequest)(nil),          // 0: flow.BreakBridgeRequest
	(*BreakBridgeResponse)(nil),         // 1: flow.BreakBridgeResponse
	(*Message)(nil),                     // 2: flow.Message
	(*StartRequest)(nil),                // 3: flow.StartRequest
	(*Error)(nil),                       // 4: flow.Error
	(*StartResponse)(nil),               // 5: flow.StartResponse
	(*BreakRequest)(nil),                // 6: flow.BreakRequest
	(*BreakResponse)(nil),               // 7: flow.BreakResponse
	(*ConfirmationMessageRequest)(nil),  // 8: flow.ConfirmationMessageRequest
	(*ConfirmationMessageResponse)(nil), // 9: flow.ConfirmationMessageResponse
	(*Message_File)(nil),                // 10: flow.Message.File
	nil,                                 // 11: flow.StartRequest.VariablesEntry
}
var file_flow_manager_proto_depIdxs = []int32{
	4,  // 0: flow.BreakBridgeResponse.error:type_name -> flow.Error
	10, // 1: flow.Message.file:type_name -> flow.Message.File
	2,  // 2: flow.StartRequest.message:type_name -> flow.Message
	11, // 3: flow.StartRequest.variables:type_name -> flow.StartRequest.VariablesEntry
	4,  // 4: flow.StartResponse.error:type_name -> flow.Error
	4,  // 5: flow.BreakResponse.error:type_name -> flow.Error
	2,  // 6: flow.ConfirmationMessageRequest.messages:type_name -> flow.Message
	4,  // 7: flow.ConfirmationMessageResponse.error:type_name -> flow.Error
	3,  // 8: flow.FlowChatServerService.Start:input_type -> flow.StartRequest
	6,  // 9: flow.FlowChatServerService.Break:input_type -> flow.BreakRequest
	0,  // 10: flow.FlowChatServerService.BreakBridge:input_type -> flow.BreakBridgeRequest
	8,  // 11: flow.FlowChatServerService.ConfirmationMessage:input_type -> flow.ConfirmationMessageRequest
	5,  // 12: flow.FlowChatServerService.Start:output_type -> flow.StartResponse
	7,  // 13: flow.FlowChatServerService.Break:output_type -> flow.BreakResponse
	1,  // 14: flow.FlowChatServerService.BreakBridge:output_type -> flow.BreakBridgeResponse
	9,  // 15: flow.FlowChatServerService.ConfirmationMessage:output_type -> flow.ConfirmationMessageResponse
	12, // [12:16] is the sub-list for method output_type
	8,  // [8:12] is the sub-list for method input_type
	8,  // [8:8] is the sub-list for extension type_name
	8,  // [8:8] is the sub-list for extension extendee
	0,  // [0:8] is the sub-list for field type_name
}

func init() { file_flow_manager_proto_init() }
func file_flow_manager_proto_init() {
	if File_flow_manager_proto != nil {
		return
	}
	if !protoimpl.UnsafeEnabled {
		file_flow_manager_proto_msgTypes[0].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*BreakBridgeRequest); i {
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
		file_flow_manager_proto_msgTypes[1].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*BreakBridgeResponse); i {
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
		file_flow_manager_proto_msgTypes[2].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*Message); i {
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
		file_flow_manager_proto_msgTypes[3].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*StartRequest); i {
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
		file_flow_manager_proto_msgTypes[4].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*Error); i {
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
		file_flow_manager_proto_msgTypes[5].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*StartResponse); i {
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
		file_flow_manager_proto_msgTypes[6].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*BreakRequest); i {
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
		file_flow_manager_proto_msgTypes[7].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*BreakResponse); i {
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
		file_flow_manager_proto_msgTypes[8].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*ConfirmationMessageRequest); i {
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
		file_flow_manager_proto_msgTypes[9].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*ConfirmationMessageResponse); i {
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
		file_flow_manager_proto_msgTypes[10].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*Message_File); i {
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
	file_flow_manager_proto_msgTypes[2].OneofWrappers = []interface{}{
		(*Message_Text)(nil),
		(*Message_File_)(nil),
	}
	type x struct{}
	out := protoimpl.TypeBuilder{
		File: protoimpl.DescBuilder{
			GoPackagePath: reflect.TypeOf(x{}).PkgPath(),
			RawDescriptor: file_flow_manager_proto_rawDesc,
			NumEnums:      0,
			NumMessages:   12,
			NumExtensions: 0,
			NumServices:   1,
		},
		GoTypes:           file_flow_manager_proto_goTypes,
		DependencyIndexes: file_flow_manager_proto_depIdxs,
		MessageInfos:      file_flow_manager_proto_msgTypes,
	}.Build()
	File_flow_manager_proto = out.File
	file_flow_manager_proto_rawDesc = nil
	file_flow_manager_proto_goTypes = nil
	file_flow_manager_proto_depIdxs = nil
}
