# Binance WebSocket Connectivity

This module provides WebSocket connectivity for Binance Spot API, allowing real-time data streaming for market data and user data.

## Features

- **Raw Stream Subscription**: Connect directly to Binance WebSocket streams using `wss://stream.binance.com:9443/ws/<streamName>`
- **Automatic Reconnection**: Exponential backoff retry mechanism with configurable attempts
- **Ping/Pong Keepalive**: Maintains connection health with periodic ping messages
- **Graceful Disconnection**: Proper cleanup and resource management
- **Callback-based Event Handling**: Flexible event processing with user-defined callbacks
- **Multiple Stream Support**: Subscribe to multiple streams simultaneously
- **Type-safe Data Parsing**: Structured data models for all WebSocket events

## Quick Start

```go
package main

import (
    "fmt"
    "log"
    "time"
    
    "github.com/BullionBear/sequex/pkg/exchange/binance"
)

func main() {
    // Create configuration
    config := binance.DefaultConfig()
    
    // Create WebSocket stream client
    client := binance.NewWSStreamClient(config)

    // Define callback for ticker data
    tickerCallback := func(data []byte) error {
        tickerData, err := binance.ParseTickerData(data)
        if err != nil {
            return err
        }
        
        fmt.Printf("Symbol: %s, Price: %.2f\n", 
            tickerData.Symbol, tickerData.LastPrice)
        return nil
    }

    // Subscribe to BTCUSDT ticker
    unsubscribe, err := client.SubscribeToTicker("BTCUSDT", tickerCallback)
    if err != nil {
        log.Fatalf("Failed to subscribe: %v", err)
    }
    defer unsubscribe()

    // Keep the program running
    time.Sleep(30 * time.Second)
}
```

## API Reference

### WSStreamClient

The main client for WebSocket stream management.

#### Constructor

```go
func NewWSStreamClient(config *Config) *WSStreamClient
```

#### Methods

##### SubscribeToKline
```go
func (c *WSStreamClient) SubscribeToKline(symbol string, interval string, callback WebSocketCallback) (func() error, error)
```
Subscribes to kline/candlestick data for a specific symbol and interval.

**Parameters:**
- `symbol`: Trading pair symbol (e.g., "BTCUSDT")
- `interval`: Kline interval (e.g., "1m", "5m", "1h", "1d")
- `callback`: Function to handle incoming data

**Returns:**
- Unsubscription function
- Error if subscription fails

##### SubscribeToTicker
```go
func (c *WSStreamClient) SubscribeToTicker(symbol string, callback WebSocketCallback) (func() error, error)
```
Subscribes to 24hr ticker data for a specific symbol.

##### SubscribeToMiniTicker
```go
func (c *WSStreamClient) SubscribeToMiniTicker(symbol string, callback WebSocketCallback) (func() error, error)
```
Subscribes to mini ticker data for a specific symbol.

##### SubscribeToAllMiniTickers
```go
func (c *WSStreamClient) SubscribeToAllMiniTickers(callback WebSocketCallback) (func() error, error)
```
Subscribes to mini ticker data for all symbols.

##### SubscribeToBookTicker
```go
func (c *WSStreamClient) SubscribeToBookTicker(symbol string, callback WebSocketCallback) (func() error, error)
```
Subscribes to book ticker data for a specific symbol.

##### SubscribeToAllBookTickers
```go
func (c *WSStreamClient) SubscribeToAllBookTickers(callback WebSocketCallback) (func() error, error)
```
Subscribes to book ticker data for all symbols.

##### SubscribeToDepth
```go
func (c *WSStreamClient) SubscribeToDepth(symbol string, levels string, callback WebSocketCallback) (func() error, error)
```
Subscribes to order book depth data for a specific symbol.

**Parameters:**
- `symbol`: Trading pair symbol
- `levels`: Depth levels (e.g., "5", "10", "20")
- `callback`: Function to handle incoming data

