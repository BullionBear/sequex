package trade

import (
	"github.com/BullionBear/sequex/internal/model/protobuf/app"
	"github.com/BullionBear/sequex/internal/nodeimpl/base"
	"github.com/BullionBear/sequex/pkg/eventbus"
	"github.com/BullionBear/sequex/pkg/exchange/binance"
	"github.com/BullionBear/sequex/pkg/log"
	"github.com/BullionBear/sequex/pkg/node"
)

const (
	EmitTradeKey = "emit_trade"
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
		shutdownC: make(chan struct{}),
		doneC:     make(chan struct{}),
	}, nil
}

func (n *TradeNode) Start() error {
	n.Logger().Info("Starting Trade node")
	return nil
}

func (n *TradeNode) Shutdown() error {

	return nil
}

func (n *TradeNode) emitTrade(shutdownC chan struct{}, doneC chan struct{}) {
	subject, err := n.GetEmit(EmitTradeKey)
	if err != nil {
		n.Logger().Error("Failed to get emit subject", log.Error(err))
		return
	}
	n.wsClient.SubscribeTrade(n.cfg.Symbol, binance.TradeSubscriptionOptions{
		OnTrade: func(trade binance.WSTrade) {
			appTrade := &app.Trade{
				Symbol:     trade.Symbol,
				Price:      trade.Price,
				Quantity:   trade.Quantity,
				Timestamp:  trade.TradeTime,
				Exchange:   n.cfg.Exchange,
				Instrument: n.cfg.Instrument,
				Side:       app.Side_SIDE_BUY,
			}
			subject.Emit(appTrade)
		},
	})
}

/*
type PublicTradeNode struct {
	*base.BaseNode
	// Configurable parameters
	cfg PublicTradeConfig

	SubscribeTrade(market string, symbol string, options TradeSubscriptionOptions)
}
*/
