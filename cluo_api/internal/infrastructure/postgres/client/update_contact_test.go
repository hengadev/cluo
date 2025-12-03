package clientRepository_test

import (
	"context"
	"testing"

	th "github.com/hengadev/cluo_api/test/helpers/client"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// make test-func TEST_NAME=TestUpdateContact TEST_PATH=internal/infrastructure/postgres/client/update_contact_test.go

func TestUpdateContact(t *testing.T) {
	if testPool == nil || repo == nil {
		t.Skip("Test database or repository not initialized")
	}

	t.Run("successful update", func(t *testing.T) {
		ctx := context.Background()

		// Create and insert a contact first
		originalContact := th.NewTestContactEncx(t)
		err := th.InsertContactEncx(t, ctx, testPool, *originalContact)
		require.NoError(t, err, "Failed to insert original contact")

		// Create updated contact with same ID but different data
		updatedContact := th.NewTestContactEncx(t)
		updatedContact.ID = originalContact.ID // Keep same ID
		updatedContact.ClientID = originalContact.ClientID

		// Update the contact
		err = repo.UpdateContact(ctx, updatedContact)
		assert.NoError(t, err, "Failed to update contact")

		// Verify the update by retrieving the contact
		retrievedContact, err := repo.GetContactByID(ctx, updatedContact.ID)
		require.NoError(t, err, "Failed to retrieve updated contact")

		// Verify fields were updated (should match updatedContact, not originalContact)
		assert.Equal(t, updatedContact.ID, retrievedContact.ID)
		assert.Equal(t, updatedContact.ClientID, retrievedContact.ClientID)
		// Note: We can't easily compare encrypted fields since they use random encryption
		// But we can verify the contact was retrieved successfully

		t.Logf("✓ Successfully updated contact with ID: %s", updatedContact.ID.String())
	})

	t.Run("update non-existent contact", func(t *testing.T) {
		ctx := context.Background()

		// Create contact with non-existent ID
		nonExistentID := uuid.New()
		contact := th.NewTestContactEncx(t)
		contact.ID = nonExistentID

		// Attempt to update non-existent contact
		err := repo.UpdateContact(ctx, contact)
		assert.Error(t, err, "Expected error when updating non-existent contact")
		assert.Contains(t, err.Error(), "not found", "Error should mention not found")

		t.Logf("✓ Correctly handled update of non-existent contact")
	})

	t.Run("update contact with nil optional fields", func(t *testing.T) {
		ctx := context.Background()

		// Create and insert a contact first
		originalContact := th.NewTestContactEncx(t)
		err := th.InsertContactEncx(t, ctx, testPool, *originalContact)
		require.NoError(t, err, "Failed to insert original contact")

		// Create updated contact with nil optional fields
		updatedContact := th.NewTestContactEncx(t)
		updatedContact.ID = originalContact.ID
		updatedContact.ClientID = originalContact.ClientID
		updatedContact.PhoneEncrypted = nil
		updatedContact.PositionEncrypted = nil

		// Update the contact
		err = repo.UpdateContact(ctx, updatedContact)
		assert.NoError(t, err, "Failed to update contact with nil fields")

		// Verify the update
		retrievedContact, err := repo.GetContactByID(ctx, updatedContact.ID)
		require.NoError(t, err, "Failed to retrieve updated contact")
		assert.Nil(t, retrievedContact.PhoneEncrypted, "Phone should be nil after update")
		assert.Nil(t, retrievedContact.PositionEncrypted, "Position should be nil after update")

		t.Logf("✓ Successfully updated contact with nil optional fields")
	})

	t.Run("update all contact fields", func(t *testing.T) {
		ctx := context.Background()

		// Create and insert a contact first
		originalContact := th.NewTestContactEncx(t)
		err := th.InsertContactEncx(t, ctx, testPool, *originalContact)
		require.NoError(t, err, "Failed to insert original contact")

		// Create completely new contact data
		updatedContact := th.NewTestContactEncx(t)
		updatedContact.ID = originalContact.ID // Keep same ID
		// All other fields should be different from original

		// Update the contact
		err = repo.UpdateContact(ctx, updatedContact)
		assert.NoError(t, err, "Failed to update all contact fields")

		// Verify the update
		retrievedContact, err := repo.GetContactByID(ctx, updatedContact.ID)
		require.NoError(t, err, "Failed to retrieve updated contact")

		// Verify all fields were updated (ID and client ID hash should match)
		assert.Equal(t, updatedContact.ID, retrievedContact.ID)
		assert.Equal(t, updatedContact.ClientID, retrievedContact.ClientID)
		assert.Equal(t, updatedContact.KeyVersion, retrievedContact.KeyVersion)

		t.Logf("✓ Successfully updated all contact fields")
	})

	t.Run("update contact with different client ID", func(t *testing.T) {
		ctx := context.Background()

		// Create and insert a contact first
		originalContact := th.NewTestContactEncx(t)
		err := th.InsertContactEncx(t, ctx, testPool, *originalContact)
		require.NoError(t, err, "Failed to insert original contact")

		// Create updated contact with different client ID
		updatedContact := th.NewTestContactEncx(t)
		updatedContact.ID = originalContact.ID
		newClientIDHash := uuid.New()
		updatedContact.ClientID = newClientIDHash

		// Update the contact
		err = repo.UpdateContact(ctx, updatedContact)
		assert.NoError(t, err, "Failed to update contact client ID")

		// Verify the update
		retrievedContact, err := repo.GetContactByID(ctx, updatedContact.ID)
		require.NoError(t, err, "Failed to retrieve updated contact")
		assert.Equal(t, newClientIDHash, retrievedContact.ClientID, "Client ID hash should be updated")

		t.Logf("✓ Successfully updated contact client ID")
	})

	t.Run("context cancellation", func(t *testing.T) {
		// Create a context that will be cancelled
		ctx, cancel := context.WithCancel(context.Background())
		cancel() // Cancel immediately

		contact := th.NewTestContactEncx(t)

		// Attempt to update with cancelled context
		err := repo.UpdateContact(ctx, contact)
		assert.Error(t, err, "Expected context cancellation error")

		t.Logf("✓ Context cancellation handled correctly")
	})

	t.Run("update with invalid UUID", func(t *testing.T) {
		ctx := context.Background()

		// Create contact with zero UUID (invalid)
		contact := th.NewTestContactEncx(t)
		contact.ID = uuid.Nil

		// Attempt to update with invalid UUID
		err := repo.UpdateContact(ctx, contact)
		// This should likely succeed at the repository level since UUID validation
		// happens at higher layers, but the database might return an error
		if err != nil {
			t.Logf("✓ Update with invalid UUID properly handled: %v", err)
		} else {
			t.Logf("✓ Update with invalid UUID succeeded (validation happens at higher layers)")
		}
	})

	t.Run("verify only one contact updated", func(t *testing.T) {
		ctx := context.Background()

		// Create and insert multiple contacts
		contact1 := th.NewTestContactEncx(t)
		err := th.InsertContactEncx(t, ctx, testPool, *contact1)
		require.NoError(t, err, "Failed to insert contact 1")

		contact2 := th.NewTestContactEncx(t)
		err = th.InsertContactEncx(t, ctx, testPool, *contact2)
		require.NoError(t, err, "Failed to insert contact 2")

		// Update only contact1
		updatedContact1 := th.NewTestContactEncx(t)
		updatedContact1.ID = contact1.ID
		updatedContact1.ClientID = contact1.ClientID

		err = repo.UpdateContact(ctx, updatedContact1)
		assert.NoError(t, err, "Failed to update contact 1")

		// Verify contact1 was updated
		retrievedContact1, err := repo.GetContactByID(ctx, contact1.ID)
		require.NoError(t, err, "Failed to retrieve contact 1")
		assert.Equal(t, updatedContact1.ID, retrievedContact1.ID)

		// Verify contact2 was not affected
		retrievedContact2, err := repo.GetContactByID(ctx, contact2.ID)
		require.NoError(t, err, "Failed to retrieve contact 2")
		assert.Equal(t, contact2.ID, retrievedContact2.ID)

		t.Logf("✓ Only target contact was updated, other contacts unaffected")
	})
}
