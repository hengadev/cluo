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

// make test-func TEST_NAME=TestCreateClient TEST_PATH=internal/infrastructure/postgres/client/create_client_test.go
func TestCreateClient(t *testing.T) {
	if testPool == nil || repo == nil {
		t.Skip("Test database or repository not initialized")
	}

	t.Run("successful creation", func(t *testing.T) {
		ctx := context.Background()

		// Create test client data using helper
		client := th.NewTestClientEncx(t)

		// Test successful client creation using the global repo
		err := repo.CreateClient(ctx, client)
		assert.NoError(t, err, "Failed to create client")

		// Verify the client was inserted by retrieving it
		retrievedClient, err := th.GetClientEncxByID(t, ctx, testPool, client.ID)
		assert.NoError(t, err, "Failed to retrieve inserted client")

		// Verify field values
		assert.Equal(t, client.ID, retrievedClient.ID, "Client ID should match")
		assert.Equal(t, client.NameHash, retrievedClient.NameHash, "Name hash should match")
		assert.Equal(t, client.TypeHash, retrievedClient.TypeHash, "Type hash should match")
		assert.Equal(t, client.KeyVersion, retrievedClient.KeyVersion, "Key version should match")
	})

	t.Run("with empty contact IDs", func(t *testing.T) {
		ctx := context.Background()

		// Create test client and set contact IDs to empty
		client := th.NewTestClientEncx(t)
		client.ContactIDsEncrypted = []byte("[]")

		// Test successful client creation with empty contact IDs
		err := repo.CreateClient(ctx, client)
		require.NoError(t, err, "Failed to create client with empty contact IDs")

		// Verify the client was inserted by retrieving it
		retrievedClient, err := th.GetClientEncxByID(t, ctx, testPool, client.ID)
		assert.NoError(t, err, "Failed to retrieve inserted client")

		// Verify that contact IDs are indeed empty
		assert.Equal(t, []byte("[]"), retrievedClient.ContactIDsEncrypted, "Expected ContactIDsEncrypted to be empty")
	})

	t.Run("duplicate ID", func(t *testing.T) {
		ctx := context.Background()

		// Create first test client
		client1 := th.NewTestClientEncx(t)

		// Insert first client using the global repo (setup - should stop if fails)
		err := repo.CreateClient(ctx, client1)
		require.NoError(t, err, "Failed to create first client")

		// Try to insert client with same ID (should fail)
		client2 := th.NewTestClientEncx(t)
		client2.ID = client1.ID // Same ID, different data

		err = repo.CreateClient(ctx, client2)
		assert.Error(t, err, "Expected error when creating client with duplicate ID")

		// Check that it's a database constraint violation (expected for duplicate ID)
		errStr := err.Error()
		assert.True(t, contains(errStr, "duplicate") || contains(errStr, "unique") || contains(errStr, "constraint"),
			"Expected constraint violation error, got: %v", err)
	})

	t.Run("empty required fields", func(t *testing.T) {
		ctx := context.Background()

		tests := []struct {
			name   string
			client *client.ClientEncx
		}{
			{
				name: "nil uuid",
				client: &client.ClientEncx{
					ID: uuid.Nil, // Invalid UUID
				},
			},
		}

		for _, tt := range tests {
			t.Run(tt.name, func(t *testing.T) {
				err := repo.CreateClient(ctx, tt.client)
				assert.Error(t, err, "Expected error for %s, but got nil", tt.name)
			})
		}
	})

	t.Run("context cancellation", func(t *testing.T) {
		// Create a context that will be cancelled
		ctx, cancel := context.WithCancel(context.Background())
		cancel() // Cancel immediately

		client := th.NewTestClientEncx(t)

		err := repo.CreateClient(ctx, client)
		assert.Error(t, err, "Expected context cancellation error, but got nil")
	})
}

