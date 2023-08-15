// Main package of a project
package main

import (
	"context"
	"log"
	"net"

	bproto "github.com/artnikel/BalanceService/proto"
	pproto "github.com/artnikel/PriceService/proto"
	"github.com/artnikel/TradingService/internal/config"
	"github.com/artnikel/TradingService/internal/handler"
	"github.com/artnikel/TradingService/internal/repository"
	"github.com/artnikel/TradingService/internal/service"
	"github.com/artnikel/TradingService/proto"
	"github.com/caarlos0/env"
	"github.com/go-playground/validator/v10"
	"github.com/jackc/pgx/v5/pgxpool"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

func connectPostgres(cfg *config.Variables) (*pgxpool.Pool, error) {
	cfgPostgres, err := pgxpool.ParseConfig(cfg.PostgresConnTrading)
	if err != nil {
		return nil, err
	}
	dbpool, err := pgxpool.NewWithConfig(context.Background(), cfgPostgres)
	if err != nil {
		return nil, err
	}
	return dbpool, nil
}

// nolint gocritic
func main() {
	pconn, err := grpc.Dial("localhost:8080", grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Fatalf("could not connect: %v", err)
	}
	bconn, err := grpc.Dial("localhost:8095", grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Fatalf("could not connect: %v", err)
	}
	defer func() {
		errConnClose := pconn.Close()
		if err != nil {
			log.Fatalf("could not close connection: %v", errConnClose)
		}
		errConnClose = bconn.Close()
		if err != nil {
			log.Fatalf("could not close connection: %v", errConnClose)
		}
	}()
	var (
		cfg config.Variables
		v   = validator.New()
	)
	if err := env.Parse(&cfg); err != nil {
		log.Fatal("could not parse config: ", err)
	}
	dbpool, errPool := connectPostgres(&cfg)
	if errPool != nil {
		log.Fatal("could not construct the pool: ", errPool)
	}
	defer dbpool.Close()

	//ctx, cancel := context.WithCancel(context.Background())
	pclient := pproto.NewPriceServiceClient(pconn)
	bclient := bproto.NewBalanceServiceClient(bconn)
	prep := repository.NewPriceRepository(pclient, dbpool)
	brep := repository.NewBalanceRepository(bclient)
	tsrv := service.NewTradingService(prep, brep)
	hndl := handler.NewEntityDeal(tsrv, v)
	lis, err := net.Listen("tcp", "localhost:8098")
	if err != nil {
		log.Fatalf("Cannot create listener: %s", err)
	}
	grpcServer := grpc.NewServer()
	proto.RegisterTradingServiceServer(grpcServer, hndl)
	err = grpcServer.Serve(lis)
	if err != nil {
		log.Fatalf("Failed to serve listener: %s", err)
	}
	go tsrv.Subscribe(context.Background())
		
}
