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
	"google.golang.org/protobuf/types/known/timestamppb"
)

// TradingService is an interface that contains methods of service for trade
type TradingService interface {
	CreatePosition(ctx context.Context, strategy string, deal *model.Deal) error
	BalanceOperation(ctx context.Context, balance *model.Balance) (float64, error)
	GetBalance(ctx context.Context, profileid uuid.UUID) (float64, error)
	ClosePositionManually(ctx context.Context, dealid uuid.UUID, profileid uuid.UUID) (decimal.Decimal, error)
	GetUnclosedPositions(ctx context.Context, profileid uuid.UUID) ([]*model.Deal, error)
	GetPrices() ([]model.Share, error)
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
func (d *EntityDeal) CreatePosition(ctx context.Context, req *proto.CreatePositionRequest) (*proto.CreatePositionResponse, error) {
	profileID, err := uuid.Parse(req.Deal.ProfileID)
	if err != nil {
		logrus.Errorf("error: %v", err)
		return &proto.CreatePositionResponse{}, fmt.Errorf("EntityDeal-GetProfit: failed to parse profile id")
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
		return &proto.CreatePositionResponse{}, fmt.Errorf("EntityDeal-GetProfit: failed to validate struct deal")
	}
	err = d.srvTrading.CreatePosition(ctx, req.Strategy, createdDeal)
	if err != nil {
		logrus.Errorf("error: %v", err)
		return &proto.CreatePositionResponse{}, fmt.Errorf("EntityDeal-GetProfit: failed to get profit")
	}
	return &proto.CreatePositionResponse{}, nil
}

// ClosePositionManually is method that calls method of Trading Service
func (d *EntityDeal) ClosePositionManually(ctx context.Context, req *proto.ClosePositionManuallyRequest) (*proto.ClosePositionManuallyResponse, error) {
	dealID, err := uuid.Parse(req.Dealid)
	if err != nil {
		logrus.Errorf("error: %v", err)
		return &proto.ClosePositionManuallyResponse{}, fmt.Errorf("EntityDeal-ClosePositionManually: failed to parse id")
	}
	err = d.validate.VarCtx(ctx, dealID, "required")
	if err != nil {
		logrus.Errorf("error: %v", err)
		return &proto.ClosePositionManuallyResponse{}, fmt.Errorf("EntityDeal-ClosePositionManually: failed to validate deal id")
	}
	profileID, err := uuid.Parse(req.Profileid)
	if err != nil {
		logrus.Errorf("error: %v", err)
		return &proto.ClosePositionManuallyResponse{}, fmt.Errorf("EntityDeal-ClosePositionManually: failed to parse id")
	}
	err = d.validate.VarCtx(ctx, profileID, "required")
	if err != nil {
		logrus.Errorf("error: %v", err)
		return &proto.ClosePositionManuallyResponse{}, fmt.Errorf("EntityDeal-ClosePositionManually: failed to validate deal id")
	}
	profit, err := d.srvTrading.ClosePositionManually(ctx, dealID, profileID)
	if err != nil {
		logrus.Errorf("error: %v", err)
		return &proto.ClosePositionManuallyResponse{}, fmt.Errorf("EntityDeal-ClosePositionManually: failed run close position")
	}
	return &proto.ClosePositionManuallyResponse{
		Profit: profit.InexactFloat64(),
	}, nil
}

// GetUnclosedPositions is method that calls method of Trading Service
func (d *EntityDeal) GetUnclosedPositions(ctx context.Context, req *proto.GetUnclosedPositionsRequest) (*proto.GetUnclosedPositionsResponse, error) {
	profileID, err := uuid.Parse(req.Profileid)
	if err != nil {
		logrus.Errorf("error: %v", err)
		return &proto.GetUnclosedPositionsResponse{}, fmt.Errorf("EntityDeal-GetUnclosedPositions: failed to parse profile id")
	}
	err = d.validate.VarCtx(ctx, profileID, "required")
	if err != nil {
		logrus.Errorf("error: %v", err)
		return &proto.GetUnclosedPositionsResponse{}, fmt.Errorf("EntityDeal-GetUnclosedPositions: failed to validate profile id")
	}
	unclosedDeals, err := d.srvTrading.GetUnclosedPositions(ctx, profileID)
	if err != nil {
		logrus.Errorf("error: %v", err)
		return &proto.GetUnclosedPositionsResponse{}, fmt.Errorf("EntityDeal-ClosePosition: failed run close position")
	}
	protoDeals := make([]*proto.Deal, len(unclosedDeals))
	for i, deal := range unclosedDeals {
		protoDeal := &proto.Deal{
			DealID:        deal.DealID.String(),
			SharesCount:   deal.SharesCount.InexactFloat64(),
			Company:       deal.Company,
			PurchasePrice: deal.PurchasePrice.InexactFloat64(),
			StopLoss:      deal.StopLoss.InexactFloat64(),
			TakeProfit:    deal.TakeProfit.InexactFloat64(),
			DealTime:      timestamppb.New(deal.DealTime),
		}
		protoDeals[i] = protoDeal
	}
	return &proto.GetUnclosedPositionsResponse{
		Deal: protoDeals,
	}, nil
}

// GetPrices is method that calls method of Trading Service
func (d *EntityDeal) GetPrices(ctx context.Context, _ *proto.GetPricesRequest) (*proto.GetPricesResponse, error) {
	shares, err := d.srvTrading.GetPrices()
	if err != nil {
		logrus.Errorf("error: %v", err)
		return &proto.GetPricesResponse{}, fmt.Errorf("EntityDeal-GetPrices: failed run close position")
	}
	protoShares := make([]*proto.TradingShare, len(shares))
	for i, share := range shares {
		protoShare := &proto.TradingShare{
			Company: share.Company,
			Price:   share.Price.InexactFloat64(),
		}
		protoShares[i] = protoShare
	}
	return &proto.GetPricesResponse{
		Share: protoShares,
	}, nil
}
