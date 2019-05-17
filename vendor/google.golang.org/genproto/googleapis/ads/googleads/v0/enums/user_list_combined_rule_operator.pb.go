// Code generated by protoc-gen-go. DO NOT EDIT.
// source: google/ads/googleads/v0/enums/user_list_combined_rule_operator.proto

package enums // import "google.golang.org/genproto/googleapis/ads/googleads/v0/enums"

import proto "github.com/golang/protobuf/proto"
import fmt "fmt"
import math "math"

// Reference imports to suppress errors if they are not otherwise used.
var _ = proto.Marshal
var _ = fmt.Errorf
var _ = math.Inf

// This is a compile-time assertion to ensure that this generated file
// is compatible with the proto package it is being compiled against.
// A compilation error at this line likely means your copy of the
// proto package needs to be updated.
const _ = proto.ProtoPackageIsVersion2 // please upgrade the proto package

// Enum describing possible user list combined rule operators.
type UserListCombinedRuleOperatorEnum_UserListCombinedRuleOperator int32

const (
	// Not specified.
	UserListCombinedRuleOperatorEnum_UNSPECIFIED UserListCombinedRuleOperatorEnum_UserListCombinedRuleOperator = 0
	// Used for return value only. Represents value unknown in this version.
	UserListCombinedRuleOperatorEnum_UNKNOWN UserListCombinedRuleOperatorEnum_UserListCombinedRuleOperator = 1
	// A AND B.
	UserListCombinedRuleOperatorEnum_AND UserListCombinedRuleOperatorEnum_UserListCombinedRuleOperator = 2
	// A AND NOT B.
	UserListCombinedRuleOperatorEnum_AND_NOT UserListCombinedRuleOperatorEnum_UserListCombinedRuleOperator = 3
)

var UserListCombinedRuleOperatorEnum_UserListCombinedRuleOperator_name = map[int32]string{
	0: "UNSPECIFIED",
	1: "UNKNOWN",
	2: "AND",
	3: "AND_NOT",
}
var UserListCombinedRuleOperatorEnum_UserListCombinedRuleOperator_value = map[string]int32{
	"UNSPECIFIED": 0,
	"UNKNOWN":     1,
	"AND":         2,
	"AND_NOT":     3,
}

func (x UserListCombinedRuleOperatorEnum_UserListCombinedRuleOperator) String() string {
	return proto.EnumName(UserListCombinedRuleOperatorEnum_UserListCombinedRuleOperator_name, int32(x))
}
func (UserListCombinedRuleOperatorEnum_UserListCombinedRuleOperator) EnumDescriptor() ([]byte, []int) {
	return fileDescriptor_user_list_combined_rule_operator_fb43ff327143f06f, []int{0, 0}
}

