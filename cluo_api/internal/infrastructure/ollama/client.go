package ollama

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net"
	"net/http"
	"net/url"
	"strings"
	"time"

	"github.com/hengadev/cluo_api/internal/app/config"
)

const (
	// generateEndpoint is the Ollama API endpoint for text generation.
	generateEndpoint = "/api/generate"
	// tagsEndpoint is the Ollama API endpoint for listing models.
	tagsEndpoint = "/api/tags"
)

// Client implements the LLMClient interface for Ollama.
type Client struct {
	baseURL    string
	model      string
	httpClient *http.Client
	timeout    time.Duration
}

// New creates a new Ollama client.
func New(cfg config.OllamaConfig) (*Client, error) {
	if !cfg.Enabled {
		return nil, nil // Disabled, not an error
	}

	if err := validateLocalhostURL(cfg.BaseURL); err != nil {
		return nil, fmt.Errorf("invalid ollama base URL: %w", err)
	}

	return &Client{
		baseURL: cfg.BaseURL,
		model:   cfg.Model,
		httpClient: &http.Client{
			Timeout: cfg.Timeout,
		},
		timeout: cfg.Timeout,
	}, nil
}

// GenerateRequest represents a request to the Ollama generate API.
type GenerateRequest struct {
	Model    string `json:"model"`
	Prompt   string `json:"prompt"`
	System   string `json:"system,omitempty"`
	Stream   bool   `json:"stream"`
	Options  *Options `json:"options,omitempty"`
}

// Options represents additional generation options.
type Options struct {
	Temperature float32 `json:"temperature,omitempty"`
	TopP        float32 `json:"top_p,omitempty"`
	TopK        int     `json:"top_k,omitempty"`
	NumPredict  int     `json:"num_predict,omitempty"`
}

// GenerateResponse represents a response from the Ollama generate API.
type GenerateResponse struct {
	Model     string    `json:"model"`
	Response  string    `json:"response"`
	Done      bool      `json:"done"`
	Context   []int     `json:"context,omitempty"`
	TotalDuration int64 `json:"total_duration,omitempty"`
	LoadDuration  int64 `json:"load_duration,omitempty"`
	PromptEvalCount int `json:"prompt_eval_count,omitempty"`
	PromptEvalDuration int64 `json:"prompt_eval_duration,omitempty"`
	EvalCount      int   `json:"eval_count,omitempty"`
	EvalDuration   int64 `json:"eval_duration,omitempty"`
	Error          string `json:"error,omitempty"`
}

// Generate sends a prompt to Ollama and returns the generated text.
func (c *Client) Generate(ctx context.Context, prompt string, systemPrompt string) (string, error) {
	if c == nil {
		return "", fmt.Errorf("ollama client is not initialized")
	}

	// Build request
	req := GenerateRequest{
		Model:  c.model,
		Prompt: prompt,
		System: systemPrompt,
		Stream: false,
		Options: &Options{
			Temperature: 0.1, // Low temperature for more deterministic output
			TopP:        0.9,
		},
	}

	reqBody, err := json.Marshal(req)
	if err != nil {
		return "", fmt.Errorf("marshal request: %w", err)
	}

	// Create HTTP request
	url := c.baseURL + generateEndpoint
	httpReq, err := http.NewRequestWithContext(ctx, http.MethodPost, url, bytes.NewReader(reqBody))
	if err != nil {
		return "", fmt.Errorf("create request: %w", err)
	}

	httpReq.Header.Set("Content-Type", "application/json")

	// Execute request
	resp, err := c.httpClient.Do(httpReq)
	if err != nil {
		return "", fmt.Errorf("execute request: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		return "", fmt.Errorf("ollama returned status %d: %s", resp.StatusCode, string(body))
	}

	// Parse response
	var genResp GenerateResponse
	if err := json.NewDecoder(resp.Body).Decode(&genResp); err != nil {
		return "", fmt.Errorf("decode response: %w", err)
	}

	if genResp.Error != "" {
		return "", fmt.Errorf("ollama error: %s", genResp.Error)
	}

	// Trim whitespace from response
	result := strings.TrimSpace(genResp.Response)

	return result, nil
}

// HealthCheck checks if Ollama is available.
func (c *Client) HealthCheck(ctx context.Context) error {
	if c == nil {
		return fmt.Errorf("ollama client is not initialized")
	}

	// Use tags endpoint to check if Ollama is running
	url := c.baseURL + tagsEndpoint

	httpReq, err := http.NewRequestWithContext(ctx, http.MethodGet, url, nil)
	if err != nil {
		return fmt.Errorf("create request: %w", err)
	}

	resp, err := c.httpClient.Do(httpReq)
	if err != nil {
		return fmt.Errorf("execute request: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("ollama returned status %d", resp.StatusCode)
	}

	return nil
}

// validateLocalhostURL validates that the URL points to localhost or a
// private network address (e.g. a sibling Docker Compose service such as
// "ollama"), never a public host — Ollama must not receive transcript data
// over the public internet.
func validateLocalhostURL(rawURL string) error {
	parsed, err := url.Parse(rawURL)
	if err != nil {
		return fmt.Errorf("invalid URL: %w", err)
	}

	hostname := strings.ToLower(parsed.Hostname())
	if hostname == "" {
		return fmt.Errorf("URL must have a hostname")
	}

	if hostname == "localhost" || strings.HasSuffix(hostname, ".localhost") {
		return nil
	}

	ip := net.ParseIP(hostname)
	if ip == nil {
		ips, err := net.LookupIP(hostname)
		if err != nil {
			return fmt.Errorf("resolve host %q: %w", hostname, err)
		}
		if len(ips) == 0 {
			return fmt.Errorf("host %q did not resolve to an address", hostname)
		}
		ip = ips[0]
	}

	if !ip.IsLoopback() && !ip.IsPrivate() {
		return fmt.Errorf("ollama must be localhost or a private network address, got: %s", hostname)
	}

	return nil
}
