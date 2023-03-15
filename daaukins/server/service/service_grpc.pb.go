// Code generated by protoc-gen-go-grpc. DO NOT EDIT.
// versions:
// - protoc-gen-go-grpc v1.2.0
// - protoc             v3.21.12
// source: service/service.proto

package service

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

// ServiceClient is the client API for Service service.
//
// For semantics around ctx use and closing/ending streaming RPCs, please refer to https://pkg.go.dev/google.golang.org/grpc/?tab=doc#ClientConn.NewStream.
type ServiceClient interface {
	HaveCapacity(ctx context.Context, in *HaveCapacityRequest, opts ...grpc.CallOption) (*HaveCapacityResponse, error)
	ScheduleLab(ctx context.Context, in *ScheduleLabRequest, opts ...grpc.CallOption) (*ScheduleLabResponse, error)
	GetLab(ctx context.Context, in *GetLabRequest, opts ...grpc.CallOption) (*GetLabResponse, error)
	GetLabs(ctx context.Context, in *GetLabsRequest, opts ...grpc.CallOption) (*GetLabsResponse, error)
	RemoveLab(ctx context.Context, in *RemoveLabRequest, opts ...grpc.CallOption) (*RemoveLabResponse, error)
}

type serviceClient struct {
	cc grpc.ClientConnInterface
}

func NewServiceClient(cc grpc.ClientConnInterface) ServiceClient {
	return &serviceClient{cc}
}

func (c *serviceClient) HaveCapacity(ctx context.Context, in *HaveCapacityRequest, opts ...grpc.CallOption) (*HaveCapacityResponse, error) {
	out := new(HaveCapacityResponse)
	err := c.cc.Invoke(ctx, "/service.service/HaveCapacity", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *serviceClient) ScheduleLab(ctx context.Context, in *ScheduleLabRequest, opts ...grpc.CallOption) (*ScheduleLabResponse, error) {
	out := new(ScheduleLabResponse)
	err := c.cc.Invoke(ctx, "/service.service/ScheduleLab", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *serviceClient) GetLab(ctx context.Context, in *GetLabRequest, opts ...grpc.CallOption) (*GetLabResponse, error) {
	out := new(GetLabResponse)
	err := c.cc.Invoke(ctx, "/service.service/GetLab", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *serviceClient) GetLabs(ctx context.Context, in *GetLabsRequest, opts ...grpc.CallOption) (*GetLabsResponse, error) {
	out := new(GetLabsResponse)
	err := c.cc.Invoke(ctx, "/service.service/GetLabs", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *serviceClient) RemoveLab(ctx context.Context, in *RemoveLabRequest, opts ...grpc.CallOption) (*RemoveLabResponse, error) {
	out := new(RemoveLabResponse)
	err := c.cc.Invoke(ctx, "/service.service/RemoveLab", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

// ServiceServer is the server API for Service service.
// All implementations must embed UnimplementedServiceServer
// for forward compatibility
type ServiceServer interface {
	HaveCapacity(context.Context, *HaveCapacityRequest) (*HaveCapacityResponse, error)
	ScheduleLab(context.Context, *ScheduleLabRequest) (*ScheduleLabResponse, error)
	GetLab(context.Context, *GetLabRequest) (*GetLabResponse, error)
	GetLabs(context.Context, *GetLabsRequest) (*GetLabsResponse, error)
	RemoveLab(context.Context, *RemoveLabRequest) (*RemoveLabResponse, error)
	mustEmbedUnimplementedServiceServer()
}

// UnimplementedServiceServer must be embedded to have forward compatible implementations.
type UnimplementedServiceServer struct {
}

func (UnimplementedServiceServer) HaveCapacity(context.Context, *HaveCapacityRequest) (*HaveCapacityResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method HaveCapacity not implemented")
}
func (UnimplementedServiceServer) ScheduleLab(context.Context, *ScheduleLabRequest) (*ScheduleLabResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method ScheduleLab not implemented")
}
func (UnimplementedServiceServer) GetLab(context.Context, *GetLabRequest) (*GetLabResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method GetLab not implemented")
}
func (UnimplementedServiceServer) GetLabs(context.Context, *GetLabsRequest) (*GetLabsResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method GetLabs not implemented")
}
func (UnimplementedServiceServer) RemoveLab(context.Context, *RemoveLabRequest) (*RemoveLabResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method RemoveLab not implemented")
}
func (UnimplementedServiceServer) mustEmbedUnimplementedServiceServer() {}

// UnsafeServiceServer may be embedded to opt out of forward compatibility for this service.
// Use of this interface is not recommended, as added methods to ServiceServer will
// result in compilation errors.
type UnsafeServiceServer interface {
	mustEmbedUnimplementedServiceServer()
}

func RegisterServiceServer(s grpc.ServiceRegistrar, srv ServiceServer) {
	s.RegisterService(&Service_ServiceDesc, srv)
}

func _Service_HaveCapacity_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(HaveCapacityRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(ServiceServer).HaveCapacity(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/service.service/HaveCapacity",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(ServiceServer).HaveCapacity(ctx, req.(*HaveCapacityRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _Service_ScheduleLab_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(ScheduleLabRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(ServiceServer).ScheduleLab(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/service.service/ScheduleLab",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(ServiceServer).ScheduleLab(ctx, req.(*ScheduleLabRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _Service_GetLab_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(GetLabRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(ServiceServer).GetLab(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/service.service/GetLab",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(ServiceServer).GetLab(ctx, req.(*GetLabRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _Service_GetLabs_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(GetLabsRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(ServiceServer).GetLabs(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/service.service/GetLabs",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(ServiceServer).GetLabs(ctx, req.(*GetLabsRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _Service_RemoveLab_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(RemoveLabRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(ServiceServer).RemoveLab(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/service.service/RemoveLab",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(ServiceServer).RemoveLab(ctx, req.(*RemoveLabRequest))
	}
	return interceptor(ctx, in, info, handler)
}

// Service_ServiceDesc is the grpc.ServiceDesc for Service service.
// It's only intended for direct use with grpc.RegisterService,
// and not to be introspected or modified (even as a copy)
var Service_ServiceDesc = grpc.ServiceDesc{
	ServiceName: "service.service",
	HandlerType: (*ServiceServer)(nil),
	Methods: []grpc.MethodDesc{
		{
			MethodName: "HaveCapacity",
			Handler:    _Service_HaveCapacity_Handler,
		},
		{
			MethodName: "ScheduleLab",
			Handler:    _Service_ScheduleLab_Handler,
		},
		{
			MethodName: "GetLab",
			Handler:    _Service_GetLab_Handler,
		},
		{
			MethodName: "GetLabs",
			Handler:    _Service_GetLabs_Handler,
		},
		{
			MethodName: "RemoveLab",
			Handler:    _Service_RemoveLab_Handler,
		},
	},
	Streams:  []grpc.StreamDesc{},
	Metadata: "service/service.proto",
}
