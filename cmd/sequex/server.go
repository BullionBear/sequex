package main

import (
	"context"
	"flag"
	"fmt"
	"log"
	"net"
	"time"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/protobuf/types/known/timestamppb"

	"encoding/json"

	"github.com/BullionBear/sequex/internal/config"
	"github.com/BullionBear/sequex/internal/payload"
	"github.com/BullionBear/sequex/internal/strategy/solvexity"
	"github.com/BullionBear/sequex/internal/tradingpipe"
	"github.com/BullionBear/sequex/pkg/message"
	"github.com/BullionBear/sequex/pkg/mq"
	"github.com/BullionBear/sequex/pkg/mq/inprocq"
	pbSequex "github.com/BullionBear/sequex/pkg/protobuf/sequex"       // Correct import path
	pbSolvexity "github.com/BullionBear/sequex/pkg/protobuf/solvexity" // Correct import path
)

// EventServiceServer implements the gRPC service
type server struct {
	pbSequex.UnimplementedSequexServiceServer
	q        mq.MessageQueue
	pipeline tradingpipe.TradingPipeline
}

// StreamEvents streams mock events to the client
func (s *server) OnEvent(ctx context.Context, in *pbSequex.Event) (*pbSequex.Ack, error) {
	msg := message.Message{
		ID:        in.Id,
		Type:      in.Type.String(),
		Source:    in.Source.String(),
		CreatedAt: in.CreatedAt.Seconds,
		Payload:   in.Payload,
	}
	s.q.Send(&msg)
	return &pbSequex.Ack{
		Id:         in.Id,
		ReceivedAt: timestamppb.New(time.Now()),
	}, nil
}

func main() {
	// -c flag to specify the configuration file
	// Define a string flag with a default value and a description
	path := flag.String("c", "default.conf", "Path to the configuration file")

	// Parse command-line flags
	flag.Parse()

	// Use the flag value
	fmt.Println("Config file:", *path)
	q := inprocq.NewInprocQueue()
	conf := config.NewDomain(*path)
	name := conf.GetConfig().Name
	// Create a new strategy
	// Set up a connection to the server.
	conn, err := grpc.NewClient(fmt.Sprintf("%s:%d", conf.GetConfig().Solvexity.Host, conf.GetConfig().Solvexity.Port),
		grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Fatalf("did not connect: %v", err)
	}
	defer conn.Close()
	client := pbSolvexity.NewSolvexityClient(conn)
	strategy := solvexity.NewSolvexity(client)
	pipeline := tradingpipe.NewTradingPipeline(name, strategy)
	// Start the gRPC server
	lis, err := net.Listen("tcp", fmt.Sprintf(":%d", conf.GetConfig().Sequex.Port))
	if err != nil {
		log.Fatalf("Failed to listen: %v", err)
	}

	grpcServer := grpc.NewServer()
	pbSequex.RegisterSequexServiceServer(grpcServer, &server{
		q:        q,
		pipeline: *pipeline,
	})
	go func() {
		for {
			msg, err := q.RecvTimeout(2 * time.Second)
			if err != nil {
				log.Printf("Error receiving message: %v\n", err)
			} else {
				log.Printf("Received message: %v\n", msg)
				switch msg.Type {
				case "KLINE_UPDATE":
					var payload payload.KLineUpdate
					err = json.Unmarshal(msg.Payload, &payload)
					if err != nil {
						log.Printf("Error unmarshalling payload: %v\n", err)
						continue
					}
					pipeline.OnKLineUpdate(payload)
				default:
					fmt.Printf("Unknown message type: %v\n", msg.Type)
				}
			}
		}
	}()

	log.Println("Server is running on port 50051")
	if err := grpcServer.Serve(lis); err != nil {
		log.Fatalf("Failed to serve: %v", err)
	}
}
