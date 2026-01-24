package config

import (
	"fmt"
	"strings"

	"github.com/hengadev/cluo_api/internal/common/envmode"
)

// ValidationError represents a configuration validation error.
type ValidationError struct {
	Field   string
	Message string
}

func (e ValidationError) Error() string {
	return fmt.Sprintf("%s: %s", e.Field, e.Message)
}

// ValidationErrors is a collection of validation errors.
type ValidationErrors []ValidationError

func (e ValidationErrors) Error() string {
	if len(e) == 0 {
		return ""
	}

	var msgs []string
	for _, err := range e {
		msgs = append(msgs, err.Error())
	}
	return fmt.Sprintf("configuration validation failed:\n  - %s", strings.Join(msgs, "\n  - "))
}

// Validate validates the entire configuration.
func (c *Config) Validate() error {
	var errs ValidationErrors

	// Validate server config
	if c.Server.Port < 1 || c.Server.Port > 65535 {
		errs = append(errs, ValidationError{
			Field:   "CLUO_SERVER_PORT",
			Message: "must be between 1 and 65535",
		})
	}

	// Validate database config (required in production)
	if c.Environment == envmode.Prod || c.Environment == envmode.Staging {
		if c.Database.Host == "" {
			errs = append(errs, ValidationError{
				Field:   "CLUO_DB_HOST",
				Message: "required in production/staging",
			})
		}
		if c.Database.Name == "" {
			errs = append(errs, ValidationError{
				Field:   "CLUO_DB_NAME",
				Message: "required in production/staging",
			})
		}
		if c.Database.User == "" {
			errs = append(errs, ValidationError{
				Field:   "CLUO_DB_USER",
				Message: "required in production/staging",
			})
		}
		if c.Database.Password == "" {
			errs = append(errs, ValidationError{
				Field:   "CLUO_DB_PASSWORD",
				Message: "required in production/staging",
			})
		}
	}

	// Validate Redis config (required in production)
	if c.Environment == envmode.Prod || c.Environment == envmode.Staging {
		if c.Redis.Host == "" {
			errs = append(errs, ValidationError{
				Field:   "CLUO_REDIS_HOST",
				Message: "required in production/staging",
			})
		}
	}

	// Validate Vault config (required in production)
	if c.Environment == envmode.Prod || c.Environment == envmode.Staging {
		if c.Vault.Address == "" {
			errs = append(errs, ValidationError{
				Field:   "VAULT_ADDR",
				Message: "required in production/staging",
			})
		}
		// Must have either token or AppRole credentials
		if c.Vault.Token == "" && !c.Vault.UseAppRole() {
			errs = append(errs, ValidationError{
				Field:   "VAULT_TOKEN or VAULT_APPROLE_*",
				Message: "Vault authentication required in production/staging",
			})
		}
	}

	// Validate observability config
	validLogLevels := map[string]bool{"debug": true, "info": true, "warn": true, "error": true}
	if !validLogLevels[c.Observability.LogLevel] {
		errs = append(errs, ValidationError{
			Field:   "CLUO_LOG_LEVEL",
			Message: "must be one of: debug, info, warn, error",
		})
	}

	validLogStyles := map[string]bool{"json": true, "text": true, "dev": true}
	if !validLogStyles[c.Observability.LogStyle] {
		errs = append(errs, ValidationError{
			Field:   "CLUO_LOG_STYLE",
			Message: "must be one of: json, text, dev",
		})
	}

	// Validate tracing config
	if c.Observability.TracingEnabled && c.Observability.TracingEndpoint == "" {
		// Only warn in development - tracing endpoint is required in prod
		if c.Environment == envmode.Prod {
			errs = append(errs, ValidationError{
				Field:   "CLUO_TRACING_ENDPOINT",
				Message: "required when tracing is enabled in production",
			})
		}
	}

	if c.Observability.TracingSampler < 0 || c.Observability.TracingSampler > 1 {
		errs = append(errs, ValidationError{
			Field:   "CLUO_TRACING_SAMPLER",
			Message: "must be between 0 and 1",
		})
	}

	if len(errs) > 0 {
		return errs
	}

	return nil
}

// IsDevelopment returns true if running in development mode.
func (c *Config) IsDevelopment() bool {
	return c.Environment == envmode.Dev
}

// IsProduction returns true if running in production mode.
func (c *Config) IsProduction() bool {
	return c.Environment == envmode.Prod
}
