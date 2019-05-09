// Code generated by protoc-gen-gogo. DO NOT EDIT.
// source: fileformat.proto

/*
	Package fileformat is a generated protocol buffer package.

	It is generated from these files:
		fileformat.proto

	It has these top-level messages:
		Body
		Meta
		Feature
		Tag
*/
package fileformat

import proto "github.com/gogo/protobuf/proto"
import fmt "fmt"
import math "math"

import binary "encoding/binary"

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

type Feature_GeomType int32

const (
	Feature_UNKNOWN Feature_GeomType = 0
	Feature_POINT   Feature_GeomType = 1
	Feature_LINE    Feature_GeomType = 2
	Feature_POLYGON Feature_GeomType = 3
)

var Feature_GeomType_name = map[int32]string{
	0: "UNKNOWN",
	1: "POINT",
	2: "LINE",
	3: "POLYGON",
}
var Feature_GeomType_value = map[string]int32{
	"UNKNOWN": 0,
	"POINT":   1,
	"LINE":    2,
	"POLYGON": 3,
}

func (x Feature_GeomType) String() string {
	return proto.EnumName(Feature_GeomType_name, int32(x))
}
func (Feature_GeomType) EnumDescriptor() ([]byte, []int) { return fileDescriptorFileformat, []int{2, 0} }

type Feature_GeomSerialization int32

const (
	Feature_WKB Feature_GeomSerialization = 0
)

var Feature_GeomSerialization_name = map[int32]string{
	0: "WKB",
}
var Feature_GeomSerialization_value = map[string]int32{
	"WKB": 0,
}

func (x Feature_GeomSerialization) String() string {
	return proto.EnumName(Feature_GeomSerialization_name, int32(x))
}
func (Feature_GeomSerialization) EnumDescriptor() ([]byte, []int) {
	return fileDescriptorFileformat, []int{2, 1}
}

type Tag_ValueType int32

const (
	Tag_STRING Tag_ValueType = 0
	Tag_INT    Tag_ValueType = 1
	Tag_DOUBLE Tag_ValueType = 2
)

var Tag_ValueType_name = map[int32]string{
	0: "STRING",
	1: "INT",
	2: "DOUBLE",
}
var Tag_ValueType_value = map[string]int32{
	"STRING": 0,
	"INT":    1,
	"DOUBLE": 2,
}

func (x Tag_ValueType) String() string {
	return proto.EnumName(Tag_ValueType_name, int32(x))
}
func (Tag_ValueType) EnumDescriptor() ([]byte, []int) { return fileDescriptorFileformat, []int{3, 0} }

type Body struct {
	Meta    *Meta      `protobuf:"bytes,1,opt,name=meta" json:"meta,omitempty"`
	Feature []*Feature `protobuf:"bytes,2,rep,name=feature" json:"feature,omitempty"`
}

func (m *Body) Reset()                    { *m = Body{} }
func (m *Body) String() string            { return proto.CompactTextString(m) }
func (*Body) ProtoMessage()               {}
func (*Body) Descriptor() ([]byte, []int) { return fileDescriptorFileformat, []int{0} }

func (m *Body) GetMeta() *Meta {
	if m != nil {
		return m.Meta
	}
	return nil
}

func (m *Body) GetFeature() []*Feature {
	if m != nil {
		return m.Feature
	}
	return nil
}

type Meta struct {
	Tags []*Tag `protobuf:"bytes,1,rep,name=tags" json:"tags,omitempty"`
}

func (m *Meta) Reset()                    { *m = Meta{} }
func (m *Meta) String() string            { return proto.CompactTextString(m) }
func (*Meta) ProtoMessage()               {}
func (*Meta) Descriptor() ([]byte, []int) { return fileDescriptorFileformat, []int{1} }

func (m *Meta) GetTags() []*Tag {
	if m != nil {
		return m.Tags
	}
	return nil
}

