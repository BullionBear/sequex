package binanceperp

import (
	"context"
	"os"
	"testing"
)

func TestGetServerTime(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping integration test in short mode.")
	}

	cfg := &Config{
		BaseURL: MainnetBaseUrl,
	}
	client := NewClient(cfg)

	resp, err := client.GetServerTime(context.Background())

	// Test error != nil (should be nil for successful request)
	if err != nil {
		t.Fatalf("GetServerTime error: %v", err)
	}

	// Test Response.Code == 0 (success)
	if resp.Code != 0 {
		t.Fatalf("expected response code 0, got %d", resp.Code)
	}

	// Test Data is marshaled correctly
	if resp.Data == nil {
		t.Fatal("response data is nil, expected server time data")
	}

	if resp.Data.ServerTime == 0 {
		t.Error("serverTime is zero, expected non-zero timestamp")
	}

	// Verify the server time is a reasonable timestamp (after year 2020)
	minTimestamp := int64(1577836800000) // Jan 1, 2020 00:00:00 UTC in milliseconds
	if resp.Data.ServerTime < minTimestamp {
		t.Errorf("serverTime %d appears to be invalid (before 2020)", resp.Data.ServerTime)
	}
}

func TestGetServerTime_InvalidBaseURL(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping integration test in short mode.")
	}

	cfg := &Config{
		BaseURL: "https://invalid-url-that-does-not-exist.com",
	}
	client := NewClient(cfg)

	_, err := client.GetServerTime(context.Background())

	// Test error != nil (should have error for invalid URL)
	if err == nil {
		t.Fatal("expected error for invalid base URL, got nil")
	}
}

func TestGetDepth(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping integration test in short mode.")
	}

	cfg := &Config{
		BaseURL: MainnetBaseUrl,
	}
	client := NewClient(cfg)

	req := GetDepthRequest{
		Symbol: "BTCUSDT",
		Limit:  5,
	}
	resp, err := client.GetDepth(context.Background(), req)

	// Test error != nil (should be nil for successful request)
	if err != nil {
		t.Fatalf("GetDepth error: %v", err)
	}

	// Test Response.Code == 0 (success)
	if resp.Code != 0 {
		t.Fatalf("expected response code 0, got %d", resp.Code)
	}

	// Test Data is marshaled correctly
	if resp.Data == nil {
		t.Fatal("response data is nil, expected order book data")
	}

	if resp.Data.LastUpdateId == 0 {
		t.Error("lastUpdateId is zero, expected non-zero value")
	}

	if len(resp.Data.Bids) == 0 {
		t.Error("bids is empty, expected at least one bid")
	}

	if len(resp.Data.Asks) == 0 {
		t.Error("asks is empty, expected at least one ask")
	}

	// Verify bid/ask format [price, quantity]
	if len(resp.Data.Bids[0]) != 2 {
		t.Errorf("expected bid to have 2 elements [price, quantity], got %d", len(resp.Data.Bids[0]))
	}

	if len(resp.Data.Asks[0]) != 2 {
		t.Errorf("expected ask to have 2 elements [price, quantity], got %d", len(resp.Data.Asks[0]))
	}
}

func TestGetDepth_DefaultLimit(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping integration test in short mode.")
	}

	cfg := &Config{
		BaseURL: MainnetBaseUrl,
	}
	client := NewClient(cfg)

	req := GetDepthRequest{
		Symbol: "BTCUSDT",
		// No limit specified, should use default
	}
	resp, err := client.GetDepth(context.Background(), req)

	// Test error != nil (should be nil for successful request)
	if err != nil {
		t.Fatalf("GetDepth error: %v", err)
	}

	// Test Response.Code == 0 (success)
	if resp.Code != 0 {
		t.Fatalf("expected response code 0, got %d", resp.Code)
	}

	// Test Data is marshaled correctly
	if resp.Data == nil {
		t.Fatal("response data is nil, expected order book data")
	}
}

