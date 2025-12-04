package clientRepository_test

import (
	"context"
	"fmt"
	"testing"

	"github.com/hengadev/cluo_api/internal/domain/client"
	th "github.com/hengadev/cluo_api/test/helpers/client"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// make test-func TEST_NAME=TestDeleteClient TEST_PATH=internal/infrastructure/postgres/client/delete_client_test.go

func TestDeleteClient(t *testing.T) {
	if testPool == nil || repo == nil {
		t.Skip("Test database or repository not initialized")
	}

	t.Run("successful deletion", func(t *testing.T) {
		ctx := context.Background()

		th.ClearClientsTable(t, ctx, testPool)

		// Create test clientEncx data using helper
		clientEncx := th.NewTestClientEncx(t)
		err := th.InsertClientEncx(t, ctx, testPool, clientEncx)
		require.NoError(t, err)

		// Verify client exists before deletion
		retrievedClientEncx, err := th.GetClientEncxByID(t, ctx, testPool, clientEncx.ID)
		assert.NoError(t, err, "Failed to retrieve client before deletion")
		assert.NotNil(t, retrievedClientEncx, "Client should exist before deletion")

		// Test successful client deletion
		err = repo.DeleteClient(ctx, clientEncx.ID)
		assert.NoError(t, err, "Failed to delete client")

		// Verify client no longer exists
		retrievedClientEncx, err = th.GetClientEncxByID(t, ctx, testPool, clientEncx.ID)
		assert.Error(t, err, "Expected error when retrieving deleted client")
		assert.Equal(t, &client.ClientEncx{}, retrievedClientEncx, "Retrieved client should be nil after deletion")
	})

	t.Run("non-existent client", func(t *testing.T) {
		ctx := context.Background()

		th.ClearClientsTable(t, ctx, testPool)

		// Try to delete a client that doesn't exist
		nonExistentID := uuid.New()

		err := repo.DeleteClient(ctx, nonExistentID)
		assert.Error(t, err, "Expected error when deleting non-existent client")
	})

	t.Run("nil UUID", func(t *testing.T) {
		ctx := context.Background()

		th.ClearClientsTable(t, ctx, testPool)

		// Try to delete with nil UUID
		err := repo.DeleteClient(ctx, uuid.Nil)
		assert.Error(t, err, "Expected error when deleting client with nil UUID")
	})

	t.Run("context cancellation", func(t *testing.T) {
		// Create a context that will be cancelled
		ctx, cancel := context.WithCancel(context.Background())
		cancel() // Cancel immediately

		clientID := uuid.New()

		err := repo.DeleteClient(ctx, clientID)
		assert.Error(t, err, "Expected context cancellation error, but got nil")
	})

	t.Run("delete multiple clients", func(t *testing.T) {
		ctx := context.Background()

		th.ClearClientsTable(t, ctx, testPool)

		// Create multiple test clientsEncx
		clientsEncx := make([]*client.ClientEncx, 3)
		for i := 0; i < 3; i++ {
			clientEncx := th.NewTestClientEncx(t)
			clientEncx.NameEncrypted = []byte(fmt.Sprintf("name_encrypted_%d", i))
			clientEncx.NameHash = fmt.Sprintf("name_hash_%d", i)
			clientsEncx[i] = clientEncx
			err := th.InsertClientEncx(t, ctx, testPool, clientEncx)
			require.NoError(t, err, "Failed to create client %d", i)
		}

		// Delete all clients
		for i, client := range clientsEncx {
			err := repo.DeleteClient(ctx, client.ID)
			assert.NoError(t, err, "Failed to delete client %d", i)

			// Verify deletion
			_, err = th.GetClientEncxByID(t, ctx, testPool, client.ID)
			assert.Error(t, err, "Client %d should no longer exist after deletion", i)
		}
	})
}
