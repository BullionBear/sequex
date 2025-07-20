# Binance Connectivity Module

This module implements **Binance Spot connectivity** for trading and data ingestion in Go, aligned with the overall `exchange` module design.

## Agent

- Your implementation should follow Test Driven Development, that is:
    - Make a function/go struct with empty implementation
    - Add an unittest, and I will review it and filled missing resources, e.g. API KEY
    - After reviewing, go to implementation
    - Run the unittest


I don't expect you implement all the stuff at once. The plan to finish the module should be follow

# Steps to Implement public APIs
1. Prepare the const.go, error.go and enum.go, if it's already defined, go to the next (no need to defined all stuff in the begining because you need to make sure the coverage is over 80%)
2. Prepare the basic function of sending request (skip if you already done)
    2.1. Define the function to send public request in request.go, without implementation
    2.2. Writing unittest of it
    2.3. Implement the function to send public request in request.go
    2.4. Run the unittest to make sure it works
3. Integrate the API path you want to implement in client.go, without implementation (skip if you already done)
    3.1. If the path is not defined, go to define it in const.go
    3.2. If the enum is not defined, go to define it in enum.go
    3.3. If the model is not defines, go to define it in models.go
    3.4. Define the function to send public request in request.go, without implementation
    3.5. Writing the unittest of it
    3.6. Implement the public API
    3.7. Run the unittest

# Steps to Implement signed APIs 
1. Define the const.go, error.go and enum.go, if it's already defined, go to the next (no need to defined all stuff in the begining because you need to make sure the coverage is over 80%)
2. Prepare the basic function of sending signed request 
    2.1. Define the function to send signed request in request.go, without implementation
    2.2. Writing unittest of it
    2.3. Implement the function to send signed request in request.go
    2.4. Run the unittest to make sure it works
3. Integrate the API path you want to implement in client.go, without implementation (skip if you already done)
    3.1. If the path is not defined, go to define it in const.go
    3.2. If the enum is not defined, go to define it in enum.go
    3.3. If the model is not defines, go to define it in models.go
    3.4. Define the function to send signed request in request.go, without implementation
    3.5. Writing the unittest of it
    3.6. Implement the signed API
    3.7. Run the unittest

# Steps to implement unsigned websocket subscription/resubscription/unsubscription

1. Prepare the const.go, error.go and enum.go, if it's already defined, go to the next (no need to defined all stuff in the begining because you need to make sure the coverage is over 80%)
2. Prepare the basic function to subscribe a raw websocket stream (skip if you already done)
    2.1. Define the websocket client include, establish connection, ping/pong, exponential backoff retry if disconnected, graceful disconnet, error handling
    2.2. Writing unittest of it
    2.3. Implement the function to subscribe public stream
    2.4. Run the unittest to make sure it works
3. Integrate the API path you want to implement in client.go, without implementation (skip if you already done)
    3.1. If the model is not defined, go to define it in ws_models.go
    3.2. Define the function to subscribe, without implementation
    3.3. Writing the unittest of it
    3.4. Implement the public subscription
    3.5. Run the unittest


# Steps to implement signed websocket subscription/resubscription/unsubscription

1. Prepare the const.go, error.go and enum.go, if it's already defined, go to the next (no need to defined all stuff in the begining because you need to make sure the coverage is over 80%)
2. Prepare the basic function to subscribe a signed websocket stream (skip if you already done)
    2.0. Implement a function to get a listen key
    2.1. Define the websocket client include, establish connection, ping/pong, exponential backoff retry if disconnected, graceful disconnet, error handling
    2.2. Writing unittest of it
    2.3. Implement the function to subscribe public stream
    2.4. Run the unittest to make sure it works
3. Integrate the API path you want to implement in client.go, without implementation (skip if you already done)
    3.1. If the model is not defined, go to define it in ws_models.go
    3.2. Define the function to subscribe, without implementation
    3.3. Writing the unittest of it
    3.4. Implement the public subscription
    3.5. Run the unittest

## Purpose

- Provide **reliable, scalable Binance Spot connectivity**.
- Align with Binance’s official API behavior while providing a clean, testable Go interface.
- Centralize:
  - Base URLs
  - API path definitions
  - WebSocket subscription rules
  - Error code mapping


