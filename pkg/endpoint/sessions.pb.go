// Code generated by protoc-gen-gogo. DO NOT EDIT.
// source: endpoints/sessions.proto

package endpoint

import (
	fmt "fmt"
	github_com_MicroOps_cn_idas_pkg_service_models "github.com/MicroOps-cn/idas/pkg/service/models"
	models "github.com/MicroOps-cn/idas/pkg/service/models"
	github_com_MicroOps_cn_idas_pkg_utils_sign "github.com/MicroOps-cn/idas/pkg/utils/sign"
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

// @sync-to-public:public/src/services/idas/enums.ts:LoginType
type LoginType int32

const (
	LoginType_normal           LoginType = 0
	LoginType_mfa_totp         LoginType = 1
	LoginType_mfa_email        LoginType = 2
	LoginType_mfa_sms          LoginType = 3
	LoginType_email            LoginType = 4
	LoginType_sms              LoginType = 5
	LoginType_enable_mfa_totp  LoginType = 10
	LoginType_enable_mfa_email LoginType = 11
	LoginType_enable_mfa_sms   LoginType = 12
)

var LoginType_name = map[int32]string{
	0:  "normal",
	1:  "mfa_totp",
	2:  "mfa_email",
	3:  "mfa_sms",
	4:  "email",
	5:  "sms",
	10: "enable_mfa_totp",
	11: "enable_mfa_email",
	12: "enable_mfa_sms",
}

var LoginType_value = map[string]int32{
	"normal":           0,
	"mfa_totp":         1,
	"mfa_email":        2,
	"mfa_sms":          3,
	"email":            4,
	"sms":              5,
	"enable_mfa_totp":  10,
	"enable_mfa_email": 11,
	"enable_mfa_sms":   12,
}

func (x LoginType) String() string {
	return proto.EnumName(LoginType_name, int32(x))
}

func (LoginType) EnumDescriptor() ([]byte, []int) {
	return fileDescriptor_56667f8afe262146, []int{0}
}

type GetSessionsRequest struct {
	BaseListRequest      `protobuf:"bytes,1,opt,name=base_list_request,json=baseListRequest,proto3,embedded=base_list_request" json:"base_list_request"`
	UserId               string   `protobuf:"bytes,2,opt,name=user_id,json=userId,proto3" json:"userId,omitempty" valid:"required,uuid"`
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_unrecognized     []byte   `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *GetSessionsRequest) Reset()         { *m = GetSessionsRequest{} }
func (m *GetSessionsRequest) String() string { return proto.CompactTextString(m) }
func (*GetSessionsRequest) ProtoMessage()    {}
func (*GetSessionsRequest) Descriptor() ([]byte, []int) {
	return fileDescriptor_56667f8afe262146, []int{0}
}
func (m *GetSessionsRequest) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_GetSessionsRequest.Unmarshal(m, b)
}
func (m *GetSessionsRequest) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_GetSessionsRequest.Marshal(b, m, deterministic)
}
func (m *GetSessionsRequest) XXX_Merge(src proto.Message) {
	xxx_messageInfo_GetSessionsRequest.Merge(m, src)
}
func (m *GetSessionsRequest) XXX_Size() int {
	return xxx_messageInfo_GetSessionsRequest.Size(m)
}
func (m *GetSessionsRequest) XXX_DiscardUnknown() {
	xxx_messageInfo_GetSessionsRequest.DiscardUnknown(m)
}

var xxx_messageInfo_GetSessionsRequest proto.InternalMessageInfo

func (m *GetSessionsRequest) GetUserId() string {
	if m != nil {
		return m.UserId
	}
	return ""
}

type SessionInfo struct {
	Id                   string   `protobuf:"bytes,1,opt,name=id,proto3" json:"id" valid:"required,uuid"`
	CreateTime           string   `protobuf:"bytes,2,opt,name=create_time,json=createTime,proto3" json:"createTime" valid:"required"`
	Expiry               string   `protobuf:"bytes,3,opt,name=expiry,proto3" json:"expiry" valid:"required"`
	LastSeen             string   `protobuf:"bytes,4,opt,name=last_seen,json=lastSeen,proto3" json:"lastSeen,omitempty"`
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_unrecognized     []byte   `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *SessionInfo) Reset()         { *m = SessionInfo{} }
func (m *SessionInfo) String() string { return proto.CompactTextString(m) }
func (*SessionInfo) ProtoMessage()    {}
func (*SessionInfo) Descriptor() ([]byte, []int) {
	return fileDescriptor_56667f8afe262146, []int{1}
}
func (m *SessionInfo) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_SessionInfo.Unmarshal(m, b)
}
func (m *SessionInfo) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_SessionInfo.Marshal(b, m, deterministic)
}
func (m *SessionInfo) XXX_Merge(src proto.Message) {
	xxx_messageInfo_SessionInfo.Merge(m, src)
}
func (m *SessionInfo) XXX_Size() int {
	return xxx_messageInfo_SessionInfo.Size(m)
}
func (m *SessionInfo) XXX_DiscardUnknown() {
	xxx_messageInfo_SessionInfo.DiscardUnknown(m)
}

var xxx_messageInfo_SessionInfo proto.InternalMessageInfo

func (m *SessionInfo) GetId() string {
	if m != nil {
		return m.Id
	}
	return ""
}

func (m *SessionInfo) GetCreateTime() string {
	if m != nil {
		return m.CreateTime
	}
	return ""
}

func (m *SessionInfo) GetExpiry() string {
	if m != nil {
		return m.Expiry
	}
	return ""
}

