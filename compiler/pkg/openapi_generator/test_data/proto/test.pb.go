// Code generated by protoc-gen-gogo. DO NOT EDIT.
// source: test_data/proto/test.proto

package proto

import (
	fmt "fmt"
	math "math"

	proto "github.com/gogo/protobuf/proto"
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

// +k8s:openapi-gen=true
type EnumValue int32

const (
	EnumValue_FIZZ EnumValue = 0
	EnumValue_BUZZ EnumValue = 1
)

var EnumValue_name = map[int32]string{
	0: "FIZZ",
	1: "BUZZ",
}

var EnumValue_value = map[string]int32{
	"FIZZ": 0,
	"BUZZ": 1,
}

func (x EnumValue) String() string {
	return proto.EnumName(EnumValue_name, int32(x))
}

func (EnumValue) EnumDescriptor() ([]byte, []int) {
	return fileDescriptor_40509a13d82aed2b, []int{0}
}

type Foo_NestedEnum int32

const (
	Foo_JEDEN Foo_NestedEnum = 0
	Foo_DWA   Foo_NestedEnum = 1
)

var Foo_NestedEnum_name = map[int32]string{
	0: "JEDEN",
	1: "DWA",
}

var Foo_NestedEnum_value = map[string]int32{
	"JEDEN": 0,
	"DWA":   1,
}

func (x Foo_NestedEnum) String() string {
	return proto.EnumName(Foo_NestedEnum_name, int32(x))
}

func (Foo_NestedEnum) EnumDescriptor() ([]byte, []int) {
	return fileDescriptor_40509a13d82aed2b, []int{0, 0}
}

// +k8s:openapi-gen=true
type Foo struct {
	DoubleValue   float64 `protobuf:"fixed64,1,opt,name=double_value,json=doubleValue,proto3" json:"double_value,omitempty"`
	FloatValue    float32 `protobuf:"fixed32,2,opt,name=float_value,json=floatValue,proto3" json:"float_value,omitempty"`
	Int32Value    int32   `protobuf:"varint,3,opt,name=int32_value,json=int32Value,proto3" json:"int32_value,omitempty"`
	Int64Value    int64   `protobuf:"varint,4,opt,name=int64_value,json=int64Value,proto3" json:"int64_value,omitempty"`
	Uint32Value   uint32  `protobuf:"varint,5,opt,name=uint32_value,json=uint32Value,proto3" json:"uint32_value,omitempty"`
	Uint64Value   uint64  `protobuf:"varint,6,opt,name=uint64_value,json=uint64Value,proto3" json:"uint64_value,omitempty"`
	Sint32Value   int32   `protobuf:"zigzag32,7,opt,name=sint32_value,json=sint32Value,proto3" json:"sint32_value,omitempty"`
	Sint64Value   int64   `protobuf:"zigzag64,8,opt,name=sint64_value,json=sint64Value,proto3" json:"sint64_value,omitempty"`
	Fixed32Value  uint32  `protobuf:"fixed32,9,opt,name=fixed32_value,json=fixed32Value,proto3" json:"fixed32_value,omitempty"`
	Fixed64Value  uint64  `protobuf:"fixed64,10,opt,name=fixed64_value,json=fixed64Value,proto3" json:"fixed64_value,omitempty"`
	Sfixed32Value int32   `protobuf:"fixed32,11,opt,name=sfixed32_value,json=sfixed32Value,proto3" json:"sfixed32_value,omitempty"`
	Sfixed64Value int64   `protobuf:"fixed64,12,opt,name=sfixed64_value,json=sfixed64Value,proto3" json:"sfixed64_value,omitempty"`
	BoolValue     bool    `protobuf:"varint,13,opt,name=bool_value,json=boolValue,proto3" json:"bool_value,omitempty"`
	StringValue   string  `protobuf:"bytes,14,opt,name=string_value,json=stringValue,proto3" json:"string_value,omitempty"`
	BytesValue    []byte  `protobuf:"bytes,15,opt,name=bytes_value,json=bytesValue,proto3" json:"bytes_value,omitempty"`
	// Mesh7CodeGenOpenAPIEnum
	EnumValue EnumValue `protobuf:"varint,17,opt,name=enum_value,json=enumValue,proto3,enum=proto.EnumValue" json:"enum_value,omitempty"`
	// Mesh7CodeGenOpenAPIEnum
	NestedEnumValue Foo_NestedEnum     `protobuf:"varint,18,opt,name=nested_enum_value,json=nestedEnumValue,proto3,enum=proto.Foo_NestedEnum" json:"nested_enum_value,omitempty"`
	NestedMessage   *Foo_NestedMessage `protobuf:"bytes,19,opt,name=nested_message,json=nestedMessage,proto3" json:"nested_message,omitempty"`
	// Types that are valid to be assigned to OneofValue:
	//	*Foo_OneofValueString
	//	*Foo_OneofValueInt
	OneofValue           isFoo_OneofValue  `protobuf_oneof:"oneof_value"`
	MapValue             map[string]string `protobuf:"bytes,22,rep,name=map_value,json=mapValue,proto3" json:"map_value,omitempty" protobuf_key:"bytes,1,opt,name=key,proto3" protobuf_val:"bytes,2,opt,name=value,proto3"`
	RepeatedValue        []int32           `protobuf:"varint,23,rep,packed,name=repeated_value,json=repeatedValue,proto3" json:"repeated_value,omitempty"`
	BarValue             *Bar              `protobuf:"bytes,24,opt,name=bar_value,json=barValue,proto3" json:"bar_value,omitempty"`
	XXX_NoUnkeyedLiteral struct{}          `json:"-"`
	XXX_unrecognized     []byte            `json:"-"`
	XXX_sizecache        int32             `json:"-"`
}

func (m *Foo) Reset()         { *m = Foo{} }
func (m *Foo) String() string { return proto.CompactTextString(m) }
func (*Foo) ProtoMessage()    {}
func (*Foo) Descriptor() ([]byte, []int) {
	return fileDescriptor_40509a13d82aed2b, []int{0}
}
func (m *Foo) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_Foo.Unmarshal(m, b)
}
func (m *Foo) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_Foo.Marshal(b, m, deterministic)
}
func (m *Foo) XXX_Merge(src proto.Message) {
	xxx_messageInfo_Foo.Merge(m, src)
}
func (m *Foo) XXX_Size() int {
	return xxx_messageInfo_Foo.Size(m)
}
func (m *Foo) XXX_DiscardUnknown() {
	xxx_messageInfo_Foo.DiscardUnknown(m)
}

