package config

import (
	"fmt"

	"github.com/spf13/viper"
)

// Config holds all configuration for the application
type Config struct {
	GO_CLI_DB string `mapstructure:"go_cli_db"`
}

// LoadConfig loads configuration from file, environment variables or defaults
func LoadConfig(_ string) (*Config, error) {
	var config Config

	viper.AutomaticEnv()
	_ = viper.BindEnv("go_cli_db", "GO_CLI_DB")

	if err := viper.Unmarshal(&config); err != nil {
		return nil, fmt.Errorf("error unmarshalling config: %w", err)
	}

	return &config, nil
}
