package config

import "os"

type Config struct {
	HTTPPort      string
	KafkaBrokers  string
	ClickhouseDSN string
}

func Load() *Config {
	return &Config{
		HTTPPort:      os.Getenv("HTTP_PORT"),
		KafkaBrokers:  os.Getenv("KAFKA_BROKERS"),
		ClickhouseDSN: os.Getenv("CLICKHOUSE_DSN"),
	}
}
