# Binance Exchange Client

A comprehensive Go client for the Binance cryptocurrency exchange, providing both REST API and WebSocket connectivity for real-time trading and market data.

## Features

- **REST API Client** - Complete trading and market data API
- **Public WebSocket** - Real-time market data streams (kline, ticker, trades, orderbook)
- **Private WebSocket** - User data streams (account updates, order execution reports)
- **Authentication** - API key and signature-based authentication
- **Error Handling** - Comprehensive error handling and logging
- **Testnet Support** - Full testnet environment support
- **Context Support** - Timeout and cancellation handling

## Quick Start

### Installation

```bash
go get github.com/BullionBear/sequex/pkg/exchange/binance
```

### Basic Configuration

```go
package main

import (
    "context"
    "log"
    "time"
    
    "github.com/BullionBear/sequex/pkg/exchange/binance"
)

func main() {
    // Create configuration
    config := &binance.Config{
        APIKey:     "your_api_key",
        APISecret:  "your_api_secret",
        BaseURL:    binance.BaseURLSpot,        // Production
        // BaseURL: binance.BaseURLSpotTestnet, // Testnet
        Timeout:    10 * time.Second,
        UseTestnet: false, // Set to true for testnet
    }
    
    // Create client
    client := binance.NewClient(config)
    
    // Use the client...
}
```

## REST API Examples

### Market Data

#### Get Server Time
```go
ctx := context.Background()
serverTime, err := client.GetServerTime(ctx)
if err != nil {
    log.Fatal(err)
}
fmt.Printf("Server time: %v\n", serverTime.GetTime())
```

#### Get Current Price
```go
// Single symbol
result, err := client.GetTickerPrice(ctx, "BTCUSDT")
if err != nil {
    log.Fatal(err)
}
if result.IsSingle() {
    ticker := result.GetSingle()
    fmt.Printf("BTC price: %s\n", ticker.Price)
}

// All symbols
result, err := client.GetTickerPrice(ctx, "")
if err != nil {
    log.Fatal(err)
}
if result.IsArray() {
    tickers := result.GetArray()
    for _, p := range tickers {
        fmt.Printf("%s: %s\n", p.Symbol, p.Price)
    }
}
```

#### Get 24hr Ticker
```go
result, err := client.GetTicker24hr(ctx, "BTCUSDT")
if err != nil {
    log.Fatal(err)
}
if result.IsSingle() {
    ticker := result.GetSingle()
    fmt.Printf("24hr Change: %s%%\n", ticker.PriceChangePercent)
    fmt.Printf("Volume: %s\n", ticker.Volume)
}
```

#### Get Kline/Candlestick Data
```go
klines, err := client.GetKlines(ctx, "BTCUSDT", binance.Interval1h, 100)
if err != nil {
    log.Fatal(err)
}
for _, kline := range klines {
    fmt.Printf("Time: %v, Open: %s, High: %s, Low: %s, Close: %s\n",
        time.Unix(kline[0].(int64)/1000, 0),
        kline[1], kline[2], kline[3], kline[4])
}
```

#### Get Order Book
```go
orderbook, err := client.GetOrderBook(ctx, "BTCUSDT", 10)
if err != nil {
    log.Fatal(err)
}
fmt.Printf("Last Update ID: %d\n", orderbook.LastUpdateId)
for i, bid := range orderbook.Bids[:5] {
    fmt.Printf("Bid %d: Price=%s, Qty=%s\n", i+1, bid[0], bid[1])
}
```

#### Get Recent Trades
```go
trades, err := client.GetTrades(ctx, "BTCUSDT", 10)
if err != nil {
    log.Fatal(err)
}
for _, trade := range trades {
    fmt.Printf("Trade: %s %s @ %s (ID: %d)\n", 
        trade.Qty, trade.Price, trade.Time, trade.Id)
}
```

### Trading Operations

#### Get Account Information
```go
account, err := client.GetAccount(ctx)
if err != nil {
    log.Fatal(err)
}
fmt.Printf("Account Type: %s\n", account.AccountType)
for _, balance := range account.Balances {
    if balance.Free != "0.00000000" || balance.Locked != "0.00000000" {
        fmt.Printf("%s: Free=%s, Locked=%s\n", 
            balance.Asset, balance.Free, balance.Locked)
    }
}
```

