package config

import (
	"fmt"
	"time"
)

// RateLimitConfig holds per-route-group rate limiting configuration.
type RateLimitConfig struct {
	LoginRPM int
	TokenRPM int
}

func loadRateLimitConfig() (RateLimitConfig, error) {
	loginRPM, err := parseIntEnv("RATE_LIMIT_LOGIN_RPM", 10)
	if err != nil {
		return RateLimitConfig{}, fmt.Errorf("parse RATE_LIMIT_LOGIN_RPM: %w", err)
	}
	if loginRPM < 1 {
		return RateLimitConfig{}, fmt.Errorf("RATE_LIMIT_LOGIN_RPM must be >= 1, got %d", loginRPM)
	}

	tokenRPM, err := parseIntEnv("RATE_LIMIT_TOKEN_RPM", 60)
	if err != nil {
		return RateLimitConfig{}, fmt.Errorf("parse RATE_LIMIT_TOKEN_RPM: %w", err)
	}
	if tokenRPM < 1 {
		return RateLimitConfig{}, fmt.Errorf("RATE_LIMIT_TOKEN_RPM must be >= 1, got %d", tokenRPM)
	}

	return RateLimitConfig{
		LoginRPM: loginRPM,
		TokenRPM: tokenRPM,
	}, nil
}

// LoginWindow returns the per-minute config as the middleware expects it.
func (c RateLimitConfig) LoginWindow() (maxRequests int, window time.Duration) {
	return c.LoginRPM, time.Minute
}

// TokenWindow returns the per-minute config as the middleware expects it.
func (c RateLimitConfig) TokenWindow() (maxRequests int, window time.Duration) {
	return c.TokenRPM, time.Minute
}
