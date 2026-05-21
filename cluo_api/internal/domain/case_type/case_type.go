package case_type

import (
	"time"

	"github.com/google/uuid"
)

type CaseType struct {
	ID        uuid.UUID
	Name      string
	CreatedAt time.Time
	UpdatedAt time.Time
}

func New(name string) *CaseType {
	now := time.Now()
	return &CaseType{ID: uuid.New(), Name: name, CreatedAt: now, UpdatedAt: now}
}