var xxx_messageInfo_Foo proto.InternalMessageInfo

type isFoo_OneofValue interface {
	isFoo_OneofValue()
}

type Foo_OneofValueString struct {
	OneofValueString string `protobuf:"bytes,20,opt,name=oneof_value_string,json=oneofValueString,proto3,oneof" json:"oneof_value_string,omitempty"`
}
type Foo_OneofValueInt struct {
	OneofValueInt int32 `protobuf:"varint,21,opt,name=oneof_value_int,json=oneofValueInt,proto3,oneof" json:"oneof_value_int,omitempty"`
}

func (*Foo_OneofValueString) isFoo_OneofValue() {}
func (*Foo_OneofValueInt) isFoo_OneofValue()    {}

func (m *Foo) GetOneofValue() isFoo_OneofValue {
	if m != nil {
		return m.OneofValue
	}
	return nil
}

func (m *Foo) GetDoubleValue() float64 {
	if m != nil {
		return m.DoubleValue
	}
	return 0
}

func (m *Foo) GetFloatValue() float32 {
	if m != nil {
		return m.FloatValue
	}
	return 0
}

func (m *Foo) GetInt32Value() int32 {
	if m != nil {
		return m.Int32Value
	}
	return 0
}

func (m *Foo) GetInt64Value() int64 {
	if m != nil {
		return m.Int64Value
	}
	return 0
}

func (m *Foo) GetUint32Value() uint32 {
	if m != nil {
		return m.Uint32Value
	}
	return 0
}

func (m *Foo) GetUint64Value() uint64 {
	if m != nil {
		return m.Uint64Value
	}
	return 0
}

func (m *Foo) GetSint32Value() int32 {
	if m != nil {
		return m.Sint32Value
	}
	return 0
}

func (m *Foo) GetSint64Value() int64 {
	if m != nil {
		return m.Sint64Value
	}
	return 0
}

