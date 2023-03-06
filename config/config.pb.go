// Code generated by protoc-gen-gogo. DO NOT EDIT.
// source: config.proto

package config

import (
	fmt "fmt"
	_ "github.com/MicroOps-cn/idas/pkg/client/email"
	github_com_MicroOps_cn_idas_pkg_client_email "github.com/MicroOps-cn/idas/pkg/client/email"
	_ "github.com/MicroOps-cn/idas/pkg/client/gorm"
	github_com_MicroOps_cn_idas_pkg_client_gorm "github.com/MicroOps-cn/idas/pkg/client/gorm"
	_ "github.com/MicroOps-cn/idas/pkg/client/ldap"
	github_com_MicroOps_cn_idas_pkg_client_ldap "github.com/MicroOps-cn/idas/pkg/client/ldap"
	_ "github.com/MicroOps-cn/idas/pkg/client/redis"
	github_com_MicroOps_cn_idas_pkg_client_redis "github.com/MicroOps-cn/idas/pkg/client/redis"
	capacity "github.com/MicroOps-cn/idas/pkg/utils/capacity"
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

type StorageRef struct {
	Storage              *Storage `protobuf:"bytes,1,opt,name=storage,proto3" json:"storage,omitempty"`
	Path                 string   `protobuf:"bytes,2,opt,name=path,proto3" json:"path,omitempty"`
	Name                 string   `protobuf:"bytes,3,opt,name=name,proto3" json:"name,omitempty"`
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_unrecognized     []byte   `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *StorageRef) Reset()         { *m = StorageRef{} }
func (m *StorageRef) String() string { return proto.CompactTextString(m) }
func (*StorageRef) ProtoMessage()    {}
func (*StorageRef) Descriptor() ([]byte, []int) {
	return fileDescriptor_3eaf2c85e69e9ea4, []int{0}
}
func (m *StorageRef) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_StorageRef.Unmarshal(m, b)
}
func (m *StorageRef) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_StorageRef.Marshal(b, m, deterministic)
}
func (m *StorageRef) XXX_Merge(src proto.Message) {
	xxx_messageInfo_StorageRef.Merge(m, src)
}
func (m *StorageRef) XXX_Size() int {
	return xxx_messageInfo_StorageRef.Size(m)
}
func (m *StorageRef) XXX_DiscardUnknown() {
	xxx_messageInfo_StorageRef.DiscardUnknown(m)
}

var xxx_messageInfo_StorageRef proto.InternalMessageInfo

func (m *StorageRef) GetStorage() *Storage {
	if m != nil {
		return m.Storage
	}
	return nil
}

func (m *StorageRef) GetPath() string {
	if m != nil {
		return m.Path
	}
	return ""
}

func (m *StorageRef) GetName() string {
	if m != nil {
		return m.Name
	}
	return ""
}

type Storage struct {
	Name string `protobuf:"bytes,1,opt,name=name,proto3" json:"name,omitempty"`
	// Types that are valid to be assigned to Source:
	//	*Storage_Ref
	//	*Storage_Mysql
	//	*Storage_Redis
	//	*Storage_Ldap
	//	*Storage_Sqlite
	Source               isStorage_Source `protobuf_oneof:"source"`
	XXX_NoUnkeyedLiteral struct{}         `json:"-"`
	XXX_unrecognized     []byte           `json:"-"`
	XXX_sizecache        int32            `json:"-"`
}

func (m *Storage) Reset()         { *m = Storage{} }
func (m *Storage) String() string { return proto.CompactTextString(m) }
func (*Storage) ProtoMessage()    {}
func (*Storage) Descriptor() ([]byte, []int) {
	return fileDescriptor_3eaf2c85e69e9ea4, []int{1}
}
func (m *Storage) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_Storage.Unmarshal(m, b)
}
func (m *Storage) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_Storage.Marshal(b, m, deterministic)
}
func (m *Storage) XXX_Merge(src proto.Message) {
	xxx_messageInfo_Storage.Merge(m, src)
}
func (m *Storage) XXX_Size() int {
	return xxx_messageInfo_Storage.Size(m)
}
func (m *Storage) XXX_DiscardUnknown() {
	xxx_messageInfo_Storage.DiscardUnknown(m)
}

var xxx_messageInfo_Storage proto.InternalMessageInfo

type isStorage_Source interface {
	isStorage_Source()
}

type Storage_Ref struct {
	Ref *StorageRef `protobuf:"bytes,10,opt,name=ref,proto3,oneof" json:"ref,omitempty"`
}
type Storage_Mysql struct {
	Mysql *github_com_MicroOps_cn_idas_pkg_client_gorm.MySQLClient `protobuf:"bytes,11,opt,name=mysql,proto3,oneof,customtype=github.com/MicroOps-cn/idas/pkg/client/gorm.MySQLClient" json:"mysql,omitempty"`
}
type Storage_Redis struct {
	Redis *github_com_MicroOps_cn_idas_pkg_client_redis.Client `protobuf:"bytes,12,opt,name=redis,proto3,oneof,customtype=github.com/MicroOps-cn/idas/pkg/client/redis.Client" json:"redis,omitempty"`
}
type Storage_Ldap struct {
	Ldap *github_com_MicroOps_cn_idas_pkg_client_ldap.Client `protobuf:"bytes,13,opt,name=ldap,proto3,oneof,customtype=github.com/MicroOps-cn/idas/pkg/client/ldap.Client" json:"ldap,omitempty"`
}
type Storage_Sqlite struct {
	Sqlite *github_com_MicroOps_cn_idas_pkg_client_gorm.SQLiteClient `protobuf:"bytes,14,opt,name=sqlite,proto3,oneof,customtype=github.com/MicroOps-cn/idas/pkg/client/gorm.SQLiteClient" json:"sqlite,omitempty"`
}

func (*Storage_Ref) isStorage_Source()    {}
func (*Storage_Mysql) isStorage_Source()  {}
func (*Storage_Redis) isStorage_Source()  {}
func (*Storage_Ldap) isStorage_Source()   {}
func (*Storage_Sqlite) isStorage_Source() {}

func (m *Storage) GetSource() isStorage_Source {
	if m != nil {
		return m.Source
	}
	return nil
}

func (m *Storage) GetName() string {
	if m != nil {
		return m.Name
	}
	return ""
}

func (m *Storage) GetRef() *StorageRef {
	if x, ok := m.GetSource().(*Storage_Ref); ok {
		return x.Ref
	}
	return nil
}

// XXX_OneofWrappers is for the internal use of the proto package.
func (*Storage) XXX_OneofWrappers() []interface{} {
	return []interface{}{
		(*Storage_Ref)(nil),
		(*Storage_Mysql)(nil),
		(*Storage_Redis)(nil),
		(*Storage_Ldap)(nil),
		(*Storage_Sqlite)(nil),
	}
}

type Storages struct {
	Default              *Storage   `protobuf:"bytes,1,opt,name=default,proto3" json:"default,omitempty"`
	Session              *Storage   `protobuf:"bytes,2,opt,name=session,proto3" json:"session,omitempty"`
	User                 []*Storage `protobuf:"bytes,3,rep,name=user,proto3" json:"user,omitempty"`
	XXX_NoUnkeyedLiteral struct{}   `json:"-"`
	XXX_unrecognized     []byte     `json:"-"`
	XXX_sizecache        int32      `json:"-"`
}

func (m *Storages) Reset()         { *m = Storages{} }
func (m *Storages) String() string { return proto.CompactTextString(m) }
func (*Storages) ProtoMessage()    {}
func (*Storages) Descriptor() ([]byte, []int) {
	return fileDescriptor_3eaf2c85e69e9ea4, []int{2}
}
func (m *Storages) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_Storages.Unmarshal(m, b)
}
func (m *Storages) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_Storages.Marshal(b, m, deterministic)
}
func (m *Storages) XXX_Merge(src proto.Message) {
	xxx_messageInfo_Storages.Merge(m, src)
}
func (m *Storages) XXX_Size() int {
	return xxx_messageInfo_Storages.Size(m)
}
func (m *Storages) XXX_DiscardUnknown() {
	xxx_messageInfo_Storages.DiscardUnknown(m)
}

var xxx_messageInfo_Storages proto.InternalMessageInfo

func (m *Storages) GetDefault() *Storage {
	if m != nil {
		return m.Default
	}
	return nil
}

func (m *Storages) GetSession() *Storage {
	if m != nil {
		return m.Session
	}
	return nil
}

func (m *Storages) GetUser() []*Storage {
	if m != nil {
		return m.User
	}
	return nil
}

type GlobalOptions struct {
	MaxUploadSize        *capacity.Capacity `protobuf:"bytes,1,opt,name=max_upload_size,json=maxUploadSize,proto3" json:"max_upload_size,omitempty"`
	MaxBodySize          *capacity.Capacity `protobuf:"bytes,2,opt,name=max_body_size,json=maxBodySize,proto3" json:"max_body_size,omitempty"`
	UploadPath           string             `protobuf:"bytes,3,opt,name=upload_path,json=uploadPath,proto3" json:"upload_path,omitempty"`
	Workspace            string             `protobuf:"bytes,4,opt,name=workspace,proto3" json:"workspace,omitempty"`
	Secret               string             `protobuf:"bytes,5,opt,name=secret,proto3" json:"secret,omitempty"`
	XXX_NoUnkeyedLiteral struct{}           `json:"-"`
	XXX_unrecognized     []byte             `json:"-"`
	XXX_sizecache        int32              `json:"-"`
}

func (m *GlobalOptions) Reset()         { *m = GlobalOptions{} }
func (m *GlobalOptions) String() string { return proto.CompactTextString(m) }
func (*GlobalOptions) ProtoMessage()    {}
func (*GlobalOptions) Descriptor() ([]byte, []int) {
	return fileDescriptor_3eaf2c85e69e9ea4, []int{3}
}
func (m *GlobalOptions) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_GlobalOptions.Unmarshal(m, b)
}
func (m *GlobalOptions) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_GlobalOptions.Marshal(b, m, deterministic)
}
func (m *GlobalOptions) XXX_Merge(src proto.Message) {
	xxx_messageInfo_GlobalOptions.Merge(m, src)
}
func (m *GlobalOptions) XXX_Size() int {
	return xxx_messageInfo_GlobalOptions.Size(m)
}
func (m *GlobalOptions) XXX_DiscardUnknown() {
	xxx_messageInfo_GlobalOptions.DiscardUnknown(m)
}

var xxx_messageInfo_GlobalOptions proto.InternalMessageInfo

func (m *GlobalOptions) GetMaxUploadSize() *capacity.Capacity {
	if m != nil {
		return m.MaxUploadSize
	}
	return nil
}

func (m *GlobalOptions) GetMaxBodySize() *capacity.Capacity {
	if m != nil {
		return m.MaxBodySize
	}
	return nil
}

func (m *GlobalOptions) GetUploadPath() string {
	if m != nil {
		return m.UploadPath
	}
	return ""
}

func (m *GlobalOptions) GetWorkspace() string {
	if m != nil {
		return m.Workspace
	}
	return ""
}

func (m *GlobalOptions) GetSecret() string {
	if m != nil {
		return m.Secret
	}
	return ""
}

type Config struct {
	Storage              *Storages                                                 `protobuf:"bytes,1,opt,name=storage,proto3" json:"storage,omitempty"`
	Global               *GlobalOptions                                            `protobuf:"bytes,2,opt,name=global,proto3" json:"global,omitempty"`
	Smtp                 *github_com_MicroOps_cn_idas_pkg_client_email.SmtpOptions `protobuf:"bytes,3,opt,name=smtp,proto3,customtype=github.com/MicroOps-cn/idas/pkg/client/email.SmtpOptions" json:"smtp,omitempty"`
	XXX_NoUnkeyedLiteral struct{}                                                  `json:"-"`
	XXX_unrecognized     []byte                                                    `json:"-"`
	XXX_sizecache        int32                                                     `json:"-"`
}

func (m *Config) Reset()         { *m = Config{} }
func (m *Config) String() string { return proto.CompactTextString(m) }
func (*Config) ProtoMessage()    {}
func (*Config) Descriptor() ([]byte, []int) {
	return fileDescriptor_3eaf2c85e69e9ea4, []int{4}
}
func (m *Config) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_Config.Unmarshal(m, b)
}
func (m *Config) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_Config.Marshal(b, m, deterministic)
}
func (m *Config) XXX_Merge(src proto.Message) {
	xxx_messageInfo_Config.Merge(m, src)
}
func (m *Config) XXX_Size() int {
	return xxx_messageInfo_Config.Size(m)
}
func (m *Config) XXX_DiscardUnknown() {
	xxx_messageInfo_Config.DiscardUnknown(m)
}

var xxx_messageInfo_Config proto.InternalMessageInfo

func (m *Config) GetStorage() *Storages {
	if m != nil {
		return m.Storage
	}
	return nil
}

func (m *Config) GetGlobal() *GlobalOptions {
	if m != nil {
		return m.Global
	}
	return nil
}

func init() {
	proto.RegisterType((*StorageRef)(nil), "idas.config.StorageRef")
	proto.RegisterType((*Storage)(nil), "idas.config.Storage")
	proto.RegisterType((*Storages)(nil), "idas.config.Storages")
	proto.RegisterType((*GlobalOptions)(nil), "idas.config.GlobalOptions")
	proto.RegisterType((*Config)(nil), "idas.config.Config")
}

func init() { proto.RegisterFile("config.proto", fileDescriptor_3eaf2c85e69e9ea4) }

var fileDescriptor_3eaf2c85e69e9ea4 = []byte{
	// 632 bytes of a gzipped FileDescriptorProto
	0x1f, 0x8b, 0x08, 0x00, 0x00, 0x00, 0x00, 0x00, 0x02, 0xff, 0x9c, 0x54, 0xcb, 0x4e, 0x1b, 0x31,
	0x14, 0x65, 0xc8, 0x03, 0xb8, 0x81, 0x3e, 0x2c, 0xda, 0x5a, 0xa8, 0x6d, 0x50, 0x56, 0x54, 0xa8,
	0x33, 0x12, 0x48, 0xa5, 0x12, 0x2c, 0xaa, 0xb0, 0x28, 0x0b, 0x10, 0x30, 0x51, 0x37, 0xdd, 0x20,
	0x67, 0xc6, 0x19, 0xac, 0xcc, 0xc4, 0x66, 0xec, 0x08, 0xc2, 0x77, 0xf4, 0x07, 0xfa, 0x5b, 0x59,
	0xa4, 0xbb, 0x7e, 0x47, 0x35, 0xd7, 0x1e, 0xd2, 0x50, 0xa0, 0xb4, 0x9b, 0xc4, 0x3e, 0xf7, 0x9c,
	0x39, 0xf6, 0x7d, 0x18, 0x96, 0x23, 0x39, 0xe8, 0x89, 0xc4, 0x57, 0xb9, 0x34, 0x92, 0x34, 0x44,
	0xcc, 0xb4, 0x6f, 0xa1, 0xb5, 0x67, 0x66, 0xa4, 0xb8, 0x0e, 0x12, 0x99, 0x67, 0x36, 0xbc, 0xf6,
	0xdc, 0x22, 0x39, 0x8f, 0x85, 0x76, 0x90, 0x23, 0xa5, 0x31, 0x53, 0x0e, 0x59, 0xb5, 0x48, 0xc4,
	0x14, 0x8b, 0x84, 0x19, 0xcd, 0x4a, 0x79, 0xc6, 0x44, 0x5a, 0x12, 0x13, 0x99, 0x48, 0x5c, 0x06,
	0xc5, 0xca, 0xa2, 0xad, 0x18, 0xa0, 0x63, 0x64, 0xce, 0x12, 0x1e, 0xf2, 0x1e, 0xf1, 0x61, 0x41,
	0xdb, 0x1d, 0xf5, 0xd6, 0xbd, 0x8d, 0xc6, 0xd6, 0xaa, 0xff, 0xdb, 0x11, 0xfd, 0x92, 0x59, 0x92,
	0x08, 0x81, 0xaa, 0x62, 0xe6, 0x9c, 0xce, 0xaf, 0x7b, 0x1b, 0x4b, 0x21, 0xae, 0x0b, 0x6c, 0xc0,
	0x32, 0x4e, 0x2b, 0x16, 0x2b, 0xd6, 0xad, 0xef, 0x55, 0x58, 0xe8, 0x4c, 0x35, 0x18, 0xf7, 0xa6,
	0x71, 0xb2, 0x09, 0x95, 0x9c, 0xf7, 0x28, 0xa0, 0xe7, 0xab, 0x3b, 0x3d, 0x79, 0xef, 0x60, 0x2e,
	0x2c, 0x58, 0x44, 0x43, 0x2d, 0x1b, 0xe9, 0x8b, 0x94, 0x36, 0x90, 0xfe, 0xd6, 0xd1, 0x53, 0xc1,
	0x07, 0xc6, 0xc7, 0xf4, 0x1d, 0x8d, 0x3a, 0xa7, 0x87, 0xc7, 0xca, 0x08, 0x39, 0xd0, 0xed, 0xdd,
	0xf1, 0xa4, 0xb9, 0x93, 0x08, 0x73, 0x3e, 0xec, 0xfa, 0x91, 0xcc, 0x82, 0x23, 0x11, 0xe5, 0xf2,
	0x58, 0xe9, 0xf7, 0xd1, 0x20, 0x28, 0x94, 0x81, 0xea, 0x27, 0x81, 0x55, 0x07, 0x53, 0xf5, 0x3e,
	0x02, 0x07, 0x73, 0xa1, 0xf5, 0x22, 0x0a, 0x6a, 0x58, 0x07, 0xba, 0x8c, 0xa6, 0xcd, 0x19, 0x53,
	0x5b, 0xa1, 0xb0, 0xf8, 0x2d, 0x5d, 0x77, 0xc6, 0x93, 0xe6, 0xf6, 0x23, 0x5d, 0xad, 0x7c, 0xea,
	0x88, 0x7b, 0xd2, 0x87, 0x6a, 0x51, 0x66, 0xba, 0x82, 0x86, 0x6f, 0x66, 0x0c, 0xb1, 0xfe, 0x87,
	0x31, 0x53, 0xa5, 0xdd, 0x87, 0xf1, 0xa4, 0xb9, 0xf5, 0x48, 0x3b, 0x14, 0xdf, 0xb8, 0xa1, 0x09,
	0xb9, 0x84, 0xba, 0xbe, 0x48, 0x85, 0xe1, 0xf4, 0xc9, 0x1d, 0xf7, 0xc3, 0xb4, 0x74, 0x4e, 0x0f,
	0x85, 0xe1, 0xa5, 0xe1, 0xde, 0x78, 0xd2, 0xfc, 0xf8, 0x2f, 0x59, 0xb5, 0xf2, 0x1b, 0x5b, 0x67,
	0xd7, 0x5e, 0x84, 0xba, 0x96, 0xc3, 0x3c, 0xe2, 0xad, 0x6f, 0x1e, 0x2c, 0xba, 0x62, 0xeb, 0xa2,
	0x11, 0x63, 0xde, 0x63, 0xc3, 0xd4, 0x3c, 0xdc, 0x88, 0x8e, 0x84, 0x8d, 0xcb, 0xb5, 0x16, 0x72,
	0x80, 0xbd, 0x78, 0x7f, 0xe3, 0x5a, 0x12, 0xd9, 0x80, 0xea, 0x50, 0xf3, 0x9c, 0x56, 0xd6, 0x2b,
	0xf7, 0x92, 0x91, 0xd1, 0xfa, 0xe9, 0xc1, 0xca, 0xe7, 0x54, 0x76, 0x59, 0xea, 0xae, 0x4e, 0x3e,
	0xc1, 0xd3, 0x8c, 0x5d, 0x9d, 0x0d, 0x55, 0x2a, 0x59, 0x7c, 0xa6, 0xc5, 0x75, 0x39, 0x2c, 0xd4,
	0x7e, 0x46, 0xf5, 0x13, 0x7f, 0x68, 0x44, 0xaa, 0xfd, 0x7d, 0x37, 0x94, 0xe1, 0x4a, 0xc6, 0xae,
	0xbe, 0x20, 0xbf, 0x23, 0xae, 0x39, 0xd9, 0x83, 0x02, 0x38, 0xeb, 0xca, 0x78, 0x64, 0xf5, 0xf3,
	0x7f, 0xd1, 0x37, 0x32, 0x76, 0xd5, 0x96, 0xf1, 0x08, 0xd5, 0x4d, 0x68, 0x38, 0x6f, 0x9c, 0x3d,
	0x3b, 0x67, 0x60, 0xa1, 0x93, 0x62, 0x02, 0x5f, 0xc3, 0xd2, 0xa5, 0xcc, 0xfb, 0x5a, 0xb1, 0x88,
	0xd3, 0x2a, 0x86, 0xa7, 0x00, 0x79, 0x09, 0x75, 0xcd, 0xa3, 0x9c, 0x1b, 0x5a, 0xc3, 0x90, 0xdb,
	0xb5, 0x7e, 0x78, 0x50, 0xdf, 0xc7, 0x0c, 0x90, 0xe0, 0xf6, 0x33, 0xf0, 0xe2, 0xae, 0x04, 0xe9,
	0xe9, 0x3b, 0xb0, 0x05, 0xf5, 0x04, 0x73, 0xe4, 0x6e, 0xb2, 0x36, 0xc3, 0x9f, 0x49, 0x5f, 0xe8,
	0x98, 0x44, 0x41, 0x55, 0x67, 0x46, 0xe1, 0xf9, 0x6f, 0x4f, 0xb1, 0x7d, 0xb7, 0x3a, 0x99, 0x51,
	0xff, 0xd3, 0x6f, 0x7f, 0xa8, 0x43, 0x74, 0x6a, 0x6f, 0x7e, 0x7d, 0xf7, 0xd0, 0x17, 0xec, 0x69,
	0x77, 0xed, 0xdf, 0x49, 0xad, 0x5b, 0xc7, 0x17, 0x72, 0xfb, 0x57, 0x00, 0x00, 0x00, 0xff, 0xff,
	0x6e, 0xdf, 0x90, 0x25, 0xb4, 0x05, 0x00, 0x00,
}
