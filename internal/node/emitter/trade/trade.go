package trade

import "github.com/BullionBear/sequex/internal/exchange"

func NewTradeEmitter(name string, connector exchange.IsolatedSpotConnector) *TradeEmitter {
	return &TradeEmitter{name: name, connector: connector}
}

type TradeEmitter struct {
	name      string
	connector exchange.IsolatedSpotConnector
}

func (e *TradeEmitter) Type() string {
	return "TradeEmitter"
}

func (e *TradeEmitter) Name() string {
	return e.name
}

func (e *TradeEmitter) Start() error {

	return nil
}

func (e *TradeEmitter) Stop() error {
	return nil
}
