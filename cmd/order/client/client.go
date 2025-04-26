package main

import (
	"context"
	"log"
	"time"

	pb "github.com/BullionBear/sequex/pkg/protobuf/order"
	pbdecimal "google.golang.org/genproto/googleapis/type/decimal"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

const (
	address = "localhost:50052"
)

func main() {
	// Set up a connection to the server.
	conn, err := grpc.NewClient(address, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Fatalf("did not connect: %v", err)
	}
	defer conn.Close()
	c := pb.NewBinanceOrderServiceClient(conn)

	// Contact the server and print out its response.
	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()
	r, err := c.PlaceLimitOrder(ctx, &pb.LimitOrderRequest{
		Account:  "scylla",
		Symbol:   "ADAUSDT",
		Side:     pb.Side_SELL,
		Quantity: &pbdecimal.Decimal{Value: "15"},
		Price:    &pbdecimal.Decimal{Value: "0.8"},
		Tif:      pb.TimeInForce_GTC,
	})
	if err != nil {
		log.Fatalf("Error send limit order %v", err)
	}
	log.Printf("Order ID: %+v", r)
}