func TestGetDepth_InvalidSymbol(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping integration test in short mode.")
	}

	cfg := &Config{
		BaseURL: MainnetBaseUrl,
	}
	client := NewClient(cfg)

	req := GetDepthRequest{
		Symbol: "INVALIDSYMBOL",
		Limit:  5,
	}
	_, err := client.GetDepth(context.Background(), req)

	// Test error != nil (should have error for invalid symbol)
	if err == nil {
		t.Fatal("expected error for invalid symbol, got nil")
	}
}

func TestGetRecentTrades(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping integration test in short mode.")
	}

	cfg := &Config{
		BaseURL: MainnetBaseUrl,
	}
	client := NewClient(cfg)

	req := GetRecentTradesRequest{
		Symbol: "BTCUSDT",
		Limit:  5,
	}
	resp, err := client.GetRecentTrades(context.Background(), req)

	// Test error != nil (should be nil for successful request)
	if err != nil {
		t.Fatalf("GetRecentTrades error: %v", err)
	}

	// Test Response.Code == 0 (success)
	if resp.Code != 0 {
		t.Fatalf("expected response code 0, got %d", resp.Code)
	}

	// Test Data is marshaled correctly
	if resp.Data == nil {
		t.Fatal("response data is nil, expected trades data")
	}

	if len(*resp.Data) == 0 {
		t.Fatal("trades list is empty, expected at least one trade")
	}

	// Verify trade structure
	trade := (*resp.Data)[0]
	if trade.Id == 0 {
		t.Error("trade ID is zero, expected non-zero value")
	}

	if trade.Price == "" {
		t.Error("trade price is empty, expected non-empty value")
	}

	if trade.Qty == "" {
		t.Error("trade quantity is empty, expected non-empty value")
	}

	if trade.QuoteQty == "" {
		t.Error("trade quote quantity is empty, expected non-empty value")
	}

	if trade.Time == 0 {
		t.Error("trade time is zero, expected non-zero timestamp")
	}
}

func TestGetRecentTrades_DefaultLimit(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping integration test in short mode.")
	}

	cfg := &Config{
		BaseURL: MainnetBaseUrl,
	}
	client := NewClient(cfg)

	req := GetRecentTradesRequest{
		Symbol: "BTCUSDT",
		// No limit specified, should use default
	}
	resp, err := client.GetRecentTrades(context.Background(), req)

	// Test error != nil (should be nil for successful request)
	if err != nil {
		t.Fatalf("GetRecentTrades error: %v", err)
	}

	// Test Response.Code == 0 (success)
	if resp.Code != 0 {
		t.Fatalf("expected response code 0, got %d", resp.Code)
	}

	// Test Data is marshaled correctly
	if resp.Data == nil {
		t.Fatal("response data is nil, expected trades data")
	}

	// Should return more trades than the limited test (default is 500)
	if len(*resp.Data) == 0 {
		t.Fatal("trades list is empty, expected at least one trade")
	}
}

func TestGetRecentTrades_InvalidSymbol(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping integration test in short mode.")
	}

	cfg := &Config{
		BaseURL: MainnetBaseUrl,
	}
	client := NewClient(cfg)

	req := GetRecentTradesRequest{
		Symbol: "INVALIDSYMBOL",
		Limit:  5,
	}
	_, err := client.GetRecentTrades(context.Background(), req)

	// Test error != nil (should have error for invalid symbol)
	if err == nil {
		t.Fatal("expected error for invalid symbol, got nil")
	}
}

