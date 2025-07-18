# Binance API Client

A comprehensive Go client for the Binance REST API with independent HTTP session implementation, HMAC-SHA256 authentication, and well-organized folder structure.

## Features

- ✅ **Independent HTTP Client**: No external dependencies like go-binance
- ✅ **HMAC-SHA256 Authentication**: Secure API request signing
- ✅ **Comprehensive API Coverage**: Market data, account info, trading operations
- ✅ **Configuration Management**: YAML-based configuration with multiple accounts
- ✅ **Error Handling**: Proper error types with retry logic
- ✅ **Testnet Support**: Sandbox environment for testing
- ✅ **Extensive Testing**: Unit tests with 62%+ coverage
- ✅ **Type Safety**: Decimal precision for financial calculations

## Quick Start

### Basic Usage

```go
package main

import (
    "context"
    "fmt"
    "log"
    
    "github.com/BullionBear/sequex/pkg/exchange/binance"
)

func main() {
    // Create configuration
    config := &binance.Config{
        APIKey:    "your_api_key",
        APISecret: "your_api_secret",
        Sandbox:   true, // Use testnet
        Timeout:   30,
    }

    // Create client
    client, err := binance.NewClient(config)
    if err != nil {
        log.Fatal(err)
    }

    ctx := context.Background()

    // Test connectivity
    if err := client.Ping(ctx); err != nil {
        log.Fatal("Connection failed:", err)
    }

    // Get market data
    ticker, err := client.GetTicker24hr(ctx, "BTCUSDT")
    if err != nil {
        log.Fatal(err)
    }
    
    fmt.Printf("BTC Price: %s USDT\n", ticker.LastPrice)
}
```

### Configuration File

Create a `config.yml` file:

```yaml
accounts:
  binance:
    - name: "trading_account"
      api_key: "your_api_key"
      api_secret: "your_api_secret"
      sandbox: true
      timeout: 30

market:
  binance:
    - "BTCUSDT"
    - "ETHUSDT"
    - "ADAUSDT"
```

Load and use configuration:

```go
// Load configuration
appConfig, err := binance.LoadConfig("config.yml")
if err != nil {
    log.Fatal(err)
}

// Get Binance config
config, err := appConfig.GetBinanceConfig()
if err != nil {
    log.Fatal(err)
}

// Create client
client, err := binance.NewClient(config)
if err != nil {
    log.Fatal(err)
}
```

## API Methods

### Public Endpoints (No Authentication Required)

```go
// Server connectivity
err := client.Ping(ctx)

// Server time
serverTime, err := client.GetServerTime(ctx)

// Exchange information
exchangeInfo, err := client.GetExchangeInfo(ctx)

// 24hr ticker statistics
ticker, err := client.GetTicker24hr(ctx, "BTCUSDT")

// Order book depth
orderBook, err := client.GetOrderBook(ctx, "BTCUSDT", 10)

// Recent trades
trades, err := client.GetRecentTrades(ctx, "BTCUSDT", 5)

// Kline/candlestick data
klines, err := client.GetKlines(ctx, "BTCUSDT", "1h", 24, nil, nil)
```

### Authenticated Endpoints (Require API Credentials)

```go
// Account information
account, err := client.GetAccount(ctx)

// Get open orders
orders, err := client.GetOpenOrders(ctx, "BTCUSDT")

// Get order status
orderID := int64(12345)
order, err := client.GetOrder(ctx, "BTCUSDT", &orderID, nil)

// Place new order
quantity := "0.001"
price := "50000"
timeInForce := "GTC"
response, err := client.CreateOrder(ctx, "BTCUSDT", "BUY", "LIMIT", &quantity, &price, &timeInForce)

// Cancel order
cancelledOrder, err := client.CancelOrder(ctx, "BTCUSDT", &orderID, nil)

// Get trade history
trades, err := client.GetTrades(ctx, "BTCUSDT", 10, nil)
```

## Configuration

### Config Structure

```go
type Config struct {
    Name      string `yaml:"name" json:"name"`
    APIKey    string `yaml:"api_key" json:"api_key"`
    APISecret string `yaml:"api_secret" json:"api_secret"`
    Sandbox   bool   `yaml:"sandbox" json:"sandbox"`
    Timeout   int    `yaml:"timeout" json:"timeout"` // in seconds
}
```

