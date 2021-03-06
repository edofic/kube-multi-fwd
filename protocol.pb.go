// Code generated by protoc-gen-go. DO NOT EDIT.
// source: protocol.proto

package fwd // import "."

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

type ProxyRequest struct {
	// Types that are valid to be assigned to Req:
	//	*ProxyRequest_Connect
	//	*ProxyRequest_Chunk
	Req                  isProxyRequest_Req `protobuf_oneof:"req"`
	XXX_NoUnkeyedLiteral struct{}           `json:"-"`
	XXX_unrecognized     []byte             `json:"-"`
	XXX_sizecache        int32              `json:"-"`
}

func (m *ProxyRequest) Reset()         { *m = ProxyRequest{} }
func (m *ProxyRequest) String() string { return proto.CompactTextString(m) }
func (*ProxyRequest) ProtoMessage()    {}
func (*ProxyRequest) Descriptor() ([]byte, []int) {
	return fileDescriptor_protocol_536ecc5858538e00, []int{0}
}
func (m *ProxyRequest) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_ProxyRequest.Unmarshal(m, b)
}
func (m *ProxyRequest) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_ProxyRequest.Marshal(b, m, deterministic)
}
func (dst *ProxyRequest) XXX_Merge(src proto.Message) {
	xxx_messageInfo_ProxyRequest.Merge(dst, src)
}
func (m *ProxyRequest) XXX_Size() int {
	return xxx_messageInfo_ProxyRequest.Size(m)
}
func (m *ProxyRequest) XXX_DiscardUnknown() {
	xxx_messageInfo_ProxyRequest.DiscardUnknown(m)
}

var xxx_messageInfo_ProxyRequest proto.InternalMessageInfo

type isProxyRequest_Req interface {
	isProxyRequest_Req()
}

type ProxyRequest_Connect struct {
	Connect *ProxyConnect `protobuf:"bytes,1,opt,name=connect,proto3,oneof"`
}
type ProxyRequest_Chunk struct {
	Chunk []byte `protobuf:"bytes,2,opt,name=chunk,proto3,oneof"`
}

func (*ProxyRequest_Connect) isProxyRequest_Req() {}
func (*ProxyRequest_Chunk) isProxyRequest_Req()   {}

func (m *ProxyRequest) GetReq() isProxyRequest_Req {
	if m != nil {
		return m.Req
	}
	return nil
}

func (m *ProxyRequest) GetConnect() *ProxyConnect {
	if x, ok := m.GetReq().(*ProxyRequest_Connect); ok {
		return x.Connect
	}
	return nil
}

func (m *ProxyRequest) GetChunk() []byte {
	if x, ok := m.GetReq().(*ProxyRequest_Chunk); ok {
		return x.Chunk
	}
	return nil
}

// XXX_OneofFuncs is for the internal use of the proto package.
func (*ProxyRequest) XXX_OneofFuncs() (func(msg proto.Message, b *proto.Buffer) error, func(msg proto.Message, tag, wire int, b *proto.Buffer) (bool, error), func(msg proto.Message) (n int), []interface{}) {
	return _ProxyRequest_OneofMarshaler, _ProxyRequest_OneofUnmarshaler, _ProxyRequest_OneofSizer, []interface{}{
		(*ProxyRequest_Connect)(nil),
		(*ProxyRequest_Chunk)(nil),
	}
}

func _ProxyRequest_OneofMarshaler(msg proto.Message, b *proto.Buffer) error {
	m := msg.(*ProxyRequest)
	// req
	switch x := m.Req.(type) {
	case *ProxyRequest_Connect:
		b.EncodeVarint(1<<3 | proto.WireBytes)
		if err := b.EncodeMessage(x.Connect); err != nil {
			return err
		}
	case *ProxyRequest_Chunk:
		b.EncodeVarint(2<<3 | proto.WireBytes)
		b.EncodeRawBytes(x.Chunk)
	case nil:
	default:
		return fmt.Errorf("ProxyRequest.Req has unexpected type %T", x)
	}
	return nil
}

func _ProxyRequest_OneofUnmarshaler(msg proto.Message, tag, wire int, b *proto.Buffer) (bool, error) {
	m := msg.(*ProxyRequest)
	switch tag {
	case 1: // req.connect
		if wire != proto.WireBytes {
			return true, proto.ErrInternalBadWireType
		}
		msg := new(ProxyConnect)
		err := b.DecodeMessage(msg)
		m.Req = &ProxyRequest_Connect{msg}
		return true, err
	case 2: // req.chunk
		if wire != proto.WireBytes {
			return true, proto.ErrInternalBadWireType
		}
		x, err := b.DecodeRawBytes(true)
		m.Req = &ProxyRequest_Chunk{x}
		return true, err
	default:
		return false, nil
	}
}

