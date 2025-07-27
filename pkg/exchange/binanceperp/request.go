package binanceperp

import (
	"bytes"
	"crypto/hmac"
	"crypto/sha256"
	"encoding/hex"
	"io"
	"net/http"
	"net/url"
	"sort"
	"strconv"
	"strings"
	"time"
)

// doUnsignedGet performs unsigned GET request (public endpoints)
func doUnsignedGet(cfg *Config, endpoint string, params map[string]string) ([]byte, int, error) {
	baseURL := strings.TrimRight(cfg.BaseURL, "/")
	fullURL := baseURL + endpoint
	if len(params) > 0 {
		q := url.Values{}
		for k, v := range params {
			q.Set(k, v)
		}
		fullURL += "?" + q.Encode()
	}
	resp, err := http.Get(fullURL)
	if err != nil {
		return nil, 0, err
	}
	defer resp.Body.Close()
	body, err := io.ReadAll(resp.Body)
	return body, resp.StatusCode, err
}

// doSignedRequest performs signed request (GET/POST/PUT/DELETE) for TRADE and USER_DATA endpoints
func doSignedRequest(cfg *Config, method, endpoint string, params map[string]string) ([]byte, int, error) {
	baseURL := strings.TrimRight(cfg.BaseURL, "/")
	fullURL := baseURL + endpoint

	// Add timestamp and recvWindow for timing security
	params["timestamp"] = strconv.FormatInt(time.Now().UnixNano()/1e6, 10)
	if _, exists := params["recvWindow"]; !exists {
		params["recvWindow"] = "5000" // Default recvWindow of 5000ms
	}

	// Build query string for signing
	queryString := buildQueryString(params)
	// Sign the query string using HMAC SHA256
	signature := signParams(queryString, cfg.APISecret)
	params["signature"] = signature

	// Prepare request
	var req *http.Request
	var err error
	if method == http.MethodGet || method == http.MethodDelete {
		// For GET/DELETE, put all params in query string
		q := buildQueryString(params)
		fullURL += "?" + q
		req, err = http.NewRequest(method, fullURL, nil)
	} else {
		// For POST/PUT, put params in request body
		q := buildQueryString(params)
		req, err = http.NewRequest(method, fullURL, bytes.NewBufferString(q))
		req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	}
	if err != nil {
		return nil, 0, err
	}
	// Set API key header
	req.Header.Set("X-MBX-APIKEY", cfg.APIKey)

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return nil, 0, err
	}
	defer resp.Body.Close()
	body, err := io.ReadAll(resp.Body)
	return body, resp.StatusCode, err
}

// doAPIKeyOnlyRequest handles requests that only need API key header (no signing)
// Used for MARKET_DATA and USER_STREAM endpoints
func doAPIKeyOnlyRequest(cfg *Config, method, endpoint string, params map[string]string) ([]byte, int, error) {
	baseURL := strings.TrimRight(cfg.BaseURL, "/")
	fullURL := baseURL + endpoint

	// Build query string from params (no timestamp or signature added)
	if len(params) > 0 {
		q := url.Values{}
		for k, v := range params {
			q.Set(k, v)
		}
		fullURL += "?" + q.Encode()
	}

	req, err := http.NewRequest(method, fullURL, nil)
	if err != nil {
		return nil, 0, err
	}
	// Set API key header
	req.Header.Set("X-MBX-APIKEY", cfg.APIKey)

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return nil, 0, err
	}
	defer resp.Body.Close()
	body, err := io.ReadAll(resp.Body)
	return body, resp.StatusCode, err
}

// buildQueryString sorts and encodes params according to Binance requirements
func buildQueryString(params map[string]string) string {
	var keys []string
	for k := range params {
		keys = append(keys, k)
	}
	sort.Strings(keys)
	var q []string
	for _, k := range keys {
		q = append(q, url.QueryEscape(k)+"="+url.QueryEscape(params[k]))
	}
	return strings.Join(q, "&")
}

// signParams creates HMAC SHA256 signature using secret key
func signParams(query, secret string) string {
	mac := hmac.New(sha256.New, []byte(secret))
	mac.Write([]byte(query))
	return hex.EncodeToString(mac.Sum(nil))
}