func (m *SessionInfo) GetLastSeen() string {
	if m != nil {
		return m.LastSeen
	}
	return ""
}

type GetSessionsResponse struct {
	BaseListResponse     `protobuf:"bytes,1,opt,name=base_list_response,json=baseListResponse,proto3,embedded=base_list_response" json:",omitempty"`
	Data                 []*SessionInfo `protobuf:"bytes,2,rep,name=data,proto3" json:"data,omitempty"`
	XXX_NoUnkeyedLiteral struct{}       `json:"-"`
	XXX_unrecognized     []byte         `json:"-"`
	XXX_sizecache        int32          `json:"-"`
}

func (m *GetSessionsResponse) Reset()         { *m = GetSessionsResponse{} }
func (m *GetSessionsResponse) String() string { return proto.CompactTextString(m) }
func (*GetSessionsResponse) ProtoMessage()    {}
func (*GetSessionsResponse) Descriptor() ([]byte, []int) {
	return fileDescriptor_56667f8afe262146, []int{2}
}
func (m *GetSessionsResponse) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_GetSessionsResponse.Unmarshal(m, b)
}
func (m *GetSessionsResponse) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_GetSessionsResponse.Marshal(b, m, deterministic)
}
func (m *GetSessionsResponse) XXX_Merge(src proto.Message) {
	xxx_messageInfo_GetSessionsResponse.Merge(m, src)
}
func (m *GetSessionsResponse) XXX_Size() int {
	return xxx_messageInfo_GetSessionsResponse.Size(m)
}
func (m *GetSessionsResponse) XXX_DiscardUnknown() {
	xxx_messageInfo_GetSessionsResponse.DiscardUnknown(m)
}

var xxx_messageInfo_GetSessionsResponse proto.InternalMessageInfo

func (m *GetSessionsResponse) GetData() []*SessionInfo {
	if m != nil {
		return m.Data
	}
	return nil
}

