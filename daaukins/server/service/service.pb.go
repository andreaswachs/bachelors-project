// Code generated by protoc-gen-go. DO NOT EDIT.
// versions:
// 	protoc-gen-go v1.29.0
// 	protoc        v3.21.12
// source: service/service.proto

package service

import (
	protoreflect "google.golang.org/protobuf/reflect/protoreflect"
	protoimpl "google.golang.org/protobuf/runtime/protoimpl"
	emptypb "google.golang.org/protobuf/types/known/emptypb"
	reflect "reflect"
	sync "sync"
)

const (
	// Verify that this generated code is sufficiently up-to-date.
	_ = protoimpl.EnforceVersion(20 - protoimpl.MinVersion)
	// Verify that runtime/protoimpl is sufficiently up-to-date.
	_ = protoimpl.EnforceVersion(protoimpl.MaxVersion - 20)
)

// HaveCapacityRequest contains the string lab
// which is YAML configuration for a lab.
// It is the server's responsibilityy to parse the YAML
// and check to see if it has capacity and provide an appropriate response.
type HaveCapacityRequest struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Lab string `protobuf:"bytes,1,opt,name=lab,proto3" json:"lab,omitempty"`
}

func (x *HaveCapacityRequest) Reset() {
	*x = HaveCapacityRequest{}
	if protoimpl.UnsafeEnabled {
		mi := &file_service_service_proto_msgTypes[0]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *HaveCapacityRequest) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*HaveCapacityRequest) ProtoMessage() {}

func (x *HaveCapacityRequest) ProtoReflect() protoreflect.Message {
	mi := &file_service_service_proto_msgTypes[0]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use HaveCapacityRequest.ProtoReflect.Descriptor instead.
func (*HaveCapacityRequest) Descriptor() ([]byte, []int) {
	return file_service_service_proto_rawDescGZIP(), []int{0}
}

func (x *HaveCapacityRequest) GetLab() string {
	if x != nil {
		return x.Lab
	}
	return ""
}

// HaveCapacityResponse contains a boolean
// which is true if the server has capacity for the lab
// and false if it does not.
type HaveCapacityResponse struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	HasCapacity bool `protobuf:"varint,1,opt,name=hasCapacity,proto3" json:"hasCapacity,omitempty"`
}

func (x *HaveCapacityResponse) Reset() {
	*x = HaveCapacityResponse{}
	if protoimpl.UnsafeEnabled {
		mi := &file_service_service_proto_msgTypes[1]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *HaveCapacityResponse) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*HaveCapacityResponse) ProtoMessage() {}

func (x *HaveCapacityResponse) ProtoReflect() protoreflect.Message {
	mi := &file_service_service_proto_msgTypes[1]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use HaveCapacityResponse.ProtoReflect.Descriptor instead.
func (*HaveCapacityResponse) Descriptor() ([]byte, []int) {
	return file_service_service_proto_rawDescGZIP(), []int{1}
}

func (x *HaveCapacityResponse) GetHasCapacity() bool {
	if x != nil {
		return x.HasCapacity
	}
	return false
}

// ScheduleLabRequest contains the string lab
// which is YAML configuration for a lab.
// It is the server's responsibility to parse the YAML
// and schedule the lab.
type ScheduleLabRequest struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Lab string `protobuf:"bytes,1,opt,name=lab,proto3" json:"lab,omitempty"`
}

func (x *ScheduleLabRequest) Reset() {
	*x = ScheduleLabRequest{}
	if protoimpl.UnsafeEnabled {
		mi := &file_service_service_proto_msgTypes[2]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *ScheduleLabRequest) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*ScheduleLabRequest) ProtoMessage() {}

