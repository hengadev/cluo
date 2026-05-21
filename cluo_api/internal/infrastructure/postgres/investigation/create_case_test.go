package investigationRepository_test

import (
	"context"
	"testing"

	"github.com/hengadev/cluo_api/internal/domain/investigation"
	caseHelpers "github.com/hengadev/cluo_api/test/helpers/investigation"
	clientHelpers "github.com/hengadev/cluo_api/test/helpers/client"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// make test-func TEST_NAME=TestCreateCase TEST_PATH=internal/infrastructure/postgres/case/create_case_test.go

func TestCreateCase(t *testing.T) {
	if testPool == nil || repo == nil {
		t.Skip("Test database or repository not initialized")
	}

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

	t.Run("successful creation", func(t *testing.T) {
		ctx := context.Background()

		caseHelpers.ClearCasesTable(t, ctx, testPool)
		clientHelpers.ClearContactsTable(t, ctx, testPool)
		clientHelpers.ClearClientsTable(t, ctx, testPool)

		// Setup client and contact in database
		clientID := setupClient(t, ctx)
		contactID := setupContact(t, ctx, clientID)

		// Create test caseEncx data using helper
		caseEncx := caseHelpers.NewTestCaseEncx(t)
		caseEncx.ClientID = clientID
		caseEncx.AssignedContactID = &contactID

		// Test successful case creation using the global repo
		err := repo.CreateCase(ctx, caseEncx)
		assert.NoError(t, err, "Failed to create case")

		// Verify the case was inserted by retrieving it
		retrievedCaseEncx, err := caseHelpers.GetCaseEncxByID(t, ctx, testPool, caseEncx.ID)
		assert.NoError(t, err, "Failed to retrieve inserted case")

		// Verify field values
		assert.Equal(t, caseEncx.ID, retrievedCaseEncx.ID, "Case ID should match")
		assert.Equal(t, caseEncx.ClientID, retrievedCaseEncx.ClientID, "Client ID should match")
		assert.Equal(t, caseEncx.AssignedContactID, retrievedCaseEncx.AssignedContactID, "Assigned Contact ID should match")
		assert.Equal(t, caseEncx.CaseTypeID, retrievedCaseEncx.CaseTypeID, "CaseTypeID should match")
		assert.Equal(t, caseEncx.ExternalReferenceEncrypted, retrievedCaseEncx.ExternalReferenceEncrypted, "ExternalReferenceEncrypted should match")
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

		caseHelpers.ClearCasesTable(t, ctx, testPool)
		clientHelpers.ClearClientsTable(t, ctx, testPool)
		clientHelpers.ClearContactsTable(t, ctx, testPool)

		// Setup client and contact in database
		clientID := setupClient(t, ctx)
		contactID := setupContact(t, ctx, clientID)

		// Create first test case
		caseEncx1 := caseHelpers.NewTestCaseEncx(t)
		caseEncx1.ClientID = clientID
		caseEncx1.AssignedContactID = &contactID

		// Insert first case using the global repo (setup - should stop if fails)
		err := repo.CreateCase(ctx, caseEncx1)
		require.NoError(t, err, "Failed to create first case")

		// Try to insert case with same ID (should fail)
		caseEncx2 := caseHelpers.NewTestCaseEncx(t)
		caseEncx2.ID = caseEncx1.ID // Same ID, different data
		caseEncx2.ClientID = clientID
		caseEncx2.AssignedContactID = &contactID

		err = repo.CreateCase(ctx, caseEncx2)
		assert.Error(t, err, "Expected error when creating case with duplicate ID")

		// Check that it's a database constraint violation (expected for duplicate ID)
		errStr := err.Error()
		assert.True(t, contains(errStr, "duplicate") || contains(errStr, "unique") || contains(errStr, "constraint"),
			"Expected constraint violation error, got: %v", err)
	})

	t.Run("empty required fields", func(t *testing.T) {
		ctx := context.Background()

		caseHelpers.ClearCasesTable(t, ctx, testPool)

		tests := []struct {
			name  string
			case_ *investigation.InvestigationEncx
		}{
			{
				name: "nil client ID",
				case_: &investigation.InvestigationEncx{
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

		caseEncx := caseHelpers.NewTestCaseEncx(t)

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
