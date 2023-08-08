package repository

import (
	"context"
	"fmt"

	pproto "github.com/artnikel/PriceService/proto"
	"github.com/google/uuid"
)

// PriceRepository represents the client of Price Service repository implementation.
type PriceRepository struct {
	client pproto.PriceServiceClient
}

// NewPriceRepository creates and returns a new instance of PriceRepository, using the provided proto.PriceServiceClient.
func NewPriceRepository(client pproto.PriceServiceClient) *PriceRepository {
	return &PriceRepository{
		client: client,
	}
}

// Subscribe call a method of PriceService.
func (p *PriceRepository) Subscribe(ctx context.Context, profileid uuid.UUID, selectedaction []string) error {
	_, err := p.client.Subscribe(ctx, &pproto.SubscribeRequest{
		Uuid:            profileid.String(),
		SelectedActions: selectedaction,
	})
	if err != nil {
		return fmt.Errorf("PriceRepository-Subscribe: error:%w", err)
	}
	return nil
}
