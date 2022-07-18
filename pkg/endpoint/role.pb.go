// Code generated by protoc-gen-gogo. DO NOT EDIT.
// source: v1/role.proto

package endpoint

import (
	fmt "fmt"
	proto "github.com/gogo/protobuf/proto"
	_ "idas/pkg/service/models"
	math "math"
)

// Reference imports to suppress errors if they are not otherwise used.
var _ = proto.Marshal
var _ = fmt.Errorf
var _ = math.Inf

// This is a compile-time assertion to ensure that this generated file
// is compatible with the proto package it is being compiled against.
// A compilation error at this line likely means your copy of the
// proto package needs to be updated.
const _ = proto.GoGoProtoPackageIsVersion3 // please upgrade the proto package

type PermissionInfo struct {
	Id                   string   `protobuf:"bytes,1,opt,name=id,proto3" json:"id" valid:"required"`
	UpdateTime           string   `protobuf:"bytes,2,opt,name=updateTime,proto3" json:"updateTime" valid:"required"`
	CreateTime           string   `protobuf:"bytes,3,opt,name=createTime,proto3" json:"createTime" valid:"required"`
	Name                 *string  `protobuf:"bytes,4,opt,name=name,proto3,customtype=string" json:"name,omitempty"`
	Path                 *string  `protobuf:"bytes,5,opt,name=path,proto3,customtype=string" json:"path,omitempty"`
	ParentId             *string  `protobuf:"bytes,6,opt,name=parentId,proto3,customtype=string" json:"parentId,omitempty"`
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_unrecognized     []byte   `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *PermissionInfo) Reset()         { *m = PermissionInfo{} }
func (m *PermissionInfo) String() string { return proto.CompactTextString(m) }
func (*PermissionInfo) ProtoMessage()    {}
func (*PermissionInfo) Descriptor() ([]byte, []int) {
	return fileDescriptor_742aa1c3ad5f8a33, []int{0}
}
func (m *PermissionInfo) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_PermissionInfo.Unmarshal(m, b)
}
func (m *PermissionInfo) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_PermissionInfo.Marshal(b, m, deterministic)
}
func (m *PermissionInfo) XXX_Merge(src proto.Message) {
	xxx_messageInfo_PermissionInfo.Merge(m, src)
}
func (m *PermissionInfo) XXX_Size() int {
	return xxx_messageInfo_PermissionInfo.Size(m)
}
func (m *PermissionInfo) XXX_DiscardUnknown() {
	xxx_messageInfo_PermissionInfo.DiscardUnknown(m)
}

var xxx_messageInfo_PermissionInfo proto.InternalMessageInfo

func (m *PermissionInfo) GetId() string {
	if m != nil {
		return m.Id
	}
	return ""
}

func (m *PermissionInfo) GetUpdateTime() string {
	if m != nil {
		return m.UpdateTime
	}
	return ""
}

func (m *PermissionInfo) GetCreateTime() string {
	if m != nil {
		return m.CreateTime
	}
	return ""
}

type RoleInfo struct {
	Id                   string            `protobuf:"bytes,1,opt,name=id,proto3" json:"id" valid:"required"`
	UpdateTime           string            `protobuf:"bytes,2,opt,name=updateTime,proto3" json:"updateTime" valid:"required"`
	CreateTime           string            `protobuf:"bytes,3,opt,name=createTime,proto3" json:"createTime" valid:"required"`
	Name                 string            `protobuf:"bytes,4,opt,name=name,proto3" json:"name"`
	Description          string            `protobuf:"bytes,5,opt,name=description,proto3" json:"description,omitempty"`
	Permission           []*PermissionInfo `protobuf:"bytes,6,rep,name=permission,proto3" json:"permission,omitempty"`
	XXX_NoUnkeyedLiteral struct{}          `json:"-"`
	XXX_unrecognized     []byte            `json:"-"`
	XXX_sizecache        int32             `json:"-"`
}

func (m *RoleInfo) Reset()         { *m = RoleInfo{} }
func (m *RoleInfo) String() string { return proto.CompactTextString(m) }
func (*RoleInfo) ProtoMessage()    {}
func (*RoleInfo) Descriptor() ([]byte, []int) {
	return fileDescriptor_742aa1c3ad5f8a33, []int{1}
}
func (m *RoleInfo) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_RoleInfo.Unmarshal(m, b)
}
func (m *RoleInfo) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_RoleInfo.Marshal(b, m, deterministic)
}
func (m *RoleInfo) XXX_Merge(src proto.Message) {
	xxx_messageInfo_RoleInfo.Merge(m, src)
}
func (m *RoleInfo) XXX_Size() int {
	return xxx_messageInfo_RoleInfo.Size(m)
}
func (m *RoleInfo) XXX_DiscardUnknown() {
	xxx_messageInfo_RoleInfo.DiscardUnknown(m)
}

var xxx_messageInfo_RoleInfo proto.InternalMessageInfo

func (m *RoleInfo) GetId() string {
	if m != nil {
		return m.Id
	}
	return ""
}

func (m *RoleInfo) GetUpdateTime() string {
	if m != nil {
		return m.UpdateTime
	}
	return ""
}

func (m *RoleInfo) GetCreateTime() string {
	if m != nil {
		return m.CreateTime
	}
	return ""
}

func (m *RoleInfo) GetName() string {
	if m != nil {
		return m.Name
	}
	return ""
}

func (m *RoleInfo) GetDescription() string {
	if m != nil {
		return m.Description
	}
	return ""
}

func (m *RoleInfo) GetPermission() []*PermissionInfo {
	if m != nil {
		return m.Permission
	}
	return nil
}

type GetRolesRequest struct {
	BaseListRequest      `protobuf:"bytes,1,opt,name=BaseListRequest,proto3,embedded=BaseListRequest" json:",omitempty"`
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_unrecognized     []byte   `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *GetRolesRequest) Reset()         { *m = GetRolesRequest{} }
func (m *GetRolesRequest) String() string { return proto.CompactTextString(m) }
func (*GetRolesRequest) ProtoMessage()    {}
func (*GetRolesRequest) Descriptor() ([]byte, []int) {
	return fileDescriptor_742aa1c3ad5f8a33, []int{2}
}
func (m *GetRolesRequest) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_GetRolesRequest.Unmarshal(m, b)
}
func (m *GetRolesRequest) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_GetRolesRequest.Marshal(b, m, deterministic)
}
func (m *GetRolesRequest) XXX_Merge(src proto.Message) {
	xxx_messageInfo_GetRolesRequest.Merge(m, src)
}
func (m *GetRolesRequest) XXX_Size() int {
	return xxx_messageInfo_GetRolesRequest.Size(m)
}
func (m *GetRolesRequest) XXX_DiscardUnknown() {
	xxx_messageInfo_GetRolesRequest.DiscardUnknown(m)
}

