package service

import (
	"context"
	"log"
	"os"
	"sync"
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
	testDeal = &model.Deal{
		DealID:        uuid.New(),
		SharesCount:   decimal.NewFromFloat(1),
		ProfileID:     testBalance.ProfileID,
		Company:       "Apple",
		PurchasePrice: decimal.NewFromFloat(200),
		StopLoss:      decimal.NewFromFloat(100),
		TakeProfit:    decimal.NewFromFloat(1000),
		DealTime:      time.Now().UTC(),
	}
	testBalanceSecond = &model.Balance{
		BalanceID: uuid.New(),
		ProfileID: uuid.New(),
		Operation: decimal.NewFromFloat(10000),
	}
	testDealSecond = &model.Deal{
		DealID:        uuid.New(),
		SharesCount:   decimal.NewFromFloat(1),
		ProfileID:     testBalanceSecond.ProfileID,
		Company:       "Xerox",
		PurchasePrice: decimal.NewFromFloat(2000),
		StopLoss:      decimal.NewFromFloat(1000),
		TakeProfit:    decimal.NewFromFloat(3000),
		DealTime:      time.Now().UTC(),
	}
	testShare = &model.Share{
		Company: "Logitech",
		Price:   decimal.NewFromFloat(333),
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
	srv := NewTradingService(trep, nil, cfg)
	strategy, err := srv.identifyStrategy(testDeal.TakeProfit, testDeal.StopLoss)
	require.NoError(t, err)
	require.Equal(t, strategy, long)
	trep.AssertExpectations(t)
}

func TestWaitForNotification(t *testing.T) {
	trep := new(mocks.PriceRepository)
	srv := NewTradingService(trep, nil, cfg)
	l := &pq.Listener{
		Notify: make(chan *pq.Notification, 1),
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
	srv := NewTradingService(trep, brep, cfg)
	trep.On("CreatePosition", mock.Anything, mock.AnythingOfType("*model.Deal")).Return(nil).Once()
	brep.On("GetBalance", mock.Anything, mock.AnythingOfType("uuid.UUID")).Return(testBalance.Operation.InexactFloat64(), nil).Once()
	brep.On("BalanceOperation", mock.Anything, mock.AnythingOfType("*model.Balance")).Return(testBalance.Operation.InexactFloat64(), nil).Once()
	srv.manager.Mu.Lock()
	srv.manager.PricesMap[testDeal.Company] = decimal.NewFromInt(150)
	srv.manager.Mu.Unlock()
	err := srv.CreatePosition(context.Background(), testDeal)
	require.NoError(t, err)
	trep.AssertExpectations(t)
	brep.AssertExpectations(t)
}

func TestCreateTwoPositions(t *testing.T) {
	brep := new(mocks.BalanceRepository)
	trep := new(mocks.PriceRepository)
	srv := NewTradingService(trep, brep, cfg)
	trep.On("CreatePosition", mock.Anything, mock.AnythingOfType("*model.Deal")).Return(nil)
	brep.On("GetBalance", mock.Anything, mock.AnythingOfType("uuid.UUID")).Return(testBalance.Operation.InexactFloat64(), nil)
	brep.On("BalanceOperation", mock.Anything, mock.AnythingOfType("*model.Balance")).Return(testBalance.Operation.InexactFloat64(), nil)
	var wg sync.WaitGroup
	wg.Add(2)
	go func() {
		defer wg.Done()
		srv.manager.Mu.Lock()
		srv.manager.Positions[testDeal.ProfileID] = make(map[uuid.UUID]model.Deal)
		srv.manager.PricesMap[testDeal.Company] = decimal.NewFromInt(250)
		srv.manager.Mu.Unlock()
		err := srv.CreatePosition(context.Background(), testDeal)
		require.NoError(t, err)
	}()
	go func() {
		defer wg.Done()
		srv.manager.Mu.Lock()
		srv.manager.Positions[testDealSecond.ProfileID] = make(map[uuid.UUID]model.Deal)
		srv.manager.PricesMap[testDealSecond.Company] = decimal.NewFromInt(1050)
		srv.manager.Mu.Unlock()
		err := srv.CreatePosition(context.Background(), testDealSecond)
		require.NoError(t, err)
	}()
	wg.Wait()
	trep.AssertExpectations(t)
	brep.AssertExpectations(t)
}

func TestGetProfit(t *testing.T) {
	trep := new(mocks.PriceRepository)
	brep := new(mocks.BalanceRepository)
	srv := NewTradingService(trep, brep, cfg)
	trep.On("ClosePosition", mock.Anything, mock.AnythingOfType("*model.Deal")).Return(nil).Once()
	trep.On("CreatePosition", mock.Anything, mock.AnythingOfType("*model.Deal")).Return(nil).Once()
	brep.On("GetBalance", mock.Anything, mock.AnythingOfType("uuid.UUID")).Return(testBalance.Operation.InexactFloat64(), nil).Once()
	brep.On("BalanceOperation", mock.Anything, mock.AnythingOfType("*model.Balance")).Return(testBalance.Operation.InexactFloat64(), nil)
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*5)
	defer cancel()
	srv.manager.Mu.Lock()
	srv.manager.PricesMap[testDeal.Company] = decimal.NewFromInt(950)
	srv.manager.Mu.Unlock()
	err := srv.CreatePosition(ctx, testDeal)
	require.NoError(t, err)
	go func() {
		srv.manager.Mu.Lock()
		srv.manager.Positions[testDeal.ProfileID][testDeal.DealID] = *testDeal
		srv.manager.PricesMap[testDeal.Company] = decimal.NewFromInt(1050)
		srv.manager.Mu.Unlock()
	}()
	srv.getProfit(ctx, *testDeal)
	trep.AssertExpectations(t)
	brep.AssertExpectations(t)
}

func TestGetProfitForTwoUsers(t *testing.T) {
	brep := new(mocks.BalanceRepository)
	trep := new(mocks.PriceRepository)
	srv := NewTradingService(trep, brep, cfg)
	trep.On("ClosePosition", mock.Anything, mock.AnythingOfType("*model.Deal")).Return(nil)
	trep.On("CreatePosition", mock.Anything, mock.AnythingOfType("*model.Deal")).Return(nil)
	brep.On("GetBalance", mock.Anything, mock.AnythingOfType("uuid.UUID")).Return(testBalance.Operation.InexactFloat64(), nil)
	brep.On("BalanceOperation", mock.Anything, mock.AnythingOfType("*model.Balance")).Return(testBalance.Operation.InexactFloat64(), nil)
	ctx1, cancel1 := context.WithTimeout(context.Background(), time.Second*5)
	defer cancel1()
	ctx2, cancel2 := context.WithTimeout(context.Background(), time.Second*5)
	defer cancel2()
	var wg sync.WaitGroup
	wg.Add(2)
	go func() {
		defer wg.Done()
		srv.manager.Mu.Lock()
		srv.manager.Positions[testDeal.ProfileID] = make(map[uuid.UUID]model.Deal)
		srv.manager.PricesMap[testDeal.Company] = decimal.NewFromInt(250)
		srv.manager.Mu.Unlock()
		err := srv.CreatePosition(ctx1, testDeal)
		require.NoError(t, err)
	}()
	go func() {
		defer wg.Done()
		srv.manager.Mu.Lock()
		srv.manager.Positions[testDealSecond.ProfileID] = make(map[uuid.UUID]model.Deal)
		srv.manager.PricesMap[testDealSecond.Company] = decimal.NewFromInt(1050)
		srv.manager.Mu.Unlock()
		err := srv.CreatePosition(ctx2, testDealSecond)
		require.NoError(t, err)
	}()
	wg.Wait()
	wg.Add(2)
	go func() {
		defer wg.Done()
		srv.manager.Mu.Lock()
		srv.manager.Positions[testDeal.ProfileID][testDeal.DealID] = *testDeal
		srv.manager.PricesMap[testDeal.Company] = decimal.NewFromInt(1050)
		srv.manager.Mu.Unlock()
		srv.getProfit(ctx1, *testDeal)
	}()
	go func() {
		defer wg.Done()
		srv.manager.Mu.Lock()
		srv.manager.Positions[testDealSecond.ProfileID][testDealSecond.DealID] = *testDealSecond
		srv.manager.PricesMap[testDealSecond.Company] = decimal.NewFromInt(950)
		srv.manager.Mu.Unlock()
		srv.getProfit(ctx2, *testDealSecond)
	}()
	wg.Wait()
	trep.AssertExpectations(t)
	brep.AssertExpectations(t)
}

func TestClosePosition(t *testing.T) {
	brep := new(mocks.BalanceRepository)
	trep := new(mocks.PriceRepository)
	srv := NewTradingService(trep, brep, cfg)
	brep.On("BalanceOperation", mock.Anything, mock.AnythingOfType("*model.Balance")).Return(testBalance.Operation.InexactFloat64(), nil).Once()
	trep.On("ClosePosition", mock.Anything, mock.AnythingOfType("*model.Deal")).Return(nil).Once()
	srv.manager.Mu.Lock()
	srv.manager.Positions[testDeal.ProfileID] = make(map[uuid.UUID]model.Deal)
	srv.manager.Positions[testDeal.ProfileID][testDeal.DealID] = *testDeal
	srv.manager.Mu.Unlock()
	sharePrice := decimal.NewFromFloat(300)
	profit, err := srv.closePosition(context.Background(), testDeal.DealID, testDeal.ProfileID, sharePrice)
	require.NoError(t, err)
	require.Equal(t, profit.InexactFloat64(), sharePrice.Sub(testDeal.PurchasePrice).InexactFloat64())
	trep.AssertExpectations(t)
	brep.AssertExpectations(t)
}

func TestCloseTwoPositions(t *testing.T) {
	brep := new(mocks.BalanceRepository)
	trep := new(mocks.PriceRepository)
	srv := NewTradingService(trep, brep, cfg)
	brep.On("BalanceOperation", mock.Anything, mock.AnythingOfType("*model.Balance")).Return(testBalance.Operation.InexactFloat64(), nil)
	trep.On("ClosePosition", mock.Anything, mock.AnythingOfType("*model.Deal")).Return(nil)
	srv.manager.Mu.Lock()
	srv.manager.Positions[testDeal.ProfileID] = make(map[uuid.UUID]model.Deal)
	srv.manager.Positions[testDeal.ProfileID][testDeal.DealID] = *testDeal
	srv.manager.Positions[testDealSecond.ProfileID] = make(map[uuid.UUID]model.Deal)
	srv.manager.Positions[testDealSecond.ProfileID][testDealSecond.DealID] = *testDealSecond
	srv.manager.Mu.Unlock()
	var (
		wg               sync.WaitGroup
		sharePriceFirst  = decimal.NewFromFloat(300)
		sharePriceSecond = decimal.NewFromFloat(200)
	)
	wg.Add(2)
	go func() {
		defer wg.Done()
		profit, err := srv.closePosition(context.Background(), testDeal.DealID, testDeal.ProfileID, sharePriceFirst)
		require.NoError(t, err)
		require.Equal(t, profit.InexactFloat64(), sharePriceFirst.Sub(testDeal.PurchasePrice).InexactFloat64())
	}()
	go func() {
		defer wg.Done()
		profit, err := srv.closePosition(context.Background(), testDealSecond.DealID, testDealSecond.ProfileID, sharePriceSecond)
		require.NoError(t, err)
		require.Equal(t, profit.InexactFloat64(), sharePriceSecond.Sub(testDealSecond.PurchasePrice).InexactFloat64())
	}()
	wg.Wait()
	trep.AssertExpectations(t)
	brep.AssertExpectations(t)
}

func TestClosePositionManually(t *testing.T) {
	brep := new(mocks.BalanceRepository)
	trep := new(mocks.PriceRepository)
	srv := NewTradingService(trep, brep, cfg)
	testBalanceForClose := &model.Balance{
		BalanceID: uuid.New(),
		ProfileID: uuid.New(),
		Operation: decimal.NewFromFloat(10000),
	}
	testDealForClose := &model.Deal{
		DealID:        uuid.New(),
		SharesCount:   decimal.NewFromFloat(1),
		ProfileID:     testBalanceForClose.ProfileID,
		Company:       "Xerox",
		PurchasePrice: decimal.NewFromFloat(1000),
		StopLoss:      decimal.NewFromFloat(500),
		TakeProfit:    decimal.NewFromFloat(3000),
		DealTime:      time.Now().UTC(),
	}
	trep.On("CreatePosition", mock.Anything, mock.AnythingOfType("*model.Deal")).Return(nil).Once()
	brep.On("GetBalance", mock.Anything, mock.AnythingOfType("uuid.UUID")).Return(testBalance.Operation.InexactFloat64(), nil).Once()
	brep.On("BalanceOperation", mock.Anything, mock.AnythingOfType("*model.Balance")).Return(testBalance.Operation.InexactFloat64(), nil).Once()
	trep.On("ClosePosition", mock.Anything, mock.AnythingOfType("*model.Deal")).Return(nil).Once()
	brep.On("BalanceOperation", mock.Anything, mock.AnythingOfType("*model.Balance")).Return(testBalance.Operation.InexactFloat64(), nil).Once()
	srv.manager.Mu.Lock()
	srv.manager.Positions[testDealForClose.ProfileID] = make(map[uuid.UUID]model.Deal)
	srv.manager.Positions[testDealForClose.ProfileID][testDealForClose.DealID] = *testDealForClose
	srv.manager.PricesMap[testDealForClose.Company] = decimal.NewFromInt(1100)
	srv.manager.Mu.Unlock()
	err := srv.CreatePosition(context.Background(), testDealForClose)
	require.NoError(t, err)
	profit, err := srv.ClosePositionManually(context.Background(), testDealForClose.DealID, testDealForClose.ProfileID)
	require.NoError(t, err)
	require.Equal(t, profit.InexactFloat64(), testProfit.InexactFloat64())
	trep.AssertExpectations(t)
	brep.AssertExpectations(t)
}

func TestSubscribe(t *testing.T) {
	trep := new(mocks.PriceRepository)
	srv := NewTradingService(trep, nil, cfg)
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*5)
	defer cancel()
	submanager := make(chan model.Share, 1)
	submanager <- *testShare
	go func() {
		err := trep.Subscribe(ctx, submanager)
		require.NoError(t, err)
	}()
	srv.manager.Mu.Lock()
	for i := 0; i < cap(submanager); i++ {
		share := <-submanager
		srv.manager.PricesMap[share.Company] = share.Price
	}
	srv.manager.Mu.Unlock()
	require.NotEmpty(t, srv.manager.PricesMap)
	require.Equal(t, testShare.Price, srv.manager.PricesMap[testShare.Company])
	trep.AssertExpectations(t)
}

func TestGetPrices(t *testing.T) {
	trep := new(mocks.PriceRepository)
	srv := NewTradingService(trep, nil, cfg)
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

func TestGetClosedPositions(t *testing.T) {
	trep := new(mocks.PriceRepository)
	srv := NewTradingService(trep, nil, cfg)
	var sliceDeals []*model.Deal
	sliceDeals = append(sliceDeals, testDeal)
	trep.On("GetClosedPositions", mock.Anything, mock.AnythingOfType("uuid.UUID")).Return(sliceDeals, nil).Once()
	closedDeals, err := srv.GetClosedPositions(context.Background(), testBalance.ProfileID)
	require.NoError(t, err)
	require.Equal(t, len(closedDeals), len(sliceDeals))
	trep.AssertExpectations(t)
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

func TestBackupUnclosedPositions(t *testing.T) {
	trep := new(mocks.PriceRepository)
	srv := NewTradingService(trep, nil, cfg)
	var sliceDeals []*model.Deal
	sliceDeals = append(sliceDeals, testDeal)
	trep.On("GetUnclosedPositionsForAll", mock.Anything).Return(sliceDeals, nil).Once()
	srv.BackupUnclosedPositions(context.Background())
	trep.AssertExpectations(t)
}

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
