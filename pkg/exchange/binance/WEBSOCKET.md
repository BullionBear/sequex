# Binance WebSocket Subscription Pattern

This document describes the new WebSocket subscription pattern implemented in the Binance module, which provides a clean, chainable interface for subscribing to various market data and user data streams.

## Overview

The new pattern uses subscription options with chainable methods, making it easy to configure callbacks for different events while maintaining type safety and clean code structure.

## Core Concepts

### Subscription Options

Each subscription type has its own options struct that inherits from the base `SubscriptionOptions`:

- `KlineSubscriptionOptions` - for kline/candlestick data
- `TickerSubscriptionOptions` - for 24hr ticker data
- `MiniTickerSubscriptionOptions` - for mini ticker data
- `BookTickerSubscriptionOptions` - for book ticker data
- `DepthSubscriptionOptions` - for order book depth data
- `TradeSubscriptionOptions` - for trade data
- `AggTradeSubscriptionOptions` - for aggregated trade data
- `UserDataSubscriptionOptions` - for user data streams

### Chainable Methods

All subscription options support the following chainable methods:

- `WithConnect(callback func())` - called when connection is established
- `WithReconnect(callback func())` - called when reconnection occurs
- `WithDisconnect(callback func())` - called when connection is lost
- `WithError(callback func(error))` - called when errors occur

Plus type-specific data callbacks:

- `WithKline(callback KlineCallback)` - for kline data
- `WithTicker(callback TickerCallback)` - for ticker data
- `WithMiniTicker(callback MiniTickerCallback)` - for mini ticker data
- `WithBookTicker(callback BookTickerCallback)` - for book ticker data
- `WithDepth(callback DepthCallback)` - for depth data
- `WithTrade(callback TradeCallback)` - for trade data
- `WithAggTrade(callback AggTradeCallback)` - for aggregated trade data
- `WithExecutionReport(callback ExecutionReportCallback)` - for execution reports
- `WithAccountUpdate(callback OutboundAccountPositionCallback)` - for account updates
- `WithBalanceUpdate(callback BalanceUpdateCallback)` - for balance updates

## Usage Examples

### Basic Kline Subscription

```go
config := &Config{
    UseTestnet: true,
}
wsClient := NewWSStreamClient(config)

klineOptions := &KlineSubscriptionOptions{}
klineOptions.WithConnect(func() {
    fmt.Println("Connected to kline stream")
}).WithKline(func(data *WSKlineData) error {
    fmt.Printf("Kline: %s %s Close: %f\n", 
        data.Symbol, data.Kline.Interval, data.Kline.ClosePrice)
    return nil
}).WithError(func(err error) {
    fmt.Printf("Error: %v\n", err)
})

unsubscribe, err := wsClient.SubscribeToKline("BTCUSDT", "1m", klineOptions)
if err != nil {
    log.Fatalf("Failed to subscribe: %v", err)
}

// ... use the data ...

unsubscribe() // Clean up when done
```

### Ticker Subscription with All Callbacks

```go
tickerOptions := &TickerSubscriptionOptions{}
tickerOptions.WithConnect(func() {
    fmt.Println("Connected to ticker stream")
}).WithReconnect(func() {
    fmt.Println("Reconnected to ticker stream")
}).WithTicker(func(data *WSTickerData) error {
    fmt.Printf("Ticker: %s Last: %f Volume: %f Change: %f%%\n", 
        data.Symbol, data.LastPrice, data.Volume, data.PriceChangePercent)
    return nil
}).WithDisconnect(func() {
    fmt.Println("Disconnected from ticker stream")
}).WithError(func(err error) {
    fmt.Printf("Ticker error: %v\n", err)
})

unsubscribe, err := wsClient.SubscribeToTicker("ETHUSDT", tickerOptions)
```

### User Data Stream Subscription

