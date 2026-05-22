package auth

import (
	"context"

	"github.com/hengadev/cluo_api/internal/domain/user"
	"github.com/google/uuid"
)

// GetCurrentUser returns the current authenticated user
func (s *AuthService) GetCurrentUser(ctx context.Context, userID uuid.UUID) (*user.CurrentUserResponse, error) {
	userEncx, err := s.userRepo.GetUserByID(ctx, userID)
	if err != nil {
		return nil, err
	}

	u, err := user.DecryptUserEncx(ctx, s.crypto, userEncx)
	if err != nil {
		return nil, err
	}

	return &user.CurrentUserResponse{
		ID:    u.ID.String(),
		Email: u.Email,
		Role:  u.Role,
	}, nil
}
