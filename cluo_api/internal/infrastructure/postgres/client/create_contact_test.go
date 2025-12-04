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

// make test-func TEST_NAME=TestCreateContact TEST_PATH=internal/infrastructure/postgres/client/create_contact_test.go

func TestCreateContact(t *testing.T) {
	if testPool == nil || repo == nil {
		t.Skip("Test database or repository not initialized")
	}

	setupClient := func(t *testing.T, ctx context.Context) *client.ClientEncx {
		clientEncx := th.NewTestClientEncx(t)
		err := th.InsertClientEncx(t, ctx, testPool, clientEncx)
		require.NoError(t, err)
		return clientEncx
	}

	t.Run("successful creation", func(t *testing.T) {
		ctx := context.Background()

		th.ClearClientsTable(t, ctx, testPool)
		th.ClearContactsTable(t, ctx, testPool)

		// Create test clientEncx data using helper
		clientEncx := setupClient(t, ctx)

		// Create test contactEncx data using helper
		contactEncx := th.NewTestContactEncx(t)
		contactEncx.ClientID = clientEncx.ID

		// Test successful contact creation using the global repo
		err := repo.CreateContact(ctx, contactEncx)
		assert.NoError(t, err, "Failed to create contact")

		// Verify the contact was inserted by retrieving it
		retrievedContactEncx, err := th.GetContactEncxByID(t, ctx, testPool, contactEncx.ID)
		assert.NoError(t, err, "Failed to retrieve inserted contact")

		// Verify field values
		assert.Equal(t, contactEncx.ID, retrievedContactEncx.ID, "Contact ID should match")
		assert.Equal(t, contactEncx.ClientID, retrievedContactEncx.ClientID, "Client hash should match")
		assert.Equal(t, contactEncx.EmailHash, retrievedContactEncx.EmailHash, "Email hash should match")
		assert.Equal(t, contactEncx.KeyVersion, retrievedContactEncx.KeyVersion, "Key version should match")
	})

	t.Run("with nil optional values", func(t *testing.T) {
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

		// Test successful contact creation with nil optional fields using the global repo
		err := repo.CreateContact(ctx, contactEncx)
		require.NoError(t, err, "Failed to create contact with nil values")

		// Verify the contact was inserted by retrieving it
		retrievedContactEncx, err := th.GetContactEncxByID(t, ctx, testPool, contactEncx.ID)
		assert.NoError(t, err, "Failed to retrieve inserted contact")

		// Verify that optional fields are indeed nil
		assert.Nil(t, retrievedContactEncx.PhoneEncrypted, "Expected PhoneEncrypted to be nil")
		assert.Nil(t, retrievedContactEncx.PositionEncrypted, "Expected PositionEncrypted to be nil")
	})

	t.Run("duplicate ID", func(t *testing.T) {
		ctx := context.Background()

		th.ClearClientsTable(t, ctx, testPool)
		th.ClearContactsTable(t, ctx, testPool)

		// Create test clientEncx data using helper
		clientEncx := setupClient(t, ctx)

		// Create first test contact
		contactEncx1 := th.NewTestContactEncx(t)
		contactEncx1.ClientID = clientEncx.ID

		// Insert first contact using the global repo (setup - should stop if fails)
		err := repo.CreateContact(ctx, contactEncx1)
		require.NoError(t, err, "Failed to create first contact")

		// Try to insert contact with same ID (should fail)
		contactEncx2 := th.NewTestContactEncx(t)
		contactEncx1.ClientID = clientEncx.ID
		contactEncx2.ID = contactEncx1.ID // Same ID, different data

		err = repo.CreateContact(ctx, contactEncx2)
		assert.Error(t, err, "Expected error when creating contact with duplicate ID")

		// Check that it's a database constraint violation (expected for duplicate ID)
		errStr := err.Error()
		assert.True(t, contains(errStr, "duplicate") || contains(errStr, "unique") || contains(errStr, "constraint"),
			"Expected constraint violation error, got: %v", err)
	})

	t.Run("empty required fields", func(t *testing.T) {
		ctx := context.Background()

		tests := []struct {
			name        string
			contactEncx *client.ContactEncx
		}{
			{
				name: "nil uuid",
				contactEncx: &client.ContactEncx{
					ID: uuid.Nil, // Invalid UUID
				},
			},
		}

		for _, tt := range tests {
			t.Run(tt.name, func(t *testing.T) {
				th.ClearClientsTable(t, ctx, testPool)
				th.ClearContactsTable(t, ctx, testPool)

				// Create test clientEncx data using helper
				clientEncx := setupClient(t, ctx)
				tt.contactEncx.ClientID = clientEncx.ID

				err := repo.CreateContact(ctx, tt.contactEncx)
				assert.Error(t, err, "Expected error for %s, but got nil", tt.name)
			})
		}
	})

	t.Run("context cancellation", func(t *testing.T) {
		// Create a context that will be cancelled
		ctx, cancel := context.WithCancel(context.Background())
		cancel() // Cancel immediately

		contactEncx := th.NewTestContactEncx(t)

		err := repo.CreateContact(ctx, contactEncx)
		assert.Error(t, err, "context cancelled")
	})
}

// contains checks if a string contains a substring
func contains(s, substr string) bool {
	return len(s) >= len(substr) && len(s) > 0 && len(substr) > 0 &&
		func() bool {
			for i := 0; i <= len(s)-len(substr); i++ {
				if s[i:i+len(substr)] == substr {
					return true
				}
			}
			return false
		}()
}
