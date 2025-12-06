package clientRepository_test

import (
	"context"
	"testing"

	th "github.com/hengadev/cluo_api/test/helpers/client"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// make test-func TEST_NAME=TestGetContactIDsForClient TEST_PATH=internal/infrastructure/postgres/client/get_contact_ids_for_client_test.go

func TestGetContactIDsForClient(t *testing.T) {
	if testPool == nil || repo == nil {
		t.Skip("Test database or repository not initialized")
	}

	t.Run("client with no contacts", func(t *testing.T) {
		ctx := context.Background()

		th.ClearClientsTable(t, ctx, testPool)
		th.ClearContactsTable(t, ctx, testPool)

		// Create test client
		clientEncx := th.NewTestClientEncx(t)
		err := th.InsertClientEncx(t, ctx, testPool, clientEncx)
		require.NoError(t, err, "Failed to create test client")

		// Test getting contact IDs for client with no contacts
		contactIDs, err := repo.GetContactIDsForClient(ctx, clientEncx.ID)
		assert.NoError(t, err, "Failed to get contact IDs for client with no contacts")
		assert.Empty(t, contactIDs, "Expected no contact IDs for client with no contacts")
	})

	t.Run("client with single contact", func(t *testing.T) {
		ctx := context.Background()

		th.ClearClientsTable(t, ctx, testPool)
		th.ClearContactsTable(t, ctx, testPool)

		// Create test client
		clientEncx := th.NewTestClientEncx(t)
		err := th.InsertClientEncx(t, ctx, testPool, clientEncx)
		require.NoError(t, err, "Failed to create test client")

		// Create test contact
		contactEncx := th.NewTestContactEncxWithClientID(t, clientEncx.ID)
		err = th.InsertContactEncx(t, ctx, testPool, contactEncx)
		require.NoError(t, err, "Failed to create test contact")

		// Test getting contact IDs
		contactIDs, err := repo.GetContactIDsForClient(ctx, clientEncx.ID)
		assert.NoError(t, err, "Failed to get contact IDs for client with single contact")
		assert.Len(t, contactIDs, 1, "Expected exactly 1 contact ID")
		assert.Equal(t, contactEncx.ID, contactIDs[0], "Expected contact ID to match")
	})

	t.Run("client with multiple contacts", func(t *testing.T) {
		ctx := context.Background()

		th.ClearClientsTable(t, ctx, testPool)
		th.ClearContactsTable(t, ctx, testPool)

		// Create test client
		clientEncx := th.NewTestClientEncx(t)
		err := th.InsertClientEncx(t, ctx, testPool, clientEncx)
		require.NoError(t, err, "Failed to create test client")

		// Create multiple test contacts
		expectedContactIDs := make([]uuid.UUID, 3)
		for i := 0; i < 3; i++ {
			contactEncx := th.NewTestContactEncxWithClientID(t, clientEncx.ID)
			err = th.InsertContactEncx(t, ctx, testPool, contactEncx)
			require.NoError(t, err, "Failed to create test contact %d", i)
			expectedContactIDs[i] = contactEncx.ID
		}

		// Test getting contact IDs
		contactIDs, err := repo.GetContactIDsForClient(ctx, clientEncx.ID)
		assert.NoError(t, err, "Failed to get contact IDs for client with multiple contacts")
		assert.Len(t, contactIDs, 3, "Expected exactly 3 contact IDs")

		// Verify all expected contact IDs are present (order should be by created_at)
		for i, expectedID := range expectedContactIDs {
			assert.Equal(t, expectedID, contactIDs[i], "Contact ID at index %d should match", i)
		}
	})

	t.Run("client does not exist", func(t *testing.T) {
		ctx := context.Background()

		th.ClearClientsTable(t, ctx, testPool)
		th.ClearContactsTable(t, ctx, testPool)

		// Use a non-existent client ID
		nonExistentClientID := uuid.New()

		// Test getting contact IDs for non-existent client
		contactIDs, err := repo.GetContactIDsForClient(ctx, nonExistentClientID)
		assert.NoError(t, err, "Should not return error for non-existent client")
		assert.Empty(t, contactIDs, "Expected no contact IDs for non-existent client")
	})

	t.Run("contacts ordered by created_at", func(t *testing.T) {
		ctx := context.Background()

		th.ClearClientsTable(t, ctx, testPool)
		th.ClearContactsTable(t, ctx, testPool)

		// Create test client
		clientEncx := th.NewTestClientEncx(t)
		err := th.InsertClientEncx(t, ctx, testPool, clientEncx)
		require.NoError(t, err, "Failed to create test client")

		// Create contacts with specific timestamps to test ordering
		contactIDsInOrder := make([]uuid.UUID, 3)
		for i := 0; i < 3; i++ {
			contactEncx := th.NewTestContactEncxWithTimestamp(t, clientEncx.ID, i)
			err = th.InsertContactEncx(t, ctx, testPool, contactEncx)
			require.NoError(t, err, "Failed to create test contact %d", i)
			contactIDsInOrder[i] = contactEncx.ID
		}

		// Test getting contact IDs
		contactIDs, err := repo.GetContactIDsForClient(ctx, clientEncx.ID)
		assert.NoError(t, err, "Failed to get contact IDs")
		assert.Len(t, contactIDs, 3, "Expected exactly 3 contact IDs")

		// Verify contacts are ordered by created_at
		for i, expectedID := range contactIDsInOrder {
			assert.Equal(t, expectedID, contactIDs[i], "Contact ID at index %d should be in chronological order", i)
		}
	})

	t.Run("contacts for different clients are separated", func(t *testing.T) {
		ctx := context.Background()

		th.ClearClientsTable(t, ctx, testPool)
		th.ClearContactsTable(t, ctx, testPool)

		// Create two test clients
		client1Encx := th.NewTestClientEncx(t)
		err := th.InsertClientEncx(t, ctx, testPool, client1Encx)
		require.NoError(t, err, "Failed to create first test client")

		client2Encx := th.NewTestClientEncx(t)
		err = th.InsertClientEncx(t, ctx, testPool, client2Encx)
		require.NoError(t, err, "Failed to create second test client")

		// Create contacts for each client
		client1Contact := th.NewTestContactEncxWithClientID(t, client1Encx.ID)
		err = th.InsertContactEncx(t, ctx, testPool, client1Contact)
		require.NoError(t, err, "Failed to create contact for first client")

		client2Contact := th.NewTestContactEncxWithClientID(t, client2Encx.ID)
		err = th.InsertContactEncx(t, ctx, testPool, client2Contact)
		require.NoError(t, err, "Failed to create contact for second client")

		// Test getting contact IDs for each client
		client1ContactIDs, err := repo.GetContactIDsForClient(ctx, client1Encx.ID)
		assert.NoError(t, err, "Failed to get contact IDs for first client")
		assert.Len(t, client1ContactIDs, 1, "Expected exactly 1 contact ID for first client")
		assert.Equal(t, client1Contact.ID, client1ContactIDs[0], "Expected first client contact ID to match")

		client2ContactIDs, err := repo.GetContactIDsForClient(ctx, client2Encx.ID)
		assert.NoError(t, err, "Failed to get contact IDs for second client")
		assert.Len(t, client2ContactIDs, 1, "Expected exactly 1 contact ID for second client")
		assert.Equal(t, client2Contact.ID, client2ContactIDs[0], "Expected second client contact ID to match")
	})

	t.Run("context cancellation", func(t *testing.T) {
		// Create a context that will be cancelled
		ctx, cancel := context.WithCancel(context.Background())
		cancel() // Cancel immediately

		clientID := uuid.New()

		// Test with cancelled context
		contactIDs, err := repo.GetContactIDsForClient(ctx, clientID)
		assert.Error(t, err, "Expected error for cancelled context")
		assert.Nil(t, contactIDs, "Expected nil contact IDs for cancelled context")
	})

	t.Run("nil client ID", func(t *testing.T) {
		ctx := context.Background()

		// Test with nil UUID
		contactIDs, err := repo.GetContactIDsForClient(ctx, uuid.Nil)
		assert.NoError(t, err, "Should not return error for nil client ID")
		assert.Empty(t, contactIDs, "Expected no contact IDs for nil client ID")
	})

	t.Run("large number of contacts", func(t *testing.T) {
		ctx := context.Background()

		th.ClearClientsTable(t, ctx, testPool)
		th.ClearContactsTable(t, ctx, testPool)

		// Create test client
		clientEncx := th.NewTestClientEncx(t)
		err := th.InsertClientEncx(t, ctx, testPool, clientEncx)
		require.NoError(t, err, "Failed to create test client")

		// Create 50 contacts
		const numContacts = 50
		for i := 0; i < numContacts; i++ {
			contactEncx := th.NewTestContactEncxWithClientID(t, clientEncx.ID)
			err = th.InsertContactEncx(t, ctx, testPool, contactEncx)
			require.NoError(t, err, "Failed to create test contact %d", i)
		}

		// Test getting all contact IDs
		contactIDs, err := repo.GetContactIDsForClient(ctx, clientEncx.ID)
		assert.NoError(t, err, "Failed to get contact IDs for client with many contacts")
		assert.Len(t, contactIDs, numContacts, "Expected exactly %d contact IDs", numContacts)
	})
}

