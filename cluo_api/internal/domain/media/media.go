package domain

import (
	"time"

	"github.com/google/uuid"
)

type MediaFile struct {
	ID          uuid.UUID
	CaseID      uuid.UUID
	URL         string    `encx:"encrypt"`
	Type        MediaType `encx:"encrypt"`
	MimeType    string    `encx:"encrypt"`
	FileName    string    `encx:"encrypt"`
	FileSize    int64
	Caption     string `encx:"encrypt"`
	IsPublished bool
	CreatedAt   time.Time
}
