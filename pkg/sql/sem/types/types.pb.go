// Code generated by protoc-gen-gogo. DO NOT EDIT.
// source: sql/sem/types/types.proto

package types

import proto "github.com/gogo/protobuf/proto"
import fmt "fmt"
import math "math"

import github_com_lib_pq_oid "github.com/lib/pq/oid"

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

type SemanticType int32

const (
	BOOL SemanticType = 0
	// INT(width)
	INT SemanticType = 1
	// FLOAT(precision)
	FLOAT SemanticType = 2
	// DECIMAL(precision, width /* scale */)
	DECIMAL   SemanticType = 3
	DATE      SemanticType = 4
	TIMESTAMP SemanticType = 5
	INTERVAL  SemanticType = 6
	// STRING(width)
	STRING      SemanticType = 7
	BYTES       SemanticType = 8
	TIMESTAMPTZ SemanticType = 9
	// Collated STRING, CHAR, and VARCHAR
	// Collated string key columns are encoded partly as a key and partly as a
	// value. The key part is the collation key, so that different strings that
	// collate equal cannot both be used as keys. The value part is the usual
	// UTF-8 encoding of the string. This creates several special cases in the
	// encode/decode logic.
	COLLATEDSTRING SemanticType = 10
	OID            SemanticType = 12
	// NULL is not supported as a table column type, however it can be
	// transferred through distsql streams.
	NULL  SemanticType = 13
	UUID  SemanticType = 14
	ARRAY SemanticType = 15
	INET  SemanticType = 16
	TIME  SemanticType = 17
	JSON  SemanticType = 18
	TUPLE SemanticType = 20
	BIT   SemanticType = 21
	// Special type used during static analysis. This should never be present
	// at execution time.
	ANY SemanticType = 100
)

var SemanticType_name = map[int32]string{
	0:   "BOOL",
	1:   "INT",
	2:   "FLOAT",
	3:   "DECIMAL",
	4:   "DATE",
	5:   "TIMESTAMP",
	6:   "INTERVAL",
	7:   "STRING",
	8:   "BYTES",
	9:   "TIMESTAMPTZ",
	10:  "COLLATEDSTRING",
	12:  "OID",
	13:  "NULL",
	14:  "UUID",
	15:  "ARRAY",
	16:  "INET",
	17:  "TIME",
	18:  "JSON",
	20:  "TUPLE",
	21:  "BIT",
	100: "ANY",
}
var SemanticType_value = map[string]int32{
	"BOOL":           0,
	"INT":            1,
	"FLOAT":          2,
	"DECIMAL":        3,
	"DATE":           4,
	"TIMESTAMP":      5,
	"INTERVAL":       6,
	"STRING":         7,
	"BYTES":          8,
	"TIMESTAMPTZ":    9,
	"COLLATEDSTRING": 10,
	"OID":            12,
	"NULL":           13,
	"UUID":           14,
	"ARRAY":          15,
	"INET":           16,
	"TIME":           17,
	"JSON":           18,
	"TUPLE":          20,
	"BIT":            21,
	"ANY":            100,
}

func (x SemanticType) Enum() *SemanticType {
	p := new(SemanticType)
	*p = x
	return p
}
func (x SemanticType) String() string {
	return proto.EnumName(SemanticType_name, int32(x))
}
func (x *SemanticType) UnmarshalJSON(data []byte) error {
	value, err := proto.UnmarshalJSONEnum(SemanticType_value, data, "SemanticType")
	if err != nil {
		return err
	}
	*x = SemanticType(value)
	return nil
}
func (SemanticType) EnumDescriptor() ([]byte, []int) {
	return fileDescriptor_types_f87d5ef8a4bdbdae, []int{0}
}

type VisibleType int32

const (
	VisibleType_NONE VisibleType = 0
)

var VisibleType_name = map[int32]string{
	0: "NONE",
}
var VisibleType_value = map[string]int32{
	"NONE": 0,
}

