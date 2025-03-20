package binance

import (
	"github.com/BullionBear/sequex/internal/feed"
	"github.com/BullionBear/sequex/internal/payload"
	"github.com/adshao/go-binance/v2"
)

var _ feed.Feed = (*BinanceFeed)(nil)

type BinanceFeed struct {
}

func NewBinanceFeed() *BinanceFeed {
	return &BinanceFeed{}
}

func (b *BinanceFeed) SubscribeKlineUpdate(symbol string, handler func(*payload.KLineUpdate)) (unsubscribe func(), err error) {
	doneC, stopC, err := binance.WsKlineServe(symbol, "1m", func(event *binance.WsKlineEvent) {
		handler(&payload.KLineUpdate{
			Symbol:    event.Kline.Symbol,
			Interval:  event.Kline.Interval,
			Timestamp: event.Time,
		})
	}, func(err error) { /* ignore */ })
	if err != nil {
		return nil, err
	}
	return func() {
		stopC <- struct{}{}
		<-doneC
	}, nil
}
