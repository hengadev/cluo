package clientRepository_test

import (
	"context"
	"testing"

	th "github.com/hengadev/cluo_api/test/helpers/client"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// make test-func TEST_NAME=TestGetContactByID TEST_PATH=internal/infrastructure/postgres/client/get_contact_by_id_test.go

func TestGetContactByID(t *testing.T) {
	if testPool == nil || repo == nil {
		t.Skip("Test database or repository not initialized")
	}

	t.Run("successful retrieval", func(t *testing.T) {
		ctx := context.Background()

		// Create test contact data using helper
		contact := th.NewTestContactEncx(t)

		// Insert contact first
		err := repo.CreateContact(ctx, contact)
		require.NoError(t, err, "Failed to create contact for test")

		// Test successful contact retrieval using the global repo
		retrievedContact, err := repo.GetContactByID(ctx, contact.ID)
		assert.NoError(t, err, "Failed to get contact by ID")
		require.NotNil(t, retrievedContact, "Retrieved contact should not be nil")

		// Verify field values
		assert.Equal(t, contact.ID, retrievedContact.ID, "Contact ID should match")
		assert.Equal(t, contact.ClientID, retrievedContact.ClientID, "Client hash should match")
		assert.Equal(t, contact.EmailHash, retrievedContact.EmailHash, "Email hash should match")
		assert.Equal(t, contact.KeyVersion, retrievedContact.KeyVersion, "Key version should match")
	})

	t.Run("non-existent ID", func(t *testing.T) {
		ctx := context.Background()

		// Use a random UUID that doesn't exist in database
		nonExistentID := uuid.New()

		// Test that non-existent ID returns an error
		retrievedContact, err := repo.GetContactByID(ctx, nonExistentID)
		assert.Error(t, err, "Expected error when getting non-existent contact")
		assert.Nil(t, retrievedContact, "Retrieved contact should be nil for non-existent ID")
	})

	t.Run("nil UUID", func(t *testing.T) {
		ctx := context.Background()

		// Test with nil UUID
		retrievedContact, err := repo.GetContactByID(ctx, uuid.Nil)
		assert.Error(t, err, "Expected error when getting contact with nil UUID")
		assert.Nil(t, retrievedContact, "Retrieved contact should be nil for nil UUID")
	})

	t.Run("contact with nil optional fields", func(t *testing.T) {
		ctx := context.Background()

		// Create test contact and set optional fields to nil
		contact := th.NewTestContactEncx(t)
		contact.PhoneEncrypted = nil
		contact.PositionEncrypted = nil

		// Insert contact
		err := repo.CreateContact(ctx, contact)
		require.NoError(t, err, "Failed to create contact with nil optional fields")

		// Retrieve the contact
		retrievedContact, err := repo.GetContactByID(ctx, contact.ID)
		assert.NoError(t, err, "Failed to get contact with nil optional fields")
		require.NotNil(t, retrievedContact, "Retrieved contact should not be nil")

		// Verify that optional fields are indeed nil
		assert.Nil(t, retrievedContact.PhoneEncrypted, "Expected PhoneEncrypted to be nil")
		assert.Nil(t, retrievedContact.PositionEncrypted, "Expected PositionEncrypted to be nil")
	})

	t.Run("context cancellation", func(t *testing.T) {
		// Create a context that will be cancelled
		ctx, cancel := context.WithCancel(context.Background())
		cancel() // Cancel immediately

		contactID := uuid.New()

		retrievedContact, err := repo.GetContactByID(ctx, contactID)
		assert.Error(t, err, "Expected context cancellation error, but got nil")
		assert.Nil(t, retrievedContact, "Retrieved contact should be nil due to context cancellation")
	})
}