### Default Configuration

```go
config := binance.DefaultConfig()
// Returns:
// Name: "default"
// Sandbox: false
// Timeout: 30
```

## Error Handling

The client provides structured error handling:

```go
ticker, err := client.GetTicker24hr(ctx, "INVALID_SYMBOL")
if err != nil {
    if apiErr, ok := err.(*binance.APIError); ok {
        fmt.Printf("API Error %d: %s\n", apiErr.Code, apiErr.Msg)
        
        // Check if error is retryable
        if binance.IsRetryableError(apiErr) {
            // Implement retry logic
        }
    } else {
        // Handle other errors (network, parsing, etc.)
        fmt.Printf("Error: %v\n", err)
    }
}
```

## Testing

### Run All Tests

```bash
# Using the test runner script
./test_runner.sh

# Or directly with go test
go test -v ./pkg/exchange/binance/...
```

### Test with Real Credentials

Set environment variables for authenticated endpoint testing:

```bash
export BINANCE_TESTNET_API_KEY="your_testnet_api_key"
export BINANCE_TESTNET_API_SECRET="your_testnet_api_secret"
./test_runner.sh
```

### Test Coverage

```bash
go test -coverprofile=coverage.out ./pkg/exchange/binance/...
go tool cover -html=coverage.out
```

## Project Structure

```
pkg/exchange/binance/
├── client.go          # Main HTTP client with API methods
├── config.go          # Configuration structures
├── config_loader.go   # YAML configuration loader
├── const.go           # Constants (URLs, etc.)
├── errors.go          # Error types and handling
├── models.go          # Data structures for API responses
├── utils.go           # Utility functions
├── *_test.go          # Unit tests
└── README.md          # This documentation
```

## Key Features

### Type Safety with Decimal Precision

All price and quantity fields use `decimal.Decimal` for precise financial calculations:

```go
ticker, _ := client.GetTicker24hr(ctx, "BTCUSDT")
price := ticker.LastPrice // decimal.Decimal type

// Safe arithmetic operations
priceInCents := price.Mul(decimal.NewFromInt(100))
```

### Flexible Configuration

Support for multiple accounts and environments:

```go
// Get specific account by name
config, err := appConfig.GetBinanceConfigByName("trading_account")

// Get configured symbols
symbols := appConfig.GetBinanceSymbols()
```

### Production Ready

- Proper context handling for timeouts and cancellation
- Rate limiting awareness with retryable error detection  
- Comprehensive error handling and logging
- Testnet support for safe development

## Example: Complete Trading Bot Setup

```go
package main

import (
    "context"
    "log"
    "time"
    
    "github.com/BullionBear/sequex/pkg/exchange/binance"
)

func main() {
    // Load configuration
    appConfig, err := binance.LoadConfig("config.yml")
    if err != nil {
        log.Fatal("Failed to load config:", err)
    }

    config, err := appConfig.GetBinanceConfig()
    if err != nil {
        log.Fatal("Failed to get Binance config:", err)
    }

    // Create client
    client, err := binance.NewClient(config)
    if err != nil {
        log.Fatal("Failed to create client:", err)
    }

    ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
    defer cancel()

    // Check connectivity
    if err := client.Ping(ctx); err != nil {
        log.Fatal("Failed to ping:", err)
    }

    // Get account info
    account, err := client.GetAccount(ctx)
    if err != nil {
        log.Fatal("Failed to get account:", err)
    }

    log.Printf("Account can trade: %v", account.CanTrade)
    log.Printf("Account has %d balances", len(account.Balances))

    // Get market data for configured symbols
    symbols := appConfig.GetBinanceSymbols()
    for _, symbol := range symbols {
        ticker, err := client.GetTicker24hr(ctx, symbol)
        if err != nil {
            log.Printf("Failed to get ticker for %s: %v", symbol, err)
            continue
        }
        
        log.Printf("%s: Price=%s, Change=%s%%", 
            ticker.Symbol, 
            ticker.LastPrice.String(), 
            ticker.PriceChangePercent.String())
    }
}
```

## Security Notes

- Never commit API credentials to version control
- Use environment variables or secure configuration management
- Always use testnet for development and testing
- Implement proper rate limiting in production applications
- Validate all input parameters before making API calls

## License

This project is part of the Sequex trading system. 