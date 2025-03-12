package main

import (
	"fmt"
	"time"

	"github.com/BullionBear/sequex/internal/metadata"
	"github.com/BullionBear/sequex/internal/strategy/sequex"
	"github.com/BullionBear/sequex/internal/tradingpipe"
	"github.com/BullionBear/sequex/pkg/message"
	"github.com/BullionBear/sequex/pkg/mq/inprocq"
	"github.com/google/uuid"
)

func main() {
	// Resource
	name := "Trading Pipeline"
	strategy := sequex.NewSequex()
	tradingPipeline := tradingpipe.NewTradingPipeline(name, strategy)
	eventQ := inprocq.New(8)
	ch, err := eventQ.Subscribe(name)
	if err != nil {
		panic(err)
	}
	go func() {
		ticker := time.NewTicker(time.Second)
		defer ticker.Stop()
		for range ticker.C {
			msg := message.Message{
				ID:        uuid.New().String(),
				Type:      "kline_update",
				CreatedAt: time.Now().UnixMilli(),
				Data:      nil,
				Metadata: metadata.KLineUpdate{
					Symbol:    "BTCUSDT",
					Interval:  "1m",
					Timestamp: time.Now().UnixMilli(),
				},
			}
			eventQ.Publish(name, msg)
		}
	}()
	for event := range ch {
		switch event.Type {
		case "kline_update":
			metadata := event.Metadata.(metadata.KLineUpdate)
			tradingPipeline.OnKLineUpdate(metadata)
		default:
			fmt.Println("Unknown event type")
		}
	}
	done := make(chan struct{})
	<-done
}
