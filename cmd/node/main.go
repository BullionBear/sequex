package main

import (
	"log"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/BullionBear/sequex/internal/config"
	rngpb "github.com/BullionBear/sequex/internal/model/protobuf/example"
	_ "github.com/BullionBear/sequex/internal/nodeimpl/example/rng" // Import to register RNG node
	"github.com/BullionBear/sequex/pkg/node"
	"github.com/nats-io/nats.go"
	"google.golang.org/protobuf/proto"
)

func main() {
	log.Println("Starting NATS Microservices Demo")

	// Load merged configuration
	cfg, err := config.LoadConfig("config/rng.yml")
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

	// Subscribe to RNG messages to see the system in action
	_, err = nc.Subscribe("rng_positive.rng.RngMessage", func(msg *nats.Msg) {
		var rngMsg rngpb.RngMessage
		if err := proto.Unmarshal(msg.Data, &rngMsg); err != nil {
			log.Printf("Failed to unmarshal RNG message: %v", err)
			return
		}
		log.Printf("RNG Positive: %d", rngMsg.Random)
	})
	if err != nil {
		log.Printf("Failed to subscribe to RNG positive messages: %v", err)
	}

	_, err = nc.Subscribe("rng_negative.rng.RngMessage", func(msg *nats.Msg) {
		var rngMsg rngpb.RngMessage
		if err := proto.Unmarshal(msg.Data, &rngMsg); err != nil {
			log.Printf("Failed to unmarshal RNG message: %v", err)
			return
		}
		log.Printf("RNG Negative: %d", rngMsg.Random)
	})
	if err != nil {
		log.Printf("Failed to subscribe to RNG negative messages: %v", err)
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

	// Query RNG positive node count
	req := &rngpb.RngCountRequest{}
	reqBytes, err := proto.Marshal(req)
	if err != nil {
		log.Printf("Failed to marshal RNG count request: %v", err)
		return
	}

	response, err := nc.Request("rpc.rng_positive", reqBytes, 5*time.Second)
	if err != nil {
		log.Printf("Failed to get RNG positive count: %v", err)
	} else {
		var resp rngpb.RngCountResponse
		if err := proto.Unmarshal(response.Data, &resp); err != nil {
			log.Printf("Failed to unmarshal RNG positive response: %v", err)
		} else {
			log.Printf("RNG Positive Stats: Count=%d", resp.NCount)
		}
	}

	// Query RNG negative node count
	response, err = nc.Request("rpc.rng_negative", reqBytes, 5*time.Second)
	if err != nil {
		log.Printf("Failed to get RNG negative count: %v", err)
	} else {
		var resp rngpb.RngCountResponse
		if err := proto.Unmarshal(response.Data, &resp); err != nil {
			log.Printf("Failed to unmarshal RNG negative response: %v", err)
		} else {
			log.Printf("RNG Negative Stats: Count=%d", resp.NCount)
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
