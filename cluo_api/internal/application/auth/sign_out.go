package auth

import (
	"context"

	"github.com/hengadev/cluo_api/internal/common/auth/session"
)

// SignOut logs out a user by removing their session
func (s *AuthService) SignOut(ctx context.Context, sessionInfo *session.SessionInfo) error {
	// Remove session from Redis
	err := s.sessionRepo.RemoveSessionByID(ctx, sessionInfo.ID)
	if err != nil {
		return err
	}

	return nil
}
