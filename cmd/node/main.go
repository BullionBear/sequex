package main

import (
	"flag"
	"log"
	"os"
	"os/signal"
	"syscall"

	"github.com/BullionBear/sequex/internal/config"
	_ "github.com/BullionBear/sequex/internal/nodeimpl/example/rng" // Import to register RNG node
	_ "github.com/BullionBear/sequex/internal/nodeimpl/example/sum" // Import to register Sum node
	"github.com/BullionBear/sequex/pkg/node"
)

func main() {
	// Parse command line arguments
	var configFile string
	flag.StringVar(&configFile, "c", "config/rng.yml", "Configuration file path")
	flag.Parse()

	log.Printf("Starting NATS Microservices with config: %s", configFile)

	// Load configuration
	cfg, err := config.LoadConfig(configFile)
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
	for _, nodeConfig := range cfg.Deployer.Nodes {
		nodeName := nodeConfig["name"].(string)
		log.Printf("Creating node: %s", nodeName)

		node, err := config.CreateNode(nodeConfig, nc)
		if err != nil {
			log.Fatalf("Failed to create node %s: %v", nodeName, err)
		}

		if err := d.RegisterNode(node); err != nil {
			log.Fatalf("Failed to register node %s: %v", nodeName, err)
		}

		log.Printf("Successfully registered node: %s", nodeName)
	}

	// Start all nodes
	for _, nodeConfig := range cfg.Deployer.Nodes {
		nodeName := nodeConfig["name"].(string)
		if err := d.Start(nodeName); err != nil {
			log.Fatalf("Failed to start node %s: %v", nodeName, err)
		}
		log.Printf("Successfully started node: %s", nodeName)
	}

	log.Println("All nodes deployed successfully")

	// Wait for shutdown signal
	waitForShutdown()

	// Stop all nodes
	for _, nodeConfig := range cfg.Deployer.Nodes {
		nodeName := nodeConfig["name"].(string)
		if err := d.Stop(nodeName); err != nil {
			log.Printf("Error stopping node %s: %v", nodeName, err)
		}
	}

	log.Println("All nodes stopped")
}

func waitForShutdown() {
	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)

	log.Println("Waiting for shutdown signal...")
	<-sigChan
	log.Println("Shutdown signal received")
}
