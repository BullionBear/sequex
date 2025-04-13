package orderbook

import (
	"context"
	"errors"
	"log"
	"math"
	"sync"
	"time"

	"github.com/adshao/go-binance/v2"
	"github.com/shopspring/decimal"
)

const MaxPriceLevels = 1000

type PriceLevel struct {
	Price decimal.Decimal `json:"price"`
	Size  decimal.Decimal `json:"size"`
}

func NewPriceLevel(price, size decimal.Decimal) PriceLevel {
	return PriceLevel{
		Price: price,
		Size:  size,
	}
}

func (pl *PriceLevel) Empty() {
	pl.Price = decimal.Zero
	pl.Size = decimal.Zero
}

func (pl *PriceLevel) Set(price, size decimal.Decimal) {
	pl.Price = price
	pl.Size = size
}

type AskBookArray struct {
	PriceLevels [MaxPriceLevels]PriceLevel // Static array with a fixed size of 100
	BestIndex   int
	tickSize    decimal.Decimal
}

func NewAskBookArray(tickSize decimal.Decimal) *AskBookArray {
	return &AskBookArray{
		PriceLevels: [MaxPriceLevels]PriceLevel{},
		BestIndex:   math.MaxInt,
		tickSize:    tickSize,
	}
}

func (oa *AskBookArray) GetBestLayer() (PriceLevel, error) {
	if oa.BestIndex >= 0 && oa.BestIndex < MaxPriceLevels {
		return oa.PriceLevels[oa.BestIndex], nil
	}
	return PriceLevel{}, errors.New("best price not available")
}

func (oa *AskBookArray) GetBook(depth int) []PriceLevel {
	if depth <= 0 || depth > MaxPriceLevels {
		return nil
	}

	book := make([]PriceLevel, 0, depth) // Create a slice with capacity but no length
	for i := 0; i < MaxPriceLevels; i++ {
		if !oa.PriceLevels[(oa.BestIndex+i)%MaxPriceLevels].Size.IsZero() {
			book = append(book, oa.PriceLevels[(oa.BestIndex+i)%MaxPriceLevels]) // Use append to add elements
			if len(book) == depth {
				break
			}
		}
	}
	return book
}

func (oa *AskBookArray) UpdateDiff(levels []PriceLevel) {
	for _, level := range levels {
		index := int(level.Price.Div(oa.tickSize).IntPart())
		oa.PriceLevels[index%MaxPriceLevels].Set(level.Price, level.Size)
		oa.BestIndex = min(oa.BestIndex, index)
		if level.Size.IsZero() {
			oa.PriceLevels[index%MaxPriceLevels].Empty()
		} else {
			oa.PriceLevels[index%MaxPriceLevels].Set(level.Price, level.Size)
		}
	}

	for i := 0; i < MaxPriceLevels; i++ {
		if !oa.PriceLevels[(oa.BestIndex+i)%MaxPriceLevels].Size.IsZero() {
			oa.BestIndex = oa.BestIndex + i
			break
		}
	}
}

func (oa *AskBookArray) UpdateAll(levels []PriceLevel) {
	for i := 0; i < MaxPriceLevels; i++ {
		oa.PriceLevels[i].Empty()
	}
	oa.BestIndex = math.MaxInt
	oa.UpdateDiff(levels)
}

type BidBookArray struct {
	PriceLevels [MaxPriceLevels]PriceLevel // Static array with a fixed size of 100
	BestIndex   int
	tickSize    decimal.Decimal
}

func NewBidBookArray(tickSize decimal.Decimal) *BidBookArray {
	return &BidBookArray{
		PriceLevels: [MaxPriceLevels]PriceLevel{},
		BestIndex:   math.MinInt,
		tickSize:    tickSize,
	}
}

func (ob *BidBookArray) GetBestLayer() (PriceLevel, error) {
	if ob.BestIndex >= 0 && ob.BestIndex < MaxPriceLevels {
		return ob.PriceLevels[ob.BestIndex], nil
	}
	return PriceLevel{}, errors.New("best price not available")
}

func (ob *BidBookArray) GetBook(depth int) []PriceLevel {
	if depth <= 0 || depth > MaxPriceLevels {
		return nil
	}

	book := make([]PriceLevel, 0, depth) // Create a slice with capacity but no length
	for i := 0; i < MaxPriceLevels; i++ {
		if !ob.PriceLevels[(ob.BestIndex-i+MaxPriceLevels)%MaxPriceLevels].Size.IsZero() {
			book = append(book, ob.PriceLevels[(ob.BestIndex-i+MaxPriceLevels)%MaxPriceLevels]) // Use append to add elements
			if len(book) == depth {
				break
			}
		}
	}
	return book
}

func (ob *BidBookArray) UpdateDiff(levels []PriceLevel) {
	for _, level := range levels {
		index := int(level.Price.Div(ob.tickSize).IntPart())
		ob.PriceLevels[index%MaxPriceLevels].Set(level.Price, level.Size)
		ob.BestIndex = max(ob.BestIndex, index)
		if level.Size.IsZero() {
			ob.PriceLevels[index%MaxPriceLevels].Empty()
		} else {
			ob.PriceLevels[index%MaxPriceLevels].Set(level.Price, level.Size)
		}
	}

	for i := 0; i < MaxPriceLevels; i++ {
		if !ob.PriceLevels[(ob.BestIndex-i)%MaxPriceLevels].Size.IsZero() {
			ob.BestIndex = ob.BestIndex - i
			break
		}
	}
}

func (ob *BidBookArray) UpdateAll(levels []PriceLevel) {
	for i := 0; i < MaxPriceLevels; i++ {
		ob.PriceLevels[i].Empty()
	}
	ob.BestIndex = math.MinInt
	ob.UpdateDiff(levels)
}

