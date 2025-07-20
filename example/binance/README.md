# Binance Connectivity Test Example

This example demonstrates comprehensive testing of Binance API connectivity and trading functionality.

## Features Tested

1. **Account Balance** - Retrieves and displays account information and balances
2. **Price Checking** - Gets current ADAUSDT price
3. **User Data Stream** - Subscribes to real-time user data updates
4. **Limit Buy Order** - Places a limit buy order at 90% of current price
5. **Order Cancellation** - Cancels the previously placed limit order
6. **Market Buy Order** - Places a market buy order for ADA
7. **Market Sell Order** - Places a market sell order for ADA
8. **Stream Unsubscription** - Properly closes the user data stream

## Prerequisites

1. **Binance API Credentials**: You need API key and secret from Binance
   - For testing: Use Binance Testnet (https://testnet.binance.vision/)
   - For production: Use Binance mainnet

2. **Environment Variables**: Set the following environment variables:
   ```bash
   export BINANCE_API_KEY="your_api_key_here"
   export BINANCE_API_SECRET="your_api_secret_here"
   ```

## Running the Example

1. **Navigate to the example directory**:
   ```bash
   cd example/binance
   ```

2. **Run the example**:
   ```bash
   go run binance_example.go
   ```

## Configuration

The example uses Binance Testnet by default for safety. To switch to production:

```go
// Change this line in the code:
config := binance.DefaultConfig()  // For production
// Instead of:
config := binance.TestnetConfig()  // For testnet
```

## Test Flow

1. **Account Balance Check**: Verifies API connectivity and displays account information
2. **Price Retrieval**: Gets current ADAUSDT price for order calculations
3. **WebSocket Connection**: Establishes user data stream for real-time updates
4. **Limit Order Test**: Places a limit buy order at 90% of current price
5. **Order Cancellation**: Demonstrates order cancellation functionality
6. **Market Orders**: Tests both market buy and sell orders
7. **Cleanup**: Properly closes the WebSocket connection

## Important Notes

- **Testnet Usage**: The example uses testnet by default to avoid real trading
- **Small Quantities**: Orders use small quantities (1 ADA) for testing
- **Error Handling**: Comprehensive error handling with detailed logging
- **WebSocket Management**: Proper connection lifecycle management
- **Order Types**: Demonstrates both limit and market orders

## Expected Output

The example will output detailed logs for each test step, including:
- Account information and balances
- Current ADAUSDT price
- WebSocket connection status
- Order placement confirmations
- Order cancellation confirmations
- Real-time data stream events

## Troubleshooting

1. **API Key Issues**: Ensure your API key has trading permissions
2. **Network Issues**: Check your internet connection and firewall settings
3. **Rate Limits**: Binance has rate limits; the example includes delays to respect them
4. **Testnet Balance**: Ensure you have test ADA and USDT in your testnet account

## Safety Warnings

⚠️ **Important**: 
- This example uses testnet by default for safety
- If you switch to production, be aware that real money will be used
- Always test thoroughly on testnet before using production
- Monitor your orders and account balance carefully
- The example uses small quantities, but real trading involves financial risk 