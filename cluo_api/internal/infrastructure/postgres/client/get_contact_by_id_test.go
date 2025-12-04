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

// make test-func TEST_NAME=TestGetContactByID TEST_PATH=internal/infrastructure/postgres/client/get_contact_by_id_test.go

func TestGetContactByID(t *testing.T) {
	if testPool == nil || repo == nil {
		t.Skip("Test database or repository not initialized")
	}

	setupClient := func(t *testing.T, ctx context.Context) *client.ClientEncx {
		clientEncx := th.NewTestClientEncx(t)
		err := th.InsertClientEncx(t, ctx, testPool, clientEncx)
		require.NoError(t, err)
		return clientEncx
	}

	t.Run("successful retrieval", func(t *testing.T) {
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

		// Test successful contact retrieval using the global repo
		retrievedContactEncx, err := repo.GetContactByID(ctx, contactEncx.ID)
		assert.NoError(t, err, "Failed to get contact by ID")
		require.NotNil(t, retrievedContactEncx, "Retrieved contact should not be nil")

		// Verify field values
		assert.Equal(t, contactEncx.ID, retrievedContactEncx.ID, "Contact ID should match")
		assert.Equal(t, contactEncx.ClientID, retrievedContactEncx.ClientID, "Client hash should match")
		assert.Equal(t, contactEncx.EmailHash, retrievedContactEncx.EmailHash, "Email hash should match")
		assert.Equal(t, contactEncx.KeyVersion, retrievedContactEncx.KeyVersion, "Key version should match")
	})

	t.Run("non-existent ID", func(t *testing.T) {
		ctx := context.Background()

		// Use a random UUID that doesn't exist in database
		nonExistentID := uuid.New()

		// Test that non-existent ID returns an error
		retrievedContactEncx, err := repo.GetContactByID(ctx, nonExistentID)
		assert.Error(t, err, "Expected error when getting non-existent contact")
		assert.Nil(t, retrievedContactEncx, "Retrieved contact should be nil for non-existent ID")
	})

	t.Run("nil UUID", func(t *testing.T) {
		ctx := context.Background()

		// Test with nil UUID
		retrievedContactEncx, err := repo.GetContactByID(ctx, uuid.Nil)
		assert.Error(t, err, "Expected error when getting contact with nil UUID")
		assert.Nil(t, retrievedContactEncx, "Retrieved contact should be nil for nil UUID")
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

		// Insert contact
		err := th.InsertContactEncx(t, ctx, testPool, contactEncx)
		require.NoError(t, err, "Failed to create contact with nil optional fields")

		// Retrieve the contact
		retrievedContactEncx, err := repo.GetContactByID(ctx, contactEncx.ID)
		assert.NoError(t, err, "Failed to get contact with nil optional fields")
		require.NotNil(t, retrievedContactEncx, "Retrieved contact should not be nil")

		// Verify that optional fields are indeed nil
		assert.Nil(t, retrievedContactEncx.PhoneEncrypted, "Expected PhoneEncrypted to be nil")
		assert.Nil(t, retrievedContactEncx.PositionEncrypted, "Expected PositionEncrypted to be nil")
	})

	t.Run("context cancellation", func(t *testing.T) {
		// Create a context that will be cancelled
		ctx, cancel := context.WithCancel(context.Background())
		cancel() // Cancel immediately

		contactID := uuid.New()

		retrievedContactEncx, err := repo.GetContactByID(ctx, contactID)
		assert.Error(t, err, "Expected context cancellation error, but got nil")
		assert.Nil(t, retrievedContactEncx, "Retrieved contact should be nil due to context cancellation")
	})
}
