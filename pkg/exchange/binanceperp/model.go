package binanceperp

// Response is the unified response wrapper for all endpoints.
type Response[T any] struct {
	Code    int    `json:"code"`
	Message string `json:"msg"`
	Data    *T     `json:"data,omitempty"`
}

// GetServerTimeResponse represents the server time response.
type GetServerTimeResponse struct {
	ServerTime int64 `json:"serverTime"`
}

// GetDepthRequest defines the parameters for getting order book depth.
type GetDepthRequest struct {
	Symbol string // required
	Limit  int    // optional, default 500; Valid limits:[5, 10, 20, 50, 100, 500, 1000]
}

// GetDepthResponse represents the order book depth response.
type GetDepthResponse struct {
	LastUpdateId int64      `json:"lastUpdateId"`
	E            int64      `json:"E"` // Message output time
	T            int64      `json:"T"` // Transaction time
	Bids         [][]string `json:"bids"`
	Asks         [][]string `json:"asks"`
}

// GetRecentTradesRequest defines the parameters for getting recent trades.
type GetRecentTradesRequest struct {
	Symbol string // required
	Limit  int    // optional, default 500; max 1000
}

// RecentTrade represents a single recent trade.
type RecentTrade struct {
	Id           int64  `json:"id"`
	Price        string `json:"price"`
	Qty          string `json:"qty"`
	QuoteQty     string `json:"quoteQty"`
	Time         int64  `json:"time"`
	IsBuyerMaker bool   `json:"isBuyerMaker"`
}

// GetAggTradesRequest defines the parameters for getting aggregate trades.
type GetAggTradesRequest struct {
	Symbol    string // required
	FromId    int64  // optional, ID to get aggregate trades from INCLUSIVE
	StartTime int64  // optional, timestamp in ms to get aggregate trades from INCLUSIVE
	EndTime   int64  // optional, timestamp in ms to get aggregate trades until INCLUSIVE
	Limit     int    // optional, default 500; max 1000
}

// AggTrade represents a single aggregate trade.
type AggTrade struct {
	AggTradeId   int64  `json:"a"` // Aggregate tradeId
	Price        string `json:"p"` // Price
	Quantity     string `json:"q"` // Quantity
	FirstTradeId int64  `json:"f"` // First tradeId
	LastTradeId  int64  `json:"l"` // Last tradeId
	Timestamp    int64  `json:"T"` // Timestamp
	IsBuyerMaker bool   `json:"m"` // Was the buyer the maker?
}

// GetKlinesRequest defines the parameters for getting kline data.
type GetKlinesRequest struct {
	Symbol    string // required
	Interval  string // required (e.g. "1m", "5m", "1h", "1d")
	StartTime int64  // optional, timestamp in ms
	EndTime   int64  // optional, timestamp in ms
	Limit     int    // optional, default 500; max 1500
}

// Kline represents a single kline/candlestick.
type Kline struct {
	OpenTime                 int64  `json:"openTime"`                 // Open time
	Open                     string `json:"open"`                     // Open price
	High                     string `json:"high"`                     // High price
	Low                      string `json:"low"`                      // Low price
	Close                    string `json:"close"`                    // Close price
	Volume                   string `json:"volume"`                   // Volume
	CloseTime                int64  `json:"closeTime"`                // Close time
	QuoteAssetVolume         string `json:"quoteAssetVolume"`         // Quote asset volume
	NumberOfTrades           int    `json:"numberOfTrades"`           // Number of trades
	TakerBuyBaseAssetVolume  string `json:"takerBuyBaseAssetVolume"`  // Taker buy base asset volume
	TakerBuyQuoteAssetVolume string `json:"takerBuyQuoteAssetVolume"` // Taker buy quote asset volume
	Ignore                   string `json:"ignore"`                   // Ignore
}

// GetMarkPriceRequest defines the parameters for getting mark price and funding rate.
type GetMarkPriceRequest struct {
	Symbol string // optional, if not provided returns all symbols
}

// MarkPrice represents mark price and funding rate data.
type MarkPrice struct {
	Symbol               string `json:"symbol"`               // Symbol
	MarkPrice            string `json:"markPrice"`            // Mark price
	IndexPrice           string `json:"indexPrice"`           // Index price
	EstimatedSettlePrice string `json:"estimatedSettlePrice"` // Estimated Settle Price
	LastFundingRate      string `json:"lastFundingRate"`      // Latest funding rate
	InterestRate         string `json:"interestRate"`         // Interest rate
	NextFundingTime      int64  `json:"nextFundingTime"`      // Next funding time
	Time                 int64  `json:"time"`                 // Timestamp
}

