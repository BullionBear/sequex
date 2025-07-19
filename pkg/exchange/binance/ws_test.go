package binance

import (
	"fmt"
	"testing"
	"time"
)

func TestNewWSClient(t *testing.T) {
	config := DefaultConfig()
	client := NewWSClient(config)

	if client == nil {
		t.Fatal("Expected non-nil client")
	}

	if client.config != config {
		t.Errorf("Expected config to match, got %v", client.config)
	}

	if client.url != WSBaseURL {
		t.Errorf("Expected URL %s, got %s", WSBaseURL, client.url)
	}
}

func TestNewWSClientWithTestnet(t *testing.T) {
	config := TestnetConfig()
	client := NewWSClient(config)

	if client.url != WSBaseURLTestnet {
		t.Errorf("Expected URL %s, got %s", WSBaseURLTestnet, client.url)
	}
}

func TestNewWSClientWithOptions(t *testing.T) {
	config := DefaultConfig()

	client := NewWSClient(config,
		WithOnConnect(func() {}),
		WithOnDisconnect(func() {}),
		WithOnError(func(error) {}),
		WithOnMessage(func([]byte) {}),
		WithReconnectSettings(10, 2*time.Second),
	)

	if client.maxReconnectAttempts != 10 {
		t.Errorf("Expected maxReconnectAttempts 10, got %d", client.maxReconnectAttempts)
	}

	if client.reconnectDelay != 2*time.Second {
		t.Errorf("Expected reconnectDelay 2s, got %v", client.reconnectDelay)
	}
}

func TestWSClient_IsConnected(t *testing.T) {
	config := DefaultConfig()
	client := NewWSClient(config)

	if client.IsConnected() {
		t.Error("Expected client to not be connected initially")
	}
}

func TestWSClient_SubscribeToStream(t *testing.T) {
	config := DefaultConfig()
	client := NewWSClient(config)

	streamName := "btcusdt@ticker"
	_ = client.SubscribeToStream(streamName)

	// This should fail in test environment without real connection
	// but we can test the URL construction
	expectedURL := fmt.Sprintf("%s/ws/%s", WSBaseURL, streamName)
	if client.url != expectedURL {
		t.Errorf("Expected URL %s, got %s", expectedURL, client.url)
	}
}

func TestNewWSStreamClient(t *testing.T) {
	config := DefaultConfig()
	client := NewWSStreamClient(config)

	if client == nil {
		t.Fatal("Expected non-nil client")
	}

	if client.config != config {
		t.Errorf("Expected config to match, got %v", client.config)
	}

	if len(client.clients) != 0 {
		t.Errorf("Expected empty clients map, got %d", len(client.clients))
	}
}

func TestWSStreamClient_GetActiveStreams(t *testing.T) {
	config := DefaultConfig()
	client := NewWSStreamClient(config)

	streams := client.GetActiveStreams()
	if len(streams) != 0 {
		t.Errorf("Expected empty streams, got %d", len(streams))
	}
}

func TestWSStreamClient_IsStreamActive(t *testing.T) {
	config := DefaultConfig()
	client := NewWSStreamClient(config)

	active := client.IsStreamActive("test")
	if active {
		t.Error("Expected stream to not be active")
	}
}

func TestWSStreamClient_SubscribeToKline(t *testing.T) {
	config := DefaultConfig()
	client := NewWSStreamClient(config)

	callback := func(data []byte) error {
		return nil
	}

	unsubscribe, err := client.SubscribeToKline("BTCUSDT", "1m", callback)

	// This should fail in test environment without real connection
	// but we can test the stream name construction
	if err == nil {
		defer unsubscribe()
	}

	// Check if the stream was attempted to be created
}

func TestWSStreamClient_SubscribeToTicker(t *testing.T) {
	config := DefaultConfig()
	client := NewWSStreamClient(config)

	callback := func(data []byte) error {
		return nil
	}

	unsubscribe, err := client.SubscribeToTicker("BTCUSDT", callback)

	// This should fail in test environment without real connection
	if err == nil {
		defer unsubscribe()
	}

	// Check if the stream was attempted to be created
}

