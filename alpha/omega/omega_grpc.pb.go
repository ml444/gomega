// Code generated by protoc-gen-go-grpc. DO NOT EDIT.
// versions:
// - protoc-gen-go-grpc v1.2.0
// - protoc             v3.21.5
// source: omega.proto

package omega

import (
	context "context"
	grpc "google.golang.org/grpc"
	codes "google.golang.org/grpc/codes"
	status "google.golang.org/grpc/status"
)

// This is a compile-time assertion to ensure that this generated file
// is compatible with the grpc package it is being compiled against.
// Requires gRPC-Go v1.32.0 or later.
const _ = grpc.SupportPackageIsVersion7

// OmegaServiceClient is the client API for OmegaService service.
//
// For semantics around ctx use and closing/ending streaming RPCs, please refer to https://pkg.go.dev/google.golang.org/grpc/?tab=doc#ClientConn.NewStream.
type OmegaServiceClient interface {
	Pub(ctx context.Context, in *PubReq, opts ...grpc.CallOption) (*Response, error)
	Sub(ctx context.Context, in *SubReq, opts ...grpc.CallOption) (*SubRsp, error)
	UpdateSubCfg(ctx context.Context, in *SubCfg, opts ...grpc.CallOption) (*Response, error)
	Consume(ctx context.Context, in *ConsumeReq, opts ...grpc.CallOption) (OmegaService_ConsumeClient, error)
}

type omegaServiceClient struct {
	cc grpc.ClientConnInterface
}

func NewOmegaServiceClient(cc grpc.ClientConnInterface) OmegaServiceClient {
	return &omegaServiceClient{cc}
}

