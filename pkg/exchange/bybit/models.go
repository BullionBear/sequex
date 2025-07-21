package bybit

import (
	"fmt"
	"strconv"
	"time"
)

// Common API Response Structure
type APIResponse struct {
	RetCode    int         `json:"retCode"`
	RetMsg     string      `json:"retMsg"`
	Result     interface{} `json:"result"`
	RetExtInfo interface{} `json:"retExtInfo"`
	Time       int64       `json:"time"`
}

// Kline Response Models
type KlineResponse struct {
	RetCode    int         `json:"retCode"`
	RetMsg     string      `json:"retMsg"`
	Result     KlineResult `json:"result"`
	RetExtInfo interface{} `json:"retExtInfo"`
	Time       int64       `json:"time"`
}

type KlineResult struct {
	Symbol   string     `json:"symbol"`
	Category string     `json:"category"`
	List     [][]string `json:"list"`
}

// KlineData represents a single kline/candlestick
type KlineData struct {
	Timestamp  time.Time
	OpenPrice  float64
	HighPrice  float64
	LowPrice   float64
	ClosePrice float64
	Volume     float64
	Turnover   float64
}

// ParseKlineData converts string array to KlineData struct
func ParseKlineData(data []string) (*KlineData, error) {
	if len(data) < 7 {
		return nil, fmt.Errorf("insufficient data for kline: expected 7 elements, got %d", len(data))
	}

	// Parse timestamp (milliseconds)
	timestamp, err := strconv.ParseInt(data[0], 10, 64)
	if err != nil {
		return nil, fmt.Errorf("failed to parse timestamp: %w", err)
	}

	// Parse prices
	openPrice, err := strconv.ParseFloat(data[1], 64)
	if err != nil {
		return nil, fmt.Errorf("failed to parse open price: %w", err)
	}

	highPrice, err := strconv.ParseFloat(data[2], 64)
	if err != nil {
		return nil, fmt.Errorf("failed to parse high price: %w", err)
	}

	lowPrice, err := strconv.ParseFloat(data[3], 64)
	if err != nil {
		return nil, fmt.Errorf("failed to parse low price: %w", err)
	}

	closePrice, err := strconv.ParseFloat(data[4], 64)
	if err != nil {
		return nil, fmt.Errorf("failed to parse close price: %w", err)
	}

	volume, err := strconv.ParseFloat(data[5], 64)
	if err != nil {
		return nil, fmt.Errorf("failed to parse volume: %w", err)
	}

	turnover, err := strconv.ParseFloat(data[6], 64)
	if err != nil {
		return nil, fmt.Errorf("failed to parse turnover: %w", err)
	}

	return &KlineData{
		Timestamp:  time.Unix(timestamp/1000, (timestamp%1000)*int64(time.Millisecond)),
		OpenPrice:  openPrice,
		HighPrice:  highPrice,
		LowPrice:   lowPrice,
		ClosePrice: closePrice,
		Volume:     volume,
		Turnover:   turnover,
	}, nil
}

// KlineRequest represents the request parameters for kline data
type KlineRequest struct {
	Category string `json:"category"`
	Symbol   string `json:"symbol"`
	Interval string `json:"interval"`
	Start    int64  `json:"start,omitempty"`
	End      int64  `json:"end,omitempty"`
	Limit    int    `json:"limit,omitempty"`
}

// ServerTimeResponse represents the server time response
type ServerTimeResponse struct {
	RetCode    int         `json:"retCode"`
	RetMsg     string      `json:"retMsg"`
	Result     ServerTime  `json:"result"`
	RetExtInfo interface{} `json:"retExtInfo"`
	Time       int64       `json:"time"`
}

type ServerTime struct {
	TimeSecond string `json:"timeSecond"`
	TimeNano   string `json:"timeNano"`
}

// TickerResponse represents the ticker response
type TickerResponse struct {
	RetCode    int         `json:"retCode"`
	RetMsg     string      `json:"retMsg"`
	Result     TickerList  `json:"result"`
	RetExtInfo interface{} `json:"retExtInfo"`
	Time       int64       `json:"time"`
}

type TickerList struct {
	Category string   `json:"category"`
	List     []Ticker `json:"list"`
}