func _ProxyRequest_OneofSizer(msg proto.Message) (n int) {
	m := msg.(*ProxyRequest)
	// req
	switch x := m.Req.(type) {
	case *ProxyRequest_Connect:
		s := proto.Size(x.Connect)
		n += 1 // tag and wire
		n += proto.SizeVarint(uint64(s))
		n += s
	case *ProxyRequest_Chunk:
		n += 1 // tag and wire
		n += proto.SizeVarint(uint64(len(x.Chunk)))
		n += len(x.Chunk)
	case nil:
	default:
		panic(fmt.Sprintf("proto: unexpected type %T in oneof", x))
	}
	return n
}

type ProxyConnect struct {
	Target               string   `protobuf:"bytes,1,opt,name=target,proto3" json:"target,omitempty"`
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_unrecognized     []byte   `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *ProxyConnect) Reset()         { *m = ProxyConnect{} }
func (m *ProxyConnect) String() string { return proto.CompactTextString(m) }
func (*ProxyConnect) ProtoMessage()    {}
func (*ProxyConnect) Descriptor() ([]byte, []int) {
	return fileDescriptor_protocol_536ecc5858538e00, []int{1}
}
func (m *ProxyConnect) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_ProxyConnect.Unmarshal(m, b)
}
func (m *ProxyConnect) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_ProxyConnect.Marshal(b, m, deterministic)
}
func (dst *ProxyConnect) XXX_Merge(src proto.Message) {
	xxx_messageInfo_ProxyConnect.Merge(dst, src)
}
func (m *ProxyConnect) XXX_Size() int {
	return xxx_messageInfo_ProxyConnect.Size(m)
}
func (m *ProxyConnect) XXX_DiscardUnknown() {
	xxx_messageInfo_ProxyConnect.DiscardUnknown(m)
}

var xxx_messageInfo_ProxyConnect proto.InternalMessageInfo

func (m *ProxyConnect) GetTarget() string {
	if m != nil {
		return m.Target
	}
	return ""
}

type ProxyConnected struct {
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_unrecognized     []byte   `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *ProxyConnected) Reset()         { *m = ProxyConnected{} }
func (m *ProxyConnected) String() string { return proto.CompactTextString(m) }
func (*ProxyConnected) ProtoMessage()    {}
func (*ProxyConnected) Descriptor() ([]byte, []int) {
	return fileDescriptor_protocol_536ecc5858538e00, []int{2}
}
func (m *ProxyConnected) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_ProxyConnected.Unmarshal(m, b)
}
func (m *ProxyConnected) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_ProxyConnected.Marshal(b, m, deterministic)
}
func (dst *ProxyConnected) XXX_Merge(src proto.Message) {
	xxx_messageInfo_ProxyConnected.Merge(dst, src)
}
func (m *ProxyConnected) XXX_Size() int {
	return xxx_messageInfo_ProxyConnected.Size(m)
}
func (m *ProxyConnected) XXX_DiscardUnknown() {
	xxx_messageInfo_ProxyConnected.DiscardUnknown(m)
}

var xxx_messageInfo_ProxyConnected proto.InternalMessageInfo

type ProxyResponse struct {
	// Types that are valid to be assigned to Res:
	//	*ProxyResponse_Connected
	//	*ProxyResponse_Chunk
	//	*ProxyResponse_Eof
	Res                  isProxyResponse_Res `protobuf_oneof:"res"`
	XXX_NoUnkeyedLiteral struct{}            `json:"-"`
	XXX_unrecognized     []byte              `json:"-"`
	XXX_sizecache        int32               `json:"-"`
}

func (m *ProxyResponse) Reset()         { *m = ProxyResponse{} }
func (m *ProxyResponse) String() string { return proto.CompactTextString(m) }
func (*ProxyResponse) ProtoMessage()    {}
func (*ProxyResponse) Descriptor() ([]byte, []int) {
	return fileDescriptor_protocol_536ecc5858538e00, []int{3}
}
func (m *ProxyResponse) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_ProxyResponse.Unmarshal(m, b)
}
func (m *ProxyResponse) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_ProxyResponse.Marshal(b, m, deterministic)
}
func (dst *ProxyResponse) XXX_Merge(src proto.Message) {
	xxx_messageInfo_ProxyResponse.Merge(dst, src)
}
func (m *ProxyResponse) XXX_Size() int {
	return xxx_messageInfo_ProxyResponse.Size(m)
}
func (m *ProxyResponse) XXX_DiscardUnknown() {
	xxx_messageInfo_ProxyResponse.DiscardUnknown(m)
}

