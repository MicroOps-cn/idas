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
	AutoLogin            bool      `protobuf:"varint,4,opt,name=auto_login,json=autoLogin,proto3" json:"autoLogin,omitempty"`
	AutoRedirect         bool      `protobuf:"varint,5,opt,name=auto_redirect,json=autoRedirect,proto3" json:"autoRedirect,omitempty"`
	Id                   string    `protobuf:"bytes,6,opt,name=id,proto3" json:"id,omitempty"`
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

func (m *GlobalLoginType) GetAutoRedirect() bool {
	if m != nil {
		return m.AutoRedirect
	}
	return false
}

func (m *GlobalLoginType) GetId() string {
	if m != nil {
		return m.Id
	}
	return ""
}

type GlobalConfig struct {
	LoginType            []*GlobalLoginType `protobuf:"bytes,1,rep,name=login_type,json=loginType,proto3" json:"loginType"`
	Title                string             `protobuf:"bytes,2,opt,name=title,proto3" json:"title,omitempty"`
	SubTitle             string             `protobuf:"bytes,3,opt,name=sub_title,json=subTitle,proto3" json:"subTitle,omitempty"`
	Logo                 string             `protobuf:"bytes,4,opt,name=logo,proto3" json:"logo,omitempty"`
	Copyright            string             `protobuf:"bytes,5,opt,name=copyright,proto3" json:"copyright,omitempty"`
	DefaultLoginType     LoginType          `protobuf:"varint,6,opt,name=DefaultLoginType,proto3,enum=idas.endpoint.LoginType" json:"defaultLoginType"`
	Version              string             `protobuf:"bytes,7,opt,name=version,proto3" json:"version,omitempty"`
	XXX_NoUnkeyedLiteral struct{}           `json:"-"`
	XXX_unrecognized     []byte             `json:"-"`
	XXX_sizecache        int32              `json:"-"`
}

func (m *GlobalConfig) Reset()         { *m = GlobalConfig{} }
func (m *GlobalConfig) String() string { return proto.CompactTextString(m) }
func (*GlobalConfig) ProtoMessage()    {}
func (*GlobalConfig) Descriptor() ([]byte, []int) {
	return fileDescriptor_997cafa9b4dd7474, []int{1}
}
func (m *GlobalConfig) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_GlobalConfig.Unmarshal(m, b)
}
func (m *GlobalConfig) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_GlobalConfig.Marshal(b, m, deterministic)
}
func (m *GlobalConfig) XXX_Merge(src proto.Message) {
	xxx_messageInfo_GlobalConfig.Merge(m, src)
}
func (m *GlobalConfig) XXX_Size() int {
	return xxx_messageInfo_GlobalConfig.Size(m)
}
func (m *GlobalConfig) XXX_DiscardUnknown() {
	xxx_messageInfo_GlobalConfig.DiscardUnknown(m)
}

var xxx_messageInfo_GlobalConfig proto.InternalMessageInfo

func (m *GlobalConfig) GetLoginType() []*GlobalLoginType {
	if m != nil {
		return m.LoginType
	}
	return nil
}

func (m *GlobalConfig) GetTitle() string {
	if m != nil {
		return m.Title
	}
	return ""
}

func (m *GlobalConfig) GetSubTitle() string {
	if m != nil {
		return m.SubTitle
	}
	return ""
}

func (m *GlobalConfig) GetLogo() string {
	if m != nil {
		return m.Logo
	}
	return ""
}

func (m *GlobalConfig) GetCopyright() string {
	if m != nil {
		return m.Copyright
	}
	return ""
}

func (m *GlobalConfig) GetDefaultLoginType() LoginType {
	if m != nil {
		return m.DefaultLoginType
	}
	return LoginType_normal
}

func (m *GlobalConfig) GetVersion() string {
	if m != nil {
		return m.Version
	}
	return ""
}

type GlobalConfigResponse struct {
	BaseResponse         `protobuf:"bytes,1,opt,name=base_response,json=baseResponse,proto3,embedded=base_response" json:",omitempty"`
	Data                 *GlobalConfig `protobuf:"bytes,2,opt,name=data,proto3" json:"data,omitempty"`
	XXX_NoUnkeyedLiteral struct{}      `json:"-"`
	XXX_unrecognized     []byte        `json:"-"`
	XXX_sizecache        int32         `json:"-"`
}

