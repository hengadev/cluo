package clientRepository_test

import (
	"context"
	"testing"
	"time"

	"github.com/hengadev/cluo_api/internal/domain/client"
	th "github.com/hengadev/cluo_api/test/helpers/client"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// make test-func TEST_NAME=TestDeleteContact TEST_PATH=internal/infrastructure/postgres/client/delete_contact_test.go

func TestDeleteContact(t *testing.T) {
	if testPool == nil || repo == nil {
		t.Skip("Test database or repository not initialized")
	}

	setupClient := func(t *testing.T, ctx context.Context) *client.ClientEncx {
		clientEncx := th.NewTestClientEncx(t)
		err := th.InsertClientEncx(t, ctx, testPool, clientEncx)
		require.NoError(t, err)
		return clientEncx
	}

	t.Run("successful deletion", func(t *testing.T) {
		ctx := context.Background()

		th.ClearClientsTable(t, ctx, testPool)
		th.ClearContactsTable(t, ctx, testPool)

		// Create test clientEncx data using helper
		clientEncx := setupClient(t, ctx)

		// Create test contactEncx data using helper
		contactEncx := th.NewTestContactEncx(t)
		contactEncx.ClientID = clientEncx.ID

		// Insert contact first
		err := th.InsertContactEncx(t, ctx, testPool, contactEncx)
		require.NoError(t, err, "Failed to create contact for test")

		// Test successful contact deletion using the global repo
		err = repo.DeleteContact(ctx, contactEncx.ID)
		assert.NoError(t, err, "Failed to delete contact")

		// Verify the contact was deleted by trying to retrieve it
		deletedContactEncx, err := th.GetContactEncxByID(t, ctx, testPool, contactEncx.ID)
		assert.Error(t, err, "Should get error when trying to retrieve deleted contact")
		assert.Equal(t, &client.ContactEncx{}, deletedContactEncx, "Deleted contact should be nil")
	})

	t.Run("non-existent contact", func(t *testing.T) {
		ctx := context.Background()

		// Use a random UUID that doesn't exist in database
		nonExistentID := uuid.New()

		// Test that non-existent contact deletion returns an error
		err := repo.DeleteContact(ctx, nonExistentID)
		assert.Error(t, err, "Expected error when deleting non-existent contact")
		assert.Contains(t, err.Error(), "record not found", "Error should indicate record not found")
	})

	t.Run("nil UUID", func(t *testing.T) {
		ctx := context.Background()

		// Test with nil UUID
		err := repo.DeleteContact(ctx, uuid.Nil)
		assert.Error(t, err, "Expected error when deleting with nil UUID")
	})

	t.Run("contact with nil optional fields", func(t *testing.T) {
		ctx := context.Background()

		th.ClearClientsTable(t, ctx, testPool)
		th.ClearContactsTable(t, ctx, testPool)

		// Create test clientEncx data using helper
		clientEncx := setupClient(t, ctx)

		// Create test contactEncx and set optional fields to nil
		contactEncx := th.NewTestContactEncx(t)
		contactEncx.ClientID = clientEncx.ID
		contactEncx.PhoneEncrypted = nil
		contactEncx.PositionEncrypted = nil

		// Insert contact with nil optional fields
		err := th.InsertContactEncx(t, ctx, testPool, contactEncx)
		require.NoError(t, err, "Failed to create contact with nil optional fields")

		// Verify the contact exists
		retrievedContactEncx, err := th.GetContactEncxByID(t, ctx, testPool, contactEncx.ID)
		assert.NoError(t, err, "Failed to retrieve inserted contact")
		assert.Nil(t, retrievedContactEncx.PhoneEncrypted, "Phone should be nil")
		assert.Nil(t, retrievedContactEncx.PositionEncrypted, "Position should be nil")

		// Test deletion of contact with nil optional fields
		err = repo.DeleteContact(ctx, contactEncx.ID)
		assert.NoError(t, err, "Failed to delete contact with nil optional fields")

		// Verify the contact was deleted
		deletedContactEncx, err := th.GetContactEncxByID(t, ctx, testPool, contactEncx.ID)
		assert.Error(t, err, "Should get error when trying to retrieve deleted contact")
		assert.Equal(t, &client.ContactEncx{}, deletedContactEncx, "Deleted contact should be nil")
	})

	t.Run("delete multiple contacts", func(t *testing.T) {
		ctx := context.Background()

		th.ClearClientsTable(t, ctx, testPool)
		th.ClearContactsTable(t, ctx, testPool)

		// Create test clientEncx data using helper
		clientEncx := setupClient(t, ctx)

		// Create multiple test contacts
		numContacts := 3
		contactEncxIDs := make([]uuid.UUID, numContacts)

		for i := 0; i < numContacts; i++ {
			contactEncx := th.NewTestContactEncx(t)
			contactEncx.ClientID = clientEncx.ID
			contactEncxIDs[i] = contactEncx.ID

			// Insert contact
			err := th.InsertContactEncx(t, ctx, testPool, contactEncx)
			require.NoError(t, err, "Failed to create contact %d", i)
		}

		// Delete all contacts
		for i, contactID := range contactEncxIDs {
			err := repo.DeleteContact(ctx, contactID)
			assert.NoError(t, err, "Failed to delete contact %d", i)
		}

		// Verify all contacts were deleted
		for i, contactID := range contactEncxIDs {
			deletedContactEncx, err := th.GetContactEncxByID(t, ctx, testPool, contactID)
			assert.Error(t, err, "Should get error when trying to retrieve deleted contact %d", i)
			assert.Equal(t, &client.ContactEncx{}, deletedContactEncx, "Deleted contact %d should be nil", i)
		}
	})

	t.Run("delete and recreate same contact", func(t *testing.T) {
		ctx := context.Background()

		th.ClearClientsTable(t, ctx, testPool)
		th.ClearContactsTable(t, ctx, testPool)

		// Create test clientEncx data using helper
		clientEncx := setupClient(t, ctx)

		// Create test contactEncx
		contactEncx := th.NewTestContactEncx(t)
		contactEncx.ClientID = clientEncx.ID

		// Insert contact
		err := th.InsertContactEncx(t, ctx, testPool, contactEncx)
		require.NoError(t, err, "Failed to create contact")

		// Delete contact
		err = repo.DeleteContact(ctx, contactEncx.ID)
		assert.NoError(t, err, "Failed to delete contact")

		// Verify the contact was deleted
		deletedContactEncx, err := th.GetContactEncxByID(t, ctx, testPool, contactEncx.ID)
		assert.Error(t, err, "Should get error when trying to retrieve deleted contact")
		assert.Equal(t, &client.ContactEncx{}, deletedContactEncx, "Deleted contact should be nil")

		// Create a new contact with the same ID (this should work since we deleted the old one)
		newContactEncx := th.NewTestContactEncx(t)
		newContactEncx.ID = contactEncx.ID // Use the same ID
		newContactEncx.ClientID = clientEncx.ID

		err = th.InsertContactEncx(t, ctx, testPool, newContactEncx)
		assert.NoError(t, err, "Failed to recreate contact with same ID")

		// Verify the new contact exists
		recreatedContactEncx, err := th.GetContactEncxByID(t, ctx, testPool, contactEncx.ID)
		assert.NoError(t, err, "Failed to retrieve recreated contact")
		assert.Equal(t, newContactEncx.ID, recreatedContactEncx.ID, "ID should match")
		assert.Equal(t, newContactEncx.ClientID, recreatedContactEncx.ClientID, "Client hash should match")
		assert.Equal(t, newContactEncx.EmailHash, recreatedContactEncx.EmailHash, "Email hash should match")
		assert.Equal(t, newContactEncx.EmailEncrypted, recreatedContactEncx.EmailEncrypted, "Email encrypted should match")
		assert.Equal(t, newContactEncx.PhoneEncrypted, recreatedContactEncx.PhoneEncrypted, "Phone encrypted should match")
		assert.Equal(t, newContactEncx.PositionEncrypted, recreatedContactEncx.PositionEncrypted, "Position encrypted should match")
		assert.Equal(t, newContactEncx.LastnameEncrypted, recreatedContactEncx.LastnameEncrypted, "Lastname encrypted hash should match")
		assert.Equal(t, newContactEncx.FirstnameEncrypted, recreatedContactEncx.FirstnameEncrypted, "Firstname encrypted hash should match")
		assert.WithinDuration(t, newContactEncx.CreatedAt, recreatedContactEncx.CreatedAt, time.Second, "CreatedAt should match")
	})

	t.Run("context cancellation", func(t *testing.T) {
		// Create a context that will be cancelled
		ctx, cancel := context.WithCancel(context.Background())
		cancel() // Cancel immediately

		contactID := uuid.New()

		err := repo.DeleteContact(ctx, contactID)
		assert.Error(t, err, "Expected context cancellation error, but got nil")
	})
}
