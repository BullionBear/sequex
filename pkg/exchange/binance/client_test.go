package binance

import (
	"context"
	"os"
	"testing"
)

func TestGetDepth(t *testing.T) {
	cfg := &Config{
		BaseURL: MainnetBaseUrl,
	}
	client := NewClient(cfg)
	ctx := context.Background()
	resp, err := client.GetDepth(ctx, "BTCUSDT", 5)
	if err != nil {
		t.Fatalf("GetDepth error: %v", err)
	}
	if resp.Code != 0 {
		t.Fatalf("unexpected response code: %d, msg: %s", resp.Code, resp.Message)
	}
	if resp.Data == nil {
		t.Fatal("resp.Data is nil")
	}
	if len(resp.Data.Bids) == 0 || len(resp.Data.Asks) == 0 {
		t.Fatal("bids or asks are empty")
	}
}

func TestGetRecentTrades(t *testing.T) {
	cfg := &Config{
		BaseURL: MainnetBaseUrl,
	}
	client := NewClient(cfg)
	ctx := context.Background()
	resp, err := client.GetRecentTrades(ctx, "BTCUSDT", 5)
	if err != nil {
		t.Fatalf("GetRecentTrades error: %v", err)
	}
	if resp.Code != 0 {
		t.Fatalf("unexpected response code: %d, msg: %s", resp.Code, resp.Message)
	}
	if resp.Data == nil {
		t.Fatal("resp.Data is nil")
	}
	if len(*resp.Data) == 0 {
		t.Fatal("no trades returned")
	}
}

func TestGetAggTrades(t *testing.T) {
	cfg := &Config{
		BaseURL: MainnetBaseUrl,
	}
	client := NewClient(cfg)
	ctx := context.Background()
	resp, err := client.GetAggTrades(ctx, "BTCUSDT", 0, 0, 0, 5)
	if err != nil {
		t.Fatalf("GetAggTrades error: %v", err)
	}
	if resp.Code != 0 {
		t.Fatalf("unexpected response code: %d, msg: %s", resp.Code, resp.Message)
	}
	if resp.Data == nil {
		t.Fatal("resp.Data is nil")
	}
	if len(*resp.Data) == 0 {
		t.Fatal("no aggregate trades returned")
	}
}

func TestGetKlines(t *testing.T) {
	cfg := &Config{
		BaseURL: MainnetBaseUrl,
	}
	client := NewClient(cfg)
	ctx := context.Background()
	resp, err := client.GetKlines(ctx, "BTCUSDT", "1m", 0, 0, "", 5)
	if err != nil {
		t.Fatalf("GetKlines error: %v", err)
	}
	if resp.Code != 0 {
		t.Fatalf("unexpected response code: %d, msg: %s", resp.Code, resp.Message)
	}
	if resp.Data == nil {
		t.Fatal("resp.Data is nil")
	}
	if len(*resp.Data) == 0 {
		t.Fatal("no klines returned")
	}
}

func TestGetPriceTicker(t *testing.T) {
	cfg := &Config{
		BaseURL: MainnetBaseUrl,
	}
	client := NewClient(cfg)
	ctx := context.Background()

	// Single symbol
	resp, err := client.GetPriceTicker(ctx, "BTCUSDT")
	if err != nil {
		t.Fatalf("GetPriceTicker error: %v", err)
	}
	if resp.Code != 0 {
		t.Fatalf("unexpected response code: %d, msg: %s", resp.Code, resp.Message)
	}
	if resp.Data == nil || len(*resp.Data) == 0 {
		t.Fatal("no price ticker returned (single)")
	}

	// Multiple symbols
	resp2, err := client.GetPriceTicker(ctx, "BTCUSDT", "ETHBTC")
	if err != nil {
		t.Fatalf("GetPriceTicker (multi) error: %v", err)
	}
	if resp2.Code != 0 {
		t.Fatalf("unexpected response code (multi): %d, msg: %s", resp2.Code, resp2.Message)
	}
	if resp2.Data == nil || len(*resp2.Data) == 0 {
		t.Fatal("no price ticker returned (multi)")
	}
}

func TestGetAccountInfo(t *testing.T) {
	apiKey := os.Getenv("BINANCE_API_KEY")
	apiSecret := os.Getenv("BINANCE_API_SECRET")
	if apiKey == "" || apiSecret == "" {
		t.Skip("BINANCE_API_KEY or BINANCE_API_SECRET not set; skipping signed request test.")
	}
	cfg := &Config{
		APIKey:    apiKey,
		APISecret: apiSecret,
		BaseURL:   MainnetBaseUrl,
	}
	client := NewClient(cfg)
	ctx := context.Background()

	resp, err := client.GetAccountInfo(ctx, GetAccountInfoRequest{})
	if err != nil {
		t.Fatalf("GetAccountInfo error: %v", err)
	}
	if resp.Code != 0 {
		t.Fatalf("unexpected response code: %d, msg: %s", resp.Code, resp.Message)
	}
	if resp.Data == nil {
		t.Fatal("resp.Data is nil")
	}
}

