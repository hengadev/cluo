package token

import (
	"crypto/rand"
	"crypto/sha256"
	"encoding/hex"
	"time"

	"github.com/google/uuid"
)

const TokenExpiryDays = 30

type Token struct {
	ID        uuid.UUID
	CaseID    uuid.UUID
	TokenHash string // SHA-256 hex of raw token, stored in DB
	ExpiresAt time.Time
	RevokedAt *time.Time
	CreatedAt time.Time
}

// GenerateRawToken creates a cryptographically random 32-byte token.
// Returns (rawToken, sha256Hex).
func GenerateRawToken() (string, string, error) {
	b := make([]byte, 32)
	if _, err := rand.Read(b); err != nil {
		return "", "", err
	}
	rawToken := hex.EncodeToString(b)
	hash := sha256.Sum256([]byte(rawToken))
	tokenHash := hex.EncodeToString(hash[:])
	return rawToken, tokenHash, nil
}

// HashToken returns the SHA-256 hex hash of a raw token string.
func HashToken(rawToken string) string {
	hash := sha256.Sum256([]byte(rawToken))
	return hex.EncodeToString(hash[:])
}

// IsValid returns true if the token is not expired and not revoked.
func (t *Token) IsValid() bool {
	if t.RevokedAt != nil {
		return false
	}
	return time.Now().Before(t.ExpiresAt)
}
