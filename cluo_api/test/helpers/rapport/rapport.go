package rapportHelpers

import (
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/hengadev/cluo_api/internal/domain/rapport"
)

// NewTestRapportEncx creates a mock RapportEncx with basic test data.
func NewTestRapportEncx(t *testing.T, caseID uuid.UUID) *rapport.RapportEncx {
	t.Helper()
	return &rapport.RapportEncx{
		ID:               uuid.New(),
		CaseID:           caseID,
		ContentEncrypted: []byte("content_encrypted"),
		DEKEncrypted:     []byte("dek_encrypted"),
		KeyVersion:       1,
		CreatedAt:        time.Now(),
		UpdatedAt:        time.Now(),
	}
}
