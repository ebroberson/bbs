// Code generated by protoc-gen-gogo. DO NOT EDIT.
// source: metric_tags.proto

package models

import proto "github.com/gogo/protobuf/proto"
import fmt "fmt"
import math "math"
import _ "github.com/gogo/protobuf/gogoproto"

import strconv "strconv"

import strings "strings"
import reflect "reflect"

import io "io"

// Reference imports to suppress errors if they are not otherwise used.
var _ = proto.Marshal
var _ = fmt.Errorf
var _ = math.Inf

// This is a compile-time assertion to ensure that this generated file
// is compatible with the proto package it is being compiled against.
// A compilation error at this line likely means your copy of the
// proto package needs to be updated.
const _ = proto.GoGoProtoPackageIsVersion2 // please upgrade the proto package

type MetricTagValue_DynamicValue int32

const (
	DynamicValueInvalid               MetricTagValue_DynamicValue = 0
	MetricTagDynamicValueIndex        MetricTagValue_DynamicValue = 1
	MetricTagDynamicValueInstanceGuid MetricTagValue_DynamicValue = 2
)

var MetricTagValue_DynamicValue_name = map[int32]string{
	0: "DynamicValueInvalid",
	1: "INDEX",
	2: "INSTANCE_GUID",
}
var MetricTagValue_DynamicValue_value = map[string]int32{
	"DynamicValueInvalid": 0,
	"INDEX":               1,
	"INSTANCE_GUID":       2,
}

func (MetricTagValue_DynamicValue) EnumDescriptor() ([]byte, []int) {
	return fileDescriptor_metric_tags_dd0f8b0400bb507b, []int{0, 0}
}

type MetricTagValue struct {
	// Note: we only expect one of the following set of fields to be
	// set.
	Static  string                      `protobuf:"bytes,1,opt,name=static,proto3" json:"static,omitempty"`
	Dynamic MetricTagValue_DynamicValue `protobuf:"varint,2,opt,name=dynamic,proto3,enum=models.MetricTagValue_DynamicValue" json:"dynamic,omitempty"`
}

func (m *MetricTagValue) Reset()      { *m = MetricTagValue{} }
func (*MetricTagValue) ProtoMessage() {}
func (*MetricTagValue) Descriptor() ([]byte, []int) {
	return fileDescriptor_metric_tags_dd0f8b0400bb507b, []int{0}
}
func (m *MetricTagValue) XXX_Unmarshal(b []byte) error {
	return m.Unmarshal(b)
}
func (m *MetricTagValue) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	if deterministic {
		return xxx_messageInfo_MetricTagValue.Marshal(b, m, deterministic)
	} else {
		b = b[:cap(b)]
		n, err := m.MarshalTo(b)
		if err != nil {
			return nil, err
		}
		return b[:n], nil
	}
}
func (dst *MetricTagValue) XXX_Merge(src proto.Message) {
	xxx_messageInfo_MetricTagValue.Merge(dst, src)
}
func (m *MetricTagValue) XXX_Size() int {
	return m.Size()
}
func (m *MetricTagValue) XXX_DiscardUnknown() {
	xxx_messageInfo_MetricTagValue.DiscardUnknown(m)
}

var xxx_messageInfo_MetricTagValue proto.InternalMessageInfo

func (m *MetricTagValue) GetStatic() string {
	if m != nil {
		return m.Static
	}
	return ""
}

func (m *MetricTagValue) GetDynamic() MetricTagValue_DynamicValue {
	if m != nil {
		return m.Dynamic
	}
	return DynamicValueInvalid
}

