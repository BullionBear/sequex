package order

import (
	"context"
	"fmt"

	"github.com/BullionBear/sequex/internal/order/ordertype"
	"github.com/adshao/go-binance/v2"
	"github.com/adshao/go-binance/v2/futures"
	"github.com/google/uuid"
)

type OrderManager struct {
	spotClient *binance.Client
	perpClient *futures.Client
	orders     map[string]ordertype.Order
}

func NewOrderManager(apiKey, apiSecret string) *OrderManager {
	return &OrderManager{
		spotClient: binance.NewClient(apiKey, apiSecret),
		perpClient: futures.NewClient(apiKey, apiSecret),
		orders:     make(map[string]ordertype.Order),
	}
}

func (om *OrderManager) Submit(order ordertype.Order) (string, error) {
	orderID := uuid.New().String()
	switch o := order.(type) {
	case *ordertype.MarketOrder:
		// Handle market order submission
		if o.Instrument == ordertype.InstrumentSpot { // Example: set instrument to SPOT
			response, err := om.spotClient.NewCreateOrderService().
				Symbol(o.Symbol).
				Side(binance.SideType(o.Side.String())).
				Type(binance.OrderTypeMarket).
				Quantity(o.Quantity.String()).
				Do(context.Background())
			if err != nil {
				return "", fmt.Errorf("failed to submit market order: %w", err)
			}
			return fmt.Sprintf("%d", response.OrderID), nil
		}
	case *ordertype.LimitOrder:
		// Handle limit order submission
	case *ordertype.StopMarketOrder:
		// Handle stop market order submission
	case *ordertype.TrailingStopMarketOrder:
		// Handle trailing stop market order submission
	case *ordertype.OneCancelsOtherOrder:
		// Handle OCO order submission
	case *ordertype.IfDoneOrder:
		// Handle IFDO order submission
	default:
		return "", fmt.Errorf("unsupported order type: %T", o)
	}
	om.orders[orderID] = order
	return orderID, nil
}
