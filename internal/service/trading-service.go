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
	Subscribe(ctx context.Context, position *model.Deal, subscribersActions chan []*model.Share) error
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
	subscribersActions := make(chan []*model.Share)
	go func() {
		err = ts.priceRep.Subscribe(ctx, deal, subscribersActions)
		if err != nil {
			log.Fatalf("TradingService-Strategies-GetBalance: error:%v", err)
		}
	}()
	taken, isAdded := false, false
	a, b := deal.StopLoss.Mul(deal.SharesCount), deal.TakeProfit.Mul(deal.SharesCount)
	if strategy == "short" {
		a, b = b, a
	}
	defer func() {
		if deal.EndDealTime.IsZero() {
			balance.Operation = deal.PurchasePrice.Mul(deal.SharesCount)
			_, err := ts.bRep.BalanceOperation(ctx, balance)
			if err != nil {
				log.Fatalf("PriceService-Strategies-BalanceOperation: error:%v", err)
			}
		}
	}()
	for {
		select {
		case actions := <-subscribersActions:
			for _, action := range actions {
				statusCmp := a.GreaterThanOrEqual(action.Price.Mul(deal.SharesCount))
				if statusCmp {
					deal.EndDealTime = time.Now().UTC()
					deal.Profit = a.Sub(deal.PurchasePrice.Mul(deal.SharesCount))
					for !isAdded {
						err := ts.priceRep.AddDeal(ctx, strategy, deal)
						if err != nil {
							return decimal.Zero, fmt.Errorf("PriceService-Strategies-AddDeal: error:%w", err)
						}
						balance.Operation = deal.Profit.Add(deal.PurchasePrice.Mul(deal.SharesCount))
						_, err = ts.bRep.BalanceOperation(ctx, balance)
						if err != nil {
							return decimal.Zero, fmt.Errorf("PriceService-Strategies-BalanceOperation: error:%w", err)
						}
						isAdded = true
						return deal.Profit, nil
					}
				}
				statusCmp = action.Price.Mul(deal.SharesCount).GreaterThanOrEqual(b)
				if statusCmp {
					deal.EndDealTime = time.Now().UTC()
					deal.Profit = b.Sub(deal.PurchasePrice.Mul(deal.SharesCount))
					for !isAdded {
						err := ts.priceRep.AddDeal(ctx, strategy, deal)
						if err != nil {
							return decimal.Zero, fmt.Errorf("PriceService-Strategies-AddDeal: error:%w", err)
						}
						balance.Operation = deal.Profit.Add(deal.PurchasePrice.Mul(deal.SharesCount))
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
					if a.Cmp(deal.PurchasePrice.Mul(deal.SharesCount)) == 1 || deal.PurchasePrice.Mul(deal.SharesCount).Cmp(b) == 1 {
						return decimal.Zero, fmt.Errorf("PriceService-Strategies: purchase price out of takeprofit/stoploss")
					}
					if decimal.NewFromFloat(balanceMoney).Cmp(deal.PurchasePrice.Mul(deal.SharesCount)) == -1 {
						return decimal.Zero, fmt.Errorf("PriceService-Strategies: not enough money")
					}
					balance.Operation = deal.PurchasePrice.Mul(deal.SharesCount).Neg()
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