func (m *Foo) GetFixed32Value() uint32 {
	if m != nil {
		return m.Fixed32Value
	}
	return 0
}

func (m *Foo) GetFixed64Value() uint64 {
	if m != nil {
		return m.Fixed64Value
	}
	return 0
}

func (m *Foo) GetSfixed32Value() int32 {
	if m != nil {
		return m.Sfixed32Value
	}
	return 0
}

func (m *Foo) GetSfixed64Value() int64 {
	if m != nil {
		return m.Sfixed64Value
	}
	return 0
}

func (m *Foo) GetBoolValue() bool {
	if m != nil {
		return m.BoolValue
	}
	return false
}

func (m *Foo) GetStringValue() string {
	if m != nil {
		return m.StringValue
	}
	return ""
}

func (m *Foo) GetBytesValue() []byte {
	if m != nil {
		return m.BytesValue
	}
	return nil
}

func (m *Foo) GetEnumValue() EnumValue {
	if m != nil {
		return m.EnumValue
	}
	return EnumValue_FIZZ
}

func (m *Foo) GetNestedEnumValue() Foo_NestedEnum {
	if m != nil {
		return m.NestedEnumValue
	}
	return Foo_JEDEN
}

func (m *Foo) GetNestedMessage() *Foo_NestedMessage {
	if m != nil {
		return m.NestedMessage
	}
	return nil
}

func (m *Foo) GetOneofValueString() string {
	if x, ok := m.GetOneofValue().(*Foo_OneofValueString); ok {
		return x.OneofValueString
	}
	return ""
}

func (m *Foo) GetOneofValueInt() int32 {
	if x, ok := m.GetOneofValue().(*Foo_OneofValueInt); ok {
		return x.OneofValueInt
	}
	return 0
}

func (m *Foo) GetMapValue() map[string]string {
	if m != nil {
		return m.MapValue
	}
	return nil
}

func (m *Foo) GetRepeatedValue() []int32 {
	if m != nil {
		return m.RepeatedValue
	}
	return nil
}

func (m *Foo) GetBarValue() *Bar {
	if m != nil {
		return m.BarValue
	}
	return nil
}

// XXX_OneofWrappers is for the internal use of the proto package.
func (*Foo) XXX_OneofWrappers() []interface{} {
	return []interface{}{
		(*Foo_OneofValueString)(nil),
		(*Foo_OneofValueInt)(nil),
	}
}

