package service

import (
	"context"
	"database/sql"
	"testing"
	"time"

	"github.com/artnikel/TradingService/internal/config"
	"github.com/artnikel/TradingService/internal/model"
	"github.com/artnikel/TradingService/internal/service/mocks"
	"github.com/google/uuid"
	"github.com/lib/pq"
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
	testStrategy = "long"
	testDeal     = &model.Deal{
		DealID:        uuid.New(),
		SharesCount:   decimal.NewFromFloat(1.5),
		ProfileID:     testBalance.ProfileID,
		Company:       "Apple",
		PurchasePrice: decimal.NewFromFloat(1350),
		StopLoss:      decimal.NewFromFloat(100),
		TakeProfit:    decimal.NewFromFloat(1000),
		DealTime:      time.Now().UTC(),
	}
	testProfit = decimal.NewFromFloat(150)
	cfg        config.Variables
)

func TestBalanceOperation(t *testing.T) {
	rep := new(mocks.BalanceRepository)
	srv := NewTradingService(nil, rep, cfg)
	rep.On("BalanceOperation", mock.Anything, mock.AnythingOfType("*model.Balance")).Return(testBalance.Operation.InexactFloat64(), nil).Once()
	_, err := srv.BalanceOperation(context.Background(), testBalance)
	require.NoError(t, err)
	rep.AssertExpectations(t)
}

func TestGetBalanceAndOperation(t *testing.T) {
	rep := new(mocks.BalanceRepository)
	srv := NewTradingService(nil, rep, cfg)
	rep.On("GetBalance", mock.Anything, mock.AnythingOfType("uuid.UUID")).Return(testBalance.Operation.InexactFloat64(), nil).Once()
	money, err := srv.GetBalance(context.Background(), testBalance.ProfileID)
	require.NoError(t, err)
	require.Equal(t, money, testBalance.Operation.InexactFloat64())
	rep.AssertExpectations(t)
}

func TestCreatePosition(t *testing.T) {
	tdeal := &model.Deal{
		Company:     "Apple",
		StopLoss:    decimal.NewFromFloat(100),
		TakeProfit:  decimal.NewFromFloat(1000),
		ProfileID:   uuid.New(),
		SharesCount: decimal.NewFromFloat(1),
	}
	brep := new(mocks.BalanceRepository)
	trep := new(mocks.PriceRepository)
	srv := NewTradingService(trep, brep, cfg)
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	db, err := sql.Open("postgres", cfg.PostgresConnTrading+" sslmode=disable")
	require.NoError(t, err)
	defer db.Close()

	tx, err := db.Begin()
	require.NoError(t, err)
	brep.On("GetBalance", mock.Anything, mock.AnythingOfType("uuid.UUID")).Return(testBalance.Operation.InexactFloat64(), nil).Once()
	trep.On("CreatePosition", mock.Anything, mock.AnythingOfType("string"), mock.AnythingOfType("*model.Deal")).Return(nil).Once()
	trep.On("Subscribe", mock.Anything, mock.AnythingOfType("chan model.Share")).Return(nil)
	brep.On("BalanceOperation", mock.Anything, mock.AnythingOfType("*model.Balance")).Return(testBalance.Operation.InexactFloat64(), nil).Once()
	reportProblem := func(ev pq.ListenerEventType, err error) {
		if err != nil {
			t.Errorf("error: %v", err.Error())
		}
	}
	listener := pq.NewListener(cfg.PostgresConnTrading+"?sslmode=disable", 5*time.Second, time.Minute, reportProblem)

	go func() {
		err := listener.Listen("events")
		if err != nil {
			t.Errorf("Failed to listen events: %s", err)
		}
		srv.Subscribe(ctx)
		srv.WaitForNotification(ctx, listener)
	}()
	err = srv.CreatePosition(ctx, testStrategy, tdeal)
	require.NoError(t, err)

	err = tx.Rollback() 
	require.NoError(t, err)

	brep.AssertExpectations(t)
	trep.AssertExpectations(t)
}

func TestGetProfit(t *testing.T) {
	brep := new(mocks.BalanceRepository)
	trep := new(mocks.PriceRepository)
	srv := NewTradingService(trep, brep, cfg)
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	trep.On("GetPositionInfoByDealID", mock.Anything, mock.AnythingOfType("uuid.UUID")).Return(*testDeal, nil).Once()
	brep.On("GetBalance", mock.Anything, mock.AnythingOfType("uuid.UUID")).Return(testBalance.Operation.InexactFloat64(), nil).Once()
	brep.On("BalanceOperation", mock.Anything, mock.AnythingOfType("*model.Balance")).Return(testBalance.Operation.InexactFloat64(), nil).Once()
	trep.On("ClosePosition", mock.Anything, mock.AnythingOfType("uuid.UUID"), mock.AnythingOfType("uuid.UUID")).Return(testProfit, nil).Once()
	trep.On("AddPosition", mock.Anything, mock.AnythingOfType("string"), mock.AnythingOfType("*model.Deal")).Return(nil).Once()
	err := srv.GetProfit(ctx, testStrategy, testDeal)
	require.NoError(t, err)
	profit, err := srv.ClosePosition(ctx, testDeal.DealID, testBalance.ProfileID)
	require.NoError(t, err)
	if !profit.Equal(testProfit) {
		t.Errorf("Expected profit %v, but got %v", testProfit, profit)
	}
	trep.AssertExpectations(t)
	brep.AssertExpectations(t)
}

func TestClosePosition(t *testing.T) {
	brep := new(mocks.BalanceRepository)
	trep := new(mocks.PriceRepository)
	srv := NewTradingService(trep, brep, cfg)
	trep.On("GetPositionInfoByDealID", mock.Anything, mock.AnythingOfType("uuid.UUID")).Return(*testDeal, nil).Once()
	trep.On("ClosePosition", mock.Anything, mock.AnythingOfType("uuid.UUID"), mock.AnythingOfType("uuid.UUID")).Return(testProfit, nil).Once()
	brep.On("BalanceOperation", mock.Anything, mock.AnythingOfType("*model.Balance")).Return(testBalance.Operation.InexactFloat64(), nil).Once()
	profit, err := srv.ClosePosition(context.Background(), testDeal.DealID, testBalance.ProfileID)
	require.NoError(t, err)
	require.Equal(t, profit, testProfit)
	trep.AssertExpectations(t)
	brep.AssertExpectations(t)
}

func TestGetUnclosedPositions(t *testing.T) {
	trep := new(mocks.PriceRepository)
	srv := NewTradingService(trep, nil, cfg)
	var sliceDeals []*model.Deal
	sliceDeals = append(sliceDeals, testDeal)
	trep.On("GetUnclosedPositions", mock.Anything, mock.AnythingOfType("uuid.UUID")).Return(sliceDeals, nil).Once()
	unclosedDeals, err := srv.GetUnclosedPositions(context.Background(), testBalance.ProfileID)
	require.NoError(t, err)
	require.Equal(t, len(unclosedDeals), len(sliceDeals))
	trep.AssertExpectations(t)
}

func TestGetPrices(t *testing.T) {
	trep := new(mocks.PriceRepository)
	srv := NewTradingService(trep, nil, cfg)
	prices, err := srv.GetPrices()
	require.NoError(t, err)
	require.Empty(t, prices)
	trep.AssertExpectations(t)
}
