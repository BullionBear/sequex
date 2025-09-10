package trade

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/BullionBear/sequex/internal/adapter"
	"github.com/BullionBear/sequex/internal/model/sqx"
	"github.com/BullionBear/sequex/pkg/exchange/binance"
	"github.com/BullionBear/sequex/pkg/logger"
)

func init() {
	binanceTradeAdapter := NewBinanceTradeAdapter()
	logger.Log.Info().Msg("Registering Binance trade adapter")
	adapter.RegisterTradeAdapter(sqx.ExchangeBinance, binanceTradeAdapter)
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
	if instrumentType != sqx.InstrumentTypeSpot {
		return nil, fmt.Errorf("instrument type not supported: %s", instrumentType)
	}
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
			// Parse the symbol to extract base and quote currencies
			// For BTCUSDT, we need to extract BTC and USDT
			// Common quote currencies: USDT, USDC, BUSD, BTC, ETH, BNB
			quoteCurrencies := []string{"USDT", "USDC", "BUSD", "BTC", "ETH", "BNB"}
			var base, quote string
			for _, qc := range quoteCurrencies {
				if strings.HasSuffix(wsTrade.Symbol, qc) {
					base = strings.TrimSuffix(wsTrade.Symbol, qc)
					quote = qc
					break
				}
			}
			if base == "" || quote == "" {
				logger.Log.Error().Msgf("Failed to parse symbol: %s", wsTrade.Symbol)
				return
			}

			trade := sqx.Trade{
				Id:             wsTrade.TradeId,
				Symbol:         sqx.NewSymbol(base, quote),
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
