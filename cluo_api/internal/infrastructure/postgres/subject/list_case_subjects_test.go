package subjectRepository_test

import (
	"context"
	"testing"

	subjectHelpers "github.com/hengadev/cluo_api/test/helpers/subject"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestListCaseSubjects(t *testing.T) {
	if testPool == nil || repo == nil {
		t.Skip("Test database or repository not initialized")
	}

	t.Run("successful list with pagination", func(t *testing.T) {
		ctx := context.Background()

		subjectHelpers.ClearCaseSubjectsTable(t, ctx, testPool)

		// Create 25 subjects
		for i := 0; i < 25; i++ {
			sEncx := subjectHelpers.NewTestSubjectEncx(t)
			err := subjectHelpers.InsertSubjectEncx(t, ctx, testPool, sEncx)
			require.NoError(t, err)
		}

		// Test first page
		subjects, total, err := repo.ListCaseSubjects(ctx, 1, 10)
		assert.NoError(t, err)
		assert.Equal(t, 25, total)
		assert.Len(t, subjects, 10)

		// Test second page
		subjects, total, err = repo.ListCaseSubjects(ctx, 2, 10)
		assert.NoError(t, err)
		assert.Equal(t, 25, total)
		assert.Len(t, subjects, 10)

		// Test third page (partial)
		subjects, total, err = repo.ListCaseSubjects(ctx, 3, 10)
		assert.NoError(t, err)
		assert.Equal(t, 25, total)
		assert.Len(t, subjects, 5)

		// Test fourth page (empty)
		subjects, total, err = repo.ListCaseSubjects(ctx, 4, 10)
		assert.NoError(t, err)
		assert.Equal(t, 25, total)
		assert.Len(t, subjects, 0)
	})

	t.Run("empty table returns empty list", func(t *testing.T) {
		ctx := context.Background()

		subjectHelpers.ClearCaseSubjectsTable(t, ctx, testPool)

		subjects, total, err := repo.ListCaseSubjects(ctx, 1, 10)
		assert.NoError(t, err)
		assert.Equal(t, 0, total)
		assert.Len(t, subjects, 0)
	})

	t.Run("context cancellation", func(t *testing.T) {
		ctx, cancel := context.WithCancel(context.Background())
		cancel()

		_, _, err := repo.ListCaseSubjects(ctx, 1, 10)
		assert.Error(t, err)
	})
}
