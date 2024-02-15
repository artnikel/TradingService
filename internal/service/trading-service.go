// Package service contains business logic of a project
package service

import (
	"context"
	"encoding/json"
	"fmt"
	"strings"
	"sync"
	"time"

	"github.com/artnikel/TradingService/internal/config"
	berrors "github.com/artnikel/TradingService/internal/errors"
	"github.com/artnikel/TradingService/internal/model"
	"github.com/google/uuid"
	"github.com/lib/pq"
	"github.com/shopspring/decimal"
	"github.com/sirupsen/logrus"
)

// PriceRepository is interface with method for reading prices
type PriceRepository interface {
	Subscribe(ctx context.Context, manager chan model.Share) error
	CreatePosition(ctx context.Context, deal *model.Deal) error
	ClosePosition(ctx context.Context, deal *model.Deal) error
	GetUnclosedPositions(ctx context.Context, profileid uuid.UUID) ([]*model.Deal, error)
	GetClosedPositions(ctx context.Context, profileid uuid.UUID) ([]*model.Deal, error)
	GetUnclosedPositionsForAll(ctx context.Context) ([]*model.Deal, error)
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
	manager  *model.Manager
	cfg      config.Variables
}

// NewTradingService accepts PriceRepository object and returnes an object of type *PriceService
func NewTradingService(priceRep PriceRepository, bRep BalanceRepository, cfg config.Variables) *TradingService {
	return &TradingService{
		priceRep: priceRep,
		bRep:     bRep,
		manager: &model.Manager{
			PricesMap: make(map[string]decimal.Decimal),
			Positions: make(map[uuid.UUID]map[uuid.UUID]model.Deal),
			Mu:        sync.RWMutex{}},
		cfg: cfg,
	}
}

// two constants that contain trading strategies implemented in the application
const (
	long  = "long"
	short = "short"
)

// identifyStrategy identifies strategy (long or short) by comparing takeprofit and stoploss
func (ts *TradingService) identifyStrategy(takeProfit, stopLoss decimal.Decimal) (string, error) {
	if takeProfit.Cmp(stopLoss) == 1 {
		return long, nil
	}
	if stopLoss.Cmp(takeProfit) == 1 {
		return short, nil
	}
	return "", fmt.Errorf("stoploss = takeprofit")
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
				logrus.WithField("Payloads", notification.Extra).Errorf("unmarshal %v", err)
				continue
			}
			if valueToParse.Action == "INSERT" {
				ts.manager.Mu.Lock()
				ts.manager.Positions[valueToParse.ProfileID][valueToParse.DealID] = valueToParse.Deal
				ts.manager.Mu.Unlock()
				go ts.getProfit(ctx, valueToParse.Deal)
			}
			if valueToParse.Action == "UPDATE" {
				ts.manager.Mu.Lock()
				delete(ts.manager.Positions[valueToParse.ProfileID], valueToParse.DealID)
				if len(ts.manager.Positions[valueToParse.ProfileID]) == 0 {
					delete(ts.manager.Positions, valueToParse.ProfileID)
				}
				ts.manager.Mu.Unlock()
			}
		case <-ctx.Done():
			return
		default:
			continue
		}
	}
}

// CreatePosition defines purchase price and insert new position in database
func (ts *TradingService) CreatePosition(ctx context.Context, deal *model.Deal) error {
	ts.manager.Mu.Lock()
	if _, ok := ts.manager.Positions[deal.ProfileID]; !ok {
		ts.manager.Positions[deal.ProfileID] = make(map[uuid.UUID]model.Deal)
	}
	deal.PurchasePrice = ts.manager.PricesMap[deal.Company]
	ts.manager.Mu.Unlock()
	balanceMoney, err := ts.bRep.GetBalance(ctx, deal.ProfileID)
	if err != nil {
		return fmt.Errorf("getBalance %w", err)
	}
	balance := &model.Balance{ProfileID: deal.ProfileID}
	stopLoss, takeProfit := deal.StopLoss.Mul(deal.SharesCount), deal.TakeProfit.Mul(deal.SharesCount)
	strategy, err := ts.identifyStrategy(deal.TakeProfit, deal.StopLoss)
	if err != nil {
		return fmt.Errorf("identifyStrategy %w", err)
	}
	if strategy == short {
		stopLoss, takeProfit = takeProfit, stopLoss
	}
	if stopLoss.Cmp(deal.PurchasePrice.Mul(deal.SharesCount)) == 1 || deal.PurchasePrice.Mul(deal.SharesCount).Cmp(takeProfit) == 1 {
		return berrors.New(berrors.PurchasePriceOut)
	}
	if decimal.NewFromFloat(balanceMoney).Cmp(deal.PurchasePrice.Mul(deal.SharesCount)) == -1 {
		return berrors.New(berrors.NotEnoughMoney)
	}
	balance.Operation = deal.PurchasePrice.Mul(deal.SharesCount).Neg()
	_, err = ts.bRep.BalanceOperation(ctx, balance)
	if err != nil {
		return fmt.Errorf("balanceOperation %w", err)
	}
	err = ts.priceRep.CreatePosition(ctx, deal)
	if err != nil {
		return fmt.Errorf("createPosition %w", err)
	}
	return nil
}

