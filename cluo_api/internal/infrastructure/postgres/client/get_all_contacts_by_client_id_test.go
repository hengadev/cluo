package clientRepository_test

import (
	"context"
	"testing"

	th "github.com/hengadev/cluo_api/test/helpers"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// make test-func TEST_NAME=TestGetAllContactsByClientID TEST_PATH=internal/infrastructure/postgres/client/get_all_contacts_by_client_id_test.go

func TestGetAllContactsByClientID(t *testing.T) {
	if testPool == nil || repo == nil {
		t.Skip("Test database or repository not initialized")
	}

	t.Run("successful retrieval with multiple contacts", func(t *testing.T) {
		ctx := context.Background()

		// Create multiple test contacts with the same client ID hash
		numContacts := 3
		clientIDHash := "test-client-hash-123"

		for i := 0; i < numContacts; i++ {
			contact := th.NewTestContactEncx(t)
			contact.ClientIDHash = clientIDHash

			// Insert contact directly into database
			err := th.InsertContactEncx(t, ctx, testPool, *contact)
			require.NoError(t, err, "Failed to insert contact %d", i)
		}

		// Test retrieval using the client ID hash string
		retrievedContacts, err := repo.GetAllContactsByClientID(ctx, clientIDHash)
		assert.NoError(t, err, "Failed to get contacts by client ID")
		assert.GreaterOrEqual(t, len(retrievedContacts), numContacts, "Should retrieve at least the contacts we created")

		t.Logf("✓ Retrieved %d contacts for client", len(retrievedContacts))
	})

	t.Run("successful retrieval with single contact", func(t *testing.T) {
		ctx := context.Background()

		// Create a single test contact
		contact := th.NewTestContactEncx(t)
		clientIDHash := "test-client-hash-456"
		contact.ClientIDHash = clientIDHash

		// Insert contact
		err := th.InsertContactEncx(t, ctx, testPool, *contact)
		require.NoError(t, err, "Failed to insert contact")

		// Test retrieval using the client ID hash string
		retrievedContacts, err := repo.GetAllContactsByClientID(ctx, clientIDHash)
		assert.NoError(t, err, "Failed to get contacts by client ID")
		assert.GreaterOrEqual(t, len(retrievedContacts), 1, "Should retrieve at least one contact")

		t.Logf("✓ Retrieved single contact for client")
	})

	t.Run("no contacts for client", func(t *testing.T) {
		ctx := context.Background()

		// Use a client ID hash that doesn't have any contacts
		clientIDHash := "test-client-hash-789"

		// Test retrieval should return empty slice
		contacts, err := repo.GetAllContactsByClientID(ctx, clientIDHash)
		assert.NoError(t, err, "Should not error when no contacts exist")
		assert.Empty(t, contacts, "Should return empty slice when no contacts exist")

		t.Logf("✓ Returned empty slice for client with no contacts")
	})

	t.Run("contacts with nil optional fields", func(t *testing.T) {
		ctx := context.Background()

		// Create contact with nil optional fields
		contact := th.NewTestContactEncx(t)
		clientIDHash := "test-client-hash-nil-fields"
		contact.ClientIDHash = clientIDHash
		contact.PhoneEncrypted = nil
		contact.PositionEncrypted = nil

		// Insert contact
		err := th.InsertContactEncx(t, ctx, testPool, *contact)
		require.NoError(t, err, "Failed to insert contact with nil fields")

		// Test retrieval using the client ID hash string
		retrievedContacts, err := repo.GetAllContactsByClientID(ctx, clientIDHash)
		assert.NoError(t, err, "Failed to get contacts by client ID")
		assert.NotEmpty(t, retrievedContacts, "Should retrieve contacts")

		t.Logf("✓ Retrieved contact with nil optional fields")
	})

	t.Run("contacts for different clients", func(t *testing.T) {
		ctx := context.Background()

		// Create contacts for different clients
		clientIDHash1 := "test-client-hash-diff-1"
		clientIDHash2 := "test-client-hash-diff-2"

		// Create contact for client 1
		contact1 := th.NewTestContactEncx(t)
		contact1.ClientIDHash = clientIDHash1
		err := th.InsertContactEncx(t, ctx, testPool, *contact1)
		require.NoError(t, err, "Failed to insert contact for client 1")

		// Create contacts for client 2
		contact2a := th.NewTestContactEncx(t)
		contact2a.ClientIDHash = clientIDHash2
		err = th.InsertContactEncx(t, ctx, testPool, *contact2a)
		require.NoError(t, err, "Failed to insert contact 2a for client 2")

		contact2b := th.NewTestContactEncx(t)
		contact2b.ClientIDHash = clientIDHash2
		err = th.InsertContactEncx(t, ctx, testPool, *contact2b)
		require.NoError(t, err, "Failed to insert contact 2b for client 2")

		// Test retrieval for client 1
		contacts1, err := repo.GetAllContactsByClientID(ctx, clientIDHash1)
		assert.NoError(t, err, "Failed to get contacts for client 1")

		// Test retrieval for client 2
		contacts2, err := repo.GetAllContactsByClientID(ctx, clientIDHash2)
		assert.NoError(t, err, "Failed to get contacts for client 2")

		t.Logf("✓ Retrieved %d contacts for client 1 and %d contacts for client 2", len(contacts1), len(contacts2))
	})

	t.Run("nil client ID", func(t *testing.T) {
		ctx := context.Background()

		// Test with nil UUID
		contacts, err := repo.GetAllContactsByClientID(ctx, "")
		assert.NoError(t, err, "Should not error for nil client ID")
		assert.Empty(t, contacts, "Should return empty slice for nil client ID")

		t.Logf("✓ Returned empty slice for nil client ID")
	})

	t.Run("context cancellation", func(t *testing.T) {
		// Create a context that will be cancelled
		ctx, cancel := context.WithCancel(context.Background())
		cancel() // Cancel immediately

		clientIDHash := "test-client-hash-cancel"

		contacts, err := repo.GetAllContactsByClientID(ctx, clientIDHash)
		assert.Error(t, err, "Expected context cancellation error")
		assert.Nil(t, contacts, "Should return nil contacts on context cancellation")

		t.Logf("✓ Context cancellation handled correctly")
	})
}
