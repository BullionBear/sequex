package bybit

import (
	"context"
	"os"
	"testing"
	"time"
)

func TestGetKline(t *testing.T) {
	// Create client with default config (mainnet)
	config := DefaultConfig()
	client := NewClient(config)

	// Create kline request with the exact parameters from the example
	req := &KlineRequest{
		Category: "inverse",
		Symbol:   "BTCUSD",
		Interval: "60",          // 1 hour interval
		Start:    1670601600000, // Start time from example
		End:      1670608800000, // End time from example
	}

	// Make the request
	ctx := context.Background()
	klineResp, err := client.GetKline(ctx, req)
	if err != nil {
		t.Fatalf("Failed to get kline data: %v", err)
	}

	// Verify response structure
	if klineResp.RetCode != 0 {
		t.Fatalf("Expected retCode 0, got %d: %s", klineResp.RetCode, klineResp.RetMsg)
	}

	if klineResp.Result.Symbol != "BTCUSD" {
		t.Errorf("Expected symbol BTCUSD, got %s", klineResp.Result.Symbol)
	}

	if klineResp.Result.Category != "inverse" {
		t.Errorf("Expected category inverse, got %s", klineResp.Result.Category)
	}

	// Verify we got some kline data
	if len(klineResp.Result.List) == 0 {
		t.Error("Expected kline data, got empty list")
	}

	t.Logf("Successfully retrieved %d kline records", len(klineResp.Result.List))

	// Test parsing kline data
	klineData, err := client.GetKlineData(ctx, req)
	if err != nil {
		t.Fatalf("Failed to get parsed kline data: %v", err)
	}

	if len(klineData) == 0 {
		t.Error("Expected parsed kline data, got empty list")
	}

	// Verify first kline data
	if len(klineData) > 0 {
		firstKline := klineData[0]
		t.Logf("First kline: Time=%s, Open=%.2f, High=%.2f, Low=%.2f, Close=%.2f, Volume=%.2f",
			firstKline.Timestamp.Format(time.RFC3339),
			firstKline.OpenPrice,
			firstKline.HighPrice,
			firstKline.LowPrice,
			firstKline.ClosePrice,
			firstKline.Volume)
	}
}

func TestGetServerTime(t *testing.T) {
	// Create client with default config
	config := DefaultConfig()
	client := NewClient(config)

	// Make the request
	ctx := context.Background()
	serverTimeResp, err := client.GetServerTime(ctx)
	if err != nil {
		t.Fatalf("Failed to get server time: %v", err)
	}

	// Verify response structure
	if serverTimeResp.RetCode != 0 {
		t.Fatalf("Expected retCode 0, got %d: %s", serverTimeResp.RetCode, serverTimeResp.RetMsg)
	}

	if serverTimeResp.Result.TimeSecond == "" {
		t.Error("Expected server time, got empty string")
	}

	t.Logf("Server time: %s", serverTimeResp.Result.TimeSecond)
}

func TestGetTickers(t *testing.T) {
	// Create client with default config
	config := DefaultConfig()
	client := NewClient(config)

	// Make the request
	ctx := context.Background()
	tickerResp, err := client.GetTickers(ctx, "inverse", "BTCUSD")
	if err != nil {
		t.Fatalf("Failed to get tickers: %v", err)
	}

	// Verify response structure
	if tickerResp.RetCode != 0 {
		t.Fatalf("Expected retCode 0, got %d: %s", tickerResp.RetCode, tickerResp.RetMsg)
	}

	if tickerResp.Result.Category != "inverse" {
		t.Errorf("Expected category inverse, got %s", tickerResp.Result.Category)
	}

	// Verify we got some ticker data
	if len(tickerResp.Result.List) == 0 {
		t.Error("Expected ticker data, got empty list")
	}

	t.Logf("Successfully retrieved %d ticker records", len(tickerResp.Result.List))

	// Log first ticker info
	if len(tickerResp.Result.List) > 0 {
		firstTicker := tickerResp.Result.List[0]
		t.Logf("First ticker: Symbol=%s, LastPrice=%s, Volume24h=%s",
			firstTicker.Symbol, firstTicker.LastPrice, firstTicker.Volume24h)
	}
}

