package binance

import (
	"context"
	"fmt"
	"io"
	"net/http"
	"strings"
	"testing"
	"time"
)

// mockHTTPClient is a mock implementation of HTTPClient for testing
type mockHTTPClient struct {
	responses map[string]*http.Response
	errors    map[string]error
}

func newMockHTTPClient() *mockHTTPClient {
	return &mockHTTPClient{
		responses: make(map[string]*http.Response),
		errors:    make(map[string]error),
	}
}

func (m *mockHTTPClient) Do(req *http.Request) (*http.Response, error) {
	url := req.URL.String()

	// Check for errors first
	if err, exists := m.errors[url]; exists {
		return nil, err
	}

	// Return response if exists
	if resp, exists := m.responses[url]; exists {
		return resp, nil
	}

	// Default 404 response
	return &http.Response{
		StatusCode: http.StatusNotFound,
		Body:       io.NopCloser(strings.NewReader(`{"code":-1000,"msg":"Not found"}`)),
	}, nil
}

func (m *mockHTTPClient) setResponse(url string, statusCode int, body string) {
	m.responses[url] = &http.Response{
		StatusCode: statusCode,
		Body:       io.NopCloser(strings.NewReader(body)),
	}
}

func (m *mockHTTPClient) setError(url string, err error) {
	m.errors[url] = err
}

// createTestClient creates a client with mock HTTP client
func createTestClient() (*Client, *mockHTTPClient) {
	config := DefaultConfig()
	config.BaseURL = "https://api.binance.com"

	client := NewClient(config)
	mockClient := newMockHTTPClient()
	client.requestService.httpClient = mockClient

	return client, mockClient
}

func TestClient_GetServerTime_Success(t *testing.T) {
	client, mockClient := createTestClient()

	// Mock successful response
	expectedTimestamp := time.Now().UnixNano() / int64(time.Millisecond)
	responseBody := fmt.Sprintf(`{"serverTime":%d}`, expectedTimestamp)
	mockClient.setResponse("https://api.binance.com/api/v3/time", http.StatusOK, responseBody)

	// Test GetServerTime
	ctx := context.Background()
	result, err := client.GetServerTime(ctx)

	// Assertions
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	if result == nil {
		t.Fatal("expected result, got nil")
	}

	if result.ServerTime != expectedTimestamp {
		t.Errorf("expected server time %d, got %d", expectedTimestamp, result.ServerTime)
	}

	// Test the GetTime() method
	expectedTime := time.Unix(0, expectedTimestamp*int64(time.Millisecond))
	actualTime := result.GetTime()

	if !actualTime.Equal(expectedTime) {
		t.Errorf("expected time %v, got %v", expectedTime, actualTime)
	}
}

func TestClient_GetServerTime_APIError(t *testing.T) {
	client, mockClient := createTestClient()

	// Mock API error response
	responseBody := `{"code":-1003,"msg":"Too many requests"}`
	mockClient.setResponse("https://api.binance.com/api/v3/time", http.StatusTooManyRequests, responseBody)

	// Test GetServerTime
	ctx := context.Background()
	result, err := client.GetServerTime(ctx)

	// Assertions
	if err == nil {
		t.Fatal("expected error, got nil")
	}

	if result != nil {
		t.Errorf("expected nil result, got %v", result)
	}

	// Check if it's an API error
	apiErr, ok := err.(*APIError)
	if !ok {
		t.Errorf("expected APIError, got %T", err)
	} else {
		if apiErr.Code != -1003 {
			t.Errorf("expected error code -1003, got %d", apiErr.Code)
		}
		if apiErr.Message != "Too many requests" {
			t.Errorf("expected error message 'Too many requests', got %s", apiErr.Message)
		}
	}
}

