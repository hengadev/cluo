package clientRepository_test

import (
	"context"
	"testing"

	"github.com/google/uuid"
	"github.com/hengadev/cluo_api/internal/domain/client"
	th "github.com/hengadev/cluo_api/test/helpers/client"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// make test-func TEST_NAME=TestGetAllContactsByClientID TEST_PATH=internal/infrastructure/postgres/client/get_all_contacts_by_client_id_test.go

func TestGetAllContactsByClientID(t *testing.T) {
	if testPool == nil || repo == nil {
		t.Skip("Test database or repository not initialized")
	}

	setupClient := func(t *testing.T, ctx context.Context) *client.ClientEncx {
		clientEncx := th.NewTestClientEncx(t)
		err := th.InsertClientEncx(t, ctx, testPool, clientEncx)
		require.NoError(t, err)
		return clientEncx
	}

	t.Run("successful retrieval with multiple contacts", func(t *testing.T) {
		ctx := context.Background()

		th.ClearClientsTable(t, ctx, testPool)
		th.ClearContactsTable(t, ctx, testPool)

		// Create test clientEncx data using helper
		clientEncx := setupClient(t, ctx)

		// Create multiple test contacts with the same client ID
		numContacts := 3

		for i := 0; i < numContacts; i++ {
			contactEncx := th.NewTestContactEncx(t)
			contactEncx.ClientID = clientEncx.ID

			// Insert contact directly into database
			err := th.InsertContactEncx(t, ctx, testPool, contactEncx)
			require.NoError(t, err, "Failed to insert contact %d", i)
		}

		// Test retrieval using the client ID
		retrievedContactsEncx, err := repo.GetAllContactsByClientID(ctx, clientEncx.ID)
		assert.NoError(t, err, "Failed to get contacts by client ID")
		assert.GreaterOrEqual(t, len(retrievedContactsEncx), numContacts, "Should retrieve at least the contacts we created")

		t.Logf("✓ Retrieved %d contacts for client", len(retrievedContactsEncx))
	})

	t.Run("successful retrieval with single contact", func(t *testing.T) {
		ctx := context.Background()

		th.ClearClientsTable(t, ctx, testPool)
		th.ClearContactsTable(t, ctx, testPool)

		// Create test clientEncx data using helper
		clientEncx := setupClient(t, ctx)

		// Create a single test contactEncx
		contactEncx := th.NewTestContactEncx(t)
		contactEncx.ClientID = clientEncx.ID

		// Insert contact
		err := th.InsertContactEncx(t, ctx, testPool, contactEncx)
		require.NoError(t, err, "Failed to insert contact")

		// Test retrieval using the client ID
		retrievedContactsEncx, err := repo.GetAllContactsByClientID(ctx, clientEncx.ID)
		assert.NoError(t, err, "Failed to get contacts by client ID")
		assert.GreaterOrEqual(t, len(retrievedContactsEncx), 1, "Should retrieve at least one contact")

		t.Logf("✓ Retrieved single contact for client")
	})

	t.Run("no contacts for client", func(t *testing.T) {
		ctx := context.Background()

		// Use a client ID hash that doesn't have any contacts
		clientID := uuid.New()

		// Test retrieval should return empty slice
		contactsEncx, err := repo.GetAllContactsByClientID(ctx, clientID)
		assert.NoError(t, err, "Should not error when no contacts exist")
		assert.Empty(t, contactsEncx, "Should return empty slice when no contacts exist")

		t.Logf("✓ Returned empty slice for client with no contacts")
	})

	t.Run("contacts with nil optional fields", func(t *testing.T) {
		ctx := context.Background()

		th.ClearClientsTable(t, ctx, testPool)
		th.ClearContactsTable(t, ctx, testPool)

		// Create test clientEncx data using helper
		clientEncx := setupClient(t, ctx)

		// Create contactEncx with nil optional fields
		contactEncx := th.NewTestContactEncx(t)
		contactEncx.ClientID = clientEncx.ID
		contactEncx.PhoneEncrypted = nil
		contactEncx.PositionEncrypted = nil

		// Insert contact
		err := th.InsertContactEncx(t, ctx, testPool, contactEncx)
		require.NoError(t, err, "Failed to insert contact with nil fields")

		// Test retrieval using the client ID
		retrievedContactsEncx, err := repo.GetAllContactsByClientID(ctx, clientEncx.ID)
		assert.NoError(t, err, "Failed to get contacts by client ID")
		assert.NotEmpty(t, retrievedContactsEncx, "Should retrieve contacts")

		t.Logf("✓ Retrieved contact with nil optional fields")
	})

	t.Run("contacts for different clients", func(t *testing.T) {
		ctx := context.Background()

		th.ClearClientsTable(t, ctx, testPool)
		th.ClearContactsTable(t, ctx, testPool)

		// Create test clientEncx data using helper
		clientEncx1 := setupClient(t, ctx)
		clientEncx2 := setupClient(t, ctx)

		// Create contact for client 1
		contactEncx1 := th.NewTestContactEncx(t)
		contactEncx1.ClientID = clientEncx1.ID
		err := th.InsertContactEncx(t, ctx, testPool, contactEncx1)
		require.NoError(t, err, "Failed to insert contact for client 1")

		// Create contacts for client 2
		contactEncx2a := th.NewTestContactEncx(t)
		contactEncx2a.ClientID = clientEncx2.ID
		err = th.InsertContactEncx(t, ctx, testPool, contactEncx2a)
		require.NoError(t, err, "Failed to insert contact 2a for client 2")

		contactEncx2b := th.NewTestContactEncx(t)
		contactEncx2b.ClientID = clientEncx2.ID
		err = th.InsertContactEncx(t, ctx, testPool, contactEncx2b)
		require.NoError(t, err, "Failed to insert contact 2b for client 2")

		// Test retrieval for client 1
		contacts1, err := repo.GetAllContactsByClientID(ctx, clientEncx1.ID)
		assert.NoError(t, err, "Failed to get contacts for client 1")

		// Test retrieval for client 2
		contacts2, err := repo.GetAllContactsByClientID(ctx, clientEncx2.ID)
		assert.NoError(t, err, "Failed to get contacts for client 2")

		t.Logf("✓ Retrieved %d contacts for client 1 and %d contacts for client 2", len(contacts1), len(contacts2))
	})

	t.Run("nil client ID", func(t *testing.T) {
		ctx := context.Background()

		// Test with nil UUID
		contactsEncx, err := repo.GetAllContactsByClientID(ctx, uuid.Nil)
		assert.NoError(t, err, "Should not error for nil client ID")
		assert.Empty(t, contactsEncx, "Should return empty slice for nil client ID")

		t.Logf("✓ Returned empty slice for nil client ID")
	})

	t.Run("context cancellation", func(t *testing.T) {
		// Create a context that will be cancelled
		ctx, cancel := context.WithCancel(context.Background())
		cancel() // Cancel immediately

		clientID := uuid.New()

		contactsEncx, err := repo.GetAllContactsByClientID(ctx, clientID)
		assert.Error(t, err, "Expected context cancellation error")
		assert.Nil(t, contactsEncx, "Should return nil contacts on context cancellation")

		t.Logf("✓ Context cancellation handled correctly")
	})
}
