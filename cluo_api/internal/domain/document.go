package domain

import (
	"time"

	"github.com/google/uuid"
)

// TODO: have the document type be an enum
// "invoice", "mandate", etc.

type Document struct {
	ID        uuid.UUID
	CaseID    uuid.UUID `encx:"encrypt"`
	URL       string    `encx:"encrypt"`
	Type      string    `encx:"encrypt"`
	CreatedAt time.Time
}