func TestClient_GetServerTime_NetworkError(t *testing.T) {
	client, mockClient := createTestClient()

	// Mock network error
	networkErr := fmt.Errorf("network connection failed")
	mockClient.setError("https://api.binance.com/api/v3/time", networkErr)

	// Test GetServerTime
	ctx := context.Background()
	result, err := client.GetServerTime(ctx)

	// Assertions
	if err == nil {
		t.Fatal("expected error, got nil")
	}

	if result != nil {
		t.Errorf("expected nil result, got %v", result)
	}

	// Check error message contains network error
	if !strings.Contains(err.Error(), "network connection failed") {
		t.Errorf("expected error to contain 'network connection failed', got %v", err)
	}
}

func TestClient_GetServerTime_InvalidJSON(t *testing.T) {
	client, mockClient := createTestClient()

	// Mock invalid JSON response
	responseBody := `{"serverTime":invalid_json}`
	mockClient.setResponse("https://api.binance.com/api/v3/time", http.StatusOK, responseBody)

	// Test GetServerTime
	ctx := context.Background()
	result, err := client.GetServerTime(ctx)

	// Assertions
	if err == nil {
		t.Fatal("expected error, got nil")
	}

	if result != nil {
		t.Errorf("expected nil result, got %v", result)
	}

	// Check error message contains parsing error
	if !strings.Contains(err.Error(), "failed to parse server time response") {
		t.Errorf("expected parsing error, got %v", err)
	}
}

func TestClient_Ping_Success(t *testing.T) {
	client, mockClient := createTestClient()

	// Mock successful ping response
	mockClient.setResponse("https://api.binance.com/api/v3/ping", http.StatusOK, `{}`)

	// Test Ping
	ctx := context.Background()
	err := client.Ping(ctx)

	// Assertions
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
}

func TestClient_Ping_Error(t *testing.T) {
	client, mockClient := createTestClient()

	// Mock error response
	responseBody := `{"code":-1000,"msg":"Server error"}`
	mockClient.setResponse("https://api.binance.com/api/v3/ping", http.StatusInternalServerError, responseBody)

	// Test Ping
	ctx := context.Background()
	err := client.Ping(ctx)

	// Assertions
	if err == nil {
		t.Fatal("expected error, got nil")
	}

	// Check if it's an API error
	apiErr, ok := err.(*APIError)
	if !ok {
		t.Errorf("expected APIError, got %T", err)
	} else {
		if apiErr.Code != -1000 {
			t.Errorf("expected error code -1000, got %d", apiErr.Code)
		}
	}
}

func TestClient_GetTickerPrice_SingleSymbol_Success(t *testing.T) {
	client, mockClient := createTestClient()

	// Mock successful response for single symbol
	responseBody := `{"symbol":"BTCUSDT","price":"45000.00"}`
	mockClient.setResponse("https://api.binance.com/api/v3/ticker/price?symbol=BTCUSDT", http.StatusOK, responseBody)

	// Test GetTickerPrice for single symbol
	ctx := context.Background()
	result, err := client.GetTickerPrice(ctx, "BTCUSDT")

	// Assertions
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	if !result.IsSingle() {
		t.Fatal("expected single ticker result")
	}

	ticker := result.GetSingle()
	if ticker.Symbol != "BTCUSDT" {
		t.Errorf("expected symbol BTCUSDT, got %s", ticker.Symbol)
	}

	if ticker.Price != "45000.00" {
		t.Errorf("expected price 45000.00, got %s", ticker.Price)
	}
}

func TestClient_GetTickerPrice_AllSymbols_Success(t *testing.T) {
	client, mockClient := createTestClient()

	// Mock successful response for all symbols
	responseBody := `[{"symbol":"BTCUSDT","price":"45000.00"},{"symbol":"ETHUSDT","price":"3000.00"}]`
	mockClient.setResponse("https://api.binance.com/api/v3/ticker/price", http.StatusOK, responseBody)

	// Test GetTickerPrice for all symbols
	ctx := context.Background()
	result, err := client.GetTickerPrice(ctx, "")

	// Assertions
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	if !result.IsArray() {
		t.Fatal("expected array ticker result")
	}

	tickers := result.GetArray()
	if len(tickers) != 2 {
		t.Errorf("expected 2 tickers, got %d", len(tickers))
	}

	if tickers[0].Symbol != "BTCUSDT" {
		t.Errorf("expected first ticker symbol BTCUSDT, got %s", tickers[0].Symbol)
	}
}

