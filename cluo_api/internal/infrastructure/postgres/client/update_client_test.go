package clientRepository_test

import (
	"context"
	"testing"
	"time"

	"github.com/hengadev/cluo_api/internal/domain/client"
	th "github.com/hengadev/cluo_api/test/helpers/client"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// make test-func TEST_NAME=TestUpdateClient TEST_PATH=internal/infrastructure/postgres/client/update_client_test.go
func TestUpdateClient(t *testing.T) {
	if testPool == nil || repo == nil {
		t.Skip("Test database or repository not initialized")
	}

	t.Run("successful update", func(t *testing.T) {
		ctx := context.Background()

		// Create test client data using helper
		client := th.NewTestClientEncx(t)

		// Insert client first (setup)
		err := repo.CreateClient(ctx, client)
		require.NoError(t, err, "Failed to create client for update test")

		// Modify client data
		client.NameEncrypted = []byte("updated_name_encrypted")
		client.NameHash = "updated_name_hash"
		client.TypeEncrypted = []byte("updated_type_encrypted")
		client.TypeHash = "updated_type_hash"
		client.ContactIDsEncrypted = []byte("updated_contact_ids_encrypted")
		client.KeyVersion = 2

		// Test successful client update
		err = repo.UpdateClient(ctx, client)
		assert.NoError(t, err, "Failed to update client")

		// Verify client was updated by retrieving it
		retrievedClient, err := th.GetClientEncxByID(t, ctx, testPool, client.ID)
		assert.NoError(t, err, "Failed to retrieve updated client")

		// Verify field values are updated
		assert.Equal(t, client.NameEncrypted, retrievedClient.NameEncrypted, "NameEncrypted should be updated")
		assert.Equal(t, client.NameHash, retrievedClient.NameHash, "NameHash should be updated")
		assert.Equal(t, client.TypeEncrypted, retrievedClient.TypeEncrypted, "TypeEncrypted should be updated")
		assert.Equal(t, client.TypeHash, retrievedClient.TypeHash, "TypeHash should be updated")
		assert.Equal(t, client.ContactIDsEncrypted, retrievedClient.ContactIDsEncrypted, "ContactIDsEncrypted should be updated")
		assert.Equal(t, client.KeyVersion, retrievedClient.KeyVersion, "KeyVersion should be updated")
	})

	t.Run("update non-existent client", func(t *testing.T) {
		ctx := context.Background()

		// Create a client that doesn't exist in the database
		nonExistentClient := th.NewTestClientEncx(t)

		err := repo.UpdateClient(ctx, nonExistentClient)
		assert.Error(t, err, "Expected error when updating non-existent client")
	})

	t.Run("partial update", func(t *testing.T) {
		ctx := context.Background()

		// Create test client
		client := th.NewTestClientEncx(t)
		err := repo.CreateClient(ctx, client)
		require.NoError(t, err, "Failed to create client for partial update test")

		// Update only specific fields
		originalNameEncrypted := client.NameEncrypted
		client.NameEncrypted = []byte("partially_updated_name")
		client.NameHash = "partially_updated_hash"
		// Keep other fields the same

		err = repo.UpdateClient(ctx, client)
		assert.NoError(t, err, "Failed to perform partial update")

		// Verify only updated fields changed
		retrievedClient, err := th.GetClientEncxByID(t, ctx, testPool, client.ID)
		assert.NoError(t, err, "Failed to retrieve partially updated client")

		assert.Equal(t, client.NameEncrypted, retrievedClient.NameEncrypted, "Name should be updated")
		assert.Equal(t, client.NameHash, retrievedClient.NameHash, "NameHash should be updated")
		assert.NotEqual(t, originalNameEncrypted, retrievedClient.NameEncrypted, "Name should be different from original")
	})

	t.Run("update with nil UUID", func(t *testing.T) {
		ctx := context.Background()

		// Create client with nil UUID
		client := th.NewTestClientEncx(t)
		client.ID = uuid.Nil

		err := repo.UpdateClient(ctx, client)
		assert.Error(t, err, "Expected error when updating client with nil UUID")
	})

	t.Run("update with empty encrypted fields", func(t *testing.T) {
		ctx := context.Background()

		// Create test client
		client := th.NewTestClientEncx(t)
		err := repo.CreateClient(ctx, client)
		require.NoError(t, err, "Failed to create client for empty fields update test")

		// Update with empty encrypted fields
		client.NameEncrypted = []byte("")
		client.TypeEncrypted = []byte("")
		client.ContactIDsEncrypted = []byte("")

		err = repo.UpdateClient(ctx, client)
		assert.NoError(t, err, "Failed to update client with empty encrypted fields")

		// Verify the update
		retrievedClient, err := th.GetClientEncxByID(t, ctx, testPool, client.ID)
		assert.NoError(t, err, "Failed to retrieve client with empty fields")

		assert.Equal(t, []byte(""), retrievedClient.NameEncrypted, "NameEncrypted should be empty")
		assert.Equal(t, []byte(""), retrievedClient.TypeEncrypted, "TypeEncrypted should be empty")
		assert.Equal(t, []byte(""), retrievedClient.ContactIDsEncrypted, "ContactIDsEncrypted should be empty")
	})

	t.Run("context cancellation", func(t *testing.T) {
		// Create a context that will be cancelled
		ctx, cancel := context.WithCancel(context.Background())
		cancel() // Cancel immediately

		client := th.NewTestClientEncx(t)

		err := repo.UpdateClient(ctx, client)
		assert.Error(t, err, "Expected context cancellation error, but got nil")
	})

	t.Run("update with different metadata", func(t *testing.T) {
		ctx := context.Background()

		// Create test client
		client := th.NewTestClientEncx(t)
		err := repo.CreateClient(ctx, client)
		require.NoError(t, err, "Failed to create client for metadata update test")

		// Update metadata field
		client.Metadata = client.EncryptionMetadata{
			KEKAlias:         "updated_alias",
			EncryptionTime:   time.Now().Unix(),
			GeneratorVersion: "2.0.0",
		}

		err = repo.UpdateClient(ctx, client)
		assert.NoError(t, err, "Failed to update client metadata")

		// Verify the metadata update
		retrievedClient, err := th.GetClientEncxByID(t, ctx, testPool, client.ID)
		assert.NoError(t, err, "Failed to retrieve client with updated metadata")

		assert.Equal(t, client.Metadata.KEKAlias, retrievedClient.Metadata.KEKAlias, "KEKAlias should be updated")
	})
}