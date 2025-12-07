package caseRepository_test

import (
	"context"
	"testing"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	ch "github.com/hengadev/cluo_api/test/helpers/case"
)

// make test-func TEST_NAME=TestGetCaseByID TEST_PATH=internal/infrastructure/postgres/case/get_case_by_id_test.go

func TestGetCaseByID(t *testing.T) {
	if testPool == nil || repo == nil {
		t.Skip("Test database or repository not initialized")
	}

	t.Run("successful retrieval", func(t *testing.T) {
		ctx := context.Background()

		ch.ClearCasesTable(t, ctx, testPool)

		// Create test case data using helper
		caseEncx := ch.NewTestCaseEncx(t)

		// Insert case using the global repo (setup - should stop if fails)
		err := repo.CreateCase(ctx, caseEncx)
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

		ch.ClearCasesTable(t, ctx, testPool)

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
		ch.ClearCasesTable(t, ctx, testPool)

		// Try to retrieve any case
		testID := uuid.New()

		retrievedCaseEncx, err := repo.GetCaseByID(ctx, testID)
		assert.Error(t, err, "Expected error when retrieving from empty table")
		assert.Nil(t, retrievedCaseEncx, "Retrieved case should be nil")
	})

	t.Run("multiple cases exist", func(t *testing.T) {
		ctx := context.Background()

		ch.ClearCasesTable(t, ctx, testPool)

		// Create multiple cases
		case1 := ch.NewTestCaseEncx(t)
		case2 := ch.NewTestCaseEncx(t)
		case3 := ch.NewTestCaseEncx(t)

		// Insert all cases
		err := repo.CreateCase(ctx, case1)
		require.NoError(t, err)
		err = repo.CreateCase(ctx, case2)
		require.NoError(t, err)
		err = repo.CreateCase(ctx, case3)
		require.NoError(t, err)

		// Retrieve each case and verify
		retrievedCase1, err := repo.GetCaseByID(ctx, case1.ID)
		assert.NoError(t, err)
		assert.Equal(t, case1.ID, retrievedCase1.ID)

		retrievedCase2, err := repo.GetCaseByID(ctx, case2.ID)
		assert.NoError(t, err)
		assert.Equal(t, case2.ID, retrievedCase2.ID)

		retrievedCase3, err := repo.GetCaseByID(ctx, case3.ID)
		assert.NoError(t, err)
		assert.Equal(t, case3.ID, retrievedCase3.ID)

		// Verify non-existent ID still returns error
		nonExistentID := uuid.New()
		retrievedCase, err := repo.GetCaseByID(ctx, nonExistentID)
		assert.Error(t, err)
		assert.Nil(t, retrievedCase)
	})

	t.Run("case with assigned contact", func(t *testing.T) {
		ctx := context.Background()

		ch.ClearCasesTable(t, ctx, testPool)

		// Create test case with assigned contact
		caseEncx := ch.NewTestCaseEncx(t)
		contactID := uuid.New()
		caseEncx.AssignedContactID = &contactID

		// Insert case
		err := repo.CreateCase(ctx, caseEncx)
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

		ch.ClearCasesTable(t, ctx, testPool)

		// Create test case without assigned contact (nil)
		caseEncx := ch.NewTestCaseEncx(t)
		caseEncx.AssignedContactID = nil

		// Insert case
		err := repo.CreateCase(ctx, caseEncx)
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
