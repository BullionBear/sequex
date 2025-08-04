package trade

import (
	"log/slog"

	"github.com/BullionBear/sequex/internal/exchange"
)

func NewTradeEmitter[C exchange.IsolatedSpotConnector](name string, symbol exchange.Symbol) *TradeEmitter[C] {
	connector, err := exchange.NewConnector[C](exchange.MarketTypeBinance, exchange.Credentials{
		APIKey:    "",
		APISecret: "",
	})
	if err != nil {
		slog.Error("trade emitter new", "error", err)
		return nil
	}
	return &TradeEmitter[C]{
		name:      name,
		connector: connector,

		symbol:      symbol,
		unsubscribe: nil,
	}
}

type TradeEmitter[C exchange.IsolatedSpotConnector] struct {
	name        string
	connector   C
	unsubscribe func()
	symbol      exchange.Symbol
}

func (e *TradeEmitter[C]) Type() string {
	return "TradeEmitter"
}

func (e *TradeEmitter[C]) Name() string {
	return e.name
}

func (e *TradeEmitter[C]) Start(opts exchange.TradeSubscriptionOptions) error {
	unsubscribe, err := e.connector.SubscribeTrades(e.symbol, opts)
	if err != nil {
		slog.Error("trade emitter subscribe trades", "error", err)
		return err
	}
	e.unsubscribe = unsubscribe
	return nil
}

func (e *TradeEmitter[C]) Stop() error {
	if e.unsubscribe != nil {
		e.unsubscribe()
	}
	return nil
}
