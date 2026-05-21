package pieceService

import (
	"context"
	"fmt"
	"mime/multipart"
	"time"

	"github.com/google/uuid"
	"github.com/hengadev/cluo_api/internal/common/errs"
	"github.com/hengadev/cluo_api/internal/domain/piece"
)

func (s *Service) UploadPiece(ctx context.Context, caseID uuid.UUID, file multipart.File, header *multipart.FileHeader, notes string) (*piece.PieceResponse, error) {
	// Validate case exists
	exists, err := s.caseRepo.ExistsCase(ctx, caseID)
	if err != nil {
		return nil, fmt.Errorf("failed to check case existence: %w", err)
	}
	if !exists {
		return nil, errs.NewRepositoryNotFoundErr(
			fmt.Errorf("case with ID %s not found", caseID),
			"case",
		)
	}

	// Get MIME type from header
	mimeType := header.Header.Get("Content-Type")
	if mimeType == "" {
		mimeType = "application/octet-stream"
	}

	// Upload file to storage; the returned URL becomes the storage key
	storageKey, err := s.storage.UploadFile(ctx, file, header.Filename, mimeType, header.Size)
	if err != nil {
		return nil, fmt.Errorf("failed to upload piece to storage: %w", err)
	}

	// Build domain Piece
	now := time.Now()
	p := &piece.Piece{
		ID:         uuid.New(),
		CaseID:     caseID,
		Filename:   header.Filename,
		StorageKey: storageKey,
		MimeType:   mimeType,
		SizeBytes:  header.Size,
		Notes:      notes,
		CreatedAt:  now,
		UpdatedAt:  now,
	}

	// Encrypt sensitive fields
	pieceEncx, err := piece.ProcessPieceEncx(ctx, s.crypto, p)
	if err != nil {
		// Best-effort cleanup of the uploaded file
		if delErr := s.storage.DeleteFile(ctx, storageKey); delErr != nil {
			return nil, fmt.Errorf("failed to encrypt piece (and failed to cleanup uploaded file): %w", err)
		}
		return nil, errs.NewNotEncryptedErr("piece", err)
	}

	// Persist metadata
	if err := s.repo.CreatePiece(ctx, pieceEncx); err != nil {
		// Best-effort cleanup
		if delErr := s.storage.DeleteFile(ctx, storageKey); delErr != nil {
			return nil, fmt.Errorf("failed to create piece record (and failed to cleanup uploaded file): %w", err)
		}
		return nil, fmt.Errorf("failed to create piece record: %w", err)
	}

	return p.ToResponse(), nil
}
