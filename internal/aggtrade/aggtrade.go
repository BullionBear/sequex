package aggtrade

import (
	"log"
	"time"

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

func (a *AggTrade) Run(symbol string) error {
	// Implement me
	aggHandler := func(event *binance.WsAggTradeEvent) {
		// Handle the aggregated trade event
		if len(a.eventChan) < cap(a.eventChan) {
			a.eventChan <- event
		} else {
			log.Printf("Event channel is full, dropping event")
		}
	}
	doneC, stopC, err := binance.WsAggTradeServe(symbol, aggHandler, func(err error) {
		log.Printf("error: %v", err)
	})
	ticker := time.NewTicker(1 * time.Second)
	if err != nil {
		return err
	}
	return nil
}