func (x VisibleType) Enum() *VisibleType {
	p := new(VisibleType)
	*p = x
	return p
}
func (x VisibleType) String() string {
	return proto.EnumName(VisibleType_name, int32(x))
}
func (x *VisibleType) UnmarshalJSON(data []byte) error {
	value, err := proto.UnmarshalJSONEnum(VisibleType_value, data, "VisibleType")
	if err != nil {
		return err
	}
	*x = VisibleType(value)
	return nil
}
func (VisibleType) EnumDescriptor() ([]byte, []int) {
	return fileDescriptor_types_f87d5ef8a4bdbdae, []int{1}
}

type InternalColumnType struct {
	SemanticType SemanticType `protobuf:"varint,1,opt,name=semantic_type,json=semanticType,enum=cockroach.sql.sem.types.SemanticType" json:"semantic_type"`
	// INT, DECIMAL, CHAR and BINARY
	Width int32 `protobuf:"varint,2,opt,name=width" json:"width"`
	// DECIMAL
	// Also FLOAT pre-2.1 (this was incorrect.)
	Precision int32 `protobuf:"varint,3,opt,name=precision" json:"precision"`
	// The length of each dimension in the array. A dimension of -1 means that
	// no bound was specified for that dimension.
	ArrayDimensions []int32 `protobuf:"varint,4,rep,name=array_dimensions,json=arrayDimensions" json:"array_dimensions,omitempty"`
	// Collated STRING, CHAR, and VARCHAR
	Locale *string `protobuf:"bytes,5,opt,name=locale" json:"locale,omitempty"`
	// Alias for any types where our internal representation is different than
	// the user specification. Examples are INT4, FLOAT4, etc. Mostly for Postgres
	// compatibility.
	ZZZ_VisibleType VisibleType `protobuf:"varint,6,opt,name=visible_type,json=visibleType,enum=cockroach.sql.sem.types.VisibleType" json:"visible_type"`
	// Only used if the kind is ARRAY.
	ArrayContents *SemanticType `protobuf:"varint,7,opt,name=array_contents,json=arrayContents,enum=cockroach.sql.sem.types.SemanticType" json:"array_contents,omitempty"`
	// Only used if the kind is TUPLE
	TupleContents []ColumnType `protobuf:"bytes,8,rep,name=tuple_contents,json=tupleContents,customtype=ColumnType" json:"tuple_contents"`
	TupleLabels   []string     `protobuf:"bytes,9,rep,name=tuple_labels,json=tupleLabels" json:"tuple_labels,omitempty"`
	// XXX_Oid is an internal field that should not be used directly by callers.
	// The Oid() method maps this and other fields to the correct Oid value for
	// the type.
	// TODO(andyk): Temporary name that will be renamed in future commit.
	ZZZ_Oid              github_com_lib_pq_oid.Oid `protobuf:"varint,10,opt,name=oid,customtype=github.com/lib/pq/oid.Oid" json:"oid"`
	XXX_NoUnkeyedLiteral struct{}                  `json:"-"`
	XXX_sizecache        int32                     `json:"-"`
}

func (m *InternalColumnType) Reset()         { *m = InternalColumnType{} }
func (m *InternalColumnType) String() string { return proto.CompactTextString(m) }
func (*InternalColumnType) ProtoMessage()    {}
func (*InternalColumnType) Descriptor() ([]byte, []int) {
	return fileDescriptor_types_f87d5ef8a4bdbdae, []int{0}
}
func (m *InternalColumnType) XXX_Unmarshal(b []byte) error {
	return m.Unmarshal(b)
}
func (m *InternalColumnType) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	b = b[:cap(b)]
	n, err := m.MarshalTo(b)
	if err != nil {
		return nil, err
	}
	return b[:n], nil
}
func (dst *InternalColumnType) XXX_Merge(src proto.Message) {
	xxx_messageInfo_InternalColumnType.Merge(dst, src)
}
func (m *InternalColumnType) XXX_Size() int {
	return m.Size()
}
func (m *InternalColumnType) XXX_DiscardUnknown() {
	xxx_messageInfo_InternalColumnType.DiscardUnknown(m)
}

