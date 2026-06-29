package mediaService

import (
	"context"
	"fmt"

	"github.com/hengadev/cluo_api/internal/common/errs"
	domain "github.com/hengadev/cluo_api/internal/domain/media"
)

func (s *Service) UpdateMedia(ctx context.Context, r *domain.UpdateMediaRequest) (*domain.MediaResponse, error) {
	// Validate request
	if err := r.Valid(ctx); err != nil {
		return nil, errs.NewInvalidValueErr(err.Error())
	}

	// Get existing media
	mediaEncx, err := s.repo.GetMediaByID(ctx, r.ID)
	if err != nil {
		return nil, fmt.Errorf("failed to get media: %w", err)
	}

	// Decrypt
	media, err := domain.DecryptMediaFileEncx(ctx, s.crypto, mediaEncx)
	if err != nil {
		return nil, errs.NewNotDecryptedErr("media", err)
	}

	// Apply updates
	if r.Caption != nil {
		media.Caption = *r.Caption
	}
	if r.IsPublished != nil {
		media.IsPublished = *r.IsPublished
	}
	if r.Purpose != nil {
		media.Purpose = domain.RecordingPurpose(*r.Purpose)
	}

	// Re-encrypt
	updatedEncx, err := domain.ProcessMediaFileEncx(ctx, s.crypto, media)
	if err != nil {
		return nil, errs.NewNotEncryptedErr("media", err)
	}

	// Persist
	if err := s.repo.UpdateMedia(ctx, updatedEncx); err != nil {
		return nil, fmt.Errorf("failed to update media: %w", err)
	}

	return media.ToResponse(), nil
}
