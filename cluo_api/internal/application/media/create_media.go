package mediaService

import (
	"context"
	"fmt"

	"github.com/google/uuid"
	"github.com/hengadev/cluo_api/internal/common/errs"
	domain "github.com/hengadev/cluo_api/internal/domain/media"
)

// createMedia is a private method used internally by UploadMedia
// It creates the media record in the database after the file has been uploaded to S3
func (s *Service) createMedia(ctx context.Context, r *domain.CreateMediaRequest) (*domain.MediaResponse, error) {
	// Validate request
	if err := r.Valid(ctx); err != nil {
		return nil, errs.NewInvalidValueErr(err.Error())
	}

	// Validate case exists
	caseID, _ := uuid.Parse(r.CaseID)
	exists, err := s.caseRepo.ExistsCase(ctx, caseID)
	if err != nil {
		return nil, fmt.Errorf("failed to check case existence: %w", err)
	}
	if !exists {
		return nil, errs.NewRepositoryNotFoundErr(
			fmt.Errorf("case with ID %s not found", r.CaseID),
			"case",
		)
	}

	// Create domain object
	media := domain.NewMedia(r)

	// Encrypt
	mediaEncx, err := domain.ProcessMediaFileEncx(ctx, s.crypto, media)
	if err != nil {
		return nil, errs.NewNotEncryptedErr("media", err)
	}

	// Persist
	if err := s.repo.CreateMedia(ctx, mediaEncx); err != nil {
		return nil, fmt.Errorf("failed to create media: %w", err)
	}

	return media.ToResponse(), nil
}
