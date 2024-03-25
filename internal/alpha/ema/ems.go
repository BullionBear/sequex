package ema

import (
	"fmt"
	"sync"

	"github.com/BullionBear/crypto-trade/internal/model"
)

type EMA struct {
	sourceChan chan model.Tick
	resultChan chan model.Tick
	period     int64
	multiplier float64
	ema        float64 // Current EMA value
	wg         sync.WaitGroup
}

func NewEMA(period int64) *EMA {
	multiplier := 2.0 / float64(period+1)
	return &EMA{
		sourceChan: make(chan model.Tick),
		resultChan: make(chan model.Tick),
		period:     period,
		multiplier: multiplier,
		ema:        0, // Initialized to 0, but will be set to the price of the first tick received
		wg:         sync.WaitGroup{},
	}
}

func (e *EMA) Name() string {
	return fmt.Sprintf("EMA(%d)", e.period)
}

func (e *EMA) SourceChannel() chan<- model.Tick {
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

func (e *EMA) OutputChannel() <-chan model.Tick {
	return e.resultChan
}

func (e *EMA) process(tick model.Tick) model.Tick {
	if e.ema == 0 { // If it's the first tick, initialize EMA with its price
		e.ema = tick.Price
	} else {
		e.ema = (tick.Price-e.ema)*e.multiplier + e.ema
	}
	tick.Price = e.ema
	tick.IsValid = true // Assuming the EMA is always valid after the first tick
	return tick
}