var xxx_messageInfo_ProxyResponse proto.InternalMessageInfo

type isProxyResponse_Res interface {
	isProxyResponse_Res()
}

type ProxyResponse_Connected struct {
	Connected *ProxyConnected `protobuf:"bytes,1,opt,name=connected,proto3,oneof"`
}
type ProxyResponse_Chunk struct {
	Chunk []byte `protobuf:"bytes,2,opt,name=chunk,proto3,oneof"`
}
type ProxyResponse_Eof struct {
	Eof string `protobuf:"bytes,3,opt,name=eof,proto3,oneof"`
}

func (*ProxyResponse_Connected) isProxyResponse_Res() {}
func (*ProxyResponse_Chunk) isProxyResponse_Res()     {}
func (*ProxyResponse_Eof) isProxyResponse_Res()       {}

func (m *ProxyResponse) GetRes() isProxyResponse_Res {
	if m != nil {
		return m.Res
	}
	return nil
}

func (m *ProxyResponse) GetConnected() *ProxyConnected {
	if x, ok := m.GetRes().(*ProxyResponse_Connected); ok {
		return x.Connected
	}
	return nil
}

func (m *ProxyResponse) GetChunk() []byte {
	if x, ok := m.GetRes().(*ProxyResponse_Chunk); ok {
		return x.Chunk
	}
	return nil
}

func (m *ProxyResponse) GetEof() string {
	if x, ok := m.GetRes().(*ProxyResponse_Eof); ok {
		return x.Eof
	}
	return ""
}

// XXX_OneofFuncs is for the internal use of the proto package.
func (*ProxyResponse) XXX_OneofFuncs() (func(msg proto.Message, b *proto.Buffer) error, func(msg proto.Message, tag, wire int, b *proto.Buffer) (bool, error), func(msg proto.Message) (n int), []interface{}) {
	return _ProxyResponse_OneofMarshaler, _ProxyResponse_OneofUnmarshaler, _ProxyResponse_OneofSizer, []interface{}{
		(*ProxyResponse_Connected)(nil),
		(*ProxyResponse_Chunk)(nil),
		(*ProxyResponse_Eof)(nil),
	}
}

func _ProxyResponse_OneofMarshaler(msg proto.Message, b *proto.Buffer) error {
	m := msg.(*ProxyResponse)
	// res
	switch x := m.Res.(type) {
	case *ProxyResponse_Connected:
		b.EncodeVarint(1<<3 | proto.WireBytes)
		if err := b.EncodeMessage(x.Connected); err != nil {
			return err
		}
	case *ProxyResponse_Chunk:
		b.EncodeVarint(2<<3 | proto.WireBytes)
		b.EncodeRawBytes(x.Chunk)
	case *ProxyResponse_Eof:
		b.EncodeVarint(3<<3 | proto.WireBytes)
		b.EncodeStringBytes(x.Eof)
	case nil:
	default:
		return fmt.Errorf("ProxyResponse.Res has unexpected type %T", x)
	}
	return nil
}

func _ProxyResponse_OneofUnmarshaler(msg proto.Message, tag, wire int, b *proto.Buffer) (bool, error) {
	m := msg.(*ProxyResponse)
	switch tag {
	case 1: // res.connected
		if wire != proto.WireBytes {
			return true, proto.ErrInternalBadWireType
		}
		msg := new(ProxyConnected)
		err := b.DecodeMessage(msg)
		m.Res = &ProxyResponse_Connected{msg}
		return true, err
	case 2: // res.chunk
		if wire != proto.WireBytes {
			return true, proto.ErrInternalBadWireType
		}
		x, err := b.DecodeRawBytes(true)
		m.Res = &ProxyResponse_Chunk{x}
		return true, err
	case 3: // res.eof
		if wire != proto.WireBytes {
			return true, proto.ErrInternalBadWireType
		}
		x, err := b.DecodeStringBytes()
		m.Res = &ProxyResponse_Eof{x}
		return true, err
	default:
		return false, nil
	}
}

