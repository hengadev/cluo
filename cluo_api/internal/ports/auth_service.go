package ports

import (
	"context"

	"github.com/hengadev/cluo_api/internal/common/auth/session"
	"github.com/hengadev/cluo_api/internal/domain/user"
	"github.com/google/uuid"
)

// AuthService defines the interface for authentication operations
type AuthService interface {
	// SignIn authenticates a user with email and password
	SignIn(ctx context.Context, req *user.SignInRequest) (*user.CreateSessionResponse, error)

	// Register creates a new user account
	Register(ctx context.Context, req *user.RegisterRequest) (*user.CreateSessionResponse, error)

	// SignOut logs out a user by removing their session
	SignOut(ctx context.Context, sessionInfo *session.SessionInfo) error

	// RefreshSession refreshes an access token using a refresh token
	RefreshSession(ctx context.Context, sessionID uuid.UUID) (*user.RefreshSessionResponse, error)

	// GetCurrentUser returns the current authenticated user
	GetCurrentUser(ctx context.Context, userID uuid.UUID) (*user.CurrentUserResponse, error)

	// UpdateCurrentUserName updates the display name of the authenticated user
	UpdateCurrentUserName(ctx context.Context, userID uuid.UUID, name string) error
}
