package service

import (
	"context"
	"log"
	"os"
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
		Operation: decimal.NewFromFloat(2000),
	}
	testStrategy = "long"
	testDeal     = &model.Deal{
		DealID:        uuid.New(),
		SharesCount:   decimal.NewFromFloat(1),
		ProfileID:     testBalance.ProfileID,
		Company:       "Apple",
		PurchasePrice: decimal.NewFromFloat(200),
		StopLoss:      decimal.NewFromFloat(100),
		TakeProfit:    decimal.NewFromFloat(1000),
		DealTime:      time.Now().UTC(),
	}
	testProfit = decimal.NewFromFloat(100)
	cfg        *config.Variables
)

func TestMain(m *testing.M) {
	var err error
	cfg, err = config.New()
	if err != nil {
		log.Fatalf("Could not parse config: err: %v", err)
	}
	exitVal := m.Run()
	os.Exit(exitVal)
}

func TestIdentifyStrategy(t *testing.T) {
	trep := new(mocks.PriceRepository)
	srv := NewTradingService(trep, nil, *cfg)
	strategy := srv.identifyStrategy(testDeal.TakeProfit, testDeal.StopLoss)
	require.Equal(t, strategy, testStrategy, "long")
	trep.AssertExpectations(t)
}

func TestWaitForNotification(t *testing.T) {
	trep := new(mocks.PriceRepository)
	srv := NewTradingService(trep, nil, *cfg)
	l := &pq.Listener{
		Notify: make(chan *pq.Notification),
	}
	defer close(l.Notify)
	notificationData := `{
		"action": "INSERT",
		"data": {
			"DealID": "` + testDeal.DealID.String() + `",
			"SharesCount": "` + testDeal.PurchasePrice.String() + `",
			"ProfileID": "` + testDeal.ProfileID.String() + `",
			"Company": "` + testDeal.Company + `",
			"PurchasePrice": "` + testDeal.PurchasePrice.String() + `",
			"TakeProfit": "` + testDeal.TakeProfit.String() + `",
			"StopLoss": "` + testDeal.StopLoss.String() + `"
		}
	}`
	l.Notify <- &pq.Notification{
		Extra: notificationData,
	}
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	go func() {
		srv.WaitForNotification(ctx, l)
	}()
	trep.AssertExpectations(t)
}

func TestCreatePosition(t *testing.T) {
	brep := new(mocks.BalanceRepository)
	trep := new(mocks.PriceRepository)
	srv := NewTradingService(trep, brep, *cfg)
	trep.On("CreatePosition", mock.Anything, mock.AnythingOfType("string"), mock.AnythingOfType("*model.Deal")).Return(nil).Once()
	brep.On("GetBalance", mock.Anything, mock.AnythingOfType("uuid.UUID")).Return(testBalance.Operation.InexactFloat64(), nil).Once()
	brep.On("BalanceOperation", mock.Anything, mock.AnythingOfType("*model.Balance")).Return(testBalance.Operation.InexactFloat64(), nil).Once()
	go func() {
		srv.manager.SubscribersShares[testDeal.ProfileID][testDeal.Company] <- model.Share{
			Company: testDeal.Company,
			Price:   (decimal.NewFromInt(250))}
	}()
	err := srv.CreatePosition(context.Background(), testDeal)
	require.NoError(t, err)
	trep.AssertExpectations(t)
	brep.AssertExpectations(t)
}

func TestGetProfit(t *testing.T) {
	trep := new(mocks.PriceRepository)
	brep := new(mocks.BalanceRepository)
	srv := NewTradingService(trep, brep, *cfg)
	trep.On("ClosePosition", mock.Anything, mock.AnythingOfType("uuid.UUID"), mock.AnythingOfType("uuid.UUID")).Return(testProfit, nil).Once()
	trep.On("CreatePosition", mock.Anything, mock.AnythingOfType("string"), mock.AnythingOfType("*model.Deal")).Return(nil).Once()
	brep.On("GetBalance", mock.Anything, mock.AnythingOfType("uuid.UUID")).Return(testBalance.Operation.InexactFloat64(), nil).Once()
	brep.On("BalanceOperation", mock.Anything, mock.AnythingOfType("*model.Balance")).Return(testBalance.Operation.InexactFloat64(), nil).Once()
	go func() {
		srv.manager.SubscribersShares[testDeal.ProfileID][testDeal.Company] <- model.Share{
			Company: testDeal.Company,
			Price:   (decimal.NewFromInt(250))}
	}()
	err := srv.CreatePosition(context.Background(), testDeal)
	require.NoError(t, err)
	go func() {
		srv.manager.SubscribersShares[testDeal.ProfileID][testDeal.Company] <- model.Share{
			Company: testDeal.Company,
			Price:   (decimal.NewFromInt(1050))}
	}()
	srv.GetProfit(context.Background(), testDeal)
	require.NotEmpty(t, testDeal.EndDealTime)
	require.NotEmpty(t, testDeal.Profit)
	trep.AssertExpectations(t)
	brep.AssertExpectations(t)
}

