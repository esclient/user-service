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

	//DatabaseURL string `mapstructure:"DB_URL"`

	InfisicalSecretKey string `mapstructure:"INFISICAL_SECRET_KEY"`
	InfisicalProjectId string `mapstructure:"INFISICAL_PROJECT_ID"`
	InfisicalEnv string `mapstructure:"INFISICAL_ENV"`
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

	if err := viper.BindEnv("INFISICAL_SECRET_KEY"); err != nil {
		log.Fatalf("viper.BindEnv INFISICAL_SECRET_KEY error: %v", err)
	}

	if err := viper.BindEnv("INFISICAL_PROJECT_ID"); err != nil {
		log.Fatalf("viper.BindEnv INFISICAL_PROJECT_ID error: %v", err)
	}

	if err := viper.BindEnv("INFISICAL_ENV"); err != nil {
		log.Fatalf("viper.BindEnv INFISICAL_ENV error: %v", err)
	}

	host := viper.GetString("HOST")
	port := viper.GetString("PORT")
	infisicalSecretKey := viper.GetString("INFISICAL_SECRET_KEY")
	infisicalProjectId := viper.GetString("INFISICAL_PROJECT_ID")
	infisicalEnv := viper.GetString("INFISICAL_ENV")

	return &Config{
		Host:                   host,
		Port:                   port,
		InfisicalSecretKey:     infisicalSecretKey,
		InfisicalProjectId: 	infisicalProjectId,
		InfisicalEnv: 	        infisicalEnv,
	}
}