var xxx_messageInfo_InternalColumnType proto.InternalMessageInfo

func init() {
	proto.RegisterType((*InternalColumnType)(nil), "cockroach.sql.sem.types.InternalColumnType")
	proto.RegisterEnum("cockroach.sql.sem.types.SemanticType", SemanticType_name, SemanticType_value)
	proto.RegisterEnum("cockroach.sql.sem.types.VisibleType", VisibleType_name, VisibleType_value)
}
func (m *InternalColumnType) Marshal() (dAtA []byte, err error) {
	size := m.Size()
	dAtA = make([]byte, size)
	n, err := m.MarshalTo(dAtA)
	if err != nil {
		return nil, err
	}
	return dAtA[:n], nil
}

func (m *InternalColumnType) MarshalTo(dAtA []byte) (int, error) {
	var i int
	_ = i
	var l int
	_ = l
	dAtA[i] = 0x8
	i++
	i = encodeVarintTypes(dAtA, i, uint64(m.SemanticType))
	dAtA[i] = 0x10
	i++
	i = encodeVarintTypes(dAtA, i, uint64(m.Width))
	dAtA[i] = 0x18
	i++
	i = encodeVarintTypes(dAtA, i, uint64(m.Precision))
	if len(m.ArrayDimensions) > 0 {
		for _, num := range m.ArrayDimensions {
			dAtA[i] = 0x20
			i++
			i = encodeVarintTypes(dAtA, i, uint64(num))
		}
	}
	if m.Locale != nil {
		dAtA[i] = 0x2a
		i++
		i = encodeVarintTypes(dAtA, i, uint64(len(*m.Locale)))
		i += copy(dAtA[i:], *m.Locale)
	}
	dAtA[i] = 0x30
	i++
	i = encodeVarintTypes(dAtA, i, uint64(m.ZZZ_VisibleType))
	if m.ArrayContents != nil {
		dAtA[i] = 0x38
		i++
		i = encodeVarintTypes(dAtA, i, uint64(*m.ArrayContents))
	}
	if len(m.TupleContents) > 0 {
		for _, msg := range m.TupleContents {
			dAtA[i] = 0x42
			i++
			i = encodeVarintTypes(dAtA, i, uint64(msg.Size()))
			n, err := msg.MarshalTo(dAtA[i:])
			if err != nil {
				return 0, err
			}
			i += n
		}
	}
	if len(m.TupleLabels) > 0 {
		for _, s := range m.TupleLabels {
			dAtA[i] = 0x4a
			i++
			l = len(s)
			for l >= 1<<7 {
				dAtA[i] = uint8(uint64(l)&0x7f | 0x80)
				l >>= 7
				i++
			}
			dAtA[i] = uint8(l)
			i++
			i += copy(dAtA[i:], s)
		}
	}
	dAtA[i] = 0x50
	i++
	i = encodeVarintTypes(dAtA, i, uint64(m.ZZZ_Oid))
	return i, nil
}