func TestListOpenOrders(t *testing.T) {
	apiKey := os.Getenv("BINANCE_API_KEY")
	apiSecret := os.Getenv("BINANCE_API_SECRET")
	if apiKey == "" || apiSecret == "" {
		t.Skip("BINANCE_API_KEY or BINANCE_API_SECRET not set; skipping signed request test.")
	}
	cfg := &Config{
		APIKey:    apiKey,
		APISecret: apiSecret,
		BaseURL:   MainnetBaseUrl,
	}
	client := NewClient(cfg)
	ctx := context.Background()

	// All symbols
	resp, err := client.ListOpenOrders(ctx, ListOpenOrdersRequest{})
	if err != nil {
		t.Fatalf("ListOpenOrders error (all): %v", err)
	}
	if resp.Code != 0 {
		t.Fatalf("unexpected response code (all): %d, msg: %s", resp.Code, resp.Message)
	}
	if resp.Data == nil {
		t.Fatal("resp.Data is nil (all)")
	}

	// Single symbol (replace with a real symbol for your account)
	resp2, err := client.ListOpenOrders(ctx, ListOpenOrdersRequest{Symbol: "BTCUSDT"})
	if err != nil {
		t.Fatalf("ListOpenOrders error (symbol): %v", err)
	}
	if resp2.Code != 0 {
		t.Fatalf("unexpected response code (symbol): %d, msg: %s", resp2.Code, resp2.Message)
	}
	if resp2.Data == nil {
		t.Fatal("resp.Data is nil (symbol)")
	}
}

func TestGetMyTrades(t *testing.T) {
	apiKey := os.Getenv("BINANCE_API_KEY")
	apiSecret := os.Getenv("BINANCE_API_SECRET")
	if apiKey == "" || apiSecret == "" {
		t.Skip("BINANCE_API_KEY or BINANCE_API_SECRET not set; skipping signed request test.")
	}
	cfg := &Config{
		APIKey:    apiKey,
		APISecret: apiSecret,
		BaseURL:   MainnetBaseUrl,
	}
	client := NewClient(cfg)
	ctx := context.Background()

	// Replace with a real symbol for your account
	resp, err := client.GetMyTrades(ctx, GetAccountTradesRequest{Symbol: "BTCUSDT", Limit: 5})
	if err != nil {
		t.Fatalf("GetMyTrades error: %v", err)
	}
	if resp.Code != 0 {
		t.Fatalf("unexpected response code: %d, msg: %s", resp.Code, resp.Message)
	}
	if resp.Data == nil {
		t.Fatal("resp.Data is nil")
	}
}

func TestStartUserDataStream(t *testing.T) {
	apiKey := os.Getenv("BINANCE_API_KEY")
	apiSecret := os.Getenv("BINANCE_API_SECRET")
	if apiKey == "" || apiSecret == "" {
		t.Skip("BINANCE_API_KEY or BINANCE_API_SECRET not set; skipping signed request test.")
	}
	cfg := &Config{
		APIKey:    apiKey,
		APISecret: apiSecret,
		BaseURL:   MainnetBaseUrl,
	}
	client := NewClient(cfg)
	ctx := context.Background()

	resp, err := client.StartUserDataStream(ctx)
	if err != nil {
		t.Fatalf("StartUserDataStream error: %v", err)
	}
	if resp.Code != 0 {
		t.Fatalf("unexpected response code: %d, msg: %s", resp.Code, resp.Message)
	}
	if resp.Data == nil {
		t.Fatal("resp.Data is nil")
	}
	if resp.Data.ListenKey == "" {
		t.Fatal("listenKey is empty")
	}
}

func TestKeepaliveUserDataStream(t *testing.T) {
	apiKey := os.Getenv("BINANCE_API_KEY")
	apiSecret := os.Getenv("BINANCE_API_SECRET")
	if apiKey == "" || apiSecret == "" {
		t.Skip("BINANCE_API_KEY or BINANCE_API_SECRET not set; skipping signed request test.")
	}
	cfg := &Config{
		APIKey:    apiKey,
		APISecret: apiSecret,
		BaseURL:   MainnetBaseUrl,
	}
	client := NewClient(cfg)
	ctx := context.Background()

	// Start a stream first to get a listenKey
	startResp, err := client.StartUserDataStream(ctx)
	if err != nil {
		t.Fatalf("StartUserDataStream error: %v", err)
	}
	if startResp.Data == nil || startResp.Data.ListenKey == "" {
		t.Fatal("failed to get listenKey")
	}

	// Test keepalive
	resp, err := client.KeepaliveUserDataStream(ctx, startResp.Data.ListenKey)
	if err != nil {
		t.Fatalf("KeepaliveUserDataStream error: %v", err)
	}
	if resp.Code != 0 {
		t.Fatalf("unexpected response code: %d, msg: %s", resp.Code, resp.Message)
	}
	if resp.Data == nil {
		t.Fatal("resp.Data is nil")
	}
}

func TestCloseUserDataStream(t *testing.T) {
	apiKey := os.Getenv("BINANCE_API_KEY")
	apiSecret := os.Getenv("BINANCE_API_SECRET")
	if apiKey == "" || apiSecret == "" {
		t.Skip("BINANCE_API_KEY or BINANCE_API_SECRET not set; skipping signed request test.")
	}
	cfg := &Config{
		APIKey:    apiKey,
		APISecret: apiSecret,
		BaseURL:   MainnetBaseUrl,
	}
	client := NewClient(cfg)
	ctx := context.Background()

	// Start a stream first to get a listenKey
	startResp, err := client.StartUserDataStream(ctx)
	if err != nil {
		t.Fatalf("StartUserDataStream error: %v", err)
	}
	if startResp.Data == nil || startResp.Data.ListenKey == "" {
		t.Fatal("failed to get listenKey")
	}

	// Test close
	resp, err := client.CloseUserDataStream(ctx, startResp.Data.ListenKey)
	if err != nil {
		t.Fatalf("CloseUserDataStream error: %v", err)
	}
	if resp.Code != 0 {
		t.Fatalf("unexpected response code: %d, msg: %s", resp.Code, resp.Message)
	}
	if resp.Data == nil {
		t.Fatal("resp.Data is nil")
	}
}
