package caseRepository_test

import (
	"context"
	"testing"

	caseDomain "github.com/hengadev/cluo_api/internal/domain/case"
	ch "github.com/hengadev/cluo_api/test/helpers/case"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// make test-func TEST_NAME=TestCreateCase TEST_PATH=internal/infrastructure/postgres/case/create_case_test.go

func TestCreateCase(t *testing.T) {
	if testPool == nil || repo == nil {
		t.Skip("Test database or repository not initialized")
	}

	t.Run("successful creation", func(t *testing.T) {
		ctx := context.Background()

		ch.ClearCasesTable(t, ctx, testPool)

		// Create test caseEncx data using helper
		caseEncx := ch.NewTestCaseEncx(t)

		// Test successful case creation using the global repo
		err := repo.CreateCase(ctx, caseEncx)
		assert.NoError(t, err, "Failed to create case")

		// Verify the case was inserted by retrieving it
		retrievedCaseEncx, err := ch.GetCaseEncxByID(t, ctx, testPool, caseEncx.ID)
		assert.NoError(t, err, "Failed to retrieve inserted case")

		// Verify field values
		assert.Equal(t, caseEncx.ID, retrievedCaseEncx.ID, "Case ID should match")
		assert.Equal(t, caseEncx.ClientID, retrievedCaseEncx.ClientID, "Client ID should match")
		assert.Equal(t, caseEncx.AssignedContactID, retrievedCaseEncx.AssignedContactID, "Assigned Contact ID should match")
		assert.Equal(t, caseEncx.KeyVersion, retrievedCaseEncx.KeyVersion, "Key version should match")
	})

	t.Run("with nil case", func(t *testing.T) {
		ctx := context.Background()

		// Test with nil case
		err := repo.CreateCase(ctx, nil)
		assert.Error(t, err, "Expected error when creating case with nil input")
	})

	t.Run("duplicate ID", func(t *testing.T) {
		ctx := context.Background()

		ch.ClearCasesTable(t, ctx, testPool)

		// Create first test case
		caseEncx1 := ch.NewTestCaseEncx(t)

		// Insert first case using the global repo (setup - should stop if fails)
		err := repo.CreateCase(ctx, caseEncx1)
		require.NoError(t, err, "Failed to create first case")

		// Try to insert case with same ID (should fail)
		caseEncx2 := ch.NewTestCaseEncx(t)
		caseEncx2.ID = caseEncx1.ID // Same ID, different data

		err = repo.CreateCase(ctx, caseEncx2)
		assert.Error(t, err, "Expected error when creating case with duplicate ID")

		// Check that it's a database constraint violation (expected for duplicate ID)
		errStr := err.Error()
		assert.True(t, contains(errStr, "duplicate") || contains(errStr, "unique") || contains(errStr, "constraint"),
			"Expected constraint violation error, got: %v", err)
	})

	t.Run("empty required fields", func(t *testing.T) {
		ctx := context.Background()

		ch.ClearCasesTable(t, ctx, testPool)

		tests := []struct {
			name   string
			case_  *caseDomain.CaseEncx
		}{
			{
				name: "nil client ID",
				case_: &caseDomain.CaseEncx{
					ID: uuid.Nil, // Invalid UUID
				},
			},
		}

		for _, tt := range tests {
			t.Run(tt.name, func(t *testing.T) {
				err := repo.CreateCase(ctx, tt.case_)
				assert.Error(t, err, "Expected error for %s, but got nil", tt.name)
			})
		}
	})

	t.Run("context cancellation", func(t *testing.T) {
		// Create a context that will be cancelled
		ctx, cancel := context.WithCancel(context.Background())
		cancel() // Cancel immediately

		caseEncx := ch.NewTestCaseEncx(t)

		err := repo.CreateCase(ctx, caseEncx)
		assert.Error(t, err, "Expected context cancellation error, but got nil")
	})
}

// contains is a helper function to check if a string contains a substring
func contains(s, substr string) bool {
	return len(s) >= len(substr) && (s == substr || len(substr) == 0 ||
		(len(s) > len(substr) && (s[:len(substr)] == substr || s[len(s)-len(substr):] == substr ||
		func() bool {
			for i := 1; i <= len(s)-len(substr); i++ {
				if s[i:i+len(substr)] == substr {
					return true
				}
			}
			return false
		}())))
}