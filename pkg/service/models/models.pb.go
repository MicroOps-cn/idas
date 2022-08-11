// Code generated by protoc-gen-gogo. DO NOT EDIT.
// source: types/models.proto

package models

import (
	fmt "fmt"
	proto "github.com/gogo/protobuf/proto"
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

type AuthMeta_Method int32

const (
	AuthMeta_basic           AuthMeta_Method = 0
	AuthMeta_signature       AuthMeta_Method = 1
	AuthMeta_token           AuthMeta_Method = 2
	AuthMeta_token_signature AuthMeta_Method = 3
)

var AuthMeta_Method_name = map[int32]string{
	0: "basic",
	1: "signature",
	2: "token",
	3: "token_signature",
}

var AuthMeta_Method_value = map[string]int32{
	"basic":           0,
	"signature":       1,
	"token":           2,
	"token_signature": 3,
}

func (x AuthMeta_Method) String() string {
	return proto.EnumName(AuthMeta_Method_name, int32(x))
}

func (AuthMeta_Method) EnumDescriptor() ([]byte, []int) {
	return fileDescriptor_8182e6f222cf83c8, []int{0, 0}
}

type AppMeta_Status int32

const (
	AppMeta_unknown AppMeta_Status = 0
	AppMeta_normal  AppMeta_Status = 1
	AppMeta_disable AppMeta_Status = 2
)

var AppMeta_Status_name = map[int32]string{
	0: "unknown",
	1: "normal",
	2: "disable",
}

var AppMeta_Status_value = map[string]int32{
	"unknown": 0,
	"normal":  1,
	"disable": 2,
}

func (x AppMeta_Status) String() string {
	return proto.EnumName(AppMeta_Status_name, int32(x))
}

func (AppMeta_Status) EnumDescriptor() ([]byte, []int) {
	return fileDescriptor_8182e6f222cf83c8, []int{1, 0}
}

type AppMeta_GrantType int32

const (
	AppMeta_none               AppMeta_GrantType = 0
	AppMeta_authorization_code AppMeta_GrantType = 1
)

var AppMeta_GrantType_name = map[int32]string{
	0: "none",
	1: "authorization_code",
}

var AppMeta_GrantType_value = map[string]int32{
	"none":               0,
	"authorization_code": 1,
}

func (x AppMeta_GrantType) String() string {
	return proto.EnumName(AppMeta_GrantType_name, int32(x))
}

func (AppMeta_GrantType) EnumDescriptor() ([]byte, []int) {
	return fileDescriptor_8182e6f222cf83c8, []int{1, 1}
}

type AppMeta_GrantMode int32

const (
	AppMeta_manual AppMeta_GrantMode = 0
	AppMeta_full   AppMeta_GrantMode = 1
)

var AppMeta_GrantMode_name = map[int32]string{
	0: "manual",
	1: "full",
}

var AppMeta_GrantMode_value = map[string]int32{
	"manual": 0,
	"full":   1,
}

func (x AppMeta_GrantMode) String() string {
	return proto.EnumName(AppMeta_GrantMode_name, int32(x))
}

func (AppMeta_GrantMode) EnumDescriptor() ([]byte, []int) {
	return fileDescriptor_8182e6f222cf83c8, []int{1, 2}
}

type RoleMeta_Type int32

const (
	RoleMeta_user   RoleMeta_Type = 0
	RoleMeta_system RoleMeta_Type = 1
)

var RoleMeta_Type_name = map[int32]string{
	0: "user",
	1: "system",
}

var RoleMeta_Type_value = map[string]int32{
	"user":   0,
	"system": 1,
}

func (x RoleMeta_Type) String() string {
	return proto.EnumName(RoleMeta_Type_name, int32(x))
}

func (RoleMeta_Type) EnumDescriptor() ([]byte, []int) {
	return fileDescriptor_8182e6f222cf83c8, []int{2, 0}
}

type UserMeta_UserStatus int32

const (
	UserMeta_unknown  UserMeta_UserStatus = 0
	UserMeta_normal   UserMeta_UserStatus = 1
	UserMeta_disable  UserMeta_UserStatus = 2
	UserMeta_inactive UserMeta_UserStatus = 3
)

var UserMeta_UserStatus_name = map[int32]string{
	0: "unknown",
	1: "normal",
	2: "disable",
	3: "inactive",
}

var UserMeta_UserStatus_value = map[string]int32{
	"unknown":  0,
	"normal":   1,
	"disable":  2,
	"inactive": 3,
}

func (x UserMeta_UserStatus) String() string {
	return proto.EnumName(UserMeta_UserStatus_name, int32(x))
}

func (UserMeta_UserStatus) EnumDescriptor() ([]byte, []int) {
	return fileDescriptor_8182e6f222cf83c8, []int{3, 0}
}

type AuthMeta struct {
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_unrecognized     []byte   `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *AuthMeta) Reset()         { *m = AuthMeta{} }
func (m *AuthMeta) String() string { return proto.CompactTextString(m) }
func (*AuthMeta) ProtoMessage()    {}
func (*AuthMeta) Descriptor() ([]byte, []int) {
	return fileDescriptor_8182e6f222cf83c8, []int{0}
}
func (m *AuthMeta) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_AuthMeta.Unmarshal(m, b)
}
func (m *AuthMeta) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_AuthMeta.Marshal(b, m, deterministic)
}
func (m *AuthMeta) XXX_Merge(src proto.Message) {
	xxx_messageInfo_AuthMeta.Merge(m, src)
}
func (m *AuthMeta) XXX_Size() int {
	return xxx_messageInfo_AuthMeta.Size(m)
}
func (m *AuthMeta) XXX_DiscardUnknown() {
	xxx_messageInfo_AuthMeta.DiscardUnknown(m)
}

var xxx_messageInfo_AuthMeta proto.InternalMessageInfo

type AppMeta struct {
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_unrecognized     []byte   `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *AppMeta) Reset()         { *m = AppMeta{} }
func (m *AppMeta) String() string { return proto.CompactTextString(m) }
func (*AppMeta) ProtoMessage()    {}
func (*AppMeta) Descriptor() ([]byte, []int) {
	return fileDescriptor_8182e6f222cf83c8, []int{1}
}
func (m *AppMeta) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_AppMeta.Unmarshal(m, b)
}
func (m *AppMeta) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_AppMeta.Marshal(b, m, deterministic)
}
func (m *AppMeta) XXX_Merge(src proto.Message) {
	xxx_messageInfo_AppMeta.Merge(m, src)
}
func (m *AppMeta) XXX_Size() int {
	return xxx_messageInfo_AppMeta.Size(m)
}
func (m *AppMeta) XXX_DiscardUnknown() {
	xxx_messageInfo_AppMeta.DiscardUnknown(m)
}

