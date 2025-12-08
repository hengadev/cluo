package caseRepository_test

import (
	"context"
	"testing"
	"time"

	"github.com/google/uuid"
	caseDomain "github.com/hengadev/cluo_api/internal/domain/case"
	ch "github.com/hengadev/cluo_api/test/helpers/case"
	"github.com/hengadev/encx"
	"github.com/stretchr/testify/require"
)

// make test-func TEST_NAME=TestList TEST_PATH=internal/infrastructure/postgres/case/list_test.go

func TestList(t *testing.T) {
	ctx := context.Background()

	t.Run("EmptyList", func(t *testing.T) {
		// Clean up before each test
		ch.ClearCasesTable(t, ctx, testPool)
		pagination := caseDomain.NewPagination()

		cases, total, err := repo.List(ctx, caseDomain.CaseFilter{}, pagination)

		require.NoError(t, err)
		require.Equal(t, 0, total)
		require.Empty(t, cases)
	})

	t.Run("ListWithCases", func(t *testing.T) {
		// Clean up before each test
		ch.ClearCasesTable(t, ctx, testPool)

		// Create test cases
		clientID1 := uuid.New()
		clientID2 := uuid.New()
		contactID1 := uuid.New()
		contactID2 := uuid.New()

		now := time.Now()

		// Insert test cases directly into database
		insertTestCase(t, ctx, clientID1, &contactID1, now.Add(-1*time.Hour))
		insertTestCase(t, ctx, clientID1, &contactID2, now.Add(-2*time.Hour))
		insertTestCase(t, ctx, clientID2, nil, now.Add(-3*time.Hour))

		// Test with default pagination
		pagination := caseDomain.NewPagination()

		cases, total, err := repo.List(ctx, caseDomain.CaseFilter{}, pagination)

		require.NoError(t, err)
		require.Equal(t, 3, total)
		require.Len(t, cases, 3)

		// Should be ordered by created_at DESC
		require.Equal(t, clientID1, cases[0].ClientID) // Most recent
		require.Equal(t, clientID1, cases[1].ClientID)
		require.Equal(t, clientID2, cases[2].ClientID) // Oldest
	})

	t.Run("FilterByClient", func(t *testing.T) {
		// Clean up before each test
		ch.ClearCasesTable(t, ctx, testPool)
		// Create test cases
		clientID1 := uuid.New()
		clientID2 := uuid.New()

		now := time.Now()

		insertTestCase(t, ctx, clientID1, nil, now.Add(-1*time.Hour))
		insertTestCase(t, ctx, clientID2, nil, now.Add(-2*time.Hour))
		insertTestCase(t, ctx, clientID1, nil, now.Add(-3*time.Hour))

		// Filter by clientID1
		filter := caseDomain.CaseFilter{
			ClientID: &clientID1,
		}
		pagination := caseDomain.NewPagination()

		cases, total, err := repo.List(ctx, filter, pagination)

		require.NoError(t, err)
		require.Equal(t, 2, total)
		require.Len(t, cases, 2)

		for _, c := range cases {
			require.Equal(t, clientID1, c.ClientID)
		}
	})

	t.Run("FilterByAssignedContact", func(t *testing.T) {
		// Clean up before each test
		ch.ClearCasesTable(t, ctx, testPool)
		// Create test cases
		clientID := uuid.New()
		contactID1 := uuid.New()
		contactID2 := uuid.New()

		now := time.Now()

		insertTestCase(t, ctx, clientID, &contactID1, now.Add(-1*time.Hour))
		insertTestCase(t, ctx, clientID, nil, now.Add(-2*time.Hour))
		insertTestCase(t, ctx, clientID, &contactID2, now.Add(-3*time.Hour))

		// Filter by assigned contact
		filter := caseDomain.CaseFilter{
			AssignedContactID: &contactID1,
		}
		pagination := caseDomain.NewPagination()

		cases, total, err := repo.List(ctx, filter, pagination)

		require.NoError(t, err)
		require.Equal(t, 1, total)
		require.Len(t, cases, 1)
		require.Equal(t, contactID1, *cases[0].AssignedContactID)
	})

	t.Run("FilterByNilAssignedContact", func(t *testing.T) {
		// Clean up before each test
		ch.ClearCasesTable(t, ctx, testPool)
		// Create test cases
		clientID := uuid.New()
		contactID1 := uuid.New()

		now := time.Now()

		insertTestCase(t, ctx, clientID, &contactID1, now.Add(-1*time.Hour))
		insertTestCase(t, ctx, clientID, nil, now.Add(-2*time.Hour))
		insertTestCase(t, ctx, clientID, nil, now.Add(-3*time.Hour))

		// Filter by nil assigned contact
		nilUUID := uuid.UUID{}
		filter := caseDomain.CaseFilter{
			AssignedContactID: &nilUUID,
		}
		pagination := caseDomain.NewPagination()

		cases, total, err := repo.List(ctx, filter, pagination)

		require.NoError(t, err)
		require.Equal(t, 2, total)
		require.Len(t, cases, 2)

		for _, c := range cases {
			require.Nil(t, c.AssignedContactID)
		}
	})

	t.Run("FilterByDateRange", func(t *testing.T) {
		// Clean up before each test
		ch.ClearCasesTable(t, ctx, testPool)
		// Create test cases with different creation dates
		clientID := uuid.New()

		now := time.Now()
		yesterday := now.Add(-24 * time.Hour)
		twoDaysAgo := now.Add(-48 * time.Hour)

		insertTestCase(t, ctx, clientID, nil, now)
		insertTestCase(t, ctx, clientID, nil, yesterday)
		insertTestCase(t, ctx, clientID, nil, twoDaysAgo)

		// Filter by date range
		fromDate := now.Add(-25 * time.Hour) // 25 hours ago to now
		filter := caseDomain.CaseFilter{
			DateCreatedFrom: &fromDate,
		}
		pagination := caseDomain.NewPagination()

		cases, total, err := repo.List(ctx, filter, pagination)

		require.NoError(t, err)
		require.Equal(t, 2, total) // yesterday and now
		require.Len(t, cases, 2)

		// Test with upper bound
		toDate := now.Add(-12 * time.Hour) // Up to 12 hours ago
		filter = caseDomain.CaseFilter{
			DateCreatedTo: &toDate,
		}

		cases, total, err = repo.List(ctx, filter, pagination)

		require.NoError(t, err)
		require.Equal(t, 2, total) // yesterday and two days ago
		require.Len(t, cases, 2)
	})

	t.Run("Pagination", func(t *testing.T) {
		// Clean up before each test
		ch.ClearCasesTable(t, ctx, testPool)
		// Create multiple test cases
		clientID := uuid.New()

		now := time.Now()
		for i := 0; i < 5; i++ {
			insertTestCase(t, ctx, clientID, nil, now.Add(time.Duration(i)*time.Hour))
		}

		// Test first page with 2 items
		pagination := caseDomain.Pagination{
			Page:     1,
			PageSize: 2,
		}

		cases, total, err := repo.List(ctx, caseDomain.CaseFilter{}, pagination)

		require.NoError(t, err)
		require.Equal(t, 5, total) // Total should be all items
		require.Len(t, cases, 2)   // Page should have 2 items

		// Test second page
		pagination.Page = 2

		cases, total, err = repo.List(ctx, caseDomain.CaseFilter{}, pagination)

		require.NoError(t, err)
		require.Equal(t, 5, total)
		require.Len(t, cases, 2)

		// Test third page (last item)
		pagination.Page = 3

		cases, total, err = repo.List(ctx, caseDomain.CaseFilter{}, pagination)

		require.NoError(t, err)
		require.Equal(t, 5, total)
		require.Len(t, cases, 1)
	})

	t.Run("InvalidPagination", func(t *testing.T) {
		// Clean up before each test
		ch.ClearCasesTable(t, ctx, testPool)
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
		ch.ClearCasesTable(t, ctx, testPool)
		// Create test cases with various properties
		clientID1 := uuid.New()
		clientID2 := uuid.New()
		contactID1 := uuid.New()

		now := time.Now()
		recent := now.Add(-1 * time.Hour)
		older := now.Add(-25 * time.Hour)

		// Recent cases for client1
		insertTestCase(t, ctx, clientID1, &contactID1, recent)
		insertTestCase(t, ctx, clientID1, nil, recent)

		// Older cases for client2
		insertTestCase(t, ctx, clientID2, nil, older)

		// Filter: client1 + recent date range
		fromDate := now.Add(-2 * time.Hour)
		filter := caseDomain.CaseFilter{
			ClientID:        &clientID1,
			DateCreatedFrom: &fromDate,
		}
		pagination := caseDomain.NewPagination()

		cases, total, err := repo.List(ctx, filter, pagination)

		require.NoError(t, err)
		require.Equal(t, 2, total)
		require.Len(t, cases, 2)

		for _, c := range cases {
			require.Equal(t, clientID1, c.ClientID)
			require.True(t, c.CreatedAt.After(fromDate))
		}
	})
}

// Helper function to insert test cases using helper functions
func insertTestCase(t *testing.T, ctx context.Context, clientID uuid.UUID, assignedContactID *uuid.UUID, createdAt time.Time) {
	t.Helper()

	caseID := uuid.New()
	caseEncx := &caseDomain.CaseEncx{
		ID:                   caseID,
		ClientID:             clientID,
		AssignedContactID:    assignedContactID,
		CreatedAt:            createdAt,
		TitleEncrypted:       []byte("encrypted_title"),
		DescriptionEncrypted: []byte("encrypted_description"),
		StatusEncrypted:      []byte("encrypted_status"),
		UpdatedAtEncrypted:   []byte("encrypted_updated_at"),
		DEKEncrypted:         []byte("encrypted_dek"),
		KeyVersion:           1,
		Metadata:             encx.EncryptionMetadata{},
	}

	err := ch.InsertCaseEncx(t, ctx, testPool, caseEncx)
	require.NoError(t, err, "Failed to insert test case")
}
