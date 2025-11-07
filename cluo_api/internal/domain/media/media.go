package domain

import (
	"time"

	"github.com/google/uuid"
)

type MediaFile struct {
	ID          uuid.UUID
	CaseID      uuid.UUID `encx:"encrypt"`
	URL         string    `encx:"encrypt"`
	Type        MediaType `encx:"encrypt"`
	Caption     string    `encx:"encrypt"`
	IsPublished bool
	CreatedAt   time.Time
}
