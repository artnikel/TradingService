package handler

import (
	"context"
	"testing"
	"time"

	"github.com/artnikel/TradingService/internal/handler/mocks"
	"github.com/artnikel/TradingService/internal/model"
	"github.com/artnikel/TradingService/proto"
	"github.com/go-playground/validator/v10"
	"github.com/google/uuid"
	"github.com/shopspring/decimal"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
	"google.golang.org/protobuf/types/known/timestamppb"
)

var (
	v            = validator.New()
	testStrategy = "short"
	testDeal     = &model.Deal{
		DealID:      uuid.New(),
		SharesCount: decimal.NewFromFloat(1.5),
		ProfileID:   uuid.New(),
		Company:     "Apple",
		StopLoss:    decimal.NewFromFloat(1500),
		TakeProfit:  decimal.NewFromFloat(1000),
		DealTime:    time.Now().UTC(),
	}
	testShare = &model.Share{
		Company: "Microsoft",
		Price:   decimal.NewFromFloat(999),
	}
)

func TestGetProfit(t *testing.T) {
	srv := new(mocks.TradingService)
	hndl := NewEntityDeal(srv, v)
	protoDeal := &proto.Deal{
		DealID:      testDeal.DealID.String(),
		SharesCount: testDeal.SharesCount.InexactFloat64(),
		ProfileID:   testDeal.ProfileID.String(),
		Company:     testDeal.Company,
		StopLoss:    testDeal.StopLoss.InexactFloat64(),
		TakeProfit:  testDeal.TakeProfit.InexactFloat64(),
		DealTime:    timestamppb.Now(),
	}
	srv.On("GetProfit", mock.Anything, mock.AnythingOfType("string"), mock.AnythingOfType("*model.Deal")).Return(testDeal.Profit, nil).Once()
	_, err := hndl.GetProfit(context.Background(), &proto.GetProfitRequest{
		Strategy: testStrategy,
		Deal:     protoDeal,
	})
	require.NoError(t, err)
	srv.AssertExpectations(t)
}

func TestClosePosition(t *testing.T) {
	srv := new(mocks.TradingService)
	hndl := NewEntityDeal(srv, v)
	srv.On("ClosePosition", mock.Anything, mock.AnythingOfType("uuid.UUID"), mock.AnythingOfType("uuid.UUID")).Return(testDeal.Profit, nil).Once()
	_, err := hndl.ClosePosition(context.Background(), &proto.ClosePositionRequest{
		Dealid:    testDeal.DealID.String(),
		Profileid: testDeal.ProfileID.String(),
	})
	require.NoError(t, err)
	srv.AssertExpectations(t)
}

func TestGetUnclosedPositions(t *testing.T) {
	srv := new(mocks.TradingService)
	hndl := NewEntityDeal(srv, v)
	var sliceDeals []*model.Deal
	sliceDeals = append(sliceDeals, testDeal)
	srv.On("GetUnclosedPositions", mock.Anything, mock.AnythingOfType("uuid.UUID")).Return(sliceDeals, nil)
	_, err := hndl.GetUnclosedPositions(context.Background(), &proto.GetUnclosedPositionsRequest{
		Profileid: testDeal.ProfileID.String(),
	})
	require.NoError(t, err)
	srv.AssertExpectations(t)
}

func TestGetPrices(t *testing.T) {
	srv := new(mocks.TradingService)
	hndl := NewEntityDeal(srv, v)
	var prices []model.Share
	prices = append(prices, *testShare)
	srv.On("GetPrices").Return(prices, nil)
	resp, err := hndl.GetPrices(context.Background(), &proto.GetPricesRequest{})
	require.NoError(t, err)
	for _, share := range resp.Share {
		for _, price := range prices {
			price.Company = share.Company
			price.Price = decimal.NewFromFloat(share.Price)
			require.Equal(t, price.Company, testShare.Company)
			require.Equal(t, price.Price, testShare.Price)
		}
	}
	srv.AssertExpectations(t)
}