var xxx_messageInfo_GetRolesRequest proto.InternalMessageInfo

type GetRolesResponse struct {
	BaseListResponse     `protobuf:"bytes,1,opt,name=BaseListResponse,proto3,embedded=BaseListResponse" json:",omitempty"`
	Data                 []*RoleInfo `protobuf:"bytes,2,rep,name=data,proto3" json:"data,omitempty"`
	XXX_NoUnkeyedLiteral struct{}    `json:"-"`
	XXX_unrecognized     []byte      `json:"-"`
	XXX_sizecache        int32       `json:"-"`
}

func (m *GetRolesResponse) Reset()         { *m = GetRolesResponse{} }
func (m *GetRolesResponse) String() string { return proto.CompactTextString(m) }
func (*GetRolesResponse) ProtoMessage()    {}
func (*GetRolesResponse) Descriptor() ([]byte, []int) {
	return fileDescriptor_742aa1c3ad5f8a33, []int{3}
}
func (m *GetRolesResponse) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_GetRolesResponse.Unmarshal(m, b)
}
func (m *GetRolesResponse) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_GetRolesResponse.Marshal(b, m, deterministic)
}
func (m *GetRolesResponse) XXX_Merge(src proto.Message) {
	xxx_messageInfo_GetRolesResponse.Merge(m, src)
}
func (m *GetRolesResponse) XXX_Size() int {
	return xxx_messageInfo_GetRolesResponse.Size(m)
}
func (m *GetRolesResponse) XXX_DiscardUnknown() {
	xxx_messageInfo_GetRolesResponse.DiscardUnknown(m)
}

var xxx_messageInfo_GetRolesResponse proto.InternalMessageInfo

func (m *GetRolesResponse) GetData() []*RoleInfo {
	if m != nil {
		return m.Data
	}
	return nil
}