func (m *GlobalConfigResponse) Reset()         { *m = GlobalConfigResponse{} }
func (m *GlobalConfigResponse) String() string { return proto.CompactTextString(m) }
func (*GlobalConfigResponse) ProtoMessage()    {}
func (*GlobalConfigResponse) Descriptor() ([]byte, []int) {
	return fileDescriptor_997cafa9b4dd7474, []int{2}
}
func (m *GlobalConfigResponse) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_GlobalConfigResponse.Unmarshal(m, b)
}
func (m *GlobalConfigResponse) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_GlobalConfigResponse.Marshal(b, m, deterministic)
}
func (m *GlobalConfigResponse) XXX_Merge(src proto.Message) {
	xxx_messageInfo_GlobalConfigResponse.Merge(m, src)
}
func (m *GlobalConfigResponse) XXX_Size() int {
	return xxx_messageInfo_GlobalConfigResponse.Size(m)
}
func (m *GlobalConfigResponse) XXX_DiscardUnknown() {
	xxx_messageInfo_GlobalConfigResponse.DiscardUnknown(m)
}

var xxx_messageInfo_GlobalConfigResponse proto.InternalMessageInfo

func (m *GlobalConfigResponse) GetData() *GlobalConfig {
	if m != nil {
		return m.Data
	}
	return nil
}

func init() {
	proto.RegisterType((*GlobalLoginType)(nil), "idas.endpoint.GlobalLoginType")
	proto.RegisterType((*GlobalConfig)(nil), "idas.endpoint.GlobalConfig")
	proto.RegisterType((*GlobalConfigResponse)(nil), "idas.endpoint.GlobalConfigResponse")
}

func init() { proto.RegisterFile("endpoints/global.proto", fileDescriptor_997cafa9b4dd7474) }

