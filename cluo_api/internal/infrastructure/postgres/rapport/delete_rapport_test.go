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

func TestDeleteRapport(t *testing.T) {
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

	t.Run("successful deletion", func(t *testing.T) {
		ctx := context.Background()

		rapportHelpers.ClearRapportsTable(t, ctx, testPool)
		caseHelpers.ClearCasesTable(t, ctx, testPool)
		clientHelpers.ClearClientsTable(t, ctx, testPool)

		caseID := setupCase(t, ctx)
		rEncx := rapportHelpers.NewTestRapportEncx(t, caseID)
		err := rapportHelpers.InsertRapportEncx(t, ctx, testPool, rEncx)
		require.NoError(t, err)

		err = repo.DeleteRapport(ctx, caseID)
		assert.NoError(t, err)

		// Verify it's gone
		retrieved, err := rapportHelpers.GetRapportEncxByCaseID(t, ctx, testPool, caseID)
		assert.Error(t, err)
		assert.Nil(t, retrieved)
	})

	t.Run("delete non-existent rapport returns not found", func(t *testing.T) {
		ctx := context.Background()

		rapportHelpers.ClearRapportsTable(t, ctx, testPool)

		err := repo.DeleteRapport(ctx, uuid.New())
		assert.Error(t, err)
	})

	t.Run("context cancellation", func(t *testing.T) {
		ctx, cancel := context.WithCancel(context.Background())
		cancel()

		err := repo.DeleteRapport(ctx, uuid.New())
		assert.Error(t, err)
	})
}
