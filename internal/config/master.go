package config

import (
	"os"

	"gopkg.in/yaml.v3"
)

type MasterConfig struct {
	App struct {
		Port int    `yaml:"port"`
		Host string `yaml:"host"`
	} `yaml:"app"`
	Logger   LoggerConfig   `yaml:"logger"`
	Nats     NATSConfig     `yaml:"nats"`
	Deployer DeployerConfig `yaml:"deployer"`
}

func LoadMasterConfig(configPath string) (*MasterConfig, error) {
	yamlFile, err := os.ReadFile(configPath)
	if err != nil {
		return nil, err
	}

	var cfg MasterConfig
	if err := yaml.Unmarshal(yamlFile, &cfg); err != nil {
		return nil, err
	}
	return &cfg, nil
}
