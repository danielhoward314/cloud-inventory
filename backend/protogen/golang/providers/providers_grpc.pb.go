// Code generated by protoc-gen-go-grpc. DO NOT EDIT.

package providers

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

// ProvidersServiceClient is the client API for ProvidersService service.
//
// For semantics around ctx use and closing/ending streaming RPCs, please refer to https://pkg.go.dev/google.golang.org/grpc/?tab=doc#ClientConn.NewStream.
type ProvidersServiceClient interface {
	List(ctx context.Context, in *ListRequest, opts ...grpc.CallOption) (*ListResponse, error)
}

type providersServiceClient struct {
	cc grpc.ClientConnInterface
}

func NewProvidersServiceClient(cc grpc.ClientConnInterface) ProvidersServiceClient {
	return &providersServiceClient{cc}
}

func (c *providersServiceClient) List(ctx context.Context, in *ListRequest, opts ...grpc.CallOption) (*ListResponse, error) {
	out := new(ListResponse)
	err := c.cc.Invoke(ctx, "/providers.ProvidersService/List", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

// ProvidersServiceServer is the server API for ProvidersService service.
// All implementations must embed UnimplementedProvidersServiceServer
// for forward compatibility
type ProvidersServiceServer interface {
	List(context.Context, *ListRequest) (*ListResponse, error)
	mustEmbedUnimplementedProvidersServiceServer()
}

// UnimplementedProvidersServiceServer must be embedded to have forward compatible implementations.
type UnimplementedProvidersServiceServer struct {
}

func (UnimplementedProvidersServiceServer) List(context.Context, *ListRequest) (*ListResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method List not implemented")
}
func (UnimplementedProvidersServiceServer) mustEmbedUnimplementedProvidersServiceServer() {}

// UnsafeProvidersServiceServer may be embedded to opt out of forward compatibility for this service.
// Use of this interface is not recommended, as added methods to ProvidersServiceServer will
// result in compilation errors.
type UnsafeProvidersServiceServer interface {
	mustEmbedUnimplementedProvidersServiceServer()
}

func RegisterProvidersServiceServer(s grpc.ServiceRegistrar, srv ProvidersServiceServer) {
	s.RegisterService(&ProvidersService_ServiceDesc, srv)
}

func _ProvidersService_List_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(ListRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(ProvidersServiceServer).List(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/providers.ProvidersService/List",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(ProvidersServiceServer).List(ctx, req.(*ListRequest))
	}
	return interceptor(ctx, in, info, handler)
}

// ProvidersService_ServiceDesc is the grpc.ServiceDesc for ProvidersService service.
// It's only intended for direct use with grpc.RegisterService,
// and not to be introspected or modified (even as a copy)
var ProvidersService_ServiceDesc = grpc.ServiceDesc{
	ServiceName: "providers.ProvidersService",
	HandlerType: (*ProvidersServiceServer)(nil),
	Methods: []grpc.MethodDesc{
		{
			MethodName: "List",
			Handler:    _ProvidersService_List_Handler,
		},
	},
	Streams:  []grpc.StreamDesc{},
	Metadata: "providers/providers.proto",
}
