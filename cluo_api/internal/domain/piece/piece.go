package piece

import (
	"time"

	"github.com/google/uuid"
)

type Piece struct {
	ID         uuid.UUID
	CaseID     uuid.UUID
	Filename   string    `encx:"encrypt"`
	StorageKey string
	MimeType   string
	SizeBytes  int64
	Notes      string    `encx:"encrypt"`
	CreatedAt  time.Time
	UpdatedAt  time.Time
}

type PieceEncx struct {
	ID                uuid.UUID `db:"id"`
	CaseID            uuid.UUID `db:"case_id"`
	FilenameEncrypted []byte    `db:"filename_encrypted"`
	StorageKey        string    `db:"storage_key"`
	MimeType          string    `db:"mime_type"`
	SizeBytes         int64     `db:"size_bytes"`
	NotesEncrypted    []byte    `db:"notes_encrypted"`
	DEKEncrypted      []byte    `db:"dek_encrypted"`
	KeyVersion        int       `db:"key_version"`
	CreatedAt         time.Time `db:"created_at"`
	UpdatedAt         time.Time `db:"updated_at"`
}