// Logical operator connecting two rules.
type UserListCombinedRuleOperatorEnum struct {
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_unrecognized     []byte   `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *UserListCombinedRuleOperatorEnum) Reset()         { *m = UserListCombinedRuleOperatorEnum{} }
func (m *UserListCombinedRuleOperatorEnum) String() string { return proto.CompactTextString(m) }
func (*UserListCombinedRuleOperatorEnum) ProtoMessage()    {}
func (*UserListCombinedRuleOperatorEnum) Descriptor() ([]byte, []int) {
	return fileDescriptor_user_list_combined_rule_operator_fb43ff327143f06f, []int{0}
}
func (m *UserListCombinedRuleOperatorEnum) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_UserListCombinedRuleOperatorEnum.Unmarshal(m, b)
}
func (m *UserListCombinedRuleOperatorEnum) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_UserListCombinedRuleOperatorEnum.Marshal(b, m, deterministic)
}
func (dst *UserListCombinedRuleOperatorEnum) XXX_Merge(src proto.Message) {
	xxx_messageInfo_UserListCombinedRuleOperatorEnum.Merge(dst, src)
}
func (m *UserListCombinedRuleOperatorEnum) XXX_Size() int {
	return xxx_messageInfo_UserListCombinedRuleOperatorEnum.Size(m)
}
func (m *UserListCombinedRuleOperatorEnum) XXX_DiscardUnknown() {
	xxx_messageInfo_UserListCombinedRuleOperatorEnum.DiscardUnknown(m)
}

var xxx_messageInfo_UserListCombinedRuleOperatorEnum proto.InternalMessageInfo

func init() {
	proto.RegisterType((*UserListCombinedRuleOperatorEnum)(nil), "google.ads.googleads.v0.enums.UserListCombinedRuleOperatorEnum")
	proto.RegisterEnum("google.ads.googleads.v0.enums.UserListCombinedRuleOperatorEnum_UserListCombinedRuleOperator", UserListCombinedRuleOperatorEnum_UserListCombinedRuleOperator_name, UserListCombinedRuleOperatorEnum_UserListCombinedRuleOperator_value)
}

func init() {
	proto.RegisterFile("google/ads/googleads/v0/enums/user_list_combined_rule_operator.proto", fileDescriptor_user_list_combined_rule_operator_fb43ff327143f06f)
}

var fileDescriptor_user_list_combined_rule_operator_fb43ff327143f06f = []byte{
	// 302 bytes of a gzipped FileDescriptorProto
	0x1f, 0x8b, 0x08, 0x00, 0x00, 0x00, 0x00, 0x00, 0x02, 0xff, 0x7c, 0x90, 0xc1, 0x6a, 0xb3, 0x40,
	0x14, 0x85, 0xff, 0x18, 0xf8, 0x03, 0x93, 0x45, 0xc5, 0x75, 0x03, 0x4d, 0x1e, 0x60, 0x14, 0xba,
	0x9b, 0xae, 0x26, 0x31, 0x0d, 0xa1, 0x65, 0x12, 0xd2, 0x6a, 0xa1, 0x08, 0x62, 0x32, 0xc3, 0x20,
	0xa8, 0x23, 0x73, 0xd5, 0x07, 0xea, 0xb2, 0x8f, 0xd2, 0x47, 0xe9, 0x03, 0x74, 0x5d, 0x1c, 0x8d,
	0xbb, 0xba, 0x91, 0xc3, 0x3d, 0xc7, 0x6f, 0xee, 0x3d, 0xc8, 0x97, 0x4a, 0xc9, 0x4c, 0xb8, 0x09,
	0x07, 0xb7, 0x93, 0xad, 0x6a, 0x3c, 0x57, 0x14, 0x75, 0x0e, 0x6e, 0x0d, 0x42, 0xc7, 0x59, 0x0a,
	0x55, 0x7c, 0x51, 0xf9, 0x39, 0x2d, 0x04, 0x8f, 0x75, 0x9d, 0x89, 0x58, 0x95, 0x42, 0x27, 0x95,
	0xd2, 0xb8, 0xd4, 0xaa, 0x52, 0xce, 0xa2, 0xfb, 0x15, 0x27, 0x1c, 0xf0, 0x40, 0xc1, 0x8d, 0x87,
	0x0d, 0x65, 0xd5, 0xa0, 0xbb, 0x00, 0x84, 0x7e, 0x4e, 0xa1, 0xda, 0xf4, 0x98, 0x53, 0x9d, 0x89,
	0x43, 0x0f, 0xd9, 0x16, 0x75, 0xbe, 0x3a, 0xa1, 0xdb, 0xb1, 0x8c, 0x73, 0x83, 0xe6, 0x01, 0x7b,
	0x39, 0x6e, 0x37, 0xfb, 0xc7, 0xfd, 0xd6, 0xb7, 0xff, 0x39, 0x73, 0x34, 0x0b, 0xd8, 0x13, 0x3b,
	0xbc, 0x31, 0x7b, 0xe2, 0xcc, 0xd0, 0x94, 0x32, 0xdf, 0xb6, 0xda, 0x29, 0x65, 0x7e, 0xcc, 0x0e,
	0xaf, 0xf6, 0x74, 0xfd, 0x33, 0x41, 0xcb, 0x8b, 0xca, 0xf1, 0xe8, 0x76, 0xeb, 0xe5, 0xd8, 0xbb,
	0xc7, 0xf6, 0xbe, 0xe3, 0xe4, 0x7d, 0xdd, 0x33, 0xa4, 0xca, 0x92, 0x42, 0x62, 0xa5, 0xa5, 0x2b,
	0x45, 0x61, 0xae, 0xbf, 0xf6, 0x56, 0xa6, 0xf0, 0x47, 0x8d, 0x0f, 0xe6, 0xfb, 0x61, 0x4d, 0x77,
	0x94, 0x7e, 0x5a, 0x8b, 0x5d, 0x87, 0xa2, 0x1c, 0x70, 0x27, 0x5b, 0x15, 0x7a, 0xb8, 0xed, 0x01,
	0xbe, 0xae, 0x7e, 0x44, 0x39, 0x44, 0x83, 0x1f, 0x85, 0x5e, 0x64, 0xfc, 0x6f, 0x6b, 0xd9, 0x0d,
	0x09, 0xa1, 0x1c, 0x08, 0x19, 0x12, 0x84, 0x84, 0x1e, 0x21, 0x26, 0x73, 0xfe, 0x6f, 0x16, 0xbb,
	0xff, 0x0d, 0x00, 0x00, 0xff, 0xff, 0x3f, 0x58, 0x29, 0x67, 0xde, 0x01, 0x00, 0x00,
}
