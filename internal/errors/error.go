// Package errors contains business errors
package errors

const (
	// PurchasePriceOut is error code if takeprofit or stoploss out of limit
	PurchasePriceOut = "PURCHASE_PRICE_OUT"
	// NotEnoughMoney is error code if user don`t have enough money
	NotEnoughMoney = "NOT_ENOUGH_MONEY"
)

// BusinessError is struct for business errors
type BusinessError struct {
	Code string
}

// New is constructor for manage business errors
func New(code string) *BusinessError {
	return &BusinessError{Code: code}
}

// Error is method for creating business errors
func (bs *BusinessError) Error() string {
	return bs.Code
}
