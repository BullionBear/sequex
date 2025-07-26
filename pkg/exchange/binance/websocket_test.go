package binance

import (
	"encoding/json"
	"testing"
	"time"
)

type AggTradeEvent struct {
	EventType     string `json:"e"`
	EventTime     int64  `json:"E"`
	Symbol        string `json:"s"`
	AggTradeID    int64  `json:"a"`
	Price         string `json:"p"`
	Quantity      string `json:"q"`
	FirstTradeID  int64  `json:"f"`
	LastTradeID   int64  `json:"l"`
	TradeTime     int64  `json:"T"`
	IsMarketMaker bool   `json:"m"`
	Ignore        bool   `json:"M"`
}

func TestBinanceWSConn_AggTradePayload(t *testing.T) {
	baseURL := MainnetWSBaseUrl
	stream := "/ws/btcusdt@aggTrade"
	conn := NewBinanceWSConn(baseURL, stream)

	// Add message handler to the connection
	msgCh := make(chan AggTradeEvent, 1)
	conn.OnMessage = func(data []byte) {
		var event AggTradeEvent
		if err := json.Unmarshal(data, &event); err == nil && event.EventType == "aggTrade" {
			select {
			case msgCh <- event:
			default:
			}
		}
	}

	if err := conn.Connect(); err != nil {
		t.Fatalf("Failed to connect: %v", err)
	}

	// Wait for message or timeout
	select {
	case event := <-msgCh:
		if event.Symbol == "BTCUSDT" || event.Symbol == "btcusdt" {
			t.Logf("Received aggTrade: %+v", event)
		} else {
			t.Fatalf("Unexpected symbol: %v", event.Symbol)
		}
	case <-time.After(5 * time.Second):
		t.Fatal("Timeout: did not receive aggTrade event in 5 seconds")
	}

	// Properly disconnect
	conn.Disconnect()
	time.Sleep(100 * time.Millisecond) // Give time for graceful shutdown
}
