// Code generated by protoc-gen-go. DO NOT EDIT.
// source: stats.proto

package main

import (
	fmt "fmt"
	proto "github.com/golang/protobuf/proto"
	_ "github.com/golang/protobuf/ptypes/empty"
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

type ResponseMessage struct {
	Res                  []byte   `protobuf:"bytes,1,opt,name=res,proto3" json:"res,omitempty"`
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_unrecognized     []byte   `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *ResponseMessage) Reset()         { *m = ResponseMessage{} }
func (m *ResponseMessage) String() string { return proto.CompactTextString(m) }
func (*ResponseMessage) ProtoMessage()    {}
func (*ResponseMessage) Descriptor() ([]byte, []int) {
	return fileDescriptor_b4756a0aec8b9d44, []int{0}
}

func (m *ResponseMessage) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_ResponseMessage.Unmarshal(m, b)
}
func (m *ResponseMessage) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_ResponseMessage.Marshal(b, m, deterministic)
}
func (m *ResponseMessage) XXX_Merge(src proto.Message) {
	xxx_messageInfo_ResponseMessage.Merge(m, src)
}
func (m *ResponseMessage) XXX_Size() int {
	return xxx_messageInfo_ResponseMessage.Size(m)
}
func (m *ResponseMessage) XXX_DiscardUnknown() {
	xxx_messageInfo_ResponseMessage.DiscardUnknown(m)
}

var xxx_messageInfo_ResponseMessage proto.InternalMessageInfo

func (m *ResponseMessage) GetRes() []byte {
	if m != nil {
		return m.Res
	}
	return nil
}

type LogItemMessage struct {
	CreatedAt            int64    `protobuf:"varint,1,opt,name=CreatedAt,proto3" json:"CreatedAt,omitempty"`
	Name                 string   `protobuf:"bytes,2,opt,name=Name,proto3" json:"Name,omitempty"`
	Amount               float32  `protobuf:"fixed32,3,opt,name=Amount,proto3" json:"Amount,omitempty"`
	Category             string   `protobuf:"bytes,4,opt,name=Category,proto3" json:"Category,omitempty"`
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_unrecognized     []byte   `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *LogItemMessage) Reset()         { *m = LogItemMessage{} }
func (m *LogItemMessage) String() string { return proto.CompactTextString(m) }
func (*LogItemMessage) ProtoMessage()    {}
func (*LogItemMessage) Descriptor() ([]byte, []int) {
	return fileDescriptor_b4756a0aec8b9d44, []int{1}
}

func (m *LogItemMessage) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_LogItemMessage.Unmarshal(m, b)
}
func (m *LogItemMessage) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_LogItemMessage.Marshal(b, m, deterministic)
}
func (m *LogItemMessage) XXX_Merge(src proto.Message) {
	xxx_messageInfo_LogItemMessage.Merge(m, src)
}
func (m *LogItemMessage) XXX_Size() int {
	return xxx_messageInfo_LogItemMessage.Size(m)
}
func (m *LogItemMessage) XXX_DiscardUnknown() {
	xxx_messageInfo_LogItemMessage.DiscardUnknown(m)
}

var xxx_messageInfo_LogItemMessage proto.InternalMessageInfo

func (m *LogItemMessage) GetCreatedAt() int64 {
	if m != nil {
		return m.CreatedAt
	}
	return 0
}

func (m *LogItemMessage) GetName() string {
	if m != nil {
		return m.Name
	}
	return ""
}

func (m *LogItemMessage) GetAmount() float32 {
	if m != nil {
		return m.Amount
	}
	return 0
}

func (m *LogItemMessage) GetCategory() string {
	if m != nil {
		return m.Category
	}
	return ""
}

type LogItemQueryMessage struct {
	LogItems             []*LogItemMessage `protobuf:"bytes,1,rep,name=LogItems,proto3" json:"LogItems,omitempty"`
	XXX_NoUnkeyedLiteral struct{}          `json:"-"`
	XXX_unrecognized     []byte            `json:"-"`
	XXX_sizecache        int32             `json:"-"`
}

func (m *LogItemQueryMessage) Reset()         { *m = LogItemQueryMessage{} }
func (m *LogItemQueryMessage) String() string { return proto.CompactTextString(m) }
func (*LogItemQueryMessage) ProtoMessage()    {}
func (*LogItemQueryMessage) Descriptor() ([]byte, []int) {
	return fileDescriptor_b4756a0aec8b9d44, []int{2}
}

func (m *LogItemQueryMessage) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_LogItemQueryMessage.Unmarshal(m, b)
}
func (m *LogItemQueryMessage) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_LogItemQueryMessage.Marshal(b, m, deterministic)
}
func (m *LogItemQueryMessage) XXX_Merge(src proto.Message) {
	xxx_messageInfo_LogItemQueryMessage.Merge(m, src)
}
func (m *LogItemQueryMessage) XXX_Size() int {
	return xxx_messageInfo_LogItemQueryMessage.Size(m)
}
func (m *LogItemQueryMessage) XXX_DiscardUnknown() {
	xxx_messageInfo_LogItemQueryMessage.DiscardUnknown(m)
}