#### Place Limit Order
```go
orderReq := &binance.NewOrderRequest{
    Symbol:      "BTCUSDT",
    Side:        "BUY",
    Type:        "LIMIT",
    TimeInForce: "GTC",
    Quantity:    "0.001",
    Price:       "50000.00",
}
order, err := client.PlaceOrder(ctx, orderReq)
if err != nil {
    log.Fatal(err)
}
fmt.Printf("Order placed: ID=%d, Status=%s\n", order.OrderId, order.Status)
```

#### Place Market Order
```go
orderReq := &binance.NewOrderRequest{
    Symbol:   "BTCUSDT",
    Side:     "BUY",
    Type:     "MARKET",
    Quantity: "0.001",
}
order, err := client.PlaceOrder(ctx, orderReq)
if err != nil {
    log.Fatal(err)
}
fmt.Printf("Market order filled: ID=%d, ExecutedQty=%s\n", 
    order.OrderId, order.ExecutedQty)
```

#### Cancel Order
```go
cancelResp, err := client.CancelOrder(ctx, "BTCUSDT", orderId)
if err != nil {
    log.Fatal(err)
}
fmt.Printf("Order canceled: ID=%d, Status=%s\n", 
    cancelResp.OrderId, cancelResp.Status)
```

#### Get Order Status
```go
order, err := client.GetOrder(ctx, "BTCUSDT", orderId)
if err != nil {
    log.Fatal(err)
}
fmt.Printf("Order: ID=%d, Status=%s, Price=%s, Qty=%s\n",
    order.OrderId, order.Status, order.Price, order.OrigQty)
```

#### Get Open Orders
```go
openOrders, err := client.GetOpenOrders(ctx, "BTCUSDT")
if err != nil {
    log.Fatal(err)
}
fmt.Printf("Open orders: %d\n", len(openOrders))
for _, order := range openOrders {
    fmt.Printf("  ID=%d, %s %s %s @ %s\n",
        order.OrderId, order.Side, order.Type, order.OrigQty, order.Price)
}
```

## WebSocket Examples

### Public WebSocket (Market Data)

#### Subscribe to Kline/Candlestick Stream
```go
// Create WebSocket client
wsClient := binance.NewWSClient(config)

// Set up event handlers
wsClient.OnKline(func(event *binance.WSKlineEvent) {
    fmt.Printf("Kline: %s %s - Open: %s, High: %s, Low: %s, Close: %s\n",
        event.Symbol, event.Kline.Interval,
        event.Kline.Open, event.Kline.High, event.Kline.Low, event.Kline.Close)
})

// Subscribe to kline stream
err := wsClient.SubscribeKline(ctx, "BTCUSDT", binance.Interval1m)
if err != nil {
    log.Fatal(err)
}

// Keep connection alive
time.Sleep(60 * time.Second)

// Unsubscribe
wsClient.UnsubscribeKline(ctx, "BTCUSDT", binance.Interval1m)
```

#### Subscribe to Multiple Streams
```go
wsClient := binance.NewWSClient(config)

// Set up handlers
wsClient.OnTicker(func(event *binance.WSTickerEvent) {
    fmt.Printf("Ticker: %s - Price: %s, Change: %s%%\n",
        event.Symbol, event.Price, event.PriceChangePercent)
})

wsClient.OnTrade(func(event *binance.WSTradeEvent) {
    fmt.Printf("Trade: %s - %s %s @ %s\n",
        event.Symbol, event.Side, event.Quantity, event.Price)
})

// Subscribe to multiple streams
streams := []string{
    binance.BuildKlineStreamName("BTCUSDT", binance.Interval1m),
    binance.BuildTickerStreamName("BTCUSDT"),
    binance.BuildTradeStreamName("BTCUSDT"),
}

err := wsClient.SubscribeMultiple(ctx, streams)
if err != nil {
    log.Fatal(err)
}

// Listen for events
time.Sleep(30 * time.Second)
```

