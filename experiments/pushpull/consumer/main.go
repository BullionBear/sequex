package main

import (
	"context"
	"encoding/json"
	"math/rand"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/BullionBear/sequex/pkg/log"
	"github.com/nats-io/nats.go"
)

const (
	streamName   = "sequex_orders"
	consumerName = "order_processor"
	subject      = "sequex.orders"
	url          = nats.DefaultURL
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
		log.String("component", "queue_consumer"),
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

	// Create consumer if it doesn't exist
	consumer, err := js.AddConsumer(streamName, &nats.ConsumerConfig{
		Name:          consumerName,
		FilterSubject: subject + ".new",
		AckPolicy:     nats.AckExplicitPolicy,
		DeliverPolicy: nats.DeliverNewPolicy,
		AckWait:       30 * time.Second,
		MaxDeliver:    3,
	})
	if err != nil {
		logger.Fatal("Failed to create consumer",
			log.String("consumer_name", consumerName),
			log.Error(err),
		)
	}

	logger.Info("Consumer created or already exists",
		log.String("consumer_name", consumerName),
		log.String("filter_subject", consumer.Config.FilterSubject),
	)

	// Create a context that can be cancelled
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	// Handle graceful shutdown
	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)

	go func() {
		<-sigChan
		logger.Info("Received shutdown signal, stopping consumer")
		cancel()
	}()

	// Subscribe to the consumer
	subscription, err := js.PullSubscribe(subject+".new", consumerName)
	if err != nil {
		logger.Fatal("Failed to create pull subscription",
			log.String("consumer_name", consumerName),
			log.Error(err),
		)
	}
	defer subscription.Unsubscribe()

	logger.Info("Subscribed to stream",
		log.String("stream_name", streamName),
		log.String("consumer_name", consumerName),
	)

	// Start consuming messages
	consumeMessages(ctx, subscription, logger)
}

func consumeMessages(ctx context.Context, subscription *nats.Subscription, logger log.Logger) {
	consumerLogger := logger.With(log.String("component", "order_consumer"))
	processedCount := 0

	for {
		select {
		case <-ctx.Done():
			consumerLogger.Info("Consumer stopped",
				log.Int("total_processed", processedCount),
			)
			return
		default:
			// Pull messages from the subscription
			messages, err := subscription.Fetch(1, nats.MaxWait(5*time.Second))
			if err != nil {
				if err == nats.ErrTimeout {
					// No messages available, continue
					continue
				}
				consumerLogger.Error("Failed to fetch messages",
					log.Error(err),
				)
				continue
			}

			for _, msg := range messages {
				processedCount++

				// Parse the order
				var order Order
				if err := json.Unmarshal(msg.Data, &order); err != nil {
					consumerLogger.Error("Failed to unmarshal order",
						log.String("message_data", string(msg.Data)),
						log.Error(err),
					)
					msg.Nak()
					continue
				}

				consumerLogger.Info("Received order",
					log.Int("order_id", order.ID),
					log.String("symbol", order.Symbol),
					log.String("side", order.Side),
					log.String("quantity", order.Quantity),
					log.String("price", order.Price),
					log.String("timestamp", order.Timestamp),
				)

				// Simulate order processing
				processingTime := time.Duration(rand.Intn(3000)+1000) * time.Millisecond
				consumerLogger.Info("Processing order",
					log.Int("order_id", order.ID),
					log.String("processing_time", processingTime.String()),
				)

				time.Sleep(processingTime)

				// Simulate processing success/failure
				if rand.Float64() < 0.95 { // 95% success rate
					// Acknowledge successful processing
					if err := msg.Ack(); err != nil {
						consumerLogger.Error("Failed to acknowledge message",
							log.Int("order_id", order.ID),
							log.Error(err),
						)
					} else {
						consumerLogger.Info("Order processed successfully",
							log.Int("order_id", order.ID),
							log.String("symbol", order.Symbol),
							log.String("side", order.Side),
						)
					}
				} else {
					// Simulate processing failure
					consumerLogger.Error("Order processing failed",
						log.Int("order_id", order.ID),
						log.String("symbol", order.Symbol),
					)

					// Negative acknowledgment - message will be redelivered
					if err := msg.Nak(); err != nil {
						consumerLogger.Error("Failed to negative acknowledge message",
							log.Int("order_id", order.ID),
							log.Error(err),
						)
					}
				}
			}
		}
	}
}
