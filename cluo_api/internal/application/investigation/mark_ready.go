package investigationService

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/google/uuid"
	"github.com/hengadev/cluo_api/internal/common/errs"
	"github.com/hengadev/cluo_api/internal/domain/investigation"
)

func (s *CaseService) MarkReady(ctx context.Context, caseID uuid.UUID) (*investigation.CaseResponse, error) {
	caseEncx, err := s.repo.GetCaseByID(ctx, caseID)
	if err != nil {
		return nil, fmt.Errorf("get case for mark-ready: %w", err)
	}

	c, err := investigation.DecryptInvestigationEncx(ctx, s.crypto, caseEncx)
	if err != nil {
		return nil, errs.NewNotDecryptedErr("case for mark-ready", err)
	}

	if c.Status != investigation.StatusInProgress {
		return nil, errs.NewConflictErr(fmt.Errorf("case status is %q; must be %q to mark ready", c.Status, investigation.StatusInProgress))
	}

	_, err = s.rapportRepo.GetRapportByCaseID(ctx, caseID)
	if err != nil {
		if errors.Is(err, errs.ErrRepositoryNotFound) {
			return nil, errs.NewUnprocessableEntityErr("no rapport exists for this case")
		}
		return nil, fmt.Errorf("check rapport for mark-ready: %w", err)
	}

	c.Status = investigation.StatusReady
	c.UpdatedAt = time.Now()

	updatedEncx, err := investigation.ProcessInvestigationEncx(ctx, s.crypto, c)
	if err != nil {
		return nil, errs.NewNotEncryptedErr("case for mark-ready", err)
	}

	if err := s.repo.UpdateCase(ctx, updatedEncx); err != nil {
		return nil, fmt.Errorf("update case for mark-ready: %w", err)
	}

	return c.ToResponse(), nil
}