## Structure

All core implementation files are located under:


**File Map:**

- `config.go` — Binance API key, secret, base URL configuration.
- `const.go` — Base URLs, REST/WS paths, header keys, query parameters.
    - https://developers.binance.com/docs/binance-spot-api-docs/enums
- `error.go` — Binance API error codes, error parsing utilities.
    - https://developers.binance.com/docs/binance-spot-api-docs/errors
- `models.go` — Request/response structures for REST endpoints.
- `request.go` — Low-level HTTP request management (signed/unsigned).
- `client.go` — High-level API wrappers for REST endpoints.
    - https://developers.binance.com/docs/binance-spot-api-docs/rest-api
- `client_test.go` — Unit tests for REST APIs.
- `ws_models.go` — Request/response structures for WebSocket channels.
- `ws.go` — Low-level WebSocket connection, subscription handling, reconnections.
- `ws_client.go` — High-level WebSocket APIs for subscribing to market data and user data streams.
- `utils.go` — Miscellaneous helper function can be re-use under binance
---

## References

### Official Documentation

- Binance REST API:
  - [Spot API Docs](https://binance-docs.github.io/apidocs/spot/en/)
  - [Spot Testnet](https://testnet.binance.vision/)
- Binance WebSocket Streams:
  - [WebSocket Streams](https://binance-docs.github.io/apidocs/spot/en/#websocket-market-streams)

### Base URLs

- **REST Production:**
  - `https://api.binance.com`
  - `https://api.binance.com`
  - `https://api-gcp.binance.com`
- **REST Testnet:**
  - `https://testnet.binance.vision`
- **WebSocket Production:**
  - `wss://stream.binance.com:9443/ws`
- **WebSocket Testnet:**
  - `wss://testnet.binance.vision/ws`

All base URLs are defined in `const.go`.


### API Paths

REST API endpoint paths are defined in `const.go` for consistency, e.g.:

- `/api/v3/account` — Get account information
- `/api/v3/order` — Create/cancel orders
- `/api/v3/myTrades` — Account trade list

Add new paths here when new endpoints are supported to keep agent tools and LLM agents aware of available actions.


### WebSocket Connection Rules

Binance WebSocket connectivity is handled under:

- `ws.go` (low-level connection)
- `ws_client.go` (high-level subscription interfaces)

**Rules:**

- Connection URL depends on production/testnet mode.
- Implement **ping/pong** keepalive.
- Binance provide two ways to subscribe/unsubscribe data: use subscribe/unsubscribe request or subscribe a raw stream, here the implemeation always use raw stream
- The public API of a websocket subscription should looks like that:

```go
func (c *WSStreamClient) SubscribeToKline(symbol string, interval string, options *KlineSubscriptionOptions) (unsubscribe func() error, error)
```

The KlineSubscriptionOptions define the chain methods for user include:
- WithConnect(callback func())
- WithReconnect(callback func())
- WithKline(callback KlineCallback)
- WithDisconnect(callback func())
- WithError(callback func(error))

So user can utilize the subscriptionOption.WithConnect(callback).WithKline(klineCallback) to focus what they really want of it.

Example usage:
```go
klineOptions := &KlineSubscriptionOptions{}
klineOptions.WithConnect(func() {
    fmt.Println("Connected to kline stream")
}).WithKline(func(data *WSKlineData) error {
    fmt.Printf("Kline: Symbol=%s, Close=%f\n", data.Symbol, data.Kline.ClosePrice)
    return nil
}).WithError(func(err error) {
    fmt.Printf("Error: %v\n", err)
})

unsubscribe, err := wsClient.SubscribeToKline("BTCUSDT", "1m", klineOptions)
```

If there are multiple event with the corresponding subscription, you can define the SubscriptionOptions with different pre-defined event handler like:
- WithExecutionReport(ExecutionReportCallback)
- WithAccountUpdate(OutboundAccountPositionCallback)
- WithBalanceUpdate(BalanceUpdateCallback)

This pattern makes developer easier to extend their implementation and provides type-safe callbacks for different data types.

