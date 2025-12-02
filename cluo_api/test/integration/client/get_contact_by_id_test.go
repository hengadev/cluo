package client_test

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"testing"
	"time"

	tu "github.com/hengadev/cluo_api/internal/common/testutils"
	"github.com/hengadev/cluo_api/internal/domain/client"
	ch "github.com/hengadev/cluo_api/test/helpers/client"

	"github.com/google/uuid"
	"github.com/hengadev/encx"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// make test-func TEST_NAME=TestGetContactByID TEST_PATH=test/integration/client/get_contact_by_id_test.go

// TestGetContactByID tests all scenarios for getting a contact by ID
func TestGetContactByID(t *testing.T) {
	t.Run("Success", func(t *testing.T) {
		ctx := context.Background()

		// Setup authentication
		accessToken := tu.SetupAdminUser(t, ctx, authCtx)
		defer tu.ClearAuthData(t, ctx, authCtx)
		defer ch.ClearContactsTable(t, ctx, testPool)

		// Generate test client ID and hash
		clientID := uuid.New()
		clientIDBytes, err := encx.SerializeValue(clientID)
		require.NoError(t, err)
		clientIDHash := crypto.HashBasic(ctx, clientIDBytes)

		// Create a client (by inserting an initial contact) so it "exists"
		err = ch.CreateTestClientWithContact(t, ctx, testPool, clientID, clientIDHash)
		require.NoError(t, err)

		// Create test contact data
		contact := client.NewContact(&client.CreateContactRequest{
			ClientID:  clientID,
			Lastname:  "DOE",
			Firstname: "Jane",
			Email:     "jane.doe@example.com",
			Phone:     "0687654321",
			Position:  "Director",
		})

		contactEncx, err := client.ProcessContactEncx(ctx, crypto, contact)
		require.NoError(t, err)

		err = ch.InsertContactEncx(t, ctx, testPool, *contactEncx)
		require.NoError(t, err)
		t.Logf("Created test contact with ID: %s", contact.ID)

		// Create HTTP request
		url := fmt.Sprintf("%s/contact/%s", testServerURL, contact.ID.String())
		req, err := http.NewRequest("GET", url, nil)
		require.NoError(t, err)

		req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", accessToken))

		// Execute request
		httpClient := &http.Client{Timeout: 10 * time.Second}
		resp, err := httpClient.Do(req)
		require.NoError(t, err)
		defer resp.Body.Close()

		// Assert response status
		assert.Equal(t, http.StatusOK, resp.StatusCode)

		// Decode response
		var response struct {
			Message string                  `json:"message"`
			Contact *client.ContactResponse `json:"contact"`
		}
		err = json.NewDecoder(resp.Body).Decode(&response)
		require.NoError(t, err)
		assert.Equal(t, "Contact retrieval by ID completed successfully", response.Message)

		// Verify contact data
		require.NotNil(t, response.Contact, "Contact should not be nil")
		assert.Equal(t, contact.ID.String(), response.Contact.ID, "Contact ID should match")
		assert.Equal(t, clientID.String(), response.Contact.ClientID, "Client ID should match")
		assert.Equal(t, "DOE", response.Contact.Lastname, "Lastname should match")
		assert.Equal(t, "Jane", response.Contact.Firstname, "Firstname should match")
		assert.Equal(t, "jane.doe@example.com", response.Contact.Email, "Email should match")
		assert.Equal(t, "0687654321", response.Contact.Phone, "Phone should match")
		assert.Equal(t, "Director", response.Contact.Position, "Position should match")

		t.Log("✓ Contact retrieved successfully")
	})

	t.Run("ContactNotFound", func(t *testing.T) {
		ctx := context.Background()

		// Setup authentication
		accessToken := tu.SetupAdminUser(t, ctx, authCtx)
		defer tu.ClearAuthData(t, ctx, authCtx)
		defer ch.ClearContactsTable(t, ctx, testPool)

		// Use a non-existent contact ID
		nonExistentContactID := uuid.New()

		// Create HTTP request
		url := fmt.Sprintf("%s/contact/%s", testServerURL, nonExistentContactID.String())
		req, err := http.NewRequest("GET", url, nil)
		require.NoError(t, err)

		req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", accessToken))

		// Execute request
		httpClient := &http.Client{Timeout: 10 * time.Second}
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
		assert.Contains(t, errorResp.Error, "contact")

		t.Log("✓ Contact not found error handled correctly")
	})

	t.Run("InvalidContactID", func(t *testing.T) {
		ctx := context.Background()

		// Setup authentication
		accessToken := tu.SetupAdminUser(t, ctx, authCtx)
		defer tu.ClearAuthData(t, ctx, authCtx)

		// Use invalid UUID format
		invalidContactID := "not-a-valid-uuid"

		// Create HTTP request
		url := fmt.Sprintf("%s/contact/%s", testServerURL, invalidContactID)
		req, err := http.NewRequest("GET", url, nil)
		require.NoError(t, err)

		req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", accessToken))

		// Execute request
		httpClient := &http.Client{Timeout: 10 * time.Second}
		resp, err := httpClient.Do(req)
		require.NoError(t, err)
		defer resp.Body.Close()

		// Assert response status - should be 400 Bad Request
		assert.Equal(t, http.StatusBadRequest, resp.StatusCode)

		t.Log("✓ Invalid contact ID format error handled correctly")
	})

	t.Run("Unauthorized", func(t *testing.T) {
		ctx := context.Background()
		defer ch.ClearContactsTable(t, ctx, testPool)

		// Create test contact data
		contactID := uuid.New()

		// Create HTTP request WITHOUT authentication
		url := fmt.Sprintf("%s/contact/%s", testServerURL, contactID.String())
		req, err := http.NewRequest("GET", url, nil)
		require.NoError(t, err)

		// NO Authorization header

		// Execute request
		httpClient := &http.Client{Timeout: 10 * time.Second}
		resp, err := httpClient.Do(req)
		require.NoError(t, err)
		defer resp.Body.Close()

		// Assert response status - should be 401 Unauthorized
		assert.Equal(t, http.StatusUnauthorized, resp.StatusCode)

		t.Log("✓ Unauthorized access blocked correctly")
	})

	t.Run("ContactWithNilOptionalFields", func(t *testing.T) {
		ctx := context.Background()

		// Setup authentication
		accessToken := tu.SetupAdminUser(t, ctx, authCtx)
		defer tu.ClearAuthData(t, ctx, authCtx)
		defer ch.ClearContactsTable(t, ctx, testPool)

		// Generate test client ID and hash
		clientID := uuid.New()
		clientIDBytes, err := encx.SerializeValue(clientID)
		require.NoError(t, err)
		clientIDHash := crypto.HashBasic(ctx, clientIDBytes)

		// Create a client (by inserting an initial contact) so it "exists"
		err = ch.CreateTestClientWithContact(t, ctx, testPool, clientID, clientIDHash)
		require.NoError(t, err)

		// Create test contact with nil optional fields
		contact := client.NewContact(&client.CreateContactRequest{
			ClientID:  clientID,
			Lastname:  "DOE",
			Firstname: "Jane",
			Email:     "jane.doe@example.com",
			Phone:     "", // Empty phone
			Position:  "", // Empty position
		})

		contactEncx, err := client.ProcessContactEncx(ctx, crypto, contact)
		require.NoError(t, err)

		// Manually set optional fields to nil in the encrypted version
		contactEncx.PhoneEncrypted = nil
		contactEncx.PositionEncrypted = nil

		err = ch.InsertContactEncx(t, ctx, testPool, *contactEncx)
		require.NoError(t, err)
		t.Logf("Created test contact with nil optional fields, ID: %s", contact.ID)

		// Create HTTP request
		url := fmt.Sprintf("%s/contact/%s", testServerURL, contact.ID.String())
		req, err := http.NewRequest("GET", url, nil)
		require.NoError(t, err)

		req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", accessToken))

		// Execute request
		httpClient := &http.Client{Timeout: 10 * time.Second}
		resp, err := httpClient.Do(req)
		require.NoError(t, err)
		defer resp.Body.Close()

		// Assert response status
		assert.Equal(t, http.StatusOK, resp.StatusCode)

		// Decode response
		var response struct {
			Message string                  `json:"message"`
			Contact *client.ContactResponse `json:"contact"`
		}
		err = json.NewDecoder(resp.Body).Decode(&response)
		require.NoError(t, err)

		// Verify contact data with empty strings for nil fields
		require.NotNil(t, response.Contact, "Contact should not be nil")
		assert.Equal(t, contact.ID.String(), response.Contact.ID, "Contact ID should match")
		assert.Equal(t, "", response.Contact.Phone, "Phone should be empty string")
		assert.Equal(t, "", response.Contact.Position, "Position should be empty string")

		t.Log("✓ Contact with nil optional fields retrieved successfully")
	})

	t.Run("ConcurrentAccess", func(t *testing.T) {
		ctx := context.Background()

		// Setup authentication
		accessToken := tu.SetupAdminUser(t, ctx, authCtx)
		defer tu.ClearAuthData(t, ctx, authCtx)
		defer ch.ClearContactsTable(t, ctx, testPool)

		// Generate test client ID and hash
		clientID := uuid.New()
		clientIDBytes, err := encx.SerializeValue(clientID)
		require.NoError(t, err)
		clientIDHash := crypto.HashBasic(ctx, clientIDBytes)

		// Create a client (by inserting an initial contact) so it "exists"
		err = ch.CreateTestClientWithContact(t, ctx, testPool, clientID, clientIDHash)
		require.NoError(t, err)

		// Create multiple test contacts
		numContacts := 5
		contactIDs := make([]uuid.UUID, numContacts)

		for i := 0; i < numContacts; i++ {
			contact := client.NewContact(&client.CreateContactRequest{
				ClientID:  clientID,
				Lastname:  fmt.Sprintf("DOE_%d", i),
				Firstname: fmt.Sprintf("Contact_%d", i),
				Email:     fmt.Sprintf("contact%d@example.com", i),
				Phone:     fmt.Sprintf("06%08d", i),
				Position:  fmt.Sprintf("Position_%d", i),
			})

			contactEncx, err := client.ProcessContactEncx(ctx, crypto, contact)
			require.NoError(t, err)

			err = ch.InsertContactEncx(t, ctx, testPool, *contactEncx)
			require.NoError(t, err)
			contactIDs[i] = contact.ID
		}

		// Create concurrent requests to get all contacts
		errChan := make(chan error, numContacts)
		responseChan := make(chan *client.ContactResponse, numContacts)

		for i, contactID := range contactIDs {
			go func(index int, id uuid.UUID) {
				// Create HTTP request
				url := fmt.Sprintf("%s/contact/%s", testServerURL, id.String())
				req, err := http.NewRequest("GET", url, nil)
				if err != nil {
					errChan <- err
					return
				}

				req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", accessToken))

				// Execute request
				httpClient := &http.Client{Timeout: 10 * time.Second}
				resp, err := httpClient.Do(req)
				if err != nil {
					errChan <- err
					return
				}
				defer resp.Body.Close()

				if resp.StatusCode != http.StatusOK {
					errChan <- fmt.Errorf("unexpected status code: %d", resp.StatusCode)
					return
				}

				// Decode response
				var response struct {
					Message string                  `json:"message"`
					Contact *client.ContactResponse `json:"contact"`
				}
				err = json.NewDecoder(resp.Body).Decode(&response)
				if err != nil {
					errChan <- err
					return
				}

				responseChan <- response.Contact
			}(i, contactID)
		}

		// Wait for all goroutines to complete
		successCount := 0
		errorCount := 0
		responses := make([]*client.ContactResponse, 0, numContacts)

		for i := 0; i < numContacts; i++ {
			select {
			case response := <-responseChan:
				successCount++
				responses = append(responses, response)
			case err := <-errChan:
				errorCount++
				t.Logf("Error in concurrent access: %v", err)
			case <-time.After(15 * time.Second):
				t.Fatal("Timeout waiting for concurrent contact access")
			}
		}

		// Verify all contacts were retrieved successfully
		assert.Equal(t, numContacts, successCount, "All contacts should be retrieved successfully")
		assert.Equal(t, 0, errorCount, "No errors should occur")
		assert.Len(t, responses, numContacts, "Should have received all contact responses")

		t.Log("✓ Concurrent contact access handled correctly")
	})

	t.Run("EmptyContactID", func(t *testing.T) {
		ctx := context.Background()

		// Setup authentication
		accessToken := tu.SetupAdminUser(t, ctx, authCtx)
		defer tu.ClearAuthData(t, ctx, authCtx)

		// Create HTTP request with empty contact ID
		url := fmt.Sprintf("%s/contact/", testServerURL)
		req, err := http.NewRequest("GET", url, nil)
		require.NoError(t, err)

		req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", accessToken))

		// Execute request
		httpClient := &http.Client{Timeout: 10 * time.Second}
		resp, err := httpClient.Do(req)
		require.NoError(t, err)
		defer resp.Body.Close()

		// Assert response status - should be 404 or 400
		assert.True(t, resp.StatusCode == http.StatusNotFound || resp.StatusCode == http.StatusBadRequest,
			"Should return 404 or 400 for empty contact ID")

		t.Log("✓ Empty contact ID handled correctly")
	})

	t.Run("VerifyDecryption", func(t *testing.T) {
		ctx := context.Background()

		// Setup authentication
		accessToken := tu.SetupAdminUser(t, ctx, authCtx)
		defer tu.ClearAuthData(t, ctx, authCtx)
		defer ch.ClearContactsTable(t, ctx, testPool)

		// Generate test client ID and hash
		clientID := uuid.New()
		clientIDBytes, err := encx.SerializeValue(clientID)
		require.NoError(t, err)
		clientIDHash := crypto.HashBasic(ctx, clientIDBytes)

		// Create a client (by inserting an initial contact) so it "exists"
		err = ch.CreateTestClientWithContact(t, ctx, testPool, clientID, clientIDHash)
		require.NoError(t, err)

		// Create test contact with sensitive data
		contact := client.NewContact(&client.CreateContactRequest{
			ClientID:  clientID,
			Lastname:  "SENSITIVE_LASTNAME",
			Firstname: "SENSITIVE_FIRSTNAME",
			Email:     "sensitive@example.com",
			Phone:     "0600000000",
			Position:  "SENSITIVE_POSITION",
		})

		contactEncx, err := client.ProcessContactEncx(ctx, crypto, contact)
		require.NoError(t, err)

		err = ch.InsertContactEncx(t, ctx, testPool, *contactEncx)
		require.NoError(t, err)
		t.Logf("Created test contact with sensitive data, ID: %s", contact.ID)

		// Create HTTP request
		url := fmt.Sprintf("%s/contact/%s", testServerURL, contact.ID.String())
		req, err := http.NewRequest("GET", url, nil)
		require.NoError(t, err)

		req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", accessToken))

		// Execute request
		httpClient := &http.Client{Timeout: 10 * time.Second}
		resp, err := httpClient.Do(req)
		require.NoError(t, err)
		defer resp.Body.Close()

		// Assert response status
		assert.Equal(t, http.StatusOK, resp.StatusCode)

		// Decode response
		var response struct {
			Message string                  `json:"message"`
			Contact *client.ContactResponse `json:"contact"`
		}
		err = json.NewDecoder(resp.Body).Decode(&response)
		require.NoError(t, err)

		// Verify sensitive data is properly decrypted
		require.NotNil(t, response.Contact, "Contact should not be nil")
		assert.Equal(t, "SENSITIVE_LASTNAME", response.Contact.Lastname, "Lastname should be decrypted correctly")
		assert.Equal(t, "SENSITIVE_FIRSTNAME", response.Contact.Firstname, "Firstname should be decrypted correctly")
		assert.Equal(t, "sensitive@example.com", response.Contact.Email, "Email should be decrypted correctly")
		assert.Equal(t, "0600000000", response.Contact.Phone, "Phone should be decrypted correctly")
		assert.Equal(t, "SENSITIVE_POSITION", response.Contact.Position, "Position should be decrypted correctly")

		t.Log("✓ Contact data properly decrypted and returned")
	})
}
