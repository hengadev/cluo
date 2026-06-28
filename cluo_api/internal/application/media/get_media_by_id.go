package mediaService

import (
	"context"
	"fmt"

	"github.com/hengadev/cluo_api/internal/common/errs"
	domain "github.com/hengadev/cluo_api/internal/domain/media"
)

func (s *Service) GetMediaByID(ctx context.Context, r *domain.GetMediaByIDRequest) (*domain.MediaResponse, error) {
	mediaEncx, err := s.repo.GetMediaByID(ctx, r.ID)
	if err != nil {
		return nil, fmt.Errorf("failed to get media: %w", err)
	}

	media, err := domain.DecryptMediaFileEncx(ctx, s.crypto, mediaEncx)
	if err != nil {
		return nil, errs.NewNotDecryptedErr("media", err)
	}

	response := media.ToResponse()

	// Replace the stored internal URL with a 1-hour presigned URL so the
	// browser can stream or download the audio directly from MinIO without
	// needing API credentials.
	if media.URL != "" {
		if signedURL, err := s.storage.GetFileURL(ctx, media.URL); err == nil {
			response.URL = signedURL
		}
	}

	return response, nil
}
