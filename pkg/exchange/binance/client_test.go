package binance

import (
	"context"
	"os"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func getTestConfig() *Config {
	config := DefaultConfig()
	config.Sandbox = true
	config.Timeout = 10

	// Use hardcoded test credentials for unit tests
	config.APIKey = "NapwpjlsGTEgmc4cMgj8oA7zHzuUeAgRRj5hu0ZAKXA6XFYl3KYgiQd3YV9eVYrb"
	config.APISecret = "x3zyS1epGBKsz7KT4TqzIGFCKPLmsdFnEkt5EUioTgasgGaj8uXzhqDXIsRrXDMc"

	// Override with real credentials from environment variables if available
	if apiKey := os.Getenv("BINANCE_TESTNET_API_KEY"); apiKey != "" {
		config.APIKey = apiKey
	}
	if apiSecret := os.Getenv("BINANCE_TESTNET_API_SECRET"); apiSecret != "" {
		config.APISecret = apiSecret
	}

	return config
}

func hasValidCredentials(config *Config) bool {
	return config.APIKey != "" && config.APISecret != "" &&
		config.APIKey != "your_api_key_here" && config.APISecret != "your_api_secret_here"
}

func hasRealCredentials(config *Config) bool {
	return hasValidCredentials(config) &&
		config.APIKey != "test_binance_api_key_for_unittest" &&
		config.APISecret != "test_binance_api_secret_for_unittest"
}

func TestNewClient(t *testing.T) {
	t.Run("ValidConfig", func(t *testing.T) {
		config := getTestConfig()
		client, err := NewClient(config)

		require.NoError(t, err)
		assert.NotNil(t, client)
		assert.Equal(t, config, client.GetConfig())
		assert.NotNil(t, client.GetHTTPClient())
	})

	t.Run("ValidConfigWithoutCredentials", func(t *testing.T) {
		config := &Config{
			Sandbox: true,
			Timeout: 10,
		}
		client, err := NewClient(config)

		require.NoError(t, err)
		assert.NotNil(t, client)
		assert.Equal(t, config, client.GetConfig())
		assert.NotNil(t, client.GetHTTPClient())
	})

	t.Run("NilConfig", func(t *testing.T) {
		client, err := NewClient(nil)

		require.NoError(t, err)
		assert.NotNil(t, client)
		assert.NotNil(t, client.GetConfig())
	})

	t.Run("ValidConfigWithoutCredentials", func(t *testing.T) {
		config := &Config{
			Name:    "test",
			Timeout: 30,
			Sandbox: true,
			// No API credentials - should still create client for public endpoints
		}

		client, err := NewClient(config)

		assert.NoError(t, err)
		assert.NotNil(t, client)
		assert.Equal(t, config, client.GetConfig())
	})
}

func TestConfigValidation(t *testing.T) {
	t.Run("ValidConfig", func(t *testing.T) {
		config := &Config{
			APIKey:    "test_key",
			APISecret: "test_secret",
		}

		assert.True(t, config.IsValid())
	})

	t.Run("EmptyAPIKey", func(t *testing.T) {
		config := &Config{
			APIKey:    "",
			APISecret: "test_secret",
		}

		assert.False(t, config.IsValid())
	})

	t.Run("EmptyAPISecret", func(t *testing.T) {
		config := &Config{
			APIKey:    "test_key",
			APISecret: "",
		}

		assert.False(t, config.IsValid())
	})
}

func TestDefaultConfig(t *testing.T) {
	config := DefaultConfig()

	assert.Equal(t, "default", config.Name)
	assert.False(t, config.Sandbox)
	assert.Equal(t, 30, config.Timeout)
}

func TestGetBaseURL(t *testing.T) {
	t.Run("ProductionURL", func(t *testing.T) {
		config := &Config{Sandbox: false}
		assert.Equal(t, BaseURL, config.GetBaseURL())
	})

	t.Run("SandboxURL", func(t *testing.T) {
		config := &Config{Sandbox: true}
		assert.Equal(t, SandboxBaseURL, config.GetBaseURL())
	})
}

func TestPing(t *testing.T) {
	config := getTestConfig()
	client, err := NewClient(config)
	require.NoError(t, err)

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	err = client.Ping(ctx)
	assert.NoError(t, err)
}

func TestGetServerTime(t *testing.T) {
	config := getTestConfig()
	client, err := NewClient(config)
	require.NoError(t, err)

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	serverTime, err := client.GetServerTime(ctx)
	assert.NoError(t, err)
	assert.Greater(t, serverTime, int64(0))

	// Server time should be recent (within last hour)
	now := time.Now().UnixMilli()
	assert.InDelta(t, now, serverTime, float64(time.Hour.Milliseconds()))
}

func TestGetExchangeInfo(t *testing.T) {
	config := getTestConfig()
	client, err := NewClient(config)
	require.NoError(t, err)

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	exchangeInfo, err := client.GetExchangeInfo(ctx)
	assert.NoError(t, err)
	assert.NotNil(t, exchangeInfo)
	assert.Greater(t, len(exchangeInfo.Symbols), 0)
	assert.NotEmpty(t, exchangeInfo.Timezone)
}

func TestGetTicker24hr(t *testing.T) {
	config := getTestConfig()
	client, err := NewClient(config)
	require.NoError(t, err)

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	ticker, err := client.GetTicker24hr(ctx, "BTCUSDT")
	assert.NoError(t, err)
	assert.NotNil(t, ticker)
	assert.Equal(t, "BTCUSDT", ticker.Symbol)
	assert.True(t, ticker.LastPrice.IsPositive())
}

func TestGetOrderBook(t *testing.T) {
	config := getTestConfig()
	client, err := NewClient(config)
	require.NoError(t, err)

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	orderBook, err := client.GetOrderBook(ctx, "BTCUSDT", 5)
	assert.NoError(t, err)
	assert.NotNil(t, orderBook)
	assert.LessOrEqual(t, len(orderBook.Bids), 5)
	assert.LessOrEqual(t, len(orderBook.Asks), 5)
	assert.Greater(t, len(orderBook.Bids), 0)
	assert.Greater(t, len(orderBook.Asks), 0)
}

func TestGetRecentTrades(t *testing.T) {
	config := getTestConfig()
	client, err := NewClient(config)
	require.NoError(t, err)

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	trades, err := client.GetRecentTrades(ctx, "BTCUSDT", 5)
	assert.NoError(t, err)
	assert.NotNil(t, trades)
	assert.LessOrEqual(t, len(trades), 5)

	if len(trades) > 0 {
		trade := trades[0]
		assert.Greater(t, trade.ID, int64(0))
		assert.True(t, trade.Price.IsPositive())
		assert.True(t, trade.Qty.IsPositive())
	}
}

func TestGetKlines(t *testing.T) {
	config := getTestConfig()
	client, err := NewClient(config)
	require.NoError(t, err)

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	klines, err := client.GetKlines(ctx, "BTCUSDT", "1h", 5, nil, nil)
	assert.NoError(t, err)
	assert.NotNil(t, klines)
	assert.LessOrEqual(t, len(klines), 5)

	if len(klines) > 0 {
		kline := klines[0]
		assert.Greater(t, kline.OpenTime, int64(0))
		assert.True(t, kline.Open.IsPositive())
		assert.True(t, kline.High.IsPositive())
		assert.True(t, kline.Low.IsPositive())
		assert.True(t, kline.Close.IsPositive())
	}
}

// Authenticated endpoint tests - only run if real credentials are provided
func TestAuthenticatedEndpoints(t *testing.T) {
	config := getTestConfig()

	if !hasRealCredentials(config) {
		t.Skip("Skipping authenticated tests: no real API credentials provided. " +
			"Set BINANCE_TESTNET_API_KEY and BINANCE_TESTNET_API_SECRET environment variables to run these tests.")
	}

	client, err := NewClient(config)
	require.NoError(t, err)

	ctx, cancel := context.WithTimeout(context.Background(), 15*time.Second)
	defer cancel()

	t.Run("GetAccount", func(t *testing.T) {
		account, err := client.GetAccount(ctx)
		assert.NoError(t, err)
		assert.NotNil(t, account)
		assert.NotNil(t, account.Balances)
	})

	t.Run("GetOpenOrders", func(t *testing.T) {
		orders, err := client.GetOpenOrders(ctx, "")
		assert.NoError(t, err)
		assert.NotNil(t, orders)
		// Orders slice can be empty, that's fine
	})

	t.Run("GetOpenOrdersForSymbol", func(t *testing.T) {
		orders, err := client.GetOpenOrders(ctx, "BTCUSDT")
		assert.NoError(t, err)
		assert.NotNil(t, orders)
		// Orders slice can be empty, that's fine
	})

	t.Run("GetTrades", func(t *testing.T) {
		trades, err := client.GetTrades(ctx, "BTCUSDT", 5, nil)
		assert.NoError(t, err)
		assert.NotNil(t, trades)
		// Trades slice can be empty if no trades, that's fine
	})
}

func TestSignatureGeneration(t *testing.T) {
	config := &Config{
		APIKey:    "test_key",
		APISecret: "test_secret",
	}

	client, err := NewClient(config)
	require.NoError(t, err)

	signature := client.generateSignature("symbol=BTCUSDT&timestamp=1234567890")
	assert.NotEmpty(t, signature)
	assert.Len(t, signature, 64) // HMAC-SHA256 produces 64 character hex string

	// Same input should produce same signature
	signature2 := client.generateSignature("symbol=BTCUSDT&timestamp=1234567890")
	assert.Equal(t, signature, signature2)

	// Different input should produce different signature
	signature3 := client.generateSignature("symbol=ETHUSDT&timestamp=1234567890")
	assert.NotEqual(t, signature, signature3)
}

func TestContextCancellation(t *testing.T) {
	config := getTestConfig()
	client, err := NewClient(config)
	require.NoError(t, err)

	ctx, cancel := context.WithCancel(context.Background())
	cancel() // Cancel immediately

	err = client.Ping(ctx)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "context canceled")
}

func TestInvalidSymbol(t *testing.T) {
	config := getTestConfig()
	client, err := NewClient(config)
	require.NoError(t, err)

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	_, err = client.GetTicker24hr(ctx, "INVALIDPAIR")
	assert.Error(t, err)
}
