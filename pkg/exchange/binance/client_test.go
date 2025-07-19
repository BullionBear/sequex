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

	ticker, ok := result.(*TickerPriceResponse)
	if !ok {
		t.Fatalf("expected *TickerPriceResponse, got %T", result)
	}

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

	tickers, ok := result.([]TickerPriceResponse)
	if !ok {
		t.Fatalf("expected []TickerPriceResponse, got %T", result)
	}

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
