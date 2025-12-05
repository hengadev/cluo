package client

import (
	"time"

	"github.com/google/uuid"
)

type ClientType string

const (
	ClientTypePerson     ClientType = "person"
	ClientTypeInsurance  ClientType = "insurance"
	ClientTypeLawyer     ClientType = "lawyer"
	ClientTypeCompany    ClientType = "company"
	ClientTypeGovernment ClientType = "government"
)

type Client struct {
	ID        uuid.UUID
	Name      string     `json:"name" encx:"encrypt,hash_basic"`
	Type      ClientType `json:"type" encx:"encrypt,hash_basic"`
	CreatedAt time.Time
}
