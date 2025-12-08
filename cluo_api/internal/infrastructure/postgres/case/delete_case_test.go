package caseRepository_test

import (
	"context"
	"testing"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	caseHelpers "github.com/hengadev/cluo_api/test/helpers/case"
	clientHelpers "github.com/hengadev/cluo_api/test/helpers/client"
)

// make test-func TEST_NAME=TestDeleteCase TEST_PATH=internal/infrastructure/postgres/case/delete_case_test.go

func TestDeleteCase(t *testing.T) {
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

	t.Run("successful deletion", func(t *testing.T) {
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

		// Insert case using the global repo
		err := caseHelpers.InsertCaseEncx(t, ctx, testPool, caseEncx)
		require.NoError(t, err, "Failed to create case")

		// Verify case exists before deletion
		retrievedCase, err := repo.GetCaseByID(ctx, caseEncx.ID)
		require.NoError(t, err, "Failed to retrieve case before deletion")
		require.NotNil(t, retrievedCase, "Case should exist before deletion")

		// Test deleting the case
		err = repo.DeleteCase(ctx, caseEncx.ID)
		assert.NoError(t, err, "Failed to delete case")

		// Verify case no longer exists
		retrievedCase, err = repo.GetCaseByID(ctx, caseEncx.ID)
		assert.Error(t, err, "Should return error when retrieving deleted case")
		assert.Nil(t, retrievedCase, "Retrieved case should be nil after deletion")
		assert.Contains(t, err.Error(), "not found", "Error should indicate case not found")
	})

	t.Run("case not found", func(t *testing.T) {
		ctx := context.Background()

		caseHelpers.ClearCasesTable(t, ctx, testPool)

		// Try to delete a non-existent case
		nonExistentID := uuid.New()
		err := repo.DeleteCase(ctx, nonExistentID)
		assert.Error(t, err, "Should return error when deleting non-existent case")
		assert.Contains(t, err.Error(), "not found", "Error should indicate case not found")
	})

	t.Run("nil UUID", func(t *testing.T) {
		ctx := context.Background()

		caseHelpers.ClearCasesTable(t, ctx, testPool)

		// Try to delete with nil UUID
		err := repo.DeleteCase(ctx, uuid.Nil)
		assert.Error(t, err, "Should return error when deleting with nil UUID")
		assert.Contains(t, err.Error(), "not found", "Error should indicate case not found")
	})

	t.Run("delete multiple cases", func(t *testing.T) {
		ctx := context.Background()

		caseHelpers.ClearCasesTable(t, ctx, testPool)
		clientHelpers.ClearClientsTable(t, ctx, testPool)
		clientHelpers.ClearContactsTable(t, ctx, testPool)

		// Setup client and contact in database
		clientID := setupClient(t, ctx)
		contactID := setupContact(t, ctx, clientID)

		// Create multiple test cases
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
		require.NoError(t, err, "Failed to create case1")
		err = caseHelpers.InsertCaseEncx(t, ctx, testPool, caseEncx2)
		require.NoError(t, err, "Failed to create case2")
		err = caseHelpers.InsertCaseEncx(t, ctx, testPool, caseEncx3)
		require.NoError(t, err, "Failed to create case3")

		// Delete case2
		err = repo.DeleteCase(ctx, caseEncx2.ID)
		assert.NoError(t, err, "Failed to delete case2")

		// Verify case2 is deleted but case1 and case3 still exist
		retrievedCase1, err := repo.GetCaseByID(ctx, caseEncx1.ID)
		assert.NoError(t, err, "Case1 should still exist")
		assert.NotNil(t, retrievedCase1, "Case1 should not be nil")

		retrievedCase2, err := repo.GetCaseByID(ctx, caseEncx2.ID)
		assert.Error(t, err, "Case2 should not exist")
		assert.Nil(t, retrievedCase2, "Case2 should be nil")

		retrievedCase3, err := repo.GetCaseByID(ctx, caseEncx3.ID)
		assert.NoError(t, err, "Case3 should still exist")
		assert.NotNil(t, retrievedCase3, "Case3 should not be nil")
	})

	t.Run("delete case with assigned contact", func(t *testing.T) {
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
		require.NoError(t, err, "Failed to create case with assigned contact")

		// Delete the case
		err = repo.DeleteCase(ctx, caseEncx.ID)
		assert.NoError(t, err, "Failed to delete case with assigned contact")

		// Verify case is deleted
		retrievedCase, err := repo.GetCaseByID(ctx, caseEncx.ID)
		assert.Error(t, err, "Case should not exist after deletion")
		assert.Nil(t, retrievedCase, "Retrieved case should be nil")
	})
}
