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
