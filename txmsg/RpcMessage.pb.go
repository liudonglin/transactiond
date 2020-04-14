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

type TransactionModel int32

const (
	// tcc事务模式
	TransactionModel_TCC TransactionModel = 0
	// lcn事务模式
	TransactionModel_LCN TransactionModel = 1
)

var TransactionModel_name = map[int32]string{
	0: "TCC",
	1: "LCN",
}

var TransactionModel_value = map[string]int32{
	"TCC": 0,
	"LCN": 1,
}

func (x TransactionModel) String() string {
	return proto.EnumName(TransactionModel_name, int32(x))
}

func (TransactionModel) EnumDescriptor() ([]byte, []int) {
	return fileDescriptor_a49baf2e5143a5e0, []int{0}
}

type TransactionState int32

const (
	// 事务回滚
	TransactionState_Rollback TransactionState = 0
	// 事务提交
	TransactionState_Commit TransactionState = 1
)

var TransactionState_name = map[int32]string{
	0: "Rollback",
	1: "Commit",
}

var TransactionState_value = map[string]int32{
	"Rollback": 0,
	"Commit":   1,
}

func (x TransactionState) String() string {
	return proto.EnumName(TransactionState_name, int32(x))
}

func (TransactionState) EnumDescriptor() ([]byte, []int) {
	return fileDescriptor_a49baf2e5143a5e0, []int{1}
}

type ResponseMessage struct {
	Code                 string   `protobuf:"bytes,1,opt,name=code,proto3" json:"code,omitempty"`
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_unrecognized     []byte   `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *ResponseMessage) Reset()         { *m = ResponseMessage{} }
func (m *ResponseMessage) String() string { return proto.CompactTextString(m) }
func (*ResponseMessage) ProtoMessage()    {}
func (*ResponseMessage) Descriptor() ([]byte, []int) {
	return fileDescriptor_a49baf2e5143a5e0, []int{0}
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

func (m *ResponseMessage) GetCode() string {
	if m != nil {
		return m.Code
	}
	return ""
}

type CreateGroupMessage struct {
	GroupId              string   `protobuf:"bytes,1,opt,name=groupId,proto3" json:"groupId,omitempty"`
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_unrecognized     []byte   `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *CreateGroupMessage) Reset()         { *m = CreateGroupMessage{} }
func (m *CreateGroupMessage) String() string { return proto.CompactTextString(m) }
func (*CreateGroupMessage) ProtoMessage()    {}
func (*CreateGroupMessage) Descriptor() ([]byte, []int) {
	return fileDescriptor_a49baf2e5143a5e0, []int{1}
}

func (m *CreateGroupMessage) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_CreateGroupMessage.Unmarshal(m, b)
}
func (m *CreateGroupMessage) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_CreateGroupMessage.Marshal(b, m, deterministic)
}
func (m *CreateGroupMessage) XXX_Merge(src proto.Message) {
	xxx_messageInfo_CreateGroupMessage.Merge(m, src)
}
func (m *CreateGroupMessage) XXX_Size() int {
	return xxx_messageInfo_CreateGroupMessage.Size(m)
}
func (m *CreateGroupMessage) XXX_DiscardUnknown() {
	xxx_messageInfo_CreateGroupMessage.DiscardUnknown(m)
}

var xxx_messageInfo_CreateGroupMessage proto.InternalMessageInfo

func (m *CreateGroupMessage) GetGroupId() string {
	if m != nil {
		return m.GroupId
	}
	return ""
}

type JoinGroupMessage struct {
	GroupId              string           `protobuf:"bytes,1,opt,name=groupId,proto3" json:"groupId,omitempty"`
	UnitId               string           `protobuf:"bytes,2,opt,name=unitId,proto3" json:"unitId,omitempty"`
	Model                TransactionModel `protobuf:"varint,3,opt,name=model,proto3,enum=txmsg.TransactionModel" json:"model,omitempty"`
	XXX_NoUnkeyedLiteral struct{}         `json:"-"`
	XXX_unrecognized     []byte           `json:"-"`
	XXX_sizecache        int32            `json:"-"`
}

func (m *JoinGroupMessage) Reset()         { *m = JoinGroupMessage{} }
func (m *JoinGroupMessage) String() string { return proto.CompactTextString(m) }
func (*JoinGroupMessage) ProtoMessage()    {}
func (*JoinGroupMessage) Descriptor() ([]byte, []int) {
	return fileDescriptor_a49baf2e5143a5e0, []int{2}
}

