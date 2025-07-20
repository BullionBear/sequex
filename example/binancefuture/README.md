# Binance Futures Example

This example demonstrates how to use the Binance Futures API client to perform various trading operations.

## Features

The example includes the following functionality:

1. **Get Account Balance** - Retrieve account information including balances and positions
2. **Get ADAUSDT Price** - Get current market price for ADAUSDT
3. **Subscribe User Data Stream** - Subscribe to real-time user data updates
4. **Send Limit Buy Order** - Place a limit buy order at 0.9 * current price
5. **Get Open Orders** - Retrieve all open orders for ADAUSDT
6. **Cancel Order** - Cancel the previously placed limit order
7. **Send Market Buy Order** - Place a market buy order
8. **Get Open Positions** - Retrieve current positions for ADAUSDT
9. **Send Market Sell Order** - Place a market sell order
10. **Unsubscribe User Data Stream** - Close the user data stream

## Prerequisites

1. **Binance Futures Account**: You need a Binance Futures account with API access
2. **API Credentials**: Generate API key and secret from your Binance Futures account
3. **Go Environment**: Make sure you have Go 1.19+ installed

## Setup

### 1. Set Environment Variables

Set your Binance Futures API credentials as environment variables:

```bash
export BINANCE_API_KEY="your_api_key_here"
export BINANCE_API_SECRET="your_api_secret_here"
```

### 2. Configure API Permissions

Make sure your API key has the following permissions enabled:
- **Enable Futures**: Required for futures trading
- **Enable Reading**: Required for account and market data
- **Enable Spot & Margin Trading**: Required for trading operations

### 3. Testnet (Recommended for Testing)

For testing, you can use Binance Futures testnet:

1. Create a testnet account at [Binance Futures Testnet](https://testnet.binancefuture.com/)
2. Generate testnet API credentials
3. Update the configuration in the example to use testnet URLs

## Running the Example

### From the project root:

```bash
go run example/binancefuture/binancefuture_example.go
```

### From the example directory:

```bash
cd example/binancefuture
go run binancefuture_example.go
```

## Example Output

The example will output detailed information for each operation:

```
Starting Binance Futures connectivity test...

=== Test 1: Get Account Balance ===
Account Type: UNIFIED
Can Trade: true
Can Withdraw: true
Can Deposit: true
Update Time: 1703123456789
Total Wallet Balance: 1000.00000000
Available Balance: 1000.00000000
Assets:
  USDT: Wallet Balance=1000.00000000, Available Balance=1000.00000000
  ADA: Wallet Balance=0.00000000, Available Balance=0.00000000
Positions:
  ADAUSDT: Position Amount=0, Entry Price=0.0000, Mark Price=0.5000, Unrealized PnL=0.00000000

=== Test 2: Get ADAUSDT Price ===
ADAUSDT Current Price: $0.5000

=== Test 3: Subscribe to User Data Stream ===
Created user data stream with listen key: abc123def456...
Connected to user data stream

=== Test 4: Send Limit Buy Order ===
Placing limit buy order: BUY 10.0 ADA at $0.4500
Limit buy order placed successfully:
  Order ID: 123456789
  Status: NEW
  Price: 0.4500
  Quantity: 10.0

...
```

## Important Notes

### Risk Warning
- **This is a trading example that places real orders**
- **Use testnet for testing to avoid real money losses**
- **Always test thoroughly before using with real funds**

### Order Quantities
- The example uses a quantity of 10 ADA for orders
- Adjust the quantity based on your available balance and risk tolerance
- Ensure you meet minimum notional requirements

### Error Handling
- The example includes comprehensive error handling
- Check the logs for any API errors or rate limit issues
- Some operations may fail if you don't have sufficient balance

### Rate Limits
- Binance Futures has rate limits on API calls
- The example includes delays between operations to respect rate limits
- For production use, implement proper rate limiting

## Configuration Options

You can modify the configuration in the example:

```go
config := binancefuture.DefaultConfig()

// Use testnet
config.BaseURL = binancefuture.BaseURLFuturesTestnet
config.WSBaseURL = binancefuture.WSBaseURLTestnet

// Set custom timeouts
config.RequestTimeout = 30 * time.Second
config.RecvWindow = 5000
```

## Troubleshooting

### Common Issues

1. **API Key Error**: Ensure your API key and secret are correct
2. **Insufficient Balance**: Make sure you have enough USDT for trading
3. **Rate Limit**: Wait and retry if you hit rate limits
4. **Network Issues**: Check your internet connection
5. **Symbol Not Found**: Ensure ADAUSDT is available for trading

### Debug Mode

Enable debug logging by setting the log level:

```go
log.SetFlags(log.LstdFlags | log.Lshortfile)
```

## API Documentation

For more information about the Binance Futures API:
- [Binance Futures API Documentation](https://binance-docs.github.io/apidocs/futures/en/)
- [WebSocket Streams](https://binance-docs.github.io/apidocs/futures/en/#websocket-market-streams)
- [User Data Streams](https://binance-docs.github.io/apidocs/futures/en/#user-data-streams) 