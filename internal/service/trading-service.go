// Package service contains business logic of a project
package service

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/artnikel/TradingService/internal/model"
	"github.com/google/uuid"
	"github.com/shopspring/decimal"
)

// PriceRepository is interface with method for reading prices
type PriceRepository interface {
	Subscribe(ctx context.Context, manager chan model.Share) error
	AddPosition(ctx context.Context, strategy string, deal *model.Deal) error
	ClosePosition(ctx context.Context, deal *model.Deal) error
}

// BalanceRepository is an interface that contains methods for user manipulation
type BalanceRepository interface {
	BalanceOperation(ctx context.Context, balance *model.Balance) (float64, error)
	GetBalance(ctx context.Context, profileid uuid.UUID) (float64, error)
}

// TradingService contains PriceRepository interface
type TradingService struct {
	priceRep  PriceRepository
	bRep      BalanceRepository
	chmanager *model.ChanManager
}

// NewTradingService accepts PriceRepository object and returnes an object of type *PriceService
func NewTradingService(priceRep PriceRepository, bRep BalanceRepository) *TradingService {
	return &TradingService{
		priceRep: priceRep,
		bRep:     bRep,
		chmanager: &model.ChanManager{
			SubscribersShares: make(map[uuid.UUID]map[string]chan model.Share)},
	}
}

// GetProfit contains 2 options of strategy - long and short. He`s added and closed position, and returns profit.
func (ts *TradingService) GetProfit(ctx context.Context, strategy string, deal *model.Deal) (decimal.Decimal, error) {
	balanceMoney, err := ts.bRep.GetBalance(ctx, deal.ProfileID)
	if err != nil {
		return decimal.Zero, fmt.Errorf("TradingService-GetProfit-GetBalance: error:%w", err)
	}
	balance := &model.Balance{ProfileID: deal.ProfileID}
	takenPurchasePrice, positionClosed := false, false
	stopLoss, takeProfit := deal.StopLoss.Mul(deal.SharesCount), deal.TakeProfit.Mul(deal.SharesCount)
	if strategy == "short" {
		stopLoss, takeProfit = takeProfit, stopLoss
	}
	defer func() {
		if deal.EndDealTime.IsZero() && takenPurchasePrice {
			balance.Operation = deal.PurchasePrice.Mul(deal.SharesCount)
			_, err := ts.bRep.BalanceOperation(ctx, balance)
			if err != nil {
				log.Fatalf("TradingService-GetProfit-BalanceOperation: error:%v", err)
			}
		}
	}()
	select {
	case share := <-ts.chmanager.SubscribersShares[deal.ProfileID][deal.Company]:
		statusCompare := stopLoss.GreaterThanOrEqual(share.Price.Mul(deal.SharesCount))
		if statusCompare {
			deal.Profit = stopLoss.Sub(deal.PurchasePrice.Mul(deal.SharesCount))
			profit, err := ts.ClosePosition(ctx, &positionClosed, deal, balance)
			if err != nil {
				return decimal.Zero, fmt.Errorf("TradingService-GetProfit-ClosePosition: error:%w", err)
			}
			return profit, nil
		}
		statusCompare = share.Price.Mul(deal.SharesCount).GreaterThanOrEqual(takeProfit)
		if statusCompare {
			deal.Profit = takeProfit.Sub(deal.PurchasePrice.Mul(deal.SharesCount))
			profit, err := ts.ClosePosition(ctx, &positionClosed, deal, balance)
			if err != nil {
				return decimal.Zero, fmt.Errorf("TradingService-GetProfit-ClosePosition: error:%w", err)
			}
			return profit, nil
		}
		for !takenPurchasePrice {
			takenPurchasePrice = true
			deal.PurchasePrice = share.Price
			if stopLoss.Cmp(deal.PurchasePrice.Mul(deal.SharesCount)) == 1 || deal.PurchasePrice.Mul(deal.SharesCount).Cmp(takeProfit) == 1 {
				return decimal.Zero, fmt.Errorf("TradingService-GetProfit: purchase price out of takeprofit/stoploss")
			}
			if decimal.NewFromFloat(balanceMoney).Cmp(deal.PurchasePrice.Mul(deal.SharesCount)) == -1 {
				return decimal.Zero, fmt.Errorf("TradingService-GetProfit: not enough money")
			}
			err := ts.priceRep.AddPosition(ctx, strategy, deal)
			if err != nil {
				return decimal.Zero, fmt.Errorf("TradingService-GetProfit-AddPosition: error:%w", err)
			}
			balance.Operation = deal.PurchasePrice.Mul(deal.SharesCount).Neg()
			_, err = ts.bRep.BalanceOperation(ctx, balance)
			if err != nil {
				return decimal.Zero, fmt.Errorf("TradingService-GetProfit-BalanceOperation: error:%w", err)
			}
		}
	case <-ctx.Done():
		return decimal.Zero, nil
	}
	return decimal.Zero, nil
}

// ClosePosition is a method that calls method of Repository and returns profit
func (ts *TradingService) ClosePosition(ctx context.Context, positionClosed *bool, deal *model.Deal, balance *model.Balance) (decimal.Decimal, error) {
	if !*positionClosed {
		deal.EndDealTime = time.Now().UTC()
		err := ts.priceRep.ClosePosition(ctx, deal)
		if err != nil {
			return decimal.Zero, fmt.Errorf("TradingService-GetProfit-ClosePosition: error:%w", err)
		}
		balance.Operation = deal.Profit.Add(deal.PurchasePrice.Mul(deal.SharesCount))
		_, err = ts.bRep.BalanceOperation(ctx, balance)
		if err != nil {
			return decimal.Zero, fmt.Errorf("TradingService-GetProfit-BalanceOperation: error:%w", err)
		}
		*positionClosed = true
		return deal.Profit, nil
	}
	return decimal.Zero, fmt.Errorf("TradingService-ClosePosition: position is closed")
}

// Subscribe is a method of TradingService that calls method of Repository as goroutine
func (ts *TradingService) Subscribe(ctx context.Context) {
	manager := make(chan model.Share, 5)
	go func() {
		err := ts.priceRep.Subscribe(ctx, manager)
		if err != nil {
			log.Fatalf("TradingService-Subscribe: error:%v", err)
		}
	}()
	for {
		for subID := range ts.chmanager.SubscribersShares {
			for share := range manager {
				select {
				case <-ctx.Done():
					return
				case ts.chmanager.SubscribersShares[subID][share.Company] <- share:
				}
			}
		}
	}
}

// BalanceOperation is a method of TradingService calls method of Repository
func (ts *TradingService) BalanceOperation(ctx context.Context, balance *model.Balance) (float64, error) {
	operation, err := ts.bRep.BalanceOperation(ctx, balance)
	if err != nil {
		return 0, fmt.Errorf("TradingService-BalanceOperation: error: %w", err)
	}
	return operation, nil
}

// GetBalance is a method of TradingService calls method of Repository
func (ts *TradingService) GetBalance(ctx context.Context, profileid uuid.UUID) (float64, error) {
	money, err := ts.bRep.GetBalance(ctx, profileid)
	if err != nil {
		return 0, fmt.Errorf("TradingService-GetBalance: error: %w", err)
	}
	return money, nil
}
