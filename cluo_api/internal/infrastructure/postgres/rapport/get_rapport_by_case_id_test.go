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

func TestGetRapportByCaseID(t *testing.T) {
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

	t.Run("successful retrieval", func(t *testing.T) {
		ctx := context.Background()

		rapportHelpers.ClearRapportsTable(t, ctx, testPool)
		caseHelpers.ClearCasesTable(t, ctx, testPool)
		clientHelpers.ClearClientsTable(t, ctx, testPool)

		caseID := setupCase(t, ctx)
		rEncx := rapportHelpers.NewTestRapportEncx(t, caseID)
		err := rapportHelpers.InsertRapportEncx(t, ctx, testPool, rEncx)
		require.NoError(t, err)

		retrieved, err := repo.GetRapportByCaseID(ctx, caseID)
		assert.NoError(t, err)
		require.NotNil(t, retrieved)
		assert.Equal(t, rEncx.ID, retrieved.ID)
		assert.Equal(t, rEncx.CaseID, retrieved.CaseID)
		assert.Equal(t, rEncx.ContentEncrypted, retrieved.ContentEncrypted)
		assert.Equal(t, rEncx.DEKEncrypted, retrieved.DEKEncrypted)
		assert.Equal(t, rEncx.KeyVersion, retrieved.KeyVersion)
	})

	t.Run("not found returns error", func(t *testing.T) {
		ctx := context.Background()

		rapportHelpers.ClearRapportsTable(t, ctx, testPool)

		retrieved, err := repo.GetRapportByCaseID(ctx, uuid.New())
		assert.Error(t, err)
		assert.Nil(t, retrieved)
	})

	t.Run("context cancellation", func(t *testing.T) {
		ctx, cancel := context.WithCancel(context.Background())
		cancel()

		_, err := repo.GetRapportByCaseID(ctx, uuid.New())
		assert.Error(t, err)
	})
}
