// Code generated by protoc-gen-go. DO NOT EDIT.
// source: 0.proto

package pb

import (
	fmt "fmt"
	proto "github.com/golang/protobuf/proto"
	_ "github.com/grpc-ecosystem/grpc-gateway/protoc-gen-swagger/options"
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
const _ = proto.ProtoPackageIsVersion3 // please upgrade the proto package

func init() { proto.RegisterFile("0.proto", fileDescriptor_b5d39afb3b422e60) }

var fileDescriptor_b5d39afb3b422e60 = []byte{
	// 341 bytes of a gzipped FileDescriptorProto
	0x1f, 0x8b, 0x08, 0x00, 0x00, 0x00, 0x00, 0x00, 0x02, 0xff, 0x64, 0x90, 0xcd, 0x4a, 0x23, 0x41,
	0x10, 0xc7, 0x77, 0x66, 0xd9, 0x0f, 0x06, 0x76, 0x09, 0xc3, 0x5e, 0x36, 0xa7, 0xda, 0xbd, 0x04,
	0x24, 0xe9, 0x49, 0x22, 0x88, 0xe4, 0x20, 0x24, 0x90, 0x44, 0x51, 0x51, 0xfc, 0x84, 0x80, 0x87,
	0x9e, 0x99, 0xca, 0x7c, 0x75, 0xba, 0x9a, 0xee, 0x9e, 0xc4, 0xe4, 0x05, 0xbc, 0x7b, 0xf7, 0x0d,
	0x04, 0x1f, 0xc4, 0xa3, 0x2f, 0x24, 0x4e, 0xf4, 0x10, 0xbc, 0xfd, 0xab, 0xf8, 0xff, 0xa8, 0x1f,
	0xe5, 0xfd, 0x68, 0x33, 0xa5, 0xc9, 0x92, 0xff, 0x55, 0xd8, 0xa2, 0xde, 0xac, 0x72, 0xd4, 0x4a,
	0x50, 0xb6, 0xcc, 0x82, 0x27, 0x09, 0xea, 0x80, 0x94, 0xcd, 0x48, 0x9a, 0x80, 0x4b, 0x49, 0x96,
	0x57, 0x79, 0x8d, 0x0c, 0x5e, 0xdc, 0xfb, 0xfe, 0xb3, 0xeb, 0xdf, 0x78, 0x7f, 0xae, 0xba, 0x6c,
	0x28, 0x73, 0x5a, 0x4e, 0x89, 0x62, 0x38, 0xd5, 0x94, 0x63, 0x64, 0xff, 0xef, 0x79, 0xff, 0x8e,
	0x30, 0x41, 0x19, 0x1b, 0xa0, 0x29, 0xd8, 0x54, 0x23, 0x42, 0x91, 0xc9, 0x24, 0xa6, 0x99, 0xf9,
	0x28, 0xf9, 0x7f, 0x53, 0x6b, 0x95, 0xe9, 0x05, 0x81, 0xb0, 0x05, 0xd3, 0x5c, 0x16, 0x22, 0x43,
	0xcd, 0x2c, 0x46, 0x69, 0xf7, 0x5b, 0x9b, 0xb5, 0x59, 0x67, 0xcb, 0x75, 0xdc, 0x6e, 0x8d, 0x2b,
	0x25, 0xb2, 0xa8, 0xba, 0x1e, 0xe4, 0x86, 0x64, 0xef, 0xd3, 0x66, 0xf2, 0xe4, 0x78, 0x8f, 0x8e,
	0xe7, 0x0d, 0x90, 0x6b, 0xd4, 0xfd, 0xd2, 0xa6, 0xfe, 0x83, 0xf3, 0xd3, 0xf5, 0xef, 0x9c, 0x8b,
	0x14, 0xe1, 0x6d, 0x26, 0x9d, 0xad, 0xaa, 0x3e, 0xa4, 0xc8, 0x63, 0xd4, 0x30, 0x2b, 0x8d, 0x85,
	0x10, 0xc1, 0xa0, 0x05, 0x4b, 0xb0, 0x26, 0x61, 0x4a, 0x42, 0xd0, 0x02, 0x63, 0x08, 0x97, 0xc0,
	0xc1, 0x28, 0x1e, 0x21, 0x70, 0x19, 0x03, 0x07, 0x4b, 0x05, 0x4a, 0x06, 0x23, 0xd2, 0x80, 0xb7,
	0x7c, 0xa6, 0x04, 0x36, 0xa1, 0xf1, 0x4e, 0xcd, 0xf7, 0x2f, 0x79, 0x98, 0x0d, 0x86, 0x07, 0x87,
	0xd9, 0xae, 0xec, 0x9c, 0xc5, 0xf3, 0xeb, 0x93, 0x7c, 0x3c, 0x2a, 0xc5, 0xf8, 0xfc, 0x78, 0x67,
	0x55, 0xca, 0xb0, 0xc1, 0xea, 0xbf, 0x36, 0x25, 0xdc, 0xb0, 0xe6, 0xfd, 0xde, 0xf0, 0xfd, 0x32,
	0x71, 0x55, 0x18, 0x7e, 0xaf, 0x9e, 0xbb, 0xfd, 0x1a, 0x00, 0x00, 0xff, 0xff, 0xc0, 0xe2, 0x98,
	0x79, 0x9a, 0x01, 0x00, 0x00,
}
