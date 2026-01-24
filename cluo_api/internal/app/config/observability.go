package config

import (
	"fmt"

	"github.com/hengadev/cluo_api/internal/common/envmode"
)

// ObservabilityConfig holds logging, metrics, and tracing configuration.
type ObservabilityConfig struct {
	// Logging
	LogLevel string
	LogStyle string

	// Metrics
	MetricsEnabled bool
	MetricsPath    string

	// Tracing
	TracingEnabled  bool
	TracingEndpoint string
	TracingSampler  float64

	// Service identification
	ServiceName    string
	ServiceVersion string
}

func loadObservabilityConfig(env envmode.Mode) (ObservabilityConfig, error) {
	// Set defaults based on environment
	defaultLogLevel := "info"
	defaultLogStyle := "json"

	switch env {
	case envmode.Dev:
		defaultLogLevel = "debug"
		defaultLogStyle = "dev"
	case envmode.Staging:
		defaultLogLevel = "debug"
		defaultLogStyle = "json"
	case envmode.Prod:
		defaultLogLevel = "info"
		defaultLogStyle = "json"
	}

	metricsEnabled, err := parseBoolEnv("CLUO_METRICS_ENABLED", true)
	if err != nil {
		return ObservabilityConfig{}, fmt.Errorf("invalid metrics enabled: %w", err)
	}

	tracingEnabled, err := parseBoolEnv("CLUO_TRACING_ENABLED", true)
	if err != nil {
		return ObservabilityConfig{}, fmt.Errorf("invalid tracing enabled: %w", err)
	}

	tracingSampler, err := parseFloatEnv("CLUO_TRACING_SAMPLER", 1.0)
	if err != nil {
		return ObservabilityConfig{}, fmt.Errorf("invalid tracing sampler: %w", err)
	}

	return ObservabilityConfig{
		LogLevel:        getEnv("CLUO_LOG_LEVEL", defaultLogLevel),
		LogStyle:        getEnv("CLUO_LOG_STYLE", defaultLogStyle),
		MetricsEnabled:  metricsEnabled,
		MetricsPath:     getEnv("CLUO_METRICS_PATH", "/metrics"),
		TracingEnabled:  tracingEnabled,
		TracingEndpoint: getEnv("CLUO_TRACING_ENDPOINT", ""),
		TracingSampler:  tracingSampler,
		ServiceName:     getEnv("CLUO_SERVICE_NAME", "cluo-api"),
		ServiceVersion:  getEnv("CLUO_SERVICE_VERSION", "unknown"),
	}, nil
}

// parseBoolEnv parses a boolean from an environment variable.
func parseBoolEnv(key string, defaultValue bool) (bool, error) {
	strVal := getEnv(key, "")
	if strVal == "" {
		return defaultValue, nil
	}

	switch strVal {
	case "true", "1", "yes", "on":
		return true, nil
	case "false", "0", "no", "off":
		return false, nil
	default:
		return false, fmt.Errorf("parse %s: invalid boolean value %q", key, strVal)
	}
}

// parseFloatEnv parses a float from an environment variable.
func parseFloatEnv(key string, defaultValue float64) (float64, error) {
	strVal := getEnv(key, "")
	if strVal == "" {
		return defaultValue, nil
	}

	var val float64
	_, err := fmt.Sscanf(strVal, "%f", &val)
	if err != nil {
		return 0, fmt.Errorf("parse %s: %w", key, err)
	}

	return val, nil
}
