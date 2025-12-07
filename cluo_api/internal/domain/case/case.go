package caseDomain

import (
	"time"

	"github.com/google/uuid"
)

type Case struct {
	ID                uuid.UUID  `db:"id"`
	Title             string     `encx:"encrypt" db:"title_encrypted"`
	Description       string     `encx:"encrypt" db:"description_encrypted"`
	ClientID          string     `db:"client_id"`
	AssignedContactID *string    `db:"assigned_contact_id"`
	Status            CaseStatus `encx:"encrypt" db:"status_encrypted"`
	CreatedAt         time.Time  `db:"created_at"`
	UpdatedAt         time.Time  `encx:"encrypt" db:"updated_at_encrypted"`
}

func (c *Case) MarkAsReleased() {
	c.Status = CaseStatusReleased
}
