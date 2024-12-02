// Code generated by protoc-gen-go-grpc. DO NOT EDIT.
// versions:
// - protoc-gen-go-grpc v1.5.1
// - protoc             v5.28.2
// source: api/dump-collector.proto

package main

import (
	context "context"
	grpc "google.golang.org/grpc"
	codes "google.golang.org/grpc/codes"
	status "google.golang.org/grpc/status"
)

// This is a compile-time assertion to ensure that this generated file
// is compatible with the grpc package it is being compiled against.
// Requires gRPC-Go v1.64.0 or later.
const _ = grpc.SupportPackageIsVersion9

const (
	DumpCollectorService_SubscribeToTrace_FullMethodName = "/dumpcollector.DumpCollectorService/SubscribeToTrace"
)

// DumpCollectorServiceClient is the client API for DumpCollectorService service.
//
// For semantics around ctx use and closing/ending streaming RPCs, please refer to https://pkg.go.dev/google.golang.org/grpc/?tab=doc#ClientConn.NewStream.
type DumpCollectorServiceClient interface {
	SubscribeToTrace(ctx context.Context, opts ...grpc.CallOption) (grpc.BidiStreamingClient[Request, TraceNotification], error)
}

type dumpCollectorServiceClient struct {
	cc grpc.ClientConnInterface
}

func NewDumpCollectorServiceClient(cc grpc.ClientConnInterface) DumpCollectorServiceClient {
	return &dumpCollectorServiceClient{cc}
}

func (c *dumpCollectorServiceClient) SubscribeToTrace(ctx context.Context, opts ...grpc.CallOption) (grpc.BidiStreamingClient[Request, TraceNotification], error) {
	cOpts := append([]grpc.CallOption{grpc.StaticMethod()}, opts...)
	stream, err := c.cc.NewStream(ctx, &DumpCollectorService_ServiceDesc.Streams[0], DumpCollectorService_SubscribeToTrace_FullMethodName, cOpts...)
	if err != nil {
		return nil, err
	}
	x := &grpc.GenericClientStream[Request, TraceNotification]{ClientStream: stream}
	return x, nil
}

// This type alias is provided for backwards compatibility with existing code that references the prior non-generic stream type by name.
type DumpCollectorService_SubscribeToTraceClient = grpc.BidiStreamingClient[Request, TraceNotification]

// DumpCollectorServiceServer is the server API for DumpCollectorService service.
// All implementations must embed UnimplementedDumpCollectorServiceServer
// for forward compatibility.
type DumpCollectorServiceServer interface {
	SubscribeToTrace(grpc.BidiStreamingServer[Request, TraceNotification]) error
	mustEmbedUnimplementedDumpCollectorServiceServer()
}

// UnimplementedDumpCollectorServiceServer must be embedded to have
// forward compatible implementations.
//
// NOTE: this should be embedded by value instead of pointer to avoid a nil
// pointer dereference when methods are called.
type UnimplementedDumpCollectorServiceServer struct{}

func (UnimplementedDumpCollectorServiceServer) SubscribeToTrace(grpc.BidiStreamingServer[Request, TraceNotification]) error {
	return status.Errorf(codes.Unimplemented, "method SubscribeToTrace not implemented")
}
func (UnimplementedDumpCollectorServiceServer) mustEmbedUnimplementedDumpCollectorServiceServer() {}
func (UnimplementedDumpCollectorServiceServer) testEmbeddedByValue()                              {}

// UnsafeDumpCollectorServiceServer may be embedded to opt out of forward compatibility for this service.
// Use of this interface is not recommended, as added methods to DumpCollectorServiceServer will
// result in compilation errors.
type UnsafeDumpCollectorServiceServer interface {
	mustEmbedUnimplementedDumpCollectorServiceServer()
}

func RegisterDumpCollectorServiceServer(s grpc.ServiceRegistrar, srv DumpCollectorServiceServer) {
	// If the following call pancis, it indicates UnimplementedDumpCollectorServiceServer was
	// embedded by pointer and is nil.  This will cause panics if an
	// unimplemented method is ever invoked, so we test this at initialization
	// time to prevent it from happening at runtime later due to I/O.
	if t, ok := srv.(interface{ testEmbeddedByValue() }); ok {
		t.testEmbeddedByValue()
	}
	s.RegisterService(&DumpCollectorService_ServiceDesc, srv)
}

func _DumpCollectorService_SubscribeToTrace_Handler(srv interface{}, stream grpc.ServerStream) error {
	return srv.(DumpCollectorServiceServer).SubscribeToTrace(&grpc.GenericServerStream[Request, TraceNotification]{ServerStream: stream})
}

// This type alias is provided for backwards compatibility with existing code that references the prior non-generic stream type by name.
type DumpCollectorService_SubscribeToTraceServer = grpc.BidiStreamingServer[Request, TraceNotification]

// DumpCollectorService_ServiceDesc is the grpc.ServiceDesc for DumpCollectorService service.
// It's only intended for direct use with grpc.RegisterService,
// and not to be introspected or modified (even as a copy)
var DumpCollectorService_ServiceDesc = grpc.ServiceDesc{
	ServiceName: "dumpcollector.DumpCollectorService",
	HandlerType: (*DumpCollectorServiceServer)(nil),
	Methods:     []grpc.MethodDesc{},
	Streams: []grpc.StreamDesc{
		{
			StreamName:    "SubscribeToTrace",
			Handler:       _DumpCollectorService_SubscribeToTrace_Handler,
			ServerStreams: true,
			ClientStreams: true,
		},
	},
	Metadata: "api/dump-collector.proto",
}
