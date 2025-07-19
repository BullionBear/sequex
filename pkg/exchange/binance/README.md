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
