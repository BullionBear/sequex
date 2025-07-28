package exchange

type MarketType string

const (
	MarketTypeSpot MarketType = "spot"
	MarketTypePerp MarketType = "perp"
)

type WalletType string

const (
	WalletTypeSpot   WalletType = "spot"
	WalletTypeMargin WalletType = "margin"
	WalletTypePerp   WalletType = "perp"
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
	TimeInForceGTC TimeInForce = "GTC"
	TimeInForceIOC TimeInForce = "IOC"
	TimeInForceFOK TimeInForce = "FOK"
)

type OrderType string

const (
	OrderTypeLimit      OrderType = "LIMIT"
	OrderTypeMarket     OrderType = "MARKET"
	OrderTypeLimitMaker OrderType = "LIMIT_MAKER"
	OrderTypeStopMarket OrderType = "STOP_MARKET"
)
