// Code generated by protoc-gen-go. DO NOT EDIT.
// versions:
// 	protoc-gen-go v1.27.1
// 	protoc        v3.19.4
// source: fs.proto

package fs

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

// Symbols defined in public import of google/protobuf/any.proto.

type Any = types.Any

type FsPath struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	FsPath string     `protobuf:"bytes,1,opt,name=fs_path,json=fsPath,proto3" json:"fs_path,omitempty"`
	Fs     *types.Any `protobuf:"bytes,2,opt,name=fs,proto3" json:"fs,omitempty"`
}

func (x *FsPath) Reset() {
	*x = FsPath{}
	if protoimpl.UnsafeEnabled {
		mi := &file_fs_proto_msgTypes[0]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *FsPath) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*FsPath) ProtoMessage() {}

func (x *FsPath) ProtoReflect() protoreflect.Message {
	mi := &file_fs_proto_msgTypes[0]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use FsPath.ProtoReflect.Descriptor instead.
func (*FsPath) Descriptor() ([]byte, []int) {
	return file_fs_proto_rawDescGZIP(), []int{0}
}

func (x *FsPath) GetFsPath() string {
	if x != nil {
		return x.FsPath
	}
	return ""
}

func (x *FsPath) GetFs() *types.Any {
	if x != nil {
		return x.Fs
	}
	return nil
}

var File_fs_proto protoreflect.FileDescriptor

var file_fs_proto_rawDesc = []byte{
	0x0a, 0x08, 0x66, 0x73, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x12, 0x0e, 0x69, 0x64, 0x61, 0x73,
	0x2e, 0x70, 0x6b, 0x67, 0x2e, 0x75, 0x74, 0x69, 0x6c, 0x73, 0x1a, 0x19, 0x67, 0x6f, 0x6f, 0x67,
	0x6c, 0x65, 0x2f, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x62, 0x75, 0x66, 0x2f, 0x61, 0x6e, 0x79, 0x2e,
	0x70, 0x72, 0x6f, 0x74, 0x6f, 0x22, 0x47, 0x0a, 0x06, 0x46, 0x73, 0x50, 0x61, 0x74, 0x68, 0x12,
	0x17, 0x0a, 0x07, 0x66, 0x73, 0x5f, 0x70, 0x61, 0x74, 0x68, 0x18, 0x01, 0x20, 0x01, 0x28, 0x09,
	0x52, 0x06, 0x66, 0x73, 0x50, 0x61, 0x74, 0x68, 0x12, 0x24, 0x0a, 0x02, 0x66, 0x73, 0x18, 0x02,
	0x20, 0x01, 0x28, 0x0b, 0x32, 0x14, 0x2e, 0x67, 0x6f, 0x6f, 0x67, 0x6c, 0x65, 0x2e, 0x70, 0x72,
	0x6f, 0x74, 0x6f, 0x62, 0x75, 0x66, 0x2e, 0x41, 0x6e, 0x79, 0x52, 0x02, 0x66, 0x73, 0x42, 0x16,
	0x5a, 0x14, 0x69, 0x64, 0x61, 0x73, 0x2f, 0x70, 0x6b, 0x67, 0x2f, 0x75, 0x74, 0x69, 0x6c, 0x73,
	0x2f, 0x66, 0x73, 0x3b, 0x66, 0x73, 0x50, 0x00, 0x62, 0x06, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x33,
}

var (
	file_fs_proto_rawDescOnce sync.Once
	file_fs_proto_rawDescData = file_fs_proto_rawDesc
)

func file_fs_proto_rawDescGZIP() []byte {
	file_fs_proto_rawDescOnce.Do(func() {
		file_fs_proto_rawDescData = protoimpl.X.CompressGZIP(file_fs_proto_rawDescData)
	})
	return file_fs_proto_rawDescData
}

var file_fs_proto_msgTypes = make([]protoimpl.MessageInfo, 1)
var file_fs_proto_goTypes = []interface{}{
	(*FsPath)(nil),    // 0: idas.pkg.utils.FsPath
	(*types.Any)(nil), // 1: google.protobuf.Any
}
var file_fs_proto_depIdxs = []int32{
	1, // 0: idas.pkg.utils.FsPath.fs:type_name -> google.protobuf.Any
	1, // [1:1] is the sub-list for method output_type
	1, // [1:1] is the sub-list for method input_type
	1, // [1:1] is the sub-list for extension type_name
	1, // [1:1] is the sub-list for extension extendee
	0, // [0:1] is the sub-list for field type_name
}

func init() { file_fs_proto_init() }
func file_fs_proto_init() {
	if File_fs_proto != nil {
		return
	}
	if !protoimpl.UnsafeEnabled {
		file_fs_proto_msgTypes[0].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*FsPath); i {
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
			RawDescriptor: file_fs_proto_rawDesc,
			NumEnums:      0,
			NumMessages:   1,
			NumExtensions: 0,
			NumServices:   0,
		},
		GoTypes:           file_fs_proto_goTypes,
		DependencyIndexes: file_fs_proto_depIdxs,
		MessageInfos:      file_fs_proto_msgTypes,
	}.Build()
	File_fs_proto = out.File
	file_fs_proto_rawDesc = nil
	file_fs_proto_goTypes = nil
	file_fs_proto_depIdxs = nil
}