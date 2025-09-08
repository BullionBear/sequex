package main

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/nats-io/nats.go"
)

func main() {
	fmt.Println("🚀 Starting JetStream Example")
	fmt.Println("==============================")
	fmt.Println("💡 Make sure to run './script.sh create' first to set up the stream and consumer")
	fmt.Println()

	// Connect to NATS server
	nc, err := nats.Connect("nats://localhost:4222")
	if err != nil {
		log.Fatal("Failed to connect to NATS:", err)
	}
	defer nc.Close()
	fmt.Println("✅ Connected to NATS server")

	// Create JetStream context
	js, err := nc.JetStream()
	if err != nil {
		log.Fatal("Failed to create JetStream context:", err)
	}
	fmt.Println("✅ JetStream context created")

	// Verify stream exists (should be created by script.sh)
	streamName := "TEST_STREAM"
	_, err = js.StreamInfo(streamName)
	if err != nil {
		log.Fatal("❌ Stream 'TEST_STREAM' does not exist. Please run './script.sh create' first to set up the JetStream environment.")
	}
	fmt.Printf("✅ Stream '%s' found\n", streamName)

	// Verify consumer exists (should be created by script.sh)
	consumerName := "TEST_CONSUMER"
	_, err = js.ConsumerInfo(streamName, consumerName)
	if err != nil {
		log.Fatal("❌ Consumer 'TEST_CONSUMER' does not exist. Please run './script.sh create' first to set up the JetStream environment.")
	}
	fmt.Printf("✅ Consumer '%s' found\n", consumerName)

	// Publish some test messages
	fmt.Println("\n📤 Publishing test messages...")
	for i := 1; i <= 5; i++ {
		subject := fmt.Sprintf("test.message.%d", i)
		message := fmt.Sprintf("Hello JetStream! Message #%d - %s", i, time.Now().Format(time.RFC3339))

		_, err := js.Publish(subject, []byte(message))
		if err != nil {
			log.Printf("Failed to publish message %d: %v", i, err)
			continue
		}
		fmt.Printf("  📨 Published to %s: %s\n", subject, message)
		time.Sleep(100 * time.Millisecond) // Small delay between messages
	}

	// Subscribe and consume messages
	fmt.Println("\n📥 Consuming messages...")
	sub, err := js.PullSubscribe("test.>", consumerName)
	if err != nil {
		log.Fatal("Failed to create pull subscription:", err)
	}
	defer sub.Unsubscribe()

	// Fetch messages with timeout
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	msgs, err := sub.Fetch(10, nats.Context(ctx))
	if err != nil {
		if err == context.DeadlineExceeded {
			fmt.Println("⏰ Timeout waiting for messages")
		} else {
			log.Printf("Failed to fetch messages: %v", err)
		}
	} else {
		fmt.Printf("📬 Received %d messages:\n", len(msgs))
		for i, msg := range msgs {
			fmt.Printf("  %d. Subject: %s\n", i+1, msg.Subject)
			fmt.Printf("     Data: %s\n", string(msg.Data))
			fmt.Printf("     Timestamp: %s\n", msg.Header.Get("Nats-Time-Stamp"))

			// Acknowledge the message
			if err := msg.Ack(); err != nil {
				log.Printf("Failed to ack message: %v", err)
			} else {
				fmt.Printf("     ✅ Acknowledged\n")
			}
			fmt.Println()
		}
	}

	// Display stream and consumer information
	fmt.Println("📊 Stream Information:")
	streamInfo, err := js.StreamInfo(streamName)
	if err != nil {
		log.Printf("Failed to get stream info: %v", err)
	} else {
		fmt.Printf("  Name: %s\n", streamInfo.Config.Name)
		fmt.Printf("  Subjects: %v\n", streamInfo.Config.Subjects)
		fmt.Printf("  Storage: %s\n", streamInfo.Config.Storage)
		fmt.Printf("  Retention: %s\n", streamInfo.Config.Retention)
		fmt.Printf("  Messages: %d\n", streamInfo.State.Msgs)
		fmt.Printf("  Bytes: %d\n", streamInfo.State.Bytes)
	}

	fmt.Println("\n📊 Consumer Information:")
	consumerInfo, err := js.ConsumerInfo(streamName, consumerName)
	if err != nil {
		log.Printf("Failed to get consumer info: %v", err)
	} else {
		fmt.Printf("  Name: %s\n", consumerInfo.Name)
		fmt.Printf("  Durable: %s\n", consumerInfo.Config.Durable)
		fmt.Printf("  Deliver Policy: %v\n", consumerInfo.Config.DeliverPolicy)
		fmt.Printf("  Ack Policy: %v\n", consumerInfo.Config.AckPolicy)
		fmt.Printf("  Num Pending: %d\n", consumerInfo.NumPending)
		fmt.Printf("  Num Delivered: %d\n", consumerInfo.Delivered.Consumer)
	}

	fmt.Println("\n🎉 JetStream example completed successfully!")
	fmt.Println("💡 Use './script.sh clean' to clean up the test stream and consumer")
}
