package config

import (
	"log"
	"os"

	"github.com/joho/godotenv"
	"github.com/spf13/viper"
)

type Config struct {
	Host string `mapstructure:"HOST"`
	Port string `mapstructure:"PORT"`

	DataBaseURL string `mapstructure:"DB_URL"`
}

func LoadConfig() *Config {
	if _, err := os.Stat(".env"); err == nil {
		godotenv.Load()
	}

	viper.AutomaticEnv()
	if err := viper.BindEnv("HOST"); err != nil {
		log.Fatalf("viper.BindEnv HOST error: %v", err)
	}

	if err := viper.BindEnv("PORT"); err != nil {
		log.Fatalf("viper.BindEnv PORT error: %v", err)
	}

	if err := viper.BindEnv("DB_URL"); err != nil {
		log.Fatalf("viper.BindEnv DB_URL error: %v", err)
	}

	host := viper.GetString("HOST")
	port := viper.GetString("PORT")
	dataBaseUrl := viper.GetString("DB_URL")

	return &Config{
		Host:        host,
		Port:        port,
		DataBaseURL: dataBaseUrl,
	}
}