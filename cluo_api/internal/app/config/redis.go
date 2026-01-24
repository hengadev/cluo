package config

import (
	"fmt"
)

// RedisConfig holds Redis configuration.
type RedisConfig struct {
	Host     string
	Port     int
	Password string
	DB       int
}

// Addr returns the Redis address in host:port format.
func (c RedisConfig) Addr() string {
	return fmt.Sprintf("%s:%d", c.Host, c.Port)
}

func loadRedisConfig() (RedisConfig, error) {
	port, err := parseIntEnv("CLUO_REDIS_PORT", 6379)
	if err != nil {
		return RedisConfig{}, fmt.Errorf("invalid port: %w", err)
	}

	db, err := parseIntEnv("CLUO_REDIS_DB", 0)
	if err != nil {
		return RedisConfig{}, fmt.Errorf("invalid db: %w", err)
	}

	return RedisConfig{
		Host:     getEnv("CLUO_REDIS_HOST", ""),
		Port:     port,
		Password: getEnv("CLUO_REDIS_PASSWORD", ""),
		DB:       db,
	}, nil
}