func (x *ScheduleLabRequest) ProtoReflect() protoreflect.Message {
	mi := &file_service_service_proto_msgTypes[2]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use ScheduleLabRequest.ProtoReflect.Descriptor instead.
func (*ScheduleLabRequest) Descriptor() ([]byte, []int) {
	return file_service_service_proto_rawDescGZIP(), []int{2}
}

func (x *ScheduleLabRequest) GetLab() string {
	if x != nil {
		return x.Lab
	}
	return ""
}

// ScheduleLabResponse
type ScheduleLabResponse struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Scheduled bool   `protobuf:"varint,1,opt,name=scheduled,proto3" json:"scheduled,omitempty"`
	Id        string `protobuf:"bytes,2,opt,name=id,proto3" json:"id,omitempty"`
}

func (x *ScheduleLabResponse) Reset() {
	*x = ScheduleLabResponse{}
	if protoimpl.UnsafeEnabled {
		mi := &file_service_service_proto_msgTypes[3]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *ScheduleLabResponse) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*ScheduleLabResponse) ProtoMessage() {}

func (x *ScheduleLabResponse) ProtoReflect() protoreflect.Message {
	mi := &file_service_service_proto_msgTypes[3]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use ScheduleLabResponse.ProtoReflect.Descriptor instead.
func (*ScheduleLabResponse) Descriptor() ([]byte, []int) {
	return file_service_service_proto_rawDescGZIP(), []int{3}
}

func (x *ScheduleLabResponse) GetScheduled() bool {
	if x != nil {
		return x.Scheduled
	}
	return false
}

func (x *ScheduleLabResponse) GetId() string {
	if x != nil {
		return x.Id
	}
	return ""
}

// GetLabRequest contains the string id
// which is the id of the lab to get.
type GetLabRequest struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Id string `protobuf:"bytes,1,opt,name=id,proto3" json:"id,omitempty"`
}

func (x *GetLabRequest) Reset() {
	*x = GetLabRequest{}
	if protoimpl.UnsafeEnabled {
		mi := &file_service_service_proto_msgTypes[4]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *GetLabRequest) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*GetLabRequest) ProtoMessage() {}

func (x *GetLabRequest) ProtoReflect() protoreflect.Message {
	mi := &file_service_service_proto_msgTypes[4]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use GetLabRequest.ProtoReflect.Descriptor instead.
func (*GetLabRequest) Descriptor() ([]byte, []int) {
	return file_service_service_proto_rawDescGZIP(), []int{4}
}

func (x *GetLabRequest) GetId() string {
	if x != nil {
		return x.Id
	}
	return ""
}

// GetLabResponse
type GetLabResponse struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Lab *LabDescription `protobuf:"bytes,1,opt,name=lab,proto3" json:"lab,omitempty"`
}

func (x *GetLabResponse) Reset() {
	*x = GetLabResponse{}
	if protoimpl.UnsafeEnabled {
		mi := &file_service_service_proto_msgTypes[5]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *GetLabResponse) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*GetLabResponse) ProtoMessage() {}

func (x *GetLabResponse) ProtoReflect() protoreflect.Message {
	mi := &file_service_service_proto_msgTypes[5]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use GetLabResponse.ProtoReflect.Descriptor instead.
func (*GetLabResponse) Descriptor() ([]byte, []int) {
	return file_service_service_proto_rawDescGZIP(), []int{5}
}

func (x *GetLabResponse) GetLab() *LabDescription {
	if x != nil {
		return x.Lab
	}
	return nil
}

type LabDescription struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Id            string `protobuf:"bytes,1,opt,name=id,proto3" json:"id,omitempty"`
	Name          string `protobuf:"bytes,2,opt,name=name,proto3" json:"name,omitempty"`
	NumChallenges int32  `protobuf:"varint,3,opt,name=num_challenges,json=numChallenges,proto3" json:"num_challenges,omitempty"`
	NumUsers      int32  `protobuf:"varint,4,opt,name=num_users,json=numUsers,proto3" json:"num_users,omitempty"`
}

