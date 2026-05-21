package rapportService

import (
	"context"
	"fmt"

	"github.com/google/uuid"
	"github.com/hengadev/cluo_api/internal/common/errs"
	"github.com/hengadev/cluo_api/internal/domain/rapport"
)

func (s *Service) GetRapportByCaseID(ctx context.Context, caseID uuid.UUID) (*rapport.RapportResponse, error) {
	rEncx, err := s.repo.GetRapportByCaseID(ctx, caseID)
	if err != nil {
		return nil, fmt.Errorf("failed to get rapport: %w", err)
	}

	r, err := rapport.DecryptRapportEncx(ctx, s.crypto, rEncx)
	if err != nil {
		return nil, errs.NewNotDecryptedErr("rapport", err)
	}

	return r.ToResponse(), nil
}
