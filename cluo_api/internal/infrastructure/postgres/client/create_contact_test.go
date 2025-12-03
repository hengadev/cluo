package clientRepository_test

import (
	"context"
	"testing"

	"github.com/hengadev/cluo_api/internal/domain/client"
	th "github.com/hengadev/cluo_api/test/helpers/client"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// make test-func TEST_NAME=TestCreateContact TEST_PATH=internal/infrastructure/postgres/client/create_contact_test.go

func TestCreateContact(t *testing.T) {
	if testPool == nil || repo == nil {
		t.Skip("Test database or repository not initialized")
	}

	t.Run("successful creation", func(t *testing.T) {
		ctx := context.Background()

		// Create test clientEncx data using helper
		clientEncx := th.NewTestClientEncx(t)
		err := th.InsertClientEncx(t, ctx, testPool, clientEncx)
		require.NoError(t, err)

		// Create test contact data using helper
		contact := th.NewTestContactEncx(t)
		contact.ClientIDHash = clientEncx.ID.String()

		// Test successful contact creation using the global repo
		err = repo.CreateContact(ctx, contact)
		assert.NoError(t, err, "Failed to create contact")

		// Verify the contact was inserted by retrieving it
		retrievedContact, err := th.GetContactEncxByID(t, ctx, testPool, contact.ID)
		assert.NoError(t, err, "Failed to retrieve inserted contact")

		// Verify field values
		assert.Equal(t, contact.ID, retrievedContact.ID, "Contact ID should match")
		assert.Equal(t, contact.ClientIDHash, retrievedContact.ClientIDHash, "Client hash should match")
		assert.Equal(t, contact.EmailHash, retrievedContact.EmailHash, "Email hash should match")
		assert.Equal(t, contact.KeyVersion, retrievedContact.KeyVersion, "Key version should match")
	})

	t.Run("with nil optional values", func(t *testing.T) {
		ctx := context.Background()

		// Create test contact and set optional fields to nil
		contact := th.NewTestContactEncx(t)
		contact.PhoneEncrypted = nil
		contact.PositionEncrypted = nil

		// Test successful contact creation with nil optional fields using the global repo
		err := repo.CreateContact(ctx, contact)
		require.NoError(t, err, "Failed to create contact with nil values")

		// Verify the contact was inserted by retrieving it
		retrievedContact, err := th.GetContactEncxByID(t, ctx, testPool, contact.ID)
		assert.NoError(t, err, "Failed to retrieve inserted contact")

		// Verify that optional fields are indeed nil
		assert.Nil(t, retrievedContact.PhoneEncrypted, "Expected PhoneEncrypted to be nil")
		assert.Nil(t, retrievedContact.PositionEncrypted, "Expected PositionEncrypted to be nil")
	})

	t.Run("duplicate ID", func(t *testing.T) {
		ctx := context.Background()

		// Create first test contact
		contact1 := th.NewTestContactEncx(t)

		// Insert first contact using the global repo (setup - should stop if fails)
		err := repo.CreateContact(ctx, contact1)
		require.NoError(t, err, "Failed to create first contact")

		// Try to insert contact with same ID (should fail)
		contact2 := th.NewTestContactEncx(t)
		contact2.ID = contact1.ID // Same ID, different data

		err = repo.CreateContact(ctx, contact2)
		assert.Error(t, err, "Expected error when creating contact with duplicate ID")

		// Check that it's a database constraint violation (expected for duplicate ID)
		errStr := err.Error()
		assert.True(t, contains(errStr, "duplicate") || contains(errStr, "unique") || contains(errStr, "constraint"),
			"Expected constraint violation error, got: %v", err)
	})

	t.Run("empty required fields", func(t *testing.T) {
		ctx := context.Background()

		tests := []struct {
			name    string
			contact *client.ContactEncx
		}{
			{
				name: "nil uuid",
				contact: &client.ContactEncx{
					ID: uuid.Nil, // Invalid UUID
				},
			},
		}

		for _, tt := range tests {
			t.Run(tt.name, func(t *testing.T) {
				err := repo.CreateContact(ctx, tt.contact)
				assert.Error(t, err, "Expected error for %s, but got nil", tt.name)
			})
		}
	})

	t.Run("context cancellation", func(t *testing.T) {
		// Create a context that will be cancelled
		ctx, cancel := context.WithCancel(context.Background())
		cancel() // Cancel immediately

		contact := th.NewTestContactEncx(t)

		err := repo.CreateContact(ctx, contact)
		assert.Error(t, err, "Expected context cancellation error, but got nil")
	})
}

// contains checks if a string contains a substring
func contains(s, substr string) bool {
	return len(s) >= len(substr) && len(s) > 0 && len(substr) > 0 &&
		func() bool {
			for i := 0; i <= len(s)-len(substr); i++ {
				if s[i:i+len(substr)] == substr {
					return true
				}
			}
			return false
		}()
}
