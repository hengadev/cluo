package subjectRepository_test

import (
	"context"
	"testing"

	subjectHelpers "github.com/hengadev/cluo_api/test/helpers/subject"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestUpdateCaseSubject(t *testing.T) {
	if testPool == nil || repo == nil {
		t.Skip("Test database or repository not initialized")
	}

	t.Run("successful update", func(t *testing.T) {
		ctx := context.Background()

		subjectHelpers.ClearCaseSubjectsTable(t, ctx, testPool)

		sEncx := subjectHelpers.NewTestSubjectEncx(t)
		err := subjectHelpers.InsertSubjectEncx(t, ctx, testPool, sEncx)
		require.NoError(t, err)

		// Modify fields
		sEncx.LastnameEncrypted = []byte("updated_lastname_encrypted")
		sEncx.LastnameHash = "updated_lastname_hash"
		sEncx.FirstnameEncrypted = []byte("updated_firstname_encrypted")
		sEncx.FirstnameHash = "updated_firstname_hash"
		sEncx.EmailEncrypted = []byte("updated_email_encrypted")
		sEncx.EmailHash = "updated_email_hash"
		sEncx.PhoneEncrypted = []byte("updated_phone_encrypted")
		sEncx.CityEncrypted = []byte("updated_city_encrypted")
		sEncx.CityHash = "updated_city_hash"
		sEncx.PostalCodeEncrypted = []byte("updated_postal_code_encrypted")
		sEncx.PostalCodeHash = "updated_postal_code_hash"
		sEncx.Address1Encrypted = []byte("updated_address1_encrypted")
		sEncx.Address1Hash = "updated_address1_hash"
		sEncx.Address2Encrypted = []byte("updated_address2_encrypted")
		sEncx.Address2Hash = "updated_address2_hash"
		sEncx.OccupationEncrypted = []byte("updated_occupation_encrypted")
		sEncx.OccupationHash = "updated_occupation_hash"
		sEncx.NotesEncrypted = []byte("updated_notes_encrypted")
		sEncx.DEKEncrypted = []byte("updated_dek_encrypted")
		sEncx.KeyVersion = 2

		err = repo.UpdateCaseSubject(ctx, sEncx)
		assert.NoError(t, err)

		retrieved, err := subjectHelpers.GetSubjectEncxByID(t, ctx, testPool, sEncx.ID)
		require.NoError(t, err)
		assert.Equal(t, []byte("updated_lastname_encrypted"), retrieved.LastnameEncrypted)
		assert.Equal(t, "updated_lastname_hash", retrieved.LastnameHash)
		assert.Equal(t, []byte("updated_firstname_encrypted"), retrieved.FirstnameEncrypted)
		assert.Equal(t, "updated_firstname_hash", retrieved.FirstnameHash)
		assert.Equal(t, []byte("updated_email_encrypted"), retrieved.EmailEncrypted)
		assert.Equal(t, "updated_email_hash", retrieved.EmailHash)
		assert.Equal(t, []byte("updated_phone_encrypted"), retrieved.PhoneEncrypted)
		assert.Equal(t, []byte("updated_city_encrypted"), retrieved.CityEncrypted)
		assert.Equal(t, "updated_city_hash", retrieved.CityHash)
		assert.Equal(t, []byte("updated_postal_code_encrypted"), retrieved.PostalCodeEncrypted)
		assert.Equal(t, "updated_postal_code_hash", retrieved.PostalCodeHash)
		assert.Equal(t, []byte("updated_address1_encrypted"), retrieved.Address1Encrypted)
		assert.Equal(t, "updated_address1_hash", retrieved.Address1Hash)
		assert.Equal(t, []byte("updated_address2_encrypted"), retrieved.Address2Encrypted)
		assert.Equal(t, "updated_address2_hash", retrieved.Address2Hash)
		assert.Equal(t, []byte("updated_occupation_encrypted"), retrieved.OccupationEncrypted)
		assert.Equal(t, "updated_occupation_hash", retrieved.OccupationHash)
		assert.Equal(t, []byte("updated_notes_encrypted"), retrieved.NotesEncrypted)
		assert.Equal(t, []byte("updated_dek_encrypted"), retrieved.DEKEncrypted)
		assert.Equal(t, 2, retrieved.KeyVersion)
	})

	t.Run("nil subject returns error", func(t *testing.T) {
		ctx := context.Background()

		err := repo.UpdateCaseSubject(ctx, nil)
		assert.Error(t, err)
	})

	t.Run("update non-existent subject returns not found", func(t *testing.T) {
		ctx := context.Background()

		subjectHelpers.ClearCaseSubjectsTable(t, ctx, testPool)

		sEncx := subjectHelpers.NewTestSubjectEncx(t)
		err := repo.UpdateCaseSubject(ctx, sEncx)
		assert.Error(t, err)
	})

	t.Run("context cancellation", func(t *testing.T) {
		ctx, cancel := context.WithCancel(context.Background())
		cancel()

		sEncx := subjectHelpers.NewTestSubjectEncx(t)
		err := repo.UpdateCaseSubject(ctx, sEncx)
		assert.Error(t, err)
	})
}
