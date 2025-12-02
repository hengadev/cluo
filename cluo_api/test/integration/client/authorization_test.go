package client_test

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"testing"
	"time"

	tu "github.com/hengadev/cluo_api/internal/common/testutils"
	ch "github.com/hengadev/cluo_api/test/helpers/client"

	"github.com/google/uuid"
	"github.com/hengadev/encx"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// make test-func TEST_NAME=TestContactAuthorization TEST_PATH=test/integration/client/authorization_test.go

// TestContactAuthorization tests role-based access control for contact endpoints
func TestContactAuthorization(t *testing.T) {
	t.Run("Administrator Full Access", func(t *testing.T) {
		ctx := context.Background()

		// Setup admin authentication
		adminToken := tu.SetupAdminUser(t, ctx, authCtx)
		defer tu.ClearAuthData(t, ctx, authCtx)
		defer ch.ClearContactsTable(t, ctx, testPool)

		// Create test client
		clientID := uuid.New()
		clientIDBytes, err := encx.SerializeValue(clientID)
		require.NoError(t, err)
		clientIDHash := crypto.HashBasic(ctx, clientIDBytes)

		// Create contacts for the client
		contact1, err := ch.CreateTestClientWithContact(t, ctx, testPool, clientID, clientIDHash)
		require.NoError(t, err)
		contact2, err := ch.CreateTestClientWithContact(t, ctx, testPool, clientID, clientIDHash)
		require.NoError(t, err)

		// Test all endpoints as admin
		testCases := []struct {
			name   string
			method string
			path   string
		}{
			{"Create Contact", "POST", fmt.Sprintf("/client/%s/contact", clientID)},
			{"Get Contact", "GET", fmt.Sprintf("/contact/%s", contact1.ID)},
			{"Update Contact", "PATCH", fmt.Sprintf("/contact/%s", contact2.ID)},
			{"Delete Contact", "DELETE", fmt.Sprintf("/contact/%s", contact1.ID)},
			{"Get All Contacts", "GET", fmt.Sprintf("/client/%s/contact", clientID)},
		}

		for _, tc := range testCases {
			t.Run(fmt.Sprintf("Admin %s", tc.name), func(t *testing.T) {
				// Create request
				var req *http.Request
				if tc.method == "POST" || tc.method == "PATCH" || tc.method == "DELETE" {
					payload := map[string]interface{}{
						"clientID":  clientID,
						"lastname":  "Test",
						"firstname": "Admin",
						"email":     "admin@test.com",
						"phone":     "+1234567890",
						"position":  "Admin",
					}
					payloadBytes, err := json.Marshal(payload)
					require.NoError(t, err)
					req, err = http.NewRequest(tc.method, fmt.Sprintf("%s%s", testServerURL, tc.path), bytes.NewBuffer(payloadBytes))
					require.NoError(t, err)
					req.Header.Set("Content-Type", "application/json")
				} else {
					req, err = http.NewRequest(tc.method, fmt.Sprintf("%s%s", testServerURL, tc.path), nil)
					require.NoError(t, err)
				}
				req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", adminToken))

				// Execute request
				httpClient := &http.Client{Timeout: 10 * time.Second}
				resp, err := httpClient.Do(req)
				require.NoError(t, err)
				defer resp.Body.Close()

				// Admin should have access to all operations
				assert.Equal(t, http.StatusOK, resp.StatusCode,
					fmt.Sprintf("Admin should access %s", tc.name))
			})
		}
	})

	t.Run("Client Limited Access", func(t *testing.T) {
		ctx := context.Background()

		// Setup client authentication
		clientToken := tu.SetupClientUser(t, ctx, authCtx)
		defer tu.ClearAuthData(t, ctx, authCtx)
		defer ch.ClearContactsTable(t, ctx, testPool)

		// Create test client
		clientID := uuid.New()
		clientIDBytes, err := encx.SerializeValue(clientID)
		require.NoError(t, err)
		clientIDHash := crypto.HashBasic(ctx, clientIDBytes)

		// Create contacts for the client
		contact1, err := ch.CreateTestClientWithContact(t, ctx, testPool, clientID, clientIDHash)
		require.NoError(t, err)
		contact2, err := ch.CreateTestClientWithContact(t, ctx, testPool, clientID, clientIDHash)
		require.NoError(t, err)

		// Test endpoints client should access
		testCases := []struct {
			name   string
			method string
			path   string
		}{
			{"Create Contact", "POST", fmt.Sprintf("/client/%s/contact", clientID)},
			{"Get Contact", "GET", fmt.Sprintf("/contact/%s", contact1.ID)},
			{"Update Contact", "PATCH", fmt.Sprintf("/contact/%s", contact2.ID)},
			{"Delete Contact", "DELETE", fmt.Sprintf("/contact/%s", contact1.ID)},
			{"Get All Contacts", "GET", fmt.Sprintf("/client/%s/contact", clientID)},
		}

		for _, tc := range testCases {
			t.Run(fmt.Sprintf("Client %s", tc.name), func(t *testing.T) {
				// Create request
				var req *http.Request
				if tc.method == "POST" || tc.method == "PATCH" || tc.method == "DELETE" {
					payload := map[string]interface{}{
						"clientID":  clientID,
						"lastname":  "Test",
						"firstname": "Client",
						"email":     "client@test.com",
						"phone":     "+1234567890",
						"position":  "Client",
					}
					payloadBytes, err := json.Marshal(payload)
					require.NoError(t, err)
					req, err = http.NewRequest(tc.method, fmt.Sprintf("%s%s", testServerURL, tc.path), bytes.NewBuffer(payloadBytes))
					require.NoError(t, err)
					req.Header.Set("Content-Type", "application/json")
				} else {
					req, err = http.NewRequest(tc.method, fmt.Sprintf("%s%s", testServerURL, tc.path), nil)
					require.NoError(t, err)
				}
				req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", clientToken))

				// Execute request
				httpClient := &http.Client{Timeout: 10 * time.Second}
				resp, err := httpClient.Do(req)
				require.NoError(t, err)
				defer resp.Body.Close()

				// Client should access own data
				assert.Equal(t, http.StatusOK, resp.StatusCode,
					fmt.Sprintf("Client should access %s", tc.name))
			})
		}

		// Test client cannot access other client's data
		t.Run("Client Cross-Client Access", func(t *testing.T) {
			// Create another client
			otherClientID := uuid.New()
			otherClientIDBytes, err := encx.SerializeValue(otherClientID)
			require.NoError(t, err)
			otherClientIDHash := crypto.HashBasic(ctx, otherClientIDBytes)

			// Create contact for other client
			otherContact, err := ch.CreateTestClientWithContact(t, ctx, testPool, otherClientID, otherClientIDHash)
			require.NoError(t, err)

			// Test client accessing other client's contact
			req, err := http.NewRequest("GET", fmt.Sprintf("%s/contact/%s", testServerURL, otherContact.ID), nil)
			require.NoError(t, err)
			req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", clientToken))

			// Execute request
			httpClient := &http.Client{Timeout: 10 * time.Second}
			resp, err := httpClient.Do(req)
			require.NoError(t, err)
			defer resp.Body.Close()

			// Should be denied
			assert.True(t, resp.StatusCode == http.StatusUnauthorized || resp.StatusCode == http.StatusForbidden,
				fmt.Sprintf("Client should not access other client's contact, status code %d", resp.StatusCode))
		})
	})

	t.Run("Guest No Access", func(t *testing.T) {
		ctx := context.Background()

		// Setup guest authentication
		guestToken := tu.SetupGuestUser(t, ctx, authCtx)
		defer tu.ClearAuthData(t, ctx, authCtx)
		defer ch.ClearContactsTable(t, ctx, testPool)

		// Create test client
		clientID := uuid.New()
		clientIDBytes, err := encx.SerializeValue(clientID)
		require.NoError(t, err)
		clientIDHash := crypto.HashBasic(ctx, clientIDBytes)

		// Create contact for the client
		contact, err := ch.CreateTestClientWithContact(t, ctx, testPool, clientID, clientIDHash)
		require.NoError(t, err)

		// Test all endpoints guests should not access
		testCases := []struct {
			name   string
			method string
			path   string
		}{
			{"Create Contact", "POST", fmt.Sprintf("/client/%s/contact", clientID)},
			{"Get Contact", "GET", fmt.Sprintf("/contact/%s", contact.ID)},
			{"Update Contact", "PATCH", fmt.Sprintf("/contact/%s", contact.ID)},
			{"Delete Contact", "DELETE", fmt.Sprintf("/contact/%s", contact.ID)},
			{"Get All Contacts", "GET", fmt.Sprintf("/client/%s/contact", clientID)},
		}

		for _, tc := range testCases {
			t.Run(fmt.Sprintf("Guest %s", tc.name), func(t *testing.T) {
				// Create request
				var req *http.Request
				if tc.method == "POST" || tc.method == "PATCH" || tc.method == "DELETE" {
					payload := map[string]interface{}{
						"clientID":  clientID,
						"lastname":  "Test",
						"firstname": "Guest",
						"email":     "guest@test.com",
						"phone":     "+1234567890",
						"position":  "Guest",
					}
					payloadBytes, err := json.Marshal(payload)
					require.NoError(t, err)
					req, err = http.NewRequest(tc.method, fmt.Sprintf("%s%s", testServerURL, tc.path), bytes.NewBuffer(payloadBytes))
					require.NoError(t, err)
					req.Header.Set("Content-Type", "application/json")
				} else {
					req, err = http.NewRequest(tc.method, fmt.Sprintf("%s%s", testServerURL, tc.path), nil)
					require.NoError(t, err)
				}
				req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", guestToken))

				// Execute request
				httpClient := &http.Client{Timeout: 10 * time.Second}
				resp, err := httpClient.Do(req)
				require.NoError(t, err)
				defer resp.Body.Close()

				// Guest should be denied access to all operations
				assert.Equal(t, http.StatusForbidden, resp.StatusCode,
					fmt.Sprintf("Guest should not access %s", tc.name))
			})
		}
	})
}
