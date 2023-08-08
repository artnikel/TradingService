module github.com/artnikel/TradingService

go 1.20

// replace (
// 	github.com/artnikel/BalanceService => /home/artyom/gofiles/BalanceService
//  	github.com/artnikel/PriceService => /home/artyom/gofiles/PriceService 
// )

require (
	github.com/artnikel/BalanceService v0.0.0-20230808133307-22305acc1cc9
	github.com/artnikel/PriceService v0.0.0-20230802190226-70778ee0be09
	github.com/caarlos0/env v3.5.0+incompatible
	github.com/golang-jwt/jwt v3.2.2+incompatible
	github.com/google/uuid v1.3.0
)

require (
	github.com/golang/protobuf v1.5.3 // indirect
	golang.org/x/net v0.12.0 // indirect
	golang.org/x/sys v0.10.0 // indirect
	golang.org/x/text v0.11.0 // indirect
	google.golang.org/genproto/googleapis/rpc v0.0.0-20230525234030-28d5490b6b19 // indirect
	google.golang.org/grpc v1.57.0 // indirect
	google.golang.org/protobuf v1.30.0 // indirect
)