type CreateRoleRequest struct {
	Name                 string   `protobuf:"bytes,1,opt,name=name,proto3" json:"name" valid:"required"`
	Description          string   `protobuf:"bytes,2,opt,name=description,proto3" json:"description,omitempty"`
	Permission           []string `protobuf:"bytes,3,rep,name=permission,proto3" json:"permission,omitempty"`
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_unrecognized     []byte   `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *CreateRoleRequest) Reset()         { *m = CreateRoleRequest{} }
func (m *CreateRoleRequest) String() string { return proto.CompactTextString(m) }
func (*CreateRoleRequest) ProtoMessage()    {}
func (*CreateRoleRequest) Descriptor() ([]byte, []int) {
	return fileDescriptor_742aa1c3ad5f8a33, []int{4}
}
func (m *CreateRoleRequest) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_CreateRoleRequest.Unmarshal(m, b)
}
func (m *CreateRoleRequest) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_CreateRoleRequest.Marshal(b, m, deterministic)
}
func (m *CreateRoleRequest) XXX_Merge(src proto.Message) {
	xxx_messageInfo_CreateRoleRequest.Merge(m, src)
}
func (m *CreateRoleRequest) XXX_Size() int {
	return xxx_messageInfo_CreateRoleRequest.Size(m)
}
func (m *CreateRoleRequest) XXX_DiscardUnknown() {
	xxx_messageInfo_CreateRoleRequest.DiscardUnknown(m)
}

var xxx_messageInfo_CreateRoleRequest proto.InternalMessageInfo

func (m *CreateRoleRequest) GetName() string {
	if m != nil {
		return m.Name
	}
	return ""
}

func (m *CreateRoleRequest) GetDescription() string {
	if m != nil {
		return m.Description
	}
	return ""
}

func (m *CreateRoleRequest) GetPermission() []string {
	if m != nil {
		return m.Permission
	}
	return nil
}

type CreateRoleResponse struct {
	BaseResponse         `protobuf:"bytes,1,opt,name=BaseResponse,proto3,embedded=BaseResponse" json:",omitempty"`
	Data                 []*RoleInfo `protobuf:"bytes,2,rep,name=data,proto3" json:"data,omitempty"`
	XXX_NoUnkeyedLiteral struct{}    `json:"-"`
	XXX_unrecognized     []byte      `json:"-"`
	XXX_sizecache        int32       `json:"-"`
}

func (m *CreateRoleResponse) Reset()         { *m = CreateRoleResponse{} }
func (m *CreateRoleResponse) String() string { return proto.CompactTextString(m) }
func (*CreateRoleResponse) ProtoMessage()    {}
func (*CreateRoleResponse) Descriptor() ([]byte, []int) {
	return fileDescriptor_742aa1c3ad5f8a33, []int{5}
}
func (m *CreateRoleResponse) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_CreateRoleResponse.Unmarshal(m, b)
}
func (m *CreateRoleResponse) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_CreateRoleResponse.Marshal(b, m, deterministic)
}
func (m *CreateRoleResponse) XXX_Merge(src proto.Message) {
	xxx_messageInfo_CreateRoleResponse.Merge(m, src)
}
func (m *CreateRoleResponse) XXX_Size() int {
	return xxx_messageInfo_CreateRoleResponse.Size(m)
}
func (m *CreateRoleResponse) XXX_DiscardUnknown() {
	xxx_messageInfo_CreateRoleResponse.DiscardUnknown(m)
}

var xxx_messageInfo_CreateRoleResponse proto.InternalMessageInfo

func (m *CreateRoleResponse) GetData() []*RoleInfo {
	if m != nil {
		return m.Data
	}
	return nil
}

type UpdateRoleRequest struct {
	Id                   string   `protobuf:"bytes,1,opt,name=id,proto3" json:"id" valid:"required"`
	Name                 string   `protobuf:"bytes,2,opt,name=name,proto3" json:"name" valid:"required"`
	Description          string   `protobuf:"bytes,3,opt,name=description,proto3" json:"description,omitempty"`
	Permission           []string `protobuf:"bytes,4,rep,name=permission,proto3" json:"permission,omitempty"`
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_unrecognized     []byte   `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *UpdateRoleRequest) Reset()         { *m = UpdateRoleRequest{} }
func (m *UpdateRoleRequest) String() string { return proto.CompactTextString(m) }
func (*UpdateRoleRequest) ProtoMessage()    {}
func (*UpdateRoleRequest) Descriptor() ([]byte, []int) {
	return fileDescriptor_742aa1c3ad5f8a33, []int{6}
}
func (m *UpdateRoleRequest) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_UpdateRoleRequest.Unmarshal(m, b)
}
func (m *UpdateRoleRequest) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_UpdateRoleRequest.Marshal(b, m, deterministic)
}
func (m *UpdateRoleRequest) XXX_Merge(src proto.Message) {
	xxx_messageInfo_UpdateRoleRequest.Merge(m, src)
}
func (m *UpdateRoleRequest) XXX_Size() int {
	return xxx_messageInfo_UpdateRoleRequest.Size(m)
}
func (m *UpdateRoleRequest) XXX_DiscardUnknown() {
	xxx_messageInfo_UpdateRoleRequest.DiscardUnknown(m)
}