func init() {
	proto.RegisterType((*MetricTagValue)(nil), "models.MetricTagValue")
	proto.RegisterEnum("models.MetricTagValue_DynamicValue", MetricTagValue_DynamicValue_name, MetricTagValue_DynamicValue_value)
}
func (x MetricTagValue_DynamicValue) String() string {
	s, ok := MetricTagValue_DynamicValue_name[int32(x)]
	if ok {
		return s
	}
	return strconv.Itoa(int(x))
}
func (this *MetricTagValue) Equal(that interface{}) bool {
	if that == nil {
		return this == nil
	}

	that1, ok := that.(*MetricTagValue)
	if !ok {
		that2, ok := that.(MetricTagValue)
		if ok {
			that1 = &that2
		} else {
			return false
		}
	}
	if that1 == nil {
		return this == nil
	} else if this == nil {
		return false
	}
	if this.Static != that1.Static {
		return false
	}
	if this.Dynamic != that1.Dynamic {
		return false
	}
	return true
}
func (this *MetricTagValue) GoString() string {
	if this == nil {
		return "nil"
	}
	s := make([]string, 0, 6)
	s = append(s, "&models.MetricTagValue{")
	s = append(s, "Static: "+fmt.Sprintf("%#v", this.Static)+",\n")
	s = append(s, "Dynamic: "+fmt.Sprintf("%#v", this.Dynamic)+",\n")
	s = append(s, "}")
	return strings.Join(s, "")
}
func valueToGoStringMetricTags(v interface{}, typ string) string {
	rv := reflect.ValueOf(v)
	if rv.IsNil() {
		return "nil"
	}
	pv := reflect.Indirect(rv).Interface()
	return fmt.Sprintf("func(v %v) *%v { return &v } ( %#v )", typ, typ, pv)
}
func (m *MetricTagValue) Marshal() (dAtA []byte, err error) {
	size := m.Size()
	dAtA = make([]byte, size)
	n, err := m.MarshalTo(dAtA)
	if err != nil {
		return nil, err
	}
	return dAtA[:n], nil
}

func (m *MetricTagValue) MarshalTo(dAtA []byte) (int, error) {
	var i int
	_ = i
	var l int
	_ = l
	if len(m.Static) > 0 {
		dAtA[i] = 0xa
		i++
		i = encodeVarintMetricTags(dAtA, i, uint64(len(m.Static)))
		i += copy(dAtA[i:], m.Static)
	}
	if m.Dynamic != 0 {
		dAtA[i] = 0x10
		i++
		i = encodeVarintMetricTags(dAtA, i, uint64(m.Dynamic))
	}
	return i, nil
}

func encodeVarintMetricTags(dAtA []byte, offset int, v uint64) int {
	for v >= 1<<7 {
		dAtA[offset] = uint8(v&0x7f | 0x80)
		v >>= 7
		offset++
	}
	dAtA[offset] = uint8(v)
	return offset + 1
}
func (m *MetricTagValue) Size() (n int) {
	if m == nil {
		return 0
	}
	var l int
	_ = l
	l = len(m.Static)
	if l > 0 {
		n += 1 + l + sovMetricTags(uint64(l))
	}
	if m.Dynamic != 0 {
		n += 1 + sovMetricTags(uint64(m.Dynamic))
	}
	return n
}

