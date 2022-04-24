// Code generated by protoc-gen-go. DO NOT EDIT.
// versions:
// 	protoc-gen-go v1.27.1
// 	protoc        v3.19.4
// source: types/ldap.proto

package ldap

import (
	protoreflect "google.golang.org/protobuf/reflect/protoreflect"
	protoimpl "google.golang.org/protobuf/runtime/protoimpl"
	reflect "reflect"
	sync "sync"
)

const (
	// Verify that this generated code is sufficiently up-to-date.
	_ = protoimpl.EnforceVersion(20 - protoimpl.MinVersion)
	// Verify that runtime/protoimpl is sufficiently up-to-date.
	_ = protoimpl.EnforceVersion(protoimpl.MaxVersion - 20)
)

type LdapOptions struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Host                string `protobuf:"bytes,1,opt,name=host,proto3" json:"host,omitempty"`
	ManagerDn           string `protobuf:"bytes,2,opt,name=manager_dn,json=managerDn,proto3" json:"manager_dn,omitempty"`
	ManagerPassword     string `protobuf:"bytes,3,opt,name=manager_password,json=managerPassword,proto3" json:"manager_password,omitempty"`
	UserSearchBase      string `protobuf:"bytes,4,opt,name=user_search_base,json=userSearchBase,proto3" json:"user_search_base,omitempty"`
	UserSearchFilter    string `protobuf:"bytes,5,opt,name=user_search_filter,json=userSearchFilter,proto3" json:"user_search_filter,omitempty"`
	GroupSearchBase     string `protobuf:"bytes,6,opt,name=group_search_base,json=groupSearchBase,proto3" json:"group_search_base,omitempty"`
	GroupSearchFilter   string `protobuf:"bytes,7,opt,name=group_search_filter,json=groupSearchFilter,proto3" json:"group_search_filter,omitempty"`
	AttrUsername        string `protobuf:"bytes,8,opt,name=attr_username,json=attrUsername,proto3" json:"attr_username,omitempty"`
	AttrEmail           string `protobuf:"bytes,9,opt,name=attr_email,json=attrEmail,proto3" json:"attr_email,omitempty"`
	AttrUserDisplayName string `protobuf:"bytes,10,opt,name=attr_user_display_name,json=attrUserDisplayName,proto3" json:"attr_user_display_name,omitempty"`
	AttrUserPhoneNo     string `protobuf:"bytes,11,opt,name=attr_user_phone_no,json=attrUserPhoneNo,proto3" json:"attr_user_phone_no,omitempty"`
}

func (x *LdapOptions) Reset() {
	*x = LdapOptions{}
	if protoimpl.UnsafeEnabled {
		mi := &file_types_ldap_proto_msgTypes[0]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *LdapOptions) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*LdapOptions) ProtoMessage() {}

