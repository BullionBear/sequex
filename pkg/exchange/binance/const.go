package binance

// Base URLs
const (
	// Production endpoints
	BaseURLSpot = "https://api.binance.com"

	// Testnet endpoints
	BaseURLSpotTestnet = "https://testnet.binance.vision"
)

// API Endpoints
const (
	// General endpoints
	EndpointServerTime   = "/api/v3/time"
	EndpointExchangeInfo = "/api/v3/exchangeInfo"
	EndpointPing         = "/api/v3/ping"

	// Market data endpoints
	EndpointTicker24hr  = "/api/v3/ticker/24hr"
	EndpointTickerPrice = "/api/v3/ticker/price"
	EndpointOrderBook   = "/api/v3/depth"
	EndpointKlines      = "/api/v3/klines"
	EndpointTrades      = "/api/v3/trades"

	// Account endpoints
	EndpointAccount    = "/api/v3/account"
	EndpointOrder      = "/api/v3/order"
	EndpointOrders     = "/api/v3/allOrders"
	EndpointOpenOrders = "/api/v3/openOrders"

	// Trading endpoints
	EndpointNewOrder    = "/api/v3/order"
	EndpointCancelOrder = "/api/v3/order"
	EndpointOrderStatus = "/api/v3/order"
	EndpointMyTrades    = "/api/v3/myTrades"
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
	WSBaseURL = "wss://stream.binance.com:9443"

	// Testnet WebSocket endpoints
	WSBaseURLTestnet = "wss://testnet.binance.vision"
)

// WebSocket Stream Names
const (
	WSStreamKline      = "kline"
	WSStreamTicker     = "ticker"
	WSStreamMiniTicker = "miniTicker"
	WSStreamBookTicker = "bookTicker"
	WSStreamDepth      = "depth"
	WSStreamTrade      = "trade"
	WSStreamAggTrade   = "aggTrade"
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
