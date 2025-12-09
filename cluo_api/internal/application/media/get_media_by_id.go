package mediaService

import (
	"context"
	"fmt"

	"github.com/hengadev/cluo_api/internal/common/errs"
	domain "github.com/hengadev/cluo_api/internal/domain/media"
)

func (s *Service) GetMediaByID(ctx context.Context, r *domain.GetMediaByIDRequest) (*domain.MediaResponse, error) {
	// Get encrypted media
	mediaEncx, err := s.repo.GetMediaByID(ctx, r.ID)
	if err != nil {
		return nil, fmt.Errorf("failed to get media: %w", err)
	}

	// Decrypt
	media, err := domain.DecryptMediaFileEncx(ctx, s.crypto, mediaEncx)
	if err != nil {
		return nil, errs.NewNotDecryptedErr("media", err)
	}

	return media.ToResponse(), nil
}