type Ticker struct {
	Symbol               string `json:"symbol"`
	LastPrice            string `json:"lastPrice"`
	IndexPrice           string `json:"indexPrice"`
	MarkPrice            string `json:"markPrice"`
	PrevPrice24h         string `json:"prevPrice24h"`
	Price24hPcnt         string `json:"price24hPcnt"`
	HighPrice24h         string `json:"highPrice24h"`
	LowPrice24h          string `json:"lowPrice24h"`
	Turnover24h          string `json:"turnover24h"`
	Volume24h            string `json:"volume24h"`
	UsdIndexPrice        string `json:"usdIndexPrice"`
	OpenInterest         string `json:"openInterest"`
	OpenInterestValue    string `json:"openInterestValue"`
	FundingRate          string `json:"fundingRate"`
	NextFundingTime      string `json:"nextFundingTime"`
	BasisRate            string `json:"basisRate"`
	DeliveryFeeRate      string `json:"deliveryFeeRate"`
	DeliveryTime         string `json:"deliveryTime"`
	OpenInterest24h      string `json:"openInterest24h"`
	OpenInterestValue24h string `json:"openInterestValue24h"`
	Basis                string `json:"basis"`
	BasisRate24h         string `json:"basisRate24h"`
	BasisValue24h        string `json:"basisValue24h"`
	FundingRate24h       string `json:"fundingRate24h"`
	MarkPrice24h         string `json:"markPrice24h"`
	IndexPrice24h        string `json:"indexPrice24h"`
	PrevPrice1h          string `json:"prevPrice1h"`
	Price1hPcnt          string `json:"price1hPcnt"`
	HighPrice1h          string `json:"highPrice1h"`
	LowPrice1h           string `json:"lowPrice1h"`
	Turnover1h           string `json:"turnover1h"`
	Volume1h             string `json:"volume1h"`
	PrevPrice30m         string `json:"prevPrice30m"`
	Price30mPcnt         string `json:"price30mPcnt"`
	HighPrice30m         string `json:"highPrice30m"`
	LowPrice30m          string `json:"lowPrice30m"`
	Turnover30m          string `json:"turnover30m"`
	Volume30m            string `json:"volume30m"`
	PrevPrice15m         string `json:"prevPrice15m"`
	Price15mPcnt         string `json:"price15mPcnt"`
	HighPrice15m         string `json:"highPrice15m"`
	LowPrice15m          string `json:"lowPrice15m"`
	Turnover15m          string `json:"turnover15m"`
	Volume15m            string `json:"volume15m"`
	PrevPrice5m          string `json:"prevPrice5m"`
	Price5mPcnt          string `json:"price5mPcnt"`
	HighPrice5m          string `json:"highPrice5m"`
	LowPrice5m           string `json:"lowPrice5m"`
	Turnover5m           string `json:"turnover5m"`
	Volume5m             string `json:"volume5m"`
	PrevPrice1m          string `json:"prevPrice1m"`
	Price1mPcnt          string `json:"price1mPcnt"`
	HighPrice1m          string `json:"highPrice1m"`
	LowPrice1m           string `json:"lowPrice1m"`
	Turnover1m           string `json:"turnover1m"`
	Volume1m             string `json:"volume1m"`
}

// Trading Models

// CreateOrderRequest represents the request to create an order
type CreateOrderRequest struct {
	Category       string `json:"category"`
	Symbol         string `json:"symbol"`
	Side           string `json:"side"`
	OrderType      string `json:"orderType"`
	Qty            string `json:"qty"`
	Price          string `json:"price,omitempty"`
	TimeInForce    string `json:"timeInForce,omitempty"`
	OrderLinkId    string `json:"orderLinkId,omitempty"`
	TakeProfit     string `json:"takeProfit,omitempty"`
	StopLoss       string `json:"stopLoss,omitempty"`
	ReduceOnly     bool   `json:"reduceOnly,omitempty"`
	CloseOnTrigger bool   `json:"closeOnTrigger,omitempty"`
}

