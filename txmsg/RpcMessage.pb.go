// Code generated by protoc-gen-go. DO NOT EDIT.
// source: RpcMessage.proto

package txmsg

import (
	context "context"
	fmt "fmt"
	proto "github.com/golang/protobuf/proto"
	grpc "google.golang.org/grpc"
	codes "google.golang.org/grpc/codes"
	status "google.golang.org/grpc/status"
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

type ActionType int32

const (
	ActionType_CreateGroup ActionType = 0
)

var ActionType_name = map[int32]string{
	0: "CreateGroup",
}

var ActionType_value = map[string]int32{
	"CreateGroup": 0,
}

func (x ActionType) String() string {
	return proto.EnumName(ActionType_name, int32(x))
}

func (ActionType) EnumDescriptor() ([]byte, []int) {
	return fileDescriptor_a49baf2e5143a5e0, []int{0}
}

type RpcMessage struct {
	GroupId              string     `protobuf:"bytes,1,opt,name=groupId,proto3" json:"groupId,omitempty"`
	ActionType           ActionType `protobuf:"varint,2,opt,name=actionType,proto3,enum=txmsg.ActionType" json:"actionType,omitempty"`
	XXX_NoUnkeyedLiteral struct{}   `json:"-"`
	XXX_unrecognized     []byte     `json:"-"`
	XXX_sizecache        int32      `json:"-"`
}

func (m *RpcMessage) Reset()         { *m = RpcMessage{} }
func (m *RpcMessage) String() string { return proto.CompactTextString(m) }
func (*RpcMessage) ProtoMessage()    {}
func (*RpcMessage) Descriptor() ([]byte, []int) {
	return fileDescriptor_a49baf2e5143a5e0, []int{0}
}

func (m *RpcMessage) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_RpcMessage.Unmarshal(m, b)
}
func (m *RpcMessage) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_RpcMessage.Marshal(b, m, deterministic)
}
func (m *RpcMessage) XXX_Merge(src proto.Message) {
	xxx_messageInfo_RpcMessage.Merge(m, src)
}
func (m *RpcMessage) XXX_Size() int {
	return xxx_messageInfo_RpcMessage.Size(m)
}
func (m *RpcMessage) XXX_DiscardUnknown() {
	xxx_messageInfo_RpcMessage.DiscardUnknown(m)
}

var xxx_messageInfo_RpcMessage proto.InternalMessageInfo

func (m *RpcMessage) GetGroupId() string {
	if m != nil {
		return m.GroupId
	}
	return ""
}

func (m *RpcMessage) GetActionType() ActionType {
	if m != nil {
		return m.ActionType
	}
	return ActionType_CreateGroup
}

func init() {
	proto.RegisterEnum("txmsg.ActionType", ActionType_name, ActionType_value)
	proto.RegisterType((*RpcMessage)(nil), "txmsg.RpcMessage")
}

func init() {
	proto.RegisterFile("RpcMessage.proto", fileDescriptor_a49baf2e5143a5e0)
}

