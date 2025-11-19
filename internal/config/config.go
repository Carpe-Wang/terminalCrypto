package config

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/spf13/viper"
)

type Config struct {
	Exchange        string          `mapstructure:"exchange"`
	Exchanges       map[string]bool `mapstructure:"exchanges"`
	RefreshInterval int             `mapstructure:"refresh_interval"`
	Display         DisplayConfig   `mapstructure:"display"`
}

type DisplayConfig struct {
	Currency      string `mapstructure:"currency"`
	DecimalPlaces int    `mapstructure:"decimal_places"`
}

// InitConfig initializes the configuration
func InitConfig() error {
	home, err := os.UserHomeDir()
	if err != nil {
		return fmt.Errorf("failed to get home directory: %w", err)
	}

	configDir := filepath.Join(home, ".terminalcrypto")
	configPath := filepath.Join(configDir, "config.yaml")

	// Create config directory if it doesn't exist
	if err := os.MkdirAll(configDir, 0700); err != nil {
		return fmt.Errorf("failed to create config directory: %w", err)
	}

	// Set config file settings
	viper.SetConfigName("config")
	viper.SetConfigType("yaml")
	viper.AddConfigPath(configDir)
	viper.AddConfigPath(".")

	// Set environment variable prefix
	viper.SetEnvPrefix("CRYPTO")
	viper.AutomaticEnv()

	// Set default values
	viper.SetDefault("exchange", "binance")
	viper.SetDefault("exchanges.binance", false)
	viper.SetDefault("exchanges.coinbase", false)
	viper.SetDefault("exchanges.okx", false)
	viper.SetDefault("refresh_interval", 5)
	viper.SetDefault("display.currency", "USDT")
	viper.SetDefault("display.decimal_places", 2)

	// Try to read existing config
	if err := viper.ReadInConfig(); err != nil {
		// If config file doesn't exist, create it with defaults
		if _, ok := err.(viper.ConfigFileNotFoundError); ok {
			if err := viper.SafeWriteConfigAs(configPath); err != nil {
				return fmt.Errorf("failed to write default config: %w", err)
			}
		} else {
			return fmt.Errorf("failed to read config: %w", err)
		}
	}

	// Set proper permissions on config file
	if err := os.Chmod(configPath, 0600); err != nil {
		return fmt.Errorf("failed to set config file permissions: %w", err)
	}

	return nil
}

// GetConfig returns the current configuration
func GetConfig() (*Config, error) {
	var cfg Config
	if err := viper.Unmarshal(&cfg); err != nil {
		return nil, fmt.Errorf("failed to unmarshal config: %w", err)
	}
	return &cfg, nil
}

// SetExchange sets the default exchange
func SetExchange(exchange string) error {
	viper.Set("exchange", exchange)
	viper.Set("exchanges."+exchange, true)
	return viper.WriteConfig()
}

// GetExchange returns the currently selected exchange
func GetExchange() string {
	return viper.GetString("exchange")
}

// IsExchangeEnabled checks if an exchange is enabled
func IsExchangeEnabled(exchange string) bool {
	return viper.GetBool("exchanges." + exchange)
}