var xxx_messageInfo_UpdateRoleRequest proto.InternalMessageInfo

func (m *UpdateRoleRequest) GetId() string {
	if m != nil {
		return m.Id
	}
	return ""
}

func (m *UpdateRoleRequest) GetName() string {
	if m != nil {
		return m.Name
	}
	return ""
}

func (m *UpdateRoleRequest) GetDescription() string {
	if m != nil {
		return m.Description
	}
	return ""
}

func (m *UpdateRoleRequest) GetPermission() []string {
	if m != nil {
		return m.Permission
	}
	return nil
}

type UpdateRoleResponse struct {
	BaseListResponse     `protobuf:"bytes,1,opt,name=BaseListResponse,proto3,embedded=BaseListResponse" json:",omitempty"`
	Data                 []*RoleInfo `protobuf:"bytes,2,rep,name=data,proto3" json:"data,omitempty"`
	XXX_NoUnkeyedLiteral struct{}    `json:"-"`
	XXX_unrecognized     []byte      `json:"-"`
	XXX_sizecache        int32       `json:"-"`
}

func (m *UpdateRoleResponse) Reset()         { *m = UpdateRoleResponse{} }
func (m *UpdateRoleResponse) String() string { return proto.CompactTextString(m) }
func (*UpdateRoleResponse) ProtoMessage()    {}
func (*UpdateRoleResponse) Descriptor() ([]byte, []int) {
	return fileDescriptor_742aa1c3ad5f8a33, []int{7}
}
func (m *UpdateRoleResponse) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_UpdateRoleResponse.Unmarshal(m, b)
}
func (m *UpdateRoleResponse) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_UpdateRoleResponse.Marshal(b, m, deterministic)
}
func (m *UpdateRoleResponse) XXX_Merge(src proto.Message) {
	xxx_messageInfo_UpdateRoleResponse.Merge(m, src)
}
func (m *UpdateRoleResponse) XXX_Size() int {
	return xxx_messageInfo_UpdateRoleResponse.Size(m)
}
func (m *UpdateRoleResponse) XXX_DiscardUnknown() {
	xxx_messageInfo_UpdateRoleResponse.DiscardUnknown(m)
}

var xxx_messageInfo_UpdateRoleResponse proto.InternalMessageInfo

func (m *UpdateRoleResponse) GetData() []*RoleInfo {
	if m != nil {
		return m.Data
	}
	return nil
}

