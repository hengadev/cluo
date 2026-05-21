package tokenService

import (
	"context"
	"fmt"

	"github.com/google/uuid"
	"github.com/hengadev/cluo_api/internal/domain/token"
)

func (s *Service) ListTokensByCaseID(ctx context.Context, caseID uuid.UUID) ([]*token.TokenResponse, error) {
	tokens, err := s.repo.ListTokensByCaseID(ctx, caseID)
	if err != nil {
		return nil, fmt.Errorf("failed to list tokens: %w", err)
	}

	responses := make([]*token.TokenResponse, 0, len(tokens))
	for _, t := range tokens {
		responses = append(responses, t.ToResponse())
	}
	return responses, nil
}
