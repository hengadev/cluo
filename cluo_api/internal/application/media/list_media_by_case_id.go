package mediaService

import (
	"context"
	"fmt"

	"github.com/google/uuid"
	"github.com/hengadev/cluo_api/internal/common/errs"
	domain "github.com/hengadev/cluo_api/internal/domain/media"
)

func (s *Service) ListMediaByCaseID(ctx context.Context, r *domain.ListMediaByCaseIDRequest) (*domain.ListMediaResponse, error) {
	// Validate request
	if err := r.Valid(ctx); err != nil {
		return nil, errs.NewInvalidValueErr(err.Error())
	}

	caseID, _ := uuid.Parse(r.CaseID)

	// Parse type filter if provided
	var typeFilter *domain.MediaType
	if r.Type != nil && *r.Type != "" {
		mt := domain.MediaType(*r.Type)
		typeFilter = &mt
	}

	// Get from repository
	mediaEncxList, total, err := s.repo.ListMediaByCaseID(ctx, caseID, typeFilter, r.Page, r.PageSize)
	if err != nil {
		return nil, fmt.Errorf("failed to list media: %w", err)
	}

	// Decrypt and filter by type at service layer
	mediaResponses := make([]*domain.MediaResponse, 0, len(mediaEncxList))
	for _, mediaEncx := range mediaEncxList {
		media, err := domain.DecryptMediaFileEncx(ctx, s.crypto, mediaEncx)
		if err != nil {
			// Skip media that can't be decrypted
			continue
		}

		// Apply type filter if provided (since type is encrypted)
		if typeFilter != nil && media.Type != *typeFilter {
			continue
		}

		mediaResponses = append(mediaResponses, media.ToResponse())
	}

	// Adjust total if we filtered at service layer
	if typeFilter != nil {
		total = len(mediaResponses)
	}

	paginationInfo := domain.NewPaginationInfo(r.Page, r.PageSize, total)

	return &domain.ListMediaResponse{
		Media:      mediaResponses,
		Pagination: paginationInfo,
	}, nil
}
