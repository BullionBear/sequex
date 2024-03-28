package macd

import (
	"testing"

	"github.com/BullionBear/crypto-trade/internal/model"
)

func TestEMAChannelCommunication(t *testing.T) {
	macd := NewMACD(2, 3, 2)

	ticksToSend := []*model.Tick{
		{Price: 1},
		{Price: 2},
		{Price: 3},
		{Price: 4},
		// Add more ticks as needed for testing
	}

	go macd.Start()
	defer macd.End()

	go func() {
		for _, tick := range ticksToSend {
			macd.SourceChannel() <- tick
		}
	}()

}
