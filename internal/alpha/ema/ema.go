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
	e.tickCount++
	e.ticks = append(e.ticks, tick.Price)
	IsValid := true
	if e.tickCount <= e.Period {
		sum := 0.0
		for _, price := range e.ticks {
			sum += price
		}
		e.ema = sum / float64(len(e.ticks))
		IsValid = false
	} else {
		if e.ema == 0 {
			e.ema = tick.Price
		} else {
			e.ema = (tick.Price-e.ema)*e.multiplier + e.ema
		}
	}

	if e.tickCount >= 1024 { // reset tickCount
		e.tickCount = e.Period
	}
	return &model.Tick{
		TradeID: tick.TradeID,
		Time:    tick.Time,
		Price:   e.ema,
		IsValid: tick.IsValid && IsValid,
	}
}
