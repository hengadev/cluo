package rapportService

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/google/uuid"
	"github.com/hengadev/cluo_api/internal/common/errs"
	"github.com/hengadev/cluo_api/internal/domain/rapport"
)

func (s *Service) CreateRapport(ctx context.Context, req *rapport.CreateRapportRequest) (*rapport.RapportResponse, error) {
	// Validate caseID
	caseID, err := uuid.Parse(req.CaseID)
	if err != nil {
		return nil, errs.NewInvalidValueErr("invalid case ID format")
	}

	// Check if case exists
	exists, err := s.caseRepo.ExistsCase(ctx, caseID)
	if err != nil {
		return nil, fmt.Errorf("failed to check case existence: %w", err)
	}
	if !exists {
		return nil, errs.NewNotFoundErr(fmt.Errorf("case with ID %s not found", req.CaseID), "case")
	}

	// Build domain object
	now := time.Now()
	r := &rapport.Rapport{
		ID:        uuid.New(),
		CaseID:    caseID,
		Content:   req.Content,
		CreatedAt: now,
		UpdatedAt: now,
	}

	// Encrypt
	rEncx, err := rapport.ProcessRapportEncx(ctx, s.crypto, r)
	if err != nil {
		return nil, errs.NewNotEncryptedErr("rapport", err)
	}

	// Persist — repo returns ErrUniqueViolation when a rapport already exists for this case
	if err := s.repo.CreateRapport(ctx, rEncx); err != nil {
		if errors.Is(err, errs.ErrUniqueViolation) {
			return nil, errs.NewAlreadyExistsError(err, "rapport")
		}
		return nil, fmt.Errorf("failed to create rapport: %w", err)
	}

	return r.ToResponse(), nil
}
