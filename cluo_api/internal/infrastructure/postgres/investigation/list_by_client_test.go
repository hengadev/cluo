package investigationRepository_test

import (
	"context"
	"testing"
	// "time"

	"github.com/google/uuid"
	"github.com/hengadev/cluo_api/internal/domain/investigation"
	caseHelpers "github.com/hengadev/cluo_api/test/helpers/investigation"
	clientHelpers "github.com/hengadev/cluo_api/test/helpers/client"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// make test-func TEST_NAME=TestListByClient TEST_PATH=internal/infrastructure/postgres/case/list_by_client_test.go

func TestListByClient(t *testing.T) {
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
		clientID := uuid.New()
		pagination := investigation.NewPagination()

		cases, total, err := repo.ListByClient(ctx, clientID, pagination)

		require.NoError(t, err)
		assert.Equal(t, 0, total)
		assert.Empty(t, cases)
	})

	t.Run("ListWithCases", func(t *testing.T) {
		// Clean up before each test
		caseHelpers.ClearCasesTable(t, ctx, testPool)
		clientHelpers.ClearClientsTable(t, ctx, testPool)
		clientHelpers.ClearContactsTable(t, ctx, testPool)

		// Setup client and contact in database
		clientID1 := setupClient(t, ctx)
		clientID2 := setupClient(t, ctx)
		contactID := setupContact(t, ctx, clientID1)

		caseEncx1 := caseHelpers.NewTestCaseEncx(t)
		caseEncx1.ClientID = clientID1
		caseEncx1.AssignedContactID = &contactID

		caseEncx2 := caseHelpers.NewTestCaseEncx(t)
		caseEncx2.ClientID = clientID1
		caseEncx2.AssignedContactID = nil

		caseEncx3 := caseHelpers.NewTestCaseEncx(t)
		caseEncx3.ClientID = clientID1
		caseEncx3.AssignedContactID = &contactID

		// Insert test cases for this client
		err := caseHelpers.InsertCaseEncx(t, ctx, testPool, caseEncx1)
		require.NoError(t, err)

		err = caseHelpers.InsertCaseEncx(t, ctx, testPool, caseEncx2)
		require.NoError(t, err)

		err = caseHelpers.InsertCaseEncx(t, ctx, testPool, caseEncx3)
		require.NoError(t, err)

		// Insert cases for different client (should not be included)

		caseEncx4 := caseHelpers.NewTestCaseEncx(t)
		caseEncx4.ClientID = clientID2
		caseEncx4.AssignedContactID = nil

		err = caseHelpers.InsertCaseEncx(t, ctx, testPool, caseEncx4)
		require.NoError(t, err)

		pagination := investigation.NewPagination()

		cases, total, err := repo.ListByClient(ctx, clientID1, pagination)

		require.NoError(t, err)
		assert.Equal(t, 3, total) // Only cases for the specified client
		assert.Len(t, cases, 3)

		// All cases should belong to the specified client
		for _, c := range cases {
			assert.Equal(t, clientID1, c.ClientID)
		}

		// Should be ordered by created_at DESC
		assert.True(t, cases[0].CreatedAt.After(cases[1].CreatedAt))
		assert.True(t, cases[1].CreatedAt.After(cases[2].CreatedAt))
	})

	t.Run("Pagination", func(t *testing.T) {
		// Clean up before each test
		caseHelpers.ClearCasesTable(t, ctx, testPool)
		clientHelpers.ClearClientsTable(t, ctx, testPool)
		clientHelpers.ClearContactsTable(t, ctx, testPool)

		// Setup client and contact in database
		clientID := setupClient(t, ctx)

		// Create 5 test cases
		for i := 0; i < 5; i++ {
			caseEncx := caseHelpers.NewTestCaseEncx(t)
			caseEncx.ClientID = clientID
			caseEncx.AssignedContactID = nil

			err := caseHelpers.InsertCaseEncx(t, ctx, testPool, caseEncx)
			require.NoError(t, err)
		}

		// Test first page with 2 items
		pagination := investigation.Pagination{
			Page:     1,
			PageSize: 2,
		}

		cases, total, err := repo.ListByClient(ctx, clientID, pagination)

		require.NoError(t, err)
		assert.Equal(t, 5, total) // Total should be all items for this client
		assert.Len(t, cases, 2)   // Page should have 2 items

		// Test second page
		pagination.Page = 2

		cases, total, err = repo.ListByClient(ctx, clientID, pagination)

		require.NoError(t, err)
		assert.Equal(t, 5, total)
		assert.Len(t, cases, 2)

		// Test third page (last item)
		pagination.Page = 3

		cases, total, err = repo.ListByClient(ctx, clientID, pagination)

		require.NoError(t, err)
		assert.Equal(t, 5, total)
		assert.Len(t, cases, 1)

		// Test page beyond results
		pagination.Page = 4

		cases, total, err = repo.ListByClient(ctx, clientID, pagination)

		require.NoError(t, err)
		assert.Equal(t, 5, total)
		assert.Empty(t, cases)
	})

	t.Run("InvalidPagination", func(t *testing.T) {
		// Clean up before each test
		caseHelpers.ClearCasesTable(t, ctx, testPool)
		clientID := uuid.New()
		pagination := investigation.Pagination{
			Page:     0, // Invalid
			PageSize: 10,
		}

		cases, total, err := repo.ListByClient(ctx, clientID, pagination)

		assert.Error(t, err)
		assert.Contains(t, err.Error(), "invalid pagination")
		assert.Empty(t, cases)
		assert.Equal(t, 0, total)
	})

	t.Run("NilClientID", func(t *testing.T) {
		// Clean up before each test
		caseHelpers.ClearCasesTable(t, ctx, testPool)
		pagination := investigation.NewPagination()

		cases, total, err := repo.ListByClient(ctx, uuid.Nil, pagination)

		require.Error(t, err)
		assert.Contains(t, err.Error(), "client ID cannot be nil")
		assert.Empty(t, cases)
		assert.Equal(t, 0, total)
	})

	t.Run("MaxPageSize", func(t *testing.T) {
		// Clean up before each test
		caseHelpers.ClearCasesTable(t, ctx, testPool)
		clientHelpers.ClearClientsTable(t, ctx, testPool)
		clientHelpers.ClearContactsTable(t, ctx, testPool)

		// Setup client and contact in database
		clientID := setupClient(t, ctx)

		// Create multiple test cases
		for i := 0; i < 25; i++ {
			caseEncx := caseHelpers.NewTestCaseEncx(t)
			caseEncx.ClientID = clientID
			caseEncx.AssignedContactID = nil
			err := caseHelpers.InsertCaseEncx(t, ctx, testPool, caseEncx)
			require.NoError(t, err)
		}

		// Test with maximum page size
		pagination := investigation.Pagination{
			Page:     1,
			PageSize: 100, // Max allowed
		}

		cases, total, err := repo.ListByClient(ctx, clientID, pagination)

		require.NoError(t, err)
		assert.Equal(t, 25, total)
		assert.Len(t, cases, 25)
	})

	t.Run("AssignedContactHandling", func(t *testing.T) {
		// Clean up before each test
		caseHelpers.ClearCasesTable(t, ctx, testPool)
		clientHelpers.ClearContactsTable(t, ctx, testPool)
		clientHelpers.ClearClientsTable(t, ctx, testPool)

		clientID := setupClient(t, ctx)
		contactID1 := setupContact(t, ctx, clientID)
		contactID2 := setupContact(t, ctx, clientID)

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

		pagination := investigation.NewPagination()

		cases, total, err := repo.ListByClient(ctx, clientID, pagination)

		require.NoError(t, err)
		assert.Equal(t, 3, total)
		assert.Len(t, cases, 3)

		// Verify assigned contacts are preserved
		assert.Equal(t, contactID1, *cases[2].AssignedContactID)
		assert.Nil(t, cases[1].AssignedContactID)
		assert.Equal(t, contactID2, *cases[0].AssignedContactID)
	})

	t.Run("MultipleClientsIsolation", func(t *testing.T) {
		// Clean up before each test
		caseHelpers.ClearCasesTable(t, ctx, testPool)
		clientHelpers.ClearContactsTable(t, ctx, testPool)
		clientHelpers.ClearClientsTable(t, ctx, testPool)

		// Create multiple clients with cases
		clientID1 := setupClient(t, ctx)
		clientID2 := setupClient(t, ctx)
		clientID3 := setupClient(t, ctx)

		// Add cases for each client
		for i := 0; i < 3; i++ {
			caseEncx1 := caseHelpers.NewTestCaseEncx(t)
			caseEncx1.ClientID = clientID1
			caseEncx1.AssignedContactID = nil

			caseEncx2 := caseHelpers.NewTestCaseEncx(t)
			caseEncx2.ClientID = clientID2
			caseEncx2.AssignedContactID = nil

			caseEncx3 := caseHelpers.NewTestCaseEncx(t)
			caseEncx3.ClientID = clientID3
			caseEncx3.AssignedContactID = nil

			err := caseHelpers.InsertCaseEncx(t, ctx, testPool, caseEncx1)
			require.NoError(t, err)

			err = caseHelpers.InsertCaseEncx(t, ctx, testPool, caseEncx2)
			require.NoError(t, err)

			err = caseHelpers.InsertCaseEncx(t, ctx, testPool, caseEncx3)
			require.NoError(t, err)
		}

		pagination := investigation.NewPagination()

		// Test each client separately
		for _, clientID := range []uuid.UUID{clientID1, clientID2, clientID3} {
			cases, total, err := repo.ListByClient(ctx, clientID, pagination)

			require.NoError(t, err)
			require.Equal(t, 3, total)
			require.Len(t, cases, 3)

			// All cases should belong to the specified client
			for _, c := range cases {
				require.Equal(t, clientID, c.ClientID)
			}
		}
	})
}
