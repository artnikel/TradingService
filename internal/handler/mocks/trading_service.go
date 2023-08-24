// Code generated by mockery v2.18.0. DO NOT EDIT.

package mocks

import (
	context "context"

	decimal "github.com/shopspring/decimal"

	mock "github.com/stretchr/testify/mock"

	model "github.com/artnikel/TradingService/internal/model"

	uuid "github.com/google/uuid"
)

// TradingService is an autogenerated mock type for the TradingService type
type TradingService struct {
	mock.Mock
}

// BalanceOperation provides a mock function with given fields: ctx, balance
func (_m *TradingService) BalanceOperation(ctx context.Context, balance *model.Balance) (float64, error) {
	ret := _m.Called(ctx, balance)

	var r0 float64
	if rf, ok := ret.Get(0).(func(context.Context, *model.Balance) float64); ok {
		r0 = rf(ctx, balance)
	} else {
		r0 = ret.Get(0).(float64)
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(context.Context, *model.Balance) error); ok {
		r1 = rf(ctx, balance)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// ClosePosition provides a mock function with given fields: ctx, dealid, profileid
func (_m *TradingService) ClosePosition(ctx context.Context, dealid uuid.UUID, profileid uuid.UUID) (decimal.Decimal, error) {
	ret := _m.Called(ctx, dealid, profileid)

	var r0 decimal.Decimal
	if rf, ok := ret.Get(0).(func(context.Context, uuid.UUID, uuid.UUID) decimal.Decimal); ok {
		r0 = rf(ctx, dealid, profileid)
	} else {
		r0 = ret.Get(0).(decimal.Decimal)
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(context.Context, uuid.UUID, uuid.UUID) error); ok {
		r1 = rf(ctx, dealid, profileid)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// GetBalance provides a mock function with given fields: ctx, profileid
func (_m *TradingService) GetBalance(ctx context.Context, profileid uuid.UUID) (float64, error) {
	ret := _m.Called(ctx, profileid)

	var r0 float64
	if rf, ok := ret.Get(0).(func(context.Context, uuid.UUID) float64); ok {
		r0 = rf(ctx, profileid)
	} else {
		r0 = ret.Get(0).(float64)
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(context.Context, uuid.UUID) error); ok {
		r1 = rf(ctx, profileid)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// GetPrices provides a mock function with given fields: ctx
func (_m *TradingService) GetPrices(ctx context.Context) ([]model.Share, error) {
	ret := _m.Called(ctx)

	var r0 []model.Share
	if rf, ok := ret.Get(0).(func(context.Context) []model.Share); ok {
		r0 = rf(ctx)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).([]model.Share)
		}
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(context.Context) error); ok {
		r1 = rf(ctx)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// GetProfit provides a mock function with given fields: ctx, strategy, deal
func (_m *TradingService) GetProfit(ctx context.Context, strategy string, deal *model.Deal) (decimal.Decimal, error) {
	ret := _m.Called(ctx, strategy, deal)

	var r0 decimal.Decimal
	if rf, ok := ret.Get(0).(func(context.Context, string, *model.Deal) decimal.Decimal); ok {
		r0 = rf(ctx, strategy, deal)
	} else {
		r0 = ret.Get(0).(decimal.Decimal)
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(context.Context, string, *model.Deal) error); ok {
		r1 = rf(ctx, strategy, deal)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// GetUnclosedPositions provides a mock function with given fields: ctx, profileid
func (_m *TradingService) GetUnclosedPositions(ctx context.Context, profileid uuid.UUID) ([]*model.Deal, error) {
	ret := _m.Called(ctx, profileid)

	var r0 []*model.Deal
	if rf, ok := ret.Get(0).(func(context.Context, uuid.UUID) []*model.Deal); ok {
		r0 = rf(ctx, profileid)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).([]*model.Deal)
		}
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(context.Context, uuid.UUID) error); ok {
		r1 = rf(ctx, profileid)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

type mockConstructorTestingTNewTradingService interface {
	mock.TestingT
	Cleanup(func())
}

// NewTradingService creates a new instance of TradingService. It also registers a testing interface on the mock and a cleanup function to assert the mocks expectations.
func NewTradingService(t mockConstructorTestingTNewTradingService) *TradingService {
	mock := &TradingService{}
	mock.Mock.Test(t)

	t.Cleanup(func() { mock.AssertExpectations(t) })

	return mock
}
