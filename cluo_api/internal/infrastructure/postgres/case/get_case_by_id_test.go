package caseRepository_test

import (
	"context"
	"testing"

	caseHelpers "github.com/hengadev/cluo_api/test/helpers/case"
	clientHelpers "github.com/hengadev/cluo_api/test/helpers/client"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// make test-func TEST_NAME=TestGetCaseByID TEST_PATH=internal/infrastructure/postgres/case/get_case_by_id_test.go

func TestGetCaseByID(t *testing.T) {
	if testPool == nil || repo == nil {
		t.Skip("Test database or repository not initialized")
	}

	setupClient := func(t *testing.T, ctx context.Context) uuid.UUID {
		clientEncx := clientHelpers.NewTestClientEncx(t)
		err := clientHelpers.InsertClientEncx(t, ctx, testPool, clientEncx)
		require.NoError(t, err)
		return clientEncx.ID
	}
	setupContact := func(t *testing.T, ctx context.Context, clientID uuid.UUID) uuid.UUID {
		contactEncx := clientHelpers.NewTestContactEncx(t)
		contactEncx.ClientID = clientID
		err := clientHelpers.InsertContactEncx(t, ctx, testPool, contactEncx)
		require.NoError(t, err)
		return contactEncx.ID
	}

	t.Run("successful retrieval", func(t *testing.T) {
		ctx := context.Background()

		caseHelpers.ClearCasesTable(t, ctx, testPool)
		clientHelpers.ClearClientsTable(t, ctx, testPool)
		clientHelpers.ClearContactsTable(t, ctx, testPool)

		// Setup client and contact in database
		clientID := setupClient(t, ctx)
		contactID := setupContact(t, ctx, clientID)

		// Create test case data using helper
		caseEncx := caseHelpers.NewTestCaseEncx(t)
		caseEncx.ClientID = clientID
		caseEncx.AssignedContactID = &contactID

		// Insert case using the global repo (setup - should stop if fails)
		err := caseHelpers.InsertCaseEncx(t, ctx, testPool, caseEncx)
		require.NoError(t, err, "Failed to create case")

		// Test retrieving the case
		retrievedCaseEncx, err := repo.GetCaseByID(ctx, caseEncx.ID)
		assert.NoError(t, err, "Failed to retrieve case")
		require.NotNil(t, retrievedCaseEncx, "Retrieved case should not be nil")

		// Verify field values
		assert.Equal(t, caseEncx.ID, retrievedCaseEncx.ID, "Case ID should match")
		assert.Equal(t, caseEncx.ClientID, retrievedCaseEncx.ClientID, "Client ID should match")
		assert.Equal(t, caseEncx.AssignedContactID, retrievedCaseEncx.AssignedContactID, "Assigned Contact ID should match")
		assert.Equal(t, caseEncx.KeyVersion, retrievedCaseEncx.KeyVersion, "Key version should match")
		assert.Equal(t, caseEncx.TitleEncrypted, retrievedCaseEncx.TitleEncrypted, "Title encrypted should match")
		assert.Equal(t, caseEncx.DescriptionEncrypted, retrievedCaseEncx.DescriptionEncrypted, "Description encrypted should match")
		assert.Equal(t, caseEncx.StatusEncrypted, retrievedCaseEncx.StatusEncrypted, "Status encrypted should match")
	})

	t.Run("case not found", func(t *testing.T) {
		ctx := context.Background()

		caseHelpers.ClearCasesTable(t, ctx, testPool)

		// Try to retrieve a non-existent case
		nonExistentID := uuid.New()

		retrievedCaseEncx, err := repo.GetCaseByID(ctx, nonExistentID)
		assert.Error(t, err, "Expected error when retrieving non-existent case")
		assert.Nil(t, retrievedCaseEncx, "Retrieved case should be nil")
	})

	t.Run("nil UUID", func(t *testing.T) {
		ctx := context.Background()

		// Try to retrieve with nil UUID (all zeros)
		retrievedCaseEncx, err := repo.GetCaseByID(ctx, uuid.Nil)
		assert.Error(t, err, "Expected error when retrieving case with nil UUID")
		assert.Nil(t, retrievedCaseEncx, "Retrieved case should be nil")
	})

	t.Run("empty table", func(t *testing.T) {
		ctx := context.Background()

		// Clear all cases
		caseHelpers.ClearCasesTable(t, ctx, testPool)

		// Try to retrieve any case
		testID := uuid.New()

		retrievedCaseEncx, err := repo.GetCaseByID(ctx, testID)
		assert.Error(t, err, "Expected error when retrieving from empty table")
		assert.Nil(t, retrievedCaseEncx, "Retrieved case should be nil")
	})

	t.Run("multiple cases exist", func(t *testing.T) {
		ctx := context.Background()

		caseHelpers.ClearCasesTable(t, ctx, testPool)
		clientHelpers.ClearClientsTable(t, ctx, testPool)
		clientHelpers.ClearContactsTable(t, ctx, testPool)

		// Setup client and contact in database
		clientID := setupClient(t, ctx)
		contactID := setupContact(t, ctx, clientID)

		// Create multiple cases
		caseEncx1 := caseHelpers.NewTestCaseEncx(t)
		caseEncx1.ClientID = clientID
		caseEncx1.AssignedContactID = &contactID
		caseEncx2 := caseHelpers.NewTestCaseEncx(t)
		caseEncx2.ClientID = clientID
		caseEncx2.AssignedContactID = &contactID
		caseEncx3 := caseHelpers.NewTestCaseEncx(t)
		caseEncx3.ClientID = clientID
		caseEncx3.AssignedContactID = &contactID

		// Insert all cases
		err := caseHelpers.InsertCaseEncx(t, ctx, testPool, caseEncx1)
		require.NoError(t, err)

		err = caseHelpers.InsertCaseEncx(t, ctx, testPool, caseEncx2)
		require.NoError(t, err)

		err = caseHelpers.InsertCaseEncx(t, ctx, testPool, caseEncx3)
		require.NoError(t, err)

		// Retrieve each case and verify
		retrievedCase1, err := repo.GetCaseByID(ctx, caseEncx1.ID)
		assert.NoError(t, err)
		assert.Equal(t, caseEncx1.ID, retrievedCase1.ID)

		retrievedCase2, err := repo.GetCaseByID(ctx, caseEncx2.ID)
		assert.NoError(t, err)
		assert.Equal(t, caseEncx2.ID, retrievedCase2.ID)

		retrievedCase3, err := repo.GetCaseByID(ctx, caseEncx3.ID)
		assert.NoError(t, err)
		assert.Equal(t, caseEncx3.ID, retrievedCase3.ID)

		// Verify non-existent ID still returns error
		nonExistentID := uuid.New()
		retrievedCase, err := repo.GetCaseByID(ctx, nonExistentID)
		assert.Error(t, err)
		assert.Nil(t, retrievedCase)
	})

	t.Run("case with assigned contact", func(t *testing.T) {
		ctx := context.Background()

		caseHelpers.ClearCasesTable(t, ctx, testPool)
		clientHelpers.ClearClientsTable(t, ctx, testPool)
		clientHelpers.ClearContactsTable(t, ctx, testPool)

		// Setup client and contact in database
		clientID := setupClient(t, ctx)
		contactID := setupContact(t, ctx, clientID)

		// Create test case with assigned contact
		caseEncx := caseHelpers.NewTestCaseEncx(t)
		caseEncx.ClientID = clientID
		caseEncx.AssignedContactID = &contactID

		// Insert case
		err := caseHelpers.InsertCaseEncx(t, ctx, testPool, caseEncx)
		require.NoError(t, err)

		// Retrieve the case
		retrievedCaseEncx, err := repo.GetCaseByID(ctx, caseEncx.ID)
		assert.NoError(t, err)
		require.NotNil(t, retrievedCaseEncx)

		// Verify assigned contact ID is preserved
		require.NotNil(t, retrievedCaseEncx.AssignedContactID, "Assigned Contact ID should not be nil")
		assert.Equal(t, contactID, *retrievedCaseEncx.AssignedContactID, "Assigned Contact ID should match")
	})

	t.Run("case without assigned contact", func(t *testing.T) {
		ctx := context.Background()

		caseHelpers.ClearCasesTable(t, ctx, testPool)
		clientHelpers.ClearClientsTable(t, ctx, testPool)
		clientHelpers.ClearContactsTable(t, ctx, testPool)

		clientID := setupClient(t, ctx)

		// Create test case without assigned contact (nil)
		caseEncx := caseHelpers.NewTestCaseEncx(t)
		caseEncx.ClientID = clientID
		caseEncx.AssignedContactID = nil

		// Insert case
		err := caseHelpers.InsertCaseEncx(t, ctx, testPool, caseEncx)
		require.NoError(t, err)

		// Retrieve the case
		retrievedCaseEncx, err := repo.GetCaseByID(ctx, caseEncx.ID)
		assert.NoError(t, err)
		require.NotNil(t, retrievedCaseEncx)

		// Verify assigned contact ID is nil
		assert.Nil(t, retrievedCaseEncx.AssignedContactID, "Assigned Contact ID should be nil")
	})

	t.Run("context cancellation", func(t *testing.T) {
		// Create a context that will be cancelled
		ctx, cancel := context.WithCancel(context.Background())
		cancel() // Cancel immediately

		testID := uuid.New()

		// Test with cancelled context
		_, err := repo.GetCaseByID(ctx, testID)
		assert.Error(t, err, "Expected context cancellation error")
	})
}
