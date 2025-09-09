package sqx

import (
	"fmt"
	"strings"

	"github.com/BullionBear/sequex/internal/model/protobuf"
)

type Exchange int

const (
	ExchangeUnknown Exchange = iota
	ExchangeBinance
	ExchangeBinancePerp
	ExchangeBybit
)

func (e Exchange) ToProtobuf() protobuf.Exchange {
	return protobuf.Exchange(e)
}

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
	InstrumentTypeMargin
	InstrumentTypePerp
	InstrumentTypeInverse
	InstrumentTypeFutures
	InstrumentTypeOption
)

func (i InstrumentType) ToProtobuf() protobuf.Instrument {
	return protobuf.Instrument(i)
}

func NewInstrumentType(instrumentType string) InstrumentType {
	switch strings.ToUpper(instrumentType) {
	case "SPOT":
		return InstrumentTypeSpot
	case "MARGIN":
		return InstrumentTypeMargin
	case "PERP":
		return InstrumentTypePerp
	case "INVERSE":
		return InstrumentTypeInverse
	case "FUTURES":
		return InstrumentTypeFutures
	case "OPTION":
		return InstrumentTypeOption
	}
	return InstrumentTypeUnknown
}

func (i InstrumentType) String() string {
	return []string{"UNKNOWN", "SPOT", "MARGIN", "PERP", "INVERSE", "FUTURES", "OPTION"}[i]
}

type Symbol struct {
	Base  string
	Quote string
}

func (s Symbol) ToProtobuf() protobuf.Symbol {
	return protobuf.Symbol{
		Base:  s.Base,
		Quote: s.Quote,
	}
}

func NewSymbol(base, quote string) Symbol {
	return Symbol{
		Base:  strings.ToUpper(base),
		Quote: strings.ToUpper(quote),
	}
}

func NewSymbolFromStr(symbol string) (Symbol, error) {
	parts := strings.Split(symbol, "-")
	if len(parts) < 2 {
		return Symbol{}, fmt.Errorf("invalid symbol: %s", symbol)
	}
	return Symbol{
		Base:  strings.ToUpper(parts[0]),
		Quote: strings.ToUpper(parts[1]),
	}, nil
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

func (s Side) ToProtobuf() protobuf.Side {
	return protobuf.Side(s)
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

func (t TimeInForce) ToProtobuf() protobuf.TimeInForce {
	return protobuf.TimeInForce(t)
}

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

type DataType int

const (
	DataTypeUnknown DataType = iota
	DataTypeTrade
	DataTypeDepth
	DataTypeOrder
)

func NewDataType(dataType string) DataType {
	switch strings.ToUpper(dataType) {
	case "TRADE":
		return DataTypeTrade
	case "DEPTH":
		return DataTypeDepth
	case "ORDER":
		return DataTypeOrder
	}
	return DataTypeUnknown
}

func (d DataType) String() string {
	return []string{"UNKNOWN", "TRADE", "DEPTH", "ORDER"}[d]
}
