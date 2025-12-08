package case_test

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"testing"
	"time"

	tu "github.com/hengadev/cluo_api/internal/common/testutils"
	caseDomain "github.com/hengadev/cluo_api/internal/domain/case"
	clientDomain "github.com/hengadev/cluo_api/internal/domain/client"
	ch "github.com/hengadev/cluo_api/test/helpers/case"
	clientHelpers "github.com/hengadev/cluo_api/test/helpers/client"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// make test-func TEST_NAME=TestListCases TEST_PATH=test/integration/case/list_cases_test.go

// TestListCases tests all scenarios for listing cases with filtering and pagination
func TestListCases(t *testing.T) {
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

	t.Run("EmptyList", func(t *testing.T) {
		ctx := context.Background()

		// Setup authentication
		accessToken := tu.SetupAdminUser(t, ctx, authCtx)
		defer tu.ClearAuthData(t, ctx, authCtx)

		// Clear all test data
		ch.ClearCasesTable(t, ctx, testPool)
		clientHelpers.ClearContactsTable(t, ctx, testPool)
		clientHelpers.ClearClientsTable(t, ctx, testPool)

		// Create request
		req := ch.NewListCasesRequest(t, ctx, testServerURL, 1, 20, nil, accessToken)

		// Make request
		resp, err := httpClient.Do(req)
		require.NoError(t, err)
		defer resp.Body.Close()

		// Verify response
		assert.Equal(t, http.StatusOK, resp.StatusCode)

		var response caseDomain.ListCasesResponse
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
		contactID2 := setupContact(t, ctx, clientID2)

		now := time.Now()

		// Create cases with different timestamps to test ordering
		case3Data := ch.NewTestCase(t)
		case3Data.ClientID = clientID1
		case3Data.AssignedContactID = &contactID1
		case3Data.Title = "Case 3"
		case3Data.CreatedAt = now.Add(-1 * time.Hour)
		case3Encx, err := caseDomain.ProcessCaseEncx(ctx, crypto, case3Data)
		require.NoError(t, err)
		err = ch.InsertCaseEncx(t, ctx, testPool, case3Encx)
		require.NoError(t, err)

		case2Data := ch.NewTestCase(t)
		case2Data.ClientID = clientID1
		case2Data.AssignedContactID = nil
		case2Data.Title = "Case 2"
		case2Data.CreatedAt = now.Add(-2 * time.Hour)
		case2Encx, err := caseDomain.ProcessCaseEncx(ctx, crypto, case2Data)
		require.NoError(t, err)
		err = ch.InsertCaseEncx(t, ctx, testPool, case2Encx)
		require.NoError(t, err)

		case1Data := ch.NewTestCase(t)
		case1Data.ClientID = clientID2
		case1Data.AssignedContactID = &contactID2
		case1Data.Title = "Case 1"
		case1Data.CreatedAt = now.Add(-3 * time.Hour)
		case1Encx, err := caseDomain.ProcessCaseEncx(ctx, crypto, case1Data)
		require.NoError(t, err)
		err = ch.InsertCaseEncx(t, ctx, testPool, case1Encx)
		require.NoError(t, err)

		// Create request
		req := ch.NewListCasesRequest(t, ctx, testServerURL, 1, 20, nil, accessToken)

		// Make request
		resp, err := httpClient.Do(req)
		require.NoError(t, err)
		defer resp.Body.Close()

		// Verify response
		assert.Equal(t, http.StatusOK, resp.StatusCode)

		var response caseDomain.ListCasesResponse
		err = json.NewDecoder(resp.Body).Decode(&response)
		require.NoError(t, err)

		assert.Len(t, response.Cases, 3)
		assert.Equal(t, 1, response.Pagination.Page)
		assert.Equal(t, 20, response.Pagination.PageSize)
		assert.Equal(t, 3, response.Pagination.TotalItems)
		assert.Equal(t, 1, response.Pagination.TotalPages)

		// Verify ordering (should be by created_at DESC)
		assert.Equal(t, case3Data.ID.String(), response.Cases[0].ID)
		assert.Equal(t, case2Data.ID.String(), response.Cases[1].ID)
		assert.Equal(t, case1Data.ID.String(), response.Cases[2].ID)
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

		// Create 5 cases
		for i := 0; i < 5; i++ {
			caseData := ch.NewTestCase(t)
			caseData.ClientID = clientID
			caseData.AssignedContactID = nil
			caseData.Title = fmt.Sprintf("Test Case %d", i+1)
			caseEncx, err := caseDomain.ProcessCaseEncx(ctx, crypto, caseData)
			require.NoError(t, err)
			err = ch.InsertCaseEncx(t, ctx, testPool, caseEncx)
			require.NoError(t, err)
		}

		// Test first page with 2 items
		req := ch.NewListCasesRequest(t, ctx, testServerURL, 1, 2, nil, accessToken)

		resp, err := httpClient.Do(req)
		require.NoError(t, err)
		defer resp.Body.Close()

		assert.Equal(t, http.StatusOK, resp.StatusCode)

		var response caseDomain.ListCasesResponse
		err = json.NewDecoder(resp.Body).Decode(&response)
		require.NoError(t, err)

		assert.Len(t, response.Cases, 2)
		assert.Equal(t, 1, response.Pagination.Page)
		assert.Equal(t, 2, response.Pagination.PageSize)
		assert.Equal(t, 5, response.Pagination.TotalItems)
		assert.Equal(t, 3, response.Pagination.TotalPages)

		// Test second page
		req = ch.NewListCasesRequest(t, ctx, testServerURL, 2, 2, nil, accessToken)

		resp, err = httpClient.Do(req)
		require.NoError(t, err)
		defer resp.Body.Close()

		assert.Equal(t, http.StatusOK, resp.StatusCode)

		err = json.NewDecoder(resp.Body).Decode(&response)
		require.NoError(t, err)

		assert.Len(t, response.Cases, 2)
		assert.Equal(t, 2, response.Pagination.Page)
		assert.Equal(t, 2, response.Pagination.PageSize)
		assert.Equal(t, 5, response.Pagination.TotalItems)
		assert.Equal(t, 3, response.Pagination.TotalPages)

		// Test third page
		req = ch.NewListCasesRequest(t, ctx, testServerURL, 3, 2, nil, accessToken)

		resp, err = httpClient.Do(req)
		require.NoError(t, err)
		defer resp.Body.Close()

		assert.Equal(t, http.StatusOK, resp.StatusCode)

		err = json.NewDecoder(resp.Body).Decode(&response)
		require.NoError(t, err)

		assert.Len(t, response.Cases, 1)
		assert.Equal(t, 3, response.Pagination.Page)
		assert.Equal(t, 2, response.Pagination.PageSize)
		assert.Equal(t, 5, response.Pagination.TotalItems)
		assert.Equal(t, 3, response.Pagination.TotalPages)
	})

	t.Run("FilterByClientId", func(t *testing.T) {
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

		// Create cases for each client
		for i := 0; i < 3; i++ {
			case1 := ch.NewTestCase(t)
			case1.ClientID = clientID1
			case1.AssignedContactID = nil
			case1Encx, err := caseDomain.ProcessCaseEncx(ctx, crypto, case1)
			require.NoError(t, err)
			err = ch.InsertCaseEncx(t, ctx, testPool, case1Encx)
			require.NoError(t, err)

			case2 := ch.NewTestCase(t)
			case2.ClientID = clientID2
			case2.AssignedContactID = nil
			case2Encx, err := caseDomain.ProcessCaseEncx(ctx, crypto, case2)
			require.NoError(t, err)
			err = ch.InsertCaseEncx(t, ctx, testPool, case2Encx)
			require.NoError(t, err)
		}

		// Filter by clientID1
		filters := map[string]string{
			"clientId": clientID1.String(),
		}
		req := ch.NewListCasesRequest(t, ctx, testServerURL, 1, 20, filters, accessToken)

		resp, err := httpClient.Do(req)
		require.NoError(t, err)
		defer resp.Body.Close()

		assert.Equal(t, http.StatusOK, resp.StatusCode)

		var response caseDomain.ListCasesResponse
		err = json.NewDecoder(resp.Body).Decode(&response)
		require.NoError(t, err)

		assert.Len(t, response.Cases, 3)

		// Verify all cases belong to clientID1
		for _, caseResp := range response.Cases {
			assert.Equal(t, clientID1.String(), caseResp.ClientID)
		}
	})

	t.Run("FilterByStatus", func(t *testing.T) {
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

		// Create cases with different statuses
		draftCase := ch.NewTestCase(t)
		draftCase.ClientID = clientID
		draftCase.AssignedContactID = nil
		draftCase.Status = caseDomain.CaseStatusDraft

		progressCase := ch.NewTestCase(t)
		progressCase.ClientID = clientID
		progressCase.AssignedContactID = nil
		progressCase.Status = caseDomain.CaseStatusInProgress

		readyCase := ch.NewTestCase(t)
		readyCase.ClientID = clientID
		readyCase.AssignedContactID = nil
		readyCase.Status = caseDomain.CaseStatusReady

		draftCaseEncx, _ := caseDomain.ProcessCaseEncx(ctx, crypto, draftCase)
		progressCaseEncx, _ := caseDomain.ProcessCaseEncx(ctx, crypto, progressCase)
		readyCaseEncx, _ := caseDomain.ProcessCaseEncx(ctx, crypto, readyCase)

		ch.InsertCaseEncx(t, ctx, testPool, draftCaseEncx)
		ch.InsertCaseEncx(t, ctx, testPool, progressCaseEncx)
		ch.InsertCaseEncx(t, ctx, testPool, readyCaseEncx)

		// Filter by draft status
		filters := map[string]string{
			"status": "draft",
		}
		req := ch.NewListCasesRequest(t, ctx, testServerURL, 1, 20, filters, accessToken)

		resp, err := httpClient.Do(req)
		require.NoError(t, err)
		defer resp.Body.Close()

		assert.Equal(t, http.StatusOK, resp.StatusCode)

		var response caseDomain.ListCasesResponse
		err = json.NewDecoder(resp.Body).Decode(&response)
		require.NoError(t, err)

		assert.Len(t, response.Cases, 1)
		assert.Equal(t, "draft", response.Cases[0].Status)
	})

	t.Run("SearchFunctionality", func(t *testing.T) {
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

		// Create cases with searchable content
		importantCase := ch.NewTestCase(t)
		importantCase.ClientID = clientID
		importantCase.AssignedContactID = nil
		importantCase.Title = "This is an important case"
		importantCase.Description = "Contains sensitive information about legal matters"

		normalCase := ch.NewTestCase(t)
		normalCase.ClientID = clientID
		normalCase.AssignedContactID = nil
		normalCase.Title = "Regular case"
		normalCase.Description = "Normal business operations"

		urgentCase := ch.NewTestCase(t)
		urgentCase.ClientID = clientID
		urgentCase.AssignedContactID = nil
		urgentCase.Title = "Urgent legal matter requires immediate attention"
		urgentCase.Description = "Time-sensitive information"

		importantCaseEncx, _ := caseDomain.ProcessCaseEncx(ctx, crypto, importantCase)
		normalCaseEncx, _ := caseDomain.ProcessCaseEncx(ctx, crypto, normalCase)
		urgentCaseEncx, _ := caseDomain.ProcessCaseEncx(ctx, crypto, urgentCase)

		ch.InsertCaseEncx(t, ctx, testPool, importantCaseEncx)
		ch.InsertCaseEncx(t, ctx, testPool, normalCaseEncx)
		ch.InsertCaseEncx(t, ctx, testPool, urgentCaseEncx)

		// Search for "urgent" in title/description
		filters := map[string]string{
			"search": "urgent",
		}
		req := ch.NewListCasesRequest(t, ctx, testServerURL, 1, 20, filters, accessToken)

		resp, err := httpClient.Do(req)
		require.NoError(t, err)
		defer resp.Body.Close()

		assert.Equal(t, http.StatusOK, resp.StatusCode)

		var response caseDomain.ListCasesResponse
		err = json.NewDecoder(resp.Body).Decode(&response)
		require.NoError(t, err)

		assert.Len(t, response.Cases, 1)
		assert.Contains(t, response.Cases[0].Title, "Urgent")
		assert.Contains(t, response.Cases[0].Description, "Time-sensitive")

		// Search for "important"
		filters = map[string]string{
			"search": "important",
		}
		req = ch.NewListCasesRequest(t, ctx, testServerURL, 1, 20, filters, accessToken)

		resp, err = httpClient.Do(req)
		require.NoError(t, err)
		defer resp.Body.Close()

		assert.Equal(t, http.StatusOK, resp.StatusCode)

		err = json.NewDecoder(resp.Body).Decode(&response)
		require.NoError(t, err)

		assert.Len(t, response.Cases, 1)
		assert.Contains(t, response.Cases[0].Title, "important")
	})

	t.Run("DateRangeFiltering", func(t *testing.T) {
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

		now := time.Now()
		oldDate := now.Add(-7 * 24 * time.Hour)
		recentDate := now.Add(-2 * 24 * time.Hour)

		// Create cases with different creation dates
		oldCase := ch.NewTestCase(t)
		oldCase.ClientID = clientID
		oldCase.AssignedContactID = nil
		oldCase.CreatedAt = oldDate
		oldCaseEncx, _ := caseDomain.ProcessCaseEncx(ctx, crypto, oldCase)
		ch.InsertCaseEncx(t, ctx, testPool, oldCaseEncx)

		recentCase := ch.NewTestCase(t)
		recentCase.ClientID = clientID
		recentCase.AssignedContactID = nil
		recentCase.CreatedAt = recentDate
		recentCaseEncx, _ := caseDomain.ProcessCaseEncx(ctx, crypto, recentCase)
		ch.InsertCaseEncx(t, ctx, testPool, recentCaseEncx)

		// Filter by creation date (last 5 days)
		fromDate := now.Add(-5 * 24 * time.Hour).Format(time.RFC3339)
		filters := map[string]string{
			"dateCreatedFrom": fromDate,
		}
		req := ch.NewListCasesRequest(t, ctx, testServerURL, 1, 20, filters, accessToken)

		resp, err := httpClient.Do(req)
		require.NoError(t, err)
		defer resp.Body.Close()

		assert.Equal(t, http.StatusOK, resp.StatusCode)

		var response caseDomain.ListCasesResponse
		err = json.NewDecoder(resp.Body).Decode(&response)
		require.NoError(t, err)

		assert.Len(t, response.Cases, 1) // Only the recent case
		assert.Equal(t, clientID.String(), response.Cases[0].ClientID)
	})

	t.Run("InvalidPagination", func(t *testing.T) {
		ctx := context.Background()

		// Setup authentication
		accessToken := tu.SetupAdminUser(t, ctx, authCtx)
		defer tu.ClearAuthData(t, ctx, authCtx)

		// Test invalid page
		req := ch.NewListCasesRequest(t, ctx, testServerURL, 0, 20, nil, accessToken)

		resp, err := httpClient.Do(req)
		require.NoError(t, err)
		defer resp.Body.Close()

		assert.Equal(t, http.StatusBadRequest, resp.StatusCode)

		// Test invalid pageSize
		req = ch.NewListCasesRequest(t, ctx, testServerURL, 1, 0, nil, accessToken)

		resp, err = httpClient.Do(req)
		require.NoError(t, err)
		defer resp.Body.Close()

		assert.Equal(t, http.StatusBadRequest, resp.StatusCode)

		// Test pageSize too large
		req = ch.NewListCasesRequest(t, ctx, testServerURL, 1, 200, nil, accessToken)

		resp, err = httpClient.Do(req)
		require.NoError(t, err)
		defer resp.Body.Close()

		assert.Equal(t, http.StatusBadRequest, resp.StatusCode)
	})

	t.Run("UnauthorizedAccess", func(t *testing.T) {
		ctx := context.Background()

		// Clear all test data
		ch.ClearCasesTable(t, ctx, testPool)
		clientHelpers.ClearContactsTable(t, ctx, testPool)
		clientHelpers.ClearClientsTable(t, ctx, testPool)

		// Create request without access token
		req := ch.NewListCasesRequest(t, ctx, testServerURL, 1, 20, nil, "")

		resp, err := httpClient.Do(req)
		require.NoError(t, err)
		defer resp.Body.Close()

		// Should return unauthorized without authentication
		assert.Equal(t, http.StatusUnauthorized, resp.StatusCode)
	})
}
