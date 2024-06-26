// Main package of a project
package main

import (
	"context"
	"fmt"
	"log"
	"net"
	"time"

	bproto "github.com/artnikel/BalanceService/proto"
	pproto "github.com/artnikel/PriceService/proto"
	"github.com/artnikel/TradingService/internal/config"
	"github.com/artnikel/TradingService/internal/handler"
	"github.com/artnikel/TradingService/internal/repository"
	"github.com/artnikel/TradingService/internal/service"
	"github.com/artnikel/TradingService/proto"
	"github.com/go-playground/validator/v10"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/lib/pq"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

func connectPostgres(connString string) (*pgxpool.Pool, error) {
	cfgPostgres, err := pgxpool.ParseConfig(connString)
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
	cfg, err := config.New()
	if err != nil {
		log.Fatalf("could not parse config: err: %v", err)
	}
	pconn, err := grpc.Dial(cfg.PriceAddress, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Fatalf("could not connect: %v", err)
	}
	bconn, err := grpc.Dial(cfg.BalanceAddress, grpc.WithTransportCredentials(insecure.NewCredentials()))
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
	v := validator.New()
	dbpool, errPool := connectPostgres(cfg.PostgresConnTrading)
	if errPool != nil {
		log.Fatalf("could not construct the pool: err: %v", errPool)
	}
	defer dbpool.Close()
	pclient := pproto.NewPriceServiceClient(pconn)
	bclient := bproto.NewBalanceServiceClient(bconn)
	prep := repository.NewPriceRepository(pclient, dbpool, cfg)
	brep := repository.NewBalanceRepository(bclient)
	tsrv := service.NewTradingService(prep, brep, cfg, )
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	go tsrv.Subscribe(ctx)
	go tsrv.BackupUnclosedPositions(ctx)
	reportProblem := func(ev pq.ListenerEventType, err error) {
		if err != nil {
			fmt.Println(err.Error())
		}
	}
	listener := pq.NewListener(cfg.PostgresConnTrading+"?sslmode=disable", 10*time.Second, time.Minute, reportProblem)
	err = listener.Listen("events")
	if err != nil {
		log.Fatalf("failed to listen events: %v", err)
	}
	fmt.Println("Trading Service started")

	go tsrv.WaitForNotification(ctx, listener)

	hndl := handler.NewEntityDeal(tsrv, v)
	lis, err := net.Listen("tcp", cfg.TradingAddress)
	if err != nil {
		log.Fatalf("cannot create listener: %v", err)
	}
	grpcServer := grpc.NewServer()
	proto.RegisterTradingServiceServer(grpcServer, hndl)
	err = grpcServer.Serve(lis)
	if err != nil {
		log.Fatalf("failed to serve listener: %v", err)
	}
}
