// Package repository is a lower level of project
package repository

import (
	"context"
	"fmt"
	"strings"

	pproto "github.com/artnikel/PriceService/proto"
	"github.com/artnikel/TradingService/internal/config"
	"github.com/artnikel/TradingService/internal/model"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/shopspring/decimal"
)

// PriceRepository represents the client of Price Service repository implementation.
type PriceRepository struct {
	pool   *pgxpool.Pool
	client pproto.PriceServiceClient
	cfg    config.Variables
}

// NewPriceRepository creates and returns a new instance of PriceRepository, using the provided proto.PriceServiceClient.
func NewPriceRepository(client pproto.PriceServiceClient, pool *pgxpool.Pool, cfg *config.Variables) *PriceRepository {
	return &PriceRepository{
		client: client,
		pool:   pool,
		cfg:    *cfg,
	}
}

// Subscribe call a method of PriceService.
func (p *PriceRepository) Subscribe(ctx context.Context, manager chan model.Share) error {
	companyShares := strings.Split(p.cfg.CompanyShares, ",")
	lastShares := make(map[string]decimal.Decimal)
	stream, err := p.client.Subscribe(ctx, &pproto.SubscribeRequest{
		Uuid:           uuid.NewString(),
		SelectedShares: companyShares,
	})
	if err != nil {
		return fmt.Errorf("subscribe %w", err)
	}
	for {
		select {
		case <-ctx.Done():
			return nil
		default:
			protoResponse, err := stream.Recv()
			if err != nil {
				if ctx.Err() != nil {
					return nil
				}
				return fmt.Errorf("stream.Recv %w", err)
			}
			for _, protoShare := range protoResponse.Shares {
				sharePrice := decimal.NewFromFloat(protoShare.Price)
				lastPrice, exists := lastShares[protoShare.Company]
				if !exists || !sharePrice.Equal(lastPrice) {
					recievedShares := model.Share{
						Company: protoShare.Company,
						Price:   sharePrice,
					}
					manager <- recievedShares
					lastShares[protoShare.Company] = sharePrice
				}
			}
		}
	}
}

// CreatePosition created in database new row with deal
func (p *PriceRepository) CreatePosition(ctx context.Context, deal *model.Deal) error {
	if deal.StopLoss.Cmp(deal.TakeProfit) == 1 {
		deal.Profit = deal.Profit.Neg()
	}
	_, err := p.pool.Exec(ctx, `INSERT INTO deal 
	(dealid, profileid, company, purchaseprice, sharescount, takeprofit, stoploss, dealtime) 
	VALUES ($1, $2, $3, $4, $5, $6, $7, $8)`,
		deal.DealID, deal.ProfileID, deal.Company, deal.PurchasePrice.InexactFloat64(), deal.SharesCount.InexactFloat64(),
		deal.TakeProfit.InexactFloat64(), deal.StopLoss.InexactFloat64(), deal.DealTime)
	if err != nil {
		return fmt.Errorf("exec %w", err)
	}
	return nil
}

// ClosePosition updated and finished row of deal
func (p *PriceRepository) ClosePosition(ctx context.Context, deal *model.Deal) error {
	_, err := p.pool.Exec(ctx, "UPDATE deal SET profit = $1, enddealtime = $2 WHERE dealid = $3",
		deal.Profit.InexactFloat64(), deal.EndDealTime, deal.DealID)
	if err != nil {
		return fmt.Errorf("exec %w", err)
	}
	return nil
}

// GetUnclosedPositions returns info about positions which not closed
func (p *PriceRepository) GetUnclosedPositions(ctx context.Context, profileid uuid.UUID) ([]*model.Deal, error) {
	var deals []*model.Deal
	rows, err := p.pool.Query(ctx, `SELECT dealid, company, purchaseprice, sharescount, takeprofit, stoploss, dealtime
	FROM deal WHERE profileid = $1 AND enddealtime IS NULL`, profileid)
	if err != nil {
		return nil, fmt.Errorf("query %w", err)
	}
	defer rows.Close()
	for rows.Next() {
		var deal model.Deal
		err := rows.Scan(&deal.DealID, &deal.Company, &deal.PurchasePrice, &deal.SharesCount, &deal.TakeProfit, &deal.StopLoss, &deal.DealTime)
		if err != nil {
			return nil, fmt.Errorf("scan %w", err)
		}
		deals = append(deals, &deal)
	}
	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("rows.Err %w", err)
	}
	return deals, nil
}

// GetClosedPositions returns info about positions which closed
func (p *PriceRepository) GetClosedPositions(ctx context.Context, profileid uuid.UUID) ([]*model.Deal, error) {
	var deals []*model.Deal
	rows, err := p.pool.Query(ctx, `SELECT dealid, company, purchaseprice, sharescount, takeprofit, stoploss, dealtime, profit, enddealtime
		FROM deal WHERE profileid = $1 AND enddealtime IS NOT NULL`, profileid)
	if err != nil {
		return nil, fmt.Errorf("query %w", err)
	}
	defer rows.Close()
	for rows.Next() {
		var deal model.Deal
		err := rows.Scan(&deal.DealID, &deal.Company, &deal.PurchasePrice, &deal.SharesCount, &deal.TakeProfit,
			&deal.StopLoss, &deal.DealTime, &deal.Profit, &deal.EndDealTime)
		if err != nil {
			return nil, fmt.Errorf("scan %w", err)
		}
		deals = append(deals, &deal)
	}
	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("rows.Err %w", err)
	}
	return deals, nil
}

// GetUnclosedPositionsForAll returns info about positions which not closed
func (p *PriceRepository) GetUnclosedPositionsForAll(ctx context.Context) ([]*model.Deal, error) {
	var deals []*model.Deal
	rows, err := p.pool.Query(ctx, `SELECT dealid, profileid, company, purchaseprice, sharescount, takeprofit, stoploss, dealtime
	FROM deal WHERE enddealtime IS NULL`)
	if err != nil {
		return nil, fmt.Errorf("query %w", err)
	}
	defer rows.Close()
	for rows.Next() {
		var deal model.Deal
		err := rows.Scan(&deal.DealID, &deal.ProfileID, &deal.Company, &deal.PurchasePrice, &deal.SharesCount, &deal.TakeProfit, &deal.StopLoss, &deal.DealTime)
		if err != nil {
			return nil, fmt.Errorf("scan %w", err)
		}
		deals = append(deals, &deal)
	}
	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("rows.Err %w", err)
	}
	return deals, nil
}