// CreateOrderResponse represents the response from creating an order
type CreateOrderResponse struct {
	RetCode    int               `json:"retCode"`
	RetMsg     string            `json:"retMsg"`
	Result     CreateOrderResult `json:"result"`
	RetExtInfo interface{}       `json:"retExtInfo"`
	Time       int64             `json:"time"`
}

type CreateOrderResult struct {
	OrderId     string `json:"orderId"`
	OrderLinkId string `json:"orderLinkId"`
}

// CancelOrderRequest represents the request to cancel an order
type CancelOrderRequest struct {
	Category    string `json:"category"`
	Symbol      string `json:"symbol"`
	OrderId     string `json:"orderId,omitempty"`
	OrderLinkId string `json:"orderLinkId,omitempty"`
}

// CancelOrderResponse represents the response from canceling an order
type CancelOrderResponse struct {
	RetCode    int               `json:"retCode"`
	RetMsg     string            `json:"retMsg"`
	Result     CancelOrderResult `json:"result"`
	RetExtInfo interface{}       `json:"retExtInfo"`
	Time       int64             `json:"time"`
}

type CancelOrderResult struct {
	OrderId      string `json:"orderId"`
	OrderLinkId  string `json:"orderLinkId"`
	Symbol       string `json:"symbol"`
	Category     string `json:"category"`
	Side         string `json:"side"`
	OrderType    string `json:"orderType"`
	Qty          string `json:"qty"`
	Price        string `json:"price"`
	TimeInForce  string `json:"timeInForce"`
	OrderStatus  string `json:"orderStatus"`
	CreatedTime  string `json:"createdTime"`
	UpdatedTime  string `json:"updatedTime"`
	AvgPrice     string `json:"avgPrice"`
	CumExecQty   string `json:"cumExecQty"`
	CumExecValue string `json:"cumExecValue"`
	CumExecFee   string `json:"cumExecFee"`
}

// GetOrderRequest represents the request to get order information
type GetOrderRequest struct {
	Category    string `json:"category"`
	Symbol      string `json:"symbol,omitempty"`
	OrderId     string `json:"orderId,omitempty"`
	OrderLinkId string `json:"orderLinkId,omitempty"`
	SettleCoin  string `json:"settleCoin,omitempty"`
	OrderFilter string `json:"orderFilter,omitempty"`
	Limit       int    `json:"limit,omitempty"`
	Cursor      string `json:"cursor,omitempty"`
}

// GetOrderResponse represents the response from getting order information
type GetOrderResponse struct {
	RetCode    int            `json:"retCode"`
	RetMsg     string         `json:"retMsg"`
	Result     GetOrderResult `json:"result"`
	RetExtInfo interface{}    `json:"retExtInfo"`
	Time       int64          `json:"time"`
}

type GetOrderListResponse struct {
	RetCode    int                `json:"retCode"`
	RetMsg     string             `json:"retMsg"`
	Result     GetOrderListResult `json:"result"`
	RetExtInfo interface{}        `json:"retExtInfo"`
	Time       int64              `json:"time"`
}

type GetOrderListResult struct {
	Category       string           `json:"category"`
	Symbol         string           `json:"symbol"`
	List           []GetOrderResult `json:"list"`
	NextPageCursor string           `json:"nextPageCursor"`
}