type DeleteSessionRequest struct {
	Id                   string   `protobuf:"bytes,1,opt,name=id,proto3" json:"id" valid:"required,uuid"`
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_unrecognized     []byte   `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *DeleteSessionRequest) Reset()         { *m = DeleteSessionRequest{} }
func (m *DeleteSessionRequest) String() string { return proto.CompactTextString(m) }
func (*DeleteSessionRequest) ProtoMessage()    {}
func (*DeleteSessionRequest) Descriptor() ([]byte, []int) {
	return fileDescriptor_56667f8afe262146, []int{3}
}
func (m *DeleteSessionRequest) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_DeleteSessionRequest.Unmarshal(m, b)
}
func (m *DeleteSessionRequest) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_DeleteSessionRequest.Marshal(b, m, deterministic)
}
func (m *DeleteSessionRequest) XXX_Merge(src proto.Message) {
	xxx_messageInfo_DeleteSessionRequest.Merge(m, src)
}
func (m *DeleteSessionRequest) XXX_Size() int {
	return xxx_messageInfo_DeleteSessionRequest.Size(m)
}
func (m *DeleteSessionRequest) XXX_DiscardUnknown() {
	xxx_messageInfo_DeleteSessionRequest.DiscardUnknown(m)
}

var xxx_messageInfo_DeleteSessionRequest proto.InternalMessageInfo

func (m *DeleteSessionRequest) GetId() string {
	if m != nil {
		return m.Id
	}
	return ""
}

type AuthenticationRequest struct {
	Username             string                                                   `protobuf:"bytes,1,opt,name=username,proto3" json:"username,omitempty"`
	Password             string                                                   `protobuf:"bytes,2,opt,name=password,proto3" json:"password,omitempty"`
	AuthMethod           models.AuthMeta_Method                                   `protobuf:"varint,3,opt,name=auth_method,json=authMethod,proto3,enum=idas.service.models.AuthMeta_Method" json:"authMethod,omitempty"`
	AuthAlgorithm        github_com_MicroOps_cn_idas_pkg_utils_sign.AuthAlgorithm `protobuf:"bytes,4,opt,name=auth_algorithm,json=authAlgorithm,proto3,customtype=github.com/MicroOps-cn/idas/pkg/utils/sign.AuthAlgorithm" json:"authAlgorithm,omitempty"`
	AuthKey              string                                                   `protobuf:"bytes,5,opt,name=auth_key,json=authKey,proto3" json:"authKey,omitempty"`
	AuthSecret           string                                                   `protobuf:"bytes,6,opt,name=auth_secret,json=authSecret,proto3" json:"authSecret,omitempty"`
	AuthSign             string                                                   `protobuf:"bytes,7,opt,name=auth_sign,json=authSign,proto3" json:"authSign,omitempty"`
	Payload              string                                                   `protobuf:"bytes,8,opt,name=payload,proto3" json:"-"`
	XXX_NoUnkeyedLiteral struct{}                                                 `json:"-"`
	XXX_unrecognized     []byte                                                   `json:"-"`
	XXX_sizecache        int32                                                    `json:"-"`
}

func (m *AuthenticationRequest) Reset()         { *m = AuthenticationRequest{} }
func (m *AuthenticationRequest) String() string { return proto.CompactTextString(m) }
func (*AuthenticationRequest) ProtoMessage()    {}
func (*AuthenticationRequest) Descriptor() ([]byte, []int) {
	return fileDescriptor_56667f8afe262146, []int{4}
}
func (m *AuthenticationRequest) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_AuthenticationRequest.Unmarshal(m, b)
}
func (m *AuthenticationRequest) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_AuthenticationRequest.Marshal(b, m, deterministic)
}
func (m *AuthenticationRequest) XXX_Merge(src proto.Message) {
	xxx_messageInfo_AuthenticationRequest.Merge(m, src)
}
func (m *AuthenticationRequest) XXX_Size() int {
	return xxx_messageInfo_AuthenticationRequest.Size(m)
}
func (m *AuthenticationRequest) XXX_DiscardUnknown() {
	xxx_messageInfo_AuthenticationRequest.DiscardUnknown(m)
}

var xxx_messageInfo_AuthenticationRequest proto.InternalMessageInfo

func (m *AuthenticationRequest) GetUsername() string {
	if m != nil {
		return m.Username
	}
	return ""
}

func (m *AuthenticationRequest) GetPassword() string {
	if m != nil {
		return m.Password
	}
	return ""
}

func (m *AuthenticationRequest) GetAuthMethod() models.AuthMeta_Method {
	if m != nil {
		return m.AuthMethod
	}
	return models.AuthMeta_basic
}

func (m *AuthenticationRequest) GetAuthKey() string {
	if m != nil {
		return m.AuthKey
	}
	return ""
}

func (m *AuthenticationRequest) GetAuthSecret() string {
	if m != nil {
		return m.AuthSecret
	}
	return ""
}

func (m *AuthenticationRequest) GetAuthSign() string {
	if m != nil {
		return m.AuthSign
	}
	return ""
}

func (m *AuthenticationRequest) GetPayload() string {
	if m != nil {
		return m.Payload
	}
	return ""
}

type UserLoginRequest struct {
	Username             string                                                `protobuf:"bytes,1,opt,name=username,proto3" json:"username,omitempty"`
	Email                string                                                `protobuf:"bytes,2,opt,name=email,proto3" json:"email,omitempty" valid:"email"`
	Phone                string                                                `protobuf:"bytes,3,opt,name=phone,proto3" json:"phone,omitempty"`
	Password             github_com_MicroOps_cn_idas_pkg_service_models.Secret `protobuf:"bytes,4,opt,name=password,proto3,customtype=github.com/MicroOps-cn/idas/pkg/service/models.Secret" json:"password,omitempty"`
	AutoLogin            bool                                                  `protobuf:"varint,5,opt,name=auto_login,json=autoLogin,proto3" json:"autoLogin,omitempty"`
	Type                 LoginType                                             `protobuf:"varint,6,opt,name=type,proto3,enum=idas.endpoint.LoginType" json:"type,omitempty"`
	Code                 string                                                `protobuf:"bytes,7,opt,name=code,proto3" json:"code,omitempty"`
	Token                string                                                `protobuf:"bytes,8,opt,name=token,proto3" json:"token,omitempty"`
	FirstCode            string                                                `protobuf:"bytes,9,opt,name=first_code,json=firstCode,proto3" json:"firstCode,omitempty"`
	SecondCode           string                                                `protobuf:"bytes,10,opt,name=second_code,json=secondCode,proto3" json:"secondCode,omitempty"`
	XXX_NoUnkeyedLiteral struct{}                                              `json:"-"`
	XXX_unrecognized     []byte                                                `json:"-"`
	XXX_sizecache        int32                                                 `json:"-"`
}

func (m *UserLoginRequest) Reset()         { *m = UserLoginRequest{} }
func (m *UserLoginRequest) String() string { return proto.CompactTextString(m) }
func (*UserLoginRequest) ProtoMessage()    {}
func (*UserLoginRequest) Descriptor() ([]byte, []int) {
	return fileDescriptor_56667f8afe262146, []int{5}
}
func (m *UserLoginRequest) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_UserLoginRequest.Unmarshal(m, b)
}
func (m *UserLoginRequest) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_UserLoginRequest.Marshal(b, m, deterministic)
}
func (m *UserLoginRequest) XXX_Merge(src proto.Message) {
	xxx_messageInfo_UserLoginRequest.Merge(m, src)
}
func (m *UserLoginRequest) XXX_Size() int {
	return xxx_messageInfo_UserLoginRequest.Size(m)
}
func (m *UserLoginRequest) XXX_DiscardUnknown() {
	xxx_messageInfo_UserLoginRequest.DiscardUnknown(m)
}

var xxx_messageInfo_UserLoginRequest proto.InternalMessageInfo

func (m *UserLoginRequest) GetUsername() string {
	if m != nil {
		return m.Username
	}
	return ""
}

func (m *UserLoginRequest) GetEmail() string {
	if m != nil {
		return m.Email
	}
	return ""
}

func (m *UserLoginRequest) GetPhone() string {
	if m != nil {
		return m.Phone
	}
	return ""
}

func (m *UserLoginRequest) GetAutoLogin() bool {
	if m != nil {
		return m.AutoLogin
	}
	return false
}

func (m *UserLoginRequest) GetType() LoginType {
	if m != nil {
		return m.Type
	}
	return LoginType_normal
}

func (m *UserLoginRequest) GetCode() string {
	if m != nil {
		return m.Code
	}
	return ""
}

func (m *UserLoginRequest) GetToken() string {
	if m != nil {
		return m.Token
	}
	return ""
}

func (m *UserLoginRequest) GetFirstCode() string {
	if m != nil {
		return m.FirstCode
	}
	return ""
}

func (m *UserLoginRequest) GetSecondCode() string {
	if m != nil {
		return m.SecondCode
	}
	return ""
}

type UserLoginResponseData struct {
	Token                string      `protobuf:"bytes,1,opt,name=token,proto3" json:"token,omitempty"`
	NextMethod           []LoginType `protobuf:"varint,2,rep,packed,name=next_method,json=nextMethod,proto3,enum=idas.endpoint.LoginType" json:"nextMethod"`
	Email                string      `protobuf:"bytes,3,opt,name=email,proto3" json:"email,omitempty"`
	PhoneNumber          string      `protobuf:"bytes,4,opt,name=phone_number,json=phoneNumber,proto3" json:"phone_number,omitempty"`
	XXX_NoUnkeyedLiteral struct{}    `json:"-"`
	XXX_unrecognized     []byte      `json:"-"`
	XXX_sizecache        int32       `json:"-"`
}

func (m *UserLoginResponseData) Reset()         { *m = UserLoginResponseData{} }
func (m *UserLoginResponseData) String() string { return proto.CompactTextString(m) }
func (*UserLoginResponseData) ProtoMessage()    {}
func (*UserLoginResponseData) Descriptor() ([]byte, []int) {
	return fileDescriptor_56667f8afe262146, []int{6}
}
func (m *UserLoginResponseData) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_UserLoginResponseData.Unmarshal(m, b)
}
func (m *UserLoginResponseData) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_UserLoginResponseData.Marshal(b, m, deterministic)
}
func (m *UserLoginResponseData) XXX_Merge(src proto.Message) {
	xxx_messageInfo_UserLoginResponseData.Merge(m, src)
}
func (m *UserLoginResponseData) XXX_Size() int {
	return xxx_messageInfo_UserLoginResponseData.Size(m)
}
func (m *UserLoginResponseData) XXX_DiscardUnknown() {
	xxx_messageInfo_UserLoginResponseData.DiscardUnknown(m)
}

var xxx_messageInfo_UserLoginResponseData proto.InternalMessageInfo

func (m *UserLoginResponseData) GetToken() string {
	if m != nil {
		return m.Token
	}
	return ""
}

func (m *UserLoginResponseData) GetNextMethod() []LoginType {
	if m != nil {
		return m.NextMethod
	}
	return nil
}

func (m *UserLoginResponseData) GetEmail() string {
	if m != nil {
		return m.Email
	}
	return ""
}

func (m *UserLoginResponseData) GetPhoneNumber() string {
	if m != nil {
		return m.PhoneNumber
	}
	return ""
}

type UserLoginResponse struct {
	BaseResponse         `protobuf:"bytes,1,opt,name=base_response,json=baseResponse,proto3,embedded=base_response" json:",omitempty"`
	Data                 *UserLoginResponseData `protobuf:"bytes,2,opt,name=data,proto3" json:"data,omitempty"`
	XXX_NoUnkeyedLiteral struct{}               `json:"-"`
	XXX_unrecognized     []byte                 `json:"-"`
	XXX_sizecache        int32                  `json:"-"`
}

func (m *UserLoginResponse) Reset()         { *m = UserLoginResponse{} }
func (m *UserLoginResponse) String() string { return proto.CompactTextString(m) }
func (*UserLoginResponse) ProtoMessage()    {}
func (*UserLoginResponse) Descriptor() ([]byte, []int) {
	return fileDescriptor_56667f8afe262146, []int{7}
}
func (m *UserLoginResponse) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_UserLoginResponse.Unmarshal(m, b)
}
func (m *UserLoginResponse) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_UserLoginResponse.Marshal(b, m, deterministic)
}
func (m *UserLoginResponse) XXX_Merge(src proto.Message) {
	xxx_messageInfo_UserLoginResponse.Merge(m, src)
}
func (m *UserLoginResponse) XXX_Size() int {
	return xxx_messageInfo_UserLoginResponse.Size(m)
}
func (m *UserLoginResponse) XXX_DiscardUnknown() {
	xxx_messageInfo_UserLoginResponse.DiscardUnknown(m)
}

var xxx_messageInfo_UserLoginResponse proto.InternalMessageInfo

func (m *UserLoginResponse) GetData() *UserLoginResponseData {
	if m != nil {
		return m.Data
	}
	return nil
}

type SendLoginCaptchaRequest struct {
	Username             string    `protobuf:"bytes,1,opt,name=username,proto3" json:"username,omitempty"`
	Type                 LoginType `protobuf:"varint,2,opt,name=type,proto3,enum=idas.endpoint.LoginType" json:"type" valid:"required"`
	Email                string    `protobuf:"bytes,3,opt,name=email,proto3" json:"email,omitempty" valid:"email"`
	Phone                string    `protobuf:"bytes,4,opt,name=phone,proto3" json:"phone,omitempty"`
	XXX_NoUnkeyedLiteral struct{}  `json:"-"`
	XXX_unrecognized     []byte    `json:"-"`
	XXX_sizecache        int32     `json:"-"`
}

func (m *SendLoginCaptchaRequest) Reset()         { *m = SendLoginCaptchaRequest{} }
func (m *SendLoginCaptchaRequest) String() string { return proto.CompactTextString(m) }
func (*SendLoginCaptchaRequest) ProtoMessage()    {}
func (*SendLoginCaptchaRequest) Descriptor() ([]byte, []int) {
	return fileDescriptor_56667f8afe262146, []int{8}
}
func (m *SendLoginCaptchaRequest) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_SendLoginCaptchaRequest.Unmarshal(m, b)
}
func (m *SendLoginCaptchaRequest) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_SendLoginCaptchaRequest.Marshal(b, m, deterministic)
}
func (m *SendLoginCaptchaRequest) XXX_Merge(src proto.Message) {
	xxx_messageInfo_SendLoginCaptchaRequest.Merge(m, src)
}
func (m *SendLoginCaptchaRequest) XXX_Size() int {
	return xxx_messageInfo_SendLoginCaptchaRequest.Size(m)
}
func (m *SendLoginCaptchaRequest) XXX_DiscardUnknown() {
	xxx_messageInfo_SendLoginCaptchaRequest.DiscardUnknown(m)
}

var xxx_messageInfo_SendLoginCaptchaRequest proto.InternalMessageInfo

func (m *SendLoginCaptchaRequest) GetUsername() string {
	if m != nil {
		return m.Username
	}
	return ""
}

func (m *SendLoginCaptchaRequest) GetType() LoginType {
	if m != nil {
		return m.Type
	}
	return LoginType_normal
}

func (m *SendLoginCaptchaRequest) GetEmail() string {
	if m != nil {
		return m.Email
	}
	return ""
}

func (m *SendLoginCaptchaRequest) GetPhone() string {
	if m != nil {
		return m.Phone
	}
	return ""
}

type SendLoginCaptchaResponseData struct {
	Token                string   `protobuf:"bytes,1,opt,name=token,proto3" json:"token,omitempty"`
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_unrecognized     []byte   `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *SendLoginCaptchaResponseData) Reset()         { *m = SendLoginCaptchaResponseData{} }
func (m *SendLoginCaptchaResponseData) String() string { return proto.CompactTextString(m) }
func (*SendLoginCaptchaResponseData) ProtoMessage()    {}
func (*SendLoginCaptchaResponseData) Descriptor() ([]byte, []int) {
	return fileDescriptor_56667f8afe262146, []int{9}
}
func (m *SendLoginCaptchaResponseData) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_SendLoginCaptchaResponseData.Unmarshal(m, b)
}
func (m *SendLoginCaptchaResponseData) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_SendLoginCaptchaResponseData.Marshal(b, m, deterministic)
}
func (m *SendLoginCaptchaResponseData) XXX_Merge(src proto.Message) {
	xxx_messageInfo_SendLoginCaptchaResponseData.Merge(m, src)
}
func (m *SendLoginCaptchaResponseData) XXX_Size() int {
	return xxx_messageInfo_SendLoginCaptchaResponseData.Size(m)
}
func (m *SendLoginCaptchaResponseData) XXX_DiscardUnknown() {
	xxx_messageInfo_SendLoginCaptchaResponseData.DiscardUnknown(m)
}

