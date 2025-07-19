package binance

import (
	"encoding/json"
	"time"
)

// WSRequest represents a WebSocket request message
type WSRequest struct {
	Method string   `json:"method"`
	Params []string `json:"params"`
	ID     int      `json:"id"`
}

// WSResponse represents a WebSocket response message
type WSResponse struct {
	Result json.RawMessage `json:"result,omitempty"`
	Error  *WSError        `json:"error,omitempty"`
	ID     int             `json:"id"`
}

// WSError represents a WebSocket error
type WSError struct {
	Code int    `json:"code"`
	Msg  string `json:"msg"`
}

// Error implements the error interface
func (e *WSError) Error() string {
	return e.Msg
}

// WSStreamMessage represents a WebSocket stream message wrapper
type WSStreamMessage struct {
	Stream string          `json:"stream"`
	Data   json.RawMessage `json:"data"`
}

// WSKlineEvent represents a kline/candlestick stream event
type WSKlineEvent struct {
	EventType string      `json:"e"` // Event type
	EventTime int64       `json:"E"` // Event time
	Symbol    string      `json:"s"` // Symbol
	Kline     WSKlineData `json:"k"` // Kline data
}

// WSKlineData represents kline/candlestick data from WebSocket
type WSKlineData struct {
	Symbol              string `json:"s"` // Symbol
	OpenTime            int64  `json:"t"` // Kline start time
	CloseTime           int64  `json:"T"` // Kline close time
	Interval            string `json:"i"` // Interval
	FirstTradeID        int64  `json:"f"` // First trade ID
	LastTradeID         int64  `json:"L"` // Last trade ID
	OpenPrice           string `json:"o"` // Open price
	ClosePrice          string `json:"c"` // Close price
	HighPrice           string `json:"h"` // High price
	LowPrice            string `json:"l"` // Low price
	BaseAssetVolume     string `json:"v"` // Base asset volume
	NumberOfTrades      int64  `json:"n"` // Number of trades
	IsClosed            bool   `json:"x"` // Is this kline closed?
	QuoteAssetVolume    string `json:"q"` // Quote asset volume
	TakerBuyBaseVolume  string `json:"V"` // Taker buy base asset volume
	TakerBuyQuoteVolume string `json:"Q"` // Taker buy quote asset volume
	Ignore              string `json:"B"` // Ignore
}

// GetOpenTime returns the kline open time as time.Time
func (k *WSKlineData) GetOpenTime() time.Time {
	return time.Unix(0, k.OpenTime*int64(time.Millisecond))
}

// GetCloseTime returns the kline close time as time.Time
func (k *WSKlineData) GetCloseTime() time.Time {
	return time.Unix(0, k.CloseTime*int64(time.Millisecond))
}

// WSTickerEvent represents a 24hr ticker statistics stream event
type WSTickerEvent struct {
	EventType              string `json:"e"` // Event type
	EventTime              int64  `json:"E"` // Event time
	Symbol                 string `json:"s"` // Symbol
	PriceChange            string `json:"p"` // Price change
	PriceChangePercent     string `json:"P"` // Price change percent
	WeightedAvgPrice       string `json:"w"` // Weighted average price
	FirstTradeBefore24hr   string `json:"x"` // First trade before the 24hr rolling window
	LastPrice              string `json:"c"` // Last price
	LastQuantity           string `json:"Q"` // Last quantity
	BestBidPrice           string `json:"b"` // Best bid price
	BestBidQuantity        string `json:"B"` // Best bid quantity
	BestAskPrice           string `json:"a"` // Best ask price
	BestAskQuantity        string `json:"A"` // Best ask quantity
	OpenPrice              string `json:"o"` // Open price
	HighPrice              string `json:"h"` // High price
	LowPrice               string `json:"l"` // Low price
	TotalTradedVolume      string `json:"v"` // Total traded base asset volume
	TotalTradedQuoteVolume string `json:"q"` // Total traded quote asset volume
	StatisticsOpenTime     int64  `json:"O"` // Statistics open time
	StatisticsCloseTime    int64  `json:"C"` // Statistics close time
	FirstTradeID           int64  `json:"F"` // First trade ID
	LastTradeID            int64  `json:"L"` // Last trade ID
	TotalTrades            int64  `json:"n"` // Total number of trades
}

