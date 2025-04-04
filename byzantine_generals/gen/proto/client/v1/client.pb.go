// Code generated by protoc-gen-go. DO NOT EDIT.
// versions:
// 	protoc-gen-go v1.36.5
// 	protoc        (unknown)
// source: proto/client/v1/client.proto

package clientv1

import (
	v1 "github.com/nicholasjackson/demo-lamport/byzantine_generals/gen/proto/common/v1"
	protoreflect "google.golang.org/protobuf/reflect/protoreflect"
	protoimpl "google.golang.org/protobuf/runtime/protoimpl"
	reflect "reflect"
	sync "sync"
	unsafe "unsafe"
)

const (
	// Verify that this generated code is sufficiently up-to-date.
	_ = protoimpl.EnforceVersion(20 - protoimpl.MinVersion)
	// Verify that runtime/protoimpl is sufficiently up-to-date.
	_ = protoimpl.EnforceVersion(protoimpl.MaxVersion - 20)
)

type ReceiveCommandRequest struct {
	state         protoimpl.MessageState `protogen:"open.v1"`
	Command       *v1.Command            `protobuf:"bytes,1,opt,name=command,proto3" json:"command,omitempty"`
	unknownFields protoimpl.UnknownFields
	sizeCache     protoimpl.SizeCache
}

func (x *ReceiveCommandRequest) Reset() {
	*x = ReceiveCommandRequest{}
	mi := &file_proto_client_v1_client_proto_msgTypes[0]
	ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
	ms.StoreMessageInfo(mi)
}

func (x *ReceiveCommandRequest) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*ReceiveCommandRequest) ProtoMessage() {}

