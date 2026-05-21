package rapportRepository_test

import (
	"context"
	"testing"

	caseHelpers "github.com/hengadev/cluo_api/test/helpers/investigation"
	clientHelpers "github.com/hengadev/cluo_api/test/helpers/client"
	rapportHelpers "github.com/hengadev/cluo_api/test/helpers/rapport"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestCreateRapport(t *testing.T) {
	if testPool == nil || repo == nil {
		t.Skip("Test database or repository not initialized")
	}

	setupCase := func(t *testing.T, ctx context.Context) uuid.UUID {
		t.Helper()
		clientEncx := clientHelpers.NewTestClientEncx(t)
		err := clientHelpers.InsertClientEncx(t, ctx, testPool, clientEncx)
		require.NoError(t, err)

		caseEncx := caseHelpers.NewTestCaseEncxWithClientID(t, clientEncx.ID)
		err = caseHelpers.InsertCaseEncx(t, ctx, testPool, caseEncx)
		require.NoError(t, err)
		return caseEncx.ID
	}

	t.Run("successful creation", func(t *testing.T) {
		ctx := context.Background()

		rapportHelpers.ClearRapportsTable(t, ctx, testPool)
		caseHelpers.ClearCasesTable(t, ctx, testPool)
		clientHelpers.ClearClientsTable(t, ctx, testPool)

		caseID := setupCase(t, ctx)
		rEncx := rapportHelpers.NewTestRapportEncx(t, caseID)

		err := repo.CreateRapport(ctx, rEncx)
		assert.NoError(t, err, "Failed to create rapport")

		retrieved, err := rapportHelpers.GetRapportEncxByCaseID(t, ctx, testPool, caseID)
		assert.NoError(t, err)
		assert.Equal(t, rEncx.ID, retrieved.ID)
		assert.Equal(t, rEncx.CaseID, retrieved.CaseID)
		assert.Equal(t, rEncx.ContentEncrypted, retrieved.ContentEncrypted)
		assert.Equal(t, rEncx.DEKEncrypted, retrieved.DEKEncrypted)
		assert.Equal(t, rEncx.KeyVersion, retrieved.KeyVersion)
	})

	t.Run("nil rapport returns error", func(t *testing.T) {
		ctx := context.Background()

		err := repo.CreateRapport(ctx, nil)
		assert.Error(t, err)
	})

	t.Run("unique constraint: one rapport per case", func(t *testing.T) {
		ctx := context.Background()

		rapportHelpers.ClearRapportsTable(t, ctx, testPool)
		caseHelpers.ClearCasesTable(t, ctx, testPool)
		clientHelpers.ClearClientsTable(t, ctx, testPool)

		caseID := setupCase(t, ctx)
		rEncx1 := rapportHelpers.NewTestRapportEncx(t, caseID)

		err := repo.CreateRapport(ctx, rEncx1)
		require.NoError(t, err)

		// Second rapport for the same case must fail
		rEncx2 := rapportHelpers.NewTestRapportEncx(t, caseID)
		err = repo.CreateRapport(ctx, rEncx2)
		assert.Error(t, err, "Expected unique constraint violation")
	})

	t.Run("context cancellation", func(t *testing.T) {
		ctx, cancel := context.WithCancel(context.Background())
		cancel()

		rEncx := rapportHelpers.NewTestRapportEncx(t, uuid.New())
		err := repo.CreateRapport(ctx, rEncx)
		assert.Error(t, err)
	})
}
