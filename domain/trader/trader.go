package trader

import "github.com/shopspring/decimal"

type Trader interface {
	// Trade is the main function to make trading decisions.
	CreateMarketOrder(symbol string, side bool, quoteQty decimal.Decimal) error // side: true for buy, false for sell
}
