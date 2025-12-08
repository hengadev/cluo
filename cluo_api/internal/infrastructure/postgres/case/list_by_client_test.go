package caseRepository_test

import (
	"context"
	"testing"
	"time"

	"github.com/google/uuid"
	caseDomain "github.com/hengadev/cluo_api/internal/domain/case"
	ch "github.com/hengadev/cluo_api/test/helpers/case"
	"github.com/stretchr/testify/require"
)

// make test-func TEST_NAME=TestListByClient TEST_PATH=internal/infrastructure/postgres/case/list_by_client_test.go

func TestListByClient(t *testing.T) {
	ctx := context.Background()

	t.Run("EmptyList", func(t *testing.T) {
		// Clean up before each test
		ch.ClearCasesTable(t, ctx, testPool)
		clientID := uuid.New()
		pagination := caseDomain.NewPagination()

		cases, total, err := repo.ListByClient(ctx, clientID, pagination)

		require.NoError(t, err)
		require.Equal(t, 0, total)
		require.Empty(t, cases)
	})

	t.Run("ListWithCases", func(t *testing.T) {
		// Clean up before each test
		ch.ClearCasesTable(t, ctx, testPool)
		clientID := uuid.New()
		contactID := uuid.New()
		now := time.Now()

		// Insert test cases for this client
		insertTestCase(t, ctx, clientID, &contactID, now.Add(-1*time.Hour))
		insertTestCase(t, ctx, clientID, nil, now.Add(-2*time.Hour))
		insertTestCase(t, ctx, clientID, &contactID, now.Add(-3*time.Hour))

		// Insert cases for different client (should not be included)
		otherClientID := uuid.New()
		insertTestCase(t, ctx, otherClientID, nil, now.Add(-1*time.Hour))

		pagination := caseDomain.NewPagination()

		cases, total, err := repo.ListByClient(ctx, clientID, pagination)

		require.NoError(t, err)
		require.Equal(t, 3, total) // Only cases for the specified client
		require.Len(t, cases, 3)

		// All cases should belong to the specified client
		for _, c := range cases {
			require.Equal(t, clientID, c.ClientID)
		}

		// Should be ordered by created_at DESC
		require.True(t, cases[0].CreatedAt.After(cases[1].CreatedAt))
		require.True(t, cases[1].CreatedAt.After(cases[2].CreatedAt))
	})

	t.Run("Pagination", func(t *testing.T) {
		// Clean up before each test
		ch.ClearCasesTable(t, ctx, testPool)
		clientID := uuid.New()
		now := time.Now()

		// Create 5 test cases
		for i := 0; i < 5; i++ {
			insertTestCase(t, ctx, clientID, nil, now.Add(time.Duration(i)*time.Hour))
		}

		// Test first page with 2 items
		pagination := caseDomain.Pagination{
			Page:     1,
			PageSize: 2,
		}

		cases, total, err := repo.ListByClient(ctx, clientID, pagination)

		require.NoError(t, err)
		require.Equal(t, 5, total) // Total should be all items for this client
		require.Len(t, cases, 2)   // Page should have 2 items

		// Test second page
		pagination.Page = 2

		cases, total, err = repo.ListByClient(ctx, clientID, pagination)

		require.NoError(t, err)
		require.Equal(t, 5, total)
		require.Len(t, cases, 2)

		// Test third page (last item)
		pagination.Page = 3

		cases, total, err = repo.ListByClient(ctx, clientID, pagination)

		require.NoError(t, err)
		require.Equal(t, 5, total)
		require.Len(t, cases, 1)

		// Test page beyond results
		pagination.Page = 4

		cases, total, err = repo.ListByClient(ctx, clientID, pagination)

		require.NoError(t, err)
		require.Equal(t, 5, total)
		require.Empty(t, cases)
	})

	t.Run("InvalidPagination", func(t *testing.T) {
		// Clean up before each test
		ch.ClearCasesTable(t, ctx, testPool)
		clientID := uuid.New()
		pagination := caseDomain.Pagination{
			Page:     0, // Invalid
			PageSize: 10,
		}

		cases, total, err := repo.ListByClient(ctx, clientID, pagination)

		require.Error(t, err)
		require.Contains(t, err.Error(), "invalid pagination")
		require.Empty(t, cases)
		require.Equal(t, 0, total)
	})

	t.Run("NilClientID", func(t *testing.T) {
		// Clean up before each test
		ch.ClearCasesTable(t, ctx, testPool)
		pagination := caseDomain.NewPagination()

		cases, total, err := repo.ListByClient(ctx, uuid.Nil, pagination)

		require.Error(t, err)
		require.Contains(t, err.Error(), "client ID cannot be nil")
		require.Empty(t, cases)
		require.Equal(t, 0, total)
	})

	t.Run("MaxPageSize", func(t *testing.T) {
		// Clean up before each test
		ch.ClearCasesTable(t, ctx, testPool)
		clientID := uuid.New()
		now := time.Now()

		// Create multiple test cases
		for i := 0; i < 25; i++ {
			insertTestCase(t, ctx, clientID, nil, now.Add(time.Duration(i)*time.Hour))
		}

		// Test with maximum page size
		pagination := caseDomain.Pagination{
			Page:     1,
			PageSize: 100, // Max allowed
		}

		cases, total, err := repo.ListByClient(ctx, clientID, pagination)

		require.NoError(t, err)
		require.Equal(t, 25, total)
		require.Len(t, cases, 25)
	})

	t.Run("AssignedContactHandling", func(t *testing.T) {
		// Clean up before each test
		ch.ClearCasesTable(t, ctx, testPool)
		clientID := uuid.New()
		contactID1 := uuid.New()
		contactID2 := uuid.New()
		now := time.Now()

		// Create cases with different assigned contacts
		insertTestCase(t, ctx, clientID, &contactID1, now.Add(-1*time.Hour))
		insertTestCase(t, ctx, clientID, nil, now.Add(-2*time.Hour))
		insertTestCase(t, ctx, clientID, &contactID2, now.Add(-3*time.Hour))

		pagination := caseDomain.NewPagination()

		cases, total, err := repo.ListByClient(ctx, clientID, pagination)

		require.NoError(t, err)
		require.Equal(t, 3, total)
		require.Len(t, cases, 3)

		// Verify assigned contacts are preserved
		require.Equal(t, contactID1, *cases[0].AssignedContactID)
		require.Nil(t, cases[1].AssignedContactID)
		require.Equal(t, contactID2, *cases[2].AssignedContactID)
	})

	t.Run("MultipleClientsIsolation", func(t *testing.T) {
		// Clean up before each test
		ch.ClearCasesTable(t, ctx, testPool)
		// Create multiple clients with cases
		clientID1 := uuid.New()
		clientID2 := uuid.New()
		clientID3 := uuid.New()

		now := time.Now()

		// Add cases for each client
		for i := 0; i < 3; i++ {
			insertTestCase(t, ctx, clientID1, nil, now.Add(time.Duration(i)*time.Hour))
			insertTestCase(t, ctx, clientID2, nil, now.Add(time.Duration(i)*time.Hour))
			insertTestCase(t, ctx, clientID3, nil, now.Add(time.Duration(i)*time.Hour))
		}

		pagination := caseDomain.NewPagination()

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
