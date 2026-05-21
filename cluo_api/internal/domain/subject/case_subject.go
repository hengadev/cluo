package subject

import (
	"time"

	"github.com/google/uuid"
)

type Subject struct {
	ID         uuid.UUID `db:"id"`
	Lastname   string    `encx:"encrypt,hash_basic" db:"lastname_encrypted"`
	Firstname  string    `encx:"encrypt,hash_basic" db:"firstname_encrypted"`
	Email      string    `encx:"encrypt,hash_basic" db:"email_encrypted"`
	Phone      string    `encx:"encrypt" db:"phone_encrypted"`
	City       string    `encx:"encrypt,hash_basic" db:"city_encrypted"`
	PostalCode string    `encx:"encrypt,hash_basic" db:"postal_code_encrypted"`
	Address1   string    `encx:"encrypt,hash_basic" db:"address1_encrypted"`
	Address2   string    `encx:"encrypt,hash_basic" db:"address2_encrypted"`
	Occupation string    `encx:"encrypt,hash_basic" db:"occupation_encrypted"`
	Notes      string    `encx:"encrypt" db:"notes_encrypted"`

	CreatedAt time.Time `db:"created_at"`
}
