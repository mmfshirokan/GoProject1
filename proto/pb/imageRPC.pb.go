// Code generated by protoc-gen-go. DO NOT EDIT.
// versions:
// 	protoc-gen-go v1.28.1
// 	protoc        v4.25.1
// source: imageRPC.proto

package pb

import (
	protoreflect "google.golang.org/protobuf/reflect/protoreflect"
	protoimpl "google.golang.org/protobuf/runtime/protoimpl"
	emptypb "google.golang.org/protobuf/types/known/emptypb"
	reflect "reflect"
	sync "sync"
)

const (
	// Verify that this generated code is sufficiently up-to-date.
	_ = protoimpl.EnforceVersion(20 - protoimpl.MinVersion)
	// Verify that runtime/protoimpl is sufficiently up-to-date.
	_ = protoimpl.EnforceVersion(protoimpl.MaxVersion - 20)
)

type RequestDownloadImage struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	AuthToken string `protobuf:"bytes,1,opt,name=authToken,proto3" json:"authToken,omitempty"`
	UserID    int64  `protobuf:"varint,2,opt,name=userID,proto3" json:"userID,omitempty"`
	ImageName string `protobuf:"bytes,3,opt,name=imageName,proto3" json:"imageName,omitempty"`
}

func (x *RequestDownloadImage) Reset() {
	*x = RequestDownloadImage{}
	if protoimpl.UnsafeEnabled {
		mi := &file_imageRPC_proto_msgTypes[0]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *RequestDownloadImage) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*RequestDownloadImage) ProtoMessage() {}

func (x *RequestDownloadImage) ProtoReflect() protoreflect.Message {
	mi := &file_imageRPC_proto_msgTypes[0]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use RequestDownloadImage.ProtoReflect.Descriptor instead.
func (*RequestDownloadImage) Descriptor() ([]byte, []int) {
	return file_imageRPC_proto_rawDescGZIP(), []int{0}
}

func (x *RequestDownloadImage) GetAuthToken() string {
	if x != nil {
		return x.AuthToken
	}
	return ""
}

func (x *RequestDownloadImage) GetUserID() int64 {
	if x != nil {
		return x.UserID
	}
	return 0
}

func (x *RequestDownloadImage) GetImageName() string {
	if x != nil {
		return x.ImageName
	}
	return ""
}

type ResponseDownloadImage struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	ImagePiece       []byte `protobuf:"bytes,1,opt,name=imagePiece,proto3" json:"imagePiece,omitempty"`
	StreamIsFinished bool   `protobuf:"varint,2,opt,name=streamIsFinished,proto3" json:"streamIsFinished,omitempty"`
}

func (x *ResponseDownloadImage) Reset() {
	*x = ResponseDownloadImage{}
	if protoimpl.UnsafeEnabled {
		mi := &file_imageRPC_proto_msgTypes[1]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *ResponseDownloadImage) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*ResponseDownloadImage) ProtoMessage() {}

func (x *ResponseDownloadImage) ProtoReflect() protoreflect.Message {
	mi := &file_imageRPC_proto_msgTypes[1]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use ResponseDownloadImage.ProtoReflect.Descriptor instead.
func (*ResponseDownloadImage) Descriptor() ([]byte, []int) {
	return file_imageRPC_proto_rawDescGZIP(), []int{1}
}

func (x *ResponseDownloadImage) GetImagePiece() []byte {
	if x != nil {
		return x.ImagePiece
	}
	return nil
}

func (x *ResponseDownloadImage) GetStreamIsFinished() bool {
	if x != nil {
		return x.StreamIsFinished
	}
	return false
}

type RequestUploadImage struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	AuthToken        *string `protobuf:"bytes,1,opt,name=authToken,proto3,oneof" json:"authToken,omitempty"`
	UserID           *int64  `protobuf:"varint,2,opt,name=userID,proto3,oneof" json:"userID,omitempty"`
	ImageName        *string `protobuf:"bytes,3,opt,name=imageName,proto3,oneof" json:"imageName,omitempty"`
	ImagePiece       []byte  `protobuf:"bytes,4,opt,name=imagePiece,proto3" json:"imagePiece,omitempty"`
	StreamIsFinished bool    `protobuf:"varint,5,opt,name=streamIsFinished,proto3" json:"streamIsFinished,omitempty"`
}

