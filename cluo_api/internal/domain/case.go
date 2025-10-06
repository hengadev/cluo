package domain

import (
	"time"

	"github.com/google/uuid"
)

type Case struct {
	ID          uuid.UUID
	Title       string `encx:"encrypt"`
	Description string `encx:"encrypt"`
	ClientID    string
	Status      string `encx:"encrypt"`
	CreatedAt   time.Time
	UpdatedAt   time.Time `encx:"encrypt"`
}

// TODO: have the status being an enum please so that I can use constant values for this
// "draft", "in_progress", "ready", "released"

func (c *Case) MarkAsReleased() {
	c.Status = "released"
}