func (x *LabDescription) Reset() {
	*x = LabDescription{}
	if protoimpl.UnsafeEnabled {
		mi := &file_service_service_proto_msgTypes[6]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *LabDescription) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*LabDescription) ProtoMessage() {}

func (x *LabDescription) ProtoReflect() protoreflect.Message {
	mi := &file_service_service_proto_msgTypes[6]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use LabDescription.ProtoReflect.Descriptor instead.
func (*LabDescription) Descriptor() ([]byte, []int) {
	return file_service_service_proto_rawDescGZIP(), []int{6}
}

func (x *LabDescription) GetId() string {
	if x != nil {
		return x.Id
	}
	return ""
}

func (x *LabDescription) GetName() string {
	if x != nil {
		return x.Name
	}
	return ""
}

func (x *LabDescription) GetNumChallenges() int32 {
	if x != nil {
		return x.NumChallenges
	}
	return 0
}

func (x *LabDescription) GetNumUsers() int32 {
	if x != nil {
		return x.NumUsers
	}
	return 0
}

// GetLabsRequest
type GetLabsRequest struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Empty *emptypb.Empty `protobuf:"bytes,1,opt,name=empty,proto3" json:"empty,omitempty"`
}

func (x *GetLabsRequest) Reset() {
	*x = GetLabsRequest{}
	if protoimpl.UnsafeEnabled {
		mi := &file_service_service_proto_msgTypes[7]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *GetLabsRequest) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*GetLabsRequest) ProtoMessage() {}

func (x *GetLabsRequest) ProtoReflect() protoreflect.Message {
	mi := &file_service_service_proto_msgTypes[7]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use GetLabsRequest.ProtoReflect.Descriptor instead.
func (*GetLabsRequest) Descriptor() ([]byte, []int) {
	return file_service_service_proto_rawDescGZIP(), []int{7}
}

func (x *GetLabsRequest) GetEmpty() *emptypb.Empty {
	if x != nil {
		return x.Empty
	}
	return nil
}

// GetLabsResponse
type GetLabsResponse struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Labs []*LabDescription `protobuf:"bytes,1,rep,name=labs,proto3" json:"labs,omitempty"`
}

func (x *GetLabsResponse) Reset() {
	*x = GetLabsResponse{}
	if protoimpl.UnsafeEnabled {
		mi := &file_service_service_proto_msgTypes[8]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *GetLabsResponse) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*GetLabsResponse) ProtoMessage() {}

func (x *GetLabsResponse) ProtoReflect() protoreflect.Message {
	mi := &file_service_service_proto_msgTypes[8]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use GetLabsResponse.ProtoReflect.Descriptor instead.
func (*GetLabsResponse) Descriptor() ([]byte, []int) {
	return file_service_service_proto_rawDescGZIP(), []int{8}
}

func (x *GetLabsResponse) GetLabs() []*LabDescription {
	if x != nil {
		return x.Labs
	}
	return nil
}

type RemoveLabRequest struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Id string `protobuf:"bytes,1,opt,name=id,proto3" json:"id,omitempty"`
}

func (x *RemoveLabRequest) Reset() {
	*x = RemoveLabRequest{}
	if protoimpl.UnsafeEnabled {
		mi := &file_service_service_proto_msgTypes[9]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *RemoveLabRequest) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*RemoveLabRequest) ProtoMessage() {}

func (x *RemoveLabRequest) ProtoReflect() protoreflect.Message {
	mi := &file_service_service_proto_msgTypes[9]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use RemoveLabRequest.ProtoReflect.Descriptor instead.
func (*RemoveLabRequest) Descriptor() ([]byte, []int) {
	return file_service_service_proto_rawDescGZIP(), []int{9}
}

func (x *RemoveLabRequest) GetId() string {
	if x != nil {
		return x.Id
	}
	return ""
}

type RemoveLabResponse struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Ok bool `protobuf:"varint,1,opt,name=ok,proto3" json:"ok,omitempty"`
}

