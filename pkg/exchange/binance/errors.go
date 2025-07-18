package binance

import "fmt"

// APIError represents Binance API error response
type APIError struct {
	Code int    `json:"code"`
	Msg  string `json:"msg"`
}

// Error implements the error interface
func (e *APIError) Error() string {
	return fmt.Sprintf("Binance API error %d: %s", e.Code, e.Msg)
}

// Common Binance API error codes
const (
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
	ErrCodeIllegalChars            = -1100
	ErrCodeTooManyParameters       = -1101
	ErrCodeMandatoryParamEmpty     = -1102
	ErrCodeUnknownParam            = -1103
	ErrCodeUnreadParameters        = -1104
	ErrCodeParamEmpty              = -1105
	ErrCodeParamNotRequired        = -1106
	ErrCodeNoDepth                 = -1112
	ErrCodeTIFNotRequired          = -1114
	ErrCodeInvalidTIF              = -1115
	ErrCodeInvalidOrderType        = -1116
	ErrCodeInvalidSide             = -1117
	ErrCodeEmptyNewClOrdID         = -1118
	ErrCodeEmptyOrgClOrdID         = -1119
	ErrCodeBadInterval             = -1120
	ErrCodeBadSymbol               = -1121
	ErrCodeInvalidListenKey        = -1125
	ErrCodeMoreThanXXHours         = -1127
	ErrCodeOptionalParamsBadCombo  = -1128
	ErrCodeInvalidParameter        = -1130
)

// IsRetryableError checks if the error is retryable
func IsRetryableError(err error) bool {
	if apiErr, ok := err.(*APIError); ok {
		switch apiErr.Code {
		case ErrCodeTooManyRequests, ErrCodeTimeout, ErrCodeDisconnected:
			return true
		}
	}
	return false
}