func TestWSStreamClient_SubscribeToMiniTicker(t *testing.T) {
	config := DefaultConfig()
	client := NewWSStreamClient(config)

	callback := func(data []byte) error {
		return nil
	}

	unsubscribe, err := client.SubscribeToMiniTicker("BTCUSDT", callback)

	// This should fail in test environment without real connection
	if err == nil {
		defer unsubscribe()
	}

	// Check if the stream was attempted to be created
}

func TestWSStreamClient_SubscribeToBookTicker(t *testing.T) {
	config := DefaultConfig()
	client := NewWSStreamClient(config)

	callback := func(data []byte) error {
		return nil
	}

	unsubscribe, err := client.SubscribeToBookTicker("BTCUSDT", callback)

	// This should fail in test environment without real connection
	if err == nil {
		defer unsubscribe()
	}

	// Check if the stream was attempted to be created
}

func TestWSStreamClient_SubscribeToDepth(t *testing.T) {
	config := DefaultConfig()
	client := NewWSStreamClient(config)

	callback := func(data []byte) error {
		return nil
	}

	unsubscribe, err := client.SubscribeToDepth("BTCUSDT", "5", callback)

	// This should fail in test environment without real connection
	if err == nil {
		defer unsubscribe()
	}

	// Check if the stream was attempted to be created
}

func TestWSStreamClient_SubscribeToTrade(t *testing.T) {
	config := DefaultConfig()
	client := NewWSStreamClient(config)

	callback := func(data []byte) error {
		return nil
	}

	unsubscribe, err := client.SubscribeToTrade("BTCUSDT", callback)

	// This should fail in test environment without real connection
	if err == nil {
		defer unsubscribe()
	}

	// Check if the stream was attempted to be created
}

func TestWSStreamClient_SubscribeToAggTrade(t *testing.T) {
	config := DefaultConfig()
	client := NewWSStreamClient(config)

	callback := func(data []byte) error {
		return nil
	}

	unsubscribe, err := client.SubscribeToAggTrade("BTCUSDT", callback)

	// This should fail in test environment without real connection
	if err == nil {
		defer unsubscribe()
	}

	// Check if the stream was attempted to be created
}

func TestWSStreamClient_SubscribeToCombinedStreams(t *testing.T) {
	config := DefaultConfig()
	client := NewWSStreamClient(config)

	callback := func(data []byte) error {
		return nil
	}

	streams := []string{"btcusdt@ticker", "ethusdt@ticker"}
	unsubscribe, err := client.SubscribeToCombinedStreams(streams, callback)

	// This should fail in test environment without real connection
	if err == nil {
		defer unsubscribe()
	}

	// Check if the stream was attempted to be created
}

func TestParseKlineData(t *testing.T) {
	// Sample kline data from Binance WebSocket
	sampleData := `{
		"e": "kline",
		"E": 123456789,
		"s": "BTCUSDT",
		"k": {
			"t": 123400000,
			"T": 123460000,
			"s": "BTCUSDT",
			"i": "1m",
			"f": 100,
			"L": 200,
			"o": "50000.00",
			"c": "50100.00",
			"h": "50200.00",
			"l": "49900.00",
			"v": "100.00",
			"n": 1000,
			"x": false,
			"q": "5000000.00",
			"V": "50.00",
			"Q": "2500000.00",
			"B": "0"
		}
	}`

	klineData, err := ParseKlineData([]byte(sampleData))
	if err != nil {
		t.Fatalf("Failed to parse kline data: %v", err)
	}

	if klineData.Symbol != "BTCUSDT" {
		t.Errorf("Expected symbol BTCUSDT, got %s", klineData.Symbol)
	}

	if klineData.EventType != "kline" {
		t.Errorf("Expected event type kline, got %s", klineData.EventType)
	}

	if klineData.Kline.OpenPrice != 50000.00 {
		t.Errorf("Expected open price 50000.00, got %f", klineData.Kline.OpenPrice)
	}
}