func sovMetricTags(x uint64) (n int) {
	for {
		n++
		x >>= 7
		if x == 0 {
			break
		}
	}
	return n
}
func sozMetricTags(x uint64) (n int) {
	return sovMetricTags(uint64((x << 1) ^ uint64((int64(x) >> 63))))
}
func (this *MetricTagValue) String() string {
	if this == nil {
		return "nil"
	}
	s := strings.Join([]string{`&MetricTagValue{`,
		`Static:` + fmt.Sprintf("%v", this.Static) + `,`,
		`Dynamic:` + fmt.Sprintf("%v", this.Dynamic) + `,`,
		`}`,
	}, "")
	return s
}
func valueToStringMetricTags(v interface{}) string {
	rv := reflect.ValueOf(v)
	if rv.IsNil() {
		return "nil"
	}
	pv := reflect.Indirect(rv).Interface()
	return fmt.Sprintf("*%v", pv)
}
func (m *MetricTagValue) Unmarshal(dAtA []byte) error {
	l := len(dAtA)
	iNdEx := 0
	for iNdEx < l {
		preIndex := iNdEx
		var wire uint64
		for shift := uint(0); ; shift += 7 {
			if shift >= 64 {
				return ErrIntOverflowMetricTags
			}
			if iNdEx >= l {
				return io.ErrUnexpectedEOF
			}
			b := dAtA[iNdEx]
			iNdEx++
			wire |= (uint64(b) & 0x7F) << shift
			if b < 0x80 {
				break
			}
		}
		fieldNum := int32(wire >> 3)
		wireType := int(wire & 0x7)
		if wireType == 4 {
			return fmt.Errorf("proto: MetricTagValue: wiretype end group for non-group")
		}
		if fieldNum <= 0 {
			return fmt.Errorf("proto: MetricTagValue: illegal tag %d (wire type %d)", fieldNum, wire)
		}
		switch fieldNum {
		case 1:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field Static", wireType)
			}
			var stringLen uint64
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowMetricTags
				}
				if iNdEx >= l {
					return io.ErrUnexpectedEOF
				}
				b := dAtA[iNdEx]
				iNdEx++
				stringLen |= (uint64(b) & 0x7F) << shift
				if b < 0x80 {
					break
				}
			}
			intStringLen := int(stringLen)
			if intStringLen < 0 {
				return ErrInvalidLengthMetricTags
			}
			postIndex := iNdEx + intStringLen
			if postIndex > l {
				return io.ErrUnexpectedEOF
			}
			m.Static = string(dAtA[iNdEx:postIndex])
			iNdEx = postIndex
		case 2:
			if wireType != 0 {
				return fmt.Errorf("proto: wrong wireType = %d for field Dynamic", wireType)
			}
			m.Dynamic = 0
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowMetricTags
				}
				if iNdEx >= l {
					return io.ErrUnexpectedEOF
				}
				b := dAtA[iNdEx]
				iNdEx++
				m.Dynamic |= (MetricTagValue_DynamicValue(b) & 0x7F) << shift
				if b < 0x80 {
					break
				}
			}
		default:
			iNdEx = preIndex
			skippy, err := skipMetricTags(dAtA[iNdEx:])
			if err != nil {
				return err
			}
			if skippy < 0 {
				return ErrInvalidLengthMetricTags
			}
			if (iNdEx + skippy) > l {
				return io.ErrUnexpectedEOF
			}
			iNdEx += skippy
		}
	}

	if iNdEx > l {
		return io.ErrUnexpectedEOF
	}
	return nil
}
func skipMetricTags(dAtA []byte) (n int, err error) {
	l := len(dAtA)
	iNdEx := 0
	for iNdEx < l {
		var wire uint64
		for shift := uint(0); ; shift += 7 {
			if shift >= 64 {
				return 0, ErrIntOverflowMetricTags
			}
			if iNdEx >= l {
				return 0, io.ErrUnexpectedEOF
			}
			b := dAtA[iNdEx]
			iNdEx++
			wire |= (uint64(b) & 0x7F) << shift
			if b < 0x80 {
				break
			}
		}
		wireType := int(wire & 0x7)
		switch wireType {
		case 0:
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return 0, ErrIntOverflowMetricTags
				}
				if iNdEx >= l {
					return 0, io.ErrUnexpectedEOF
				}
				iNdEx++
				if dAtA[iNdEx-1] < 0x80 {
					break
				}
			}
			return iNdEx, nil
		case 1:
			iNdEx += 8
			return iNdEx, nil
		case 2:
			var length int
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return 0, ErrIntOverflowMetricTags
				}
				if iNdEx >= l {
					return 0, io.ErrUnexpectedEOF
				}
				b := dAtA[iNdEx]
				iNdEx++
				length |= (int(b) & 0x7F) << shift
				if b < 0x80 {
					break
				}
			}
			iNdEx += length
			if length < 0 {
				return 0, ErrInvalidLengthMetricTags
			}
			return iNdEx, nil
		case 3:
			for {
				var innerWire uint64
				var start int = iNdEx
				for shift := uint(0); ; shift += 7 {
					if shift >= 64 {
						return 0, ErrIntOverflowMetricTags
					}
					if iNdEx >= l {
						return 0, io.ErrUnexpectedEOF
					}
					b := dAtA[iNdEx]
					iNdEx++
					innerWire |= (uint64(b) & 0x7F) << shift
					if b < 0x80 {
						break
					}
				}
				innerWireType := int(innerWire & 0x7)
				if innerWireType == 4 {
					break
				}
				next, err := skipMetricTags(dAtA[start:])
				if err != nil {
					return 0, err
				}
				iNdEx = start + next
			}
			return iNdEx, nil
		case 4:
			return iNdEx, nil
		case 5:
			iNdEx += 4
			return iNdEx, nil
		default:
			return 0, fmt.Errorf("proto: illegal wireType %d", wireType)
		}
	}
	panic("unreachable")
}

var (
	ErrInvalidLengthMetricTags = fmt.Errorf("proto: negative length found during unmarshaling")
	ErrIntOverflowMetricTags   = fmt.Errorf("proto: integer overflow")
)

