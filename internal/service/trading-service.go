// Package service contains business logic of a project
package service

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"strings"
	"time"

	"github.com/artnikel/TradingService/internal/config"
	"github.com/artnikel/TradingService/internal/model"
	"github.com/google/uuid"
	"github.com/lib/pq"
	"github.com/shopspring/decimal"
	"github.com/sirupsen/logrus"
)

// PriceRepository is interface with method for reading prices
type PriceRepository interface {
	Subscribe(ctx context.Context, manager chan model.Share) error
	CreatePosition(ctx context.Context, strategy string, deal *model.Deal) error
	ClosePosition(ctx context.Context, deal *model.Deal) error
	GetUnclosedPositions(ctx context.Context, profileid uuid.UUID) ([]*model.Deal, error)
	GetUnclosedPositionsForAll(ctx context.Context) ([]*model.Deal, error)
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
	cfg       config.Variables
}

// NewTradingService accepts PriceRepository object and returnes an object of type *PriceService
func NewTradingService(priceRep PriceRepository, bRep BalanceRepository, cfg config.Variables) *TradingService {
	return &TradingService{
		priceRep: priceRep,
		bRep:     bRep,
		chmanager: &model.ChanManager{
			SubscribersShares: make(map[uuid.UUID]map[string]chan model.Share),
			PricesMap:         make(map[string]decimal.Decimal),
			Positions:         make(map[uuid.UUID]model.Deal)},
		cfg: cfg,
	}
}

// identifyStrategy identifies strategy (long or short) by takeprofit and stoploss
func (ts *TradingService) identifyStrategy(takeProfit, stopLoss decimal.Decimal) string {
	if takeProfit.Cmp(stopLoss) == 1 {
		return "long"
	}
	if stopLoss.Cmp(takeProfit) == 1 {
		return "short"
	}
	return ""
}

// WaitForNotification listens changes of positions database
func (ts *TradingService) WaitForNotification(ctx context.Context, l *pq.Listener) {
	var valueToParse = struct {
		Action     string `json:"action"`
		model.Deal `json:"data"`
	}{}
	for {
		select {
		case notification := <-l.Notify:
			err := json.Unmarshal([]byte(notification.Extra), &valueToParse)
			logrus.Info(notification)
			if err != nil {
				logrus.WithField("Payloads", notification.Extra).Errorf("TradingService-WaitForNotification: %v", err)
				continue
			}
			if valueToParse.Action == "INSERT" {
				ts.chmanager.Mu.Lock()
				ts.chmanager.Positions[valueToParse.DealID] = valueToParse.Deal
				ts.chmanager.Mu.Unlock()
				strategy := ts.identifyStrategy(valueToParse.TakeProfit, valueToParse.StopLoss)
				go func() {
					err = ts.GetProfit(ctx, strategy, &valueToParse.Deal)
					if err != nil {
						log.Fatalf("TradingService-CreatePosition-GetProfit: error:%v", err)
					}
				}()
			}
			if valueToParse.Action == "UPDATE" {
				ts.chmanager.Mu.Lock()
				delete(ts.chmanager.Positions, valueToParse.DealID)
				delete(ts.chmanager.SubscribersShares[valueToParse.ProfileID], valueToParse.Company)
				delete(ts.chmanager.SubscribersShares, valueToParse.ProfileID)
				ts.chmanager.Mu.Unlock()
			}
		case <-ctx.Done():
			return
		default:
			continue
		}
	}
}

