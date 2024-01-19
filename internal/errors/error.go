package errors

const (
	PurchasePriceOut = "PURCHASE_PRICE_OUT"
	NotEnoughMoney   = "NOT_ENOUGH_MONEY"
)

type BusinessError struct {
	Code string
}

func New(code string) *BusinessError {
	return &BusinessError{Code: code}
}

func (bs *BusinessError) Error() string {
	return bs.Code
}
