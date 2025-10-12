package domain

import (
	"time"

	"github.com/google/uuid"
)

type Client struct {
	ID        uuid.UUID
	Name      string     `encx:"encrypt"`
	Email     string     `encx:"hash_basic"`
	Type      ClientType `encx:"encrypt"`
	CreatedAt time.Time
}
