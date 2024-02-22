// Simple service definition for Mangle deductive databases.

// Code generated by protoc-gen-go-grpc. DO NOT EDIT.
// versions:
// - protoc-gen-go-grpc v1.3.0
// - protoc             v4.25.3
// source: mangle.proto

package proto

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

const (
	Mangle_Query_FullMethodName = "/mangle.Mangle/Query"
)

// MangleClient is the client API for Mangle service.
//
// For semantics around ctx use and closing/ending streaming RPCs, please refer to https://pkg.go.dev/google.golang.org/grpc/?tab=doc#ClientConn.NewStream.
type MangleClient interface {
	// The server answers a query with a stream of responses.
	// It is possible that the list of results is empty.
	// In case of errors, no answers are sent and a QueryError
	// message is included in status response metadata.
	Query(ctx context.Context, in *QueryRequest, opts ...grpc.CallOption) (Mangle_QueryClient, error)
}

type mangleClient struct {
	cc grpc.ClientConnInterface
}

func NewMangleClient(cc grpc.ClientConnInterface) MangleClient {
	return &mangleClient{cc}
}

func (c *mangleClient) Query(ctx context.Context, in *QueryRequest, opts ...grpc.CallOption) (Mangle_QueryClient, error) {
	stream, err := c.cc.NewStream(ctx, &Mangle_ServiceDesc.Streams[0], Mangle_Query_FullMethodName, opts...)
	if err != nil {
		return nil, err
	}
	x := &mangleQueryClient{stream}
	if err := x.ClientStream.SendMsg(in); err != nil {
		return nil, err
	}
	if err := x.ClientStream.CloseSend(); err != nil {
		return nil, err
	}
	return x, nil
}

type Mangle_QueryClient interface {
	Recv() (*QueryAnswer, error)
	grpc.ClientStream
}

type mangleQueryClient struct {
	grpc.ClientStream
}

func (x *mangleQueryClient) Recv() (*QueryAnswer, error) {
	m := new(QueryAnswer)
	if err := x.ClientStream.RecvMsg(m); err != nil {
		return nil, err
	}
	return m, nil
}

// MangleServer is the server API for Mangle service.
// All implementations must embed UnimplementedMangleServer
// for forward compatibility
type MangleServer interface {
	// The server answers a query with a stream of responses.
	// It is possible that the list of results is empty.
	// In case of errors, no answers are sent and a QueryError
	// message is included in status response metadata.
	Query(*QueryRequest, Mangle_QueryServer) error
	mustEmbedUnimplementedMangleServer()
}

// UnimplementedMangleServer must be embedded to have forward compatible implementations.
type UnimplementedMangleServer struct {
}

func (UnimplementedMangleServer) Query(*QueryRequest, Mangle_QueryServer) error {
	return status.Errorf(codes.Unimplemented, "method Query not implemented")
}
func (UnimplementedMangleServer) mustEmbedUnimplementedMangleServer() {}

// UnsafeMangleServer may be embedded to opt out of forward compatibility for this service.
// Use of this interface is not recommended, as added methods to MangleServer will
// result in compilation errors.
type UnsafeMangleServer interface {
	mustEmbedUnimplementedMangleServer()
}

func RegisterMangleServer(s grpc.ServiceRegistrar, srv MangleServer) {
	s.RegisterService(&Mangle_ServiceDesc, srv)
}

func _Mangle_Query_Handler(srv interface{}, stream grpc.ServerStream) error {
	m := new(QueryRequest)
	if err := stream.RecvMsg(m); err != nil {
		return err
	}
	return srv.(MangleServer).Query(m, &mangleQueryServer{stream})
}

type Mangle_QueryServer interface {
	Send(*QueryAnswer) error
	grpc.ServerStream
}

type mangleQueryServer struct {
	grpc.ServerStream
}

func (x *mangleQueryServer) Send(m *QueryAnswer) error {
	return x.ServerStream.SendMsg(m)
}

// Mangle_ServiceDesc is the grpc.ServiceDesc for Mangle service.
// It's only intended for direct use with grpc.RegisterService,
// and not to be introspected or modified (even as a copy)
var Mangle_ServiceDesc = grpc.ServiceDesc{
	ServiceName: "mangle.Mangle",
	HandlerType: (*MangleServer)(nil),
	Methods:     []grpc.MethodDesc{},
	Streams: []grpc.StreamDesc{
		{
			StreamName:    "Query",
			Handler:       _Mangle_Query_Handler,
			ServerStreams: true,
		},
	},
	Metadata: "mangle.proto",
}
