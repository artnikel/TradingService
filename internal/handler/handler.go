// Package handler is the top level of the application and it contains request handlers
package handler

import (
	"context"
	"fmt"
	"time"

	"github.com/artnikel/TradingService/internal/model"
	"github.com/artnikel/TradingService/proto"
	"github.com/go-playground/validator/v10"
	"github.com/google/uuid"
	"github.com/shopspring/decimal"
	"github.com/sirupsen/logrus"
)

// TradingService is an interface that contains methods of service for trade
type TradingService interface {
	GetProfit(ctx context.Context, strategy string, deal *model.Deal) (decimal.Decimal, error)
	BalanceOperation(ctx context.Context, balance *model.Balance) (float64, error)
	GetBalance(ctx context.Context, profileid uuid.UUID) (float64, error)
}

// EntityDeal contains Balance Service interface
type EntityDeal struct {
	srvTrading TradingService
	validate   *validator.Validate
	proto.UnimplementedTradingServiceServer
}

// NewEntityDeal accepts Trading Service interface and returns an object of *EntityDeal
func NewEntityDeal(srvTrading TradingService, validate *validator.Validate) *EntityDeal {
	return &EntityDeal{srvTrading: srvTrading, validate: validate}
}

// GetProfit is method that calls method of Trading Service
func (d *EntityDeal) GetProfit(ctx context.Context, req *proto.GetProfitRequest) (*proto.GetProfitResponse, error) {
	profileID, err := uuid.Parse(req.Deal.ProfileID)
	if err != nil {
		logrus.Errorf("error: %v", err)
		return &proto.GetProfitResponse{}, fmt.Errorf("EntityDeal-Strategies: failed to parse id")
	}
	createdDeal := &model.Deal{
		DealID:      uuid.New(),
		SharesCount: decimal.NewFromFloat(req.Deal.SharesCount),
		ProfileID:   profileID,
		Company:     req.Deal.Company,
		StopLoss:    decimal.NewFromFloat(req.Deal.StopLoss),
		TakeProfit:  decimal.NewFromFloat(req.Deal.TakeProfit),
		DealTime:    time.Now().UTC(),
	}
	err = d.validate.StructCtx(ctx, createdDeal)
	if err != nil {
		logrus.Errorf("error: %v", err)
		return &proto.GetProfitResponse{}, fmt.Errorf("EntityDeal-Strategies: failed to validate struct deal")
	}
	profit, err := d.srvTrading.GetProfit(ctx, req.Strategy, createdDeal)
	if err != nil {
		logrus.Errorf("error: %v", err)
		return &proto.GetProfitResponse{}, fmt.Errorf("EntityDeal-Strategies: failed run method strategies")
	}
	return &proto.GetProfitResponse{Profit: profit.InexactFloat64()}, nil
}