func TestGetAggTrades(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping integration test in short mode.")
	}

	cfg := &Config{
		BaseURL: MainnetBaseUrl,
	}
	client := NewClient(cfg)

	req := GetAggTradesRequest{
		Symbol: "BTCUSDT",
		Limit:  5,
	}
	resp, err := client.GetAggTrades(context.Background(), req)

	// Test error != nil (should be nil for successful request)
	if err != nil {
		t.Fatalf("GetAggTrades error: %v", err)
	}

	// Test Response.Code == 0 (success)
	if resp.Code != 0 {
		t.Fatalf("expected response code 0, got %d", resp.Code)
	}

	// Test Data is marshaled correctly
	if resp.Data == nil {
		t.Fatal("response data is nil, expected aggregate trades data")
	}

	if len(*resp.Data) == 0 {
		t.Fatal("aggregate trades list is empty, expected at least one trade")
	}

	// Verify aggregate trade structure
	trade := (*resp.Data)[0]
	if trade.AggTradeId == 0 {
		t.Error("aggregate trade ID is zero, expected non-zero value")
	}

	if trade.Price == "" {
		t.Error("trade price is empty, expected non-empty value")
	}

	if trade.Quantity == "" {
		t.Error("trade quantity is empty, expected non-empty value")
	}

	if trade.FirstTradeId == 0 {
		t.Error("first trade ID is zero, expected non-zero value")
	}

	if trade.LastTradeId == 0 {
		t.Error("last trade ID is zero, expected non-zero value")
	}

	if trade.Timestamp == 0 {
		t.Error("trade timestamp is zero, expected non-zero timestamp")
	}
}

func TestGetAggTrades_DefaultLimit(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping integration test in short mode.")
	}

	cfg := &Config{
		BaseURL: MainnetBaseUrl,
	}
	client := NewClient(cfg)

	req := GetAggTradesRequest{
		Symbol: "BTCUSDT",
		// No limit specified, should use default
	}
	resp, err := client.GetAggTrades(context.Background(), req)

	// Test error != nil (should be nil for successful request)
	if err != nil {
		t.Fatalf("GetAggTrades error: %v", err)
	}

	// Test Response.Code == 0 (success)
	if resp.Code != 0 {
		t.Fatalf("expected response code 0, got %d", resp.Code)
	}

	// Test Data is marshaled correctly
	if resp.Data == nil {
		t.Fatal("response data is nil, expected aggregate trades data")
	}

	// Should return more trades than the limited test (default is 500)
	if len(*resp.Data) == 0 {
		t.Fatal("aggregate trades list is empty, expected at least one trade")
	}
}

func TestGetAggTrades_WithFromId(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping integration test in short mode.")
	}

	cfg := &Config{
		BaseURL: MainnetBaseUrl,
	}
	client := NewClient(cfg)

	req := GetAggTradesRequest{
		Symbol: "BTCUSDT",
		FromId: 1000000000, // Use a reasonable fromId
		Limit:  3,
	}
	resp, err := client.GetAggTrades(context.Background(), req)

	// Test error != nil (should be nil for successful request)
	if err != nil {
		t.Fatalf("GetAggTrades error: %v", err)
	}

	// Test Response.Code == 0 (success)
	if resp.Code != 0 {
		t.Fatalf("expected response code 0, got %d", resp.Code)
	}

	// Test Data is marshaled correctly
	if resp.Data == nil {
		t.Fatal("response data is nil, expected aggregate trades data")
	}
}

func TestGetAggTrades_InvalidSymbol(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping integration test in short mode.")
	}

	cfg := &Config{
		BaseURL: MainnetBaseUrl,
	}
	client := NewClient(cfg)

	req := GetAggTradesRequest{
		Symbol: "INVALIDSYMBOL",
		Limit:  5,
	}
	_, err := client.GetAggTrades(context.Background(), req)

	// Test error != nil (should have error for invalid symbol)
	if err == nil {
		t.Fatal("expected error for invalid symbol, got nil")
	}
}

