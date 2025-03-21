package main

import (
	"flag"
	"fmt"
	"log"
	"net"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/protobuf/types/known/timestamppb"

	"encoding/json"

	"github.com/BullionBear/sequex/internal/config"
	"github.com/BullionBear/sequex/internal/payload"
	"github.com/BullionBear/sequex/internal/strategy"
	"github.com/BullionBear/sequex/internal/strategy/solvexity"
	"github.com/BullionBear/sequex/pkg/mq"
	pbSequex "github.com/BullionBear/sequex/pkg/protobuf/sequex"       // Correct import path
	pbSolvexity "github.com/BullionBear/sequex/pkg/protobuf/solvexity" // Correct import path
)

// EventServiceServer implements the gRPC service
type SequexServer struct {
	pbSequex.UnimplementedSequexServiceServer
	st strategy.Strategy
}

func NewSequexServer(q mq.MessageQueue, st strategy.Strategy) *SequexServer {
	return &SequexServer{
		st: st,
	}
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
		case pbSequex.EventType_KLINE_UPDATE:
			sendChan <- &pbSequex.Event{
				Id:        event.Id,
				Type:      pbSequex.EventType_KLINE_ACK,
				Source:    pbSequex.EventSource_SEQUEX,
				CreatedAt: timestamppb.Now(),
				Payload:   []byte("Kline update received"),
			}
			var payload payload.KLineUpdate
			err := json.Unmarshal(event.Payload, &payload)
			if err != nil {
				log.Printf("Error unmarshalling payload: %v\n", err)
				continue
			}
			err = s.st.OnKLineUpdate(payload)
			if err != nil {
				log.Printf("Error processing KlineUpdate: %v\n", err)
				sendChan <- &pbSequex.Event{
					Id:        event.Id,
					Type:      pbSequex.EventType_KLINE_FAILED,
					Source:    pbSequex.EventSource_SEQUEX,
					CreatedAt: timestamppb.Now(),
					Payload:   []byte(err.Error()),
				}
				continue
			}
			sendChan <- &pbSequex.Event{
				Id:        event.Id,
				Type:      pbSequex.EventType_KLINE_FINISHED,
				Source:    pbSequex.EventSource_SEQUEX,
				CreatedAt: timestamppb.Now(),
				Payload:   []byte("Kline update processed"),
			}
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
	conf := config.NewDomain(*path)
	// Resource

	// Set up trading pipeline.
	conn, err := grpc.NewClient(fmt.Sprintf("%s:%d", conf.GetConfig().Solvexity.Host, conf.GetConfig().Solvexity.Port),
		grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Fatalf("did not connect: %v", err)
	}
	defer conn.Close()
	client := pbSolvexity.NewSolvexityClient(conn)
	strategy := solvexity.NewSolvexity(client)

	// Start the gRPC server
	lis, err := net.Listen("tcp", fmt.Sprintf(":%d", conf.GetConfig().Sequex.Port))
	if err != nil {
		log.Fatalf("Failed to listen: %v", err)
	}
	grpcServer := grpc.NewServer()
	pbSequex.RegisterSequexServiceServer(grpcServer, &SequexServer{
		st: strategy,
	})
	// event source
	/*
		feed := binance.NewBinanceFeed()
		unsubscribe, err := feed.SubscribeKlineUpdate(symbol, func(klineUpdate *payload.KLineUpdate) {
			payload, err := json.Marshal(klineUpdate)
			if err != nil {
				log.Printf("Error marshalling payload: %v\n", err)
				return
			}
			msg := message.Message{
				ID:        "1",
				Type:      "KLINE_UPDATE",
				Source:    "BINANCE",
				CreatedAt: time.Now().Unix(),
				Payload:   payload,
			}
			q.Send(&msg)
		})
		if err != nil {
			log.Fatalf("Failed to subscribe to KlineUpdate: %v", err)
		}
		defer unsubscribe()
	*/

	log.Println("Server is running on port 50051")
	if err := grpcServer.Serve(lis); err != nil {
		log.Fatalf("Failed to serve: %v", err)
	}
}