// nolint gocrtitc
// getProfit monitors prices and compares them to takeprofit and stoploss
func (ts *TradingService) getProfit(ctx context.Context, deal model.Deal) {
	ts.manager.Mu.RLock()
	if _, ok := ts.manager.Positions[deal.ProfileID]; !ok {
		logrus.Errorf("value in map Positions with key profileid: %s is not exist", deal.ProfileID.String())
	}
	if _, ok := ts.manager.Positions[deal.ProfileID][deal.DealID]; !ok {
		logrus.Errorf("value in map Positions with key dealID: %s is not exist", deal.DealID.String())
	}
	ts.manager.Mu.RUnlock()
	stopLoss, takeProfit := deal.StopLoss.Mul(deal.SharesCount), deal.TakeProfit.Mul(deal.SharesCount)
	strategy, err := ts.identifyStrategy(deal.TakeProfit, deal.StopLoss)
	if err != nil {
		logrus.Errorf("identifyStrategy %v", err)
	}
	if strategy == short {
		stopLoss, takeProfit = takeProfit, stopLoss
	}
	for {
		select {
		case <-ctx.Done():
			return
		default:
			ts.manager.Mu.RLock()
			share := ts.manager.PricesMap[deal.Company]
			if _, ok := ts.manager.Positions[deal.ProfileID][deal.DealID]; !ok {
				ts.manager.Mu.RUnlock()
				continue
			}
			ts.manager.Mu.RUnlock()
			// if strategy == long {
			// 	fmt.Println("(Long) Deal ID: ", deal.DealID, " Profit :", share.Mul(deal.SharesCount).
			// 		Sub(deal.PurchasePrice.Mul(deal.SharesCount)))
			// }
			// if strategy == short {
			// 	fmt.Println("(Short) Deal ID: ", deal.DealID, " Profit :", deal.PurchasePrice.Mul(deal.SharesCount).
			// 		Sub(share.Mul(deal.SharesCount)))
			// }
			if stopLoss.GreaterThanOrEqual(share.Mul(deal.SharesCount)) || share.Mul(deal.SharesCount).GreaterThanOrEqual(takeProfit) {
				profit, err := ts.closePosition(ctx, deal.DealID, deal.ProfileID, share)
				if err != nil {
					logrus.Errorf("closePosition %v", err)
				}
				fmt.Println("Position closed with profit: ", profit)
				return
			}
		}
	}
}

// closePosition is a method that calls method of Repository and returns profit
func (ts *TradingService) closePosition(ctx context.Context, dealid, profileid uuid.UUID, sharePrice decimal.Decimal) (decimal.Decimal, error) {
	var balance model.Balance
	balance.ProfileID = profileid
	ts.manager.Mu.RLock()
	if _, ok := ts.manager.Positions[profileid]; !ok {
		return decimal.Zero, fmt.Errorf("key of map Positions: %s is not exist", profileid.String())
	}
	if _, ok := ts.manager.Positions[profileid][dealid]; !ok {
		return decimal.Zero, fmt.Errorf("value in map Positions with key dealid: %s is not exist", dealid.String())
	}
	deal := ts.manager.Positions[profileid][dealid]
	ts.manager.Mu.RUnlock()
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
		return decimal.Zero, fmt.Errorf("closePosition %w", err)
	}
	balance.Operation = deal.Profit.Add(deal.PurchasePrice.Mul(deal.SharesCount))
	_, err = ts.bRep.BalanceOperation(ctx, &balance)
	if err != nil {
		return decimal.Zero, fmt.Errorf("balanceOperation %w", err)
	}
	return deal.Profit, nil
}

