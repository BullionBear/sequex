package config

import (
	"fmt"
	"log"
	"os"

	"github.com/BullionBear/sequex/internal/node"

	"github.com/nats-io/nats.go"
	"gopkg.in/yaml.v3"
)

// Config represents the merged configuration structure
type Config struct {
	NATS     NATSConfig     `yaml:"nats"`
	Deployer DeployerConfig `yaml:"deployer"`
}

// NATSConfig represents NATS connection configuration
type NATSConfig struct {
	URL string `yaml:"url"`
}

// DeployerConfig represents deployer configuration
type DeployerConfig struct {
	Nodes map[string]NodeConfig `yaml:"nodes"`
}

// NodeConfig represents individual node configuration
type NodeConfig map[string]interface{}

// LoadConfig loads the merged configuration from file
func LoadConfig(configPath string) (*Config, error) {
	data, err := os.ReadFile(configPath)
	if err != nil {
		return nil, fmt.Errorf("failed to read config file: %w", err)
	}

	var config Config
	if err := yaml.Unmarshal(data, &config); err != nil {
		return nil, fmt.Errorf("failed to parse config file: %w", err)
	}

	return &config, nil
}

// CreateNATSConnection creates a single NATS connection for the entire process
func CreateNATSConnection(natsURL string) (*nats.Conn, error) {
	nc, err := nats.Connect(natsURL)
	if err != nil {
		return nil, fmt.Errorf("failed to connect to NATS: %w", err)
	}
	log.Printf("Connected to NATS at: %s", natsURL)
	return nc, nil
}

// CreateNode creates a node based on its type name
func CreateNode(nodeName string, nodeConfig NodeConfig, nc *nats.Conn) (node.Node, error) {
	// Extract node type from configuration
	nodeType, ok := nodeConfig["type"].(string)
	if !ok {
		return nil, fmt.Errorf("missing 'type' field in node configuration for node '%s'", nodeName)
	}

	// Remove type from config to avoid passing it to the node
	delete(nodeConfig, "type")

	switch nodeType {
	case "random-generator":
		return createRandomGeneratorNode(nodeName, nodeConfig, nc)
	case "sum-node":
		return createSumNode(nodeName, nodeConfig, nc)
	default:
		return nil, fmt.Errorf("unknown node type '%s' for node '%s'", nodeType, nodeName)
	}
}

// createRandomGeneratorNode creates a random generator node
func createRandomGeneratorNode(nodeName string, config NodeConfig, nc *nats.Conn) (node.Node, error) {
	randomNode := randomgenerator.NewRandomGeneratorNode()

	// Set the node name to the unique identifier
	randomNode.SetName(nodeName)

	// Inject NATS connection
	randomNode.SetNATSConnection(nc)

	// Set up the node directly with the configuration
	if err := setupRandomGeneratorNode(randomNode, config); err != nil {
		return nil, fmt.Errorf("failed to setup random generator: %w", err)
	}

	return randomNode, nil
}

// setupRandomGeneratorNode configures the random generator node
func setupRandomGeneratorNode(randomNode *randomgenerator.RandomGeneratorNode, config NodeConfig) error {
	// Extract configuration values
	low, ok := config["low"].(int)
	if !ok {
		return fmt.Errorf("invalid or missing 'low' configuration")
	}

	high, ok := config["high"].(int)
	if !ok {
		return fmt.Errorf("invalid or missing 'high' configuration")
	}

	interval, ok := config["interval"].(int)
	if !ok {
		return fmt.Errorf("invalid or missing 'interval' configuration")
	}

	publishSubject, ok := config["publish_subject"].(string)
	if !ok {
		return fmt.Errorf("invalid or missing 'publish_subject' configuration")
	}

	rpcSubject, ok := config["rpc_subject"].(string)
	if !ok {
		return fmt.Errorf("invalid or missing 'rpc_subject' configuration")
	}

	// Create type-safe configuration
	nodeConfig := randomgenerator.RandomGeneratorConfig{
		Low:            low,
		High:           high,
		Interval:       interval,
		PublishSubject: publishSubject,
		RPCSubject:     rpcSubject,
	}

	// Set the configuration
	if err := randomNode.SetConfiguration(nodeConfig); err != nil {
		return fmt.Errorf("failed to set configuration: %w", err)
	}

	return nil
}

// createSumNode creates a sum node
func createSumNode(nodeName string, config NodeConfig, nc *nats.Conn) (node.Node, error) {
	sumNode := sumnode.NewSumNode()

	// Set the node name to the unique identifier
	sumNode.SetName(nodeName)

	// Inject NATS connection
	sumNode.SetNATSConnection(nc)

	// Set up the node directly with the configuration
	if err := setupSumNode(sumNode, config); err != nil {
		return nil, fmt.Errorf("failed to setup sum node: %w", err)
	}

	return sumNode, nil
}

// setupSumNode configures the sum node
func setupSumNode(sumNode *sumnode.SumNode, config NodeConfig) error {
	// Extract configuration values
	subscribeTopics, ok := config["subscribe_topics"].(string)
	if !ok {
		return fmt.Errorf("invalid or missing 'subscribe_topics' configuration")
	}

	resultSubject, ok := config["result_subject"].(string)
	if !ok {
		return fmt.Errorf("invalid or missing 'result_subject' configuration")
	}

	rpcSubject, ok := config["rpc_subject"].(string)
	if !ok {
		return fmt.Errorf("invalid or missing 'rpc_subject' configuration")
	}

	// Create type-safe configuration
	nodeConfig := sumnode.SumNodeConfig{
		SubscribeTopics: subscribeTopics,
		ResultSubject:   resultSubject,
		RPCSubject:      rpcSubject,
	}

	// Set the configuration
	if err := sumNode.SetConfiguration(nodeConfig); err != nil {
		return fmt.Errorf("failed to set configuration: %w", err)
	}

	return nil
}
