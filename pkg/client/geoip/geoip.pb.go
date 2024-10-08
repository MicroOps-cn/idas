// Code generated by protoc-gen-gogo. DO NOT EDIT.
// source: types/geoip.proto

package geoip

import (
	fmt "fmt"
	github_com_MicroOps_cn_fuck_sets "github.com/MicroOps-cn/fuck/sets"
	proto "github.com/gogo/protobuf/proto"
	types "github.com/gogo/protobuf/types"
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

// Duration from public import google/protobuf/duration.proto
type Duration = types.Duration

type CustomGeoOptions struct {
	Name                 string                                  `protobuf:"bytes,1,opt,name=name,proto3" json:"name,omitempty"`
	Subnets              github_com_MicroOps_cn_fuck_sets.IPNets `protobuf:"bytes,2,opt,name=subnets,proto3,customtype=github.com/MicroOps-cn/fuck/sets.IPNets" json:"subnets"`
	XXX_NoUnkeyedLiteral struct{}                                `json:"-"`
	XXX_unrecognized     []byte                                  `json:"-"`
	XXX_sizecache        int32                                   `json:"-"`
}

func (m *CustomGeoOptions) Reset()         { *m = CustomGeoOptions{} }
func (m *CustomGeoOptions) String() string { return proto.CompactTextString(m) }
func (*CustomGeoOptions) ProtoMessage()    {}
func (*CustomGeoOptions) Descriptor() ([]byte, []int) {
	return fileDescriptor_02369b33dee0b6ba, []int{0}
}
func (m *CustomGeoOptions) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_CustomGeoOptions.Unmarshal(m, b)
}
func (m *CustomGeoOptions) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_CustomGeoOptions.Marshal(b, m, deterministic)
}
func (m *CustomGeoOptions) XXX_Merge(src proto.Message) {
	xxx_messageInfo_CustomGeoOptions.Merge(m, src)
}
func (m *CustomGeoOptions) XXX_Size() int {
	return xxx_messageInfo_CustomGeoOptions.Size(m)
}
func (m *CustomGeoOptions) XXX_DiscardUnknown() {
	xxx_messageInfo_CustomGeoOptions.DiscardUnknown(m)
}

var xxx_messageInfo_CustomGeoOptions proto.InternalMessageInfo

func (m *CustomGeoOptions) GetName() string {
	if m != nil {
		return m.Name
	}
	return ""
}

type GeoIPOptions struct {
	Path                 string              `protobuf:"bytes,1,opt,name=path,proto3" json:"path,omitempty"`
	Custom               []*CustomGeoOptions `protobuf:"bytes,2,rep,name=custom,proto3" json:"custom,omitempty"`
	XXX_NoUnkeyedLiteral struct{}            `json:"-"`
	XXX_unrecognized     []byte              `json:"-"`
	XXX_sizecache        int32               `json:"-"`
}

func (m *GeoIPOptions) Reset()         { *m = GeoIPOptions{} }
func (m *GeoIPOptions) String() string { return proto.CompactTextString(m) }
func (*GeoIPOptions) ProtoMessage()    {}
func (*GeoIPOptions) Descriptor() ([]byte, []int) {
	return fileDescriptor_02369b33dee0b6ba, []int{1}
}
func (m *GeoIPOptions) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_GeoIPOptions.Unmarshal(m, b)
}
func (m *GeoIPOptions) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_GeoIPOptions.Marshal(b, m, deterministic)
}
func (m *GeoIPOptions) XXX_Merge(src proto.Message) {
	xxx_messageInfo_GeoIPOptions.Merge(m, src)
}
func (m *GeoIPOptions) XXX_Size() int {
	return xxx_messageInfo_GeoIPOptions.Size(m)
}
func (m *GeoIPOptions) XXX_DiscardUnknown() {
	xxx_messageInfo_GeoIPOptions.DiscardUnknown(m)
}

var xxx_messageInfo_GeoIPOptions proto.InternalMessageInfo

func (m *GeoIPOptions) GetPath() string {
	if m != nil {
		return m.Path
	}
	return ""
}

func (m *GeoIPOptions) GetCustom() []*CustomGeoOptions {
	if m != nil {
		return m.Custom
	}
	return nil
}

