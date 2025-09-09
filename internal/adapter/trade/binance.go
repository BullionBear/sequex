package trade

import (
	"fmt"
	"strconv"

	"github.com/BullionBear/sequex/internal/adapter"
	"github.com/BullionBear/sequex/internal/model/sqx"
	"github.com/BullionBear/sequex/pkg/exchange/binance"
	"github.com/BullionBear/sequex/pkg/logger"
)

func init() {
	binanceTradeAdapter := NewBinanceTradeAdapter()
	logger.Log.Info().Msg("Registering Binance trade adapter")
	adapter.RegisterTradeAdapter(sqx.ExchangeBinance, sqx.DataTypeTrade, binanceTradeAdapter)
}

type BinanceTradeAdapter struct {
	wsClient *binance.WSClient
}

func NewBinanceTradeAdapter() *BinanceTradeAdapter {
	return &BinanceTradeAdapter{
		wsClient: binance.NewWSClient(binance.NewMainnetWSConfig("", "")),
	}
}

func (a *BinanceTradeAdapter) Subscribe(symbol sqx.Symbol, instrumentType sqx.InstrumentType, callback adapter.TradeCallback) (func(), error) {
	binanceSymbol := fmt.Sprintf("%s%s", symbol.Base, symbol.Quote)
	return a.wsClient.SubscribeTrade(binanceSymbol, binance.TradeSubscriptionOptions{
		OnTrade: func(wsTrade binance.WSTrade) {
			logger.Log.Info().Msgf("Received trade: %+v", wsTrade)
			takerSide := sqx.SideBuy
			if wsTrade.IsBuyerMaker {
				takerSide = sqx.SideSell
			}
			price, err := strconv.ParseFloat(wsTrade.Price, 64)
			if err != nil {
				logger.Log.Error().Err(err).Msgf("Failed to parse price: %s", wsTrade.Price)
				return
			}
			quantity, err := strconv.ParseFloat(wsTrade.Quantity, 64)
			if err != nil {
				logger.Log.Error().Err(err).Msgf("Failed to parse quantity: %s", wsTrade.Quantity)
				return
			}
			trade := sqx.Trade{
				Id:             wsTrade.TradeId,
				Symbol:         sqx.NewSymbol(wsTrade.Symbol, wsTrade.Symbol),
				Exchange:       sqx.ExchangeBinance,
				InstrumentType: sqx.InstrumentTypeSpot,
				TakerSide:      takerSide,
				Price:          price,
				Quantity:       quantity,
				Timestamp:      wsTrade.TradeTime,
			}

			callback(trade)
		},
	})
}