#### Subscribe to Order Book Stream
```go
wsClient := binance.NewWSClient(config)

wsClient.OnDepth(func(event *binance.WSDepthEvent) {
    fmt.Printf("Order Book Update: %s\n", event.Symbol)
    fmt.Printf("  Bids: %d, Asks: %d\n", len(event.Bids), len(event.Asks))
    
    // Show top 5 bids and asks
    for i := 0; i < 5 && i < len(event.Bids); i++ {
        fmt.Printf("  Bid %d: %s @ %s\n", i+1, event.Bids[i].Quantity, event.Bids[i].Price)
    }
    for i := 0; i < 5 && i < len(event.Asks); i++ {
        fmt.Printf("  Ask %d: %s @ %s\n", i+1, event.Asks[i].Quantity, event.Asks[i].Price)
    }
})

err := wsClient.SubscribeDepth(ctx, "BTCUSDT", 5)
if err != nil {
    log.Fatal(err)
}

time.Sleep(30 * time.Second)
```

### Private WebSocket (User Data Stream)

#### Real-Time Account and Order Updates
```go
// Create user data stream client
userStream := binance.NewUserDataStreamClient(config)

// Set up event handlers
userStream.OnAccountUpdate(func(event *binance.WSAccountUpdate) {
    fmt.Printf("Account Update: %d balances updated\n", len(event.Balances))
    for _, balance := range event.Balances {
        if balance.Free != "0.00000000" || balance.Locked != "0.00000000" {
            fmt.Printf("  %s: Free=%s, Locked=%s\n", 
                balance.Asset, balance.Free, balance.Locked)
        }
    }
})

userStream.OnExecutionReport(func(event *binance.WSExecutionReport) {
    fmt.Printf("Order Update: %s %s %s OrderID:%d Status:%s\n", 
        event.Symbol, event.Side, event.OrderType, event.OrderID, event.CurrentOrderStatus)
    
    if event.IsNewOrder() {
        fmt.Println("  → New order created")
    }
    if event.IsTrade() {
        fmt.Printf("  → Trade executed: %s @ %s\n", 
            event.LastExecutedQuantity, event.LastExecutedPrice)
    }
    if event.IsCanceled() {
        fmt.Println("  → Order canceled")
    }
    if event.IsFilled() {
        fmt.Println("  → Order completely filled")
    }
})

userStream.OnBalanceUpdate(func(event *binance.WSBalanceUpdate) {
    fmt.Printf("Balance Update: %s Delta=%s\n", event.Asset, event.BalanceDelta)
})

// Connect to user data stream
err := userStream.Connect(ctx)
if err != nil {
    log.Fatal(err)
}

fmt.Println("✅ Connected to user data stream")

// Listen for events
time.Sleep(60 * time.Second)

// Disconnect
err = userStream.Disconnect()
if err != nil {
    log.Printf("Failed to disconnect: %v", err)
}
```

#### Manual Listen Key Management
```go
// Create listen key manually
streamResp, err := client.CreateUserDataStream(ctx)
if err != nil {
    log.Fatal(err)
}

fmt.Printf("Created listen key: %s\n", streamResp.ListenKey[:8]+"...")

// Keep alive (should be done every 30-60 minutes)
err = client.KeepAliveUserDataStream(ctx, streamResp.ListenKey)
if err != nil {
    log.Printf("Failed to keep alive: %v", err)
}

// Close when done
err = client.CloseUserDataStream(ctx, streamResp.ListenKey)
if err != nil {
    log.Printf("Failed to close: %v", err)
}
```

## Complete Trading Bot Example

