package main

import (
	"context"
	"encoding/json"
	"log"
	"sync"
	"time"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/protobuf/types/known/timestamppb"

	"github.com/BullionBear/sequex/internal/feed/binance"
	"github.com/BullionBear/sequex/internal/payload"
	pbSequex "github.com/BullionBear/sequex/pkg/protobuf/sequex" // Correct import path)
	"github.com/google/uuid"
)

func main() {
	symbol := "BTCUSDT"
	conn, err := grpc.NewClient("localhost:50051", grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Fatalf("did not connect: %v", err)
	}
	defer conn.Close()

	client := pbSequex.NewSequexServiceClient(conn)
	stream, err := client.OnEvent(context.Background())
	if err != nil {
		log.Fatalf("error creating stream: %v", err)
	}

	// Use wait group to properly wait for goroutines
	var wg sync.WaitGroup
	wg.Add(1)

	go func() {
		receiveEvents(stream)
		wg.Done()
	}()

	feed := binance.NewBinanceFeed()
	unsubscribe, err := feed.SubscribeKlineUpdate(symbol, func(klineUpdate *payload.KLineUpdate) {
		payload, err := json.Marshal(klineUpdate)
		if err != nil {
			log.Printf("Error marshalling payload: %v\n", err)
			return
		}
		sendEvent(stream, pbSequex.EventType_KLINE_UPDATE, payload)
	})
	if err != nil {
		log.Fatalf("Failed to subscribe to KlineUpdate: %v", err)
	}
	defer unsubscribe()

	time.Sleep(20 * time.Second)
	// Proper close send
	if err := stream.CloseSend(); err != nil {
		log.Printf("Error closing send: %v", err)
	}

	// Wait for receiver to complete
	wg.Wait()
}

func sendEvent(stream pbSequex.SequexService_OnEventClient, eventType pbSequex.EventType, payload []byte) {
	event := &pbSequex.Event{
		Id:        uuid.New().String(),
		Type:      eventType,
		Source:    pbSequex.EventSource_SEQUEX,
		CreatedAt: timestamppb.Now(),
		Payload:   payload,
	}

	if err := stream.Send(event); err != nil {
		log.Printf("error sending event: %v", err)
		return
	}
	log.Printf("Sent event: %s (%s)", event.Type, event.Id)
}

func receiveEvents(stream pbSequex.SequexService_OnEventClient) {
	for {
		event, err := stream.Recv()
		if err != nil {
			log.Printf("Error receiving event: %v", err)
			return // Exit on any error to avoid infinite loop
		}
		log.Printf("Received response: %s for event %s", event.Type, event.Id)
	}
}
