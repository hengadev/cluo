package case_test

import (
	"context"
	"net/http"
	"testing"
	"time"

	tu "github.com/hengadev/cluo_api/internal/common/testutils"
	caseDomain "github.com/hengadev/cluo_api/internal/domain/case"
	"github.com/hengadev/cluo_api/internal/domain/client"
	ch "github.com/hengadev/cluo_api/test/helpers/case"
	clientHelpers "github.com/hengadev/cluo_api/test/helpers/client"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// make test-func TEST_NAME=TestDeleteCase TEST_PATH=test/integration/case/delete_case_test.go

func TestDeleteCase(t *testing.T) {
	httpClient := &http.Client{Timeout: 10 * time.Second}

	setupClient := func(t *testing.T, ctx context.Context) uuid.UUID {
		c := clientHelpers.NewTestClient(t)
		clientEncx, err := client.ProcessClientEncx(ctx, crypto, c)
		require.NoError(t, err)
		err = clientRepo.CreateClient(ctx, clientEncx)
		require.NoError(t, err)
		return clientEncx.ID
	}

	setupCase := func(t *testing.T, ctx context.Context, clientID uuid.UUID, assignedContactID *uuid.UUID) *caseDomain.Case {
		c := ch.NewTestCase(t)
		c.ClientID = clientID
		if assignedContactID != nil {
			c.AssignedContactID = assignedContactID
		}
		return c
	}

	t.Run("Success", func(t *testing.T) {
		ctx := context.Background()

		// Setup authentication
		adminToken := tu.SetupAdminUser(t, ctx, authCtx)
		defer tu.ClearAuthData(t, ctx, authCtx)
		defer ch.ClearCasesTable(t, ctx, testPool)

		// Create test client and case
		clientID := setupClient(t, ctx)
		c := setupCase(t, ctx, clientID, nil)

		caseEncx, err := caseDomain.ProcessCaseEncx(ctx, crypto, c)
		require.NoError(t, err)
		err = caseRepo.CreateCase(ctx, caseEncx)
		require.NoError(t, err)

		// Create HTTP request using the test helper
		req := ch.NewDeleteCaseRequest(t, ctx, testServerURL, caseEncx.ID, adminToken)

		// Execute request
		resp, err := httpClient.Do(req)
		require.NoError(t, err)
		defer resp.Body.Close()

		// Verify response
		assert.Equal(t, http.StatusNoContent, resp.StatusCode)

		// Verify case is actually deleted
		retrievedCase, err := caseRepo.GetCaseByID(ctx, caseEncx.ID)
		assert.Error(t, err, "Should return error when retrieving deleted case")
		assert.Nil(t, retrievedCase, "Retrieved case should be nil")
		assert.Contains(t, err.Error(), "not found", "Error should indicate case not found")

		t.Log("✓ Case deleted successfully")
	})

	t.Run("SuccessWithAssignedContact", func(t *testing.T) {
		ctx := context.Background()

		// Setup authentication
		adminToken := tu.SetupAdminUser(t, ctx, authCtx)
		defer tu.ClearAuthData(t, ctx, authCtx)
		defer ch.ClearCasesTable(t, ctx, testPool)
		defer clientHelpers.ClearClientsTable(t, ctx, testPool)

		// Create test client with contact
		clientID := setupClient(t, ctx)
		contactID := uuid.New()

		c := setupCase(t, ctx, clientID, &contactID)

		caseEncx, err := caseDomain.ProcessCaseEncx(ctx, crypto, c)
		require.NoError(t, err)
		err = caseRepo.CreateCase(ctx, caseEncx)
		require.NoError(t, err)

		// Create HTTP request
		req := ch.NewDeleteCaseRequest(t, ctx, testServerURL, caseEncx.ID, adminToken)

		// Execute request
		resp, err := httpClient.Do(req)
		require.NoError(t, err)
		defer resp.Body.Close()

		// Verify response
		assert.Equal(t, http.StatusNoContent, resp.StatusCode)

		// Verify case is deleted
		retrievedCase, err := caseRepo.GetCaseByID(ctx, caseEncx.ID)
		assert.Error(t, err, "Should return error when retrieving deleted case")
		assert.Nil(t, retrievedCase, "Retrieved case should be nil")

		t.Log("✓ Case with assigned contact deleted successfully")
	})

	t.Run("CaseNotFound", func(t *testing.T) {
		ctx := context.Background()

		// Setup authentication
		adminToken := tu.SetupAdminUser(t, ctx, authCtx)
		defer tu.ClearAuthData(t, ctx, authCtx)
		defer ch.ClearCasesTable(t, ctx, testPool)

		// Use non-existent case ID
		nonExistentID := uuid.New()

		// Create HTTP request
		req := ch.NewDeleteCaseRequest(t, ctx, testServerURL, nonExistentID, adminToken)

		// Execute request
		resp, err := httpClient.Do(req)
		require.NoError(t, err)
		defer resp.Body.Close()

		// Verify response
		assert.Equal(t, http.StatusNotFound, resp.StatusCode)

		t.Log("✓ Case not found error handled correctly")
	})

	t.Run("InvalidUUID", func(t *testing.T) {
		ctx := context.Background()

		// Setup authentication
		adminToken := tu.SetupAdminUser(t, ctx, authCtx)
		defer tu.ClearAuthData(t, ctx, authCtx)

		// Create HTTP request with invalid UUID (using UUID.Nil)
		req := ch.NewDeleteCaseRequest(t, ctx, testServerURL, uuid.Nil, adminToken)

		// Execute request
		resp, err := httpClient.Do(req)
		require.NoError(t, err)
		defer resp.Body.Close()

		// Verify response - should be 404 since UUID.Nil is valid format but case doesn't exist
		assert.Equal(t, http.StatusNotFound, resp.StatusCode)

		t.Log("✓ Invalid UUID request handled correctly")
	})

	t.Run("Unauthorized", func(t *testing.T) {
		ctx := context.Background()

		defer ch.ClearCasesTable(t, ctx, testPool)

		// Create test case
		clientID := setupClient(t, ctx)
		c := setupCase(t, ctx, clientID, nil)

		caseEncx, err := caseDomain.ProcessCaseEncx(ctx, crypto, c)
		require.NoError(t, err)
		err = caseRepo.CreateCase(ctx, caseEncx)
		require.NoError(t, err)

		// Create HTTP request without authentication
		req := ch.NewDeleteCaseRequest(t, ctx, testServerURL, caseEncx.ID, "")

		// Execute request
		resp, err := httpClient.Do(req)
		require.NoError(t, err)
		defer resp.Body.Close()

		// Verify response
		assert.Equal(t, http.StatusUnauthorized, resp.StatusCode)

		// Verify case still exists
		retrievedCase, err := caseRepo.GetCaseByID(ctx, caseEncx.ID)
		assert.NoError(t, err, "Case should still exist when unauthorized")
		assert.NotNil(t, retrievedCase, "Retrieved case should not be nil")

		t.Log("✓ Unauthorized error handled correctly")
	})

	t.Run("MultipleCases", func(t *testing.T) {
		ctx := context.Background()

		// Setup authentication
		adminToken := tu.SetupAdminUser(t, ctx, authCtx)
		defer tu.ClearAuthData(t, ctx, authCtx)
		defer ch.ClearCasesTable(t, ctx, testPool)

		// Create multiple test cases
		clientID := setupClient(t, ctx)
		case1 := setupCase(t, ctx, clientID, nil)
		case2 := setupCase(t, ctx, clientID, nil)
		case3 := setupCase(t, ctx, clientID, nil)

		case1Encx, err := caseDomain.ProcessCaseEncx(ctx, crypto, case1)
		require.NoError(t, err)
		case2Encx, err := caseDomain.ProcessCaseEncx(ctx, crypto, case2)
		require.NoError(t, err)
		case3Encx, err := caseDomain.ProcessCaseEncx(ctx, crypto, case3)
		require.NoError(t, err)

		err = caseRepo.CreateCase(ctx, case1Encx)
		require.NoError(t, err)
		err = caseRepo.CreateCase(ctx, case2Encx)
		require.NoError(t, err)
		err = caseRepo.CreateCase(ctx, case3Encx)
		require.NoError(t, err)

		// Delete case2
		req := ch.NewDeleteCaseRequest(t, ctx, testServerURL, case2Encx.ID, adminToken)
		resp, err := httpClient.Do(req)
		require.NoError(t, err)
		defer resp.Body.Close()
		assert.Equal(t, http.StatusNoContent, resp.StatusCode)

		// Verify case2 is deleted but case1 and case3 still exist
		retrievedCase1, err := caseRepo.GetCaseByID(ctx, case1Encx.ID)
		assert.NoError(t, err, "Case1 should still exist")
		assert.NotNil(t, retrievedCase1, "Case1 should not be nil")

		retrievedCase2, err := caseRepo.GetCaseByID(ctx, case2Encx.ID)
		assert.Error(t, err, "Case2 should not exist")
		assert.Nil(t, retrievedCase2, "Case2 should be nil")

		retrievedCase3, err := caseRepo.GetCaseByID(ctx, case3Encx.ID)
		assert.NoError(t, err, "Case3 should still exist")
		assert.NotNil(t, retrievedCase3, "Case3 should not be nil")

		t.Log("✓ Multiple cases handled correctly")
	})

	t.Run("EmptyResponse", func(t *testing.T) {
		ctx := context.Background()

		// Setup authentication
		adminToken := tu.SetupAdminUser(t, ctx, authCtx)
		defer tu.ClearAuthData(t, ctx, authCtx)
		defer ch.ClearCasesTable(t, ctx, testPool)

		// Create test case
		clientID := setupClient(t, ctx)
		c := setupCase(t, ctx, clientID, nil)

		caseEncx, err := caseDomain.ProcessCaseEncx(ctx, crypto, c)
		require.NoError(t, err)
		err = caseRepo.CreateCase(ctx, caseEncx)
		require.NoError(t, err)

		// Create HTTP request
		req := ch.NewDeleteCaseRequest(t, ctx, testServerURL, caseEncx.ID, adminToken)

		// Execute request
		resp, err := httpClient.Do(req)
		require.NoError(t, err)
		defer resp.Body.Close()

		// Verify response
		assert.Equal(t, http.StatusNoContent, resp.StatusCode)

		// Verify response body is empty
		body := make([]byte, 1)
		n, err := resp.Body.Read(body)
		assert.Equal(t, 0, n, "Response body should be empty")
		assert.Equal(t, err.Error(), "EOF", "Should reach EOF when reading empty body")

		t.Log("✓ Empty response body handled correctly")
	})
}
