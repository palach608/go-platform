package config

import (
	"os"
)

type Config struct {
	DBDSN     string
	HTTPPort  string
	JWTSecret string
}

func Load() *Config {
	return &Config{
		DBDSN:     getEnv("DB_DSN", "host=localhost user=postgres password=pass dbname=auth_db port=5432 sslmode=disable"),
		HTTPPort:  getEnv("HTTP_PORT", "8081"),
		JWTSecret: getEnv("JWT_SECRET", "very-secret-key"),
	}
}

func getEnv(key, fallback string) string {
	if value, ok := os.LookupEnv(key); ok {
		return value
	}
	return fallback
}
