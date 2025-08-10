package main

import (
	"encoding/json"
	"log"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/BullionBear/sequex/internal/config"
	"github.com/BullionBear/sequex/pkg/node"
	"github.com/nats-io/nats.go"
)

func main() {
	log.Println("Starting NATS Microservices Demo")

	// Load merged configuration
	cfg, err := config.LoadConfig("configs/merged-config.yaml")
	if err != nil {
		log.Fatalf("Failed to load configuration: %v", err)
	}

	// Create a single NATS connection for the entire process
	nc, err := config.CreateNATSConnection(cfg.NATS.URL)
	if err != nil {
		log.Fatalf("Failed to connect to NATS: %v", err)
	}
	defer nc.Close()

	// Create deployer
	d := node.NewDeployer()

	// Create and register nodes based on configuration
	for nodeName, nodeConfig := range cfg.Deployer.Nodes {
		log.Printf("Creating node: %s", nodeName)

		node, err := config.CreateNode(nodeName, nodeConfig, nc)
		if err != nil {
			log.Fatalf("Failed to create node %s: %v", nodeName, err)
		}

		if err := d.RegisterNode(node); err != nil {
			log.Fatalf("Failed to register node %s: %v", nodeName, err)
		}

		log.Printf("Successfully registered node: %s", nodeName)
	}

	// Start all nodes
	for nodeName := range cfg.Deployer.Nodes {
		if err := d.Start(nodeName); err != nil {
			log.Fatalf("Failed to start node %s: %v", nodeName, err)
		}
		log.Printf("Successfully started node: %s", nodeName)
	}

	log.Println("All nodes deployed successfully")

	// Subscribe to results to see the system in action
	_, err = nc.Subscribe("sum.result", func(msg *nats.Msg) {
		var result map[string]interface{}
		if err := json.Unmarshal(msg.Data, &result); err != nil {
			log.Printf("Failed to unmarshal result: %v", err)
			return
		}
		log.Printf("Sum Result: Sum=%v, Count=%v, Last Number=%v",
			result["sum"], result["count"], result["last_number"])
	})
	if err != nil {
		log.Printf("Failed to subscribe to results: %v", err)
	}

	// Demonstrate RPC calls
	go demonstrateRPC(nc)

	// Wait for shutdown signal
	waitForShutdown()

	// Stop all nodes
	for nodeName := range cfg.Deployer.Nodes {
		if err := d.Stop(nodeName); err != nil {
			log.Printf("Error stopping node %s: %v", nodeName, err)
		}
	}

	log.Println("All nodes stopped")
}

func demonstrateRPC(nc *nats.Conn) {
	time.Sleep(10 * time.Second) // Wait for some activity

	// Query random generator count
	response, err := nc.Request("random.generator.count", nil, 5*time.Second)
	if err != nil {
		log.Printf("Failed to get random generator count: %v", err)
	} else {
		var result map[string]interface{}
		if err := json.Unmarshal(response.Data, &result); err != nil {
			log.Printf("Failed to unmarshal random generator response: %v", err)
		} else {
			log.Printf("Random Generator Stats: Count=%v", result["count"])
		}
	}

	// Query sum node stats
	response, err = nc.Request("sum.node.count", nil, 5*time.Second)
	if err != nil {
		log.Printf("Failed to get sum node stats: %v", err)
	} else {
		var result map[string]interface{}
		if err := json.Unmarshal(response.Data, &result); err != nil {
			log.Printf("Failed to unmarshal sum node response: %v", err)
		} else {
			log.Printf("Sum Node Stats: Count=%v, Sum=%v", result["count"], result["sum"])
		}
	}
}

func waitForShutdown() {
	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)

	log.Println("Waiting for shutdown signal...")
	<-sigChan
	log.Println("Shutdown signal received")
}
