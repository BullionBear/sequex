package orderbook

import (
	"context"
	"errors"
	"fmt"

	"sync"
	"sync/atomic"
	"time"

	"github.com/BullionBear/sequex/pkg/log"
	"github.com/adshao/go-binance/v2"
	"github.com/adshao/go-binance/v2/futures"
	"github.com/emirpasic/gods/maps/treemap"
	"github.com/google/uuid"
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

type BookArray struct {
	PriceLevels treemap.Map
}

func NewBookArray() *BookArray {
	return &BookArray{
		PriceLevels: *treemap.NewWith(decimalComparator),
	}
}

func (ba *BookArray) GetBestLayer(isAsk bool) (PriceLevel, error) {
	if ba.PriceLevels.Empty() {
		return PriceLevel{}, errors.New("best price not available")
	}
	if isAsk {
		bestPrice, bestSize := ba.PriceLevels.Min()
		return NewPriceLevel(bestPrice.(decimal.Decimal), bestSize.(decimal.Decimal)), nil
	} else {
		bestPrice, bestSize := ba.PriceLevels.Max()
		return NewPriceLevel(bestPrice.(decimal.Decimal), bestSize.(decimal.Decimal)), nil
	}
}

func (ba *BookArray) GetBook(depth int, isAsk bool) []PriceLevel {
	book := make([]PriceLevel, 0, depth)
	it := ba.PriceLevels.Iterator()
	count := 0
	if isAsk {
		for it.Next() {
			book = append(book, NewPriceLevel(it.Key().(decimal.Decimal), it.Value().(decimal.Decimal)))
			count++
			if count >= depth {
				break
			}
		}
	} else {
		for it.End(); it.Prev(); {
			book = append(book, NewPriceLevel(it.Key().(decimal.Decimal), it.Value().(decimal.Decimal)))
			count++
			if count >= depth {
				break
			}
		}
	}
	return book
}

func (ba *BookArray) UpdateDiff(levels []PriceLevel) {
	for _, level := range levels {
		if level.Size.IsZero() {
			ba.PriceLevels.Remove(level.Price)
		} else {
			ba.PriceLevels.Put(level.Price, level.Size)
		}
	}
}

func (ba *BookArray) UpdateAll(levels []PriceLevel) {
	ba.PriceLevels.Clear()
	for _, level := range levels {
		if level.Size.IsZero() {
			continue
		}
		ba.PriceLevels.Put(level.Price, level.Size)
	}
}

type AskBookArray struct {
	*BookArray
}

func NewAskBookArray() *AskBookArray {
	return &AskBookArray{
		BookArray: NewBookArray(),
	}
}

type BidBookArray struct {
	*BookArray
}

func NewBidBookArray() *BidBookArray {
	return &BidBookArray{
		BookArray: NewBookArray(),
	}
}

// BinanceOrderBook implements io.Closer so callers can defer ob.Close().
type BinanceOrderBook struct {
	/* ======= public, read‑only fields ======= */
	Symbol string
	Asks   AskBookArray
	Bids   BidBookArray

	/* ======= streaming internals ======= */
	timestamp    int64
	lastUpdateID int64
	createdAt    int64

	/* ======= streaming internals ======= */
	eventCh chan *binance.WsDepthEvent // buffered, never nil
	stopC   chan struct{}              // <- signal to underlying WS service
	doneC   chan struct{}              // <- closed by WS service on exit

	subscriber sync.Map // map[depth int]func(ask, bid []PriceLevel)

	/* ======= coordination ======= */
	ctx    context.Context    // global cancel point for all goroutines
	cancel context.CancelFunc // paired with ctx
	wg     sync.WaitGroup     // waits for internal goroutines

	/* ======= optional metrics ======= */
	logger          *log.Logger // optional logger
	numUpdateCall   int64       // accessed atomically
	numSnapshotCall int64       // accessed atomically
}

/*
NewBinanceOrderBook allocates every resource the instance owns.
Nothing is started yet, so it’s safe to create many instances cheaply.
*/
func NewBinanceOrderBook(symbol string, bufferSize int, logger *log.Logger) *BinanceOrderBook {
	ctx, cancel := context.WithCancel(context.Background())

	return &BinanceOrderBook{
		Symbol:    symbol,
		Asks:      *NewAskBookArray(),
		Bids:      *NewBidBookArray(),
		createdAt: time.Now().UnixMilli(),
		eventCh:   make(chan *binance.WsDepthEvent, bufferSize),
		ctx:       ctx,
		cancel:    cancel,
		logger:    logger,
	}
}

/*
Start dials the Binance stream and launches exactly one listener goroutine.
It can be called once; subsequent calls return an error.
*/
func (ob *BinanceOrderBook) Start(spd UpdateSpeed) error {
	if ob.stopC != nil {
		return errors.New("orderbook already started")
	}

	// ----- 1. wire handlers -----
	wsDepthHandler := func(ev *binance.WsDepthEvent) {
		select {
		case ob.eventCh <- ev:
		default:
			ob.logger.Warn("depth‑event dropped (buffer full)")
		}
	}
	errHandler := func(err error) { ob.logger.Error("WS error: %s", err.Error()) }

	// ----- 2. start WS -----
	var err error
	switch spd {
	case UpdateSpeed100Ms:
		ob.doneC, ob.stopC, err = binance.WsDepthServe100Ms(ob.Symbol, wsDepthHandler, errHandler)
	case UpdateSpeed1s:
		ob.doneC, ob.stopC, err = binance.WsDepthServe(ob.Symbol, wsDepthHandler, errHandler)
	default:
		return fmt.Errorf("unknown update speed %v", spd)
	}
	if err != nil {
		return err
	}

	// ----- 3. launch listener -----
	ob.wg.Add(1)
	go ob.listen()

	return nil
}

/*
listen consumes the depth events and terminates when:
  - Binance closes the socket (doneC closed), or
  - The caller invokes Close() (ctx cancelled).

It is the only place that reads from eventCh, so when it returns we can
safely close(eventCh) to free memory.
*/
func (ob *BinanceOrderBook) listen() {
	defer ob.wg.Done()
	defer close(ob.eventCh)

	client := binance.NewClient("", "")

	for {
		select {
		case <-ob.ctx.Done():
			return // caller asked us to quit
		case <-ob.doneC: // WS layer finished
			return
		case ev := <-ob.eventCh: // main path
			ob.handleDepthEvent(client, ev)
			bestAskLayer, err := ob.Asks.GetBestLayer(true)
			// get best ask layer
			if err != nil {
				continue
			}
			bestBidLayer, err := ob.Bids.GetBestLayer(false)
			if err != nil {
				continue
			}
			ob.Bids.GetBestLayer(false)
			ob.publishBestDepth(bestAskLayer, bestBidLayer)
		}
	}
}

func (ob *BinanceOrderBook) handleDepthEvent(cl *binance.Client, ev *binance.WsDepthEvent) {
	ob.logger.Info("FUID:%d  localUID:%d  LUID:%d", ev.FirstUpdateID, ob.lastUpdateID, ev.LastUpdateID)

	switch {
	// ──────────────────────────────────────────────────────────────
	// 1. Normal in‑sequence diff
	case ev.FirstUpdateID <= ob.lastUpdateID && ob.lastUpdateID <= ev.LastUpdateID:
		ob.partialUpdate(ev)
		atomic.AddInt64(&ob.numUpdateCall, 1)

	// ──────────────────────────────────────────────────────────────
	// 2. Entire message is older than what we already have → drop
	case ev.LastUpdateID < ob.lastUpdateID:
		return

	// ──────────────────────────────────────────────────────────────
	// 3. We missed one or more updates → fetch fresh snapshot, then
	//    (optionally) re‑apply the current diff if it now fits.
	case ev.FirstUpdateID > ob.lastUpdateID:
		snap, err := cl.NewDepthService().
			Symbol(ob.Symbol).
			Limit(MaxSpotLayerRequest).
			Do(context.Background())

		atomic.AddInt64(&ob.numSnapshotCall, 1)

		if err != nil {
			ob.logger.Error("snapshot error: %s", err.Error())
			return
		}

		ob.totalUpdate(snap)

		if ev.FirstUpdateID <= ob.lastUpdateID && ob.lastUpdateID <= ev.LastUpdateID {
			ob.partialUpdate(ev)
			atomic.AddInt64(&ob.numUpdateCall, 1)
		}

	// ──────────────────────────────────────────────────────────────
	default:
		ob.logger.Error("unexpected depth‑event ordering")
	}
}

func (ob *BinanceOrderBook) SubscribeBestDepth(handler func(ask, bid PriceLevel)) func() {
	id := uuid.New().String()
	ob.subscriber.Store(id, handler)
	return func() {
		ob.subscriber.Delete(id)
	}
}

func (ob *BinanceOrderBook) publishBestDepth(ask, bid PriceLevel) {
	ob.subscriber.Range(func(key, value interface{}) bool {
		if handler, ok := value.(func(ask, bid PriceLevel)); ok {
			handler(ask, bid)
		}
		return true
	})
}

/*
Close is idempotent and blocks until every goroutine and the underlying
WebSocket have finished.  After it returns the instance is unusable.
*/
func (ob *BinanceOrderBook) Close() error {
	ob.cancel() // 1. stop internal goroutines

	// 2. tell Binance stream to shut down (non‑blocking if already closed)
	if ob.stopC != nil {
		select {
		case ob.stopC <- struct{}{}:
		default:
		}
	}

	// 3. wait for listener goroutine
	ob.wg.Wait()
	return nil
}

/* ======= utility helpers ======= */

func (ob *BinanceOrderBook) Summary() {
	ob.logger.Info("NumUpdateCall: %d", atomic.LoadInt64(&ob.numUpdateCall))
	ob.logger.Info("NumSnapshotCall: %d", atomic.LoadInt64(&ob.numSnapshotCall))
}

func (ob *BinanceOrderBook) GetDepth(depth int) ([]PriceLevel, []PriceLevel, error) {
	asks := ob.Asks.GetBook(depth, true)
	bids := ob.Bids.GetBook(depth, false)
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
	/* ======= public, read‑only fields ======= */
	Symbol string
	Asks   AskBookArray
	Bids   BidBookArray

	/* ======= streaming internals ======= */
	timestamp    int64
	lastUpdateID int64
	createdAt    int64

	/* ======= streaming internals ======= */
	eventCh chan *futures.WsDepthEvent
	stopC   chan struct{}
	doneC   chan struct{}

	/* ======= coordination ======= */
	ctx    context.Context    // global cancel point for all goroutines
	cancel context.CancelFunc // paired with ctx
	wg     sync.WaitGroup     // waits for internal goroutines

	logger          *log.Logger // optional logger
	NumUpdateCall   int
	NumSnapshotCall int
}

func NewBinancePerpOrderBook(symbol string, bufferSize int, logger *log.Logger) *BinancePerpOrderBook {
	ctx, cancel := context.WithCancel(context.Background())
	return &BinancePerpOrderBook{
		Symbol:       symbol,
		Asks:         *NewAskBookArray(),
		Bids:         *NewBidBookArray(),
		timestamp:    0,
		lastUpdateID: 0,
		createdAt:    time.Now().UnixMilli(),

		eventCh: make(chan *futures.WsDepthEvent, bufferSize),
		stopC:   nil,
		doneC:   nil,

		ctx:    ctx,
		cancel: cancel,
		logger: logger,

		NumUpdateCall:   0,
		NumSnapshotCall: 0,
	}
}

func (ob *BinancePerpOrderBook) Start(spd UpdateSpeed) error { // blocking call
	if ob.stopC != nil {
		return errors.New("orderbook already started")
	}
	wsDepthHandler := func(ev *futures.WsDepthEvent) {
		select {
		case ob.eventCh <- ev:
		default:
			ob.logger.Info("[%s-PERP] depth‑event dropped (buffer full)", ob.Symbol)
		}
	}
	errHandler := func(err error) { ob.logger.Error("WS error: %+v", err) }

	var err error

	// Start the WebSocket connection
	switch spd {
	case UpdateSpeed100Ms:
		ob.doneC, ob.stopC, err = futures.WsDiffDepthServeWithRate(ob.Symbol, 100*time.Millisecond, wsDepthHandler, errHandler)
	case UpdateSpeed250Ms:
		ob.doneC, ob.stopC, err = futures.WsDiffDepthServeWithRate(ob.Symbol, 250*time.Millisecond, wsDepthHandler, errHandler)
	case UpdateSpeed500Ms:
		ob.doneC, ob.stopC, err = futures.WsDiffDepthServeWithRate(ob.Symbol, 500*time.Millisecond, wsDepthHandler, errHandler)
	default:
		return errors.New("invalid update speed")
	}

	if err != nil {
		return err
	}

	// ----- 3. launch listener -----
	ob.wg.Add(1)
	go ob.listen()

	return nil
}

/*
listen consumes the depth events and terminates when:
  - Binance closes the socket (doneC closed), or
  - The caller invokes Close() (ctx cancelled).

It is the only place that reads from eventCh, so when it returns we can
safely close(eventCh) to free memory.
*/
func (ob *BinancePerpOrderBook) listen() {
	defer ob.wg.Done()
	defer close(ob.eventCh)

	client := futures.NewClient("", "")

	for {
		select {
		case <-ob.ctx.Done():
			return // caller asked us to quit
		case <-ob.doneC: // WS layer finished
			return
		case ev := <-ob.eventCh: // main path
			ob.handleDepthEvent(client, ev)
		}
	}
}

func (ob *BinancePerpOrderBook) handleDepthEvent(cl *futures.Client, ev *futures.WsDepthEvent) {
	switch {

	case ev.LastUpdateID < ob.lastUpdateID:
		// 4. Drop any event where u is < lastUpdateId in the snapshot
		return

	case ev.PrevLastUpdateID == ob.lastUpdateID:
		// 6. While listening to the stream, each new event's pu should be equal to the previous event's u, otherwise initialize the process from step 3.
		ob.partialUpdate(ev)
		ob.NumUpdateCall++
	default:
		snapshot, err := cl.NewDepthService().
			Symbol(ob.Symbol).
			Limit(MaxPerpLayerRequest).
			Do(context.Background())
		ob.NumSnapshotCall++
		if err != nil {
			ob.logger.Error("Error getting snapshot: %s", err.Error())
			return
		}
		ob.totalUpdate(snapshot)
		for event := range ob.eventCh {
			if event.FirstUpdateID <= ob.lastUpdateID && ob.lastUpdateID <= event.LastUpdateID {
				// 5.The first processed event should have U <= lastUpdateId**AND**u >= lastUpdateId
				ob.partialUpdate(event)
				break
			} else if event.LastUpdateID < ob.lastUpdateID {
				return
			} else {
				ob.logger.Error("Unexpected event state")
				return
			}
		}
	}
}

func (ob *BinancePerpOrderBook) Summary() {
	ob.logger.Info("NumUpdateCall: %d", ob.NumUpdateCall)
	ob.logger.Info("NumSnapshotCall: %d", ob.NumSnapshotCall)
}

/*
Close is idempotent and blocks until every goroutine and the underlying
WebSocket have finished.  After it returns the instance is unusable.
*/
func (ob *BinancePerpOrderBook) Close() error {
	ob.cancel() // 1. stop internal goroutines

	// 2. tell Binance stream to shut down (non‑blocking if already closed)
	if ob.stopC != nil {
		select {
		case ob.stopC <- struct{}{}:
		default:
		}
	}

	// 3. wait for listener goroutine
	ob.wg.Wait()
	return nil
}

func (ob *BinancePerpOrderBook) GetDepth(depth int) ([]PriceLevel, []PriceLevel, error) {
	asks := ob.Asks.GetBook(depth, true)
	bids := ob.Bids.GetBook(depth, false)
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
	ob.logger.Info("Snapshot LastUpdateID: %d", ob.lastUpdateID)
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