type Feature struct {
	Geomtype   Feature_GeomType          `protobuf:"varint,1,opt,name=geomtype,proto3,enum=Feature_GeomType" json:"geomtype,omitempty"`
	Geomserial Feature_GeomSerialization `protobuf:"varint,2,opt,name=geomserial,proto3,enum=Feature_GeomSerialization" json:"geomserial,omitempty"`
	Geom       []byte                    `protobuf:"bytes,3,opt,name=geom,proto3" json:"geom,omitempty"`
	// geometry bbox
	Left   float64 `protobuf:"fixed64,4,opt,name=left,proto3" json:"left,omitempty"`
	Right  float64 `protobuf:"fixed64,5,opt,name=right,proto3" json:"right,omitempty"`
	Top    float64 `protobuf:"fixed64,6,opt,name=top,proto3" json:"top,omitempty"`
	Bottom float64 `protobuf:"fixed64,7,opt,name=bottom,proto3" json:"bottom,omitempty"`
	Tags   []*Tag  `protobuf:"bytes,8,rep,name=tags" json:"tags,omitempty"`
}

func (m *Feature) Reset()                    { *m = Feature{} }
func (m *Feature) String() string            { return proto.CompactTextString(m) }
func (*Feature) ProtoMessage()               {}
func (*Feature) Descriptor() ([]byte, []int) { return fileDescriptorFileformat, []int{2} }

func (m *Feature) GetGeomtype() Feature_GeomType {
	if m != nil {
		return m.Geomtype
	}
	return Feature_UNKNOWN
}

func (m *Feature) GetGeomserial() Feature_GeomSerialization {
	if m != nil {
		return m.Geomserial
	}
	return Feature_WKB
}

func (m *Feature) GetGeom() []byte {
	if m != nil {
		return m.Geom
	}
	return nil
}

func (m *Feature) GetLeft() float64 {
	if m != nil {
		return m.Left
	}
	return 0
}

func (m *Feature) GetRight() float64 {
	if m != nil {
		return m.Right
	}
	return 0
}

func (m *Feature) GetTop() float64 {
	if m != nil {
		return m.Top
	}
	return 0
}

func (m *Feature) GetBottom() float64 {
	if m != nil {
		return m.Bottom
	}
	return 0
}

func (m *Feature) GetTags() []*Tag {
	if m != nil {
		return m.Tags
	}
	return nil
}

type Tag struct {
	Key   string        `protobuf:"bytes,1,opt,name=key,proto3" json:"key,omitempty"`
	Value []byte        `protobuf:"bytes,2,opt,name=value,proto3" json:"value,omitempty"`
	Type  Tag_ValueType `protobuf:"varint,3,opt,name=type,proto3,enum=Tag_ValueType" json:"type,omitempty"`
}

func (m *Tag) Reset()                    { *m = Tag{} }
func (m *Tag) String() string            { return proto.CompactTextString(m) }
func (*Tag) ProtoMessage()               {}
func (*Tag) Descriptor() ([]byte, []int) { return fileDescriptorFileformat, []int{3} }

func (m *Tag) GetKey() string {
	if m != nil {
		return m.Key
	}
	return ""
}

func (m *Tag) GetValue() []byte {
	if m != nil {
		return m.Value
	}
	return nil
}

func (m *Tag) GetType() Tag_ValueType {
	if m != nil {
		return m.Type
	}
	return Tag_STRING
}

func init() {
	proto.RegisterType((*Body)(nil), "Body")
	proto.RegisterType((*Meta)(nil), "Meta")
	proto.RegisterType((*Feature)(nil), "Feature")
	proto.RegisterType((*Tag)(nil), "Tag")
	proto.RegisterEnum("Feature_GeomType", Feature_GeomType_name, Feature_GeomType_value)
	proto.RegisterEnum("Feature_GeomSerialization", Feature_GeomSerialization_name, Feature_GeomSerialization_value)
	proto.RegisterEnum("Tag_ValueType", Tag_ValueType_name, Tag_ValueType_value)
}
func (m *Body) Marshal() (dAtA []byte, err error) {
	size := m.Size()
	dAtA = make([]byte, size)
	n, err := m.MarshalTo(dAtA)
	if err != nil {
		return nil, err
	}
	return dAtA[:n], nil
}

