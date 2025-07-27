package binanceperp

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
)

// Client is the Binance Perpetual Futures API client.
type Client struct {
	cfg *Config
}

// NewClient creates a new Binance Perpetual Futures API client.
func NewClient(cfg *Config) *Client {
	return &Client{cfg: cfg}
}

// GetServerTime tests connectivity to the Rest API and gets the current server time.
func (c *Client) GetServerTime(ctx context.Context) (Response[GetServerTimeResponse], error) {
	body, status, err := doUnsignedGet(c.cfg, PathGetServerTime, nil)
	if err != nil {
		return Response[GetServerTimeResponse]{}, err
	}
	if status != http.StatusOK {
		return Response[GetServerTimeResponse]{Code: status, Message: string(body)}, fmt.Errorf("http error: %d", status)
	}
	var resp GetServerTimeResponse
	if err := json.Unmarshal(body, &resp); err != nil {
		return Response[GetServerTimeResponse]{}, err
	}
	return Response[GetServerTimeResponse]{Code: 0, Message: "success", Data: &resp}, nil
}
