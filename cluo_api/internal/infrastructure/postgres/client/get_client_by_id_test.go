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

		// Create test client data using helper
		client := th.NewTestClientEncx(t)

		// Insert client first (setup)
		err := repo.CreateClient(ctx, client)
		require.NoError(t, err, "Failed to create client for retrieval test")

		// Test successful client retrieval
		retrievedClient, err := repo.GetClientByID(ctx, client.ID)
		assert.NoError(t, err, "Failed to retrieve client")
		assert.NotNil(t, retrievedClient, "Retrieved client should not be nil")

		// Verify all field values match
		assert.Equal(t, client.ID, retrievedClient.ID, "Client ID should match")
		assert.Equal(t, client.NameEncrypted, retrievedClient.NameEncrypted, "NameEncrypted should match")
		assert.Equal(t, client.NameHash, retrievedClient.NameHash, "NameHash should match")
		assert.Equal(t, client.TypeEncrypted, retrievedClient.TypeEncrypted, "TypeEncrypted should match")
		assert.Equal(t, client.TypeHash, retrievedClient.TypeHash, "TypeHash should match")
		assert.Equal(t, client.ContactIDsEncrypted, retrievedClient.ContactIDsEncrypted, "ContactIDsEncrypted should match")
		assert.Equal(t, client.DEKEncrypted, retrievedClient.DEKEncrypted, "DEKEncrypted should match")
		assert.Equal(t, client.KeyVersion, retrievedClient.KeyVersion, "KeyVersion should match")
	})

	t.Run("non-existent client", func(t *testing.T) {
		ctx := context.Background()

		// Try to retrieve a client that doesn't exist
		nonExistentID := uuid.New()

		retrievedClient, err := repo.GetClientByID(ctx, nonExistentID)
		assert.Error(t, err, "Expected error when retrieving non-existent client")
		assert.Nil(t, retrievedClient, "Retrieved client should be nil for non-existent client")
	})

	t.Run("nil UUID", func(t *testing.T) {
		ctx := context.Background()

		// Try to retrieve with nil UUID
		retrievedClient, err := repo.GetClientByID(ctx, uuid.Nil)
		assert.Error(t, err, "Expected error when retrieving client with nil UUID")
		assert.Nil(t, retrievedClient, "Retrieved client should be nil for nil UUID")
	})

	t.Run("client with empty encrypted fields", func(t *testing.T) {
		ctx := context.Background()

		// Create client with empty encrypted fields
		client := th.NewTestClientEncx(t)
		client.NameEncrypted = []byte("")
		client.TypeEncrypted = []byte("")
		client.ContactIDsEncrypted = []byte("")

		// Insert client
		err := repo.CreateClient(ctx, client)
		require.NoError(t, err, "Failed to create client with empty encrypted fields")

		// Retrieve client
		retrievedClient, err := repo.GetClientByID(ctx, client.ID)
		assert.NoError(t, err, "Failed to retrieve client with empty encrypted fields")
		assert.NotNil(t, retrievedClient, "Retrieved client should not be nil")

		// Verify empty fields are preserved
		assert.Equal(t, []byte(""), retrievedClient.NameEncrypted, "NameEncrypted should be empty")
		assert.Equal(t, []byte(""), retrievedClient.TypeEncrypted, "TypeEncrypted should be empty")
		assert.Equal(t, []byte(""), retrievedClient.ContactIDsEncrypted, "ContactIDsEncrypted should be empty")
	})

	t.Run("client with large encrypted fields", func(t *testing.T) {
		ctx := context.Background()

		// Create client with large encrypted fields
		client := th.NewTestClientEncx(t)
		largeData := make([]byte, 10000) // 10KB of data
		for i := range largeData {
			largeData[i] = byte(i % 256)
		}

		client.NameEncrypted = largeData
		client.ContactIDsEncrypted = largeData

		// Insert client
		err := repo.CreateClient(ctx, client)
		require.NoError(t, err, "Failed to create client with large encrypted fields")

		// Retrieve client
		retrievedClient, err := repo.GetClientByID(ctx, client.ID)
		assert.NoError(t, err, "Failed to retrieve client with large encrypted fields")
		assert.NotNil(t, retrievedClient, "Retrieved client should not be nil")

		// Verify large fields are preserved
		assert.Equal(t, len(largeData), len(retrievedClient.NameEncrypted), "NameEncrypted length should match")
		assert.Equal(t, len(largeData), len(retrievedClient.ContactIDsEncrypted), "ContactIDsEncrypted length should match")
		assert.Equal(t, largeData, retrievedClient.NameEncrypted, "NameEncrypted content should match")
		assert.Equal(t, largeData, retrievedClient.ContactIDsEncrypted, "ContactIDsEncrypted content should match")
	})

	t.Run("context cancellation", func(t *testing.T) {
		// Create a context that will be cancelled
		ctx, cancel := context.WithCancel(context.Background())
		cancel() // Cancel immediately

		clientID := uuid.New()

		retrievedClient, err := repo.GetClientByID(ctx, clientID)
		assert.Error(t, err, "Expected context cancellation error, but got nil")
		assert.Nil(t, retrievedClient, "Retrieved client should be nil for cancelled context")
	})

	t.Run("multiple clients retrieval", func(t *testing.T) {
		ctx := context.Background()

		// Create multiple test clients
		clients := make([]*client.ClientEncx, 5)
		clientIDs := make([]uuid.UUID, 5)

		for i := 0; i < 5; i++ {
			clients[i] = th.NewTestClientEncx(t)
			clientIDs[i] = clients[i].ID
			err := repo.CreateClient(ctx, clients[i])
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
