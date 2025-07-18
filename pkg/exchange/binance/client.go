package binance

import (
	"context"
	"crypto/hmac"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strconv"
	"strings"
	"time"
)

// Client represents the Binance API client
type Client struct {
	config     *Config
	httpClient *http.Client
	baseURL    string
}

// NewClient creates a new Binance API client
func NewClient(config *Config) (*Client, error) {
	if config == nil {
		config = DefaultConfig()
	}

	// For public endpoints, API credentials are not required
	// They are only required for authenticated endpoints

	httpClient := &http.Client{
		Timeout: time.Duration(config.Timeout) * time.Second,
	}

	client := &Client{
		config:     config,
		httpClient: httpClient,
		baseURL:    config.GetBaseURL(),
	}

	return client, nil
}

// GetConfig returns the client configuration
func (c *Client) GetConfig() *Config {
	return c.config
}

// GetHTTPClient returns the underlying HTTP client
func (c *Client) GetHTTPClient() *http.Client {
	return c.httpClient
}

// generateSignature creates HMAC SHA256 signature for authenticated requests
func (c *Client) generateSignature(queryString string) string {
	mac := hmac.New(sha256.New, []byte(c.config.APISecret))
	mac.Write([]byte(queryString))
	return hex.EncodeToString(mac.Sum(nil))
}

// buildAuthenticatedURL builds URL with signature for authenticated requests
func (c *Client) buildAuthenticatedURL(endpoint string, params url.Values) string {
	// Add timestamp
	params.Set("timestamp", strconv.FormatInt(time.Now().UnixMilli(), 10))

	// Create query string
	queryString := params.Encode()

	// Generate signature
	signature := c.generateSignature(queryString)
	params.Set("signature", signature)

	return fmt.Sprintf("%s%s?%s", c.baseURL, endpoint, params.Encode())
}

// doRequest performs HTTP request with proper headers
func (c *Client) doRequest(ctx context.Context, method, endpoint string, params url.Values, authenticated bool) ([]byte, error) {
	var reqURL string

	if authenticated {
		if !c.config.IsValid() {
			return nil, fmt.Errorf("authenticated request requires valid API credentials")
		}
		reqURL = c.buildAuthenticatedURL(endpoint, params)
	} else {
		if params != nil && len(params) > 0 {
			reqURL = fmt.Sprintf("%s%s?%s", c.baseURL, endpoint, params.Encode())
		} else {
			reqURL = fmt.Sprintf("%s%s", c.baseURL, endpoint)
		}
	}

	req, err := http.NewRequestWithContext(ctx, method, reqURL, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}

	// Set headers
	req.Header.Set("Content-Type", "application/json")
	if authenticated {
		req.Header.Set("X-MBX-APIKEY", c.config.APIKey)
	}

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("failed to execute request: %w", err)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read response body: %w", err)
	}

	if resp.StatusCode != http.StatusOK {
		var apiErr APIError
		if json.Unmarshal(body, &apiErr) == nil {
			return nil, &apiErr
		}
		return nil, fmt.Errorf("API error: status %d, body: %s", resp.StatusCode, string(body))
	}

	return body, nil
}

// doPostRequest performs POST request with form data
func (c *Client) doPostRequest(ctx context.Context, endpoint string, params url.Values) ([]byte, error) {
	if !c.config.IsValid() {
		return nil, fmt.Errorf("POST request requires valid API credentials")
	}

	// Add timestamp
	params.Set("timestamp", strconv.FormatInt(time.Now().UnixMilli(), 10))

	// Create query string for signature
	queryString := params.Encode()

	// Generate signature
	signature := c.generateSignature(queryString)
	params.Set("signature", signature)

	reqURL := fmt.Sprintf("%s%s", c.baseURL, endpoint)

	req, err := http.NewRequestWithContext(ctx, "POST", reqURL, strings.NewReader(params.Encode()))
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}

	// Set headers
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	req.Header.Set("X-MBX-APIKEY", c.config.APIKey)

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("failed to execute request: %w", err)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read response body: %w", err)
	}

	if resp.StatusCode != http.StatusOK {
		var apiErr APIError
		if json.Unmarshal(body, &apiErr) == nil {
			return nil, &apiErr
		}
		return nil, fmt.Errorf("API error: status %d, body: %s", resp.StatusCode, string(body))
	}

	return body, nil
}

// Ping tests connectivity to the REST API
func (c *Client) Ping(ctx context.Context) error {
	_, err := c.doRequest(ctx, "GET", "/api/v3/ping", nil, false)
	return err
}

