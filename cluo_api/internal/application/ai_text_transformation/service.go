package aiTextTransformation

import (
	"context"
	"fmt"
	"log/slog"
	"time"

	"github.com/hengadev/cluo_api/internal/domain/ai"
	"github.com/hengadev/cluo_api/internal/ports"
)

// Service implements the TextTransformationService interface.
type Service struct {
	llmClient  ports.LLMClient
	maxInputLen int
	logger     *slog.Logger
}

// New creates a new text transformation service.
func New(llmClient ports.LLMClient, maxInputLen int, logger *slog.Logger) *Service {
	return &Service{
		llmClient:   llmClient,
		maxInputLen: maxInputLen,
		logger:      logger.With("component", "ai_text_transformation"),
	}
}

// TransformText applies a transformation to input text.
func (s *Service) TransformText(ctx context.Context, req *ports.TransformTextRequest) (*ai.TextTransformation, error) {
	// Validate input
	if err := s.validateRequest(req); err != nil {
		return nil, fmt.Errorf("validation failed: %w", err)
	}

	// Get the system prompt for the transformation type
	systemPrompt := req.TransformationType.Prompt()

	s.logger.Debug("transforming text",
		"type", req.TransformationType,
		"inputLength", len(req.InputText))

	// Call LLM
	startTime := time.Now()
	outputText, err := s.llmClient.Generate(ctx, req.InputText, systemPrompt)
	processingTimeMs := time.Since(startTime).Milliseconds()

	if err != nil {
		s.logger.Error("LLM generation failed", "error", err)
		return nil, fmt.Errorf("LLM generation failed: %w", err)
	}

	if outputText == "" {
		return nil, ai.ErrInvalidResponse
	}

	s.logger.Debug("text transformation completed",
		"inputLength", len(req.InputText),
		"outputLength", len(outputText),
		"processingTimeMs", processingTimeMs)

	// Create result entity
	result := ai.NewTextTransformation(
		req.InputText,
		outputText,
		req.TransformationType,
		"llm", // ModelUsed - could be more specific from the LLM client
		processingTimeMs,
	)

	return result, nil
}

// validateRequest validates the transformation request.
func (s *Service) validateRequest(req *ports.TransformTextRequest) error {
	if req.InputText == "" {
		return ai.ErrEmptyInput
	}

	if len(req.InputText) > s.maxInputLen {
		return fmt.Errorf("%w: maximum length is %d", ai.ErrTextTooLong, s.maxInputLen)
	}

	if !req.TransformationType.Valid() {
		return fmt.Errorf("%w: %s", ai.ErrTransformationTypeInvalid, req.TransformationType)
	}

	return nil
}

// HealthCheck checks if the LLM service is available.
func (s *Service) HealthCheck(ctx context.Context) error {
	if s.llmClient == nil {
		return fmt.Errorf("LLM client not initialized")
	}
	return s.llmClient.HealthCheck(ctx)
}
