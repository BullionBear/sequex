package aggtrade

import (
	"github.com/adshao/go-binance/v2"
	influxdb2 "github.com/influxdata/influxdb-client-go/v2"
)

type AggTrade struct {
	client    influxdb2.Client
	eventChan chan *binance.WsAggTradeEvent
}

func NewAggTrade(client influxdb2.Client, limit int) *AggTrade {
	return &AggTrade{
		client:    client,
		eventChan: make(chan *binance.WsAggTradeEvent, limit),
	}
}
