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

func TestUpdateRapport(t *testing.T) {
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

	t.Run("successful update", func(t *testing.T) {
		ctx := context.Background()

		rapportHelpers.ClearRapportsTable(t, ctx, testPool)
		caseHelpers.ClearCasesTable(t, ctx, testPool)
		clientHelpers.ClearClientsTable(t, ctx, testPool)

		caseID := setupCase(t, ctx)
		rEncx := rapportHelpers.NewTestRapportEncx(t, caseID)
		err := rapportHelpers.InsertRapportEncx(t, ctx, testPool, rEncx)
		require.NoError(t, err)

		// Modify content
		rEncx.ContentEncrypted = []byte("updated_content_encrypted")
		rEncx.DEKEncrypted = []byte("updated_dek_encrypted")
		rEncx.KeyVersion = 2

		err = repo.UpdateRapport(ctx, rEncx)
		assert.NoError(t, err)

		retrieved, err := rapportHelpers.GetRapportEncxByCaseID(t, ctx, testPool, caseID)
		require.NoError(t, err)
		assert.Equal(t, []byte("updated_content_encrypted"), retrieved.ContentEncrypted)
		assert.Equal(t, []byte("updated_dek_encrypted"), retrieved.DEKEncrypted)
		assert.Equal(t, 2, retrieved.KeyVersion)
	})

	t.Run("nil rapport returns error", func(t *testing.T) {
		ctx := context.Background()

		err := repo.UpdateRapport(ctx, nil)
		assert.Error(t, err)
	})

	t.Run("update non-existent rapport returns not found", func(t *testing.T) {
		ctx := context.Background()

		rapportHelpers.ClearRapportsTable(t, ctx, testPool)

		rEncx := rapportHelpers.NewTestRapportEncx(t, uuid.New())
		err := repo.UpdateRapport(ctx, rEncx)
		assert.Error(t, err)
	})

	t.Run("context cancellation", func(t *testing.T) {
		ctx, cancel := context.WithCancel(context.Background())
		cancel()

		rEncx := rapportHelpers.NewTestRapportEncx(t, uuid.New())
		err := repo.UpdateRapport(ctx, rEncx)
		assert.Error(t, err)
	})
}
