// Code generated by protoc-gen-gogo. DO NOT EDIT.
// source: config.proto

package config

import (
	fmt "fmt"
	proto "github.com/gogo/protobuf/proto"
	math "math"
	email "pkg/client/email"
	gorm "pkg/client/gorm"
	ldap "pkg/client/ldap"
	redis "pkg/client/redis"
	capacity "pkg/utils/capacity"
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
	Storage              *Storage `protobuf:"bytes,1,opt,name=Storage,proto3" json:"Storage,omitempty"`
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
	Mysql *gorm.MySQLOptions `protobuf:"bytes,11,opt,name=mysql,proto3,oneof" json:"mysql,omitempty"`
}
type Storage_Redis struct {
	Redis *redis.RedisOptions `protobuf:"bytes,12,opt,name=redis,proto3,oneof" json:"redis,omitempty"`
}
type Storage_Ldap struct {
	Ldap *ldap.LdapOptions `protobuf:"bytes,13,opt,name=ldap,proto3,oneof" json:"ldap,omitempty"`
}
type Storage_Sqlite struct {
	Sqlite *gorm.SQLiteOptions `protobuf:"bytes,14,opt,name=sqlite,proto3,oneof" json:"sqlite,omitempty"`
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

func (m *Storage) GetMysql() *gorm.MySQLOptions {
	if x, ok := m.GetSource().(*Storage_Mysql); ok {
		return x.Mysql
	}
	return nil
}

func (m *Storage) GetRedis() *redis.RedisOptions {
	if x, ok := m.GetSource().(*Storage_Redis); ok {
		return x.Redis
	}
	return nil
}

func (m *Storage) GetLdap() *ldap.LdapOptions {
	if x, ok := m.GetSource().(*Storage_Ldap); ok {
		return x.Ldap
	}
	return nil
}

func (m *Storage) GetSqlite() *gorm.SQLiteOptions {
	if x, ok := m.GetSource().(*Storage_Sqlite); ok {
		return x.Sqlite
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

type Capacity struct {
	Capacity             uint64   `protobuf:"varint,1,opt,name=capacity,proto3" json:"capacity,omitempty"`
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_unrecognized     []byte   `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *Capacity) Reset()         { *m = Capacity{} }
func (m *Capacity) String() string { return proto.CompactTextString(m) }
func (*Capacity) ProtoMessage()    {}
func (*Capacity) Descriptor() ([]byte, []int) {
	return fileDescriptor_3eaf2c85e69e9ea4, []int{3}
}
func (m *Capacity) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_Capacity.Unmarshal(m, b)
}
func (m *Capacity) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_Capacity.Marshal(b, m, deterministic)
}
func (m *Capacity) XXX_Merge(src proto.Message) {
	xxx_messageInfo_Capacity.Merge(m, src)
}
func (m *Capacity) XXX_Size() int {
	return xxx_messageInfo_Capacity.Size(m)
}
func (m *Capacity) XXX_DiscardUnknown() {
	xxx_messageInfo_Capacity.DiscardUnknown(m)
}

var xxx_messageInfo_Capacity proto.InternalMessageInfo

func (m *Capacity) GetCapacity() uint64 {
	if m != nil {
		return m.Capacity
	}
	return 0
}

type GlobalOptions struct {
	MaxUploadSize        *capacity.Capacity `protobuf:"bytes,1,opt,name=max_upload_size,json=maxUploadSize,proto3" json:"max_upload_size,omitempty"`
	MaxBodySize          *capacity.Capacity `protobuf:"bytes,2,opt,name=max_body_size,json=maxBodySize,proto3" json:"max_body_size,omitempty"`
	UploadPath           string             `protobuf:"bytes,3,opt,name=upload_path,json=uploadPath,proto3" json:"upload_path,omitempty"`
	Workspace            string             `protobuf:"bytes,4,opt,name=workspace,proto3" json:"workspace,omitempty"`
	XXX_NoUnkeyedLiteral struct{}           `json:"-"`
	XXX_unrecognized     []byte             `json:"-"`
	XXX_sizecache        int32              `json:"-"`
}

func (m *GlobalOptions) Reset()         { *m = GlobalOptions{} }
func (m *GlobalOptions) String() string { return proto.CompactTextString(m) }
func (*GlobalOptions) ProtoMessage()    {}
func (*GlobalOptions) Descriptor() ([]byte, []int) {
	return fileDescriptor_3eaf2c85e69e9ea4, []int{4}
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

type Config struct {
	Storage              *Storages          `protobuf:"bytes,1,opt,name=storage,proto3" json:"storage,omitempty"`
	Global               *GlobalOptions     `protobuf:"bytes,2,opt,name=global,proto3" json:"global,omitempty"`
	Smtp                 *email.SmtpOptions `protobuf:"bytes,3,opt,name=smtp,proto3" json:"smtp,omitempty"`
	XXX_NoUnkeyedLiteral struct{}           `json:"-"`
	XXX_unrecognized     []byte             `json:"-"`
	XXX_sizecache        int32              `json:"-"`
}

func (m *Config) Reset()         { *m = Config{} }
func (m *Config) String() string { return proto.CompactTextString(m) }
func (*Config) ProtoMessage()    {}
func (*Config) Descriptor() ([]byte, []int) {
	return fileDescriptor_3eaf2c85e69e9ea4, []int{5}
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

func (m *Config) GetSmtp() *email.SmtpOptions {
	if m != nil {
		return m.Smtp
	}
	return nil
}

func init() {
	proto.RegisterType((*StorageRef)(nil), "idas.config.StorageRef")
	proto.RegisterType((*Storage)(nil), "idas.config.Storage")
	proto.RegisterType((*Storages)(nil), "idas.config.Storages")
	proto.RegisterType((*Capacity)(nil), "idas.config.Capacity")
	proto.RegisterType((*GlobalOptions)(nil), "idas.config.GlobalOptions")
	proto.RegisterType((*Config)(nil), "idas.config.Config")
}

func init() { proto.RegisterFile("config.proto", fileDescriptor_3eaf2c85e69e9ea4) }

var fileDescriptor_3eaf2c85e69e9ea4 = []byte{
	// 545 bytes of a gzipped FileDescriptorProto
	0x1f, 0x8b, 0x08, 0x00, 0x00, 0x00, 0x00, 0x00, 0x02, 0xff, 0x84, 0x93, 0xd1, 0x6e, 0xd3, 0x30,
	0x14, 0x86, 0x97, 0x36, 0xeb, 0xba, 0xd3, 0x95, 0x81, 0x35, 0x84, 0x55, 0x01, 0xab, 0x72, 0x81,
	0x2a, 0x21, 0x65, 0x52, 0x27, 0x81, 0x10, 0x5c, 0xa0, 0xee, 0x82, 0x5d, 0x14, 0xc1, 0x52, 0x71,
	0xc3, 0x4d, 0xe5, 0x36, 0x6e, 0xb1, 0x96, 0xd4, 0x5e, 0xec, 0x88, 0x66, 0xcf, 0xc1, 0x0b, 0xf0,
	0x10, 0xbc, 0x06, 0xcf, 0x84, 0x72, 0xec, 0xd0, 0x06, 0x75, 0xec, 0xa6, 0x4d, 0xfe, 0xf3, 0xfd,
	0x3d, 0xa7, 0xbf, 0x8f, 0xe1, 0x68, 0x2e, 0x57, 0x0b, 0xb1, 0x0c, 0x55, 0x26, 0x8d, 0x24, 0x1d,
	0x11, 0x33, 0x1d, 0x5a, 0xa9, 0xf7, 0xd0, 0x14, 0x8a, 0xeb, 0xb3, 0xa5, 0xcc, 0x52, 0x5b, 0xee,
	0x3d, 0xb2, 0x4a, 0xc6, 0x63, 0xa1, 0x9d, 0xe4, 0xa0, 0x24, 0x66, 0xca, 0x29, 0x27, 0x56, 0x99,
	0x33, 0xc5, 0xe6, 0xc2, 0x14, 0x75, 0x2b, 0x4f, 0x99, 0x48, 0xac, 0x14, 0xc4, 0x00, 0x13, 0x23,
	0x33, 0xb6, 0xe4, 0x11, 0x5f, 0x90, 0x10, 0x0e, 0xdc, 0x1b, 0xf5, 0xfa, 0xde, 0xa0, 0x33, 0x3c,
	0x09, 0xb7, 0x86, 0x09, 0x2b, 0xb2, 0x82, 0x08, 0x01, 0x5f, 0x31, 0xf3, 0x8d, 0x36, 0xfa, 0xde,
	0xe0, 0x30, 0xc2, 0xe7, 0x52, 0x5b, 0xb1, 0x94, 0xd3, 0xa6, 0xd5, 0xca, 0xe7, 0xe0, 0x57, 0x03,
	0xb6, 0x3d, 0x58, 0xf7, 0x36, 0x75, 0xf2, 0x12, 0x9a, 0x19, 0x5f, 0x50, 0xc0, 0x9e, 0x4f, 0x76,
	0xf6, 0xe4, 0x8b, 0xcb, 0xbd, 0xa8, 0xa4, 0xc8, 0x2b, 0xd8, 0x4f, 0x0b, 0x7d, 0x93, 0xd0, 0x0e,
	0xe2, 0xcf, 0x1d, 0x9e, 0x08, 0xbe, 0x32, 0x21, 0x06, 0xf5, 0xb1, 0x98, 0x5c, 0x8d, 0x3f, 0x29,
	0x23, 0xe4, 0x4a, 0x5f, 0xee, 0x45, 0x16, 0x27, 0xaf, 0x61, 0x1f, 0x43, 0xa3, 0x47, 0xe8, 0x3b,
	0xad, 0xf9, 0x6c, 0x9c, 0x51, 0xf9, 0xb9, 0x65, 0x44, 0x95, 0x9c, 0x83, 0x5f, 0x46, 0x4b, 0xbb,
	0xe8, 0x7b, 0x56, 0xf3, 0x61, 0xe6, 0xe3, 0x98, 0xa9, 0x8d, 0x0b, 0x61, 0xf2, 0x06, 0x5a, 0xfa,
	0x26, 0x11, 0x86, 0xd3, 0x07, 0x3b, 0xda, 0xe1, 0x98, 0x93, 0xab, 0xb1, 0x30, 0x7c, 0x63, 0x74,
	0x86, 0x51, 0x1b, 0x5a, 0x5a, 0xe6, 0xd9, 0x9c, 0x07, 0x3f, 0x3c, 0x68, 0xbb, 0x00, 0x74, 0x79,
	0x38, 0x31, 0x5f, 0xb0, 0x3c, 0x31, 0xff, 0x3f, 0x1c, 0x07, 0x95, 0xbc, 0xe6, 0x5a, 0x0b, 0xb9,
	0xc2, 0xf3, 0xb9, 0x93, 0x77, 0x10, 0x19, 0x80, 0x9f, 0x6b, 0x9e, 0xd1, 0x66, 0xbf, 0x79, 0x27,
	0x8c, 0x44, 0xf0, 0x02, 0xda, 0x17, 0x6e, 0xb3, 0x48, 0x0f, 0xda, 0xd5, 0x96, 0xe1, 0x58, 0x7e,
	0xf4, 0xf7, 0x3d, 0xf8, 0xed, 0x41, 0xf7, 0x43, 0x22, 0x67, 0x2c, 0x71, 0x7f, 0x92, 0xbc, 0x87,
	0xe3, 0x94, 0xad, 0xa7, 0xb9, 0x4a, 0x24, 0x8b, 0xa7, 0x5a, 0xdc, 0x56, 0x8b, 0x46, 0x6d, 0x3b,
	0x75, 0xbd, 0x0c, 0x73, 0x23, 0x12, 0x1d, 0x56, 0x0d, 0xa2, 0x6e, 0xca, 0xd6, 0x5f, 0x90, 0x9f,
	0x88, 0x5b, 0x4e, 0xde, 0x41, 0x29, 0x4c, 0x67, 0x32, 0x2e, 0xac, 0xbf, 0x71, 0x8f, 0xbf, 0x93,
	0xb2, 0xf5, 0x48, 0xc6, 0x05, 0xba, 0x4f, 0xa1, 0xe3, 0x7a, 0xe3, 0xde, 0xda, 0x1d, 0x05, 0x2b,
	0x7d, 0x2e, 0xb7, 0xf7, 0x29, 0x1c, 0x7e, 0x97, 0xd9, 0xb5, 0x56, 0x6c, 0xce, 0xa9, 0x8f, 0xe5,
	0x8d, 0x10, 0xfc, 0xf4, 0xa0, 0x75, 0x81, 0x89, 0x90, 0x33, 0x38, 0xd0, 0xb5, 0xab, 0xf2, 0x78,
	0x57, 0x60, 0x3a, 0xaa, 0x28, 0x32, 0x84, 0xd6, 0x12, 0xb3, 0x70, 0x13, 0xf7, 0x6a, 0x7c, 0x2d,
	0xa6, 0xc8, 0x91, 0x64, 0x08, 0xbe, 0x4e, 0x8d, 0xc2, 0x39, 0xff, 0xdd, 0x74, 0x7b, 0x8b, 0x27,
	0xa9, 0xa9, 0x56, 0x2f, 0x42, 0x76, 0x74, 0xfc, 0xb5, 0x6b, 0x7f, 0xf3, 0xad, 0xfd, 0x9a, 0xb5,
	0xf0, 0xa6, 0x9f, 0xff, 0x09, 0x00, 0x00, 0xff, 0xff, 0xec, 0xd0, 0x5e, 0xf6, 0x66, 0x04, 0x00,
	0x00,
}
