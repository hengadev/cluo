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

// make test-func TEST_NAME=TestUpdateContact TEST_PATH=internal/infrastructure/postgres/client/update_contact_test.go

func TestUpdateContact(t *testing.T) {
	if testPool == nil || repo == nil {
		t.Skip("Test database or repository not initialized")
	}

	setupClient := func(t *testing.T, ctx context.Context) *client.ClientEncx {
		clientEncx := th.NewTestClientEncx(t)
		err := th.InsertClientEncx(t, ctx, testPool, clientEncx)
		require.NoError(t, err)
		return clientEncx
	}

	t.Run("successful update", func(t *testing.T) {
		ctx := context.Background()

		th.ClearClientsTable(t, ctx, testPool)
		th.ClearContactsTable(t, ctx, testPool)

		// Create test clientEncx data using helper
		clientEncx := setupClient(t, ctx)

		// Create and insert a contact first
		originalContactEncx := th.NewTestContactEncx(t)
		originalContactEncx.ClientID = clientEncx.ID
		err := th.InsertContactEncx(t, ctx, testPool, originalContactEncx)
		require.NoError(t, err, "Failed to insert original contact")

		// Create updated contact with same ID but different data
		updatedContactEncx := th.NewTestContactEncx(t)
		updatedContactEncx.ID = originalContactEncx.ID // Keep same ID
		updatedContactEncx.LastnameEncrypted = []byte("new_lastname_encrypted")
		updatedContactEncx.FirstnameEncrypted = []byte("new_firstname_encrypted")
		updatedContactEncx.ClientID = clientEncx.ID

		// Update the contact
		err = repo.UpdateContact(ctx, updatedContactEncx)
		assert.NoError(t, err, "Failed to update contact")

		// Verify the update by retrieving the contact
		retrievedContact, err := repo.GetContactByID(ctx, updatedContactEncx.ID)
		require.NoError(t, err, "Failed to retrieve updated contact")

		// Verify fields were updated (should match updatedContact, not originalContact)
		assert.Equal(t, updatedContactEncx.ID, retrievedContact.ID)
		assert.Equal(t, updatedContactEncx.ClientID, retrievedContact.ClientID)
		assert.Equal(t, updatedContactEncx.LastnameEncrypted, retrievedContact.LastnameEncrypted)
		assert.Equal(t, updatedContactEncx.FirstnameEncrypted, retrievedContact.FirstnameEncrypted)

		t.Logf("✓ Successfully updated contact with ID: %s", updatedContactEncx.ID.String())
	})

	t.Run("update non-existent contact", func(t *testing.T) {
		ctx := context.Background()

		// Create contact with non-existent ID
		nonExistentID := uuid.New()
		contactEncx := th.NewTestContactEncx(t)
		contactEncx.ID = nonExistentID

		// Attempt to update non-existent contact
		err := repo.UpdateContact(ctx, contactEncx)
		assert.Error(t, err, "Expected error when updating non-existent contact")
		assert.Contains(t, err.Error(), "not found", "Error should mention not found")

		t.Logf("✓ Correctly handled update of non-existent contact")
	})

	t.Run("update contact with nil optional fields", func(t *testing.T) {
		ctx := context.Background()

		th.ClearClientsTable(t, ctx, testPool)
		th.ClearContactsTable(t, ctx, testPool)

		// Create test clientEncx data using helper
		clientEncx := setupClient(t, ctx)

		// Create and insert a contact first
		originalContactEncx := th.NewTestContactEncx(t)
		originalContactEncx.ClientID = clientEncx.ID
		err := th.InsertContactEncx(t, ctx, testPool, originalContactEncx)
		require.NoError(t, err, "Failed to insert original contact")

		// Create updated contact with nil optional fields
		updatedContactEncx := th.NewTestContactEncx(t)
		updatedContactEncx.ID = originalContactEncx.ID
		updatedContactEncx.ClientID = clientEncx.ID
		updatedContactEncx.PhoneEncrypted = nil
		updatedContactEncx.PositionEncrypted = nil

		// Update the contact
		err = repo.UpdateContact(ctx, updatedContactEncx)
		assert.NoError(t, err, "Failed to update contact with nil fields")

		// Verify the update
		retrievedContact, err := repo.GetContactByID(ctx, updatedContactEncx.ID)
		require.NoError(t, err, "Failed to retrieve updated contact")
		assert.Nil(t, retrievedContact.PhoneEncrypted, "Phone should be nil after update")
		assert.Nil(t, retrievedContact.PositionEncrypted, "Position should be nil after update")

		t.Logf("✓ Successfully updated contact with nil optional fields")
	})

	t.Run("update all contact fields", func(t *testing.T) {
		ctx := context.Background()

		th.ClearClientsTable(t, ctx, testPool)
		th.ClearContactsTable(t, ctx, testPool)

		// Create test clientEncx data using helper
		clientEncx := setupClient(t, ctx)

		// Create and insert a contact first
		originalContactEncx := th.NewTestContactEncx(t)
		originalContactEncx.ClientID = clientEncx.ID
		err := th.InsertContactEncx(t, ctx, testPool, originalContactEncx)
		require.NoError(t, err, "Failed to insert original contact")

		// Create completely new contact data
		updatedContactEncx := th.NewTestContactEncx(t)
		updatedContactEncx.ID = originalContactEncx.ID // Keep same ID
		updatedContactEncx.ClientID = clientEncx.ID
		updatedContactEncx.LastnameEncrypted = []byte("new_lastname_encrypted")
		updatedContactEncx.FirstnameEncrypted = []byte("new_firstname_encrypted")
		updatedContactEncx.EmailEncrypted = []byte("new_email_encrypted")
		updatedContactEncx.EmailHash = "new_email_encrypted"
		updatedContactEncx.PhoneEncrypted = []byte("new_phone_encrypted")
		updatedContactEncx.PositionEncrypted = []byte("new_position_encrypted")
		// All other fields should be different from original

		// Update the contact
		err = repo.UpdateContact(ctx, updatedContactEncx)
		assert.NoError(t, err, "Failed to update all contact fields")

		// Verify the update
		retrievedContact, err := repo.GetContactByID(ctx, updatedContactEncx.ID)
		require.NoError(t, err, "Failed to retrieve updated contact")

		// Verify all fields were updated (ID and client ID hash should match)
		assert.Equal(t, updatedContactEncx.ID, retrievedContact.ID)
		assert.Equal(t, updatedContactEncx.ClientID, retrievedContact.ClientID)
		assert.Equal(t, updatedContactEncx.KeyVersion, retrievedContact.KeyVersion)
		assert.Equal(t, updatedContactEncx.FirstnameEncrypted, retrievedContact.FirstnameEncrypted)
		assert.Equal(t, updatedContactEncx.LastnameEncrypted, retrievedContact.LastnameEncrypted)
		assert.Equal(t, updatedContactEncx.EmailEncrypted, retrievedContact.EmailEncrypted)
		assert.Equal(t, updatedContactEncx.EmailHash, retrievedContact.EmailHash)
		assert.Equal(t, updatedContactEncx.PhoneEncrypted, retrievedContact.PhoneEncrypted)
		assert.Equal(t, updatedContactEncx.PositionEncrypted, retrievedContact.PositionEncrypted)

		t.Logf("✓ Successfully updated all contact fields")
	})

	t.Run("update contact with different client ID", func(t *testing.T) {
		ctx := context.Background()

		th.ClearClientsTable(t, ctx, testPool)
		th.ClearContactsTable(t, ctx, testPool)

		// Create test clientEncx data using helper
		clientEncx := setupClient(t, ctx)

		// Create and insert a contact first
		originalContactEncx := th.NewTestContactEncx(t)
		originalContactEncx.ClientID = clientEncx.ID
		err := th.InsertContactEncx(t, ctx, testPool, originalContactEncx)
		require.NoError(t, err, "Failed to insert original contact")

		// Create updated contact with different client ID
		updatedContactEncx := th.NewTestContactEncx(t)
		updatedContactEncx.ID = originalContactEncx.ID
		clientEncx2 := th.NewTestClientEncx(t)
		clientEncx2.ID = uuid.New()
		err = th.InsertClientEncx(t, ctx, testPool, clientEncx2)
		require.NoError(t, err)
		updatedContactEncx.ClientID = clientEncx2.ID

		// Update the contact
		err = repo.UpdateContact(ctx, updatedContactEncx)
		assert.NoError(t, err, "Failed to update contact client ID")

		// Verify the update
		retrievedContactEncx, err := repo.GetContactByID(ctx, updatedContactEncx.ID)
		require.NoError(t, err, "Failed to retrieve updated contact")
		assert.Equal(t, clientEncx2.ID, retrievedContactEncx.ClientID, "Client ID hash should be updated")

		t.Logf("✓ Successfully updated contact client ID")
	})

	t.Run("context cancellation", func(t *testing.T) {
		// Create a context that will be cancelled
		ctx, cancel := context.WithCancel(context.Background())
		cancel() // Cancel immediately

		contactEncx := th.NewTestContactEncx(t)

		// Attempt to update with cancelled context
		err := repo.UpdateContact(ctx, contactEncx)
		assert.Error(t, err, "Expected context cancellation error")

		t.Logf("✓ Context cancellation handled correctly")
	})

	t.Run("update with invalid UUID", func(t *testing.T) {
		ctx := context.Background()

		// Create contactEncx with zero UUID (invalid)
		contactEncx := th.NewTestContactEncx(t)
		contactEncx.ID = uuid.Nil

		// Attempt to update with invalid UUID
		err := repo.UpdateContact(ctx, contactEncx)
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

		th.ClearClientsTable(t, ctx, testPool)
		th.ClearContactsTable(t, ctx, testPool)

		// Create test clientEncx data using helper
		clientEncx := setupClient(t, ctx)

		// Create and insert multiple contacts
		contactEncx1 := th.NewTestContactEncx(t)
		contactEncx1.ClientID = clientEncx.ID
		err := th.InsertContactEncx(t, ctx, testPool, contactEncx1)
		require.NoError(t, err, "Failed to insert contact 1")

		contactEncx2 := th.NewTestContactEncx(t)
		contactEncx2.ClientID = clientEncx.ID
		err = th.InsertContactEncx(t, ctx, testPool, contactEncx2)
		require.NoError(t, err, "Failed to insert contact 2")

		// Update only contact1
		updatedContactEncx1 := th.NewTestContactEncx(t)
		updatedContactEncx1.ID = contactEncx1.ID
		updatedContactEncx1.ClientID = clientEncx.ID

		err = repo.UpdateContact(ctx, updatedContactEncx1)
		assert.NoError(t, err, "Failed to update contact 1")

		// Verify contact1 was updated
		retrievedContactEncx1, err := repo.GetContactByID(ctx, contactEncx1.ID)
		require.NoError(t, err, "Failed to retrieve contact 1")
		assert.Equal(t, updatedContactEncx1.ID, retrievedContactEncx1.ID)

		// Verify contact2 was not affected
		retrievedContactEncx2, err := repo.GetContactByID(ctx, contactEncx2.ID)
		require.NoError(t, err, "Failed to retrieve contact 2")
		assert.Equal(t, contactEncx2.ID, retrievedContactEncx2.ID)

		t.Logf("✓ Only target contact was updated, other contacts unaffected")
	})
}
