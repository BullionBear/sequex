package sqx

import (
	"fmt"
	"strings"
)

type Exchange int

const (
	ExchangeUnknown Exchange = iota
	ExchangeBinance
	ExchangeBinancePerp
	ExchangeBybit
)

func (e Exchange) String() string {
	return []string{"UNKNOWN", "BINANCE", "BINANCE_PERP", "BYBIT"}[e]
}

func NewExchange(exchange string) Exchange {
	switch strings.ToUpper(exchange) {
	case "BINANCE":
		return ExchangeBinance
	case "BINANCE_PERP":
		return ExchangeBinancePerp
	case "BYBIT":
		return ExchangeBybit
	}
	return ExchangeUnknown
}

type InstrumentType int

const (
	InstrumentTypeUnknown InstrumentType = iota
	InstrumentTypeSpot
	InstrumentTypePerp
	InstrumentTypeInverse
)

func NewInstrumentType(instrumentType string) InstrumentType {
	switch strings.ToUpper(instrumentType) {
	case "SPOT":
		return InstrumentTypeSpot
	case "PERP":
		return InstrumentTypePerp
	case "INVERSE":
		return InstrumentTypeInverse
	}
	return InstrumentTypeUnknown
}

func (i InstrumentType) String() string {
	return []string{"UNKNOWN", "SPOT", "PERP", "INVERSE"}[i]
}

type Symbol struct {
	Base  string
	Quote string
}

func NewSymbol(base, quote string) Symbol {
	return Symbol{
		Base:  strings.ToUpper(base),
		Quote: strings.ToUpper(quote),
	}
}

func (s Symbol) String() string {
	return fmt.Sprintf("%s-%s", s.Base, s.Quote)
}

type Side int

const (
	SideUnknown Side = iota
	SideBuy
	SideSell
)

func NewSide(side string) Side {
	switch strings.ToUpper(side) {
	case "BUY":
		return SideBuy
	case "SELL":
		return SideSell
	}
	return SideUnknown
}

func (s Side) String() string {
	return []string{"UNKNOWN", "BUY", "SELL"}[s]
}

type TimeInForce int

const (
	TimeInForceUnknown TimeInForce = iota
	TimeInForceGTC
	TimeInForceIOC
	TimeInForceFOK
)

func NewTimeInForce(timeInForce string) TimeInForce {
	switch strings.ToUpper(timeInForce) {
	case "GTC":
		return TimeInForceGTC
	case "IOC":
		return TimeInForceIOC
	case "FOK":
		return TimeInForceFOK
	}
	return TimeInForceUnknown
}

func (t TimeInForce) String() string {
	return []string{"UNKNOWN", "GTC", "IOC", "FOK"}[t]
}
