package clientRepository_test

import (
	"context"
	"testing"
	"time"

	th "github.com/hengadev/cluo_api/test/helpers/client"

	"github.com/google/uuid"
	"github.com/hengadev/encx"
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

		th.ClearClientsTable(t, ctx, testPool)

		// Create test clientEncx data using helper
		clientEncx := th.NewTestClientEncx(t)

		// Insert client first (setup)
		err := th.InsertClientEncx(t, ctx, testPool, clientEncx)
		require.NoError(t, err, "Failed to create client for update test")

		// Modify client data
		clientEncx.NameEncrypted = []byte("updated_name_encrypted")
		clientEncx.NameHash = "updated_name_hash"
		clientEncx.TypeEncrypted = []byte("updated_type_encrypted")
		clientEncx.TypeHash = "updated_type_hash"
		clientEncx.KeyVersion = 2

		// Test successful client update
		err = repo.UpdateClient(ctx, clientEncx)
		assert.NoError(t, err, "Failed to update client")

		// Verify client was updated by retrieving it
		retrievedClientEncx, err := th.GetClientEncxByID(t, ctx, testPool, clientEncx.ID)
		assert.NoError(t, err, "Failed to retrieve updated client")

		// Verify field values are updated
		assert.Equal(t, clientEncx.NameEncrypted, retrievedClientEncx.NameEncrypted, "NameEncrypted should be updated")
		assert.Equal(t, clientEncx.NameHash, retrievedClientEncx.NameHash, "NameHash should be updated")
		assert.Equal(t, clientEncx.TypeEncrypted, retrievedClientEncx.TypeEncrypted, "TypeEncrypted should be updated")
		assert.Equal(t, clientEncx.TypeHash, retrievedClientEncx.TypeHash, "TypeHash should be updated")
		assert.Equal(t, clientEncx.KeyVersion, retrievedClientEncx.KeyVersion, "KeyVersion should be updated")
	})

	t.Run("update non-existent client", func(t *testing.T) {
		ctx := context.Background()

		th.ClearClientsTable(t, ctx, testPool)

		// Create a client that doesn't exist in the database
		nonExistentClientEncx := th.NewTestClientEncx(t)

		err := repo.UpdateClient(ctx, nonExistentClientEncx)
		assert.Error(t, err, "Expected error when updating non-existent client")
	})

	t.Run("partial update", func(t *testing.T) {
		ctx := context.Background()

		th.ClearClientsTable(t, ctx, testPool)

		// Create test clientEncx
		clientEncx := th.NewTestClientEncx(t)
		err := repo.CreateClient(ctx, clientEncx)
		require.NoError(t, err, "Failed to create client for partial update test")

		// Update only specific fields
		originalNameEncrypted := clientEncx.NameEncrypted
		clientEncx.NameEncrypted = []byte("partially_updated_name")
		clientEncx.NameHash = "partially_updated_hash"
		// Keep other fields the same

		err = repo.UpdateClient(ctx, clientEncx)
		assert.NoError(t, err, "Failed to perform partial update")

		// Verify only updated fields changed
		retrievedClientEncx, err := th.GetClientEncxByID(t, ctx, testPool, clientEncx.ID)
		assert.NoError(t, err, "Failed to retrieve partially updated client")

		assert.Equal(t, clientEncx.NameEncrypted, retrievedClientEncx.NameEncrypted, "Name should be updated")
		assert.Equal(t, clientEncx.NameHash, retrievedClientEncx.NameHash, "NameHash should be updated")
		assert.NotEqual(t, originalNameEncrypted, retrievedClientEncx.NameEncrypted, "Name should be different from original")
	})

	t.Run("update with nil UUID", func(t *testing.T) {
		ctx := context.Background()

		th.ClearClientsTable(t, ctx, testPool)

		// Create clientEncx with nil UUID
		clientEncx := th.NewTestClientEncx(t)
		clientEncx.ID = uuid.Nil

		err := repo.UpdateClient(ctx, clientEncx)
		assert.Error(t, err, "Expected error when updating client with nil UUID")
	})

	t.Run("update with empty encrypted fields", func(t *testing.T) {
		ctx := context.Background()

		th.ClearClientsTable(t, ctx, testPool)
		// Create test clientEncx
		clientEncx := th.NewTestClientEncx(t)
		err := repo.CreateClient(ctx, clientEncx)
		require.NoError(t, err, "Failed to create client for empty fields update test")

		// Update with empty encrypted fields
		clientEncx.NameEncrypted = []byte("")
		clientEncx.TypeEncrypted = []byte("")

		err = repo.UpdateClient(ctx, clientEncx)
		assert.NoError(t, err, "Failed to update client with empty encrypted fields")

		// Verify the update
		retrievedClientEncx, err := th.GetClientEncxByID(t, ctx, testPool, clientEncx.ID)
		assert.NoError(t, err, "Failed to retrieve client with empty fields")

		assert.Equal(t, []byte(""), retrievedClientEncx.NameEncrypted, "NameEncrypted should be empty")
		assert.Equal(t, []byte(""), retrievedClientEncx.TypeEncrypted, "TypeEncrypted should be empty")
	})

	t.Run("context cancellation", func(t *testing.T) {
		// Create a context that will be cancelled
		ctx, cancel := context.WithCancel(context.Background())
		cancel() // Cancel immediately

		clientEncx := th.NewTestClientEncx(t)

		err := repo.UpdateClient(ctx, clientEncx)
		assert.Error(t, err, "Expected context cancellation error, but got nil")
	})

	t.Run("update with different metadata", func(t *testing.T) {
		ctx := context.Background()

		th.ClearClientsTable(t, ctx, testPool)

		// Create test clientEncx
		clientEncx := th.NewTestClientEncx(t)
		err := repo.CreateClient(ctx, clientEncx)
		require.NoError(t, err, "Failed to create client for metadata update test")

		// Update metadata field
		clientEncx.Metadata = encx.EncryptionMetadata{
			KEKAlias:         "updated_alias",
			EncryptionTime:   time.Now().Unix(),
			GeneratorVersion: "2.0.0",
		}

		err = repo.UpdateClient(ctx, clientEncx)
		assert.NoError(t, err, "Failed to update client metadata")

		// Verify the metadata update
		retrievedClientEncx, err := th.GetClientEncxByID(t, ctx, testPool, clientEncx.ID)
		assert.NoError(t, err, "Failed to retrieve client with updated metadata")

		assert.Equal(t, clientEncx.Metadata.KEKAlias, retrievedClientEncx.Metadata.KEKAlias, "KEKAlias should be updated")
	})
}
