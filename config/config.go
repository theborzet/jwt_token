package config

import (
	"os"

	"github.com/joho/godotenv"
)

type Config struct {
	Port           string
	DBHost         string
	DBPort         string
	DBUser         string
	DBPass         string
	DBName         string
	ExternalApiURL string
}

func LoadConfig() (*Config, error) {
	err := godotenv.Load()
	if err != nil {
		return nil, err
	}

	return &Config{
		Port:           ":" + os.Getenv("PORT"),
		DBHost:         os.Getenv("DB_HOST"),
		DBPort:         os.Getenv("DB_PORT"),
		DBUser:         os.Getenv("DB_USER"),
		DBPass:         os.Getenv("DB_PASS"),
		DBName:         os.Getenv("DB_NAME"),
		ExternalApiURL: os.Getenv("EXTERNAL_API_URL"),
	}, nil
}
