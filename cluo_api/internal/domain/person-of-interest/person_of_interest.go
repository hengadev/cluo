package personofinterest

import (
	"time"

	"github.com/google/uuid"
)

type PersonOfInterest struct {
	ID         uuid.UUID
	CaseID     []uuid.UUID
	Lastname   string `encx:"encrypt,hash_basic"`
	Firstname  string `encx:"encrypt,hash_basic"`
	Email      string `encx:"encrypt,hash_basic"`
	Phone      string `encx:"encrypt"`
	City       string `encx:"encrypt,hash_basic"`
	PostalCode string `encx:"encrypt,hash_basic"`
	Address1   string `encx:"encrypt,hash_basic"`
	Address2   string `encx:"encrypt,hash_basic"`
	Occupation string `encx:"encrypt,hash_basic"`
	Notes      string `encx:"encrypt"`

	Roles []PersonRole

	CreatedAt time.Time `json:"created_at"`
}
