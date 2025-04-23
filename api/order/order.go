package orderapi

import (
	"context"

	"github.com/BullionBear/sequex/internal/order"
	pb "github.com/BullionBear/sequex/pkg/protobuf/order" // Correct import path
	"github.com/google/uuid"
	"github.com/shopspring/decimal"
)

type BinaceOrderService struct {
	pb.UnimplementedBinanceOrderServiceServer
	orderManger *order.BinanceOrderManager
}

func NewBinanceOrderService(orderManager *order.BinanceOrderManager) *BinaceOrderService {
	return &BinaceOrderService{
		orderManger: orderManager,
	}
}

func (s *BinaceOrderService) PlaceMarketOrder(ctx context.Context, req *pb.MarketOrderRequest) (*pb.OrderResponse, error) {
	qty, _ := decimal.NewFromString(req.Quantity.String())
	order := order.MarketOrder{
		ClientOrderID: uuid.New().String(),
		Symbol:        req.Symbol,
		Side:          convertSide(req.Side),
		Quantity:      qty,
	}
	orderID, err := s.orderManger.MarketOrder(order)
	if err != nil {
		return nil, err
	}
	return &pb.OrderResponse{
		OrderId: orderID,
	}, nil
}

func (s *BinaceOrderService) PlaceLimitOrder(ctx context.Context, req *pb.LimitOrderRequest) (*pb.OrderResponse, error) {
	qty, _ := decimal.NewFromString(req.Quantity.String())
	price, _ := decimal.NewFromString(req.Price.String())
	order := order.LimitOrder{
		ClientOrderID: uuid.New().String(),
		Symbol:        req.Symbol,
		Side:          order.Side(req.Side),
		Quantity:      qty,
		Price:         price,
		TimeInForce:   convertTimeInForce(req.Tif),
	}
	orderID, err := s.orderManger.LimitOrder(order)
	if err != nil {
		return nil, err
	}
	return &pb.OrderResponse{
		OrderId: orderID,
	}, nil
}

func (s *BinaceOrderService) PlaceStopMarketOrder(ctx context.Context, req *pb.StopMarketOrderRequest) (*pb.OrderResponse, error) {
	qty, _ := decimal.NewFromString(req.Quantity.String())
	stopPrice, _ := decimal.NewFromString(req.StopPrice.String())
	order := order.StopMarketOrder{
		ClientOrderID: uuid.New().String(),
		Symbol:        req.Symbol,
		Side:          convertSide(req.Side),
		Quantity:      qty,
		StopPrice:     stopPrice,
	}
	orderID, err := s.orderManger.StopMarketOrder(order)
	if err != nil {
		return nil, err
	}
	return &pb.OrderResponse{
		OrderId: orderID,
	}, nil
}
