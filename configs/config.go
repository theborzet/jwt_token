package config

import (
	"os"
	"time"

	"github.com/joho/godotenv"
	"gopkg.in/yaml.v3"
)

type Config struct {
	Port   string
	DBHost string
	DBPort string
	DBUser string
	DBPass string
	DBName string
	Token  struct {
		AccessTokenLifetime  time.Duration `yaml:"access_token_lifetime"`
		RefreshTokenLifetime time.Duration `yaml:"refresh_token_lifetime"`
		SigningKey           string
	}
}

func LoadConfig() (*Config, error) {
	err := godotenv.Load()
	if err != nil {
		return nil, err
	}

	var cfg Config

	yamlFile, err := os.ReadFile("configs/config.yaml")
	if err != nil {
		return nil, err
	}

	err = yaml.Unmarshal(yamlFile, &cfg)
	if err != nil {
		return nil, err
	}

	cfg.Port = os.Getenv("PORT")
	cfg.DBHost = os.Getenv("DB_HOST")
	cfg.DBPort = os.Getenv("DB_PORT")
	cfg.DBUser = os.Getenv("DB_USER")
	cfg.DBPass = os.Getenv("DB_PASS")
	cfg.DBName = os.Getenv("DB_NAME")
	cfg.Token.SigningKey = os.Getenv("JWT_SIGNING_KEY")

	return &cfg, nil
}
