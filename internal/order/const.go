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
	SideUnknown Side = iota // Unknown
	SideBuy                 // BUY
	SideSell                // SELL
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

type OrderStatus int

const (
	OrderStatusUnknown         OrderStatus = iota // Unknown
	OrderStatusLocalPending                       // LOCAL_PENDING
	OrderStatusNew                                // NEW
	OrderStatusPartiallyFilled                    // PARTIALLY_FILLED
	OrderStatusFilled                             // FILLED
	OrderStatusPendingCancel                      // PENDING_CANCEL
	OrderStatusCanceled                           // CANCELED
	OrderStatusRejected                           // REJECTED
	OrderStatusExpired                            // EXPIRED
)

func (s OrderStatus) String() string {
	switch s {
	case OrderStatusLocalPending:
		return "LOCAL_PENDING"
	case OrderStatusNew:
		return "NEW"
	case OrderStatusPartiallyFilled:
		return "PARTIALLY_FILLED"
	case OrderStatusFilled:
		return "FILLED"
	case OrderStatusCanceled:
		return "CANCELED"
	case OrderStatusRejected:
		return "REJECTED"
	case OrderStatusExpired:
		return "EXPIRED"
	default:
		return fmt.Sprintf("Unknown OrderStatus (%d)", s)
	}
}

// OrderType represents the type of an order
type OrderType int

const (
	// Simple order
	OrderTypeLimit              OrderType = iota // LIMIT
	OrderTypeMarket                              // MARKET
	OrderTypeStopMarket                          // STOP_MARKET
	OrderTypeLimitMaker                          // LIMIT_MAKER
	OrderTypeTrailingStopMarket                  // TRAILING_STOP_MARKET

	// Complex order
	OrderTypeOCO  // OCO (One Cancels Other)
	OrderTypeIFDO // IFDO (If Done Order)
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

// TimeInForce represents the time in force for an order
type TimeInForce int

const (
	TimeInForceUnknown TimeInForce = iota // Unknown
	TimeInForceGTC                        // GTC (Good Till Cancelled)
	TimeInForceIOC                        // IOC (Immediate Or Cancel)
	TimeInForceFOK                        // FOK (Fill Or Kill)
)

func (t TimeInForce) String() string {
	switch t {
	case TimeInForceGTC:
		return "GTC"
	case TimeInForceIOC:
		return "IOC"
	case TimeInForceFOK:
		return "FOK"
	default:
		return fmt.Sprintf("Unknown TimeInForce (%d)", t)
	}
}
