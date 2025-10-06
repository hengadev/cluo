package domain

import (
	"time"

	"github.com/google/uuid"
)

type AccessToken struct {
	ID        uuid.UUID
	ClientID  uuid.UUID `json:"clientID"`
	CaseID    uuid.UUID `encx:"encrypt"`
	Token     string    `encx:"encrypt"`
	ExpiresAt time.Time
	CreatedAt time.Time
}

func (t *AccessToken) IsExpired() bool {
	return time.Now().After(t.ExpiresAt)
}
