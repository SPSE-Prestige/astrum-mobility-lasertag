package config

import (
	"fmt"
	"os"
	"strconv"
	"strings"
	"time"
)

type Config struct {
	// Server
	ServerPort  string
	CORSOrigins []string

	// PostgreSQL
	PostgresHost string
	PostgresPort string
	PostgresUser string
	PostgresPass string
	PostgresDB   string

	// Redis (reserved for future cache layer)
	RedisAddr string

	// MQTT
	MQTTBroker string

	// Auth
	SessionTTL time.Duration

	// Background tasks
	HeartbeatCheckInterval time.Duration
	DeviceOfflineTimeout   time.Duration
	AutoEndCheckInterval   time.Duration
	SessionCleanupInterval time.Duration

	// HTTP
	ReadTimeout     time.Duration
	WriteTimeout    time.Duration
	IdleTimeout     time.Duration
	ShutdownTimeout time.Duration

	// Rate limiting
	RateLimitRPS   float64 // requests per second for login
	RateLimitBurst int
}

func Load() (*Config, error) {
	cfg := &Config{
		ServerPort:   getEnv("SERVER_PORT", "8080"),
		CORSOrigins:  strings.Split(getEnv("CORS_ORIGINS", "*"), ","),
		PostgresHost: getEnv("POSTGRES_HOST", "localhost"),
		PostgresPort: getEnv("POSTGRES_PORT", "5432"),
		PostgresUser: getEnv("POSTGRES_USER", "lasertag"),
		PostgresPass: getEnv("POSTGRES_PASSWORD", "lasertag"),
		PostgresDB:   getEnv("POSTGRES_DB", "lasertag"),
		RedisAddr:    getEnv("REDIS_ADDR", "localhost:6379"),
		MQTTBroker:   getEnv("MQTT_BROKER", "tcp://localhost:1883"),

		SessionTTL:             getDurationEnv("SESSION_TTL", 24*time.Hour),
		HeartbeatCheckInterval: getDurationEnv("HEARTBEAT_CHECK_INTERVAL", 10*time.Second),
		DeviceOfflineTimeout:   getDurationEnv("DEVICE_OFFLINE_TIMEOUT", 30*time.Second),
		AutoEndCheckInterval:   getDurationEnv("AUTO_END_CHECK_INTERVAL", 5*time.Second),
		SessionCleanupInterval: getDurationEnv("SESSION_CLEANUP_INTERVAL", 1*time.Hour),

		ReadTimeout:     getDurationEnv("HTTP_READ_TIMEOUT", 15*time.Second),
		WriteTimeout:    getDurationEnv("HTTP_WRITE_TIMEOUT", 15*time.Second),
		IdleTimeout:     getDurationEnv("HTTP_IDLE_TIMEOUT", 60*time.Second),
		ShutdownTimeout: getDurationEnv("SHUTDOWN_TIMEOUT", 10*time.Second),

		RateLimitRPS:   getFloatEnv("RATE_LIMIT_RPS", 5),
		RateLimitBurst: getIntEnv("RATE_LIMIT_BURST", 10),
	}

	if err := cfg.validate(); err != nil {
		return nil, err
	}
	return cfg, nil
}

func (c *Config) validate() error {
	port, err := strconv.Atoi(c.ServerPort)
	if err != nil || port < 1 || port > 65535 {
		return fmt.Errorf("invalid SERVER_PORT: %q", c.ServerPort)
	}
	if c.PostgresHost == "" {
		return fmt.Errorf("POSTGRES_HOST is required")
	}
	if c.PostgresDB == "" {
		return fmt.Errorf("POSTGRES_DB is required")
	}
	if c.MQTTBroker == "" {
		return fmt.Errorf("MQTT_BROKER is required")
	}
	return nil
}

func getEnv(key, fallback string) string {
	if v := os.Getenv(key); v != "" {
		return v
	}
	return fallback
}

func getDurationEnv(key string, fallback time.Duration) time.Duration {
	v := os.Getenv(key)
	if v == "" {
		return fallback
	}
	d, err := time.ParseDuration(v)
	if err != nil {
		return fallback
	}
	return d
}

func getIntEnv(key string, fallback int) int {
	v := os.Getenv(key)
	if v == "" {
		return fallback
	}
	i, err := strconv.Atoi(v)
	if err != nil {
		return fallback
	}
	return i
}

func getFloatEnv(key string, fallback float64) float64 {
	v := os.Getenv(key)
	if v == "" {
		return fallback
	}
	f, err := strconv.ParseFloat(v, 64)
	if err != nil {
		return fallback
	}
	return f
}