type DeleteRoleRequest struct {
	Id                   string   `protobuf:"bytes,1,opt,name=id,proto3" json:"id" valid:"required"`
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_unrecognized     []byte   `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *DeleteRoleRequest) Reset()         { *m = DeleteRoleRequest{} }
func (m *DeleteRoleRequest) String() string { return proto.CompactTextString(m) }
func (*DeleteRoleRequest) ProtoMessage()    {}
func (*DeleteRoleRequest) Descriptor() ([]byte, []int) {
	return fileDescriptor_742aa1c3ad5f8a33, []int{8}
}
func (m *DeleteRoleRequest) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_DeleteRoleRequest.Unmarshal(m, b)
}
func (m *DeleteRoleRequest) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_DeleteRoleRequest.Marshal(b, m, deterministic)
}
func (m *DeleteRoleRequest) XXX_Merge(src proto.Message) {
	xxx_messageInfo_DeleteRoleRequest.Merge(m, src)
}
func (m *DeleteRoleRequest) XXX_Size() int {
	return xxx_messageInfo_DeleteRoleRequest.Size(m)
}
func (m *DeleteRoleRequest) XXX_DiscardUnknown() {
	xxx_messageInfo_DeleteRoleRequest.DiscardUnknown(m)
}

var xxx_messageInfo_DeleteRoleRequest proto.InternalMessageInfo

func (m *DeleteRoleRequest) GetId() string {
	if m != nil {
		return m.Id
	}
	return ""
}

type GetPermissionsRequest struct {
	BaseListRequest      `protobuf:"bytes,1,opt,name=BaseListRequest,proto3,embedded=BaseListRequest" json:",omitempty"`
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_unrecognized     []byte   `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *GetPermissionsRequest) Reset()         { *m = GetPermissionsRequest{} }
func (m *GetPermissionsRequest) String() string { return proto.CompactTextString(m) }
func (*GetPermissionsRequest) ProtoMessage()    {}
func (*GetPermissionsRequest) Descriptor() ([]byte, []int) {
	return fileDescriptor_742aa1c3ad5f8a33, []int{9}
}
func (m *GetPermissionsRequest) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_GetPermissionsRequest.Unmarshal(m, b)
}
func (m *GetPermissionsRequest) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_GetPermissionsRequest.Marshal(b, m, deterministic)
}
func (m *GetPermissionsRequest) XXX_Merge(src proto.Message) {
	xxx_messageInfo_GetPermissionsRequest.Merge(m, src)
}
func (m *GetPermissionsRequest) XXX_Size() int {
	return xxx_messageInfo_GetPermissionsRequest.Size(m)
}
func (m *GetPermissionsRequest) XXX_DiscardUnknown() {
	xxx_messageInfo_GetPermissionsRequest.DiscardUnknown(m)
}

var xxx_messageInfo_GetPermissionsRequest proto.InternalMessageInfo

type GetPermissionsResponse struct {
	BaseListResponse     `protobuf:"bytes,1,opt,name=BaseListResponse,proto3,embedded=BaseListResponse" json:",omitempty"`
	Data                 []*PermissionInfo `protobuf:"bytes,2,rep,name=data,proto3" json:"data,omitempty"`
	XXX_NoUnkeyedLiteral struct{}          `json:"-"`
	XXX_unrecognized     []byte            `json:"-"`
	XXX_sizecache        int32             `json:"-"`
}

func (m *GetPermissionsResponse) Reset()         { *m = GetPermissionsResponse{} }
func (m *GetPermissionsResponse) String() string { return proto.CompactTextString(m) }
func (*GetPermissionsResponse) ProtoMessage()    {}
func (*GetPermissionsResponse) Descriptor() ([]byte, []int) {
	return fileDescriptor_742aa1c3ad5f8a33, []int{10}
}
func (m *GetPermissionsResponse) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_GetPermissionsResponse.Unmarshal(m, b)
}
func (m *GetPermissionsResponse) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_GetPermissionsResponse.Marshal(b, m, deterministic)
}
func (m *GetPermissionsResponse) XXX_Merge(src proto.Message) {
	xxx_messageInfo_GetPermissionsResponse.Merge(m, src)
}
func (m *GetPermissionsResponse) XXX_Size() int {
	return xxx_messageInfo_GetPermissionsResponse.Size(m)
}
func (m *GetPermissionsResponse) XXX_DiscardUnknown() {
	xxx_messageInfo_GetPermissionsResponse.DiscardUnknown(m)
}

var xxx_messageInfo_GetPermissionsResponse proto.InternalMessageInfo

func (m *GetPermissionsResponse) GetData() []*PermissionInfo {
	if m != nil {
		return m.Data
	}
	return nil
}

func init() {
	proto.RegisterType((*PermissionInfo)(nil), "idas.endpoint.PermissionInfo")
	proto.RegisterType((*RoleInfo)(nil), "idas.endpoint.RoleInfo")
	proto.RegisterType((*GetRolesRequest)(nil), "idas.endpoint.GetRolesRequest")
	proto.RegisterType((*GetRolesResponse)(nil), "idas.endpoint.GetRolesResponse")
	proto.RegisterType((*CreateRoleRequest)(nil), "idas.endpoint.CreateRoleRequest")
	proto.RegisterType((*CreateRoleResponse)(nil), "idas.endpoint.CreateRoleResponse")
	proto.RegisterType((*UpdateRoleRequest)(nil), "idas.endpoint.UpdateRoleRequest")
	proto.RegisterType((*UpdateRoleResponse)(nil), "idas.endpoint.UpdateRoleResponse")
	proto.RegisterType((*DeleteRoleRequest)(nil), "idas.endpoint.DeleteRoleRequest")
	proto.RegisterType((*GetPermissionsRequest)(nil), "idas.endpoint.GetPermissionsRequest")
	proto.RegisterType((*GetPermissionsResponse)(nil), "idas.endpoint.GetPermissionsResponse")
}