type BinanceOrderBook struct {
	Symbol       string
	Asks         AskBookArray
	Bids         BidBookArray
	timestamp    int64
	lastUpdateID int64
	createdAt    int64
	// eventChan is used to send events to the order book
	eventChan chan *binance.WsDepthEvent
	stopChan  chan struct{}
	doneChan  chan struct{}
}

func NewBinanceOrderBook(symbol string, tickSize decimal.Decimal, bufferSize int) *BinanceOrderBook {
	return &BinanceOrderBook{
		Symbol:       symbol,
		Asks:         *NewAskBookArray(tickSize),
		Bids:         *NewBidBookArray(tickSize),
		timestamp:    0,
		lastUpdateID: 0,
		createdAt:    time.Now().UnixMilli(),

		eventChan: make(chan *binance.WsDepthEvent, bufferSize),
		stopChan:  nil,
	}
}

func (ob *BinanceOrderBook) Run() error { // blocking call
	wsDepthHandler := func(event *binance.WsDepthEvent) {
		if len(ob.eventChan) < cap(ob.eventChan) {
			ob.eventChan <- event
		} else {
			log.Printf("Event channel is full, dropping event")
		}
	}
	errHandler := func(err error) {
		log.Printf("Error in WebSocket: %+v", err)
	}

	// Initialize stopChan and doneChan properly
	stopChan := make(chan struct{})
	doneChan := make(chan struct{})
	ob.stopChan = stopChan
	ob.doneChan = doneChan

	// Start the WebSocket connection
	doneC, stopC, err := binance.WsDepthServe(ob.Symbol, wsDepthHandler, errHandler)
	if err != nil {
		return err
	}

	// Assign the stop channel from the WebSocket to the order book
	ob.stopChan = stopC

	// Wait for the WebSocket to finish
	go func() {
		<-doneC
		close(doneChan)
	}()

	go func() {
		client := binance.NewClient("", "")
		for {
			select {
			case <-ob.stopChan:
				log.Printf("Stopping order book")
				ob.doneChan <- struct{}{}
				return
			case event := <-ob.eventChan:
				if event.FirstUpdateID <= ob.lastUpdateID && ob.lastUpdateID <= event.LastUpdateID {
					ob.partialUpdate(event)
				} else if event.LastUpdateID < ob.lastUpdateID {
					// outdated event, ignore
					continue
				} else if event.FirstUpdateID > ob.lastUpdateID {
					snapshot, err := client.NewDepthService().Symbol(ob.Symbol).Limit(1000).Do(context.Background())
					if err != nil {
						log.Printf("Error getting snapshot: %+v", err)
						continue
					}
					ob.totalUpdate(snapshot)
					if event.FirstUpdateID <= ob.lastUpdateID && ob.lastUpdateID <= event.LastUpdateID {
						ob.partialUpdate(event)
					}
				} else {
					log.Printf("Unexpected event state")
				}
			}
		}
	}()

	return nil
}

func (ob *BinanceOrderBook) Close() {
	if ob.stopChan != nil {
		close(ob.stopChan)
	}
}

func (ob *BinanceOrderBook) GetDepth(depth int) ([]PriceLevel, []PriceLevel, error) {
	if depth <= 0 || depth > MaxPriceLevels {
		return nil, nil, errors.New("depth must be between 1 and 1000")
	}
	asks := ob.Asks.GetBook(depth)
	bids := ob.Bids.GetBook(depth)
	return asks, bids, nil
}

func (ob *BinanceOrderBook) partialUpdate(event *binance.WsDepthEvent) {
	ob.timestamp = event.Time
	ob.lastUpdateID = event.LastUpdateID
	wg := sync.WaitGroup{}
	wg.Add(2)
	go func() {
		pxlv := make([]PriceLevel, len(event.Asks))
		for i := 0; i < len(event.Asks); i++ {
			pxlv[i] = NewPriceLevel(
				decimal.RequireFromString(event.Asks[i].Price),
				decimal.RequireFromString(event.Asks[i].Quantity),
			)
		}
		ob.Asks.UpdateDiff(pxlv)
		wg.Done()
	}()
	go func() {
		pxlv := make([]PriceLevel, len(event.Bids))
		for i := 0; i < len(event.Bids); i++ {
			pxlv[i] = NewPriceLevel(
				decimal.RequireFromString(event.Bids[i].Price),
				decimal.RequireFromString(event.Bids[i].Quantity),
			)
		}
		ob.Bids.UpdateDiff(pxlv)
		wg.Done()
	}()
	wg.Wait()
}

func (ob *BinanceOrderBook) totalUpdate(response *binance.DepthResponse) {
	ob.timestamp = time.Now().UnixMilli()
	ob.lastUpdateID = response.LastUpdateID
	wg := sync.WaitGroup{}
	wg.Add(2)
	go func() {
		pxlv := make([]PriceLevel, len(response.Asks))
		for i := 0; i < len(response.Asks); i++ {
			pxlv[i] = NewPriceLevel(
				decimal.RequireFromString(response.Asks[i].Price),
				decimal.RequireFromString(response.Asks[i].Quantity),
			)
		}
		ob.Asks.UpdateAll(pxlv)
		wg.Done()
	}()
	go func() {
		pxlv := make([]PriceLevel, len(response.Bids))
		for i := 0; i < len(response.Bids); i++ {
			pxlv[i] = NewPriceLevel(
				decimal.RequireFromString(response.Bids[i].Price),
				decimal.RequireFromString(response.Bids[i].Quantity),
			)
		}
		ob.Bids.UpdateAll(pxlv)
		wg.Done()
	}()
	wg.Wait()
}
