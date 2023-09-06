package repository

import (
	"context"
	"fmt"
	"os"
	"os/exec"
	"strings"
	"testing"
	"time"

	sproto "github.com/artnikel/PriceService/proto"
	"github.com/artnikel/PriceService/proto/mocks"
	"github.com/artnikel/TradingService/internal/config"
	"github.com/artnikel/TradingService/internal/model"
	"github.com/caarlos0/env"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/ory/dockertest"
	"github.com/shopspring/decimal"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
)

var (
	pg       *PriceRepository
	testDeal = &model.Deal{
		DealID:        uuid.New(),
		SharesCount:   decimal.NewFromFloat(1.5),
		ProfileID:     uuid.New(),
		Company:       "Apple",
		PurchasePrice: decimal.NewFromFloat(157.8),
		StopLoss:      decimal.NewFromFloat(100),
		TakeProfit:    decimal.NewFromFloat(200),
		DealTime:      time.Now().UTC(),
	}
	closeDeal = &model.Deal{
		Profit:      decimal.NewFromFloat(42.2),
		EndDealTime: time.Now().UTC(),
	}
	cfg config.Variables
)

func SetupTestPostgres() (*pgxpool.Pool, func(), error) {
	pool, err := dockertest.NewPool("")
	if err != nil {
		return nil, nil, fmt.Errorf("could not construct pool: %w", err)
	}
	resource, err := pool.Run("postgres", "latest", []string{
		"POSTGRES_USER=dealuser",
		"POSTGRES_PASSWORD=dealpassword",
		"POSTGRES_DB=dealdb"})
	if err != nil {
		return nil, nil, fmt.Errorf("could not start resource: %w", err)
	}
	err = RunMigrations(resource.GetPort("5432/tcp"))
	if err != nil {
		return nil, nil, fmt.Errorf("failed to run migrations: %w", err)
	}
	dbURL := fmt.Sprintf("postgres://dealuser:dealpassword@localhost:%s/dealdb", resource.GetPort("5432/tcp"))
	cfg, err := pgxpool.ParseConfig(dbURL)
	if err != nil {
		return nil, nil, fmt.Errorf("failed to parse dbURL: %w", err)
	}
	dbpool, err := pgxpool.NewWithConfig(context.Background(), cfg)
	if err != nil {
		return nil, nil, fmt.Errorf("failed to connect pgxpool: %w", err)
	}
	cleanup := func() {
		dbpool.Close()
		pool.Purge(resource)
	}
	return dbpool, cleanup, nil
}

func RunMigrations(port string) error {
	cmd := exec.Command("flyway", "-url=jdbc:postgresql://localhost:"+port+"/dealdb", "-user=dealuser", "-password=dealpassword", "-locations=filesystem:../../migrations", "-connectRetries=10", "migrate")
	err := cmd.Run()
	if err != nil {
		return fmt.Errorf("failed to run migrations: %w", err)
	}
	return nil
}

func TestMain(m *testing.M) {
	dbpool, cleanupPostgres, err := SetupTestPostgres()
	if err != nil {
		fmt.Println("Could not construct the pool: ", err)
		cleanupPostgres()
		os.Exit(1)
	}
	pg = NewPriceRepository(nil, dbpool, cfg)
	exitVal := m.Run()
	cleanupPostgres()
	os.Exit(exitVal)
}

func TestSubscribe(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	mockClient := new(mocks.PriceServiceClient)

	mockStream := new(mocks.PriceService_SubscribeClient)

	expectedShares := []*sproto.Shares{
		{Company: "Apple", Price: 100},
		{Company: "Xerox", Price: 200},
	}
	mockStream.On("Recv").Return(&sproto.SubscribeResponse{Shares: expectedShares}, nil)

	mockClient.On("Subscribe", mock.Anything, mock.Anything).Return(mockStream, nil)

	var cfg config.Variables
	err := env.Parse(&cfg)
	require.NoError(t, err)

	lenOfReadableShares := len(strings.Split(cfg.CompanyShares, ","))

	r := NewPriceRepository(mockClient, nil, cfg)

	subscribersShares := make(chan model.Share, lenOfReadableShares)

	go r.Subscribe(ctx, subscribersShares)

	select {
	case firstShare := <-subscribersShares:
		require.Equal(t, len(expectedShares), len(subscribersShares)+1)

		require.Equal(t, firstShare.Company, expectedShares[0].Company)
		require.Equal(t, firstShare.Price.InexactFloat64(), expectedShares[0].Price)

		secondShare := <-subscribersShares
		require.Equal(t, secondShare.Company, expectedShares[1].Company)
		require.Equal(t, secondShare.Price.InexactFloat64(), expectedShares[1].Price)

	case <-time.After(time.Second * 5):
		t.Error("Timed out while waiting for shares")
	}
	mockClient.AssertExpectations(t)
	mockStream.AssertExpectations(t)
}

func TestAddGetPosition(t *testing.T) {
	var tempDeal *model.Deal
	err := pg.CreatePosition(context.Background(), testDeal)
	require.NoError(t, err)
	unclosedDeals, err := pg.GetUnclosedPositions(context.Background(), testDeal.ProfileID)
	require.NoError(t, err)
	require.Equal(t, len(unclosedDeals), 1)
	for _, unclosedDeal := range unclosedDeals {
		tempDeal = &model.Deal{
			DealID:        unclosedDeal.DealID,
			Company:       unclosedDeal.Company,
			PurchasePrice: unclosedDeal.PurchasePrice,
			SharesCount:   unclosedDeal.SharesCount,
			StopLoss:      unclosedDeal.StopLoss,
			TakeProfit:    unclosedDeal.TakeProfit,
			DealTime:      unclosedDeal.DealTime,
		}
	}
	require.Equal(t, tempDeal.DealID, testDeal.DealID)
	require.Equal(t, tempDeal.Company, testDeal.Company)
	require.Equal(t, tempDeal.PurchasePrice, testDeal.PurchasePrice)
	require.Equal(t, tempDeal.SharesCount, testDeal.SharesCount)
	require.Equal(t, tempDeal.StopLoss, testDeal.StopLoss)
	require.Equal(t, tempDeal.TakeProfit, testDeal.TakeProfit)
}

func TestClosePosition(t *testing.T) {
	testDeal.DealID = uuid.New()
	testDeal.ProfileID = uuid.New()
	err := pg.CreatePosition(context.Background(), testDeal)
	require.NoError(t, err)
	closeDeal.DealID = testDeal.DealID
	err = pg.ClosePosition(context.Background(), closeDeal)
	require.NoError(t, err)
	unclosedDeals, err := pg.GetUnclosedPositions(context.Background(), testDeal.ProfileID)
	require.NoError(t, err)
	require.Equal(t, len(unclosedDeals), 0)
}

func TestGetUnclosedPositions(t *testing.T) {
	testDeal.DealID = uuid.New()
	err := pg.CreatePosition(context.Background(), testDeal)
	require.NoError(t, err)
	unclosedDeals, err := pg.GetUnclosedPositions(context.Background(), testDeal.ProfileID)
	require.NoError(t, err)
	countUnclosedBefore := len(unclosedDeals)
	closeDeal.DealID = testDeal.DealID
	err = pg.ClosePosition(context.Background(), closeDeal)
	require.NoError(t, err)
	unclosedDeals, err = pg.GetUnclosedPositions(context.Background(), testDeal.ProfileID)
	require.NoError(t, err)
	countUnclosedAfter := len(unclosedDeals)
	require.Equal(t, countUnclosedBefore-countUnclosedAfter, 1)
}
