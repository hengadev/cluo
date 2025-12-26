package case_test

import (
	"bytes"
	"context"
	"encoding/json"
	"net/http"
	"testing"
	"time"

	"github.com/hengadev/cluo_api/internal/common/auth/cookies"
	"github.com/hengadev/cluo_api/internal/common/errs"
	tu "github.com/hengadev/cluo_api/internal/common/testutils"
	caseDomain "github.com/hengadev/cluo_api/internal/domain/case"
	"github.com/hengadev/cluo_api/internal/domain/client"
	ch "github.com/hengadev/cluo_api/test/helpers/case"
	clientHelpers "github.com/hengadev/cluo_api/test/helpers/client"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// make test-func TEST_NAME=TestCreateCase TEST_PATH=test/integration/case/create_case_test.go

// TestCreateCase tests all scenarios for creating a case
func TestCreateCase(t *testing.T) {
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

	t.Run("Success", func(t *testing.T) {
		ctx := context.Background()

		// Setup authentication
		accessToken := tu.SetupAdminUser(t, ctx, authCtx)
		defer tu.ClearAuthData(t, ctx, authCtx)
		defer ch.ClearCasesTable(t, ctx, testPool)
		defer clientHelpers.ClearClientsTable(t, ctx, testPool)

		// Generate test client
		clientID := setupClient(t, ctx)

		// Prepare request payload
		externalRef := "EXT-REF-123"
		payload := caseDomain.CreateCaseRequest{
			Title:             "Test Case Title",
			Description:       "Test case description for integration testing",
			ClientID:          clientID.String(),
			ExternalReference: &externalRef,
			CaseType:          "Theft",
			Status:            "draft",
		}

		// Create HTTP request using the test helper
		req := ch.NewCreateCaseRequest(t, ctx, testServerURL, payload, accessToken)

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
		assert.Equal(t, payload.Title, response.Title)
		assert.Equal(t, payload.Description, response.Description)
		assert.Equal(t, payload.ClientID, response.ClientID)
		assert.Equal(t, payload.ExternalReference, response.ExternalReference)
		assert.Equal(t, payload.CaseType, response.CaseType)
		assert.Equal(t, payload.Status, response.Status)
		assert.NotEmpty(t, response.ID)
		assert.NotEmpty(t, response.CreatedAt)
		assert.NotEmpty(t, response.UpdatedAt)

		t.Log("✓ Case created successfully")
	})

	t.Run("SuccessWithAssignedContact", func(t *testing.T) {
		ctx := context.Background()

		// Setup authentication
		accessToken := tu.SetupAdminUser(t, ctx, authCtx)
		defer tu.ClearAuthData(t, ctx, authCtx)
		defer ch.ClearCasesTable(t, ctx, testPool)
		defer clientHelpers.ClearClientsTable(t, ctx, testPool)
		defer clientHelpers.ClearContactsTable(t, ctx, testPool)

		// Generate test client and contact
		clientID := setupClient(t, ctx)
		contactID := setupContact(t, ctx, clientID)
		contactIDStr := contactID.String()

		// Prepare request payload with assigned contact
		payload := caseDomain.CreateCaseRequest{
			Title:             "Test Case with Contact",
			Description:       "Test case with assigned contact",
			ClientID:          clientID.String(),
			AssignedContactID: &contactIDStr,
			CaseType:          "Accident",
			Status:            "draft",
		}

		// Create HTTP request using the test helper
		req := ch.NewCreateCaseRequest(t, ctx, testServerURL, payload, accessToken)

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
		assert.Equal(t, payload.Title, response.Title)
		assert.Equal(t, payload.Description, response.Description)
		assert.Equal(t, payload.ClientID, response.ClientID)
		assert.Equal(t, payload.AssignedContactID, response.AssignedContactID)
		assert.Equal(t, payload.CaseType, response.CaseType)
		assert.Equal(t, payload.Status, response.Status)

		t.Log("✓ Case created successfully with assigned contact")
	})

	t.Run("ClientNotFound", func(t *testing.T) {
		ctx := context.Background()

		// Setup authentication
		accessToken := tu.SetupAdminUser(t, ctx, authCtx)
		defer tu.ClearAuthData(t, ctx, authCtx)
		defer ch.ClearCasesTable(t, ctx, testPool)
		defer clientHelpers.ClearClientsTable(t, ctx, testPool)

		// Use a non-existent client ID
		nonExistentClientID := uuid.New()

		// Prepare request payload
		payload := caseDomain.CreateCaseRequest{
			Title:       "Test Case",
			Description: "Test case description",
			ClientID:    nonExistentClientID.String(),
			CaseType:    "Theft",
			Status:      "draft",
		}

		// Create HTTP request using the test helper
		req := ch.NewCreateCaseRequest(t, ctx, testServerURL, payload, accessToken)

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
		assert.Contains(t, errorResp.Error, errs.ErrRepositoryNotFound.Error())

		t.Log("✓ Client not found error handled correctly")
	})

	t.Run("ContactNotFound", func(t *testing.T) {
		ctx := context.Background()

		// Setup authentication
		accessToken := tu.SetupAdminUser(t, ctx, authCtx)
		defer tu.ClearAuthData(t, ctx, authCtx)
		defer ch.ClearCasesTable(t, ctx, testPool)
		defer clientHelpers.ClearClientsTable(t, ctx, testPool)

		// Generate test client
		clientID := setupClient(t, ctx)

		// Use a non-existent contact ID
		nonExistentContactID := uuid.New()
		contactIDStr := nonExistentContactID.String()

		// Prepare request payload with non-existent contact
		payload := caseDomain.CreateCaseRequest{
			Title:             "Test Case",
			Description:       "Test case description",
			ClientID:          clientID.String(),
			AssignedContactID: &contactIDStr,
			CaseType:          "Theft",
			Status:            "draft",
		}

		// Create HTTP request using the test helper
		req := ch.NewCreateCaseRequest(t, ctx, testServerURL, payload, accessToken)

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
		assert.Contains(t, errorResp.Error, errs.ErrRepositoryNotFound.Error())

		t.Log("✓ Contact not found error handled correctly")
	})

	t.Run("Unauthorized", func(t *testing.T) {
		ctx := context.Background()

		defer ch.ClearCasesTable(t, ctx, testPool)
		defer clientHelpers.ClearClientsTable(t, ctx, testPool)

		// Generate test client
		clientID := setupClient(t, ctx)

		// Prepare request payload
		payload := caseDomain.CreateCaseRequest{
			Title:       "Test Case",
			Description: "Test case description",
			ClientID:    clientID.String(),
			CaseType:    "Theft",
			Status:      "draft",
		}

		// Create HTTP request without authentication token
		req := ch.NewCreateCaseRequest(t, ctx, testServerURL, payload, "")

		// Execute request
		resp, err := httpClient.Do(req)
		require.NoError(t, err)
		defer resp.Body.Close()

		// Assert response status - should be 401 Unauthorized
		assert.Equal(t, http.StatusUnauthorized, resp.StatusCode)

		t.Log("✓ Unauthorized error handled correctly")
	})

	t.Run("InvalidClientID", func(t *testing.T) {
		ctx := context.Background()

		// Setup authentication
		accessToken := tu.SetupAdminUser(t, ctx, authCtx)
		defer tu.ClearAuthData(t, ctx, authCtx)
		defer ch.ClearCasesTable(t, ctx, testPool)

		// Prepare request payload with invalid client ID
		payload := caseDomain.CreateCaseRequest{
			Title:       "Test Case",
			Description: "Test case description",
			ClientID:    "invalid-uuid-format",
			CaseType:    "Theft",
			Status:      "draft",
		}

		// Create HTTP request using the test helper
		req := ch.NewCreateCaseRequest(t, ctx, testServerURL, payload, accessToken)

		// Execute request
		resp, err := httpClient.Do(req)
		require.NoError(t, err)
		defer resp.Body.Close()

		// Assert response status - should be 400 Bad Request
		assert.Equal(t, http.StatusBadRequest, resp.StatusCode)

		t.Log("✓ Invalid client ID error handled correctly")
	})

	t.Run("InvalidPayload", func(t *testing.T) {
		ctx := context.Background()

		// Setup authentication
		accessToken := tu.SetupAdminUser(t, ctx, authCtx)
		defer tu.ClearAuthData(t, ctx, authCtx)
		defer ch.ClearCasesTable(t, ctx, testPool)

		// Create HTTP request with invalid JSON
		invalidJSON := `{
			"title": "",
			"description": "Test case description",
			"client_id": "invalid-uuid",
			"status": "invalid-status"
		}`

		req, err := http.NewRequestWithContext(
			ctx,
			http.MethodPost,
			testServerURL+"/cases",
			bytes.NewReader([]byte(invalidJSON)),
		)
		require.NoError(t, err)
		req.Header.Set("Content-Type", "application/json")

		if accessToken != "" {
			cookie := &http.Cookie{
				Name:  cookies.AccessTokenCookieName,
				Value: accessToken,
			}
			req.AddCookie(cookie)
		}

		// Execute request
		resp, err := httpClient.Do(req)
		require.NoError(t, err)
		defer resp.Body.Close()

		// Assert response status - should be 400 Bad Request
		assert.Equal(t, http.StatusBadRequest, resp.StatusCode)

		t.Log("✓ Invalid payload error handled correctly")
	})
}