// GetServerTime returns the server time
func (c *Client) GetServerTime(ctx context.Context) (int64, error) {
	body, err := c.doRequest(ctx, "GET", "/api/v3/time", nil, false)
	if err != nil {
		return 0, err
	}

	var result struct {
		ServerTime int64 `json:"serverTime"`
	}

	if err := json.Unmarshal(body, &result); err != nil {
		return 0, fmt.Errorf("failed to unmarshal response: %w", err)
	}

	return result.ServerTime, nil
}

// GetExchangeInfo returns current exchange trading rules and symbol information
func (c *Client) GetExchangeInfo(ctx context.Context) (*ExchangeInfo, error) {
	body, err := c.doRequest(ctx, "GET", "/api/v3/exchangeInfo", nil, false)
	if err != nil {
		return nil, err
	}

	var result ExchangeInfo
	if err := json.Unmarshal(body, &result); err != nil {
		return nil, fmt.Errorf("failed to unmarshal response: %w", err)
	}

	return &result, nil
}

// GetTicker24hr returns 24hr ticker price change statistics for a symbol
func (c *Client) GetTicker24hr(ctx context.Context, symbol string) (*Ticker24hr, error) {
	params := url.Values{}
	if symbol != "" {
		params.Set("symbol", symbol)
	}

	body, err := c.doRequest(ctx, "GET", "/api/v3/ticker/24hr", params, false)
	if err != nil {
		return nil, err
	}

	var result Ticker24hr
	if err := json.Unmarshal(body, &result); err != nil {
		return nil, fmt.Errorf("failed to unmarshal response: %w", err)
	}

	return &result, nil
}

// GetOrderBook returns order book depth for a symbol
func (c *Client) GetOrderBook(ctx context.Context, symbol string, limit int) (*OrderBook, error) {
	params := url.Values{}
	params.Set("symbol", symbol)
	if limit > 0 {
		params.Set("limit", strconv.Itoa(limit))
	}

	body, err := c.doRequest(ctx, "GET", "/api/v3/depth", params, false)
	if err != nil {
		return nil, err
	}

	var result OrderBook
	if err := json.Unmarshal(body, &result); err != nil {
		return nil, fmt.Errorf("failed to unmarshal response: %w", err)
	}

	return &result, nil
}

// GetRecentTrades returns recent trades for a symbol
func (c *Client) GetRecentTrades(ctx context.Context, symbol string, limit int) ([]Trade, error) {
	params := url.Values{}
	params.Set("symbol", symbol)
	if limit > 0 {
		params.Set("limit", strconv.Itoa(limit))
	}

	body, err := c.doRequest(ctx, "GET", "/api/v3/trades", params, false)
	if err != nil {
		return nil, err
	}

	var result []Trade
	if err := json.Unmarshal(body, &result); err != nil {
		return nil, fmt.Errorf("failed to unmarshal response: %w", err)
	}

	return result, nil
}

// GetKlines returns kline/candlestick data for a symbol
func (c *Client) GetKlines(ctx context.Context, symbol, interval string, limit int, startTime, endTime *int64) ([]Kline, error) {
	params := url.Values{}
	params.Set("symbol", symbol)
	params.Set("interval", interval)

	if limit > 0 {
		params.Set("limit", strconv.Itoa(limit))
	}
	if startTime != nil {
		params.Set("startTime", strconv.FormatInt(*startTime, 10))
	}
	if endTime != nil {
		params.Set("endTime", strconv.FormatInt(*endTime, 10))
	}

	body, err := c.doRequest(ctx, "GET", "/api/v3/klines", params, false)
	if err != nil {
		return nil, err
	}

	var rawResult [][]interface{}
	if err := json.Unmarshal(body, &rawResult); err != nil {
		return nil, fmt.Errorf("failed to unmarshal response: %w", err)
	}

	result := make([]Kline, len(rawResult))
	for i, klineData := range rawResult {
		if len(klineData) < 12 {
			continue
		}

		result[i] = Kline{
			OpenTime:                 int64(klineData[0].(float64)),
			Open:                     mustParseDecimal(klineData[1].(string)),
			High:                     mustParseDecimal(klineData[2].(string)),
			Low:                      mustParseDecimal(klineData[3].(string)),
			Close:                    mustParseDecimal(klineData[4].(string)),
			Volume:                   mustParseDecimal(klineData[5].(string)),
			CloseTime:                int64(klineData[6].(float64)),
			QuoteAssetVolume:         mustParseDecimal(klineData[7].(string)),
			NumberOfTrades:           int(klineData[8].(float64)),
			TakerBuyBaseAssetVolume:  mustParseDecimal(klineData[9].(string)),
			TakerBuyQuoteAssetVolume: mustParseDecimal(klineData[10].(string)),
		}
	}

	return result, nil
}

