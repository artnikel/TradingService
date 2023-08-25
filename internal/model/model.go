// Package model contains models of using entities
package model

import (
	"sync"
	"time"

	"github.com/google/uuid"
	"github.com/shopspring/decimal"
)

// Share is a struct for shares entity
type Share struct {
	Company string
	Price   decimal.Decimal
}

// Deal is a struct for creating new deals
type Deal struct {
	DealID        uuid.UUID
	SharesCount   decimal.Decimal `json:"sharescount"`
	ProfileID     uuid.UUID       `json:"profileid" validate:"required"`
	Company       string          `json:"company" validate:"required"`
	PurchasePrice decimal.Decimal
	StopLoss      decimal.Decimal `json:"stoploss" validate:"required"`
	TakeProfit    decimal.Decimal `json:"takeprofit" validate:"required"`
	DealTime      time.Time
	EndDealTime   time.Time
	Profit        decimal.Decimal
}

// Balance contains an info about the balance and will be written in a balance table
type Balance struct {
	BalanceID uuid.UUID       `json:"balanceid" validate:"required,uuid"`
	ProfileID uuid.UUID       `json:"profileid" validate:"required,uuid"`
	Operation decimal.Decimal `json:"operation" validate:"required,gt=0" form:"operation"`
}

// ChanManager contains custom map to work with shares of subscribers
type ChanManager struct {
	SubscribersShares map[uuid.UUID]map[string]chan Share
	PricesMap         map[string]decimal.Decimal
	Mu                sync.RWMutex
}
