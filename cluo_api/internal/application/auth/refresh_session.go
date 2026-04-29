package auth

import (
	"context"
	"encoding/json"
	"time"

	"github.com/hengadev/cluo_api/internal/common/auth/session"
	"github.com/hengadev/cluo_api/internal/common/errs"
	"github.com/hengadev/cluo_api/internal/domain/user"
	"github.com/google/uuid"
)

// RefreshSession refreshes an access token using a refresh token
func (s *AuthService) RefreshSession(ctx context.Context, sessionID uuid.UUID) (*user.RefreshSessionResponse, error) {
	// Find session by ID
	sessionData, err := s.sessionRepo.FindSessionByID(ctx, sessionID)
	if err != nil {
		return nil, err
	}

	// Decode session
	sessionEncx, err := session.DecodeSession(sessionData)
	if err != nil {
		return nil, err
	}

	// Decrypt session
	sess, err := session.DecryptSessionEncx(ctx, s.crypto, sessionEncx)
	if err != nil {
		return nil, err
	}

	// Check if session is expired
	if time.Now().After(sess.ExpiresAt) {
		return nil, errs.ErrExpiredToken
	}

	// Capture the old refresh token hash before generating new tokens
	// (sess.RefreshToken is empty after decryption — it's hash_secure, one-way)
	oldRefreshTokenHash := sessionEncx.RefreshTokenHash

	// Generate new tokens
	now := time.Now()
	accessTokenExpiry := now.Add(session.ActiveSessionDuration)
	refreshTokenExpiry := now.Add(session.ActiveSessionDuration)

	accessToken, err := session.GenerateToken()
	if err != nil {
		return nil, err
	}
	refreshToken, err := session.GenerateToken()
	if err != nil {
		return nil, err
	}

	// Update session with new plaintext tokens — ProcessSessionEncx will hash them
	sess.AccessToken = accessToken
	sess.RefreshToken = refreshToken
	sess.ExpiresAt = refreshTokenExpiry

	// Encrypt updated session (Session has json:"-" on all fields — must go through SessionEncx)
	updatedEncx, err := session.ProcessSessionEncx(ctx, s.crypto, sess)
	if err != nil {
		return nil, err
	}
	updatedSessionData, err := json.Marshal(updatedEncx)
	if err != nil {
		return nil, err
	}

	// Update in Redis — get new token hashes from the encrypted struct
	err = s.sessionRepo.RefreshTokenPair(ctx, oldRefreshTokenHash, updatedEncx.AccessTokenHash, updatedEncx.RefreshTokenHash, sess.ID, updatedSessionData, session.ActiveSessionDuration, session.ActiveSessionDuration)
	if err != nil {
		return nil, err
	}

	return &user.RefreshSessionResponse{
		AccessToken:        accessToken,
		RefreshToken:       refreshToken,
		AccessTokenExpiry:  accessTokenExpiry,
		RefreshTokenExpiry: refreshTokenExpiry,
	}, nil
}
