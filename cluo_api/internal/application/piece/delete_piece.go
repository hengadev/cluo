package pieceService

import (
	"context"
	"fmt"

	"github.com/google/uuid"
	"github.com/hengadev/cluo_api/internal/common/errs"
	"github.com/hengadev/cluo_api/internal/domain/piece"
)

func (s *Service) DeletePiece(ctx context.Context, caseID uuid.UUID, pieceID uuid.UUID) error {
	// Retrieve piece to get storage key
	pieceEncx, err := s.repo.GetPieceByID(ctx, pieceID)
	if err != nil {
		return fmt.Errorf("failed to get piece: %w", err)
	}

	// Decrypt to recover the storage key (StorageKey is plain, but we still go through DecryptPieceEncx)
	p, err := piece.DecryptPieceEncx(ctx, s.crypto, pieceEncx)
	if err != nil {
		return errs.NewNotDecryptedErr("piece", err)
	}

	// Delete from DB first
	if err := s.repo.DeletePiece(ctx, pieceID); err != nil {
		return fmt.Errorf("failed to delete piece: %w", err)
	}

	// Delete from storage (best effort)
	if err := s.storage.DeleteFile(ctx, p.StorageKey); err != nil {
		fmt.Printf("Warning: failed to delete piece file from storage: %v\n", err)
	}

	return nil
}
