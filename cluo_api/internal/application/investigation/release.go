package investigationService

import (
	"context"
	"fmt"
	"time"

	"github.com/google/uuid"
	"github.com/hengadev/cluo_api/internal/common/ctxutil"
	"github.com/hengadev/cluo_api/internal/common/errs"
	"github.com/hengadev/cluo_api/internal/domain/investigation"
)

func (s *CaseService) Release(ctx context.Context, caseID uuid.UUID) (*investigation.ReleaseResponse, error) {
	logger, err := ctxutil.GetLoggerFromContext(ctx)
	if err != nil {
		return nil, fmt.Errorf("get logger for release: %w", err)
	}

	caseEncx, err := s.repo.GetCaseByID(ctx, caseID)
	if err != nil {
		return nil, fmt.Errorf("get case for release: %w", err)
	}

	c, err := investigation.DecryptInvestigationEncx(ctx, s.crypto, caseEncx)
	if err != nil {
		return nil, errs.NewNotDecryptedErr("case for release", err)
	}

	if c.Status == investigation.StatusInProgress {
		return nil, errs.NewConflictErr(fmt.Errorf("case status is %q; must be %q or %q to release", c.Status, investigation.StatusReady, investigation.StatusReleased))
	}

	if c.Status == investigation.StatusReady {
		c.Status = investigation.StatusReleased
		c.UpdatedAt = time.Now()

		updatedEncx, err := investigation.ProcessInvestigationEncx(ctx, s.crypto, c)
		if err != nil {
			return nil, errs.NewNotEncryptedErr("case for release", err)
		}

		if err := s.repo.UpdateCase(ctx, updatedEncx); err != nil {
			return nil, fmt.Errorf("update case status for release: %w", err)
		}
	}

	tokenResp, err := s.tokenService.CreateToken(ctx, caseID)
	if err != nil {
		return nil, fmt.Errorf("create token for release: %w", err)
	}

	portalURL := "/token/" + tokenResp.RawToken
	logger.InfoContext(ctx, "portal URL", "url", portalURL)

	return &investigation.ReleaseResponse{
		CaseID:    caseID.String(),
		TokenID:   tokenResp.ID,
		RawToken:  tokenResp.RawToken,
		PortalURL: portalURL,
		ExpiresAt: tokenResp.ExpiresAt,
	}, nil
}
