package domain

import (
	"time"

	"github.com/google/uuid"
)

// TODO: make an enum for the type Client (person, lawyer, insurance)

type Client struct {
	ID        uuid.UUID
	Name      string `encx:"encrypt"`
	Email     string `encx:"hash_basic"`
	Type      string `encx:"encrypt"`
	CreatedAt time.Time
}
