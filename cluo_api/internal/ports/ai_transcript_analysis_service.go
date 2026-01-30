package ports

import (
	"context"

	"github.com/google/uuid"
	"github.com/hengadev/cluo_api/internal/domain/ai"
)

// TranscriptAnalysisService defines the interface for transcript analysis operations.
type TranscriptAnalysisService interface {
	// AnalyzeTranscript performs analysis on a transcript.
	AnalyzeTranscript(ctx context.Context, req *AnalyzeTranscriptRequest) (*ai.TranscriptAnalysis, error)

	// GetAnalysis retrieves an analysis by ID.
	GetAnalysis(ctx context.Context, id uuid.UUID) (*ai.TranscriptAnalysis, error)

	// GetAnalysisByTranscriptionID retrieves analysis by transcription ID.
	GetAnalysisByTranscriptionID(ctx context.Context, transcriptionID uuid.UUID) (*ai.TranscriptAnalysis, error)

	// ListAnalyses retrieves analyses with pagination.
	ListAnalyses(ctx context.Context, req *ListAnalysesRequest) ([]*ai.TranscriptAnalysis, int, error)

	// DeleteAnalysis deletes an analysis.
	DeleteAnalysis(ctx context.Context, id uuid.UUID) error
}

// AnalyzeTranscriptRequest defines the request for analyzing a transcript.
type AnalyzeTranscriptRequest struct {
	TranscriptionID uuid.UUID
}

// ListAnalysesRequest defines the request for listing analyses.
type ListAnalysesRequest struct {
	TranscriptionID *uuid.UUID
	Limit           int
	Offset          int
}
