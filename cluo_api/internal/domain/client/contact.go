package client

import (
	"time"

	"github.com/google/uuid"
)

type Contact struct {
	ID        uuid.UUID `json:"-"`
	ClientID  uuid.UUID `json:"-"`
	Lastname  string    `json:"lastname" encx:"encrypt"`
	Firstname string    `json:"firstname" encx:"encrypt"`
	Email     string    `json:"email" encx:"encrypt,hash_basic"`
	Phone     string    `json:"phone" encx:"encrypt"`
	Position  string    `json:"position" encx:"encrypt"`
	CreatedAt time.Time
}
