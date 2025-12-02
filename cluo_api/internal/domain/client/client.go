package client

import (
	"time"

	"github.com/google/uuid"
)

type ClientType string

const (
	ClientTypePerson    ClientType = "person"
	ClientTypeInsurance ClientType = "insurance"
	ClientTypeLawyer    ClientType = "lawyer"
	// etc.
)

type Client struct {
	ID         uuid.UUID
	Name       string   `json:"name" encx:"encrypt,hash_basic"`
	Type       string   `json:"type" encx:"encrypt,hash_basic"`
	ContactIDs []string `json:"contacts" encx:"encrypt"` // the list of contact IDs
	CreatedAt  time.Time
}
