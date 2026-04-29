package auth

import (
	"context"
	"time"

	"github.com/hengadev/cluo_api/internal/common/auth/session"
	"github.com/hengadev/cluo_api/internal/common/contracts/identity"
	"github.com/hengadev/cluo_api/internal/common/errs"
	"github.com/hengadev/cluo_api/internal/domain/user"
	"github.com/hengadev/encx"
	"github.com/google/uuid"
)

// Register creates a new user account
func (s *AuthService) Register(ctx context.Context, req *user.RegisterRequest) (*user.CreateSessionResponse, error) {
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

	// Check if user already exists
	exists, err := s.userRepo.ExistsByEmailHash(ctx, emailHash)
	if err != nil {
		return nil, err
	}

	if exists {
		return nil, errs.ErrAlreadyExists
	}

	// Create new user with plaintext password — ProcessUserEncx will hash it via hash_secure tag
	userID := uuid.New()
	newUser := &user.User{
		ID:        userID,
		Email:     req.Email,
		Password:  req.Password,
		Role:      identity.Client.String(),
		CreatedAt: time.Now(),
	}

	// Encrypt user
	userEncx, err := user.ProcessUserEncx(ctx, s.crypto, newUser)
	if err != nil {
		return nil, err
	}

	// Save user to database
	err = s.userRepo.CreateUser(ctx, userEncx)
	if err != nil {
		return nil, err
	}

	// Create session
	return s.createSession(ctx, userID, identity.Client, session.SessionActive)
}
