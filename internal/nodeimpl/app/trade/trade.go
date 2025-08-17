package trade

import (
	"encoding/json"
	"strconv"

	"github.com/BullionBear/sequex/internal/model/protobuf/app"
	pbCommon "github.com/BullionBear/sequex/internal/model/protobuf/common"
	"github.com/BullionBear/sequex/internal/nodeimpl/base"
	"github.com/BullionBear/sequex/pkg/eventbus"
	"github.com/BullionBear/sequex/pkg/exchange/binance"
	"github.com/BullionBear/sequex/pkg/log"
	"github.com/BullionBear/sequex/pkg/node"
	"google.golang.org/protobuf/proto"
)

const (
	EmitTradeKey = "emit_trade"

	RpcReqMetadataKey   = "req_metadata"
	RpcReqParametersKey = "req_parameters"
	RpcReqStatusKey     = "req_status"
)

type TradeParams struct {
	Exchange   string `json:"exchange"`
	Instrument string `json:"instrument"`
	Symbol     string `json:"symbol"`
}

type TradeNode struct {
	*base.BaseNode
	// Configurable parameters
	cfg TradeParams

	wsClient  *binance.WSClient
	currentId int64
	shutdownC chan struct{}
	doneC     chan struct{}
}

func init() {
	node.RegisterNode("trade", NewTradeNode)
}

func NewTradeNode(name string, eb *eventbus.EventBus, config *node.NodeConfig, logger log.Logger) (node.Node, error) {
	baseNode := base.NewBaseNode(name, eb, config, logger)

	cfg := TradeParams{
		Exchange:   config.Params["exchange"].(string),
		Instrument: config.Params["instrument"].(string),
		Symbol:     config.Params["symbol"].(string),
	}

	return &TradeNode{
		BaseNode:  baseNode,
		cfg:       cfg,
		wsClient:  binance.NewWSClient(&binance.WSConfig{}),
		currentId: 0,
		shutdownC: make(chan struct{}),
		doneC:     make(chan struct{}),
	}, nil
}

func (n *TradeNode) Start() error {
	n.Logger().Info("Starting Trade node")
	go n.emitTrade(n.shutdownC, n.doneC)
	if metadata, err := n.GetRpc(RpcReqMetadataKey); err != nil {
		return err
	} else {
		n.EventBus().RegisterRPC(metadata, func() proto.Message {
			return &pbCommon.MetadataRequest{}
		}, func(event proto.Message) proto.Message {
			if req, ok := event.(*pbCommon.MetadataRequest); ok {
				return n.RequestMetadata(req)
			}
			return &pbCommon.MetadataResponse{
				Id:      -1,
				Code:    pbCommon.ErrorCode_ERROR_CODE_INVALID_REQUEST,
				Message: "Invalid request",
			}
		})
	}
	if parameters, err := n.GetRpc(RpcReqParametersKey); err != nil {
		return err
	} else {
		n.EventBus().RegisterRPC(parameters, func() proto.Message {
			return &pbCommon.ParametersRequest{}
		}, func(event proto.Message) proto.Message {
			if req, ok := event.(*pbCommon.ParametersRequest); ok {
				return n.RequestParameters(req)
			}
			return &pbCommon.ParametersResponse{
				Id:      -1,
				Code:    pbCommon.ErrorCode_ERROR_CODE_INVALID_REQUEST,
				Message: "Invalid request",
			}
		})
	}
	if status, err := n.GetRpc(RpcReqStatusKey); err != nil {
		return err
	} else {
		n.EventBus().RegisterRPC(status, func() proto.Message {
			return &pbCommon.StatusRequest{}
		}, func(event proto.Message) proto.Message {
			if req, ok := event.(*pbCommon.StatusRequest); ok {
				return n.RequestStatus(req)
			}
			return &pbCommon.StatusResponse{
				Id:      -1,
				Code:    pbCommon.ErrorCode_ERROR_CODE_INVALID_REQUEST,
				Message: "Invalid request",
			}
		})
	}
	return nil
}

func (n *TradeNode) Shutdown() error {
	n.Logger().Info("Shutting down Trade node")
	close(n.shutdownC)
	<-n.doneC
	return nil
}

func (n *TradeNode) emitTrade(shutdownC chan struct{}, doneC chan struct{}) {
	subject, err := n.GetEmit(EmitTradeKey)
	if err != nil {
		n.Logger().Error("Failed to get emit subject", log.Error(err))
		return
	}
	unsubscribe, err := n.wsClient.SubscribeTrade(n.cfg.Symbol, binance.TradeSubscriptionOptions{
		OnTrade: func(trade binance.WSTrade) {
			baseAsset, err := binance.GetBaseAsset(trade.Symbol)
			if err != nil {
				n.Logger().Error("Failed to get base asset", log.Error(err))
				return
			}
			quoteAsset, err := binance.GetQuoteAsset(trade.Symbol)
			if err != nil {
				n.Logger().Error("Failed to get quote asset", log.Error(err))
				return
			}
			symbol := app.Symbol{
				Base:  baseAsset,
				Quote: quoteAsset,
			}
			price, err := strconv.ParseFloat(trade.Price, 64)
			if err != nil {
				n.Logger().Error("Failed to parse price", log.Error(err))
				return
			}
			quantity, err := strconv.ParseFloat(trade.Quantity, 64)
			if err != nil {
				n.Logger().Error("Failed to parse quantity", log.Error(err))
				return
			}
			takerSide := app.Side_SIDE_BUY
			if trade.IsBuyerMaker {
				takerSide = app.Side_SIDE_SELL
			}
			appTrade := &app.Trade{
				Id:         trade.TradeId,
				Exchange:   app.Exchange_EXCHANGE_BINANCE,
				Instrument: app.Instrument_INSTRUMENT_SPOT,
				Symbol:     &symbol,
				Price:      price,
				Side:       takerSide,
				Quantity:   quantity,
				Timestamp:  trade.TradeTime,
			}
			n.currentId++
			n.EventBus().Emit(subject, appTrade)
			n.Logger().Infof("Emitting trade %d", appTrade.Id)
		},
		OnError: func(err error) {
			n.Logger().Error("Failed to subscribe to trade", log.Error(err))
		},
		OnDisconnect: func() {
			n.Logger().Info("Disconnected from trade")
		},
		OnConnect: func() {
			n.Logger().Info("Connected to trade")
		},
		OnReconnect: func() {
			n.Logger().Info("Reconnected to trade")
		},
	})
	if err != nil {
		n.Logger().Error("Failed to subscribe to trade", log.Error(err))
		return
	}
	<-shutdownC
	unsubscribe()
	close(doneC)
}

func (n *TradeNode) RequestParameters(req *pbCommon.ParametersRequest) *pbCommon.ParametersResponse {
	jsonBytes, err := json.Marshal(n.cfg)
	if err != nil {
		n.Logger().Error("Failed to marshal parameters", log.Error(err))
		return &pbCommon.ParametersResponse{
			Id:      -1,
			Code:    pbCommon.ErrorCode_ERROR_CODE_SERIALIZATION_ERROR,
			Message: "Failed to json marshal parameters",
		}
	}
	return &pbCommon.ParametersResponse{
		Id:         req.Id,
		Code:       pbCommon.ErrorCode_ERROR_CODE_OK,
		Message:    "",
		Parameters: jsonBytes,
	}
}

func (n *TradeNode) RequestStatus(req *pbCommon.StatusRequest) *pbCommon.StatusResponse {
	jsonBytes, err := json.Marshal(map[string]any{
		"currentId": n.currentId,
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
