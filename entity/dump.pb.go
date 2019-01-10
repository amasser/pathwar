// Code generated by protoc-gen-gogo. DO NOT EDIT.
// source: entity/dump.proto

package entity // import "pathwar.pw/entity"

import proto "github.com/gogo/protobuf/proto"
import fmt "fmt"
import math "math"
import _ "github.com/gogo/protobuf/gogoproto"

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

type Dump struct {
	UserSessions []*UserSession `protobuf:"bytes,1,rep,name=user_sessions,json=userSessions" json:"user_sessions,omitempty"`
	Levels       []*Level       `protobuf:"bytes,2,rep,name=levels" json:"levels,omitempty"`
}

func (m *Dump) Reset()         { *m = Dump{} }
func (m *Dump) String() string { return proto.CompactTextString(m) }
func (*Dump) ProtoMessage()    {}
func (*Dump) Descriptor() ([]byte, []int) {
	return fileDescriptor_dump_44160b0d34e37f65, []int{0}
}
func (m *Dump) XXX_Unmarshal(b []byte) error {
	return m.Unmarshal(b)
}
func (m *Dump) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	if deterministic {
		return xxx_messageInfo_Dump.Marshal(b, m, deterministic)
	} else {
		b = b[:cap(b)]
		n, err := m.MarshalTo(b)
		if err != nil {
			return nil, err
		}
		return b[:n], nil
	}
}
func (dst *Dump) XXX_Merge(src proto.Message) {
	xxx_messageInfo_Dump.Merge(dst, src)
}
func (m *Dump) XXX_Size() int {
	return m.Size()
}
func (m *Dump) XXX_DiscardUnknown() {
	xxx_messageInfo_Dump.DiscardUnknown(m)
}

var xxx_messageInfo_Dump proto.InternalMessageInfo

func (m *Dump) GetUserSessions() []*UserSession {
	if m != nil {
		return m.UserSessions
	}
	return nil
}

func (m *Dump) GetLevels() []*Level {
	if m != nil {
		return m.Levels
	}
	return nil
}

func init() {
	proto.RegisterType((*Dump)(nil), "pathwar.entity.Dump")
}
func (m *Dump) Marshal() (dAtA []byte, err error) {
	size := m.Size()
	dAtA = make([]byte, size)
	n, err := m.MarshalTo(dAtA)
	if err != nil {
		return nil, err
	}
	return dAtA[:n], nil
}

func (m *Dump) MarshalTo(dAtA []byte) (int, error) {
	var i int
	_ = i
	var l int
	_ = l
	if len(m.UserSessions) > 0 {
		for _, msg := range m.UserSessions {
			dAtA[i] = 0xa
			i++
			i = encodeVarintDump(dAtA, i, uint64(msg.Size()))
			n, err := msg.MarshalTo(dAtA[i:])
			if err != nil {
				return 0, err
			}
			i += n
		}
	}
	if len(m.Levels) > 0 {
		for _, msg := range m.Levels {
			dAtA[i] = 0x12
			i++
			i = encodeVarintDump(dAtA, i, uint64(msg.Size()))
			n, err := msg.MarshalTo(dAtA[i:])
			if err != nil {
				return 0, err
			}
			i += n
		}
	}
	return i, nil
}

func encodeVarintDump(dAtA []byte, offset int, v uint64) int {
	for v >= 1<<7 {
		dAtA[offset] = uint8(v&0x7f | 0x80)
		v >>= 7
		offset++
	}
	dAtA[offset] = uint8(v)
	return offset + 1
}
func (m *Dump) Size() (n int) {
	if m == nil {
		return 0
	}
	var l int
	_ = l
	if len(m.UserSessions) > 0 {
		for _, e := range m.UserSessions {
			l = e.Size()
			n += 1 + l + sovDump(uint64(l))
		}
	}
	if len(m.Levels) > 0 {
		for _, e := range m.Levels {
			l = e.Size()
			n += 1 + l + sovDump(uint64(l))
		}
	}
	return n
}

