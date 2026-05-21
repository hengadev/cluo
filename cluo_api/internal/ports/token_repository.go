package ports

import (
	"context"

	"github.com/google/uuid"
	"github.com/hengadev/cluo_api/internal/domain/token"
)

type TokenRepository interface {
	CreateToken(ctx context.Context, t *token.Token) error
	GetTokenByHash(ctx context.Context, tokenHash string) (*token.Token, error)
	ListTokensByCaseID(ctx context.Context, caseID uuid.UUID) ([]*token.Token, error)
	RevokeToken(ctx context.Context, tokenID uuid.UUID) error
}
