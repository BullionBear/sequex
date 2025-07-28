package trade

import "github.com/nats-io/nats.go"

func NewTradeEmitter(nc *nats.Conn, name string) *TradeEmitter {
	return &TradeEmitter{nc: nc, name: name}
}

type TradeEmitter struct {
	nc   *nats.Conn
	name string
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
