package order

import "github.com/adshao/go-binance/v2"

func fromSide(side Side) binance.SideType {
	switch side {
	case SideBuy:
		return "BUY"
	case SideSell:
		return "SELL"
	default:
		return "UNKNOWN"
	}
}

func toSide(side binance.SideType) Side {
	switch side {
	case "BUY":
		return SideBuy
	case "SELL":
		return SideSell
	default:
		return SideUnknown
	}
}

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