var fileDescriptor_997cafa9b4dd7474 = []byte{
	// 506 bytes of a gzipped FileDescriptorProto
	0x1f, 0x8b, 0x08, 0x00, 0x00, 0x00, 0x00, 0x00, 0x02, 0xff, 0x7c, 0x53, 0xcd, 0x6e, 0xd3, 0x40,
	0x10, 0xae, 0xd3, 0xb4, 0x8d, 0xa7, 0x4d, 0xa9, 0x16, 0x2b, 0xac, 0x02, 0xc2, 0x51, 0x4e, 0x39,
	0x80, 0x2d, 0x12, 0xa9, 0x17, 0x0e, 0x48, 0x06, 0x89, 0x4b, 0x11, 0xc8, 0xea, 0x01, 0xf5, 0x12,
	0xf9, 0x67, 0xeb, 0xae, 0x70, 0xbc, 0x96, 0x77, 0x8d, 0x94, 0xc7, 0xe1, 0x6d, 0xb8, 0xc1, 0x13,
	0xf8, 0x01, 0xfc, 0x08, 0x9c, 0xd0, 0x8e, 0xe3, 0xc4, 0x4d, 0xa4, 0x5e, 0xac, 0xd9, 0x6f, 0xe6,
	0x9b, 0xfd, 0xe6, 0xdb, 0x31, 0x8c, 0x58, 0x16, 0xe7, 0x82, 0x67, 0x4a, 0xba, 0x49, 0x2a, 0xc2,
	0x20, 0x75, 0xf2, 0x42, 0x28, 0x41, 0x86, 0x3c, 0x0e, 0xa4, 0xd3, 0x26, 0xc7, 0x56, 0x22, 0x12,
	0x81, 0x19, 0x57, 0x47, 0x4d, 0xd1, 0x98, 0xee, 0xc8, 0x92, 0x49, 0xc9, 0x45, 0x26, 0x37, 0x19,
	0x6b, 0x97, 0x09, 0x03, 0xc9, 0x1a, 0x74, 0xfa, 0xcf, 0x80, 0x67, 0x9f, 0xf1, 0x96, 0x1b, 0x91,
	0xf0, 0xec, 0x76, 0x9d, 0x33, 0x72, 0x0d, 0x7d, 0xb5, 0xce, 0x19, 0x35, 0x26, 0xc6, 0xec, 0x72,
	0x4e, 0x9d, 0x47, 0xf7, 0x3a, 0xdb, 0x3a, 0x6f, 0x50, 0x57, 0x36, 0x56, 0xfa, 0xf8, 0x25, 0x04,
	0xfa, 0x59, 0xb0, 0x62, 0xb4, 0x37, 0x31, 0x66, 0xa6, 0x8f, 0xb1, 0xc6, 0x78, 0x24, 0x32, 0x7a,
	0xdc, 0x60, 0x3a, 0x26, 0xd7, 0x00, 0x41, 0xa9, 0xc4, 0x32, 0xd5, 0x9d, 0x68, 0x7f, 0x62, 0xcc,
	0x06, 0xde, 0x8b, 0xba, 0xb2, 0x9f, 0x6b, 0x14, 0xdb, 0xbf, 0x11, 0x2b, 0xae, 0xd8, 0x2a, 0x57,
	0x6b, 0xdf, 0xdc, 0x82, 0xe4, 0x03, 0x0c, 0x91, 0x57, 0xb0, 0x98, 0x17, 0x2c, 0x52, 0xf4, 0x04,
	0xa9, 0xe3, 0xba, 0xb2, 0x47, 0x3a, 0xe1, 0x6f, 0xf0, 0x0e, 0xfb, 0xa2, 0x8b, 0x93, 0x4b, 0xe8,
	0xf1, 0x98, 0x9e, 0xa2, 0x94, 0x1e, 0x8f, 0xa7, 0x7f, 0x7a, 0x70, 0xd1, 0x0c, 0xff, 0x51, 0x64,
	0xf7, 0x3c, 0x21, 0x37, 0x00, 0x28, 0x6a, 0xb9, 0x99, 0xff, 0x78, 0x76, 0x3e, 0x7f, 0xbd, 0x37,
	0xff, 0x9e, 0x5b, 0xde, 0xb0, 0xae, 0x6c, 0x33, 0x6d, 0x8f, 0xfe, 0x2e, 0x24, 0x16, 0x9c, 0x28,
	0xae, 0xd2, 0xd6, 0x90, 0xe6, 0x40, 0x16, 0x60, 0xca, 0x32, 0x5c, 0x36, 0x19, 0xb4, 0xc5, 0x1b,
	0xd5, 0x95, 0x4d, 0x64, 0x19, 0xde, 0x6a, 0xac, 0xa3, 0x7e, 0xd0, 0x62, 0xda, 0xc6, 0x54, 0x24,
	0x02, 0xcd, 0x32, 0x7d, 0x8c, 0xc9, 0x2b, 0x30, 0x23, 0x91, 0xaf, 0x0b, 0x9e, 0x3c, 0x34, 0x56,
	0x98, 0xfe, 0x0e, 0x20, 0x77, 0x70, 0xf5, 0x89, 0xdd, 0x07, 0x65, 0xaa, 0xb6, 0x52, 0x71, 0xf2,
	0xa7, 0x1e, 0xd4, 0xaa, 0x2b, 0xfb, 0x2a, 0xde, 0x63, 0xf9, 0x07, 0x7d, 0x08, 0x85, 0xb3, 0x9f,
	0xac, 0xd0, 0xcb, 0x45, 0xcf, 0xf0, 0xde, 0xf6, 0x38, 0xfd, 0x65, 0x80, 0xd5, 0x75, 0xd4, 0x67,
	0x32, 0x17, 0x99, 0x64, 0xe4, 0x3b, 0x0c, 0xf5, 0xd6, 0x2d, 0x8b, 0x0d, 0x80, 0xcb, 0x75, 0x3e,
	0x7f, 0xb9, 0xa7, 0xc5, 0x0b, 0x24, 0x6b, 0x39, 0xde, 0xe8, 0x77, 0x65, 0x1f, 0xfd, 0xad, 0x6c,
	0xa3, 0xae, 0x6c, 0xe8, 0x3e, 0x6a, 0xd8, 0xa9, 0x22, 0x2e, 0xf4, 0xe3, 0x40, 0x05, 0x68, 0xf2,
	0x61, 0xc3, 0x47, 0x62, 0xb0, 0xd0, 0x5b, 0xdc, 0xbd, 0x4b, 0xb8, 0x7a, 0x28, 0x43, 0x27, 0x12,
	0x2b, 0xf7, 0x0b, 0x8f, 0x0a, 0xf1, 0x35, 0x97, 0x6f, 0xa3, 0xcc, 0xd5, 0x54, 0x37, 0xff, 0x91,
	0xb8, 0x2d, 0xfd, 0x7d, 0x1b, 0x7c, 0x3b, 0x0a, 0x4f, 0xf1, 0x87, 0x59, 0xfc, 0x0f, 0x00, 0x00,
	0xff, 0xff, 0xce, 0x95, 0x3d, 0x60, 0x9f, 0x03, 0x00, 0x00,
}
