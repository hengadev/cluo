package subjectRepository_test

import (
	"context"
	"testing"

	"github.com/google/uuid"
	subjectHelpers "github.com/hengadev/cluo_api/test/helpers/subject"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestDeleteCaseSubject(t *testing.T) {
	if testPool == nil || repo == nil {
		t.Skip("Test database or repository not initialized")
	}

	t.Run("successful deletion", func(t *testing.T) {
		ctx := context.Background()

		subjectHelpers.ClearCaseSubjectsTable(t, ctx, testPool)

		sEncx := subjectHelpers.NewTestSubjectEncx(t)
		err := subjectHelpers.InsertSubjectEncx(t, ctx, testPool, sEncx)
		require.NoError(t, err)

		err = repo.DeleteCaseSubject(ctx, sEncx.ID)
		assert.NoError(t, err)

		// Verify it's gone
		retrieved, err := subjectHelpers.GetSubjectEncxByID(t, ctx, testPool, sEncx.ID)
		assert.Error(t, err)
		assert.Nil(t, retrieved)
	})

	t.Run("delete non-existent subject returns not found", func(t *testing.T) {
		ctx := context.Background()

		subjectHelpers.ClearCaseSubjectsTable(t, ctx, testPool)

		err := repo.DeleteCaseSubject(ctx, uuid.New())
		assert.Error(t, err)
	})

	t.Run("context cancellation", func(t *testing.T) {
		ctx, cancel := context.WithCancel(context.Background())
		cancel()

		err := repo.DeleteCaseSubject(ctx, uuid.New())
		assert.Error(t, err)
	})
}
