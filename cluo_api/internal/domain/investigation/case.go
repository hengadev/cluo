package investigation

import (
	"fmt"
	"time"

	"github.com/google/uuid"
)

type Investigation struct {
	ID                uuid.UUID  `db:"id"`
	Title             string     `encx:"encrypt" db:"title_encrypted"`
	Description       string     `encx:"encrypt" db:"description_encrypted"`
	ClientID          uuid.UUID  `db:"client_id"`
	AssignedContactID *uuid.UUID `db:"assigned_contact_id"`
	CaseSubjectID     *uuid.UUID `db:"case_subject_id"`
	ExternalReference *string    `encx:"encrypt,hash_basic" db:"external_reference_encrypted"`
	CaseType          string     `db:"case_type"`
	Status            Status `encx:"encrypt" db:"status_encrypted"`
	Placename         string     `encx:"encrypt,hash_basic" db:"placename_encrypted"`
	Address1          string     `encx:"encrypt,hash_basic" db:"address1_encrypted"`
	Address2          string     `encx:"encrypt,hash_basic" db:"address2_encrypted"`
	City              string     `encx:"encrypt,hash_basic" db:"city_encrypted"`
	PostalCode        string     `encx:"encrypt,hash_basic" db:"postal_code_encrypted"`
	Country           string     `encx:"encrypt,hash_basic" db:"country_encrypted"`
	Latitude          *string    `encx:"encrypt,hash_basic" db:"latitude_encrypted"`
	Longitude         *string    `encx:"encrypt,hash_basic" db:"longitude_encrypted"`
	LocationType      string     `encx:"encrypt,hash_basic" db:"location_type_encrypted"`
	LocationNotes     string     `encx:"encrypt,hash_basic" db:"location_notes_encrypted"`
	CreatedAt         time.Time  `db:"created_at"`
	UpdatedAt         time.Time  `encx:"encrypt" db:"updated_at_encrypted"`
}

func (c *Investigation) MarkAsReleased() {
	c.Status = StatusReleased
}

// Filter represents filtering options for case queries.
type Filter struct {
	ClientID          *uuid.UUID  `json:"client_id,omitempty"`
	Status            *Status `json:"status,omitempty"`
	AssignedContactID *uuid.UUID  `json:"assigned_contact_id,omitempty"`
	CaseSubjectID     *uuid.UUID  `json:"case_subject_id,omitempty"`
	CaseType          *string     `json:"case_type,omitempty"`
	City              *string     `json:"city,omitempty"`
	PostalCode        *string     `json:"postal_code,omitempty"`
	Country           *string     `json:"country,omitempty"`
	DateCreatedFrom   *time.Time  `json:"date_created_from,omitempty"`
	DateCreatedTo     *time.Time  `json:"date_created_to,omitempty"`
	DateUpdatedFrom   *time.Time  `json:"date_updated_from,omitempty"`
	DateUpdatedTo     *time.Time  `json:"date_updated_to,omitempty"`
	Search            *string     `json:"search,omitempty"`
	// Hash fields (populated by application layer for repository filtering)
	CityHash       *string `json:"-"`
	PostalCodeHash *string `json:"-"`
	CountryHash    *string `json:"-"`
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
