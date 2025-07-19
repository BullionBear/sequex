# Binance Futures Public REST API Client

This package provides a Go client for the Binance Futures public REST API, following the same patterns as the Binance Spot implementation.

## Features

- **Public REST API endpoints** for Binance Futures
- **Testnet support** for development and testing
- **Comprehensive error handling** with Binance-specific error codes
- **Structured logging** with slog
- **Test-driven development** with full test coverage
- **Type-safe responses** with proper JSON unmarshaling

## Supported Endpoints

### General Endpoints
- `GET /fapi/v1/time` - Get server time
- `GET /fapi/v1/ping` - Test connectivity
- `GET /fapi/v1/exchangeInfo` - Get exchange information and symbol filters

### Market Data Endpoints
- `GET /fapi/v1/ticker/price` - Get symbol price ticker
- `GET /fapi/v1/ticker/24hr` - Get 24hr ticker price change statistics
- `GET /fapi/v1/depth` - Get order book
- `GET /fapi/v1/klines` - Get kline/candlestick data
- `GET /fapi/v1/trades` - Get recent trades
- `GET /fapi/v1/premiumIndex` - Get mark price
- `GET /fapi/v1/fundingRate` - Get funding rate
- `GET /fapi/v1/openInterest` - Get open interest

## Quick Start

### Public Endpoints

```go
package main

import (
    "context"
    "fmt"
    "log"
    "time"
    
    "github.com/BullionBear/sequex/pkg/exchange/binancefuture"
)

func main() {
    // Create client with testnet configuration
    config := binancefuture.TestnetConfig()
    client := binancefuture.NewClient(config)
    
    ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
    defer cancel()
    
    // Get server time
    serverTime, err := client.GetServerTime(ctx)
    if err != nil {
        log.Fatalf("Failed to get server time: %v", err)
    }
    fmt.Printf("Server time: %v\n", serverTime.GetTime())
    
    // Get BTCUSDT price
    result, err := client.GetTickerPrice(ctx, "BTCUSDT")
    if err != nil {
        log.Fatalf("Failed to get ticker price: %v", err)
    }
    
    if result.IsSingle() {
        ticker := result.GetSingle()
        fmt.Printf("BTCUSDT price: %s\n", ticker.Price)
    }
}
```

### Signed Endpoints

