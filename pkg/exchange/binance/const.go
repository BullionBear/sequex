package binance

// Mainnet REST API base URLs
const (
	MainnetBaseUrl    = "https://api.binance.com/api"
	MainnetBaseGCPUrl = "https://api-gcp.binance.com/api"
	MainnetBaseUrl1   = "https://api1.binance.com/api"
	MainnetBaseUrl2   = "https://api2.binance.com/api"
	MainnetBaseUrl3   = "https://api3.binance.com/api"
	MainnetBaseUrl4   = "https://api4.binance.com/api"
)

// Mainnet WebSocket base URLs
const (
	MainnetWSBaseUrl     = "wss://stream.binance.com/ws"
	MainnetWSBaseUrl9443 = "wss://stream.binance.com:9443/ws"
)

// Testnet REST API base URL
const TestnetBaseUrl = "https://testnet.binance.vision/api"

// Testnet WebSocket base URLs
const (
	TestnetWSBaseUrl     = "wss://stream.testnet.binance.vision/ws"
	TestnetWSBaseUrl9443 = "wss://stream.testnet.binance.vision:9443/ws"
)

// Paths
const (
	PathCreateOrder      = "/v3/order"
	PathGetDepth         = "/v3/depth"
	PathGetRecentTrades  = "/v3/trades"
	PathGetAggTrades     = "/v3/aggTrades"
	PathGetKlines        = "/v3/klines"
	PathGetPriceTicker   = "/v3/ticker/price"
	PathGetExchangeInfo  = "/v3/exchangeInfo"
	PathCancelOrder      = "/v3/order"
	PathCancelAllOrders  = "/v3/openOrders"
	PathQueryOrder       = "/v3/order"
	PathGetAccountInfo   = "/v3/account"
	PathListOpenOrders   = "/v3/openOrders"
	PathGetAccountTrades = "/v3/myTrades"
	PathUserDataStream   = "/v3/userDataStream"
)
