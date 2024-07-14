package config

import (
	"os"

	"gopkg.in/yaml.v2"
)

type Config struct {
	Port    string        `yaml:"port"`
	DB      DBConfig      `yaml:"db"`
	Binance BinanceConfig `yaml:"binance"`
}

type DBConfig struct {
	Host string `yaml:"host"`
	Port string `yaml:"port"`
	User string `yaml:"user"`
	Pass string `yaml:"pass"`
	Name string `yaml:"name"`
}

type BinanceConfig struct {
	APIURL string `yaml:"api_url"`
}

func LoadConfig() (*Config, error) {

	config, err := os.ReadFile("./configs/config.yaml")
	if err != nil {
		return nil, err
	}

	var cfg Config

	if err := yaml.Unmarshal(config, &cfg); err != nil {
		return nil, err
	}

	return &cfg, nil

}
