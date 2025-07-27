package binanceperp

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
	resp, status, err := doUnsignedGet(cfg, "/fapi/v1/time", nil)
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
	resp, status, err := doUnsignedGet(cfg, "/fapi/v1/depth", params)
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
	params := map[string]string{}
	resp, status, err := doSignedRequest(cfg, "GET", "/fapi/v2/account", params)
	if err != nil {
		t.Fatalf("doSignedRequest error: %v", err)
	}
	if status != 200 {
		t.Fatalf("expected status 200, got %d", status)
	}
	var data struct {
		Assets []struct {
			Asset                  string `json:"asset"`
			WalletBalance          string `json:"walletBalance"`
			UnrealizedProfit       string `json:"unrealizedProfit"`
			MarginBalance          string `json:"marginBalance"`
			MaintMargin            string `json:"maintMargin"`
			InitialMargin          string `json:"initialMargin"`
			PositionInitialMargin  string `json:"positionInitialMargin"`
			OpenOrderInitialMargin string `json:"openOrderInitialMargin"`
		} `json:"assets"`
		Positions []struct {
			Symbol                 string `json:"symbol"`
			InitialMargin          string `json:"initialMargin"`
			MaintMargin            string `json:"maintMargin"`
			UnrealizedProfit       string `json:"unrealizedProfit"`
			PositionInitialMargin  string `json:"positionInitialMargin"`
			OpenOrderInitialMargin string `json:"openOrderInitialMargin"`
			Leverage               string `json:"leverage"`
			Isolated               bool   `json:"isolated"`
			EntryPrice             string `json:"entryPrice"`
			MaxNotional            string `json:"maxNotional"`
			PositionSide           string `json:"positionSide"`
			PositionAmt            string `json:"positionAmt"`
		} `json:"positions"`
	}
	err = json.Unmarshal(resp, &data)
	if err != nil {
		t.Fatalf("json unmarshal error: %v", err)
	}
	if len(data.Assets) == 0 {
		t.Error("assets is empty, expected at least one asset entry")
	}
}

// Test unhappy path - invalid symbol for unsigned request
func TestDoUnsignedGet_InvalidSymbol(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping integration test in short mode.")
	}
	cfg := &Config{
		BaseURL: MainnetBaseUrl,
	}
	params := map[string]string{
		"symbol": "INVALIDSYMBOL",
		"limit":  "5",
	}
	_, status, err := doUnsignedGet(cfg, "/fapi/v1/depth", params)
	if err != nil {
		t.Fatalf("doUnsignedGet error: %v", err)
	}
	if status == 200 {
		t.Error("expected non-200 status for invalid symbol, got 200")
	}
}

// Test unhappy path - invalid credentials for signed request
func TestDoSignedRequest_InvalidCredentials(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping integration test in short mode.")
	}
	cfg := &Config{
		APIKey:    "invalid_api_key",
		APISecret: "invalid_api_secret",
		BaseURL:   MainnetBaseUrl,
	}
	params := map[string]string{}
	_, status, err := doSignedRequest(cfg, "GET", "/fapi/v2/account", params)
	if err != nil {
		t.Fatalf("doSignedRequest error: %v", err)
	}
	if status == 200 {
		t.Error("expected non-200 status for invalid credentials, got 200")
	}
}

// Test buildQueryString function
func TestBuildQueryString(t *testing.T) {
	params := map[string]string{
		"symbol":      "BTCUSDT",
		"side":        "BUY",
		"type":        "LIMIT",
		"quantity":    "1",
		"price":       "9000",
		"timeInForce": "GTC",
		"timestamp":   "1591702613943",
	}
	result := buildQueryString(params)
	expected := "price=9000&quantity=1&side=BUY&symbol=BTCUSDT&timeInForce=GTC&timestamp=1591702613943&type=LIMIT"
	if result != expected {
		t.Errorf("buildQueryString failed.\nExpected: %s\nGot: %s", expected, result)
	}
}

// Test signParams function
func TestSignParams(t *testing.T) {
	query := "symbol=BTCUSDT&side=BUY&type=LIMIT&quantity=1&price=9000&timeInForce=GTC&recvWindow=5000&timestamp=1591702613943"
	secret := "2b5eb11e18796d12d88f13dc27dbbd02c2cc51ff7059765ed9821957d82bb4d9"
	result := signParams(query, secret)
	expected := "3c661234138461fcc7a7d8746c6558c9842d4e10870d2ecbedf7777cad694af9"
	if result != expected {
		t.Errorf("signParams failed.\nExpected: %s\nGot: %s", expected, result)
	}
}