```go
// Create client with API credentials
config := binancefuture.TestnetConfig()
config.APIKey = "your-api-key"
config.APISecret = "your-api-secret"
client := binancefuture.NewClient(config)

// Get account information
account, err := client.GetAccount(ctx)
if err != nil {
    log.Fatalf("Failed to get account: %v", err)
}
fmt.Printf("Total wallet balance: %s\n", account.TotalWalletBalance)

// Place a limit order
orderReq := &binancefuture.NewOrderRequest{
    Symbol:      "BTCUSDT",
    Side:        binancefuture.SideBuy,
    Type:        binancefuture.OrderTypeLimit,
    TimeInForce: binancefuture.TimeInForceGTC,
    Quantity:    "0.001",
    Price:       "50000",
}

order, err := client.PlaceOrder(ctx, orderReq)
if err != nil {
    log.Fatalf("Failed to place order: %v", err)
}
fmt.Printf("Order placed: ID=%d, Status=%s\n", order.OrderId, order.Status)

// Get position risk information
positionRisks, err := client.GetPositionRisk(ctx, "")
if err != nil {
    log.Fatalf("Failed to get position risk: %v", err)
}
fmt.Printf("Found %d position risk entries\n", len(positionRisks))

// Get and change position side mode
positionSide, err := client.GetPositionSide(ctx)
if err != nil {
    log.Fatalf("Failed to get position side: %v", err)
}
fmt.Printf("Current position side mode: dualSidePosition=%t\n", positionSide.DualSidePosition)

// Get and change leverage
leverage, err := client.GetLeverage(ctx, "BTCUSDT")
if err != nil {
    log.Fatalf("Failed to get leverage: %v", err)
}
fmt.Printf("Current leverage for BTCUSDT: %dx\n", leverage.Leverage)

## Configuration

### Testnet Configuration
```go
config := binancefuture.TestnetConfig()
client := binancefuture.NewClient(config)
```

### Production Configuration
```go
config := binancefuture.DefaultConfig()
client := binancefuture.NewClient(config)
```

### Custom Configuration
```go
config := &binancefuture.Config{
    BaseURL: "https://fapi.binance.com",
    Timeout: 30 * time.Second,
    UseTestnet: false,
}
client := binancefuture.NewClient(config)
```

## API Methods

### General
- `GetServerTime(ctx)` - Get server time
- `Ping(ctx)` - Test connectivity
- `GetExchangeInfo(ctx)` - Get exchange information

### Market Data
- `GetTickerPrice(ctx, symbol)` - Get price ticker (single or all symbols)
- `GetTicker24hr(ctx, symbol)` - Get 24hr ticker (single or all symbols)
- `GetKlines(ctx, symbol, interval, limit)` - Get candlestick data
- `GetOrderBook(ctx, symbol, limit)` - Get order book
- `GetTrades(ctx, symbol, limit)` - Get recent trades
- `GetMarkPrice(ctx, symbol)` - Get mark price
- `GetFundingRate(ctx, symbol, limit)` - Get funding rate history
- `GetOpenInterest(ctx, symbol)` - Get open interest

### Account & Trading (Signed Endpoints)
- `GetAccount(ctx)` - Get account information and positions
- `GetPositionRisk(ctx, symbol)` - Get position risk information
- `GetPositionSide(ctx)` - Get current position side mode
- `ChangePositionSide(ctx, dualSidePosition)` - Change position side mode
- `GetLeverage(ctx, symbol)` - Get current leverage for a symbol
- `ChangeLeverage(ctx, symbol, leverage)` - Change leverage for a symbol
- `PlaceOrder(ctx, req)` - Place a new order
- `GetOrder(ctx, symbol, orderID)` - Get order information
- `CancelOrder(ctx, symbol, orderID)` - Cancel an order
- `GetOpenOrders(ctx, symbol)` - Get open orders
- `GetUserTrades(ctx, symbol, limit)` - Get user trade history

## Response Types

### TickerPriceResult
Handles both single symbol and all symbols responses:
```go
result, err := client.GetTickerPrice(ctx, "BTCUSDT")
if result.IsSingle() {
    ticker := result.GetSingle()
    fmt.Printf("Price: %s\n", ticker.Price)
} else {
    tickers := result.GetArray()
    fmt.Printf("Found %d tickers\n", len(tickers))
}
```

### KlineResponse
Array of candlestick data with custom JSON unmarshaling:
```go
klines, err := client.GetKlines(ctx, "BTCUSDT", binancefuture.Interval1h, 10)
for _, kline := range *klines {
    fmt.Printf("Open: %s, High: %s, Low: %s, Close: %s\n",
        kline.Open, kline.High, kline.Low, kline.Close)
}
```

## Error Handling

The client returns structured errors that implement the `error` interface:

```go
result, err := client.GetTickerPrice(ctx, "INVALID")
if err != nil {
    if apiErr, ok := err.(*binancefuture.APIError); ok {
        fmt.Printf("API Error: %d - %s\n", apiErr.Code, apiErr.Message)
    } else {
        fmt.Printf("Network error: %v\n", err)
    }
}
```

## Testing

Run the test suite:
```bash
go test ./pkg/exchange/binancefuture -v
```

The tests use the Binance Futures testnet and verify all public endpoints work correctly.

## Constants

### Kline Intervals
- `Interval1m`, `Interval3m`, `Interval5m`, `Interval15m`, `Interval30m`
- `Interval1h`, `Interval2h`, `Interval4h`, `Interval6h`, `Interval8h`, `Interval12h`
- `Interval1d`, `Interval3d`, `Interval1w`, `Interval1M`

### Order Types
- `OrderTypeLimit`, `OrderTypeMarket`, `OrderTypeStopLoss`, etc.

### Order Sides
- `SideBuy`, `SideSell`

## Base URLs

- **Production**: `https://fapi.binance.com`
- **Testnet**: `https://testnet.binancefuture.com`

## Rate Limits

The client respects Binance's rate limits. For detailed rate limit information, see the [Binance Futures API documentation](https://binance-docs.github.io/apidocs/futures/en/).

## Contributing

When adding new endpoints:

1. Add endpoint constants to `const.go`
2. Add response types to `models.go`
3. Implement the method in `client.go`
4. Add comprehensive tests in `client_test.go`
5. Add examples in `example_test.go`

Follow the existing patterns and ensure all tests pass before submitting changes. 