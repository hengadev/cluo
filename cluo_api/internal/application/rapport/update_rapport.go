package rapportService

import (
	"context"
	"fmt"
	"time"

	"github.com/hengadev/cluo_api/internal/common/errs"
	"github.com/hengadev/cluo_api/internal/domain/rapport"
)

func (s *Service) UpdateRapport(ctx context.Context, req *rapport.UpdateRapportRequest) (*rapport.RapportResponse, error) {
	// Get existing encrypted rapport to preserve ID and CreatedAt
	existing, err := s.repo.GetRapportByCaseID(ctx, req.CaseID)
	if err != nil {
		return nil, fmt.Errorf("failed to get existing rapport: %w", err)
	}

	// Build updated domain object preserving immutable fields
	r := &rapport.Rapport{
		ID:        existing.ID,
		CaseID:    req.CaseID,
		Content:   req.Content,
		CreatedAt: existing.CreatedAt,
		UpdatedAt: time.Now(),
	}

	// Re-encrypt with new content
	rEncx, err := rapport.ProcessRapportEncx(ctx, s.crypto, r)
	if err != nil {
		return nil, errs.NewNotEncryptedErr("rapport", err)
	}

	if err := s.repo.UpdateRapport(ctx, rEncx); err != nil {
		return nil, fmt.Errorf("failed to update rapport: %w", err)
	}

	return r.ToResponse(), nil
}
