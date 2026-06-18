package aiTranscriptAnalysis

import (
	"context"
	"encoding/json"
	"fmt"
	"log/slog"
	"strings"
	"time"

	"github.com/google/uuid"
	"github.com/hengadev/cluo_api/internal/domain/ai"
	"github.com/hengadev/cluo_api/internal/ports"
	"github.com/hengadev/encx"
)

// Service implements the TranscriptAnalysisService interface.
type Service struct {
	llmClient         ports.LLMClient
	transcriptionRepo ports.TranscriptionRepository
	analysisRepo      ports.TranscriptAnalysisRepository
	crypto            encx.CryptoService
	logger            *slog.Logger
}

// New creates a new transcript analysis service.
func New(
	llmClient ports.LLMClient,
	transcriptionRepo ports.TranscriptionRepository,
	analysisRepo ports.TranscriptAnalysisRepository,
	crypto encx.CryptoService,
	logger *slog.Logger,
) *Service {
	return &Service{
		llmClient:         llmClient,
		transcriptionRepo: transcriptionRepo,
		analysisRepo:      analysisRepo,
		crypto:            crypto,
		logger:            logger.With("component", "ai_transcript_analysis"),
	}
}

// AnalyzeTranscript performs analysis on a transcript.
func (s *Service) AnalyzeTranscript(ctx context.Context, req *ports.AnalyzeTranscriptRequest) (*ai.TranscriptAnalysis, error) {
	// Get the encrypted transcription
	transcriptionEncx, err := s.transcriptionRepo.GetByID(ctx, req.TranscriptionID)
	if err != nil {
		return nil, fmt.Errorf("failed to get transcription: %w", err)
	}

	if transcriptionEncx == nil {
		return nil, fmt.Errorf("transcription not found")
	}

	// Decrypt the transcription
	transcription, err := ai.DecryptTranscriptionEncx(ctx, s.crypto, transcriptionEncx)
	if err != nil {
		return nil, fmt.Errorf("failed to decrypt transcription: %w", err)
	}

	s.logger.Debug("analyzing transcript",
		"transcriptionID", req.TranscriptionID,
		"transcriptLength", len(transcription.Transcript))

	// Build the analysis prompt
	prompt := s.buildAnalysisPrompt(transcription.Transcript)

	// System prompt for transcript analysis
	systemPrompt := `You are an expert analyst for investigative transcripts. Your task is to:
1. Identify key findings and important details
2. Provide a concise summary
3. Determine the overall sentiment (positive, neutral, negative, mixed)
4. Extract main topics as a JSON array
5. Suggest relevant actions based on the content

Respond in JSON format with the following structure:
{
  "keyFindings": "bullet points of key findings",
  "summary": "concise summary",
  "sentiment": "positive|neutral|negative|mixed",
  "topics": ["topic1", "topic2", ...],
  "suggestedActions": "actionable recommendations"
}`

	// Call LLM
	startTime := time.Now()
	response, err := s.llmClient.Generate(ctx, prompt, systemPrompt)
	processingTimeMs := time.Since(startTime).Milliseconds()

	if err != nil {
		s.logger.Error("LLM analysis failed", "error", err)
		return nil, fmt.Errorf("LLM analysis failed: %w", err)
	}

	// Parse response
	analysis, err := s.parseAnalysisResponse(response)
	if err != nil {
		s.logger.Error("failed to parse analysis response", "error", err)
		return nil, fmt.Errorf("failed to parse analysis: %w", err)
	}

	s.logger.Debug("transcript analysis completed",
		"transcriptionID", req.TranscriptionID,
		"sentiment", analysis.Sentiment,
		"processingTimeMs", processingTimeMs)

	// Create result entity
	result := ai.NewTranscriptAnalysis(
		req.TranscriptionID,
		analysis.KeyFindings,
		analysis.Summary,
		analysis.Sentiment,
		analysis.TopicsJSON,
		analysis.SuggestedActions,
		"llm",
		processingTimeMs,
	)

	// Persist the analysis (encrypted). A transcription has at most one
	// analysis, so re-analyzing replaces the previous result.
	resultEncx, err := ai.ProcessTranscriptAnalysisEncx(ctx, s.crypto, result)
	if err != nil {
		return nil, fmt.Errorf("failed to encrypt analysis: %w", err)
	}
	if err := s.analysisRepo.Create(ctx, resultEncx); err != nil {
		return nil, fmt.Errorf("failed to persist analysis: %w", err)
	}

	return result, nil
}

// GetAnalysis retrieves an analysis by ID.
func (s *Service) GetAnalysis(ctx context.Context, id uuid.UUID) (*ai.TranscriptAnalysis, error) {
	analysisEncx, err := s.analysisRepo.GetByID(ctx, id)
	if err != nil {
		return nil, fmt.Errorf("failed to get analysis: %w", err)
	}
	return ai.DecryptTranscriptAnalysisEncx(ctx, s.crypto, analysisEncx)
}