type GetOrderResult struct {
	Category            string `json:"category"`
	Symbol              string `json:"symbol"`
	OrderId             string `json:"orderId"`
	OrderLinkId         string `json:"orderLinkId"`
	BlockTradeId        string `json:"blockTradeId"`
	Side                string `json:"side"`
	OrderType           string `json:"orderType"`
	Qty                 string `json:"qty"`
	Price               string `json:"price"`
	TimeInForce         string `json:"timeInForce"`
	OrderStatus         string `json:"orderStatus"`
	CreatedTime         string `json:"createdTime"`
	UpdatedTime         string `json:"updatedTime"`
	AvgPrice            string `json:"avgPrice"`
	CumExecQty          string `json:"cumExecQty"`
	CumExecValue        string `json:"cumExecValue"`
	CumExecFee          string `json:"cumExecFee"`
	TakeProfit          string `json:"takeProfit"`
	StopLoss            string `json:"stopLoss"`
	TrailingStop        string `json:"trailingStop"`
	PositionIdx         int    `json:"positionIdx"`
	LastPriceOnCreated  string `json:"lastPriceOnCreated"`
	ReduceOnly          bool   `json:"reduceOnly"`
	CloseOnTrigger      bool   `json:"closeOnTrigger"`
	Leverage            string `json:"leverage"`
	BasePrice           string `json:"basePrice"`
	TriggerPrice        string `json:"triggerPrice"`
	TriggerDirection    int    `json:"triggerDirection"`
	TriggerBy           string `json:"triggerBy"`
	TpslMode            string `json:"tpslMode"`
	TpLimitPrice        string `json:"tpLimitPrice"`
	SlLimitPrice        string `json:"slLimitPrice"`
	TpTriggerBy         string `json:"tpTriggerBy"`
	SlTriggerBy         string `json:"slTriggerBy"`
	TpOrderType         string `json:"tpOrderType"`
	SlOrderType         string `json:"slOrderType"`
	TpSize              string `json:"tpSize"`
	SlSize              string `json:"slSize"`
	TpTakeProfit        string `json:"tpTakeProfit"`
	SlStopLoss          string `json:"slStopLoss"`
	TpTriggerPrice      string `json:"tpTriggerPrice"`
	SlTriggerPrice      string `json:"slTriggerPrice"`
	TpLimitPrice2       string `json:"tpLimitPrice2"`
	SlLimitPrice2       string `json:"slLimitPrice2"`
	TpTriggerBy2        string `json:"tpTriggerBy2"`
	SlTriggerBy2        string `json:"slTriggerBy2"`
	TpOrderType2        string `json:"tpOrderType2"`
	SlOrderType2        string `json:"slOrderType2"`
	TpSize2             string `json:"tpSize2"`
	SlSize2             string `json:"slSize2"`
	TpTakeProfit2       string `json:"tpTakeProfit2"`
	SlStopLoss2         string `json:"slStopLoss2"`
	TpTriggerPrice2     string `json:"tpTriggerPrice2"`
	SlTriggerPrice2     string `json:"slTriggerPrice2"`
	PlaceType           string `json:"placeType"`
	Iv                  string `json:"iv"`
	MarketUnit          string `json:"marketUnit"`
	ContractType        string `json:"contractType"`
	ContractValue       string `json:"contractValue"`
	CategoryType        string `json:"categoryType"`
	SmPnl               string `json:"smPnl"`
	MmPnl               string `json:"mmPnl"`
	Gap                 string `json:"gap"`
	Rebate              string `json:"rebate"`
	FromRp              string `json:"fromRp"`
	RpId                string `json:"rpId"`
	RpPnl               string `json:"rpPnl"`
	StopOrderType       string `json:"stopOrderType"`
	OcoTriggerBy        string `json:"ocoTriggerBy"`
	CancelType          string `json:"cancelType"`
	CanceledBy          string `json:"canceledBy"`
	CancelReason        string `json:"cancelReason"`
	CancelReasonCode    string `json:"cancelReasonCode"`
	CancelReasonMessage string `json:"cancelReasonMessage"`
}

// AccountResponse represents the account information response
type AccountResponse struct {
	RetCode    int           `json:"retCode"`
	RetMsg     string        `json:"retMsg"`
	Result     AccountResult `json:"result"`
	RetExtInfo interface{}   `json:"retExtInfo"`
	Time       int64         `json:"time"`
}

type AccountResult struct {
	List []AccountInfo `json:"list"`
}

