package caseDomain

import (
	"fmt"
	"time"

	"github.com/google/uuid"
)

type Case struct {
	ID                uuid.UUID  `db:"id"`
	Title             string     `encx:"encrypt" db:"title_encrypted"`
	Description       string     `encx:"encrypt" db:"description_encrypted"`
	ClientID          uuid.UUID  `db:"client_id"`
	AssignedContactID *uuid.UUID `db:"assigned_contact_id"`
	ExternalReference *string    `encx:"encrypt" db:"external_reference"`
	CaseType          string     `db:"case_type"`
	Status            CaseStatus `encx:"encrypt" db:"status_encrypted"`
	CreatedAt         time.Time  `db:"created_at"`
	UpdatedAt         time.Time  `encx:"encrypt" db:"updated_at_encrypted"`
}

func (c *Case) MarkAsReleased() {
	c.Status = CaseStatusReleased
}

// CaseFilter represents filtering options for case queries.
type CaseFilter struct {
	ClientID          *uuid.UUID  `json:"client_id,omitempty"`
	Status            *CaseStatus `json:"status,omitempty"`
	AssignedContactID *uuid.UUID  `json:"assigned_contact_id,omitempty"`
	CaseType          *string     `json:"case_type,omitempty"`
	DateCreatedFrom   *time.Time  `json:"date_created_from,omitempty"`
	DateCreatedTo     *time.Time  `json:"date_created_to,omitempty"`
	DateUpdatedFrom   *time.Time  `json:"date_updated_from,omitempty"`
	DateUpdatedTo     *time.Time  `json:"date_updated_to,omitempty"`
	Search            *string     `json:"search,omitempty"`
}

// Pagination represents pagination parameters.
type Pagination struct {
	Page     int `json:"page" validate:"min=1"`
	PageSize int `json:"page_size" validate:"min=1,max=100"`
}

// NewPagination creates a new pagination with default values.
func NewPagination() Pagination {
	return Pagination{
		Page:     1,
		PageSize: 20,
	}
}

// GetOffset calculates the offset for database queries.
func (p Pagination) GetOffset() int {
	return (p.Page - 1) * p.PageSize
}

// Validate performs validation on the pagination parameters.
func (p Pagination) Validate() error {
	if p.Page < 1 {
		return fmt.Errorf("page must be at least 1")
	}
	if p.PageSize < 1 {
		return fmt.Errorf("page size must be at least 1")
	}
	if p.PageSize > 100 {
		return fmt.Errorf("page size cannot exceed 100")
	}
	return nil
}
