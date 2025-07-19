package binancefuture

import (
	"context"
	"testing"
	"time"
)

func TestNewClient(t *testing.T) {
	config := TestnetConfig()
	client := NewClient(config)

	if client == nil {
		t.Fatal("client should not be nil")
	}

	if client.config != config {
		t.Error("client config should match the provided config")
	}

	if client.requestService == nil {
		t.Error("request service should not be nil")
	}
}

func TestGetServerTime(t *testing.T) {
	config := TestnetConfig()
	client := NewClient(config)

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	serverTime, err := client.GetServerTime(ctx)
	if err != nil {
		t.Fatalf("failed to get server time: %v", err)
	}

	if serverTime == nil {
		t.Fatal("server time response should not be nil")
	}

	if serverTime.ServerTime <= 0 {
		t.Error("server time should be positive")
	}

	// Verify the time is reasonable (within 1 hour of current time)
	now := time.Now()
	serverTimeGo := serverTime.GetTime()
	diff := now.Sub(serverTimeGo)
	if diff < -time.Hour || diff > time.Hour {
		t.Errorf("server time seems unreasonable: server=%v, local=%v, diff=%v", serverTimeGo, now, diff)
	}
}

func TestPing(t *testing.T) {
	config := TestnetConfig()
	client := NewClient(config)

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	err := client.Ping(ctx)
	if err != nil {
		t.Fatalf("ping failed: %v", err)
	}
}

func TestGetExchangeInfo(t *testing.T) {
	config := TestnetConfig()
	client := NewClient(config)

	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	exchangeInfo, err := client.GetExchangeInfo(ctx)
	if err != nil {
		t.Fatalf("failed to get exchange info: %v", err)
	}

	if exchangeInfo == nil {
		t.Fatal("exchange info response should not be nil")
	}

	if len(exchangeInfo.Symbols) == 0 {
		t.Error("exchange info should contain symbols")
	}

	// Check for some common symbols
	foundBTC := false
	foundETH := false
	for _, symbol := range exchangeInfo.Symbols {
		if symbol.Symbol == "BTCUSDT" {
			foundBTC = true
		}
		if symbol.Symbol == "ETHUSDT" {
			foundETH = true
		}
	}

	if !foundBTC {
		t.Error("should find BTCUSDT symbol")
	}
	if !foundETH {
		t.Error("should find ETHUSDT symbol")
	}
}

func TestGetTickerPrice(t *testing.T) {
	config := TestnetConfig()
	client := NewClient(config)

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	// Test single symbol
	result, err := client.GetTickerPrice(ctx, "BTCUSDT")
	if err != nil {
		t.Fatalf("failed to get ticker price: %v", err)
	}

	if !result.IsSingle() {
		t.Error("should return single ticker for specific symbol")
	}

	if result.Single == nil {
		t.Fatal("single ticker should not be nil")
	}

	if result.Single.Symbol != "BTCUSDT" {
		t.Errorf("expected symbol BTCUSDT, got %s", result.Single.Symbol)
	}

	if result.Single.Price == "" {
		t.Error("price should not be empty")
	}
}

func TestGetTicker24hr(t *testing.T) {
	config := TestnetConfig()
	client := NewClient(config)

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	// Test single symbol
	result, err := client.GetTicker24hr(ctx, "BTCUSDT")
	if err != nil {
		t.Fatalf("failed to get 24hr ticker: %v", err)
	}

	if !result.IsSingle() {
		t.Error("should return single ticker for specific symbol")
	}

	if result.Single == nil {
		t.Fatal("single ticker should not be nil")
	}

	if result.Single.Symbol != "BTCUSDT" {
		t.Errorf("expected symbol BTCUSDT, got %s", result.Single.Symbol)
	}
}

func TestGetKlines(t *testing.T) {
	config := TestnetConfig()
	client := NewClient(config)

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	klines, err := client.GetKlines(ctx, "BTCUSDT", Interval1h, 10)
	if err != nil {
		t.Fatalf("failed to get klines: %v", err)
	}

	if klines == nil {
		t.Fatal("klines response should not be nil")
	}

	if len(*klines) == 0 {
		t.Error("klines should not be empty")
	}

	if len(*klines) > 10 {
		t.Errorf("expected max 10 klines, got %d", len(*klines))
	}

	// Check first kline structure
	firstKline := (*klines)[0]
	if firstKline.OpenTime <= 0 {
		t.Error("open time should be positive")
	}
	if firstKline.Open == "" {
		t.Error("open price should not be empty")
	}
}

