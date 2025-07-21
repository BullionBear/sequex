package binance

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

// unsigned GET request (public endpoints)
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

// signed request (GET/POST/PUT/DELETE)
func doSignedRequest(cfg *Config, method, endpoint string, params map[string]string) ([]byte, int, error) {
	baseURL := strings.TrimRight(cfg.BaseURL, "/")
	fullURL := baseURL + endpoint

	// Add timestamp and recvWindow
	params["timestamp"] = strconv.FormatInt(time.Now().UnixNano()/1e6, 10)
	if cfg.RecvWindow > 0 {
		params["recvWindow"] = strconv.FormatInt(cfg.RecvWindow, 10)
	}

	// Build query string
	queryString := buildQueryString(params)
	// Sign
	signature := signParams(queryString, cfg.APISecret)
	params["signature"] = signature

	// Prepare request
	var req *http.Request
	var err error
	if method == http.MethodGet || method == http.MethodDelete {
		q := buildQueryString(params)
		fullURL += "?" + q
		req, err = http.NewRequest(method, fullURL, nil)
	} else {
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

// buildQueryString sorts and encodes params
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

// signParams creates HMAC SHA256 signature
func signParams(query, secret string) string {
	mac := hmac.New(sha256.New, []byte(secret))
	mac.Write([]byte(query))
	return hex.EncodeToString(mac.Sum(nil))
}
