package client_test

import (
	"context"
	"encoding/json"
	"net/http"
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/hengadev/cluo_api/internal/common/errs"
	tu "github.com/hengadev/cluo_api/internal/common/testutils"
	clientDomain "github.com/hengadev/cluo_api/internal/domain/client"
	ch "github.com/hengadev/cluo_api/test/helpers/client"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// make test-func TEST_NAME=TestUpdateClient TEST_PATH=test/integration/client/update_client_test.go

// TestUpdateClient tests all scenarios for the update client endpoint
func TestUpdateClient(t *testing.T) {
	ctx := context.Background()
	httpClient := &http.Client{Timeout: 10 * time.Second}

	t.Run("Success Cases", func(t *testing.T) {
		t.Run("Administrator updates client name successfully", func(t *testing.T) {
			// Setup administrator authentication
			adminToken := tu.SetupAdminUser(t, ctx, authCtx)
			defer tu.ClearAuthData(t, ctx, authCtx)
			defer ch.ClearClientsTable(t, ctx, testPool)

			// Create test client data
			testClient := ch.NewTestClient(t)
			testClient.Name = "Original Name"
			testClient.Type = clientDomain.ClientTypePerson

			// Encrypt and insert client into database
			clientEncx, err := clientDomain.ProcessClientEncx(ctx, crypto, testClient)
			require.NoError(t, err)
			err = ch.InsertClientEncx(t, ctx, testPool, clientEncx)
			require.NoError(t, err)

			// Create update request with new name
			newName := "Updated Name"
			updateRequest := clientDomain.UpdateClientRequest{
				ID:   testClient.ID,
				Name: &newName,
			}

			// Create HTTP request using the test helper
			req := ch.NewUpdateClientRequest(t, ctx, testServerURL, testClient.ID.String(), updateRequest, adminToken)

			// Execute request
			resp, err := httpClient.Do(req)
			require.NoError(t, err)
			defer resp.Body.Close()

			// Verify response
			assert.Equal(t, http.StatusOK, resp.StatusCode)

			// Parse response body
			var response clientDomain.ClientResponse
			err = json.NewDecoder(resp.Body).Decode(&response)
			require.NoError(t, err)

			// Verify response data
			assert.Equal(t, testClient.ID.String(), response.ID)
			assert.Equal(t, newName, response.Name)
			assert.Equal(t, string(testClient.Type), response.Type)
		})

		t.Run("Administrator updates client type successfully", func(t *testing.T) {
			// Setup administrator authentication
			adminToken := tu.SetupAdminUser(t, ctx, authCtx)
			defer tu.ClearAuthData(t, ctx, authCtx)
			defer ch.ClearClientsTable(t, ctx, testPool)

			// Create test client data
			testClient := ch.NewTestClient(t)
			testClient.Name = "Test Client"
			testClient.Type = clientDomain.ClientTypePerson

			// Encrypt and insert client into database
			clientEncx, err := clientDomain.ProcessClientEncx(ctx, crypto, testClient)
			require.NoError(t, err)
			err = ch.InsertClientEncx(t, ctx, testPool, clientEncx)
			require.NoError(t, err)

			// Create update request with new type
			newType := string(clientDomain.ClientTypeInsurance)
			updateRequest := clientDomain.UpdateClientRequest{
				ID:   testClient.ID,
				Type: &newType,
			}

			// Create HTTP request using the test helper
			req := ch.NewUpdateClientRequest(t, ctx, testServerURL, testClient.ID.String(), updateRequest, adminToken)

			// Execute request
			resp, err := httpClient.Do(req)
			require.NoError(t, err)
			defer resp.Body.Close()

			// Verify response
			assert.Equal(t, http.StatusOK, resp.StatusCode)

			// Parse response body
			var response clientDomain.ClientResponse
			err = json.NewDecoder(resp.Body).Decode(&response)
			require.NoError(t, err)

			// Verify response data
			assert.Equal(t, testClient.ID.String(), response.ID)
			assert.Equal(t, testClient.Name, response.Name)
			assert.Equal(t, newType, response.Type)
		})

		t.Run("Administrator updates both name and type successfully", func(t *testing.T) {
			// Setup administrator authentication
			adminToken := tu.SetupAdminUser(t, ctx, authCtx)
			defer tu.ClearAuthData(t, ctx, authCtx)
			defer ch.ClearClientsTable(t, ctx, testPool)

			// Create test client data
			testClient := ch.NewTestClient(t)
			testClient.Name = "Original Name"
			testClient.Type = clientDomain.ClientTypePerson

			// Encrypt and insert client into database
			clientEncx, err := clientDomain.ProcessClientEncx(ctx, crypto, testClient)
			require.NoError(t, err)
			err = ch.InsertClientEncx(t, ctx, testPool, clientEncx)
			require.NoError(t, err)

			// Create update request with both new name and type
			newName := "Completely New Name"
			newType := string(clientDomain.ClientTypeLawyer)
			updateRequest := clientDomain.UpdateClientRequest{
				ID:   testClient.ID,
				Name: &newName,
				Type: &newType,
			}

			// Create HTTP request using the test helper
			req := ch.NewUpdateClientRequest(t, ctx, testServerURL, testClient.ID.String(), updateRequest, adminToken)

			// Execute request
			resp, err := httpClient.Do(req)
			require.NoError(t, err)
			defer resp.Body.Close()

			// Verify response
			assert.Equal(t, http.StatusOK, resp.StatusCode)

			// Parse response body
			var response clientDomain.ClientResponse
			err = json.NewDecoder(resp.Body).Decode(&response)
			require.NoError(t, err)

			// Verify response data
			assert.Equal(t, testClient.ID.String(), response.ID)
			assert.Equal(t, newName, response.Name)
			assert.Equal(t, newType, response.Type)
		})

		t.Run("Administrator updates with no changes (empty update)", func(t *testing.T) {
			// Setup administrator authentication
			adminToken := tu.SetupAdminUser(t, ctx, authCtx)
			defer tu.ClearAuthData(t, ctx, authCtx)
			defer ch.ClearClientsTable(t, ctx, testPool)

			// Create test client data
			testClient := ch.NewTestClient(t)
			testClient.Name = "Test Client"
			testClient.Type = clientDomain.ClientTypePerson

			// Encrypt and insert client into database
			clientEncx, err := clientDomain.ProcessClientEncx(ctx, crypto, testClient)
			require.NoError(t, err)
			err = ch.InsertClientEncx(t, ctx, testPool, clientEncx)
			require.NoError(t, err)

			// Create update request with no changes
			updateRequest := clientDomain.UpdateClientRequest{
				ID: testClient.ID,
			}

			// Create HTTP request using the test helper
			req := ch.NewUpdateClientRequest(t, ctx, testServerURL, testClient.ID.String(), updateRequest, adminToken)

			// Execute request
			resp, err := httpClient.Do(req)
			require.NoError(t, err)
			defer resp.Body.Close()

			// Verify response
			assert.Equal(t, http.StatusOK, resp.StatusCode)

			// Parse response body
			var response clientDomain.ClientResponse
			err = json.NewDecoder(resp.Body).Decode(&response)
			require.NoError(t, err)

			// Verify response data (should be unchanged)
			assert.Equal(t, testClient.ID.String(), response.ID)
			assert.Equal(t, testClient.Name, response.Name)
			assert.Equal(t, string(testClient.Type), response.Type)
		})
	})

	t.Run("Validation Errors", func(t *testing.T) {
		// Setup administrator authentication
		adminToken := tu.SetupAdminUser(t, ctx, authCtx)
		defer tu.ClearAuthData(t, ctx, authCtx)
		defer ch.ClearClientsTable(t, ctx, testPool)

		// Create test client data
		testClient := ch.NewTestClient(t)

		// Encrypt and insert client into database
		clientEncx, err := clientDomain.ProcessClientEncx(ctx, crypto, testClient)
		require.NoError(t, err)
		err = ch.InsertClientEncx(t, ctx, testPool, clientEncx)
		require.NoError(t, err)

		testCases := []struct {
			name     string
			request  clientDomain.UpdateClientRequest
			errorMsg string
		}{
			{
				name: "Empty Name",
				request: clientDomain.UpdateClientRequest{
					ID:   testClient.ID,
					Name: func() *string { s := ""; return &s }(),
				},
				errorMsg: "client name is required",
			},
			{
				name: "Empty Type",
				request: clientDomain.UpdateClientRequest{
					ID:   testClient.ID,
					Type: func() *string { s := ""; return &s }(),
				},
				errorMsg: "client type is required",
			},
			{
				name: "Name Too Long",
				request: clientDomain.UpdateClientRequest{
					ID:   testClient.ID,
					Name: func() *string { s := string(make([]byte, 101)); return &s }(),
				},
				errorMsg: "name must be less than 100 characters",
			},
			{
				name: "Type Too Long",
				request: clientDomain.UpdateClientRequest{
					ID:   testClient.ID,
					Type: func() *string { s := string(make([]byte, 51)); return &s }(),
				},
				errorMsg: "type must be less than 50 characters",
			},
		}

		for _, tc := range testCases {
			t.Run(tc.name, func(t *testing.T) {
				// Create HTTP request using the test helper
				req := ch.NewUpdateClientRequest(t, ctx, testServerURL, testClient.ID.String(), tc.request, adminToken)

				// Execute request
				resp, err := httpClient.Do(req)
				require.NoError(t, err)
				defer resp.Body.Close()

				// Verify bad request response
				assert.Equal(t, http.StatusBadRequest, resp.StatusCode)

				// Parse response body to verify error message
				var response map[string]interface{}
				err = json.NewDecoder(resp.Body).Decode(&response)
				require.NoError(t, err)

				// Check if error message contains expected validation error
				if message, ok := response["error"].(string); ok {
					assert.Contains(t, message, tc.errorMsg)
				}
			})
		}
	})

	t.Run("Error Cases", func(t *testing.T) {
		t.Run("Client not found", func(t *testing.T) {
			// Setup administrator authentication
			adminToken := tu.SetupAdminUser(t, ctx, authCtx)
			defer tu.ClearAuthData(t, ctx, authCtx)
			defer ch.ClearClientsTable(t, ctx, testPool)

			// Use non-existent client ID
			nonExistentID := uuid.New()
			newName := "Updated Name"
			updateRequest := clientDomain.UpdateClientRequest{
				ID:   nonExistentID,
				Name: &newName,
			}

			// Create HTTP request using the test helper
			req := ch.NewUpdateClientRequest(t, ctx, testServerURL, nonExistentID.String(), updateRequest, adminToken)

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
			newName := "Updated Name"
			updateRequest := clientDomain.UpdateClientRequest{
				ID:   uuid.UUID{}, // Empty UUID for invalid ID case
				Name: &newName,
			}

			// Create HTTP request using the test helper
			req := ch.NewUpdateClientRequest(t, ctx, testServerURL, invalidID, updateRequest, adminToken)

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
	})

	t.Run("Authentication/Authorization", func(t *testing.T) {
		defer ch.ClearClientsTable(t, ctx, testPool)

		// Create test client data for authorization tests
		testClient := ch.NewTestClient(t)
		clientEncx, err := clientDomain.ProcessClientEncx(ctx, crypto, testClient)
		require.NoError(t, err)
		err = ch.InsertClientEncx(t, ctx, testPool, clientEncx)
		require.NoError(t, err)

		newName := "Updated Name"
		updateRequest := clientDomain.UpdateClientRequest{
			ID:   testClient.ID,
			Name: &newName,
		}

		t.Run("No Authentication", func(t *testing.T) {
			// Create HTTP request without authentication
			req := ch.NewUpdateClientRequest(t, ctx, testServerURL, testClient.ID.String(), updateRequest, "")

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
			req := ch.NewUpdateClientRequest(t, ctx, testServerURL, testClient.ID.String(), updateRequest, clientToken)

			// Execute request
			resp, err := httpClient.Do(req)
			require.NoError(t, err)
			defer resp.Body.Close()

			// Client should be forbidden from updating clients
			assert.Equal(t, http.StatusForbidden, resp.StatusCode)
		})

		t.Run("Guest Authentication", func(t *testing.T) {
			// Setup guest authentication
			guestToken := tu.SetupGuestUser(t, ctx, authCtx)
			defer tu.ClearAuthData(t, ctx, authCtx)

			// Create HTTP request with guest authentication
			req := ch.NewUpdateClientRequest(t, ctx, testServerURL, testClient.ID.String(), updateRequest, guestToken)

			// Execute request
			resp, err := httpClient.Do(req)
			require.NoError(t, err)
			defer resp.Body.Close()

			// Guest should be forbidden from updating clients
			assert.Equal(t, http.StatusForbidden, resp.StatusCode)
		})
	})
}