// GetAnalysisByTranscriptionID retrieves analysis by transcription ID.
func (s *Service) GetAnalysisByTranscriptionID(ctx context.Context, transcriptionID uuid.UUID) (*ai.TranscriptAnalysis, error) {
	analysisEncx, err := s.analysisRepo.GetByTranscriptionID(ctx, transcriptionID)
	if err != nil {
		return nil, fmt.Errorf("failed to get analysis: %w", err)
	}
	return ai.DecryptTranscriptAnalysisEncx(ctx, s.crypto, analysisEncx)
}

// ListAnalyses retrieves analyses with pagination.
func (s *Service) ListAnalyses(ctx context.Context, req *ports.ListAnalysesRequest) ([]*ai.TranscriptAnalysis, int, error) {
	analysesEncx, total, err := s.analysisRepo.List(ctx, req.TranscriptionID, req.Limit, req.Offset)
	if err != nil {
		return nil, 0, fmt.Errorf("failed to list analyses: %w", err)
	}

	analyses := make([]*ai.TranscriptAnalysis, len(analysesEncx))
	for i, encx := range analysesEncx {
		analysis, err := ai.DecryptTranscriptAnalysisEncx(ctx, s.crypto, encx)
		if err != nil {
			return nil, 0, fmt.Errorf("failed to decrypt analysis: %w", err)
		}
		analyses[i] = analysis
	}

	return analyses, total, nil
}

// DeleteAnalysis deletes an analysis.
func (s *Service) DeleteAnalysis(ctx context.Context, id uuid.UUID) error {
	return s.analysisRepo.Delete(ctx, id)
}

// buildAnalysisPrompt builds the prompt for transcript analysis.
func (s *Service) buildAnalysisPrompt(transcript string) string {
	// Truncate transcript if too long for the context
	maxLength := 8000 // Leave room for the system prompt
	if len(transcript) > maxLength {
		transcript = transcript[:maxLength] + "..."
	}

	return fmt.Sprintf("Please analyze the following transcript:\n\n%s", transcript)
}

// parseAnalysisResponse parses the LLM response into analysis components.
type parsedAnalysis struct {
	KeyFindings      string
	Summary          string
	Sentiment        string
	TopicsJSON       string
	SuggestedActions string
}

func (s *Service) parseAnalysisResponse(response string) (*parsedAnalysis, error) {
	// Try to extract JSON from the response
	response = strings.TrimSpace(response)

	// Remove markdown code blocks if present
	if strings.HasPrefix(response, "```json") {
		response = strings.TrimPrefix(response, "```json")
		response = strings.TrimSuffix(response, "```")
		response = strings.TrimSpace(response)
	} else if strings.HasPrefix(response, "```") {
		response = strings.TrimPrefix(response, "```")
		response = strings.TrimSuffix(response, "```")
		response = strings.TrimSpace(response)
	}

	// Parse JSON
	var result map[string]interface{}
	if err := json.Unmarshal([]byte(response), &result); err != nil {
		// If JSON parsing fails, try to extract components using regex or heuristics
		return s.parseAnalysisFallback(response)
	}

	analysis := &parsedAnalysis{}

	// Extract fields
	if v, ok := result["keyFindings"].(string); ok {
		analysis.KeyFindings = v
	}
	if v, ok := result["summary"].(string); ok {
		analysis.Summary = v
	}
	if v, ok := result["sentiment"].(string); ok {
		analysis.Sentiment = v
	}
	if v, ok := result["topics"].([]interface{}); ok {
		topics := make([]string, len(v))
		for i, t := range v {
			if str, ok := t.(string); ok {
				topics[i] = str
			}
		}
		topicsJSON, _ := json.Marshal(topics)
		analysis.TopicsJSON = string(topicsJSON)
	}
	if v, ok := result["suggestedActions"].(string); ok {
		analysis.SuggestedActions = v
	}

	// Validate sentiment
	if !ai.Sentiment(analysis.Sentiment).Valid() {
		analysis.Sentiment = "neutral"
	}

	return analysis, nil
}

// parseAnalysisFallback is a fallback parser when JSON parsing fails.
func (s *Service) parseAnalysisFallback(response string) (*parsedAnalysis, error) {
	// This is a simple fallback that attempts to extract structured data
	// In a production system, you might want to use a more sophisticated parser
	// or ask the LLM to retry with a strict JSON format

	return &parsedAnalysis{
		KeyFindings:      response,
		Summary:          "Analysis available in key findings",
		Sentiment:        "neutral",
		TopicsJSON:       "[]",
		SuggestedActions: "Review key findings",
	}, nil
}

// HealthCheck checks if the analysis service is healthy.
func (s *Service) HealthCheck(ctx context.Context) error {
	if s.llmClient == nil {
		return fmt.Errorf("LLM client not initialized")
	}
	return s.llmClient.HealthCheck(ctx)
}
