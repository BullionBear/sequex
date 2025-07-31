package exchange

type MarketType string

const (
	MarketTypeBinance     MarketType = "binance"
	MarketTypeBinancePerp MarketType = "binance_perp"
	MarketTypeBybit       MarketType = "bybit"
)

type WalletType string

const (
	WalletTypeSpot    WalletType = "spot"
	WalletTypeMargin  WalletType = "margin"
	WalletTypePerp    WalletType = "perp"
	WalletTypeUnified WalletType = "unified"
)

type MarginType string

const (
	MarginTypeCrossed  MarginType = "crossed"
	MarginTypeIsolated MarginType = "isolated"
)

type PositionType string

const (
	PositionTypeHedge  PositionType = "hedge"
	PositionTypeOneWay PositionType = "oneway"
)

type TimeInForce string

const (
	TimeInForceUnknown TimeInForce = "UNKNOWN"
	TimeInForceGTC     TimeInForce = "GTC"
	TimeInForceIOC     TimeInForce = "IOC"
	TimeInForceFOK     TimeInForce = "FOK"
)

type OrderType string

const (
	OrderTypeUnknown    OrderType = "UNKNOWN"
	OrderTypeLimit      OrderType = "LIMIT"
	OrderTypeMarket     OrderType = "MARKET"
	OrderTypeLimitMaker OrderType = "LIMIT_MAKER"
	OrderTypeStopMarket OrderType = "STOP_MARKET"
)

type OrderSide string

const (
	OrderSideUnknown OrderSide = "UNKNOWN"
	OrderSideBuy     OrderSide = "BUY"
	OrderSideSell    OrderSide = "SELL"
)

type OrderStatus string

const (
	OrderStatusUnknown         OrderStatus = "UNKNOWN"
	OrderStatusNew             OrderStatus = "NEW"
	OrderStatusPartiallyFilled OrderStatus = "PARTIALLY_FILLED"
	OrderStatusFilled          OrderStatus = "FILLED"
	OrderStatusCanceled        OrderStatus = "CANCELED"
	OrderStatusRejected        OrderStatus = "REJECTED"
)
