package main

import (
	"context"
	"fmt"
	"math/rand"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/BullionBear/sequex/pkg/log"
	"github.com/nats-io/nats.go"
)

const (
	streamName = "sequex_orders"
	subject    = "sequex.orders"
	url        = nats.DefaultURL
)

// Order represents a trading order
type Order struct {
	ID        int    `json:"id"`
	Symbol    string `json:"symbol"`
	Side      string `json:"side"`
	Quantity  string `json:"quantity"`
	Price     string `json:"price"`
	Timestamp string `json:"timestamp"`
}

func main() {
	// Initialize structured logger
	logger := log.New(
		log.WithLevel(log.LevelInfo),
		log.WithEncoder(log.NewTextEncoder()),
		log.WithOutput(os.Stdout),
	)

	// Connect to NATS
	nc, err := nats.Connect(url)
	if err != nil {
		logger.Fatal("Failed to connect to NATS",
			log.String("nats_url", url),
			log.Error(err),
		)
	}
	defer nc.Close()

	logger.Info("Connected to NATS",
		log.String("nats_url", url),
		log.String("component", "queue_producer"),
	)

	// Create JetStream context
	js, err := nc.JetStream()
	if err != nil {
		logger.Fatal("Failed to create JetStream context",
			log.Error(err),
		)
	}

	logger.Info("JetStream context created")

	// Create stream if it doesn't exist
	stream, err := js.AddStream(&nats.StreamConfig{
		Name:      streamName,
		Subjects:  []string{subject + ".*"},
		Storage:   nats.FileStorage,
		Retention: nats.LimitsPolicy,
		MaxAge:    24 * time.Hour,
		MaxMsgs:   10000,
	})
	if err != nil {
		logger.Fatal("Failed to create stream",
			log.String("stream_name", streamName),
			log.Error(err),
		)
	}

	logger.Info("Stream created or already exists",
		log.String("stream_name", streamName),
		log.Int("subjects", len(stream.Config.Subjects)),
	)

	// Create a context that can be cancelled
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	// Handle graceful shutdown
	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)

	go func() {
		<-sigChan
		logger.Info("Received shutdown signal, stopping producer")
		cancel()
	}()

	// Start producing orders
	produceOrders(ctx, js, logger)
}

func produceOrders(ctx context.Context, js nats.JetStreamContext, logger log.Logger) {
	ticker := time.NewTicker(2 * time.Second)
	defer ticker.Stop()

	orderCount := 0
	symbols := []string{"BTCUSDT", "ETHUSDT", "ADAUSDT", "DOTUSDT", "LINKUSDT"}
	sides := []string{"BUY", "SELL"}
	producerLogger := logger.With(log.String("component", "order_producer"))

	for {
		select {
		case <-ctx.Done():
			producerLogger.Info("Producer stopped",
				log.Int("total_orders", orderCount),
			)
			return
		case <-ticker.C:
			orderCount++

			// Create a random order
			order := Order{
				ID:        orderCount,
				Symbol:    symbols[rand.Intn(len(symbols))],
				Side:      sides[rand.Intn(len(sides))],
				Quantity:  fmt.Sprintf("%.4f", rand.Float64()*10+0.1),
				Price:     fmt.Sprintf("%.2f", rand.Float64()*50000+1000),
				Timestamp: time.Now().Format(time.RFC3339),
			}

			// Convert order to JSON
			orderData := fmt.Sprintf(`{"id":%d,"symbol":"%s","side":"%s","quantity":"%s","price":"%s","timestamp":"%s"}`,
				order.ID, order.Symbol, order.Side, order.Quantity, order.Price, order.Timestamp)

			// Publish to JetStream
			ack, err := js.Publish(subject+".new", []byte(orderData))
			if err != nil {
				producerLogger.Error("Failed to publish order",
					log.Int("order_id", order.ID),
					log.String("subject", subject+".new"),
					log.String("order", orderData),
					log.Error(err),
				)
				continue
			}

			producerLogger.Info("Order published successfully",
				log.Int("order_id", order.ID),
				log.String("symbol", order.Symbol),
				log.String("side", order.Side),
				log.String("quantity", order.Quantity),
				log.String("price", order.Price),
				log.String("stream_sequence", fmt.Sprintf("%d", ack.Sequence)),
				log.String("stream_name", ack.Stream),
			)
		}
	}
}
