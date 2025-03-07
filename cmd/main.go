package main

import (
	"context"
	"fmt"
	"math/rand"
	"sync"
	"time"
)

// 1. Market Data Structures
type MarketData struct {
	Symbol    string
	Price     float64
	Volume    float64
	Timestamp time.Time
}

type Candle struct {
	Symbol     string
	Open       float64
	High       float64
	Low        float64
	Close      float64
	Volume     float64
	Resolution string
	StartTime  time.Time
}

// 2. Data Ingestion Component
type MarketDataFetcher interface {
	Fetch() <-chan MarketData
}

type MockMarketDataFetcher struct {
	Symbol string
}

func (m *MockMarketDataFetcher) Fetch() <-chan MarketData {
	out := make(chan MarketData)
	go func() {
		for {
			time.Sleep(500 * time.Millisecond)
			out <- MarketData{
				Symbol:    m.Symbol,
				Price:     rand.Float64() * 100,
				Volume:    rand.Float64() * 1000,
				Timestamp: time.Now(),
			}
		}
	}()
	return out
}

// 3. Data Processing Component
type DataProcessor struct {
	rawData <-chan MarketData
}

func NewDataProcessor(input <-chan MarketData) *DataProcessor {
	return &DataProcessor{rawData: input}
}

func (p *DataProcessor) Process() <-chan Candle {
	out := make(chan Candle)
	go func() {
		defer close(out)
		for data := range p.rawData {
			// Simple processing: create 1-second candles
			// In real implementation, this would aggregate data
			candle := Candle{
				Symbol:     data.Symbol,
				Open:       data.Price,
				High:       data.Price,
				Low:        data.Price,
				Close:      data.Price,
				Volume:     data.Volume,
				Resolution: "1s",
				StartTime:  data.Timestamp,
			}
			out <- candle
		}
	}()
	return out
}

// 4. Trading Strategy Component
type IStrategy interface {
	GenerateSignal(Candle) *Order
}

type MomentumStrategy struct{}

func (s *MomentumStrategy) GenerateSignal(candle Candle) *Order {
	// Simple momentum strategy
	if candle.Close > candle.Open {
		return &Order{
			Symbol:   candle.Symbol,
			Side:     "BUY",
			Quantity: 100,
			Price:    candle.Close,
		}
	}
	return nil
}

// 5. Risk Management Component
type IRiskManager interface {
	ValidateOrder(*Order) bool
}

type BasicRiskManager struct {
	maxPositionSize float64
}

func (r *BasicRiskManager) ValidateOrder(order *Order) bool {
	return order.Quantity <= r.maxPositionSize
}

// 6. Order Management Component
type Order struct {
	Symbol   string
	Side     string
	Quantity float64
	Price    float64
}

type OrderManager struct {
	orderSender IOrderSender
}

type IOrderSender interface {
	SendOrder(*Order) error
}

type MockOrderSender struct{}

func (s *MockOrderSender) SendOrder(order *Order) error {
	fmt.Printf("Order sent: %+v\n", order)
	return nil
}

// TradingPipelineManager manages multiple trading pipelines
type TradingPipelineManager struct {
	pipelines map[string]*TradingPipeline
	mu        sync.RWMutex
	wg        sync.WaitGroup
	ctx       context.Context
	cancel    context.CancelFunc
}

// NewTradingPipelineManager creates a new manager
func NewTradingPipelineManager() *TradingPipelineManager {
	ctx, cancel := context.WithCancel(context.Background())
	return &TradingPipelineManager{
		pipelines: make(map[string]*TradingPipeline),
		ctx:       ctx,
		cancel:    cancel,
	}
}

// PipelineConfig contains configuration for a trading pipeline
type PipelineConfig struct {
	Symbol          string
	Strategy        IStrategy
	RiskManager     IRiskManager
	OrderSender     IOrderSender
	DataFetcher     MarketDataFetcher
	Resolution      string
	MaxPositionSize float64
}

// AddPipeline creates and starts a new trading pipeline
func (m *TradingPipelineManager) AddPipeline(config *PipelineConfig) error {
	m.mu.Lock()
	defer m.mu.Unlock()

	if _, exists := m.pipelines[config.Symbol]; exists {
		return fmt.Errorf("pipeline for symbol %s already exists", config.Symbol)
	}

	// Create pipeline components
	processor := NewDataProcessor(config.DataFetcher.Fetch())

	pipeline := &TradingPipeline{
		fetcher:     config.DataFetcher,
		processor:   processor,
		strategy:    config.Strategy,
		riskManager: config.RiskManager,
		orderMgr: &OrderManager{
			orderSender: config.OrderSender,
		},
		status: PipelineStatusStopped,
		config: config,
	}

	m.pipelines[config.Symbol] = pipeline
	m.startPipeline(pipeline)
	return nil
}

