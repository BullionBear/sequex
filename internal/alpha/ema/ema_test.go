package ema

import (
	"testing"
	"time"

	"github.com/BullionBear/crypto-trade/internal/model"
)

type expectedTick struct {
	TradeID int64
	Time    int64
	Price   float64
	IsValid bool
}

func TestEMAChannelCommunication(t *testing.T) {
	cases := []struct {
		name          string
		period        int64
		ticksToSend   []*model.Tick
		expectedTicks []expectedTick
	}{
		{
			name:   "Basic EMA with multiple ticks",
			period: 3,
			ticksToSend: []*model.Tick{
				{TradeID: 1, Time: 100, Price: 10, IsValid: true},
				{TradeID: 2, Time: 200, Price: 20, IsValid: true},
				{TradeID: 3, Time: 300, Price: 30, IsValid: true},
				{TradeID: 4, Time: 400, Price: 40, IsValid: true},
				// Add more ticks as needed
			},
			expectedTicks: []expectedTick{
				{TradeID: 1, Time: 100, Price: 10, IsValid: false},
				{TradeID: 2, Time: 200, Price: 15, IsValid: false},
				{TradeID: 3, Time: 300, Price: 20, IsValid: true},
				{TradeID: 4, Time: 400, Price: 30, IsValid: true},
			},
		},
		// Add more test cases as needed
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			ema := NewEMA(tc.period)
			ema.Start()
			defer ema.End()

			go func() {
				for _, tick := range tc.ticksToSend {
					ema.SourceChannel() <- tick
				}
			}()

			var receivedTicks []*model.Tick
			timeout := time.After(2 * time.Second)

			keepCollecting := true
			for keepCollecting {
				select {
				case tick := <-ema.OutputChannel():
					receivedTicks = append(receivedTicks, tick)
					if len(receivedTicks) == len(tc.ticksToSend) {
						keepCollecting = false
					}
				case <-timeout:
					t.Fatal("Timeout waiting for ticks to be processed")
					keepCollecting = false
				}
			}

			if len(receivedTicks) != len(tc.expectedTicks) {
				t.Fatalf("Expected to receive %d EMAs, got %d", len(tc.expectedTicks), len(receivedTicks))
			}

			// Verify each tick against expected outcomes
			for i, expected := range tc.expectedTicks {
				actual := receivedTicks[i]
				if actual.TradeID != expected.TradeID || actual.Time != expected.Time || actual.Price != expected.Price || actual.IsValid != expected.IsValid {
					t.Errorf("Tick %d mismatched. Expected %+v, got %+v", i, expected, actual)
				}
			}
		})
	}
}
