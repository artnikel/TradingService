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
	Subscribe(ctx context.Context, manager chan []*model.Share) error
	AddDeal(ctx context.Context, strategy string, deal *model.Deal) error
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
		priceRep:  priceRep,
		bRep:      bRep,
		chmanager: &model.ChanManager{SubscribersShares: make(chan []*model.Share)},
	}
}

// nolint gocognit
// GetProfit contains 2 options of strategy - long and short
func (ts *TradingService) GetProfit(ctx context.Context, strategy string, deal *model.Deal) (decimal.Decimal, error) {
	ts.chmanager.Mu.Lock()
	defer ts.chmanager.Mu.Unlock()
	balanceMoney, err := ts.bRep.GetBalance(ctx, deal.ProfileID)
	if err != nil {
		return decimal.Zero, fmt.Errorf("PriceService-GetProfit-GetBalance: error:%w", err)
	}
	balance := &model.Balance{
		BalanceID: uuid.New(),
		ProfileID: deal.ProfileID,
	}
	taken, isAdded := false, false
	stopLoss, takeProfit := deal.StopLoss.Mul(deal.SharesCount), deal.TakeProfit.Mul(deal.SharesCount)
	if strategy == "short" {
		stopLoss, takeProfit = takeProfit, stopLoss
	}
	defer func() {
		if deal.EndDealTime.IsZero() && taken {
			balance.Operation = deal.PurchasePrice.Mul(deal.SharesCount)
			_, err := ts.bRep.BalanceOperation(ctx, balance)
			if err != nil {
				log.Fatalf("PriceService-GetProfit-BalanceOperation: error:%v", err)
			}
		}
	}()
	for {
		select {
		case shares := <-ts.chmanager.SubscribersShares:
			for _, share := range shares {
				if share.Company == deal.Company {
					statusCompare := stopLoss.GreaterThanOrEqual(share.Price.Mul(deal.SharesCount))
					if statusCompare {
						deal.EndDealTime = time.Now().UTC()
						deal.Profit = stopLoss.Sub(deal.PurchasePrice.Mul(deal.SharesCount))
						for !isAdded {
							err := ts.priceRep.AddDeal(ctx, strategy, deal)
							if err != nil {
								return decimal.Zero, fmt.Errorf("PriceService-GetProfit-AddDeal: error:%w", err)
							}
							balance.Operation = deal.Profit.Add(deal.PurchasePrice.Mul(deal.SharesCount))
							_, err = ts.bRep.BalanceOperation(ctx, balance)
							if err != nil {
								return decimal.Zero, fmt.Errorf("PriceService-GetProfit-BalanceOperation: error:%w", err)
							}
							isAdded = true
							return deal.Profit, nil
						}
					}
					statusCompare = share.Price.Mul(deal.SharesCount).GreaterThanOrEqual(takeProfit)
					if statusCompare {
						deal.EndDealTime = time.Now().UTC()
						deal.Profit = takeProfit.Sub(deal.PurchasePrice.Mul(deal.SharesCount))
						for !isAdded {
							err := ts.priceRep.AddDeal(ctx, strategy, deal)
							if err != nil {
								return decimal.Zero, fmt.Errorf("PriceService-GetProfit-AddDeal: error:%w", err)
							}
							balance.Operation = deal.Profit.Add(deal.PurchasePrice.Mul(deal.SharesCount))
							_, err = ts.bRep.BalanceOperation(ctx, balance)
							if err != nil {
								return decimal.Zero, fmt.Errorf("PriceService-GetProfit-BalanceOperation: error:%w", err)
							}
							isAdded = true
							return deal.Profit, nil
						}
					}
					for !taken {
						deal.PurchasePrice = share.Price
						taken = true
						if stopLoss.Cmp(deal.PurchasePrice.Mul(deal.SharesCount)) == 1 || deal.PurchasePrice.Mul(deal.SharesCount).Cmp(takeProfit) == 1 {
							return decimal.Zero, fmt.Errorf("PriceService-GetProfit: purchase price out of takeprofit/stoploss")
						}
						if decimal.NewFromFloat(balanceMoney).Cmp(deal.PurchasePrice.Mul(deal.SharesCount)) == -1 {
							return decimal.Zero, fmt.Errorf("PriceService-GetProfit: not enough money")
						}
						balance.Operation = deal.PurchasePrice.Mul(deal.SharesCount).Neg()
						_, err := ts.bRep.BalanceOperation(ctx, balance)
						if err != nil {
							return decimal.Zero, fmt.Errorf("PriceService-GetProfit-BalanceOperation: error:%w", err)
						}
					}
				}
			}
		case <-ctx.Done():
			return decimal.Zero, nil
		}
	}
}

// Subscribe is a method of TradingService that calls method of Repository as goroutine
func (ts *TradingService) Subscribe(ctx context.Context) {
	manager := make(chan []*model.Share)
	go func() {
		err := ts.priceRep.Subscribe(ctx, manager)
		if err != nil {
			log.Fatalf("TradingService-Subscribe: error:%v", err)
		}
	}()
	defer close(manager)
	for _, shares := range <-manager {
		ts.chmanager.SubscribersShares <- []*model.Share{shares}
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
