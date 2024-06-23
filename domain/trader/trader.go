package trader

import "github.com/shopspring/decimal"

type Trader interface {
	// Trade is the main function to make trading decisions.
	CreateMarketOrder(symbol string, quoteQty decimal.Decimal) error
}
