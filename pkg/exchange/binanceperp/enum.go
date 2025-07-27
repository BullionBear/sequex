package binanceperp

// SymbolType represents the symbol type.
const (
	SymbolTypeFuture = "FUTURE"
)

// ContractType represents the contract type.
const (
	ContractTypePerpetual           = "PERPETUAL"
	ContractTypeCurrentMonth        = "CURRENT_MONTH"
	ContractTypeNextMonth           = "NEXT_MONTH"
	ContractTypeCurrentQuarter      = "CURRENT_QUARTER"
	ContractTypeNextQuarter         = "NEXT_QUARTER"
	ContractTypePerpetualDelivering = "PERPETUAL_DELIVERING"
)

// ContractStatus represents the contract status.
const (
	ContractStatusPendingTrading = "PENDING_TRADING"
	ContractStatusTrading        = "TRADING"
	ContractStatusPreDelivering  = "PRE_DELIVERING"
	ContractStatusDelivering     = "DELIVERING"
	ContractStatusDelivered      = "DELIVERED"
	ContractStatusPreSettle      = "PRE_SETTLE"
	ContractStatusSettling       = "SETTLING"
	ContractStatusClose          = "CLOSE"
)

// OrderStatus represents the status of an order.
const (
	OrderStatusNew             = "NEW"
	OrderStatusPartiallyFilled = "PARTIALLY_FILLED"
	OrderStatusFilled          = "FILLED"
	OrderStatusCanceled        = "CANCELED"
	OrderStatusRejected        = "REJECTED"
	OrderStatusExpired         = "EXPIRED"
	OrderStatusExpiredInMatch  = "EXPIRED_IN_MATCH"
)

// OrderType represents the type of an order.
const (
	OrderTypeLimit              = "LIMIT"
	OrderTypeMarket             = "MARKET"
	OrderTypeStop               = "STOP"
	OrderTypeStopMarket         = "STOP_MARKET"
	OrderTypeTakeProfit         = "TAKE_PROFIT"
	OrderTypeTakeProfitMarket   = "TAKE_PROFIT_MARKET"
	OrderTypeTrailingStopMarket = "TRAILING_STOP_MARKET"
)

// OrderSide represents the side of an order (buy or sell).
const (
	OrderSideBuy  = "BUY"
	OrderSideSell = "SELL"
)

// PositionSide represents the position side.
const (
	PositionSideBoth  = "BOTH"
	PositionSideLong  = "LONG"
	PositionSideShort = "SHORT"
)

// TimeInForce represents how long an order will remain active.
const (
	TimeInForceGTC = "GTC" // Good Till Cancel (GTC order validity is 1 year from placement)
	TimeInForceIOC = "IOC" // Immediate or Cancel
	TimeInForceFOK = "FOK" // Fill or Kill
	TimeInForceGTX = "GTX" // Good Till Crossing (Post Only)
	TimeInForceGTD = "GTD" // Good Till Date
)

// WorkingType represents the working type.
const (
	WorkingTypeMarkPrice     = "MARK_PRICE"
	WorkingTypeContractPrice = "CONTRACT_PRICE"
)

// NewOrderRespType represents the response type for new orders.
const (
	NewOrderRespTypeAck    = "ACK"
	NewOrderRespTypeResult = "RESULT"
)
