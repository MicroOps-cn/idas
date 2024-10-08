// Code generated by protoc-gen-gogo. DO NOT EDIT.
// source: config.proto

package config

import (
	fmt "fmt"
	github_com_MicroOps_cn_fuck_clients_gorm "github.com/MicroOps-cn/fuck/clients/gorm"
	github_com_MicroOps_cn_fuck_clients_redis "github.com/MicroOps-cn/fuck/clients/redis"
	github_com_MicroOps_cn_fuck_clients_tracing "github.com/MicroOps-cn/fuck/clients/tracing"
	github_com_MicroOps_cn_fuck_sets "github.com/MicroOps-cn/fuck/sets"
	github_com_MicroOps_cn_fuck_wrapper "github.com/MicroOps-cn/fuck/wrapper"
	_ "github.com/MicroOps-cn/idas/pkg/client/email"
	github_com_MicroOps_cn_idas_pkg_client_email "github.com/MicroOps-cn/idas/pkg/client/email"
	_ "github.com/MicroOps-cn/idas/pkg/client/geoip"
	github_com_MicroOps_cn_idas_pkg_client_geoip "github.com/MicroOps-cn/idas/pkg/client/geoip"
	_ "github.com/MicroOps-cn/idas/pkg/client/ldap"
	github_com_MicroOps_cn_idas_pkg_client_ldap "github.com/MicroOps-cn/idas/pkg/client/ldap"
	oauth2 "github.com/MicroOps-cn/idas/pkg/client/oauth2"
	capacity "github.com/MicroOps-cn/idas/pkg/utils/capacity"
	github_com_MicroOps_cn_idas_pkg_utils_jwt "github.com/MicroOps-cn/idas/pkg/utils/jwt"
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

// @sync-to-public:public/src/services/idas/enums.ts:PasswordComplexity
type PasswordComplexity int32

const (
	PasswordComplexity_unsafe    PasswordComplexity = 0
	PasswordComplexity_general   PasswordComplexity = 1
	PasswordComplexity_safe      PasswordComplexity = 2
	PasswordComplexity_very_safe PasswordComplexity = 3
)

var PasswordComplexity_name = map[int32]string{
	0: "unsafe",
	1: "general",
	2: "safe",
	3: "very_safe",
}

var PasswordComplexity_value = map[string]int32{
	"unsafe":    0,
	"general":   1,
	"safe":      2,
	"very_safe": 3,
}

func (x PasswordComplexity) String() string {
	return proto.EnumName(PasswordComplexity_name, int32(x))
}

func (PasswordComplexity) EnumDescriptor() ([]byte, []int) {
	return fileDescriptor_3eaf2c85e69e9ea4, []int{0}
}

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

type CustomType struct {
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_unrecognized     []byte   `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *CustomType) Reset()         { *m = CustomType{} }
func (m *CustomType) String() string { return proto.CompactTextString(m) }
func (*CustomType) ProtoMessage()    {}
func (*CustomType) Descriptor() ([]byte, []int) {
	return fileDescriptor_3eaf2c85e69e9ea4, []int{1}
}
func (m *CustomType) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_CustomType.Unmarshal(m, b)
}
func (m *CustomType) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_CustomType.Marshal(b, m, deterministic)
}
func (m *CustomType) XXX_Merge(src proto.Message) {
	xxx_messageInfo_CustomType.Merge(m, src)
}
func (m *CustomType) XXX_Size() int {
	return xxx_messageInfo_CustomType.Size(m)
}
func (m *CustomType) XXX_DiscardUnknown() {
	xxx_messageInfo_CustomType.DiscardUnknown(m)
}

var xxx_messageInfo_CustomType proto.InternalMessageInfo

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
	return fileDescriptor_3eaf2c85e69e9ea4, []int{2}
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
	Mysql *github_com_MicroOps_cn_fuck_clients_gorm.MySQLClient `protobuf:"bytes,11,opt,name=mysql,proto3,oneof,customtype=github.com/MicroOps-cn/fuck/clients/gorm.MySQLClient" json:"mysql,omitempty"`
}
type Storage_Redis struct {
	Redis *github_com_MicroOps_cn_fuck_clients_redis.Client `protobuf:"bytes,12,opt,name=redis,proto3,oneof,customtype=github.com/MicroOps-cn/fuck/clients/redis.Client" json:"redis,omitempty"`
}
type Storage_Ldap struct {
	Ldap *github_com_MicroOps_cn_idas_pkg_client_ldap.Client `protobuf:"bytes,13,opt,name=ldap,proto3,oneof,customtype=github.com/MicroOps-cn/idas/pkg/client/ldap.Client" json:"ldap,omitempty"`
}
type Storage_Sqlite struct {
	Sqlite *github_com_MicroOps_cn_fuck_clients_gorm.SQLiteClient `protobuf:"bytes,14,opt,name=sqlite,proto3,oneof,customtype=github.com/MicroOps-cn/fuck/clients/gorm.SQLiteClient" json:"sqlite,omitempty"`
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
	Default              *Storage                                             `protobuf:"bytes,1,opt,name=default,proto3" json:"default,omitempty"`
	Session              *Storage                                             `protobuf:"bytes,2,opt,name=session,proto3" json:"session,omitempty"`
	User                 *Storage                                             `protobuf:"bytes,3,opt,name=user,proto3" json:"user,omitempty"`
	Logging              *Storage                                             `protobuf:"bytes,4,opt,name=logging,proto3" json:"logging,omitempty"`
	Geoip                *github_com_MicroOps_cn_idas_pkg_client_geoip.Client `protobuf:"bytes,15,opt,name=geoip,proto3,customtype=github.com/MicroOps-cn/idas/pkg/client/geoip.Client" json:"geoip,omitempty"`
	XXX_NoUnkeyedLiteral struct{}                                             `json:"-"`
	XXX_unrecognized     []byte                                               `json:"-"`
	XXX_sizecache        int32                                                `json:"-"`
}

func (m *Storages) Reset()         { *m = Storages{} }
func (m *Storages) String() string { return proto.CompactTextString(m) }
func (*Storages) ProtoMessage()    {}
func (*Storages) Descriptor() ([]byte, []int) {
	return fileDescriptor_3eaf2c85e69e9ea4, []int{3}
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

func (m *Storages) GetUser() *Storage {
	if m != nil {
		return m.User
	}
	return nil
}

func (m *Storages) GetLogging() *Storage {
	if m != nil {
		return m.Logging
	}
	return nil
}

type GlobalOptions struct {
	MaxUploadSize *capacity.Capacity `protobuf:"bytes,1,opt,name=max_upload_size,json=maxUploadSize,proto3" json:"max_upload_size,omitempty"`
	MaxBodySize   *capacity.Capacity `protobuf:"bytes,2,opt,name=max_body_size,json=maxBodySize,proto3" json:"max_body_size,omitempty"`
	UploadPath    string             `protobuf:"bytes,3,opt,name=upload_path,json=uploadPath,proto3" json:"upload_path,omitempty"`
	Workspace     string             `protobuf:"bytes,4,opt,name=workspace,proto3" json:"workspace,omitempty"`
	// Deprecated: use security.secret
	Secret *string `protobuf:"bytes,5,opt,name=secret,proto3,customtype=string" json:"secret,omitempty"`
	// Deprecated: use security.jwt_secret
	JwtSecret            *string           `protobuf:"bytes,6,opt,name=jwt_secret,json=jwtSecret,proto3,customtype=string" json:"jwt_secret,omitempty"`
	AppName              string            `protobuf:"bytes,7,opt,name=app_name,json=appName,proto3" json:"app_name,omitempty"`
	Oauth2               []*oauth2.Options `protobuf:"bytes,8,rep,name=oauth2,proto3" json:"oauth2,omitempty"`
	DisableLoginForm     bool              `protobuf:"varint,9,opt,name=disable_login_form,json=disableLoginForm,proto3" json:"disable_login_form,omitempty"`
	Title                string            `protobuf:"bytes,10,opt,name=title,proto3" json:"title,omitempty"`
	SubTitle             string            `protobuf:"bytes,11,opt,name=sub_title,json=subTitle,proto3" json:"sub_title,omitempty"`
	Logo                 string            `protobuf:"bytes,12,opt,name=logo,proto3" json:"logo,omitempty"`
	Copyright            string            `protobuf:"bytes,13,opt,name=copyright,proto3" json:"copyright,omitempty"`
	AdminEmail           string            `protobuf:"bytes,14,opt,name=admin_email,json=adminEmail,proto3" json:"admin_email,omitempty"`
	DefaultLoginType     string            `protobuf:"bytes,15,opt,name=default_login_type,json=defaultLoginType,proto3" json:"default_login_type,omitempty"`
	XXX_NoUnkeyedLiteral struct{}          `json:"-"`
	XXX_unrecognized     []byte            `json:"-"`
	XXX_sizecache        int32             `json:"-"`
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

func (m *GlobalOptions) GetAppName() string {
	if m != nil {
		return m.AppName
	}
	return ""
}

func (m *GlobalOptions) GetOauth2() []*oauth2.Options {
	if m != nil {
		return m.Oauth2
	}
	return nil
}

func (m *GlobalOptions) GetDisableLoginForm() bool {
	if m != nil {
		return m.DisableLoginForm
	}
	return false
}

func (m *GlobalOptions) GetTitle() string {
	if m != nil {
		return m.Title
	}
	return ""
}

func (m *GlobalOptions) GetSubTitle() string {
	if m != nil {
		return m.SubTitle
	}
	return ""
}

func (m *GlobalOptions) GetLogo() string {
	if m != nil {
		return m.Logo
	}
	return ""
}

func (m *GlobalOptions) GetCopyright() string {
	if m != nil {
		return m.Copyright
	}
	return ""
}

func (m *GlobalOptions) GetAdminEmail() string {
	if m != nil {
		return m.AdminEmail
	}
	return ""
}

func (m *GlobalOptions) GetDefaultLoginType() string {
	if m != nil {
		return m.DefaultLoginType
	}
	return ""
}

type RateLimit struct {
	Name                 github_com_MicroOps_cn_fuck_wrapper.OneOrMore[string] `protobuf:"bytes,1,opt,name=name,proto3,customtype=github.com/MicroOps-cn/fuck/wrapper.OneOrMore[string]" json:"name"`
	Allower              Limiter                                               `protobuf:"bytes,2,opt,name=allower,proto3,customtype=Limiter" json:"-"`
	Limit                string                                                `protobuf:"bytes,3,opt,name=limit,proto3" json:"limit,omitempty"`
	Burst                int32                                                 `protobuf:"varint,4,opt,name=burst,proto3" json:"burst,omitempty"`
	XXX_NoUnkeyedLiteral struct{}                                              `json:"-"`
	XXX_unrecognized     []byte                                                `json:"-"`
	XXX_sizecache        int32                                                 `json:"-"`
}

func (m *RateLimit) Reset()         { *m = RateLimit{} }
func (m *RateLimit) String() string { return proto.CompactTextString(m) }
func (*RateLimit) ProtoMessage()    {}
func (*RateLimit) Descriptor() ([]byte, []int) {
	return fileDescriptor_3eaf2c85e69e9ea4, []int{5}
}
func (m *RateLimit) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_RateLimit.Unmarshal(m, b)
}
func (m *RateLimit) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_RateLimit.Marshal(b, m, deterministic)
}
func (m *RateLimit) XXX_Merge(src proto.Message) {
	xxx_messageInfo_RateLimit.Merge(m, src)
}
func (m *RateLimit) XXX_Size() int {
	return xxx_messageInfo_RateLimit.Size(m)
}
func (m *RateLimit) XXX_DiscardUnknown() {
	xxx_messageInfo_RateLimit.DiscardUnknown(m)
}

var xxx_messageInfo_RateLimit proto.InternalMessageInfo

func (m *RateLimit) GetLimit() string {
	if m != nil {
		return m.Limit
	}
	return ""
}

func (m *RateLimit) GetBurst() int32 {
	if m != nil {
		return m.Burst
	}
	return 0
}

type Security struct {
	TrustIp   github_com_MicroOps_cn_fuck_sets.IPNets `protobuf:"bytes,1,opt,name=trust_ip,json=trustIp,proto3,customtype=github.com/MicroOps-cn/fuck/sets.IPNets" json:"trust_ip"`
	RateLimit []*RateLimit                            `protobuf:"bytes,2,rep,name=rate_limit,json=rateLimit,proto3" json:"rate_limit,omitempty"`
	Secret    string                                  `protobuf:"bytes,3,opt,name=secret,proto3" json:"secret,omitempty"`
	//@deprecated
	JwtSecret            string                                               `protobuf:"bytes,4,opt,name=jwt_secret,json=jwtSecret,proto3" json:"jwt_secret,omitempty"`
	Jwt                  *github_com_MicroOps_cn_idas_pkg_utils_jwt.JWTConfig `protobuf:"bytes,5,opt,name=jwt,proto3,customtype=github.com/MicroOps-cn/idas/pkg/utils/jwt.JWTConfig" json:"jwt,omitempty"`
	XXX_NoUnkeyedLiteral struct{}                                             `json:"-"`
	XXX_unrecognized     []byte                                               `json:"-"`
	XXX_sizecache        int32                                                `json:"-"`
}

func (m *Security) Reset()         { *m = Security{} }
func (m *Security) String() string { return proto.CompactTextString(m) }
func (*Security) ProtoMessage()    {}
func (*Security) Descriptor() ([]byte, []int) {
	return fileDescriptor_3eaf2c85e69e9ea4, []int{6}
}
func (m *Security) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_Security.Unmarshal(m, b)
}
func (m *Security) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_Security.Marshal(b, m, deterministic)
}
func (m *Security) XXX_Merge(src proto.Message) {
	xxx_messageInfo_Security.Merge(m, src)
}
func (m *Security) XXX_Size() int {
	return xxx_messageInfo_Security.Size(m)
}
func (m *Security) XXX_DiscardUnknown() {
	xxx_messageInfo_Security.DiscardUnknown(m)
}

var xxx_messageInfo_Security proto.InternalMessageInfo

func (m *Security) GetRateLimit() []*RateLimit {
	if m != nil {
		return m.RateLimit
	}
	return nil
}

func (m *Security) GetSecret() string {
	if m != nil {
		return m.Secret
	}
	return ""
}

func (m *Security) GetJwtSecret() string {
	if m != nil {
		return m.JwtSecret
	}
	return ""
}

type Config struct {
	Storage              *Storages                                                 `protobuf:"bytes,1,opt,name=storage,proto3" json:"storage,omitempty"`
	Global               *GlobalOptions                                            `protobuf:"bytes,2,opt,name=global,proto3" json:"global,omitempty"`
	Smtp                 *github_com_MicroOps_cn_idas_pkg_client_email.SmtpOptions `protobuf:"bytes,3,opt,name=smtp,proto3,customtype=github.com/MicroOps-cn/idas/pkg/client/email.SmtpOptions" json:"smtp,omitempty"`
	Security             *Security                                                 `protobuf:"bytes,4,opt,name=security,proto3" json:"security,omitempty"`
	Trace                github_com_MicroOps_cn_fuck_clients_tracing.TraceOptions  `protobuf:"bytes,5,opt,name=trace,proto3,customtype=github.com/MicroOps-cn/fuck/clients/tracing.TraceOptions" json:"trace"`
	XXX_NoUnkeyedLiteral struct{}                                                  `json:"-"`
	XXX_unrecognized     []byte                                                    `json:"-"`
	XXX_sizecache        int32                                                     `json:"-"`
}

func (m *Config) Reset()         { *m = Config{} }
func (m *Config) String() string { return proto.CompactTextString(m) }
func (*Config) ProtoMessage()    {}
func (*Config) Descriptor() ([]byte, []int) {
	return fileDescriptor_3eaf2c85e69e9ea4, []int{7}
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

func (m *Config) GetSecurity() *Security {
	if m != nil {
		return m.Security
	}
	return nil
}

type RuntimeSecurityConfig struct {
	ForceEnableMfa              bool               `protobuf:"varint,1,opt,name=force_enable_mfa,json=forceEnableMfa,proto3" json:"forceEnableMfa"`
	PasswordComplexity          PasswordComplexity `protobuf:"varint,2,opt,name=password_complexity,json=passwordComplexity,proto3,enum=idas.config.PasswordComplexity" json:"passwordComplexity"`
	PasswordMinLength           uint32             `protobuf:"varint,3,opt,name=password_min_length,json=passwordMinLength,proto3" json:"passwordMinLength"`
	PasswordExpireTime          uint32             `protobuf:"varint,4,opt,name=password_expire_time,json=passwordExpireTime,proto3" json:"passwordExpireTime"`
	PasswordFailedLockThreshold uint32             `protobuf:"varint,5,opt,name=password_failed_lock_threshold,json=passwordFailedLockThreshold,proto3" json:"passwordFailedLockThreshold"`
	PasswordFailedLockDuration  uint32             `protobuf:"varint,6,opt,name=password_failed_lock_duration,json=passwordFailedLockDuration,proto3" json:"passwordFailedLockDuration"`
	PasswordHistory             uint32             `protobuf:"varint,7,opt,name=password_history,json=passwordHistory,proto3" json:"passwordHistory"`
	AccountInactiveLock         uint32             `protobuf:"varint,8,opt,name=account_inactive_lock,json=accountInactiveLock,proto3" json:"accountInactiveLock"`
	LoginSessionInactivityTime  uint32             `protobuf:"varint,9,opt,name=login_session_inactivity_time,json=loginSessionInactivityTime,proto3" json:"loginSessionInactivityTime"`
	LoginSessionMaxTime         uint32             `protobuf:"varint,10,opt,name=login_session_max_time,json=loginSessionMaxTime,proto3" json:"loginSessionMaxTime"`
	XXX_NoUnkeyedLiteral        struct{}           `json:"-"`
	XXX_unrecognized            []byte             `json:"-"`
	XXX_sizecache               int32              `json:"-"`
}

func (m *RuntimeSecurityConfig) Reset()         { *m = RuntimeSecurityConfig{} }
func (m *RuntimeSecurityConfig) String() string { return proto.CompactTextString(m) }
func (*RuntimeSecurityConfig) ProtoMessage()    {}
func (*RuntimeSecurityConfig) Descriptor() ([]byte, []int) {
	return fileDescriptor_3eaf2c85e69e9ea4, []int{8}
}
func (m *RuntimeSecurityConfig) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_RuntimeSecurityConfig.Unmarshal(m, b)
}
func (m *RuntimeSecurityConfig) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_RuntimeSecurityConfig.Marshal(b, m, deterministic)
}
func (m *RuntimeSecurityConfig) XXX_Merge(src proto.Message) {
	xxx_messageInfo_RuntimeSecurityConfig.Merge(m, src)
}
func (m *RuntimeSecurityConfig) XXX_Size() int {
	return xxx_messageInfo_RuntimeSecurityConfig.Size(m)
}
func (m *RuntimeSecurityConfig) XXX_DiscardUnknown() {
	xxx_messageInfo_RuntimeSecurityConfig.DiscardUnknown(m)
}

var xxx_messageInfo_RuntimeSecurityConfig proto.InternalMessageInfo

func (m *RuntimeSecurityConfig) GetForceEnableMfa() bool {
	if m != nil {
		return m.ForceEnableMfa
	}
	return false
}

func (m *RuntimeSecurityConfig) GetPasswordComplexity() PasswordComplexity {
	if m != nil {
		return m.PasswordComplexity
	}
	return PasswordComplexity_unsafe
}

func (m *RuntimeSecurityConfig) GetPasswordMinLength() uint32 {
	if m != nil {
		return m.PasswordMinLength
	}
	return 0
}

func (m *RuntimeSecurityConfig) GetPasswordExpireTime() uint32 {
	if m != nil {
		return m.PasswordExpireTime
	}
	return 0
}

func (m *RuntimeSecurityConfig) GetPasswordFailedLockThreshold() uint32 {
	if m != nil {
		return m.PasswordFailedLockThreshold
	}
	return 0
}

func (m *RuntimeSecurityConfig) GetPasswordFailedLockDuration() uint32 {
	if m != nil {
		return m.PasswordFailedLockDuration
	}
	return 0
}

func (m *RuntimeSecurityConfig) GetPasswordHistory() uint32 {
	if m != nil {
		return m.PasswordHistory
	}
	return 0
}

func (m *RuntimeSecurityConfig) GetAccountInactiveLock() uint32 {
	if m != nil {
		return m.AccountInactiveLock
	}
	return 0
}

func (m *RuntimeSecurityConfig) GetLoginSessionInactivityTime() uint32 {
	if m != nil {
		return m.LoginSessionInactivityTime
	}
	return 0
}

func (m *RuntimeSecurityConfig) GetLoginSessionMaxTime() uint32 {
	if m != nil {
		return m.LoginSessionMaxTime
	}
	return 0
}

type RuntimeConfig struct {
	Security             *RuntimeSecurityConfig `protobuf:"bytes,1,opt,name=security,proto3" json:"security,omitempty"`
	XXX_NoUnkeyedLiteral struct{}               `json:"-"`
	XXX_unrecognized     []byte                 `json:"-"`
	XXX_sizecache        int32                  `json:"-"`
}

func (m *RuntimeConfig) Reset()         { *m = RuntimeConfig{} }
func (m *RuntimeConfig) String() string { return proto.CompactTextString(m) }
func (*RuntimeConfig) ProtoMessage()    {}
func (*RuntimeConfig) Descriptor() ([]byte, []int) {
	return fileDescriptor_3eaf2c85e69e9ea4, []int{9}
}
func (m *RuntimeConfig) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_RuntimeConfig.Unmarshal(m, b)
}
func (m *RuntimeConfig) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_RuntimeConfig.Marshal(b, m, deterministic)
}
func (m *RuntimeConfig) XXX_Merge(src proto.Message) {
	xxx_messageInfo_RuntimeConfig.Merge(m, src)
}
func (m *RuntimeConfig) XXX_Size() int {
	return xxx_messageInfo_RuntimeConfig.Size(m)
}
func (m *RuntimeConfig) XXX_DiscardUnknown() {
	xxx_messageInfo_RuntimeConfig.DiscardUnknown(m)
}

var xxx_messageInfo_RuntimeConfig proto.InternalMessageInfo

func (m *RuntimeConfig) GetSecurity() *RuntimeSecurityConfig {
	if m != nil {
		return m.Security
	}
	return nil
}

func init() {
	proto.RegisterEnum("idas.config.PasswordComplexity", PasswordComplexity_name, PasswordComplexity_value)
	proto.RegisterType((*StorageRef)(nil), "idas.config.StorageRef")
	proto.RegisterType((*CustomType)(nil), "idas.config.custom_type")
	proto.RegisterType((*Storage)(nil), "idas.config.Storage")
	proto.RegisterType((*Storages)(nil), "idas.config.Storages")
	proto.RegisterType((*GlobalOptions)(nil), "idas.config.GlobalOptions")
	proto.RegisterType((*RateLimit)(nil), "idas.config.RateLimit")
	proto.RegisterType((*Security)(nil), "idas.config.Security")
	proto.RegisterType((*Config)(nil), "idas.config.Config")
	proto.RegisterType((*RuntimeSecurityConfig)(nil), "idas.config.RuntimeSecurityConfig")
	proto.RegisterType((*RuntimeConfig)(nil), "idas.config.RuntimeConfig")
}

func init() { proto.RegisterFile("config.proto", fileDescriptor_3eaf2c85e69e9ea4) }

var fileDescriptor_3eaf2c85e69e9ea4 = []byte{
	// 1510 bytes of a gzipped FileDescriptorProto
	0x1f, 0x8b, 0x08, 0x00, 0x00, 0x00, 0x00, 0x00, 0x02, 0xff, 0x8c, 0x57, 0x41, 0x73, 0x1b, 0xb7,
	0x15, 0xb6, 0x24, 0x8a, 0x22, 0xc1, 0xd0, 0x56, 0x60, 0xd9, 0x61, 0xed, 0xda, 0xeb, 0xe1, 0xa5,
	0x4a, 0xd3, 0x2c, 0x5b, 0x39, 0x69, 0xdd, 0xa9, 0x9b, 0xc9, 0xd0, 0xb5, 0x63, 0xa5, 0x54, 0x24,
	0x83, 0x6a, 0x3b, 0xd3, 0x4e, 0x67, 0x07, 0xdc, 0x05, 0x97, 0x90, 0x76, 0x17, 0x1b, 0x00, 0x1b,
	0x8a, 0x39, 0xf7, 0xcf, 0xf4, 0xde, 0x53, 0x7f, 0x41, 0x6f, 0x3d, 0x76, 0x46, 0x87, 0xfd, 0x01,
	0xfc, 0x15, 0x1d, 0x3c, 0x80, 0xa4, 0x68, 0xd1, 0x94, 0x4e, 0x5a, 0xbc, 0xf7, 0xe1, 0x7d, 0x0f,
	0x78, 0xef, 0x7d, 0xa0, 0xd0, 0x47, 0xa1, 0xc8, 0x86, 0x3c, 0xf6, 0x73, 0x29, 0xb4, 0xc0, 0x0d,
	0x1e, 0x51, 0xe5, 0x5b, 0xd3, 0xa3, 0x5d, 0x3d, 0xc9, 0x99, 0xea, 0x24, 0x11, 0xcd, 0xad, 0xfb,
	0xd1, 0x9e, 0xb5, 0x84, 0x34, 0xa7, 0x21, 0xd7, 0x13, 0x67, 0xfd, 0xd8, 0x5a, 0x59, 0x4a, 0x79,
	0xb2, 0x6c, 0x8a, 0x99, 0xe0, 0xb3, 0xbd, 0xd8, 0x9a, 0x04, 0x2d, 0xf4, 0xe8, 0x60, 0x16, 0x2f,
	0x16, 0xb1, 0x80, 0xcf, 0x8e, 0xf9, 0xb2, 0xd6, 0x76, 0x84, 0x50, 0x5f, 0x0b, 0x49, 0x63, 0x46,
	0xd8, 0x10, 0xfb, 0x68, 0x47, 0xd9, 0x55, 0x6b, 0xe3, 0xd9, 0xc6, 0x7e, 0xe3, 0x60, 0xcf, 0xbf,
	0x92, 0xa4, 0x3f, 0x43, 0xce, 0x40, 0x18, 0xa3, 0x4a, 0x4e, 0xf5, 0xa8, 0xb5, 0xf9, 0x6c, 0x63,
	0xbf, 0x4e, 0xe0, 0xdb, 0xd8, 0x32, 0x9a, 0xb2, 0xd6, 0x96, 0xb5, 0x99, 0xef, 0x76, 0x13, 0x35,
	0xc2, 0x42, 0x69, 0x91, 0x06, 0x26, 0xb1, 0xf6, 0x3f, 0x2a, 0x68, 0xa7, 0xbf, 0x08, 0x01, 0xf0,
	0x8d, 0x05, 0x1c, 0x7f, 0x86, 0xb6, 0x24, 0x1b, 0xb6, 0x10, 0xa4, 0xf0, 0xc9, 0xca, 0x14, 0xd8,
	0xf0, 0xed, 0x1d, 0x62, 0x50, 0xf8, 0x0c, 0x6d, 0xa7, 0x13, 0xf5, 0x7d, 0xd2, 0x6a, 0x00, 0xbc,
	0xb5, 0x04, 0xbf, 0xc2, 0xda, 0x7d, 0x71, 0x59, 0x7a, 0x5f, 0xc4, 0x5c, 0x8f, 0x8a, 0x81, 0x1f,
	0x8a, 0xb4, 0x73, 0xc4, 0x43, 0x29, 0x8e, 0x73, 0xf5, 0x79, 0x98, 0x75, 0x86, 0x45, 0x78, 0xde,
	0x09, 0x13, 0xce, 0x32, 0xad, 0x3a, 0xb1, 0x90, 0xa9, 0x7f, 0x34, 0xe9, 0xbf, 0xeb, 0xbd, 0x02,
	0xcb, 0xdb, 0x3b, 0xc4, 0x52, 0xe0, 0x21, 0xda, 0x96, 0x2c, 0xe2, 0xaa, 0xf5, 0xd1, 0x0d, 0x5c,
	0x5f, 0x5c, 0x96, 0xde, 0x2f, 0x6f, 0xc3, 0x05, 0xe1, 0xfc, 0x05, 0x0f, 0xac, 0xf1, 0x39, 0xaa,
	0x98, 0x4e, 0x68, 0x35, 0x81, 0xe6, 0x89, 0xa3, 0x01, 0x90, 0x0f, 0x2d, 0xd2, 0x8b, 0x68, 0x7e,
	0x9c, 0x6b, 0x2e, 0x32, 0xd5, 0xfd, 0xf5, 0x65, 0xe9, 0x1d, 0x7c, 0x80, 0xcb, 0x6c, 0xec, 0xe4,
	0xe7, 0xb1, 0xe3, 0xb3, 0xfd, 0x35, 0x67, 0x03, 0x12, 0x9c, 0xa2, 0xaa, 0xfa, 0x3e, 0xe1, 0x9a,
	0xb5, 0xee, 0xde, 0x70, 0xaa, 0xdf, 0x5e, 0x96, 0xde, 0x97, 0xb7, 0xbe, 0xc1, 0xfe, 0xbb, 0x1e,
	0xd7, 0x6c, 0x4e, 0xe6, 0x48, 0xba, 0x35, 0x54, 0x55, 0xa2, 0x90, 0x21, 0x6b, 0xff, 0x6b, 0x13,
	0xd5, 0x5c, 0x3d, 0x95, 0x69, 0xbd, 0x88, 0x0d, 0x69, 0x91, 0xe8, 0xf5, 0xad, 0xe7, 0x40, 0xd0,
	0xaa, 0x4c, 0x29, 0x2e, 0x32, 0xe8, 0xbe, 0x0f, 0xb7, 0xaa, 0x05, 0xe1, 0x7d, 0x54, 0x29, 0x14,
	0x93, 0xd0, 0x96, 0x1f, 0x02, 0x03, 0xc2, 0x44, 0x4e, 0x44, 0x1c, 0xf3, 0x2c, 0x6e, 0x55, 0xd6,
	0x45, 0x76, 0x20, 0x9c, 0xa1, 0x6d, 0x98, 0xbd, 0xd6, 0x3d, 0x40, 0x7b, 0x4b, 0xd5, 0xb2, 0x53,
	0xf9, 0x0d, 0x13, 0x87, 0x27, 0xb3, 0x7a, 0xfd, 0xe6, 0xb2, 0xf4, 0x9e, 0xdf, 0xb2, 0x5e, 0x76,
	0xbb, 0xbd, 0x43, 0x62, 0x69, 0xda, 0xff, 0xad, 0xa0, 0xe6, 0x37, 0x89, 0x18, 0xd0, 0xc4, 0x45,
	0xc4, 0x5f, 0xa3, 0x7b, 0x29, 0xbd, 0x08, 0x8a, 0x3c, 0x11, 0x34, 0x0a, 0x14, 0xff, 0x71, 0x36,
	0xbe, 0xae, 0x94, 0xf9, 0x79, 0xec, 0x17, 0x9a, 0x27, 0xca, 0x7f, 0xe5, 0xd4, 0x84, 0x34, 0x53,
	0x7a, 0xf1, 0x27, 0xc0, 0xf7, 0xf9, 0x8f, 0x0c, 0xbf, 0x44, 0xc6, 0x10, 0x0c, 0x44, 0x34, 0xb1,
	0xfb, 0x37, 0x6f, 0xd8, 0xdf, 0x48, 0xe9, 0x45, 0x57, 0x44, 0x13, 0xd8, 0xed, 0xa1, 0x86, 0xe3,
	0x06, 0x35, 0xb0, 0x93, 0x8f, 0xac, 0xe9, 0xc4, 0x68, 0xc2, 0x4f, 0x51, 0x7d, 0x2c, 0xe4, 0xb9,
	0xca, 0x69, 0xc8, 0xe0, 0x52, 0xeb, 0x64, 0x61, 0xc0, 0x6d, 0x54, 0x55, 0x2c, 0x94, 0x4c, 0xb7,
	0xb6, 0x8d, 0xab, 0x8b, 0x2e, 0x4b, 0xaf, 0xaa, 0xb4, 0xe4, 0x59, 0x4c, 0x9c, 0x07, 0x7f, 0x8a,
	0xd0, 0xd9, 0x58, 0x07, 0x0e, 0x57, 0xbd, 0x86, 0xab, 0x9f, 0x8d, 0x75, 0xdf, 0x42, 0x7f, 0x82,
	0x6a, 0x34, 0xcf, 0x03, 0x50, 0x95, 0x1d, 0xe0, 0xda, 0xa1, 0x79, 0xfe, 0x9d, 0x11, 0x96, 0xe7,
	0xa8, 0x6a, 0x35, 0xb1, 0x55, 0x7b, 0xb6, 0xb5, 0xdf, 0x38, 0x78, 0xbc, 0x54, 0x2b, 0x27, 0x97,
	0xee, 0x56, 0x89, 0x83, 0xe2, 0x5f, 0x20, 0x1c, 0x71, 0x45, 0x07, 0x09, 0x0b, 0x12, 0x11, 0xf3,
	0x2c, 0x18, 0x0a, 0x99, 0xb6, 0xea, 0xcf, 0x36, 0xf6, 0x6b, 0x64, 0xd7, 0x79, 0x7a, 0xc6, 0xf1,
	0x46, 0xc8, 0x14, 0xef, 0xa1, 0x6d, 0xcd, 0x75, 0xc2, 0x40, 0xbd, 0xea, 0xc4, 0x2e, 0xf0, 0x63,
	0x54, 0x57, 0xc5, 0x20, 0xb0, 0x9e, 0x06, 0x78, 0x6a, 0xaa, 0x18, 0x9c, 0x82, 0x13, 0xa3, 0x4a,
	0x22, 0x62, 0x01, 0xa2, 0x52, 0x27, 0xf0, 0x6d, 0x6e, 0x2c, 0x14, 0xf9, 0x44, 0xf2, 0x78, 0xa4,
	0x41, 0x06, 0xea, 0x64, 0x61, 0x30, 0x17, 0x4e, 0xa3, 0x94, 0x67, 0x01, 0xbc, 0x03, 0x30, 0xb7,
	0x75, 0x82, 0xc0, 0xf4, 0xda, 0x58, 0x20, 0x67, 0x3b, 0x28, 0x2e, 0x67, 0x33, 0xbf, 0xd0, 0xa0,
	0x75, 0xb2, 0xeb, 0x3c, 0x90, 0xf3, 0xa9, 0xd1, 0xe3, 0x7f, 0x6f, 0xa0, 0x3a, 0xa1, 0x9a, 0xf5,
	0x78, 0xca, 0x35, 0x7e, 0x77, 0x55, 0x91, 0xbb, 0xbf, 0xff, 0x4f, 0xe9, 0xdd, 0xb9, 0x69, 0xee,
	0xc7, 0x92, 0xe6, 0x39, 0x93, 0xfe, 0x71, 0xc6, 0x8e, 0xe5, 0x91, 0x90, 0xec, 0x6f, 0xb6, 0x2c,
	0x7f, 0x77, 0x82, 0xee, 0xa3, 0x1d, 0x9a, 0x24, 0x62, 0xcc, 0xa4, 0x7d, 0x2a, 0xba, 0x7b, 0x2e,
	0xea, 0x0e, 0x50, 0x32, 0x39, 0x2d, 0xbd, 0x8d, 0xcf, 0xc9, 0x0c, 0x64, 0x2e, 0x31, 0x31, 0x0e,
	0xd7, 0x4a, 0x76, 0x61, 0xac, 0x83, 0x42, 0x2a, 0x0d, 0x1d, 0xb4, 0x4d, 0xec, 0xa2, 0xfd, 0x4f,
	0xa3, 0x22, 0x2c, 0x2c, 0x24, 0xd7, 0x13, 0xfc, 0x2d, 0xaa, 0x69, 0x59, 0x28, 0x1d, 0xf0, 0xdc,
	0xe5, 0xdf, 0x71, 0x4c, 0x3f, 0x5b, 0x97, 0xbf, 0x62, 0x5a, 0xf9, 0x87, 0x27, 0xdf, 0x31, 0xad,
	0xc8, 0x0e, 0x04, 0x38, 0xcc, 0xf1, 0x97, 0x08, 0x49, 0xaa, 0x59, 0x60, 0x33, 0xd9, 0x84, 0x86,
	0x79, 0xb8, 0x24, 0x05, 0xf3, 0x3b, 0x23, 0x75, 0x39, 0xbf, 0xbe, 0x87, 0xf3, 0x6e, 0xb6, 0xc9,
	0xcf, 0x3a, 0xf8, 0xc9, 0x52, 0x07, 0xbb, 0x21, 0x58, 0x74, 0x6d, 0x88, 0xb6, 0xce, 0xc6, 0x76,
	0x02, 0xd6, 0x49, 0xf0, 0xad, 0xc4, 0x03, 0x66, 0xb5, 0x73, 0x36, 0xd6, 0xfe, 0xb7, 0x7f, 0x39,
	0x7d, 0x05, 0x41, 0x88, 0x89, 0xde, 0x9e, 0x6e, 0xa2, 0xaa, 0x5d, 0xe3, 0xce, 0xfb, 0x4f, 0xfd,
	0x83, 0x55, 0x2a, 0xa7, 0x16, 0x6f, 0xfd, 0x01, 0xaa, 0xc6, 0xa0, 0x3a, 0x4e, 0x1b, 0x1e, 0x2d,
	0xe1, 0x97, 0x04, 0x89, 0x38, 0x24, 0xce, 0x51, 0x45, 0xa5, 0x3a, 0x77, 0xa2, 0xfb, 0x74, 0x69,
	0xda, 0xec, 0x4f, 0x98, 0x7e, 0xaa, 0xe7, 0x0f, 0xd9, 0xcb, 0xcb, 0xd2, 0x7b, 0x71, 0x4b, 0x61,
	0xbc, 0xb6, 0x9b, 0x00, 0x13, 0xfe, 0x15, 0xaa, 0x29, 0xd7, 0x0c, 0x4e, 0xbd, 0xdf, 0x3b, 0x97,
	0x73, 0x92, 0x39, 0x0c, 0xff, 0x19, 0x6d, 0x6b, 0x69, 0x84, 0xc9, 0xaa, 0xcf, 0xd7, 0xae, 0x61,
	0x5e, 0xdc, 0xe6, 0xa1, 0x33, 0x1b, 0x79, 0x16, 0xfb, 0xa7, 0x26, 0xc0, 0x2c, 0x13, 0x1b, 0xae,
	0xfd, 0xbf, 0x2a, 0x7a, 0x40, 0x8a, 0x4c, 0xf3, 0x94, 0xcd, 0x58, 0xdd, 0xdd, 0xbf, 0x44, 0xbb,
	0x43, 0x21, 0x43, 0x16, 0xb0, 0x0c, 0x64, 0x25, 0x1d, 0x52, 0x28, 0x42, 0xad, 0x8b, 0xa7, 0xa5,
	0x77, 0x17, 0x7c, 0xaf, 0xc1, 0x75, 0x34, 0xa4, 0xe4, 0xbd, 0x35, 0x1e, 0xa1, 0xfb, 0x39, 0x55,
	0x6a, 0x2c, 0x64, 0x14, 0x84, 0x22, 0xcd, 0x13, 0x76, 0x61, 0x4e, 0x6b, 0xaa, 0x72, 0x77, 0xfe,
	0xfa, 0xd8, 0xd3, 0x9e, 0x38, 0xdc, 0xab, 0x39, 0xac, 0xfb, 0x70, 0x5a, 0x7a, 0x38, 0xbf, 0x66,
	0x27, 0x2b, 0x6c, 0xf8, 0xf5, 0x15, 0x26, 0xa3, 0x36, 0x09, 0xcb, 0x62, 0xa7, 0xef, 0xcd, 0xee,
	0x83, 0x69, 0xe9, 0x7d, 0x3c, 0x73, 0x1f, 0xf1, 0xac, 0x07, 0x4e, 0x72, 0xdd, 0x84, 0xdf, 0xa2,
	0xbd, 0x79, 0x18, 0x76, 0x91, 0x73, 0xc9, 0x02, 0x73, 0x29, 0x50, 0x9f, 0xe6, 0x72, 0x42, 0xaf,
	0xc1, 0x7d, 0xca, 0x53, 0x46, 0x56, 0xd8, 0x70, 0x84, 0x9e, 0xce, 0x23, 0x0d, 0x29, 0x4f, 0x58,
	0x14, 0x24, 0x22, 0x3c, 0x0f, 0xf4, 0x48, 0x32, 0x35, 0x12, 0x49, 0x04, 0x35, 0x6c, 0x76, 0xbd,
	0x69, 0xe9, 0x3d, 0x9e, 0x21, 0xdf, 0x00, 0xb0, 0x27, 0xc2, 0xf3, 0xd3, 0x19, 0x8c, 0xac, 0x73,
	0x62, 0x8a, 0x9e, 0xac, 0x64, 0x89, 0x0a, 0x49, 0x4d, 0x85, 0xe1, 0xf9, 0x69, 0x76, 0x9f, 0x4e,
	0x4b, 0xef, 0xd1, 0xf5, 0x38, 0x7f, 0x70, 0x28, 0xb2, 0xc6, 0x87, 0xbf, 0x42, 0xbb, 0x73, 0x8a,
	0x11, 0x37, 0x23, 0x36, 0x81, 0xb7, 0xaa, 0xd9, 0xbd, 0x3f, 0x2d, 0xbd, 0x7b, 0x33, 0xdf, 0x5b,
	0xeb, 0x22, 0xef, 0x1b, 0xf0, 0x1f, 0xd1, 0x03, 0x1a, 0x86, 0xa2, 0xc8, 0x74, 0xc0, 0x33, 0x1a,
	0x6a, 0xfe, 0x03, 0x83, 0x1c, 0x5b, 0x35, 0x08, 0xf2, 0xc9, 0xb4, 0xf4, 0xee, 0x3b, 0xc0, 0xa1,
	0xf3, 0x1b, 0x7e, 0xb2, 0xca, 0x68, 0xce, 0x6b, 0x1f, 0x09, 0xf7, 0x5b, 0x69, 0x16, 0x92, 0xeb,
	0x89, 0x2d, 0x54, 0x7d, 0x71, 0x5e, 0x00, 0xf6, 0x2d, 0xee, 0x70, 0x0e, 0x83, 0x82, 0xad, 0xf1,
	0xe1, 0x1e, 0x7a, 0xb8, 0x4c, 0x61, 0x7e, 0x6d, 0x40, 0x6c, 0xb4, 0x48, 0xf8, 0xea, 0xfe, 0x23,
	0x7a, 0x01, 0x41, 0x57, 0x19, 0xdb, 0xc7, 0xa8, 0xe9, 0x06, 0xcb, 0x0d, 0xd4, 0x57, 0x57, 0xa6,
	0xde, 0xaa, 0x59, 0x7b, 0x59, 0xa8, 0x57, 0x8d, 0xe1, 0x42, 0x02, 0x7e, 0xfe, 0x06, 0xe1, 0xeb,
	0xa3, 0x82, 0x11, 0xaa, 0x16, 0x99, 0xa2, 0x43, 0xb6, 0x7b, 0x07, 0x37, 0xd0, 0x4e, 0xcc, 0x32,
	0x26, 0x69, 0xb2, 0xbb, 0x81, 0x6b, 0xa8, 0x02, 0xe6, 0x4d, 0xdc, 0x44, 0xf5, 0x1f, 0x98, 0x9c,
	0x04, 0xb0, 0xdc, 0xea, 0x7e, 0xf6, 0xd7, 0x4f, 0xd7, 0xe9, 0x97, 0xcd, 0xe6, 0x77, 0xf6, 0xcf,
	0xc9, 0xf6, 0xa0, 0x0a, 0xff, 0x83, 0x3d, 0xff, 0x7f, 0x00, 0x00, 0x00, 0xff, 0xff, 0x83, 0x02,
	0x20, 0x00, 0x18, 0x0e, 0x00, 0x00,
}
