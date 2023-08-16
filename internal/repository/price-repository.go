// Package repository is a lower level of project
package repository

import (
	"context"
	"fmt"

	pproto "github.com/artnikel/PriceService/proto"
	"github.com/artnikel/TradingService/internal/model"
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
	stream, err := p.client.Subscribe(ctx, &pproto.SubscribeRequest{
		Uuid:           uuid.NewString(),
		SelectedShares: []string{"Apple", "Microsoft", "Xerox", "Samsung", "Logitech"},
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
			manager <- model.Share{
				Company: protoShare.Company,
				Price:   decimal.NewFromFloat(protoShare.Price),
			}
		}
	}
}

// AddPosition created in database new row with deal
func (p *PriceRepository) AddPosition(ctx context.Context, strategy string, deal *model.Deal) error {
	if strategy == "short" {
		deal.Profit = deal.Profit.Neg()
	}
	_, err := p.pool.Exec(ctx, `INSERT INTO deal 
	(dealid, profileid, company, purchaseprice, sharescount, stoploss, dealtime) 
	VALUES ($1, $2, $3, $4, $5, $6, $7)`,
		deal.DealID, deal.ProfileID, deal.Company, deal.PurchasePrice.InexactFloat64(),
		deal.SharesCount.InexactFloat64(), deal.StopLoss.InexactFloat64(), deal.DealTime)
	if err != nil {
		return fmt.Errorf("PriceRepository-BalanceOperation: error in method p.pool.Exec(): %w", err)
	}
	return nil
}

// ClosePosition updated and finished row of deal
func (p *PriceRepository) ClosePosition(ctx context.Context, deal *model.Deal) error {
	_, err := p.pool.Exec(ctx, "UDPATE deal SET takeprofit = $1, enddealtime = $2 WHERE dealid = $3",
		deal.TakeProfit.InexactFloat64(), deal.EndDealTime, deal.DealID)
	if err != nil {
		return fmt.Errorf("PriceRepository-BalanceOperation: error in method p.pool.Exec(): %w", err)
	}
	return nil
}
