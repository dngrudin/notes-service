package config

import (
	"log"
	"os"
	"time"

	"github.com/ilyakaznacheev/cleanenv"
)

type Config struct {
	HTTPServer HTTPServer `yaml:"http_server"`
}

type HTTPServer struct {
	Address         string        `yaml:"address" env-default:"0.0.0.0:8080"`
	Timeout         time.Duration `yaml:"timeout" env-default:"5s"`
	IdleTimeout     time.Duration `yaml:"idle_timeout" env-default:"60s"`
	ShutdownTimeout time.Duration `yaml:"shutdown_timeout" env-default:"10s"`
}

func MustLoad(configPath string) *Config {
	if _, err := os.Stat(configPath); err != nil {
		log.Fatalf("Error opening config file: %s", err)
	}

	var cfg Config
	err := cleanenv.ReadConfig(configPath, &cfg)
	if err != nil {
		log.Fatalf("Error reading config file: %s", err)
	}

	return &cfg
}
