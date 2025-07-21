package bybit

// Base URLs
const (
	// Production endpoints
	BaseURLMainnet = "https://api.bybit.com"
	BaseURLBytick  = "https://api.bytick.com"

	// Testnet endpoints
	BaseURLTestnet = "https://api-testnet.bybit.com"
)

// API Endpoints
const (
	// General endpoints
	EndpointServerTime   = "/v5/market/time"
	EndpointExchangeInfo = "/v5/market/instruments-info"
	EndpointPing         = "/v5/market/time"

	// Market data endpoints
	EndpointTicker24hr  = "/v5/market/tickers"
	EndpointTickerPrice = "/v5/market/tickers"
	EndpointOrderBook   = "/v5/market/orderbook"
	EndpointKlines      = "/v5/market/kline"
	EndpointTrades      = "/v5/market/recent-trade"

	// Account endpoints (signed)
	EndpointAccount    = "/v5/account/wallet-balance"
	EndpointPosition   = "/v5/position/list"
	EndpointOrder      = "/v5/order/realtime"
	EndpointOrders     = "/v5/order/history"
	EndpointOpenOrders = "/v5/order/realtime"

	// Trading endpoints (signed)
	EndpointNewOrder    = "/v5/order/create"
	EndpointCancelOrder = "/v5/order/cancel"
	EndpointOrderStatus = "/v5/order/realtime"
	EndpointMyTrades    = "/v5/execution/list"

	// User data stream endpoints
	EndpointUserDataStream = "/v5/user/query-api"
)

// HTTP Headers
const (
	HeaderAPIKey     = "X-BAPI-API-KEY"
	HeaderSignature  = "X-BAPI-SIGN"
	HeaderTimestamp  = "X-BAPI-TIMESTAMP"
	HeaderRecvWindow = "X-BAPI-RECV-WINDOW"
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
	SideBuy  = "Buy"
	SideSell = "Sell"
)

// Order Types
const (
	OrderTypeLimit           = "Limit"
	OrderTypeMarket          = "Market"
	OrderTypeStopLoss        = "StopLoss"
	OrderTypeStopLossLimit   = "StopLossLimit"
	OrderTypeTakeProfit      = "TakeProfit"
	OrderTypeTakeProfitLimit = "TakeProfitLimit"
	OrderTypeLimitMaker      = "LimitMaker"
)

// Time in Force
const (
	TimeInForceGTC = "GTC" // Good Till Cancel
	TimeInForceIOC = "IOC" // Immediate or Cancel
	TimeInForceFOK = "FOK" // Fill or Kill
)

// Order Status
const (
	OrderStatusCreated         = "Created"
	OrderStatusNew             = "New"
	OrderStatusPartiallyFilled = "PartiallyFilled"
	OrderStatusFilled          = "Filled"
	OrderStatusCancelled       = "Cancelled"
	OrderStatusPendingCancel   = "PendingCancel"
	OrderStatusRejected        = "Rejected"
	OrderStatusExpired         = "Expired"
)

// Order Filter (UTA 2.0)
const (
	OrderFilterActive      = "ActiveOrder"
	OrderFilterStop        = "StopOrder"
	OrderFilterTpsl        = "TpslOrder"
	OrderFilterOco         = "OcoOrder"
	OrderFilterConditional = "ConditionalOrder"
)

// Place Type (UTA 2.0)
const (
	PlaceTypeIvRequest = "iv_request"
	PlaceTypeUta       = "uta"
)

// Position Side
const (
	PositionSideLong  = "Long"
	PositionSideShort = "Short"
)

// Category
const (
	CategorySpot    = "spot"
	CategoryLinear  = "linear"
	CategoryInverse = "inverse"
	CategoryOption  = "option"
)

// Kline Intervals
const (
	Interval1m  = "1"
	Interval3m  = "3"
	Interval5m  = "5"
	Interval15m = "15"
	Interval30m = "30"
	Interval1h  = "60"
	Interval2h  = "120"
	Interval4h  = "240"
	Interval6h  = "360"
	Interval8h  = "480"
	Interval12h = "720"
	Interval1d  = "D"
	Interval1w  = "W"
	Interval1M  = "M"
)

// WebSocket URLs
const (
	// Production WebSocket endpoints
	WSBaseURLSpot    = "wss://stream.bybit.com/v5/public/spot"
	WSBaseURLLinear  = "wss://stream.bybit.com/v5/public/linear"
	WSBaseURLPrivate = "wss://stream.bybit.com/v5/private"

	// Testnet WebSocket endpoints
	WSBaseURLSpotTestnet    = "wss://stream-testnet.bybit.com/v5/public/spot"
	WSBaseURLLinearTestnet  = "wss://stream-testnet.bybit.com/v5/public/linear"
	WSBaseURLPrivateTestnet = "wss://stream-testnet.bybit.com/v5/private"
)

// WebSocket Stream Names
const (
	WSStreamKline       = "kline"
	WSStreamTicker      = "tickers"
	WSStreamOrderBook   = "orderbook"
	WSStreamTrade       = "publicTrade"
	WSStreamLiquidation = "liquidation"
)

// WebSocket Methods
const (
	WSMethodSubscribe         = "subscribe"
	WSMethodUnsubscribe       = "unsubscribe"
	WSMethodListSubscriptions = "list_subscriptions"
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
	WSEventOrder     = "order"
	WSEventExecution = "execution"
	WSEventPosition  = "position"
	WSEventWallet    = "wallet"
)