// WSTradeEvent represents an individual trade stream event
type WSTradeEvent struct {
	EventType     string `json:"e"` // Event type
	EventTime     int64  `json:"E"` // Event time
	Symbol        string `json:"s"` // Symbol
	TradeID       int64  `json:"t"` // Trade ID
	Price         string `json:"p"` // Price
	Quantity      string `json:"q"` // Quantity
	BuyerOrderID  int64  `json:"b"` // Buyer order ID
	SellerOrderID int64  `json:"a"` // Seller order ID
	TradeTime     int64  `json:"T"` // Trade time
	IsBuyerMaker  bool   `json:"m"` // Is the buyer the market maker?
	Ignore        bool   `json:"M"` // Ignore
}

// WSDepthEvent represents a partial book depth stream event
type WSDepthEvent struct {
	EventType     string     `json:"e"` // Event type
	EventTime     int64      `json:"E"` // Event time
	Symbol        string     `json:"s"` // Symbol
	FirstUpdateID int64      `json:"U"` // First update ID in event
	FinalUpdateID int64      `json:"u"` // Final update ID in event
	Bids          [][]string `json:"b"` // Bids to be updated
	Asks          [][]string `json:"a"` // Asks to be updated
}

// WSBookTickerEvent represents individual symbol book ticker stream event
type WSBookTickerEvent struct {
	UpdateID     int64  `json:"u"` // Order book updateId
	Symbol       string `json:"s"` // Symbol
	BestBidPrice string `json:"b"` // Best bid price
	BestBidQty   string `json:"B"` // Best bid quantity
	BestAskPrice string `json:"a"` // Best ask price
	BestAskQty   string `json:"A"` // Best ask quantity
}

// WSAggTradeEvent represents aggregate trade stream event
type WSAggTradeEvent struct {
	EventType        string `json:"e"` // Event type
	EventTime        int64  `json:"E"` // Event time
	Symbol           string `json:"s"` // Symbol
	AggregateTradeID int64  `json:"a"` // Aggregate trade ID
	Price            string `json:"p"` // Price
	Quantity         string `json:"q"` // Quantity
	FirstTradeID     int64  `json:"f"` // First trade ID
	LastTradeID      int64  `json:"l"` // Last trade ID
	TradeTime        int64  `json:"T"` // Trade time
	IsBuyerMaker     bool   `json:"m"` // Is the buyer the market maker?
	Ignore           bool   `json:"M"` // Ignore
}

// BuildKlineStreamName builds a kline stream name for subscription
func BuildKlineStreamName(symbol, interval string) string {
	return NormalizeSymbol(symbol) + "@" + WSStreamKline + "_" + interval
}

// BuildTickerStreamName builds a ticker stream name for subscription
func BuildTickerStreamName(symbol string) string {
	return NormalizeSymbol(symbol) + "@" + WSStreamTicker
}

// BuildTradeStreamName builds a trade stream name for subscription
func BuildTradeStreamName(symbol string) string {
	return NormalizeSymbol(symbol) + "@" + WSStreamTrade
}

// BuildDepthStreamName builds a depth stream name for subscription
func BuildDepthStreamName(symbol string, levels int) string {
	if levels == 0 {
		return NormalizeSymbol(symbol) + "@" + WSStreamDepth
	}
	return NormalizeSymbol(symbol) + "@" + WSStreamDepth + string(rune('0'+levels))
}

// BuildBookTickerStreamName builds a book ticker stream name for subscription
func BuildBookTickerStreamName(symbol string) string {
	return NormalizeSymbol(symbol) + "@" + WSStreamBookTicker
}

// BuildAggTradeStreamName builds an aggregate trade stream name for subscription
func BuildAggTradeStreamName(symbol string) string {
	return NormalizeSymbol(symbol) + "@" + WSStreamAggTrade
}
