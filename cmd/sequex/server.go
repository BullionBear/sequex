package main

import (
	"context"
	"log"
	"net"
	"time"

	"google.golang.org/grpc"
	"google.golang.org/protobuf/types/known/timestamppb"

	"github.com/BullionBear/sequex/internal/strategy/sequex"
	"github.com/BullionBear/sequex/internal/tradingpipe"
	"github.com/BullionBear/sequex/pkg/mq"
	"github.com/BullionBear/sequex/pkg/mq/inprocq"
	pb "github.com/BullionBear/sequex/pkg/protobuf/sequex" // Correct import path
)

// EventServiceServer implements the gRPC service
type server struct {
	pb.UnimplementedSequexServiceServer
	q        mq.MessageQueue
	pipeline tradingpipe.TradingPipeline
}

// StreamEvents streams mock events to the client
func (s *server) OnEvent(ctx context.Context, in *pb.Event) (*pb.Ack, error) {

	return &pb.Ack{
		Id:         in.Id,
		ReceivedAt: timestamppb.New(time.Now()),
	}, nil
}

func main() {
	q := inprocq.NewInprocQueue()
	name := "Sequex"
	strategy := sequex.NewSequex()
	pipeline := tradingpipe.NewTradingPipeline(name, strategy)
	// Start the gRPC server
	lis, err := net.Listen("tcp", ":50051")
	if err != nil {
		log.Fatalf("Failed to listen: %v", err)
	}

	grpcServer := grpc.NewServer()
	pb.RegisterSequexServiceServer(grpcServer, &server{
		q:        q,
		pipeline: *pipeline,
	})

	log.Println("Server is running on port 50051")
	if err := grpcServer.Serve(lis); err != nil {
		log.Fatalf("Failed to serve: %v", err)
	}
}