func (m *Body) MarshalTo(dAtA []byte) (int, error) {
	var i int
	_ = i
	var l int
	_ = l
	if m.Meta != nil {
		dAtA[i] = 0xa
		i++
		i = encodeVarintFileformat(dAtA, i, uint64(m.Meta.Size()))
		n1, err := m.Meta.MarshalTo(dAtA[i:])
		if err != nil {
			return 0, err
		}
		i += n1
	}
	if len(m.Feature) > 0 {
		for _, msg := range m.Feature {
			dAtA[i] = 0x12
			i++
			i = encodeVarintFileformat(dAtA, i, uint64(msg.Size()))
			n, err := msg.MarshalTo(dAtA[i:])
			if err != nil {
				return 0, err
			}
			i += n
		}
	}
	return i, nil
}

func (m *Meta) Marshal() (dAtA []byte, err error) {
	size := m.Size()
	dAtA = make([]byte, size)
	n, err := m.MarshalTo(dAtA)
	if err != nil {
		return nil, err
	}
	return dAtA[:n], nil
}

func (m *Meta) MarshalTo(dAtA []byte) (int, error) {
	var i int
	_ = i
	var l int
	_ = l
	if len(m.Tags) > 0 {
		for _, msg := range m.Tags {
			dAtA[i] = 0xa
			i++
			i = encodeVarintFileformat(dAtA, i, uint64(msg.Size()))
			n, err := msg.MarshalTo(dAtA[i:])
			if err != nil {
				return 0, err
			}
			i += n
		}
	}
	return i, nil
}

func (m *Feature) Marshal() (dAtA []byte, err error) {
	size := m.Size()
	dAtA = make([]byte, size)
	n, err := m.MarshalTo(dAtA)
	if err != nil {
		return nil, err
	}
	return dAtA[:n], nil
}

func (m *Feature) MarshalTo(dAtA []byte) (int, error) {
	var i int
	_ = i
	var l int
	_ = l
	if m.Geomtype != 0 {
		dAtA[i] = 0x8
		i++
		i = encodeVarintFileformat(dAtA, i, uint64(m.Geomtype))
	}
	if m.Geomserial != 0 {
		dAtA[i] = 0x10
		i++
		i = encodeVarintFileformat(dAtA, i, uint64(m.Geomserial))
	}
	if len(m.Geom) > 0 {
		dAtA[i] = 0x1a
		i++
		i = encodeVarintFileformat(dAtA, i, uint64(len(m.Geom)))
		i += copy(dAtA[i:], m.Geom)
	}
	if m.Left != 0 {
		dAtA[i] = 0x21
		i++
		binary.LittleEndian.PutUint64(dAtA[i:], uint64(math.Float64bits(float64(m.Left))))
		i += 8
	}
	if m.Right != 0 {
		dAtA[i] = 0x29
		i++
		binary.LittleEndian.PutUint64(dAtA[i:], uint64(math.Float64bits(float64(m.Right))))
		i += 8
	}
	if m.Top != 0 {
		dAtA[i] = 0x31
		i++
		binary.LittleEndian.PutUint64(dAtA[i:], uint64(math.Float64bits(float64(m.Top))))
		i += 8
	}
	if m.Bottom != 0 {
		dAtA[i] = 0x39
		i++
		binary.LittleEndian.PutUint64(dAtA[i:], uint64(math.Float64bits(float64(m.Bottom))))
		i += 8
	}
	if len(m.Tags) > 0 {
		for _, msg := range m.Tags {
			dAtA[i] = 0x42
			i++
			i = encodeVarintFileformat(dAtA, i, uint64(msg.Size()))
			n, err := msg.MarshalTo(dAtA[i:])
			if err != nil {
				return 0, err
			}
			i += n
		}
	}
	return i, nil
}

func (m *Tag) Marshal() (dAtA []byte, err error) {
	size := m.Size()
	dAtA = make([]byte, size)
	n, err := m.MarshalTo(dAtA)
	if err != nil {
		return nil, err
	}
	return dAtA[:n], nil
}

func (m *Tag) MarshalTo(dAtA []byte) (int, error) {
	var i int
	_ = i
	var l int
	_ = l
	if len(m.Key) > 0 {
		dAtA[i] = 0xa
		i++
		i = encodeVarintFileformat(dAtA, i, uint64(len(m.Key)))
		i += copy(dAtA[i:], m.Key)
	}
	if len(m.Value) > 0 {
		dAtA[i] = 0x12
		i++
		i = encodeVarintFileformat(dAtA, i, uint64(len(m.Value)))
		i += copy(dAtA[i:], m.Value)
	}
	if m.Type != 0 {
		dAtA[i] = 0x18
		i++
		i = encodeVarintFileformat(dAtA, i, uint64(m.Type))
	}
	return i, nil
}