func TestGetOrderBook(t *testing.T) {
	config := TestnetConfig()
	client := NewClient(config)

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	orderBook, err := client.GetOrderBook(ctx, "BTCUSDT", 10)
	if err != nil {
		t.Fatalf("failed to get order book: %v", err)
	}

	if orderBook == nil {
		t.Fatal("order book response should not be nil")
	}

	if orderBook.LastUpdateId <= 0 {
		t.Error("last update id should be positive")
	}

	if len(orderBook.Bids) == 0 {
		t.Error("bids should not be empty")
	}

	if len(orderBook.Asks) == 0 {
		t.Error("asks should not be empty")
	}

	// Check bid structure
	if len(orderBook.Bids) > 0 {
		bid := orderBook.Bids[0]
		if len(bid) != 2 {
			t.Errorf("bid should have 2 elements (price, quantity), got %d", len(bid))
		}
	}

	// Check ask structure
	if len(orderBook.Asks) > 0 {
		ask := orderBook.Asks[0]
		if len(ask) != 2 {
			t.Errorf("ask should have 2 elements (price, quantity), got %d", len(ask))
		}
	}
}

func TestGetTrades(t *testing.T) {
	config := TestnetConfig()
	client := NewClient(config)

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	trades, err := client.GetTrades(ctx, "BTCUSDT", 10)
	if err != nil {
		t.Fatalf("failed to get trades: %v", err)
	}

	if trades == nil {
		t.Fatal("trades response should not be nil")
	}

	if len(trades) == 0 {
		t.Error("trades should not be empty")
	}

	if len(trades) > 10 {
		t.Errorf("expected max 10 trades, got %d", len(trades))
	}

	// Check first trade structure
	firstTrade := trades[0]
	if firstTrade.Id <= 0 {
		t.Error("trade id should be positive")
	}
	if firstTrade.Price == "" {
		t.Error("trade price should not be empty")
	}
	if firstTrade.Qty == "" {
		t.Error("trade quantity should not be empty")
	}
}

func TestGetMarkPrice(t *testing.T) {
	config := TestnetConfig()
	client := NewClient(config)

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	markPrices, err := client.GetMarkPrice(ctx, "BTCUSDT")
	if err != nil {
		t.Fatalf("failed to get mark price: %v", err)
	}

	if markPrices == nil {
		t.Fatal("mark prices response should not be nil")
	}

	if len(markPrices) == 0 {
		t.Error("mark prices should not be empty")
	}

	// Check first mark price structure
	firstMarkPrice := markPrices[0]
	if firstMarkPrice.Symbol != "BTCUSDT" {
		t.Errorf("expected symbol BTCUSDT, got %s", firstMarkPrice.Symbol)
	}
	if firstMarkPrice.MarkPrice == "" {
		t.Error("mark price should not be empty")
	}
}

func TestGetFundingRate(t *testing.T) {
	config := TestnetConfig()
	client := NewClient(config)

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	fundingRates, err := client.GetFundingRate(ctx, "BTCUSDT", 5)
	if err != nil {
		t.Fatalf("failed to get funding rate: %v", err)
	}

	if fundingRates == nil {
		t.Fatal("funding rates response should not be nil")
	}

	if len(fundingRates) == 0 {
		t.Error("funding rates should not be empty")
	}

	if len(fundingRates) > 5 {
		t.Errorf("expected max 5 funding rates, got %d", len(fundingRates))
	}

	// Check first funding rate structure
	firstFundingRate := fundingRates[0]
	if firstFundingRate.Symbol != "BTCUSDT" {
		t.Errorf("expected symbol BTCUSDT, got %s", firstFundingRate.Symbol)
	}
	if firstFundingRate.FundingRate == "" {
		t.Error("funding rate should not be empty")
	}
}

func TestGetOpenInterest(t *testing.T) {
	config := TestnetConfig()
	client := NewClient(config)

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	openInterest, err := client.GetOpenInterest(ctx, "BTCUSDT")
	if err != nil {
		t.Fatalf("failed to get open interest: %v", err)
	}

	if openInterest == nil {
		t.Fatal("open interest response should not be nil")
	}

	if openInterest.Symbol != "BTCUSDT" {
		t.Errorf("expected symbol BTCUSDT, got %s", openInterest.Symbol)
	}

	if openInterest.OpenInterest == "" {
		t.Error("open interest should not be empty")
	}
}
