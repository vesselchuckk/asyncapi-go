package config

import (
	"fmt"
	"github.com/caarlos0/env/v11"
	"github.com/joho/godotenv"
	log "github.com/sirupsen/logrus"
)

type Env string

const (
	Env_Test Env = "test"
	Env_Dev  Env = "dev"
)

type Config struct {
	ServerPort       string `env:"SERVER_PORT"`
	ServerHost       string `env:"SERVER_HOST"`
	DatabaseName     string `env:"DB_NAME"`
	DatabaseHost     string `env:"DB_HOST"`
	DatabasePort     string `env:"DB_PORT"`
	DatabaseUser     string `env:"DB_USER"`
	DatabasePassword string `env:"DB_PASSWORD"`
	JwtSecret        string `env:"JWT_SECRET"`
	Env              Env    `env:"ENV" envDefault:"dev"`
}

func (c *Config) DatabaseUrl() string {
	err := godotenv.Load()
	if err != nil {
		log.Fatalf("Error loading .env: %w", err)
	}

	return fmt.Sprintf("postgresql://%s:%s@%s:%s/%s?sslmode=disable",
		c.DatabaseUser,
		c.DatabasePassword,
		c.DatabaseHost,
		c.DatabasePort,
		c.DatabaseName,
	)
}

func New() (*Config, error) {
	err := godotenv.Load()
	if err != nil {
		log.Fatalf("Error loading .env: %w", err)
	}

	var cfg Config
	cfg, err = env.ParseAs[Config]()
	if err != nil {
		return nil, fmt.Errorf("Error getting config: %w", err)
	}

	return &cfg, nil
}
