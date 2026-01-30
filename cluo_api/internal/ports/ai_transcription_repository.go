package ports

import (
	"context"

	"github.com/google/uuid"
	"github.com/hengadev/cluo_api/internal/domain/ai"
)

// TranscriptionRepository defines the interface for transcription persistence.
// Works with TranscriptionEncx (encrypted) type. Encryption/decryption is handled
// at the service layer using ProcessTranscriptionEncx/DecryptTranscriptionEncx.
type TranscriptionRepository interface {
	// Create creates a new transcription (encrypted).
	Create(ctx context.Context, transcription *ai.TranscriptionEncx) error

	// GetByID retrieves a transcription by ID (encrypted).
	GetByID(ctx context.Context, id uuid.UUID) (*ai.TranscriptionEncx, error)

	// GetByMediaFileID retrieves a transcription by media file ID (encrypted).
	GetByMediaFileID(ctx context.Context, mediaFileID uuid.UUID) (*ai.TranscriptionEncx, error)

	// GetByJobID retrieves a transcription by job ID (encrypted).
	GetByJobID(ctx context.Context, jobID uuid.UUID) (*ai.TranscriptionEncx, error)

	// Delete deletes a transcription by ID.
	Delete(ctx context.Context, id uuid.UUID) error
}
