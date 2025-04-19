package order

import (
	"github.com/BullionBear/sequex/internal/order/ordertype"
	"github.com/google/uuid"
)

type OrderManager struct {
	orders map[string]ordertype.Order
}

func NewOrderManager() *OrderManager {
	return &OrderManager{
		orders: make(map[string]ordertype.Order),
	}
}

func (om *OrderManager) Submit(order ordertype.Order) (string, error) {
	orderID := uuid.New().String()
	om.orders[orderID] = order
	return orderID, nil
}
