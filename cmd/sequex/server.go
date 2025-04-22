package main

import (
	"flag"
	"fmt"
	"log"

	"google.golang.org/protobuf/types/known/timestamppb"

	"github.com/BullionBear/sequex/pkg/mq"
	pbSequex "github.com/BullionBear/sequex/pkg/protobuf/sequex"
)

// EventServiceServer implements the gRPC service
type SequexServer struct {
	pbSequex.UnimplementedSequexServiceServer
}

func NewSequexServer(q mq.MessageQueue) *SequexServer {
	return &SequexServer{}
}

// StreamEvents streams mock events to the client
func (s *SequexServer) OnEvent(stream pbSequex.SequexService_OnEventServer) error {
	sendChan := make(chan *pbSequex.Event, 100)
	defer close(sendChan)

	go func() {
		for event := range sendChan {
			if err := stream.Send(event); err != nil {
				log.Printf("Error sending event: %v", err)
			}
		}
	}()

	for {
		event, err := stream.Recv()
		if err != nil {
			return err
		}

		log.Printf("Received event: %s (%s) from %s", event.Type, event.Id, event.Source)
		switch event.Type {
		case pbSequex.EventType_ORDER_UPDATE:
			sendChan <- &pbSequex.Event{
				Id:        event.Id,
				Type:      pbSequex.EventType_ORDER_ACK,
				Source:    pbSequex.EventSource_SEQUEX,
				CreatedAt: timestamppb.Now(),
				Payload:   []byte("Order update received"),
			}
		case pbSequex.EventType_EXECUTION_UPDATE:
			sendChan <- &pbSequex.Event{
				Id:        event.Id,
				Type:      pbSequex.EventType_EXECUTION_ACK,
				Source:    pbSequex.EventSource_SEQUEX,
				CreatedAt: timestamppb.Now(),
				Payload:   []byte("Execution update received"),
			}
		default:
			log.Printf("Undefine event type: %s", event.Type)
		}
	}
}

func main() {
	// -c flag to specify the configuration file
	path := flag.String("c", "default.conf", "Path to the configuration file")
	flag.Parse()

	// Use the flag value
	fmt.Println("Config file:", *path)
	/*
		conf := config.NewDomain(*path)
		// Resource

		// Set up trading pipeline.
		conn, err := grpc.NewClient(fmt.Sprintf("%s:%d", conf.GetConfig().Solvexity.Host, conf.GetConfig().Solvexity.Port),
			grpc.WithTransportCredentials(insecure.NewCredentials()))
		if err != nil {
			log.Fatalf("did not connect: %v", err)
		}
		defer conn.Close()

		// Start the gRPC server
		lis, err := net.Listen("tcp", fmt.Sprintf(":%d", conf.GetConfig().Sequex.Port))
		if err != nil {
			log.Fatalf("Failed to listen: %v", err)
		}
		grpcServer := grpc.NewServer()
		pbSequex.RegisterSequexServiceServer(grpcServer, &SequexServer{})
		log.Println("Server is running on port 50051")
		if err := grpcServer.Serve(lis); err != nil {
			log.Fatalf("Failed to serve: %v", err)
		}
	*/
}
