package bar

import (
	"encoding/json"
	"fmt"
	"math"
	"sort"

	pbCommon "github.com/BullionBear/sequex/internal/model/protobuf/common"
	"github.com/BullionBear/sequex/pkg/log"
	"google.golang.org/protobuf/proto"

	"github.com/BullionBear/sequex/internal/model/protobuf/app"
	"github.com/BullionBear/sequex/internal/nodeimpl/base"
	"github.com/BullionBear/sequex/pkg/eventbus"
	"github.com/BullionBear/sequex/pkg/node"
	"gonum.org/v1/gonum/floats"
	"gonum.org/v1/gonum/stat"
)

const (
	EmitBarKey = "emit_bar"

	OnTradeKey          = "on_trade"
	RpcReqMetadataKey   = "req_metadata"
	RpcReqParametersKey = "req_parameters"
	RpcReqStatusKey     = "req_status"
)

type TimeBarParams struct {
	Interval int64 `json:"interval"`
}

type TimeBarNode struct {
	*base.BaseNode
	cfg TimeBarParams
	// Variables
	isFirstBar       bool
	currentTimeframe int64
	tradeBuffer      []*app.Trade
}

func init() {
	node.RegisterNode("timebar", NewTimeBarNode)
}

func NewTimeBarNode(name string, eb *eventbus.EventBus, config *node.NodeConfig, logger log.Logger) (node.Node, error) {
	baseNode := base.NewBaseNode(name, eb, config, logger)
	interval, ok := config.Params["interval"].(int)
	if !ok {
		return nil, fmt.Errorf("interval is not an int")
	}
	cfg := TimeBarParams{
		Interval: int64(interval),
	}
	return &TimeBarNode{
		BaseNode:         baseNode,
		cfg:              cfg,
		isFirstBar:       true,
		currentTimeframe: 0,
		tradeBuffer:      make([]*app.Trade, 0),
	}, nil
}

func (n *TimeBarNode) Start() error {
	subject, err := n.GetOn(OnTradeKey)
	if err != nil {
		n.Logger().Error("Failed to get on subject", log.Error(err))
		return err
	}
	n.EventBus().On(subject, func() proto.Message {
		return &app.Trade{}
	}, func(event proto.Message) {
		n.onTrade(event.(*app.Trade))
	})
	return nil
}

func (n *TimeBarNode) Shutdown() error {
	n.Logger().Info("Shutting down timebar node")
	return nil
}

func (n *TimeBarNode) onTrade(trade *app.Trade) {
	n.Logger().Infof("Received trade %d", trade.Id)
	if n.currentTimeframe == 0 {
		n.currentTimeframe = trade.Timestamp / n.cfg.Interval
		return
	}
	if n.currentTimeframe == trade.Timestamp/n.cfg.Interval {
		n.tradeBuffer = append(n.tradeBuffer, trade)
		return
	}
	n.currentTimeframe = trade.Timestamp / n.cfg.Interval
	priceBuffer := make([]float64, len(n.tradeBuffer))
	for _, trade := range n.tradeBuffer {
		priceBuffer = append(priceBuffer, trade.Price)
	}
	baseBuffer := make([]float64, len(n.tradeBuffer))
	for _, trade := range n.tradeBuffer {
		baseBuffer = append(baseBuffer, trade.Quantity)
	}
	quoteBuffer := make([]float64, len(n.tradeBuffer))
	for _, trade := range n.tradeBuffer {
		quoteBuffer = append(quoteBuffer, trade.Price*trade.Quantity)
	}
	sortedPriceBuffer := append([]float64(nil), priceBuffer...)
	sort.Float64s(sortedPriceBuffer)
	bar := app.Bar{
		Symbol:        trade.Symbol,
		Instrument:    trade.Instrument,
		Exchange:      trade.Exchange,
		StartSeq:      n.tradeBuffer[0].Id,
		EndSeq:        n.tradeBuffer[len(n.tradeBuffer)-1].Id,
		NextSeq:       trade.Id,
		StartTime:     n.tradeBuffer[0].Timestamp,
		EndTime:       trade.Timestamp,
		Open:          priceBuffer[0],
		High:          sortedPriceBuffer[len(sortedPriceBuffer)-1],
		Low:           sortedPriceBuffer[0],
		Close:         priceBuffer[len(priceBuffer)-1],
		Mean:          stat.Mean(priceBuffer, nil),
		Std:           math.Sqrt(stat.Variance(priceBuffer, nil)),
		Median:        stat.Quantile(0.5, stat.Empirical, sortedPriceBuffer, nil),
		FirstQuartile: stat.Quantile(0.25, stat.Empirical, sortedPriceBuffer, nil),
		ThirdQuartile: stat.Quantile(0.75, stat.Empirical, sortedPriceBuffer, nil),
		VolumeBase:    floats.Sum(baseBuffer),
		VolumeQuote:   floats.Sum(quoteBuffer),
		Count:         int64(len(n.tradeBuffer)),
	}
	n.tradeBuffer = make([]*app.Trade, 0)
	subject, err := n.GetEmit(EmitBarKey)
	if err != nil {
		n.Logger().Error("Failed to get emit subject", log.Error(err))
		return
	}
	if n.isFirstBar {
		n.isFirstBar = false
		return
	}
	n.EventBus().Emit(subject, &bar)
}

func (n *TimeBarNode) RequestParameters(req *pbCommon.ParametersRequest) *pbCommon.ParametersResponse {
	jsonBytes, err := json.Marshal(n.cfg)
	if err != nil {
		n.Logger().Error("Failed to marshal parameters", log.Error(err))
		return &pbCommon.ParametersResponse{
			Id:         req.Id,
			Code:       pbCommon.ErrorCode_ERROR_CODE_SERIALIZATION_ERROR,
			Message:    "Failed to json marshal parameters",
			Parameters: jsonBytes,
		}
	}
	return &pbCommon.ParametersResponse{
		Id:         req.Id,
		Code:       pbCommon.ErrorCode_ERROR_CODE_OK,
		Message:    "",
		Parameters: jsonBytes,
	}
}

func (n *TimeBarNode) RequestStatus(req *pbCommon.StatusRequest) *pbCommon.StatusResponse {
	jsonBytes, err := json.Marshal(map[string]any{
		"is_first_bar":        n.isFirstBar,
		"current_open_time":   n.currentTimeframe * n.cfg.Interval,
		"trade_buffer_length": len(n.tradeBuffer),
	})
	if err != nil {
		n.Logger().Error("Failed to marshal status", log.Error(err))
		return &pbCommon.StatusResponse{
			Id:      req.Id,
			Code:    pbCommon.ErrorCode_ERROR_CODE_SERIALIZATION_ERROR,
			Message: "Failed to json marshal status",
		}
	}
	return &pbCommon.StatusResponse{
		Id:      req.Id,
		Code:    pbCommon.ErrorCode_ERROR_CODE_OK,
		Message: "",
		Status:  jsonBytes,
	}
}
