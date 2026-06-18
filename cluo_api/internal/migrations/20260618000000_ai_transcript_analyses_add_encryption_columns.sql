-- +goose Up
-- SQL in this section is executed when the migration is applied.

-- The transcript analysis persistence layer was never wired up, so this table
-- has been unused since it was created. Add the DEK/key-version columns every
-- other encx-encrypted table carries, needed to decrypt rows after they're written.

ALTER TABLE ai.transcript_analyses
    ADD COLUMN dek_encrypted BYTEA,
    ADD COLUMN key_version INTEGER;

UPDATE ai.transcript_analyses SET dek_encrypted = ''::bytea, key_version = 0 WHERE dek_encrypted IS NULL;

ALTER TABLE ai.transcript_analyses
    ALTER COLUMN dek_encrypted SET NOT NULL,
    ALTER COLUMN key_version SET NOT NULL,
    ADD CONSTRAINT chk_analyses_key_version_positive CHECK (key_version >= 0);

-- +goose Down
-- SQL in this section is executed when the migration is rolled back.

ALTER TABLE ai.transcript_analyses
    DROP CONSTRAINT IF EXISTS chk_analyses_key_version_positive,
    DROP COLUMN IF EXISTS dek_encrypted,
    DROP COLUMN IF EXISTS key_version;
