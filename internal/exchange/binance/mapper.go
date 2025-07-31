package binance

import (
	"fmt"
	"log"

	"github.com/BullionBear/sequex/internal/exchange"
	"github.com/BullionBear/sequex/pkg/exchange/binance"
)

func toExchangeTimeInForce(timeInForce string) exchange.TimeInForce {
	switch timeInForce {
	case binance.TimeInForceGTC:
		return exchange.TimeInForceGTC
	case binance.TimeInForceIOC:
		return exchange.TimeInForceIOC
	case binance.TimeInForceFOK:
		return exchange.TimeInForceFOK
	default:
		log.Printf("Unknown time in force: %s", timeInForce)
		return exchange.TimeInForceUnknown
	}
}

func toBinanceTimeInForce(timeInForce exchange.TimeInForce) string {
	switch timeInForce {
	case exchange.TimeInForceGTC:
		return binance.TimeInForceGTC
	case exchange.TimeInForceIOC:
		return binance.TimeInForceIOC
	case exchange.TimeInForceFOK:
		return binance.TimeInForceFOK
	default:
		log.Printf("Unknown time in force: %s", timeInForce)
		return ""
	}
}

func toExchangeSymbol(symbol string) (exchange.Symbol, error) {
	binanceSymbol, err := GetSymbol(symbol)
	if err != nil {
		return exchange.Symbol{}, err
	}
	return exchange.Symbol{
		Base:  binanceSymbol.BaseAsset,
		Quote: binanceSymbol.QuoteAsset,
	}, nil
}

func toBianceSymbol(symbol exchange.Symbol) string {
	return fmt.Sprintf("%s%s", symbol.Base, symbol.Quote)
}

func toExchangeOrderStatus(status string) exchange.OrderStatus {
	switch status {
	case binance.OrderStatusNew:
		return exchange.OrderStatusNew
	case binance.OrderStatusPartiallyFilled:
		return exchange.OrderStatusPartiallyFilled
	case binance.OrderStatusFilled:
		return exchange.OrderStatusFilled
	case binance.OrderStatusCanceled:
		return exchange.OrderStatusCanceled
	case binance.OrderStatusRejected:
		return exchange.OrderStatusRejected
	default:
		log.Printf("Unknown order status: %s", status)
		return exchange.OrderStatusUnknown
	}
}

func toExchangeOrderType(orderType string) exchange.OrderType {
	switch orderType {
	case binance.OrderTypeLimit:
		return exchange.OrderTypeLimit
	case binance.OrderTypeMarket:
		return exchange.OrderTypeMarket
	case binance.OrderTypeLimitMaker:
		return exchange.OrderTypeLimitMaker
	case binance.OrderTypeStopLoss:
		return exchange.OrderTypeStopMarket
	default:
		log.Printf("Unknown order type: %s", orderType)
		return exchange.OrderTypeUnknown
	}
}

func toExchangeOrderSide(side string) exchange.OrderSide {
	switch side {
	case binance.OrderSideBuy:
		return exchange.OrderSideBuy
	case binance.OrderSideSell:
		return exchange.OrderSideSell
	default:
		log.Printf("Unknown order side: %s", side)
		return exchange.OrderSideUnknown
	}
}

func toBinanceOrderSide(side exchange.OrderSide) string {
	switch side {
	case exchange.OrderSideBuy:
		return binance.OrderSideBuy
	case exchange.OrderSideSell:
		return binance.OrderSideSell
	default:
		log.Printf("Unknown order side: %s", side)
		return ""
	}
}
