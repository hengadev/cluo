package clientRepository_test

import (
	"context"
	"testing"

	"github.com/hengadev/cluo_api/internal/domain/client"
	ch "github.com/hengadev/cluo_api/test/helpers/client"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// make test-func TEST_NAME=TestExistsContact TEST_PATH=internal/infrastructure/postgres/client/exists_contact_test.go

func TestExistsContact(t *testing.T) {
	if testPool == nil || repo == nil {
		t.Skip("Test database or repository not initialized")
	}

	setupClient := func(t *testing.T, ctx context.Context) *client.ClientEncx {
		clientEncx := ch.NewTestClientEncx(t)
		err := ch.InsertClientEncx(t, ctx, testPool, clientEncx)
		require.NoError(t, err)
		return clientEncx
	}

	t.Run("contact exists", func(t *testing.T) {
		ctx := context.Background()

		ch.ClearContactsTable(t, ctx, testPool)

		// Create test clientEncx data using helper
		clientEncx := setupClient(t, ctx)

		// Create test contact
		contactEncx := ch.NewTestContactEncx(t)
		contactEncx.ClientID = clientEncx.ID

		// Insert contact using helper
		err := ch.InsertContactEncx(t, ctx, testPool, contactEncx)
		require.NoError(t, err, "Failed to insert contact")

		// Test existence check
		exists, err := repo.ExistsContact(ctx, contactEncx.ID)
		assert.NoError(t, err, "Failed to check if contact exists")
		assert.True(t, exists, "Contact should exist")
	})

	t.Run("contact does not exist", func(t *testing.T) {
		ctx := context.Background()

		ch.ClearContactsTable(t, ctx, testPool)

		// Test with non-existent UUID
		nonExistentID := uuid.New()

		exists, err := repo.ExistsContact(ctx, nonExistentID)
		assert.NoError(t, err, "Failed to check if contact exists")
		assert.False(t, exists, "Contact should not exist")
	})

	t.Run("nil UUID", func(t *testing.T) {
		ctx := context.Background()

		// Test with nil UUID
		exists, err := repo.ExistsContact(ctx, uuid.Nil)
		assert.NoError(t, err, "Failed to check if contact exists with nil UUID")
		assert.False(t, exists, "Contact should not exist for nil UUID")
	})

	t.Run("empty table", func(t *testing.T) {
		ctx := context.Background()

		// Clear all contacts
		ch.ClearContactsTable(t, ctx, testPool)

		// Test with any UUID
		testID := uuid.New()

		exists, err := repo.ExistsContact(ctx, testID)
		assert.NoError(t, err, "Failed to check if contact exists in empty table")
		assert.False(t, exists, "Contact should not exist in empty table")
	})

	t.Run("multiple contacts exist", func(t *testing.T) {
		ctx := context.Background()

		ch.ClearContactsTable(t, ctx, testPool)

		// Create test clientEncx data using helper
		clientEncx := setupClient(t, ctx)

		// Create multiple contacts
		contact1 := ch.NewTestContactEncx(t)
		contact1.ClientID = clientEncx.ID
		contact2 := ch.NewTestContactEncx(t)
		contact2.ClientID = clientEncx.ID
		contact3 := ch.NewTestContactEncx(t)
		contact3.ClientID = clientEncx.ID

		// Insert all contacts
		err := ch.InsertContactEncx(t, ctx, testPool, contact1)
		require.NoError(t, err)
		err = ch.InsertContactEncx(t, ctx, testPool, contact2)
		require.NoError(t, err)
		err = ch.InsertContactEncx(t, ctx, testPool, contact3)
		require.NoError(t, err)

		// Test each contact exists
		exists, err := repo.ExistsContact(ctx, contact1.ID)
		assert.NoError(t, err)
		assert.True(t, exists)

		exists, err = repo.ExistsContact(ctx, contact2.ID)
		assert.NoError(t, err)
		assert.True(t, exists)

		exists, err = repo.ExistsContact(ctx, contact3.ID)
		assert.NoError(t, err)
		assert.True(t, exists)

		// Test non-existent ID still returns false
		nonExistentID := uuid.New()
		exists, err = repo.ExistsContact(ctx, nonExistentID)
		assert.NoError(t, err)
		assert.False(t, exists)
	})

}