func encodeVarintTypes(dAtA []byte, offset int, v uint64) int {
	for v >= 1<<7 {
		dAtA[offset] = uint8(v&0x7f | 0x80)
		v >>= 7
		offset++
	}
	dAtA[offset] = uint8(v)
	return offset + 1
}
func (m *InternalColumnType) Size() (n int) {
	if m == nil {
		return 0
	}
	var l int
	_ = l
	n += 1 + sovTypes(uint64(m.SemanticType))
	n += 1 + sovTypes(uint64(m.Width))
	n += 1 + sovTypes(uint64(m.Precision))
	if len(m.ArrayDimensions) > 0 {
		for _, e := range m.ArrayDimensions {
			n += 1 + sovTypes(uint64(e))
		}
	}
	if m.Locale != nil {
		l = len(*m.Locale)
		n += 1 + l + sovTypes(uint64(l))
	}
	n += 1 + sovTypes(uint64(m.ZZZ_VisibleType))
	if m.ArrayContents != nil {
		n += 1 + sovTypes(uint64(*m.ArrayContents))
	}
	if len(m.TupleContents) > 0 {
		for _, e := range m.TupleContents {
			l = e.Size()
			n += 1 + l + sovTypes(uint64(l))
		}
	}
	if len(m.TupleLabels) > 0 {
		for _, s := range m.TupleLabels {
			l = len(s)
			n += 1 + l + sovTypes(uint64(l))
		}
	}
	n += 1 + sovTypes(uint64(m.ZZZ_Oid))
	return n
}

