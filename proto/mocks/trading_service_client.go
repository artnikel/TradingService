// Code generated by mockery v2.18.0. DO NOT EDIT.

package mocks

import (
	context "context"

	grpc "google.golang.org/grpc"

	mock "github.com/stretchr/testify/mock"

	proto "github.com/artnikel/TradingService/proto"
)

// TradingServiceClient is an autogenerated mock type for the TradingServiceClient type
type TradingServiceClient struct {
	mock.Mock
}

// ClosePositionManually provides a mock function with given fields: ctx, in, opts
func (_m *TradingServiceClient) ClosePositionManually(ctx context.Context, in *proto.ClosePositionManuallyRequest, opts ...grpc.CallOption) (*proto.ClosePositionManuallyResponse, error) {
	_va := make([]interface{}, len(opts))
	for _i := range opts {
		_va[_i] = opts[_i]
	}
	var _ca []interface{}
	_ca = append(_ca, ctx, in)
	_ca = append(_ca, _va...)
	ret := _m.Called(_ca...)

	var r0 *proto.ClosePositionManuallyResponse
	if rf, ok := ret.Get(0).(func(context.Context, *proto.ClosePositionManuallyRequest, ...grpc.CallOption) *proto.ClosePositionManuallyResponse); ok {
		r0 = rf(ctx, in, opts...)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*proto.ClosePositionManuallyResponse)
		}
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(context.Context, *proto.ClosePositionManuallyRequest, ...grpc.CallOption) error); ok {
		r1 = rf(ctx, in, opts...)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// CreatePosition provides a mock function with given fields: ctx, in, opts
func (_m *TradingServiceClient) CreatePosition(ctx context.Context, in *proto.CreatePositionRequest, opts ...grpc.CallOption) (*proto.CreatePositionResponse, error) {
	_va := make([]interface{}, len(opts))
	for _i := range opts {
		_va[_i] = opts[_i]
	}
	var _ca []interface{}
	_ca = append(_ca, ctx, in)
	_ca = append(_ca, _va...)
	ret := _m.Called(_ca...)

	var r0 *proto.CreatePositionResponse
	if rf, ok := ret.Get(0).(func(context.Context, *proto.CreatePositionRequest, ...grpc.CallOption) *proto.CreatePositionResponse); ok {
		r0 = rf(ctx, in, opts...)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*proto.CreatePositionResponse)
		}
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(context.Context, *proto.CreatePositionRequest, ...grpc.CallOption) error); ok {
		r1 = rf(ctx, in, opts...)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// GetClosedPositions provides a mock function with given fields: ctx, in, opts
func (_m *TradingServiceClient) GetClosedPositions(ctx context.Context, in *proto.GetClosedPositionsRequest, opts ...grpc.CallOption) (*proto.GetClosedPositionsResponse, error) {
	_va := make([]interface{}, len(opts))
	for _i := range opts {
		_va[_i] = opts[_i]
	}
	var _ca []interface{}
	_ca = append(_ca, ctx, in)
	_ca = append(_ca, _va...)
	ret := _m.Called(_ca...)

	var r0 *proto.GetClosedPositionsResponse
	if rf, ok := ret.Get(0).(func(context.Context, *proto.GetClosedPositionsRequest, ...grpc.CallOption) *proto.GetClosedPositionsResponse); ok {
		r0 = rf(ctx, in, opts...)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*proto.GetClosedPositionsResponse)
		}
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(context.Context, *proto.GetClosedPositionsRequest, ...grpc.CallOption) error); ok {
		r1 = rf(ctx, in, opts...)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// GetPrices provides a mock function with given fields: ctx, in, opts
func (_m *TradingServiceClient) GetPrices(ctx context.Context, in *proto.GetPricesRequest, opts ...grpc.CallOption) (*proto.GetPricesResponse, error) {
	_va := make([]interface{}, len(opts))
	for _i := range opts {
		_va[_i] = opts[_i]
	}
	var _ca []interface{}
	_ca = append(_ca, ctx, in)
	_ca = append(_ca, _va...)
	ret := _m.Called(_ca...)

	var r0 *proto.GetPricesResponse
	if rf, ok := ret.Get(0).(func(context.Context, *proto.GetPricesRequest, ...grpc.CallOption) *proto.GetPricesResponse); ok {
		r0 = rf(ctx, in, opts...)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*proto.GetPricesResponse)
		}
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(context.Context, *proto.GetPricesRequest, ...grpc.CallOption) error); ok {
		r1 = rf(ctx, in, opts...)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// GetUnclosedPositions provides a mock function with given fields: ctx, in, opts
func (_m *TradingServiceClient) GetUnclosedPositions(ctx context.Context, in *proto.GetUnclosedPositionsRequest, opts ...grpc.CallOption) (*proto.GetUnclosedPositionsResponse, error) {
	_va := make([]interface{}, len(opts))
	for _i := range opts {
		_va[_i] = opts[_i]
	}
	var _ca []interface{}
	_ca = append(_ca, ctx, in)
	_ca = append(_ca, _va...)
	ret := _m.Called(_ca...)

	var r0 *proto.GetUnclosedPositionsResponse
	if rf, ok := ret.Get(0).(func(context.Context, *proto.GetUnclosedPositionsRequest, ...grpc.CallOption) *proto.GetUnclosedPositionsResponse); ok {
		r0 = rf(ctx, in, opts...)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*proto.GetUnclosedPositionsResponse)
		}
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(context.Context, *proto.GetUnclosedPositionsRequest, ...grpc.CallOption) error); ok {
		r1 = rf(ctx, in, opts...)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

type mockConstructorTestingTNewTradingServiceClient interface {
	mock.TestingT
	Cleanup(func())
}

// NewTradingServiceClient creates a new instance of TradingServiceClient. It also registers a testing interface on the mock and a cleanup function to assert the mocks expectations.
func NewTradingServiceClient(t mockConstructorTestingTNewTradingServiceClient) *TradingServiceClient {
	mock := &TradingServiceClient{}
	mock.Mock.Test(t)

	t.Cleanup(func() { mock.AssertExpectations(t) })

	return mock
}
