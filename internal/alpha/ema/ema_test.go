package ema

import (
	"testing"
	"time"

	"github.com/BullionBear/crypto-trade/internal/model"
)

func TestEMAChannelCommunication(t *testing.T) {
	period := int64(3)
	ema := NewEMA(period) // Assuming a constructor similar to NewSMA

	ticksToSend := []*model.Tick{
		{Price: 1},
		{Price: 2},
		{Price: 3},
		{Price: 4},
		// Add more ticks as needed for testing
	}

	expectedEMA := 3.0 // Placeholder value, calculate based on your EMA formula

	go ema.Start()
	defer ema.End()

	go func() {
		for _, tick := range ticksToSend {
			ema.SourceChannel() <- tick
		}
	}()

	var receivedTicks []model.Tick
	timeout := time.After(2 * time.Second)

collectLoop:
	for {
		select {
		case tick := <-ema.OutputChannel():
			receivedTicks = append(receivedTicks, *tick)
			if len(receivedTicks) == len(ticksToSend) {
				break collectLoop
			}
		case <-timeout:
			t.Fatal("Timeout waiting for ticks to be processed")
		}
	}

	// Verify the EMA calculation of the last tick
	lastTick := receivedTicks[len(receivedTicks)-1]
	if lastTick.Price != expectedEMA {
		t.Errorf("Expected the last tick to have EMA %v, got %v", expectedEMA, lastTick.Price)
	}
}