func sovTypes(x uint64) (n int) {
	for {
		n++
		x >>= 7
		if x == 0 {
			break
		}
	}
	return n
}
func sozTypes(x uint64) (n int) {
	return sovTypes(uint64((x << 1) ^ uint64((int64(x) >> 63))))
}
func (m *InternalColumnType) Unmarshal(dAtA []byte) error {
	l := len(dAtA)
	iNdEx := 0
	for iNdEx < l {
		preIndex := iNdEx
		var wire uint64
		for shift := uint(0); ; shift += 7 {
			if shift >= 64 {
				return ErrIntOverflowTypes
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
			return fmt.Errorf("proto: InternalColumnType: wiretype end group for non-group")
		}
		if fieldNum <= 0 {
			return fmt.Errorf("proto: InternalColumnType: illegal tag %d (wire type %d)", fieldNum, wire)
		}
		switch fieldNum {
		case 1:
			if wireType != 0 {
				return fmt.Errorf("proto: wrong wireType = %d for field SemanticType", wireType)
			}
			m.SemanticType = 0
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowTypes
				}
				if iNdEx >= l {
					return io.ErrUnexpectedEOF
				}
				b := dAtA[iNdEx]
				iNdEx++
				m.SemanticType |= (SemanticType(b) & 0x7F) << shift
				if b < 0x80 {
					break
				}
			}
		case 2:
			if wireType != 0 {
				return fmt.Errorf("proto: wrong wireType = %d for field Width", wireType)
			}
			m.Width = 0
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowTypes
				}
				if iNdEx >= l {
					return io.ErrUnexpectedEOF
				}
				b := dAtA[iNdEx]
				iNdEx++
				m.Width |= (int32(b) & 0x7F) << shift
				if b < 0x80 {
					break
				}
			}
		case 3:
			if wireType != 0 {
				return fmt.Errorf("proto: wrong wireType = %d for field Precision", wireType)
			}
			m.Precision = 0
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowTypes
				}
				if iNdEx >= l {
					return io.ErrUnexpectedEOF
				}
				b := dAtA[iNdEx]
				iNdEx++
				m.Precision |= (int32(b) & 0x7F) << shift
				if b < 0x80 {
					break
				}
			}
		case 4:
			if wireType == 0 {
				var v int32
				for shift := uint(0); ; shift += 7 {
					if shift >= 64 {
						return ErrIntOverflowTypes
					}
					if iNdEx >= l {
						return io.ErrUnexpectedEOF
					}
					b := dAtA[iNdEx]
					iNdEx++
					v |= (int32(b) & 0x7F) << shift
					if b < 0x80 {
						break
					}
				}
				m.ArrayDimensions = append(m.ArrayDimensions, v)
			} else if wireType == 2 {
				var packedLen int
				for shift := uint(0); ; shift += 7 {
					if shift >= 64 {
						return ErrIntOverflowTypes
					}
					if iNdEx >= l {
						return io.ErrUnexpectedEOF
					}
					b := dAtA[iNdEx]
					iNdEx++
					packedLen |= (int(b) & 0x7F) << shift
					if b < 0x80 {
						break
					}
				}
				if packedLen < 0 {
					return ErrInvalidLengthTypes
				}
				postIndex := iNdEx + packedLen
				if postIndex > l {
					return io.ErrUnexpectedEOF
				}
				var elementCount int
				var count int
				for _, integer := range dAtA {
					if integer < 128 {
						count++
					}
				}
				elementCount = count
				if elementCount != 0 && len(m.ArrayDimensions) == 0 {
					m.ArrayDimensions = make([]int32, 0, elementCount)
				}
				for iNdEx < postIndex {
					var v int32
					for shift := uint(0); ; shift += 7 {
						if shift >= 64 {
							return ErrIntOverflowTypes
						}
						if iNdEx >= l {
							return io.ErrUnexpectedEOF
						}
						b := dAtA[iNdEx]
						iNdEx++
						v |= (int32(b) & 0x7F) << shift
						if b < 0x80 {
							break
						}
					}
					m.ArrayDimensions = append(m.ArrayDimensions, v)
				}
			} else {
				return fmt.Errorf("proto: wrong wireType = %d for field ArrayDimensions", wireType)
			}
		case 5:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field Locale", wireType)
			}
			var stringLen uint64
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowTypes
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
				return ErrInvalidLengthTypes
			}
			postIndex := iNdEx + intStringLen
			if postIndex > l {
				return io.ErrUnexpectedEOF
			}
			s := string(dAtA[iNdEx:postIndex])
			m.Locale = &s
			iNdEx = postIndex
		case 6:
			if wireType != 0 {
				return fmt.Errorf("proto: wrong wireType = %d for field ZZZ_VisibleType", wireType)
			}
			m.ZZZ_VisibleType = 0
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowTypes
				}
				if iNdEx >= l {
					return io.ErrUnexpectedEOF
				}
				b := dAtA[iNdEx]
				iNdEx++
				m.ZZZ_VisibleType |= (VisibleType(b) & 0x7F) << shift
				if b < 0x80 {
					break
				}
			}
		case 7:
			if wireType != 0 {
				return fmt.Errorf("proto: wrong wireType = %d for field ArrayContents", wireType)
			}
			var v SemanticType
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowTypes
				}
				if iNdEx >= l {
					return io.ErrUnexpectedEOF
				}
				b := dAtA[iNdEx]
				iNdEx++
				v |= (SemanticType(b) & 0x7F) << shift
				if b < 0x80 {
					break
				}
			}
			m.ArrayContents = &v
		case 8:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field TupleContents", wireType)
			}
			var byteLen int
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowTypes
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
				return ErrInvalidLengthTypes
			}
			postIndex := iNdEx + byteLen
			if postIndex > l {
				return io.ErrUnexpectedEOF
			}
			var v ColumnType
			m.TupleContents = append(m.TupleContents, v)
			if err := m.TupleContents[len(m.TupleContents)-1].Unmarshal(dAtA[iNdEx:postIndex]); err != nil {
				return err
			}
			iNdEx = postIndex
		case 9:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field TupleLabels", wireType)
			}
			var stringLen uint64
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowTypes
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
				return ErrInvalidLengthTypes
			}
			postIndex := iNdEx + intStringLen
			if postIndex > l {
				return io.ErrUnexpectedEOF
			}
			m.TupleLabels = append(m.TupleLabels, string(dAtA[iNdEx:postIndex]))
			iNdEx = postIndex
		case 10:
			if wireType != 0 {
				return fmt.Errorf("proto: wrong wireType = %d for field ZZZ_Oid", wireType)
			}
			m.ZZZ_Oid = 0
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowTypes
				}
				if iNdEx >= l {
					return io.ErrUnexpectedEOF
				}
				b := dAtA[iNdEx]
				iNdEx++
				m.ZZZ_Oid |= (github_com_lib_pq_oid.Oid(b) & 0x7F) << shift
				if b < 0x80 {
					break
				}
			}
		default:
			iNdEx = preIndex
			skippy, err := skipTypes(dAtA[iNdEx:])
			if err != nil {
				return err
			}
			if skippy < 0 {
				return ErrInvalidLengthTypes
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
func skipTypes(dAtA []byte) (n int, err error) {
	l := len(dAtA)
	iNdEx := 0
	for iNdEx < l {
		var wire uint64
		for shift := uint(0); ; shift += 7 {
			if shift >= 64 {
				return 0, ErrIntOverflowTypes
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
					return 0, ErrIntOverflowTypes
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
					return 0, ErrIntOverflowTypes
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
				return 0, ErrInvalidLengthTypes
			}
			return iNdEx, nil
		case 3:
			for {
				var innerWire uint64
				var start int = iNdEx
				for shift := uint(0); ; shift += 7 {
					if shift >= 64 {
						return 0, ErrIntOverflowTypes
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
				next, err := skipTypes(dAtA[start:])
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
	ErrInvalidLengthTypes = fmt.Errorf("proto: negative length found during unmarshaling")
	ErrIntOverflowTypes   = fmt.Errorf("proto: integer overflow")
)

func init() { proto.RegisterFile("sql/sem/types/types.proto", fileDescriptor_types_f87d5ef8a4bdbdae) }

var fileDescriptor_types_f87d5ef8a4bdbdae = []byte{
	// 704 bytes of a gzipped FileDescriptorProto
	0x1f, 0x8b, 0x08, 0x00, 0x00, 0x00, 0x00, 0x00, 0x02, 0xff, 0x94, 0x52, 0xcd, 0x8e, 0xe3, 0x44,
	0x18, 0x4c, 0xc7, 0x6d, 0xa7, 0xdd, 0xf9, 0xfb, 0xe8, 0x5d, 0x58, 0xef, 0x1c, 0x9c, 0xd0, 0x02,
	0x29, 0xac, 0x90, 0x23, 0x71, 0xe3, 0x82, 0xe4, 0x24, 0x06, 0x19, 0x79, 0xec, 0x91, 0xe3, 0xac,
	0x34, 0xb9, 0x0c, 0x89, 0x6d, 0xcd, 0x58, 0x38, 0x76, 0x26, 0xf6, 0x2c, 0x9a, 0x37, 0xe0, 0xc8,
	0x2b, 0x20, 0x78, 0x07, 0x24, 0x9e, 0x20, 0x7b, 0xe3, 0x88, 0x38, 0x44, 0x10, 0x5e, 0x04, 0xb5,
	0x13, 0x4d, 0x72, 0x99, 0x03, 0x17, 0xeb, 0xeb, 0xfa, 0xaa, 0xca, 0x55, 0x76, 0xd3, 0xd7, 0xc5,
	0x7d, 0x3a, 0x2c, 0xe2, 0xd5, 0xb0, 0x7c, 0x5c, 0xc7, 0xc5, 0xe1, 0x69, 0xac, 0x37, 0x79, 0x99,
	0xb3, 0x57, 0x61, 0x1e, 0x7e, 0xbf, 0xc9, 0x17, 0xe1, 0x9d, 0x51, 0xdc, 0xa7, 0x46, 0x11, 0xaf,
	0x8c, 0x6a, 0x7d, 0xf1, 0xf2, 0x36, 0xbf, 0xcd, 0x2b, 0xce, 0x50, 0x4c, 0x07, 0x3a, 0xff, 0x0d,
	0x53, 0x66, 0x67, 0x65, 0xbc, 0xc9, 0x16, 0xe9, 0x38, 0x4f, 0x1f, 0x56, 0x59, 0xf0, 0xb8, 0x8e,
	0xd9, 0x15, 0x6d, 0x17, 0xf1, 0x6a, 0x91, 0x95, 0x49, 0x78, 0x23, 0xe4, 0x1a, 0xea, 0xa3, 0x41,
	0xe7, 0x8b, 0x4f, 0x8d, 0x67, 0xdc, 0x8d, 0xe9, 0x91, 0x2d, 0xd4, 0x23, 0xbc, 0xdd, 0xf5, 0x6a,
	0x7e, 0xab, 0x38, 0xc3, 0xd8, 0x05, 0x95, 0x7f, 0x48, 0xa2, 0xf2, 0x4e, 0xab, 0xf7, 0xd1, 0x40,
	0x3e, 0x52, 0x0e, 0x10, 0xe3, 0x54, 0x5d, 0x6f, 0xe2, 0x30, 0x29, 0x92, 0x3c, 0xd3, 0xa4, 0xb3,
	0xfd, 0x09, 0x66, 0x9f, 0x51, 0x58, 0x6c, 0x36, 0x8b, 0xc7, 0x9b, 0x28, 0x59, 0xc5, 0x99, 0x80,
	0x0a, 0x0d, 0xf7, 0xa5, 0x81, 0xec, 0x77, 0x2b, 0x7c, 0xf2, 0x04, 0xb3, 0x8f, 0xa8, 0x92, 0xe6,
	0xe1, 0x22, 0x8d, 0x35, 0xb9, 0x8f, 0x06, 0xaa, 0x7f, 0x3c, 0xb1, 0xef, 0x68, 0xeb, 0x5d, 0x52,
	0x24, 0xcb, 0x34, 0x3e, 0x74, 0x52, 0xaa, 0x4e, 0x9f, 0x3c, 0xdb, 0xe9, 0xed, 0x81, 0x5c, 0x55,
	0x7a, 0x25, 0xf2, 0xec, 0x77, 0xbd, 0xee, 0x7c, 0x3e, 0xbf, 0x39, 0x5b, 0xf8, 0xcd, 0x77, 0xa7,
	0x03, 0x73, 0x68, 0xe7, 0x10, 0x32, 0xcc, 0xb3, 0x32, 0xce, 0xca, 0x42, 0x6b, 0xfc, 0x8f, 0xef,
	0xe6, 0xb7, 0x2b, 0xf1, 0xf8, 0xa8, 0x65, 0x5f, 0xd2, 0x4e, 0xf9, 0xb0, 0x4e, 0xe3, 0x93, 0x1b,
	0xe9, 0x4b, 0x83, 0xd6, 0x88, 0x89, 0x2c, 0x7f, 0xed, 0x7a, 0xf4, 0xf4, 0xc3, 0xfc, 0x76, 0xc5,
	0x7c, 0x92, 0x7e, 0x4c, 0x5b, 0x07, 0x69, 0xba, 0x58, 0xc6, 0x69, 0xa1, 0xa9, 0x7d, 0x69, 0xa0,
	0xfa, 0xcd, 0x0a, 0x73, 0x2a, 0x88, 0x7d, 0x45, 0xa5, 0x3c, 0x89, 0x34, 0xda, 0x47, 0x83, 0xf6,
	0xe8, 0xf3, 0xa3, 0xe5, 0xeb, 0xdb, 0xa4, 0xbc, 0x7b, 0x58, 0x1a, 0x61, 0xbe, 0x1a, 0xa6, 0xc9,
	0x72, 0xb8, 0xbe, 0x1f, 0xe6, 0x49, 0x64, 0x78, 0x49, 0xb4, 0xdf, 0xf5, 0x1a, 0xa2, 0xbb, 0x97,
	0x44, 0xbe, 0x10, 0xbe, 0xf9, 0xb9, 0x4e, 0x5b, 0xe7, 0xe9, 0x19, 0xa1, 0x78, 0xe4, 0x79, 0x0e,
	0xd4, 0x58, 0x83, 0x4a, 0xb6, 0x1b, 0x00, 0x62, 0x2a, 0x95, 0xbf, 0x76, 0x3c, 0x33, 0x80, 0x3a,
	0x6b, 0xd2, 0xc6, 0xc4, 0x1a, 0xdb, 0x97, 0xa6, 0x03, 0x92, 0xa0, 0x4e, 0xcc, 0xc0, 0x02, 0xcc,
	0xda, 0x54, 0x0d, 0xec, 0x4b, 0x6b, 0x1a, 0x98, 0x97, 0x57, 0x20, 0xb3, 0x16, 0x25, 0xb6, 0x1b,
	0x58, 0xfe, 0x5b, 0xd3, 0x01, 0x85, 0x51, 0xaa, 0x4c, 0x03, 0xdf, 0x76, 0xbf, 0x81, 0x86, 0xb0,
	0x1a, 0x5d, 0x07, 0xd6, 0x14, 0x08, 0xeb, 0xd2, 0xe6, 0x93, 0x26, 0x98, 0x83, 0xca, 0x18, 0xed,
	0x8c, 0x3d, 0xc7, 0x31, 0x03, 0x6b, 0x72, 0xe4, 0x53, 0x91, 0xc1, 0xb3, 0x27, 0xd0, 0x12, 0xef,
	0x72, 0x67, 0x8e, 0x03, 0x6d, 0x31, 0xcd, 0x66, 0xf6, 0x04, 0x3a, 0xc2, 0xcc, 0xf4, 0x7d, 0xf3,
	0x1a, 0xba, 0x02, 0xb4, 0x5d, 0x2b, 0x00, 0x10, 0x93, 0xb0, 0x85, 0x0f, 0xc4, 0xf4, 0xed, 0xd4,
	0x73, 0x81, 0x09, 0x62, 0x30, 0xbb, 0x72, 0x2c, 0x78, 0x29, 0x0c, 0x47, 0x76, 0x00, 0x1f, 0x8a,
	0xc1, 0x74, 0xaf, 0x21, 0xba, 0xc0, 0x3f, 0xfe, 0xa2, 0xd7, 0x38, 0x26, 0x4d, 0x68, 0x72, 0x4c,
	0x5e, 0xc0, 0x0b, 0xae, 0x90, 0x2d, 0x82, 0x2d, 0xe2, 0x0a, 0x79, 0x8f, 0xe0, 0x3d, 0x7a, 0x13,
	0xd1, 0xe6, 0xd9, 0x5d, 0xa9, 0xa2, 0x78, 0xae, 0x05, 0xb5, 0x8b, 0xee, 0xef, 0xbf, 0xea, 0xe7,
	0x2b, 0x8e, 0x09, 0x02, 0xc4, 0x31, 0xa9, 0x43, 0x9d, 0x63, 0x22, 0x81, 0xc4, 0x31, 0xc1, 0x80,
	0x39, 0x26, 0x32, 0xc8, 0x1c, 0x13, 0x05, 0x14, 0x8e, 0x49, 0x03, 0x1a, 0x1c, 0x13, 0x02, 0x84,
	0x63, 0xa2, 0x82, 0xca, 0x31, 0xa1, 0x40, 0x47, 0xbd, 0xed, 0x3f, 0x7a, 0x6d, 0xbb, 0xd7, 0xd1,
	0x1f, 0x7b, 0x1d, 0xfd, 0xb9, 0xd7, 0xd1, 0xdf, 0x7b, 0x1d, 0xfd, 0xf4, 0xaf, 0x5e, 0x9b, 0xcb,
	0xd5, 0x25, 0xfb, 0x2f, 0x00, 0x00, 0xff, 0xff, 0x83, 0xa9, 0x88, 0xaa, 0x2f, 0x04, 0x00, 0x00,
}