package clientHelpers

import (
	"testing"
	"time"

	"github.com/hengadev/cluo_api/internal/domain/client"

	"github.com/google/uuid"
)

// NewTestContact creates a Contact domain object with basic test data (plaintext fields only)
func NewTestContact(t *testing.T) *client.Contact {
	t.Helper()
	return &client.Contact{
		ID:        uuid.New(),
		ClientID:  uuid.New(),
		Lastname:  "DOE",
		Firstname: "John",
		Email:     "john.doe@example.com",
		Phone:     "0612345678",
		Position:  "position",
		CreatedAt: time.Now(),
	}
}

// NewTestContactEncx creates a mock ContactEncx domain object with basic test data (plaintext fields only)
func NewTestContactEncx(t *testing.T) *client.ContactEncx {
	t.Helper()
	return &client.ContactEncx{
		ID:                 uuid.New(),
		CreatedAt:          time.Now(),
		ClientID:           uuid.New(),
		LastnameEncrypted:  []byte("lastname_encrypted"),
		FirstnameEncrypted: []byte("lastname_encrypted"),
		EmailEncrypted:     []byte("email_encrypted"),
		EmailHash:          "email_hash",
		PhoneEncrypted:     []byte("phone_encrypted"),
		PositionEncrypted:  []byte("position_encrypted"),
		DEKEncrypted:       []byte("dek_encrypted"),
		KeyVersion:         1,
	}
}