func _ProxyResponse_OneofSizer(msg proto.Message) (n int) {
	m := msg.(*ProxyResponse)
	// res
	switch x := m.Res.(type) {
	case *ProxyResponse_Connected:
		s := proto.Size(x.Connected)
		n += 1 // tag and wire
		n += proto.SizeVarint(uint64(s))
		n += s
	case *ProxyResponse_Chunk:
		n += 1 // tag and wire
		n += proto.SizeVarint(uint64(len(x.Chunk)))
		n += len(x.Chunk)
	case *ProxyResponse_Eof:
		n += 1 // tag and wire
		n += proto.SizeVarint(uint64(len(x.Eof)))
		n += len(x.Eof)
	case nil:
	default:
		panic(fmt.Sprintf("proto: unexpected type %T in oneof", x))
	}
	return n
}

func init() {
	proto.RegisterType((*ProxyRequest)(nil), "fwd.ProxyRequest")
	proto.RegisterType((*ProxyConnect)(nil), "fwd.ProxyConnect")
	proto.RegisterType((*ProxyConnected)(nil), "fwd.ProxyConnected")
	proto.RegisterType((*ProxyResponse)(nil), "fwd.ProxyResponse")
}

func init() { proto.RegisterFile("protocol.proto", fileDescriptor_protocol_536ecc5858538e00) }

var fileDescriptor_protocol_536ecc5858538e00 = []byte{
	// 239 bytes of a gzipped FileDescriptorProto
	0x1f, 0x8b, 0x08, 0x00, 0x00, 0x00, 0x00, 0x00, 0x02, 0xff, 0x74, 0x90, 0xc1, 0x4b, 0xc3, 0x30,
	0x14, 0xc6, 0x13, 0x4b, 0x37, 0xf6, 0x9c, 0x43, 0x9f, 0x20, 0xc5, 0xd3, 0xc8, 0x41, 0x7a, 0xb1,
	0xc8, 0xe6, 0x4d, 0xbc, 0x6c, 0x97, 0x1e, 0x25, 0x47, 0xf1, 0xa2, 0xc9, 0xab, 0x82, 0xd2, 0x6c,
	0x4d, 0x46, 0xf4, 0xbf, 0x97, 0x26, 0xa9, 0x16, 0xc4, 0x5b, 0xbe, 0xbc, 0x1f, 0x1f, 0x3f, 0x3e,
	0x58, 0xec, 0x3a, 0xe3, 0x8c, 0x32, 0x1f, 0x55, 0x78, 0x60, 0xd6, 0x78, 0x2d, 0x9e, 0x60, 0xfe,
	0xd0, 0x99, 0xcf, 0x2f, 0x49, 0xfb, 0x03, 0x59, 0x87, 0xd7, 0x30, 0x55, 0xa6, 0x6d, 0x49, 0xb9,
	0x82, 0x2f, 0x79, 0x79, 0xbc, 0x3a, 0xab, 0x1a, 0xaf, 0xab, 0xc0, 0x6c, 0xe3, 0xa1, 0x66, 0x72,
	0x60, 0xf0, 0x02, 0x72, 0xf5, 0x76, 0x68, 0xdf, 0x8b, 0xa3, 0x25, 0x2f, 0xe7, 0x35, 0x93, 0x31,
	0x6e, 0x72, 0xc8, 0x3a, 0xda, 0x8b, 0xab, 0xd4, 0xbe, 0xfd, 0xc1, 0x27, 0xee, 0xb9, 0x7b, 0xa5,
	0x58, 0x3e, 0x93, 0x29, 0x89, 0x53, 0x58, 0x8c, 0x39, 0xd2, 0xc2, 0xc3, 0x49, 0xf2, 0xb2, 0x3b,
	0xd3, 0x5a, 0xc2, 0x35, 0xcc, 0xd4, 0x70, 0x4d, 0x6a, 0xe7, 0x7f, 0xd4, 0x48, 0xd7, 0x4c, 0xfe,
	0x72, 0xff, 0xe9, 0x21, 0x42, 0x46, 0xa6, 0x29, 0xb2, 0x5e, 0xa2, 0x66, 0xb2, 0x0f, 0x51, 0xd9,
	0xae, 0xee, 0x21, 0x0f, 0x8d, 0x78, 0x3b, 0x3c, 0x46, 0x0b, 0xa4, 0x95, 0x2e, 0x71, 0xfc, 0x15,
	0x05, 0x05, 0x2b, 0xf9, 0x0d, 0xdf, 0x4c, 0x1f, 0xf3, 0xea, 0xae, 0xf1, 0xfa, 0x65, 0x12, 0x46,
	0x5e, 0x7f, 0x07, 0x00, 0x00, 0xff, 0xff, 0x1f, 0x6c, 0x34, 0xb8, 0x76, 0x01, 0x00, 0x00,
}
