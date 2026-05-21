package token

import (
	"testing"
	"time"

	"github.com/google/uuid"
)

func TestIsValid_RevokedToken(t *testing.T) {
	revokedAt := time.Now().Add(-time.Hour)
	tok := &Token{
		ID:        uuid.New(),
		CaseID:    uuid.New(),
		TokenHash: "abc",
		ExpiresAt: time.Now().Add(24 * time.Hour),
		RevokedAt: &revokedAt,
		CreatedAt: time.Now().Add(-time.Hour * 2),
	}
	if tok.IsValid() {
		t.Error("expected IsValid() to return false for a revoked token")
	}
}

func TestIsValid_ExpiredToken(t *testing.T) {
	tok := &Token{
		ID:        uuid.New(),
		CaseID:    uuid.New(),
		TokenHash: "abc",
		ExpiresAt: time.Now().Add(-time.Hour),
		RevokedAt: nil,
		CreatedAt: time.Now().Add(-time.Hour * 2),
	}
	if tok.IsValid() {
		t.Error("expected IsValid() to return false for an expired token")
	}
}

func TestIsValid_FreshToken(t *testing.T) {
	tok := &Token{
		ID:        uuid.New(),
		CaseID:    uuid.New(),
		TokenHash: "abc",
		ExpiresAt: time.Now().Add(TokenExpiryDays * 24 * time.Hour),
		RevokedAt: nil,
		CreatedAt: time.Now(),
	}
	if !tok.IsValid() {
		t.Error("expected IsValid() to return true for a fresh, unrevoked token")
	}
}
