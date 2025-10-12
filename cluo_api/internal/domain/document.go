package domain

import (
	"time"

	"github.com/google/uuid"
)

// TODO: complete the different document type fields to reflect the template created.

type Document struct {
	ID        uuid.UUID
	CaseID    uuid.UUID    `encx:"encrypt"`
	Type      DocumentType `encx:"encrypt"`
	CreatedAt time.Time
}

// The shared interface — all concrete document types implement this
type Documentable interface {
	GetID() uuid.UUID
	GetType() DocumentType
	GetCaseID() uuid.UUID
	Validate() error
	// ToDTO() DocumentDTO
}

type Invoice struct {
	Document
	Amount      float64
	Description string
	DueDate     time.Time
	Issuer      string
	Recipient   string
}

type Mandate struct {
	Document
	ClientName string
	Scope      string
	StartDate  time.Time
	EndDate    time.Time
}

type Report struct {
	Document
	Summary         string
	Findings        string
	Recommendations string
}
