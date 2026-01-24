package config

import (
	"fmt"
	"strconv"
	"time"
)

// ServerConfig holds HTTP server configuration.
type ServerConfig struct {
	Host            string
	Port            int
	ReadTimeout     time.Duration
	WriteTimeout    time.Duration
	IdleTimeout     time.Duration
	ShutdownTimeout time.Duration
}

// Addr returns the server address in host:port format.
func (c ServerConfig) Addr() string {
	return fmt.Sprintf("%s:%d", c.Host, c.Port)
}

func loadServerConfig() (ServerConfig, error) {
	port, err := parseIntEnv("CLUO_SERVER_PORT", 8080)
	if err != nil {
		return ServerConfig{}, fmt.Errorf("invalid port: %w", err)
	}

	readTimeout, err := parseDurationEnv("CLUO_SERVER_READ_TIMEOUT", 15*time.Second)
	if err != nil {
		return ServerConfig{}, fmt.Errorf("invalid read timeout: %w", err)
	}

	writeTimeout, err := parseDurationEnv("CLUO_SERVER_WRITE_TIMEOUT", 15*time.Second)
	if err != nil {
		return ServerConfig{}, fmt.Errorf("invalid write timeout: %w", err)
	}

	idleTimeout, err := parseDurationEnv("CLUO_SERVER_IDLE_TIMEOUT", 60*time.Second)
	if err != nil {
		return ServerConfig{}, fmt.Errorf("invalid idle timeout: %w", err)
	}

	shutdownTimeout, err := parseDurationEnv("CLUO_SERVER_SHUTDOWN_TIMEOUT", 30*time.Second)
	if err != nil {
		return ServerConfig{}, fmt.Errorf("invalid shutdown timeout: %w", err)
	}

	return ServerConfig{
		Host:            getEnv("CLUO_SERVER_HOST", "0.0.0.0"),
		Port:            port,
		ReadTimeout:     readTimeout,
		WriteTimeout:    writeTimeout,
		IdleTimeout:     idleTimeout,
		ShutdownTimeout: shutdownTimeout,
	}, nil
}

// parseIntEnv parses an integer from an environment variable.
func parseIntEnv(key string, defaultValue int) (int, error) {
	strVal := getEnv(key, "")
	if strVal == "" {
		return defaultValue, nil
	}

	val, err := strconv.Atoi(strVal)
	if err != nil {
		return 0, fmt.Errorf("parse %s: %w", key, err)
	}

	return val, nil
}

// parseDurationEnv parses a duration from an environment variable.
func parseDurationEnv(key string, defaultValue time.Duration) (time.Duration, error) {
	strVal := getEnv(key, "")
	if strVal == "" {
		return defaultValue, nil
	}

	val, err := time.ParseDuration(strVal)
	if err != nil {
		return 0, fmt.Errorf("parse %s: %w", key, err)
	}

	return val, nil
}
