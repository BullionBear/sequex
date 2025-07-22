package binance

import (
	"context"
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