func encodeVarintFileformat(dAtA []byte, offset int, v uint64) int {
	for v >= 1<<7 {
		dAtA[offset] = uint8(v&0x7f | 0x80)
		v >>= 7
		offset++
	}
	dAtA[offset] = uint8(v)
	return offset + 1
}
func (m *Body) Size() (n int) {
	var l int
	_ = l
	if m.Meta != nil {
		l = m.Meta.Size()
		n += 1 + l + sovFileformat(uint64(l))
	}
	if len(m.Feature) > 0 {
		for _, e := range m.Feature {
			l = e.Size()
			n += 1 + l + sovFileformat(uint64(l))
		}
	}
	return n
}

func (m *Meta) Size() (n int) {
	var l int
	_ = l
	if len(m.Tags) > 0 {
		for _, e := range m.Tags {
			l = e.Size()
			n += 1 + l + sovFileformat(uint64(l))
		}
	}
	return n
}

func (m *Feature) Size() (n int) {
	var l int
	_ = l
	if m.Geomtype != 0 {
		n += 1 + sovFileformat(uint64(m.Geomtype))
	}
	if m.Geomserial != 0 {
		n += 1 + sovFileformat(uint64(m.Geomserial))
	}
	l = len(m.Geom)
	if l > 0 {
		n += 1 + l + sovFileformat(uint64(l))
	}
	if m.Left != 0 {
		n += 9
	}
	if m.Right != 0 {
		n += 9
	}
	if m.Top != 0 {
		n += 9
	}
	if m.Bottom != 0 {
		n += 9
	}
	if len(m.Tags) > 0 {
		for _, e := range m.Tags {
			l = e.Size()
			n += 1 + l + sovFileformat(uint64(l))
		}
	}
	return n
}

func (m *Tag) Size() (n int) {
	var l int
	_ = l
	l = len(m.Key)
	if l > 0 {
		n += 1 + l + sovFileformat(uint64(l))
	}
	l = len(m.Value)
	if l > 0 {
		n += 1 + l + sovFileformat(uint64(l))
	}
	if m.Type != 0 {
		n += 1 + sovFileformat(uint64(m.Type))
	}
	return n
}