func sovDump(x uint64) (n int) {
	for {
		n++
		x >>= 7
		if x == 0 {
			break
		}
	}
	return n
}
func sozDump(x uint64) (n int) {
	return sovDump(uint64((x << 1) ^ uint64((int64(x) >> 63))))
}
func (m *Dump) Unmarshal(dAtA []byte) error {
	l := len(dAtA)
	iNdEx := 0
	for iNdEx < l {
		preIndex := iNdEx
		var wire uint64
		for shift := uint(0); ; shift += 7 {
			if shift >= 64 {
				return ErrIntOverflowDump
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
			return fmt.Errorf("proto: Dump: wiretype end group for non-group")
		}
		if fieldNum <= 0 {
			return fmt.Errorf("proto: Dump: illegal tag %d (wire type %d)", fieldNum, wire)
		}
		switch fieldNum {
		case 1:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field UserSessions", wireType)
			}
			var msglen int
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowDump
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
				return ErrInvalidLengthDump
			}
			postIndex := iNdEx + msglen
			if postIndex > l {
				return io.ErrUnexpectedEOF
			}
			m.UserSessions = append(m.UserSessions, &UserSession{})
			if err := m.UserSessions[len(m.UserSessions)-1].Unmarshal(dAtA[iNdEx:postIndex]); err != nil {
				return err
			}
			iNdEx = postIndex
		case 2:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field Levels", wireType)
			}
			var msglen int
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowDump
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
				return ErrInvalidLengthDump
			}
			postIndex := iNdEx + msglen
			if postIndex > l {
				return io.ErrUnexpectedEOF
			}
			m.Levels = append(m.Levels, &Level{})
			if err := m.Levels[len(m.Levels)-1].Unmarshal(dAtA[iNdEx:postIndex]); err != nil {
				return err
			}
			iNdEx = postIndex
		default:
			iNdEx = preIndex
			skippy, err := skipDump(dAtA[iNdEx:])
			if err != nil {
				return err
			}
			if skippy < 0 {
				return ErrInvalidLengthDump
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
func skipDump(dAtA []byte) (n int, err error) {
	l := len(dAtA)
	iNdEx := 0
	for iNdEx < l {
		var wire uint64
		for shift := uint(0); ; shift += 7 {
			if shift >= 64 {
				return 0, ErrIntOverflowDump
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
					return 0, ErrIntOverflowDump
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
					return 0, ErrIntOverflowDump
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
				return 0, ErrInvalidLengthDump
			}
			return iNdEx, nil
		case 3:
			for {
				var innerWire uint64
				var start int = iNdEx
				for shift := uint(0); ; shift += 7 {
					if shift >= 64 {
						return 0, ErrIntOverflowDump
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
				next, err := skipDump(dAtA[start:])
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
	ErrInvalidLengthDump = fmt.Errorf("proto: negative length found during unmarshaling")
	ErrIntOverflowDump   = fmt.Errorf("proto: integer overflow")
)

func init() { proto.RegisterFile("entity/dump.proto", fileDescriptor_dump_44160b0d34e37f65) }

var fileDescriptor_dump_44160b0d34e37f65 = []byte{
	// 215 bytes of a gzipped FileDescriptorProto
	0x1f, 0x8b, 0x08, 0x00, 0x00, 0x00, 0x00, 0x00, 0x02, 0xff, 0xe2, 0x12, 0x4c, 0xcd, 0x2b, 0xc9,
	0x2c, 0xa9, 0xd4, 0x4f, 0x29, 0xcd, 0x2d, 0xd0, 0x2b, 0x28, 0xca, 0x2f, 0xc9, 0x17, 0xe2, 0x2b,
	0x48, 0x2c, 0xc9, 0x28, 0x4f, 0x2c, 0xd2, 0x83, 0x48, 0x49, 0xe9, 0xa6, 0x67, 0x96, 0x64, 0x94,
	0x26, 0xe9, 0x25, 0xe7, 0xe7, 0xea, 0xa7, 0xe7, 0xa7, 0xe7, 0xeb, 0x83, 0x95, 0x25, 0x95, 0xa6,
	0x81, 0x79, 0x60, 0x0e, 0x98, 0x05, 0xd1, 0x2e, 0x25, 0x09, 0x35, 0xb1, 0xb4, 0x38, 0xb5, 0x28,
	0xbe, 0x38, 0xb5, 0xb8, 0x38, 0x33, 0x3f, 0x0f, 0x2a, 0x25, 0x04, 0x95, 0xca, 0x49, 0x2d, 0x4b,
	0xcd, 0x81, 0x88, 0x29, 0x95, 0x73, 0xb1, 0xb8, 0x94, 0xe6, 0x16, 0x08, 0x39, 0x70, 0xf1, 0x22,
	0xeb, 0x28, 0x96, 0x60, 0x54, 0x60, 0xd6, 0xe0, 0x36, 0x92, 0xd6, 0x43, 0x75, 0x8d, 0x5e, 0x68,
	0x71, 0x6a, 0x51, 0x30, 0x44, 0x4d, 0x10, 0x4f, 0x29, 0x82, 0x53, 0x2c, 0xa4, 0xcb, 0xc5, 0x06,
	0x36, 0xb8, 0x58, 0x82, 0x09, 0xac, 0x55, 0x14, 0x5d, 0xab, 0x0f, 0x48, 0x36, 0x08, 0xaa, 0xc8,
	0x49, 0xfb, 0xc4, 0x23, 0x39, 0xc6, 0x0b, 0x8f, 0xe4, 0x18, 0x1f, 0x3c, 0x92, 0x63, 0x9c, 0xf0,
	0x58, 0x8e, 0xe1, 0xc2, 0x63, 0x39, 0x86, 0x1b, 0x8f, 0xe5, 0x18, 0xa2, 0x04, 0x61, 0xfa, 0x0a,
	0xca, 0xf5, 0x21, 0x5a, 0x93, 0xd8, 0xc0, 0x8e, 0x35, 0x06, 0x04, 0x00, 0x00, 0xff, 0xff, 0x88,
	0x69, 0x00, 0x05, 0x2f, 0x01, 0x00, 0x00,
}
