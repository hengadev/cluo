package client_test

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"testing"
	"time"

	"github.com/hengadev/cluo_api/internal/common/auth/cookies"
	"github.com/hengadev/cluo_api/internal/common/errs"
	tu "github.com/hengadev/cluo_api/internal/common/testutils"
	"github.com/hengadev/cluo_api/internal/domain/client"
	ch "github.com/hengadev/cluo_api/test/helpers/client"

	"github.com/google/uuid"
	"github.com/hengadev/encx"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// make test-func TEST_NAME=TestCreateContact TEST_PATH=test/integration/client/create_contact_test.go

// TestCreateContact tests all scenarios for creating a contact
func TestCreateContact(t *testing.T) {
	setupClient := func(t *testing.T, ctx context.Context) uuid.UUID {
		c := ch.NewTestClient(t)
		clientEncx, err := client.ProcessClientEncx(ctx, crypto, c)
		require.NoError(t, err)
		err = ch.InsertClientEncx(t, ctx, testPool, clientEncx)
		require.NoError(t, err)
		return c.ID
	}

	t.Run("Success", func(t *testing.T) {
		ctx := context.Background()

		// Setup authentication
		accessToken := tu.SetupAdminUser(t, ctx, authCtx)
		defer tu.ClearAuthData(t, ctx, authCtx)
		defer ch.ClearContactsTable(t, ctx, testPool)
		defer ch.ClearClientsTable(t, ctx, testPool)

		// Generate test client
		clientID := setupClient(t, ctx)

		// Prepare request payload
		payload := client.CreateContactRequest{
			Lastname:  "DOE",
			Firstname: "Jane",
			Email:     "jane.doe@example.com",
			Phone:     "0687654321",
			Position:  "Director",
		}

		// Create HTTP request using the test helper
		req := ch.NewCreateContactRequest(t, ctx, testServerURL, clientID.String(), payload, accessToken)

		// Execute request
		httpClient := &http.Client{Timeout: 10 * time.Second}
		resp, err := httpClient.Do(req)
		require.NoError(t, err)
		defer resp.Body.Close()

		// Assert response status
		assert.Equal(t, http.StatusOK, resp.StatusCode)

		// Decode response
		var response *client.ContactResponse
		err = json.NewDecoder(resp.Body).Decode(&response)
		require.NoError(t, err)
		assert.Equal(t, payload.Lastname, response.Lastname)
		assert.Equal(t, payload.Firstname, response.Firstname)
		assert.Equal(t, payload.Email, response.Email)
		assert.Equal(t, payload.Phone, response.Phone)
		assert.Equal(t, payload.Position, response.Position)
		assert.Equal(t, clientID.String(), response.ClientID)

		t.Log("✓ Contact created successfully")
	})

	t.Run("ClientNotFound", func(t *testing.T) {
		ctx := context.Background()

		// Setup authentication
		accessToken := tu.SetupAdminUser(t, ctx, authCtx)
		defer tu.ClearAuthData(t, ctx, authCtx)
		defer ch.ClearContactsTable(t, ctx, testPool)
		defer ch.ClearClientsTable(t, ctx, testPool)

		// Use a non-existent client ID
		nonExistentClientID := uuid.New()

		// Prepare request payload
		payload := client.CreateContactRequest{
			Lastname:  "DOE",
			Firstname: "John",
			Email:     "john.doe@example.com",
			Phone:     "0612345678",
			Position:  "Manager",
		}

		// Create HTTP request using the test helper
		req := ch.NewCreateContactRequest(t, ctx, testServerURL, nonExistentClientID.String(), payload, accessToken)

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
		assert.Contains(t, errorResp.Error, errs.ErrRepositoryNotFound.Error())

		t.Log("✓ Client not found error handled correctly")
	})

	t.Run("InvalidPayload", func(t *testing.T) {
		ctx := context.Background()

		// Setup authentication
		accessToken := tu.SetupAdminUser(t, ctx, authCtx)
		defer tu.ClearAuthData(t, ctx, authCtx)

		clientID := setupClient(t, ctx)

		// Create HTTP request with invalid JSON manually since helper requires valid struct
		invalidJSON := []byte(`{"lastname": "DOE", "firstname": "John", invalid}`)
		req, err := http.NewRequest("POST", testServerURL+"/client/"+clientID.String()+"/contact", bytes.NewBuffer(invalidJSON))
		require.NoError(t, err)

		req.Header.Set("Content-Type", "application/json")

		// Add authentication cookie manually
		cookie := &http.Cookie{
			Name:  cookies.AccessTokenCookieName,
			Value: accessToken,
		}
		req.AddCookie(cookie)

		// Execute request
		httpClient := &http.Client{Timeout: 10 * time.Second}
		resp, err := httpClient.Do(req)
		require.NoError(t, err)
		defer resp.Body.Close()

		// Assert response status - should be 400 Bad Request
		assert.Equal(t, http.StatusBadRequest, resp.StatusCode)

		// Decode error response
		var errorResp struct {
			Error string `json:"error"`
		}
		err = json.NewDecoder(resp.Body).Decode(&errorResp)
		require.NoError(t, err)

		assert.Contains(t, errorResp.Error, "invalid")

		t.Log("✓ Invalid JSON payload error handled correctly")
	})

	t.Run("MissingRequiredFields", func(t *testing.T) {
		ctx := context.Background()

		// Setup authentication
		accessToken := tu.SetupAdminUser(t, ctx, authCtx)
		defer tu.ClearAuthData(t, ctx, authCtx)
		defer ch.ClearContactsTable(t, ctx, testPool)
		defer ch.ClearClientsTable(t, ctx, testPool)

		// Generate test client ID
		clientID := setupClient(t, ctx)

		testCases := []struct {
			name    string
			payload client.CreateContactRequest
		}{
			{
				name: "MissingLastname",
				payload: client.CreateContactRequest{
					ClientID:  clientID,
					Firstname: "John",
					Email:     "john@example.com",
					Phone:     "0612345678",
					Position:  "Manager",
				},
			},
			{
				name: "MissingFirstname",
				payload: client.CreateContactRequest{
					ClientID: clientID,
					Lastname: "DOE",
					Email:    "john@example.com",
					Phone:    "0612345678",
					Position: "Manager",
				},
			},
			{
				name: "MissingEmail",
				payload: client.CreateContactRequest{
					ClientID:  clientID,
					Lastname:  "DOE",
					Firstname: "John",
					Phone:     "0612345678",
					Position:  "Manager",
				},
			},
		}

		for _, tc := range testCases {
			t.Run(tc.name, func(t *testing.T) {
				// Create HTTP request using the test helper
				req := ch.NewCreateContactRequest(t, ctx, testServerURL, clientID.String(), tc.payload, accessToken)

				// Execute request
				httpClient := &http.Client{Timeout: 10 * time.Second}
				resp, err := httpClient.Do(req)
				require.NoError(t, err)
				defer resp.Body.Close()

				// Assert response status - should be 400 Bad Request
				assert.Equal(t, http.StatusBadRequest, resp.StatusCode)

				t.Logf("✓ %s validation error handled correctly", tc.name)
			})
		}
	})

	t.Run("Unauthorized", func(t *testing.T) {
		ctx := context.Background()
		defer ch.ClearContactsTable(t, ctx, testPool)
		defer ch.ClearClientsTable(t, ctx, testPool)

		clientID := uuid.New()

		// Prepare request payload
		payload := client.CreateContactRequest{
			ClientID:  clientID,
			Lastname:  "DOE",
			Firstname: "John",
			Email:     "john.doe@example.com",
			Phone:     "0612345678",
			Position:  "Manager",
		}

		// Create HTTP request WITHOUT authentication using the test helper
		req := ch.NewCreateContactRequest(t, ctx, testServerURL, clientID.String(), payload, "")

		// Execute request
		httpClient := &http.Client{Timeout: 10 * time.Second}
		resp, err := httpClient.Do(req)
		require.NoError(t, err)
		defer resp.Body.Close()

		// Assert response status - should be 401 Unauthorized
		assert.Equal(t, http.StatusUnauthorized, resp.StatusCode)

		t.Log("✓ Unauthorized access blocked correctly")
	})

	t.Run("InvalidClientID", func(t *testing.T) {
		ctx := context.Background()

		// Setup authentication
		accessToken := tu.SetupAdminUser(t, ctx, authCtx)
		defer tu.ClearAuthData(t, ctx, authCtx)

		invalidClientID := uuid.UUID{}

		// Prepare request payload
		payload := client.CreateContactRequest{
			ClientID:  invalidClientID,
			Lastname:  "DOE",
			Firstname: "John",
			Email:     "john.doe@example.com",
			Phone:     "0612345678",
			Position:  "Manager",
		}

		// Create HTTP request using the test helper
		req := ch.NewCreateContactRequest(t, ctx, testServerURL, invalidClientID.String(), payload, accessToken)

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
		assert.Contains(t, errorResp.Error, errs.ErrRepositoryNotFound.Error())

		t.Log("✓ Invalid client ID format error handled correctly")
	})

	t.Run("VerifyEncryption", func(t *testing.T) {
		ctx := context.Background()

		// Setup authentication
		accessToken := tu.SetupAdminUser(t, ctx, authCtx)
		defer tu.ClearAuthData(t, ctx, authCtx)
		defer ch.ClearContactsTable(t, ctx, testPool)
		defer ch.ClearClientsTable(t, ctx, testPool)

		// Generate test client
		clientID := setupClient(t, ctx)

		// Sensitive data that should be encrypted
		sensitiveData := client.CreateContactRequest{
			ClientID:  clientID,
			Lastname:  "SENSITIVE_LASTNAME",
			Firstname: "SENSITIVE_FIRSTNAME",
			Email:     "sensitive@example.com",
			Phone:     "0600000000",
			Position:  "SENSITIVE_POSITION",
		}

		// Create HTTP request using the test helper
		req := ch.NewCreateContactRequest(t, ctx, testServerURL, clientID.String(), sensitiveData, accessToken)

		// Execute request
		httpClient := &http.Client{Timeout: 10 * time.Second}
		resp, err := httpClient.Do(req)
		require.NoError(t, err)
		defer resp.Body.Close()

		// Assert success
		assert.Equal(t, http.StatusOK, resp.StatusCode)

		// Calculate email hash
		emailBytes, err := encx.SerializeValue(sensitiveData.Email)
		require.NoError(t, err)
		emailHash := crypto.HashBasic(ctx, emailBytes)

		// Retrieve the contact from database
		contactEncx, err := ch.GetContactEncxByEmailHash(t, ctx, testPool, emailHash)
		require.NoError(t, err)

		// Verify data is encrypted (should not match plaintext)
		assert.NotContains(t, string(contactEncx.LastnameEncrypted), "SENSITIVE_LASTNAME",
			"Lastname should be encrypted, not plaintext")
		assert.NotContains(t, string(contactEncx.FirstnameEncrypted), "SENSITIVE_FIRSTNAME",
			"Firstname should be encrypted, not plaintext")
		assert.NotContains(t, string(contactEncx.EmailEncrypted), "sensitive@example.com",
			"Email should be encrypted, not plaintext")
		assert.NotContains(t, string(contactEncx.PhoneEncrypted), "0600000000",
			"Phone should be encrypted, not plaintext")
		assert.NotContains(t, string(contactEncx.PositionEncrypted), "SENSITIVE_POSITION",
			"Position should be encrypted, not plaintext")

		// Verify hashes are set
		assert.NotEmpty(t, contactEncx.ClientID, "Client ID hash should be set")
		assert.NotEmpty(t, contactEncx.EmailHash, "Email hash should be set")

		// Verify encrypted fields are not empty
		assert.NotEmpty(t, contactEncx.LastnameEncrypted, "Lastname should be encrypted")
		assert.NotEmpty(t, contactEncx.FirstnameEncrypted, "Firstname should be encrypted")
		assert.NotEmpty(t, contactEncx.EmailEncrypted, "Email should be encrypted")
		assert.NotEmpty(t, contactEncx.PhoneEncrypted, "Phone should be encrypted")
		assert.NotEmpty(t, contactEncx.PositionEncrypted, "Position should be encrypted")
		assert.NotEmpty(t, contactEncx.DEKEncrypted, "DEK should be encrypted")

		t.Log("✓ Contact data properly encrypted and stored")
	})

	t.Run("ConcurrentCreation", func(t *testing.T) {
		ctx := context.Background()

		// Setup authentication
		accessToken := tu.SetupAdminUser(t, ctx, authCtx)
		defer tu.ClearAuthData(t, ctx, authCtx)
		defer ch.ClearContactsTable(t, ctx, testPool)
		defer ch.ClearClientsTable(t, ctx, testPool)

		// Generate test client
		clientID := setupClient(t, ctx)

		// Create multiple contacts concurrently
		numContacts := 5
		errChan := make(chan error, numContacts)
		doneChan := make(chan bool, numContacts)

		for i := 0; i < numContacts; i++ {
			go func(index int) {
				payload := client.CreateContactRequest{
					ClientID:  clientID,
					Lastname:  fmt.Sprintf("DOE_%d", index),
					Firstname: fmt.Sprintf("Contact_%d", index),
					Email:     fmt.Sprintf("contact%d@example.com", index),
					Phone:     fmt.Sprintf("06%08d", index),
					Position:  fmt.Sprintf("Position_%d", index),
				}

				// Create request manually in goroutine since NewCreateContactRequest uses t.Helper()
				jsonBody, err := json.Marshal(payload)
				if err != nil {
					errChan <- err
					return
				}

				req, err := http.NewRequest(
					"POST",
					testServerURL+"/client/"+clientID.String()+"/contact",
					bytes.NewBuffer(jsonBody),
				)
				if err != nil {
					errChan <- err
					return
				}

				req.Header.Set("Content-Type", "application/json")
				cookie := &http.Cookie{
					Name:  cookies.AccessTokenCookieName,
					Value: accessToken,
				}
				req.AddCookie(cookie)

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

				doneChan <- true
			}(i)
		}

		// Wait for all goroutines to complete
		successCount := 0
		errorCount := 0
		for i := 0; i < numContacts; i++ {
			select {
			case <-doneChan:
				successCount++
			case err := <-errChan:
				errorCount++
				t.Logf("Error in concurrent creation: %v", err)
			case <-time.After(15 * time.Second):
				t.Fatal("Timeout waiting for concurrent contact creation")
			}
		}

		// Verify all contacts were created successfully
		assert.Equal(t, numContacts, successCount, "All contacts should be created successfully")
		assert.Equal(t, 0, errorCount, "No errors should occur")

		// Verify count in database (initial contact + numContacts)
		count, err := ch.CountContactsByClientID(t, ctx, testPool, clientID)
		require.NoError(t, err)
		assert.Equal(t, numContacts, count, "Should have all newly created contacts")

		t.Log("✓ Concurrent contact creation handled correctly")
	})

	t.Run("MissingClientIDInURL", func(t *testing.T) {
		ctx := context.Background()

		// Setup authentication
		accessToken := tu.SetupAdminUser(t, ctx, authCtx)
		defer tu.ClearAuthData(t, ctx, authCtx)

		// Prepare request payload
		payload := client.CreateContactRequest{
			ClientID:  uuid.New(), // Still need a valid UUID for the payload
			Lastname:  "DOE",
			Firstname: "John",
			Email:     "john.doe@example.com",
			Phone:     "0612345678",
			Position:  "Manager",
		}

		// Create HTTP request manually with empty client ID in URL
		jsonBody, err := json.Marshal(payload)
		require.NoError(t, err)

		req, err := http.NewRequest("POST", testServerURL+"/client//contact", bytes.NewBuffer(jsonBody))
		require.NoError(t, err)

		req.Header.Set("Content-Type", "application/json")
		cookie := &http.Cookie{
			Name:  cookies.AccessTokenCookieName,
			Value: accessToken,
		}
		req.AddCookie(cookie)

		// Execute request
		httpClient := &http.Client{Timeout: 10 * time.Second}
		resp, err := httpClient.Do(req)
		require.NoError(t, err)
		defer resp.Body.Close()

		// Assert response status - should be 404 or 400
		assert.True(t, resp.StatusCode == http.StatusNotFound || resp.StatusCode == http.StatusBadRequest,
			"Should return 404 or 400 for missing client ID in URL")

		t.Log("✓ Missing client ID in URL handled correctly")
	})

	t.Run("DatabaseConnectionFailure", func(t *testing.T) {
		// This test would require mocking the database connection
		// or using a special test mode that can simulate database failures
		// For now, we skip this test as it requires infrastructure changes
		t.Skip("Skipping database connection failure test - requires infrastructure for simulating DB failures")
	})
}
