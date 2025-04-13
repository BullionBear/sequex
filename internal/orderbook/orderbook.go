package orderbook

import (
	"context"
	"errors"
	"log"
	"sync"
	"time"

	"github.com/adshao/go-binance/v2"
	"github.com/adshao/go-binance/v2/futures"
	"github.com/emirpasic/gods/maps/treemap"
	"github.com/shopspring/decimal"
)

const MaxSpotLayerRequest = 5000
const MaxPerpLayerRequest = 1000

type UpdateSpeed int

const (
	UpdateSpeed100Ms UpdateSpeed = iota
	UpdateSpeed250Ms
	UpdateSpeed500Ms
	UpdateSpeed1s
)

func decimalComparator(a, b interface{}) int {
	d1 := a.(decimal.Decimal)
	d2 := b.(decimal.Decimal)
	return d1.Cmp(d2)
}

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
	PriceLevels treemap.Map
}

func NewAskBookArray() *AskBookArray {
	return &AskBookArray{
		PriceLevels: *treemap.NewWith(decimalComparator),
	}
}

func (oa *AskBookArray) GetBestLayer() (PriceLevel, error) {
	if oa.PriceLevels.Empty() {
		return PriceLevel{}, errors.New("best price not available")
	}
	bestPrice, bestSize := oa.PriceLevels.Min()
	return NewPriceLevel(bestPrice.(decimal.Decimal), bestSize.(decimal.Decimal)), nil
}

func (oa *AskBookArray) GetBook(depth int) []PriceLevel {
	book := make([]PriceLevel, 0, depth) // Create a slice with capacity but no length
	it := oa.PriceLevels.Iterator()
	count := 0
	for it.Next() {
		book = append(book, NewPriceLevel(it.Key().(decimal.Decimal), it.Value().(decimal.Decimal))) // Use append to add elements
		count++
		if count >= depth {
			break
		}
	}
	return book
}

func (oa *AskBookArray) UpdateDiff(levels []PriceLevel) {
	for _, level := range levels {
		if level.Size.IsZero() {
			oa.PriceLevels.Remove(level.Price)
		} else {
			oa.PriceLevels.Put(level.Price, level.Size)
		}
	}
}

func (oa *AskBookArray) UpdateAll(levels []PriceLevel) {
	oa.PriceLevels.Clear()
	for _, level := range levels {
		if level.Size.IsZero() {
			continue
		}
		oa.PriceLevels.Put(level.Price, level.Size)
	}
}

type BidBookArray struct {
	PriceLevels treemap.Map
	BestIndex   int
	tickSize    decimal.Decimal
}

func NewBidBookArray() *BidBookArray {
	return &BidBookArray{
		PriceLevels: *treemap.NewWith(decimalComparator),
	}
}

func (ob *BidBookArray) GetBestLayer() (PriceLevel, error) {
	if ob.PriceLevels.Empty() {
		return PriceLevel{}, errors.New("best price not available")
	}
	bestPrice, bestSize := ob.PriceLevels.Max()
	return NewPriceLevel(bestPrice.(decimal.Decimal), bestSize.(decimal.Decimal)), nil
}

func (ob *BidBookArray) GetBook(depth int) []PriceLevel {
	book := make([]PriceLevel, 0, depth) // Create a slice with capacity but no length
	it := ob.PriceLevels.Iterator()
	count := 0
	for it.End(); it.Prev(); {
		book = append(book, NewPriceLevel(it.Key().(decimal.Decimal), it.Value().(decimal.Decimal))) // Use append to add elements
		count++
		if count >= depth {
			break
		}
	}
	return book
}

func (ob *BidBookArray) UpdateDiff(levels []PriceLevel) {
	for _, level := range levels {
		if level.Size.IsZero() {
			ob.PriceLevels.Remove(level.Price)
		} else {
			ob.PriceLevels.Put(level.Price, level.Size)
		}
	}
}

func (ob *BidBookArray) UpdateAll(levels []PriceLevel) {
	ob.PriceLevels.Clear()
	for _, level := range levels {
		if level.Size.IsZero() {
			continue
		}
		ob.PriceLevels.Put(level.Price, level.Size)
	}
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

func NewBinanceOrderBook(symbol string, bufferSize int) *BinanceOrderBook {
	return &BinanceOrderBook{
		Symbol:       symbol,
		Asks:         *NewAskBookArray(),
		Bids:         *NewBidBookArray(),
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

func NewBinancePerpOrderBook(symbol string, bufferSize int) *BinancePerpOrderBook {
	return &BinancePerpOrderBook{
		Symbol:       symbol,
		Asks:         *NewAskBookArray(),
		Bids:         *NewBidBookArray(),
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
