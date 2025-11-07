package client

import (
	"time"

	"github.com/google/uuid"
)

type ClientType string

// NOTE: this is just to set the first values, add that to a migration
const (
	ClientTypePerson    ClientType = "person"
	ClientTypeInsurance ClientType = "insurance"
	ClientTypeLawyer    ClientType = "lawyer"
	// etc.
)

type Client struct {
	ID   uuid.UUID
	Name string `json:"name" encx:"encrypt"`
	Type string `json:"type" encx:"encrypt"`
	// Contacts  []string `json:"contacts" encx:"encrypt"` // the list of contact IDs
	CreatedAt time.Time
}

// client : cabinet avocat
// contact : someone  that works within the client