func (x *ReceiveCommandRequest) ProtoReflect() protoreflect.Message {
	mi := &file_proto_client_v1_client_proto_msgTypes[0]
	if x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use ReceiveCommandRequest.ProtoReflect.Descriptor instead.
func (*ReceiveCommandRequest) Descriptor() ([]byte, []int) {
	return file_proto_client_v1_client_proto_rawDescGZIP(), []int{0}
}

func (x *ReceiveCommandRequest) GetCommand() *v1.Command {
	if x != nil {
		return x.Command
	}
	return nil
}

var File_proto_client_v1_client_proto protoreflect.FileDescriptor

var file_proto_client_v1_client_proto_rawDesc = string([]byte{
	0x0a, 0x1c, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x2f, 0x63, 0x6c, 0x69, 0x65, 0x6e, 0x74, 0x2f, 0x76,
	0x31, 0x2f, 0x63, 0x6c, 0x69, 0x65, 0x6e, 0x74, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x12, 0x0f,
	0x70, 0x72, 0x6f, 0x74, 0x6f, 0x2e, 0x63, 0x6f, 0x6d, 0x6d, 0x6f, 0x6e, 0x2e, 0x76, 0x31, 0x1a,
	0x1c, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x2f, 0x63, 0x6f, 0x6d, 0x6d, 0x6f, 0x6e, 0x2f, 0x76, 0x31,
	0x2f, 0x63, 0x6f, 0x6d, 0x6d, 0x6f, 0x6e, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x22, 0x4b, 0x0a,
	0x15, 0x52, 0x65, 0x63, 0x65, 0x69, 0x76, 0x65, 0x43, 0x6f, 0x6d, 0x6d, 0x61, 0x6e, 0x64, 0x52,
	0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x12, 0x32, 0x0a, 0x07, 0x63, 0x6f, 0x6d, 0x6d, 0x61, 0x6e,
	0x64, 0x18, 0x01, 0x20, 0x01, 0x28, 0x0b, 0x32, 0x18, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x2e,
	0x63, 0x6f, 0x6d, 0x6d, 0x6f, 0x6e, 0x2e, 0x76, 0x31, 0x2e, 0x43, 0x6f, 0x6d, 0x6d, 0x61, 0x6e,
	0x64, 0x52, 0x07, 0x63, 0x6f, 0x6d, 0x6d, 0x61, 0x6e, 0x64, 0x32, 0xb7, 0x01, 0x0a, 0x0f, 0x47,
	0x65, 0x6e, 0x65, 0x72, 0x61, 0x6c, 0x73, 0x53, 0x65, 0x72, 0x76, 0x69, 0x63, 0x65, 0x12, 0x48,
	0x0a, 0x05, 0x52, 0x65, 0x73, 0x65, 0x74, 0x12, 0x1d, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x2e,
	0x63, 0x6f, 0x6d, 0x6d, 0x6f, 0x6e, 0x2e, 0x76, 0x31, 0x2e, 0x45, 0x6d, 0x70, 0x74, 0x79, 0x52,
	0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x1a, 0x1e, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x2e, 0x63,
	0x6f, 0x6d, 0x6d, 0x6f, 0x6e, 0x2e, 0x76, 0x31, 0x2e, 0x45, 0x6d, 0x70, 0x74, 0x79, 0x52, 0x65,
	0x73, 0x70, 0x6f, 0x6e, 0x73, 0x65, 0x22, 0x00, 0x12, 0x5a, 0x0a, 0x0e, 0x52, 0x65, 0x63, 0x65,
	0x69, 0x76, 0x65, 0x43, 0x6f, 0x6d, 0x6d, 0x61, 0x6e, 0x64, 0x12, 0x26, 0x2e, 0x70, 0x72, 0x6f,
	0x74, 0x6f, 0x2e, 0x63, 0x6f, 0x6d, 0x6d, 0x6f, 0x6e, 0x2e, 0x76, 0x31, 0x2e, 0x52, 0x65, 0x63,
	0x65, 0x69, 0x76, 0x65, 0x43, 0x6f, 0x6d, 0x6d, 0x61, 0x6e, 0x64, 0x52, 0x65, 0x71, 0x75, 0x65,
	0x73, 0x74, 0x1a, 0x1e, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x2e, 0x63, 0x6f, 0x6d, 0x6d, 0x6f,
	0x6e, 0x2e, 0x76, 0x31, 0x2e, 0x45, 0x6d, 0x70, 0x74, 0x79, 0x52, 0x65, 0x73, 0x70, 0x6f, 0x6e,
	0x73, 0x65, 0x22, 0x00, 0x42, 0x59, 0x5a, 0x57, 0x67, 0x69, 0x74, 0x68, 0x75, 0x62, 0x2e, 0x63,
	0x6f, 0x6d, 0x2f, 0x6e, 0x69, 0x63, 0x68, 0x6f, 0x6c, 0x61, 0x73, 0x6a, 0x61, 0x63, 0x6b, 0x73,
	0x6f, 0x6e, 0x2f, 0x64, 0x65, 0x6d, 0x6f, 0x2d, 0x6c, 0x61, 0x6d, 0x70, 0x6f, 0x72, 0x74, 0x2f,
	0x62, 0x79, 0x7a, 0x61, 0x6e, 0x74, 0x69, 0x6e, 0x65, 0x5f, 0x67, 0x65, 0x6e, 0x65, 0x72, 0x61,
	0x6c, 0x73, 0x2f, 0x67, 0x65, 0x6e, 0x2f, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x2f, 0x63, 0x6c, 0x69,
	0x65, 0x6e, 0x74, 0x2f, 0x76, 0x31, 0x3b, 0x63, 0x6c, 0x69, 0x65, 0x6e, 0x74, 0x76, 0x31, 0x62,
	0x06, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x33,
})

var (
	file_proto_client_v1_client_proto_rawDescOnce sync.Once
	file_proto_client_v1_client_proto_rawDescData []byte
)

func file_proto_client_v1_client_proto_rawDescGZIP() []byte {
	file_proto_client_v1_client_proto_rawDescOnce.Do(func() {
		file_proto_client_v1_client_proto_rawDescData = protoimpl.X.CompressGZIP(unsafe.Slice(unsafe.StringData(file_proto_client_v1_client_proto_rawDesc), len(file_proto_client_v1_client_proto_rawDesc)))
	})
	return file_proto_client_v1_client_proto_rawDescData
}

var file_proto_client_v1_client_proto_msgTypes = make([]protoimpl.MessageInfo, 1)
var file_proto_client_v1_client_proto_goTypes = []any{
	(*ReceiveCommandRequest)(nil), // 0: proto.common.v1.ReceiveCommandRequest
	(*v1.Command)(nil),            // 1: proto.common.v1.Command
	(*v1.EmptyRequest)(nil),       // 2: proto.common.v1.EmptyRequest
	(*v1.EmptyResponse)(nil),      // 3: proto.common.v1.EmptyResponse
}
var file_proto_client_v1_client_proto_depIdxs = []int32{
	1, // 0: proto.common.v1.ReceiveCommandRequest.command:type_name -> proto.common.v1.Command
	2, // 1: proto.common.v1.GeneralsService.Reset:input_type -> proto.common.v1.EmptyRequest
	0, // 2: proto.common.v1.GeneralsService.ReceiveCommand:input_type -> proto.common.v1.ReceiveCommandRequest
	3, // 3: proto.common.v1.GeneralsService.Reset:output_type -> proto.common.v1.EmptyResponse
	3, // 4: proto.common.v1.GeneralsService.ReceiveCommand:output_type -> proto.common.v1.EmptyResponse
	3, // [3:5] is the sub-list for method output_type
	1, // [1:3] is the sub-list for method input_type
	1, // [1:1] is the sub-list for extension type_name
	1, // [1:1] is the sub-list for extension extendee
	0, // [0:1] is the sub-list for field type_name
}

func init() { file_proto_client_v1_client_proto_init() }
func file_proto_client_v1_client_proto_init() {
	if File_proto_client_v1_client_proto != nil {
		return
	}
	type x struct{}
	out := protoimpl.TypeBuilder{
		File: protoimpl.DescBuilder{
			GoPackagePath: reflect.TypeOf(x{}).PkgPath(),
			RawDescriptor: unsafe.Slice(unsafe.StringData(file_proto_client_v1_client_proto_rawDesc), len(file_proto_client_v1_client_proto_rawDesc)),
			NumEnums:      0,
			NumMessages:   1,
			NumExtensions: 0,
			NumServices:   1,
		},
		GoTypes:           file_proto_client_v1_client_proto_goTypes,
		DependencyIndexes: file_proto_client_v1_client_proto_depIdxs,
		MessageInfos:      file_proto_client_v1_client_proto_msgTypes,
	}.Build()
	File_proto_client_v1_client_proto = out.File
	file_proto_client_v1_client_proto_goTypes = nil
	file_proto_client_v1_client_proto_depIdxs = nil
}