var xxx_messageInfo_AppMeta proto.InternalMessageInfo

type RoleMeta struct {
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_unrecognized     []byte   `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *RoleMeta) Reset()         { *m = RoleMeta{} }
func (m *RoleMeta) String() string { return proto.CompactTextString(m) }
func (*RoleMeta) ProtoMessage()    {}
func (*RoleMeta) Descriptor() ([]byte, []int) {
	return fileDescriptor_8182e6f222cf83c8, []int{2}
}
func (m *RoleMeta) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_RoleMeta.Unmarshal(m, b)
}
func (m *RoleMeta) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_RoleMeta.Marshal(b, m, deterministic)
}
func (m *RoleMeta) XXX_Merge(src proto.Message) {
	xxx_messageInfo_RoleMeta.Merge(m, src)
}
func (m *RoleMeta) XXX_Size() int {
	return xxx_messageInfo_RoleMeta.Size(m)
}
func (m *RoleMeta) XXX_DiscardUnknown() {
	xxx_messageInfo_RoleMeta.DiscardUnknown(m)
}

var xxx_messageInfo_RoleMeta proto.InternalMessageInfo

type UserMeta struct {
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_unrecognized     []byte   `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *UserMeta) Reset()         { *m = UserMeta{} }
func (m *UserMeta) String() string { return proto.CompactTextString(m) }
func (*UserMeta) ProtoMessage()    {}
func (*UserMeta) Descriptor() ([]byte, []int) {
	return fileDescriptor_8182e6f222cf83c8, []int{3}
}
func (m *UserMeta) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_UserMeta.Unmarshal(m, b)
}
func (m *UserMeta) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_UserMeta.Marshal(b, m, deterministic)
}
func (m *UserMeta) XXX_Merge(src proto.Message) {
	xxx_messageInfo_UserMeta.Merge(m, src)
}
func (m *UserMeta) XXX_Size() int {
	return xxx_messageInfo_UserMeta.Size(m)
}
func (m *UserMeta) XXX_DiscardUnknown() {
	xxx_messageInfo_UserMeta.DiscardUnknown(m)
}

var xxx_messageInfo_UserMeta proto.InternalMessageInfo

func init() {
	proto.RegisterEnum("idas.service.models.AuthMeta_Method", AuthMeta_Method_name, AuthMeta_Method_value)
	proto.RegisterEnum("idas.service.models.AppMeta_Status", AppMeta_Status_name, AppMeta_Status_value)
	proto.RegisterEnum("idas.service.models.AppMeta_GrantType", AppMeta_GrantType_name, AppMeta_GrantType_value)
	proto.RegisterEnum("idas.service.models.AppMeta_GrantMode", AppMeta_GrantMode_name, AppMeta_GrantMode_value)
	proto.RegisterEnum("idas.service.models.RoleMeta_Type", RoleMeta_Type_name, RoleMeta_Type_value)
	proto.RegisterEnum("idas.service.models.UserMeta_UserStatus", UserMeta_UserStatus_name, UserMeta_UserStatus_value)
	proto.RegisterType((*AuthMeta)(nil), "idas.service.models.AuthMeta")
	proto.RegisterType((*AppMeta)(nil), "idas.service.models.AppMeta")
	proto.RegisterType((*RoleMeta)(nil), "idas.service.models.RoleMeta")
	proto.RegisterType((*UserMeta)(nil), "idas.service.models.UserMeta")
}

