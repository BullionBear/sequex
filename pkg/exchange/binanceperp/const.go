package binanceperp

// Mainnet REST API base URL
const MainnetBaseUrl = "https://fapi.binance.com"

// Mainnet WebSocket base URL
const MainnetWSBaseUrl = "wss://fstream.binance.com"

// Testnet REST API base URL
const TestnetBaseUrl = "https://testnet.binancefuture.com"

// Testnet WebSocket base URL
const TestnetWSBaseUrl = "wss://fstream.binancefuture.com"

// Paths
const (
	PathGetServerTime         = "/fapi/v1/time"
	PathGetDepth              = "/fapi/v1/depth"
	PathGetRecentTrades       = "/fapi/v1/trades"
	PathGetAggTrades          = "/fapi/v1/aggTrades"
	PathGetKlines             = "/fapi/v1/klines"
	PathGetMarkPrice          = "/fapi/v1/premiumIndex"
	PathGetPriceTicker        = "/fapi/v2/ticker/price"
	PathGetBookTicker         = "/fapi/v1/ticker/bookTicker"
	PathGetAccountBalance     = "/fapi/v3/balance"
	PathCreateOrder           = "/fapi/v1/order"
	PathCancelOrder           = "/fapi/v1/order"
	PathQueryOrder            = "/fapi/v1/order"
	PathQueryCurrentOpenOrder = "/fapi/v1/openOrder"
	PathGetMyTrades           = "/fapi/v1/userTrades"
	PathGetPositions          = "/fapi/v2/positionRisk"
	PathCancelAllOrders       = "/fapi/v1/allOpenOrders"
)
