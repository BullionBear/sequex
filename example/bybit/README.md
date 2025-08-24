# Bybit API Example

This example demonstrates how to use the Bybit API client to retrieve market data.

## Features Demonstrated

- **Server Time**: Get current server time
- **Kline Data**: Retrieve candlestick data for BTCUSD inverse perpetual
- **Parsed Kline Data**: Get structured kline data with parsed timestamps and prices
- **Ticker Information**: Get current ticker data including prices and volumes
- **Recent Data**: Get recent kline data with limit parameter
- **Account Information**: Get account balance and PnL information
- **Trading**: Create, query, and cancel orders (limit and market orders)

## Running the Example

```bash
go run bybit_example.go
```

## Expected Output

The example will output:
- Current server time
- Kline data for BTCUSD inverse perpetual
- Parsed kline data with timestamps and OHLCV values
- Current ticker information
- Recent kline data (last 10 records)

## API Endpoints Used

- `GET /v5/market/time` - Server time
- `GET /v5/market/kline` - Kline/candlestick data
- `GET /v5/market/tickers` - Ticker information
- `GET /v5/account/wallet-balance` - Account information (signed)
- `POST /v5/order/create` - Create order (signed)
- `POST /v5/order/cancel` - Cancel order (signed)
- `GET /v5/order/realtime` - Get order information (signed)

## Configuration

The example uses the testnet configuration with API credentials:
- Base URL: `https://api-testnet.bybit.com`
- API credentials required for signed endpoints
- 30-second timeout

**Note**: The example uses testnet API credentials from environment variables:
- `BYBIT_TESTNET_API_KEY`
- `BYBIT_TESTNET_API_SECRET`

## Parameters

### Kline Request Parameters
- `category`: "inverse" (for inverse perpetual contracts)
- `symbol`: "BTCUSD" (Bitcoin/USD pair)
- `interval`: "60" (1 hour intervals)
- `start`: Start timestamp in milliseconds
- `end`: End timestamp in milliseconds
- `limit`: Maximum number of records to return

### Supported Categories
- `spot` - Spot trading
- `linear` - USDT/USDC perpetual
- `inverse` - Inverse perpetual
- `option` - Options trading

### Supported Intervals
- `1`, `3`, `5`, `15`, `30` (minutes)
- `60`, `120`, `240`, `360`, `480`, `720` (hours)
- `D`, `W`, `M` (day, week, month)

### Trading Parameters

#### Order Types
- `Market` - Market order (executed immediately at market price)
- `Limit` - Limit order (executed at specified price or better)

#### Order Sides
- `Buy` - Buy order
- `Sell` - Sell order

#### Time in Force
- `GTC` - Good Till Cancel
- `IOC` - Immediate or Cancel
- `FOK` - Fill or Kill

#### Quantity Format
- **Inverse Perpetual**: Quantity is in USD contract value (e.g., "100" = $100 contract value)
- **Linear Perpetual**: Quantity is in contract units (e.g., "0.01" = 0.01 BTC)
- **Spot**: Quantity is in base currency (e.g., "0.01" = 0.01 BTC) 