func init() { proto.RegisterFile("types/models.proto", fileDescriptor_8182e6f222cf83c8) }

var fileDescriptor_8182e6f222cf83c8 = []byte{
	// 308 bytes of a gzipped FileDescriptorProto
	0x1f, 0x8b, 0x08, 0x00, 0x00, 0x00, 0x00, 0x00, 0x02, 0xff, 0x94, 0x91, 0xbf, 0x4e, 0xf3, 0x30,
	0x14, 0xc5, 0x93, 0xb6, 0x5f, 0x9a, 0xde, 0x0f, 0x84, 0xe5, 0x4a, 0x4c, 0x08, 0x81, 0x27, 0x16,
	0xd2, 0x81, 0x91, 0x85, 0x76, 0x61, 0xa1, 0x0c, 0xfc, 0x59, 0x58, 0x2a, 0x37, 0xb9, 0xb4, 0x56,
	0x1d, 0xdf, 0xc8, 0xbe, 0x2e, 0x2a, 0xaf, 0xc0, 0x4b, 0xa3, 0xa4, 0x45, 0xcc, 0x4c, 0xbe, 0xd2,
	0xef, 0xf8, 0xe8, 0xe8, 0x1c, 0x90, 0xbc, 0x6b, 0x30, 0x4c, 0x6a, 0xaa, 0xd0, 0x86, 0xa2, 0xf1,
	0xc4, 0x24, 0xc7, 0xa6, 0xd2, 0xa1, 0x08, 0xe8, 0xb7, 0xa6, 0xc4, 0x62, 0x8f, 0xd4, 0x23, 0xe4,
	0xd3, 0xc8, 0xeb, 0x39, 0xb2, 0x56, 0x33, 0xc8, 0xe6, 0xc8, 0x6b, 0xaa, 0xe4, 0x08, 0xfe, 0x2d,
	0x75, 0x30, 0xa5, 0x48, 0xe4, 0x31, 0x8c, 0x82, 0x59, 0x39, 0xcd, 0xd1, 0xa3, 0x48, 0x5b, 0xc2,
	0xb4, 0x41, 0x27, 0x7a, 0x72, 0x0c, 0x27, 0xdd, 0xb9, 0xf8, 0xe5, 0x7d, 0xf5, 0x95, 0xc2, 0x70,
	0xda, 0x34, 0x9d, 0x5f, 0x01, 0xd9, 0x33, 0x6b, 0x8e, 0x41, 0xfe, 0x87, 0x61, 0x74, 0x1b, 0x47,
	0x1f, 0x4e, 0x24, 0x12, 0x20, 0x73, 0xe4, 0x6b, 0x6d, 0x45, 0xda, 0x82, 0xca, 0x04, 0xbd, 0xb4,
	0x28, 0x7a, 0xea, 0x1a, 0x46, 0xf7, 0x5e, 0x3b, 0x7e, 0xd9, 0x35, 0x28, 0x73, 0x18, 0x38, 0x72,
	0x28, 0x12, 0x79, 0x0a, 0x52, 0x47, 0x5e, 0x93, 0x37, 0x9f, 0x9a, 0x0d, 0xb9, 0x45, 0x49, 0x15,
	0x8a, 0x54, 0x5d, 0x1e, 0xe4, 0x73, 0xaa, 0xb0, 0x35, 0xad, 0xb5, 0x8b, 0xda, 0x8a, 0xa4, 0xfd,
	0xfa, 0x1e, 0xad, 0x15, 0xa9, 0xba, 0x82, 0xfc, 0x89, 0x2c, 0x76, 0x69, 0xce, 0x60, 0xf0, 0x63,
	0x1c, 0x03, 0xfa, 0x7d, 0x90, 0xb0, 0x0b, 0x8c, 0xb5, 0x48, 0xd5, 0x03, 0xe4, 0xaf, 0x01, 0x7d,
	0xa7, 0xbc, 0x03, 0x68, 0xef, 0xbf, 0x64, 0x97, 0x47, 0x90, 0x1b, 0xa7, 0x4b, 0x36, 0x5b, 0x14,
	0xfd, 0xd9, 0xc5, 0xdb, 0x79, 0x5b, 0xf6, 0xa4, 0xd9, 0xac, 0x26, 0x87, 0xc2, 0x0f, 0x5b, 0xdc,
	0xee, 0x9f, 0x65, 0xd6, 0x6d, 0x72, 0xf3, 0x1d, 0x00, 0x00, 0xff, 0xff, 0x72, 0xaa, 0x10, 0x91,
	0xa9, 0x01, 0x00, 0x00,
}