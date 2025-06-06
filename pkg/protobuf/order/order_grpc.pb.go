// Code generated by protoc-gen-go-grpc. DO NOT EDIT.
// versions:
// - protoc-gen-go-grpc v1.5.1
// - protoc             (unknown)
// source: protobuf/order/order.proto

package order

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
	BinanceOrderService_PlaceMarketOrder_FullMethodName     = "/order.BinanceOrderService/PlaceMarketOrder"
	BinanceOrderService_PlaceLimitOrder_FullMethodName      = "/order.BinanceOrderService/PlaceLimitOrder"
	BinanceOrderService_PlaceStopMarketOrder_FullMethodName = "/order.BinanceOrderService/PlaceStopMarketOrder"
)

// BinanceOrderServiceClient is the client API for BinanceOrderService service.
//
// For semantics around ctx use and closing/ending streaming RPCs, please refer to https://pkg.go.dev/google.golang.org/grpc/?tab=doc#ClientConn.NewStream.
//
// Spot market orders (e.g. BTC/USDT)
type BinanceOrderServiceClient interface {
	PlaceMarketOrder(ctx context.Context, in *MarketOrderRequest, opts ...grpc.CallOption) (*OrderResponse, error)
	PlaceLimitOrder(ctx context.Context, in *LimitOrderRequest, opts ...grpc.CallOption) (*OrderResponse, error)
	PlaceStopMarketOrder(ctx context.Context, in *StopMarketOrderRequest, opts ...grpc.CallOption) (*OrderResponse, error)
}

type binanceOrderServiceClient struct {
	cc grpc.ClientConnInterface
}

func NewBinanceOrderServiceClient(cc grpc.ClientConnInterface) BinanceOrderServiceClient {
	return &binanceOrderServiceClient{cc}
}