type AccountInfo struct {
	TotalWalletBalance      string `json:"totalWalletBalance"`
	TotalUnrealizedPnl      string `json:"totalUnrealizedPnl"`
	TotalRealizedPnl        string `json:"totalRealizedPnl"`
	TotalMarginBalance      string `json:"totalMarginBalance"`
	TotalInitialMargin      string `json:"totalInitialMargin"`
	TotalMaintenanceMargin  string `json:"totalMaintenanceMargin"`
	TotalPositionMargin     string `json:"totalPositionMargin"`
	TotalOrderMargin        string `json:"totalOrderMargin"`
	TotalAvailableBalance   string `json:"totalAvailableBalance"`
	AccountType             string `json:"accountType"`
	AccountLTV              string `json:"accountLTV"`
	AccountMMR              string `json:"accountMMR"`
	AccountIMR              string `json:"accountIMR"`
	TotalOpenOrder          string `json:"totalOpenOrder"`
	TotalOpenOrderBuy       string `json:"totalOpenOrderBuy"`
	TotalOpenOrderSell      string `json:"totalOpenOrderSell"`
	TotalOpenOrderBuyCost   string `json:"totalOpenOrderBuyCost"`
	TotalOpenOrderSellCost  string `json:"totalOpenOrderSellCost"`
	TotalDeposit            string `json:"totalDeposit"`
	TotalWithdraw           string `json:"totalWithdraw"`
	TotalTransferIn         string `json:"totalTransferIn"`
	TotalTransferOut        string `json:"totalTransferOut"`
	TotalFee                string `json:"totalFee"`
	TotalPnl                string `json:"totalPnl"`
	TotalPnl24h             string `json:"totalPnl24h"`
	TotalPnl7d              string `json:"totalPnl7d"`
	TotalPnl30d             string `json:"totalPnl30d"`
	TotalPnl365d            string `json:"totalPnl365d"`
	TotalPnlYtd             string `json:"totalPnlYtd"`
	TotalPnlQtd             string `json:"totalPnlQtd"`
	TotalPnlMtd             string `json:"totalPnlMtd"`
	TotalPnlWtd             string `json:"totalPnlWtd"`
	TotalPnlDtd             string `json:"totalPnlDtd"`
	TotalPnlHtd             string `json:"totalPnlHtd"`
	TotalPnlMtd2            string `json:"totalPnlMtd2"`
	TotalPnlWtd2            string `json:"totalPnlWtd2"`
	TotalPnlDtd2            string `json:"totalPnlDtd2"`
	TotalPnlHtd2            string `json:"totalPnlHtd2"`
	TotalPnlYtd2            string `json:"totalPnlYtd2"`
	TotalPnlQtd2            string `json:"totalPnlQtd2"`
	TotalPnl365d2           string `json:"totalPnl365d2"`
	TotalPnl7d2             string `json:"totalPnl7d2"`
	TotalPnl24h2            string `json:"totalPnl24h2"`
	TotalPnl2               string `json:"totalPnl2"`
	TotalFee2               string `json:"totalFee2"`
	TotalTransferOut2       string `json:"totalTransferOut2"`
	TotalTransferIn2        string `json:"totalTransferIn2"`
	TotalWithdraw2          string `json:"totalWithdraw2"`
	TotalDeposit2           string `json:"totalDeposit2"`
	TotalOpenOrderSellCost2 string `json:"totalOpenOrderSellCost2"`
	TotalOpenOrderBuyCost2  string `json:"totalOpenOrderBuyCost2"`
	TotalOpenOrderSell2     string `json:"totalOpenOrderSell2"`
	TotalOpenOrderBuy2      string `json:"totalOpenOrderBuy2"`
	TotalOpenOrder2         string `json:"totalOpenOrder2"`
	AccountIMR2             string `json:"accountIMR2"`
	AccountMMR2             string `json:"accountMMR2"`
	AccountLTV2             string `json:"accountLTV2"`
	AccountType2            string `json:"accountType2"`
	TotalAvailableBalance2  string `json:"totalAvailableBalance2"`
	TotalOrderMargin2       string `json:"totalOrderMargin2"`
	TotalPositionMargin2    string `json:"totalPositionMargin2"`
	TotalMaintenanceMargin2 string `json:"totalMaintenanceMargin2"`
	TotalInitialMargin2     string `json:"totalInitialMargin2"`
	TotalMarginBalance2     string `json:"totalMarginBalance2"`
	TotalRealizedPnl2       string `json:"totalRealizedPnl2"`
	TotalUnrealizedPnl2     string `json:"totalUnrealizedPnl2"`
	TotalWalletBalance2     string `json:"totalWalletBalance2"`
}

// GetAccountRequest represents the request to get account information
type GetAccountRequest struct {
	AccountType string `json:"accountType,omitempty"`
	Coin        string `json:"coin,omitempty"`
}
