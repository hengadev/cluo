package document

import (
	"time"

	"github.com/google/uuid"
)

// DocumentBase contains metadata shared by all document types.
type DocumentBase struct {
	ID        uuid.UUID      `json:"id" db:"id"`
	CaseID    uuid.UUID      `json:"case_id" db:"case_id" encx:"encrypt"`
	ClientID  uuid.UUID      `json:"client_id" db:"client_id" encx:"encrypt"`
	Status    DocumentStatus `json:"status" db:"status"`
	CreatedAt time.Time      `json:"created_at" db:"created_at"`
	UpdatedAt time.Time      `json:"updated_at" db:"updated_at"`
}

// GetID returns the document's unique identifier.
func (d DocumentBase) GetID() uuid.UUID {
	return d.ID
}

// GetCaseID returns the document's associated case ID.
func (d DocumentBase) GetCaseID() uuid.UUID {
	return d.CaseID
}

// GetStatus returns the current status of the document.
func (d DocumentBase) GetStatus() DocumentStatus {
	return d.Status
}

// SetStatus updates the document status and timestamps.
func (d *DocumentBase) SetStatus(status DocumentStatus) {
	d.Status = status
	d.UpdatedAt = time.Now()
}

// UpdateTimestamp updates the document's updated_at timestamp.
func (d *DocumentBase) UpdateTimestamp() {
	d.UpdatedAt = time.Now()
}

// NewDocumentBase creates a new DocumentBase with default values.
func NewDocumentBase(caseID, clientID uuid.UUID) DocumentBase {
	now := time.Now()
	return DocumentBase{
		ID:        uuid.New(),
		CaseID:    caseID,
		ClientID:  clientID,
		Status:    DocumentStatusDraft,
		CreatedAt: now,
		UpdatedAt: now,
	}
}