func init() {
	proto.RegisterType((*CustomGeoOptions)(nil), "idas.client.geoip.CustomGeoOptions")
	proto.RegisterType((*GeoIPOptions)(nil), "idas.client.geoip.GeoIPOptions")
}

func init() { proto.RegisterFile("types/geoip.proto", fileDescriptor_02369b33dee0b6ba) }

var fileDescriptor_02369b33dee0b6ba = []byte{
	// 270 bytes of a gzipped FileDescriptorProto
	0x1f, 0x8b, 0x08, 0x00, 0x00, 0x00, 0x00, 0x00, 0x02, 0xff, 0x74, 0x90, 0xb1, 0x4e, 0xc3, 0x30,
	0x10, 0x86, 0x9b, 0x82, 0x8a, 0x30, 0x0c, 0x34, 0x62, 0x88, 0x3a, 0xd0, 0xaa, 0x0c, 0x74, 0xc1,
	0x96, 0x0a, 0x5b, 0xb7, 0x32, 0x54, 0x19, 0xa0, 0x51, 0x46, 0x16, 0x94, 0xb8, 0xae, 0x6b, 0xb5,
	0xc9, 0x99, 0xdc, 0x79, 0xe0, 0x0d, 0x79, 0x06, 0x86, 0x3e, 0x0b, 0xb2, 0x43, 0x25, 0x04, 0xea,
	0x62, 0xfd, 0xf6, 0xe7, 0xbb, 0xfb, 0xef, 0x67, 0x7d, 0xfa, 0xb0, 0x0a, 0x85, 0x56, 0x60, 0x2c,
	0xb7, 0x0d, 0x10, 0xc4, 0x7d, 0xb3, 0x2a, 0x90, 0xcb, 0x9d, 0x51, 0x35, 0xf1, 0x00, 0x06, 0xd7,
	0x1a, 0x34, 0x04, 0x2a, 0xbc, 0x6a, 0x3f, 0x0e, 0x6e, 0x34, 0x80, 0xde, 0x29, 0x11, 0x6e, 0xa5,
	0x5b, 0x8b, 0x95, 0x6b, 0x0a, 0x32, 0x50, 0xb7, 0x7c, 0xfc, 0xce, 0xae, 0x9e, 0x1c, 0x12, 0x54,
	0x0b, 0x05, 0x4b, 0xeb, 0x01, 0xc6, 0x31, 0x3b, 0xad, 0x8b, 0x4a, 0x25, 0xd1, 0x28, 0x9a, 0x9c,
	0xe7, 0x41, 0xc7, 0x29, 0x3b, 0x43, 0x57, 0xd6, 0x8a, 0x30, 0xe9, 0xfa, 0xe7, 0xb9, 0xf8, 0xdc,
	0x0f, 0x3b, 0x5f, 0xfb, 0xe1, 0x9d, 0x36, 0xb4, 0x71, 0x25, 0x97, 0x50, 0x89, 0x67, 0x23, 0x1b,
	0x58, 0x5a, 0xbc, 0x97, 0xb5, 0x58, 0x3b, 0xb9, 0x15, 0xa8, 0x08, 0x79, 0x9a, 0xbd, 0x28, 0xc2,
	0xfc, 0x50, 0x3f, 0x7e, 0x63, 0x97, 0x0b, 0x05, 0x69, 0xf6, 0x6b, 0x9c, 0x2d, 0x68, 0x73, 0x18,
	0xe7, 0x75, 0x3c, 0x63, 0x3d, 0x19, 0x6c, 0x25, 0xdd, 0xd1, 0xc9, 0xe4, 0x62, 0x7a, 0xcb, 0xff,
	0x2d, 0xcc, 0xff, 0xfa, 0xce, 0x7f, 0x4a, 0xe6, 0x8f, 0xaf, 0xd3, 0x23, 0xa6, 0x7c, 0x13, 0x61,
	0xb7, 0x5a, 0xb4, 0x8d, 0xda, 0x48, 0x67, 0xe1, 0xcc, 0x3a, 0x59, 0x54, 0xf6, 0x42, 0x28, 0x0f,
	0xdf, 0x01, 0x00, 0x00, 0xff, 0xff, 0x45, 0x0d, 0x0a, 0x13, 0x72, 0x01, 0x00, 0x00,
}
