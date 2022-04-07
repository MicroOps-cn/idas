// Code generated by protoc-gen-go. DO NOT EDIT.
// versions:
// 	protoc-gen-go v1.27.1
// 	protoc        v3.19.4
// source: pkg/utils/capacity/capacity.proto

package capacity

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

type Capacity struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Capacity int64 `protobuf:"varint,1,opt,name=capacity,proto3" json:"capacity,omitempty"`
}

func (x *Capacity) Reset() {
	*x = Capacity{}
	if protoimpl.UnsafeEnabled {
		mi := &file_pkg_utils_capacity_capacity_proto_msgTypes[0]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *Capacity) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*Capacity) ProtoMessage() {}

func (x *Capacity) ProtoReflect() protoreflect.Message {
	mi := &file_pkg_utils_capacity_capacity_proto_msgTypes[0]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use Capacity.ProtoReflect.Descriptor instead.
func (*Capacity) Descriptor() ([]byte, []int) {
	return file_pkg_utils_capacity_capacity_proto_rawDescGZIP(), []int{0}
}

func (x *Capacity) GetCapacity() int64 {
	if x != nil {
		return x.Capacity
	}
	return 0
}

var File_pkg_utils_capacity_capacity_proto protoreflect.FileDescriptor

var file_pkg_utils_capacity_capacity_proto_rawDesc = []byte{
	0x0a, 0x21, 0x70, 0x6b, 0x67, 0x2f, 0x75, 0x74, 0x69, 0x6c, 0x73, 0x2f, 0x63, 0x61, 0x70, 0x61,
	0x63, 0x69, 0x74, 0x79, 0x2f, 0x63, 0x61, 0x70, 0x61, 0x63, 0x69, 0x74, 0x79, 0x2e, 0x70, 0x72,
	0x6f, 0x74, 0x6f, 0x12, 0x0e, 0x69, 0x64, 0x61, 0x73, 0x2e, 0x70, 0x6b, 0x67, 0x2e, 0x75, 0x74,
	0x69, 0x6c, 0x73, 0x22, 0x26, 0x0a, 0x08, 0x43, 0x61, 0x70, 0x61, 0x63, 0x69, 0x74, 0x79, 0x12,
	0x1a, 0x0a, 0x08, 0x63, 0x61, 0x70, 0x61, 0x63, 0x69, 0x74, 0x79, 0x18, 0x01, 0x20, 0x01, 0x28,
	0x03, 0x52, 0x08, 0x63, 0x61, 0x70, 0x61, 0x63, 0x69, 0x74, 0x79, 0x42, 0x22, 0x5a, 0x20, 0x69,
	0x64, 0x61, 0x73, 0x2f, 0x70, 0x6b, 0x67, 0x2f, 0x75, 0x74, 0x69, 0x6c, 0x73, 0x2f, 0x63, 0x61,
	0x70, 0x61, 0x63, 0x69, 0x74, 0x79, 0x3b, 0x63, 0x61, 0x70, 0x61, 0x63, 0x69, 0x74, 0x79, 0x62,
	0x06, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x33,
}

var (
	file_pkg_utils_capacity_capacity_proto_rawDescOnce sync.Once
	file_pkg_utils_capacity_capacity_proto_rawDescData = file_pkg_utils_capacity_capacity_proto_rawDesc
)

func file_pkg_utils_capacity_capacity_proto_rawDescGZIP() []byte {
	file_pkg_utils_capacity_capacity_proto_rawDescOnce.Do(func() {
		file_pkg_utils_capacity_capacity_proto_rawDescData = protoimpl.X.CompressGZIP(file_pkg_utils_capacity_capacity_proto_rawDescData)
	})
	return file_pkg_utils_capacity_capacity_proto_rawDescData
}

var file_pkg_utils_capacity_capacity_proto_msgTypes = make([]protoimpl.MessageInfo, 1)
var file_pkg_utils_capacity_capacity_proto_goTypes = []interface{}{
	(*Capacity)(nil), // 0: idas.pkg.utils.Capacity
}
var file_pkg_utils_capacity_capacity_proto_depIdxs = []int32{
	0, // [0:0] is the sub-list for method output_type
	0, // [0:0] is the sub-list for method input_type
	0, // [0:0] is the sub-list for extension type_name
	0, // [0:0] is the sub-list for extension extendee
	0, // [0:0] is the sub-list for field type_name
}

func init() { file_pkg_utils_capacity_capacity_proto_init() }
func file_pkg_utils_capacity_capacity_proto_init() {
	if File_pkg_utils_capacity_capacity_proto != nil {
		return
	}
	if !protoimpl.UnsafeEnabled {
		file_pkg_utils_capacity_capacity_proto_msgTypes[0].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*Capacity); i {
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
			RawDescriptor: file_pkg_utils_capacity_capacity_proto_rawDesc,
			NumEnums:      0,
			NumMessages:   1,
			NumExtensions: 0,
			NumServices:   0,
		},
		GoTypes:           file_pkg_utils_capacity_capacity_proto_goTypes,
		DependencyIndexes: file_pkg_utils_capacity_capacity_proto_depIdxs,
		MessageInfos:      file_pkg_utils_capacity_capacity_proto_msgTypes,
	}.Build()
	File_pkg_utils_capacity_capacity_proto = out.File
	file_pkg_utils_capacity_capacity_proto_rawDesc = nil
	file_pkg_utils_capacity_capacity_proto_goTypes = nil
	file_pkg_utils_capacity_capacity_proto_depIdxs = nil
}