package order

import "github.com/adshao/go-binance/v2"

func toOrderStatus(status binance.OrderStatusType) OrderStatus {
	switch status {
	case "NEW":
		return OrderStatusNew
	case "PARTIALLY_FILLED":
		return OrderStatusPartiallyFilled
	case "FILLED":
		return OrderStatusFilled
	case "PENDING_CANCEL":
		return OrderStatusPendingCancel
	case "CANCELED":
		return OrderStatusCanceled
	case "REJECTED":
		return OrderStatusRejected
	case "EXPIRED":
		return OrderStatusExpired
	default:
		return OrderStatusUnknown
	}
}
