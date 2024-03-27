package sma

import (
	"testing"
	"time"

	"github.com/BullionBear/crypto-trade/internal/model"
)

func TestSMAChannelCommunication(t *testing.T) {
	period := int64(3)
	sma := NewSMA(period)

	// Test data
	ticksToSend := []*model.Tick{
		{TradeID: 1, Time: 100, Price: 10},
		{TradeID: 2, Time: 200, Price: 20},
		{TradeID: 3, Time: 300, Price: 30},
		{TradeID: 4, Time: 400, Price: 40},
	}

	// the SMA for the last tick (price 40) should consider the prices [20, 30, 40].
	expectedLastPrice := (20.0 + 30.0 + 40.0) / 3.0

	sma.Start()
	defer sma.End()

	go func() {
		for _, tick := range ticksToSend {
			sma.SourceChannel() <- tick
		}
	}()

	// Collect the processed ticks
	var receivedTicks []*model.Tick
	timeout := time.After(2 * time.Second) // Set a timeout for safety

	keepCollecting := true
	for keepCollecting {
		select {
		case tick := <-sma.OutputChannel():
			receivedTicks = append(receivedTicks, tick)
			if len(receivedTicks) == len(ticksToSend) {
				keepCollecting = false
			}
		case <-timeout:
			t.Fatal("Timeout waiting for ticks to be processed")
			keepCollecting = false
		}
	}

	// Verify the number of received ticks matches the number sent
	if len(receivedTicks) != len(ticksToSend) {
		t.Errorf("Expected to receive %d ticks, got %d", len(ticksToSend), len(receivedTicks))
	}

	// Verify the SMA calculation of the last tick
	lastTick := receivedTicks[len(receivedTicks)-1]
	if !lastTick.IsValid || lastTick.Price != expectedLastPrice {
		t.Errorf("Expected the last tick to have price %v, got %v", expectedLastPrice, lastTick.Price)
	}
}