var xxx_messageInfo_SendLoginCaptchaResponseData proto.InternalMessageInfo

func (m *SendLoginCaptchaResponseData) GetToken() string {
	if m != nil {
		return m.Token
	}
	return ""
}

type SendLoginCaptchaResponse struct {
	BaseResponse         `protobuf:"bytes,1,opt,name=base_response,json=baseResponse,proto3,embedded=base_response" json:",omitempty"`
	Data                 *SendLoginCaptchaResponseData `protobuf:"bytes,2,opt,name=data,proto3" json:"data,omitempty"`
	XXX_NoUnkeyedLiteral struct{}                      `json:"-"`
	XXX_unrecognized     []byte                        `json:"-"`
	XXX_sizecache        int32                         `json:"-"`
}

func (m *SendLoginCaptchaResponse) Reset()         { *m = SendLoginCaptchaResponse{} }
func (m *SendLoginCaptchaResponse) String() string { return proto.CompactTextString(m) }
func (*SendLoginCaptchaResponse) ProtoMessage()    {}
func (*SendLoginCaptchaResponse) Descriptor() ([]byte, []int) {
	return fileDescriptor_56667f8afe262146, []int{10}
}
func (m *SendLoginCaptchaResponse) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_SendLoginCaptchaResponse.Unmarshal(m, b)
}
func (m *SendLoginCaptchaResponse) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_SendLoginCaptchaResponse.Marshal(b, m, deterministic)
}
func (m *SendLoginCaptchaResponse) XXX_Merge(src proto.Message) {
	xxx_messageInfo_SendLoginCaptchaResponse.Merge(m, src)
}
func (m *SendLoginCaptchaResponse) XXX_Size() int {
	return xxx_messageInfo_SendLoginCaptchaResponse.Size(m)
}
func (m *SendLoginCaptchaResponse) XXX_DiscardUnknown() {
	xxx_messageInfo_SendLoginCaptchaResponse.DiscardUnknown(m)
}