func TestGetKlines(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping integration test in short mode.")
	}

	cfg := &Config{
		BaseURL: MainnetBaseUrl,
	}
	client := NewClient(cfg)

	req := GetKlinesRequest{
		Symbol:   "BTCUSDT",
		Interval: "1m",
		Limit:    5,
	}
	resp, err := client.GetKlines(context.Background(), req)

	// Test error != nil (should be nil for successful request)
	if err != nil {
		t.Fatalf("GetKlines error: %v", err)
	}

	// Test Response.Code == 0 (success)
	if resp.Code != 0 {
		t.Fatalf("expected response code 0, got %d", resp.Code)
	}

	// Test Data is marshaled correctly
	if resp.Data == nil {
		t.Fatal("response data is nil, expected klines data")
	}

	if len(*resp.Data) == 0 {
		t.Fatal("klines list is empty, expected at least one kline")
	}

	// Verify kline structure
	kline := (*resp.Data)[0]
	if kline.OpenTime == 0 {
		t.Error("kline open time is zero, expected non-zero timestamp")
	}

	if kline.Open == "" {
		t.Error("kline open price is empty, expected non-empty value")
	}

	if kline.High == "" {
		t.Error("kline high price is empty, expected non-empty value")
	}

	if kline.Low == "" {
		t.Error("kline low price is empty, expected non-empty value")
	}

	if kline.Close == "" {
		t.Error("kline close price is empty, expected non-empty value")
	}

	if kline.Volume == "" {
		t.Error("kline volume is empty, expected non-empty value")
	}

	if kline.CloseTime == 0 {
		t.Error("kline close time is zero, expected non-zero timestamp")
	}

	if kline.QuoteAssetVolume == "" {
		t.Error("kline quote asset volume is empty, expected non-empty value")
	}

	if kline.NumberOfTrades == 0 {
		t.Error("kline number of trades is zero, expected non-zero value")
	}
}

func TestGetKlines_DefaultLimit(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping integration test in short mode.")
	}

	cfg := &Config{
		BaseURL: MainnetBaseUrl,
	}
	client := NewClient(cfg)

	req := GetKlinesRequest{
		Symbol:   "BTCUSDT",
		Interval: "1m",
		// No limit specified, should use default
	}
	resp, err := client.GetKlines(context.Background(), req)

	// Test error != nil (should be nil for successful request)
	if err != nil {
		t.Fatalf("GetKlines error: %v", err)
	}

	// Test Response.Code == 0 (success)
	if resp.Code != 0 {
		t.Fatalf("expected response code 0, got %d", resp.Code)
	}

	// Test Data is marshaled correctly
	if resp.Data == nil {
		t.Fatal("response data is nil, expected klines data")
	}

	// Should return more klines than the limited test (default is 500)
	if len(*resp.Data) == 0 {
		t.Fatal("klines list is empty, expected at least one kline")
	}
}

func TestGetKlines_DifferentIntervals(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping integration test in short mode.")
	}

	cfg := &Config{
		BaseURL: MainnetBaseUrl,
	}
	client := NewClient(cfg)

	intervals := []string{"1m", "5m", "1h", "1d"}
	for _, interval := range intervals {
		t.Run(interval, func(t *testing.T) {
			req := GetKlinesRequest{
				Symbol:   "BTCUSDT",
				Interval: interval,
				Limit:    3,
			}
			resp, err := client.GetKlines(context.Background(), req)

			// Test error != nil (should be nil for successful request)
			if err != nil {
				t.Fatalf("GetKlines error for interval %s: %v", interval, err)
			}

			// Test Response.Code == 0 (success)
			if resp.Code != 0 {
				t.Fatalf("expected response code 0 for interval %s, got %d", interval, resp.Code)
			}

			// Test Data is marshaled correctly
			if resp.Data == nil {
				t.Fatalf("response data is nil for interval %s, expected klines data", interval)
			}

			if len(*resp.Data) == 0 {
				t.Fatalf("klines list is empty for interval %s, expected at least one kline", interval)
			}
		})
	}
}

func TestGetKlines_InvalidSymbol(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping integration test in short mode.")
	}

	cfg := &Config{
		BaseURL: MainnetBaseUrl,
	}
	client := NewClient(cfg)

	req := GetKlinesRequest{
		Symbol:   "INVALIDSYMBOL",
		Interval: "1m",
		Limit:    5,
	}
	_, err := client.GetKlines(context.Background(), req)

	// Test error != nil (should have error for invalid symbol)
	if err == nil {
		t.Fatal("expected error for invalid symbol, got nil")
	}
}