// +k8s:openapi-gen=true
type Foo_NestedMessage struct {
	NestedMessageValue   string   `protobuf:"bytes,16,opt,name=nested_message_value,json=nestedMessageValue,proto3" json:"nested_message_value,omitempty"`
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_unrecognized     []byte   `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *Foo_NestedMessage) Reset()         { *m = Foo_NestedMessage{} }
func (m *Foo_NestedMessage) String() string { return proto.CompactTextString(m) }
func (*Foo_NestedMessage) ProtoMessage()    {}
func (*Foo_NestedMessage) Descriptor() ([]byte, []int) {
	return fileDescriptor_40509a13d82aed2b, []int{0, 0}
}
func (m *Foo_NestedMessage) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_Foo_NestedMessage.Unmarshal(m, b)
}
func (m *Foo_NestedMessage) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_Foo_NestedMessage.Marshal(b, m, deterministic)
}
func (m *Foo_NestedMessage) XXX_Merge(src proto.Message) {
	xxx_messageInfo_Foo_NestedMessage.Merge(m, src)
}
func (m *Foo_NestedMessage) XXX_Size() int {
	return xxx_messageInfo_Foo_NestedMessage.Size(m)
}
func (m *Foo_NestedMessage) XXX_DiscardUnknown() {
	xxx_messageInfo_Foo_NestedMessage.DiscardUnknown(m)
}

var xxx_messageInfo_Foo_NestedMessage proto.InternalMessageInfo

func (m *Foo_NestedMessage) GetNestedMessageValue() string {
	if m != nil {
		return m.NestedMessageValue
	}
	return ""
}

// +k8s:openapi-gen=true
type Bar struct {
	// Mesh7CodeGenOpenAPIEnum
	EnumValue EnumValue `protobuf:"varint,1,opt,name=enum_value,json=enumValue,proto3,enum=proto.EnumValue" json:"enum_value,omitempty"`
	// Types that are valid to be assigned to OneofValue:
	//	*Bar_OneofValueString
	//	*Bar_OneofValueInt
	OneofValue           isBar_OneofValue `protobuf_oneof:"oneof_value"`
	XXX_NoUnkeyedLiteral struct{}         `json:"-"`
	XXX_unrecognized     []byte           `json:"-"`
	XXX_sizecache        int32            `json:"-"`
}

func (m *Bar) Reset()         { *m = Bar{} }
func (m *Bar) String() string { return proto.CompactTextString(m) }
func (*Bar) ProtoMessage()    {}
func (*Bar) Descriptor() ([]byte, []int) {
	return fileDescriptor_40509a13d82aed2b, []int{1}
}
func (m *Bar) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_Bar.Unmarshal(m, b)
}
func (m *Bar) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_Bar.Marshal(b, m, deterministic)
}
func (m *Bar) XXX_Merge(src proto.Message) {
	xxx_messageInfo_Bar.Merge(m, src)
}
func (m *Bar) XXX_Size() int {
	return xxx_messageInfo_Bar.Size(m)
}
func (m *Bar) XXX_DiscardUnknown() {
	xxx_messageInfo_Bar.DiscardUnknown(m)
}

var xxx_messageInfo_Bar proto.InternalMessageInfo

type isBar_OneofValue interface {
	isBar_OneofValue()
}

type Bar_OneofValueString struct {
	OneofValueString string `protobuf:"bytes,2,opt,name=oneof_value_string,json=oneofValueString,proto3,oneof" json:"oneof_value_string,omitempty"`
}
type Bar_OneofValueInt struct {
	OneofValueInt int32 `protobuf:"varint,3,opt,name=oneof_value_int,json=oneofValueInt,proto3,oneof" json:"oneof_value_int,omitempty"`
}

func (*Bar_OneofValueString) isBar_OneofValue() {}
func (*Bar_OneofValueInt) isBar_OneofValue()    {}

func (m *Bar) GetOneofValue() isBar_OneofValue {
	if m != nil {
		return m.OneofValue
	}
	return nil
}

func (m *Bar) GetEnumValue() EnumValue {
	if m != nil {
		return m.EnumValue
	}
	return EnumValue_FIZZ
}

func (m *Bar) GetOneofValueString() string {
	if x, ok := m.GetOneofValue().(*Bar_OneofValueString); ok {
		return x.OneofValueString
	}
	return ""
}

func (m *Bar) GetOneofValueInt() int32 {
	if x, ok := m.GetOneofValue().(*Bar_OneofValueInt); ok {
		return x.OneofValueInt
	}
	return 0
}

// XXX_OneofWrappers is for the internal use of the proto package.
func (*Bar) XXX_OneofWrappers() []interface{} {
	return []interface{}{
		(*Bar_OneofValueString)(nil),
		(*Bar_OneofValueInt)(nil),
	}
}

func init() {
	proto.RegisterEnum("proto.EnumValue", EnumValue_name, EnumValue_value)
	proto.RegisterEnum("proto.Foo_NestedEnum", Foo_NestedEnum_name, Foo_NestedEnum_value)
	proto.RegisterType((*Foo)(nil), "proto.Foo")
	proto.RegisterMapType((map[string]string)(nil), "proto.Foo.MapValueEntry")
	proto.RegisterType((*Foo_NestedMessage)(nil), "proto.Foo.NestedMessage")
	proto.RegisterType((*Bar)(nil), "proto.Bar")
}

func init() { proto.RegisterFile("test_data/proto/test.proto", fileDescriptor_40509a13d82aed2b) }

var fileDescriptor_40509a13d82aed2b = []byte{
	// 645 bytes of a gzipped FileDescriptorProto
	0x1f, 0x8b, 0x08, 0x00, 0x00, 0x00, 0x00, 0x00, 0x02, 0xff, 0x94, 0x93, 0x5f, 0x6b, 0xdb, 0x30,
	0x10, 0xc0, 0xa3, 0xb8, 0x69, 0xec, 0x73, 0x9c, 0xb8, 0x5a, 0xbb, 0x85, 0xc2, 0xa8, 0xda, 0x31,
	0x26, 0xc6, 0x48, 0x47, 0xda, 0x95, 0xb1, 0x3d, 0x8c, 0x86, 0xa6, 0xb4, 0x83, 0xf6, 0x41, 0x65,
	0x1b, 0xe4, 0x25, 0x28, 0x8b, 0x52, 0xc2, 0x12, 0x2b, 0xd8, 0x72, 0x59, 0x3f, 0xcd, 0x3e, 0xdc,
	0xbe, 0xc8, 0xb0, 0xa4, 0xf8, 0x4f, 0x57, 0x18, 0x7b, 0xca, 0xe9, 0xf4, 0xd3, 0x2f, 0x3e, 0xdd,
	0x09, 0x76, 0x95, 0x48, 0xd4, 0x78, 0xca, 0x15, 0x3f, 0x5c, 0xc5, 0x52, 0xc9, 0xc3, 0x6c, 0xdd,
	0xd3, 0x21, 0x6e, 0xe8, 0x9f, 0x83, 0xdf, 0x2e, 0x38, 0xe7, 0x52, 0xe2, 0x7d, 0x68, 0x4d, 0x65,
	0x3a, 0x59, 0x88, 0xf1, 0x1d, 0x5f, 0xa4, 0xa2, 0x8b, 0x08, 0xa2, 0x88, 0xf9, 0x26, 0xf7, 0x35,
	0x4b, 0xe1, 0x3d, 0xf0, 0x67, 0x0b, 0xc9, 0x95, 0x25, 0xea, 0x04, 0xd1, 0x3a, 0x03, 0x9d, 0xca,
	0x81, 0x79, 0xa4, 0x8e, 0xfa, 0x16, 0x70, 0x08, 0xa2, 0x0d, 0x06, 0x3a, 0x55, 0x06, 0x4e, 0x8e,
	0x2d, 0xb0, 0x41, 0x10, 0x75, 0x34, 0x70, 0x72, 0x6c, 0x80, 0x7d, 0x68, 0xa5, 0x65, 0x45, 0x83,
	0x20, 0x1a, 0x30, 0x3f, 0x2d, 0x39, 0x2c, 0x92, 0x4b, 0x36, 0x09, 0xa2, 0x1b, 0x06, 0x29, 0x59,
	0x92, 0xb2, 0xa5, 0x49, 0x10, 0xdd, 0x62, 0x7e, 0x52, 0xb5, 0x24, 0x65, 0x8b, 0x4b, 0x10, 0xc5,
	0x06, 0x59, 0x5b, 0x5e, 0x40, 0x30, 0x9b, 0xff, 0x14, 0xd3, 0x5c, 0xe3, 0x11, 0x44, 0x9b, 0xac,
	0x65, 0x93, 0x55, 0x28, 0x17, 0x01, 0x41, 0x74, 0xd3, 0x42, 0x6b, 0xd3, 0x4b, 0x68, 0x27, 0x55,
	0x95, 0x4f, 0x10, 0xed, 0xb0, 0x20, 0xa9, 0xb8, 0x72, 0x2c, 0x97, 0xb5, 0x08, 0xa2, 0xe1, 0x1a,
	0x5b, 0xdb, 0x9e, 0x03, 0x4c, 0xa4, 0x5c, 0x58, 0x24, 0x20, 0x88, 0xba, 0xcc, 0xcb, 0x32, 0x45,
	0x65, 0x2a, 0x9e, 0x47, 0xb7, 0x16, 0x68, 0x13, 0x44, 0x3d, 0xe6, 0x9b, 0x5c, 0xde, 0x86, 0xc9,
	0xbd, 0x12, 0x89, 0x25, 0x3a, 0x04, 0xd1, 0x16, 0x03, 0x9d, 0x32, 0xc0, 0x21, 0x80, 0x88, 0xd2,
	0xa5, 0xdd, 0xdf, 0x22, 0x88, 0xb6, 0xfb, 0xa1, 0x99, 0x9b, 0xde, 0x30, 0x4a, 0x97, 0x9a, 0x62,
	0x9e, 0x58, 0x87, 0xf8, 0x14, 0xb6, 0x22, 0x91, 0x28, 0x31, 0x1d, 0x97, 0xce, 0x61, 0x7d, 0x6e,
	0xc7, 0x9e, 0x3b, 0x97, 0xb2, 0x77, 0xad, 0x99, 0xcc, 0xc0, 0x3a, 0x51, 0x1e, 0x1b, 0xc5, 0x27,
	0x68, 0x5b, 0xc5, 0x52, 0x24, 0x09, 0xbf, 0x15, 0xdd, 0x27, 0x04, 0x51, 0xbf, 0xdf, 0xfd, 0xeb,
	0xfc, 0x95, 0xd9, 0x67, 0x41, 0x54, 0x5e, 0xe2, 0x1e, 0x60, 0x19, 0x09, 0x39, 0x33, 0xff, 0x3e,
	0x36, 0x05, 0x77, 0xb7, 0xb3, 0xf2, 0x2f, 0x6a, 0x2c, 0xd4, 0x7b, 0xfa, 0xaf, 0x6e, 0xf4, 0x0e,
	0xa6, 0xd0, 0x29, 0xf3, 0xf3, 0x48, 0x75, 0x77, 0xb2, 0x89, 0xbd, 0xa8, 0xb1, 0xa0, 0x80, 0x2f,
	0x23, 0x85, 0xdf, 0x81, 0xb7, 0xe4, 0x2b, 0x5b, 0xd5, 0x53, 0xe2, 0x3c, 0xf8, 0xaa, 0x2b, 0xbe,
	0xd2, 0xe8, 0x30, 0x52, 0xf1, 0x3d, 0x73, 0x97, 0x76, 0x99, 0xf5, 0x33, 0x16, 0x2b, 0xc1, 0xb3,
	0x9a, 0xcc, 0xd9, 0x67, 0xc4, 0xa1, 0x0d, 0x16, 0xac, 0xb3, 0x06, 0x7b, 0x05, 0xde, 0x84, 0xc7,
	0x96, 0xe8, 0xea, 0x9a, 0xc1, 0xda, 0x07, 0x3c, 0x66, 0xee, 0x84, 0xc7, 0x1a, 0xdc, 0x3d, 0x85,
	0xa0, 0x72, 0x01, 0xf8, 0x2d, 0x6c, 0x57, 0xaf, 0xcc, 0x4a, 0x42, 0xdd, 0x72, 0x5c, 0xb9, 0x1e,
	0xa3, 0xf8, 0x08, 0x41, 0xe5, 0x6b, 0x71, 0x08, 0xce, 0x0f, 0x71, 0xaf, 0x5f, 0xbb, 0xc7, 0xb2,
	0x10, 0x6f, 0x43, 0xa3, 0x78, 0xdf, 0x1e, 0x33, 0x8b, 0x0f, 0xf5, 0xf7, 0xe8, 0x80, 0x00, 0x14,
	0x0d, 0xc4, 0x1e, 0x34, 0x3e, 0x0f, 0xcf, 0x86, 0xd7, 0x61, 0x0d, 0x37, 0xc1, 0x39, 0xfb, 0x76,
	0x1a, 0xa2, 0x41, 0x00, 0x7e, 0xe9, 0x4a, 0x99, 0x1b, 0x8b, 0x44, 0xc4, 0x77, 0x62, 0x7a, 0xf0,
	0x0b, 0x81, 0x33, 0xe0, 0xf1, 0x83, 0xc1, 0x42, 0xff, 0x1e, 0xac, 0xc7, 0x9b, 0x5a, 0xff, 0x9f,
	0xa6, 0x3a, 0x8f, 0x36, 0xf5, 0xc1, 0xb7, 0xbe, 0xde, 0x03, 0xaf, 0x98, 0x45, 0x17, 0x36, 0xce,
	0x2f, 0x47, 0xa3, 0xb0, 0x96, 0x45, 0x83, 0x2f, 0xa3, 0x51, 0x88, 0xfa, 0x6f, 0xa0, 0x79, 0x23,
	0xe2, 0xbb, 0xf9, 0xf7, 0xec, 0x89, 0x35, 0xcf, 0xe4, 0x8d, 0x4a, 0x67, 0x33, 0x0c, 0xc5, 0x1c,
	0xec, 0x96, 0xe2, 0xc9, 0xa6, 0x0e, 0x8f, 0xfe, 0x04, 0x00, 0x00, 0xff, 0xff, 0x30, 0x9d, 0x30,
	0xb6, 0x82, 0x05, 0x00, 0x00,
}