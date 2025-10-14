package config

import (
	"errors"
	"os"
)

type Config struct {
	Host        string `mapstructure:"HOST"`
	Port        string `mapstructure:"PORT"`
	DatabaseURL string `mapstructure:"DATABASE_URL"`
}

func LoadConfig() (*Config, error) {
	host, exists := os.LookupEnv("HOST")
	if !exists {
		return nil, errors.New("HOST environment variable is required")
	}

	port, exists := os.LookupEnv("PORT")
	if !exists {
		return nil, errors.New("PORT environment variable is required")
	}

	databaseURL, exists := os.LookupEnv("DATABASE_URL")
	if !exists {
		return nil, errors.New("DATABASE_URL environment variable is required")
	}

	return &Config{
		Host:        host,
		Port:        port,
		DatabaseURL: databaseURL,
	}, nil
}
