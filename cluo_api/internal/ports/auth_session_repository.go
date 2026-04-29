package ports

import (
	"context"
	"time"

	"github.com/google/uuid"
)

// AuthSessionRepository defines the interface for auth session operations
// This is separate from the middleware's session.SessionRepository which only has read operations
type AuthSessionRepository interface {
	// CreateSession creates a new session with access and refresh tokens
	CreateSession(ctx context.Context, sessionID uuid.UUID, accessTokenHash, refreshTokenHash, userIDHash string, sessionEncoded []byte, accessTTL, refreshTTL time.Duration) error

	// FindSessionByID retrieves session data by session ID
	FindSessionByID(ctx context.Context, sessionID uuid.UUID) ([]byte, error)

	// RefreshTokenPair refreshes the access and refresh tokens for a session
	RefreshTokenPair(ctx context.Context, oldRefreshTokenHash, newAccessTokenHash, newRefreshTokenHash string, sessionID uuid.UUID, updatedSessionData []byte, accessTTL, refreshTTL time.Duration) error

	// RemoveSessionByID removes a session by its ID
	RemoveSessionByID(ctx context.Context, sessionID uuid.UUID) error
}