func (c *binanceOrderServiceClient) PlaceMarketOrder(ctx context.Context, in *MarketOrderRequest, opts ...grpc.CallOption) (*OrderResponse, error) {
	cOpts := append([]grpc.CallOption{grpc.StaticMethod()}, opts...)
	out := new(OrderResponse)
	err := c.cc.Invoke(ctx, BinanceOrderService_PlaceMarketOrder_FullMethodName, in, out, cOpts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *binanceOrderServiceClient) PlaceLimitOrder(ctx context.Context, in *LimitOrderRequest, opts ...grpc.CallOption) (*OrderResponse, error) {
	cOpts := append([]grpc.CallOption{grpc.StaticMethod()}, opts...)
	out := new(OrderResponse)
	err := c.cc.Invoke(ctx, BinanceOrderService_PlaceLimitOrder_FullMethodName, in, out, cOpts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *binanceOrderServiceClient) PlaceStopMarketOrder(ctx context.Context, in *StopMarketOrderRequest, opts ...grpc.CallOption) (*OrderResponse, error) {
	cOpts := append([]grpc.CallOption{grpc.StaticMethod()}, opts...)
	out := new(OrderResponse)
	err := c.cc.Invoke(ctx, BinanceOrderService_PlaceStopMarketOrder_FullMethodName, in, out, cOpts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

// BinanceOrderServiceServer is the server API for BinanceOrderService service.
// All implementations must embed UnimplementedBinanceOrderServiceServer
// for forward compatibility.
//
// Spot market orders (e.g. BTC/USDT)
type BinanceOrderServiceServer interface {
	PlaceMarketOrder(context.Context, *MarketOrderRequest) (*OrderResponse, error)
	PlaceLimitOrder(context.Context, *LimitOrderRequest) (*OrderResponse, error)
	PlaceStopMarketOrder(context.Context, *StopMarketOrderRequest) (*OrderResponse, error)
	mustEmbedUnimplementedBinanceOrderServiceServer()
}

// UnimplementedBinanceOrderServiceServer must be embedded to have
// forward compatible implementations.
//
// NOTE: this should be embedded by value instead of pointer to avoid a nil
// pointer dereference when methods are called.
type UnimplementedBinanceOrderServiceServer struct{}

func (UnimplementedBinanceOrderServiceServer) PlaceMarketOrder(context.Context, *MarketOrderRequest) (*OrderResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method PlaceMarketOrder not implemented")
}
func (UnimplementedBinanceOrderServiceServer) PlaceLimitOrder(context.Context, *LimitOrderRequest) (*OrderResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method PlaceLimitOrder not implemented")
}
func (UnimplementedBinanceOrderServiceServer) PlaceStopMarketOrder(context.Context, *StopMarketOrderRequest) (*OrderResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method PlaceStopMarketOrder not implemented")
}
func (UnimplementedBinanceOrderServiceServer) mustEmbedUnimplementedBinanceOrderServiceServer() {}
func (UnimplementedBinanceOrderServiceServer) testEmbeddedByValue()                             {}

// UnsafeBinanceOrderServiceServer may be embedded to opt out of forward compatibility for this service.
// Use of this interface is not recommended, as added methods to BinanceOrderServiceServer will
// result in compilation errors.
type UnsafeBinanceOrderServiceServer interface {
	mustEmbedUnimplementedBinanceOrderServiceServer()
}

func RegisterBinanceOrderServiceServer(s grpc.ServiceRegistrar, srv BinanceOrderServiceServer) {
	// If the following call pancis, it indicates UnimplementedBinanceOrderServiceServer was
	// embedded by pointer and is nil.  This will cause panics if an
	// unimplemented method is ever invoked, so we test this at initialization
	// time to prevent it from happening at runtime later due to I/O.
	if t, ok := srv.(interface{ testEmbeddedByValue() }); ok {
		t.testEmbeddedByValue()
	}
	s.RegisterService(&BinanceOrderService_ServiceDesc, srv)
}

func _BinanceOrderService_PlaceMarketOrder_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(MarketOrderRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(BinanceOrderServiceServer).PlaceMarketOrder(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: BinanceOrderService_PlaceMarketOrder_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(BinanceOrderServiceServer).PlaceMarketOrder(ctx, req.(*MarketOrderRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _BinanceOrderService_PlaceLimitOrder_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(LimitOrderRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(BinanceOrderServiceServer).PlaceLimitOrder(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: BinanceOrderService_PlaceLimitOrder_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(BinanceOrderServiceServer).PlaceLimitOrder(ctx, req.(*LimitOrderRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _BinanceOrderService_PlaceStopMarketOrder_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(StopMarketOrderRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(BinanceOrderServiceServer).PlaceStopMarketOrder(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: BinanceOrderService_PlaceStopMarketOrder_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(BinanceOrderServiceServer).PlaceStopMarketOrder(ctx, req.(*StopMarketOrderRequest))
	}
	return interceptor(ctx, in, info, handler)
}

// BinanceOrderService_ServiceDesc is the grpc.ServiceDesc for BinanceOrderService service.
// It's only intended for direct use with grpc.RegisterService,
// and not to be introspected or modified (even as a copy)
var BinanceOrderService_ServiceDesc = grpc.ServiceDesc{
	ServiceName: "order.BinanceOrderService",
	HandlerType: (*BinanceOrderServiceServer)(nil),
	Methods: []grpc.MethodDesc{
		{
			MethodName: "PlaceMarketOrder",
			Handler:    _BinanceOrderService_PlaceMarketOrder_Handler,
		},
		{
			MethodName: "PlaceLimitOrder",
			Handler:    _BinanceOrderService_PlaceLimitOrder_Handler,
		},
		{
			MethodName: "PlaceStopMarketOrder",
			Handler:    _BinanceOrderService_PlaceStopMarketOrder_Handler,
		},
	},
	Streams:  []grpc.StreamDesc{},
	Metadata: "protobuf/order/order.proto",
}

const (
	BinancePerpOrderService_PlaceMarketOrder_FullMethodName     = "/order.BinancePerpOrderService/PlaceMarketOrder"
	BinancePerpOrderService_PlaceLimitOrder_FullMethodName      = "/order.BinancePerpOrderService/PlaceLimitOrder"
	BinancePerpOrderService_PlaceStopMarketOrder_FullMethodName = "/order.BinancePerpOrderService/PlaceStopMarketOrder"
)

// BinancePerpOrderServiceClient is the client API for BinancePerpOrderService service.
//
// For semantics around ctx use and closing/ending streaming RPCs, please refer to https://pkg.go.dev/google.golang.org/grpc/?tab=doc#ClientConn.NewStream.
//
// Perpetual futures orders (e.g. BTCUSDT‑PERP)
type BinancePerpOrderServiceClient interface {
	PlaceMarketOrder(ctx context.Context, in *MarketOrderRequest, opts ...grpc.CallOption) (*OrderResponse, error)
	PlaceLimitOrder(ctx context.Context, in *LimitOrderRequest, opts ...grpc.CallOption) (*OrderResponse, error)
	PlaceStopMarketOrder(ctx context.Context, in *StopMarketOrderRequest, opts ...grpc.CallOption) (*OrderResponse, error)
}

type binancePerpOrderServiceClient struct {
	cc grpc.ClientConnInterface
}

func NewBinancePerpOrderServiceClient(cc grpc.ClientConnInterface) BinancePerpOrderServiceClient {
	return &binancePerpOrderServiceClient{cc}
}

func (c *binancePerpOrderServiceClient) PlaceMarketOrder(ctx context.Context, in *MarketOrderRequest, opts ...grpc.CallOption) (*OrderResponse, error) {
	cOpts := append([]grpc.CallOption{grpc.StaticMethod()}, opts...)
	out := new(OrderResponse)
	err := c.cc.Invoke(ctx, BinancePerpOrderService_PlaceMarketOrder_FullMethodName, in, out, cOpts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *binancePerpOrderServiceClient) PlaceLimitOrder(ctx context.Context, in *LimitOrderRequest, opts ...grpc.CallOption) (*OrderResponse, error) {
	cOpts := append([]grpc.CallOption{grpc.StaticMethod()}, opts...)
	out := new(OrderResponse)
	err := c.cc.Invoke(ctx, BinancePerpOrderService_PlaceLimitOrder_FullMethodName, in, out, cOpts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *binancePerpOrderServiceClient) PlaceStopMarketOrder(ctx context.Context, in *StopMarketOrderRequest, opts ...grpc.CallOption) (*OrderResponse, error) {
	cOpts := append([]grpc.CallOption{grpc.StaticMethod()}, opts...)
	out := new(OrderResponse)
	err := c.cc.Invoke(ctx, BinancePerpOrderService_PlaceStopMarketOrder_FullMethodName, in, out, cOpts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

// BinancePerpOrderServiceServer is the server API for BinancePerpOrderService service.
// All implementations must embed UnimplementedBinancePerpOrderServiceServer
// for forward compatibility.
//
// Perpetual futures orders (e.g. BTCUSDT‑PERP)
type BinancePerpOrderServiceServer interface {
	PlaceMarketOrder(context.Context, *MarketOrderRequest) (*OrderResponse, error)
	PlaceLimitOrder(context.Context, *LimitOrderRequest) (*OrderResponse, error)
	PlaceStopMarketOrder(context.Context, *StopMarketOrderRequest) (*OrderResponse, error)
	mustEmbedUnimplementedBinancePerpOrderServiceServer()
}

// UnimplementedBinancePerpOrderServiceServer must be embedded to have
// forward compatible implementations.
//
// NOTE: this should be embedded by value instead of pointer to avoid a nil
// pointer dereference when methods are called.
type UnimplementedBinancePerpOrderServiceServer struct{}

func (UnimplementedBinancePerpOrderServiceServer) PlaceMarketOrder(context.Context, *MarketOrderRequest) (*OrderResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method PlaceMarketOrder not implemented")
}
func (UnimplementedBinancePerpOrderServiceServer) PlaceLimitOrder(context.Context, *LimitOrderRequest) (*OrderResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method PlaceLimitOrder not implemented")
}
func (UnimplementedBinancePerpOrderServiceServer) PlaceStopMarketOrder(context.Context, *StopMarketOrderRequest) (*OrderResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method PlaceStopMarketOrder not implemented")
}
func (UnimplementedBinancePerpOrderServiceServer) mustEmbedUnimplementedBinancePerpOrderServiceServer() {
}
func (UnimplementedBinancePerpOrderServiceServer) testEmbeddedByValue() {}

// UnsafeBinancePerpOrderServiceServer may be embedded to opt out of forward compatibility for this service.
// Use of this interface is not recommended, as added methods to BinancePerpOrderServiceServer will
// result in compilation errors.
type UnsafeBinancePerpOrderServiceServer interface {
	mustEmbedUnimplementedBinancePerpOrderServiceServer()
}

func RegisterBinancePerpOrderServiceServer(s grpc.ServiceRegistrar, srv BinancePerpOrderServiceServer) {
	// If the following call pancis, it indicates UnimplementedBinancePerpOrderServiceServer was
	// embedded by pointer and is nil.  This will cause panics if an
	// unimplemented method is ever invoked, so we test this at initialization
	// time to prevent it from happening at runtime later due to I/O.
	if t, ok := srv.(interface{ testEmbeddedByValue() }); ok {
		t.testEmbeddedByValue()
	}
	s.RegisterService(&BinancePerpOrderService_ServiceDesc, srv)
}

func _BinancePerpOrderService_PlaceMarketOrder_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(MarketOrderRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(BinancePerpOrderServiceServer).PlaceMarketOrder(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: BinancePerpOrderService_PlaceMarketOrder_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(BinancePerpOrderServiceServer).PlaceMarketOrder(ctx, req.(*MarketOrderRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _BinancePerpOrderService_PlaceLimitOrder_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(LimitOrderRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(BinancePerpOrderServiceServer).PlaceLimitOrder(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: BinancePerpOrderService_PlaceLimitOrder_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(BinancePerpOrderServiceServer).PlaceLimitOrder(ctx, req.(*LimitOrderRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _BinancePerpOrderService_PlaceStopMarketOrder_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(StopMarketOrderRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(BinancePerpOrderServiceServer).PlaceStopMarketOrder(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: BinancePerpOrderService_PlaceStopMarketOrder_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(BinancePerpOrderServiceServer).PlaceStopMarketOrder(ctx, req.(*StopMarketOrderRequest))
	}
	return interceptor(ctx, in, info, handler)
}

// BinancePerpOrderService_ServiceDesc is the grpc.ServiceDesc for BinancePerpOrderService service.
// It's only intended for direct use with grpc.RegisterService,
// and not to be introspected or modified (even as a copy)
var BinancePerpOrderService_ServiceDesc = grpc.ServiceDesc{
	ServiceName: "order.BinancePerpOrderService",
	HandlerType: (*BinancePerpOrderServiceServer)(nil),
	Methods: []grpc.MethodDesc{
		{
			MethodName: "PlaceMarketOrder",
			Handler:    _BinancePerpOrderService_PlaceMarketOrder_Handler,
		},
		{
			MethodName: "PlaceLimitOrder",
			Handler:    _BinancePerpOrderService_PlaceLimitOrder_Handler,
		},
		{
			MethodName: "PlaceStopMarketOrder",
			Handler:    _BinancePerpOrderService_PlaceStopMarketOrder_Handler,
		},
	},
	Streams:  []grpc.StreamDesc{},
	Metadata: "protobuf/order/order.proto",
}
