package binanceperp

// Response is the unified response wrapper for all endpoints.
type Response[T any] struct {
	Code    int    `json:"code"`
	Message string `json:"msg"`
	Data    *T     `json:"data,omitempty"`
}

// GetServerTimeResponse represents the server time response.
type GetServerTimeResponse struct {
	ServerTime int64 `json:"serverTime"`
}
