package binance

import (
	"github.com/BullionBear/sequex/internal/exchange"
	"github.com/BullionBear/sequex/pkg/exchange/binance"
)

func toExchangeGTC(timeInForce string) exchange.TimeInForce {
	switch timeInForce {
	case binance.TimeInForceGTC:
		return exchange.TimeInForceGTC
	case binance.TimeInForceIOC:
		return exchange.TimeInForceIOC
	case binance.TimeInForceFOK:
		return exchange.TimeInForceFOK
	default:
		return exchange.TimeInForceUnknown
	}
}

func toExchangeSymbol(symbol string) exchange.Symbol {

	return exchange.Symbol{
		Base:  symbol,
		Quote: symbol,
	}
}
