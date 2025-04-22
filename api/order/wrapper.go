package api

import (
	"github.com/BullionBear/sequex/internal/order"
	pb "github.com/BullionBear/sequex/pkg/protobuf/order" // Correct import path
)

func convertSide(side pb.Side) order.Side {
	switch side {
	case pb.Side_BUY:
		return order.SideBuy
	case pb.Side_SELL:
		return order.SideSell
	default:
		return order.SideUnknown
	}
}

func convertTimeInForce(tif pb.TimeInForce) order.TimeInForce {
	switch tif {
	case pb.TimeInForce_GTC:
		return order.TimeInForceGTC
	case pb.TimeInForce_IOC:
		return order.TimeInForceIOC
	case pb.TimeInForce_FOK:
		return order.TimeInForceFOK
	default:
		return order.TimeInForceUnknown
	}
}