// GetPriceTickerRequest defines the parameters for getting price ticker.
type GetPriceTickerRequest struct {
	Symbol string // optional, if not provided returns all symbols
}

// PriceTicker represents symbol price ticker data.
type PriceTicker struct {
	Symbol string `json:"symbol"` // Symbol
	Price  string `json:"price"`  // Price
	Time   int64  `json:"time"`   // Transaction time
}

// GetBookTickerRequest defines the parameters for getting book ticker.
type GetBookTickerRequest struct {
	Symbol string // optional, if not provided returns all symbols
}

// BookTicker represents symbol order book ticker data (best bid/ask).
type BookTicker struct {
	Symbol   string `json:"symbol"`   // Symbol
	BidPrice string `json:"bidPrice"` // Best bid price
	BidQty   string `json:"bidQty"`   // Best bid quantity
	AskPrice string `json:"askPrice"` // Best ask price
	AskQty   string `json:"askQty"`   // Best ask quantity
	Time     int64  `json:"time"`     // Transaction time
}

// GetAccountBalanceRequest defines the parameters for getting account balance.
type GetAccountBalanceRequest struct {
	RecvWindow int64 // optional, default 5000
}

// AccountBalance represents account balance information for a single asset.
type AccountBalance struct {
	AccountAlias       string `json:"accountAlias"`       // Unique account code
	Asset              string `json:"asset"`              // Asset name
	Balance            string `json:"balance"`            // Wallet balance
	CrossWalletBalance string `json:"crossWalletBalance"` // Crossed wallet balance
	CrossUnPnl         string `json:"crossUnPnl"`         // Unrealized profit of crossed positions
	AvailableBalance   string `json:"availableBalance"`   // Available balance
	MaxWithdrawAmount  string `json:"maxWithdrawAmount"`  // Maximum amount for transfer out
	MarginAvailable    bool   `json:"marginAvailable"`    // Whether the asset can be used as margin in Multi-Assets mode
	UpdateTime         int64  `json:"updateTime"`         // Update timestamp
}

// CreateOrderRequest defines the parameters for creating a new order.
type CreateOrderRequest struct {
	Symbol                  string // required
	Side                    string // required (BUY/SELL)
	PositionSide            string // optional, default BOTH for One-way Mode
	Type                    string // required (LIMIT/MARKET/etc)
	TimeInForce             string // optional
	Quantity                string // optional, cannot be sent with closePosition=true
	ReduceOnly              string // optional, "true" or "false", default "false"
	Price                   string // optional
	NewClientOrderId        string // optional
	StopPrice               string // optional, used with STOP/STOP_MARKET or TAKE_PROFIT/TAKE_PROFIT_MARKET
	ClosePosition           string // optional, "true" or "false", Close-All
	ActivationPrice         string // optional, used with TRAILING_STOP_MARKET
	CallbackRate            string // optional, used with TRAILING_STOP_MARKET, min 0.1, max 10
	WorkingType             string // optional, "MARK_PRICE" or "CONTRACT_PRICE", default "CONTRACT_PRICE"
	PriceProtect            string // optional, "TRUE" or "FALSE", default "FALSE"
	NewOrderRespType        string // optional, "ACK" or "RESULT", default "ACK"
	PriceMatch              string // optional, OPPONENT/QUEUE variations
	SelfTradePreventionMode string // optional, EXPIRE_TAKER/EXPIRE_MAKER/EXPIRE_BOTH/NONE
	GoodTillDate            int64  // optional, order cancel time for timeInForce GTD
	RecvWindow              int64  // optional, default 5000
}

// CreateOrderResponse represents the response from creating an order.
type CreateOrderResponse struct {
	ClientOrderId           string `json:"clientOrderId"`
	CumQty                  string `json:"cumQty"`
	CumQuote                string `json:"cumQuote"`
	ExecutedQty             string `json:"executedQty"`
	OrderId                 int64  `json:"orderId"`
	AvgPrice                string `json:"avgPrice"`
	OrigQty                 string `json:"origQty"`
	Price                   string `json:"price"`
	ReduceOnly              bool   `json:"reduceOnly"`
	Side                    string `json:"side"`
	PositionSide            string `json:"positionSide"`
	Status                  string `json:"status"`
	StopPrice               string `json:"stopPrice"`
	ClosePosition           bool   `json:"closePosition"`
	Symbol                  string `json:"symbol"`
	TimeInForce             string `json:"timeInForce"`
	Type                    string `json:"type"`
	OrigType                string `json:"origType"`
	ActivatePrice           string `json:"activatePrice,omitempty"` // only with TRAILING_STOP_MARKET
	PriceRate               string `json:"priceRate,omitempty"`     // only with TRAILING_STOP_MARKET
	UpdateTime              int64  `json:"updateTime"`
	WorkingType             string `json:"workingType"`
	PriceProtect            bool   `json:"priceProtect"`
	PriceMatch              string `json:"priceMatch"`
	SelfTradePreventionMode string `json:"selfTradePreventionMode"`
	GoodTillDate            int64  `json:"goodTillDate,omitempty"` // only with GTD orders
}