func (m *JoinGroupMessage) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_JoinGroupMessage.Unmarshal(m, b)
}
func (m *JoinGroupMessage) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_JoinGroupMessage.Marshal(b, m, deterministic)
}
func (m *JoinGroupMessage) XXX_Merge(src proto.Message) {
	xxx_messageInfo_JoinGroupMessage.Merge(m, src)
}
func (m *JoinGroupMessage) XXX_Size() int {
	return xxx_messageInfo_JoinGroupMessage.Size(m)
}
func (m *JoinGroupMessage) XXX_DiscardUnknown() {
	xxx_messageInfo_JoinGroupMessage.DiscardUnknown(m)
}

var xxx_messageInfo_JoinGroupMessage proto.InternalMessageInfo

func (m *JoinGroupMessage) GetGroupId() string {
	if m != nil {
		return m.GroupId
	}
	return ""
}

func (m *JoinGroupMessage) GetUnitId() string {
	if m != nil {
		return m.UnitId
	}
	return ""
}

func (m *JoinGroupMessage) GetModel() TransactionModel {
	if m != nil {
		return m.Model
	}
	return TransactionModel_TCC
}

type NotifyGroupMessage struct {
	GroupId              string           `protobuf:"bytes,1,opt,name=groupId,proto3" json:"groupId,omitempty"`
	State                TransactionState `protobuf:"varint,2,opt,name=state,proto3,enum=txmsg.TransactionState" json:"state,omitempty"`
	XXX_NoUnkeyedLiteral struct{}         `json:"-"`
	XXX_unrecognized     []byte           `json:"-"`
	XXX_sizecache        int32            `json:"-"`
}

func (m *NotifyGroupMessage) Reset()         { *m = NotifyGroupMessage{} }
func (m *NotifyGroupMessage) String() string { return proto.CompactTextString(m) }
func (*NotifyGroupMessage) ProtoMessage()    {}
func (*NotifyGroupMessage) Descriptor() ([]byte, []int) {
	return fileDescriptor_a49baf2e5143a5e0, []int{3}
}

func (m *NotifyGroupMessage) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_NotifyGroupMessage.Unmarshal(m, b)
}
func (m *NotifyGroupMessage) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_NotifyGroupMessage.Marshal(b, m, deterministic)
}
func (m *NotifyGroupMessage) XXX_Merge(src proto.Message) {
	xxx_messageInfo_NotifyGroupMessage.Merge(m, src)
}
func (m *NotifyGroupMessage) XXX_Size() int {
	return xxx_messageInfo_NotifyGroupMessage.Size(m)
}
func (m *NotifyGroupMessage) XXX_DiscardUnknown() {
	xxx_messageInfo_NotifyGroupMessage.DiscardUnknown(m)
}

var xxx_messageInfo_NotifyGroupMessage proto.InternalMessageInfo

func (m *NotifyGroupMessage) GetGroupId() string {
	if m != nil {
		return m.GroupId
	}
	return ""
}

func (m *NotifyGroupMessage) GetState() TransactionState {
	if m != nil {
		return m.State
	}
	return TransactionState_Rollback
}

type NotifyUnitMessage struct {
	GroupId              string           `protobuf:"bytes,1,opt,name=groupId,proto3" json:"groupId,omitempty"`
	State                TransactionState `protobuf:"varint,2,opt,name=state,proto3,enum=txmsg.TransactionState" json:"state,omitempty"`
	XXX_NoUnkeyedLiteral struct{}         `json:"-"`
	XXX_unrecognized     []byte           `json:"-"`
	XXX_sizecache        int32            `json:"-"`
}

func (m *NotifyUnitMessage) Reset()         { *m = NotifyUnitMessage{} }
func (m *NotifyUnitMessage) String() string { return proto.CompactTextString(m) }
func (*NotifyUnitMessage) ProtoMessage()    {}
func (*NotifyUnitMessage) Descriptor() ([]byte, []int) {
	return fileDescriptor_a49baf2e5143a5e0, []int{4}
}

func (m *NotifyUnitMessage) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_NotifyUnitMessage.Unmarshal(m, b)
}
func (m *NotifyUnitMessage) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_NotifyUnitMessage.Marshal(b, m, deterministic)
}
func (m *NotifyUnitMessage) XXX_Merge(src proto.Message) {
	xxx_messageInfo_NotifyUnitMessage.Merge(m, src)
}
func (m *NotifyUnitMessage) XXX_Size() int {
	return xxx_messageInfo_NotifyUnitMessage.Size(m)
}
func (m *NotifyUnitMessage) XXX_DiscardUnknown() {
	xxx_messageInfo_NotifyUnitMessage.DiscardUnknown(m)
}

