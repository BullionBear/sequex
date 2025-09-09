package trade

import (
	"fmt"

	"github.com/BullionBear/sequex/internal/adapter"
	"github.com/BullionBear/sequex/internal/model/sqx"
	"github.com/BullionBear/sequex/pkg/exchange/binance"
	"github.com/BullionBear/sequex/pkg/logger"
)

func init() {
	binanceTradeAdapter := NewBinanceTradeAdapter()
	logger.Log.Info().Msg("Registering Binance trade adapter")
	adapter.RegisterAdapter(sqx.ExchangeBinance, sqx.DataTypeTrade, binanceTradeAdapter)
}

type BinanceTradeAdapter struct {
	wsClient *binance.WSClient
}

func NewBinanceTradeAdapter() *BinanceTradeAdapter {
	return &BinanceTradeAdapter{
		wsClient: binance.NewWSClient(binance.NewMainnetWSConfig("", "")),
	}
}

func (a *BinanceTradeAdapter) Subscribe(symbol sqx.Symbol, instrumentType sqx.InstrumentType, callback adapter.Callback) (func(), error) {
	binanceSymbol := fmt.Sprintf("%s%s", symbol.Base, symbol.Quote)
	return a.wsClient.SubscribeTrade(binanceSymbol, binance.TradeSubscriptionOptions{
		OnTrade: func(trade binance.WSTrade) {
			logger.Log.Info().Msgf("Received trade: %+v", trade)
		},
	})
}
