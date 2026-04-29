package config

import (
	"os"
	"time"

	"github.com/joho/godotenv"
)

type Config struct {
	HTTPPort      string
	DBDSN         string
	JWTSecret     string
	JWTExpiration time.Duration
}

func Load() (*Config, error) {
	if err := godotenv.Load(); err != nil {
		return nil, err
	}

	return &Config{
		HTTPPort:      os.Getenv("HTTP_PORT"),
		DBDSN:         os.Getenv("DB_DSN"),
		JWTSecret:     os.Getenv("JWT_SECRET"),
		JWTExpiration: 24 * time.Hour,
	}, nil
}