func (x *RequestUploadImage) Reset() {
	*x = RequestUploadImage{}
	if protoimpl.UnsafeEnabled {
		mi := &file_imageRPC_proto_msgTypes[2]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *RequestUploadImage) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*RequestUploadImage) ProtoMessage() {}

func (x *RequestUploadImage) ProtoReflect() protoreflect.Message {
	mi := &file_imageRPC_proto_msgTypes[2]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use RequestUploadImage.ProtoReflect.Descriptor instead.
func (*RequestUploadImage) Descriptor() ([]byte, []int) {
	return file_imageRPC_proto_rawDescGZIP(), []int{2}
}

func (x *RequestUploadImage) GetAuthToken() string {
	if x != nil && x.AuthToken != nil {
		return *x.AuthToken
	}
	return ""
}

func (x *RequestUploadImage) GetUserID() int64 {
	if x != nil && x.UserID != nil {
		return *x.UserID
	}
	return 0
}

func (x *RequestUploadImage) GetImageName() string {
	if x != nil && x.ImageName != nil {
		return *x.ImageName
	}
	return ""
}

func (x *RequestUploadImage) GetImagePiece() []byte {
	if x != nil {
		return x.ImagePiece
	}
	return nil
}

func (x *RequestUploadImage) GetStreamIsFinished() bool {
	if x != nil {
		return x.StreamIsFinished
	}
	return false
}

var File_imageRPC_proto protoreflect.FileDescriptor

var file_imageRPC_proto_rawDesc = []byte{
	0x0a, 0x0e, 0x69, 0x6d, 0x61, 0x67, 0x65, 0x52, 0x50, 0x43, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f,
	0x12, 0x02, 0x70, 0x62, 0x1a, 0x1b, 0x67, 0x6f, 0x6f, 0x67, 0x6c, 0x65, 0x2f, 0x70, 0x72, 0x6f,
	0x74, 0x6f, 0x62, 0x75, 0x66, 0x2f, 0x65, 0x6d, 0x70, 0x74, 0x79, 0x2e, 0x70, 0x72, 0x6f, 0x74,
	0x6f, 0x22, 0x6a, 0x0a, 0x14, 0x72, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x44, 0x6f, 0x77, 0x6e,
	0x6c, 0x6f, 0x61, 0x64, 0x49, 0x6d, 0x61, 0x67, 0x65, 0x12, 0x1c, 0x0a, 0x09, 0x61, 0x75, 0x74,
	0x68, 0x54, 0x6f, 0x6b, 0x65, 0x6e, 0x18, 0x01, 0x20, 0x01, 0x28, 0x09, 0x52, 0x09, 0x61, 0x75,
	0x74, 0x68, 0x54, 0x6f, 0x6b, 0x65, 0x6e, 0x12, 0x16, 0x0a, 0x06, 0x75, 0x73, 0x65, 0x72, 0x49,
	0x44, 0x18, 0x02, 0x20, 0x01, 0x28, 0x03, 0x52, 0x06, 0x75, 0x73, 0x65, 0x72, 0x49, 0x44, 0x12,
	0x1c, 0x0a, 0x09, 0x69, 0x6d, 0x61, 0x67, 0x65, 0x4e, 0x61, 0x6d, 0x65, 0x18, 0x03, 0x20, 0x01,
	0x28, 0x09, 0x52, 0x09, 0x69, 0x6d, 0x61, 0x67, 0x65, 0x4e, 0x61, 0x6d, 0x65, 0x22, 0x63, 0x0a,
	0x15, 0x72, 0x65, 0x73, 0x70, 0x6f, 0x6e, 0x73, 0x65, 0x44, 0x6f, 0x77, 0x6e, 0x6c, 0x6f, 0x61,
	0x64, 0x49, 0x6d, 0x61, 0x67, 0x65, 0x12, 0x1e, 0x0a, 0x0a, 0x69, 0x6d, 0x61, 0x67, 0x65, 0x50,
	0x69, 0x65, 0x63, 0x65, 0x18, 0x01, 0x20, 0x01, 0x28, 0x0c, 0x52, 0x0a, 0x69, 0x6d, 0x61, 0x67,
	0x65, 0x50, 0x69, 0x65, 0x63, 0x65, 0x12, 0x2a, 0x0a, 0x10, 0x73, 0x74, 0x72, 0x65, 0x61, 0x6d,
	0x49, 0x73, 0x46, 0x69, 0x6e, 0x69, 0x73, 0x68, 0x65, 0x64, 0x18, 0x02, 0x20, 0x01, 0x28, 0x08,
	0x52, 0x10, 0x73, 0x74, 0x72, 0x65, 0x61, 0x6d, 0x49, 0x73, 0x46, 0x69, 0x6e, 0x69, 0x73, 0x68,
	0x65, 0x64, 0x22, 0xea, 0x01, 0x0a, 0x12, 0x72, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x55, 0x70,
	0x6c, 0x6f, 0x61, 0x64, 0x49, 0x6d, 0x61, 0x67, 0x65, 0x12, 0x21, 0x0a, 0x09, 0x61, 0x75, 0x74,
	0x68, 0x54, 0x6f, 0x6b, 0x65, 0x6e, 0x18, 0x01, 0x20, 0x01, 0x28, 0x09, 0x48, 0x00, 0x52, 0x09,
	0x61, 0x75, 0x74, 0x68, 0x54, 0x6f, 0x6b, 0x65, 0x6e, 0x88, 0x01, 0x01, 0x12, 0x1b, 0x0a, 0x06,
	0x75, 0x73, 0x65, 0x72, 0x49, 0x44, 0x18, 0x02, 0x20, 0x01, 0x28, 0x03, 0x48, 0x01, 0x52, 0x06,
	0x75, 0x73, 0x65, 0x72, 0x49, 0x44, 0x88, 0x01, 0x01, 0x12, 0x21, 0x0a, 0x09, 0x69, 0x6d, 0x61,
	0x67, 0x65, 0x4e, 0x61, 0x6d, 0x65, 0x18, 0x03, 0x20, 0x01, 0x28, 0x09, 0x48, 0x02, 0x52, 0x09,
	0x69, 0x6d, 0x61, 0x67, 0x65, 0x4e, 0x61, 0x6d, 0x65, 0x88, 0x01, 0x01, 0x12, 0x1e, 0x0a, 0x0a,
	0x69, 0x6d, 0x61, 0x67, 0x65, 0x50, 0x69, 0x65, 0x63, 0x65, 0x18, 0x04, 0x20, 0x01, 0x28, 0x0c,
	0x52, 0x0a, 0x69, 0x6d, 0x61, 0x67, 0x65, 0x50, 0x69, 0x65, 0x63, 0x65, 0x12, 0x2a, 0x0a, 0x10,
	0x73, 0x74, 0x72, 0x65, 0x61, 0x6d, 0x49, 0x73, 0x46, 0x69, 0x6e, 0x69, 0x73, 0x68, 0x65, 0x64,
	0x18, 0x05, 0x20, 0x01, 0x28, 0x08, 0x52, 0x10, 0x73, 0x74, 0x72, 0x65, 0x61, 0x6d, 0x49, 0x73,
	0x46, 0x69, 0x6e, 0x69, 0x73, 0x68, 0x65, 0x64, 0x42, 0x0c, 0x0a, 0x0a, 0x5f, 0x61, 0x75, 0x74,
	0x68, 0x54, 0x6f, 0x6b, 0x65, 0x6e, 0x42, 0x09, 0x0a, 0x07, 0x5f, 0x75, 0x73, 0x65, 0x72, 0x49,
	0x44, 0x42, 0x0c, 0x0a, 0x0a, 0x5f, 0x69, 0x6d, 0x61, 0x67, 0x65, 0x4e, 0x61, 0x6d, 0x65, 0x32,
	0x94, 0x01, 0x0a, 0x05, 0x49, 0x6d, 0x61, 0x67, 0x65, 0x12, 0x48, 0x0a, 0x0d, 0x44, 0x6f, 0x77,
	0x6e, 0x6c, 0x6f, 0x61, 0x64, 0x49, 0x6d, 0x61, 0x67, 0x65, 0x12, 0x18, 0x2e, 0x70, 0x62, 0x2e,
	0x72, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x44, 0x6f, 0x77, 0x6e, 0x6c, 0x6f, 0x61, 0x64, 0x49,
	0x6d, 0x61, 0x67, 0x65, 0x1a, 0x19, 0x2e, 0x70, 0x62, 0x2e, 0x72, 0x65, 0x73, 0x70, 0x6f, 0x6e,
	0x73, 0x65, 0x44, 0x6f, 0x77, 0x6e, 0x6c, 0x6f, 0x61, 0x64, 0x49, 0x6d, 0x61, 0x67, 0x65, 0x22,
	0x00, 0x30, 0x01, 0x12, 0x41, 0x0a, 0x0b, 0x55, 0x70, 0x6c, 0x6f, 0x61, 0x64, 0x49, 0x6d, 0x61,
	0x67, 0x65, 0x12, 0x16, 0x2e, 0x70, 0x62, 0x2e, 0x72, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x55,
	0x70, 0x6c, 0x6f, 0x61, 0x64, 0x49, 0x6d, 0x61, 0x67, 0x65, 0x1a, 0x16, 0x2e, 0x67, 0x6f, 0x6f,
	0x67, 0x6c, 0x65, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x62, 0x75, 0x66, 0x2e, 0x45, 0x6d, 0x70,
	0x74, 0x79, 0x22, 0x00, 0x28, 0x01, 0x42, 0x05, 0x5a, 0x03, 0x2f, 0x70, 0x62, 0x62, 0x06, 0x70,
	0x72, 0x6f, 0x74, 0x6f, 0x33,
}

var (
	file_imageRPC_proto_rawDescOnce sync.Once
	file_imageRPC_proto_rawDescData = file_imageRPC_proto_rawDesc
)

func file_imageRPC_proto_rawDescGZIP() []byte {
	file_imageRPC_proto_rawDescOnce.Do(func() {
		file_imageRPC_proto_rawDescData = protoimpl.X.CompressGZIP(file_imageRPC_proto_rawDescData)
	})
	return file_imageRPC_proto_rawDescData
}

var file_imageRPC_proto_msgTypes = make([]protoimpl.MessageInfo, 3)
var file_imageRPC_proto_goTypes = []interface{}{
	(*RequestDownloadImage)(nil),  // 0: pb.requestDownloadImage
	(*ResponseDownloadImage)(nil), // 1: pb.responseDownloadImage
	(*RequestUploadImage)(nil),    // 2: pb.requestUploadImage
	(*emptypb.Empty)(nil),         // 3: google.protobuf.Empty
}
var file_imageRPC_proto_depIdxs = []int32{
	0, // 0: pb.Image.DownloadImage:input_type -> pb.requestDownloadImage
	2, // 1: pb.Image.UploadImage:input_type -> pb.requestUploadImage
	1, // 2: pb.Image.DownloadImage:output_type -> pb.responseDownloadImage
	3, // 3: pb.Image.UploadImage:output_type -> google.protobuf.Empty
	2, // [2:4] is the sub-list for method output_type
	0, // [0:2] is the sub-list for method input_type
	0, // [0:0] is the sub-list for extension type_name
	0, // [0:0] is the sub-list for extension extendee
	0, // [0:0] is the sub-list for field type_name
}

func init() { file_imageRPC_proto_init() }
func file_imageRPC_proto_init() {
	if File_imageRPC_proto != nil {
		return
	}
	if !protoimpl.UnsafeEnabled {
		file_imageRPC_proto_msgTypes[0].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*RequestDownloadImage); i {
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
		file_imageRPC_proto_msgTypes[1].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*ResponseDownloadImage); i {
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
		file_imageRPC_proto_msgTypes[2].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*RequestUploadImage); i {
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
	file_imageRPC_proto_msgTypes[2].OneofWrappers = []interface{}{}
	type x struct{}
	out := protoimpl.TypeBuilder{
		File: protoimpl.DescBuilder{
			GoPackagePath: reflect.TypeOf(x{}).PkgPath(),
			RawDescriptor: file_imageRPC_proto_rawDesc,
			NumEnums:      0,
			NumMessages:   3,
			NumExtensions: 0,
			NumServices:   1,
		},
		GoTypes:           file_imageRPC_proto_goTypes,
		DependencyIndexes: file_imageRPC_proto_depIdxs,
		MessageInfos:      file_imageRPC_proto_msgTypes,
	}.Build()
	File_imageRPC_proto = out.File
	file_imageRPC_proto_rawDesc = nil
	file_imageRPC_proto_goTypes = nil
	file_imageRPC_proto_depIdxs = nil
}
