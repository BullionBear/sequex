package binancefuture

// Base URLs
const (
	// Production endpoints
	BaseURLFutures = "https://fapi.binance.com"

	// Testnet endpoints
	BaseURLFuturesTestnet = "https://testnet.binancefuture.com"
)

// API Endpoints
const (
	// General endpoints
	EndpointServerTime   = "/fapi/v1/time"
	EndpointExchangeInfo = "/fapi/v1/exchangeInfo"
	EndpointPing         = "/fapi/v1/ping"

	// Market data endpoints
	EndpointTicker24hr   = "/fapi/v1/ticker/24hr"
	EndpointTickerPrice  = "/fapi/v1/ticker/price"
	EndpointOrderBook    = "/fapi/v1/depth"
	EndpointKlines       = "/fapi/v1/klines"
	EndpointTrades       = "/fapi/v1/trades"
	EndpointMarkPrice    = "/fapi/v1/premiumIndex"
	EndpointFundingRate  = "/fapi/v1/fundingRate"
	EndpointOpenInterest = "/fapi/v1/openInterest"

	// Account endpoints (signed)
	EndpointAccount      = "/fapi/v2/account"
	EndpointPositionRisk = "/fapi/v3/positionRisk"
	EndpointOrder        = "/fapi/v1/order"
	EndpointOrders       = "/fapi/v1/allOrders"
	EndpointOpenOrders   = "/fapi/v1/openOrders"

	// Trading endpoints (signed)
	EndpointNewOrder    = "/fapi/v1/order"
	EndpointCancelOrder = "/fapi/v1/order"
	EndpointOrderStatus = "/fapi/v1/order"
	EndpointMyTrades    = "/fapi/v1/userTrades"

	// User data stream endpoints
	EndpointUserDataStream = "/fapi/v1/listenKey"
)

// HTTP Headers
const (
	HeaderAPIKey     = "X-MBX-APIKEY"
	HeaderSignature  = "signature"
	HeaderTimestamp  = "timestamp"
	HeaderRecvWindow = "recvWindow"
)

// HTTP Methods
const (
	MethodGET    = "GET"
	MethodPOST   = "POST"
	MethodPUT    = "PUT"
	MethodDELETE = "DELETE"
)

// Security Types
const (
	SecurityTypeNone       = "NONE"        // No authentication required
	SecurityTypeTradeKey   = "TRADE"       // API key required
	SecurityTypeMarketData = "MARKET_DATA" // API key required for market data
	SecurityTypeSigned     = "SIGNED"      // Signature required
)

// Order Sides
const (
	SideBuy  = "BUY"
	SideSell = "SELL"
)

// Order Types
const (
	OrderTypeLimit           = "LIMIT"
	OrderTypeMarket          = "MARKET"
	OrderTypeStopLoss        = "STOP_LOSS"
	OrderTypeStopLossLimit   = "STOP_LOSS_LIMIT"
	OrderTypeTakeProfit      = "TAKE_PROFIT"
	OrderTypeTakeProfitLimit = "TAKE_PROFIT_LIMIT"
	OrderTypeLimitMaker      = "LIMIT_MAKER"
)

// Time in Force
const (
	TimeInForceGTC = "GTC" // Good Till Cancel
	TimeInForceIOC = "IOC" // Immediate or Cancel
	TimeInForceFOK = "FOK" // Fill or Kill
	TimeInForceGTX = "GTX" // Good Till Crossing (Post Only)
)

// Order Status
const (
	OrderStatusNew             = "NEW"
	OrderStatusPartiallyFilled = "PARTIALLY_FILLED"
	OrderStatusFilled          = "FILLED"
	OrderStatusCanceled        = "CANCELED"
	OrderStatusPendingCancel   = "PENDING_CANCEL"
	OrderStatusRejected        = "REJECTED"
	OrderStatusExpired         = "EXPIRED"
)

// Position Side
const (
	PositionSideLong  = "LONG"
	PositionSideShort = "SHORT"
)

// Working Type
const (
	WorkingTypeMarkPrice     = "MARK_PRICE"
	WorkingTypeContractPrice = "CONTRACT_PRICE"
)

// New Order Response Type
const (
	NewOrderRespTypeAck    = "ACK"
	NewOrderRespTypeResult = "RESULT"
	NewOrderRespTypeFull   = "FULL"
)

// Kline Intervals
const (
	Interval1m  = "1m"
	Interval3m  = "3m"
	Interval5m  = "5m"
	Interval15m = "15m"
	Interval30m = "30m"
	Interval1h  = "1h"
	Interval2h  = "2h"
	Interval4h  = "4h"
	Interval6h  = "6h"
	Interval8h  = "8h"
	Interval12h = "12h"
	Interval1d  = "1d"
	Interval3d  = "3d"
	Interval1w  = "1w"
	Interval1M  = "1M"
)

// WebSocket URLs
const (
	// Production WebSocket endpoints
	WSBaseURL = "wss://fstream.binance.com"

	// Testnet WebSocket endpoints
	WSBaseURLTestnet = "wss://stream.binancefuture.com"
)

// WebSocket Stream Names
const (
	WSStreamKline       = "kline"
	WSStreamTicker      = "ticker"
	WSStreamMiniTicker  = "miniTicker"
	WSStreamBookTicker  = "bookTicker"
	WSStreamDepth       = "depth"
	WSStreamTrade       = "trade"
	WSStreamAggTrade    = "aggTrade"
	WSStreamMarkPrice   = "markPrice"
	WSStreamFundingRate = "fundingRate"
)

// WebSocket Methods
const (
	WSMethodSubscribe         = "SUBSCRIBE"
	WSMethodUnsubscribe       = "UNSUBSCRIBE"
	WSMethodListSubscriptions = "LIST_SUBSCRIPTIONS"
)

// WebSocket Message Types
const (
	WSMsgTypeSubscribe   = "subscribe"
	WSMsgTypeUnsubscribe = "unsubscribe"
	WSMsgTypeResult      = "result"
	WSMsgTypeError       = "error"
	WSMsgTypeStream      = "stream"
)

// User Data Stream Event Types
const (
	WSEventAccountUpdate   = "outboundAccountPosition"
	WSEventBalanceUpdate   = "balanceUpdate"
	WSEventExecutionReport = "executionReport"
	WSEventListStatus      = "listStatus"
)
