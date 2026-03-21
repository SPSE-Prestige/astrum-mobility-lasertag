package config

import (
	"os"
)

type Config struct {
	ServerPort   string
	PostgresHost string
	PostgresPort string
	PostgresUser string
	PostgresPass string
	PostgresDB   string
	RedisAddr    string
	MQTTBroker   string
}

func Load() *Config {
	return &Config{
		ServerPort:   getEnv("SERVER_PORT", "8080"),
		PostgresHost: getEnv("POSTGRES_HOST", "localhost"),
		PostgresPort: getEnv("POSTGRES_PORT", "5432"),
		PostgresUser: getEnv("POSTGRES_USER", "lasertag"),
		PostgresPass: getEnv("POSTGRES_PASSWORD", "lasertag"),
		PostgresDB:   getEnv("POSTGRES_DB", "lasertag"),
		RedisAddr:    getEnv("REDIS_ADDR", "localhost:6379"),
		MQTTBroker:   getEnv("MQTT_BROKER", "tcp://localhost:1883"),
	}
}

func getEnv(key, fallback string) string {
	if v := os.Getenv(key); v != "" {
		return v
	}
	return fallback
}
