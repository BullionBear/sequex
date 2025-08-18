package trade

import (
	"fmt"
	"strconv"

	"github.com/BullionBear/sequex/internal/nodeimpl/app/share"
	"github.com/BullionBear/sequex/pkg/exchange/binance"
)

func init() {
	RegisterAdapter(share.ExchangeBinance, share.InstrumentSpot, NewBinanceSubscribeTradeAdapter())
}

type BinanceSubscribeTradeAdapter struct {
	wsClient *binance.WSClient
}

func NewBinanceSubscribeTradeAdapter() *BinanceSubscribeTradeAdapter {
	wsClient := binance.NewWSClient(&binance.WSConfig{})
	return &BinanceSubscribeTradeAdapter{
		wsClient: wsClient,
	}
}

func (a *BinanceSubscribeTradeAdapter) Subscribe(symbol share.Symbol, options TradeSubscriptionOptions) (func(), error) {
	unsubscribe, err := a.wsClient.SubscribeTrade(fmt.Sprintf("%s%s", symbol.Base, symbol.Quote), binance.TradeSubscriptionOptions{
		OnTrade: func(trade binance.WSTrade) {
			side := share.SideBuy
			if trade.IsBuyerMaker {
				side = share.SideSell
			}
			price, err := strconv.ParseFloat(trade.Price, 64)
			if err != nil {
				options.OnError(err)
				return
			}
			quantity, err := strconv.ParseFloat(trade.Quantity, 64)
			if err != nil {
				options.OnError(err)
				return
			}
			options.OnTrade(Trade{
				Symbol:    symbol,
				ID:        int64(trade.TradeId),
				Price:     price,
				Qty:       quantity,
				Time:      trade.TradeTime,
				TakerSide: side,
			})
		},
		OnConnect:    options.OnConnect,
		OnReconnect:  options.OnReconnect,
		OnDisconnect: options.OnDisconnect,
	})
	if err != nil {
		return nil, err
	}
	return unsubscribe, nil
}
