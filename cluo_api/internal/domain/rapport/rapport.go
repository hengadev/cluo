package rapport

import (
	"time"

	"github.com/google/uuid"
)

// Rapport is the final investigation report. Content is an opaque encrypted blob.
type Rapport struct {
	ID        uuid.UUID
	CaseID    uuid.UUID
	Content   []byte    // raw TipTap JSON bytes (plaintext)
	CreatedAt time.Time
	UpdatedAt time.Time
}

type RapportEncx struct {
	ID               uuid.UUID `db:"id"`
	CaseID           uuid.UUID `db:"case_id"`
	ContentEncrypted []byte    `db:"content_encrypted"`
	DEKEncrypted     []byte    `db:"dek_encrypted"`
	KeyVersion       int       `db:"key_version"`
	CreatedAt        time.Time `db:"created_at"`
	UpdatedAt        time.Time `db:"updated_at"`
}
