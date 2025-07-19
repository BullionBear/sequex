package binance

import (
	"encoding/json"
	"fmt"
)

// APIError represents a Binance API error response
type APIError struct {
	Code    int    `json:"code"`
	Message string `json:"msg"`
}

// Error implements the error interface
func (e *APIError) Error() string {
	return fmt.Sprintf("binance api error: code=%d, msg=%s", e.Code, e.Message)
}

// Common Binance API error codes
const (
	// General errors
	ErrCodeUnknown                 = -1000
	ErrCodeDisconnected            = -1001
	ErrCodeUnauthorized            = -1002
	ErrCodeTooManyRequests         = -1003
	ErrCodeUnexpectedResp          = -1006
	ErrCodeTimeout                 = -1007
	ErrCodeUnknownOrderComposition = -1014
	ErrCodeTooManyOrders           = -1015
	ErrCodeServiceShuttingDown     = -1016
	ErrCodeUnsupportedOperation    = -1020
	ErrCodeInvalidTimestamp        = -1021
	ErrCodeInvalidSignature        = -1022

	// Request errors
	ErrCodeIllegalChars           = -1100
	ErrCodeTooManyParameters      = -1101
	ErrCodeMandatoryParamEmpty    = -1102
	ErrCodeUnknownParam           = -1103
	ErrCodeUnreadParameters       = -1104
	ErrCodeParamEmpty             = -1105
	ErrCodeParamNotRequired       = -1106
	ErrCodeInvalidParam           = -1111
	ErrCodeInvalidSymbol          = -1121
	ErrCodeInvalidListenKey       = -1125
	ErrCodeMoreThanXXHours        = -1127
	ErrCodeOptionalParamsBadCombo = -1128
	ErrCodeInvalidParameter       = -1130

	// Order errors
	ErrCodeNewOrderRejected = -2010
	ErrCodeCancelRejected   = -2011
	ErrCodeNoSuchOrder      = -2013
	ErrCodeBadAPIKeyFmt     = -2014
	ErrCodeRejectedMBXKey   = -2015
	ErrCodeNoTradingSymbol  = -2016
)

// ErrorCodeMessages maps error codes to human-readable messages
var ErrorCodeMessages = map[int]string{
	ErrCodeUnknown:                 "Unknown error occurred while processing the request",
	ErrCodeDisconnected:            "Internal error; unable to process your request. Please try again",
	ErrCodeUnauthorized:            "You are not authorized to execute this request",
	ErrCodeTooManyRequests:         "Too many requests queued",
	ErrCodeUnexpectedResp:          "Unexpected response format from internal service",
	ErrCodeTimeout:                 "Timeout waiting for response from internal service",
	ErrCodeUnknownOrderComposition: "Unknown order composition",
	ErrCodeTooManyOrders:           "Too many orders",
	ErrCodeServiceShuttingDown:     "This service is no longer available",
	ErrCodeUnsupportedOperation:    "This operation is not supported",
	ErrCodeInvalidTimestamp:        "Timestamp for this request is outside of the recvWindow",
	ErrCodeInvalidSignature:        "Signature for this request is not valid",
	ErrCodeIllegalChars:            "Illegal characters found in a parameter",
	ErrCodeTooManyParameters:       "Too many parameters sent for this endpoint",
	ErrCodeMandatoryParamEmpty:     "A mandatory parameter was not sent, was empty/null, or malformed",
	ErrCodeUnknownParam:            "An unknown parameter was sent",
	ErrCodeUnreadParameters:        "Not all sent parameters were read",
	ErrCodeParamEmpty:              "A parameter was empty",
	ErrCodeParamNotRequired:        "A parameter was sent when not required",
	ErrCodeInvalidParam:            "Invalid parameter",
	ErrCodeInvalidSymbol:           "Invalid symbol",
	ErrCodeInvalidListenKey:        "Invalid listen key",
	ErrCodeMoreThanXXHours:         "Lookup interval is too big",
	ErrCodeOptionalParamsBadCombo:  "Combination of optional parameters invalid",
	ErrCodeInvalidParameter:        "Invalid data sent for a parameter",
	ErrCodeNewOrderRejected:        "New order rejected",
	ErrCodeCancelRejected:          "Cancel order rejected",
	ErrCodeNoSuchOrder:             "No such order",
	ErrCodeBadAPIKeyFmt:            "API-key format invalid",
	ErrCodeRejectedMBXKey:          "Invalid API-key, IP, or permissions for action",
	ErrCodeNoTradingSymbol:         "No trading symbol",
}

// ParseAPIError attempts to parse an API error from the response body
func ParseAPIError(body []byte) *APIError {
	var apiErr APIError
	if err := json.Unmarshal(body, &apiErr); err != nil {
		// If we can't parse the error, return a generic one
		return &APIError{
			Code:    ErrCodeUnknown,
			Message: string(body),
		}
	}
	return &apiErr
}

// IsRetryableError determines if an error is retryable
func IsRetryableError(err error) bool {
	if apiErr, ok := err.(*APIError); ok {
		switch apiErr.Code {
		case ErrCodeTooManyRequests, ErrCodeTimeout, ErrCodeDisconnected:
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