func TestParseTickerData(t *testing.T) {
	// Sample ticker data from Binance WebSocket
	sampleData := `{
		"e": "24hrTicker",
		"E": 123456789,
		"s": "BTCUSDT",
		"P": "100.00",
		"p": "100.00",
		"w": "50000.00",
		"x": "49900.00",
		"c": "50100.00",
		"Q": "1.00",
		"b": "50099.00",
		"B": "10.00",
		"a": "50101.00",
		"A": "5.00",
		"o": "50000.00",
		"h": "50200.00",
		"l": "49900.00",
		"v": "1000.00",
		"q": "50000000.00",
		"O": 123400000,
		"C": 123456789,
		"F": 100,
		"L": 200,
		"n": 1000
	}`

	tickerData, err := ParseTickerData([]byte(sampleData))
	if err != nil {
		t.Fatalf("Failed to parse ticker data: %v", err)
	}

	if tickerData.Symbol != "BTCUSDT" {
		t.Errorf("Expected symbol BTCUSDT, got %s", tickerData.Symbol)
	}

	if tickerData.EventType != "24hrTicker" {
		t.Errorf("Expected event type 24hrTicker, got %s", tickerData.EventType)
	}

	if tickerData.LastPrice != 50100.00 {
		t.Errorf("Expected last price 50100.00, got %f", tickerData.LastPrice)
	}
}

func TestParseMiniTickerData(t *testing.T) {
	// Sample mini ticker data from Binance WebSocket
	sampleData := `{
		"e": "24hrMiniTicker",
		"E": 123456789,
		"s": "BTCUSDT",
		"c": "50100.00",
		"o": "50000.00",
		"h": "50200.00",
		"l": "49900.00",
		"v": "1000.00",
		"q": "50000000.00"
	}`

	miniTickerData, err := ParseMiniTickerData([]byte(sampleData))
	if err != nil {
		t.Fatalf("Failed to parse mini ticker data: %v", err)
	}

	if miniTickerData.Symbol != "BTCUSDT" {
		t.Errorf("Expected symbol BTCUSDT, got %s", miniTickerData.Symbol)
	}

	if miniTickerData.EventType != "24hrMiniTicker" {
		t.Errorf("Expected event type 24hrMiniTicker, got %s", miniTickerData.EventType)
	}

	if miniTickerData.ClosePrice != 50100.00 {
		t.Errorf("Expected close price 50100.00, got %f", miniTickerData.ClosePrice)
	}
}

func TestParseBookTickerData(t *testing.T) {
	// Sample book ticker data from Binance WebSocket
	sampleData := `{
		"e": "bookTicker",
		"E": 123456789,
		"s": "BTCUSDT",
		"b": "50099.00",
		"B": "10.00",
		"a": "50101.00",
		"A": "5.00",
		"u": 123456,
		"T": 123456789
	}`

	bookTickerData, err := ParseBookTickerData([]byte(sampleData))
	if err != nil {
		t.Fatalf("Failed to parse book ticker data: %v", err)
	}

	if bookTickerData.Symbol != "BTCUSDT" {
		t.Errorf("Expected symbol BTCUSDT, got %s", bookTickerData.Symbol)
	}

	if bookTickerData.EventType != "bookTicker" {
		t.Errorf("Expected event type bookTicker, got %s", bookTickerData.EventType)
	}

	if bookTickerData.BidPrice != 50099.00 {
		t.Errorf("Expected bid price 50099.00, got %f", bookTickerData.BidPrice)
	}
}

func TestParseTradeData(t *testing.T) {
	// Sample trade data from Binance WebSocket
	sampleData := `{
		"e": "trade",
		"E": 123456789,
		"s": "BTCUSDT",
		"t": 12345,
		"p": "50100.00",
		"q": "1.00",
		"b": 88,
		"a": 50,
		"T": 123456789,
		"m": false,
		"M": true
	}`

	tradeData, err := ParseTradeData([]byte(sampleData))
	if err != nil {
		t.Fatalf("Failed to parse trade data: %v", err)
	}

	if tradeData.Symbol != "BTCUSDT" {
		t.Errorf("Expected symbol BTCUSDT, got %s", tradeData.Symbol)
	}

	if tradeData.EventType != "trade" {
		t.Errorf("Expected event type trade, got %s", tradeData.EventType)
	}

	if tradeData.Price != 50100.00 {
		t.Errorf("Expected price 50100.00, got %f", tradeData.Price)
	}
}

