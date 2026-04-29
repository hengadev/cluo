package auth

import (
	"context"
	"encoding/json"
	"time"

	"github.com/hengadev/cluo_api/internal/common/auth/session"
	"github.com/hengadev/cluo_api/internal/common/contracts/identity"
	"github.com/hengadev/cluo_api/internal/common/errs"
	"github.com/hengadev/cluo_api/internal/domain/user"
	"github.com/hengadev/encx"
	"github.com/google/uuid"
)

// SignIn authenticates a user with email and password
func (s *AuthService) SignIn(ctx context.Context, req *user.SignInRequest) (*user.CreateSessionResponse, error) {
	// Validate request
	if err := req.Valid(); err != nil {
		return nil, err
	}

	// Hash email for lookup
	emailBytes, err := encx.SerializeValue(req.Email)
	if err != nil {
		return nil, err
	}
	emailHash := s.crypto.HashBasic(ctx, emailBytes)

	// Get user by email hash
	userEncx, err := s.userRepo.GetUserByEmailHash(ctx, emailHash)
	if err != nil {
		return nil, err
	}

	// Compare password against stored hash (password is hash_secure — not decryptable)
	passwordMatch, err := s.crypto.CompareSecureHashAndValue(ctx, req.Password, userEncx.PasswordHashSecure)
	if err != nil {
		return nil, err
	}
	if !passwordMatch {
		return nil, errs.ErrUnauthorized
	}

	// Decrypt user to read role
	u, err := user.DecryptUserEncx(ctx, s.crypto, userEncx)
	if err != nil {
		return nil, err
	}

	// Parse role
	role, err := identity.ParseRole(u.Role)
	if err != nil {
		return nil, err
	}

	// Create session
	return s.createSession(ctx, u.ID, role, session.SessionActive)
}

// createSession creates a new session for a user
func (s *AuthService) createSession(ctx context.Context, userID uuid.UUID, role identity.Role, state session.SessionState) (*user.CreateSessionResponse, error) {
	now := time.Now()
	accessTokenExpiry := now.Add(session.ActiveSessionDuration)
	refreshTokenExpiry := now.Add(session.ActiveSessionDuration)

	sessionID := uuid.New()
	accessToken, err := session.GenerateToken()
	if err != nil {
		return nil, err
	}
	refreshToken, err := session.GenerateToken()
	if err != nil {
		return nil, err
	}

	// Serialize and hash tokens
	accessTokenBytes, err := encx.SerializeValue(accessToken)
	if err != nil {
		return nil, err
	}
	accessTokenHash := s.crypto.HashBasic(ctx, accessTokenBytes)

	refreshTokenBytes, err := encx.SerializeValue(refreshToken)
	if err != nil {
		return nil, err
	}
	refreshTokenHash := s.crypto.HashBasic(ctx, refreshTokenBytes)

	// Create session struct
	sess := &session.Session{
		ID:           sessionID,
		UserID:       userID,
		Role:         role,
		State:        state,
		CreatedAt:    now,
		ExpiresAt:    accessTokenExpiry,
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
	}

	// Encrypt session (Session has json:"-" on all fields — must encrypt to SessionEncx first)
	sessEncx, err := session.ProcessSessionEncx(ctx, s.crypto, sess)
	if err != nil {
		return nil, err
	}
	sessionEncoded, err := json.Marshal(sessEncx)
	if err != nil {
		return nil, err
	}

	// Hash user ID for Redis indexing
	userIDBytes, err := encx.SerializeValue(userID)
	if err != nil {
		return nil, err
	}
	userIDHash := s.crypto.HashBasic(ctx, userIDBytes)

	// Store session
	err = s.sessionRepo.CreateSession(ctx, sessionID, accessTokenHash, refreshTokenHash, userIDHash, sessionEncoded, session.ActiveSessionDuration, session.ActiveSessionDuration)
	if err != nil {
		return nil, err
	}

	return &user.CreateSessionResponse{
		AccessToken:        accessToken,
		RefreshToken:       refreshToken,
		AccessTokenExpiry:  accessTokenExpiry,
		RefreshTokenExpiry: refreshTokenExpiry,
	}, nil
}