var fileDescriptor_a49baf2e5143a5e0 = []byte{
	// 172 bytes of a gzipped FileDescriptorProto
	0x1f, 0x8b, 0x08, 0x00, 0x00, 0x00, 0x00, 0x00, 0x02, 0xff, 0xe2, 0x12, 0x08, 0x2a, 0x48, 0xf6,
	0x4d, 0x2d, 0x2e, 0x4e, 0x4c, 0x4f, 0xd5, 0x2b, 0x28, 0xca, 0x2f, 0xc9, 0x17, 0x62, 0x2d, 0xa9,
	0xc8, 0x2d, 0x4e, 0x57, 0x8a, 0xe4, 0xe2, 0x42, 0x48, 0x09, 0x49, 0x70, 0xb1, 0xa7, 0x17, 0xe5,
	0x97, 0x16, 0x78, 0xa6, 0x48, 0x30, 0x2a, 0x30, 0x6a, 0x70, 0x06, 0xc1, 0xb8, 0x42, 0x86, 0x5c,
	0x5c, 0x89, 0xc9, 0x25, 0x99, 0xf9, 0x79, 0x21, 0x95, 0x05, 0xa9, 0x12, 0x4c, 0x0a, 0x8c, 0x1a,
	0x7c, 0x46, 0x82, 0x7a, 0x60, 0x33, 0xf4, 0x1c, 0xe1, 0x12, 0x41, 0x48, 0x8a, 0xb4, 0x64, 0xb9,
	0xb8, 0x10, 0x32, 0x42, 0xfc, 0x5c, 0xdc, 0xce, 0x45, 0xa9, 0x89, 0x25, 0xa9, 0xee, 0x20, 0x13,
	0x05, 0x18, 0x8c, 0x5c, 0xb8, 0x78, 0x7d, 0x13, 0xf3, 0x12, 0xd3, 0x53, 0x83, 0x53, 0x8b, 0xca,
	0x32, 0x93, 0x53, 0x85, 0x8c, 0xb9, 0xb8, 0x8b, 0x53, 0xf3, 0x52, 0x60, 0x6e, 0x81, 0x99, 0x8e,
	0x70, 0x9e, 0x14, 0xa6, 0x90, 0x13, 0x53, 0x00, 0x63, 0x12, 0x1b, 0xd8, 0x47, 0xc6, 0x80, 0x00,
	0x00, 0x00, 0xff, 0xff, 0x51, 0x38, 0x4b, 0x2b, 0xe5, 0x00, 0x00, 0x00,
}

// Reference imports to suppress errors if they are not otherwise used.
var _ context.Context
var _ grpc.ClientConnInterface

// This is a compile-time assertion to ensure that this generated file
// is compatible with the grpc package it is being compiled against.
const _ = grpc.SupportPackageIsVersion6

// ManageServiceClient is the client API for ManageService service.
//
// For semantics around ctx use and closing/ending streaming RPCs, please refer to https://godoc.org/google.golang.org/grpc#ClientConn.NewStream.
type ManageServiceClient interface {
	SendMessage(ctx context.Context, in *RpcMessage, opts ...grpc.CallOption) (*RpcMessage, error)
}

type manageServiceClient struct {
	cc grpc.ClientConnInterface
}

func NewManageServiceClient(cc grpc.ClientConnInterface) ManageServiceClient {
	return &manageServiceClient{cc}
}

func (c *manageServiceClient) SendMessage(ctx context.Context, in *RpcMessage, opts ...grpc.CallOption) (*RpcMessage, error) {
	out := new(RpcMessage)
	err := c.cc.Invoke(ctx, "/txmsg.ManageService/sendMessage", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

// ManageServiceServer is the server API for ManageService service.
type ManageServiceServer interface {
	SendMessage(context.Context, *RpcMessage) (*RpcMessage, error)
}

// UnimplementedManageServiceServer can be embedded to have forward compatible implementations.
type UnimplementedManageServiceServer struct {
}

func (*UnimplementedManageServiceServer) SendMessage(ctx context.Context, req *RpcMessage) (*RpcMessage, error) {
	return nil, status.Errorf(codes.Unimplemented, "method SendMessage not implemented")
}

func RegisterManageServiceServer(s *grpc.Server, srv ManageServiceServer) {
	s.RegisterService(&_ManageService_serviceDesc, srv)
}

func _ManageService_SendMessage_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(RpcMessage)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(ManageServiceServer).SendMessage(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/txmsg.ManageService/SendMessage",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(ManageServiceServer).SendMessage(ctx, req.(*RpcMessage))
	}
	return interceptor(ctx, in, info, handler)
}

var _ManageService_serviceDesc = grpc.ServiceDesc{
	ServiceName: "txmsg.ManageService",
	HandlerType: (*ManageServiceServer)(nil),
	Methods: []grpc.MethodDesc{
		{
			MethodName: "sendMessage",
			Handler:    _ManageService_SendMessage_Handler,
		},
	},
	Streams:  []grpc.StreamDesc{},
	Metadata: "RpcMessage.proto",
}