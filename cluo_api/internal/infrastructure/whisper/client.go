package whisper

import (
	"context"
	"encoding/json"
	"fmt"
	"log/slog"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"time"

	"github.com/hengadev/cluo_api/internal/app/config"
	"github.com/hengadev/cluo_api/internal/ports"
)

// Client implements the WhisperClient interface using whisper.cpp.
type Client struct {
	binary    string
	modelPath string
	model     string
	language  string
	threads   int
	timeout   time.Duration
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
		language:  cfg.Language,
		threads:   cfg.Threads,
		timeout:   cfg.Timeout,
	}, nil
}

// Transcribe transcribes an audio file and returns the transcript.
func (c *Client) Transcribe(ctx context.Context, audioPath string) (*ports.WhisperResult, error) {
	// Create a context with timeout
	timeoutCtx, cancel := context.WithTimeout(ctx, c.timeout)
	defer cancel()

	// whisper-cli only handles WAV (16 kHz, mono, 16-bit PCM) natively.
	// Convert any other format (webm, ogg, mp3, …) via ffmpeg first.
	inputPath := audioPath
	if strings.ToLower(filepath.Ext(audioPath)) != ".wav" {
		wavPath := strings.TrimSuffix(audioPath, filepath.Ext(audioPath)) + "_converted.wav"
		defer os.Remove(wavPath)

		ffmpegStart := time.Now()
		conv := exec.CommandContext(timeoutCtx, "ffmpeg",
			"-y",                // overwrite if exists
			"-i", audioPath,     // input
			"-ar", "16000",      // 16 kHz sample rate
			"-ac", "1",          // mono
			"-c:a", "pcm_s16le", // 16-bit PCM
			wavPath,
		)
		out, err := conv.CombinedOutput()
		ffmpegElapsed := time.Since(ffmpegStart)
		if err != nil {
			return nil, fmt.Errorf("ffmpeg conversion failed after %s: %w, output: %s", ffmpegElapsed, err, string(out))
		}
		slog.Info("ffmpeg conversion done", "elapsed", ffmpegElapsed, "audioPath", audioPath, "output", strings.TrimSpace(string(out)))
		inputPath = wavPath
	}

	// whisper-cli writes its JSON result to "<inputPath_without_ext>.json"
	// (it strips the input extension before appending .json), so we must
	// do the same when computing the sidecar path.
	jsonPath := strings.TrimSuffix(inputPath, filepath.Ext(inputPath)) + ".json"
	defer os.Remove(jsonPath)

	args := []string{
		"-m", c.modelPath + "/" + c.model + ".ggml",
		"-f", inputPath,
		"-l", c.language,
		"-t", fmt.Sprintf("%d", c.threads),
		"--output-json",
	}

	slog.Info("starting whisper transcription", "audioPath", inputPath, "model", c.model, "threads", c.threads)
	whisperStart := time.Now()
	cmd := exec.CommandContext(timeoutCtx, c.binary, args...)
	output, err := cmd.CombinedOutput()
	whisperElapsed := time.Since(whisperStart)
	if err != nil {
		return nil, fmt.Errorf("whisper command failed after %s: %w, output: %s", whisperElapsed, err, string(output))
	}
	slog.Info("whisper transcription done", "elapsed", whisperElapsed, "audioPath", inputPath)

	raw, err := os.ReadFile(jsonPath)
	if err != nil {
		return nil, fmt.Errorf("read whisper JSON output: %w, command output: %s", err, string(output))
	}

	result, err := c.parseOutput(raw)
	if err != nil {
		return nil, fmt.Errorf("parse whisper output: %w", err)
	}

	return result, nil
}

// whisperJSONOutput represents the JSON file whisper-cli writes alongside
// the input audio when run with --output-json.
type whisperJSONOutput struct {
	Result struct {
		Language string `json:"language"`
	} `json:"result"`
	Transcription []whisperSegment `json:"transcription"`
}

type whisperSegment struct {
	Offsets struct {
		From int64 `json:"from"`
		To   int64 `json:"to"`
	} `json:"offsets"`
	Text string `json:"text"`
}

// parseOutput parses the JSON file produced by whisper.cpp.
func (c *Client) parseOutput(raw []byte) (*ports.WhisperResult, error) {
	var whisperOut whisperJSONOutput
	if err := json.Unmarshal(raw, &whisperOut); err != nil {
		return nil, fmt.Errorf("unmarshal JSON: %w", err)
	}

	var textBuilder strings.Builder
	var durationMs int64
	for _, segment := range whisperOut.Transcription {
		textBuilder.WriteString(segment.Text)
		if segment.Offsets.To > durationMs {
			durationMs = segment.Offsets.To
		}
	}

	confidence := c.calculateConfidence(whisperOut.Transcription, textBuilder.Len())

	return &ports.WhisperResult{
		Transcript:      strings.TrimSpace(textBuilder.String()),
		ConfidenceScore: confidence,
		Language:        whisperOut.Result.Language,
		Duration:        durationMs,
	}, nil
}

// calculateConfidence calculates a confidence score based on segment data.
func (c *Client) calculateConfidence(segments []whisperSegment, totalLength int) float32 {
	// If no segments, return default confidence
	if len(segments) == 0 {
		return 0.8 // Default confidence
	}

	if totalLength == 0 {
		return 0.0
	}

	// More segments with reasonable text = higher confidence
	avgSegmentLength := totalLength / len(segments)

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

