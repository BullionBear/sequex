# Binance Futures User Data Streams

This document describes the user data stream events available in Binance Futures WebSocket API.

## Overview

User data streams provide real-time updates about account activities, including:
- Account balance changes
- Position updates
- Order status changes
- Trade executions
- Margin calls
- Configuration changes

## Event Types

### 1. Listen Key Expired (`listenKeyExpired`)

Triggered when the listen key expires and needs to be renewed.

**Payload:**
```json
{
    "e": "listenKeyExpired",    // event type
    "E": "1736996475556",       // event time
    "listenKey":"WsCMN0a4KHUPTQuX6IUnqEZfB1inxmv1qR4kbf1LuEjur5VdbzqvyxqG9TSjVVxv"
}
```

**Go Struct:**
```go
type WSListenKeyExpiredEvent struct {
    EventType string `json:"e"`
    EventTime int64  `json:"E"`
    ListenKey string `json:"listenKey"`
}
```

### 2. Account Update (`ACCOUNT_UPDATE`)

Triggered when account balances or positions change due to trading activities.

**Payload:**
```json
{
  "e": "ACCOUNT_UPDATE",        // Event Type
  "E": 1564745798939,           // Event Time
  "T": 1564745798938,           // Transaction Time
  "a": {                        // Update Data
    "m":"ORDER",                // Event reason type
    "B":[                       // Balances
      {
        "a":"USDT",             // Asset
        "wb":"122624.12345678", // Wallet Balance
        "cw":"100.12345678",    // Cross Wallet Balance
        "bc":"50.12345678"      // Balance Change except PnL and Commission
      }
    ],
    "P":[                       // Positions
      {
        "s":"BTCUSDT",          // Symbol
        "pa":"0",               // Position Amount
        "ep":"0.00000",         // Entry Price
        "bep":"0",              // Breakeven Price
        "cr":"200",             // (Pre-fee) Accumulated Realized
        "up":"0",               // Unrealized PnL
        "mt":"isolated",        // Margin Type
        "iw":"0.00000000",      // Isolated Wallet (if isolated position)
        "ps":"BOTH"             // Position Side
      }
    ]
  }
}
```

**Go Structs:**
```go
type WSAccountUpdateEvent struct {
    EventType       string              `json:"e"`
    EventTime       int64               `json:"E"`
    TransactionTime int64               `json:"T"`
    UpdateData      WSAccountUpdateData `json:"a"`
}

type WSAccountUpdateData struct {
    EventReasonType string           `json:"m"`
    Balances        []WSBalanceData  `json:"B"`
    Positions       []WSPositionData `json:"P"`
}

type WSBalanceData struct {
    Asset              string `json:"a"`
    WalletBalance      string `json:"wb"`
    CrossWalletBalance string `json:"cw"`
    BalanceChange      string `json:"bc"`
}

type WSPositionData struct {
    Symbol              string `json:"s"`
    PositionAmount      string `json:"pa"`
    EntryPrice          string `json:"ep"`
    BreakevenPrice      string `json:"bep"`
    AccumulatedRealized string `json:"cr"`
    UnrealizedPnL       string `json:"up"`
    MarginType          string `json:"mt"`
    IsolatedWallet      string `json:"iw"`
    PositionSide        string `json:"ps"`
}
```

### 3. Margin Call (`MARGIN_CALL`)

Triggered when a margin call occurs due to insufficient margin.

**Payload:**
```json
{
    "e":"MARGIN_CALL",          // Event Type
    "E":1587727187525,          // Event Time
    "cw":"3.16812045",          // Cross Wallet Balance. Only pushed with crossed position margin call
    "p":[                       // Position(s) of Margin Call
      {
        "s":"ETHUSDT",          // Symbol
        "ps":"LONG",            // Position Side
        "pa":"1.327",           // Position Amount
        "mt":"CROSSED",         // Margin Type
        "iw":"0",               // Isolated Wallet (if isolated position)
        "mp":"187.17127",       // Mark Price
        "up":"-1.166074",       // Unrealized PnL
        "mm":"1.614445"         // Maintenance Margin Required
      }
    ]
}
```

**Go Structs:**
```go
type WSMarginCallEvent struct {
    EventType          string                 `json:"e"`
    EventTime          int64                  `json:"E"`
    CrossWalletBalance string                 `json:"cw"`
    Positions          []WSMarginCallPosition `json:"p"`
}

type WSMarginCallPosition struct {
    Symbol                    string `json:"s"`
    PositionSide              string `json:"ps"`
    PositionAmount            string `json:"pa"`
    MarginType                string `json:"mt"`
    IsolatedWallet            string `json:"iw"`
    MarkPrice                 string `json:"mp"`
    UnrealizedPnL             string `json:"up"`
    MaintenanceMarginRequired string `json:"mm"`
}
```