// GetAccount returns account information
func (c *Client) GetAccount(ctx context.Context) (*Account, error) {
	params := url.Values{}

	body, err := c.doRequest(ctx, "GET", "/api/v3/account", params, true)
	if err != nil {
		return nil, err
	}

	var result Account
	if err := json.Unmarshal(body, &result); err != nil {
		return nil, fmt.Errorf("failed to unmarshal response: %w", err)
	}

	return &result, nil
}

// GetOrder returns order information
func (c *Client) GetOrder(ctx context.Context, symbol string, orderID *int64, origClientOrderID *string) (*Order, error) {
	params := url.Values{}
	params.Set("symbol", symbol)

	if orderID != nil {
		params.Set("orderId", strconv.FormatInt(*orderID, 10))
	}
	if origClientOrderID != nil {
		params.Set("origClientOrderId", *origClientOrderID)
	}

	body, err := c.doRequest(ctx, "GET", "/api/v3/order", params, true)
	if err != nil {
		return nil, err
	}

	var result Order
	if err := json.Unmarshal(body, &result); err != nil {
		return nil, fmt.Errorf("failed to unmarshal response: %w", err)
	}

	return &result, nil
}

// GetOpenOrders returns all open orders for a symbol
func (c *Client) GetOpenOrders(ctx context.Context, symbol string) ([]Order, error) {
	params := url.Values{}
	if symbol != "" {
		params.Set("symbol", symbol)
	}

	body, err := c.doRequest(ctx, "GET", "/api/v3/openOrders", params, true)
	if err != nil {
		return nil, err
	}

	var result []Order
	if err := json.Unmarshal(body, &result); err != nil {
		return nil, fmt.Errorf("failed to unmarshal response: %w", err)
	}

	return result, nil
}

// CreateOrder places a new order
func (c *Client) CreateOrder(ctx context.Context, symbol, side, orderType string, quantity, price *string, timeInForce *string) (*NewOrderResponse, error) {
	params := url.Values{}
	params.Set("symbol", symbol)
	params.Set("side", side)
	params.Set("type", orderType)

	if quantity != nil {
		params.Set("quantity", *quantity)
	}
	if price != nil {
		params.Set("price", *price)
	}
	if timeInForce != nil {
		params.Set("timeInForce", *timeInForce)
	}

	body, err := c.doPostRequest(ctx, "/api/v3/order", params)
	if err != nil {
		return nil, err
	}

	var result NewOrderResponse
	if err := json.Unmarshal(body, &result); err != nil {
		return nil, fmt.Errorf("failed to unmarshal response: %w", err)
	}

	return &result, nil
}

// CancelOrder cancels an active order
func (c *Client) CancelOrder(ctx context.Context, symbol string, orderID *int64, origClientOrderID *string) (*Order, error) {
	params := url.Values{}
	params.Set("symbol", symbol)

	if orderID != nil {
		params.Set("orderId", strconv.FormatInt(*orderID, 10))
	}
	if origClientOrderID != nil {
		params.Set("origClientOrderId", *origClientOrderID)
	}

	body, err := c.doRequest(ctx, "DELETE", "/api/v3/order", params, true)
	if err != nil {
		return nil, err
	}

	var result Order
	if err := json.Unmarshal(body, &result); err != nil {
		return nil, fmt.Errorf("failed to unmarshal response: %w", err)
	}

	return &result, nil
}

// GetTrades returns trades for a specific account and symbol
func (c *Client) GetTrades(ctx context.Context, symbol string, limit int, fromID *int64) ([]Trade, error) {
	params := url.Values{}
	params.Set("symbol", symbol)

	if limit > 0 {
		params.Set("limit", strconv.Itoa(limit))
	}
	if fromID != nil {
		params.Set("fromId", strconv.FormatInt(*fromID, 10))
	}

	body, err := c.doRequest(ctx, "GET", "/api/v3/myTrades", params, true)
	if err != nil {
		return nil, err
	}

	var result []Trade
	if err := json.Unmarshal(body, &result); err != nil {
		return nil, fmt.Errorf("failed to unmarshal response: %w", err)
	}

	return result, nil
}
