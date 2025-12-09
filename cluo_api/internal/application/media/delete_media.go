package mediaService

import (
	"context"
	"fmt"

	"github.com/hengadev/cluo_api/internal/common/errs"
	domain "github.com/hengadev/cluo_api/internal/domain/media"
)

func (s *Service) DeleteMedia(ctx context.Context, r *domain.DeleteMediaRequest) error {
	// Get the media to retrieve the file URL
	mediaEncx, err := s.repo.GetMediaByID(ctx, r.ID)
	if err != nil {
		return fmt.Errorf("failed to get media: %w", err)
	}

	// Decrypt to get the URL
	media, err := domain.DecryptMediaFileEncx(ctx, s.crypto, mediaEncx)
	if err != nil {
		return errs.NewNotDecryptedErr("media", err)
	}

	// Delete from database first
	if err := s.repo.DeleteMedia(ctx, r.ID); err != nil {
		return fmt.Errorf("failed to delete media: %w", err)
	}

	// Delete file from storage (best effort - don't fail if storage deletion fails)
	if err := s.storage.DeleteFile(ctx, media.URL); err != nil {
		// Log the error but don't fail the operation
		// In production, you might want to queue this for retry
		fmt.Printf("Warning: failed to delete file from storage: %v\n", err)
	}

	return nil
}