### 4. Order Trade Update (`ORDER_TRADE_UPDATE`)

Triggered when order status changes or trades are executed.

**Payload:**
```json
{
  "e":"ORDER_TRADE_UPDATE",    // Event Type
  "E":1568879465651,           // Event Time
  "T":1568879465650,           // Transaction Time
  "o":{
    "s":"BTCUSDT",             // Symbol
    "c":"TEST",                // Client Order Id
    "S":"SELL",                // Side
    "o":"TRAILING_STOP_MARKET", // Order Type
    "f":"GTC",                 // Time in Force
    "q":"0.001",               // Original Quantity
    "p":"0",                   // Original Price
    "ap":"0",                  // Average Price
    "sp":"7103.04",            // Stop Price
    "x":"NEW",                 // Execution Type
    "X":"NEW",                 // Order Status
    "i":8886774,               // Order Id
    "l":"0",                   // Order Last Filled Quantity
    "z":"0",                   // Order Filled Accumulated Quantity
    "L":"0",                   // Last Filled Price
    "N":"USDT",                // Commission Asset
    "n":"0",                   // Commission
    "T":1568879465650,         // Order Trade Time
    "t":0,                     // Trade Id
    "b":"0",                   // Bids Notional
    "a":"9.91",                // Ask Notional
    "m":false,                 // Is this trade the maker side?
    "R":false,                 // Is this reduce only
    "wt":"CONTRACT_PRICE",     // Stop Price Working Type
    "ot":"TRAILING_STOP_MARKET", // Original Order Type
    "ps":"LONG",               // Position Side
    "cp":false,                // If Close-All, pushed with conditional order
    "AP":"7476.89",            // Activation Price, only pushed with TRAILING_STOP_MARKET order
    "cr":"5.0",                // Callback Rate, only pushed with TRAILING_STOP_MARKET order
    "pP": false,               // If price protection is turned on
    "si": 0,                   // ignore
    "ss": 0,                   // ignore
    "rp":"0",                  // Realized Profit of the trade
    "V":"EXPIRE_TAKER",        // STP mode
    "pm":"OPPONENT",           // Price match mode
    "gtd":0                    // TIF GTD order auto cancel time
  }
}
```

**Go Structs:**
```go
type WSOrderTradeUpdateEvent struct {
    EventType       string                  `json:"e"`
    EventTime       int64                   `json:"E"`
    TransactionTime int64                   `json:"T"`
    Order           WSOrderTradeUpdateOrder `json:"o"`
}

type WSOrderTradeUpdateOrder struct {
    Symbol                    string `json:"s"`
    ClientOrderID             string `json:"c"`
    Side                      string `json:"S"`
    OrderType                 string `json:"o"`
    TimeInForce               string `json:"f"`
    OriginalQuantity          string `json:"q"`
    OriginalPrice             string `json:"p"`
    AveragePrice              string `json:"ap"`
    StopPrice                 string `json:"sp"`
    ExecutionType             string `json:"x"`
    OrderStatus               string `json:"X"`
    OrderID                   int64  `json:"i"`
    LastFilledQuantity        string `json:"l"`
    FilledAccumulatedQuantity string `json:"z"`
    LastFilledPrice           string `json:"L"`
    CommissionAsset           string `json:"N"`
    Commission                string `json:"n"`
    OrderTradeTime            int64  `json:"T"`
    TradeID                   int64  `json:"t"`
    BidsNotional              string `json:"b"`
    AsksNotional              string `json:"a"`
    IsMaker                   bool   `json:"m"`
    IsReduceOnly              bool   `json:"R"`
    WorkingType               string `json:"wt"`
    OriginalOrderType         string `json:"ot"`
    PositionSide              string `json:"ps"`
    IsCloseAll                bool   `json:"cp"`
    ActivationPrice           string `json:"AP"`
    CallbackRate              string `json:"cr"`
    PriceProtection           bool   `json:"pP"`
    RealizedProfit            string `json:"rp"`
    STPMode                   string `json:"V"`
    PriceMatchMode            string `json:"pm"`
    GoodTillDate              int64  `json:"gtd"`
    Ignore1                   int    `json:"si,omitempty"` // ignore field
    Ignore2                   int    `json:"ss,omitempty"` // ignore field
}
```

### 5. Trade Lite (`TRADE_LITE`)

A simplified trade event that provides basic trade information.

