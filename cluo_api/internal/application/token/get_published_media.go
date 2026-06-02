package tokenService

import (
	"context"
	"fmt"

	domainMedia "github.com/hengadev/cluo_api/internal/domain/media"
)

func (s *Service) GetPublishedMediaByToken(ctx context.Context, rawToken string) ([]*domainMedia.MediaResponse, error) {
	caseID, err := s.ValidateToken(ctx, rawToken)
	if err != nil {
		return nil, err
	}

	// Fetch all media for the case (no type filter, large page size to get all)
	mediaEncxList, _, err := s.mediaRepo.ListMediaByCaseID(ctx, caseID, nil, 1, 10000)
	if err != nil {
		return nil, fmt.Errorf("failed to list media: %w", err)
	}

	responses := make([]*domainMedia.MediaResponse, 0, len(mediaEncxList))
	for _, mediaEncx := range mediaEncxList {
		// IsPublished is a plain (non-encrypted) field — check before decrypting
		if !mediaEncx.IsPublished {
			continue
		}

		media, err := domainMedia.DecryptMediaFileEncx(ctx, s.crypto, mediaEncx)
		if err != nil {
			// Skip media that can't be decrypted
			continue
		}

		// Exclude audio files from the client portal
		if media.Type == domainMedia.MediaTypeAudio {
			continue
		}

		responses = append(responses, media.ToResponse())
	}

	return responses, nil
}
