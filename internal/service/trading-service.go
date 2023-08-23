// Package service contains business logic of a project
package service

import (
	"context"
	"fmt"
	"log"
	"strings"
	"time"

	"github.com/artnikel/TradingService/internal/config"
	"github.com/artnikel/TradingService/internal/model"
	"github.com/caarlos0/env"
	"github.com/google/uuid"
	"github.com/shopspring/decimal"
)

// PriceRepository is interface with method for reading prices
type PriceRepository interface {
	Subscribe(ctx context.Context, manager chan model.Share) error
	AddPosition(ctx context.Context, strategy string, deal *model.Deal) error
	ClosePosition(ctx context.Context, deal *model.Deal) error
	GetPositionInfoByDealID(ctx context.Context, dealid uuid.UUID) (model.Deal, error)
	GetUnclosedPositions(ctx context.Context, profileid uuid.UUID) ([]*model.Deal, error)
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
	if _, ok := ts.chmanager.SubscribersShares[deal.ProfileID]; !ok {
		ts.chmanager.SubscribersShares[deal.ProfileID] = make(map[string]chan model.Share)
		ts.chmanager.SubscribersShares[deal.ProfileID][deal.Company] = make(chan model.Share)
	}
	balanceMoney, err := ts.bRep.GetBalance(ctx, deal.ProfileID)
	if err != nil {
		return decimal.Zero, fmt.Errorf("TradingService-GetProfit-GetBalance: error:%w", err)
	}
	balance := &model.Balance{ProfileID: deal.ProfileID}
	takenPurchasePrice := false
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
	for {
		select {
		case share := <-ts.chmanager.SubscribersShares[deal.ProfileID][deal.Company]:
			if !takenPurchasePrice {
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
			if stopLoss.GreaterThanOrEqual(share.Price.Mul(deal.SharesCount)) || share.Price.Mul(deal.SharesCount).GreaterThanOrEqual(takeProfit) {
				profit, err := ts.ClosePosition(ctx, deal.DealID, deal.ProfileID)
				if err != nil {
					return decimal.Zero, fmt.Errorf("TradingService-GetProfit-ClosePosition: error:%w", err)
				}
				return profit, nil
			}
		case <-ctx.Done():
			return decimal.Zero, nil
		}
	}
}

// ClosePosition is a method that calls method of Repository and returns profit
func (ts *TradingService) ClosePosition(ctx context.Context, dealid, profileid uuid.UUID) (decimal.Decimal, error) {
	var balance model.Balance
	balance.ProfileID = profileid
	deal, err := ts.priceRep.GetPositionInfoByDealID(ctx, dealid)
	if err != nil {
		return decimal.Zero, fmt.Errorf("TradingService-ClosePosition-GetPositionInfoByDealID: error:%w", err)
	}
	if _, ok := ts.chmanager.SubscribersShares[profileid]; !ok {
		ts.chmanager.SubscribersShares[profileid] = make(map[string]chan model.Share)
		ts.chmanager.SubscribersShares[profileid][deal.Company] = make(chan model.Share)
	}
	for {
		select {
		case <-ctx.Done():
			return decimal.Zero, nil
		case share := <-ts.chmanager.SubscribersShares[profileid][deal.Company]:
			deal.EndDealTime = time.Now().UTC()
			deal.DealID = dealid
			if deal.TakeProfit.Cmp(deal.StopLoss) == 1 {
				deal.Profit = share.Price.Mul(deal.SharesCount).Sub(deal.PurchasePrice.Mul(deal.SharesCount))
			}
			if deal.StopLoss.Cmp(deal.TakeProfit) == 1 {
				deal.Profit = deal.PurchasePrice.Mul(deal.SharesCount).Sub(share.Price.Mul(deal.SharesCount))
			}
			err = ts.priceRep.ClosePosition(ctx, &deal)
			if err != nil {
				return decimal.Zero, fmt.Errorf("TradingService-GetProfit-ClosePosition: error:%w", err)
			}
			balance.Operation = deal.Profit.Add(deal.PurchasePrice.Mul(deal.SharesCount))
			_, err = ts.bRep.BalanceOperation(ctx, &balance)
			if err != nil {
				return decimal.Zero, fmt.Errorf("TradingService-GetProfit-BalanceOperation: error:%w", err)
			}
			return deal.Profit, nil
		}
	}
}

// Subscribe is a method of TradingService that calls method of Repository as goroutine
func (ts *TradingService) Subscribe(ctx context.Context) {
	var cfg config.Variables
	if err := env.Parse(&cfg); err != nil {
		log.Fatal("could not parse config: ", err)
	}
	companyShares := strings.Split(cfg.CompanyShares, ",")
	manager := make(chan model.Share, len(companyShares))
	go func() {
		err := ts.priceRep.Subscribe(ctx, manager)
		if err != nil {
			log.Fatalf("TradingService-Subscribe: error:%v", err)
		}
	}()
	for {
		for subID, subShares := range ts.chmanager.SubscribersShares {
			for subCompany := range subShares {
				select {
				case <-ctx.Done():
					return
				case share, ok := <-manager:
					if !ok {
						break
					}
					if share.Company == subCompany {
						ts.chmanager.SubscribersShares[subID][subCompany] <- share
					}
				default:
					continue
				}
			}
		}
	}
}

func (ts *TradingService) GetPrices(ctx context.Context) ([]model.Share, error) {
	var cfg config.Variables
	if err := env.Parse(&cfg); err != nil {
		log.Fatal("could not parse config: ", err)
	}
	var shares []model.Share
	companyShares := strings.Split(cfg.CompanyShares, ",")
	priceUUID := uuid.New()
	ts.chmanager.Mu.Lock() 
	defer ts.chmanager.Mu.Unlock()

	if _, ok := ts.chmanager.SubscribersShares[priceUUID]; !ok {
		ts.chmanager.SubscribersShares[priceUUID] = make(map[string]chan model.Share)
	}

	for _, company := range companyShares {
		ts.chmanager.SubscribersShares[priceUUID][company] = make(chan model.Share)
	}
	for i := 0; i < len(companyShares); {
		select {
		case <-ctx.Done():
			return nil, nil
		case share := <-ts.chmanager.SubscribersShares[priceUUID][companyShares[i]]:
			shares = append(shares, share)
			i++
		default:
			continue
		}
	}
	return shares, nil
}

// GetUnclosedPositions is a method of TradingService calls method of Repository
func (ts *TradingService) GetUnclosedPositions(ctx context.Context, profileid uuid.UUID) ([]*model.Deal, error) {
	unclosedDeals, err := ts.priceRep.GetUnclosedPositions(ctx, profileid)
	if err != nil {
		return nil, fmt.Errorf("TradingService-GetUnclosedPositions: error:%w", err)
	}
	return unclosedDeals, nil
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
