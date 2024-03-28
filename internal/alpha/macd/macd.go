package macd

import (
	"fmt"
	"sync"

	"github.com/BullionBear/crypto-trade/internal/alpha/ema"
	"github.com/BullionBear/crypto-trade/internal/model"
)

type MACD struct {
	sourceChan chan *model.Tick
	resultChan chan *model.MACD
	fastEMA    ema.EMA
	slowEMA    ema.EMA
	signalEMA  ema.EMA
	wg         sync.WaitGroup
}

func NewMACD(fastPeriod int64, slowPeriod int64, signalPeriod int64) *MACD {
	return &MACD{
		sourceChan: make(chan *model.Tick),
		resultChan: make(chan *model.MACD),
		fastEMA:    *ema.NewEMA(fastPeriod),
		slowEMA:    *ema.NewEMA(slowPeriod),
		signalEMA:  *ema.NewEMA(signalPeriod),
		wg:         sync.WaitGroup{},
	}
}

func (m *MACD) Name() string {
	return fmt.Sprintf("MACD(%d, %d, %d)", m.fastEMA.Period, m.slowEMA.Period, m.signalEMA.Period)
}

func (m *MACD) Start() {
	m.wg.Add(1)
	m.fastEMA.Start()
	defer m.fastEMA.End()
	m.slowEMA.Start()
	defer m.slowEMA.End()
	m.signalEMA.Start()
	defer m.signalEMA.End()
	go func() {
		defer m.wg.Done()
		for tick := range m.sourceChan {
			data := m.process(tick)
			m.resultChan <- data
		}
	}()
}

func (m *MACD) End() {
	close(m.sourceChan)
	m.wg.Wait()
}

func (m *MACD) SourceChannel() chan<- *model.Tick {
	return m.sourceChan
}

func (m *MACD) OutputChannel() <-chan *model.MACD {
	return m.resultChan
}

func (m *MACD) process(tick *model.Tick) *model.MACD {
	m.fastEMA.SourceChannel() <- tick
	m.slowEMA.SourceChannel() <- tick
	fastEMA := <-m.fastEMA.OutputChannel()
	slowEMA := <-m.slowEMA.OutputChannel()
	difTick := &model.Tick{
		TradeID: tick.TradeID,
		Time:    tick.Time,
		Price:   fastEMA.Price - slowEMA.Price,
		IsValid: fastEMA.IsValid && slowEMA.IsValid,
	}

	m.signalEMA.SourceChannel() <- difTick
	macdSignal := <-m.signalEMA.OutputChannel()

	return &model.MACD{
		TradeID: tick.TradeID,
		Time:    tick.Time,
		DIF:     difTick.Price,
		DEM:     macdSignal.Price,
		IsValid: difTick.IsValid && macdSignal.IsValid,
	}
}
