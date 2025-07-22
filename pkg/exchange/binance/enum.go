package binance

// OrderStatus represents the status of an order.
const (
	OrderStatusNew             = "NEW"
	OrderStatusPendingNew      = "PENDING_NEW"
	OrderStatusPartiallyFilled = "PARTIALLY_FILLED"
	OrderStatusFilled          = "FILLED"
	OrderStatusCanceled        = "CANCELED"
	OrderStatusPendingCancel   = "PENDING_CANCEL"
	OrderStatusRejected        = "REJECTED"
	OrderStatusExpired         = "EXPIRED"
	OrderStatusExpiredInMatch  = "EXPIRED_IN_MATCH"
)

// OrderType represents the type of an order.
const (
	OrderTypeLimit           = "LIMIT"
	OrderTypeMarket          = "MARKET"
	OrderTypeStopLoss        = "STOP_LOSS"
	OrderTypeStopLossLimit   = "STOP_LOSS_LIMIT"
	OrderTypeTakeProfit      = "TAKE_PROFIT"
	OrderTypeTakeProfitLimit = "TAKE_PROFIT_LIMIT"
	OrderTypeLimitMaker      = "LIMIT_MAKER"
)

// OrderSide represents the side of an order (buy or sell).
const (
	OrderSideBuy  = "BUY"
	OrderSideSell = "SELL"
)

// TimeInForce represents how long an order will remain active.
const (
	TimeInForceGTC = "GTC" // Good Til Canceled
	TimeInForceIOC = "IOC" // Immediate Or Cancel
	TimeInForceFOK = "FOK" // Fill or Kill
)

const (
	NewOrderRespTypeAck    = "ACK"
	NewOrderRespTypeResult = "RESULT"
	NewOrderRespTypeFull   = "FULL"
)
