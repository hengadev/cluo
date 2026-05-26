package tokenService

import (
	"context"
	"fmt"
	"time"

	"github.com/google/uuid"
	"github.com/hengadev/cluo_api/internal/common/errs"
	"github.com/hengadev/cluo_api/internal/domain/token"
)

func (s *Service) CreateToken(ctx context.Context, caseID uuid.UUID) (*token.CreateTokenResponse, error) {
	exists, err := s.caseRepo.ExistsCase(ctx, caseID)
	if err != nil {
		return nil, fmt.Errorf("failed to verify case existence: %w", err)
	}
	if !exists {
		return nil, errs.NewNotFoundErr(fmt.Errorf("case %s not found", caseID), "case")
	}

	rawToken, tokenHash, err := token.GenerateRawToken()
	if err != nil {
		return nil, fmt.Errorf("failed to generate token: %w", err)
	}

	now := time.Now()
	t := &token.Token{
		ID:        uuid.New(),
		CaseID:    caseID,
		TokenHash: tokenHash,
		ExpiresAt: now.Add(token.TokenExpiryDays * 24 * time.Hour),
		RevokedAt: nil,
		CreatedAt: now,
	}

	if err := s.repo.CreateToken(ctx, t); err != nil {
		return nil, errs.NewNotCreatedErr(err, "token")
	}

	// Dispatch portal access email asynchronously.
	portalPublicURL := s.portalPublicURL
	go func() {
		bgCtx := context.Background()
		email, err := s.resolveClientEmail(bgCtx, caseID)
		if err != nil {
			s.logger.ErrorContext(bgCtx, "Failed to resolve client email for token notification",
				"error", err,
				"case_id", caseID,
			)
			return
		}
		if email == "" {
			s.logger.WarnContext(bgCtx, "No client email found; skipping token notification",
				"case_id", caseID,
			)
			return
		}

		subject := "Accès au portail client"
		portalURL := portalPublicURL + "/token/" + rawToken
		bodyHTML := fmt.Sprintf(`
			<html><body>
			<p>Bonjour,</p>
			<p>Un accès au portail client a été créé pour votre dossier.</p>
			<p><a href="%s">Accéder au portail</a></p>
			<p>Ce lien expirera le %s.</p>
			</body></html>
		`, portalURL, t.ExpiresAt.Format("02/01/2006"))

		if err := s.emailService.Send(bgCtx, email, subject, bodyHTML); err != nil {
			s.logger.ErrorContext(bgCtx, "Failed to send portal token email",
				"error", err,
				"to", email,
				"case_id", caseID,
			)
		}
	}()

	return &token.CreateTokenResponse{
		ID:        t.ID.String(),
		CaseID:    t.CaseID.String(),
		RawToken:  rawToken,
		ExpiresAt: t.ExpiresAt,
		CreatedAt: t.CreatedAt,
	}, nil
}
