package binance

import (
	"encoding/json"
	"os"
	"testing"
)

func TestDoUnsignedGet_ServerTime(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping integration test in short mode.")
	}
	cfg := &Config{
		BaseURL: MainnetBaseUrl,
	}
	resp, status, err := doUnsignedGet(cfg, "/v3/time", nil)
	if err != nil {
		t.Fatalf("doUnsignedGet error: %v", err)
	}
	if status != 200 {
		t.Fatalf("expected status 200, got %d", status)
	}
	var data struct {
		ServerTime int64 `json:"serverTime"`
	}
	err = json.Unmarshal(resp, &data)
	if err != nil {
		t.Fatalf("json unmarshal error: %v", err)
	}
	if data.ServerTime == 0 {
		t.Error("serverTime is zero, expected non-zero value")
	}
}

func TestDoUnsignedGet_OrderBook(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping integration test in short mode.")
	}
	cfg := &Config{
		BaseURL: MainnetBaseUrl,
	}
	params := map[string]string{
		"symbol": "BTCUSDT",
		"limit":  "5",
	}
	resp, status, err := doUnsignedGet(cfg, "/v3/depth", params)
	if err != nil {
		t.Fatalf("doUnsignedGet error: %v", err)
	}
	if status != 200 {
		t.Fatalf("expected status 200, got %d", status)
	}
	var data struct {
		LastUpdateId int        `json:"lastUpdateId"`
		Bids         [][]string `json:"bids"`
		Asks         [][]string `json:"asks"`
	}
	err = json.Unmarshal(resp, &data)
	if err != nil {
		t.Fatalf("json unmarshal error: %v", err)
	}
	if len(data.Bids) == 0 {
		t.Error("bids is empty, expected at least one bid")
	}
	if len(data.Asks) == 0 {
		t.Error("asks is empty, expected at least one ask")
	}
}

func TestDoSignedRequest_AccountInfo(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping integration test in short mode.")
	}
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
	params := map[string]string{}
	resp, status, err := doSignedRequest(cfg, "GET", "/v3/account", params)
	if err != nil {
		t.Fatalf("doSignedRequest error: %v", err)
	}
	if status != 200 {
		t.Fatalf("expected status 200, got %d", status)
	}
	var data struct {
		Balances []struct {
			Asset  string `json:"asset"`
			Free   string `json:"free"`
			Locked string `json:"locked"`
		} `json:"balances"`
	}
	err = json.Unmarshal(resp, &data)
	if err != nil {
		t.Fatalf("json unmarshal error: %v", err)
	}
	if len(data.Balances) == 0 {
		t.Error("balances is empty, expected at least one balance entry")
	}
}

func TestDoSignedRequest_PostTestOrderWithCommissionRates(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping integration test in short mode.")
	}
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
	params := map[string]string{
		"symbol":      "ADAUSDT",
		"side":        "BUY",
		"type":        "LIMIT",
		"quantity":    "15",
		"price":       "0.5",
		"timeInForce": "GTC",
	}
	_, status, err := doSignedRequest(cfg, "POST", "/v3/order/test", params)
	if err != nil {
		t.Fatalf("doSignedRequest error: %v", err)
	}
	if status != 200 {
		t.Fatalf("expected status 200, got %d", status)
	}
}
