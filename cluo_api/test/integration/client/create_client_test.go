package client_test

import (
	"bytes"
	"context"
	"encoding/json"
	"net/http"
	"testing"
	"time"

	"github.com/hengadev/cluo_api/internal/common/auth/cookies"
	tu "github.com/hengadev/cluo_api/internal/common/testutils"
	clientDomain "github.com/hengadev/cluo_api/internal/domain/client"
	"github.com/hengadev/cluo_api/internal/interface/client"
	ch "github.com/hengadev/cluo_api/test/helpers/client"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// make test-func TEST_NAME=TestCreateClient TEST_PATH=test/integration/client/create_client_test.go

// TestCreateClient tests all scenarios for the create client endpoint
func TestCreateClient(t *testing.T) {
	ctx := context.Background()
	httpClient := &http.Client{Timeout: 10 * time.Second}

	t.Run("Success Cases", func(t *testing.T) {
		t.Run("Administrator creates client successfully", func(t *testing.T) {
			// Setup administrator authentication
			adminToken := tu.SetupAdminUser(t, ctx, authCtx)
			defer tu.ClearAuthData(t, ctx, authCtx)
			defer ch.ClearClientsTable(t, ctx, testPool)

			clientName := "Admin Test Client"
			clientType := "Business"

			// Create a valid client request
			request := clientDomain.CreateClientRequest{
				Name: clientName,
				Type: clientType,
			}

			// Create HTTP request using the test helper
			req := ch.NewCreateClientRequest(t, ctx, testServerURL, request, adminToken)

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

			// Verify response message
			assert.Equal(t, response.Name, clientName)
			assert.Equal(t, response.Type, clientType)
		})
	})

	t.Run("Validation Errors", func(t *testing.T) {
		// Setup administrator authentication
		adminToken := tu.SetupAdminUser(t, ctx, authCtx)
		defer tu.ClearAuthData(t, ctx, authCtx)
		defer ch.ClearClientsTable(t, ctx, testPool)

		testCases := []struct {
			name     string
			request  clientDomain.CreateClientRequest
			errorMsg string
		}{
			{
				name: "Empty Name",
				request: clientDomain.CreateClientRequest{
					Name: "",
					Type: "Individual",
				},
				errorMsg: "name is required",
			},
			{
				name: "Empty Type",
				request: clientDomain.CreateClientRequest{
					Name: "Test Client",
					Type: "",
				},
				errorMsg: "type is required",
			},
			{
				name: "Name Too Long",
				request: clientDomain.CreateClientRequest{
					Name: string(make([]byte, 101)), // 101 characters
					Type: "Individual",
				},
				errorMsg: "name must be less than 100 characters",
			},
			{
				name: "Type Too Long",
				request: clientDomain.CreateClientRequest{
					Name: "Test Client",
					Type: string(make([]byte, 51)), // 51 characters
				},
				errorMsg: "type must be less than 50 characters",
			},
		}

		for _, tc := range testCases {
			t.Run(tc.name, func(t *testing.T) {
				// Create HTTP request using the test helper
				req := ch.NewCreateClientRequest(t, ctx, testServerURL, tc.request, adminToken)

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

	t.Run("Authentication/Authorization", func(t *testing.T) {
		defer ch.ClearClientsTable(t, ctx, testPool)

		// Create a valid client request
		request := clientDomain.CreateClientRequest{
			Name: "Test Client",
			Type: "Individual",
		}

		t.Run("No Authentication", func(t *testing.T) {
			// Create HTTP request without authentication
			req := ch.NewCreateClientRequest(t, ctx, testServerURL, request, "")

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
			req := ch.NewCreateClientRequest(t, ctx, testServerURL, request, clientToken)

			// Execute request
			resp, err := httpClient.Do(req)
			require.NoError(t, err)
			defer resp.Body.Close()

			// Client should be forbidden from creating clients
			assert.Equal(t, http.StatusForbidden, resp.StatusCode)
		})

		t.Run("Guest Authentication", func(t *testing.T) {
			// Setup guest authentication
			guestToken := tu.SetupGuestUser(t, ctx, authCtx)
			defer tu.ClearAuthData(t, ctx, authCtx)

			// Create HTTP request with guest authentication
			req := ch.NewCreateClientRequest(t, ctx, testServerURL, request, guestToken)

			// Execute request
			resp, err := httpClient.Do(req)
			require.NoError(t, err)
			defer resp.Body.Close()

			// Guest should be forbidden from creating clients
			assert.Equal(t, http.StatusForbidden, resp.StatusCode)
		})
	})

	t.Run("Invalid JSON Payloads", func(t *testing.T) {
		// Setup administrator authentication
		adminToken := tu.SetupAdminUser(t, ctx, authCtx)
		defer tu.ClearAuthData(t, ctx, authCtx)
		defer ch.ClearClientsTable(t, ctx, testPool)

		t.Run("Invalid JSON Syntax", func(t *testing.T) {
			// Create request with invalid JSON
			invalidJSON := `{"id": "invalid-uuid", "name": "test", "type": "individual"` // Missing closing brace

			req, err := http.NewRequestWithContext(
				ctx,
				http.MethodPost,
				testServerURL+clientHandler.CreateClientEndpoint,
				bytes.NewReader([]byte(invalidJSON)),
			)
			require.NoError(t, err)

			req.Header.Set("Content-Type", "application/json")

			// Add authentication cookie manually
			cookie := &http.Cookie{
				Name:  cookies.AccessTokenCookieName,
				Value: adminToken,
			}
			req.AddCookie(cookie)

			// Execute request
			resp, err := httpClient.Do(req)
			require.NoError(t, err)
			defer resp.Body.Close()

			// Should return bad request for invalid JSON
			assert.Equal(t, http.StatusBadRequest, resp.StatusCode)
		})

		t.Run("Unknown Fields in JSON", func(t *testing.T) {
			// Create request with unknown fields
			payload := map[string]interface{}{
				"name":    "Test Client",
				"type":    "Individual",
				"unknown": "field", // This field should be rejected
				"extra":   123,     // This field should be rejected
			}

			payloadBytes, err := json.Marshal(payload)
			require.NoError(t, err)

			req, err := http.NewRequestWithContext(
				ctx,
				http.MethodPost,
				testServerURL+clientHandler.CreateClientEndpoint,
				bytes.NewReader(payloadBytes),
			)
			require.NoError(t, err)

			req.Header.Set("Content-Type", "application/json")

			// Add authentication cookie manually
			cookie := &http.Cookie{
				Name:  cookies.AccessTokenCookieName,
				Value: adminToken,
			}
			req.AddCookie(cookie)

			// Execute request
			resp, err := httpClient.Do(req)
			require.NoError(t, err)
			defer resp.Body.Close()

			// Should return bad request for unknown fields
			assert.Equal(t, http.StatusBadRequest, resp.StatusCode)
		})
	})
}