```go
package main

import (
    "context"
    "fmt"
    "log"
    "time"
    
    "github.com/BullionBear/sequex/pkg/exchange/binance"
)

func main() {
    // Configuration
    config := &binance.Config{
        APIKey:     "your_api_key",
        APISecret:  "your_api_secret",
        BaseURL:    binance.BaseURLSpotTestnet, // Use testnet for testing
        Timeout:    10 * time.Second,
        UseTestnet: true,
    }
    
    // Create clients
    restClient := binance.NewClient(config)
    wsClient := binance.NewWSClient(config)
    userStream := binance.NewUserDataStreamClient(config)
    
    ctx := context.Background()
    
    // Get account info
    account, err := restClient.GetAccount(ctx)
    if err != nil {
        log.Fatal(err)
    }
    fmt.Printf("Account: %s\n", account.AccountType)
    
    // Set up WebSocket handlers
    wsClient.OnKline(func(event *binance.WSKlineEvent) {
        fmt.Printf("Kline: %s %s - Close: %s\n", 
            event.Symbol, event.Kline.Interval, event.Kline.Close)
    })
    
    userStream.OnExecutionReport(func(event *binance.WSExecutionReport) {
        fmt.Printf("Order: %d %s -> %s\n", 
            event.OrderID, event.CurrentExecutionType, event.CurrentOrderStatus)
    })
    
    // Subscribe to market data
    err = wsClient.SubscribeKline(ctx, "BTCUSDT", binance.Interval1m)
    if err != nil {
        log.Fatal(err)
    }
    
    // Connect to user data stream
    err = userStream.Connect(ctx)
    if err != nil {
        log.Fatal(err)
    }
    
    // Place a test order
    orderReq := &binance.NewOrderRequest{
        Symbol:      "BTCUSDT",
        Side:        "BUY",
        Type:        "LIMIT",
        TimeInForce: "GTC",
        Quantity:    "0.001",
        Price:       "50000.00", // Low price, won't fill
    }
    
    order, err := restClient.PlaceOrder(ctx, orderReq)
    if err != nil {
        log.Fatal(err)
    }
    fmt.Printf("Order placed: ID=%d\n", order.OrderId)
    
    // Wait for events
    time.Sleep(10 * time.Second)
    
    // Cancel order
    _, err = restClient.CancelOrder(ctx, "BTCUSDT", order.OrderId)
    if err != nil {
        log.Printf("Failed to cancel order: %v", err)
    }
    
    // Clean up
    userStream.Disconnect()
    fmt.Println("Trading bot completed")
}
```

## Configuration Options

### Environment Variables
```bash
export BINANCE_API_KEY="your_api_key"
export BINANCE_API_SECRET="your_api_secret"
export BINANCE_USE_TESTNET="true"
```

### Configuration Struct
```go
type Config struct {
    APIKey     string        // Your Binance API key
    APISecret  string        // Your Binance API secret
    BaseURL    string        // API base URL
    Timeout    time.Duration // HTTP request timeout
    UseTestnet bool          // Use testnet environment
}
```

## Error Handling

The client provides comprehensive error handling:

```go
// Check for specific error types
if err != nil {
    if apiErr, ok := err.(*binance.APIError); ok {
        switch apiErr.Code {
        case -1001: // Internal error
            log.Printf("Binance internal error: %s", apiErr.Message)
        case -1003: // Too many requests
            log.Printf("Rate limit exceeded: %s", apiErr.Message)
        case -2010: // Insufficient balance
            log.Printf("Insufficient balance: %s", apiErr.Message)
        default:
            log.Printf("API error %d: %s", apiErr.Code, apiErr.Message)
        }
    } else {
        log.Printf("Network error: %v", err)
    }
}
```

## Testing

Run the test suite:

```bash
# Run all tests
go test -v ./...

# Run specific test categories
go test -v -run="TestClient" ./...           # REST API tests
go test -v -run="TestWSClient" ./...         # WebSocket tests
go test -v -run="TestUserDataStream" ./...   # User data stream tests

# Run with testnet credentials
BINANCE_API_KEY="your_testnet_key" \
BINANCE_API_SECRET="your_testnet_secret" \
go test -v ./...
```

## Best Practices

1. **Use Context** - Always provide context with timeouts for API calls
2. **Handle Errors** - Check for specific error types and handle appropriately
3. **Rate Limiting** - Respect Binance's rate limits (1200 requests/minute)
4. **WebSocket Reconnection** - The client handles reconnection automatically
5. **Listen Key Management** - User data streams require keep-alive every 30-60 minutes
6. **Testnet First** - Always test with testnet before using production
7. **Secure Credentials** - Never hardcode API keys, use environment variables

## Rate Limits

- **REST API**: 1200 requests per minute per IP
- **WebSocket**: No rate limits, but connection limits apply
- **User Data Stream**: 1 stream per API key

## Support

For issues and questions:
- Check the test files for usage examples
- Review the Binance API documentation
- Ensure you're using the correct environment (testnet vs production)

## License

This package is part of the Sequex project. See the main project license for details. 