##### SubscribeToTrade
```go
func (c *WSStreamClient) SubscribeToTrade(symbol string, callback WebSocketCallback) (func() error, error)
```
Subscribes to trade data for a specific symbol.

##### SubscribeToAggTrade
```go
func (c *WSStreamClient) SubscribeToAggTrade(symbol string, callback WebSocketCallback) (func() error, error)
```
Subscribes to aggregated trade data for a specific symbol.

##### SubscribeToCombinedStreams
```go
func (c *WSStreamClient) SubscribeToCombinedStreams(streams []string, callback WebSocketCallback) (func() error, error)
```
Subscribes to multiple streams at once.

**Parameters:**
- `streams`: Array of stream names (e.g., ["btcusdt@ticker", "ethusdt@ticker"])
- `callback`: Function to handle incoming data

##### UnsubscribeFromAllStreams
```go
func (c *WSStreamClient) UnsubscribeFromAllStreams() error
```
Unsubscribes from all active streams.

##### GetActiveStreams
```go
func (c *WSStreamClient) GetActiveStreams() []string
```
Returns a list of currently active stream names.

##### IsStreamActive
```go
func (c *WSStreamClient) IsStreamActive(streamName string) bool
```
Checks if a specific stream is currently active.

##### Close
```go
func (c *WSStreamClient) Close() error
```
Closes all WebSocket connections.

### Data Models

#### WSKlineData
Represents kline/candlestick data from WebSocket.

```go
type WSKlineData struct {
    Symbol            string  `json:"s"`
    Kline             WSKline `json:"k"`
    EventTime         int64   `json:"E"`
    EventType         string  `json:"e"`
    FirstTradeID      int64   `json:"f"`
    LastTradeID       int64   `json:"L"`
    IsKlineClosed     bool    `json:"x"`
    QuoteVolume       float64 `json:"q,string"`
    ActiveBuyVolume   float64 `json:"V,string"`
    ActiveBuyQuoteVol float64 `json:"Q,string"`
}
```

#### WSTickerData
Represents 24hr ticker data from WebSocket.

```go
type WSTickerData struct {
    Symbol             string  `json:"s"`
    PriceChange        float64 `json:"P,string"`
    PriceChangePercent float64 `json:"P,string"`
    WeightedAvgPrice   float64 `json:"w,string"`
    PrevClosePrice     float64 `json:"x,string"`
    LastPrice          float64 `json:"c,string"`
    LastQty            float64 `json:"Q,string"`
    BidPrice           float64 `json:"b,string"`
    BidQty             float64 `json:"B,string"`
    AskPrice           float64 `json:"a,string"`
    AskQty             float64 `json:"A,string"`
    OpenPrice          float64 `json:"o,string"`
    HighPrice          float64 `json:"h,string"`
    LowPrice           float64 `json:"l,string"`
    Volume             float64 `json:"v,string"`
    QuoteVolume        float64 `json:"q,string"`
    OpenTime           int64   `json:"O"`
    CloseTime          int64   `json:"C"`
    FirstID            int64   `json:"F"`
    LastID             int64   `json:"L"`
    Count              int64   `json:"n"`
    EventTime          int64   `json:"E"`
    EventType          string  `json:"e"`
}
```

#### WSTradeData
Represents trade data from WebSocket.

```go
type WSTradeData struct {
    Symbol        string  `json:"s"`
    ID            int64   `json:"t"`
    Price         float64 `json:"p,string"`
    Quantity      float64 `json:"q,string"`
    BuyerOrderID  int64   `json:"b"`
    SellerOrderID int64   `json:"a"`
    TradeTime     int64   `json:"T"`
    IsBuyerMaker  bool    `json:"m"`
    Ignore        bool    `json:"M"`
    EventTime     int64   `json:"E"`
    EventType     string  `json:"e"`
}
```

### Data Parsing Functions

