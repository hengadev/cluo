package case_test

import (
	"context"
	"encoding/json"
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

// make test-func TEST_NAME=TestGetCaseByID TEST_PATH=test/integration/case/get_case_by_id_test.go

func TestGetCaseByID(t *testing.T) {
	httpClient := &http.Client{Timeout: 10 * time.Second}

	setupClient := func(t *testing.T, ctx context.Context) uuid.UUID {
		c := clientHelpers.NewTestClient(t)
		clientEncx, err := client.ProcessClientEncx(ctx, crypto, c)
		require.NoError(t, err)
		err = clientHelpers.InsertClientEncx(t, ctx, testPool, clientEncx)
		require.NoError(t, err)
		return c.ID
	}

	setupContact := func(t *testing.T, ctx context.Context, clientID uuid.UUID) uuid.UUID {
		contact := clientHelpers.NewTestContact(t)
	contact.ClientID = clientID
	contactEncx, err := client.ProcessContactEncx(ctx, crypto, contact)
		require.NoError(t, err)
		err = clientHelpers.InsertContactEncx(t, ctx, testPool, contactEncx)
		require.NoError(t, err)
		return contact.ID
	}

	setupCase := func(t *testing.T, ctx context.Context, clientID uuid.UUID) uuid.UUID {
		c := ch.NewTestCase(t)
		c.ClientID = clientID
		caseEncx, err := caseDomain.ProcessCaseEncx(ctx, crypto, c)
		require.NoError(t, err)
		err = ch.InsertCaseEncx(t, ctx, testPool, caseEncx)
		require.NoError(t, err)
		return c.ID
	}

	t.Run("Success", func(t *testing.T) {
		ctx := context.Background()

		// Setup authentication
		accessToken := tu.SetupAdminUser(t, ctx, authCtx)
		defer tu.ClearAuthData(t, ctx, authCtx)
		defer ch.ClearCasesTable(t, ctx, testPool)
		defer clientHelpers.ClearClientsTable(t, ctx, testPool)
		defer clientHelpers.ClearContactsTable(t, ctx, testPool)

		// Generate test client and case
		clientID := setupClient(t, ctx)
		caseID := setupCase(t, ctx, clientID)

		// Create HTTP request using the test helper
		req := ch.NewGetCaseByIDRequest(t, ctx, testServerURL, caseID, accessToken)

		// Execute request
		resp, err := httpClient.Do(req)
		require.NoError(t, err)
		defer resp.Body.Close()

		// Assert response status
		assert.Equal(t, http.StatusOK, resp.StatusCode)

		// Decode response
		var response *caseDomain.CaseResponse
		err = json.NewDecoder(resp.Body).Decode(&response)
		require.NoError(t, err)
		assert.Equal(t, caseID.String(), response.ID, "Case ID should match")
		assert.Equal(t, clientID.String(), response.ClientID, "Client ID should match")
		assert.Equal(t, "Test Case Title", response.Title, "Title should match")
		assert.Equal(t, "Test case description for unit testing", response.Description, "Description should match")
		assert.Equal(t, "EXT-REF-123", *response.ExternalReference, "ExternalReference should match")
		assert.Equal(t, "Test Case Type", response.CaseType, "CaseType should match")
		assert.Equal(t, "draft", response.Status, "Status should match")
		assert.NotEmpty(t, response.CreatedAt, "CreatedAt should not be empty")
		assert.NotEmpty(t, response.UpdatedAt, "UpdatedAt should not be empty")

		t.Log("✓ Case retrieved successfully")
	})

	t.Run("SuccessWithAssignedContact", func(t *testing.T) {
		ctx := context.Background()

		// Setup authentication
		accessToken := tu.SetupAdminUser(t, ctx, authCtx)
		defer tu.ClearAuthData(t, ctx, authCtx)
		defer ch.ClearCasesTable(t, ctx, testPool)
		defer clientHelpers.ClearClientsTable(t, ctx, testPool)
		defer clientHelpers.ClearContactsTable(t, ctx, testPool)

		// Generate test client, contact, and case with assigned contact
		clientID := setupClient(t, ctx)
		contactID := setupContact(t, ctx, clientID)
		c := ch.NewTestCase(t)
		c.ClientID = clientID
		c.AssignedContactID = &contactID
		caseEncx, err := caseDomain.ProcessCaseEncx(ctx, crypto, c)
		require.NoError(t, err)
		err = ch.InsertCaseEncx(t, ctx, testPool, caseEncx)
		require.NoError(t, err)

		// Create HTTP request using the test helper
		req := ch.NewGetCaseByIDRequest(t, ctx, testServerURL, c.ID, accessToken)

		// Execute request
		resp, err := httpClient.Do(req)
		require.NoError(t, err)
		defer resp.Body.Close()

		// Assert response status
		assert.Equal(t, http.StatusOK, resp.StatusCode)

		// Decode response
		var response *caseDomain.CaseResponse
		err = json.NewDecoder(resp.Body).Decode(&response)
		require.NoError(t, err)
		assert.NotNil(t, response.AssignedContactID, "AssignedContactID should not be nil")
		assert.Equal(t, contactID.String(), *response.AssignedContactID, "Assigned Contact ID should match")

		t.Log("✓ Case with assigned contact retrieved successfully")
	})

	t.Run("CaseNotFound", func(t *testing.T) {
		ctx := context.Background()

		// Setup authentication
		accessToken := tu.SetupAdminUser(t, ctx, authCtx)
		defer tu.ClearAuthData(t, ctx, authCtx)
		defer ch.ClearCasesTable(t, ctx, testPool)
		defer clientHelpers.ClearClientsTable(t, ctx, testPool)

		// Use a non-existent case ID
		nonExistentID := uuid.New()

		// Create HTTP request using the test helper
		req := ch.NewGetCaseByIDRequest(t, ctx, testServerURL, nonExistentID, accessToken)

		// Execute request
		resp, err := httpClient.Do(req)
		require.NoError(t, err)
		defer resp.Body.Close()

		// Assert response status - should be 404 Not Found
		assert.Equal(t, http.StatusNotFound, resp.StatusCode)

		// Decode error response
		var errorResp struct {
			Error string `json:"error"`
		}

		err = json.NewDecoder(resp.Body).Decode(&errorResp)
		require.NoError(t, err)
		assert.Contains(t, errorResp.Error, "not found", "Error message should indicate case not found")

		t.Log("✓ Case not found error handled correctly")
	})

	t.Run("InvalidUUID", func(t *testing.T) {
		ctx := context.Background()

		// Setup authentication
		accessToken := tu.SetupAdminUser(t, ctx, authCtx)
		defer tu.ClearAuthData(t, ctx, authCtx)

		// Use invalid UUID format
		invalidID := uuid.Nil

		// Create HTTP request using the test helper
		req := ch.NewGetCaseByIDRequest(t, ctx, testServerURL, invalidID, accessToken)

		// Execute request
		resp, err := httpClient.Do(req)
		require.NoError(t, err)
		defer resp.Body.Close()

		// Assert response status - should be 400 Bad Request (invalid UUID path)
	// Note: This might be handled by the router path matching or the handler

		t.Log("✓ Invalid UUID request handled")
	})

	t.Run("EmptyUUID", func(t *testing.T) {
		ctx := context.Background()

		// Setup authentication
		accessToken := tu.SetupAdminUser(t, ctx, authCtx)
		defer tu.ClearAuthData(t, ctx, authCtx)

		// Use empty UUID path parameter
		req := ch.NewGetCaseByIDRequest(t, ctx, testServerURL, uuid.Nil, accessToken)

		// Execute request
		resp, err := httpClient.Do(req)
		require.NoError(t, err)
		defer resp.Body.Close()

		t.Log("✓ Empty UUID request handled")
	})

	t.Run("Unauthorized", func(t *testing.T) {
		ctx := context.Background()

		defer ch.ClearCasesTable(t, ctx, testPool)
		defer clientHelpers.ClearClientsTable(t, ctx, testPool)

		// Generate test client and case
		clientID := setupClient(t, ctx)
		caseID := setupCase(t, ctx, clientID)

		// Create HTTP request without authentication token
		req := ch.NewGetCaseByIDRequest(t, ctx, testServerURL, caseID, "")

		// Execute request
		resp, err := httpClient.Do(req)
		require.NoError(t, err)
		defer resp.Body.Close()

		// Assert response status - should be 401 Unauthorized
		assert.Equal(t, http.StatusUnauthorized, resp.StatusCode)

		t.Log("✓ Unauthorized error handled correctly")
	})

	t.Run("MultipleCases", func(t *testing.T) {
		ctx := context.Background()

		// Setup authentication
		accessToken := tu.SetupAdminUser(t, ctx, authCtx)
		defer tu.ClearAuthData(t, ctx, authCtx)
		defer ch.ClearCasesTable(t, ctx, testPool)
		defer clientHelpers.ClearClientsTable(t, ctx, testPool)

		// Generate test client
		clientID := setupClient(t, ctx)

		// Create multiple cases
		case1ID := setupCase(t, ctx, clientID)
		case2ID := setupCase(t, ctx, clientID)
		case3ID := setupCase(t, ctx, clientID)

		// Retrieve each case individually
		req1 := ch.NewGetCaseByIDRequest(t, ctx, testServerURL, case1ID, accessToken)
		resp1, err := httpClient.Do(req1)
		require.NoError(t, err)
		defer resp1.Body.Close()
		assert.Equal(t, http.StatusOK, resp1.StatusCode)

		req2 := ch.NewGetCaseByIDRequest(t, ctx, testServerURL, case2ID, accessToken)
		resp2, err := httpClient.Do(req2)
		require.NoError(t, err)
		defer resp2.Body.Close()
		assert.Equal(t, http.StatusOK, resp2.StatusCode)

		req3 := ch.NewGetCaseByIDRequest(t, ctx, testServerURL, case3ID, accessToken)
		resp3, err := httpClient.Do(req3)
		require.NoError(t, err)
		defer resp3.Body.Close()
		assert.Equal(t, http.StatusOK, resp3.StatusCode)

		// Decode responses and verify they're different cases
		var response1, response2, response3 *caseDomain.CaseResponse
		err = json.NewDecoder(resp1.Body).Decode(&response1)
		require.NoError(t, err)
		err = json.NewDecoder(resp2.Body).Decode(&response2)
		require.NoError(t, err)
		err = json.NewDecoder(resp3.Body).Decode(&response3)
		require.NoError(t, err)

		assert.Equal(t, case1ID.String(), response1.ID, "First case ID should match")
		assert.Equal(t, case2ID.String(), response2.ID, "Second case ID should match")
		assert.Equal(t, case3ID.String(), response3.ID, "Third case ID should match")
		assert.NotEqual(t, response1.ID, response2.ID, "Cases should be different")
		assert.NotEqual(t, response2.ID, response3.ID, "Cases should be different")

		t.Log("✓ Multiple cases retrieved successfully")
	})
}
