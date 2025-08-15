package config

import (
	"fmt"
	"os"

	"github.com/BullionBear/sequex/pkg/eventbus"
	"github.com/BullionBear/sequex/pkg/log"
	"github.com/BullionBear/sequex/pkg/node"

	"gopkg.in/yaml.v3"
)

// Config represents the merged configuration structure
type SrvConfig struct {
	App struct {
		Host string `yaml:"host"`
		Port int    `yaml:"port"`
	} `yaml:"app"`
	EventBus struct {
		Url string `yaml:"url"`
	} `yaml:"eventbus"`
	Node NodeConfig `yaml:"node"`
}

type NodeConfig struct {
	Name   string            `yaml:"name"`
	Type   string            `yaml:"type"`
	Params map[string]any    `yaml:"params,omitempty"`
	On     map[string]string `yaml:"on,omitempty"`
	Emit   map[string]string `yaml:"emit,omitempty"`
	Rpc    map[string]string `yaml:"rpc,omitempty"`
}

// LoadConfig loads the merged configuration from file
func LoadConfig[T any](configPath string) (*T, error) {
	data, err := os.ReadFile(configPath)
	if err != nil {
		return nil, fmt.Errorf("failed to read config file: %w", err)
	}

	var config T
	if err := yaml.Unmarshal(data, &config); err != nil {
		return nil, fmt.Errorf("failed to parse config file: %w", err)
	}

	return &config, nil
}

// CreateNode creates a node based on its type name
func CreateNode(nodeConfig NodeConfig, eventbus *eventbus.EventBus, logger log.Logger) (node.Node, error) {
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

	// Create the node
	node, err := node.CreateNode(nodeType, eventbus, configMap, logger)
	if err != nil {
		return nil, err
	}

	return node, nil
}
