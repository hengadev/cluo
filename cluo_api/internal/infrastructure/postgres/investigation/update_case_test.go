package investigationRepository_test

import (
	"context"
	"testing"
	"time"

	caseHelpers "github.com/hengadev/cluo_api/test/helpers/investigation"
	clientHelpers "github.com/hengadev/cluo_api/test/helpers/client"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// make test-func TEST_NAME=TestUpdateCase TEST_PATH=internal/infrastructure/postgres/case/update_case_test.go

func TestUpdateCase(t *testing.T) {
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

	t.Run("successful update", func(t *testing.T) {
		ctx := context.Background()

		caseHelpers.ClearCasesTable(t, ctx, testPool)
		clientHelpers.ClearContactsTable(t, ctx, testPool)
		clientHelpers.ClearClientsTable(t, ctx, testPool)

		clientID := setupClient(t, ctx)
		contactID := setupContact(t, ctx, clientID)

		// Create test case data using helper
		caseEncx := caseHelpers.NewTestCaseEncx(t)
		caseEncx.ClientID = clientID
		caseEncx.AssignedContactID = &contactID

		// Insert case using the global repo
		err := caseHelpers.InsertCaseEncx(t, ctx, testPool, caseEncx)
		require.NoError(t, err, "Failed to create case")

		// Modify the case for update
		updatedTime := time.Now()
		caseEncx.ClientID = setupClient(t, ctx)
		caseEncx.UpdatedAtEncrypted = []byte(updatedTime.Format(time.RFC3339))
		caseEncx.TitleEncrypted = []byte("updated_title_encrypted")
		caseEncx.DescriptionEncrypted = []byte("updated_description_encrypted")
		caseEncx.ExternalReferenceEncrypted = []byte("updated_external_ref_encrypted")
		caseEncx.CaseType = "Updated Case Type"
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
		assert.Equal(t, caseEncx.ExternalReferenceEncrypted, retrievedCase.ExternalReferenceEncrypted, "ExternalReferenceEncrypted should be updated")
		assert.Equal(t, caseEncx.CaseType, retrievedCase.CaseType, "CaseType should be updated")
		assert.Equal(t, caseEncx.StatusEncrypted, retrievedCase.StatusEncrypted, "Status encrypted should be updated")
		assert.Equal(t, caseEncx.UpdatedAtEncrypted, retrievedCase.UpdatedAtEncrypted, "Updated at should be updated")
		assert.Equal(t, caseEncx.KeyVersion, retrievedCase.KeyVersion, "Key version should be updated")
	})

	t.Run("case not found", func(t *testing.T) {
		ctx := context.Background()

		caseHelpers.ClearCasesTable(t, ctx, testPool)

		// Try to update a non-existent case
		nonExistentID := uuid.New()
		caseEncx := caseHelpers.NewTestCaseEncx(t)
		caseEncx.ID = nonExistentID

		err := repo.UpdateCase(ctx, caseEncx)
		assert.Error(t, err, "Should return error when updating non-existent case")
		assert.Contains(t, err.Error(), "not found", "Error should indicate case not found")
	})

	t.Run("nil case", func(t *testing.T) {
		ctx := context.Background()

		caseHelpers.ClearCasesTable(t, ctx, testPool)

		// Try to update with nil case
		err := repo.UpdateCase(ctx, nil)
		assert.Error(t, err, "Should return error when updating nil case")
		assert.Contains(t, err.Error(), "case cannot be nil", "Error should mention case cannot be nil")
	})

	t.Run("update with nil assigned contact", func(t *testing.T) {
		ctx := context.Background()

		caseHelpers.ClearCasesTable(t, ctx, testPool)
		clientHelpers.ClearContactsTable(t, ctx, testPool)
		clientHelpers.ClearClientsTable(t, ctx, testPool)

		// Setup client and contact in database
		clientID := setupClient(t, ctx)
		contactID := setupContact(t, ctx, clientID)

		// Create test case with assigned contact
		caseEncx := caseHelpers.NewTestCaseEncx(t)
		caseEncx.ClientID = clientID
		caseEncx.AssignedContactID = &contactID

		err := caseHelpers.InsertCaseEncx(t, ctx, testPool, caseEncx)
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

		caseHelpers.ClearCasesTable(t, ctx, testPool)
		clientHelpers.ClearContactsTable(t, ctx, testPool)
		clientHelpers.ClearClientsTable(t, ctx, testPool)

		// Setup client in database
		clientID := setupClient(t, ctx)

		// Create test case without assigned contact
		caseEncx := caseHelpers.NewTestCaseEncx(t)
		caseEncx.ClientID = clientID
		caseEncx.AssignedContactID = nil

		err := caseHelpers.InsertCaseEncx(t, ctx, testPool, caseEncx)
		require.NoError(t, err, "Failed to create case")

		// Update case to add assigned contact
		newContactID := setupContact(t, ctx, clientID)
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

		caseHelpers.ClearCasesTable(t, ctx, testPool)
		clientHelpers.ClearContactsTable(t, ctx, testPool)
		clientHelpers.ClearClientsTable(t, ctx, testPool)

		// Setup client and contact in database
		clientID := setupClient(t, ctx)

		// Create test case
		caseEncx := caseHelpers.NewTestCaseEncx(t)
		caseEncx.ClientID = clientID
		caseEncx.AssignedContactID = nil
		originalDescription := caseEncx.DescriptionEncrypted
		originalStatus := caseEncx.StatusEncrypted

		err := caseHelpers.InsertCaseEncx(t, ctx, testPool, caseEncx)
		require.NoError(t, err, "Failed to create case")

		// Update only some fields
		caseEncx.ClientID = setupClient(t, ctx)
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

		caseHelpers.ClearCasesTable(t, ctx, testPool)
		clientHelpers.ClearContactsTable(t, ctx, testPool)
		clientHelpers.ClearClientsTable(t, ctx, testPool)

		// Setup client and contact in database
		clientID := setupClient(t, ctx)

		// Create multiple test cases
		caseEncx1 := caseHelpers.NewTestCaseEncx(t)
		caseEncx1.ClientID = clientID
		caseEncx1.AssignedContactID = nil

		caseEncx2 := caseHelpers.NewTestCaseEncx(t)
		caseEncx2.ClientID = clientID
		caseEncx2.AssignedContactID = nil

		err := caseHelpers.InsertCaseEncx(t, ctx, testPool, caseEncx1)
		require.NoError(t, err, "Failed to create case1")
		err = caseHelpers.InsertCaseEncx(t, ctx, testPool, caseEncx2)
		require.NoError(t, err, "Failed to create case2")

		// Update only case1
		caseEncx1.TitleEncrypted = []byte("case1_updated")
		caseEncx1.KeyVersion = 3

		err = repo.UpdateCase(ctx, caseEncx1)
		assert.NoError(t, err, "Failed to update case1")

		// Verify case1 was updated
		retrievedCase1, err := repo.GetCaseByID(ctx, caseEncx1.ID)
		assert.NoError(t, err, "Failed to retrieve case1")
		assert.Equal(t, caseEncx1.TitleEncrypted, retrievedCase1.TitleEncrypted, "Case1 title should be updated")
		assert.Equal(t, caseEncx1.KeyVersion, retrievedCase1.KeyVersion, "Case1 key version should be updated")

		// Verify case2 remains unchanged
		retrievedCase2, err := repo.GetCaseByID(ctx, caseEncx2.ID)
		assert.NoError(t, err, "Failed to retrieve case2")
		assert.NotEqual(t, caseEncx1.TitleEncrypted, retrievedCase2.TitleEncrypted, "Case2 title should be different from case1")
		assert.Equal(t, caseEncx2.KeyVersion, retrievedCase2.KeyVersion, "Case2 key version should remain unchanged")
	})
}