// CancelOrderRequest defines the parameters for canceling an order.
type CancelOrderRequest struct {
	Symbol            string // required
	OrderId           int64  // optional, either orderId or origClientOrderId must be sent
	OrigClientOrderId string // optional, either orderId or origClientOrderId must be sent
	RecvWindow        int64  // optional, default 5000
}

// CancelOrderResponse represents the response from canceling an order.
type CancelOrderResponse struct {
	ClientOrderId           string `json:"clientOrderId"`
	CumQty                  string `json:"cumQty"`
	CumQuote                string `json:"cumQuote"`
	ExecutedQty             string `json:"executedQty"`
	OrderId                 int64  `json:"orderId"`
	OrigQty                 string `json:"origQty"`
	OrigType                string `json:"origType"`
	Price                   string `json:"price"`
	ReduceOnly              bool   `json:"reduceOnly"`
	Side                    string `json:"side"`
	PositionSide            string `json:"positionSide"`
	Status                  string `json:"status"`
	StopPrice               string `json:"stopPrice"`
	ClosePosition           bool   `json:"closePosition"`
	Symbol                  string `json:"symbol"`
	TimeInForce             string `json:"timeInForce"`
	Type                    string `json:"type"`
	ActivatePrice           string `json:"activatePrice,omitempty"` // only with TRAILING_STOP_MARKET
	PriceRate               string `json:"priceRate,omitempty"`     // only with TRAILING_STOP_MARKET
	UpdateTime              int64  `json:"updateTime"`
	WorkingType             string `json:"workingType"`
	PriceProtect            bool   `json:"priceProtect"`
	PriceMatch              string `json:"priceMatch"`
	SelfTradePreventionMode string `json:"selfTradePreventionMode"`
	GoodTillDate            int64  `json:"goodTillDate,omitempty"` // only with GTD orders
}

// CancelAllOrdersRequest defines the parameters for canceling all open orders.
type CancelAllOrdersRequest struct {
	Symbol     string // required
	RecvWindow int64  // optional, default 5000
}

// CancelAllOrdersResponse represents the response from canceling all open orders.
type CancelAllOrdersResponse struct {
	Code int    `json:"code"`
	Msg  string `json:"msg"`
}

// QueryOrderRequest defines the parameters for querying an order.
type QueryOrderRequest struct {
	Symbol            string // required
	OrderId           int64  // optional, either orderId or origClientOrderId must be sent
	OrigClientOrderId string // optional, either orderId or origClientOrderId must be sent
	RecvWindow        int64  // optional, default 5000
}

// QueryOrderResponse represents the response from querying an order.
type QueryOrderResponse struct {
	AvgPrice                string `json:"avgPrice"`
	ClientOrderId           string `json:"clientOrderId"`
	CumQuote                string `json:"cumQuote"`
	ExecutedQty             string `json:"executedQty"`
	OrderId                 int64  `json:"orderId"`
	OrigQty                 string `json:"origQty"`
	OrigType                string `json:"origType"`
	Price                   string `json:"price"`
	ReduceOnly              bool   `json:"reduceOnly"`
	Side                    string `json:"side"`
	PositionSide            string `json:"positionSide"`
	Status                  string `json:"status"`
	StopPrice               string `json:"stopPrice"`
	ClosePosition           bool   `json:"closePosition"`
	Symbol                  string `json:"symbol"`
	Time                    int64  `json:"time"` // order time
	TimeInForce             string `json:"timeInForce"`
	Type                    string `json:"type"`
	ActivatePrice           string `json:"activatePrice,omitempty"` // only with TRAILING_STOP_MARKET
	PriceRate               string `json:"priceRate,omitempty"`     // only with TRAILING_STOP_MARKET
	UpdateTime              int64  `json:"updateTime"`              // update time
	WorkingType             string `json:"workingType"`
	PriceProtect            bool   `json:"priceProtect"`
	PriceMatch              string `json:"priceMatch"`
	SelfTradePreventionMode string `json:"selfTradePreventionMode"`
	GoodTillDate            int64  `json:"goodTillDate,omitempty"` // only with GTD orders
}

