// Package model contains models of using entities
package model

import (
	"time"

	"github.com/google/uuid"
	"github.com/shopspring/decimal"
)

// Action is a struct for actions entity
type Action struct {
	Company string
	Price   decimal.Decimal
}

// Deal is a struct for creating new deals
type Deal struct {
	DealID        uuid.UUID
	ActionsCount  decimal.Decimal `json:"actionscount"`
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
	BalanceID uuid.UUID       `json:"balanceid" validate:"required,uuid"`                  // id of balance operation - each operation have new id
	ProfileID uuid.UUID       `json:"profileid" validate:"required,uuid"`                  // same value as ID in struct User
	Operation decimal.Decimal `json:"operation" validate:"required,gt=0" form:"operation"` // sum of money to be deposit or withdraw
}
