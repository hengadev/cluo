package client_test

import (
	"context"
	"encoding/json"
	"net/http"
	"testing"
	"time"

	"github.com/google/uuid"
	tu "github.com/hengadev/cluo_api/internal/common/testutils"
	"github.com/hengadev/cluo_api/internal/domain/client"
	ch "github.com/hengadev/cluo_api/test/helpers/client"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// make test-func TEST_NAME=TestGetAllContactsByClientID TEST_PATH=test/integration/client/get_all_contacts_by_client_id_test.go

// TestGetAllContactsByClientID tests all scenarios for the GetAllContactsByClientID endpoint
func TestGetAllContactsByClientID(t *testing.T) {
	httpClient := &http.Client{Timeout: 10 * time.Second}

	setupClientWithContacts := func(t *testing.T, ctx context.Context, numContacts int) (*client.ClientEncx, []uuid.UUID) {
		// Create client
		c := ch.NewTestClient(t)
		clientEncx, err := client.ProcessClientEncx(ctx, crypto, c)
		require.NoError(t, err)
		err = ch.InsertClientEncx(t, ctx, testPool, clientEncx)
		require.NoError(t, err)

		// Create contacts
		contactIDs := make([]uuid.UUID, numContacts)
		for i := 0; i < numContacts; i++ {
			contact := ch.NewTestContact(t)
			contact.ClientID = clientEncx.ID
			contactEncx, err := client.ProcessContactEncx(ctx, crypto, contact)
			require.NoError(t, err)
			err = ch.InsertContactEncx(t, ctx, testPool, contactEncx)
			require.NoError(t, err)
			contactIDs[i] = contactEncx.ID
		}

		return clientEncx, contactIDs
	}

	t.Run("Success - Client with Multiple Contacts", func(t *testing.T) {
		ctx := context.Background()

		// Setup administrator authentication
		accessToken := tu.SetupAdminUser(t, ctx, authCtx)
		defer tu.ClearAuthData(t, ctx, authCtx)
		defer ch.ClearContactsTable(t, ctx, testPool)
		defer ch.ClearClientsTable(t, ctx, testPool)

		// Setup client with multiple contacts
		clientEncx, contactIDs := setupClientWithContacts(t, ctx, 5)

		// Create HTTP request
		req := ch.NewGetAllContactsByClientIDRequest(t, ctx, testServerURL, clientEncx.ID.String(), accessToken)

		// Execute request
		resp, err := httpClient.Do(req)
		require.NoError(t, err)
		defer resp.Body.Close()

		// Assert response status
		assert.Equal(t, http.StatusOK, resp.StatusCode)

		// Decode response
		var response []client.ContactResponse
		err = json.NewDecoder(resp.Body).Decode(&response)
		require.NoError(t, err)

		// Verify response data
		assert.Len(t, response, 5, "Expected 5 contacts")

		// Verify all contact IDs are present
		responseContactIDs := make(map[string]bool)
		for _, contact := range response {
			responseContactIDs[contact.ID] = true
			assert.Equal(t, clientEncx.ID.String(), contact.ClientID, "Contact should belong to the correct client")
			assert.NotEmpty(t, contact.Lastname, "Contact should have a lastname")
			assert.NotEmpty(t, contact.Firstname, "Contact should have a firstname")
			assert.NotEmpty(t, contact.Email, "Contact should have an email")
		}

		// Verify all expected contact IDs are in the response
		for _, expectedID := range contactIDs {
			assert.True(t, responseContactIDs[expectedID.String()], "Expected contact ID %s should be in response", expectedID)
		}

		t.Log("✓ Retrieved all contacts for client successfully")
	})

	t.Run("Success - Client with No Contacts", func(t *testing.T) {
		ctx := context.Background()

		// Setup administrator authentication
		accessToken := tu.SetupAdminUser(t, ctx, authCtx)
		defer tu.ClearAuthData(t, ctx, authCtx)
		defer ch.ClearContactsTable(t, ctx, testPool)
		defer ch.ClearClientsTable(t, ctx, testPool)

		// Create client without contacts
		c := ch.NewTestClient(t)
		clientEncx, err := client.ProcessClientEncx(ctx, crypto, c)
		require.NoError(t, err)
		err = ch.InsertClientEncx(t, ctx, testPool, clientEncx)
		require.NoError(t, err)

		// Create HTTP request
		req := ch.NewGetAllContactsByClientIDRequest(t, ctx, testServerURL, clientEncx.ID.String(), accessToken)

		// Execute request
		resp, err := httpClient.Do(req)
		require.NoError(t, err)
		defer resp.Body.Close()

		// Assert response status
		assert.Equal(t, http.StatusOK, resp.StatusCode)

		// Decode response
		var response []client.ContactResponse
		err = json.NewDecoder(resp.Body).Decode(&response)
		require.NoError(t, err)

		// Verify response data
		assert.Empty(t, response, "Expected no contacts for client")

		t.Log("✓ Retrieved empty contact list for client successfully")
	})

	t.Run("Success - Single Contact", func(t *testing.T) {
		ctx := context.Background()

		// Setup administrator authentication
		accessToken := tu.SetupAdminUser(t, ctx, authCtx)
		defer tu.ClearAuthData(t, ctx, authCtx)
		defer ch.ClearContactsTable(t, ctx, testPool)
		defer ch.ClearClientsTable(t, ctx, testPool)

		// Setup client with single contact
		clientEncx, contactIDs := setupClientWithContacts(t, ctx, 1)

		// Create HTTP request
		req := ch.NewGetAllContactsByClientIDRequest(t, ctx, testServerURL, clientEncx.ID.String(), accessToken)

		// Execute request
		resp, err := httpClient.Do(req)
		require.NoError(t, err)
		defer resp.Body.Close()

		// Assert response status
		assert.Equal(t, http.StatusOK, resp.StatusCode)

		// Decode response
		var response []client.ContactResponse
		err = json.NewDecoder(resp.Body).Decode(&response)
		require.NoError(t, err)

		// Verify response data
		assert.Len(t, response, 1, "Expected exactly 1 contact")
		assert.Equal(t, contactIDs[0].String(), response[0].ID, "Contact ID should match")

		t.Log("✓ Retrieved single contact for client successfully")
	})

	t.Run("Error - Client Not Found", func(t *testing.T) {
		ctx := context.Background()

		// Setup administrator authentication
		accessToken := tu.SetupAdminUser(t, ctx, authCtx)
		defer tu.ClearAuthData(t, ctx, authCtx)
		defer ch.ClearContactsTable(t, ctx, testPool)
		defer ch.ClearClientsTable(t, ctx, testPool)

		// Use a non-existent client ID
		nonExistentClientID := uuid.New()

		// Create HTTP request
		req := ch.NewGetAllContactsByClientIDRequest(t, ctx, testServerURL, nonExistentClientID.String(), accessToken)

		// Execute request
		resp, err := httpClient.Do(req)
		require.NoError(t, err)
		defer resp.Body.Close()

		// Assert response status - should return empty contacts list for non-existent client
		assert.Equal(t, http.StatusOK, resp.StatusCode)

		// Decode response
		var response []client.ContactResponse
		err = json.NewDecoder(resp.Body).Decode(&response)
		require.NoError(t, err)

		// Should return empty contacts list
		assert.Empty(t, response, "Expected no contacts for non-existent client")

		t.Log("✓ Non-existent client handled correctly (returns empty list)")
	})

	t.Run("Error - Invalid Client ID", func(t *testing.T) {
		ctx := context.Background()

		// Setup administrator authentication
		accessToken := tu.SetupAdminUser(t, ctx, authCtx)
		defer tu.ClearAuthData(t, ctx, authCtx)

		invalidClientID := "invalid-uuid-format"

		// Create HTTP request with invalid client ID
		req := ch.NewGetAllContactsByClientIDRequest(t, ctx, testServerURL, invalidClientID, accessToken)

		// Execute request
		resp, err := httpClient.Do(req)
		require.NoError(t, err)
		defer resp.Body.Close()

		// Assert response status - should be 400 Bad Request (due to invalid UUID format)
		assert.Equal(t, http.StatusBadRequest, resp.StatusCode)

		t.Log("✓ Invalid client ID error handled correctly")
	})

	t.Run("Success - Contacts for Different Clients Are Separated", func(t *testing.T) {
		ctx := context.Background()

		// Setup administrator authentication
		accessToken := tu.SetupAdminUser(t, ctx, authCtx)
		defer tu.ClearAuthData(t, ctx, authCtx)
		defer ch.ClearContactsTable(t, ctx, testPool)
		defer ch.ClearClientsTable(t, ctx, testPool)

		// Create two clients with contacts
		client1Encx, _ := setupClientWithContacts(t, ctx, 2)
		client2Encx, _ := setupClientWithContacts(t, ctx, 3)

		// Test first client's contacts
		req1 := ch.NewGetAllContactsByClientIDRequest(t, ctx, testServerURL, client1Encx.ID.String(), accessToken)
		resp1, err := httpClient.Do(req1)
		require.NoError(t, err)
		defer resp1.Body.Close()

		assert.Equal(t, http.StatusOK, resp1.StatusCode)

		var response1 []client.ContactResponse
		err = json.NewDecoder(resp1.Body).Decode(&response1)
		require.NoError(t, err)

		assert.Len(t, response1, 2, "Expected 2 contacts for first client")

		// Verify all contacts belong to first client
		for _, contact := range response1 {
			assert.Equal(t, client1Encx.ID.String(), contact.ClientID, "Contact should belong to first client")
		}

		// Test second client's contacts
		req2 := ch.NewGetAllContactsByClientIDRequest(t, ctx, testServerURL, client2Encx.ID.String(), accessToken)
		resp2, err := httpClient.Do(req2)
		require.NoError(t, err)
		defer resp2.Body.Close()

		assert.Equal(t, http.StatusOK, resp2.StatusCode)

		var response2 []client.ContactResponse
		err = json.NewDecoder(resp2.Body).Decode(&response2)
		require.NoError(t, err)

		assert.Len(t, response2, 3, "Expected 3 contacts for second client")

		// Verify all contacts belong to second client
		for _, contact := range response2 {
			assert.Equal(t, client2Encx.ID.String(), contact.ClientID, "Contact should belong to second client")
		}

		t.Log("✓ Contacts for different clients are properly separated")
	})

	t.Run("Authentication/Authorization", func(t *testing.T) {
		ctx := context.Background()

		// Setup client with contacts
		clientEncx, _ := setupClientWithContacts(t, ctx, 2)

		t.Run("No Authentication", func(t *testing.T) {
			// Create HTTP request without authentication
			req := ch.NewGetAllContactsByClientIDRequest(t, ctx, testServerURL, clientEncx.ID.String(), "")

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
			req := ch.NewGetAllContactsByClientIDRequest(t, ctx, testServerURL, clientEncx.ID.String(), clientToken)

			// Execute request
			resp, err := httpClient.Do(req)
			require.NoError(t, err)
			defer resp.Body.Close()

			// Client should be forbidden from getting contacts
			assert.Equal(t, http.StatusForbidden, resp.StatusCode)
			t.Log("✓ Client authentication blocked correctly")
		})

		t.Run("Guest Authentication", func(t *testing.T) {
			// Setup guest authentication
			guestToken := tu.SetupGuestUser(t, ctx, authCtx)
			defer tu.ClearAuthData(t, ctx, authCtx)

			// Create HTTP request with guest authentication
			req := ch.NewGetAllContactsByClientIDRequest(t, ctx, testServerURL, clientEncx.ID.String(), guestToken)

			// Execute request
			resp, err := httpClient.Do(req)
			require.NoError(t, err)
			defer resp.Body.Close()

			// Guest should be forbidden from getting contacts
			assert.Equal(t, http.StatusForbidden, resp.StatusCode)
			t.Log("✓ Guest authentication blocked correctly")
		})

		// Clean up the test data
		ch.ClearContactsTable(t, ctx, testPool)
		ch.ClearClientsTable(t, ctx, testPool)
	})

	t.Run("Context Cancellation", func(t *testing.T) {
		ctx := context.Background()

		// Setup administrator authentication
		accessToken := tu.SetupAdminUser(t, ctx, authCtx)
		defer tu.ClearAuthData(t, ctx, authCtx)

		// Create a context that will be cancelled
		cancelledCtx, cancel := context.WithCancel(context.Background())
		cancel() // Cancel immediately

		// Use a non-existent client ID to avoid modifying actual data
		clientID := uuid.New()

		// Create HTTP request with cancelled context
		req := ch.NewGetAllContactsByClientIDRequest(t, cancelledCtx, testServerURL, clientID.String(), accessToken)

		// Execute request
		_, err := httpClient.Do(req)
		assert.Error(t, err)

		t.Log("✓ Context cancellation handled correctly")
	})

	t.Run("Performance with Large Dataset", func(t *testing.T) {
		ctx := context.Background()

		// Setup administrator authentication
		accessToken := tu.SetupAdminUser(t, ctx, authCtx)
		defer tu.ClearAuthData(t, ctx, authCtx)
		defer ch.ClearContactsTable(t, ctx, testPool)
		defer ch.ClearClientsTable(t, ctx, testPool)

		// Setup client with many contacts
		clientEncx, _ := setupClientWithContacts(t, ctx, 100)

		// Execute request and measure time
		start := time.Now()
		req := ch.NewGetAllContactsByClientIDRequest(t, ctx, testServerURL, clientEncx.ID.String(), accessToken)

		resp, err := httpClient.Do(req)
		duration := time.Since(start)

		require.NoError(t, err)
		defer resp.Body.Close()

		// Assert response status
		assert.Equal(t, http.StatusOK, resp.StatusCode)

		// Decode response
		var response []client.ContactResponse
		err = json.NewDecoder(resp.Body).Decode(&response)
		require.NoError(t, err)

		// Verify response contains all contacts
		assert.Len(t, response, 100, "Expected 100 contacts")

		// Log performance
		t.Logf("✓ Retrieved 100 contacts in %v", duration)

		// Performance assertion - should complete within reasonable time
		assert.Less(t, duration, 2*time.Second, "Request should complete within 2 seconds")
	})
}
