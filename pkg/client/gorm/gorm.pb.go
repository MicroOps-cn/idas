// Code generated by protoc-gen-go. DO NOT EDIT.
// versions:
// 	protoc-gen-go v1.26.0
// 	protoc        v3.10.1
// source: types/gorm.proto

package gorm

import (
	types "github.com/gogo/protobuf/types"
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

// Symbols defined in public import of google/protobuf/duration.proto.

type Duration = types.Duration

type MySQLOptions struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Host                  string          `protobuf:"bytes,1,opt,name=host,proto3" json:"host,omitempty"`
	Username              string          `protobuf:"bytes,2,opt,name=username,proto3" json:"username,omitempty"`
	Password              string          `protobuf:"bytes,3,opt,name=password,proto3" json:"password,omitempty"`
	Schema                string          `protobuf:"bytes,4,opt,name=schema,proto3" json:"schema,omitempty"`
	MaxIdle               int32           `protobuf:"varint,5,opt,name=max_idle,json=maxIdle,proto3" json:"max_idle,omitempty"`
	MaxIdleConnections    int32           `protobuf:"varint,6,opt,name=max_idle_connections,json=maxIdleConnections,proto3" json:"max_idle_connections,omitempty"`
	MaxOpenConnections    int32           `protobuf:"varint,7,opt,name=max_open_connections,json=maxOpenConnections,proto3" json:"max_open_connections,omitempty"`
	MaxConnectionLifeTime *types.Duration `protobuf:"bytes,8,opt,name=max_connection_lifeTime,json=maxConnectionLifeTime,proto3" json:"max_connection_lifeTime,omitempty"`
	Charset               string          `protobuf:"bytes,9,opt,name=charset,proto3" json:"charset,omitempty"`
	Collation             string          `protobuf:"bytes,10,opt,name=collation,proto3" json:"collation,omitempty"`
	TablePrefix           string          `protobuf:"bytes,11,opt,name=table_prefix,json=tablePrefix,proto3" json:"table_prefix,omitempty"`
}

func (x *MySQLOptions) Reset() {
	*x = MySQLOptions{}
	if protoimpl.UnsafeEnabled {
		mi := &file_types_gorm_proto_msgTypes[0]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *MySQLOptions) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*MySQLOptions) ProtoMessage() {}

func (x *MySQLOptions) ProtoReflect() protoreflect.Message {
	mi := &file_types_gorm_proto_msgTypes[0]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use MySQLOptions.ProtoReflect.Descriptor instead.
func (*MySQLOptions) Descriptor() ([]byte, []int) {
	return file_types_gorm_proto_rawDescGZIP(), []int{0}
}

func (x *MySQLOptions) GetHost() string {
	if x != nil {
		return x.Host
	}
	return ""
}

func (x *MySQLOptions) GetUsername() string {
	if x != nil {
		return x.Username
	}
	return ""
}

func (x *MySQLOptions) GetPassword() string {
	if x != nil {
		return x.Password
	}
	return ""
}

func (x *MySQLOptions) GetSchema() string {
	if x != nil {
		return x.Schema
	}
	return ""
}

func (x *MySQLOptions) GetMaxIdle() int32 {
	if x != nil {
		return x.MaxIdle
	}
	return 0
}

func (x *MySQLOptions) GetMaxIdleConnections() int32 {
	if x != nil {
		return x.MaxIdleConnections
	}
	return 0
}

func (x *MySQLOptions) GetMaxOpenConnections() int32 {
	if x != nil {
		return x.MaxOpenConnections
	}
	return 0
}

func (x *MySQLOptions) GetMaxConnectionLifeTime() *types.Duration {
	if x != nil {
		return x.MaxConnectionLifeTime
	}
	return nil
}

func (x *MySQLOptions) GetCharset() string {
	if x != nil {
		return x.Charset
	}
	return ""
}

func (x *MySQLOptions) GetCollation() string {
	if x != nil {
		return x.Collation
	}
	return ""
}

func (x *MySQLOptions) GetTablePrefix() string {
	if x != nil {
		return x.TablePrefix
	}
	return ""
}

type SQLiteOptions struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Path string `protobuf:"bytes,1,opt,name=path,proto3" json:"path,omitempty"`
}

func (x *SQLiteOptions) Reset() {
	*x = SQLiteOptions{}
	if protoimpl.UnsafeEnabled {
		mi := &file_types_gorm_proto_msgTypes[1]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *SQLiteOptions) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*SQLiteOptions) ProtoMessage() {}