func init() { proto.RegisterFile("v1/role.proto", fileDescriptor_742aa1c3ad5f8a33) }

var fileDescriptor_742aa1c3ad5f8a33 = []byte{
	// 588 bytes of a gzipped FileDescriptorProto
	0x1f, 0x8b, 0x08, 0x00, 0x00, 0x00, 0x00, 0x00, 0x02, 0xff, 0xd4, 0x95, 0xcf, 0x6e, 0xd3, 0x40,
	0x10, 0xc6, 0x6b, 0x27, 0x44, 0xe9, 0x84, 0xd2, 0x64, 0x05, 0xc5, 0x0a, 0x21, 0x8e, 0x7c, 0x40,
	0x11, 0xa0, 0x84, 0x96, 0x1b, 0x08, 0x21, 0x05, 0xa4, 0xaa, 0x12, 0x87, 0xca, 0x02, 0x21, 0x71,
	0xc2, 0xcd, 0x0e, 0x61, 0x45, 0xec, 0x75, 0xbd, 0x9b, 0x4a, 0x79, 0x02, 0xde, 0x81, 0x33, 0xe2,
	0xc4, 0x01, 0x89, 0x33, 0x77, 0x8e, 0x9c, 0x39, 0xf8, 0x01, 0x72, 0xe4, 0x09, 0x90, 0xd7, 0x76,
	0xfd, 0x27, 0xa9, 0xa2, 0xf6, 0x50, 0xc1, 0xcd, 0x99, 0xf9, 0xe6, 0x1b, 0xfb, 0x37, 0xb3, 0x1b,
	0xd8, 0x3a, 0xd9, 0x1d, 0x06, 0x7c, 0x8a, 0x03, 0x3f, 0xe0, 0x92, 0x93, 0x2d, 0x46, 0x1d, 0x31,
	0x40, 0x8f, 0xfa, 0x9c, 0x79, 0xb2, 0x7d, 0x7d, 0xc2, 0x27, 0x5c, 0x65, 0x86, 0xd1, 0x53, 0x2c,
	0x6a, 0x47, 0x35, 0x47, 0x8e, 0x48, 0x6a, 0xda, 0x44, 0xce, 0x7d, 0x14, 0x43, 0x97, 0x53, 0x9c,
	0x8a, 0x38, 0x66, 0x7d, 0xd5, 0xe1, 0xda, 0x21, 0x06, 0x2e, 0x13, 0x82, 0x71, 0xef, 0xc0, 0x7b,
	0xc7, 0xc9, 0x5d, 0xd0, 0x19, 0x35, 0xb4, 0x9e, 0xd6, 0xdf, 0x1c, 0xb5, 0x17, 0xa1, 0xa9, 0x33,
	0xfa, 0x27, 0x34, 0x9b, 0x27, 0xce, 0x94, 0xd1, 0x47, 0x56, 0x80, 0xc7, 0x33, 0x16, 0x20, 0xb5,
	0x6c, 0x9d, 0x51, 0x32, 0x02, 0x98, 0xf9, 0xd4, 0x91, 0xf8, 0x92, 0xb9, 0x68, 0xe8, 0xaa, 0xc6,
	0x5a, 0x84, 0x66, 0x2e, 0xba, 0xb2, 0x36, 0x97, 0x8f, 0x3c, 0xc6, 0x01, 0xa6, 0x1e, 0x95, 0xcc,
	0x23, 0x8b, 0xae, 0xf6, 0xc8, 0xf2, 0xa4, 0x0b, 0x55, 0xcf, 0x71, 0xd1, 0xa8, 0xaa, 0x6a, 0xf8,
	0x1d, 0x9a, 0x35, 0x21, 0x03, 0xe6, 0x4d, 0x6c, 0x15, 0x8f, 0xf2, 0xbe, 0x23, 0xdf, 0x1b, 0x57,
	0x96, 0xf3, 0x51, 0x9c, 0xdc, 0x81, 0xba, 0xef, 0x04, 0xe8, 0xc9, 0x03, 0x6a, 0xd4, 0x96, 0x34,
	0xa7, 0x39, 0xeb, 0x87, 0x0e, 0x75, 0x9b, 0x4f, 0xf1, 0xbf, 0x05, 0xd5, 0x29, 0x80, 0xaa, 0x2f,
	0x42, 0x53, 0xfd, 0x4e, 0x30, 0xf5, 0xa0, 0x41, 0x51, 0x8c, 0x03, 0xe6, 0x4b, 0xc6, 0xbd, 0x98,
	0x96, 0x9d, 0x0f, 0x91, 0x27, 0x00, 0xfe, 0xe9, 0xba, 0x18, 0xb5, 0x5e, 0xa5, 0xdf, 0xd8, 0xbb,
	0x3d, 0x28, 0x2c, 0xe3, 0xa0, 0xb8, 0x4f, 0x76, 0xae, 0xc0, 0x12, 0xb0, 0xbd, 0x8f, 0x32, 0x22,
	0x28, 0x6c, 0x3c, 0x9e, 0xa1, 0x90, 0xe4, 0x2d, 0x6c, 0x8f, 0x1c, 0x81, 0x2f, 0x98, 0x90, 0x49,
	0x48, 0x21, 0x6d, 0xec, 0x75, 0x4b, 0xb6, 0x25, 0xd5, 0x68, 0xe7, 0x67, 0x68, 0x6e, 0xfc, 0x0a,
	0x4d, 0x2d, 0x42, 0x70, 0x9f, 0xbb, 0x4c, 0xa2, 0xeb, 0xcb, 0xb9, 0x5d, 0xb6, 0xb3, 0x3e, 0x6b,
	0xd0, 0xcc, 0xba, 0x0a, 0x9f, 0x7b, 0x02, 0xc9, 0x18, 0x9a, 0x99, 0x2e, 0x8e, 0x25, 0x7d, 0xcd,
	0x33, 0xfb, 0xc6, 0xb2, 0x33, 0x1b, 0x2f, 0x19, 0x92, 0x7b, 0x50, 0xa5, 0x8e, 0x74, 0x0c, 0x5d,
	0x71, 0xba, 0x59, 0x32, 0x4e, 0x17, 0xc9, 0x56, 0x22, 0xeb, 0xa3, 0x06, 0xad, 0x67, 0x6a, 0x52,
	0x51, 0x22, 0xc5, 0xf3, 0x20, 0x19, 0x58, 0xbc, 0x66, 0x9d, 0x74, 0x60, 0x2b, 0x07, 0xbd, 0x72,
	0x88, 0xfa, 0xf2, 0x10, 0xbb, 0x85, 0x21, 0x56, 0x7a, 0x95, 0xfe, 0x66, 0x61, 0x4a, 0x9f, 0x34,
	0x20, 0xf9, 0x37, 0x49, 0xbe, 0xe6, 0x35, 0x5c, 0x8d, 0xbe, 0xb0, 0x84, 0xeb, 0xd6, 0x0a, 0x5c,
	0x6b, 0x51, 0x15, 0x8c, 0xce, 0x87, 0xe9, 0xbb, 0x06, 0xad, 0x57, 0xea, 0x50, 0xe4, 0x31, 0x9d,
	0xe7, 0x2c, 0xa6, 0x48, 0xf5, 0x8b, 0x22, 0xad, 0xac, 0x43, 0x5a, 0x5d, 0x42, 0xfa, 0x45, 0x03,
	0x92, 0x7f, 0xeb, 0x7f, 0x76, 0x0b, 0x9f, 0x42, 0xeb, 0x39, 0x4e, 0xf1, 0xc2, 0x74, 0xad, 0x39,
	0xdc, 0xd8, 0x47, 0x99, 0xdd, 0x01, 0x97, 0x78, 0xd0, 0xbf, 0x69, 0xb0, 0x53, 0xee, 0x7d, 0x99,
	0xa0, 0x77, 0x0b, 0xa0, 0xd7, 0x5c, 0x8b, 0x4a, 0x3a, 0xea, 0xbc, 0x69, 0x47, 0xaa, 0xa1, 0xff,
	0x61, 0x32, 0x4c, 0x95, 0x8f, 0xd3, 0x87, 0xc3, 0x8d, 0xa3, 0x9a, 0xfa, 0x9b, 0x7e, 0xf8, 0x37,
	0x00, 0x00, 0xff, 0xff, 0x8e, 0x02, 0x65, 0x6e, 0xff, 0x07, 0x00, 0x00,
}
