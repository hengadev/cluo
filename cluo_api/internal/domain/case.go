package domain

import (
	"time"

	"github.com/google/uuid"
)

type Case struct {
	ID          uuid.UUID
	Title       string     `encx:"encrypt"`
	Description string     `encx:"encrypt"`
	ClientID    string
	Status      CaseStatus `encx:"encrypt"`
	CreatedAt   time.Time
	UpdatedAt   time.Time `encx:"encrypt"`
}

func (c *Case) MarkAsReleased() {
	c.Status = CaseStatusReleased
}
