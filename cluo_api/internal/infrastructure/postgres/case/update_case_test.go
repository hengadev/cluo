package caseRepository_test

import (
	"context"
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	ch "github.com/hengadev/cluo_api/test/helpers/case"
)

// make test-func TEST_NAME=TestUpdateCase TEST_PATH=internal/infrastructure/postgres/case/update_case_test.go

func TestUpdateCase(t *testing.T) {
	if testPool == nil || repo == nil {
		t.Skip("Test database or repository not initialized")
	}

	t.Run("successful update", func(t *testing.T) {
		ctx := context.Background()

		ch.ClearCasesTable(t, ctx, testPool)

		// Create test case data using helper
		caseEncx := ch.NewTestCaseEncx(t)

		// Insert case using the global repo
		err := repo.CreateCase(ctx, caseEncx)
		require.NoError(t, err, "Failed to create case")

		// Modify the case for update
		updatedTime := time.Now()
		caseEncx.ClientID = uuid.New()
		caseEncx.UpdatedAtEncrypted = []byte(updatedTime.Format(time.RFC3339))
		caseEncx.TitleEncrypted = []byte("updated_title_encrypted")
		caseEncx.DescriptionEncrypted = []byte("updated_description_encrypted")
		caseEncx.StatusEncrypted = []byte("updated_status_encrypted")
		caseEncx.KeyVersion = 2

		// Test updating the case
		err = repo.UpdateCase(ctx, caseEncx)
		assert.NoError(t, err, "Failed to update case")

		// Verify the case was updated
		retrievedCase, err := repo.GetCaseByID(ctx, caseEncx.ID)
		assert.NoError(t, err, "Failed to retrieve updated case")
		require.NotNil(t, retrievedCase, "Retrieved case should not be nil")

		// Verify updated field values
		assert.Equal(t, caseEncx.ID, retrievedCase.ID, "Case ID should match")
		assert.Equal(t, caseEncx.ClientID, retrievedCase.ClientID, "Client ID should be updated")
		assert.Equal(t, caseEncx.TitleEncrypted, retrievedCase.TitleEncrypted, "Title encrypted should be updated")
		assert.Equal(t, caseEncx.DescriptionEncrypted, retrievedCase.DescriptionEncrypted, "Description encrypted should be updated")
		assert.Equal(t, caseEncx.StatusEncrypted, retrievedCase.StatusEncrypted, "Status encrypted should be updated")
		assert.Equal(t, caseEncx.UpdatedAtEncrypted, retrievedCase.UpdatedAtEncrypted, "Updated at should be updated")
		assert.Equal(t, caseEncx.KeyVersion, retrievedCase.KeyVersion, "Key version should be updated")
	})

	t.Run("case not found", func(t *testing.T) {
		ctx := context.Background()

		ch.ClearCasesTable(t, ctx, testPool)

		// Try to update a non-existent case
		nonExistentID := uuid.New()
		caseEncx := ch.NewTestCaseEncx(t)
		caseEncx.ID = nonExistentID

		err := repo.UpdateCase(ctx, caseEncx)
		assert.Error(t, err, "Should return error when updating non-existent case")
		assert.Contains(t, err.Error(), "not found", "Error should indicate case not found")
	})

	t.Run("nil case", func(t *testing.T) {
		ctx := context.Background()

		ch.ClearCasesTable(t, ctx, testPool)

		// Try to update with nil case
		err := repo.UpdateCase(ctx, nil)
		assert.Error(t, err, "Should return error when updating nil case")
		assert.Contains(t, err.Error(), "case cannot be nil", "Error should mention case cannot be nil")
	})

	t.Run("update with nil assigned contact", func(t *testing.T) {
		ctx := context.Background()

		ch.ClearCasesTable(t, ctx, testPool)

		// Create test case with assigned contact
		caseEncx := ch.NewTestCaseEncx(t)
		err := repo.CreateCase(ctx, caseEncx)
		require.NoError(t, err, "Failed to create case")

		// Update case to remove assigned contact
		caseEncx.AssignedContactID = nil
		caseEncx.TitleEncrypted = []byte("updated_title_encrypted")

		err = repo.UpdateCase(ctx, caseEncx)
		assert.NoError(t, err, "Failed to update case with nil assigned contact")

		// Verify the case was updated
		retrievedCase, err := repo.GetCaseByID(ctx, caseEncx.ID)
		assert.NoError(t, err, "Failed to retrieve updated case")
		assert.Nil(t, retrievedCase.AssignedContactID, "Assigned contact should be nil")
		assert.Equal(t, caseEncx.TitleEncrypted, retrievedCase.TitleEncrypted, "Title should be updated")
	})

	t.Run("update with new assigned contact", func(t *testing.T) {
		ctx := context.Background()

		ch.ClearCasesTable(t, ctx, testPool)

		// Create test case without assigned contact
		caseEncx := ch.NewTestCaseEncx(t)
		caseEncx.AssignedContactID = nil
		err := repo.CreateCase(ctx, caseEncx)
		require.NoError(t, err, "Failed to create case")

		// Update case to add assigned contact
		newContactID := uuid.New()
		caseEncx.AssignedContactID = &newContactID
		caseEncx.DescriptionEncrypted = []byte("updated_description_encrypted")

		err = repo.UpdateCase(ctx, caseEncx)
		assert.NoError(t, err, "Failed to update case with new assigned contact")

		// Verify the case was updated
		retrievedCase, err := repo.GetCaseByID(ctx, caseEncx.ID)
		assert.NoError(t, err, "Failed to retrieve updated case")
		require.NotNil(t, retrievedCase.AssignedContactID, "Assigned contact should not be nil")
		assert.Equal(t, *caseEncx.AssignedContactID, *retrievedCase.AssignedContactID, "Assigned contact should be updated")
		assert.Equal(t, caseEncx.DescriptionEncrypted, retrievedCase.DescriptionEncrypted, "Description should be updated")
	})

	t.Run("partial update", func(t *testing.T) {
		ctx := context.Background()

		ch.ClearCasesTable(t, ctx, testPool)

		// Create test case
		caseEncx := ch.NewTestCaseEncx(t)
		originalDescription := caseEncx.DescriptionEncrypted
		originalStatus := caseEncx.StatusEncrypted

		err := repo.CreateCase(ctx, caseEncx)
		require.NoError(t, err, "Failed to create case")

		// Update only some fields
		caseEncx.ClientID = uuid.New()
		caseEncx.TitleEncrypted = []byte("partially_updated_title")
		caseEncx.UpdatedAtEncrypted = []byte(time.Now().Format(time.RFC3339))
		// Keep description and status unchanged

		err = repo.UpdateCase(ctx, caseEncx)
		assert.NoError(t, err, "Failed to partially update case")

		// Verify the case was partially updated
		retrievedCase, err := repo.GetCaseByID(ctx, caseEncx.ID)
		assert.NoError(t, err, "Failed to retrieve updated case")
		require.NotNil(t, retrievedCase, "Retrieved case should not be nil")

		// Verify updated fields
		assert.Equal(t, caseEncx.ClientID, retrievedCase.ClientID, "Client ID should be updated")
		assert.Equal(t, caseEncx.TitleEncrypted, retrievedCase.TitleEncrypted, "Title should be updated")
		assert.Equal(t, caseEncx.UpdatedAtEncrypted, retrievedCase.UpdatedAtEncrypted, "Updated at should be updated")

		// Verify unchanged fields are still the same
		assert.Equal(t, originalDescription, retrievedCase.DescriptionEncrypted, "Description should remain unchanged")
		assert.Equal(t, originalStatus, retrievedCase.StatusEncrypted, "Status should remain unchanged")
	})

	t.Run("multiple cases update", func(t *testing.T) {
		ctx := context.Background()

		ch.ClearCasesTable(t, ctx, testPool)

		// Create multiple test cases
		case1 := ch.NewTestCaseEncx(t)
		case2 := ch.NewTestCaseEncx(t)

		err := repo.CreateCase(ctx, case1)
		require.NoError(t, err, "Failed to create case1")
		err = repo.CreateCase(ctx, case2)
		require.NoError(t, err, "Failed to create case2")

		// Update only case1
		case1.TitleEncrypted = []byte("case1_updated")
		case1.KeyVersion = 3

		err = repo.UpdateCase(ctx, case1)
		assert.NoError(t, err, "Failed to update case1")

		// Verify case1 was updated
		retrievedCase1, err := repo.GetCaseByID(ctx, case1.ID)
		assert.NoError(t, err, "Failed to retrieve case1")
		assert.Equal(t, case1.TitleEncrypted, retrievedCase1.TitleEncrypted, "Case1 title should be updated")
		assert.Equal(t, case1.KeyVersion, retrievedCase1.KeyVersion, "Case1 key version should be updated")

		// Verify case2 remains unchanged
		retrievedCase2, err := repo.GetCaseByID(ctx, case2.ID)
		assert.NoError(t, err, "Failed to retrieve case2")
		assert.NotEqual(t, case1.TitleEncrypted, retrievedCase2.TitleEncrypted, "Case2 title should be different from case1")
		assert.Equal(t, case2.KeyVersion, retrievedCase2.KeyVersion, "Case2 key version should remain unchanged")
	})
}

