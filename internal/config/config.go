package config

import (
	"encoding/json"
	"os"
	"sync"
)

type Config struct {
	Name   string `json:"name"`
	Sequex struct {
		Host string `json:"host"`
		Port int    `json:"port"`
	} `json:"sequex"`
	Solvexity struct {
		Host string `json:"host"`
		Port int    `json:"port"`
	} `json:"solvexity"`
}

type Domain struct {
	rwMux sync.RWMutex
	c     *Config
}

func NewDomain(path string) *Domain {
	config, err := loadConfig(path)
	if err != nil {
		panic(err) // Handle the error
	}
	return &Domain{
		c: config,
	}
}

func loadConfig(path string) (*Config, error) {
	file, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	var config Config
	decoder := json.NewDecoder(file)
	err = decoder.Decode(&config)
	if err != nil {
		return nil, err
	}

	return &config, nil
}

func (d *Domain) GetConfig() *Config {
	d.rwMux.RLock()
	defer d.rwMux.RUnlock()
	return d.c
}