func (x *LdapOptions) ProtoReflect() protoreflect.Message {
	mi := &file_types_ldap_proto_msgTypes[0]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use LdapOptions.ProtoReflect.Descriptor instead.
func (*LdapOptions) Descriptor() ([]byte, []int) {
	return file_types_ldap_proto_rawDescGZIP(), []int{0}
}

func (x *LdapOptions) GetHost() string {
	if x != nil {
		return x.Host
	}
	return ""
}

func (x *LdapOptions) GetManagerDn() string {
	if x != nil {
		return x.ManagerDn
	}
	return ""
}

func (x *LdapOptions) GetManagerPassword() string {
	if x != nil {
		return x.ManagerPassword
	}
	return ""
}

func (x *LdapOptions) GetUserSearchBase() string {
	if x != nil {
		return x.UserSearchBase
	}
	return ""
}

func (x *LdapOptions) GetUserSearchFilter() string {
	if x != nil {
		return x.UserSearchFilter
	}
	return ""
}

func (x *LdapOptions) GetGroupSearchBase() string {
	if x != nil {
		return x.GroupSearchBase
	}
	return ""
}

func (x *LdapOptions) GetGroupSearchFilter() string {
	if x != nil {
		return x.GroupSearchFilter
	}
	return ""
}

func (x *LdapOptions) GetAttrUsername() string {
	if x != nil {
		return x.AttrUsername
	}
	return ""
}

func (x *LdapOptions) GetAttrEmail() string {
	if x != nil {
		return x.AttrEmail
	}
	return ""
}

func (x *LdapOptions) GetAttrUserDisplayName() string {
	if x != nil {
		return x.AttrUserDisplayName
	}
	return ""
}

func (x *LdapOptions) GetAttrUserPhoneNo() string {
	if x != nil {
		return x.AttrUserPhoneNo
	}
	return ""
}

var File_types_ldap_proto protoreflect.FileDescriptor

var file_types_ldap_proto_rawDesc = []byte{
	0x0a, 0x10, 0x74, 0x79, 0x70, 0x65, 0x73, 0x2f, 0x6c, 0x64, 0x61, 0x70, 0x2e, 0x70, 0x72, 0x6f,
	0x74, 0x6f, 0x12, 0x10, 0x69, 0x64, 0x61, 0x73, 0x2e, 0x63, 0x6c, 0x69, 0x65, 0x6e, 0x74, 0x2e,
	0x6c, 0x64, 0x61, 0x70, 0x22, 0xc5, 0x03, 0x0a, 0x0b, 0x4c, 0x64, 0x61, 0x70, 0x4f, 0x70, 0x74,
	0x69, 0x6f, 0x6e, 0x73, 0x12, 0x12, 0x0a, 0x04, 0x68, 0x6f, 0x73, 0x74, 0x18, 0x01, 0x20, 0x01,
	0x28, 0x09, 0x52, 0x04, 0x68, 0x6f, 0x73, 0x74, 0x12, 0x1d, 0x0a, 0x0a, 0x6d, 0x61, 0x6e, 0x61,
	0x67, 0x65, 0x72, 0x5f, 0x64, 0x6e, 0x18, 0x02, 0x20, 0x01, 0x28, 0x09, 0x52, 0x09, 0x6d, 0x61,
	0x6e, 0x61, 0x67, 0x65, 0x72, 0x44, 0x6e, 0x12, 0x29, 0x0a, 0x10, 0x6d, 0x61, 0x6e, 0x61, 0x67,
	0x65, 0x72, 0x5f, 0x70, 0x61, 0x73, 0x73, 0x77, 0x6f, 0x72, 0x64, 0x18, 0x03, 0x20, 0x01, 0x28,
	0x09, 0x52, 0x0f, 0x6d, 0x61, 0x6e, 0x61, 0x67, 0x65, 0x72, 0x50, 0x61, 0x73, 0x73, 0x77, 0x6f,
	0x72, 0x64, 0x12, 0x28, 0x0a, 0x10, 0x75, 0x73, 0x65, 0x72, 0x5f, 0x73, 0x65, 0x61, 0x72, 0x63,
	0x68, 0x5f, 0x62, 0x61, 0x73, 0x65, 0x18, 0x04, 0x20, 0x01, 0x28, 0x09, 0x52, 0x0e, 0x75, 0x73,
	0x65, 0x72, 0x53, 0x65, 0x61, 0x72, 0x63, 0x68, 0x42, 0x61, 0x73, 0x65, 0x12, 0x2c, 0x0a, 0x12,
	0x75, 0x73, 0x65, 0x72, 0x5f, 0x73, 0x65, 0x61, 0x72, 0x63, 0x68, 0x5f, 0x66, 0x69, 0x6c, 0x74,
	0x65, 0x72, 0x18, 0x05, 0x20, 0x01, 0x28, 0x09, 0x52, 0x10, 0x75, 0x73, 0x65, 0x72, 0x53, 0x65,
	0x61, 0x72, 0x63, 0x68, 0x46, 0x69, 0x6c, 0x74, 0x65, 0x72, 0x12, 0x2a, 0x0a, 0x11, 0x67, 0x72,
	0x6f, 0x75, 0x70, 0x5f, 0x73, 0x65, 0x61, 0x72, 0x63, 0x68, 0x5f, 0x62, 0x61, 0x73, 0x65, 0x18,
	0x06, 0x20, 0x01, 0x28, 0x09, 0x52, 0x0f, 0x67, 0x72, 0x6f, 0x75, 0x70, 0x53, 0x65, 0x61, 0x72,
	0x63, 0x68, 0x42, 0x61, 0x73, 0x65, 0x12, 0x2e, 0x0a, 0x13, 0x67, 0x72, 0x6f, 0x75, 0x70, 0x5f,
	0x73, 0x65, 0x61, 0x72, 0x63, 0x68, 0x5f, 0x66, 0x69, 0x6c, 0x74, 0x65, 0x72, 0x18, 0x07, 0x20,
	0x01, 0x28, 0x09, 0x52, 0x11, 0x67, 0x72, 0x6f, 0x75, 0x70, 0x53, 0x65, 0x61, 0x72, 0x63, 0x68,
	0x46, 0x69, 0x6c, 0x74, 0x65, 0x72, 0x12, 0x23, 0x0a, 0x0d, 0x61, 0x74, 0x74, 0x72, 0x5f, 0x75,
	0x73, 0x65, 0x72, 0x6e, 0x61, 0x6d, 0x65, 0x18, 0x08, 0x20, 0x01, 0x28, 0x09, 0x52, 0x0c, 0x61,
	0x74, 0x74, 0x72, 0x55, 0x73, 0x65, 0x72, 0x6e, 0x61, 0x6d, 0x65, 0x12, 0x1d, 0x0a, 0x0a, 0x61,
	0x74, 0x74, 0x72, 0x5f, 0x65, 0x6d, 0x61, 0x69, 0x6c, 0x18, 0x09, 0x20, 0x01, 0x28, 0x09, 0x52,
	0x09, 0x61, 0x74, 0x74, 0x72, 0x45, 0x6d, 0x61, 0x69, 0x6c, 0x12, 0x33, 0x0a, 0x16, 0x61, 0x74,
	0x74, 0x72, 0x5f, 0x75, 0x73, 0x65, 0x72, 0x5f, 0x64, 0x69, 0x73, 0x70, 0x6c, 0x61, 0x79, 0x5f,
	0x6e, 0x61, 0x6d, 0x65, 0x18, 0x0a, 0x20, 0x01, 0x28, 0x09, 0x52, 0x13, 0x61, 0x74, 0x74, 0x72,
	0x55, 0x73, 0x65, 0x72, 0x44, 0x69, 0x73, 0x70, 0x6c, 0x61, 0x79, 0x4e, 0x61, 0x6d, 0x65, 0x12,
	0x2b, 0x0a, 0x12, 0x61, 0x74, 0x74, 0x72, 0x5f, 0x75, 0x73, 0x65, 0x72, 0x5f, 0x70, 0x68, 0x6f,
	0x6e, 0x65, 0x5f, 0x6e, 0x6f, 0x18, 0x0b, 0x20, 0x01, 0x28, 0x09, 0x52, 0x0f, 0x61, 0x74, 0x74,
	0x72, 0x55, 0x73, 0x65, 0x72, 0x50, 0x68, 0x6f, 0x6e, 0x65, 0x4e, 0x6f, 0x42, 0x1b, 0x5a, 0x19,
	0x69, 0x64, 0x61, 0x73, 0x2f, 0x70, 0x6b, 0x67, 0x2f, 0x63, 0x6c, 0x69, 0x65, 0x6e, 0x74, 0x2f,
	0x6c, 0x64, 0x61, 0x70, 0x3b, 0x6c, 0x64, 0x61, 0x70, 0x62, 0x06, 0x70, 0x72, 0x6f, 0x74, 0x6f,
	0x33,
}

var (
	file_types_ldap_proto_rawDescOnce sync.Once
	file_types_ldap_proto_rawDescData = file_types_ldap_proto_rawDesc
)

func file_types_ldap_proto_rawDescGZIP() []byte {
	file_types_ldap_proto_rawDescOnce.Do(func() {
		file_types_ldap_proto_rawDescData = protoimpl.X.CompressGZIP(file_types_ldap_proto_rawDescData)
	})
	return file_types_ldap_proto_rawDescData
}

var file_types_ldap_proto_msgTypes = make([]protoimpl.MessageInfo, 1)
var file_types_ldap_proto_goTypes = []interface{}{
	(*LdapOptions)(nil), // 0: idas.client.ldap.LdapOptions
}
var file_types_ldap_proto_depIdxs = []int32{
	0, // [0:0] is the sub-list for method output_type
	0, // [0:0] is the sub-list for method input_type
	0, // [0:0] is the sub-list for extension type_name
	0, // [0:0] is the sub-list for extension extendee
	0, // [0:0] is the sub-list for field type_name
}

func init() { file_types_ldap_proto_init() }
func file_types_ldap_proto_init() {
	if File_types_ldap_proto != nil {
		return
	}
	if !protoimpl.UnsafeEnabled {
		file_types_ldap_proto_msgTypes[0].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*LdapOptions); i {
			case 0:
				return &v.state
			case 1:
				return &v.sizeCache
			case 2:
				return &v.unknownFields
			default:
				return nil
			}
		}
	}
	type x struct{}
	out := protoimpl.TypeBuilder{
		File: protoimpl.DescBuilder{
			GoPackagePath: reflect.TypeOf(x{}).PkgPath(),
			RawDescriptor: file_types_ldap_proto_rawDesc,
			NumEnums:      0,
			NumMessages:   1,
			NumExtensions: 0,
			NumServices:   0,
		},
		GoTypes:           file_types_ldap_proto_goTypes,
		DependencyIndexes: file_types_ldap_proto_depIdxs,
		MessageInfos:      file_types_ldap_proto_msgTypes,
	}.Build()
	File_types_ldap_proto = out.File
	file_types_ldap_proto_rawDesc = nil
	file_types_ldap_proto_goTypes = nil
	file_types_ldap_proto_depIdxs = nil
}
