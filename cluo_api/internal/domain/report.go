package domain

import (
	"encoding/json"
	"time"

	"github.com/google/uuid"
)

type Report struct {
	ID        uuid.UUID
	CaseID    uuid.UUID       `encx:"encrypt"`
	Content   json.RawMessage `encx:"encrypt"`
	CreatedAt time.Time
	UpdatedAt time.Time `encx:"encrypt"`
}
