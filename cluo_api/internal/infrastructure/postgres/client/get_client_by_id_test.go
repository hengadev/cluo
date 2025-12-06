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

// make test-func TEST_NAME=TestGetClientByID TEST_PATH=internal/infrastructure/postgres/client/get_client_by_id_test.go

func TestGetClientByID(t *testing.T) {
	if testPool == nil || repo == nil {
		t.Skip("Test database or repository not initialized")
	}

	t.Run("successful retrieval", func(t *testing.T) {
		ctx := context.Background()

		// Create test clientEncx data using helper
		clientEncx := th.NewTestClientEncx(t)

		// Insert client first (setup)
		err := th.InsertClientEncx(t, ctx, testPool, clientEncx)
		require.NoError(t, err, "Failed to create client for retrieval test")

		// Test successful client retrieval
		retrievedClientEncx, err := repo.GetClientByID(ctx, clientEncx.ID)
		assert.NoError(t, err, "Failed to retrieve client")
		assert.NotNil(t, retrievedClientEncx, "Retrieved client should not be nil")

		// Verify all field values match
		assert.Equal(t, clientEncx.ID, retrievedClientEncx.ID, "Client ID should match")
		assert.Equal(t, clientEncx.NameEncrypted, retrievedClientEncx.NameEncrypted, "NameEncrypted should match")
		assert.Equal(t, clientEncx.NameHash, retrievedClientEncx.NameHash, "NameHash should match")
		assert.Equal(t, clientEncx.TypeEncrypted, retrievedClientEncx.TypeEncrypted, "TypeEncrypted should match")
		assert.Equal(t, clientEncx.TypeHash, retrievedClientEncx.TypeHash, "TypeHash should match")
		assert.Equal(t, clientEncx.DEKEncrypted, retrievedClientEncx.DEKEncrypted, "DEKEncrypted should match")
		assert.Equal(t, clientEncx.KeyVersion, retrievedClientEncx.KeyVersion, "KeyVersion should match")
	})

	t.Run("non-existent client", func(t *testing.T) {
		ctx := context.Background()

		// Try to retrieve a client that doesn't exist
		nonExistentID := uuid.New()

		retrievedClientEncx, err := repo.GetClientByID(ctx, nonExistentID)
		assert.Error(t, err, "Expected error when retrieving non-existent client")
		assert.Nil(t, retrievedClientEncx, "Retrieved client should be nil for non-existent client")
	})

	t.Run("nil UUID", func(t *testing.T) {
		ctx := context.Background()

		// Try to retrieve with nil UUID
		retrievedClientEncx, err := repo.GetClientByID(ctx, uuid.Nil)
		assert.Error(t, err, "Expected error when retrieving client with nil UUID")
		assert.Nil(t, retrievedClientEncx, "Retrieved client should be nil for nil UUID")
	})

	t.Run("client with empty encrypted fields", func(t *testing.T) {
		ctx := context.Background()

		// Create clientEncx with empty encrypted fields
		clientEncx := th.NewTestClientEncx(t)
		clientEncx.NameEncrypted = []byte("")
		clientEncx.TypeEncrypted = []byte("")

		// Insert client
		err := th.InsertClientEncx(t, ctx, testPool, clientEncx)
		require.NoError(t, err, "Failed to create client with empty encrypted fields")

		// Retrieve client
		retrievedClientEncx, err := repo.GetClientByID(ctx, clientEncx.ID)
		assert.NoError(t, err, "Failed to retrieve client with empty encrypted fields")
		assert.NotNil(t, retrievedClientEncx, "Retrieved client should not be nil")

		// Verify empty fields are preserved
		assert.Equal(t, []byte(""), retrievedClientEncx.NameEncrypted, "NameEncrypted should be empty")
		assert.Equal(t, []byte(""), retrievedClientEncx.TypeEncrypted, "TypeEncrypted should be empty")
	})

	t.Run("client with large encrypted fields", func(t *testing.T) {
		ctx := context.Background()

		// Create clientEncx with large encrypted fields
		clientEncx := th.NewTestClientEncx(t)
		largeData := make([]byte, 10000) // 10KB of data
		for i := range largeData {
			largeData[i] = byte(i % 256)
		}

		clientEncx.NameEncrypted = largeData

		// Insert client
		err := th.InsertClientEncx(t, ctx, testPool, clientEncx)
		require.NoError(t, err, "Failed to create client with large encrypted fields")

		// Retrieve client
		retrievedClientEncx, err := repo.GetClientByID(ctx, clientEncx.ID)
		assert.NoError(t, err, "Failed to retrieve client with large encrypted fields")
		assert.NotNil(t, retrievedClientEncx, "Retrieved client should not be nil")

		// Verify large fields are preserved
		assert.Equal(t, len(largeData), len(retrievedClientEncx.NameEncrypted), "NameEncrypted length should match")
		assert.Equal(t, largeData, retrievedClientEncx.NameEncrypted, "NameEncrypted content should match")
	})

	t.Run("context cancellation", func(t *testing.T) {
		// Create a context that will be cancelled
		ctx, cancel := context.WithCancel(context.Background())
		cancel() // Cancel immediately

		clientID := uuid.New()

		retrievedClientEncx, err := repo.GetClientByID(ctx, clientID)
		assert.Error(t, err, "Expected context cancellation error, but got nil")
		assert.Nil(t, retrievedClientEncx, "Retrieved client should be nil for cancelled context")
	})

	t.Run("multiple clients retrieval", func(t *testing.T) {
		ctx := context.Background()

		// Create multiple test clientsEncx
		clientsEncx := make([]*client.ClientEncx, 5)
		clientIDs := make([]uuid.UUID, 5)

		for i := 0; i < 5; i++ {
			clientsEncx[i] = th.NewTestClientEncx(t)
			clientIDs[i] = clientsEncx[i].ID
			err := th.InsertClientEncx(t, ctx, testPool, clientsEncx[i])
			require.NoError(t, err, "Failed to create client %d", i)
		}

		// Retrieve each client and verify
		for i, clientID := range clientIDs {
			retrievedClient, err := repo.GetClientByID(ctx, clientID)
			assert.NoError(t, err, "Failed to retrieve client %d", i)
			assert.NotNil(t, retrievedClient, "Retrieved client %d should not be nil", i)
			assert.Equal(t, clientID, retrievedClient.ID, "Client %d ID should match", i)
		}
	})
}