var xxx_messageInfo_LogItemQueryMessage proto.InternalMessageInfo

func (m *LogItemQueryMessage) GetLogItems() []*LogItemMessage {
	if m != nil {
		return m.LogItems
	}
	return nil
}

func init() {
	proto.RegisterType((*ResponseMessage)(nil), "main.ResponseMessage")
	proto.RegisterType((*LogItemMessage)(nil), "main.LogItemMessage")
	proto.RegisterType((*LogItemQueryMessage)(nil), "main.LogItemQueryMessage")
}

func init() { proto.RegisterFile("stats.proto", fileDescriptor_b4756a0aec8b9d44) }

var fileDescriptor_b4756a0aec8b9d44 = []byte{
	// 278 bytes of a gzipped FileDescriptorProto
	0x1f, 0x8b, 0x08, 0x00, 0x00, 0x00, 0x00, 0x00, 0x02, 0xff, 0x9c, 0x90, 0xdd, 0x4a, 0xc3, 0x40,
	0x10, 0x85, 0x4d, 0x53, 0x4b, 0x3b, 0x15, 0x95, 0xd5, 0x86, 0x10, 0xbd, 0x08, 0xf1, 0x26, 0x57,
	0x5b, 0xa9, 0x4f, 0x10, 0x8b, 0xd4, 0x82, 0x15, 0x5c, 0x7d, 0x81, 0x2d, 0x8e, 0xb1, 0x90, 0xcd,
	0x86, 0xec, 0x44, 0xc8, 0xbb, 0xf8, 0xb0, 0xb2, 0xf9, 0x69, 0xd1, 0x4b, 0xef, 0x66, 0xce, 0x7e,
	0x9c, 0x39, 0x7b, 0x60, 0x6a, 0x48, 0x92, 0xe1, 0x45, 0xa9, 0x49, 0xb3, 0xa1, 0x92, 0xbb, 0x3c,
	0xb8, 0x4a, 0xb5, 0x4e, 0x33, 0x9c, 0x37, 0xda, 0xb6, 0xfa, 0x98, 0xa3, 0x2a, 0xa8, 0x6e, 0x91,
	0xe8, 0x06, 0xce, 0x04, 0x9a, 0x42, 0xe7, 0x06, 0x37, 0x68, 0x8c, 0x4c, 0x91, 0x9d, 0x83, 0x5b,
	0xa2, 0xf1, 0x9d, 0xd0, 0x89, 0x4f, 0x84, 0x1d, 0xa3, 0x2f, 0x38, 0x7d, 0xd2, 0xe9, 0x9a, 0x50,
	0xf5, 0xcc, 0x35, 0x4c, 0x96, 0x25, 0x4a, 0xc2, 0xf7, 0x84, 0x1a, 0xd2, 0x15, 0x07, 0x81, 0x31,
	0x18, 0x3e, 0x4b, 0x85, 0xfe, 0x20, 0x74, 0xe2, 0x89, 0x68, 0x66, 0xe6, 0xc1, 0x28, 0x51, 0xba,
	0xca, 0xc9, 0x77, 0x43, 0x27, 0x1e, 0x88, 0x6e, 0x63, 0x01, 0x8c, 0x97, 0x92, 0x30, 0xd5, 0x65,
	0xed, 0x0f, 0x1b, 0x7e, 0xbf, 0x47, 0x2b, 0xb8, 0xe8, 0xee, 0xbe, 0x54, 0x58, 0xd6, 0xfd, 0xf1,
	0x5b, 0x18, 0x77, 0xb2, 0x4d, 0xe9, 0xc6, 0xd3, 0xc5, 0x25, 0xb7, 0x3f, 0xe5, 0xbf, 0x43, 0x8a,
	0x3d, 0xb5, 0xf8, 0x76, 0xe0, 0xf8, 0xd5, 0x16, 0xc3, 0x1e, 0x61, 0xb6, 0x42, 0x4a, 0xb2, 0xec,
	0x6d, 0xa7, 0xf0, 0xbe, 0xde, 0xe8, 0x9c, 0x3e, 0xed, 0x0b, 0xf3, 0x78, 0x5b, 0x13, 0xef, 0x6b,
	0xe2, 0x0f, 0xb6, 0xa6, 0x60, 0xd6, 0x5a, 0xff, 0x29, 0x29, 0x3a, 0x62, 0x6b, 0xf0, 0x0e, 0x4e,
	0x7d, 0xe4, 0x7f, 0x59, 0x6d, 0x47, 0x0d, 0x78, 0xf7, 0x13, 0x00, 0x00, 0xff, 0xff, 0xcf, 0xf2,
	0x18, 0xec, 0xbd, 0x01, 0x00, 0x00,
}