**Payload:**
```json
{
  "e":"TRADE_LITE",            // Event Type
  "E":1721895408092,           // Event Time
  "T":1721895408214,           // Transaction Time                          
  "s":"BTCUSDT",               // Symbol
  "q":"0.001",                 // Original Quantity
  "p":"0",                     // Original Price
  "m":false,                   // Is this trade the maker side?
  "c":"z8hcUoOsqEdKMeKPSABslD", // Client Order Id
  "S":"BUY",                   // Side
  "L":"64089.20",              // Last Filled Price
  "l":"0.040",                 // Order Last Filled Quantity
  "t":109100866,               // Trade Id
  "i":8886774                  // Order Id
}
```

**Go Struct:**
```go
type WSTradeLiteEvent struct {
    EventType          string `json:"e"`
    EventTime          int64  `json:"E"`
    TransactionTime    int64  `json:"T"`
    Symbol             string `json:"s"`
    Quantity           string `json:"q"`
    Price              string `json:"p"`
    IsMaker            bool   `json:"m"`
    ClientOrderID      string `json:"c"`
    Side               string `json:"S"`
    LastFilledPrice    string `json:"L"`
    LastFilledQuantity string `json:"l"`
    TradeID            int64  `json:"t"`
    OrderID            int64  `json:"i"`
}
```

### 6. Account Config Update (`ACCOUNT_CONFIG_UPDATE`)

Triggered when account configuration changes (e.g., leverage changes).

**Payload:**
```json
{
    "e":"ACCOUNT_CONFIG_UPDATE", // Event Type
    "E":1611646737479,           // Event Time
    "T":1611646737476,           // Transaction Time
    "ac":{
        "s":"BTCUSDT",           // symbol
        "l":25                   // leverage
    }
}
```

**Go Structs:**
```go
type WSAccountConfigUpdateEvent struct {
    EventType       string          `json:"e"`
    EventTime       int64           `json:"E"`
    TransactionTime int64           `json:"T"`
    AccountConfig   WSAccountConfig `json:"ac"`
}

type WSAccountConfig struct {
    Symbol   string `json:"s"`
    Leverage int    `json:"l"`
}
```

## Usage Example

```go
// Create user data subscription options
userDataOptions := binancefuture.NewUserDataSubscriptionOptions()

// Set up callbacks for different event types
userDataOptions.WithConnect(func() {
    log.Println("Connected to user data stream")
}).WithDisconnect(func() {
    log.Println("Disconnected from user data stream")
}).WithError(func(err error) {
    log.Printf("User data stream error: %v", err)
}).WithAccountUpdate(func(data *binancefuture.WSOutboundAccountPosition) error {
    log.Printf("Account update received: %+v", data)
    return nil
}).WithBalanceUpdate(func(data *binancefuture.WSBalanceUpdate) error {
    log.Printf("Balance update received: %+v", data)
    return nil
}).WithExecutionReport(func(data *binancefuture.WSExecutionReport) error {
    log.Printf("Execution report received: %+v", data)
    return nil
}).WithListenKeyExpired(func(data *binancefuture.WSListenKeyExpiredEvent) error {
    log.Printf("Listen key expired: %+v", data)
    return nil
}).WithMarginCall(func(data *binancefuture.WSMarginCallEvent) error {
    log.Printf("Margin call received: %+v", data)
    return nil
}).WithOrderTradeUpdate(func(data *binancefuture.WSOrderTradeUpdateEvent) error {
    log.Printf("Order trade update received: %+v", data)
    return nil
}).WithTradeLite(func(data *binancefuture.WSTradeLiteEvent) error {
    log.Printf("Trade lite received: %+v", data)
    return nil
}).WithAccountConfigUpdate(func(data *binancefuture.WSAccountConfigUpdateEvent) error {
    log.Printf("Account config update received: %+v", data)
    return nil
})

// Subscribe to user data stream
unsubscribe, err := wsClient.SubscribeToUserDataStream(userDataOptions)
if err != nil {
    log.Fatalf("Failed to subscribe to user data stream: %v", err)
}
defer unsubscribe()
```

## Important Notes

1. **Listen Key Management**: Listen keys expire after 60 minutes. You should implement automatic renewal or reconnection logic.

2. **Event Ordering**: Events are not guaranteed to be received in chronological order. Use the `EventTime` field for ordering if needed.

3. **Rate Limits**: User data streams have rate limits. Implement proper error handling for rate limit errors.

4. **Reconnection**: Implement automatic reconnection logic for network interruptions.

5. **Error Handling**: Always handle errors in callback functions to prevent stream disconnections.

## Special Client Order IDs

Some events may contain special client order IDs:
- `autoclose-*`: Liquidation order
- `adl_autoclose`: ADL auto close order  
- `settlement_autoclose-*`: Settlement order for delisting or delivery

## References

- [Binance Futures API Documentation](https://binance-docs.github.io/apidocs/futures/en/)
- [User Data Streams](https://binance-docs.github.io/apidocs/futures/en/#user-data-streams) 