// startPipeline starts an individual pipeline
func (m *TradingPipelineManager) startPipeline(p *TradingPipeline) {
	m.wg.Add(1)
	go func() {
		defer m.wg.Done()
		p.Run(m.ctx)
	}()
}

// StopPipeline stops a specific pipeline
func (m *TradingPipelineManager) StopPipeline(symbol string) error {
	m.mu.Lock()
	defer m.mu.Unlock()

	p, exists := m.pipelines[symbol]
	if !exists {
		return fmt.Errorf("pipeline for symbol %s not found", symbol)
	}

	p.Stop()
	delete(m.pipelines, symbol)
	return nil
}

// Shutdown stops all pipelines gracefully
func (m *TradingPipelineManager) Shutdown() {
	m.cancel()
	m.wg.Wait()
}

// GetPipelineStatus returns the status of all pipelines
func (m *TradingPipelineManager) GetPipelineStatus() map[string]PipelineStatus {
	m.mu.RLock()
	defer m.mu.RUnlock()

	status := make(map[string]PipelineStatus)
	for symbol, p := range m.pipelines {
		status[symbol] = p.status
	}
	return status
}

// Updated TradingPipeline with status management
type TradingPipeline struct {
	fetcher     MarketDataFetcher
	processor   *DataProcessor
	strategy    IStrategy
	riskManager IRiskManager
	orderMgr    *OrderManager
	config      *PipelineConfig
	status      PipelineStatus
}

type PipelineStatus int

const (
	PipelineStatusStopped PipelineStatus = iota
	PipelineStatusRunning
	PipelineStatusError
)

func (p *TradingPipeline) Run(ctx context.Context) {
	p.status = PipelineStatusRunning
	defer func() {
		if r := recover(); r != nil {
			p.status = PipelineStatusError
		}
	}()

	rawData := p.fetcher.Fetch()
	fmt.Printf("Starting pipeline for symbol %+v\n", rawData)
	processedData := p.processor.Process()

	for {
		select {
		case <-ctx.Done():
			p.status = PipelineStatusStopped
			return
		case candle, ok := <-processedData:
			if !ok {
				p.status = PipelineStatusStopped
				return
			}
			order := p.strategy.GenerateSignal(candle)
			if order != nil && p.riskManager.ValidateOrder(order) {
				p.orderMgr.orderSender.SendOrder(order)
			}
		}
	}
}

func (p *TradingPipeline) Stop() {
	// Implement any cleanup logic needed
	p.status = PipelineStatusStopped
}

// Example usage
func main() {
	manager := NewTradingPipelineManager()

	// Configuration for BTC pipeline
	btcConfig := &PipelineConfig{
		Symbol:   "BTC-USD",
		Strategy: &MomentumStrategy{},
		RiskManager: &BasicRiskManager{
			maxPositionSize: 1000,
		},
		OrderSender: &MockOrderSender{},
		DataFetcher: &MockMarketDataFetcher{Symbol: "BTC-USD"},
	}

	// Configuration for ETH pipeline
	ethConfig := &PipelineConfig{
		Symbol:   "ETH-USD",
		Strategy: &MomentumStrategy{},
		RiskManager: &BasicRiskManager{
			maxPositionSize: 500,
		},
		OrderSender: &MockOrderSender{},
		DataFetcher: &MockMarketDataFetcher{Symbol: "ETH-USD"},
	}

	// Add pipelines
	manager.AddPipeline(btcConfig)
	manager.AddPipeline(ethConfig)

	// Start monitoring goroutine
	go func() {
		for {
			status := manager.GetPipelineStatus()
			fmt.Println("Pipeline Status:")
			for symbol, s := range status {
				fmt.Printf("%s: %v\n", symbol, s)
			}
			time.Sleep(2 * time.Second)
		}
	}()

	// Run for 30 seconds
	time.Sleep(30 * time.Second)

	// Graceful shutdown
	manager.Shutdown()
	fmt.Println("All pipelines stopped")
}
