package investigation_test

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"testing"
	"time"

	"github.com/hengadev/cluo_api/internal/common/auth/cookies"
	tu "github.com/hengadev/cluo_api/internal/common/testutils"
	"github.com/hengadev/cluo_api/internal/domain/investigation"
	clientDomain "github.com/hengadev/cluo_api/internal/domain/client"
	ch "github.com/hengadev/cluo_api/test/helpers/investigation"
	clientHelpers "github.com/hengadev/cluo_api/test/helpers/client"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// make test-func TEST_NAME=TestListCasesByClient TEST_PATH=test/integration/case/list_cases_by_client_test.go

// TestListCasesByClient tests all scenarios for listing cases by client ID
func TestListCasesByClient(t *testing.T) {
	httpClient := &http.Client{Timeout: 10 * time.Second}

	setupClient := func(t *testing.T, ctx context.Context) uuid.UUID {
		c := clientHelpers.NewTestClient(t)
		clientEncx, err := clientDomain.ProcessClientEncx(ctx, crypto, c)
		require.NoError(t, err)
		err = clientHelpers.InsertClientEncx(t, ctx, testPool, clientEncx)
		require.NoError(t, err)
		return c.ID
	}

	setupContact := func(t *testing.T, ctx context.Context, clientID uuid.UUID) uuid.UUID {
		contact := clientHelpers.NewTestContact(t)
		contact.ClientID = clientID
		contactEncx, err := clientDomain.ProcessContactEncx(ctx, crypto, contact)
		require.NoError(t, err)
		err = clientHelpers.InsertContactEncx(t, ctx, testPool, contactEncx)
		require.NoError(t, err)
		return contact.ID
	}

	setupCase := func(t *testing.T, ctx context.Context, clientID uuid.UUID, assignedContactID *uuid.UUID) *investigation.Investigation {
		caseData := ch.NewTestCase(t)
		caseData.ClientID = clientID
		caseData.AssignedContactID = assignedContactID
		caseEncx, err := investigation.ProcessInvestigationEncx(ctx, crypto, caseData)
		require.NoError(t, err)
		err = ch.InsertCaseEncx(t, ctx, testPool, caseEncx)
		require.NoError(t, err)
		return caseData
	}

	t.Run("EmptyList", func(t *testing.T) {
		ctx := context.Background()

		// Setup authentication
		accessToken := tu.SetupAdminUser(t, ctx, authCtx)
		defer tu.ClearAuthData(t, ctx, authCtx)

		// Clear all test data
		ch.ClearCasesTable(t, ctx, testPool)
		clientHelpers.ClearContactsTable(t, ctx, testPool)
		clientHelpers.ClearClientsTable(t, ctx, testPool)

		// Setup client (but no cases)
		clientID := setupClient(t, ctx)

		// Create request
		req := ch.NewListCasesByClientRequest(t, ctx, testServerURL, clientID, 1, 20, accessToken)

		// Make request
		resp, err := httpClient.Do(req)
		require.NoError(t, err)
		defer resp.Body.Close()

		// Verify response
		assert.Equal(t, http.StatusOK, resp.StatusCode)

		var response investigation.ListCasesResponse
		err = json.NewDecoder(resp.Body).Decode(&response)
		require.NoError(t, err)

		assert.Empty(t, response.Cases)
		assert.Equal(t, 1, response.Pagination.Page)
		assert.Equal(t, 20, response.Pagination.PageSize)
		assert.Equal(t, 0, response.Pagination.TotalItems)
		assert.Equal(t, 1, response.Pagination.TotalPages)
	})

	t.Run("ListWithCases", func(t *testing.T) {
		ctx := context.Background()

		// Setup authentication
		accessToken := tu.SetupAdminUser(t, ctx, authCtx)
		defer tu.ClearAuthData(t, ctx, authCtx)

		// Clear all test data
		ch.ClearCasesTable(t, ctx, testPool)
		clientHelpers.ClearContactsTable(t, ctx, testPool)
		clientHelpers.ClearClientsTable(t, ctx, testPool)

		// Setup test data
		clientID1 := setupClient(t, ctx)
		clientID2 := setupClient(t, ctx)
		contactID1 := setupContact(t, ctx, clientID1)

		now := time.Now()

		// Create cases for client1 with different timestamps
		case3Data := ch.NewTestCase(t)
		case3Data.ClientID = clientID1
		case3Data.AssignedContactID = &contactID1
		case3Data.Title = "Case 3"
		case3Data.CreatedAt = now.Add(-1 * time.Hour)
		case3Encx, err := investigation.ProcessInvestigationEncx(ctx, crypto, case3Data)
		require.NoError(t, err)
		err = ch.InsertCaseEncx(t, ctx, testPool, case3Encx)
		require.NoError(t, err)

		case2Data := ch.NewTestCase(t)
		case2Data.ClientID = clientID1
		case2Data.AssignedContactID = nil
		case2Data.Title = "Case 2"
		case2Data.CreatedAt = now.Add(-2 * time.Hour)
		case2Encx, err := investigation.ProcessInvestigationEncx(ctx, crypto, case2Data)
		require.NoError(t, err)
		err = ch.InsertCaseEncx(t, ctx, testPool, case2Encx)
		require.NoError(t, err)

		case1Data := ch.NewTestCase(t)
		case1Data.ClientID = clientID2
		case1Data.AssignedContactID = nil
		case1Data.Title = "Case 1"
		case1Data.CreatedAt = now.Add(-3 * time.Hour)
		case1Encx, err := investigation.ProcessInvestigationEncx(ctx, crypto, case1Data)
		require.NoError(t, err)
		err = ch.InsertCaseEncx(t, ctx, testPool, case1Encx)
		require.NoError(t, err)

		// Create request for client1
		req := ch.NewListCasesByClientRequest(t, ctx, testServerURL, clientID1, 1, 20, accessToken)

		// Make request
		resp, err := httpClient.Do(req)
		require.NoError(t, err)
		defer resp.Body.Close()

		// Verify response
		assert.Equal(t, http.StatusOK, resp.StatusCode)

		var response investigation.ListCasesResponse
		err = json.NewDecoder(resp.Body).Decode(&response)
		require.NoError(t, err)

		assert.Len(t, response.Cases, 2) // Only cases for client1
		assert.Equal(t, 1, response.Pagination.Page)
		assert.Equal(t, 20, response.Pagination.PageSize)
		assert.Equal(t, 2, response.Pagination.TotalItems)
		assert.Equal(t, 1, response.Pagination.TotalPages)

		// Verify all cases belong to client1
		for _, caseResp := range response.Cases {
			assert.Equal(t, clientID1.String(), caseResp.ClientID)
			assert.NotEmpty(t, caseResp.CaseType, "CaseType should be present")
		}

		// Verify ordering (should be by created_at DESC)
		assert.Equal(t, case3Data.ID.String(), response.Cases[0].ID)
		assert.Equal(t, case2Data.ID.String(), response.Cases[1].ID)
	})

	t.Run("MultipleClientsIsolation", func(t *testing.T) {
		ctx := context.Background()

		// Setup authentication
		accessToken := tu.SetupAdminUser(t, ctx, authCtx)
		defer tu.ClearAuthData(t, ctx, authCtx)

		// Clear all test data
		ch.ClearCasesTable(t, ctx, testPool)
		clientHelpers.ClearContactsTable(t, ctx, testPool)
		clientHelpers.ClearClientsTable(t, ctx, testPool)

		// Setup test data
		clientID1 := setupClient(t, ctx)
		clientID2 := setupClient(t, ctx)
		clientID3 := setupClient(t, ctx)

		// Create cases for each client
		for i := 0; i < 5; i++ {
			setupCase(t, ctx, clientID1, nil)
			setupCase(t, ctx, clientID2, nil)
			setupCase(t, ctx, clientID3, nil)
		}

		// Request cases for client1
		req1 := ch.NewListCasesByClientRequest(t, ctx, testServerURL, clientID1, 1, 10, accessToken)

		resp1, err := httpClient.Do(req1)
		require.NoError(t, err)
		defer resp1.Body.Close()

		assert.Equal(t, http.StatusOK, resp1.StatusCode)

		var response1 investigation.ListCasesResponse
		err = json.NewDecoder(resp1.Body).Decode(&response1)
		require.NoError(t, err)

		assert.Len(t, response1.Cases, 5)
		assert.Equal(t, 5, response1.Pagination.TotalItems)

		// Verify all cases belong to client1
		for _, caseResp := range response1.Cases {
			assert.Equal(t, clientID1.String(), caseResp.ClientID)
		}

		// Request cases for client2
		req2 := ch.NewListCasesByClientRequest(t, ctx, testServerURL, clientID2, 1, 10, accessToken)

		resp2, err := httpClient.Do(req2)
		require.NoError(t, err)
		defer resp2.Body.Close()

		assert.Equal(t, http.StatusOK, resp2.StatusCode)

		var response2 investigation.ListCasesResponse
		err = json.NewDecoder(resp2.Body).Decode(&response2)
		require.NoError(t, err)

		assert.Len(t, response2.Cases, 5)
		assert.Equal(t, 5, response2.Pagination.TotalItems)

		// Verify all cases belong to client2
		for _, caseResp := range response2.Cases {
			assert.Equal(t, clientID2.String(), caseResp.ClientID)
		}

		// Request cases for client3
		req3 := ch.NewListCasesByClientRequest(t, ctx, testServerURL, clientID3, 1, 10, accessToken)

		resp3, err := httpClient.Do(req3)
		require.NoError(t, err)
		defer resp3.Body.Close()

		assert.Equal(t, http.StatusOK, resp3.StatusCode)

		var response3 investigation.ListCasesResponse
		err = json.NewDecoder(resp3.Body).Decode(&response3)
		require.NoError(t, err)

		assert.Len(t, response3.Cases, 5)
		assert.Equal(t, 5, response3.Pagination.TotalItems)

		// Verify all cases belong to client3
		for _, caseResp := range response3.Cases {
			assert.Equal(t, clientID3.String(), caseResp.ClientID)
		}
	})

	t.Run("Pagination", func(t *testing.T) {
		ctx := context.Background()

		// Setup authentication
		accessToken := tu.SetupAdminUser(t, ctx, authCtx)
		defer tu.ClearAuthData(t, ctx, authCtx)

		// Clear all test data
		ch.ClearCasesTable(t, ctx, testPool)
		clientHelpers.ClearContactsTable(t, ctx, testPool)
		clientHelpers.ClearClientsTable(t, ctx, testPool)

		// Setup test data
		clientID := setupClient(t, ctx)

		// Create 7 cases for pagination testing
		for i := 0; i < 7; i++ {
			caseData := ch.NewTestCase(t)
			caseData.ClientID = clientID
			caseData.AssignedContactID = nil
			caseData.Title = fmt.Sprintf("Test Case %d", i+1)
			caseEncx, err := investigation.ProcessInvestigationEncx(ctx, crypto, caseData)
			require.NoError(t, err)
			err = ch.InsertCaseEncx(t, ctx, testPool, caseEncx)
			require.NoError(t, err)
		}

		// Test first page with 3 items
		req := ch.NewListCasesByClientRequest(t, ctx, testServerURL, clientID, 1, 3, accessToken)

		resp, err := httpClient.Do(req)
		require.NoError(t, err)
		defer resp.Body.Close()

		assert.Equal(t, http.StatusOK, resp.StatusCode)

		var response investigation.ListCasesResponse
		err = json.NewDecoder(resp.Body).Decode(&response)
		require.NoError(t, err)

		assert.Len(t, response.Cases, 3)
		assert.Equal(t, 1, response.Pagination.Page)
		assert.Equal(t, 3, response.Pagination.PageSize)
		assert.Equal(t, 7, response.Pagination.TotalItems)
		assert.Equal(t, 3, response.Pagination.TotalPages)

		// Test second page with 3 items
		req = ch.NewListCasesByClientRequest(t, ctx, testServerURL, clientID, 2, 3, accessToken)

		resp, err = httpClient.Do(req)
		require.NoError(t, err)
		defer resp.Body.Close()

		assert.Equal(t, http.StatusOK, resp.StatusCode)

		err = json.NewDecoder(resp.Body).Decode(&response)
		require.NoError(t, err)

		assert.Len(t, response.Cases, 3)
		assert.Equal(t, 2, response.Pagination.Page)
		assert.Equal(t, 3, response.Pagination.PageSize)
		assert.Equal(t, 7, response.Pagination.TotalItems)
		assert.Equal(t, 3, response.Pagination.TotalPages)

		// Test third page with 1 item
		req = ch.NewListCasesByClientRequest(t, ctx, testServerURL, clientID, 3, 3, accessToken)

		resp, err = httpClient.Do(req)
		require.NoError(t, err)
		defer resp.Body.Close()

		assert.Equal(t, http.StatusOK, resp.StatusCode)

		err = json.NewDecoder(resp.Body).Decode(&response)
		require.NoError(t, err)

		assert.Len(t, response.Cases, 1)
		assert.Equal(t, 3, response.Pagination.Page)
		assert.Equal(t, 3, response.Pagination.PageSize)
		assert.Equal(t, 7, response.Pagination.TotalItems)
		assert.Equal(t, 3, response.Pagination.TotalPages)

		// Test fourth page (should be empty)
		req = ch.NewListCasesByClientRequest(t, ctx, testServerURL, clientID, 4, 3, accessToken)

		resp, err = httpClient.Do(req)
		require.NoError(t, err)
		defer resp.Body.Close()

		assert.Equal(t, http.StatusOK, resp.StatusCode)

		err = json.NewDecoder(resp.Body).Decode(&response)
		require.NoError(t, err)

		assert.Empty(t, response.Cases)
		assert.Equal(t, 4, response.Pagination.Page)
		assert.Equal(t, 3, response.Pagination.PageSize)
		assert.Equal(t, 7, response.Pagination.TotalItems)
		assert.Equal(t, 3, response.Pagination.TotalPages)
	})

	t.Run("AssignedContactHandling", func(t *testing.T) {
		ctx := context.Background()

		// Setup authentication
		accessToken := tu.SetupAdminUser(t, ctx, authCtx)
		defer tu.ClearAuthData(t, ctx, authCtx)

		// Clear all test data
		ch.ClearCasesTable(t, ctx, testPool)
		clientHelpers.ClearContactsTable(t, ctx, testPool)
		clientHelpers.ClearClientsTable(t, ctx, testPool)

		// Setup test data
		clientID := setupClient(t, ctx)
		contactID1 := setupContact(t, ctx, clientID)
		contactID2 := setupContact(t, ctx, clientID)

		// Create cases with different assigned contacts
		caseWithContact1 := ch.NewTestCase(t)
		caseWithContact1.ClientID = clientID
		caseWithContact1.AssignedContactID = &contactID1
		caseWithContact1.Title = "Case with Contact 1"
		caseWithContact1Encx, err := investigation.ProcessInvestigationEncx(ctx, crypto, caseWithContact1)
		require.NoError(t, err)
		err = ch.InsertCaseEncx(t, ctx, testPool, caseWithContact1Encx)
		require.NoError(t, err)

		caseWithoutContact := ch.NewTestCase(t)
		caseWithoutContact.ClientID = clientID
		caseWithoutContact.AssignedContactID = nil
		caseWithoutContact.Title = "Case without Contact"
		caseWithoutContactEncx, err := investigation.ProcessInvestigationEncx(ctx, crypto, caseWithoutContact)
		require.NoError(t, err)
		err = ch.InsertCaseEncx(t, ctx, testPool, caseWithoutContactEncx)
		require.NoError(t, err)

		caseWithContact2 := ch.NewTestCase(t)
		caseWithContact2.ClientID = clientID
		caseWithContact2.AssignedContactID = &contactID2
		caseWithContact2.Title = "Case with Contact 2"
		caseWithContact2Encx, err := investigation.ProcessInvestigationEncx(ctx, crypto, caseWithContact2)
		require.NoError(t, err)
		err = ch.InsertCaseEncx(t, ctx, testPool, caseWithContact2Encx)
		require.NoError(t, err)

		// Create request
		req := ch.NewListCasesByClientRequest(t, ctx, testServerURL, clientID, 1, 20, accessToken)

		// Make request
		resp, err := httpClient.Do(req)
		require.NoError(t, err)
		defer resp.Body.Close()

		// Verify response
		assert.Equal(t, http.StatusOK, resp.StatusCode)

		var response investigation.ListCasesResponse
		err = json.NewDecoder(resp.Body).Decode(&response)
		require.NoError(t, err)

		assert.Len(t, response.Cases, 3)
		assert.Equal(t, 3, response.Pagination.TotalItems)

		// Verify assigned contact handling
		case1 := findCaseByTitle(response.Cases, "Case with Contact 1")
		case2 := findCaseByTitle(response.Cases, "Case without Contact")
		case3 := findCaseByTitle(response.Cases, "Case with Contact 2")

		require.NotNil(t, case1, "Should find case with Contact 1")
		assert.NotNil(t, case2, "Should find case without Contact")
		assert.NotNil(t, case3, "Should find case with Contact 2")

		assert.Equal(t, contactID1.String(), *case1.AssignedContactID)
		assert.Nil(t, case2.AssignedContactID, "Should have no assigned contact")
		assert.Equal(t, contactID2.String(), *case3.AssignedContactID)
	})

	t.Run("InvalidClientID", func(t *testing.T) {
		ctx := context.Background()

		// Setup authentication
		accessToken := tu.SetupAdminUser(t, ctx, authCtx)
		defer tu.ClearAuthData(t, ctx, authCtx)

		// Test with invalid UUID
		invalidClientID := uuid.New().String() + "invalid"

		req, err := http.NewRequestWithContext(
			ctx,
			http.MethodGet,
			testServerURL+"/clients/"+invalidClientID+"/cases",
			nil,
		)
		require.NoError(t, err)

		if accessToken != "" {
			cookie := &http.Cookie{
				Name:  cookies.AccessTokenCookieName,
				Value: accessToken,
			}
			req.AddCookie(cookie)
		}

		// Make request
		resp, err := httpClient.Do(req)
		require.NoError(t, err)
		defer resp.Body.Close()

		// Should return bad request due to invalid UUID format
		assert.Equal(t, http.StatusBadRequest, resp.StatusCode)
	})

	t.Run("NonExistentClient", func(t *testing.T) {
		ctx := context.Background()

		// Setup authentication
		accessToken := tu.SetupAdminUser(t, ctx, authCtx)
		defer tu.ClearAuthData(t, ctx, authCtx)

		// Test with non-existent client UUID
		nonExistentClientID := uuid.New()

		req := ch.NewListCasesByClientRequest(t, ctx, testServerURL, nonExistentClientID, 1, 20, accessToken)

		// Make request
		resp, err := httpClient.Do(req)
		require.NoError(t, err)
		defer resp.Body.Close()

		// Should return success but with empty list
		assert.Equal(t, http.StatusOK, resp.StatusCode)

		var response investigation.ListCasesResponse
		err = json.NewDecoder(resp.Body).Decode(&response)
		require.NoError(t, err)

		assert.Empty(t, response.Cases)
		assert.Equal(t, 0, response.Pagination.TotalItems)
	})

	t.Run("UnauthorizedAccess", func(t *testing.T) {
		ctx := context.Background()

		// Clear all test data
		ch.ClearCasesTable(t, ctx, testPool)
		clientHelpers.ClearContactsTable(t, ctx, testPool)
		clientHelpers.ClearClientsTable(t, ctx, testPool)

		// Setup client
		clientID := setupClient(t, ctx)

		// Create request without access token
		req := ch.NewListCasesByClientRequest(t, ctx, testServerURL, clientID, 1, 20, "")

		// Make request
		resp, err := httpClient.Do(req)
		require.NoError(t, err)
		defer resp.Body.Close()

		// Should return unauthorized without authentication
		assert.Equal(t, http.StatusUnauthorized, resp.StatusCode)
	})

	t.Run("InvalidPaginationParameters", func(t *testing.T) {
		ctx := context.Background()

		// Setup authentication
		accessToken := tu.SetupAdminUser(t, ctx, authCtx)
		defer tu.ClearAuthData(t, ctx, authCtx)

		// Clear all test data
		ch.ClearCasesTable(t, ctx, testPool)
		clientHelpers.ClearContactsTable(t, ctx, testPool)
		clientHelpers.ClearClientsTable(t, ctx, testPool)

		// Setup client
		clientID := setupClient(t, ctx)

		// Test invalid page
		req := ch.NewListCasesByClientRequest(t, ctx, testServerURL, clientID, 0, 20, accessToken)

		resp, err := httpClient.Do(req)
		require.NoError(t, err)
		defer resp.Body.Close()

		assert.Equal(t, http.StatusBadRequest, resp.StatusCode)

		// Test invalid pageSize
		req = ch.NewListCasesByClientRequest(t, ctx, testServerURL, clientID, 1, 0, accessToken)

		resp, err = httpClient.Do(req)
		require.NoError(t, err)
		defer resp.Body.Close()

		assert.Equal(t, http.StatusBadRequest, resp.StatusCode)

		// Test pageSize too large
		req = ch.NewListCasesByClientRequest(t, ctx, testServerURL, clientID, 1, 200, accessToken)

		resp, err = httpClient.Do(req)
		require.NoError(t, err)
		defer resp.Body.Close()

		assert.Equal(t, http.StatusBadRequest, resp.StatusCode)
	})
}

// Helper function to find a case by title in the response
func findCaseByTitle(cases []*investigation.CaseResponse, title string) *investigation.CaseResponse {
	for _, c := range cases {
		if c.Title == title {
			return c
		}
	}
	return nil
}
