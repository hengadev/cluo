package domain

import (
	"time"

	"github.com/google/uuid"
)

// TODO: have the roles being an enum also
// "admin", "investigator"

type User struct {
	ID        uuid.UUID
	Email     string `encx:"hash_basic"`
	Password  string `encx:"hash_secure"`
	Role      string `encx:"encrypt"`
	CreatedAt time.Time
}
