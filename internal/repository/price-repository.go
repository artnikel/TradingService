// Package repository is a lower level of project
package repository

import (
	"context"
	"fmt"
	"log"
	"strings"

	pproto "github.com/artnikel/PriceService/proto"
	"github.com/artnikel/TradingService/internal/config"
	"github.com/artnikel/TradingService/internal/model"
	"github.com/caarlos0/env"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/shopspring/decimal"
)

// PriceRepository represents the client of Price Service repository implementation.
type PriceRepository struct {
	pool   *pgxpool.Pool
	client pproto.PriceServiceClient
}

// NewPriceRepository creates and returns a new instance of PriceRepository, using the provided proto.PriceServiceClient.
func NewPriceRepository(client pproto.PriceServiceClient, pool *pgxpool.Pool) *PriceRepository {
	return &PriceRepository{
		client: client,
		pool:   pool,
	}
}

// Subscribe call a method of PriceService.
func (p *PriceRepository) Subscribe(ctx context.Context, manager chan model.Share) error {
	var cfg config.Variables
	if err := env.Parse(&cfg); err != nil {
		log.Fatal("could not parse config: ", err)
	}
	companyShares := strings.Split(cfg.CompanyShares, ",")
	stream, err := p.client.Subscribe(ctx, &pproto.SubscribeRequest{
		Uuid:           uuid.NewString(),
		SelectedShares: companyShares,
	})
	if err != nil {
		return fmt.Errorf("PriceRepository-Subscribe: error:%w", err)
	}
	for {
		protoResponse, err := stream.Recv()
		if err != nil {
			if ctx.Err() != nil {
				return nil
			}
			return fmt.Errorf("PriceRepository-Subscribe: error:%w", err)
		}
		for _, protoShare := range protoResponse.Shares {
			recievedShares := model.Share{
				Company: protoShare.Company,
				Price:   decimal.NewFromFloat(protoShare.Price),
			}
			manager <- recievedShares
		}
	}
}

// AddPosition created in database new row with deal
func (p *PriceRepository) AddPosition(ctx context.Context, strategy string, deal *model.Deal) error {
	if strategy == "short" {
		deal.Profit = deal.Profit.Neg()
	}
	_, err := p.pool.Exec(ctx, `INSERT INTO deal 
	(dealid, profileid, company, purchaseprice, sharescount, takeprofit, stoploss, dealtime) 
	VALUES ($1, $2, $3, $4, $5, $6, $7, $8)`,
		deal.DealID, deal.ProfileID, deal.Company, deal.PurchasePrice.InexactFloat64(), deal.SharesCount.InexactFloat64(),
		deal.TakeProfit.InexactFloat64(), deal.StopLoss.InexactFloat64(), deal.DealTime)
	if err != nil {
		return fmt.Errorf("PriceRepository-BalanceOperation: error in method p.pool.Exec(): %w", err)
	}
	return nil
}

// ClosePosition updated and finished row of deal
func (p *PriceRepository) ClosePosition(ctx context.Context, deal *model.Deal) error {
	_, err := p.pool.Exec(ctx, "UPDATE deal SET profit = $1, enddealtime = $2 WHERE dealid = $3",
		deal.Profit.InexactFloat64(), deal.EndDealTime, deal.DealID)
	if err != nil {
		return fmt.Errorf("PriceRepository-BalanceOperation: error in method p.pool.Exec(): %w", err)
	}
	return nil
}

func (p *PriceRepository) GetUnclosedPositions(ctx context.Context, profileid uuid.UUID) ([]*model.Deal, error) {
	var deals []*model.Deal
	rows, err := p.pool.Query(ctx, `SELECT (dealid, company, purchaseprice, sharescount, takeprofit, stoploss, dealtime)
	FROM deal WHERE profileid = $1 AND enddealtime IS NULL`, profileid)
	if err != nil {
		return nil, fmt.Errorf("PriceRepository-GetUnclosedPositions: error in method p.pool.Query(): %w", err)
	}
	defer rows.Close()
	for rows.Next() {
		var deal model.Deal
		err := rows.Scan(&deal.DealID, &deal.Company, &deal.PurchasePrice, &deal.SharesCount, &deal.TakeProfit, &deal.StopLoss, &deal.DealTime)
		if err != nil {
			return nil, fmt.Errorf("PriceRepository-GetUnclosedPositions error in method rows.Scan(): %w", err)
		}
		deals = append(deals, &deal)
	}
	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("PriceRepository-GetUnclosedPositions error iterating rows: %w", err)
	}
	return deals, nil
}
