package tokenService

import (
	"context"
	"fmt"

	"github.com/google/uuid"
	"github.com/hengadev/cluo_api/internal/common/errs"
	domainMedia "github.com/hengadev/cluo_api/internal/domain/media"
)

func (s *Service) GetPublishedMediaByIDAndToken(ctx context.Context, rawToken string, mediaID uuid.UUID) (*domainMedia.MediaResponse, error) {
	caseID, err := s.ValidateToken(ctx, rawToken)
	if err != nil {
		return nil, err
	}

	mediaEncx, err := s.mediaRepo.GetMediaByID(ctx, mediaID)
	if err != nil {
		return nil, fmt.Errorf("failed to get media: %w", err)
	}

	if mediaEncx.CaseID != caseID {
		return nil, errs.NewNotFoundErr(fmt.Errorf("media %s not in case %s", mediaID, caseID), "media")
	}

	if !mediaEncx.IsPublished {
		return nil, errs.NewNotFoundErr(fmt.Errorf("media %s is not published", mediaID), "media")
	}

	media, err := domainMedia.DecryptMediaFileEncx(ctx, s.crypto, mediaEncx)
	if err != nil {
		return nil, fmt.Errorf("failed to decrypt media: %w", err)
	}

	return media.ToResponse(), nil
}
