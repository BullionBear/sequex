package main

import (
	"fmt"
	"log"
	"net"
	"time"

	"google.golang.org/grpc"
	"google.golang.org/protobuf/types/known/timestamppb"

	pb "github.com/BullionBear/sequex/pkg/protobuf/sequex" // Correct import path
)

// EventServiceServer implements the gRPC service
type SequexServiceServer struct {
	pb.UnimplementedSequexServiceServer
}

// StreamEvents streams mock events to the client
func (s *SequexServiceServer) StreamEvents(req *pb.Empty, stream pb.SequexService_StreamEventsServer) error {
	for i := 0; i < 10; i++ { // Stream 10 events
		event := &pb.Event{
			Id:        fmt.Sprintf("event-%d", i+1), // Generate a unique ID
			Type:      pb.EventType_KLINE_UPDATE,    // Set event type
			Source:    pb.EventSource_SEQUEX,        // Set event source
			CreatedAt: timestamppb.Now(),            // Set current timestamp
		}

		// Send the event to the client
		if err := stream.Send(event); err != nil {
			return err
		}

		// Simulate a delay between events
		time.Sleep(1 * time.Second)
	}

	return nil
}

func main() {
	// Start the gRPC server
	lis, err := net.Listen("tcp", ":50051")
	if err != nil {
		log.Fatalf("Failed to listen: %v", err)
	}

	grpcServer := grpc.NewServer()
	pb.RegisterSequexServiceServer(grpcServer, &SequexServiceServer{})

	log.Println("Server is running on port 50051")
	if err := grpcServer.Serve(lis); err != nil {
		log.Fatalf("Failed to serve: %v", err)
	}
}