// ClosePositionManually is a method that calls from user before closing by takeprofit/stoploss
func (ts *TradingService) ClosePositionManually(ctx context.Context, dealid, profileid uuid.UUID) (decimal.Decimal, error) {
	ts.manager.Mu.RLock()
	deal := ts.manager.Positions[profileid][dealid]
	ts.manager.Mu.RUnlock()
	for {
		select {
		case <-ctx.Done():
			return decimal.Zero, nil
		default:
			ts.manager.Mu.RLock()
			share := ts.manager.PricesMap[deal.Company]
			ts.manager.Mu.RUnlock()
			profit, err := ts.closePosition(ctx, dealid, profileid, share)
			if err != nil {
				return decimal.Zero, fmt.Errorf("closePosition %w", err)
			}
			return profit, nil
		}
	}
}

// Subscribe is a method of TradingService that calls method of Repository as goroutine
func (ts *TradingService) Subscribe(ctx context.Context) {
	companyShares := strings.Split(ts.cfg.CompanyShares, ",")
	submanager := make(chan model.Share, len(companyShares))
	go func() {
		err := ts.priceRep.Subscribe(ctx, submanager)
		if err != nil {
			logrus.Errorf("subscribe %v", err)
		}
	}()
	for {
		select
		{
		case <-ctx.Done():
			return
		case share := <-submanager:
			ts.manager.Mu.Lock()
			ts.manager.PricesMap[share.Company] = share.Price
			ts.manager.Mu.Unlock()
		}
	}
}

// GetPrices is method that read actual prices of shares
func (ts *TradingService) GetPrices() ([]model.Share, error) {
	var shares []model.Share
	companyShares := strings.Split(ts.cfg.CompanyShares, ",")
	ts.manager.Mu.RLock()
	defer ts.manager.Mu.RUnlock()
	for _, company := range companyShares {
		if price, exists := ts.manager.PricesMap[company]; exists {
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
		return nil, fmt.Errorf("getUnclosedPositions %w", err)
	}
	return unclosedDeals, nil
}

// GetClosedPositions is a method of TradingService calls method of Repository
func (ts *TradingService) GetClosedPositions(ctx context.Context, profileid uuid.UUID) ([]*model.Deal, error) {
	closedDeals, err := ts.priceRep.GetClosedPositions(ctx, profileid)
	if err != nil {
		return nil, fmt.Errorf("getClosedPositions %w", err)
	}
	return closedDeals, nil
}

// BackupUnclosedPositions is a method of TradingService calls method of Repository
func (ts *TradingService) BackupUnclosedPositions(ctx context.Context) {
	unclosedDeals, err := ts.priceRep.GetUnclosedPositionsForAll(ctx)
	if err != nil {
		logrus.Errorf("getUnclosedPositionsForAll %v", err)
	}
	for _, deal := range unclosedDeals {
		ts.manager.Mu.Lock()
		ts.manager.Positions[deal.ProfileID] = make(map[uuid.UUID]model.Deal)
		ts.manager.Positions[deal.ProfileID][deal.DealID] = *deal
		ts.manager.Mu.Unlock()
		go ts.getProfit(ctx, *deal)
	}
}

// BalanceOperation is a method of TradingService calls method of Repository
func (ts *TradingService) BalanceOperation(ctx context.Context, balance *model.Balance) (float64, error) {
	operation, err := ts.bRep.BalanceOperation(ctx, balance)
	if err != nil {
		return 0, fmt.Errorf("balanceOperation %w", err)
	}
	return operation, nil
}

// GetBalance is a method of TradingService calls method of Repository
func (ts *TradingService) GetBalance(ctx context.Context, profileid uuid.UUID) (float64, error) {
	money, err := ts.bRep.GetBalance(ctx, profileid)
	if err != nil {
		return 0, fmt.Errorf("getBalance %w", err)
	}
	return money, nil
}
