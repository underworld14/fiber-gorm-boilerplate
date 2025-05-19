package config

import (
	"fmt"
	"strings"

	"github.com/spf13/viper"
)

// Config stores all configuration of the application
type Config struct {
	Environment string `mapstructure:"ENVIRONMENT"`
	DBDriver    string `mapstructure:"DB_DRIVER"`
	DBSource    string `mapstructure:"DB_SOURCE"`
	ServerPort  string `mapstructure:"SERVER_PORT"`
	LogLevel    string `mapstructure:"LOG_LEVEL"`
	JWTSecret   string `mapstructure:"JWT_SECRET"`
}

// LoadConfig reads configuration from file or environment variables
func LoadConfig() (config Config, err error) {
	// Set defaults
	viper.SetDefault("ENVIRONMENT", "development")
	viper.SetDefault("DB_DRIVER", "sqlite")
	viper.SetDefault("DB_SOURCE", "app.db")
	viper.SetDefault("SERVER_PORT", "3000")
	viper.SetDefault("LOG_LEVEL", "info")
	viper.SetDefault("JWT_SECRET", "very-secret")

	// Look for .env file
	viper.SetConfigName(".env")
	viper.SetConfigType("env")
	viper.AddConfigPath(".")

	// Read from .env file if it exists
	if err := viper.ReadInConfig(); err != nil {
		if _, ok := err.(viper.ConfigFileNotFoundError); !ok {
			return config, fmt.Errorf("error reading config file: %w", err)
		}
		// Config file not found, will use defaults and env vars
	}

	// Read from environment variables
	viper.AutomaticEnv()
	viper.SetEnvKeyReplacer(strings.NewReplacer(".", "_"))

	// Unmarshal config
	err = viper.Unmarshal(&config)
	return
}
