package config

import (
	"encoding/json"
	"os"
)

// Config represents the application configuration.
// Config struct matches the structure of the JSON configuration.
type Config struct {
	Strategy struct {
		Symbol   string `json:"symbol"`
		Alpha    string `json:"alpha"`
		BuyPlane struct {
			Normal []int `json:"normal"`
			Shift  []int `json:"shift"`
		} `json:"buy_plane"`
		SellPlane struct {
			Normal []int `json:"normal"`
			Shift  []int `json:"shift"`
		} `json:"sell_plane"`
	} `json:"strategy"`
	Account struct {
		ApiKey    string `json:"api_key"`
		ApiSecret string `json:"api_secret"`
	} `json:"account"`
}

// LoadConfig loads configuration from the specified file path.
func LoadConfig(path string) (*Config, error) {
	var config Config
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}
	if err := json.Unmarshal(data, &config); err != nil {
		return nil, err
	}
	return &config, nil
}