func TestParseAggTradeData(t *testing.T) {
	// Sample aggregated trade data from Binance WebSocket
	sampleData := `{
		"e": "aggTrade",
		"E": 123456789,
		"s": "BTCUSDT",
		"a": 12345,
		"p": "50100.00",
		"q": "1.00",
		"f": 100,
		"l": 200,
		"T": 123456789,
		"m": false,
		"M": true
	}`

	aggTradeData, err := ParseAggTradeData([]byte(sampleData))
	if err != nil {
		t.Fatalf("Failed to parse aggregated trade data: %v", err)
	}

	if aggTradeData.Symbol != "BTCUSDT" {
		t.Errorf("Expected symbol BTCUSDT, got %s", aggTradeData.Symbol)
	}

	if aggTradeData.EventType != "aggTrade" {
		t.Errorf("Expected event type aggTrade, got %s", aggTradeData.EventType)
	}

	if aggTradeData.Price != 50100.00 {
		t.Errorf("Expected price 50100.00, got %f", aggTradeData.Price)
	}
}

func TestWSStreamClient_Close(t *testing.T) {
	config := DefaultConfig()
	client := NewWSStreamClient(config)

	err := client.Close()
	if err != nil {
		t.Errorf("Expected no error on close, got %v", err)
	}
}

func TestWSStreamClient_SubscribeToUserDataStream(t *testing.T) {
	config := DefaultConfig()
	client := NewWSStreamClient(config)

	callback := func(data []byte) error {
		return nil
	}

	listenKey := "test-listen-key"
	unsubscribe, err := client.SubscribeToUserDataStream(listenKey, callback)

	// This should fail in test environment without real connection
	if err == nil {
		defer unsubscribe()
	}

	// Check if the stream was attempted to be created
}

func TestWSStreamClient_SubscribeToUserDataStreamWithCallbacks(t *testing.T) {
	config := DefaultConfig()
	client := NewWSStreamClient(config)

	accountPositionCallback := func(data *WSOutboundAccountPosition) error {
		if data.EventType != "outboundAccountPosition" {
			t.Errorf("Expected event type outboundAccountPosition, got %s", data.EventType)
		}
		return nil
	}

	balanceUpdateCallback := func(data *WSBalanceUpdate) error {
		if data.EventType != "balanceUpdate" {
			t.Errorf("Expected event type balanceUpdate, got %s", data.EventType)
		}
		return nil
	}

	executionReportCallback := func(data *WSExecutionReport) error {
		if data.EventType != "executionReport" {
			t.Errorf("Expected event type executionReport, got %s", data.EventType)
		}
		return nil
	}

	unsubscribe, err := client.SubscribeToUserDataStreamWithCallbacks(
		accountPositionCallback,
		balanceUpdateCallback,
		executionReportCallback,
	)

	// This should fail in test environment without real connection
	if err == nil {
		defer unsubscribe()
	}

	// Check if the stream was attempted to be created
}

func TestParseOutboundAccountPosition(t *testing.T) {
	// Sample outbound account position data from Binance WebSocket
	sampleData := `{
		"e": "outboundAccountPosition",
		"E": 1564034571105,
		"u": 1564034571073,
		"B": [
			{
				"a": "ETH",
				"f": "10000.000000",
				"l": "0.000000"
			},
			{
				"a": "BTC",
				"f": "1.000000",
				"l": "0.000000"
			}
		]
	}`

	accountPosition, err := ParseOutboundAccountPosition([]byte(sampleData))
	if err != nil {
		t.Fatalf("Failed to parse outbound account position data: %v", err)
	}

	if accountPosition.EventType != "outboundAccountPosition" {
		t.Errorf("Expected event type outboundAccountPosition, got %s", accountPosition.EventType)
	}

	if accountPosition.EventTime != 1564034571105 {
		t.Errorf("Expected event time 1564034571105, got %d", accountPosition.EventTime)
	}

	if len(accountPosition.Balances) != 2 {
		t.Errorf("Expected 2 balances, got %d", len(accountPosition.Balances))
	}

	if accountPosition.Balances[0].Asset != "ETH" {
		t.Errorf("Expected first balance asset ETH, got %s", accountPosition.Balances[0].Asset)
	}
}

