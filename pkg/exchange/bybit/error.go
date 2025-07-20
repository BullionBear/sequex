package bybit

import (
	"encoding/json"
	"fmt"
)

// APIError represents a Bybit API error response
type APIError struct {
	Code    int    `json:"code"`
	Message string `json:"msg"`
}

// Error implements the error interface
func (e *APIError) Error() string {
	return fmt.Sprintf("bybit api error: code=%d, msg=%s", e.Code, e.Message)
}

// Common Bybit API error codes (HTTP status codes)
const (
	// HTTP 400 - Bad request
	ErrCodeBadRequest = 400

	// HTTP 401 - Invalid request
	ErrCodeInvalidRequest = 401

	// HTTP 403 - Forbidden request
	ErrCodeForbidden = 403

	// HTTP 404 - Cannot find path
	ErrCodeNotFound = 404

	// HTTP 429 - System level frequency protection
	ErrCodeTooManyRequests = 429
)

// ErrorCodeMessages maps error codes to human-readable messages
var ErrorCodeMessages = map[int]string{
	ErrCodeBadRequest:      "Bad request. Need to send the request with GET / POST (must be capitalized)",
	ErrCodeInvalidRequest:  "Invalid request. 1. Need to use the correct key to access; 2. Need to put authentication params in the request header",
	ErrCodeForbidden:       "Forbidden request. Possible causes: 1. IP rate limit breached; 2. You send GET request with an empty json body; 3. You are using U.S IP",
	ErrCodeNotFound:        "Cannot find path. Possible causes: 1. Wrong path; 2. Category value does not match account mode",
	ErrCodeTooManyRequests: "System level frequency protection. Please retry when encounter this",
}

// ParseAPIError attempts to parse an API error from the response body
func ParseAPIError(body []byte) *APIError {
	var apiErr APIError
	if err := json.Unmarshal(body, &apiErr); err != nil {
		// If we can't parse the error, return a generic one
		return &APIError{
			Code:    ErrCodeBadRequest,
			Message: string(body),
		}
	}
	return &apiErr
}

// IsRetryableError determines if an error is retryable
func IsRetryableError(err error) bool {
	if apiErr, ok := err.(*APIError); ok {
		switch apiErr.Code {
		case ErrCodeTooManyRequests:
			return true
		}
	}
	return false
}

// GetErrorMessage returns a human-readable error message for a given error code
func GetErrorMessage(code int) string {
	if msg, exists := ErrorCodeMessages[code]; exists {
		return msg
	}
	return "Unknown error code"
}