```go
// Create listen key (requires API credentials)
restClient := NewClient(config)
userDataStream, err := restClient.CreateUserDataStream(context.Background())
if err != nil {
    log.Fatalf("Failed to create user data stream: %v", err)
}

userDataOptions := &UserDataSubscriptionOptions{}
userDataOptions.WithConnect(func() {
    fmt.Println("Connected to user data stream")
}).WithExecutionReport(func(data *WSExecutionReport) error {
    fmt.Printf("Order: %s %s %s Status: %s\n", 
        data.Symbol, data.Side, data.OrderPrice, data.CurrentOrderStatus)
    return nil
}).WithAccountUpdate(func(data *WSOutboundAccountPosition) error {
    fmt.Printf("Account update: %d balances\n", len(data.Balances))
    return nil
}).WithBalanceUpdate(func(data *WSBalanceUpdate) error {
    fmt.Printf("Balance: %s changed by %s\n", data.Asset, data.BalanceDelta)
    return nil
}).WithError(func(err error) {
    fmt.Printf("User data error: %v\n", err)
})

unsubscribe, err := wsClient.SubscribeToUserDataStream(userDataStream.ListenKey, userDataOptions)

// Keep the stream alive
go func() {
    ticker := time.NewTicker(30 * time.Minute)
    defer ticker.Stop()
    for range ticker.C {
        restClient.KeepAliveUserDataStream(context.Background(), userDataStream.ListenKey)
    }
}()

// Cleanup
defer func() {
    unsubscribe()
    restClient.CloseUserDataStream(context.Background(), userDataStream.ListenKey)
}()
```



### All Mini Tickers

```go
allMiniTickerOptions := &MiniTickerSubscriptionOptions{}
allMiniTickerOptions.WithConnect(func() {
    fmt.Println("Connected to all mini tickers")
}).WithMiniTicker(func(data *WSMiniTickerData) error {
    fmt.Printf("Mini: %s Close: %f Volume: %f\n", 
        data.Symbol, data.ClosePrice, data.Volume)
    return nil
}).WithError(func(err error) {
    fmt.Printf("Mini ticker error: %v\n", err)
})

unsubscribe, err := wsClient.SubscribeToAllMiniTickers(allMiniTickerOptions)
```

## Available Subscription Methods

### Market Data Streams

- `SubscribeToKline(symbol, interval, options)` - Kline/candlestick data
- `SubscribeToTicker(symbol, options)` - 24hr ticker data
- `SubscribeToMiniTicker(symbol, options)` - Mini ticker data
- `SubscribeToAllMiniTickers(options)` - All mini tickers
- `SubscribeToBookTicker(symbol, options)` - Book ticker data
- `SubscribeToAllBookTickers(options)` - All book tickers
- `SubscribeToDepth(symbol, levels, options)` - Order book depth data
- `SubscribeToTrade(symbol, options)` - Trade data
- `SubscribeToAggTrade(symbol, options)` - Aggregated trade data

### User Data Streams

- `SubscribeToUserDataStream(listenKey, options)` - User data stream

### Utility Methods

- `UnsubscribeFromAllStreams()` - Unsubscribe from all streams
- `GetActiveStreams()` - Get list of active streams
- `IsStreamActive(streamName)` - Check if stream is active
- `Close()` - Close all connections

## Error Handling

All callback functions that return errors should handle them appropriately:

```go
klineOptions.WithKline(func(data *WSKlineData) error {
    // Process the data
    if err := processKlineData(data); err != nil {
        // Log the error but don't panic
        log.Printf("Error processing kline data: %v", err)
        return err // This will be logged by the WebSocket client
    }
    return nil
})
```

## Best Practices

1. **Always provide error callbacks** to handle connection and parsing errors
2. **Use defer statements** to ensure proper cleanup of subscriptions
3. **Keep user data streams alive** with periodic keep-alive calls
4. **Handle reconnection gracefully** by providing reconnect callbacks
5. **Use type-safe callbacks** for all data handling
6. **Close connections properly** when your application shuts down

## API Design Philosophy

The new subscription pattern is designed to:

1. **Hide implementation details** - Users don't need to handle raw WebSocket callbacks
2. **Provide type safety** - All callbacks are strongly typed for specific data structures
3. **Enable chainable configuration** - Easy to configure multiple callbacks in a readable way
4. **Focus on business logic** - Developers can focus on what they want to do with the data, not how to handle the WebSocket connection

## Type Definitions

All callback types are defined in `ws_models.go`:

```go
type KlineCallback func(data *WSKlineData) error
type TickerCallback func(data *WSTickerData) error
type MiniTickerCallback func(data *WSMiniTickerData) error
type BookTickerCallback func(data *WSBookTickerData) error
type DepthCallback func(data *WSDepthData) error
type TradeCallback func(data *WSTradeData) error
type AggTradeCallback func(data *WSAggTradeData) error
type ExecutionReportCallback func(data *WSExecutionReport) error
type OutboundAccountPositionCallback func(data *WSOutboundAccountPosition) error
type BalanceUpdateCallback func(data *WSBalanceUpdate) error
```

This pattern provides a clean, extensible, and type-safe way to handle WebSocket subscriptions in the Binance module. 