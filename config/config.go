package config

import (
	"fmt"
	"strings"

	"github.com/spf13/viper"
)

// Config holds all configuration for the application
type Config struct {
	Database DatabaseConfig `mapstructure:"database"`
}

// DatabaseConfig holds database configuration
type DatabaseConfig struct {
	Host     string `mapstructure:"host"`
	Port     int    `mapstructure:"port"`
	User     string `mapstructure:"user"`
	Password string `mapstructure:"password"`
	DBName   string `mapstructure:"dbname"`
	SSLMode  string `mapstructure:"sslmode"`
}

// LoadConfig loads configuration from file, environment variables or defaults
func LoadConfig(configPath string) (*Config, error) {
	var config Config

	// Set defaults
	viper.SetDefault("database.host", "localhost")
	viper.SetDefault("database.port", 5434)
	viper.SetDefault("database.user", "postgres")
	viper.SetDefault("database.password", "postgres")
	viper.SetDefault("database.dbname", "go-env-cli")
	viper.SetDefault("database.sslmode", "disable")

	// Set config file path if provided
	if configPath != "" {
		viper.SetConfigFile(configPath)
	} else {
		// Search config in home directory with name ".go-env-cli" (without extension)
		viper.AddConfigPath(".")
		viper.AddConfigPath("$HOME")
		viper.SetConfigName(".go-env-cli")
	}

	// Read the config file
	if err := viper.ReadInConfig(); err != nil {
		// It's ok if config file doesn't exist
		if _, ok := err.(viper.ConfigFileNotFoundError); !ok {
			return nil, fmt.Errorf("error reading config file: %w", err)
		}
	}

	// Set up environment variable prefixes
	viper.SetEnvPrefix("GO_ENV_CLI")
	viper.SetEnvKeyReplacer(strings.NewReplacer(".", "_"))
	viper.AutomaticEnv()

	// Unmarshal config
	if err := viper.Unmarshal(&config); err != nil {
		return nil, fmt.Errorf("error unmarshalling config: %w", err)
	}

	// Override with environment variables from provided .env file, if any
	if configPath != "" && strings.HasSuffix(configPath, ".env") {
		if err := viper.BindEnv("database.host", "DB_HOST"); err != nil {
			return nil, fmt.Errorf("error binding env var DB_HOST: %w", err)
		}
		if err := viper.BindEnv("database.port", "DB_PORT"); err != nil {
			return nil, fmt.Errorf("error binding env var DB_PORT: %w", err)
		}
		if err := viper.BindEnv("database.user", "DB_USER"); err != nil {
			return nil, fmt.Errorf("error binding env var DB_USER: %w", err)
		}
		if err := viper.BindEnv("database.password", "DB_PASSWORD"); err != nil {
			return nil, fmt.Errorf("error binding env var DB_PASSWORD: %w", err)
		}
		if err := viper.BindEnv("database.dbname", "DB_NAME"); err != nil {
			return nil, fmt.Errorf("error binding env var DB_NAME: %w", err)
		}
	}

	return &config, nil
}
