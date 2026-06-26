package config

import (
	"fmt"
	"net"
	"net/url"
	"strings"
	"time"
)

// AIConfig holds AI services configuration.
type AIConfig struct {
	Ollama  OllamaConfig
	Whisper WhisperConfig
	Worker  WorkerConfig
}

// OllamaConfig holds Ollama (local LLM) configuration.
type OllamaConfig struct {
	Enabled    bool
	BaseURL    string
	Model      string
	Timeout    time.Duration
	MaxInputLen int
}

// WhisperConfig holds Whisper.cpp (local speech-to-text) configuration.
type WhisperConfig struct {
	Enabled     bool
	Binary      string
	ModelPath   string
	Model       string
	Language    string
	Threads     int
	Timeout     time.Duration
	DeleteAudio bool
}

// WorkerConfig holds background worker configuration for async transcription.
type WorkerConfig struct {
	Concurrency  int
	QueueSize    int
	RetryAttempts int
	RetryDelay   time.Duration
}

func loadAIConfig() (AIConfig, error) {
	ollama, err := loadOllamaConfig()
	if err != nil {
		return AIConfig{}, fmt.Errorf("ollama config: %w", err)
	}

	whisper, err := loadWhisperConfig()
	if err != nil {
		return AIConfig{}, fmt.Errorf("whisper config: %w", err)
	}

	worker, err := loadWorkerConfig()
	if err != nil {
		return AIConfig{}, fmt.Errorf("worker config: %w", err)
	}

	return AIConfig{
		Ollama:  ollama,
		Whisper: whisper,
		Worker:  worker,
	}, nil
}

func loadOllamaConfig() (OllamaConfig, error) {
	enabled, err := parseBoolEnv("CLUO_AI_OLLAMA_ENABLED", true)
	if err != nil {
		return OllamaConfig{}, fmt.Errorf("invalid ollama enabled: %w", err)
	}

	baseURL := getEnv("CLUO_AI_OLLAMA_BASE_URL", "http://localhost:11434")
	if err := validateLocalhostURL(baseURL); err != nil {
		return OllamaConfig{}, fmt.Errorf("invalid ollama base URL: %w", err)
	}

	timeout, err := parseDurationEnv("CLUO_AI_OLLAMA_TIMEOUT", 30*time.Second)
	if err != nil {
		return OllamaConfig{}, err
	}

	maxInputLen, err := parseIntEnv("CLUO_AI_OLLAMA_MAX_INPUT_LEN", 10000)
	if err != nil {
		return OllamaConfig{}, err
	}

	return OllamaConfig{
		Enabled:     enabled,
		BaseURL:     baseURL,
		Model:       getEnv("CLUO_AI_OLLAMA_MODEL", "llama3.2"),
		Timeout:     timeout,
		MaxInputLen: maxInputLen,
	}, nil
}

func loadWhisperConfig() (WhisperConfig, error) {
	enabled, err := parseBoolEnv("CLUO_AI_WHISPER_ENABLED", true)
	if err != nil {
		return WhisperConfig{}, fmt.Errorf("invalid whisper enabled: %w", err)
	}

	binary := getEnv("CLUO_AI_WHISPER_BINARY", "/usr/local/bin/whisper")

	timeout, err := parseDurationEnv("CLUO_AI_WHISPER_TIMEOUT", 10*time.Minute)
	if err != nil {
		return WhisperConfig{}, err
	}

	threads, err := parseIntEnv("CLUO_AI_WHISPER_THREADS", 4)
	if err != nil {
		return WhisperConfig{}, err
	}

	deleteAudio, err := parseBoolEnv("CLUO_AI_WHISPER_DELETE_AUDIO", true)
	if err != nil {
		return WhisperConfig{}, fmt.Errorf("invalid delete audio: %w", err)
	}

	return WhisperConfig{
		Enabled:     enabled,
		Binary:      binary,
		ModelPath:   getEnv("CLUO_AI_WHISPER_MODEL_PATH", "/models/whisper"),
		Model:       getEnv("CLUO_AI_WHISPER_MODEL", "base"),
		Language:    getEnv("CLUO_AI_WHISPER_LANGUAGE", "fr"),
		Threads:     threads,
		Timeout:     timeout,
		DeleteAudio: deleteAudio,
	}, nil
}

func loadWorkerConfig() (WorkerConfig, error) {
	concurrency, err := parseIntEnv("CLUO_AI_WORKER_CONCURRENCY", 2)
	if err != nil {
		return WorkerConfig{}, err
	}

	queueSize, err := parseIntEnv("CLUO_AI_WORKER_QUEUE_SIZE", 100)
	if err != nil {
		return WorkerConfig{}, err
	}

	retryAttempts, err := parseIntEnv("CLUO_AI_WORKER_RETRY_ATTEMPTS", 3)
	if err != nil {
		return WorkerConfig{}, err
	}

	retryDelay, err := parseDurationEnv("CLUO_AI_WORKER_RETRY_DELAY", 30*time.Second)
	if err != nil {
		return WorkerConfig{}, err
	}

	return WorkerConfig{
		Concurrency:   concurrency,
		QueueSize:     queueSize,
		RetryAttempts: retryAttempts,
		RetryDelay:    retryDelay,
	}, nil
}

// validateLocalhostURL validates that the URL points to localhost or a
// private network address (e.g. a sibling Docker Compose service such as
// "ollama"), never a public host — AI services must not receive transcript
// data over the public internet.
func validateLocalhostURL(rawURL string) error {
	u, err := url.Parse(rawURL)
	if err != nil {
		return fmt.Errorf("parse URL: %w", err)
	}

	host := strings.ToLower(u.Hostname())
	if host == "" {
		return fmt.Errorf("URL must have a hostname")
	}

	if host == "localhost" || strings.HasSuffix(host, ".localhost") {
		return nil
	}

	ip := net.ParseIP(host)
	if ip == nil {
		ips, err := net.LookupIP(host)
		if err != nil {
			return fmt.Errorf("resolve host %q: %w", host, err)
		}
		if len(ips) == 0 {
			return fmt.Errorf("host %q did not resolve to an address", host)
		}
		ip = ips[0]
	}

	if !ip.IsLoopback() && !ip.IsPrivate() {
		return fmt.Errorf("AI service must be localhost or a private network address, got: %s", host)
	}

	return nil
}
