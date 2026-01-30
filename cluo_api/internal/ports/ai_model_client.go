package ports

import (
	"context"
)

// LLMClient defines the interface for Large Language Model operations (e.g., Ollama).
type LLMClient interface {
	// Generate sends a prompt to the LLM and returns the generated text.
	Generate(ctx context.Context, prompt string, systemPrompt string) (string, error)

	// HealthCheck checks if the LLM service is available.
	HealthCheck(ctx context.Context) error
}

// WhisperClient defines the interface for Whisper speech-to-text operations.
type WhisperClient interface {
	// Transcribe transcribes an audio file and returns the transcript.
	Transcribe(ctx context.Context, audioPath string) (*WhisperResult, error)

	// HealthCheck checks if the Whisper service is available.
	HealthCheck(ctx context.Context) error
}

// WhisperResult represents the result of a transcription.
type WhisperResult struct {
	Transcript      string
	ConfidenceScore float32
	Language        string
	Duration        int64 // milliseconds
}
