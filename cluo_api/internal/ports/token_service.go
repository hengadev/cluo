package ports

import (
	"context"

	"github.com/google/uuid"
	"github.com/hengadev/cluo_api/internal/domain/investigation"
	domainMedia "github.com/hengadev/cluo_api/internal/domain/media"
	"github.com/hengadev/cluo_api/internal/domain/token"
)

type TokenService interface {
	CreateToken(ctx context.Context, caseID uuid.UUID) (*token.CreateTokenResponse, error)
	ValidateToken(ctx context.Context, rawToken string) (uuid.UUID, error) // returns caseID or error
	ListTokensByCaseID(ctx context.Context, caseID uuid.UUID) ([]*token.TokenResponse, error)
	RevokeToken(ctx context.Context, tokenID uuid.UUID) error
	GetCaseSummaryByToken(ctx context.Context, rawToken string) (*investigation.PortalCaseResponse, error)
	GetPublishedMediaByToken(ctx context.Context, rawToken string) ([]*domainMedia.MediaResponse, error)
}
