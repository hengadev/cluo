package subjectRepository_test

import (
	"context"
	"testing"

	"github.com/google/uuid"
	subjectHelpers "github.com/hengadev/cluo_api/test/helpers/subject"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestGetCaseSubjectByID(t *testing.T) {
	if testPool == nil || repo == nil {
		t.Skip("Test database or repository not initialized")
	}

	t.Run("successful retrieval", func(t *testing.T) {
		ctx := context.Background()

		subjectHelpers.ClearCaseSubjectsTable(t, ctx, testPool)

		sEncx := subjectHelpers.NewTestSubjectEncx(t)
		err := subjectHelpers.InsertSubjectEncx(t, ctx, testPool, sEncx)
		require.NoError(t, err)

		retrieved, err := repo.GetCaseSubjectByID(ctx, sEncx.ID)
		assert.NoError(t, err)
		assert.Equal(t, sEncx.ID, retrieved.ID)
		assert.Equal(t, sEncx.LastnameEncrypted, retrieved.LastnameEncrypted)
		assert.Equal(t, sEncx.LastnameHash, retrieved.LastnameHash)
		assert.Equal(t, sEncx.FirstnameEncrypted, retrieved.FirstnameEncrypted)
		assert.Equal(t, sEncx.FirstnameHash, retrieved.FirstnameHash)
		assert.Equal(t, sEncx.EmailEncrypted, retrieved.EmailEncrypted)
		assert.Equal(t, sEncx.EmailHash, retrieved.EmailHash)
		assert.Equal(t, sEncx.PhoneEncrypted, retrieved.PhoneEncrypted)
		assert.Equal(t, sEncx.CityEncrypted, retrieved.CityEncrypted)
		assert.Equal(t, sEncx.CityHash, retrieved.CityHash)
		assert.Equal(t, sEncx.PostalCodeEncrypted, retrieved.PostalCodeEncrypted)
		assert.Equal(t, sEncx.PostalCodeHash, retrieved.PostalCodeHash)
		assert.Equal(t, sEncx.Address1Encrypted, retrieved.Address1Encrypted)
		assert.Equal(t, sEncx.Address1Hash, retrieved.Address1Hash)
		assert.Equal(t, sEncx.Address2Encrypted, retrieved.Address2Encrypted)
		assert.Equal(t, sEncx.Address2Hash, retrieved.Address2Hash)
		assert.Equal(t, sEncx.OccupationEncrypted, retrieved.OccupationEncrypted)
		assert.Equal(t, sEncx.OccupationHash, retrieved.OccupationHash)
		assert.Equal(t, sEncx.NotesEncrypted, retrieved.NotesEncrypted)
		assert.Equal(t, sEncx.DEKEncrypted, retrieved.DEKEncrypted)
		assert.Equal(t, sEncx.KeyVersion, retrieved.KeyVersion)
	})

	t.Run("get non-existent subject returns not found", func(t *testing.T) {
		ctx := context.Background()

		subjectHelpers.ClearCaseSubjectsTable(t, ctx, testPool)

		_, err := repo.GetCaseSubjectByID(ctx, uuid.New())
		assert.Error(t, err)
	})

	t.Run("context cancellation", func(t *testing.T) {
		ctx, cancel := context.WithCancel(context.Background())
		cancel()

		_, err := repo.GetCaseSubjectByID(ctx, uuid.New())
		assert.Error(t, err)
	})
}
