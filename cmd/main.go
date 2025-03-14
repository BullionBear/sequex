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
	pipeline := tradingpipe.NewTradingPipeline(name, strategy)
	q := inprocq.NewInprocQueue()

	go func() {
		for {
			msg := &message.Message{
				ID:   uuid.New().String(),
				Type: "kline_update",
				Metadata: metadata.KLineUpdate{
					Symbol:    "BTCUSDT",
					Interval:  "1m",
					Timestamp: time.Now().Unix(),
				},
			}
			q.Send(msg)
			time.Sleep(1 * time.Second)
		}
	}()

	go func() {
		for {
			msg, err := q.RecvTimeout(2 * time.Second)
			if err != nil {
				fmt.Printf("Error receiving message: %v\n", err)
			} else {
				fmt.Printf("Received message: %v\n", msg)
				switch msg.Type {
				case "kline_update":
					metadata := msg.Metadata.(metadata.KLineUpdate)
					pipeline.OnKLineUpdate(metadata)

				default:
					fmt.Printf("Unknown message type: %v\n", msg.Type)
				}
			}
		}
	}()

	done := make(chan struct{})
	<-done

}
