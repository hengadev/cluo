package domain

import (
	"time"

	"github.com/google/uuid"
)

// TODO: make the type an enum
// "image" | "video"

// TODO: add something to precise if that thing is published for the client to see

type MediaFile struct {
	ID        uuid.UUID
	CaseID    uuid.UUID `encx:"encrypt"`
	URL       string    `encx:"encrypt"`
	Type      string    `encx:"encrypt"`
	Caption   string    `encx:"encrypt"`
	CreatedAt time.Time
}
