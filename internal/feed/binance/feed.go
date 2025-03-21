package binance

import (
	"log"
	"strconv"

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
		k := event.Kline

		openPx, _ := strconv.ParseFloat(k.Open, 64)
		highPx, _ := strconv.ParseFloat(k.High, 64)
		lowPx, _ := strconv.ParseFloat(k.Low, 64)
		closePx, _ := strconv.ParseFloat(k.Close, 64)
		baseVolume, _ := strconv.ParseFloat(k.Volume, 64)
		quoteVolume, _ := strconv.ParseFloat(k.QuoteVolume, 64)
		takerBaseVolume, _ := strconv.ParseFloat(k.ActiveBuyVolume, 64)
		takerQuoteVolume, _ := strconv.ParseFloat(k.ActiveBuyQuoteVolume, 64)

		handler(&payload.KLineUpdate{
			Symbol:                   k.Symbol,
			Interval:                 k.Interval,
			OpenTime:                 k.StartTime,
			CloseTime:                k.EndTime,
			EventTime:                event.Time, // Event time is the time the event was received
			OpenPx:                   openPx,
			HighPx:                   highPx,
			LowPx:                    lowPx,
			ClosePx:                  closePx,
			NumberOfTrades:           int(k.TradeNum),
			BaseAssetVolume:          baseVolume,
			QuoteAssetVolume:         quoteVolume,
			TakerBuyBaseAssetVolume:  takerBaseVolume,
			TakerBuyQuoteAssetVolume: takerQuoteVolume,
			IsClosed:                 k.IsFinal,
		})
	}, func(err error) { log.Printf("error: %v", err) })
	if err != nil {
		return nil, err
	}
	return func() {
		stopC <- struct{}{}
		<-doneC
	}, nil
}