func TestGetKlines_InvalidInterval(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping integration test in short mode.")
	}

	cfg := &Config{
		BaseURL: MainnetBaseUrl,
	}
	client := NewClient(cfg)

	req := GetKlinesRequest{
		Symbol:   "BTCUSDT",
		Interval: "INVALIDINTERVAL",
		Limit:    5,
	}
	_, err := client.GetKlines(context.Background(), req)

	// Test error != nil (should have error for invalid interval)
	if err == nil {
		t.Fatal("expected error for invalid interval, got nil")
	}
}

func TestGetMarkPrice(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping integration test in short mode.")
	}

	cfg := &Config{
		BaseURL: MainnetBaseUrl,
	}
	client := NewClient(cfg)

	req := GetMarkPriceRequest{
		Symbol: "BTCUSDT",
	}
	resp, err := client.GetMarkPrice(context.Background(), req)

	// Test error != nil (should be nil for successful request)
	if err != nil {
		t.Fatalf("GetMarkPrice error: %v", err)
	}

	// Test Response.Code == 0 (success)
	if resp.Code != 0 {
		t.Fatalf("expected response code 0, got %d", resp.Code)
	}

	// Test Data is marshaled correctly
	if resp.Data == nil {
		t.Fatal("response data is nil, expected mark price data")
	}

	if len(*resp.Data) == 0 {
		t.Fatal("mark price list is empty, expected at least one entry")
	}

	// Verify mark price structure
	markPrice := (*resp.Data)[0]
	if markPrice.Symbol == "" {
		t.Error("mark price symbol is empty, expected non-empty value")
	}

	if markPrice.MarkPrice == "" {
		t.Error("mark price is empty, expected non-empty value")
	}

	if markPrice.IndexPrice == "" {
		t.Error("index price is empty, expected non-empty value")
	}

	if markPrice.EstimatedSettlePrice == "" {
		t.Error("estimated settle price is empty, expected non-empty value")
	}

	if markPrice.LastFundingRate == "" {
		t.Error("last funding rate is empty, expected non-empty value")
	}

	if markPrice.InterestRate == "" {
		t.Error("interest rate is empty, expected non-empty value")
	}

	if markPrice.NextFundingTime == 0 {
		t.Error("next funding time is zero, expected non-zero timestamp")
	}

	if markPrice.Time == 0 {
		t.Error("time is zero, expected non-zero timestamp")
	}
}

func TestGetMarkPrice_AllSymbols(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping integration test in short mode.")
	}

	cfg := &Config{
		BaseURL: MainnetBaseUrl,
	}
	client := NewClient(cfg)

	req := GetMarkPriceRequest{
		// No symbol specified, should return all symbols
	}
	resp, err := client.GetMarkPrice(context.Background(), req)

	// Test error != nil (should be nil for successful request)
	if err != nil {
		t.Fatalf("GetMarkPrice error: %v", err)
	}

	// Test Response.Code == 0 (success)
	if resp.Code != 0 {
		t.Fatalf("expected response code 0, got %d", resp.Code)
	}

	// Test Data is marshaled correctly
	if resp.Data == nil {
		t.Fatal("response data is nil, expected mark price data")
	}

	// Should return multiple symbols when no symbol is specified
	if len(*resp.Data) == 0 {
		t.Fatal("mark price list is empty, expected multiple entries")
	}

	// Verify we got multiple symbols
	if len(*resp.Data) < 2 {
		t.Logf("Warning: only got %d mark price entries, expected multiple symbols", len(*resp.Data))
	}
}

func TestGetMarkPrice_InvalidSymbol(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping integration test in short mode.")
	}

	cfg := &Config{
		BaseURL: MainnetBaseUrl,
	}
	client := NewClient(cfg)

	req := GetMarkPriceRequest{
		Symbol: "INVALIDSYMBOL",
	}
	_, err := client.GetMarkPrice(context.Background(), req)

	// Test error != nil (should have error for invalid symbol)
	if err == nil {
		t.Fatal("expected error for invalid symbol, got nil")
	}
}

