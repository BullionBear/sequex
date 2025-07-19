package binance

import (
	"bytes"
	"context"
	"crypto/hmac"
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strconv"
	"time"
)

// HTTPClient interface for HTTP requests (useful for testing)
type HTTPClient interface {
	Do(req *http.Request) (*http.Response, error)
}

// RequestService handles HTTP requests to Binance API
type RequestService struct {
	config     *Config
	httpClient HTTPClient
}

// NewRequestService creates a new RequestService
func NewRequestService(config *Config) *RequestService {
	httpClient := &http.Client{
		Timeout: config.Timeout,
	}

	return &RequestService{
		config:     config,
		httpClient: httpClient,
	}
}

// RequestOptions holds options for making HTTP requests
type RequestOptions struct {
	Method       string
	Endpoint     string
	QueryParams  url.Values
	Body         []byte
	SecurityType string
}

// DoRequest performs an HTTP request to Binance API
func (r *RequestService) DoRequest(ctx context.Context, opts *RequestOptions) ([]byte, error) {
	// Build URL
	fullURL := r.config.BaseURL + opts.Endpoint

	// Add query parameters
	if len(opts.QueryParams) > 0 {
		fullURL += "?" + opts.QueryParams.Encode()
	}

	// Create HTTP request
	var body io.Reader
	if opts.Body != nil {
		body = bytes.NewReader(opts.Body)
	}

	req, err := http.NewRequestWithContext(ctx, opts.Method, fullURL, body)
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}

	// Set headers
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("User-Agent", "sequex-binance-client/1.0")

	// Handle authentication based on security type
	if err := r.setAuthHeaders(req, opts); err != nil {
		return nil, fmt.Errorf("failed to set auth headers: %w", err)
	}

	// Perform request
	resp, err := r.httpClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("failed to perform request: %w", err)
	}
	defer resp.Body.Close()

	// Read response body
	respBody, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read response body: %w", err)
	}

	// Check for API errors
	if resp.StatusCode != http.StatusOK {
		apiErr := ParseAPIError(respBody)
		return nil, apiErr
	}

	return respBody, nil
}

// setAuthHeaders sets authentication headers based on security type
func (r *RequestService) setAuthHeaders(req *http.Request, opts *RequestOptions) error {
	switch opts.SecurityType {
	case SecurityTypeNone:
		// No authentication required
		return nil

	case SecurityTypeTradeKey, SecurityTypeMarketData:
		// API key required
		if r.config.APIKey == "" {
			return fmt.Errorf("API key required for security type %s", opts.SecurityType)
		}
		req.Header.Set(HeaderAPIKey, r.config.APIKey)
		return nil

	case SecurityTypeSigned:
		// Signature required
		if r.config.APIKey == "" || r.config.APISecret == "" {
			return fmt.Errorf("API key and secret required for signed requests")
		}

		// Set API key header
		req.Header.Set(HeaderAPIKey, r.config.APIKey)

		// Add timestamp to query parameters
		timestamp := strconv.FormatInt(time.Now().UnixNano()/int64(time.Millisecond), 10)
		query := req.URL.Query()
		query.Set(HeaderTimestamp, timestamp)

		// Add recvWindow if not set (default 5000ms)
		if !query.Has(HeaderRecvWindow) {
			query.Set(HeaderRecvWindow, "5000")
		}

		// Create signature
		signature := r.createSignature(query.Encode())
		query.Set(HeaderSignature, signature)

		// Update request URL with signed parameters
		req.URL.RawQuery = query.Encode()

		return nil

	default:
		return fmt.Errorf("unknown security type: %s", opts.SecurityType)
	}
}

// createSignature creates HMAC SHA256 signature for signed requests
func (r *RequestService) createSignature(queryString string) string {
	mac := hmac.New(sha256.New, []byte(r.config.APISecret))
	mac.Write([]byte(queryString))
	return hex.EncodeToString(mac.Sum(nil))
}

// DoUnsignedRequest performs an unsigned request (for public endpoints)
func (r *RequestService) DoUnsignedRequest(ctx context.Context, method, endpoint string, params url.Values) ([]byte, error) {
	opts := &RequestOptions{
		Method:       method,
		Endpoint:     endpoint,
		QueryParams:  params,
		SecurityType: SecurityTypeNone,
	}
	return r.DoRequest(ctx, opts)
}

// DoSignedRequest performs a signed request (for private endpoints)
func (r *RequestService) DoSignedRequest(ctx context.Context, method, endpoint string, params url.Values) ([]byte, error) {
	opts := &RequestOptions{
		Method:       method,
		Endpoint:     endpoint,
		QueryParams:  params,
		SecurityType: SecurityTypeSigned,
	}
	return r.DoRequest(ctx, opts)
}

// DoAPIKeyRequest performs a request with API key (for market data endpoints)
func (r *RequestService) DoAPIKeyRequest(ctx context.Context, method, endpoint string, params url.Values) ([]byte, error) {
	opts := &RequestOptions{
		Method:       method,
		Endpoint:     endpoint,
		QueryParams:  params,
		SecurityType: SecurityTypeMarketData,
	}
	return r.DoRequest(ctx, opts)
}
