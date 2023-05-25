// Code generated by protoc-gen-gogo. DO NOT EDIT.
// source: endpoints/global.proto

package endpoint

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

type GlobalLoginType struct {
	Type                 LoginType `protobuf:"varint,1,opt,name=type,proto3,enum=idas.endpoint.LoginType" json:"type"`
	Name                 string    `protobuf:"bytes,2,opt,name=name,proto3" json:"name,omitempty"`
	Icon                 string    `protobuf:"bytes,3,opt,name=icon,proto3" json:"icon,omitempty"`
	AutoLogin            bool      `protobuf:"varint,4,opt,name=auto_login,json=autoLogin,proto3" json:"auto_login,omitempty"`
	Id                   string    `protobuf:"bytes,5,opt,name=id,proto3" json:"id,omitempty"`
	XXX_NoUnkeyedLiteral struct{}  `json:"-"`
	XXX_unrecognized     []byte    `json:"-"`
	XXX_sizecache        int32     `json:"-"`
}

func (m *GlobalLoginType) Reset()         { *m = GlobalLoginType{} }
func (m *GlobalLoginType) String() string { return proto.CompactTextString(m) }
func (*GlobalLoginType) ProtoMessage()    {}
func (*GlobalLoginType) Descriptor() ([]byte, []int) {
	return fileDescriptor_997cafa9b4dd7474, []int{0}
}
func (m *GlobalLoginType) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_GlobalLoginType.Unmarshal(m, b)
}
func (m *GlobalLoginType) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_GlobalLoginType.Marshal(b, m, deterministic)
}
func (m *GlobalLoginType) XXX_Merge(src proto.Message) {
	xxx_messageInfo_GlobalLoginType.Merge(m, src)
}
func (m *GlobalLoginType) XXX_Size() int {
	return xxx_messageInfo_GlobalLoginType.Size(m)
}
func (m *GlobalLoginType) XXX_DiscardUnknown() {
	xxx_messageInfo_GlobalLoginType.DiscardUnknown(m)
}

var xxx_messageInfo_GlobalLoginType proto.InternalMessageInfo

func (m *GlobalLoginType) GetType() LoginType {
	if m != nil {
		return m.Type
	}
	return LoginType_normal
}

func (m *GlobalLoginType) GetName() string {
	if m != nil {
		return m.Name
	}
	return ""
}

func (m *GlobalLoginType) GetIcon() string {
	if m != nil {
		return m.Icon
	}
	return ""
}

func (m *GlobalLoginType) GetAutoLogin() bool {
	if m != nil {
		return m.AutoLogin
	}
	return false
}

func (m *GlobalLoginType) GetId() string {
	if m != nil {
		return m.Id
	}
	return ""
}

type GlobalLoginTypeResponse struct {
	BaseResponse         `protobuf:"bytes,1,opt,name=base_response,json=baseResponse,proto3,embedded=base_response" json:",omitempty"`
	Data                 []*GlobalLoginType `protobuf:"bytes,2,rep,name=data,proto3" json:"data,omitempty"`
	XXX_NoUnkeyedLiteral struct{}           `json:"-"`
	XXX_unrecognized     []byte             `json:"-"`
	XXX_sizecache        int32              `json:"-"`
}

func (m *GlobalLoginTypeResponse) Reset()         { *m = GlobalLoginTypeResponse{} }
func (m *GlobalLoginTypeResponse) String() string { return proto.CompactTextString(m) }
func (*GlobalLoginTypeResponse) ProtoMessage()    {}
func (*GlobalLoginTypeResponse) Descriptor() ([]byte, []int) {
	return fileDescriptor_997cafa9b4dd7474, []int{1}
}
func (m *GlobalLoginTypeResponse) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_GlobalLoginTypeResponse.Unmarshal(m, b)
}
func (m *GlobalLoginTypeResponse) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_GlobalLoginTypeResponse.Marshal(b, m, deterministic)
}
func (m *GlobalLoginTypeResponse) XXX_Merge(src proto.Message) {
	xxx_messageInfo_GlobalLoginTypeResponse.Merge(m, src)
}
func (m *GlobalLoginTypeResponse) XXX_Size() int {
	return xxx_messageInfo_GlobalLoginTypeResponse.Size(m)
}
func (m *GlobalLoginTypeResponse) XXX_DiscardUnknown() {
	xxx_messageInfo_GlobalLoginTypeResponse.DiscardUnknown(m)
}

var xxx_messageInfo_GlobalLoginTypeResponse proto.InternalMessageInfo

func (m *GlobalLoginTypeResponse) GetData() []*GlobalLoginType {
	if m != nil {
		return m.Data
	}
	return nil
}

func init() {
	proto.RegisterType((*GlobalLoginType)(nil), "idas.endpoint.GlobalLoginType")
	proto.RegisterType((*GlobalLoginTypeResponse)(nil), "idas.endpoint.GlobalLoginTypeResponse")
}

