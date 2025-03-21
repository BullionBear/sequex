package payload

type KLine struct {
	Symbol                   string  `json:"symbol"`
	Interval                 string  `json:"interval"`
	OpenTime                 int64   `json:"open_time"`
	CloseTime                int64   `json:"close_time"`
	OpenPx                   float64 `json:"open_px"`
	HighPx                   float64 `json:"high_px"`
	LowPx                    float64 `json:"low_px"`
	ClosePx                  float64 `json:"close_px"`
	NumberOfTrades           int     `json:"number_of_trades"`
	BaseAssetVolume          float64 `json:"base_asset_volume"`
	QuoteAssetVolume         float64 `json:"quote_asset_volume"`
	TakerBuyBaseAssetVolume  float64 `json:"taker_buy_base_asset_volume"`
	TakerBuyQuoteAssetVolume float64 `json:"taker_buy_quote_asset_volume"`
	IsClosed                 bool    `json:"is_closed"`
}