// QueryCurrentOpenOrderRequest defines the parameters for querying a current open order.
// This is identical to QueryOrderRequest but represents a different endpoint.
type QueryCurrentOpenOrderRequest struct {
	Symbol            string // required
	OrderId           int64  // optional, either orderId or origClientOrderId must be sent
	OrigClientOrderId string // optional, either orderId or origClientOrderId must be sent
	RecvWindow        int64  // optional, default 5000
}

// QueryCurrentOpenOrderResponse represents the response from querying a current open order.
// This is identical to QueryOrderResponse but represents a different endpoint that only returns open orders.
type QueryCurrentOpenOrderResponse = QueryOrderResponse

// GetMyTradesRequest defines the parameters for getting account trade list.
type GetMyTradesRequest struct {
	Symbol     string // required
	OrderId    int64  // optional, can only be used in combination with symbol
	StartTime  int64  // optional, timestamp in ms
	EndTime    int64  // optional, timestamp in ms
	FromId     int64  // optional, trade id to fetch from. Default gets most recent trades
	Limit      int    // optional, default 500; max 1000
	RecvWindow int64  // optional, default 5000
}

// MyTrade represents a single trade from the user's trade history.
type MyTrade struct {
	Buyer           bool   `json:"buyer"`           // Whether the user is the buyer
	Commission      string `json:"commission"`      // Commission amount
	CommissionAsset string `json:"commissionAsset"` // Commission asset
	Id              int64  `json:"id"`              // Trade ID
	Maker           bool   `json:"maker"`           // Whether the user is the maker
	OrderId         int64  `json:"orderId"`         // Order ID
	Price           string `json:"price"`           // Trade price
	Qty             string `json:"qty"`             // Trade quantity
	QuoteQty        string `json:"quoteQty"`        // Quote quantity
	RealizedPnl     string `json:"realizedPnl"`     // Realized PnL
	Side            string `json:"side"`            // Trade side (BUY/SELL)
	PositionSide    string `json:"positionSide"`    // Position side (LONG/SHORT/BOTH)
	Symbol          string `json:"symbol"`          // Trading symbol
	Time            int64  `json:"time"`            // Trade timestamp
}

// GetPositionsRequest defines the parameters for getting position information.
type GetPositionsRequest struct {
	MarginAsset string // optional, if neither marginAsset nor symbol is sent, positions of all symbols will be returned
	Symbol      string // optional, if neither marginAsset nor symbol is sent, positions of all symbols will be returned
	RecvWindow  int64  // optional, default 5000
}

// Position represents position information for a single symbol.
type Position struct {
	Symbol           string `json:"symbol"`           // Symbol
	PositionAmt      string `json:"positionAmt"`      // Position amount
	EntryPrice       string `json:"entryPrice"`       // Entry price
	BreakEvenPrice   string `json:"breakEvenPrice"`   // Break-even price
	MarkPrice        string `json:"markPrice"`        // Mark price
	UnRealizedProfit string `json:"unRealizedProfit"` // Unrealized profit
	LiquidationPrice string `json:"liquidationPrice"` // Liquidation price
	Leverage         string `json:"leverage"`         // Leverage
	MaxQty           string `json:"maxQty"`           // Maximum quantity of base asset
	MarginType       string `json:"marginType"`       // Margin type (cross/isolated)
	IsolatedMargin   string `json:"isolatedMargin"`   // Isolated margin
	IsAutoAddMargin  string `json:"isAutoAddMargin"`  // Auto add margin ("true"/"false")
	PositionSide     string `json:"positionSide"`     // Position side (BOTH/LONG/SHORT)
	UpdateTime       int64  `json:"updateTime"`       // Update timestamp
}

// StartUserDataStreamResponse represents the response for starting a user data stream.
type StartUserDataStreamResponse struct {
	ListenKey string `json:"listenKey"`
}

// KeepaliveUserDataStreamResponse represents the response for keepalive user data stream.
type KeepaliveUserDataStreamResponse struct {
	ListenKey string `json:"listenKey"`
}

// CloseUserDataStreamResponse represents the empty response for closing user data stream.
type CloseUserDataStreamResponse struct{}
