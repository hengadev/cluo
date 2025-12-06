package clientRepository_test

import (
	"context"
	"testing"
	"time"

	"github.com/hengadev/cluo_api/internal/domain/client"
	th "github.com/hengadev/cluo_api/test/helpers/client"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// make test-func TEST_NAME=TestGetAllClients TEST_PATH=internal/infrastructure/postgres/client/get_all_clients_test.go

func TestGetAllClients(t *testing.T) {
	if testPool == nil || repo == nil {
		t.Skip("Test database or repository not initialized")
	}

	t.Run("empty database", func(t *testing.T) {
		ctx := context.Background()

		th.ClearClientsTable(t, ctx, testPool)

		// Test retrieving all clientsEncx when database is empty
		clientsEncx, err := repo.GetAllClients(ctx)
		assert.NoError(t, err, "Failed to retrieve all clients from empty database")
		assert.Empty(t, clientsEncx, "Expected no clients from empty database")
	})

	t.Run("single client", func(t *testing.T) {
		ctx := context.Background()

		th.ClearClientsTable(t, ctx, testPool)

		// Create test clientEncx
		clientEncx := th.NewTestClientEncx(t)
		// err := repo.CreateClient(ctx, client)
		err := th.InsertClientEncx(t, ctx, testPool, clientEncx)
		require.NoError(t, err, "Failed to create client for single client test")

		// Retrieve all clientsEncx
		clientsEncx, err := repo.GetAllClients(ctx)
		assert.NoError(t, err, "Failed to retrieve all clients")
		assert.Len(t, clientsEncx, 1, "Expected exactly 1 client")

		// Verify the client data
		retrievedClientEncx := clientsEncx[0]
		assert.Equal(t, clientEncx.ID, retrievedClientEncx.ID, "Client ID should match")
		assert.Equal(t, clientEncx.NameHash, retrievedClientEncx.NameHash, "NameHash should match")
		assert.Equal(t, clientEncx.TypeHash, retrievedClientEncx.TypeHash, "TypeHash should match")
	})

	t.Run("multiple clients", func(t *testing.T) {
		ctx := context.Background()

		th.ClearClientsTable(t, ctx, testPool)

		// Create multiple test clients with delays to ensure different CreatedAt times
		createdClientsEncx := make([]*client.ClientEncx, 3)
		for i := 0; i < 3; i++ {
			createdClientsEncx[i] = th.NewTestClientEncx(t)
			err := th.InsertClientEncx(t, ctx, testPool, createdClientsEncx[i])
			require.NoError(t, err, "Failed to create client %d", i)

			// Small delay to ensure different creation times
			time.Sleep(10 * time.Millisecond)
		}

		// Retrieve all clients
		allClientsEncx, err := repo.GetAllClients(ctx)
		assert.NoError(t, err, "Failed to retrieve all clients")
		assert.Len(t, allClientsEncx, 3, "Expected exactly 3 clients")

		// Verify clients are ordered by CreatedAt DESC (most recent first)
		for i := 0; i < len(allClientsEncx)-1; i++ {
			current := allClientsEncx[i].CreatedAt
			next := allClientsEncx[i+1].CreatedAt
			assert.True(t, current.After(next) || current.Equal(next),
				"Client %d should be newer or same age as client %d", i, i+1)
		}

		// Verify all created clients are in the result
		clientIDs := make(map[string]bool)
		for _, client := range allClientsEncx {
			clientIDs[client.ID.String()] = true
		}

		for _, expectedClient := range createdClientsEncx {
			assert.True(t, clientIDs[expectedClient.ID.String()],
				"Expected client ID %s should be in the results", expectedClient.ID)
		}
	})

	t.Run("many clients", func(t *testing.T) {
		ctx := context.Background()

		th.ClearClientsTable(t, ctx, testPool)

		// Create many test clients
		numClients := 20
		createdClientsEncx := make([]*client.ClientEncx, numClients)
		for i := 0; i < numClients; i++ {
			createdClientsEncx[i] = th.NewTestClientEncx(t)
			err := th.InsertClientEncx(t, ctx, testPool, createdClientsEncx[i])
			require.NoError(t, err, "Failed to create client %d", i)
		}

		// Retrieve all clients
		allClientsEncx, err := repo.GetAllClients(ctx)
		assert.NoError(t, err, "Failed to retrieve all clients")
		assert.Len(t, allClientsEncx, numClients, "Expected exactly %d clients", numClients)
	})

	t.Run("clients with various field states", func(t *testing.T) {
		ctx := context.Background()

		// Create clients with different field states
		testCases := []struct {
			name   string
			modify func(*client.ClientEncx)
		}{
			{
				name: "client with empty name",
				modify: func(c *client.ClientEncx) {
					c.NameEncrypted = []byte("")
				},
			},
			{
				name: "client with empty type",
				modify: func(c *client.ClientEncx) {
					c.TypeEncrypted = []byte("")
				},
			},
			{
				name: "client with high key version",
				modify: func(c *client.ClientEncx) {
					c.KeyVersion = 999
				},
			},
		}

		for _, tc := range testCases {
			t.Run(tc.name, func(t *testing.T) {
				th.ClearClientsTable(t, ctx, testPool)

				clientEncx := th.NewTestClientEncx(t)
				tc.modify(clientEncx)

				err := th.InsertClientEncx(t, ctx, testPool, clientEncx)
				require.NoError(t, err, "Failed to create test client for %s", tc.name)

				// Retrieve all clients and find our specific client
				allClientsEncx, err := repo.GetAllClients(ctx)
				assert.NoError(t, err, "Failed to retrieve all clients")

				// Find our client in the results
				var foundClient *client.ClientEncx
				for _, c := range allClientsEncx {
					if c.ID == clientEncx.ID {
						foundClient = c
						break
					}
				}

				assert.NotNil(t, foundClient, "Should find created client in GetAllClients results")
				if foundClient != nil {
					// Verify the specific modifications are preserved
					tc.modify(clientEncx)
					assert.Equal(t, clientEncx.NameEncrypted, foundClient.NameEncrypted, "NameEncrypted should match for %s", tc.name)
					assert.Equal(t, clientEncx.TypeEncrypted, foundClient.TypeEncrypted, "TypeEncrypted should match for %s", tc.name)
					assert.Equal(t, clientEncx.KeyVersion, foundClient.KeyVersion, "KeyVersion should match for %s", tc.name)
				}
			})
		}
	})

	t.Run("context cancellation", func(t *testing.T) {
		// Create a context that will be cancelled
		ctx, cancel := context.WithCancel(context.Background())
		cancel() // Cancel immediately

		clientsEncx, err := repo.GetAllClients(ctx)
		assert.Error(t, err, "Expected context cancellation error, but got nil")
		assert.Nil(t, clientsEncx, "Retrieved clients should be nil for cancelled context")
	})

	t.Run("concurrent access", func(t *testing.T) {
		ctx := context.Background()

		th.ClearClientsTable(t, ctx, testPool)

		// Create some test clients first
		for i := 0; i < 3; i++ {
			clientEncx := th.NewTestClientEncx(t)
			err := th.InsertClientEncx(t, ctx, testPool, clientEncx)
			require.NoError(t, err, "Failed to create client %d", i)
		}

		// Test concurrent GetAllClients calls
		const numGoroutines = 5
		results := make(chan []*client.ClientEncx, numGoroutines)
		errors := make(chan error, numGoroutines)

		for i := 0; i < numGoroutines; i++ {
			go func() {
				clientsEncx, err := repo.GetAllClients(ctx)
				if err != nil {
					errors <- err
					return
				}
				results <- clientsEncx
			}()
		}

		// Collect results
		for i := 0; i < numGoroutines; i++ {
			select {
			case err := <-errors:
				assert.NoError(t, err, "Concurrent GetAllClients should not error")
			case clients := <-results:
				assert.NotEmpty(t, clients, "Concurrent GetAllClients should return clients")
			case <-time.After(5 * time.Second):
				t.Fatal("Concurrent GetAllClients timed out")
			}
		}
	})
}