func (x *SQLiteOptions) ProtoReflect() protoreflect.Message {
	mi := &file_types_gorm_proto_msgTypes[1]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use SQLiteOptions.ProtoReflect.Descriptor instead.
func (*SQLiteOptions) Descriptor() ([]byte, []int) {
	return file_types_gorm_proto_rawDescGZIP(), []int{1}
}

func (x *SQLiteOptions) GetPath() string {
	if x != nil {
		return x.Path
	}
	return ""
}

var File_types_gorm_proto protoreflect.FileDescriptor

var file_types_gorm_proto_rawDesc = []byte{
	0x0a, 0x10, 0x74, 0x79, 0x70, 0x65, 0x73, 0x2f, 0x67, 0x6f, 0x72, 0x6d, 0x2e, 0x70, 0x72, 0x6f,
	0x74, 0x6f, 0x12, 0x10, 0x69, 0x64, 0x61, 0x73, 0x2e, 0x63, 0x6c, 0x69, 0x65, 0x6e, 0x74, 0x2e,
	0x67, 0x6f, 0x72, 0x6d, 0x1a, 0x1e, 0x67, 0x6f, 0x6f, 0x67, 0x6c, 0x65, 0x2f, 0x70, 0x72, 0x6f,
	0x74, 0x6f, 0x62, 0x75, 0x66, 0x2f, 0x64, 0x75, 0x72, 0x61, 0x74, 0x69, 0x6f, 0x6e, 0x2e, 0x70,
	0x72, 0x6f, 0x74, 0x6f, 0x22, 0x9f, 0x03, 0x0a, 0x0c, 0x4d, 0x79, 0x53, 0x51, 0x4c, 0x4f, 0x70,
	0x74, 0x69, 0x6f, 0x6e, 0x73, 0x12, 0x12, 0x0a, 0x04, 0x68, 0x6f, 0x73, 0x74, 0x18, 0x01, 0x20,
	0x01, 0x28, 0x09, 0x52, 0x04, 0x68, 0x6f, 0x73, 0x74, 0x12, 0x1a, 0x0a, 0x08, 0x75, 0x73, 0x65,
	0x72, 0x6e, 0x61, 0x6d, 0x65, 0x18, 0x02, 0x20, 0x01, 0x28, 0x09, 0x52, 0x08, 0x75, 0x73, 0x65,
	0x72, 0x6e, 0x61, 0x6d, 0x65, 0x12, 0x1a, 0x0a, 0x08, 0x70, 0x61, 0x73, 0x73, 0x77, 0x6f, 0x72,
	0x64, 0x18, 0x03, 0x20, 0x01, 0x28, 0x09, 0x52, 0x08, 0x70, 0x61, 0x73, 0x73, 0x77, 0x6f, 0x72,
	0x64, 0x12, 0x16, 0x0a, 0x06, 0x73, 0x63, 0x68, 0x65, 0x6d, 0x61, 0x18, 0x04, 0x20, 0x01, 0x28,
	0x09, 0x52, 0x06, 0x73, 0x63, 0x68, 0x65, 0x6d, 0x61, 0x12, 0x19, 0x0a, 0x08, 0x6d, 0x61, 0x78,
	0x5f, 0x69, 0x64, 0x6c, 0x65, 0x18, 0x05, 0x20, 0x01, 0x28, 0x05, 0x52, 0x07, 0x6d, 0x61, 0x78,
	0x49, 0x64, 0x6c, 0x65, 0x12, 0x30, 0x0a, 0x14, 0x6d, 0x61, 0x78, 0x5f, 0x69, 0x64, 0x6c, 0x65,
	0x5f, 0x63, 0x6f, 0x6e, 0x6e, 0x65, 0x63, 0x74, 0x69, 0x6f, 0x6e, 0x73, 0x18, 0x06, 0x20, 0x01,
	0x28, 0x05, 0x52, 0x12, 0x6d, 0x61, 0x78, 0x49, 0x64, 0x6c, 0x65, 0x43, 0x6f, 0x6e, 0x6e, 0x65,
	0x63, 0x74, 0x69, 0x6f, 0x6e, 0x73, 0x12, 0x30, 0x0a, 0x14, 0x6d, 0x61, 0x78, 0x5f, 0x6f, 0x70,
	0x65, 0x6e, 0x5f, 0x63, 0x6f, 0x6e, 0x6e, 0x65, 0x63, 0x74, 0x69, 0x6f, 0x6e, 0x73, 0x18, 0x07,
	0x20, 0x01, 0x28, 0x05, 0x52, 0x12, 0x6d, 0x61, 0x78, 0x4f, 0x70, 0x65, 0x6e, 0x43, 0x6f, 0x6e,
	0x6e, 0x65, 0x63, 0x74, 0x69, 0x6f, 0x6e, 0x73, 0x12, 0x51, 0x0a, 0x17, 0x6d, 0x61, 0x78, 0x5f,
	0x63, 0x6f, 0x6e, 0x6e, 0x65, 0x63, 0x74, 0x69, 0x6f, 0x6e, 0x5f, 0x6c, 0x69, 0x66, 0x65, 0x54,
	0x69, 0x6d, 0x65, 0x18, 0x08, 0x20, 0x01, 0x28, 0x0b, 0x32, 0x19, 0x2e, 0x67, 0x6f, 0x6f, 0x67,
	0x6c, 0x65, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x62, 0x75, 0x66, 0x2e, 0x44, 0x75, 0x72, 0x61,
	0x74, 0x69, 0x6f, 0x6e, 0x52, 0x15, 0x6d, 0x61, 0x78, 0x43, 0x6f, 0x6e, 0x6e, 0x65, 0x63, 0x74,
	0x69, 0x6f, 0x6e, 0x4c, 0x69, 0x66, 0x65, 0x54, 0x69, 0x6d, 0x65, 0x12, 0x18, 0x0a, 0x07, 0x63,
	0x68, 0x61, 0x72, 0x73, 0x65, 0x74, 0x18, 0x09, 0x20, 0x01, 0x28, 0x09, 0x52, 0x07, 0x63, 0x68,
	0x61, 0x72, 0x73, 0x65, 0x74, 0x12, 0x1c, 0x0a, 0x09, 0x63, 0x6f, 0x6c, 0x6c, 0x61, 0x74, 0x69,
	0x6f, 0x6e, 0x18, 0x0a, 0x20, 0x01, 0x28, 0x09, 0x52, 0x09, 0x63, 0x6f, 0x6c, 0x6c, 0x61, 0x74,
	0x69, 0x6f, 0x6e, 0x12, 0x21, 0x0a, 0x0c, 0x74, 0x61, 0x62, 0x6c, 0x65, 0x5f, 0x70, 0x72, 0x65,
	0x66, 0x69, 0x78, 0x18, 0x0b, 0x20, 0x01, 0x28, 0x09, 0x52, 0x0b, 0x74, 0x61, 0x62, 0x6c, 0x65,
	0x50, 0x72, 0x65, 0x66, 0x69, 0x78, 0x22, 0x23, 0x0a, 0x0d, 0x53, 0x51, 0x4c, 0x69, 0x74, 0x65,
	0x4f, 0x70, 0x74, 0x69, 0x6f, 0x6e, 0x73, 0x12, 0x12, 0x0a, 0x04, 0x70, 0x61, 0x74, 0x68, 0x18,
	0x01, 0x20, 0x01, 0x28, 0x09, 0x52, 0x04, 0x70, 0x61, 0x74, 0x68, 0x42, 0x1b, 0x5a, 0x19, 0x69,
	0x64, 0x61, 0x73, 0x2f, 0x70, 0x6b, 0x67, 0x2f, 0x63, 0x6c, 0x69, 0x65, 0x6e, 0x74, 0x2f, 0x67,
	0x6f, 0x72, 0x6d, 0x3b, 0x67, 0x6f, 0x72, 0x6d, 0x50, 0x00, 0x62, 0x06, 0x70, 0x72, 0x6f, 0x74,
	0x6f, 0x33,
}

var (
	file_types_gorm_proto_rawDescOnce sync.Once
	file_types_gorm_proto_rawDescData = file_types_gorm_proto_rawDesc
)

func file_types_gorm_proto_rawDescGZIP() []byte {
	file_types_gorm_proto_rawDescOnce.Do(func() {
		file_types_gorm_proto_rawDescData = protoimpl.X.CompressGZIP(file_types_gorm_proto_rawDescData)
	})
	return file_types_gorm_proto_rawDescData
}

var file_types_gorm_proto_msgTypes = make([]protoimpl.MessageInfo, 2)
var file_types_gorm_proto_goTypes = []interface{}{
	(*MySQLOptions)(nil),   // 0: idas.client.gorm.MySQLOptions
	(*SQLiteOptions)(nil),  // 1: idas.client.gorm.SQLiteOptions
	(*types.Duration)(nil), // 2: google.protobuf.Duration
}
var file_types_gorm_proto_depIdxs = []int32{
	2, // 0: idas.client.gorm.MySQLOptions.max_connection_lifeTime:type_name -> google.protobuf.Duration
	1, // [1:1] is the sub-list for method output_type
	1, // [1:1] is the sub-list for method input_type
	1, // [1:1] is the sub-list for extension type_name
	1, // [1:1] is the sub-list for extension extendee
	0, // [0:1] is the sub-list for field type_name
}

func init() { file_types_gorm_proto_init() }
func file_types_gorm_proto_init() {
	if File_types_gorm_proto != nil {
		return
	}
	if !protoimpl.UnsafeEnabled {
		file_types_gorm_proto_msgTypes[0].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*MySQLOptions); i {
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
		file_types_gorm_proto_msgTypes[1].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*SQLiteOptions); i {
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
			RawDescriptor: file_types_gorm_proto_rawDesc,
			NumEnums:      0,
			NumMessages:   2,
			NumExtensions: 0,
			NumServices:   0,
		},
		GoTypes:           file_types_gorm_proto_goTypes,
		DependencyIndexes: file_types_gorm_proto_depIdxs,
		MessageInfos:      file_types_gorm_proto_msgTypes,
	}.Build()
	File_types_gorm_proto = out.File
	file_types_gorm_proto_rawDesc = nil
	file_types_gorm_proto_goTypes = nil
	file_types_gorm_proto_depIdxs = nil
}