var xxx_messageInfo_NotifyUnitMessage proto.InternalMessageInfo

func (m *NotifyUnitMessage) GetGroupId() string {
	if m != nil {
		return m.GroupId
	}
	return ""
}

func (m *NotifyUnitMessage) GetState() TransactionState {
	if m != nil {
		return m.State
	}
	return TransactionState_Rollback
}

func init() {
	proto.RegisterEnum("txmsg.TransactionModel", TransactionModel_name, TransactionModel_value)
	proto.RegisterEnum("txmsg.TransactionState", TransactionState_name, TransactionState_value)
	proto.RegisterType((*ResponseMessage)(nil), "txmsg.ResponseMessage")
	proto.RegisterType((*CreateGroupMessage)(nil), "txmsg.CreateGroupMessage")
	proto.RegisterType((*JoinGroupMessage)(nil), "txmsg.JoinGroupMessage")
	proto.RegisterType((*NotifyGroupMessage)(nil), "txmsg.NotifyGroupMessage")
	proto.RegisterType((*NotifyUnitMessage)(nil), "txmsg.NotifyUnitMessage")
}

func init() {
	proto.RegisterFile("RpcMessage.proto", fileDescriptor_a49baf2e5143a5e0)
}

var fileDescriptor_a49baf2e5143a5e0 = []byte{
	// 351 bytes of a gzipped FileDescriptorProto
	0x1f, 0x8b, 0x08, 0x00, 0x00, 0x00, 0x00, 0x00, 0x02, 0xff, 0xac, 0x93, 0x4f, 0x4f, 0xf2, 0x40,
	0x10, 0xc6, 0x29, 0xbc, 0xc0, 0xcb, 0xf0, 0xfe, 0x59, 0xe7, 0x80, 0x95, 0x13, 0x69, 0x34, 0x21,
	0x44, 0x1a, 0x83, 0x57, 0x0f, 0x84, 0x1e, 0x0c, 0x46, 0x88, 0x29, 0x78, 0xd3, 0xc3, 0xd2, 0xae,
	0xcd, 0x2a, 0xdd, 0x6d, 0xba, 0x8b, 0xd1, 0x0f, 0xe2, 0xf7, 0x35, 0x2d, 0x2d, 0x08, 0x48, 0xc2,
	0xc1, 0xdb, 0x4c, 0x9e, 0xa7, 0xcf, 0xcc, 0xe4, 0xd7, 0x05, 0xe2, 0x46, 0xde, 0x88, 0x29, 0x45,
	0x03, 0x66, 0x47, 0xb1, 0xd4, 0x12, 0xcb, 0xfa, 0x2d, 0x54, 0x81, 0x75, 0x06, 0xff, 0x5d, 0xa6,
	0x22, 0x29, 0x14, 0xcb, 0x74, 0x44, 0xf8, 0xe5, 0x49, 0x9f, 0x99, 0x46, 0xcb, 0x68, 0xd7, 0xdc,
	0xb4, 0xb6, 0x6c, 0x40, 0x27, 0x66, 0x54, 0xb3, 0xeb, 0x58, 0x2e, 0xa2, 0xdc, 0x69, 0x42, 0x35,
	0x48, 0xfa, 0xa1, 0x9f, 0x99, 0xf3, 0xd6, 0x52, 0x40, 0x6e, 0x24, 0x17, 0x87, 0xb9, 0xb1, 0x01,
	0x95, 0x85, 0xe0, 0x7a, 0xe8, 0x9b, 0xc5, 0x54, 0xc8, 0x3a, 0xec, 0x42, 0x39, 0x94, 0x3e, 0x9b,
	0x9b, 0xa5, 0x96, 0xd1, 0xfe, 0xd7, 0x3b, 0xb6, 0xd3, 0x9d, 0xed, 0x69, 0x4c, 0x85, 0xa2, 0x9e,
	0xe6, 0x52, 0x8c, 0x12, 0xd9, 0x5d, 0xba, 0xac, 0x47, 0xc0, 0xb1, 0xd4, 0xfc, 0xe9, 0xfd, 0xc0,
	0xb1, 0x5d, 0x28, 0x2b, 0x4d, 0x35, 0x4b, 0xa7, 0x7e, 0x1b, 0x3f, 0x49, 0x64, 0x77, 0xe9, 0xb2,
	0x1e, 0xe0, 0x68, 0x19, 0x7f, 0x2f, 0xb8, 0xfe, 0xe9, 0xf4, 0xce, 0x29, 0x90, 0xed, 0xbb, 0xb0,
	0x0a, 0xa5, 0xa9, 0xe3, 0x90, 0x42, 0x52, 0xdc, 0x3a, 0x63, 0x62, 0x74, 0xce, 0x37, 0x5c, 0x69,
	0x00, 0xfe, 0x81, 0xdf, 0xae, 0x9c, 0xcf, 0x67, 0xd4, 0x7b, 0x21, 0x05, 0x04, 0xa8, 0x38, 0x32,
	0x0c, 0xb9, 0x26, 0x46, 0xef, 0xa3, 0x08, 0x7f, 0x47, 0x54, 0xd0, 0x80, 0x4d, 0x58, 0xfc, 0xca,
	0x3d, 0x86, 0x7d, 0xa8, 0x7b, 0x6b, 0x8e, 0x78, 0x92, 0x2d, 0xb5, 0xcb, 0xb6, 0xd9, 0xc8, 0xa4,
	0xed, 0xbf, 0xe3, 0x0a, 0x6a, 0xcf, 0x39, 0x59, 0xcc, 0x8f, 0xda, 0x66, 0xbd, 0xf7, 0xeb, 0x3e,
	0xd4, 0xc5, 0x1a, 0xd1, 0x6a, 0xfe, 0x2e, 0xb6, 0xbd, 0x09, 0x03, 0x00, 0xb1, 0xa2, 0x80, 0xe6,
	0x46, 0xc0, 0x17, 0x30, 0xcd, 0xbd, 0xca, 0x85, 0x31, 0x28, 0xde, 0x19, 0xb3, 0x4a, 0xfa, 0x0c,
	0x2e, 0x3f, 0x03, 0x00, 0x00, 0xff, 0xff, 0x96, 0x53, 0x26, 0x8b, 0x1a, 0x03, 0x00, 0x00,
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
	CreateGroup(ctx context.Context, in *CreateGroupMessage, opts ...grpc.CallOption) (*ResponseMessage, error)
	JoinGroup(ctx context.Context, in *JoinGroupMessage, opts ...grpc.CallOption) (*ResponseMessage, error)
	NotifyGroup(ctx context.Context, in *NotifyGroupMessage, opts ...grpc.CallOption) (*ResponseMessage, error)
	NotifyUnit(ctx context.Context, in *NotifyUnitMessage, opts ...grpc.CallOption) (ManageService_NotifyUnitClient, error)
}

