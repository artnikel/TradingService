package service

import (
	"context"
	"testing"
	"time"

	"github.com/artnikel/TradingService/internal/model"
	"github.com/artnikel/TradingService/internal/service/mocks"
	"github.com/google/uuid"
	"github.com/shopspring/decimal"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
)

var (
	testBalance = &model.Balance{
		BalanceID: uuid.New(),
		ProfileID: uuid.New(),
		Operation: decimal.NewFromFloat(100.9),
	}
	testStrategy = "short"
	testDeal     = &model.Deal{
		DealID:       uuid.New(),
		ActionsCount: decimal.NewFromFloat(1.5),
		ProfileID:    testBalance.ProfileID,
		Company:      "Apple",
		StopLoss:     decimal.NewFromFloat(1500),
		TakeProfit:   decimal.NewFromFloat(1000),
		DealTime:     time.Now().UTC(),
	}
)

func TestBalanceOperation(t *testing.T) {
	rep := new(mocks.BalanceRepository)
	srv := NewTradingService(nil, rep)
	rep.On("BalanceOperation", mock.Anything, mock.AnythingOfType("*model.Balance")).Return(testBalance.Operation.InexactFloat64(), nil).Once()
	_, err := srv.BalanceOperation(context.Background(), testBalance)
	require.NoError(t, err)
	rep.AssertExpectations(t)
}

func TestGetBalanceAndOperation(t *testing.T) {
	rep := new(mocks.BalanceRepository)
	srv := NewTradingService(nil, rep)
	rep.On("GetBalance", mock.Anything, mock.AnythingOfType("uuid.UUID")).Return(testBalance.Operation.InexactFloat64(), nil).Once()
	money, err := srv.GetBalance(context.Background(), testBalance.ProfileID)
	require.NoError(t, err)
	require.Equal(t, money, testBalance.Operation.InexactFloat64())
	rep.AssertExpectations(t)
}

func TestStrategies(t *testing.T) {
	prep := new(mocks.PriceRepository)
	brep := new(mocks.BalanceRepository)
	srv := NewTradingService(prep, brep)
	brep.On("BalanceOperation", mock.Anything, mock.AnythingOfType("*model.Balance")).Return(testBalance.Operation.InexactFloat64(), nil).Once()
	prep.On("Subscribe", mock.Anything, mock.AnythingOfType("*model.Deal"), mock.AnythingOfType("chan []*model.Action")).Return(nil).Once()
	brep.On("GetBalance", mock.Anything, mock.AnythingOfType("uuid.UUID")).Return(testBalance.Operation.InexactFloat64(), nil).Once()
	prep.On("AddDeal", mock.Anything, mock.AnythingOfType("string"), mock.AnythingOfType("*model.Deal")).
		Return(nil).Once()
	_, err := srv.Strategies(context.Background(), testStrategy, testDeal)
	require.NoError(t, err)
	brep.AssertExpectations(t)
	prep.AssertExpectations(t)
}
