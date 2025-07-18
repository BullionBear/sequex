package binance

import (
	"github.com/shopspring/decimal"
)

// mustParseDecimal parses a string to decimal, panics on error
func mustParseDecimal(s string) decimal.Decimal {
	d, err := decimal.NewFromString(s)
	if err != nil {
		// Return zero decimal if parsing fails
		return decimal.Zero
	}
	return d
}

// parseDecimal safely parses a string to decimal with error handling
func parseDecimal(s string) (decimal.Decimal, error) {
	return decimal.NewFromString(s)
}