type manageServiceClient struct {
	cc grpc.ClientConnInterface
}

func NewManageServiceClient(cc grpc.ClientConnInterface) ManageServiceClient {
	return &manageServiceClient{cc}
}

func (c *manageServiceClient) CreateGroup(ctx context.Context, in *CreateGroupMessage, opts ...grpc.CallOption) (*ResponseMessage, error) {
	out := new(ResponseMessage)
	err := c.cc.Invoke(ctx, "/txmsg.ManageService/createGroup", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *manageServiceClient) JoinGroup(ctx context.Context, in *JoinGroupMessage, opts ...grpc.CallOption) (*ResponseMessage, error) {
	out := new(ResponseMessage)
	err := c.cc.Invoke(ctx, "/txmsg.ManageService/joinGroup", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *manageServiceClient) NotifyGroup(ctx context.Context, in *NotifyGroupMessage, opts ...grpc.CallOption) (*ResponseMessage, error) {
	out := new(ResponseMessage)
	err := c.cc.Invoke(ctx, "/txmsg.ManageService/notifyGroup", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *manageServiceClient) NotifyUnit(ctx context.Context, in *NotifyUnitMessage, opts ...grpc.CallOption) (ManageService_NotifyUnitClient, error) {
	stream, err := c.cc.NewStream(ctx, &_ManageService_serviceDesc.Streams[0], "/txmsg.ManageService/notifyUnit", opts...)
	if err != nil {
		return nil, err
	}
	x := &manageServiceNotifyUnitClient{stream}
	if err := x.ClientStream.SendMsg(in); err != nil {
		return nil, err
	}
	if err := x.ClientStream.CloseSend(); err != nil {
		return nil, err
	}
	return x, nil
}

type ManageService_NotifyUnitClient interface {
	Recv() (*NotifyUnitMessage, error)
	grpc.ClientStream
}

type manageServiceNotifyUnitClient struct {
	grpc.ClientStream
}

func (x *manageServiceNotifyUnitClient) Recv() (*NotifyUnitMessage, error) {
	m := new(NotifyUnitMessage)
	if err := x.ClientStream.RecvMsg(m); err != nil {
		return nil, err
	}
	return m, nil
}

// ManageServiceServer is the server API for ManageService service.
type ManageServiceServer interface {
	CreateGroup(context.Context, *CreateGroupMessage) (*ResponseMessage, error)
	JoinGroup(context.Context, *JoinGroupMessage) (*ResponseMessage, error)
	NotifyGroup(context.Context, *NotifyGroupMessage) (*ResponseMessage, error)
	NotifyUnit(*NotifyUnitMessage, ManageService_NotifyUnitServer) error
}

// UnimplementedManageServiceServer can be embedded to have forward compatible implementations.
type UnimplementedManageServiceServer struct {
}

func (*UnimplementedManageServiceServer) CreateGroup(ctx context.Context, req *CreateGroupMessage) (*ResponseMessage, error) {
	return nil, status.Errorf(codes.Unimplemented, "method CreateGroup not implemented")
}
func (*UnimplementedManageServiceServer) JoinGroup(ctx context.Context, req *JoinGroupMessage) (*ResponseMessage, error) {
	return nil, status.Errorf(codes.Unimplemented, "method JoinGroup not implemented")
}
func (*UnimplementedManageServiceServer) NotifyGroup(ctx context.Context, req *NotifyGroupMessage) (*ResponseMessage, error) {
	return nil, status.Errorf(codes.Unimplemented, "method NotifyGroup not implemented")
}
func (*UnimplementedManageServiceServer) NotifyUnit(req *NotifyUnitMessage, srv ManageService_NotifyUnitServer) error {
	return status.Errorf(codes.Unimplemented, "method NotifyUnit not implemented")
}

func RegisterManageServiceServer(s *grpc.Server, srv ManageServiceServer) {
	s.RegisterService(&_ManageService_serviceDesc, srv)
}

func _ManageService_CreateGroup_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(CreateGroupMessage)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(ManageServiceServer).CreateGroup(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/txmsg.ManageService/CreateGroup",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(ManageServiceServer).CreateGroup(ctx, req.(*CreateGroupMessage))
	}
	return interceptor(ctx, in, info, handler)
}

func _ManageService_JoinGroup_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(JoinGroupMessage)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(ManageServiceServer).JoinGroup(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/txmsg.ManageService/JoinGroup",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(ManageServiceServer).JoinGroup(ctx, req.(*JoinGroupMessage))
	}
	return interceptor(ctx, in, info, handler)
}

func _ManageService_NotifyGroup_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(NotifyGroupMessage)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(ManageServiceServer).NotifyGroup(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/txmsg.ManageService/NotifyGroup",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(ManageServiceServer).NotifyGroup(ctx, req.(*NotifyGroupMessage))
	}
	return interceptor(ctx, in, info, handler)
}

