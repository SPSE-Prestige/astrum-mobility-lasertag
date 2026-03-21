package config

import (
	"os"
	"strconv"
	"time"
)

type Config struct {
	Server   ServerConfig
	Postgres PostgresConfig
	Redis    RedisConfig
	JWT      JWTConfig
	Game     GameDefaults
}

type ServerConfig struct {
	Port string
}

type PostgresConfig struct {
	Host     string
	Port     string
	User     string
	Password string
	DBName   string
	SSLMode  string
}

type RedisConfig struct {
	Addr     string
	Password string
	DB       int
}

type JWTConfig struct {
	Secret     string
	Expiration time.Duration
}

type GameDefaults struct {
	DefaultHP         int
	RespawnDelay      time.Duration
	MaxPlayersPerGame int
}

func (p PostgresConfig) DSN() string {
	return "host=" + p.Host +
		" port=" + p.Port +
		" user=" + p.User +
		" password=" + p.Password +
		" dbname=" + p.DBName +
		" sslmode=" + p.SSLMode
}

func Load() *Config {
	return &Config{
		Server: ServerConfig{
			Port: envOrDefault("SERVER_PORT", "8080"),
		},
		Postgres: PostgresConfig{
			Host:     envOrDefault("POSTGRES_HOST", "localhost"),
			Port:     envOrDefault("POSTGRES_PORT", "5432"),
			User:     envOrDefault("POSTGRES_USER", "lasertag"),
			Password: envOrDefault("POSTGRES_PASSWORD", "lasertag"),
			DBName:   envOrDefault("POSTGRES_DB", "lasertag"),
			SSLMode:  envOrDefault("POSTGRES_SSLMODE", "disable"),
		},
		Redis: RedisConfig{
			Addr:     envOrDefault("REDIS_ADDR", "localhost:6379"),
			Password: envOrDefault("REDIS_PASSWORD", ""),
			DB:       envOrDefaultInt("REDIS_DB", 0),
		},
		JWT: JWTConfig{
			Secret:     envOrDefault("JWT_SECRET", "change-me-in-production"),
			Expiration: envOrDefaultDuration("JWT_EXPIRATION", 24*time.Hour),
		},
		Game: GameDefaults{
			DefaultHP:         envOrDefaultInt("GAME_DEFAULT_HP", 100),
			RespawnDelay:      envOrDefaultDuration("GAME_RESPAWN_DELAY", 5*time.Second),
			MaxPlayersPerGame: envOrDefaultInt("GAME_MAX_PLAYERS", 20),
		},
	}
}

func envOrDefault(key, fallback string) string {
	if val := os.Getenv(key); val != "" {
		return val
	}
	return fallback
}

func envOrDefaultInt(key string, fallback int) int {
	if val := os.Getenv(key); val != "" {
		if parsed, err := strconv.Atoi(val); err == nil {
			return parsed
		}
	}
	return fallback
}

func envOrDefaultDuration(key string, fallback time.Duration) time.Duration {
	if val := os.Getenv(key); val != "" {
		if parsed, err := time.ParseDuration(val); err == nil {
			return parsed
		}
	}
	return fallback
}