func TestClient_GetExchangeInfo_Success(t *testing.T) {
	client, mockClient := createTestClient()

	// Mock successful response
	responseBody := `{
		"timezone":"UTC",
		"serverTime":1635724800000,
		"rateLimits":[],
		"symbols":[
			{
				"symbol":"BTCUSDT",
				"status":"TRADING",
				"baseAsset":"BTC",
				"baseAssetPrecision":8,
				"quoteAsset":"USDT",
				"quoteAssetPrecision":8,
				"orderTypes":["LIMIT","MARKET"],
				"icebergAllowed":true,
				"ocoAllowed":true,
				"isSpotTradingAllowed":true,
				"isMarginTradingAllowed":true,
				"filters":[],
				"permissions":["SPOT"]
			}
		]
	}`
	mockClient.setResponse("https://api.binance.com/api/v3/exchangeInfo", http.StatusOK, responseBody)

	// Test GetExchangeInfo
	ctx := context.Background()
	result, err := client.GetExchangeInfo(ctx)

	// Assertions
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	if result == nil {
		t.Fatal("expected result, got nil")
	}

	if result.Timezone != "UTC" {
		t.Errorf("expected timezone UTC, got %s", result.Timezone)
	}

	if len(result.Symbols) != 1 {
		t.Errorf("expected 1 symbol, got %d", len(result.Symbols))
	}

	if result.Symbols[0].Symbol != "BTCUSDT" {
		t.Errorf("expected symbol BTCUSDT, got %s", result.Symbols[0].Symbol)
	}
}

// Test with context cancellation
func TestClient_GetServerTime_ContextCancellation(t *testing.T) {
	client, mockClient := createTestClient()

	// Create cancelled context
	ctx, cancel := context.WithCancel(context.Background())
	cancel() // Cancel immediately

	// Mock the client to return context cancelled error
	mockClient.setError("https://api.binance.com/api/v3/time", context.Canceled)

	// Test GetServerTime with cancelled context
	result, err := client.GetServerTime(ctx)

	// Assertions
	if err == nil {
		t.Fatal("expected error due to context cancellation, got nil")
	}

	if result != nil {
		t.Errorf("expected nil result, got %v", result)
	}

	// Check error is context related
	if !strings.Contains(err.Error(), "context canceled") {
		t.Errorf("expected context cancellation error, got %v", err)
	}
}

func TestClient_GetKlines_Success(t *testing.T) {
	client, mockClient := createTestClient()

	// Mock successful response for klines
	responseBody := `[
		[1635724800000,"45000.00","45100.00","44900.00","45050.00","100.5",1635724859999,"4525025.00",1000,"50.25","50.25"],
		[1635724860000,"45050.00","45200.00","45000.00","45150.00","150.2",1635724919999,"6780300.00",1200,"75.10","75.10"]
	]`
	mockClient.setResponse("https://api.binance.com/api/v3/klines?interval=1m&limit=2&symbol=BTCUSDT", http.StatusOK, responseBody)

	// Test GetKlines
	ctx := context.Background()
	result, err := client.GetKlines(ctx, "BTCUSDT", "1m", 2)

	// Assertions
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	if result == nil {
		t.Fatal("expected result, got nil")
	}

	if len(*result) != 2 {
		t.Errorf("expected 2 klines, got %d", len(*result))
	}

	// Check first kline
	firstKline := (*result)[0]
	if firstKline.Open != "45000.00" {
		t.Errorf("expected open 45000.00, got %s", firstKline.Open)
	}

	if firstKline.High != "45100.00" {
		t.Errorf("expected high 45100.00, got %s", firstKline.High)
	}

	if firstKline.Low != "44900.00" {
		t.Errorf("expected low 44900.00, got %s", firstKline.Low)
	}

	if firstKline.Close != "45050.00" {
		t.Errorf("expected close 45050.00, got %s", firstKline.Close)
	}

	if firstKline.Volume != "100.5" {
		t.Errorf("expected volume 100.5, got %s", firstKline.Volume)
	}
}