func TestGetPriceTicker(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping integration test in short mode.")
	}

	cfg := &Config{
		BaseURL: MainnetBaseUrl,
	}
	client := NewClient(cfg)

	req := GetPriceTickerRequest{
		Symbol: "BTCUSDT",
	}
	resp, err := client.GetPriceTicker(context.Background(), req)

	// Test error != nil (should be nil for successful request)
	if err != nil {
		t.Fatalf("GetPriceTicker error: %v", err)
	}

	// Test Response.Code == 0 (success)
	if resp.Code != 0 {
		t.Fatalf("expected response code 0, got %d", resp.Code)
	}

	// Test Data is marshaled correctly
	if resp.Data == nil {
		t.Fatal("response data is nil, expected price ticker data")
	}

	if len(*resp.Data) == 0 {
		t.Fatal("price ticker list is empty, expected at least one entry")
	}

	// Verify price ticker structure
	ticker := (*resp.Data)[0]
	if ticker.Symbol == "" {
		t.Error("ticker symbol is empty, expected non-empty value")
	}

	if ticker.Price == "" {
		t.Error("ticker price is empty, expected non-empty value")
	}

	if ticker.Time == 0 {
		t.Error("ticker time is zero, expected non-zero timestamp")
	}

	// Verify the symbol matches what we requested
	if ticker.Symbol != "BTCUSDT" {
		t.Errorf("expected symbol BTCUSDT, got %s", ticker.Symbol)
	}
}

func TestGetPriceTicker_AllSymbols(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping integration test in short mode.")
	}

	cfg := &Config{
		BaseURL: MainnetBaseUrl,
	}
	client := NewClient(cfg)

	req := GetPriceTickerRequest{
		// No symbol specified, should return all symbols
	}
	resp, err := client.GetPriceTicker(context.Background(), req)

	// Test error != nil (should be nil for successful request)
	if err != nil {
		t.Fatalf("GetPriceTicker error: %v", err)
	}

	// Test Response.Code == 0 (success)
	if resp.Code != 0 {
		t.Fatalf("expected response code 0, got %d", resp.Code)
	}

	// Test Data is marshaled correctly
	if resp.Data == nil {
		t.Fatal("response data is nil, expected price ticker data")
	}

	// Should return multiple symbols when no symbol is specified
	if len(*resp.Data) == 0 {
		t.Fatal("price ticker list is empty, expected multiple entries")
	}

	// Verify we got multiple symbols
	if len(*resp.Data) < 2 {
		t.Logf("Warning: only got %d price ticker entries, expected multiple symbols", len(*resp.Data))
	}

	// Verify each ticker has valid data
	for i, ticker := range *resp.Data {
		if i >= 3 { // Only check first 3 for performance
			break
		}
		if ticker.Symbol == "" {
			t.Errorf("ticker %d symbol is empty, expected non-empty value", i)
		}
		if ticker.Price == "" {
			t.Errorf("ticker %d price is empty, expected non-empty value", i)
		}
		if ticker.Time == 0 {
			t.Errorf("ticker %d time is zero, expected non-zero timestamp", i)
		}
	}
}

func TestGetPriceTicker_InvalidSymbol(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping integration test in short mode.")
	}

	cfg := &Config{
		BaseURL: MainnetBaseUrl,
	}
	client := NewClient(cfg)

	req := GetPriceTickerRequest{
		Symbol: "INVALIDSYMBOL",
	}
	_, err := client.GetPriceTicker(context.Background(), req)

	// Test error != nil (should have error for invalid symbol)
	if err == nil {
		t.Fatal("expected error for invalid symbol, got nil")
	}
}

