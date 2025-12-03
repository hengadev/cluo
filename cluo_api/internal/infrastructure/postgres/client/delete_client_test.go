package clientRepository_test

import (
	"context"
	"testing"

	th "github.com/hengadev/cluo_api/test/helpers/client"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestDeleteClient(t *testing.T) {
	if testPool == nil || repo == nil {
		t.Skip("Test database or repository not initialized")
	}

	t.Run("successful deletion", func(t *testing.T) {
		ctx := context.Background()

		// Create test client data using helper
		client := th.NewTestClientEncx(t)

		// Insert client first (setup)
		err := repo.CreateClient(ctx, client)
		require.NoError(t, err, "Failed to create client for deletion test")

		// Verify client exists before deletion
		retrievedClient, err := th.GetClientEncxByID(t, ctx, testPool, client.ID)
		assert.NoError(t, err, "Failed to retrieve client before deletion")
		assert.NotNil(t, retrievedClient, "Client should exist before deletion")

		// Test successful client deletion
		err = repo.DeleteClient(ctx, client.ID)
		assert.NoError(t, err, "Failed to delete client")

		// Verify client no longer exists
		retrievedClient, err = th.GetClientEncxByID(t, ctx, testPool, client.ID)
		assert.Error(t, err, "Expected error when retrieving deleted client")
		assert.Nil(t, retrievedClient, "Retrieved client should be nil after deletion")
	})

	t.Run("non-existent client", func(t *testing.T) {
		ctx := context.Background()

		// Try to delete a client that doesn't exist
		nonExistentID := uuid.New()

		err := repo.DeleteClient(ctx, nonExistentID)
		assert.Error(t, err, "Expected error when deleting non-existent client")
	})

	t.Run("nil UUID", func(t *testing.T) {
		ctx := context.Background()

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

		// Create multiple test clients
		clients := make([]*client.ClientEncx, 3)
		for i := 0; i < 3; i++ {
			clients[i] = th.NewTestClientEncx(t)
			err := repo.CreateClient(ctx, clients[i])
			require.NoError(t, err, "Failed to create client %d", i)
		}

		// Delete all clients
		for i, client := range clients {
			err := repo.DeleteClient(ctx, client.ID)
			assert.NoError(t, err, "Failed to delete client %d", i)

			// Verify deletion
			_, err = th.GetClientEncxByID(t, ctx, testPool, client.ID)
			assert.Error(t, err, "Client %d should no longer exist after deletion", i)
		}
	})
}