func (c *omegaServiceClient) Pub(ctx context.Context, in *PubReq, opts ...grpc.CallOption) (*Response, error) {
	out := new(Response)
	err := c.cc.Invoke(ctx, "/omega.OmegaService/Pub", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *omegaServiceClient) Sub(ctx context.Context, in *SubReq, opts ...grpc.CallOption) (*SubRsp, error) {
	out := new(SubRsp)
	err := c.cc.Invoke(ctx, "/omega.OmegaService/Sub", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *omegaServiceClient) UpdateSubCfg(ctx context.Context, in *SubCfg, opts ...grpc.CallOption) (*Response, error) {
	out := new(Response)
	err := c.cc.Invoke(ctx, "/omega.OmegaService/UpdateSubCfg", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *omegaServiceClient) Consume(ctx context.Context, in *ConsumeReq, opts ...grpc.CallOption) (OmegaService_ConsumeClient, error) {
	stream, err := c.cc.NewStream(ctx, &OmegaService_ServiceDesc.Streams[0], "/omega.OmegaService/Consume", opts...)
	if err != nil {
		return nil, err
	}
	x := &omegaServiceConsumeClient{stream}
	if err := x.ClientStream.SendMsg(in); err != nil {
		return nil, err
	}
	if err := x.ClientStream.CloseSend(); err != nil {
		return nil, err
	}
	return x, nil
}

type OmegaService_ConsumeClient interface {
	Recv() (*ConsumeRsp, error)
	grpc.ClientStream
}

type omegaServiceConsumeClient struct {
	grpc.ClientStream
}

func (x *omegaServiceConsumeClient) Recv() (*ConsumeRsp, error) {
	m := new(ConsumeRsp)
	if err := x.ClientStream.RecvMsg(m); err != nil {
		return nil, err
	}
	return m, nil
}

// OmegaServiceServer is the server API for OmegaService service.
// All implementations must embed UnimplementedOmegaServiceServer
// for forward compatibility
type OmegaServiceServer interface {
	Pub(context.Context, *PubReq) (*Response, error)
	Sub(context.Context, *SubReq) (*SubRsp, error)
	UpdateSubCfg(context.Context, *SubCfg) (*Response, error)
	Consume(*ConsumeReq, OmegaService_ConsumeServer) error
	mustEmbedUnimplementedOmegaServiceServer()
}

// UnimplementedOmegaServiceServer must be embedded to have forward compatible implementations.
type UnimplementedOmegaServiceServer struct {
}

func (UnimplementedOmegaServiceServer) Pub(context.Context, *PubReq) (*Response, error) {
	return nil, status.Errorf(codes.Unimplemented, "method Pub not implemented")
}
func (UnimplementedOmegaServiceServer) Sub(context.Context, *SubReq) (*SubRsp, error) {
	return nil, status.Errorf(codes.Unimplemented, "method Sub not implemented")
}
func (UnimplementedOmegaServiceServer) UpdateSubCfg(context.Context, *SubCfg) (*Response, error) {
	return nil, status.Errorf(codes.Unimplemented, "method UpdateSubCfg not implemented")
}
func (UnimplementedOmegaServiceServer) Consume(*ConsumeReq, OmegaService_ConsumeServer) error {
	return status.Errorf(codes.Unimplemented, "method Consume not implemented")
}
func (UnimplementedOmegaServiceServer) mustEmbedUnimplementedOmegaServiceServer() {}

// UnsafeOmegaServiceServer may be embedded to opt out of forward compatibility for this service.
// Use of this interface is not recommended, as added methods to OmegaServiceServer will
// result in compilation errors.
type UnsafeOmegaServiceServer interface {
	mustEmbedUnimplementedOmegaServiceServer()
}

func RegisterOmegaServiceServer(s grpc.ServiceRegistrar, srv OmegaServiceServer) {
	s.RegisterService(&OmegaService_ServiceDesc, srv)
}

func _OmegaService_Pub_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(PubReq)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(OmegaServiceServer).Pub(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/omega.OmegaService/Pub",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(OmegaServiceServer).Pub(ctx, req.(*PubReq))
	}
	return interceptor(ctx, in, info, handler)
}

func _OmegaService_Sub_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(SubReq)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(OmegaServiceServer).Sub(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/omega.OmegaService/Sub",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(OmegaServiceServer).Sub(ctx, req.(*SubReq))
	}
	return interceptor(ctx, in, info, handler)
}

func _OmegaService_UpdateSubCfg_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(SubCfg)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(OmegaServiceServer).UpdateSubCfg(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/omega.OmegaService/UpdateSubCfg",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(OmegaServiceServer).UpdateSubCfg(ctx, req.(*SubCfg))
	}
	return interceptor(ctx, in, info, handler)
}

func _OmegaService_Consume_Handler(srv interface{}, stream grpc.ServerStream) error {
	m := new(ConsumeReq)
	if err := stream.RecvMsg(m); err != nil {
		return err
	}
	return srv.(OmegaServiceServer).Consume(m, &omegaServiceConsumeServer{stream})
}

type OmegaService_ConsumeServer interface {
	Send(*ConsumeRsp) error
	grpc.ServerStream
}

type omegaServiceConsumeServer struct {
	grpc.ServerStream
}

func (x *omegaServiceConsumeServer) Send(m *ConsumeRsp) error {
	return x.ServerStream.SendMsg(m)
}

// OmegaService_ServiceDesc is the grpc.ServiceDesc for OmegaService service.
// It's only intended for direct use with grpc.RegisterService,
// and not to be introspected or modified (even as a copy)
var OmegaService_ServiceDesc = grpc.ServiceDesc{
	ServiceName: "omega.OmegaService",
	HandlerType: (*OmegaServiceServer)(nil),
	Methods: []grpc.MethodDesc{
		{
			MethodName: "Pub",
			Handler:    _OmegaService_Pub_Handler,
		},
		{
			MethodName: "Sub",
			Handler:    _OmegaService_Sub_Handler,
		},
		{
			MethodName: "UpdateSubCfg",
			Handler:    _OmegaService_UpdateSubCfg_Handler,
		},
	},
	Streams: []grpc.StreamDesc{
		{
			StreamName:    "Consume",
			Handler:       _OmegaService_Consume_Handler,
			ServerStreams: true,
		},
	},
	Metadata: "omega.proto",
}
