package auth

import (
	"github.com/hengadev/cluo_api/internal/ports"
	"github.com/hengadev/encx"
)

// AuthService implements the authentication operations
type AuthService struct {
	userRepo    ports.UserRepository
	sessionRepo ports.AuthSessionRepository
	crypto      encx.CryptoService
}

// New creates a new AuthService
func New(userRepo ports.UserRepository, sessionRepo ports.AuthSessionRepository, crypto encx.CryptoService) ports.AuthService {
	return &AuthService{
		userRepo:    userRepo,
		sessionRepo: sessionRepo,
		crypto:      crypto,
	}
}
