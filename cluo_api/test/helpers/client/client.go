package clientHelpers

import (
	"testing"
	"time"

	"github.com/hengadev/cluo_api/internal/domain/client"

	"github.com/google/uuid"
)

// NewTestContact creates a Contact domain object with basic test data (plaintext fields only)
func NewTestClient(t *testing.T) *client.Client {
	t.Helper()

	return &client.Client{
		ID:         uuid.New(),
		Name:       "client_name",
		Type:       client.ClientTypePerson,
		ContactIDs: []string{},
		CreatedAt:  time.Now(),
	}
}

// NewTestClientEncx creates a mock ClientEncx domain object with basic test data (plaintext fields only)
func NewTestClientEncx(t *testing.T) *client.ClientEncx {
	t.Helper()
	return &client.ClientEncx{
		ID:                  uuid.New(),
		NameEncrypted:       []byte("name_encrypted"),
		NameHash:            "name_hash",
		TypeEncrypted:       []byte("type_encrypted"),
		TypeHash:            "type_hash",
		ContactIDsEncrypted: []byte("contact_ids_encrypted"),
		CreatedAt:           time.Now(),
		DEKEncrypted:        []byte("dek_encrypted"),
		KeyVersion:          1,
	}
}