func _ManageService_NotifyUnit_Handler(srv interface{}, stream grpc.ServerStream) error {
	m := new(NotifyUnitMessage)
	if err := stream.RecvMsg(m); err != nil {
		return err
	}
	return srv.(ManageServiceServer).NotifyUnit(m, &manageServiceNotifyUnitServer{stream})
}

type ManageService_NotifyUnitServer interface {
	Send(*NotifyUnitMessage) error
	grpc.ServerStream
}

type manageServiceNotifyUnitServer struct {
	grpc.ServerStream
}

func (x *manageServiceNotifyUnitServer) Send(m *NotifyUnitMessage) error {
	return x.ServerStream.SendMsg(m)
}

var _ManageService_serviceDesc = grpc.ServiceDesc{
	ServiceName: "txmsg.ManageService",
	HandlerType: (*ManageServiceServer)(nil),
	Methods: []grpc.MethodDesc{
		{
			MethodName: "createGroup",
			Handler:    _ManageService_CreateGroup_Handler,
		},
		{
			MethodName: "joinGroup",
			Handler:    _ManageService_JoinGroup_Handler,
		},
		{
			MethodName: "notifyGroup",
			Handler:    _ManageService_NotifyGroup_Handler,
		},
	},
	Streams: []grpc.StreamDesc{
		{
			StreamName:    "notifyUnit",
			Handler:       _ManageService_NotifyUnit_Handler,
			ServerStreams: true,
		},
	},
	Metadata: "RpcMessage.proto",
}
