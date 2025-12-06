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

// make test-func TEST_NAME=TestUpdateContact TEST_PATH=test/integration/client/update_contact_test.go

// TestUpdateContact tests all scenarios for the UpdateContact endpoint
func TestUpdateContact(t *testing.T) {
	httpClient := &http.Client{Timeout: 10 * time.Second}

	setupClientWithContact := func(t *testing.T, ctx context.Context) (*client.ClientEncx, *client.ContactEncx) {
		c := ch.NewTestClient(t)
		clientEncx, err := client.ProcessClientEncx(ctx, crypto, c)
		require.NoError(t, err)
		err = ch.InsertClientEncx(t, ctx, testPool, clientEncx)
		require.NoError(t, err)

		contact := ch.NewTestContact(t)
		contact.ClientID = clientEncx.ID
		contactEncx, err := client.ProcessContactEncx(ctx, crypto, contact)
		require.NoError(t, err)
		err = ch.InsertContactEncx(t, ctx, testPool, contactEncx)
		require.NoError(t, err)

		return clientEncx, contactEncx
	}

	t.Run("Success - Partial Update", func(t *testing.T) {
		ctx := context.Background()

		// Setup administrator authentication
		accessToken := tu.SetupAdminUser(t, ctx, authCtx)
		defer tu.ClearAuthData(t, ctx, authCtx)
		defer ch.ClearContactsTable(t, ctx, testPool)
		defer ch.ClearClientsTable(t, ctx, testPool)

		// Setup client and contact
		clientEncx, contactEncx := setupClientWithContact(t, ctx)

		// Create update request with only email and phone
		updateReq := client.UpdateContactRequest{
			ID:       contactEncx.ID,
			ClientID: clientEncx.ID.String(),
			Email:    func() *string { s := "updated.email@example.com"; return &s }(),
			Phone:    func() *string { s := "0612345679"; return &s }(),
		}

		// Create HTTP request
		req := ch.NewUpdateContactRequest(t, ctx, testServerURL, contactEncx.ID.String(), updateReq, accessToken)

		// Execute request
		resp, err := httpClient.Do(req)
		require.NoError(t, err)
		defer resp.Body.Close()

		// Assert response status
		assert.Equal(t, http.StatusOK, resp.StatusCode)

		// Decode response
		var response client.ContactResponse
		err = json.NewDecoder(resp.Body).Decode(&response)
		require.NoError(t, err)

		// Verify response data
		assert.Equal(t, contactEncx.ID.String(), response.ID)
		assert.Equal(t, "updated.email@example.com", response.Email)
		assert.Equal(t, "0612345679", response.Phone)
		// Verify non-updated fields remain the same
		assert.NotEmpty(t, response.Lastname)  // Should have original value
		assert.NotEmpty(t, response.Firstname) // Should have original value

		t.Log("✓ Contact partially updated successfully")
	})

	t.Run("Success - Full Update", func(t *testing.T) {
		ctx := context.Background()

		// Setup administrator authentication
		accessToken := tu.SetupAdminUser(t, ctx, authCtx)
		defer tu.ClearAuthData(t, ctx, authCtx)
		defer ch.ClearContactsTable(t, ctx, testPool)
		defer ch.ClearClientsTable(t, ctx, testPool)

		// Setup client and contact
		clientEncx, contactEncx := setupClientWithContact(t, ctx)

		// Create update request with all fields
		updateReq := client.UpdateContactRequest{
			ID:        contactEncx.ID,
			ClientID:  clientEncx.ID.String(),
			Lastname:  func() *string { s := "UpdatedLastname"; return &s }(),
			Firstname: func() *string { s := "UpdatedFirstname"; return &s }(),
			Email:     func() *string { s := "fully.updated@example.com"; return &s }(),
			Phone:     func() *string { s := "0612345679"; return &s }(),
			Position:  func() *string { s := "Updated Position"; return &s }(),
		}

		// Create HTTP request
		req := ch.NewUpdateContactRequest(t, ctx, testServerURL, contactEncx.ID.String(), updateReq, accessToken)

		// Execute request
		resp, err := httpClient.Do(req)
		require.NoError(t, err)
		defer resp.Body.Close()

		// Assert response status
		assert.Equal(t, http.StatusOK, resp.StatusCode)

		// Decode response
		var response client.ContactResponse
		err = json.NewDecoder(resp.Body).Decode(&response)
		require.NoError(t, err)

		// Verify all fields are updated
		assert.Equal(t, contactEncx.ID.String(), response.ID)
		assert.Equal(t, "UpdatedLastname", response.Lastname)
		assert.Equal(t, "UpdatedFirstname", response.Firstname)
		assert.Equal(t, "fully.updated@example.com", response.Email)
		assert.Equal(t, "0612345679", response.Phone)
		assert.Equal(t, "Updated Position", response.Position)

		t.Log("✓ Contact fully updated successfully")
	})

	t.Run("Success - Nil Pointers Don't Modify Fields", func(t *testing.T) {
		ctx := context.Background()

		// Setup administrator authentication
		accessToken := tu.SetupAdminUser(t, ctx, authCtx)
		defer tu.ClearAuthData(t, ctx, authCtx)
		defer ch.ClearContactsTable(t, ctx, testPool)
		defer ch.ClearClientsTable(t, ctx, testPool)

		// Setup client and contact
		clientEncx, contactEncx := setupClientWithContact(t, ctx)

		// Get original contact data for comparison
		originalContactEncx, err := ch.GetContactEncxByID(t, ctx, testPool, contactEncx.ID)
		require.NoError(t, err)
		originalResponse, err := client.DecryptContactEncx(ctx, crypto, originalContactEncx)
		require.NoError(t, err)

		// Create update request with only some fields (others nil)
		updateReq := client.UpdateContactRequest{
			ID:       contactEncx.ID,
			ClientID: clientEncx.ID.String(),
			Email:    func() *string { s := "partial.update@example.com"; return &s }(),
			// Lastname, Firstname, Phone, Position are nil
		}

		// Create HTTP request
		req := ch.NewUpdateContactRequest(t, ctx, testServerURL, contactEncx.ID.String(), updateReq, accessToken)

		// Execute request
		resp, err := httpClient.Do(req)
		require.NoError(t, err)
		defer resp.Body.Close()

		// Assert response status
		assert.Equal(t, http.StatusOK, resp.StatusCode)

		// Decode response
		var response client.ContactResponse
		err = json.NewDecoder(resp.Body).Decode(&response)
		require.NoError(t, err)

		// Verify only specified field is updated, others remain the same
		assert.Equal(t, "partial.update@example.com", response.Email)
		assert.Equal(t, originalResponse.Lastname, response.Lastname)
		assert.Equal(t, originalResponse.Firstname, response.Firstname)
		assert.Equal(t, originalResponse.Phone, response.Phone)
		assert.Equal(t, originalResponse.Position, response.Position)

		t.Log("✓ Nil pointers correctly preserved existing field values")
	})

	t.Run("Error - Contact Not Found", func(t *testing.T) {
		ctx := context.Background()

		// Setup administrator authentication
		accessToken := tu.SetupAdminUser(t, ctx, authCtx)
		defer tu.ClearAuthData(t, ctx, authCtx)
		defer ch.ClearContactsTable(t, ctx, testPool)
		defer ch.ClearClientsTable(t, ctx, testPool)

		// Use a non-existent contact ID
		nonExistentContactID := uuid.New()
		clientID := uuid.New()

		// Create update request
		updateReq := client.UpdateContactRequest{
			ID:       nonExistentContactID,
			ClientID: clientID.String(),
			Email:    func() *string { s := "test@example.com"; return &s }(),
		}

		// Create HTTP request
		req := ch.NewUpdateContactRequest(t, ctx, testServerURL, nonExistentContactID.String(), updateReq, accessToken)

		// Execute request
		resp, err := httpClient.Do(req)
		require.NoError(t, err)
		defer resp.Body.Close()

		// Assert response status - should be 404 Not Found
		assert.Equal(t, http.StatusNotFound, resp.StatusCode)

		t.Log("✓ Contact not found error handled correctly")
	})

	t.Run("Error - Invalid Contact ID", func(t *testing.T) {
		ctx := context.Background()

		// Setup administrator authentication
		accessToken := tu.SetupAdminUser(t, ctx, authCtx)
		defer tu.ClearAuthData(t, ctx, authCtx)

		invalidContactID := "invalid-uuid-format"

		// Create update request
		updateReq := client.UpdateContactRequest{
			ID:       uuid.New(), // Valid UUID in request
			ClientID: uuid.New().String(),
			Email:    func() *string { s := "test@example.com"; return &s }(),
		}

		// Create HTTP request with invalid contact ID in URL
		req := ch.NewUpdateContactRequest(t, ctx, testServerURL, invalidContactID, updateReq, accessToken)

		// Execute request
		resp, err := httpClient.Do(req)
		require.NoError(t, err)
		defer resp.Body.Close()

		// Assert response status - should be 400 Bad Request or 404 Not Found
		assert.True(t, resp.StatusCode == http.StatusBadRequest || resp.StatusCode == http.StatusNotFound)

		t.Log("✓ Invalid contact ID error handled correctly")
	})

	t.Run("Error - Invalid Email Format", func(t *testing.T) {
		ctx := context.Background()

		// Setup administrator authentication
		accessToken := tu.SetupAdminUser(t, ctx, authCtx)
		defer tu.ClearAuthData(t, ctx, authCtx)
		defer ch.ClearContactsTable(t, ctx, testPool)
		defer ch.ClearClientsTable(t, ctx, testPool)

		// Setup client and contact
		clientEncx, contactEncx := setupClientWithContact(t, ctx)

		// Create update request with invalid email
		updateReq := client.UpdateContactRequest{
			ID:       contactEncx.ID,
			ClientID: clientEncx.ID.String(),
			Email:    func() *string { s := "invalid-email-format"; return &s }(),
		}

		// Create HTTP request
		req := ch.NewUpdateContactRequest(t, ctx, testServerURL, contactEncx.ID.String(), updateReq, accessToken)

		// Execute request
		resp, err := httpClient.Do(req)
		require.NoError(t, err)
		defer resp.Body.Close()

		// Assert response status - should be 400 Bad Request
		assert.Equal(t, http.StatusBadRequest, resp.StatusCode)

		t.Log("✓ Invalid email format error handled correctly")
	})

	t.Run("Error - Mismatched Contact IDs", func(t *testing.T) {
		ctx := context.Background()

		// Setup administrator authentication
		accessToken := tu.SetupAdminUser(t, ctx, authCtx)
		defer tu.ClearAuthData(t, ctx, authCtx)
		defer ch.ClearContactsTable(t, ctx, testPool)
		defer ch.ClearClientsTable(t, ctx, testPool)

		// Setup client and contact
		clientEncx, contactEncx := setupClientWithContact(t, ctx)
		differentContactID := uuid.New()

		// Create update request with different ID than in URL
		updateReq := client.UpdateContactRequest{
			ID:       differentContactID,
			ClientID: clientEncx.ID.String(),
			Email:    func() *string { s := "test@example.com"; return &s }(),
		}

		// Create HTTP request with contactEncx.ID in URL but differentContactID in request
		req := ch.NewUpdateContactRequest(t, ctx, testServerURL, contactEncx.ID.String(), updateReq, accessToken)

		// Execute request
		resp, err := httpClient.Do(req)
		require.NoError(t, err)
		defer resp.Body.Close()

		// Assert response status - should be 404 Not Found
		assert.Equal(t, http.StatusNotFound, resp.StatusCode)

		t.Log("✓ Mismatched contact IDs error handled correctly")
	})

	t.Run("Authentication/Authorization", func(t *testing.T) {
		ctx := context.Background()

		// Setup client and contact
		clientEncx, contactEncx := setupClientWithContact(t, ctx)

		t.Run("No Authentication", func(t *testing.T) {
			// Create update request
			updateReq := client.UpdateContactRequest{
				ID:       contactEncx.ID,
				ClientID: clientEncx.ID.String(),
				Email:    func() *string { s := "no.auth@example.com"; return &s }(),
			}

			// Create HTTP request without authentication
			req := ch.NewUpdateContactRequest(t, ctx, testServerURL, contactEncx.ID.String(), updateReq, "")

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

			// Create update request
			updateReq := client.UpdateContactRequest{
				ID:       contactEncx.ID,
				ClientID: clientEncx.ID.String(),
				Email:    func() *string { s := "client.auth@example.com"; return &s }(),
			}

			// Create HTTP request with client authentication
			req := ch.NewUpdateContactRequest(t, ctx, testServerURL, contactEncx.ID.String(), updateReq, clientToken)

			// Execute request
			resp, err := httpClient.Do(req)
			require.NoError(t, err)
			defer resp.Body.Close()

			// Client should be forbidden from updating contacts
			assert.Equal(t, http.StatusForbidden, resp.StatusCode)
			t.Log("✓ Client authentication blocked correctly")
		})

		t.Run("Guest Authentication", func(t *testing.T) {
			// Setup guest authentication
			guestToken := tu.SetupGuestUser(t, ctx, authCtx)
			defer tu.ClearAuthData(t, ctx, authCtx)

			// Create update request
			updateReq := client.UpdateContactRequest{
				ID:       contactEncx.ID,
				ClientID: clientEncx.ID.String(),
				Email:    func() *string { s := "guest.auth@example.com"; return &s }(),
			}

			// Create HTTP request with guest authentication
			req := ch.NewUpdateContactRequest(t, ctx, testServerURL, contactEncx.ID.String(), updateReq, guestToken)

			// Execute request
			resp, err := httpClient.Do(req)
			require.NoError(t, err)
			defer resp.Body.Close()

			// Guest should be forbidden from updating contacts
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

		// Use a non-existent contact ID to avoid modifying actual data
		contactID := uuid.New()
		clientID := uuid.New()

		// Create update request
		updateReq := client.UpdateContactRequest{
			ID:       contactID,
			ClientID: clientID.String(),
			Email:    func() *string { s := "cancelled@example.com"; return &s }(),
		}

		// Create HTTP request with cancelled context
		req := ch.NewUpdateContactRequest(t, cancelledCtx, testServerURL, contactID.String(), updateReq, accessToken)

		// Execute request
		_, err := httpClient.Do(req)
		assert.Error(t, err)

		t.Log("✓ Context cancellation handled correctly")
	})
}
