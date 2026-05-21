-- +goose Up
CREATE TABLE IF NOT EXISTS cases.case_access_tokens (
    id UUID PRIMARY KEY,
    case_id UUID NOT NULL,
    token_hash VARCHAR(64) NOT NULL,
    expires_at TIMESTAMPTZ NOT NULL,
    revoked_at TIMESTAMPTZ NULL,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    CONSTRAINT fk_tokens_case_id
        FOREIGN KEY (case_id) REFERENCES cases.cases(id) ON DELETE CASCADE,
    CONSTRAINT uq_tokens_hash UNIQUE (token_hash)
);
CREATE INDEX IF NOT EXISTS idx_tokens_case_id ON cases.case_access_tokens(case_id);
CREATE INDEX IF NOT EXISTS idx_tokens_hash ON cases.case_access_tokens(token_hash);

-- +goose Down
DROP TABLE IF EXISTS cases.case_access_tokens;
