// Code generated by protoc-gen-gogo. DO NOT EDIT.
// source: types/tls.proto

package tls

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

type TLSOptions struct {
	CAFile               string   `protobuf:"bytes,1,opt,name=ca_file,json=caFile,proto3" json:"ca_file,omitempty"`
	CertFile             string   `protobuf:"bytes,2,opt,name=cert_file,json=certFile,proto3" json:"cert_file,omitempty"`
	KeyFile              string   `protobuf:"bytes,3,opt,name=key_file,json=keyFile,proto3" json:"key_file,omitempty"`
	ServerName           string   `protobuf:"bytes,4,opt,name=server_name,json=serverName,proto3" json:"server_name,omitempty"`
	InsecureSkipVerify   bool     `protobuf:"varint,5,opt,name=insecure_skip_verify,json=insecureSkipVerify,proto3" json:"insecure_skip_verify,omitempty"`
	MinVersion           string   `protobuf:"bytes,6,opt,name=min_version,json=minVersion,proto3" json:"min_version,omitempty"`
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_unrecognized     []byte   `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *TLSOptions) Reset()         { *m = TLSOptions{} }
func (m *TLSOptions) String() string { return proto.CompactTextString(m) }
func (*TLSOptions) ProtoMessage()    {}
func (*TLSOptions) Descriptor() ([]byte, []int) {
	return fileDescriptor_9c06a1c6ca94056c, []int{0}
}
func (m *TLSOptions) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_TLSOptions.Unmarshal(m, b)
}
func (m *TLSOptions) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_TLSOptions.Marshal(b, m, deterministic)
}
func (m *TLSOptions) XXX_Merge(src proto.Message) {
	xxx_messageInfo_TLSOptions.Merge(m, src)
}
func (m *TLSOptions) XXX_Size() int {
	return xxx_messageInfo_TLSOptions.Size(m)
}
func (m *TLSOptions) XXX_DiscardUnknown() {
	xxx_messageInfo_TLSOptions.DiscardUnknown(m)
}

var xxx_messageInfo_TLSOptions proto.InternalMessageInfo

func (m *TLSOptions) GetCAFile() string {
	if m != nil {
		return m.CAFile
	}
	return ""
}

func (m *TLSOptions) GetCertFile() string {
	if m != nil {
		return m.CertFile
	}
	return ""
}

func (m *TLSOptions) GetKeyFile() string {
	if m != nil {
		return m.KeyFile
	}
	return ""
}

func (m *TLSOptions) GetServerName() string {
	if m != nil {
		return m.ServerName
	}
	return ""
}

func (m *TLSOptions) GetInsecureSkipVerify() bool {
	if m != nil {
		return m.InsecureSkipVerify
	}
	return false
}

func (m *TLSOptions) GetMinVersion() string {
	if m != nil {
		return m.MinVersion
	}
	return ""
}

func init() {
	proto.RegisterType((*TLSOptions)(nil), "idas.client.tls.TLSOptions")
}

func init() { proto.RegisterFile("types/tls.proto", fileDescriptor_9c06a1c6ca94056c) }

var fileDescriptor_9c06a1c6ca94056c = []byte{
	// 282 bytes of a gzipped FileDescriptorProto
	0x1f, 0x8b, 0x08, 0x00, 0x00, 0x00, 0x00, 0x00, 0x02, 0xff, 0x34, 0xd0, 0xbd, 0x4e, 0xc3, 0x30,
	0x14, 0x05, 0x60, 0xc2, 0x4f, 0x9a, 0x9a, 0xa1, 0x52, 0xd4, 0x21, 0xc0, 0xd0, 0x0a, 0x96, 0x2e,
	0xc4, 0x48, 0x0c, 0x08, 0x31, 0x51, 0x24, 0x26, 0xa0, 0x28, 0x45, 0x1d, 0x58, 0x22, 0xd7, 0xdc,
	0x86, 0xab, 0x38, 0xb6, 0x65, 0xbb, 0x95, 0xf2, 0xb0, 0x30, 0xf0, 0x24, 0xc8, 0xb6, 0x3a, 0xd9,
	0x3a, 0xdf, 0xd1, 0x1d, 0x0e, 0x19, 0xb9, 0x5e, 0x83, 0xa5, 0x4e, 0xd8, 0x52, 0x1b, 0xe5, 0x54,
	0x3e, 0xc2, 0x2f, 0x66, 0x4b, 0x2e, 0x10, 0xa4, 0x2b, 0x9d, 0xb0, 0xe7, 0xe3, 0x46, 0x35, 0x2a,
	0x18, 0xf5, 0xbf, 0x58, 0xbb, 0xfc, 0x49, 0x08, 0xf9, 0x78, 0x59, 0x2e, 0xb4, 0x43, 0x25, 0x6d,
	0x7e, 0x45, 0x06, 0x9c, 0xd5, 0x1b, 0x14, 0x50, 0x24, 0xd3, 0x64, 0x36, 0x9c, 0x93, 0xbf, 0xdf,
	0x49, 0xfa, 0xf4, 0xf8, 0x8c, 0x02, 0xaa, 0x94, 0x33, 0xff, 0xe6, 0x17, 0x64, 0xc8, 0xc1, 0xb8,
	0x58, 0x3b, 0xf4, 0xb5, 0x2a, 0xf3, 0x41, 0xc0, 0x33, 0x92, 0xb5, 0xd0, 0x47, 0x3b, 0x0a, 0x36,
	0x68, 0xa1, 0x0f, 0x34, 0x21, 0xa7, 0x16, 0xcc, 0x0e, 0x4c, 0x2d, 0x59, 0x07, 0xc5, 0x71, 0x50,
	0x12, 0xa3, 0x37, 0xd6, 0x41, 0x7e, 0x43, 0xc6, 0x28, 0x2d, 0xf0, 0xad, 0x81, 0xda, 0xb6, 0xa8,
	0xeb, 0x1d, 0x18, 0xdc, 0xf4, 0xc5, 0xc9, 0x34, 0x99, 0x65, 0x55, 0xbe, 0xb7, 0x65, 0x8b, 0x7a,
	0x15, 0xc4, 0x9f, 0xec, 0x50, 0xfa, 0x9e, 0x45, 0x25, 0x8b, 0x34, 0x9e, 0xec, 0x50, 0xae, 0x62,
	0x32, 0xbf, 0xff, 0xbc, 0x6b, 0xd0, 0x7d, 0x6f, 0xd7, 0x25, 0x57, 0x1d, 0x7d, 0x45, 0x6e, 0xd4,
	0x42, 0xdb, 0x6b, 0x2e, 0xa9, 0xdf, 0x87, 0xea, 0xb6, 0xa1, 0x71, 0x23, 0x8a, 0xd2, 0x81, 0x91,
	0x4c, 0xf8, 0x0d, 0x1f, 0x9c, 0xb0, 0xef, 0x07, 0xeb, 0x34, 0x6c, 0x74, 0xfb, 0x1f, 0x00, 0x00,
	0xff, 0xff, 0x5c, 0xcb, 0xaa, 0xed, 0x5d, 0x01, 0x00, 0x00,
}