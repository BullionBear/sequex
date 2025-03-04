package main

import (
	"fmt"
	"math/rand"
	"time"
)

// Event Types
const (
	DataEvent      = "data"
	SignalEvent    = "signal"
	ExecutionEvent = "execution"
	LogEvent       = "log"
)

// Event Struct
type Event struct {
	Type    string
	Payload interface{}
}

// Event Bus
type EventBus struct {
	channels map[string]chan Event
}

func NewEventBus() *EventBus {
	return &EventBus{
		channels: make(map[string]chan Event),
	}
}

func (eb *EventBus) Subscribe(eventType string) chan Event {
	ch := make(chan Event, 10) // Buffered channel to handle bursts
	eb.channels[eventType] = ch
	return ch
}

func (eb *EventBus) Publish(event Event) {
	if ch, ok := eb.channels[event.Type]; ok {
		ch <- event
	}
}

// Data Collector
type DataCollector struct {
	bus *EventBus
}

func (dc *DataCollector) Run() {
	for {
		data := map[string]float64{
			"timestamp": float64(time.Now().Unix()),
			"open":      rand.Float64()*5000 + 30000,
			"high":      rand.Float64()*5000 + 35000,
			"low":       rand.Float64()*5000 + 30000,
			"close":     rand.Float64()*5000 + 30000,
			"volume":    rand.Float64()*400 + 100,
		}
		fmt.Printf("[DataCollector] Collected data: %v\n", data)
		dc.bus.Publish(Event{Type: DataEvent, Payload: data})
		time.Sleep(5 * time.Second)
	}
}

// Strategy
type Strategy struct {
	bus         *EventBus
	dataChannel chan Event
}

func (s *Strategy) Run() {
	for event := range s.dataChannel {
		_ = event.Payload.(map[string]float64)
		signals := []string{"buy", "sell", "hold"}
		signal := signals[rand.Intn(len(signals))]
		fmt.Printf("[Strategy] Generated signal: %s\n", signal)
		s.bus.Publish(Event{Type: SignalEvent, Payload: signal})
	}
}

// Trade Executor
type TradeExecutor struct {
	bus           *EventBus
	signalChannel chan Event
}

func (te *TradeExecutor) Run() {
	for event := range te.signalChannel {
		signal := event.Payload.(string)
		if signal != "hold" {
			fmt.Printf("[TradeExecutor] Executing %s order.\n", signal)
			te.bus.Publish(Event{Type: ExecutionEvent, Payload: fmt.Sprintf("%s order executed", signal)})
		} else {
			fmt.Println("[TradeExecutor] No action taken.")
		}
	}
}

// Logger
type Logger struct {
	logChannel chan Event
}

func (l *Logger) Run() {
	for event := range l.logChannel {
		fmt.Printf("[Logger] %s\n", event.Payload)
	}
}

// Pipeline Coordinator
type TradingPipeline struct {
	bus           *EventBus
	dataCollector *DataCollector
	strategy      *Strategy
	executor      *TradeExecutor
	logger        *Logger
}

func NewTradingPipeline() *TradingPipeline {
	bus := NewEventBus()

	pipeline := &TradingPipeline{
		bus:           bus,
		dataCollector: &DataCollector{bus: bus},
		strategy: &Strategy{
			bus:         bus,
			dataChannel: bus.Subscribe(DataEvent),
		},
		executor: &TradeExecutor{
			bus:           bus,
			signalChannel: bus.Subscribe(SignalEvent),
		},
		logger: &Logger{
			logChannel: bus.Subscribe(ExecutionEvent),
		},
	}

	return pipeline
}

func (tp *TradingPipeline) Run() {
	go tp.dataCollector.Run()
	go tp.strategy.Run()
	go tp.executor.Run()
	go tp.logger.Run()

	select {} // Keep main alive
}

func main() {
	rand.Seed(time.Now().UnixNano())
	pipeline := NewTradingPipeline()
	fmt.Println("Starting event-driven trading pipeline. Press Ctrl+C to stop.")
	pipeline.Run()
}