func TestGetAccount(t *testing.T) {
	// Skip if no API credentials
	apiKey := os.Getenv("BYBIT_TESTNET_API_KEY")
	apiSecret := os.Getenv("BYBIT_TESTNET_API_SECRET")
	if apiKey == "" || apiSecret == "" {
		t.Skip("Skipping test - no API credentials provided")
	}

	// Create client with testnet config and API credentials
	config := TestnetConfig()
	config = config.WithAPIKey(apiKey).WithAPISecret(apiSecret)
	client := NewClient(config)

	// Make the request
	ctx := context.Background()
	accountReq := &GetAccountRequest{
		AccountType: AccountTypeUnified,
	}
	accountResp, err := client.GetAccount(ctx, accountReq)
	if err != nil {
		t.Fatalf("Failed to get account: %v", err)
	}

	// Verify response structure
	if accountResp.RetCode != 0 {
		t.Fatalf("Expected retCode 0, got %d: %s", accountResp.RetCode, accountResp.RetMsg)
	}

	// Verify we got account data
	if len(accountResp.Result.List) == 0 {
		t.Error("Expected account data, got empty list")
	}

	t.Logf("Successfully retrieved account information")
	if len(accountResp.Result.List) > 0 {
		account := accountResp.Result.List[0]
		t.Logf("Total Wallet Balance: %s", account.TotalWalletBalance)
		t.Logf("Total Available Balance: %s", account.TotalAvailableBalance)
		t.Logf("Total Equity: %s", account.TotalEquity)
		t.Logf("Account Type: %s", account.AccountType)
		if len(account.Coin) > 0 {
			coin := account.Coin[0]
			t.Logf("Coin: %s, Wallet Balance: %s", coin.Coin, coin.WalletBalance)
		}
	}
}

func TestCreateOrder(t *testing.T) {
	// Skip if no API credentials
	apiKey := os.Getenv("BYBIT_TESTNET_API_KEY")
	apiSecret := os.Getenv("BYBIT_TESTNET_API_SECRET")
	if apiKey == "" || apiSecret == "" {
		t.Skip("Skipping test - no API credentials provided")
	}

	// Create client with testnet config and API credentials
	config := TestnetConfig()
	config = config.WithAPIKey(apiKey).WithAPISecret(apiSecret)
	client := NewClient(config)

	// Create a small limit order (10 USD contract value)
	ctx := context.Background()
	createOrderReq := &CreateOrderRequest{
		Category:    "inverse",
		Symbol:      "BTCUSD",
		Side:        "Buy",
		OrderType:   "Limit",
		Qty:         "10",     // 10 USD contract value
		Price:       "100000", // Set price well below current market price
		TimeInForce: "GTC",
	}

	createOrderResp, err := client.CreateOrder(ctx, createOrderReq)
	if err != nil {
		// If we get insufficient balance error, that's expected for testnet
		if err.Error() == "API error: code=110007, msg=Insufficient available balance" {
			t.Logf("Expected insufficient balance error: %v", err)
			return
		}
		t.Fatalf("Failed to create order: %v", err)
	}

	// Verify response structure
	if createOrderResp.RetCode != 0 {
		t.Fatalf("Expected retCode 0, got %d: %s", createOrderResp.RetCode, createOrderResp.RetMsg)
	}

	t.Logf("Order created successfully: %s", createOrderResp.Result.OrderId)

	// Test getting order information (UTA 2.0)
	getOrderReq := &GetOrderRequest{
		Category: "inverse",
		Symbol:   "BTCUSD",
		OrderId:  createOrderResp.Result.OrderId,
	}

	getOrderResp, err := client.GetOrder(ctx, getOrderReq)
	if err != nil {
		t.Fatalf("Failed to get order: %v", err)
	}

	if getOrderResp.RetCode != 0 {
		t.Fatalf("Expected retCode 0, got %d: %s", getOrderResp.RetCode, getOrderResp.RetMsg)
	}

	// Verify we got order data in the list
	if len(getOrderResp.Result.List) == 0 {
		t.Error("Expected order data, got empty list")
	} else {
		order := getOrderResp.Result.List[0]
		t.Logf("Order retrieved successfully: %s", order.OrderId)
	}

	// Test getting single order
	getSingleOrderReq := &GetOrderRequest{
		Category: "inverse",
		Symbol:   "BTCUSD",
		OrderId:  createOrderResp.Result.OrderId,
	}

	getSingleOrderResp, err := client.GetSingleOrder(ctx, getSingleOrderReq)
	if err != nil {
		t.Fatalf("Failed to get single order: %v", err)
	}

	if getSingleOrderResp.RetCode != 0 {
		t.Fatalf("Expected retCode 0, got %d: %s", getSingleOrderResp.RetCode, getSingleOrderResp.RetMsg)
	}

	t.Logf("Single order retrieved successfully: %s", getSingleOrderResp.Result.OrderId)

	// Test canceling the order
	cancelOrderReq := &CancelOrderRequest{
		Category: "inverse",
		Symbol:   "BTCUSD",
		OrderId:  createOrderResp.Result.OrderId,
	}

	cancelOrderResp, err := client.CancelOrder(ctx, cancelOrderReq)
	if err != nil {
		t.Fatalf("Failed to cancel order: %v", err)
	}

	if cancelOrderResp.RetCode != 0 {
		t.Fatalf("Expected retCode 0, got %d: %s", cancelOrderResp.RetCode, cancelOrderResp.RetMsg)
	}

	t.Logf("Order canceled successfully: %s", cancelOrderResp.Result.OrderId)
}
