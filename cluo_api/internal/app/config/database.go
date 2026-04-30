package config

import (
	"fmt"
	"net/url"
	"time"
)

// DatabaseConfig holds PostgreSQL database configuration.
type DatabaseConfig struct {
	Host            string
	Port            int
	Name            string
	User            string
	Password        string
	SSLMode         string
	MaxOpenConns    int
	MaxIdleConns    int
	ConnMaxLifetime time.Duration
	ConnMaxIdleTime time.Duration
}

// DSN returns the PostgreSQL connection string.
func (c DatabaseConfig) DSN() string {
	return fmt.Sprintf(
		"host=%s port=%d dbname=%s user=%s password=%s sslmode=%s",
		c.Host, c.Port, c.Name, c.User, c.Password, c.SSLMode,
	)
}

// ConnectionString returns the PostgreSQL connection URL.
func (c DatabaseConfig) ConnectionString() string {
	return fmt.Sprintf(
		"postgres://%s:%s@%s:%d/%s?sslmode=%s",
		c.User, url.QueryEscape(c.Password), c.Host, c.Port, c.Name, c.SSLMode,
	)
}

func loadDatabaseConfig() (DatabaseConfig, error) {
	port, err := parseIntEnv("CLUO_DB_PORT", 5432)
	if err != nil {
		return DatabaseConfig{}, fmt.Errorf("invalid port: %w", err)
	}

	maxOpenConns, err := parseIntEnv("CLUO_DB_MAX_OPEN_CONNS", 25)
	if err != nil {
		return DatabaseConfig{}, fmt.Errorf("invalid max open conns: %w", err)
	}

	maxIdleConns, err := parseIntEnv("CLUO_DB_MAX_IDLE_CONNS", 5)
	if err != nil {
		return DatabaseConfig{}, fmt.Errorf("invalid max idle conns: %w", err)
	}

	connMaxLifetime, err := parseDurationEnv("CLUO_DB_CONN_MAX_LIFETIME", 5*time.Minute)
	if err != nil {
		return DatabaseConfig{}, fmt.Errorf("invalid conn max lifetime: %w", err)
	}

	connMaxIdleTime, err := parseDurationEnv("CLUO_DB_CONN_MAX_IDLE_TIME", 1*time.Minute)
	if err != nil {
		return DatabaseConfig{}, fmt.Errorf("invalid conn max idle time: %w", err)
	}

	return DatabaseConfig{
		Host:            getEnv("CLUO_DB_HOST", ""),
		Port:            port,
		Name:            getEnv("CLUO_DB_NAME", ""),
		User:            getEnv("CLUO_DB_USER", ""),
		Password:        getEnv("CLUO_DB_PASSWORD", ""),
		SSLMode:         getEnv("CLUO_DB_SSL_MODE", "require"),
		MaxOpenConns:    maxOpenConns,
		MaxIdleConns:    maxIdleConns,
		ConnMaxLifetime: connMaxLifetime,
		ConnMaxIdleTime: connMaxIdleTime,
	}, nil
}
