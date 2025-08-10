package config

import (
	"fmt"
	"log"
	"os"

	"github.com/BullionBear/sequex/pkg/node"

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
	Nodes []NodeConfig `yaml:"nodes"`
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
func CreateNode(nodeConfig NodeConfig, nc *nats.Conn) (node.Node, error) {
	// Extract the node name
	nodeName, ok := nodeConfig["name"].(string)
	if !ok {
		return nil, fmt.Errorf("node name not found in config")
	}

	// Extract the node type
	nodeType, ok := nodeConfig["type"].(string)
	if !ok {
		return nil, fmt.Errorf("node type not found for node %s", nodeName)
	}

	// Extract the actual configuration
	config, ok := nodeConfig["config"]
	if !ok {
		return nil, fmt.Errorf("config not found for node %s", nodeName)
	}

	// Convert config to NodeConfig
	configMap, ok := config.(NodeConfig)
	if !ok {
		return nil, fmt.Errorf("invalid config format for node %s, got type %T", nodeName, config)
	}

	// Add the node name to the config
	configMap["name"] = nodeName

	return node.CreateNode(nodeType, nc, configMap)
}