func TestGetBookTicker(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping integration test in short mode.")
	}

	cfg := &Config{
		BaseURL: MainnetBaseUrl,
	}
	client := NewClient(cfg)

	req := GetBookTickerRequest{
		Symbol: "BTCUSDT",
	}
	resp, err := client.GetBookTicker(context.Background(), req)

	// Test error != nil (should be nil for successful request)
	if err != nil {
		t.Fatalf("GetBookTicker error: %v", err)
	}

	// Test Response.Code == 0 (success)
	if resp.Code != 0 {
		t.Fatalf("expected response code 0, got %d", resp.Code)
	}

	// Test Data is marshaled correctly
	if resp.Data == nil {
		t.Fatal("response data is nil, expected book ticker data")
	}

	if len(*resp.Data) == 0 {
		t.Fatal("book ticker list is empty, expected at least one entry")
	}

	// Verify book ticker structure
	ticker := (*resp.Data)[0]
	if ticker.Symbol == "" {
		t.Error("ticker symbol is empty, expected non-empty value")
	}

	if ticker.BidPrice == "" {
		t.Error("bid price is empty, expected non-empty value")
	}

	if ticker.BidQty == "" {
		t.Error("bid quantity is empty, expected non-empty value")
	}

	if ticker.AskPrice == "" {
		t.Error("ask price is empty, expected non-empty value")
	}

	if ticker.AskQty == "" {
		t.Error("ask quantity is empty, expected non-empty value")
	}

	if ticker.Time == 0 {
		t.Error("ticker time is zero, expected non-zero timestamp")
	}

	// Verify the symbol matches what we requested
	if ticker.Symbol != "BTCUSDT" {
		t.Errorf("expected symbol BTCUSDT, got %s", ticker.Symbol)
	}
}

func TestGetBookTicker_AllSymbols(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping integration test in short mode.")
	}

	cfg := &Config{
		BaseURL: MainnetBaseUrl,
	}
	client := NewClient(cfg)

	req := GetBookTickerRequest{
		// No symbol specified, should return all symbols
	}
	resp, err := client.GetBookTicker(context.Background(), req)

	// Test error != nil (should be nil for successful request)
	if err != nil {
		t.Fatalf("GetBookTicker error: %v", err)
	}

	// Test Response.Code == 0 (success)
	if resp.Code != 0 {
		t.Fatalf("expected response code 0, got %d", resp.Code)
	}

	// Test Data is marshaled correctly
	if resp.Data == nil {
		t.Fatal("response data is nil, expected book ticker data")
	}

	// Should return multiple symbols when no symbol is specified
	if len(*resp.Data) == 0 {
		t.Fatal("book ticker list is empty, expected multiple entries")
	}

	// Verify we got multiple symbols
	if len(*resp.Data) < 2 {
		t.Logf("Warning: only got %d book ticker entries, expected multiple symbols", len(*resp.Data))
	}

	// Verify each ticker has valid data
	for i, ticker := range *resp.Data {
		if i >= 3 { // Only check first 3 for performance
			break
		}
		if ticker.Symbol == "" {
			t.Errorf("ticker %d symbol is empty, expected non-empty value", i)
		}
		if ticker.BidPrice == "" {
			t.Errorf("ticker %d bid price is empty, expected non-empty value", i)
		}
		if ticker.BidQty == "" {
			t.Errorf("ticker %d bid quantity is empty, expected non-empty value", i)
		}
		if ticker.AskPrice == "" {
			t.Errorf("ticker %d ask price is empty, expected non-empty value", i)
		}
		if ticker.AskQty == "" {
			t.Errorf("ticker %d ask quantity is empty, expected non-empty value", i)
		}
		if ticker.Time == 0 {
			t.Errorf("ticker %d time is zero, expected non-zero timestamp", i)
		}
	}
}

func TestGetBookTicker_InvalidSymbol(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping integration test in short mode.")
	}

	cfg := &Config{
		BaseURL: MainnetBaseUrl,
	}
	client := NewClient(cfg)

	req := GetBookTickerRequest{
		Symbol: "INVALIDSYMBOL",
	}
	_, err := client.GetBookTicker(context.Background(), req)

	// Test error != nil (should have error for invalid symbol)
	if err == nil {
		t.Fatal("expected error for invalid symbol, got nil")
	}
}