func sovFileformat(x uint64) (n int) {
	for {
		n++
		x >>= 7
		if x == 0 {
			break
		}
	}
	return n
}
func sozFileformat(x uint64) (n int) {
	return sovFileformat(uint64((x << 1) ^ uint64((int64(x) >> 63))))
}
func (m *Body) Unmarshal(dAtA []byte) error {
	l := len(dAtA)
	iNdEx := 0
	for iNdEx < l {
		preIndex := iNdEx
		var wire uint64
		for shift := uint(0); ; shift += 7 {
			if shift >= 64 {
				return ErrIntOverflowFileformat
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
			return fmt.Errorf("proto: Body: wiretype end group for non-group")
		}
		if fieldNum <= 0 {
			return fmt.Errorf("proto: Body: illegal tag %d (wire type %d)", fieldNum, wire)
		}
		switch fieldNum {
		case 1:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field Meta", wireType)
			}
			var msglen int
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowFileformat
				}
				if iNdEx >= l {
					return io.ErrUnexpectedEOF
				}
				b := dAtA[iNdEx]
				iNdEx++
				msglen |= (int(b) & 0x7F) << shift
				if b < 0x80 {
					break
				}
			}
			if msglen < 0 {
				return ErrInvalidLengthFileformat
			}
			postIndex := iNdEx + msglen
			if postIndex > l {
				return io.ErrUnexpectedEOF
			}
			if m.Meta == nil {
				m.Meta = &Meta{}
			}
			if err := m.Meta.Unmarshal(dAtA[iNdEx:postIndex]); err != nil {
				return err
			}
			iNdEx = postIndex
		case 2:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field Feature", wireType)
			}
			var msglen int
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowFileformat
				}
				if iNdEx >= l {
					return io.ErrUnexpectedEOF
				}
				b := dAtA[iNdEx]
				iNdEx++
				msglen |= (int(b) & 0x7F) << shift
				if b < 0x80 {
					break
				}
			}
			if msglen < 0 {
				return ErrInvalidLengthFileformat
			}
			postIndex := iNdEx + msglen
			if postIndex > l {
				return io.ErrUnexpectedEOF
			}
			m.Feature = append(m.Feature, &Feature{})
			if err := m.Feature[len(m.Feature)-1].Unmarshal(dAtA[iNdEx:postIndex]); err != nil {
				return err
			}
			iNdEx = postIndex
		default:
			iNdEx = preIndex
			skippy, err := skipFileformat(dAtA[iNdEx:])
			if err != nil {
				return err
			}
			if skippy < 0 {
				return ErrInvalidLengthFileformat
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
func (m *Meta) Unmarshal(dAtA []byte) error {
	l := len(dAtA)
	iNdEx := 0
	for iNdEx < l {
		preIndex := iNdEx
		var wire uint64
		for shift := uint(0); ; shift += 7 {
			if shift >= 64 {
				return ErrIntOverflowFileformat
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
			return fmt.Errorf("proto: Meta: wiretype end group for non-group")
		}
		if fieldNum <= 0 {
			return fmt.Errorf("proto: Meta: illegal tag %d (wire type %d)", fieldNum, wire)
		}
		switch fieldNum {
		case 1:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field Tags", wireType)
			}
			var msglen int
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowFileformat
				}
				if iNdEx >= l {
					return io.ErrUnexpectedEOF
				}
				b := dAtA[iNdEx]
				iNdEx++
				msglen |= (int(b) & 0x7F) << shift
				if b < 0x80 {
					break
				}
			}
			if msglen < 0 {
				return ErrInvalidLengthFileformat
			}
			postIndex := iNdEx + msglen
			if postIndex > l {
				return io.ErrUnexpectedEOF
			}
			m.Tags = append(m.Tags, &Tag{})
			if err := m.Tags[len(m.Tags)-1].Unmarshal(dAtA[iNdEx:postIndex]); err != nil {
				return err
			}
			iNdEx = postIndex
		default:
			iNdEx = preIndex
			skippy, err := skipFileformat(dAtA[iNdEx:])
			if err != nil {
				return err
			}
			if skippy < 0 {
				return ErrInvalidLengthFileformat
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
func (m *Feature) Unmarshal(dAtA []byte) error {
	l := len(dAtA)
	iNdEx := 0
	for iNdEx < l {
		preIndex := iNdEx
		var wire uint64
		for shift := uint(0); ; shift += 7 {
			if shift >= 64 {
				return ErrIntOverflowFileformat
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
			return fmt.Errorf("proto: Feature: wiretype end group for non-group")
		}
		if fieldNum <= 0 {
			return fmt.Errorf("proto: Feature: illegal tag %d (wire type %d)", fieldNum, wire)
		}
		switch fieldNum {
		case 1:
			if wireType != 0 {
				return fmt.Errorf("proto: wrong wireType = %d for field Geomtype", wireType)
			}
			m.Geomtype = 0
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowFileformat
				}
				if iNdEx >= l {
					return io.ErrUnexpectedEOF
				}
				b := dAtA[iNdEx]
				iNdEx++
				m.Geomtype |= (Feature_GeomType(b) & 0x7F) << shift
				if b < 0x80 {
					break
				}
			}
		case 2:
			if wireType != 0 {
				return fmt.Errorf("proto: wrong wireType = %d for field Geomserial", wireType)
			}
			m.Geomserial = 0
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowFileformat
				}
				if iNdEx >= l {
					return io.ErrUnexpectedEOF
				}
				b := dAtA[iNdEx]
				iNdEx++
				m.Geomserial |= (Feature_GeomSerialization(b) & 0x7F) << shift
				if b < 0x80 {
					break
				}
			}
		case 3:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field Geom", wireType)
			}
			var byteLen int
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowFileformat
				}
				if iNdEx >= l {
					return io.ErrUnexpectedEOF
				}
				b := dAtA[iNdEx]
				iNdEx++
				byteLen |= (int(b) & 0x7F) << shift
				if b < 0x80 {
					break
				}
			}
			if byteLen < 0 {
				return ErrInvalidLengthFileformat
			}
			postIndex := iNdEx + byteLen
			if postIndex > l {
				return io.ErrUnexpectedEOF
			}
			m.Geom = append(m.Geom[:0], dAtA[iNdEx:postIndex]...)
			if m.Geom == nil {
				m.Geom = []byte{}
			}
			iNdEx = postIndex
		case 4:
			if wireType != 1 {
				return fmt.Errorf("proto: wrong wireType = %d for field Left", wireType)
			}
			var v uint64
			if (iNdEx + 8) > l {
				return io.ErrUnexpectedEOF
			}
			v = uint64(binary.LittleEndian.Uint64(dAtA[iNdEx:]))
			iNdEx += 8
			m.Left = float64(math.Float64frombits(v))
		case 5:
			if wireType != 1 {
				return fmt.Errorf("proto: wrong wireType = %d for field Right", wireType)
			}
			var v uint64
			if (iNdEx + 8) > l {
				return io.ErrUnexpectedEOF
			}
			v = uint64(binary.LittleEndian.Uint64(dAtA[iNdEx:]))
			iNdEx += 8
			m.Right = float64(math.Float64frombits(v))
		case 6:
			if wireType != 1 {
				return fmt.Errorf("proto: wrong wireType = %d for field Top", wireType)
			}
			var v uint64
			if (iNdEx + 8) > l {
				return io.ErrUnexpectedEOF
			}
			v = uint64(binary.LittleEndian.Uint64(dAtA[iNdEx:]))
			iNdEx += 8
			m.Top = float64(math.Float64frombits(v))
		case 7:
			if wireType != 1 {
				return fmt.Errorf("proto: wrong wireType = %d for field Bottom", wireType)
			}
			var v uint64
			if (iNdEx + 8) > l {
				return io.ErrUnexpectedEOF
			}
			v = uint64(binary.LittleEndian.Uint64(dAtA[iNdEx:]))
			iNdEx += 8
			m.Bottom = float64(math.Float64frombits(v))
		case 8:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field Tags", wireType)
			}
			var msglen int
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowFileformat
				}
				if iNdEx >= l {
					return io.ErrUnexpectedEOF
				}
				b := dAtA[iNdEx]
				iNdEx++
				msglen |= (int(b) & 0x7F) << shift
				if b < 0x80 {
					break
				}
			}
			if msglen < 0 {
				return ErrInvalidLengthFileformat
			}
			postIndex := iNdEx + msglen
			if postIndex > l {
				return io.ErrUnexpectedEOF
			}
			m.Tags = append(m.Tags, &Tag{})
			if err := m.Tags[len(m.Tags)-1].Unmarshal(dAtA[iNdEx:postIndex]); err != nil {
				return err
			}
			iNdEx = postIndex
		default:
			iNdEx = preIndex
			skippy, err := skipFileformat(dAtA[iNdEx:])
			if err != nil {
				return err
			}
			if skippy < 0 {
				return ErrInvalidLengthFileformat
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
func (m *Tag) Unmarshal(dAtA []byte) error {
	l := len(dAtA)
	iNdEx := 0
	for iNdEx < l {
		preIndex := iNdEx
		var wire uint64
		for shift := uint(0); ; shift += 7 {
			if shift >= 64 {
				return ErrIntOverflowFileformat
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
			return fmt.Errorf("proto: Tag: wiretype end group for non-group")
		}
		if fieldNum <= 0 {
			return fmt.Errorf("proto: Tag: illegal tag %d (wire type %d)", fieldNum, wire)
		}
		switch fieldNum {
		case 1:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field Key", wireType)
			}
			var stringLen uint64
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowFileformat
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
				return ErrInvalidLengthFileformat
			}
			postIndex := iNdEx + intStringLen
			if postIndex > l {
				return io.ErrUnexpectedEOF
			}
			m.Key = string(dAtA[iNdEx:postIndex])
			iNdEx = postIndex
		case 2:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field Value", wireType)
			}
			var byteLen int
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowFileformat
				}
				if iNdEx >= l {
					return io.ErrUnexpectedEOF
				}
				b := dAtA[iNdEx]
				iNdEx++
				byteLen |= (int(b) & 0x7F) << shift
				if b < 0x80 {
					break
				}
			}
			if byteLen < 0 {
				return ErrInvalidLengthFileformat
			}
			postIndex := iNdEx + byteLen
			if postIndex > l {
				return io.ErrUnexpectedEOF
			}
			m.Value = append(m.Value[:0], dAtA[iNdEx:postIndex]...)
			if m.Value == nil {
				m.Value = []byte{}
			}
			iNdEx = postIndex
		case 3:
			if wireType != 0 {
				return fmt.Errorf("proto: wrong wireType = %d for field Type", wireType)
			}
			m.Type = 0
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowFileformat
				}
				if iNdEx >= l {
					return io.ErrUnexpectedEOF
				}
				b := dAtA[iNdEx]
				iNdEx++
				m.Type |= (Tag_ValueType(b) & 0x7F) << shift
				if b < 0x80 {
					break
				}
			}
		default:
			iNdEx = preIndex
			skippy, err := skipFileformat(dAtA[iNdEx:])
			if err != nil {
				return err
			}
			if skippy < 0 {
				return ErrInvalidLengthFileformat
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
func skipFileformat(dAtA []byte) (n int, err error) {
	l := len(dAtA)
	iNdEx := 0
	for iNdEx < l {
		var wire uint64
		for shift := uint(0); ; shift += 7 {
			if shift >= 64 {
				return 0, ErrIntOverflowFileformat
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
					return 0, ErrIntOverflowFileformat
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
					return 0, ErrIntOverflowFileformat
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
				return 0, ErrInvalidLengthFileformat
			}
			return iNdEx, nil
		case 3:
			for {
				var innerWire uint64
				var start int = iNdEx
				for shift := uint(0); ; shift += 7 {
					if shift >= 64 {
						return 0, ErrIntOverflowFileformat
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
				next, err := skipFileformat(dAtA[start:])
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
	ErrInvalidLengthFileformat = fmt.Errorf("proto: negative length found during unmarshaling")
	ErrIntOverflowFileformat   = fmt.Errorf("proto: integer overflow")
)

func init() { proto.RegisterFile("fileformat.proto", fileDescriptorFileformat) }

var fileDescriptorFileformat = []byte{
	// 418 bytes of a gzipped FileDescriptorProto
	0x1f, 0x8b, 0x08, 0x00, 0x00, 0x00, 0x00, 0x00, 0x02, 0xff, 0x64, 0x52, 0xc1, 0x6e, 0xd3, 0x40,
	0x14, 0xcc, 0xda, 0x1b, 0xdb, 0x79, 0xa9, 0xaa, 0xed, 0x0a, 0xa1, 0x05, 0xa1, 0xc8, 0xda, 0x53,
	0x0e, 0xe0, 0x43, 0x38, 0xc1, 0x31, 0x22, 0x44, 0x51, 0x83, 0x5d, 0x6d, 0x5d, 0x2a, 0x8e, 0x5b,
	0xb1, 0x31, 0x16, 0x31, 0x1b, 0xb9, 0x5b, 0xa4, 0x70, 0xe1, 0x17, 0xf8, 0x2c, 0x6e, 0xf0, 0x09,
	0x28, 0xfc, 0x08, 0xda, 0xe7, 0x26, 0xa2, 0xca, 0x6d, 0x66, 0xde, 0x1b, 0xef, 0x9b, 0x91, 0x81,
	0xad, 0xea, 0xb5, 0x59, 0xd9, 0xb6, 0xd1, 0x2e, 0xdb, 0xb4, 0xd6, 0x59, 0x39, 0x03, 0x3a, 0xb5,
	0x1f, 0xb7, 0xfc, 0x09, 0xd0, 0xc6, 0x38, 0x2d, 0x48, 0x4a, 0xc6, 0xc3, 0x49, 0x3f, 0x7b, 0x67,
	0x9c, 0x56, 0x28, 0x71, 0x09, 0xf1, 0xca, 0x68, 0x77, 0xd7, 0x1a, 0x11, 0xa4, 0xe1, 0x78, 0x38,
	0x49, 0xb2, 0xb7, 0x1d, 0x57, 0xfb, 0x81, 0x4c, 0x81, 0x7a, 0x07, 0x17, 0x40, 0x9d, 0xae, 0x6e,
	0x05, 0xc1, 0x45, 0x9a, 0x95, 0xba, 0x52, 0xa8, 0xc8, 0x5f, 0x01, 0xc4, 0xf7, 0x36, 0xfe, 0x02,
	0x92, 0xca, 0xd8, 0xc6, 0x6d, 0x37, 0x06, 0x1f, 0x3c, 0x9d, 0x9c, 0xed, 0x3f, 0x99, 0xcd, 0x8d,
	0x6d, 0xca, 0xed, 0xc6, 0xa8, 0xc3, 0x0a, 0x7f, 0x0d, 0xe0, 0xf1, 0xad, 0x69, 0x6b, 0xbd, 0x16,
	0x01, 0x1a, 0x9e, 0x3e, 0x30, 0x5c, 0xe2, 0xa8, 0xfe, 0xa6, 0x5d, 0x6d, 0xbf, 0xa8, 0xff, 0xb6,
	0x39, 0x07, 0xea, 0x99, 0x08, 0x53, 0x32, 0x3e, 0x51, 0x88, 0xbd, 0xb6, 0x36, 0x2b, 0x27, 0x68,
	0x4a, 0xc6, 0x44, 0x21, 0xe6, 0x8f, 0xa0, 0xdf, 0xd6, 0xd5, 0x27, 0x27, 0xfa, 0x28, 0x76, 0x84,
	0x33, 0x08, 0x9d, 0xdd, 0x88, 0x08, 0x35, 0x0f, 0xf9, 0x63, 0x88, 0x6e, 0xac, 0x73, 0xb6, 0x11,
	0x31, 0x8a, 0xf7, 0xec, 0x10, 0x3c, 0x39, 0x0a, 0xfe, 0x0a, 0x92, 0x7d, 0x26, 0x3e, 0x84, 0xf8,
	0x2a, 0x3f, 0xcf, 0x8b, 0xeb, 0x9c, 0xf5, 0xf8, 0x00, 0xfa, 0x17, 0xc5, 0x22, 0x2f, 0x19, 0xe1,
	0x09, 0xd0, 0xe5, 0x22, 0x9f, 0xb1, 0xc0, 0x6f, 0x5c, 0x14, 0xcb, 0x0f, 0xf3, 0x22, 0x67, 0xa1,
	0x7c, 0x06, 0x67, 0x47, 0xe9, 0x78, 0x0c, 0xe1, 0xf5, 0xf9, 0x94, 0xf5, 0xe4, 0x77, 0x08, 0x4b,
	0x5d, 0xf9, 0x1b, 0x3f, 0x9b, 0x2d, 0xf6, 0x38, 0x50, 0x1e, 0xfa, 0x2c, 0x5f, 0xf5, 0xfa, 0xce,
	0x60, 0x55, 0x27, 0xaa, 0x23, 0x5c, 0x02, 0xc5, 0xc2, 0x43, 0xec, 0xef, 0xd4, 0x5f, 0x98, 0xbd,
	0xf7, 0x13, 0x6c, 0x1b, 0x67, 0xf2, 0x39, 0x0c, 0x0e, 0x12, 0x07, 0x88, 0x2e, 0x4b, 0xb5, 0xc8,
	0xe7, 0xac, 0xe7, 0x1f, 0xed, 0x2e, 0x05, 0x88, 0xde, 0x14, 0x57, 0xd3, 0xe5, 0x8c, 0x05, 0x53,
	0xf6, 0x73, 0x37, 0x22, 0xbf, 0x77, 0x23, 0xf2, 0x67, 0x37, 0x22, 0x3f, 0xfe, 0x8e, 0x7a, 0x37,
	0x11, 0xfe, 0x54, 0x2f, 0xff, 0x05, 0x00, 0x00, 0xff, 0xff, 0xc6, 0xcf, 0xe0, 0xac, 0x68, 0x02,
	0x00, 0x00,
}
