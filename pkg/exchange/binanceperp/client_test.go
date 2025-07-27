package binanceperp

import (
	"context"
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
