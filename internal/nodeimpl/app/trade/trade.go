package trade

import (
	"encoding/json"
	"fmt"
	"strings"

	"github.com/BullionBear/sequex/internal/model/protobuf/app"
	pbCommon "github.com/BullionBear/sequex/internal/model/protobuf/common"
	"github.com/BullionBear/sequex/internal/nodeimpl/app/share"
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
	Exchange   share.Exchange   `json:"exchange"`
	Instrument share.Instrument `json:"instrument"`
	Symbol     share.Symbol     `json:"symbol"`
}

func (p *TradeParams) ToAppSymbol() *app.Symbol {
	return &app.Symbol{
		Base:  p.Symbol.Base,
		Quote: p.Symbol.Quote,
	}
}

func (p *TradeParams) ToInstrument() app.Instrument {
	switch p.Instrument {
	case share.InstrumentSpot:
		return app.Instrument_INSTRUMENT_SPOT
	case share.InstrumentPerp:
		return app.Instrument_INSTRUMENT_PERP
	}
	return app.Instrument_INSTRUMENT_UNSPECIFIED
}

func (p *TradeParams) ToExchange() app.Exchange {
	switch p.Exchange {
	case share.ExchangeBinance:
		return app.Exchange_EXCHANGE_BINANCE
	case share.ExchangeBinancePerp:
		return app.Exchange_EXCHANGE_BINANCE_PERP
	}
	return app.Exchange_EXCHANGE_UNSPECIFIED
}

type TradeNode struct {
	*base.BaseNode
	// Configurable parameters
	cfg        TradeParams
	restClient *binance.Client
	wsClient   *binance.WSClient
	// Variables
	currentId    int64
	nConnected   int64
	nReconnected int64
	nError       int64
	shutdownC    chan struct{}
	doneC        chan struct{}
}

func init() {
	node.RegisterNode("trade", NewTradeNode)
}

func NewTradeNode(name string, eb *eventbus.EventBus, config *node.NodeConfig, logger log.Logger) (node.Node, error) {
	baseNode := base.NewBaseNode(name, eb, config, logger)
	var exchange share.Exchange
	switch config.Params["exchange"].(string) {
	case "binance":
		exchange = share.ExchangeBinance
	case "binance_perp":
		exchange = share.ExchangeBinancePerp
	case "bybit":
		exchange = share.ExchangeBybit
	default:
		return nil, fmt.Errorf("invalid exchange: %s", config.Params["exchange"].(string))
	}
	var instrument share.Instrument
	switch config.Params["instrument"].(string) {
	case "spot":
		instrument = share.InstrumentSpot
	case "perp":
		instrument = share.InstrumentPerp
	}
	symbol := config.Params["symbol"].(string)
	if len(strings.Split(symbol, "-")) != 2 {
		return nil, fmt.Errorf("invalid symbol: %s", symbol)
	}
	base := strings.Split(symbol, "-")[0]
	quote := strings.Split(symbol, "-")[1]

	cfg := TradeParams{
		Exchange:   exchange,
		Instrument: instrument,
		Symbol:     share.Symbol{Base: base, Quote: quote},
	}

	wsClient := binance.NewWSClient(&binance.WSConfig{})
	restClient := wsClient.GetRestClient()

	return &TradeNode{
		BaseNode:     baseNode,
		cfg:          cfg,
		restClient:   restClient,
		wsClient:     wsClient,
		currentId:    0,
		nConnected:   0,
		nReconnected: 0,
		nError:       0,
		shutdownC:    make(chan struct{}),
		doneC:        make(chan struct{}),
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
	adapter, err := CreateAdapter(n.cfg.Exchange, n.cfg.Instrument)
	if err != nil {
		n.Logger().Error("Failed to create adapter", log.Error(err))
		return
	}
	unsubscribe, err := adapter.Subscribe(n.cfg.Symbol, TradeSubscriptionOptions{
		OnTrade: n.onTrade,
		OnError: func(err error) {
			n.nError++
			n.Logger().Error("Failed to subscribe to trade", log.Error(err))
		},
		OnDisconnect: func() {
			n.Logger().Info("Disconnected from trade")
		},
		OnConnect: func() {
			n.nConnected++
			n.Logger().Info("Connected to trade")
		},
		OnReconnect: func() {
			n.nReconnected++
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

func (n *TradeNode) onTrade(trade Trade) {
	subject, err := n.GetEmit(EmitTradeKey)
	if err != nil {
		n.Logger().Error("Failed to get emit subject", log.Error(err))
		return
	}
	side := app.Side_SIDE_BUY
	if trade.TakerSide == share.SideSell {
		side = app.Side_SIDE_SELL
	}
	appTrade := &app.Trade{
		Id:         trade.ID,
		Exchange:   n.cfg.ToExchange(),
		Instrument: n.cfg.ToInstrument(),
		Symbol:     n.cfg.ToAppSymbol(),
		Price:      trade.Price,
		Side:       side,
		Quantity:   trade.Qty,
		Timestamp:  trade.Time,
	}
	n.currentId = trade.ID
	n.EventBus().Emit(subject, appTrade)
	n.Logger().Infof("Emitting trade %d", appTrade.Id)
}

func (n *TradeNode) RequestParameters(req *pbCommon.ParametersRequest) *pbCommon.ParametersResponse {
	jsonBytes, err := json.Marshal(map[string]any{
		"exchange":   n.cfg.Exchange,
		"instrument": n.cfg.Instrument,
		"symbol":     n.cfg.Symbol,
	})
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
