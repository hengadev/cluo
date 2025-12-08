package caseRepository_test

import (
	"context"
	"testing"
	"time"

	caseDomain "github.com/hengadev/cluo_api/internal/domain/case"
	caseHelpers "github.com/hengadev/cluo_api/test/helpers/case"
	clientHelpers "github.com/hengadev/cluo_api/test/helpers/client"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// make test-func TEST_NAME=TestList TEST_PATH=internal/infrastructure/postgres/case/list_test.go

func TestList(t *testing.T) {
	ctx := context.Background()

	setupClient := func(t *testing.T, ctx context.Context) uuid.UUID {
		clientEncx := clientHelpers.NewTestClientEncx(t)
		err := clientHelpers.InsertClientEncx(t, ctx, testPool, clientEncx)
		require.NoError(t, err)
		return clientEncx.ID
	}
	setupContact := func(t *testing.T, ctx context.Context, clientID uuid.UUID) uuid.UUID {
		contactEncx := clientHelpers.NewTestContactEncx(t)
		contactEncx.ClientID = clientID
		err := clientHelpers.InsertContactEncx(t, ctx, testPool, contactEncx)
		require.NoError(t, err)
		return contactEncx.ID
	}

	t.Run("EmptyList", func(t *testing.T) {
		// Clean up before each test
		caseHelpers.ClearCasesTable(t, ctx, testPool)
		clientHelpers.ClearContactsTable(t, ctx, testPool)
		clientHelpers.ClearClientsTable(t, ctx, testPool)

		pagination := caseDomain.NewPagination()

		cases, total, err := repo.List(ctx, caseDomain.CaseFilter{}, pagination)

		require.NoError(t, err)
		assert.Equal(t, 0, total)
		assert.Empty(t, cases)
	})

	t.Run("ListWithCases", func(t *testing.T) {
		// Clean up before each test
		caseHelpers.ClearCasesTable(t, ctx, testPool)
		clientHelpers.ClearContactsTable(t, ctx, testPool)
		clientHelpers.ClearClientsTable(t, ctx, testPool)

		// Create test cases
		clientID1 := setupClient(t, ctx)
		clientID2 := setupClient(t, ctx)
		contactID1 := setupContact(t, ctx, clientID1)
		contactID2 := setupContact(t, ctx, clientID2)

		caseEncx1 := caseHelpers.NewTestCaseEncx(t)
		caseEncx1.ClientID = clientID1
		caseEncx1.AssignedContactID = &contactID1

		caseEncx2 := caseHelpers.NewTestCaseEncx(t)
		caseEncx2.ClientID = clientID1
		caseEncx2.AssignedContactID = &contactID2

		caseEncx3 := caseHelpers.NewTestCaseEncx(t)
		caseEncx3.ClientID = clientID2
		caseEncx3.AssignedContactID = nil

		err := caseHelpers.InsertCaseEncx(t, ctx, testPool, caseEncx1)
		require.NoError(t, err)
		err = caseHelpers.InsertCaseEncx(t, ctx, testPool, caseEncx2)
		require.NoError(t, err)
		err = caseHelpers.InsertCaseEncx(t, ctx, testPool, caseEncx3)
		require.NoError(t, err)

		// Test with default pagination
		pagination := caseDomain.NewPagination()

		cases, total, err := repo.List(ctx, caseDomain.CaseFilter{}, pagination)

		require.NoError(t, err)
		assert.Equal(t, 3, total)
		assert.Len(t, cases, 3)

		// Should be ordered by created_at DESC
		assert.Equal(t, clientID1, cases[2].ClientID) // Most recent
		assert.Equal(t, clientID1, cases[1].ClientID)
		assert.Equal(t, clientID2, cases[0].ClientID) // Oldest
	})

	t.Run("FilterByClient", func(t *testing.T) {
		// Clean up before each test
		caseHelpers.ClearCasesTable(t, ctx, testPool)
		clientHelpers.ClearContactsTable(t, ctx, testPool)
		clientHelpers.ClearClientsTable(t, ctx, testPool)

		// Create test cases
		clientID1 := setupClient(t, ctx)
		clientID2 := setupClient(t, ctx)

		// Create test caseEncxs data using helper
		caseEncx1 := caseHelpers.NewTestCaseEncx(t)
		caseEncx1.ClientID = clientID1
		caseEncx1.AssignedContactID = nil

		caseEncx2 := caseHelpers.NewTestCaseEncx(t)
		caseEncx2.ClientID = clientID2
		caseEncx2.AssignedContactID = nil

		caseEncx3 := caseHelpers.NewTestCaseEncx(t)
		caseEncx3.ClientID = clientID1
		caseEncx3.AssignedContactID = nil

		err := caseHelpers.InsertCaseEncx(t, ctx, testPool, caseEncx1)
		require.NoError(t, err)

		err = caseHelpers.InsertCaseEncx(t, ctx, testPool, caseEncx2)
		require.NoError(t, err)

		err = caseHelpers.InsertCaseEncx(t, ctx, testPool, caseEncx3)
		require.NoError(t, err)

		// Filter by clientID1
		filter := caseDomain.CaseFilter{
			ClientID: &clientID1,
		}
		pagination := caseDomain.NewPagination()

		cases, total, err := repo.List(ctx, filter, pagination)

		require.NoError(t, err)
		assert.Equal(t, 2, total)
		assert.Len(t, cases, 2)

		for _, c := range cases {
			require.Equal(t, clientID1, c.ClientID)
		}
	})

	t.Run("FilterByAssignedContact", func(t *testing.T) {
		// Clean up before each test
		caseHelpers.ClearCasesTable(t, ctx, testPool)
		clientHelpers.ClearContactsTable(t, ctx, testPool)
		clientHelpers.ClearClientsTable(t, ctx, testPool)

		// Create test cases
		clientID := setupClient(t, ctx)
		contactID1 := setupContact(t, ctx, clientID)
		contactID2 := setupContact(t, ctx, clientID)

		// Create test caseEncx1 data using helper
		caseEncx1 := caseHelpers.NewTestCaseEncx(t)
		caseEncx1.ClientID = clientID
		caseEncx1.AssignedContactID = &contactID1

		caseEncx2 := caseHelpers.NewTestCaseEncx(t)
		caseEncx2.ClientID = clientID
		caseEncx2.AssignedContactID = nil

		caseEncx3 := caseHelpers.NewTestCaseEncx(t)
		caseEncx3.ClientID = clientID
		caseEncx3.AssignedContactID = &contactID2

		err := caseHelpers.InsertCaseEncx(t, ctx, testPool, caseEncx1)
		require.NoError(t, err)

		err = caseHelpers.InsertCaseEncx(t, ctx, testPool, caseEncx2)
		require.NoError(t, err)

		err = caseHelpers.InsertCaseEncx(t, ctx, testPool, caseEncx3)
		require.NoError(t, err)

		// Filter by assigned contact
		filter := caseDomain.CaseFilter{
			AssignedContactID: &contactID1,
		}
		pagination := caseDomain.NewPagination()

		cases, total, err := repo.List(ctx, filter, pagination)

		require.NoError(t, err)
		assert.Equal(t, 1, total)
		assert.Len(t, cases, 1)
		assert.Equal(t, contactID1, *cases[0].AssignedContactID)
	})

	t.Run("FilterByNilAssignedContact", func(t *testing.T) {
		// Clean up before each test
		caseHelpers.ClearCasesTable(t, ctx, testPool)
		clientHelpers.ClearContactsTable(t, ctx, testPool)
		clientHelpers.ClearClientsTable(t, ctx, testPool)

		// Create test cases
		clientID := setupClient(t, ctx)
		contactID1 := setupContact(t, ctx, clientID)

		// Create test caseEncx1 data using helper
		caseEncx1 := caseHelpers.NewTestCaseEncx(t)
		caseEncx1.ClientID = clientID
		caseEncx1.AssignedContactID = &contactID1

		caseEncx2 := caseHelpers.NewTestCaseEncx(t)
		caseEncx2.ClientID = clientID
		caseEncx2.AssignedContactID = nil

		caseEncx3 := caseHelpers.NewTestCaseEncx(t)
		caseEncx3.ClientID = clientID
		caseEncx3.AssignedContactID = nil

		err := caseHelpers.InsertCaseEncx(t, ctx, testPool, caseEncx1)
		require.NoError(t, err)

		err = caseHelpers.InsertCaseEncx(t, ctx, testPool, caseEncx2)
		require.NoError(t, err)

		err = caseHelpers.InsertCaseEncx(t, ctx, testPool, caseEncx3)
		require.NoError(t, err)

		// Filter by nil assigned contact
		nilUUID := uuid.UUID{}
		filter := caseDomain.CaseFilter{
			AssignedContactID: &nilUUID,
		}
		pagination := caseDomain.NewPagination()

		cases, total, err := repo.List(ctx, filter, pagination)

		require.NoError(t, err)
		assert.Equal(t, 2, total)
		assert.Len(t, cases, 2)

		for _, c := range cases {
			require.Nil(t, c.AssignedContactID)
		}
	})

	t.Run("FilterByDateRange", func(t *testing.T) {
		// Clean up before each test
		caseHelpers.ClearCasesTable(t, ctx, testPool)
		clientHelpers.ClearContactsTable(t, ctx, testPool)
		clientHelpers.ClearClientsTable(t, ctx, testPool)

		// Create test cases with different creation dates
		clientID := setupClient(t, ctx)

		now := time.Now()
		yesterday := now.Add(-24 * time.Hour)
		twoDaysAgo := now.Add(-48 * time.Hour)

		caseEncx1 := caseHelpers.NewTestCaseEncx(t)
		caseEncx1.ClientID = clientID
		caseEncx1.CreatedAt = now
		caseEncx1.AssignedContactID = nil

		caseEncx2 := caseHelpers.NewTestCaseEncx(t)
		caseEncx2.ClientID = clientID
		caseEncx2.CreatedAt = yesterday
		caseEncx2.AssignedContactID = nil

		caseEncx3 := caseHelpers.NewTestCaseEncx(t)
		caseEncx3.ClientID = clientID
		caseEncx3.CreatedAt = twoDaysAgo
		caseEncx3.AssignedContactID = nil

		err := caseHelpers.InsertCaseEncx(t, ctx, testPool, caseEncx1)
		require.NoError(t, err)

		err = caseHelpers.InsertCaseEncx(t, ctx, testPool, caseEncx2)
		require.NoError(t, err)

		err = caseHelpers.InsertCaseEncx(t, ctx, testPool, caseEncx3)
		require.NoError(t, err)

		// Filter by date range
		fromDate := now.Add(-25 * time.Hour) // 25 hours ago to now
		filter := caseDomain.CaseFilter{
			DateCreatedFrom: &fromDate,
		}
		pagination := caseDomain.NewPagination()

		cases, total, err := repo.List(ctx, filter, pagination)

		require.NoError(t, err)
		assert.Equal(t, 2, total) // yesterday and now
		assert.Len(t, cases, 2)

		// Test with upper bound
		toDate := now.Add(-12 * time.Hour) // Up to 12 hours ago
		filter = caseDomain.CaseFilter{
			DateCreatedTo: &toDate,
		}

		cases, total, err = repo.List(ctx, filter, pagination)

		require.NoError(t, err)
		assert.Equal(t, 2, total) // yesterday and two days ago
		assert.Len(t, cases, 2)
	})

	t.Run("Pagination", func(t *testing.T) {
		// Clean up before each test
		caseHelpers.ClearCasesTable(t, ctx, testPool)
		clientHelpers.ClearContactsTable(t, ctx, testPool)
		clientHelpers.ClearClientsTable(t, ctx, testPool)

		// Create multiple test cases
		clientID := setupClient(t, ctx)

		for i := 0; i < 5; i++ {
			// Create test caseEncx data using helper
			caseEncx := caseHelpers.NewTestCaseEncx(t)
			caseEncx.ClientID = clientID
			caseEncx.AssignedContactID = nil
			err := caseHelpers.InsertCaseEncx(t, ctx, testPool, caseEncx)
			require.NoError(t, err)
		}

		// Test first page with 2 items
		pagination := caseDomain.Pagination{
			Page:     1,
			PageSize: 2,
		}

		cases, total, err := repo.List(ctx, caseDomain.CaseFilter{}, pagination)

		require.NoError(t, err)
		assert.Equal(t, 5, total) // Total should be all items
		assert.Len(t, cases, 2)   // Page should have 2 items

		// Test second page
		pagination.Page = 2

		cases, total, err = repo.List(ctx, caseDomain.CaseFilter{}, pagination)

		require.NoError(t, err)
		assert.Equal(t, 5, total)
		assert.Len(t, cases, 2)

		// Test third page (last item)
		pagination.Page = 3

		cases, total, err = repo.List(ctx, caseDomain.CaseFilter{}, pagination)

		require.NoError(t, err)
		assert.Equal(t, 5, total)
		assert.Len(t, cases, 1)
	})

	t.Run("InvalidPagination", func(t *testing.T) {
		// Clean up before each test
		caseHelpers.ClearCasesTable(t, ctx, testPool)
		pagination := caseDomain.Pagination{
			Page:     0, // Invalid
			PageSize: 10,
		}

		cases, total, err := repo.List(ctx, caseDomain.CaseFilter{}, pagination)

		require.Error(t, err)
		require.Contains(t, err.Error(), "invalid pagination")
		require.Empty(t, cases)
		require.Equal(t, 0, total)
	})

	t.Run("CombinedFilters", func(t *testing.T) {
		// Clean up before each test
		caseHelpers.ClearCasesTable(t, ctx, testPool)
		clientHelpers.ClearContactsTable(t, ctx, testPool)
		clientHelpers.ClearClientsTable(t, ctx, testPool)

		// Create test cases with various properties
		clientID1 := setupClient(t, ctx)
		clientID2 := setupClient(t, ctx)
		contactID1 := setupContact(t, ctx, clientID1)

		now := time.Now()
		recent := now.Add(-1 * time.Hour)
		older := now.Add(-25 * time.Hour)

		// Create test caseEncx1 data using helper
		caseEncx1 := caseHelpers.NewTestCaseEncx(t)
		caseEncx1.ClientID = clientID1
		caseEncx1.CreatedAt = recent
		caseEncx1.AssignedContactID = &contactID1

		caseEncx2 := caseHelpers.NewTestCaseEncx(t)
		caseEncx2.ClientID = clientID1
		caseEncx2.CreatedAt = recent
		caseEncx2.AssignedContactID = nil

		caseEncx3 := caseHelpers.NewTestCaseEncx(t)
		caseEncx3.ClientID = clientID2
		caseEncx3.CreatedAt = older
		caseEncx3.AssignedContactID = nil

		// Recent cases for client1
		err := caseHelpers.InsertCaseEncx(t, ctx, testPool, caseEncx1)
		require.NoError(t, err)

		err = caseHelpers.InsertCaseEncx(t, ctx, testPool, caseEncx2)
		require.NoError(t, err)

		// Older cases for client2
		err = caseHelpers.InsertCaseEncx(t, ctx, testPool, caseEncx3)
		require.NoError(t, err)

		// Filter: client1 + recent date range
		fromDate := now.Add(-2 * time.Hour)
		filter := caseDomain.CaseFilter{
			ClientID:        &clientID1,
			DateCreatedFrom: &fromDate,
		}
		pagination := caseDomain.NewPagination()

		cases, total, err := repo.List(ctx, filter, pagination)

		require.NoError(t, err)
		assert.Equal(t, 2, total)
		assert.Len(t, cases, 2)

		for _, c := range cases {
			assert.Equal(t, clientID1, c.ClientID)
			assert.True(t, c.CreatedAt.After(fromDate))
		}
	})
}
