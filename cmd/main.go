package main

import (
	"fmt"
	"math/rand"
	"sync"
	"time"
)

// --- Event Bus Implementation ---

type EventType string

const (
	MarketDataEvent EventType = "MarketDataEvent"
	BuySignalEvent  EventType = "BuySignalEvent"
	SellSignalEvent EventType = "SellSignalEvent"
)

type Event struct {
	Type    EventType
	Payload interface{}
}

type EventBus struct {
	subscribers map[EventType][]chan Event
	mu          sync.RWMutex
}

func NewEventBus() *EventBus {
	return &EventBus{
		subscribers: make(map[EventType][]chan Event),
	}
}

func (eb *EventBus) Subscribe(eventType EventType) <-chan Event {
	eb.mu.Lock()
	defer eb.mu.Unlock()

	ch := make(chan Event)
	eb.subscribers[eventType] = append(eb.subscribers[eventType], ch)
	return ch
}

func (eb *EventBus) Publish(event Event) {
	eb.mu.RLock()
	defer eb.mu.RUnlock()

	if subs, found := eb.subscribers[event.Type]; found {
		for _, sub := range subs {
			go func(s chan Event) {
				s <- event
			}(sub)
		}
	}
}

// --- Domain Models ---

type MarketData struct {
	Time  time.Time
	Price float64
}

type Position struct {
	Asset      string
	Amount     float64
	EntryPrice float64
}

type Portfolio struct {
	Cash      float64
	Positions map[string]Position
}

func NewPortfolio(initialCash float64) *Portfolio {
	return &Portfolio{
		Cash:      initialCash,
		Positions: make(map[string]Position),
	}
}

func (p *Portfolio) Buy(asset string, price, amount float64) {
	cost := price * amount
	if p.Cash >= cost {
		p.Cash -= cost
		p.Positions[asset] = Position{
			Asset:      asset,
			Amount:     amount,
			EntryPrice: price,
		}
		fmt.Printf("Bought %.2f of %s at %.2f\n", amount, asset, price)
	} else {
		fmt.Println("Insufficient cash to buy.")
	}
}

func (p *Portfolio) Sell(asset string, price float64) {
	if position, ok := p.Positions[asset]; ok {
		p.Cash += price * position.Amount
		delete(p.Positions, asset)
		fmt.Printf("Sold %.2f of %s at %.2f\n", position.Amount, asset, price)
	} else {
		fmt.Println("No position to sell.")
	}
}

// --- Trading Pipeline ---

type TradingPipeline struct {
	EventBus   *EventBus
	Portfolio  *Portfolio
	MarketData []MarketData
}

func NewTradingPipeline(initialCash float64) *TradingPipeline {
	return &TradingPipeline{
		EventBus:  NewEventBus(),
		Portfolio: NewPortfolio(initialCash),
	}
}

func (tp *TradingPipeline) fetchMarketData() {
	for i := 0; i < 20; i++ {
		data := MarketData{
			Time:  time.Now().Add(time.Duration(i) * time.Minute),
			Price: 100 + rand.Float64()*10,
		}
		tp.MarketData = append(tp.MarketData, data)
		tp.EventBus.Publish(Event{
			Type:    MarketDataEvent,
			Payload: data,
		})
		time.Sleep(100 * time.Millisecond)
	}
}

func simpleMovingAverage(data []MarketData, period int) float64 {
	if len(data) < period {
		return 0
	}
	sum := 0.0
	for i := len(data) - period; i < len(data); i++ {
		sum += data[i].Price
	}
	return sum / float64(period)
}

func (tp *TradingPipeline) runStrategy() {
	ch := tp.EventBus.Subscribe(MarketDataEvent)
	go func() {
		for event := range ch {
			marketData := event.Payload.(MarketData)
			tp.MarketData = append(tp.MarketData, marketData)

			shortPeriod := 5
			longPeriod := 10

			if len(tp.MarketData) >= longPeriod {
				shortMA := simpleMovingAverage(tp.MarketData, shortPeriod)
				longMA := simpleMovingAverage(tp.MarketData, longPeriod)
				latestPrice := marketData.Price

				fmt.Printf("Short MA: %.2f, Long MA: %.2f, Latest Price: %.2f\n", shortMA, longMA, latestPrice)

				if shortMA > longMA {
					tp.EventBus.Publish(Event{Type: BuySignalEvent, Payload: marketData})
				} else if shortMA < longMA {
					tp.EventBus.Publish(Event{Type: SellSignalEvent, Payload: marketData})
				}
			}
		}
	}()
}

func (tp *TradingPipeline) handleOrders() {
	buyCh := tp.EventBus.Subscribe(BuySignalEvent)
	sellCh := tp.EventBus.Subscribe(SellSignalEvent)

	go func() {
		for {
			select {
			case event := <-buyCh:
				data := event.Payload.(MarketData)
				if len(tp.Portfolio.Positions) == 0 {
					tp.Portfolio.Buy("BTC", data.Price, 1)
				}
			case event := <-sellCh:
				data := event.Payload.(MarketData)
				if len(tp.Portfolio.Positions) > 0 {
					tp.Portfolio.Sell("BTC", data.Price)
				}
			}
		}
	}()
}

func (tp *TradingPipeline) Run() {
	go tp.fetchMarketData()
	tp.runStrategy()
	tp.handleOrders()

	time.Sleep(3 * time.Second) // Wait for all events to process
	fmt.Printf("Final Portfolio: Cash: %.2f, Positions: %+v\n", tp.Portfolio.Cash, tp.Portfolio.Positions)
}

// --- Main Function ---

func main() {
	rand.Seed(time.Now().UnixNano())
	tradingPipeline := NewTradingPipeline(1000)
	tradingPipeline.Run()
}