func init() { proto.RegisterFile("endpoints/global.proto", fileDescriptor_997cafa9b4dd7474) }

var fileDescriptor_997cafa9b4dd7474 = []byte{
	// 331 bytes of a gzipped FileDescriptorProto
	0x1f, 0x8b, 0x08, 0x00, 0x00, 0x00, 0x00, 0x00, 0x02, 0xff, 0x5c, 0x51, 0x3d, 0x4e, 0xc3, 0x30,
	0x18, 0x6d, 0xd2, 0x80, 0x5a, 0x97, 0x16, 0xc9, 0xaa, 0x8a, 0x55, 0x04, 0x89, 0x3a, 0x65, 0x80,
	0x44, 0xa4, 0x12, 0x0b, 0x5b, 0x16, 0x16, 0x10, 0x28, 0x62, 0x40, 0x2c, 0x95, 0x93, 0x58, 0xc1,
	0xa2, 0xf1, 0x67, 0xd5, 0xee, 0xd0, 0xdb, 0xb0, 0x71, 0x15, 0x46, 0x4e, 0x90, 0x03, 0xf4, 0x14,
	0x28, 0xee, 0x1f, 0x64, 0x89, 0x9e, 0xde, 0xdf, 0xf7, 0x22, 0xa3, 0x11, 0x13, 0xb9, 0x04, 0x2e,
	0xb4, 0x0a, 0x8b, 0x39, 0xa4, 0x74, 0x1e, 0xc8, 0x05, 0x68, 0xc0, 0x7d, 0x9e, 0x53, 0x15, 0xec,
	0xc4, 0xf1, 0xb0, 0x80, 0x02, 0x8c, 0x12, 0xd6, 0x68, 0x63, 0x1a, 0x93, 0x43, 0x58, 0x31, 0xa5,
	0x38, 0x08, 0xb5, 0x55, 0x86, 0x07, 0x25, 0xa5, 0x8a, 0x6d, 0xd8, 0xc9, 0xa7, 0x85, 0x4e, 0xef,
	0xcd, 0x95, 0x07, 0x28, 0xb8, 0x78, 0x59, 0x49, 0x86, 0x6f, 0x91, 0xa3, 0x57, 0x92, 0x11, 0xcb,
	0xb3, 0xfc, 0x41, 0x44, 0x82, 0x7f, 0x77, 0x83, 0xbd, 0x2f, 0xee, 0xac, 0x2b, 0xd7, 0x38, 0x13,
	0xf3, 0xc5, 0x18, 0x39, 0x82, 0x96, 0x8c, 0xd8, 0x9e, 0xe5, 0x77, 0x13, 0x83, 0x6b, 0x8e, 0x67,
	0x20, 0x48, 0x7b, 0xc3, 0xd5, 0x18, 0x5f, 0x20, 0x44, 0x97, 0x1a, 0x66, 0xf3, 0xba, 0x89, 0x38,
	0x9e, 0xe5, 0x77, 0x92, 0x6e, 0xcd, 0x98, 0x6a, 0x3c, 0x40, 0x36, 0xcf, 0xc9, 0x91, 0x09, 0xd8,
	0x3c, 0x9f, 0x7c, 0x59, 0xe8, 0xac, 0x31, 0x31, 0x61, 0x4a, 0x82, 0x50, 0x0c, 0xbf, 0xa2, 0x7e,
	0xfd, 0x33, 0xb3, 0xc5, 0x96, 0x30, 0x9b, 0x7b, 0xd1, 0x79, 0x63, 0x73, 0x4c, 0xd5, 0x3e, 0x13,
	0x8f, 0xbe, 0x2b, 0xb7, 0xf5, 0x53, 0xb9, 0xd6, 0xba, 0x72, 0xd1, 0x15, 0x94, 0x5c, 0xb3, 0x52,
	0xea, 0x55, 0x72, 0x92, 0xfe, 0x71, 0xe1, 0x08, 0x39, 0x39, 0xd5, 0x94, 0xd8, 0x5e, 0xdb, 0xef,
	0x45, 0x97, 0x8d, 0xc2, 0xe6, 0x1e, 0xe3, 0x8d, 0xa7, 0x6f, 0x37, 0x05, 0xd7, 0xef, 0xcb, 0x34,
	0xc8, 0xa0, 0x0c, 0x1f, 0x79, 0xb6, 0x80, 0x27, 0xa9, 0xae, 0x33, 0x11, 0xd6, 0xe9, 0x50, 0x7e,
	0x14, 0xe1, 0xae, 0xe1, 0x6e, 0x07, 0x9e, 0x5b, 0xe9, 0xb1, 0x79, 0x8a, 0xe9, 0x6f, 0x00, 0x00,
	0x00, 0xff, 0xff, 0xfe, 0xee, 0x86, 0xd4, 0xf9, 0x01, 0x00, 0x00,
}