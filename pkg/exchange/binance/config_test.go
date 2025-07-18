package binance

import (
	"io/ioutil"
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestLoadConfig(t *testing.T) {
	// Create a temporary config file
	configContent := `
accounts:
  binance:
    - name: "test_account"
      api_key: "test_api_key"
      api_secret: "test_api_secret"
      sandbox: true
      timeout: 15
    - name: "prod_account"
      api_key: "prod_api_key"
      api_secret: "prod_api_secret"
      sandbox: false
      timeout: 30

market:
  binance:
    - "BTCUSDT"
    - "ETHUSDT"
    - "ADAUSDT"
`

	tmpDir, err := ioutil.TempDir("", "binance_test")
	require.NoError(t, err)
	defer os.RemoveAll(tmpDir)

	configFile := filepath.Join(tmpDir, "test_config.yml")
	err = ioutil.WriteFile(configFile, []byte(configContent), 0644)
	require.NoError(t, err)

	t.Run("ValidConfigFile", func(t *testing.T) {
		config, err := LoadConfig(configFile)
		assert.NoError(t, err)
		assert.NotNil(t, config)

		// Check binance accounts
		binanceAccounts, exists := config.Accounts["binance"]
		assert.True(t, exists)
		assert.Len(t, binanceAccounts, 2)

		// Check first account
		assert.Equal(t, "test_account", binanceAccounts[0].Name)
		assert.Equal(t, "test_api_key", binanceAccounts[0].APIKey)
		assert.Equal(t, "test_api_secret", binanceAccounts[0].APISecret)
		assert.True(t, binanceAccounts[0].Sandbox)
		assert.Equal(t, 15, binanceAccounts[0].Timeout)

		// Check second account
		assert.Equal(t, "prod_account", binanceAccounts[1].Name)
		assert.False(t, binanceAccounts[1].Sandbox)
		assert.Equal(t, 30, binanceAccounts[1].Timeout)

		// Check market symbols
		symbols := config.GetBinanceSymbols()
		assert.Len(t, symbols, 3)
		assert.Contains(t, symbols, "BTCUSDT")
		assert.Contains(t, symbols, "ETHUSDT")
		assert.Contains(t, symbols, "ADAUSDT")
	})

	t.Run("NonExistentFile", func(t *testing.T) {
		config, err := LoadConfig("non_existent_file.yml")
		assert.Error(t, err)
		assert.Nil(t, config)
	})

	t.Run("InvalidYAML", func(t *testing.T) {
		invalidFile := filepath.Join(tmpDir, "invalid.yml")
		err := ioutil.WriteFile(invalidFile, []byte("invalid: yaml: content: ["), 0644)
		require.NoError(t, err)

		config, err := LoadConfig(invalidFile)
		assert.Error(t, err)
		assert.Nil(t, config)
	})
}

func TestAppConfig_GetBinanceConfig(t *testing.T) {
	config := &AppConfig{
		Accounts: map[string][]Config{
			"binance": {
				{
					Name:      "account1",
					APIKey:    "key1",
					APISecret: "secret1",
				},
				{
					Name:      "account2",
					APIKey:    "key2",
					APISecret: "secret2",
				},
			},
		},
	}

	t.Run("ValidConfig", func(t *testing.T) {
		binanceConfig, err := config.GetBinanceConfig()
		assert.NoError(t, err)
		assert.NotNil(t, binanceConfig)
		assert.Equal(t, "account1", binanceConfig.Name)
		assert.Equal(t, "key1", binanceConfig.APIKey)
	})

	t.Run("NoAccounts", func(t *testing.T) {
		emptyConfig := &AppConfig{
			Accounts: map[string][]Config{},
		}

		binanceConfig, err := emptyConfig.GetBinanceConfig()
		assert.Error(t, err)
		assert.Nil(t, binanceConfig)
		assert.Contains(t, err.Error(), "no Binance configuration found")
	})

	t.Run("EmptyBinanceAccounts", func(t *testing.T) {
		emptyBinanceConfig := &AppConfig{
			Accounts: map[string][]Config{
				"binance": {},
			},
		}

		binanceConfig, err := emptyBinanceConfig.GetBinanceConfig()
		assert.Error(t, err)
		assert.Nil(t, binanceConfig)
	})
}

func TestAppConfig_GetBinanceConfigByName(t *testing.T) {
	config := &AppConfig{
		Accounts: map[string][]Config{
			"binance": {
				{
					Name:      "test_account",
					APIKey:    "test_key",
					APISecret: "test_secret",
				},
				{
					Name:      "prod_account",
					APIKey:    "prod_key",
					APISecret: "prod_secret",
				},
			},
		},
	}

	t.Run("ExistingAccount", func(t *testing.T) {
		binanceConfig, err := config.GetBinanceConfigByName("test_account")
		assert.NoError(t, err)
		assert.NotNil(t, binanceConfig)
		assert.Equal(t, "test_account", binanceConfig.Name)
		assert.Equal(t, "test_key", binanceConfig.APIKey)
	})

	t.Run("NonExistentAccount", func(t *testing.T) {
		binanceConfig, err := config.GetBinanceConfigByName("non_existent")
		assert.Error(t, err)
		assert.Nil(t, binanceConfig)
		assert.Contains(t, err.Error(), "not found")
	})

	t.Run("NoBinanceAccounts", func(t *testing.T) {
		emptyConfig := &AppConfig{
			Accounts: map[string][]Config{},
		}

		binanceConfig, err := emptyConfig.GetBinanceConfigByName("any_name")
		assert.Error(t, err)
		assert.Nil(t, binanceConfig)
	})
}

func TestAppConfig_GetBinanceSymbols(t *testing.T) {
	t.Run("WithSymbols", func(t *testing.T) {
		config := &AppConfig{
			Market: map[string][]string{
				"binance": {"BTCUSDT", "ETHUSDT", "ADAUSDT"},
			},
		}

		symbols := config.GetBinanceSymbols()
		assert.Len(t, symbols, 3)
		assert.Contains(t, symbols, "BTCUSDT")
		assert.Contains(t, symbols, "ETHUSDT")
		assert.Contains(t, symbols, "ADAUSDT")
	})

	t.Run("NoSymbols", func(t *testing.T) {
		config := &AppConfig{
			Market: map[string][]string{},
		}

		symbols := config.GetBinanceSymbols()
		assert.Empty(t, symbols)
	})

	t.Run("NilMarket", func(t *testing.T) {
		config := &AppConfig{}

		symbols := config.GetBinanceSymbols()
		assert.Empty(t, symbols)
	})
}
