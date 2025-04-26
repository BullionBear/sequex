package orderapi

import (
	"context"

	"github.com/BullionBear/sequex/pkg/log"

	"github.com/BullionBear/sequex/internal/order"
	pb "github.com/BullionBear/sequex/pkg/protobuf/order" // Correct import path
	"github.com/shopspring/decimal"
)

type BinaceOrderService struct {
	pb.UnimplementedBinanceOrderServiceServer
	orderManger *order.BinanceOrderManager
	logger      *log.Logger
}

func NewBinanceOrderService(orderManager *order.BinanceOrderManager, logger *log.Logger) *BinaceOrderService {
	return &BinaceOrderService{
		orderManger: orderManager,
		logger:      logger,
	}
}

func (s *BinaceOrderService) PlaceMarketOrder(ctx context.Context, req *pb.MarketOrderRequest) (*pb.OrderResponse, error) {
	s.logger.Info("Receive MarketOrderRequest %+v", req)
	qty, err := decimal.NewFromString(req.Quantity.String())
	if err != nil {
		s.logger.Error("Error parsing quantity '%s': %s", req.Quantity.String(), err)
		return nil, err
	}
	orderResp, err := s.orderManger.PlaceMarketOrder(req.Account, req.Symbol, qty)
	if err != nil {
		return nil, err
	}
	return &pb.OrderResponse{
		SequexId: orderResp.SequexID,
		Status:   orderResp.Status.String(),
	}, nil
}

func (s *BinaceOrderService) PlaceLimitOrder(ctx context.Context, req *pb.LimitOrderRequest) (*pb.OrderResponse, error) {
	s.logger.Info("Receive LimitOrderRequest %+v", req)
	qtyStr := req.Quantity.GetValue()
	qty, err := decimal.NewFromString(qtyStr)
	if err != nil {
		s.logger.Error("Error parsing quantity '%s': %s", qtyStr, err)
		return nil, err
	}
	priceStr := req.Price.GetValue()
	price, err := decimal.NewFromString(priceStr)
	if err != nil {
		s.logger.Error("Error parsing price '%s': %s", priceStr, err)
		return nil, err
	}
	orderResp, err := s.orderManger.PlaceLimitOrder(req.Account, req.Symbol, qty, price)
	if err != nil {
		return nil, err
	}
	return &pb.OrderResponse{
		SequexId: orderResp.SequexID,
		Status:   orderResp.Status.String(),
	}, nil
}

func (s *BinaceOrderService) PlaceStopMarketOrder(ctx context.Context, req *pb.StopMarketOrderRequest) (*pb.OrderResponse, error) {
	s.logger.Info("Receive StopMarketOrderRequest %+v", req)
	qtyStr := req.Quantity.GetValue()
	qty, err := decimal.NewFromString(qtyStr)
	if err != nil {
		s.logger.Error("Error parsing quantity '%s': %s", qtyStr, err)
		return nil, err
	}
	stopPriceStr := req.StopPrice.GetValue()
	stopPrice, err := decimal.NewFromString(stopPriceStr)
	orderID, err := s.orderManger.PlaceStopMarketOrder(req.Account, req.Symbol, qty, stopPrice)
	if err != nil {
		return nil, err
	}

	return &pb.OrderResponse{
		SequexId: orderID.SequexID,
		Status:   orderID.Status.String(),
	}, nil
}
