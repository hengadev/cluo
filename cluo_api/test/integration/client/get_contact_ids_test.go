package client_test

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"testing"
	"time"

	"github.com/google/uuid"
	tu "github.com/hengadev/cluo_api/internal/common/testutils"
	"github.com/hengadev/cluo_api/internal/domain/client"
	clientHandler "github.com/hengadev/cluo_api/internal/interface/client"
	ch "github.com/hengadev/cluo_api/test/helpers/client"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// make test-func TEST_NAME=TestGetContactIDsForClient TEST_PATH=test/integration/client/get_contact_ids_test.go

// TestGetContactIDsForClient tests all scenarios for the GetContactIDsForClient endpoint
func TestGetContactIDsForClient(t *testing.T) {
	httpClient := &http.Client{Timeout: 10 * time.Second}

	setupClient := func(t *testing.T, ctx context.Context) uuid.UUID {
		c := ch.NewTestClient(t)
		clientEncx, err := client.ProcessClientEncx(ctx, crypto, c)
		require.NoError(t, err)
		err = ch.InsertClientEncx(t, ctx, testPool, clientEncx)
		require.NoError(t, err)
		return c.ID
	}

	setupContacts := func(t *testing.T, ctx context.Context, clientID uuid.UUID) []uuid.UUID {
		contactIDs := make([]uuid.UUID, 3)

		// Create test contacts
		contacts := []*client.Contact{
			client.NewContact(&client.CreateContactRequest{
				ClientID:  clientID,
				Lastname:  "Smith",
				Firstname: "John",
				Email:     "john.smith@example.com",
				Phone:     "0612345678",
				Position:  "Manager",
			}),
			client.NewContact(&client.CreateContactRequest{
				ClientID:  clientID,
				Lastname:  "Doe",
				Firstname: "Jane",
				Email:     "jane.doe@example.com",
				Phone:     "0623456789",
				Position:  "Developer",
			}),
			client.NewContact(&client.CreateContactRequest{
				ClientID:  clientID,
				Lastname:  "Wilson",
				Firstname: "Bob",
				Email:     "bob.wilson@example.com",
				Phone:     "0634567890",
				Position:  "Analyst",
			}),
		}

		for i, contact := range contacts {
			contactEncx, err := client.ProcessContactEncx(ctx, crypto, contact)
			require.NoError(t, err)
			err = ch.InsertContactEncx(t, ctx, testPool, contactEncx)
			require.NoError(t, err)
			contactIDs[i] = contact.ID
		}

		return contactIDs
	}

	t.Run("Success", func(t *testing.T) {
		ctx := context.Background()

		// Setup administrator authentication
		accessToken := tu.SetupAdminUser(t, ctx, authCtx)
		defer tu.ClearAuthData(t, ctx, authCtx)
		defer ch.ClearContactsTable(t, ctx, testPool)
		defer ch.ClearClientsTable(t, ctx, testPool)

		// Setup client and contacts
		clientID := setupClient(t, ctx)
		expectedContactIDs := setupContacts(t, ctx, clientID)

		// Create HTTP request using the test helper
		req := ch.NewGetContactIDsForClientRequest(t, ctx, testServerURL, clientID.String(), accessToken)

		// Execute request
		resp, err := httpClient.Do(req)
		require.NoError(t, err)
		defer resp.Body.Close()

		// Assert response status
		assert.Equal(t, http.StatusOK, resp.StatusCode)

		// Decode response
		var response *clientHandler.GetContactIDsResponse
		err = json.NewDecoder(resp.Body).Decode(&response)
		require.NoError(t, err)

		// Verify response data
		assert.Equal(t, clientID.String(), response.ClientID)
		assert.Equal(t, 3, response.Count)
		assert.Equal(t, 3, len(response.ContactIDs))

		// Verify all expected contact IDs are present
		for _, expectedID := range expectedContactIDs {
			assert.Contains(t, response.ContactIDs, expectedID.String())
		}

		t.Log("✓ Contact IDs retrieved successfully")
	})

	t.Run("Client with No Contacts", func(t *testing.T) {
		ctx := context.Background()

		// Setup administrator authentication
		accessToken := tu.SetupAdminUser(t, ctx, authCtx)
		defer tu.ClearAuthData(t, ctx, authCtx)
		defer ch.ClearContactsTable(t, ctx, testPool)
		defer ch.ClearClientsTable(t, ctx, testPool)

		// Setup client without contacts
		clientID := setupClient(t, ctx)

		// Create HTTP request
		req := ch.NewGetContactIDsForClientRequest(t, ctx, testServerURL, clientID.String(), accessToken)

		// Execute request
		resp, err := httpClient.Do(req)
		require.NoError(t, err)
		defer resp.Body.Close()

		// Assert response status
		assert.Equal(t, http.StatusOK, resp.StatusCode)

		// Decode response
		var response *clientHandler.GetContactIDsResponse
		err = json.NewDecoder(resp.Body).Decode(&response)
		require.NoError(t, err)

		// Verify response data
		assert.Equal(t, clientID.String(), response.ClientID)
		assert.Equal(t, 0, response.Count)
		assert.Empty(t, response.ContactIDs)

		t.Log("✓ Empty contact IDs list returned correctly")
	})

	t.Run("Client Not Found", func(t *testing.T) {
		ctx := context.Background()

		// Setup administrator authentication
		accessToken := tu.SetupAdminUser(t, ctx, authCtx)
		defer tu.ClearAuthData(t, ctx, authCtx)
		defer ch.ClearContactsTable(t, ctx, testPool)
		defer ch.ClearClientsTable(t, ctx, testPool)

		// Use a non-existent client ID
		nonExistentClientID := uuid.New()

		// Create HTTP request
		req := ch.NewGetContactIDsForClientRequest(t, ctx, testServerURL, nonExistentClientID.String(), accessToken)

		// Execute request
		resp, err := httpClient.Do(req)
		require.NoError(t, err)
		defer resp.Body.Close()

		// Assert response status - should be 404 Not Found
		assert.Equal(t, http.StatusNotFound, resp.StatusCode)

		t.Log("✓ Client not found error handled correctly")
	})

	t.Run("Invalid Client ID", func(t *testing.T) {
		ctx := context.Background()

		// Setup administrator authentication
		accessToken := tu.SetupAdminUser(t, ctx, authCtx)
		defer tu.ClearAuthData(t, ctx, authCtx)

		invalidClientID := "invalid-uuid-format"

		// Create HTTP request
		req := ch.NewGetContactIDsForClientRequest(t, ctx, testServerURL, invalidClientID, accessToken)

		// Execute request
		resp, err := httpClient.Do(req)
		require.NoError(t, err)
		defer resp.Body.Close()

		// Assert response status - should be 400 Not Found (due to invalid UUID format)
		assert.Equal(t, http.StatusBadRequest, resp.StatusCode)

		t.Log("✓ Invalid client ID error handled correctly")
	})

	t.Run("Authentication/Authorization", func(t *testing.T) {
		ctx := context.Background()

		// Setup client and contacts
		clientID := setupClient(t, ctx)
		setupContacts(t, ctx, clientID)

		t.Run("No Authentication", func(t *testing.T) {
			// Create HTTP request without authentication
			req := ch.NewGetContactIDsForClientRequest(t, ctx, testServerURL, clientID.String(), "")

			// Execute request
			resp, err := httpClient.Do(req)
			require.NoError(t, err)
			defer resp.Body.Close()

			// Should be unauthorized
			assert.Equal(t, http.StatusUnauthorized, resp.StatusCode)
			t.Log("✓ No authentication blocked correctly")
		})

		t.Run("Client Authentication", func(t *testing.T) {
			// Setup client authentication
			clientToken := tu.SetupClientUser(t, ctx, authCtx)
			defer tu.ClearAuthData(t, ctx, authCtx)

			// Create HTTP request with client authentication
			req := ch.NewGetContactIDsForClientRequest(t, ctx, testServerURL, clientID.String(), clientToken)

			// Execute request
			resp, err := httpClient.Do(req)
			require.NoError(t, err)
			defer resp.Body.Close()

			// Client should be forbidden from accessing contact IDs
			assert.Equal(t, http.StatusForbidden, resp.StatusCode)
			t.Log("✓ Client authentication blocked correctly")
		})

		t.Run("Guest Authentication", func(t *testing.T) {
			// Setup guest authentication
			guestToken := tu.SetupGuestUser(t, ctx, authCtx)
			defer tu.ClearAuthData(t, ctx, authCtx)

			// Create HTTP request with guest authentication
			req := ch.NewGetContactIDsForClientRequest(t, ctx, testServerURL, clientID.String(), guestToken)

			// Execute request
			resp, err := httpClient.Do(req)
			require.NoError(t, err)
			defer resp.Body.Close()

			// Guest should be forbidden from accessing contact IDs
			assert.Equal(t, http.StatusForbidden, resp.StatusCode)
			t.Log("✓ Guest authentication blocked correctly")
		})
	})

	t.Run("Performance with Many Contacts", func(t *testing.T) {
		ctx := context.Background()

		// Setup administrator authentication
		accessToken := tu.SetupAdminUser(t, ctx, authCtx)
		defer tu.ClearAuthData(t, ctx, authCtx)
		defer ch.ClearContactsTable(t, ctx, testPool)
		defer ch.ClearClientsTable(t, ctx, testPool)

		// Setup client with many contacts
		clientID := setupClient(t, ctx)

		// Create 50 contacts for performance testing
		for i := 0; i < 50; i++ {
			contact := client.NewContact(&client.CreateContactRequest{
				ClientID:  clientID,
				Lastname:  fmt.Sprintf("TestContact_%d", i),
				Firstname: fmt.Sprintf("First_%d", i),
				Email:     fmt.Sprintf("contact%d@test.com", i),
				Phone:     fmt.Sprintf("06%08d", i),
				Position:  fmt.Sprintf("Position_%d", i),
			})

			contactEncx, err := client.ProcessContactEncx(ctx, crypto, contact)
			require.NoError(t, err)
			err = ch.InsertContactEncx(t, ctx, testPool, contactEncx)
			require.NoError(t, err)
		}

		// Create HTTP request
		req := ch.NewGetContactIDsForClientRequest(t, ctx, testServerURL, clientID.String(), accessToken)

		// Execute request and measure time
		start := time.Now()
		resp, err := httpClient.Do(req)
		duration := time.Since(start)

		require.NoError(t, err)
		defer resp.Body.Close()

		// Assert response status
		assert.Equal(t, http.StatusOK, resp.StatusCode)

		// Decode response
		var response *clientHandler.GetContactIDsResponse
		err = json.NewDecoder(resp.Body).Decode(&response)
		require.NoError(t, err)

		// Verify response data
		assert.Equal(t, clientID.String(), response.ClientID)
		assert.Equal(t, 50, response.Count)
		assert.Equal(t, 50, len(response.ContactIDs))

		// Log performance
		t.Logf("✓ Retrieved %d contact IDs in %v", response.Count, duration)

		// Performance assertion - should complete within reasonable time
		assert.Less(t, duration, 5*time.Second, "Request should complete within 5 seconds")
	})
}
