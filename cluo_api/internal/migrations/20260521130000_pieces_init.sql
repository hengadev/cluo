-- +goose Up
CREATE TABLE IF NOT EXISTS cases.pieces (
    id UUID PRIMARY KEY,
    case_id UUID NOT NULL,
    filename_encrypted BYTEA NOT NULL,
    storage_key TEXT NOT NULL,
    mime_type TEXT NOT NULL,
    size_bytes BIGINT NOT NULL,
    notes_encrypted BYTEA,
    dek_encrypted BYTEA NOT NULL,
    key_version INTEGER NOT NULL,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    CONSTRAINT fk_pieces_case_id
        FOREIGN KEY (case_id) REFERENCES cases.cases(id) ON DELETE CASCADE,
    CONSTRAINT chk_pieces_key_version_positive CHECK (key_version >= 0),
    CONSTRAINT chk_pieces_size_positive CHECK (size_bytes > 0)
);
CREATE INDEX IF NOT EXISTS idx_pieces_case_id ON cases.pieces(case_id);
CREATE INDEX IF NOT EXISTS idx_pieces_storage_key ON cases.pieces(storage_key);

-- +goose Down
DROP TABLE IF EXISTS cases.pieces;