// CreatePosition insert new position in database
func (ts *TradingService) CreatePosition(ctx context.Context, strategy string, deal *model.Deal) error {
	ts.chmanager.Mu.Lock()
	if _, ok := ts.chmanager.SubscribersShares[deal.ProfileID]; !ok {
		ts.chmanager.SubscribersShares[deal.ProfileID] = make(map[string]chan model.Share)
	}
	if _, ok := ts.chmanager.SubscribersShares[deal.ProfileID][deal.Company]; !ok {
		ts.chmanager.SubscribersShares[deal.ProfileID][deal.Company] = make(chan model.Share)
	}
	ts.chmanager.Mu.Unlock()

	balanceMoney, err := ts.bRep.GetBalance(ctx, deal.ProfileID)
	if err != nil {
		return fmt.Errorf("TradingService-CreatePosition-GetBalance: error:%w", err)
	}
	balance := &model.Balance{ProfileID: deal.ProfileID}
	stopLoss, takeProfit := deal.StopLoss.Mul(deal.SharesCount), deal.TakeProfit.Mul(deal.SharesCount)
	if strategy == "short" {
		stopLoss, takeProfit = takeProfit, stopLoss
	}
	select {
	case share := <-ts.chmanager.SubscribersShares[deal.ProfileID][deal.Company]:
		deal.PurchasePrice = share.Price
		if stopLoss.Cmp(deal.PurchasePrice.Mul(deal.SharesCount)) == 1 || deal.PurchasePrice.Mul(deal.SharesCount).Cmp(takeProfit) == 1 {
			return fmt.Errorf("TradingService-CreatePosition: purchase price out of takeprofit/stoploss")
		}
		if decimal.NewFromFloat(balanceMoney).Cmp(deal.PurchasePrice.Mul(deal.SharesCount)) == -1 {
			return fmt.Errorf("TradingService-CreatePosition: not enough money")
		}
		err := ts.priceRep.CreatePosition(ctx, strategy, deal)
		if err != nil {
			return fmt.Errorf("TradingService-CreatePosition: error:%w", err)
		}
		balance.Operation = deal.PurchasePrice.Mul(deal.SharesCount).Neg()
		_, err = ts.bRep.BalanceOperation(ctx, balance)
		if err != nil {
			return fmt.Errorf("TradingService-CreatePosition-BalanceOperation: error:%w", err)
		}
		return nil
	case <-ctx.Done():
		return nil
	}
}

// GetProfit monitors prices and compares them to takeprofit and stoploss
func (ts *TradingService) GetProfit(ctx context.Context, strategy string, deal *model.Deal) error {
	ts.chmanager.Mu.RLock()
	if _, ok := ts.chmanager.SubscribersShares[deal.ProfileID]; !ok {
		return fmt.Errorf("TradingService-GetProfit: value in map SubscribersShares with key profileid is not exist")
	}
	if _, ok := ts.chmanager.SubscribersShares[deal.ProfileID][deal.Company]; !ok {
		return fmt.Errorf("TradingService-GetProfit: value in map SubscribersShares with key company is not exist")
	}
	ts.chmanager.Mu.RUnlock()

	stopLoss, takeProfit := deal.StopLoss.Mul(deal.SharesCount), deal.TakeProfit.Mul(deal.SharesCount)
	if strategy == "short" {
		stopLoss, takeProfit = takeProfit, stopLoss
	}
	for {
		select {
		case share := <-ts.chmanager.SubscribersShares[deal.ProfileID][deal.Company]:
			ts.chmanager.Mu.RLock()
			if _, ok := ts.chmanager.Positions[deal.DealID]; !ok {
				ts.chmanager.Mu.RUnlock()
				break
			}
			ts.chmanager.Mu.RUnlock()

			if strategy == "long" {
				fmt.Println("Deal ID: ", deal.DealID, " Profit :", share.Price.Mul(deal.SharesCount).Sub(deal.PurchasePrice.Mul(deal.SharesCount)))
			}
			if strategy == "short" {
				fmt.Println("Deal ID: ", deal.DealID, " Profit :", deal.PurchasePrice.Mul(deal.SharesCount).Sub(share.Price.Mul(deal.SharesCount)))
			}
			if stopLoss.GreaterThanOrEqual(share.Price.Mul(deal.SharesCount)) || share.Price.Mul(deal.SharesCount).GreaterThanOrEqual(takeProfit) {
				profit, err := ts.ClosePosition(ctx, deal.DealID, deal.ProfileID, share.Price)
				if err != nil {
					return fmt.Errorf("TradingService-GetProfit-ClosePosition: error:%w", err)
				}
				fmt.Println("Position closed with profit: ", profit)
				return nil
			}
		case <-ctx.Done():
			return nil
		}
	}
}

