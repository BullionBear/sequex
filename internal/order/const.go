package order

import "fmt"

// Instrument represents the type of financial instrument
type Instrument int

const (
	InstrumentSpot Instrument = iota // SPOT
	InstrumentPerp                   // PERP
)

func (i Instrument) String() string {
	switch i {
	case InstrumentSpot:
		return "SPOT"
	case InstrumentPerp:
		return "PERP"
	default:
		return fmt.Sprintf("Unknown Instrument (%d)", i)
	}
}

// Side represents the side of an order (buy or sell)
type Side int

const (
	SideBuy  Side = iota // BUY
	SideSell             // SELL
)

func (s Side) String() string {
	switch s {
	case SideBuy:
		return "BUY"
	case SideSell:
		return "SELL"
	default:
		return fmt.Sprintf("Unknown Side (%d)", s)
	}
}

// OrderType represents the type of an order
type OrderType int

const (
	OrderTypeLimit      OrderType = iota // LIMIT
	OrderTypeMarket                      // MARKET
	OrderTypeStopMarket                  // STOP_MARKET
	OrderTypeLimitMaker                  // LIMIT_MAKER
)

func (o OrderType) String() string {
	switch o {
	case OrderTypeLimit:
		return "LIMIT"
	case OrderTypeMarket:
		return "MARKET"
	case OrderTypeStopMarket:
		return "STOP_MARKET"
	case OrderTypeLimitMaker:
		return "LIMIT_MAKER"
	default:
		return fmt.Sprintf("Unknown OrderType (%d)", o)
	}
}
