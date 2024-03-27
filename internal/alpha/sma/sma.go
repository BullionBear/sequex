package sma

import (
	"fmt"
	"sync"

	"github.com/BullionBear/crypto-trade/internal/model"
)

type SMA struct {
	sourceChan chan *model.Tick
	resultChan chan *model.Tick
	period     int64
	window     []float64
	runningSum float64
	wg         sync.WaitGroup
}

func NewSMA(period int64) *SMA {
	return &SMA{
		sourceChan: make(chan *model.Tick),
		resultChan: make(chan *model.Tick),
		period:     period,
		window:     make([]float64, 0, period),
		runningSum: 0,
		wg:         sync.WaitGroup{},
	}
}

func (s *SMA) Name() string {
	return fmt.Sprintf("SMA(%d)", s.period)
}

func (s *SMA) SourceChannel() chan<- *model.Tick {
	return s.sourceChan
}

func (s *SMA) Start() {
	s.wg.Add(1) // Indicate that a goroutine has started
	go func() {
		defer s.wg.Done() // Signal that this goroutine has finished on return
		for tick := range s.sourceChan {
			data := s.process(tick)
			s.resultChan <- data
		}
	}()
}

func (s *SMA) End() {
	close(s.sourceChan) // Signal to stop sending data
	s.wg.Wait()         // Wait for the processing to complete
}

func (s *SMA) OutputChannel() <-chan *model.Tick {
	return s.resultChan
}

func (s *SMA) process(tick *model.Tick) *model.Tick {
	s.runningSum += tick.Price
	s.window = append(s.window, tick.Price)

	// If window size exceeds the period, adjust the running sum and window size
	if int64(len(s.window)) > s.period {
		s.runningSum -= s.window[0]
		s.window = s.window[1:]
	}

	// Compute the SMA and update the tick if we have a full period of data
	isValid := tick.IsValid
	if int64(len(s.window)) == s.period {
		isValid = isValid && true
	} else {
		isValid = false
	}
	return &model.Tick{
		TradeID: tick.TradeID,
		Time:    tick.Time,
		Price:   s.runningSum / float64(s.period),
		IsValid: isValid,
	}
}
