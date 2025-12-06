package client_test

import (
	"context"
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

// make test-func TEST_NAME=TestDeleteContact TEST_PATH=test/integration/client/delete_contact_test.go

// TestDeleteContact tests all scenarios for the DeleteContact endpoint
func TestDeleteContact(t *testing.T) {
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

	t.Run("Success", func(t *testing.T) {
		ctx := context.Background()

		// Setup administrator authentication
		accessToken := tu.SetupAdminUser(t, ctx, authCtx)
		defer tu.ClearAuthData(t, ctx, authCtx)
		defer ch.ClearContactsTable(t, ctx, testPool)
		defer ch.ClearClientsTable(t, ctx, testPool)

		// Setup client and contact
		_, contactEncx := setupClientWithContact(t, ctx)

		// Verify contact exists before deletion
		retrievedContact, err := ch.GetContactEncxByID(t, ctx, testPool, contactEncx.ID)
		require.NoError(t, err, "Contact should exist before deletion")
		assert.Equal(t, contactEncx.ID, retrievedContact.ID, "Contact ID should match")

		// Create HTTP request to delete the contact
		req := ch.NewDeleteContactRequest(t, ctx, testServerURL, contactEncx.ID.String(), accessToken)

		// Execute request
		resp, err := httpClient.Do(req)
		require.NoError(t, err)
		defer resp.Body.Close()

		// Assert response status - should be 200 OK
		assert.Equal(t, http.StatusOK, resp.StatusCode)

		// Verify contact no longer exists in database
		retrievedContact, err = ch.GetContactEncxByID(t, ctx, testPool, contactEncx.ID)
		assert.Error(t, err, "Contact should no longer exist after deletion")
		assert.Equal(t, &client.ContactEncx{}, retrievedContact, "Retrieved contact should be zero after deletion")

		t.Log("✓ Contact deleted successfully")
	})

	t.Run("Contact Not Found", func(t *testing.T) {
		ctx := context.Background()

		// Setup administrator authentication
		accessToken := tu.SetupAdminUser(t, ctx, authCtx)
		defer tu.ClearAuthData(t, ctx, authCtx)
		defer ch.ClearContactsTable(t, ctx, testPool)
		defer ch.ClearClientsTable(t, ctx, testPool)

		// Use a non-existent contact ID
		nonExistentContactID := uuid.New()

		// Create HTTP request to delete non-existent contact
		req := ch.NewDeleteContactRequest(t, ctx, testServerURL, nonExistentContactID.String(), accessToken)

		// Execute request
		resp, err := httpClient.Do(req)
		require.NoError(t, err)
		defer resp.Body.Close()

		// Assert response status - should be 404 Not Found
		assert.Equal(t, http.StatusNotFound, resp.StatusCode)

		t.Log("✓ Contact not found error handled correctly")
	})

	t.Run("Invalid Contact ID", func(t *testing.T) {
		ctx := context.Background()

		// Setup administrator authentication
		accessToken := tu.SetupAdminUser(t, ctx, authCtx)
		defer tu.ClearAuthData(t, ctx, authCtx)

		invalidContactID := "invalid-uuid-format"

		// Create HTTP request with invalid contact ID
		req := ch.NewDeleteContactRequest(t, ctx, testServerURL, invalidContactID, accessToken)

		// Execute request
		resp, err := httpClient.Do(req)
		require.NoError(t, err)
		defer resp.Body.Close()

		// Assert response status - should be 400 Bad Request (due to invalid UUID format)
		assert.Equal(t, http.StatusBadRequest, resp.StatusCode)

		t.Log("✓ Invalid contact ID error handled correctly")
	})

	t.Run("Authentication/Authorization", func(t *testing.T) {
		ctx := context.Background()

		// Setup client and contact
		_, contactEncx := setupClientWithContact(t, ctx)

		t.Run("No Authentication", func(t *testing.T) {
			// Create HTTP request without authentication
			req := ch.NewDeleteContactRequest(t, ctx, testServerURL, contactEncx.ID.String(), "")

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
			req := ch.NewDeleteContactRequest(t, ctx, testServerURL, contactEncx.ID.String(), clientToken)

			// Execute request
			resp, err := httpClient.Do(req)
			require.NoError(t, err)
			defer resp.Body.Close()

			// Client should be forbidden from deleting contacts
			assert.Equal(t, http.StatusForbidden, resp.StatusCode)
			t.Log("✓ Client authentication blocked correctly")
		})

		t.Run("Guest Authentication", func(t *testing.T) {
			// Setup guest authentication
			guestToken := tu.SetupGuestUser(t, ctx, authCtx)
			defer tu.ClearAuthData(t, ctx, authCtx)

			// Create HTTP request with guest authentication
			req := ch.NewDeleteContactRequest(t, ctx, testServerURL, contactEncx.ID.String(), guestToken)

			// Execute request
			resp, err := httpClient.Do(req)
			require.NoError(t, err)
			defer resp.Body.Close()

			// Guest should be forbidden from deleting contacts
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

		// Create HTTP request with cancelled context
		req := ch.NewDeleteContactRequest(t, cancelledCtx, testServerURL, contactID.String(), accessToken)

		// Execute request
		_, err := httpClient.Do(req)
		assert.Error(t, err)
	})

	t.Run("Multiple Contact Deletion", func(t *testing.T) {
		ctx := context.Background()

		// Setup administrator authentication
		accessToken := tu.SetupAdminUser(t, ctx, authCtx)
		defer tu.ClearAuthData(t, ctx, authCtx)
		defer ch.ClearContactsTable(t, ctx, testPool)
		defer ch.ClearClientsTable(t, ctx, testPool)

		// Setup client with multiple contacts
		clientEncx, _ := setupClientWithContact(t, ctx)
		contactIDs := make([]uuid.UUID, 3)

		for i := 0; i < 3; i++ {
			contact := ch.NewTestContact(t)
			contact.ClientID = clientEncx.ID
			contactEncx, err := client.ProcessContactEncx(ctx, crypto, contact)
			require.NoError(t, err)
			err = ch.InsertContactEncx(t, ctx, testPool, contactEncx)
			require.NoError(t, err)
			contactIDs[i] = contactEncx.ID
		}

		// Verify all contacts exist before deletion
		for _, contactID := range contactIDs {
			retrievedContact, err := ch.GetContactEncxByID(t, ctx, testPool, contactID)
			require.NoError(t, err, "Contact %s should exist before deletion", contactID)
			assert.Equal(t, contactID, retrievedContact.ID, "Contact ID should match")
		}

		// Delete each contact
		for _, contactID := range contactIDs {
			req := ch.NewDeleteContactRequest(t, ctx, testServerURL, contactID.String(), accessToken)

			resp, err := httpClient.Do(req)
			require.NoError(t, err)
			defer resp.Body.Close()

			assert.Equal(t, http.StatusOK, resp.StatusCode, "Contact deletion should succeed")
		}

		// Verify all contacts no longer exist
		for _, contactID := range contactIDs {
			retrievedContact, err := ch.GetContactEncxByID(t, ctx, testPool, contactID)
			assert.Error(t, err, "Contact %s should no longer exist after deletion", contactID)
			assert.Equal(t, &client.ContactEncx{}, retrievedContact, "Retrieved contact should be nil after deletion")
		}

		t.Log("✓ Multiple contacts deleted successfully")
	})

	t.Run("Performance with Large Dataset", func(t *testing.T) {
		ctx := context.Background()

		// Setup administrator authentication
		accessToken := tu.SetupAdminUser(t, ctx, authCtx)
		defer tu.ClearAuthData(t, ctx, authCtx)
		defer ch.ClearContactsTable(t, ctx, testPool)
		defer ch.ClearClientsTable(t, ctx, testPool)

		// Setup client with one contact to delete
		clientEncx, contactEncx := setupClientWithContact(t, ctx)

		// Insert 100 other contacts to test performance
		for i := 0; i < 100; i++ {
			contact := ch.NewTestContact(t)
			contact.ClientID = clientEncx.ID
			testContactEncx, err := client.ProcessContactEncx(ctx, crypto, contact)
			require.NoError(t, err)
			err = ch.InsertContactEncx(t, ctx, testPool, testContactEncx)
			require.NoError(t, err)
		}

		// Execute delete request and measure time
		start := time.Now()
		req := ch.NewDeleteContactRequest(t, ctx, testServerURL, contactEncx.ID.String(), accessToken)

		resp, err := httpClient.Do(req)
		duration := time.Since(start)

		require.NoError(t, err)
		defer resp.Body.Close()

		// Assert response status
		assert.Equal(t, http.StatusOK, resp.StatusCode)

		// Log performance
		t.Logf("✓ Deleted contact in %v (dataset with 100 other contacts)", duration)

		// Performance assertion - should complete within reasonable time
		assert.Less(t, duration, 2*time.Second, "Delete request should complete within 2 seconds")
	})
}

