package client_test

import (
	"context"
	"encoding/json"
	"net/http"
	"testing"
	"time"

	"github.com/hengadev/cluo_api/internal/common/errs"
	tu "github.com/hengadev/cluo_api/internal/common/testutils"
	clientDomain "github.com/hengadev/cluo_api/internal/domain/client"
	ch "github.com/hengadev/cluo_api/test/helpers/client"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// make test-func TEST_NAME=TestDeleteClient TEST_PATH=test/integration/client/delete_client_test.go

// TestDeleteClient tests all scenarios for the delete client endpoint
func TestDeleteClient(t *testing.T) {
	ctx := context.Background()
	httpClient := &http.Client{Timeout: 10 * time.Second}

	t.Run("Success Cases", func(t *testing.T) {
		t.Run("Administrator deletes any client successfully", func(t *testing.T) {
			// Setup administrator authentication
			adminToken := tu.SetupAdminUser(t, ctx, authCtx)
			defer tu.ClearAuthData(t, ctx, authCtx)
			defer ch.ClearClientsTable(t, ctx, testPool)

			// Create test client data
			testClient := ch.NewTestClient(t)

			testClientEncx, err := clientDomain.ProcessClientEncx(ctx, crypto, testClient)
			require.NoError(t, err)

			// Insert client into database
			err = ch.InsertClientEncx(t, ctx, testPool, testClientEncx)
			require.NoError(t, err)

			// Verify client exists before deletion
			_, err = ch.GetClientEncxByID(t, ctx, testPool, testClient.ID)
			require.NoError(t, err)

			// Create HTTP request using the test helper
			req := ch.NewDeleteClientRequest(t, ctx, testServerURL, testClient.ID.String(), adminToken)

			// Execute request
			resp, err := httpClient.Do(req)
			require.NoError(t, err)
			defer resp.Body.Close()

			// Verify response
			assert.Equal(t, http.StatusOK, resp.StatusCode)

			// Parse response body
			var response map[string]interface{}
			err = json.NewDecoder(resp.Body).Decode(&response)
			require.NoError(t, err)

			// Verify response contains success message
			if message, ok := response["message"].(string); ok {
				assert.Contains(t, message, "successfully")
			}

			// Verify client no longer exists in database
			_, err = ch.GetClientEncxByID(t, ctx, testPool, testClient.ID)
			assert.Error(t, err) // Should return error when client not found
		})
	})

	t.Run("Error Cases", func(t *testing.T) {
		t.Run("Client not found", func(t *testing.T) {
			// Setup administrator authentication
			adminToken := tu.SetupAdminUser(t, ctx, authCtx)
			defer tu.ClearAuthData(t, ctx, authCtx)
			defer ch.ClearClientsTable(t, ctx, testPool)

			// Use non-existent client ID
			nonExistentID := uuid.New()

			// Create HTTP request using the test helper
			req := ch.NewDeleteClientRequest(t, ctx, testServerURL, nonExistentID.String(), adminToken)

			// Execute request
			resp, err := httpClient.Do(req)
			require.NoError(t, err)
			defer resp.Body.Close()

			// Verify response
			assert.Equal(t, http.StatusNotFound, resp.StatusCode)

			// Parse response body to verify error message
			var response map[string]interface{}
			err = json.NewDecoder(resp.Body).Decode(&response)
			require.NoError(t, err)

			// Check if error message contains expected error
			if message, ok := response["error"].(string); ok {
				assert.Contains(t, message, errs.ErrRepositoryNotFound.Error())
			}
		})

		t.Run("Invalid client ID format", func(t *testing.T) {
			// Setup administrator authentication
			adminToken := tu.SetupAdminUser(t, ctx, authCtx)
			defer tu.ClearAuthData(t, ctx, authCtx)

			// Use invalid UUID format
			invalidID := "invalid-uuid-format"

			// Create HTTP request using the test helper
			req := ch.NewDeleteClientRequest(t, ctx, testServerURL, invalidID, adminToken)

			// Execute request
			resp, err := httpClient.Do(req)
			require.NoError(t, err)
			defer resp.Body.Close()

			// Verify response - should return bad request for invalid UUID
			assert.Equal(t, http.StatusBadRequest, resp.StatusCode)

			// Parse response body to verify error message
			var response map[string]interface{}
			err = json.NewDecoder(resp.Body).Decode(&response)
			require.NoError(t, err)

			// Check if error message contains expected validation error
			if message, ok := response["error"].(string); ok {
				assert.Contains(t, message, errs.ErrInvalidValue.Error())
			}
		})

		t.Run("Client already deleted", func(t *testing.T) {
			// Setup administrator authentication
			adminToken := tu.SetupAdminUser(t, ctx, authCtx)
			defer tu.ClearAuthData(t, ctx, authCtx)
			defer ch.ClearClientsTable(t, ctx, testPool)

			// Create test client data
			testClient := ch.NewTestClient(t)

			testClientEncx, err := clientDomain.ProcessClientEncx(ctx, crypto, testClient)
			require.NoError(t, err)

			// Insert client into database
			err = ch.InsertClientEncx(t, ctx, testPool, testClientEncx)
			require.NoError(t, err)

			// First delete request
			req := ch.NewDeleteClientRequest(t, ctx, testServerURL, testClient.ID.String(), adminToken)
			resp, err := httpClient.Do(req)
			require.NoError(t, err)
			defer resp.Body.Close()
			assert.Equal(t, http.StatusOK, resp.StatusCode)

			// Second delete request (client already deleted)
			req2 := ch.NewDeleteClientRequest(t, ctx, testServerURL, testClient.ID.String(), adminToken)
			resp2, err := httpClient.Do(req2)
			require.NoError(t, err)
			defer resp2.Body.Close()

			// Should return not found for already deleted client
			assert.Equal(t, http.StatusNotFound, resp2.StatusCode)

			// Parse response body to verify error message
			var response map[string]interface{}
			err = json.NewDecoder(resp2.Body).Decode(&response)
			require.NoError(t, err)

			// Check if error message contains expected error
			if message, ok := response["error"].(string); ok {
				assert.Contains(t, message, errs.ErrRepositoryNotFound.Error())
			}
		})
	})

	t.Run("Authentication/Authorization", func(t *testing.T) {
		defer ch.ClearClientsTable(t, ctx, testPool)

		// Create test client data for authorization tests
		testClient := ch.NewTestClient(t)
		testClientEncx, err := clientDomain.ProcessClientEncx(ctx, crypto, testClient)
		require.NoError(t, err)
		err = ch.InsertClientEncx(t, ctx, testPool, testClientEncx)
		require.NoError(t, err)

		t.Run("No Authentication", func(t *testing.T) {
			// Create HTTP request without authentication
			req := ch.NewDeleteClientRequest(t, ctx, testServerURL, testClient.ID.String(), "")

			// Execute request
			resp, err := httpClient.Do(req)
			require.NoError(t, err)
			defer resp.Body.Close()

			// Should be unauthorized
			assert.Equal(t, http.StatusUnauthorized, resp.StatusCode)
		})

		t.Run("Client Authentication", func(t *testing.T) {
			// Setup client authentication
			clientToken := tu.SetupClientUser(t, ctx, authCtx)
			defer tu.ClearAuthData(t, ctx, authCtx)

			// Create HTTP request with client authentication
			req := ch.NewDeleteClientRequest(t, ctx, testServerURL, testClient.ID.String(), clientToken)

			// Execute request
			resp, err := httpClient.Do(req)
			require.NoError(t, err)
			defer resp.Body.Close()

			// Client should be forbidden from deleting clients
			assert.Equal(t, http.StatusForbidden, resp.StatusCode)
		})

		t.Run("Guest Authentication", func(t *testing.T) {
			// Setup guest authentication
			guestToken := tu.SetupGuestUser(t, ctx, authCtx)
			defer tu.ClearAuthData(t, ctx, authCtx)

			// Create HTTP request with guest authentication
			req := ch.NewDeleteClientRequest(t, ctx, testServerURL, testClient.ID.String(), guestToken)

			// Execute request
			resp, err := httpClient.Do(req)
			require.NoError(t, err)
			defer resp.Body.Close()

			// Guest should be forbidden from deleting clients
			assert.Equal(t, http.StatusForbidden, resp.StatusCode)
		})
	})
}
