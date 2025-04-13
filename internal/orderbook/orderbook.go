package orderbook

import (
	"context"
	"errors"
	"log"
	"math"
	"sync"
	"time"

	"github.com/adshao/go-binance/v2"
	"github.com/adshao/go-binance/v2/futures"
	"github.com/shopspring/decimal"
)

const MaxPriceLevels = 65535 // Request weight = 50 for safe request, max: 5000
const MaxSpotLayerRequest = 5000
const MaxPerpLayerRequest = 1000

type UpdateSpeed int

const (
	UpdateSpeed100Ms UpdateSpeed = iota
	UpdateSpeed250Ms
	UpdateSpeed500Ms
	UpdateSpeed1s
)

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
	if oa.BestIndex == math.MaxInt {
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
	if ob.BestIndex == math.MinInt {
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

	NumUpdateCall   int
	NumSnapshotCall int
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

		NumUpdateCall:   0,
		NumSnapshotCall: 0,
	}
}

func (ob *BinanceOrderBook) Summary() {
	log.Printf("NumUpdateCall: %d", ob.NumUpdateCall)
	log.Printf("NumSnapshotCall: %d", ob.NumSnapshotCall)
}

func (ob *BinanceOrderBook) Run(updateSpeed UpdateSpeed) error { // blocking call
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

	doneC := make(chan struct{})
	var err error

	// Start the WebSocket connection
	switch updateSpeed {
	case UpdateSpeed100Ms:
		doneC, ob.stopChan, err = binance.WsDepthServe100Ms(ob.Symbol, wsDepthHandler, errHandler)
	case UpdateSpeed1s:
		doneC, ob.stopChan, err = binance.WsDepthServe(ob.Symbol, wsDepthHandler, errHandler)
	default:
		return errors.New("invalid update speed")
	}

	if err != nil {
		return err
	}

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
				log.Printf("FirstUpdateID: %d, localUpdateID: %d, LastUpdateID: %d", event.FirstUpdateID, ob.lastUpdateID, event.LastUpdateID)
				if event.FirstUpdateID <= ob.lastUpdateID && ob.lastUpdateID <= event.LastUpdateID {
					ob.partialUpdate(event)
					ob.NumUpdateCall++
				} else if event.LastUpdateID < ob.lastUpdateID {
					// outdated event, ignore
					continue
				} else if event.FirstUpdateID > ob.lastUpdateID {
					snapshot, err := client.NewDepthService().Symbol(ob.Symbol).Limit(MaxSpotLayerRequest).Do(context.Background())
					ob.NumSnapshotCall++
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
	ob.lastUpdateID = event.LastUpdateID + 1
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
	ob.lastUpdateID = response.LastUpdateID + 1
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

type BinancePerpOrderBook struct {
	Symbol       string
	Asks         AskBookArray
	Bids         BidBookArray
	timestamp    int64
	lastUpdateID int64
	createdAt    int64
	// eventChan is used to send events to the order book
	eventChan chan *futures.WsDepthEvent
	stopChan  chan struct{}
	doneChan  chan struct{}

	NumUpdateCall   int
	NumSnapshotCall int
}

func NewBinancePerpOrderBook(symbol string, tickSize decimal.Decimal, bufferSize int) *BinancePerpOrderBook {
	return &BinancePerpOrderBook{
		Symbol:       symbol,
		Asks:         *NewAskBookArray(tickSize),
		Bids:         *NewBidBookArray(tickSize),
		timestamp:    0,
		lastUpdateID: 0,
		createdAt:    time.Now().UnixMilli(),

		eventChan: make(chan *futures.WsDepthEvent, bufferSize),
		stopChan:  nil,
		doneChan:  nil,

		NumUpdateCall:   0,
		NumSnapshotCall: 0,
	}
}

func (ob *BinancePerpOrderBook) Summary() {
	log.Printf("NumUpdateCall: %d", ob.NumUpdateCall)
	log.Printf("NumSnapshotCall: %d", ob.NumSnapshotCall)
}

func (ob *BinancePerpOrderBook) Run(updateSpeed UpdateSpeed) error { // blocking call
	wsDepthHandler := func(event *futures.WsDepthEvent) {
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

	doneC := make(chan struct{})
	var err error

	// Start the WebSocket connection
	switch updateSpeed {
	case UpdateSpeed100Ms:
		doneC, ob.stopChan, err = futures.WsDiffDepthServeWithRate(ob.Symbol, 100*time.Millisecond, wsDepthHandler, errHandler)
	case UpdateSpeed250Ms:
		doneC, ob.stopChan, err = futures.WsDiffDepthServeWithRate(ob.Symbol, 250*time.Millisecond, wsDepthHandler, errHandler)
	case UpdateSpeed500Ms:
		doneC, ob.stopChan, err = futures.WsDiffDepthServeWithRate(ob.Symbol, 500*time.Millisecond, wsDepthHandler, errHandler)
	default:
		return errors.New("invalid update speed")
	}

	if err != nil {
		return err
	}

	// Wait for the WebSocket to finish
	go func() {
		<-doneC
		close(doneChan)
	}()

	go func() {
		client := futures.NewClient("", "")
		for {
			select {
			case <-ob.stopChan:
				log.Printf("Stopping order book")
				ob.doneChan <- struct{}{}
				return
			case event := <-ob.eventChan:
				log.Printf("PrevLastUpdateID: %d, FirstUpdateID: %d, localUpdateID: %d, LastUpdateID: %d", event.PrevLastUpdateID, event.FirstUpdateID, ob.lastUpdateID, event.LastUpdateID)
				if event.LastUpdateID < ob.lastUpdateID {
					// 4. Drop any event where u is < lastUpdateId in the snapshot
					continue
				} else if event.PrevLastUpdateID == ob.lastUpdateID {
					// 6. While listening to the stream, each new event's pu should be equal to the previous event's u, otherwise initialize the process from step 3.
					ob.partialUpdate(event)
					ob.NumUpdateCall++
				} else {
					snapshot, err := client.NewDepthService().Symbol(ob.Symbol).Limit(MaxPerpLayerRequest).Do(context.Background())
					ob.NumSnapshotCall++
					if err != nil {
						log.Printf("Error getting snapshot: %+v", err)
						continue
					}
					ob.totalUpdate(snapshot)
					for event := range ob.eventChan {
						if event.FirstUpdateID <= ob.lastUpdateID && ob.lastUpdateID <= event.LastUpdateID {
							// 5.The first processed event should have U <= lastUpdateId**AND**u >= lastUpdateId
							ob.partialUpdate(event)
							break
						} else if event.LastUpdateID < ob.lastUpdateID {
							// 4. Drop any event where u is < lastUpdateId in the snapshot
							continue
						} else {
							log.Printf("Unexpected event state")
							break
						}
					}
				}
			}
		}
	}()

	return nil
}

func (ob *BinancePerpOrderBook) Close() {
	if ob.stopChan != nil {
		close(ob.stopChan)
	}
}

func (ob *BinancePerpOrderBook) GetDepth(depth int) ([]PriceLevel, []PriceLevel, error) {
	if depth <= 0 || depth > MaxPriceLevels {
		return nil, nil, errors.New("depth must be between 1 and 1000")
	}
	asks := ob.Asks.GetBook(depth)
	bids := ob.Bids.GetBook(depth)
	return asks, bids, nil
}

func (ob *BinancePerpOrderBook) partialUpdate(event *futures.WsDepthEvent) {
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

func (ob *BinancePerpOrderBook) totalUpdate(response *futures.DepthResponse) {
	ob.timestamp = response.TradeTime
	ob.lastUpdateID = response.LastUpdateID
	log.Printf("Snapshot LastUpdateID: %d", ob.lastUpdateID)
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