#### ParseKlineData
```go
func ParseKlineData(data []byte) (*WSKlineData, error)
```
Parses kline data from WebSocket message.

#### ParseTickerData
```go
func ParseTickerData(data []byte) (*WSTickerData, error)
```
Parses ticker data from WebSocket message.

#### ParseTradeData
```go
func ParseTradeData(data []byte) (*WSTradeData, error)
```
Parses trade data from WebSocket message.

#### ParseMiniTickerData
```go
func ParseMiniTickerData(data []byte) (*WSMiniTickerData, error)
```
Parses mini ticker data from WebSocket message.

#### ParseBookTickerData
```go
func ParseBookTickerData(data []byte) (*WSBookTickerData, error)
```
Parses book ticker data from WebSocket message.

#### ParseDepthData
```go
func ParseDepthData(data []byte) (*WSDepthData, error)
```
Parses depth data from WebSocket message.

#### ParseAggTradeData
```go
func ParseAggTradeData(data []byte) (*WSAggTradeData, error)
```
Parses aggregated trade data from WebSocket message.

## Configuration

### Production vs Testnet

```go
// Production
config := binance.DefaultConfig()

// Testnet
config := binance.TestnetConfig()
```

### WebSocket URLs

- **Production**: `wss://stream.binance.com:9443`
- **Testnet**: `wss://testnet.binance.vision`

## Stream Names

### Individual Streams
- `btcusdt@ticker` - 24hr ticker for BTCUSDT
- `btcusdt@kline_1m` - 1-minute klines for BTCUSDT
- `btcusdt@trade` - Trades for BTCUSDT
- `btcusdt@depth5` - Order book depth (5 levels) for BTCUSDT
- `btcusdt@bookTicker` - Book ticker for BTCUSDT
- `btcusdt@miniTicker` - Mini ticker for BTCUSDT
- `btcusdt@aggTrade` - Aggregated trades for BTCUSDT

### All Symbols Streams
- `!miniTicker@arr` - Mini ticker for all symbols
- `!bookTicker` - Book ticker for all symbols

### Combined Streams
Multiple streams can be combined using `/` separator:
- `btcusdt@ticker/ethusdt@ticker` - Both BTCUSDT and ETHUSDT tickers

## Error Handling

The WebSocket client includes comprehensive error handling:

1. **Connection Errors**: Automatic reconnection with exponential backoff
2. **Message Parsing Errors**: Graceful handling of malformed messages
3. **Callback Errors**: Logged but don't interrupt the connection
4. **Network Errors**: Automatic retry with configurable limits

## Best Practices

1. **Always use defer for unsubscription**:
   ```go
   unsubscribe, err := client.SubscribeToTicker("BTCUSDT", callback)
   if err != nil {
       return err
   }
   defer unsubscribe()
   ```

2. **Handle callback errors gracefully**:
   ```go
   callback := func(data []byte) error {
       tickerData, err := binance.ParseTickerData(data)
       if err != nil {
           log.Printf("Parse error: %v", err)
           return err // This won't break the connection
       }
       // Process data...
       return nil
   }
   ```

3. **Use signal handling for graceful shutdown**:
   ```go
   sigChan := make(chan os.Signal, 1)
   signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)
   
   // ... subscribe to streams ...
   
   <-sigChan
   client.Close()
   ```

4. **Monitor connection health**:
   ```go
   if !client.IsStreamActive("btcusdt@ticker") {
       // Handle inactive stream
   }
   ```

## Examples

See `cmd/ws_example.go` for a complete working example that demonstrates:
- Multiple stream subscriptions
- Data parsing
- Error handling
- Graceful shutdown

## Testing

Run the WebSocket tests:
```bash
go test ./pkg/exchange/binance -v -run TestWS
```

## Dependencies

- `github.com/gorilla/websocket` - WebSocket implementation
- Standard Go libraries for JSON parsing, time handling, etc. 