func (x *RemoveLabResponse) Reset() {
	*x = RemoveLabResponse{}
	if protoimpl.UnsafeEnabled {
		mi := &file_service_service_proto_msgTypes[10]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *RemoveLabResponse) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*RemoveLabResponse) ProtoMessage() {}

func (x *RemoveLabResponse) ProtoReflect() protoreflect.Message {
	mi := &file_service_service_proto_msgTypes[10]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use RemoveLabResponse.ProtoReflect.Descriptor instead.
func (*RemoveLabResponse) Descriptor() ([]byte, []int) {
	return file_service_service_proto_rawDescGZIP(), []int{10}
}

func (x *RemoveLabResponse) GetOk() bool {
	if x != nil {
		return x.Ok
	}
	return false
}

var File_service_service_proto protoreflect.FileDescriptor

var file_service_service_proto_rawDesc = []byte{
	0x0a, 0x15, 0x73, 0x65, 0x72, 0x76, 0x69, 0x63, 0x65, 0x2f, 0x73, 0x65, 0x72, 0x76, 0x69, 0x63,
	0x65, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x12, 0x07, 0x73, 0x65, 0x72, 0x76, 0x69, 0x63, 0x65,
	0x1a, 0x1b, 0x67, 0x6f, 0x6f, 0x67, 0x6c, 0x65, 0x2f, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x62, 0x75,
	0x66, 0x2f, 0x65, 0x6d, 0x70, 0x74, 0x79, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x22, 0x27, 0x0a,
	0x13, 0x48, 0x61, 0x76, 0x65, 0x43, 0x61, 0x70, 0x61, 0x63, 0x69, 0x74, 0x79, 0x52, 0x65, 0x71,
	0x75, 0x65, 0x73, 0x74, 0x12, 0x10, 0x0a, 0x03, 0x6c, 0x61, 0x62, 0x18, 0x01, 0x20, 0x01, 0x28,
	0x09, 0x52, 0x03, 0x6c, 0x61, 0x62, 0x22, 0x38, 0x0a, 0x14, 0x48, 0x61, 0x76, 0x65, 0x43, 0x61,
	0x70, 0x61, 0x63, 0x69, 0x74, 0x79, 0x52, 0x65, 0x73, 0x70, 0x6f, 0x6e, 0x73, 0x65, 0x12, 0x20,
	0x0a, 0x0b, 0x68, 0x61, 0x73, 0x43, 0x61, 0x70, 0x61, 0x63, 0x69, 0x74, 0x79, 0x18, 0x01, 0x20,
	0x01, 0x28, 0x08, 0x52, 0x0b, 0x68, 0x61, 0x73, 0x43, 0x61, 0x70, 0x61, 0x63, 0x69, 0x74, 0x79,
	0x22, 0x26, 0x0a, 0x12, 0x53, 0x63, 0x68, 0x65, 0x64, 0x75, 0x6c, 0x65, 0x4c, 0x61, 0x62, 0x52,
	0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x12, 0x10, 0x0a, 0x03, 0x6c, 0x61, 0x62, 0x18, 0x01, 0x20,
	0x01, 0x28, 0x09, 0x52, 0x03, 0x6c, 0x61, 0x62, 0x22, 0x43, 0x0a, 0x13, 0x53, 0x63, 0x68, 0x65,
	0x64, 0x75, 0x6c, 0x65, 0x4c, 0x61, 0x62, 0x52, 0x65, 0x73, 0x70, 0x6f, 0x6e, 0x73, 0x65, 0x12,
	0x1c, 0x0a, 0x09, 0x73, 0x63, 0x68, 0x65, 0x64, 0x75, 0x6c, 0x65, 0x64, 0x18, 0x01, 0x20, 0x01,
	0x28, 0x08, 0x52, 0x09, 0x73, 0x63, 0x68, 0x65, 0x64, 0x75, 0x6c, 0x65, 0x64, 0x12, 0x0e, 0x0a,
	0x02, 0x69, 0x64, 0x18, 0x02, 0x20, 0x01, 0x28, 0x09, 0x52, 0x02, 0x69, 0x64, 0x22, 0x1f, 0x0a,
	0x0d, 0x47, 0x65, 0x74, 0x4c, 0x61, 0x62, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x12, 0x0e,
	0x0a, 0x02, 0x69, 0x64, 0x18, 0x01, 0x20, 0x01, 0x28, 0x09, 0x52, 0x02, 0x69, 0x64, 0x22, 0x3b,
	0x0a, 0x0e, 0x47, 0x65, 0x74, 0x4c, 0x61, 0x62, 0x52, 0x65, 0x73, 0x70, 0x6f, 0x6e, 0x73, 0x65,
	0x12, 0x29, 0x0a, 0x03, 0x6c, 0x61, 0x62, 0x18, 0x01, 0x20, 0x01, 0x28, 0x0b, 0x32, 0x17, 0x2e,
	0x73, 0x65, 0x72, 0x76, 0x69, 0x63, 0x65, 0x2e, 0x4c, 0x61, 0x62, 0x44, 0x65, 0x73, 0x63, 0x72,
	0x69, 0x70, 0x74, 0x69, 0x6f, 0x6e, 0x52, 0x03, 0x6c, 0x61, 0x62, 0x22, 0x78, 0x0a, 0x0e, 0x4c,
	0x61, 0x62, 0x44, 0x65, 0x73, 0x63, 0x72, 0x69, 0x70, 0x74, 0x69, 0x6f, 0x6e, 0x12, 0x0e, 0x0a,
	0x02, 0x69, 0x64, 0x18, 0x01, 0x20, 0x01, 0x28, 0x09, 0x52, 0x02, 0x69, 0x64, 0x12, 0x12, 0x0a,
	0x04, 0x6e, 0x61, 0x6d, 0x65, 0x18, 0x02, 0x20, 0x01, 0x28, 0x09, 0x52, 0x04, 0x6e, 0x61, 0x6d,
	0x65, 0x12, 0x25, 0x0a, 0x0e, 0x6e, 0x75, 0x6d, 0x5f, 0x63, 0x68, 0x61, 0x6c, 0x6c, 0x65, 0x6e,
	0x67, 0x65, 0x73, 0x18, 0x03, 0x20, 0x01, 0x28, 0x05, 0x52, 0x0d, 0x6e, 0x75, 0x6d, 0x43, 0x68,
	0x61, 0x6c, 0x6c, 0x65, 0x6e, 0x67, 0x65, 0x73, 0x12, 0x1b, 0x0a, 0x09, 0x6e, 0x75, 0x6d, 0x5f,
	0x75, 0x73, 0x65, 0x72, 0x73, 0x18, 0x04, 0x20, 0x01, 0x28, 0x05, 0x52, 0x08, 0x6e, 0x75, 0x6d,
	0x55, 0x73, 0x65, 0x72, 0x73, 0x22, 0x3e, 0x0a, 0x0e, 0x47, 0x65, 0x74, 0x4c, 0x61, 0x62, 0x73,
	0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x12, 0x2c, 0x0a, 0x05, 0x65, 0x6d, 0x70, 0x74, 0x79,
	0x18, 0x01, 0x20, 0x01, 0x28, 0x0b, 0x32, 0x16, 0x2e, 0x67, 0x6f, 0x6f, 0x67, 0x6c, 0x65, 0x2e,
	0x70, 0x72, 0x6f, 0x74, 0x6f, 0x62, 0x75, 0x66, 0x2e, 0x45, 0x6d, 0x70, 0x74, 0x79, 0x52, 0x05,
	0x65, 0x6d, 0x70, 0x74, 0x79, 0x22, 0x3e, 0x0a, 0x0f, 0x47, 0x65, 0x74, 0x4c, 0x61, 0x62, 0x73,
	0x52, 0x65, 0x73, 0x70, 0x6f, 0x6e, 0x73, 0x65, 0x12, 0x2b, 0x0a, 0x04, 0x6c, 0x61, 0x62, 0x73,
	0x18, 0x01, 0x20, 0x03, 0x28, 0x0b, 0x32, 0x17, 0x2e, 0x73, 0x65, 0x72, 0x76, 0x69, 0x63, 0x65,
	0x2e, 0x4c, 0x61, 0x62, 0x44, 0x65, 0x73, 0x63, 0x72, 0x69, 0x70, 0x74, 0x69, 0x6f, 0x6e, 0x52,
	0x04, 0x6c, 0x61, 0x62, 0x73, 0x22, 0x22, 0x0a, 0x10, 0x52, 0x65, 0x6d, 0x6f, 0x76, 0x65, 0x4c,
	0x61, 0x62, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x12, 0x0e, 0x0a, 0x02, 0x69, 0x64, 0x18,
	0x01, 0x20, 0x01, 0x28, 0x09, 0x52, 0x02, 0x69, 0x64, 0x22, 0x23, 0x0a, 0x11, 0x52, 0x65, 0x6d,
	0x6f, 0x76, 0x65, 0x4c, 0x61, 0x62, 0x52, 0x65, 0x73, 0x70, 0x6f, 0x6e, 0x73, 0x65, 0x12, 0x0e,
	0x0a, 0x02, 0x6f, 0x6b, 0x18, 0x01, 0x20, 0x01, 0x28, 0x08, 0x52, 0x02, 0x6f, 0x6b, 0x32, 0xe7,
	0x02, 0x0a, 0x07, 0x73, 0x65, 0x72, 0x76, 0x69, 0x63, 0x65, 0x12, 0x4d, 0x0a, 0x0c, 0x48, 0x61,
	0x76, 0x65, 0x43, 0x61, 0x70, 0x61, 0x63, 0x69, 0x74, 0x79, 0x12, 0x1c, 0x2e, 0x73, 0x65, 0x72,
	0x76, 0x69, 0x63, 0x65, 0x2e, 0x48, 0x61, 0x76, 0x65, 0x43, 0x61, 0x70, 0x61, 0x63, 0x69, 0x74,
	0x79, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x1a, 0x1d, 0x2e, 0x73, 0x65, 0x72, 0x76, 0x69,
	0x63, 0x65, 0x2e, 0x48, 0x61, 0x76, 0x65, 0x43, 0x61, 0x70, 0x61, 0x63, 0x69, 0x74, 0x79, 0x52,
	0x65, 0x73, 0x70, 0x6f, 0x6e, 0x73, 0x65, 0x22, 0x00, 0x12, 0x4a, 0x0a, 0x0b, 0x53, 0x63, 0x68,
	0x65, 0x64, 0x75, 0x6c, 0x65, 0x4c, 0x61, 0x62, 0x12, 0x1b, 0x2e, 0x73, 0x65, 0x72, 0x76, 0x69,
	0x63, 0x65, 0x2e, 0x53, 0x63, 0x68, 0x65, 0x64, 0x75, 0x6c, 0x65, 0x4c, 0x61, 0x62, 0x52, 0x65,
	0x71, 0x75, 0x65, 0x73, 0x74, 0x1a, 0x1c, 0x2e, 0x73, 0x65, 0x72, 0x76, 0x69, 0x63, 0x65, 0x2e,
	0x53, 0x63, 0x68, 0x65, 0x64, 0x75, 0x6c, 0x65, 0x4c, 0x61, 0x62, 0x52, 0x65, 0x73, 0x70, 0x6f,
	0x6e, 0x73, 0x65, 0x22, 0x00, 0x12, 0x3b, 0x0a, 0x06, 0x47, 0x65, 0x74, 0x4c, 0x61, 0x62, 0x12,
	0x16, 0x2e, 0x73, 0x65, 0x72, 0x76, 0x69, 0x63, 0x65, 0x2e, 0x47, 0x65, 0x74, 0x4c, 0x61, 0x62,
	0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x1a, 0x17, 0x2e, 0x73, 0x65, 0x72, 0x76, 0x69, 0x63,
	0x65, 0x2e, 0x47, 0x65, 0x74, 0x4c, 0x61, 0x62, 0x52, 0x65, 0x73, 0x70, 0x6f, 0x6e, 0x73, 0x65,
	0x22, 0x00, 0x12, 0x3e, 0x0a, 0x07, 0x47, 0x65, 0x74, 0x4c, 0x61, 0x62, 0x73, 0x12, 0x17, 0x2e,
	0x73, 0x65, 0x72, 0x76, 0x69, 0x63, 0x65, 0x2e, 0x47, 0x65, 0x74, 0x4c, 0x61, 0x62, 0x73, 0x52,
	0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x1a, 0x18, 0x2e, 0x73, 0x65, 0x72, 0x76, 0x69, 0x63, 0x65,
	0x2e, 0x47, 0x65, 0x74, 0x4c, 0x61, 0x62, 0x73, 0x52, 0x65, 0x73, 0x70, 0x6f, 0x6e, 0x73, 0x65,
	0x22, 0x00, 0x12, 0x44, 0x0a, 0x09, 0x52, 0x65, 0x6d, 0x6f, 0x76, 0x65, 0x4c, 0x61, 0x62, 0x12,
	0x19, 0x2e, 0x73, 0x65, 0x72, 0x76, 0x69, 0x63, 0x65, 0x2e, 0x52, 0x65, 0x6d, 0x6f, 0x76, 0x65,
	0x4c, 0x61, 0x62, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x1a, 0x1a, 0x2e, 0x73, 0x65, 0x72,
	0x76, 0x69, 0x63, 0x65, 0x2e, 0x52, 0x65, 0x6d, 0x6f, 0x76, 0x65, 0x4c, 0x61, 0x62, 0x52, 0x65,
	0x73, 0x70, 0x6f, 0x6e, 0x73, 0x65, 0x22, 0x00, 0x42, 0x3c, 0x5a, 0x3a, 0x67, 0x69, 0x74, 0x68,
	0x75, 0x62, 0x2e, 0x63, 0x6f, 0x6d, 0x2f, 0x61, 0x6e, 0x64, 0x72, 0x65, 0x61, 0x73, 0x77, 0x61,
	0x63, 0x68, 0x73, 0x2f, 0x62, 0x61, 0x63, 0x68, 0x65, 0x6c, 0x6f, 0x72, 0x73, 0x2d, 0x70, 0x72,
	0x6f, 0x6a, 0x65, 0x63, 0x74, 0x2f, 0x64, 0x61, 0x61, 0x75, 0x6b, 0x69, 0x6e, 0x73, 0x2f, 0x73,
	0x65, 0x72, 0x76, 0x69, 0x63, 0x65, 0x62, 0x06, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x33,
}

var (
	file_service_service_proto_rawDescOnce sync.Once
	file_service_service_proto_rawDescData = file_service_service_proto_rawDesc
)

func file_service_service_proto_rawDescGZIP() []byte {
	file_service_service_proto_rawDescOnce.Do(func() {
		file_service_service_proto_rawDescData = protoimpl.X.CompressGZIP(file_service_service_proto_rawDescData)
	})
	return file_service_service_proto_rawDescData
}

var file_service_service_proto_msgTypes = make([]protoimpl.MessageInfo, 11)
var file_service_service_proto_goTypes = []interface{}{
	(*HaveCapacityRequest)(nil),  // 0: service.HaveCapacityRequest
	(*HaveCapacityResponse)(nil), // 1: service.HaveCapacityResponse
	(*ScheduleLabRequest)(nil),   // 2: service.ScheduleLabRequest
	(*ScheduleLabResponse)(nil),  // 3: service.ScheduleLabResponse
	(*GetLabRequest)(nil),        // 4: service.GetLabRequest
	(*GetLabResponse)(nil),       // 5: service.GetLabResponse
	(*LabDescription)(nil),       // 6: service.LabDescription
	(*GetLabsRequest)(nil),       // 7: service.GetLabsRequest
	(*GetLabsResponse)(nil),      // 8: service.GetLabsResponse
	(*RemoveLabRequest)(nil),     // 9: service.RemoveLabRequest
	(*RemoveLabResponse)(nil),    // 10: service.RemoveLabResponse
	(*emptypb.Empty)(nil),        // 11: google.protobuf.Empty
}
var file_service_service_proto_depIdxs = []int32{
	6,  // 0: service.GetLabResponse.lab:type_name -> service.LabDescription
	11, // 1: service.GetLabsRequest.empty:type_name -> google.protobuf.Empty
	6,  // 2: service.GetLabsResponse.labs:type_name -> service.LabDescription
	0,  // 3: service.service.HaveCapacity:input_type -> service.HaveCapacityRequest
	2,  // 4: service.service.ScheduleLab:input_type -> service.ScheduleLabRequest
	4,  // 5: service.service.GetLab:input_type -> service.GetLabRequest
	7,  // 6: service.service.GetLabs:input_type -> service.GetLabsRequest
	9,  // 7: service.service.RemoveLab:input_type -> service.RemoveLabRequest
	1,  // 8: service.service.HaveCapacity:output_type -> service.HaveCapacityResponse
	3,  // 9: service.service.ScheduleLab:output_type -> service.ScheduleLabResponse
	5,  // 10: service.service.GetLab:output_type -> service.GetLabResponse
	8,  // 11: service.service.GetLabs:output_type -> service.GetLabsResponse
	10, // 12: service.service.RemoveLab:output_type -> service.RemoveLabResponse
	8,  // [8:13] is the sub-list for method output_type
	3,  // [3:8] is the sub-list for method input_type
	3,  // [3:3] is the sub-list for extension type_name
	3,  // [3:3] is the sub-list for extension extendee
	0,  // [0:3] is the sub-list for field type_name
}

func init() { file_service_service_proto_init() }
func file_service_service_proto_init() {
	if File_service_service_proto != nil {
		return
	}
	if !protoimpl.UnsafeEnabled {
		file_service_service_proto_msgTypes[0].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*HaveCapacityRequest); i {
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
		file_service_service_proto_msgTypes[1].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*HaveCapacityResponse); i {
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
		file_service_service_proto_msgTypes[2].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*ScheduleLabRequest); i {
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
		file_service_service_proto_msgTypes[3].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*ScheduleLabResponse); i {
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
		file_service_service_proto_msgTypes[4].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*GetLabRequest); i {
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
		file_service_service_proto_msgTypes[5].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*GetLabResponse); i {
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
		file_service_service_proto_msgTypes[6].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*LabDescription); i {
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
		file_service_service_proto_msgTypes[7].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*GetLabsRequest); i {
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
		file_service_service_proto_msgTypes[8].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*GetLabsResponse); i {
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
		file_service_service_proto_msgTypes[9].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*RemoveLabRequest); i {
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
		file_service_service_proto_msgTypes[10].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*RemoveLabResponse); i {
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
			RawDescriptor: file_service_service_proto_rawDesc,
			NumEnums:      0,
			NumMessages:   11,
			NumExtensions: 0,
			NumServices:   1,
		},
		GoTypes:           file_service_service_proto_goTypes,
		DependencyIndexes: file_service_service_proto_depIdxs,
		MessageInfos:      file_service_service_proto_msgTypes,
	}.Build()
	File_service_service_proto = out.File
	file_service_service_proto_rawDesc = nil
	file_service_service_proto_goTypes = nil
	file_service_service_proto_depIdxs = nil
}
