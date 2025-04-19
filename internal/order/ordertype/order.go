package ordertype

import "github.com/shopspring/decimal"

type Order interface {
	GetType() OrderType
}

var _ Order = (*MarketOrder)(nil)
var _ Order = (*LimitOrder)(nil)
var _ Order = (*StopMarketOrder)(nil)

type MarketOrder struct {
	ID         string          `json:"id"`
	Instrument Instrument      `json:"instrument"`
	Symbol     string          `json:"symbol"`
	Side       Side            `json:"side"`
	Quantity   decimal.Decimal `json:"quantity"`
}

func (m MarketOrder) GetType() OrderType {
	return OrderTypeMarket
}

type LimitOrder struct {
	ID          string          `json:"id"`
	Instrument  Instrument      `json:"instrument"`
	Symbol      string          `json:"symbol"`
	Side        Side            `json:"side"`
	Quantity    decimal.Decimal `json:"quantity"`
	Price       decimal.Decimal `json:"price"`
	TimeInForce TimeInForce     `json:"time_in_force"`
}

func (l LimitOrder) GetType() OrderType {
	return OrderTypeLimit
}

type StopMarketOrder struct {
	Instrument Instrument      `json:"instrument"`
	Symbol     string          `json:"symbol"`
	Side       Side            `json:"side"`
	Quantity   decimal.Decimal `json:"quantity"`
	StopPrice  decimal.Decimal `json:"stop_price"`
}

func (s StopMarketOrder) GetType() OrderType {
	return OrderTypeStopMarket
}

type OneCancelsOtherOrder struct {
	Orders []Order `json:"orders"`
}

func (o *OneCancelsOtherOrder) GetType() OrderType {
	return OrderTypeOCO
}

type IfDoneOrder struct {
	Orders            []Order `json:"orders"`
	currentOrderIndex int
}

func (o *IfDoneOrder) GetType() OrderType {
	return OrderTypeIFDO
}

func (o *IfDoneOrder) GetCurrentOrder() Order {
	if o.currentOrderIndex >= len(o.Orders) {
		return nil
	}
	return o.Orders[o.currentOrderIndex]
}

func (o *IfDoneOrder) ToNext() int {
	o.currentOrderIndex++
	if o.currentOrderIndex >= len(o.Orders) {
		o.currentOrderIndex = -1
	}
	return o.currentOrderIndex
}
