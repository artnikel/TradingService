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
	Subscribe(ctx context.Context, position *model.Deal, subscribersActions chan []*model.Action) error
	AddDeal(ctx context.Context, strategy string, deal *model.Deal) error
}

// BalanceRepository is an interface that contains methods for user manipulation
type BalanceRepository interface {
	BalanceOperation(ctx context.Context, balance *model.Balance) (float64, error)
	GetBalance(ctx context.Context, profileid uuid.UUID) (float64, error)
}

// TradingService contains PriceRepository interface
type TradingService struct {
	priceRep PriceRepository
	bRep     BalanceRepository
}

// NewTradingService accepts PriceRepository object and returnes an object of type *PriceService
func NewTradingService(priceRep PriceRepository, bRep BalanceRepository) *TradingService {
	return &TradingService{
		priceRep: priceRep,
		bRep:     bRep,
	}
}

// nolint gocognit
// Strategies contains 2 options of strategy - long and short
func (ts *TradingService) Strategies(ctx context.Context, strategy string, deal *model.Deal) (decimal.Decimal, error) {
	balanceMoney, err := ts.bRep.GetBalance(ctx, deal.ProfileID)
	if err != nil {
		return decimal.Zero, fmt.Errorf("TradingService-Strategies-GetBalance: error:%w", err)
	}
	balance := &model.Balance{
		BalanceID: uuid.New(),
		ProfileID: deal.ProfileID,
	}
	subscribersActions := make(chan []*model.Action)
	go func() {
		err = ts.priceRep.Subscribe(ctx, deal, subscribersActions)
		if err != nil {
			log.Fatalf("TradingService-Strategies-GetBalance: error:%v", err)
		}
	}()
	taken, isAdded := false, false
	a, b := deal.StopLoss.Mul(deal.ActionsCount), deal.TakeProfit.Mul(deal.ActionsCount)
	if strategy == "short" {
		a, b = b, a
	}
	for {
		select {
		case actions := <-subscribersActions:
			for _, action := range actions {
				//fmt.Printf("Received action: Company=%s, Price=%.2f\n", action.Company, action.Price.InexactFloat64())
				statusCmp := a.GreaterThanOrEqual(action.Price.Mul(deal.ActionsCount))
				if statusCmp {
					deal.EndDealTime = time.Now().UTC()
					deal.Profit = a.Sub(deal.PurchasePrice.Mul(deal.ActionsCount))
					for !isAdded {
						err := ts.priceRep.AddDeal(ctx, strategy, deal)
						if err != nil {
							return decimal.Zero, fmt.Errorf("PriceService-Strategies-AddDeal: error:%w", err)
						}
						balance.Operation = deal.Profit
						_, err = ts.bRep.BalanceOperation(ctx, balance)
						if err != nil {
							return decimal.Zero, fmt.Errorf("PriceService-Strategies-BalanceOperation: error:%w", err)
						}
						isAdded = true
						return deal.Profit, nil
					}
				}
				statusCmp = action.Price.Mul(deal.ActionsCount).GreaterThanOrEqual(b)
				if statusCmp {
					deal.EndDealTime = time.Now().UTC()
					deal.Profit = b.Sub(deal.PurchasePrice.Mul(deal.ActionsCount))
					for !isAdded {
						err := ts.priceRep.AddDeal(ctx, strategy, deal)
						if err != nil {
							return decimal.Zero, fmt.Errorf("PriceService-Strategies-AddDeal: error:%w", err)
						}
						balance.Operation = deal.Profit
						_, err = ts.bRep.BalanceOperation(ctx, balance)
						if err != nil {
							return decimal.Zero, fmt.Errorf("PriceService-Strategies-BalanceOperation: error:%w", err)
						}
						isAdded = true
						return deal.Profit, nil
					}
				}
				for !taken {
					deal.Company = action.Company
					deal.PurchasePrice = action.Price
					taken = true
					if deal.PurchasePrice.Mod(deal.ActionsCount).Cmp(a) == -1 || deal.PurchasePrice.Mod(deal.ActionsCount).Cmp(b) == 1 {
						return decimal.Zero, fmt.Errorf("PriceService-Strategies: purchase price out of takeprofit/stoploss")
					}
					if decimal.NewFromFloat(balanceMoney).Cmp(deal.PurchasePrice.Mul(deal.ActionsCount)) == -1 {
						return decimal.Zero, fmt.Errorf("PriceService-Strategies: not enough money")
					}
					balance.Operation = deal.PurchasePrice.Mul(deal.ActionsCount).Neg()
					_, err := ts.bRep.BalanceOperation(ctx, balance)
					if err != nil {
						return decimal.Zero, fmt.Errorf("PriceService-Strategies-BalanceOperation: error:%w", err)
					}
				}
			}
		case <-ctx.Done():
			return decimal.Zero, nil
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
