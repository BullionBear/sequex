package trade

import (
	"log/slog"

	"github.com/BullionBear/sequex/internal/exchange"
)

func NewTradeEmitter(name string, connector exchange.IsolatedSpotConnector, symbol exchange.Symbol) *TradeEmitter {
	return &TradeEmitter{
		name:      name,
		connector: connector,

		symbol:      symbol,
		unsubscribe: nil,
	}
}

type TradeEmitter struct {
	name        string
	connector   exchange.IsolatedSpotConnector
	unsubscribe func()
	symbol      exchange.Symbol
}

func (e *TradeEmitter) Type() string {
	return "TradeEmitter"
}

func (e *TradeEmitter) Name() string {
	return e.name
}

func (e *TradeEmitter) Start() error {
	unsubscribe, err := e.connector.SubscribeTrades(e.symbol, exchange.TradeSubscriptionOptions{
		OnConnect:    func() { slog.Info("trade emitter connect") },
		OnReconnect:  func() { slog.Info("trade emitter reconnect") },
		OnError:      func(err error) { slog.Error("trade emitter error", "error", err) },
		OnTrade:      func(trade exchange.Trade) { slog.Info("trade emitter trade", "trade", trade) },
		OnDisconnect: func() { slog.Info("trade emitter disconnect") },
	})
	if err != nil {
		slog.Error("trade emitter subscribe trades", "error", err)
		return err
	}
	e.unsubscribe = unsubscribe
	return nil
}

func (e *TradeEmitter) Stop() error {
	if e.unsubscribe != nil {
		e.unsubscribe()
	}
	return nil
}