func TestGetAccountBalance(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping integration test in short mode.")
	}

	apiKey := os.Getenv("BINANCEPERP_API_KEY")
	apiSecret := os.Getenv("BINANCEPERP_API_SECRET")
	if apiKey == "" || apiSecret == "" {
		t.Skip("BINANCEPERP_API_KEY or BINANCEPERP_API_SECRET not set; skipping signed request test.")
	}

	cfg := &Config{
		APIKey:    apiKey,
		APISecret: apiSecret,
		BaseURL:   MainnetBaseUrl,
	}
	client := NewClient(cfg)

	req := GetAccountBalanceRequest{
		RecvWindow: 5000,
	}
	resp, err := client.GetAccountBalance(context.Background(), req)

	// Test error != nil (should be nil for successful request)
	if err != nil {
		t.Fatalf("GetAccountBalance error: %v", err)
	}

	// Test Response.Code == 0 (success)
	if resp.Code != 0 {
		t.Fatalf("expected response code 0, got %d", resp.Code)
	}

	// Test Data is marshaled correctly
	if resp.Data == nil {
		t.Fatal("response data is nil, expected account balance data")
	}

	if len(*resp.Data) == 0 {
		t.Fatal("account balance list is empty, expected at least one balance entry")
	}

	// Verify account balance structure
	balance := (*resp.Data)[0]
	if balance.AccountAlias == "" {
		t.Error("account alias is empty, expected non-empty value")
	}

	if balance.Asset == "" {
		t.Error("asset is empty, expected non-empty value")
	}

	if balance.Balance == "" {
		t.Error("balance is empty, expected non-empty value")
	}

	if balance.CrossWalletBalance == "" {
		t.Error("cross wallet balance is empty, expected non-empty value")
	}

	if balance.AvailableBalance == "" {
		t.Error("available balance is empty, expected non-empty value")
	}

	if balance.MaxWithdrawAmount == "" {
		t.Error("max withdraw amount is empty, expected non-empty value")
	}

	if balance.UpdateTime == 0 {
		t.Error("update time is zero, expected non-zero timestamp")
	}
}

func TestGetAccountBalance_DefaultRecvWindow(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping integration test in short mode.")
	}

	apiKey := os.Getenv("BINANCEPERP_API_KEY")
	apiSecret := os.Getenv("BINANCEPERP_API_SECRET")
	if apiKey == "" || apiSecret == "" {
		t.Skip("BINANCEPERP_API_KEY or BINANCEPERP_API_SECRET not set; skipping signed request test.")
	}

	cfg := &Config{
		APIKey:    apiKey,
		APISecret: apiSecret,
		BaseURL:   MainnetBaseUrl,
	}
	client := NewClient(cfg)

	req := GetAccountBalanceRequest{
		// No recvWindow specified, should use default
	}
	resp, err := client.GetAccountBalance(context.Background(), req)

	// Test error != nil (should be nil for successful request)
	if err != nil {
		t.Fatalf("GetAccountBalance error: %v", err)
	}

	// Test Response.Code == 0 (success)
	if resp.Code != 0 {
		t.Fatalf("expected response code 0, got %d", resp.Code)
	}

	// Test Data is marshaled correctly
	if resp.Data == nil {
		t.Fatal("response data is nil, expected account balance data")
	}
}

func TestGetAccountBalance_InvalidCredentials(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping integration test in short mode.")
	}

	cfg := &Config{
		APIKey:    "invalid_api_key",
		APISecret: "invalid_api_secret",
		BaseURL:   MainnetBaseUrl,
	}
	client := NewClient(cfg)

	req := GetAccountBalanceRequest{}
	_, err := client.GetAccountBalance(context.Background(), req)

	// Test error != nil (should have error for invalid credentials)
	if err == nil {
		t.Fatal("expected error for invalid credentials, got nil")
	}
}