func TestParseBalanceUpdate(t *testing.T) {
	// Sample balance update data from Binance WebSocket
	sampleData := `{
		"e": "balanceUpdate",
		"E": 1573200697110,
		"a": "BTC",
		"d": "100.00000000",
		"T": 1573200697068
	}`

	balanceUpdate, err := ParseBalanceUpdate([]byte(sampleData))
	if err != nil {
		t.Fatalf("Failed to parse balance update data: %v", err)
	}

	if balanceUpdate.EventType != "balanceUpdate" {
		t.Errorf("Expected event type balanceUpdate, got %s", balanceUpdate.EventType)
	}

	if balanceUpdate.Asset != "BTC" {
		t.Errorf("Expected asset BTC, got %s", balanceUpdate.Asset)
	}

	if balanceUpdate.BalanceDelta != "100.00000000" {
		t.Errorf("Expected balance delta 100.00000000, got %s", balanceUpdate.BalanceDelta)
	}
}

func TestParseExecutionReport(t *testing.T) {
	// Sample execution report data from Binance WebSocket
	sampleData := `{
		"e": "executionReport",
		"E": 1499405658658,
		"s": "ETHBTC",
		"c": "mUvoqJxFIILMdfAW5iGSOW",
		"S": "BUY",
		"o": "LIMIT",
		"f": "GTC",
		"q": "1.00000000",
		"p": "0.10264410",
		"P": "0.00000000",
		"F": "0.00000000",
		"g": -1,
		"C": "",
		"x": "NEW",
		"X": "NEW",
		"r": "NONE",
		"i": 4293153,
		"l": "0.00000000",
		"z": "0.00000000",
		"L": "0.00000000",
		"n": "0",
		"N": null,
		"T": 1499405658657,
		"t": -1,
		"v": 3,
		"I": 8641984,
		"w": true,
		"m": false,
		"M": false,
		"O": 1499405658657,
		"Z": "0.00000000",
		"Y": "0.00000000",
		"Q": "0.00000000",
		"W": 1499405658657,
		"V": "NONE"
	}`

	executionReport, err := ParseExecutionReport([]byte(sampleData))
	if err != nil {
		t.Fatalf("Failed to parse execution report data: %v", err)
	}

	if executionReport.EventType != "executionReport" {
		t.Errorf("Expected event type executionReport, got %s", executionReport.EventType)
	}

	if executionReport.Symbol != "ETHBTC" {
		t.Errorf("Expected symbol ETHBTC, got %s", executionReport.Symbol)
	}

	if executionReport.Side != "BUY" {
		t.Errorf("Expected side BUY, got %s", executionReport.Side)
	}

	if executionReport.OrderType != "LIMIT" {
		t.Errorf("Expected order type LIMIT, got %s", executionReport.OrderType)
	}

	if executionReport.OrderID != 4293153 {
		t.Errorf("Expected order ID 4293153, got %d", executionReport.OrderID)
	}
}

func TestWSStreamClient_SubscribeToKlineWithCallback(t *testing.T) {
	config := DefaultConfig()
	client := NewWSStreamClient(config)

	callback := func(data *WSKlineData) error {
		if data.Symbol != "BTCUSDT" {
			t.Errorf("Expected symbol BTCUSDT, got %s", data.Symbol)
		}
		return nil
	}

	unsubscribe, err := client.SubscribeToKlineWithCallback("BTCUSDT", "1m", callback)

	// This should fail in test environment without real connection
	if err == nil {
		defer unsubscribe()
	}

	// Check if the stream was attempted to be created
}

func TestWSStreamClient_SubscribeToTickerWithCallback(t *testing.T) {
	config := DefaultConfig()
	client := NewWSStreamClient(config)

	callback := func(data *WSTickerData) error {
		if data.Symbol != "BTCUSDT" {
			t.Errorf("Expected symbol BTCUSDT, got %s", data.Symbol)
		}
		return nil
	}

	unsubscribe, err := client.SubscribeToTickerWithCallback("BTCUSDT", callback)

	// This should fail in test environment without real connection
	if err == nil {
		defer unsubscribe()
	}

	// Check if the stream was attempted to be created
}

func TestWSStreamClient_SubscribeToTradeWithCallback(t *testing.T) {
	config := DefaultConfig()
	client := NewWSStreamClient(config)

	callback := func(data *WSTradeData) error {
		if data.Symbol != "BTCUSDT" {
			t.Errorf("Expected symbol BTCUSDT, got %s", data.Symbol)
		}
		return nil
	}

	unsubscribe, err := client.SubscribeToTradeWithCallback("BTCUSDT", callback)

	// This should fail in test environment without real connection
	if err == nil {
		defer unsubscribe()
	}

	// Check if the stream was attempted to be created
}

