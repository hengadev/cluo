package whisper

import (
	"context"
	"encoding/json"
	"fmt"
	"os/exec"
	"strings"
	"time"

	"github.com/hengadev/cluo_api/internal/app/config"
	"github.com/hengadev/cluo_api/internal/ports"
)

const (
	// defaultLanguage is the default language for transcription.
	defaultLanguage = "en"
)

// Client implements the WhisperClient interface using whisper.cpp.
type Client struct {
	binary      string
	modelPath   string
	model       string
	timeout     time.Duration
}

// New creates a new Whisper client.
func New(cfg config.WhisperConfig) (*Client, error) {
	if !cfg.Enabled {
		return nil, nil // Disabled, not an error
	}

	// Validate binary exists
	if _, err := exec.LookPath(cfg.Binary); err != nil {
		return nil, fmt.Errorf("whisper binary not found at %s: %w", cfg.Binary, err)
	}

	return &Client{
		binary:    cfg.Binary,
		modelPath: cfg.ModelPath,
		model:     cfg.Model,
		timeout:   cfg.Timeout,
	}, nil
}

// Transcribe transcribes an audio file and returns the transcript.
func (c *Client) Transcribe(ctx context.Context, audioPath string) (*ports.WhisperResult, error) {
	// Create a context with timeout
	timeoutCtx, cancel := context.WithTimeout(ctx, c.timeout)
	defer cancel()

	// Build whisper command
	args := []string{
		"-m", c.modelPath + "/" + c.model + ".ggml",
		"-f", audioPath,
		"-l", defaultLanguage,
		"--output-json", // Output in JSON format
	}

	// Run whisper command
	cmd := exec.CommandContext(timeoutCtx, c.binary, args...)

	// Capture combined output (stdout + stderr)
	output, err := cmd.CombinedOutput()
	if err != nil {
		return nil, fmt.Errorf("whisper command failed: %w, output: %s", err, string(output))
	}

	// Parse JSON output
	result, err := c.parseOutput(string(output))
	if err != nil {
		return nil, fmt.Errorf("parse whisper output: %w", err)
	}

	return result, nil
}

// whisperJSONOutput represents the JSON output from whisper.cpp.
type whisperJSONOutput struct {
	Text      string             `json:"text"`
	Segments []whisperSegment    `json:"segments,omitempty"`
	Language  string             `json:"language,omitempty"`
	Duration  float64            `json:"duration,omitempty"`
}

type whisperSegment struct {
	ID        int     `json:"id"`
	Start     float64 `json:"start"`
	End       float64 `json:"end"`
	Text      string  `json:"text"`
}

// parseOutput parses the output from whisper.cpp.
func (c *Client) parseOutput(output string) (*ports.WhisperResult, error) {
	// The JSON output is typically in the last few lines
	// Split by lines and find the JSON object
	lines := strings.Split(output, "\n")

	var jsonOutput string
	// Look for JSON object (starts with {) and collect multi-line JSON
	for i, line := range lines {
		trimmed := strings.TrimSpace(line)
		if strings.HasPrefix(trimmed, "{") {
			// Collect lines starting from the opening brace
			var builder strings.Builder
			for j := i; j < len(lines); j++ {
				builder.WriteString(lines[j])
				builder.WriteString("\n")
				candidate := builder.String()
				if strings.Count(candidate, "{") <= strings.Count(candidate, "}") {
					jsonOutput = candidate
					break
				}
			}
			break
		}
	}

	if jsonOutput == "" {
		// Fallback: try to parse the entire output as JSON
		jsonOutput = output
	}

	var whisperOut whisperJSONOutput
	if err := json.Unmarshal([]byte(jsonOutput), &whisperOut); err != nil {
		return nil, fmt.Errorf("unmarshal JSON: %w", err)
	}

	// Build result
	durationMs := int64(whisperOut.Duration * 1000)
	confidence := c.calculateConfidence(whisperOut)

	return &ports.WhisperResult{
		Transcript:      strings.TrimSpace(whisperOut.Text),
		ConfidenceScore: confidence,
		Language:        whisperOut.Language,
		Duration:        durationMs,
	}, nil
}

// calculateConfidence calculates a confidence score based on segment data.
func (c *Client) calculateConfidence(output whisperJSONOutput) float32 {
	// If no segments, return default confidence
	if len(output.Segments) == 0 {
		return 0.8 // Default confidence
	}

	// Simple heuristic: confidence based on text length and segment count
	// In a real implementation, you might use actual token probabilities if available
	totalLength := len(output.Text)
	if totalLength == 0 {
		return 0.0
	}

	// More segments with reasonable text = higher confidence
	segmentCount := len(output.Segments)
	avgSegmentLength := totalLength / segmentCount

	// Heuristic: good transcriptions have segments of 5-50 characters
	if avgSegmentLength >= 5 && avgSegmentLength <= 50 {
		return 0.9
	}
	if avgSegmentLength > 0 && avgSegmentLength < 100 {
		return 0.8
	}

	return 0.7
}

// HealthCheck checks if whisper binary is available.
func (c *Client) HealthCheck(ctx context.Context) error {
	// Check if binary exists
	if _, err := exec.LookPath(c.binary); err != nil {
		return fmt.Errorf("whisper binary not found: %w", err)
	}

	// Try to run whisper with --help flag to verify it works
	cmd := exec.CommandContext(ctx, c.binary, "--help")
	output, err := cmd.CombinedOutput()
	if err != nil {
		return fmt.Errorf("whisper health check failed: %w, output: %s", err, string(output))
	}

	// Check if help output contains expected text
	if !strings.Contains(string(output), "whisper") && !strings.Contains(string(output), "usage") {
		return fmt.Errorf("whisper binary output unexpected: %s", string(output))
	}

	return nil
}

