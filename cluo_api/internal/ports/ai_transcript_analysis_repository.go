package ports

import (
	"context"

	"github.com/google/uuid"
	"github.com/hengadev/cluo_api/internal/domain/ai"
)

// TranscriptAnalysisRepository defines the interface for transcript analysis persistence.
// Works with TranscriptAnalysisEncx (encrypted) type. Encryption/decryption is handled
// at the service layer using ProcessTranscriptAnalysisEncx/DecryptTranscriptAnalysisEncx.
type TranscriptAnalysisRepository interface {
	// Create creates or replaces the analysis for a transcription (encrypted).
	Create(ctx context.Context, analysis *ai.TranscriptAnalysisEncx) error

	// GetByID retrieves an analysis by ID (encrypted).
	GetByID(ctx context.Context, id uuid.UUID) (*ai.TranscriptAnalysisEncx, error)

	// GetByTranscriptionID retrieves an analysis by transcription ID (encrypted).
	GetByTranscriptionID(ctx context.Context, transcriptionID uuid.UUID) (*ai.TranscriptAnalysisEncx, error)

	// List retrieves analyses with pagination, optionally filtered by transcription ID (encrypted).
	List(ctx context.Context, transcriptionID *uuid.UUID, limit, offset int) ([]*ai.TranscriptAnalysisEncx, int, error)

	// Delete deletes an analysis by ID.
	Delete(ctx context.Context, id uuid.UUID) error
}
