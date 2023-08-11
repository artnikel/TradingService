// Package repository is a lower level of project
package repository

import (
	"context"
	"fmt"

	pproto "github.com/artnikel/PriceService/proto"
	"github.com/artnikel/TradingService/internal/model"
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
func (p *PriceRepository) Subscribe(ctx context.Context, deal *model.Deal, subscribersActions chan []*model.Action) error {
	var selectedActions []string
	selectedActions = append(selectedActions, deal.Company)
	stream, err := p.client.Subscribe(ctx, &pproto.SubscribeRequest{
		Uuid:            deal.ProfileID.String(),
		SelectedActions: selectedActions,
	})
	if err != nil {
		return fmt.Errorf("PriceRepository-Subscribe: error:%w", err)
	}
	defer close(subscribersActions)
	for {
		protoResponse, err := stream.Recv()
		if err != nil {
			if ctx.Err() != nil {
				return nil
			}
			return fmt.Errorf("PriceRepository-Subscribe: error:%w", err)
		}
		recievedActions := make([]*model.Action, 0)
		for _, protoAction := range protoResponse.Actions {
			recievedActions = append(recievedActions, &model.Action{
				Company: protoAction.Company,
				Price:   decimal.NewFromFloat(protoAction.Price),
			})
		}
		select {
		case <-ctx.Done():
			return nil
		case subscribersActions <- recievedActions:
		}
	}
}

// AddDeal created in database new row with deal
func (p *PriceRepository) AddDeal(ctx context.Context, strategy string, deal *model.Deal) error {
	if strategy == "short" {
		deal.Profit = deal.Profit.Neg()
	}
	_, err := p.pool.Exec(ctx, `INSERT INTO deal 
	(dealid, profileid, company, purchaseprice, actionscount, stoploss, takeprofit, dealtime, enddealtime, profit) 
	VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10)`,
		deal.DealID, deal.ProfileID, deal.Company, deal.PurchasePrice.InexactFloat64(), deal.ActionsCount.InexactFloat64(),
		deal.StopLoss.InexactFloat64(), deal.TakeProfit.InexactFloat64(), deal.DealTime, deal.EndDealTime, deal.Profit)
	if err != nil {
		return fmt.Errorf("PriceRepository-BalanceOperation: error in method p.pool.Exec(): %w", err)
	}
	return nil
}
