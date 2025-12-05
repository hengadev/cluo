package client_test

import (
	"context"
	"encoding/json"
	"net/http"
	"testing"
	"time"

	tu "github.com/hengadev/cluo_api/internal/common/testutils"
	clientDomain "github.com/hengadev/cluo_api/internal/domain/client"
	ch "github.com/hengadev/cluo_api/test/helpers/client"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// make test-func TEST_NAME=TestGetAllClients TEST_PATH=test/integration/client/get_all_clients_test.go

// TestGetAllClients tests all scenarios for the get all clients endpoint
func TestGetAllClients(t *testing.T) {
	ctx := context.Background()
	httpClient := &http.Client{Timeout: 10 * time.Second}

	t.Run("Success Cases", func(t *testing.T) {
		t.Run("Administrator retrieves all clients successfully", func(t *testing.T) {
			// Setup administrator authentication
			adminToken := tu.SetupAdminUser(t, ctx, authCtx)
			defer tu.ClearAuthData(t, ctx, authCtx)
			defer ch.ClearClientsTable(t, ctx, testPool)

			// Create multiple test clients
			testClients := []*clientDomain.Client{
				ch.NewTestClient(t),
				ch.NewTestClient(t),
				ch.NewTestClient(t),
			}

			// Set unique names for each client
			testClients[0].Name = "Client One"
			testClients[1].Name = "Client Two"
			testClients[2].Name = "Client Three"

			// Encrypt and insert clients into database
			for _, client := range testClients {
				clientEncx, err := clientDomain.ProcessClientEncx(ctx, crypto, client)
				require.NoError(t, err)
				err = ch.InsertClientEncx(t, ctx, testPool, clientEncx)
				require.NoError(t, err)
			}

			// Create HTTP request using the test helper
			req := ch.NewGetAllClientsRequest(t, ctx, testServerURL, adminToken)

			// Execute request
			resp, err := httpClient.Do(req)
			require.NoError(t, err)
			defer resp.Body.Close()

			// Verify response
			assert.Equal(t, http.StatusOK, resp.StatusCode)

			// Parse response body
			var response []*clientDomain.ClientResponse
			err = json.NewDecoder(resp.Body).Decode(&response)
			require.NoError(t, err)

			// Verify response contains all clients
			assert.Len(t, response, 3)

			// Verify client data (order may vary)
			responseNames := make(map[string]bool)
			for _, clientResp := range response {
				responseNames[clientResp.Name] = true
				assert.NotEmpty(t, clientResp.ID)
				assert.NotEmpty(t, clientResp.Type)
			}

			// Verify all our test clients are present
			assert.True(t, responseNames["Client One"])
			assert.True(t, responseNames["Client Two"])
			assert.True(t, responseNames["Client Three"])
		})

		t.Run("Administrator retrieves empty client list successfully", func(t *testing.T) {
			// Setup administrator authentication
			adminToken := tu.SetupAdminUser(t, ctx, authCtx)
			defer tu.ClearAuthData(t, ctx, authCtx)
			defer ch.ClearClientsTable(t, ctx, testPool)

			// Create HTTP request using the test helper
			req := ch.NewGetAllClientsRequest(t, ctx, testServerURL, adminToken)

			// Execute request
			resp, err := httpClient.Do(req)
			require.NoError(t, err)
			defer resp.Body.Close()

			// Verify response
			assert.Equal(t, http.StatusOK, resp.StatusCode)

			// Parse response body
			var response []clientDomain.ClientResponse
			err = json.NewDecoder(resp.Body).Decode(&response)
			require.NoError(t, err)

			// Verify empty list
			assert.Empty(t, response)
		})

		t.Run("Administrator retrieves single client successfully", func(t *testing.T) {
			// Setup administrator authentication
			adminToken := tu.SetupAdminUser(t, ctx, authCtx)
			defer tu.ClearAuthData(t, ctx, authCtx)
			defer ch.ClearClientsTable(t, ctx, testPool)

			// Create single test client
			testClient := ch.NewTestClient(t)
			testClient.Name = "Single Client"

			// Encrypt and insert client into database
			clientEncx, err := clientDomain.ProcessClientEncx(ctx, crypto, testClient)
			require.NoError(t, err)
			err = ch.InsertClientEncx(t, ctx, testPool, clientEncx)
			require.NoError(t, err)

			// Create HTTP request using the test helper
			req := ch.NewGetAllClientsRequest(t, ctx, testServerURL, adminToken)

			// Execute request
			resp, err := httpClient.Do(req)
			require.NoError(t, err)
			defer resp.Body.Close()

			// Verify response
			assert.Equal(t, http.StatusOK, resp.StatusCode)

			// Parse response body
			var response []clientDomain.ClientResponse
			err = json.NewDecoder(resp.Body).Decode(&response)
			require.NoError(t, err)

			// Verify single client
			assert.Len(t, response, 1)
			assert.Equal(t, "Single Client", response[0].Name)
			assert.Equal(t, testClient.ID.String(), response[0].ID)
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
			req := ch.NewGetAllClientsRequest(t, ctx, testServerURL, "")

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
			req := ch.NewGetAllClientsRequest(t, ctx, testServerURL, clientToken)

			// Execute request
			resp, err := httpClient.Do(req)
			require.NoError(t, err)
			defer resp.Body.Close()

			// Client should be forbidden from accessing all clients
			assert.Equal(t, http.StatusForbidden, resp.StatusCode)
		})

		t.Run("Guest Authentication", func(t *testing.T) {
			// Setup guest authentication
			guestToken := tu.SetupGuestUser(t, ctx, authCtx)
			defer tu.ClearAuthData(t, ctx, authCtx)

			// Create HTTP request with guest authentication
			req := ch.NewGetAllClientsRequest(t, ctx, testServerURL, guestToken)

			// Execute request
			resp, err := httpClient.Do(req)
			require.NoError(t, err)
			defer resp.Body.Close()

			// Guest should be forbidden from accessing all clients
			assert.Equal(t, http.StatusForbidden, resp.StatusCode)
		})
	})
}

