package model

import "github.com/google/uuid"

// Balance contains an info about the balance and will be written in a balance table
type Balance struct {
	BalanceID uuid.UUID `json:"balanceid" validate:"required,uuid"`                  // id of balance operation - each operation have new id
	ProfileID uuid.UUID `json:"profileid" validate:"required,uuid"`                  // same value as ID in struct User
	Operation float64   `json:"operation" validate:"required,gt=0" form:"operation"` // sum of money to be deposit or withdraw
}
