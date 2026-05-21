package tokenService

import (
	"context"
	"fmt"

	"github.com/google/uuid"
	"github.com/hengadev/cluo_api/internal/common/errs"
)

func (s *Service) RevokeToken(ctx context.Context, tokenID uuid.UUID) error {
	if err := s.repo.RevokeToken(ctx, tokenID); err != nil {
		return errs.NewNotUpdatedErr(fmt.Errorf("revoke token: %w", err), "token")
	}
	return nil
}