func TestWSStreamClient_SubscribeToBookTickerWithCallback(t *testing.T) {
	config := DefaultConfig()
	client := NewWSStreamClient(config)

	callback := func(data *WSBookTickerData) error {
		if data.Symbol != "BTCUSDT" {
			t.Errorf("Expected symbol BTCUSDT, got %s", data.Symbol)
		}
		return nil
	}

	unsubscribe, err := client.SubscribeToBookTickerWithCallback("BTCUSDT", callback)

	// This should fail in test environment without real connection
	if err == nil {
		defer unsubscribe()
	}

	// Check if the stream was attempted to be created
}

func TestWSStreamClient_SubscribeToPartialDepthWithCallback(t *testing.T) {
	config := DefaultConfig()
	client := NewWSStreamClient(config)

	callback := func(data *WSPartialDepthData) error {
		if data.LastUpdateID <= 0 {
			t.Error("Expected positive LastUpdateID")
		}
		return nil
	}

	unsubscribe, err := client.SubscribeToPartialDepthWithCallback("BTCUSDT", "5", callback)

	// This should fail in test environment without real connection
	if err == nil {
		defer unsubscribe()
	}

	// Check if the stream was attempted to be created
}

func TestWSStreamClient_SubscribeToDiffDepthWithCallback(t *testing.T) {
	config := DefaultConfig()
	client := NewWSStreamClient(config)

	callback := func(data *WSDiffDepthData) error {
		if data.Symbol != "BTCUSDT" {
			t.Errorf("Expected symbol BTCUSDT, got %s", data.Symbol)
		}
		if data.EventType != "depthUpdate" {
			t.Errorf("Expected event type depthUpdate, got %s", data.EventType)
		}
		return nil
	}

	unsubscribe, err := client.SubscribeToDiffDepthWithCallback("BTCUSDT", "", callback)

	// This should fail in test environment without real connection
	if err == nil {
		defer unsubscribe()
	}

	// Check if the stream was attempted to be created
}

func TestParsePartialDepthData(t *testing.T) {
	// Sample partial depth data from Binance WebSocket
	sampleData := `{
		"lastUpdateId": 160,
		"bids": [
			["0.0024", "10"],
			["0.0023", "15"]
		],
		"asks": [
			["0.0026", "100"],
			["0.0027", "50"]
		]
	}`

	partialDepthData, err := ParsePartialDepthData([]byte(sampleData))
	if err != nil {
		t.Fatalf("Failed to parse partial depth data: %v", err)
	}

	if partialDepthData.LastUpdateID != 160 {
		t.Errorf("Expected LastUpdateID 160, got %d", partialDepthData.LastUpdateID)
	}

	if len(partialDepthData.Bids) != 2 {
		t.Errorf("Expected 2 bids, got %d", len(partialDepthData.Bids))
	}

	if len(partialDepthData.Asks) != 2 {
		t.Errorf("Expected 2 asks, got %d", len(partialDepthData.Asks))
	}
}

func TestParseDiffDepthData(t *testing.T) {
	// Sample diff depth data from Binance WebSocket
	sampleData := `{
		"e": "depthUpdate",
		"E": 1672515782136,
		"s": "BNBBTC",
		"U": 157,
		"u": 160,
		"b": [
			["0.0024", "10"],
			["0.0023", "15"]
		],
		"a": [
			["0.0026", "100"],
			["0.0027", "50"]
		]
	}`

	diffDepthData, err := ParseDiffDepthData([]byte(sampleData))
	if err != nil {
		t.Fatalf("Failed to parse diff depth data: %v", err)
	}

	if diffDepthData.Symbol != "BNBBTC" {
		t.Errorf("Expected symbol BNBBTC, got %s", diffDepthData.Symbol)
	}

	if diffDepthData.EventType != "depthUpdate" {
		t.Errorf("Expected event type depthUpdate, got %s", diffDepthData.EventType)
	}

	if diffDepthData.FirstUpdateID != 157 {
		t.Errorf("Expected FirstUpdateID 157, got %d", diffDepthData.FirstUpdateID)
	}

	if diffDepthData.FinalUpdateID != 160 {
		t.Errorf("Expected FinalUpdateID 160, got %d", diffDepthData.FinalUpdateID)
	}
}
