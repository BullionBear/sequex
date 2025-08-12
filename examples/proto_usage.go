package main

import (
	"fmt"

	"github.com/BullionBear/sequex/internal/model/protobuf/app/order"
	"github.com/BullionBear/sequex/internal/model/protobuf/app/shared"
)

func main() {
	// Create a user info using the shared type
	user := &shared.UserInfo{
		UserId:    "user123",
		Username:  "john_doe",
		CreatedAt: 1640995200, // 2022-01-01 00:00:00 UTC
	}

	// Create an order using the imported shared types
	orderMsg := &order.Order{
		OrderId:   "order456",
		Symbol:    "BTCUSDT",
		Price:     50000.0,
		Quantity:  0.1,
		Side:      shared.OrderSide_ORDER_SIDE_BUY,
		User:      user,
		Timestamp: 1640995200,
	}

	// Create an order response
	response := &order.OrderResponse{
		OrderId:   orderMsg.OrderId,
		Status:    order.OrderStatus_ORDER_STATUS_FILLED,
		Message:   "Order executed successfully",
		Timestamp: 1640995200,
	}

	// Print the data
	fmt.Printf("User: %s (%s)\n", user.Username, user.UserId)
	fmt.Printf("Order: %s %s %.2f @ %.2f\n",
		orderMsg.OrderId,
		orderMsg.Symbol,
		orderMsg.Quantity,
		orderMsg.Price)
	fmt.Printf("Side: %s\n", orderMsg.Side.String())
	fmt.Printf("Status: %s - %s\n", response.Status.String(), response.Message)

	// Demonstrate enum usage
	fmt.Printf("\nEnum values:\n")
	fmt.Printf("OrderSide BUY: %d\n", shared.OrderSide_ORDER_SIDE_BUY)
	fmt.Printf("OrderSide SELL: %d\n", shared.OrderSide_ORDER_SIDE_SELL)
	fmt.Printf("OrderStatus FILLED: %d\n", order.OrderStatus_ORDER_STATUS_FILLED)
}