var xxx_messageInfo_SendLoginCaptchaResponse proto.InternalMessageInfo

func (m *SendLoginCaptchaResponse) GetData() *SendLoginCaptchaResponseData {
	if m != nil {
		return m.Data
	}
	return nil
}

func init() {
	proto.RegisterEnum("idas.endpoint.LoginType", LoginType_name, LoginType_value)
	proto.RegisterType((*GetSessionsRequest)(nil), "idas.endpoint.GetSessionsRequest")
	proto.RegisterType((*SessionInfo)(nil), "idas.endpoint.SessionInfo")
	proto.RegisterType((*GetSessionsResponse)(nil), "idas.endpoint.GetSessionsResponse")
	proto.RegisterType((*DeleteSessionRequest)(nil), "idas.endpoint.DeleteSessionRequest")
	proto.RegisterType((*AuthenticationRequest)(nil), "idas.endpoint.AuthenticationRequest")
	proto.RegisterType((*UserLoginRequest)(nil), "idas.endpoint.UserLoginRequest")
	proto.RegisterType((*UserLoginResponseData)(nil), "idas.endpoint.UserLoginResponseData")
	proto.RegisterType((*UserLoginResponse)(nil), "idas.endpoint.UserLoginResponse")
	proto.RegisterType((*SendLoginCaptchaRequest)(nil), "idas.endpoint.SendLoginCaptchaRequest")
	proto.RegisterType((*SendLoginCaptchaResponseData)(nil), "idas.endpoint.SendLoginCaptchaResponseData")
	proto.RegisterType((*SendLoginCaptchaResponse)(nil), "idas.endpoint.SendLoginCaptchaResponse")
}

