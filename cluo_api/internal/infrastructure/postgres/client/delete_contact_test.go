package clientRepository_test

import (
	"context"
	"testing"

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

	t.Run("successful deletion", func(t *testing.T) {
		ctx := context.Background()

		// Create test contact data using helper
		contact := th.NewTestContactEncx(t)

		// Insert contact first
		err := repo.CreateContact(ctx, contact)
		require.NoError(t, err, "Failed to create contact for test")

		// Verify the contact exists
		retrievedContact, err := th.GetContactEncxByID(t, ctx, testPool, contact.ID)
		assert.NoError(t, err, "Failed to retrieve inserted contact")
		assert.NotNil(t, retrievedContact, "Contact should exist before deletion")

		// Test successful contact deletion using the global repo
		err = repo.DeleteContact(ctx, contact.ID)
		assert.NoError(t, err, "Failed to delete contact")

		// Verify the contact was deleted by trying to retrieve it
		deletedContact, err := th.GetContactEncxByID(t, ctx, testPool, contact.ID)
		assert.Error(t, err, "Should get error when trying to retrieve deleted contact")
		assert.Nil(t, deletedContact, "Deleted contact should be nil")
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

		// Create test contact and set optional fields to nil
		contact := th.NewTestContactEncx(t)
		contact.PhoneEncrypted = nil
		contact.PositionEncrypted = nil

		// Insert contact with nil optional fields
		err := repo.CreateContact(ctx, contact)
		require.NoError(t, err, "Failed to create contact with nil optional fields")

		// Verify the contact exists
		retrievedContact, err := th.GetContactEncxByID(t, ctx, testPool, contact.ID)
		assert.NoError(t, err, "Failed to retrieve inserted contact")
		assert.Nil(t, retrievedContact.PhoneEncrypted, "Phone should be nil")
		assert.Nil(t, retrievedContact.PositionEncrypted, "Position should be nil")

		// Test deletion of contact with nil optional fields
		err = repo.DeleteContact(ctx, contact.ID)
		assert.NoError(t, err, "Failed to delete contact with nil optional fields")

		// Verify the contact was deleted
		deletedContact, err := th.GetContactEncxByID(t, ctx, testPool, contact.ID)
		assert.Error(t, err, "Should get error when trying to retrieve deleted contact")
		assert.Nil(t, deletedContact, "Deleted contact should be nil")
	})

	t.Run("delete multiple contacts", func(t *testing.T) {
		ctx := context.Background()

		// Create multiple test contacts
		numContacts := 3
		contactIDs := make([]uuid.UUID, numContacts)

		for i := 0; i < numContacts; i++ {
			contact := th.NewTestContactEncx(t)
			contactIDs[i] = contact.ID

			// Insert contact
			err := repo.CreateContact(ctx, contact)
			require.NoError(t, err, "Failed to create contact %d", i)
		}

		// Verify all contacts exist
		for i, contactID := range contactIDs {
			retrievedContact, err := th.GetContactEncxByID(t, ctx, testPool, contactID)
			assert.NoError(t, err, "Failed to retrieve contact %d", i)
			assert.NotNil(t, retrievedContact, "Contact %d should exist", i)
		}

		// Delete all contacts
		for i, contactID := range contactIDs {
			err := repo.DeleteContact(ctx, contactID)
			assert.NoError(t, err, "Failed to delete contact %d", i)
		}

		// Verify all contacts were deleted
		for i, contactID := range contactIDs {
			deletedContact, err := th.GetContactEncxByID(t, ctx, testPool, contactID)
			assert.Error(t, err, "Should get error when trying to retrieve deleted contact %d", i)
			assert.Nil(t, deletedContact, "Deleted contact %d should be nil", i)
		}
	})

	t.Run("delete and recreate same contact", func(t *testing.T) {
		ctx := context.Background()

		// Create test contact
		contact := th.NewTestContactEncx(t)

		// Insert contact
		err := repo.CreateContact(ctx, contact)
		require.NoError(t, err, "Failed to create contact")

		// Verify the contact exists
		retrievedContact, err := th.GetContactEncxByID(t, ctx, testPool, contact.ID)
		assert.NoError(t, err, "Failed to retrieve inserted contact")
		assert.NotNil(t, retrievedContact, "Contact should exist before deletion")

		// Delete contact
		err = repo.DeleteContact(ctx, contact.ID)
		assert.NoError(t, err, "Failed to delete contact")

		// Verify the contact was deleted
		deletedContact, err := th.GetContactEncxByID(t, ctx, testPool, contact.ID)
		assert.Error(t, err, "Should get error when trying to retrieve deleted contact")
		assert.Nil(t, deletedContact, "Deleted contact should be nil")

		// Create a new contact with the same ID (this should work since we deleted the old one)
		newContact := th.NewTestContactEncx(t)
		newContact.ID = contact.ID // Use the same ID

		err = repo.CreateContact(ctx, newContact)
		assert.NoError(t, err, "Failed to recreate contact with same ID")

		// Verify the new contact exists
		recreatedContact, err := th.GetContactEncxByID(t, ctx, testPool, contact.ID)
		assert.NoError(t, err, "Failed to retrieve recreated contact")
		assert.NotNil(t, recreatedContact, "Recreated contact should exist")
		assert.Equal(t, newContact.ClientIDHash, recreatedContact.ClientIDHash, "Client hash should match")
		assert.Equal(t, newContact.EmailHash, recreatedContact.EmailHash, "Email hash should match")
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

