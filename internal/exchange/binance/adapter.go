package binance

import (
	"context"

	"github.com/BullionBear/sequex/internal/exchange"
	"github.com/BullionBear/sequex/pkg/exchange/binance"
)

var _ exchange.Connector = (*BinanceExchangeAdapter)(nil)

func NewBinanceAdapter(cfg exchange.Config) *BinanceExchangeAdapter {
	wsClient := binance.NewWSClient(&binance.WSConfig{
		APIKey:      cfg.Credentials.APIKey,
		APISecret:   cfg.Credentials.APISecret,
		BaseWsURL:   binance.MainnetWSBaseUrl,
		BaseRestURL: binance.MainnetBaseUrl,
	})
	restClient := wsClient.GetRestClient()
	return &BinanceExchangeAdapter{cfg: cfg, restClient: restClient, wsClient: wsClient}
}

type BinanceExchangeAdapter struct {
	cfg        exchange.Config
	restClient *binance.Client
	wsClient   *binance.WSClient
}

func (a *BinanceExchangeAdapter) GetBalance(ctx context.Context) (exchange.Response[[]exchange.Balance], error) {
	resp, err := a.restClient.GetAccountInfo(ctx, binance.GetAccountInfoRequest{
		OmitZeroBalances: true,
		RecvWindow:       5000,
	})
	if err != nil {
		return exchange.Response[[]exchange.Balance]{}, err
	}
	if resp.Code != 0 {
		return exchange.Response[[]exchange.Balance]{
			Code:    resp.Code,
			Message: resp.Message,
		}, nil
	}
	balances := make([]exchange.Balance, len(resp.Data.Balances))
	for i, balance := range resp.Data.Balances {
		balances[i] = exchange.Balance{
			Asset:  balance.Asset,
			Free:   balance.Free,
			Locked: balance.Locked,
		}
	}
	return exchange.Response[[]exchange.Balance]{
		Code:    200,
		Message: "OK",
		Data:    balances,
	}, nil
}