func init() { proto.RegisterFile("metric_tags.proto", fileDescriptor_metric_tags_dd0f8b0400bb507b) }

var fileDescriptor_metric_tags_dd0f8b0400bb507b = []byte{
	// 296 bytes of a gzipped FileDescriptorProto
	0x1f, 0x8b, 0x08, 0x00, 0x00, 0x00, 0x00, 0x00, 0x02, 0xff, 0xe2, 0x12, 0xcc, 0x4d, 0x2d, 0x29,
	0xca, 0x4c, 0x8e, 0x2f, 0x49, 0x4c, 0x2f, 0xd6, 0x2b, 0x28, 0xca, 0x2f, 0xc9, 0x17, 0x62, 0xcb,
	0xcd, 0x4f, 0x49, 0xcd, 0x29, 0x96, 0xd2, 0x4d, 0xcf, 0x2c, 0xc9, 0x28, 0x4d, 0xd2, 0x4b, 0xce,
	0xcf, 0xd5, 0x4f, 0xcf, 0x4f, 0xcf, 0xd7, 0x07, 0x4b, 0x27, 0x95, 0xa6, 0x81, 0x79, 0x60, 0x0e,
	0x98, 0x05, 0xd1, 0xa6, 0xf4, 0x8d, 0x91, 0x8b, 0xcf, 0x17, 0x6c, 0x58, 0x48, 0x62, 0x7a, 0x58,
	0x62, 0x4e, 0x69, 0xaa, 0x90, 0x18, 0x17, 0x5b, 0x71, 0x49, 0x62, 0x49, 0x66, 0xb2, 0x04, 0xa3,
	0x02, 0xa3, 0x06, 0x67, 0x10, 0x94, 0x27, 0x64, 0xcb, 0xc5, 0x9e, 0x52, 0x99, 0x97, 0x98, 0x9b,
	0x99, 0x2c, 0xc1, 0xa4, 0xc0, 0xa8, 0xc1, 0x67, 0xa4, 0xac, 0x07, 0xb1, 0x53, 0x0f, 0xd5, 0x00,
	0x3d, 0x17, 0x88, 0x2a, 0x30, 0x27, 0x08, 0xa6, 0x47, 0xa9, 0x87, 0x91, 0x8b, 0x07, 0x59, 0x46,
	0x48, 0x9c, 0x4b, 0x18, 0x99, 0xef, 0x99, 0x57, 0x96, 0x98, 0x93, 0x99, 0x22, 0xc0, 0x20, 0xa4,
	0xc9, 0xc5, 0xea, 0xe9, 0xe7, 0xe2, 0x1a, 0x21, 0xc0, 0x28, 0x25, 0xd7, 0x35, 0x57, 0x41, 0x0a,
	0x6e, 0x3c, 0xaa, 0xf2, 0x94, 0xd4, 0x0a, 0x21, 0x0b, 0x2e, 0x5e, 0x4f, 0xbf, 0xe0, 0x10, 0x47,
	0x3f, 0x67, 0xd7, 0x78, 0xf7, 0x50, 0x4f, 0x17, 0x01, 0x26, 0x29, 0xd5, 0xae, 0xb9, 0x0a, 0x8a,
	0x38, 0xb4, 0x14, 0x97, 0x24, 0xe6, 0x25, 0xa7, 0xba, 0x97, 0x66, 0xa6, 0x38, 0x99, 0x5c, 0x78,
	0x28, 0xc7, 0x70, 0xe3, 0xa1, 0x1c, 0xc3, 0x87, 0x87, 0x72, 0x8c, 0x0d, 0x8f, 0xe4, 0x18, 0x57,
	0x3c, 0x92, 0x63, 0x3c, 0xf1, 0x48, 0x8e, 0xf1, 0xc2, 0x23, 0x39, 0xc6, 0x07, 0x8f, 0xe4, 0x18,
	0x5f, 0x3c, 0x92, 0x63, 0xf8, 0xf0, 0x48, 0x8e, 0x71, 0xc2, 0x63, 0x39, 0x86, 0x0b, 0x8f, 0xe5,
	0x18, 0x6e, 0x3c, 0x96, 0x63, 0x48, 0x62, 0x03, 0x87, 0x9a, 0x31, 0x20, 0x00, 0x00, 0xff, 0xff,
	0x9b, 0x83, 0x77, 0xc8, 0x81, 0x01, 0x00, 0x00,
}
