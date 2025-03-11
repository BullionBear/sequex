package main

import (
	"fmt"

	"github.com/BullionBear/sequex/internal/metadata"
	"github.com/BullionBear/sequex/internal/tradingpipe"
	"github.com/BullionBear/sequex/pkg/mq/inprocq"
)

func main() {
	// Resource
	tradingPipeline := tradingpipe.NewTradingPipeline("Trading Pipeline")
	eventQ := inprocq.New(8)
	ch, err := eventQ.Subscribe("event")
	if err != nil {
		panic(err)
	}
	for event := range ch {
		switch event.Type {
		case "kline_update":
			metadata := event.Data.(metadata.KLineUpdate)
			tradingPipeline.OnKLineUpdate(metadata)
		default:
			fmt.Println("Unknown event type")
		}
	}
	done := make(chan struct{})
	<-done
}