// ClosePosition is a method that calls method of Repository and returns profit
func (ts *TradingService) ClosePosition(ctx context.Context, dealid, profileid uuid.UUID, sharePrice decimal.Decimal) (decimal.Decimal, error) {
	var balance model.Balance
	balance.ProfileID = profileid

	ts.chmanager.Mu.RLock()
	if _, ok := ts.chmanager.Positions[dealid]; !ok {
		return decimal.Zero, fmt.Errorf("TradingService-ClosePosition: key of map Positions is not exist")
	}
	deal := ts.chmanager.Positions[dealid]
	if _, ok := ts.chmanager.SubscribersShares[profileid]; !ok {
		return decimal.Zero, fmt.Errorf("TradingService-ClosePosition: value in map SubscribersShares with key profileid is not exist")
	}
	if _, ok := ts.chmanager.SubscribersShares[profileid][deal.Company]; !ok {
		return decimal.Zero, fmt.Errorf("TradingService-ClosePosition: value in map SubscribersShares with key company is not exist")
	}
	ts.chmanager.Mu.RUnlock()

	deal.EndDealTime = time.Now().UTC()
	deal.DealID = dealid
	if deal.TakeProfit.Cmp(deal.StopLoss) == 1 {
		deal.Profit = sharePrice.Mul(deal.SharesCount).Sub(deal.PurchasePrice.Mul(deal.SharesCount))
	}
	if deal.StopLoss.Cmp(deal.TakeProfit) == 1 {
		deal.Profit = deal.PurchasePrice.Mul(deal.SharesCount).Sub(sharePrice.Mul(deal.SharesCount))
	}
	err := ts.priceRep.ClosePosition(ctx, &deal)
	if err != nil {
		return decimal.Zero, fmt.Errorf("TradingService-ClosePosition: error:%w", err)
	}
	balance.Operation = deal.Profit.Add(deal.PurchasePrice.Mul(deal.SharesCount))
	_, err = ts.bRep.BalanceOperation(ctx, &balance)
	if err != nil {
		return decimal.Zero, fmt.Errorf("TradingService-ClosePosition-BalanceOperation: error:%w", err)
	}
	return deal.Profit, nil
}

func (ts *TradingService) ClosePositionManually(ctx context.Context, dealid, profileid uuid.UUID) (decimal.Decimal, error) {
	deal := ts.chmanager.Positions[dealid]
	for {
		select {
		case <-ctx.Done():
			return decimal.Zero, nil
		case share := <-ts.chmanager.SubscribersShares[profileid][deal.Company]:
			profit, err := ts.ClosePosition(ctx, dealid, profileid, share.Price)
			if err != nil {
				return decimal.Zero, fmt.Errorf("TradingService-ClosePositionManually-ClosePosition: error: %w", err)
			}
			return profit, nil
		}
	}
}

// Subscribe is a method of TradingService that calls method of Repository as goroutine
func (ts *TradingService) Subscribe(ctx context.Context) {
	companyShares := strings.Split(ts.cfg.CompanyShares, ",")
	manager := make(chan model.Share, len(companyShares))
	go func() {
		err := ts.priceRep.Subscribe(ctx, manager)
		if err != nil {
			log.Fatalf("TradingService-Subscribe: error:%v", err)
		}
	}()
	for {
		select {
		case <-ctx.Done():
			return
		case share := <-manager:
			ts.chmanager.Mu.Lock()
			ts.chmanager.PricesMap[share.Company] = share.Price
			ts.chmanager.Mu.Unlock()
			for _, subShares := range ts.chmanager.SubscribersShares {
				if v, ok := subShares[share.Company]; ok {
					v <- share
				}
			}
		default:
			continue
		}
	}
}

// GetPrices is method that read actual prices of shares
func (ts *TradingService) GetPrices() ([]model.Share, error) {
	var shares []model.Share
	companyShares := strings.Split(ts.cfg.CompanyShares, ",")
	ts.chmanager.Mu.Lock()
	defer ts.chmanager.Mu.Unlock()
	for _, company := range companyShares {
		if price, exists := ts.chmanager.PricesMap[company]; exists {
			shares = append(shares, model.Share{
				Company: company,
				Price:   price,
			})
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

// BackupUnclosedPositions is a method of TradingService calls method of Repository
func (ts *TradingService) BackupUnclosedPositions(ctx context.Context) error {
	unclosedDeals, err := ts.priceRep.GetUnclosedPositionsForAll(ctx)
	if err != nil {
		return fmt.Errorf("TradingService-BackupUnclosedPositions-GetUnclosedPositionsForAll: error:%w", err)
	}
	for _, unclosedDeal := range unclosedDeals {
		deal := unclosedDeal
		ts.chmanager.Mu.Lock()
		ts.chmanager.Positions[deal.DealID] = *deal
		ts.chmanager.SubscribersShares[deal.ProfileID] = make(map[string]chan model.Share)
		ts.chmanager.SubscribersShares[deal.ProfileID][deal.Company] = make(chan model.Share)
		ts.chmanager.Mu.Unlock()
		strategy := ts.identifyStrategy(deal.TakeProfit, deal.StopLoss)
		go func() {
			err := ts.GetProfit(ctx, strategy, deal)
			if err != nil {
				log.Fatalf("TradingService-BackupUnclosedPositions-GetProfit: err:%v", err)
			}
		}()
	}
	return nil
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