func TestClosePosition(t *testing.T) {
	brep := new(mocks.BalanceRepository)
	trep := new(mocks.PriceRepository)
	srv := NewTradingService(trep, brep, *cfg)
	brep.On("BalanceOperation", mock.Anything, mock.AnythingOfType("*model.Balance")).Return(testBalance.Operation.InexactFloat64(), nil).Once()
	trep.On("ClosePosition", mock.Anything, mock.AnythingOfType("*model.Deal")).Return(nil).Once()
	srv.manager.Mu.Lock()
	srv.manager.Positions[testDeal.DealID] = *testDeal
	srv.manager.SubscribersShares[testDeal.ProfileID] = make(map[string]chan model.Share)
	srv.manager.SubscribersShares[testDeal.ProfileID][testDeal.Company] = make(chan model.Share)
	srv.manager.Mu.Unlock()
	sharePrice := decimal.NewFromFloat(300)
	profit, err := srv.ClosePosition(context.Background(), testDeal.DealID, testDeal.ProfileID, sharePrice)
	require.NoError(t, err)
	require.Equal(t, profit.InexactFloat64(), testProfit.InexactFloat64(), sharePrice.Sub(testDeal.PurchasePrice).InexactFloat64())
	trep.AssertExpectations(t)
	brep.AssertExpectations(t)
}

func TestClosePositionManually(t *testing.T) {
	brep := new(mocks.BalanceRepository)
	trep := new(mocks.PriceRepository)
	srv := NewTradingService(trep, brep, *cfg)
	trep.On("ClosePosition", mock.Anything, mock.AnythingOfType("*model.Deal")).Return(nil).Once()
	brep.On("BalanceOperation", mock.Anything, mock.AnythingOfType("*model.Balance")).Return(testBalance.Operation.InexactFloat64(), nil).Once()
	srv.manager.Mu.Lock()
	srv.manager.Positions[testDeal.DealID] = *testDeal
	srv.manager.SubscribersShares[testDeal.ProfileID] = make(map[string]chan model.Share)
	srv.manager.SubscribersShares[testDeal.ProfileID][testDeal.Company] = make(chan model.Share)
	srv.manager.Mu.Unlock()
	go func() {
		srv.manager.SubscribersShares[testDeal.ProfileID][testDeal.Company] <- model.Share{
			Company: testDeal.Company,
			Price:   (decimal.NewFromInt(300))}
	}()
	profit, err := srv.ClosePositionManually(context.Background(), testDeal.DealID, testDeal.ProfileID)
	require.NoError(t, err)
	require.Equal(t, profit.InexactFloat64(), testProfit.InexactFloat64())
	trep.AssertExpectations(t)
	brep.AssertExpectations(t)
}

func TestSubscribe(t *testing.T) {
	trep := new(mocks.PriceRepository)
	srv := NewTradingService(trep, nil, *cfg)
	done := make(chan struct{})
	trep.On("Subscribe", mock.Anything, mock.AnythingOfType("chan model.Share")).Return(nil).Once()
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*1)
	defer cancel()
	go func() {
		srv.Subscribe(ctx)
		done <- struct{}{}
	}()
	<-done
	trep.AssertExpectations(t)
}

func TestGetPrices(t *testing.T) {
	trep := new(mocks.PriceRepository)
	srv := NewTradingService(trep, nil, *cfg)
	srv.manager.Mu.Lock()
	srv.manager.PricesMap[testDeal.Company] = testDeal.PurchasePrice
	srv.manager.Mu.Unlock()
	shares, err := srv.GetPrices()
	require.NoError(t, err)
	for _, share := range shares {
		require.Equal(t, share.Company, testDeal.Company)
		require.Equal(t, share.Price.InexactFloat64(), testDeal.PurchasePrice.InexactFloat64())
	}
	trep.AssertExpectations(t)
}

func TestGetUnclosedPositions(t *testing.T) {
	trep := new(mocks.PriceRepository)
	srv := NewTradingService(trep, nil, *cfg)
	var sliceDeals []*model.Deal
	sliceDeals = append(sliceDeals, testDeal)
	trep.On("GetUnclosedPositions", mock.Anything, mock.AnythingOfType("uuid.UUID")).Return(sliceDeals, nil).Once()
	unclosedDeals, err := srv.GetUnclosedPositions(context.Background(), testBalance.ProfileID)
	require.NoError(t, err)
	require.Equal(t, len(unclosedDeals), len(sliceDeals))
	trep.AssertExpectations(t)
}

func TestBackupUnclosedPositions(t *testing.T) {
	trep := new(mocks.PriceRepository)
	srv := NewTradingService(trep, nil, *cfg)
	var sliceDeals []*model.Deal
	sliceDeals = append(sliceDeals, testDeal)
	trep.On("GetUnclosedPositionsForAll", mock.Anything).Return(sliceDeals, nil).Once()
	srv.BackupUnclosedPositions(context.Background())
	trep.AssertExpectations(t)
}

func TestBalanceOperation(t *testing.T) {
	rep := new(mocks.BalanceRepository)
	srv := NewTradingService(nil, rep, *cfg)
	rep.On("BalanceOperation", mock.Anything, mock.AnythingOfType("*model.Balance")).Return(testBalance.Operation.InexactFloat64(), nil).Once()
	_, err := srv.BalanceOperation(context.Background(), testBalance)
	require.NoError(t, err)
	rep.AssertExpectations(t)
}

func TestGetBalanceAndOperation(t *testing.T) {
	rep := new(mocks.BalanceRepository)
	srv := NewTradingService(nil, rep, *cfg)
	rep.On("GetBalance", mock.Anything, mock.AnythingOfType("uuid.UUID")).Return(testBalance.Operation.InexactFloat64(), nil).Once()
	money, err := srv.GetBalance(context.Background(), testBalance.ProfileID)
	require.NoError(t, err)
	require.Equal(t, money, testBalance.Operation.InexactFloat64())
	rep.AssertExpectations(t)
}
