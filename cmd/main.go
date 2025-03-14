package main

import (
	"github.com/BullionBear/sequex/internal/strategy/sequex"
	"github.com/BullionBear/sequex/internal/tradingpipe"
	"github.com/BullionBear/sequex/pkg/mq/inprocq"
)

func main() {
	// Resource
	name := "Trading Pipeline"
	strategy := sequex.NewSequex()
	_ = tradingpipe.NewTradingPipeline(name, strategy)
	_ = inprocq.New(8)
	/*
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
	*/
}
