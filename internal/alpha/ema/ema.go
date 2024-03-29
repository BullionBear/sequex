package ema

import (
	"fmt"
	"sync"

	"github.com/BullionBear/crypto-trade/internal/model"
)

type EMA struct {
	sourceChan chan *model.Tick
	resultChan chan *model.Tick
	Period     int64
	multiplier float64
	ema        float64
	ticks      []float64
	tickCount  int64
	wg         sync.WaitGroup
}

func NewEMA(period int64) *EMA {
	multiplier := 2.0 / float64(period+1)
	return &EMA{
		sourceChan: make(chan *model.Tick),
		resultChan: make(chan *model.Tick),
		Period:     period,
		multiplier: multiplier,
		ema:        0,
		ticks:      make([]float64, 0, period),
		tickCount:  0,
		wg:         sync.WaitGroup{},
	}
}

func (e *EMA) Name() string {
	return fmt.Sprintf("EMA(%d)", e.Period)
}

func (e *EMA) SourceChannel() chan<- *model.Tick {
	return e.sourceChan
}

func (e *EMA) Start() {
	e.wg.Add(1)
	go func() {
		defer e.wg.Done()
		for tick := range e.sourceChan {
			data := e.process(tick)
			e.resultChan <- data
		}
	}()
}

func (e *EMA) End() {
	close(e.sourceChan)
	e.wg.Wait()
}

func (e *EMA) OutputChannel() <-chan *model.Tick {
	return e.resultChan
}

func (e *EMA) process(tick *model.Tick) *model.Tick {
	// Only process valid ticks.
	if !tick.IsValid {
		e.ema = 0
		e.ticks = e.ticks[0:]
		e.tickCount = 0
		return &model.Tick{
			TradeID: tick.TradeID,
			Time:    tick.Time,
			Price:   0,
			IsValid: false,
		}
	}

	// Increment tick count for every valid tick.
	e.tickCount++

	// Calculate the initial EMA using simple moving average (SMA) approach for the first 'Period' ticks.
	if e.tickCount <= e.Period {
		// Aggregate ticks for initial SMA calculation.
		e.ticks = append(e.ticks, tick.Price)
		sum := 0.0
		for _, price := range e.ticks {
			sum += price
		}
		e.ema = sum / float64(len(e.ticks))

		// Return the tick with modified price (EMA) and IsValid flag.
		return &model.Tick{
			TradeID: tick.TradeID,
			Time:    tick.Time,
			Price:   e.ema,
			IsValid: e.tickCount == e.Period, // Only valid after 'Period' ticks.
		}
	}

	// For subsequent ticks, use the EMA formula.
	e.ema = ((tick.Price - e.ema) * e.multiplier) + e.ema

	// Return the processed tick.
	return &model.Tick{
		TradeID: tick.TradeID,
		Time:    tick.Time,
		Price:   e.ema,
		IsValid: true, // All subsequent ticks after the initial period are valid.
	}
}
