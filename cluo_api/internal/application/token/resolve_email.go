package tokenService

import (
	"context"
	"fmt"

	"github.com/google/uuid"
	"github.com/hengadev/cluo_api/internal/domain/client"
	"github.com/hengadev/cluo_api/internal/domain/investigation"
)

// resolveClientEmail looks up the first contact email for the client associated
// with the given case. Returns an empty string if no contact is found.
func (s *Service) resolveClientEmail(ctx context.Context, caseID uuid.UUID) (string, error) {
	// Get the case to find the client ID
	caseEncx, err := s.caseRepo.GetCaseByID(ctx, caseID)
	if err != nil {
		return "", fmt.Errorf("get case for email resolution: %w", err)
	}

	c, err := investigation.DecryptInvestigationEncx(ctx, s.crypto, caseEncx)
	if err != nil {
		return "", fmt.Errorf("decrypt case for email resolution: %w", err)
	}

	if c.ClientID == uuid.Nil {
		return "", nil
	}

	// Get contacts for this client
	contactEncxs, err := s.clientRepo.GetAllContactsByClientID(ctx, c.ClientID)
	if err != nil {
		return "", fmt.Errorf("get contacts for email resolution: %w", err)
	}

	for _, ce := range contactEncxs {
		contact, err := client.DecryptContactEncx(ctx, s.crypto, ce)
		if err != nil {
			continue
		}
		if contact.Email != "" {
			return contact.Email, nil
		}
	}

	return "", nil
}
