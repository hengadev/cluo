-- +goose Up
CREATE TABLE IF NOT EXISTS cases.rapports (
    id UUID PRIMARY KEY,
    case_id UUID NOT NULL,
    content_encrypted BYTEA NOT NULL,
    dek_encrypted BYTEA NOT NULL,
    key_version INTEGER NOT NULL,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    CONSTRAINT fk_rapports_case_id
        FOREIGN KEY (case_id) REFERENCES cases.cases(id) ON DELETE CASCADE,
    CONSTRAINT uq_rapports_case_id UNIQUE (case_id),
    CONSTRAINT chk_rapports_key_version_positive CHECK (key_version >= 0)
);
CREATE INDEX IF NOT EXISTS idx_rapports_case_id ON cases.rapports(case_id);

-- +goose Down
DROP TABLE IF EXISTS cases.rapports;