func TestClient_GetOrderBook_Success(t *testing.T) {
	client, mockClient := createTestClient()

	// Mock successful response for order book
	responseBody := `{
		"lastUpdateId": 123456789,
		"bids": [
			["45000.00", "1.5"],
			["44999.00", "2.0"]
		],
		"asks": [
			["45001.00", "1.0"],
			["45002.00", "2.5"]
		]
	}`
	mockClient.setResponse("https://api.binance.com/api/v3/depth?limit=5&symbol=BTCUSDT", http.StatusOK, responseBody)

	// Test GetOrderBook
	ctx := context.Background()
	result, err := client.GetOrderBook(ctx, "BTCUSDT", 5)

	// Assertions
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	if result == nil {
		t.Fatal("expected result, got nil")
	}

	if result.LastUpdateId != 123456789 {
		t.Errorf("expected lastUpdateId 123456789, got %d", result.LastUpdateId)
	}

	if len(result.Bids) != 2 {
		t.Errorf("expected 2 bids, got %d", len(result.Bids))
	}

	if len(result.Asks) != 2 {
		t.Errorf("expected 2 asks, got %d", len(result.Asks))
	}

	// Check first bid
	if result.Bids[0][0] != "45000.00" {
		t.Errorf("expected first bid price 45000.00, got %s", result.Bids[0][0])
	}

	if result.Bids[0][1] != "1.5" {
		t.Errorf("expected first bid quantity 1.5, got %s", result.Bids[0][1])
	}

	// Check first ask
	if result.Asks[0][0] != "45001.00" {
		t.Errorf("expected first ask price 45001.00, got %s", result.Asks[0][0])
	}

	if result.Asks[0][1] != "1.0" {
		t.Errorf("expected first ask quantity 1.0, got %s", result.Asks[0][1])
	}
}

func TestClient_GetTrades_Success(t *testing.T) {
	client, mockClient := createTestClient()

	// Mock successful response for trades
	responseBody := `[
		{
			"id": 12345,
			"price": "45000.00",
			"qty": "1.5",
			"quoteQty": "67500.00",
			"time": 1635724800000,
			"isBuyerMaker": false,
			"isBestMatch": true
		},
		{
			"id": 12346,
			"price": "45001.00",
			"qty": "0.5",
			"quoteQty": "22500.50",
			"time": 1635724801000,
			"isBuyerMaker": true,
			"isBestMatch": false
		}
	]`
	mockClient.setResponse("https://api.binance.com/api/v3/trades?limit=2&symbol=BTCUSDT", http.StatusOK, responseBody)

	// Test GetTrades
	ctx := context.Background()
	result, err := client.GetTrades(ctx, "BTCUSDT", 2)

	// Assertions
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	if result == nil {
		t.Fatal("expected result, got nil")
	}

	if len(result) != 2 {
		t.Errorf("expected 2 trades, got %d", len(result))
	}

	// Check first trade
	firstTrade := result[0]
	if firstTrade.Id != 12345 {
		t.Errorf("expected trade id 12345, got %d", firstTrade.Id)
	}

	if firstTrade.Price != "45000.00" {
		t.Errorf("expected trade price 45000.00, got %s", firstTrade.Price)
	}

	if firstTrade.Qty != "1.5" {
		t.Errorf("expected trade qty 1.5, got %s", firstTrade.Qty)
	}

	if firstTrade.IsBuyerMaker != false {
		t.Errorf("expected isBuyerMaker false, got %v", firstTrade.IsBuyerMaker)
	}

	if firstTrade.IsBestMatch != true {
		t.Errorf("expected isBestMatch true, got %v", firstTrade.IsBestMatch)
	}
}
