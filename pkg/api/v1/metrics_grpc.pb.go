// Code generated by protoc-gen-go-grpc. DO NOT EDIT.
// versions:
// - protoc-gen-go-grpc v1.3.0
// - protoc             v4.24.4
// source: metrics.proto

package api

import (
	context "context"
	status "google.golang.org/genproto/googleapis/rpc/status"
	grpc "google.golang.org/grpc"
	codes "google.golang.org/grpc/codes"
	status1 "google.golang.org/grpc/status"
)

// This is a compile-time assertion to ensure that this generated file
// is compatible with the grpc package it is being compiled against.
// Requires gRPC-Go v1.32.0 or later.
const _ = grpc.SupportPackageIsVersion7

const (
	MetricsServiceV1_MetricsBatchV1_FullMethodName = "/api.v1.MetricsServiceV1/MetricsBatchV1"
)

// MetricsServiceV1Client is the client API for MetricsServiceV1 service.
//
// For semantics around ctx use and closing/ending streaming RPCs, please refer to https://pkg.go.dev/google.golang.org/grpc/?tab=doc#ClientConn.NewStream.
type MetricsServiceV1Client interface {
	MetricsBatchV1(ctx context.Context, in *MetricsBatchRequest, opts ...grpc.CallOption) (*status.Status, error)
}

type metricsServiceV1Client struct {
	cc grpc.ClientConnInterface
}

func NewMetricsServiceV1Client(cc grpc.ClientConnInterface) MetricsServiceV1Client {
	return &metricsServiceV1Client{cc}
}

func (c *metricsServiceV1Client) MetricsBatchV1(ctx context.Context, in *MetricsBatchRequest, opts ...grpc.CallOption) (*status.Status, error) {
	out := new(status.Status)
	err := c.cc.Invoke(ctx, MetricsServiceV1_MetricsBatchV1_FullMethodName, in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

// MetricsServiceV1Server is the server API for MetricsServiceV1 service.
// All implementations must embed UnimplementedMetricsServiceV1Server
// for forward compatibility
type MetricsServiceV1Server interface {
	MetricsBatchV1(context.Context, *MetricsBatchRequest) (*status.Status, error)
	mustEmbedUnimplementedMetricsServiceV1Server()
}

// UnimplementedMetricsServiceV1Server must be embedded to have forward compatible implementations.
type UnimplementedMetricsServiceV1Server struct {
}

func (UnimplementedMetricsServiceV1Server) MetricsBatchV1(context.Context, *MetricsBatchRequest) (*status.Status, error) {
	return nil, status1.Errorf(codes.Unimplemented, "method MetricsBatchV1 not implemented")
}
func (UnimplementedMetricsServiceV1Server) mustEmbedUnimplementedMetricsServiceV1Server() {}

// UnsafeMetricsServiceV1Server may be embedded to opt out of forward compatibility for this service.
// Use of this interface is not recommended, as added methods to MetricsServiceV1Server will
// result in compilation errors.
type UnsafeMetricsServiceV1Server interface {
	mustEmbedUnimplementedMetricsServiceV1Server()
}

func RegisterMetricsServiceV1Server(s grpc.ServiceRegistrar, srv MetricsServiceV1Server) {
	s.RegisterService(&MetricsServiceV1_ServiceDesc, srv)
}

func _MetricsServiceV1_MetricsBatchV1_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(MetricsBatchRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(MetricsServiceV1Server).MetricsBatchV1(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: MetricsServiceV1_MetricsBatchV1_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(MetricsServiceV1Server).MetricsBatchV1(ctx, req.(*MetricsBatchRequest))
	}
	return interceptor(ctx, in, info, handler)
}

// MetricsServiceV1_ServiceDesc is the grpc.ServiceDesc for MetricsServiceV1 service.
// It's only intended for direct use with grpc.RegisterService,
// and not to be introspected or modified (even as a copy)
var MetricsServiceV1_ServiceDesc = grpc.ServiceDesc{
	ServiceName: "api.v1.MetricsServiceV1",
	HandlerType: (*MetricsServiceV1Server)(nil),
	Methods: []grpc.MethodDesc{
		{
			MethodName: "MetricsBatchV1",
			Handler:    _MetricsServiceV1_MetricsBatchV1_Handler,
		},
	},
	Streams:  []grpc.StreamDesc{},
	Metadata: "metrics.proto",
}