func init() { proto.RegisterFile("endpoints/sessions.proto", fileDescriptor_56667f8afe262146) }

var fileDescriptor_56667f8afe262146 = []byte{
	// 1133 bytes of a gzipped FileDescriptorProto
	0x1f, 0x8b, 0x08, 0x00, 0x00, 0x00, 0x00, 0x00, 0x02, 0xff, 0xb4, 0x56, 0xc1, 0x6f, 0x1b, 0xc5,
	0x17, 0xee, 0xda, 0x4e, 0x6c, 0x3f, 0x27, 0xa9, 0x33, 0x71, 0xda, 0x55, 0x7e, 0x55, 0xb7, 0xbf,
	0x55, 0xa5, 0x56, 0xd0, 0xda, 0x90, 0x42, 0x55, 0xe0, 0x80, 0xea, 0x56, 0x54, 0x85, 0x16, 0xd0,
	0xa6, 0x48, 0x08, 0x0e, 0xab, 0xb1, 0xf7, 0xc5, 0x1e, 0x65, 0x77, 0x67, 0xd9, 0x19, 0x97, 0xfa,
	0xca, 0xdf, 0xc0, 0xff, 0xc0, 0x85, 0x0b, 0x48, 0xdc, 0x39, 0x72, 0xe4, 0xcc, 0x61, 0x25, 0xae,
	0x3e, 0xa1, 0x9c, 0x38, 0xa2, 0x99, 0xd9, 0xb5, 0xd7, 0x4e, 0xd2, 0x84, 0x03, 0xa7, 0xec, 0x7c,
	0xef, 0xbd, 0x99, 0xf7, 0xde, 0xf7, 0xbd, 0xe7, 0x80, 0x8d, 0x71, 0x90, 0x70, 0x16, 0x4b, 0xd1,
	0x13, 0x28, 0x04, 0xe3, 0xb1, 0xe8, 0x26, 0x29, 0x97, 0x9c, 0x6c, 0xb2, 0x80, 0x8a, 0x6e, 0x61,
	0xde, 0xeb, 0x8c, 0xf8, 0x88, 0x6b, 0x4b, 0x4f, 0x7d, 0x19, 0xa7, 0xbd, 0xce, 0x22, 0x7c, 0x40,
	0x05, 0xe6, 0xe8, 0x4e, 0xc4, 0x03, 0x0c, 0x45, 0xcf, 0xfc, 0x31, 0xa0, 0xfb, 0xb3, 0x05, 0xe4,
	0x09, 0xca, 0x83, 0xfc, 0x15, 0x0f, 0xbf, 0x99, 0xa0, 0x90, 0xe4, 0x05, 0x6c, 0xab, 0x48, 0x3f,
	0x64, 0x42, 0xfa, 0xa9, 0x01, 0x6d, 0xeb, 0x86, 0x75, 0xbb, 0xb5, 0x7f, 0xbd, 0xbb, 0x94, 0x42,
	0xb7, 0x4f, 0x05, 0x3e, 0x63, 0x42, 0xe6, 0xa1, 0xfd, 0xc6, 0x6f, 0x99, 0x73, 0xe9, 0xf7, 0xcc,
	0xb1, 0xbc, 0xcb, 0x83, 0x65, 0x13, 0xf9, 0x08, 0xea, 0x13, 0x81, 0xa9, 0xcf, 0x02, 0xbb, 0x72,
	0xc3, 0xba, 0xdd, 0xec, 0xdf, 0x9d, 0x65, 0x4e, 0x5b, 0x41, 0x4f, 0x83, 0x3b, 0x3c, 0x62, 0x12,
	0xa3, 0x44, 0x4e, 0x8f, 0x33, 0x67, 0xf7, 0x25, 0x0d, 0x59, 0xf0, 0xbe, 0xab, 0x5e, 0x66, 0x29,
	0x06, 0x77, 0x26, 0x13, 0x16, 0xb8, 0xde, 0xba, 0x71, 0x75, 0xff, 0xb2, 0xa0, 0x95, 0x67, 0xfc,
	0x34, 0x3e, 0xe4, 0xa4, 0x07, 0x15, 0x16, 0xe8, 0xf4, 0x9a, 0x7d, 0x67, 0x96, 0x39, 0x15, 0x16,
	0x9c, 0x7d, 0x49, 0x85, 0x05, 0xe4, 0x11, 0xb4, 0x86, 0x29, 0x52, 0x89, 0xbe, 0x64, 0x11, 0xe6,
	0xc9, 0xb8, 0xb3, 0xcc, 0x01, 0x03, 0xbf, 0x60, 0x11, 0x1e, 0x67, 0x4e, 0x7b, 0xe5, 0x06, 0xd7,
	0x2b, 0xd9, 0xc9, 0x7d, 0x58, 0xc7, 0x57, 0x09, 0x4b, 0xa7, 0x76, 0x55, 0xc7, 0x5f, 0x9f, 0x65,
	0x4e, 0x8e, 0x9c, 0x1a, 0x9b, 0xdb, 0xc8, 0x3d, 0x68, 0x86, 0x54, 0x48, 0x5f, 0x20, 0xc6, 0x76,
	0x4d, 0x87, 0x5e, 0x99, 0x65, 0x0e, 0x51, 0xe0, 0x01, 0x62, 0xbc, 0xe8, 0x84, 0xd7, 0x28, 0x30,
	0xf7, 0x47, 0x0b, 0x76, 0x96, 0x78, 0x12, 0x09, 0x8f, 0x05, 0x12, 0x04, 0x52, 0x26, 0xca, 0xa0,
	0x39, 0x53, 0xce, 0x99, 0x4c, 0x19, 0xb7, 0xfe, 0x95, 0x82, 0x2a, 0x55, 0x79, 0xe9, 0xd9, 0xf6,
	0x60, 0xc5, 0x93, 0x74, 0xa1, 0x16, 0x50, 0x49, 0xed, 0xca, 0x8d, 0xea, 0xed, 0xd6, 0xfe, 0xde,
	0xca, 0xc5, 0x25, 0x2e, 0x3c, 0xed, 0xe7, 0x3e, 0x81, 0xce, 0x63, 0x0c, 0x51, 0x62, 0x6e, 0x2a,
	0x14, 0xf0, 0x6f, 0x99, 0x72, 0xff, 0xae, 0xc2, 0xee, 0xc3, 0x89, 0x1c, 0x63, 0x2c, 0xd9, 0x90,
	0xca, 0xd2, 0x55, 0x7b, 0xd0, 0x50, 0x72, 0x88, 0x69, 0x64, 0xea, 0x6d, 0x7a, 0xf3, 0xb3, 0xb2,
	0x25, 0x54, 0x88, 0x6f, 0x79, 0x9a, 0x2b, 0xcd, 0x9b, 0x9f, 0x89, 0x0f, 0x2d, 0x3a, 0x91, 0x63,
	0x3f, 0x42, 0x39, 0xe6, 0x81, 0xe6, 0x6e, 0x6b, 0xff, 0xa6, 0xa9, 0x48, 0x60, 0xfa, 0x92, 0x0d,
	0xb1, 0x9b, 0x8f, 0x88, 0x7a, 0xf8, 0x39, 0x4a, 0xda, 0x7d, 0xae, 0x7d, 0xfb, 0xf6, 0x2c, 0x73,
	0x3a, 0xd4, 0x80, 0x63, 0x5e, 0x92, 0xac, 0x07, 0x0b, 0x94, 0x7c, 0x67, 0xc1, 0x96, 0x7e, 0x81,
	0x86, 0x23, 0x9e, 0x32, 0x39, 0x8e, 0x72, 0x96, 0xbf, 0x56, 0xed, 0xfe, 0x23, 0x73, 0x1e, 0x8c,
	0x98, 0x1c, 0x4f, 0x06, 0xdd, 0x21, 0x8f, 0x7a, 0xcf, 0xd9, 0x30, 0xe5, 0x9f, 0x25, 0xe2, 0xee,
	0x30, 0xee, 0xa9, 0x14, 0x7a, 0xc9, 0xd1, 0xa8, 0x37, 0x91, 0x2c, 0x14, 0x3d, 0xc1, 0x46, 0xb1,
	0x4e, 0xe1, 0x61, 0x71, 0xcf, 0x2c, 0x73, 0xae, 0xd2, 0x32, 0x50, 0xca, 0x60, 0x73, 0xc9, 0x40,
	0xde, 0x82, 0x86, 0xce, 0xe1, 0x08, 0xa7, 0xf6, 0x9a, 0x7e, 0x7d, 0x77, 0x96, 0x39, 0xdb, 0x0a,
	0xfb, 0x04, 0xa7, 0xa5, 0xb8, 0x7a, 0x0e, 0x91, 0xf7, 0xf2, 0xbe, 0x08, 0x1c, 0xa6, 0x28, 0xed,
	0x75, 0x1d, 0x34, 0xaf, 0xf8, 0x40, 0xa3, 0xab, 0x15, 0x1b, 0x54, 0x29, 0xda, 0x84, 0xb2, 0x51,
	0x6c, 0xd7, 0x17, 0x8a, 0xd6, 0x2e, 0x6c, 0xb4, 0xa4, 0xe8, 0x02, 0x23, 0x0e, 0xd4, 0x13, 0x3a,
	0x0d, 0x39, 0x0d, 0xec, 0x86, 0x0e, 0x59, 0x9b, 0x65, 0x8e, 0x75, 0xd7, 0x2b, 0x50, 0xf7, 0xcf,
	0x2a, 0xb4, 0xbf, 0x10, 0x98, 0x3e, 0xe3, 0x23, 0x76, 0x21, 0xd6, 0x6f, 0xc1, 0x1a, 0x46, 0x94,
	0x85, 0xf9, 0x3c, 0x6f, 0x1f, 0x67, 0xce, 0x66, 0xae, 0x2c, 0x8d, 0xbb, 0x9e, 0xb1, 0x93, 0x0e,
	0xac, 0x25, 0x63, 0x1e, 0xa3, 0x19, 0x5c, 0xcf, 0x1c, 0x08, 0x2f, 0x89, 0xc6, 0x10, 0x76, 0x90,
	0x13, 0xf6, 0xee, 0x79, 0x84, 0xe5, 0xba, 0x29, 0x56, 0xab, 0x69, 0x8b, 0xea, 0x40, 0x71, 0x61,
	0xb9, 0x03, 0x73, 0x25, 0xde, 0x07, 0xd5, 0x44, 0xee, 0x87, 0xaa, 0x40, 0xcd, 0x52, 0xa3, 0x7f,
	0x75, 0x96, 0x39, 0x3b, 0x0a, 0xd5, 0x55, 0x97, 0xc2, 0x9a, 0x73, 0x90, 0xdc, 0x81, 0x9a, 0x9c,
	0x26, 0xa8, 0x29, 0xda, 0xda, 0xb7, 0x57, 0x86, 0x51, 0xfb, 0xbc, 0x98, 0x26, 0xe8, 0x69, 0x2f,
	0x42, 0xa0, 0x36, 0xe4, 0x01, 0x1a, 0x5e, 0x3c, 0xfd, 0xad, 0x1a, 0x20, 0xf9, 0x11, 0xc6, 0xa6,
	0xf3, 0x9e, 0x39, 0xa8, 0x7c, 0x0e, 0x59, 0x2a, 0xa4, 0xaf, 0xfd, 0x9b, 0xba, 0x05, 0x3a, 0x1f,
	0x8d, 0x3e, 0xe2, 0x01, 0x96, 0xf3, 0x99, 0x83, 0x4a, 0x39, 0x02, 0x87, 0x3c, 0x0e, 0x4c, 0x20,
	0x2c, 0x94, 0x63, 0xe0, 0x95, 0x48, 0x58, 0xa0, 0xee, 0x4f, 0x16, 0xec, 0x96, 0x38, 0x36, 0xdb,
	0xe6, 0x31, 0x95, 0x74, 0x91, 0xa2, 0x55, 0x4e, 0xf1, 0x29, 0xb4, 0x62, 0x7c, 0x25, 0x8b, 0xe1,
	0x55, 0xeb, 0xe8, 0x35, 0x1d, 0xe8, 0x6f, 0xa9, 0xc5, 0xa6, 0x02, 0xcc, 0x68, 0x7a, 0xa5, 0x6f,
	0xf5, 0x80, 0x51, 0x4b, 0x2e, 0x02, 0x23, 0x8d, 0xff, 0xc3, 0x86, 0x56, 0x83, 0x1f, 0x4f, 0xa2,
	0x01, 0xa6, 0x46, 0x08, 0x5e, 0x4b, 0x63, 0x9f, 0x6a, 0xc8, 0xfd, 0xc1, 0x82, 0xed, 0x13, 0x39,
	0x93, 0x2f, 0x61, 0x53, 0x2f, 0xe2, 0x95, 0x1d, 0xfc, 0xbf, 0x53, 0x76, 0xf0, 0xb9, 0xfb, 0x77,
	0x63, 0x50, 0xf2, 0x22, 0x0f, 0xe6, 0xbb, 0x57, 0x5d, 0x78, 0x73, 0xe5, 0xc2, 0x53, 0xbb, 0x97,
	0x6f, 0xe1, 0x5f, 0x2d, 0xb8, 0x7a, 0x80, 0x71, 0xa0, 0xed, 0x8f, 0x68, 0x22, 0x87, 0x63, 0x7a,
	0x91, 0x41, 0xfa, 0x38, 0x17, 0x58, 0xe5, 0xf5, 0x02, 0xeb, 0x5f, 0x9b, 0x65, 0x8e, 0xf6, 0x3c,
	0xf5, 0xf7, 0xce, 0xc8, 0xef, 0xd6, 0x52, 0x9b, 0x2f, 0x32, 0x94, 0xb5, 0xd2, 0x50, 0xba, 0xef,
	0xc0, 0xb5, 0x93, 0x15, 0x9c, 0x27, 0x13, 0xf7, 0x17, 0x0b, 0xec, 0xb3, 0xc2, 0xfe, 0x43, 0xa6,
	0x3e, 0x5c, 0x62, 0xea, 0xcd, 0x13, 0xbf, 0x92, 0x67, 0xd7, 0x61, 0x08, 0x7b, 0xe3, 0x7b, 0x0b,
	0x9a, 0xf3, 0xf6, 0x12, 0x80, 0xf5, 0x98, 0xa7, 0x11, 0x0d, 0xdb, 0x97, 0xc8, 0x06, 0x34, 0xa2,
	0x43, 0xea, 0x4b, 0x2e, 0x93, 0xb6, 0x45, 0x36, 0xa1, 0xa9, 0x4e, 0xba, 0x71, 0xed, 0x0a, 0x69,
	0x41, 0x5d, 0x1d, 0x45, 0x24, 0xda, 0x55, 0xd2, 0xcc, 0x1b, 0xde, 0xae, 0x91, 0x3a, 0x54, 0x15,
	0xb6, 0x46, 0x76, 0xe0, 0x32, 0xc6, 0x74, 0x10, 0xa2, 0x3f, 0xbf, 0x04, 0x48, 0x07, 0xda, 0x25,
	0xd0, 0xc4, 0xb4, 0x08, 0x81, 0xad, 0x12, 0xaa, 0xc2, 0x37, 0xfa, 0xf7, 0xbe, 0x7a, 0xfb, 0xbc,
	0x0d, 0x58, 0x54, 0xf9, 0x41, 0xf1, 0xf1, 0xf9, 0xa5, 0xc1, 0xba, 0xfe, 0x17, 0xf3, 0xde, 0x3f,
	0x01, 0x00, 0x00, 0xff, 0xff, 0x93, 0x55, 0x45, 0x47, 0xce, 0x0a, 0x00, 0x00,
}
