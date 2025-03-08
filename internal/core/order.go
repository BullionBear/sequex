package core

import "github.com/shopspring/decimal"

type Order struct {
	OrderID   string          `json:"order_id"`
	Symbol    string          `json:"symbol"`
	Price     decimal.Decimal `json:"price"`
	Size      decimal.Decimal `json:"size"`
	Side      string          `json:"side"`
	CreatedAt int64           `json:"created_at"`
}
