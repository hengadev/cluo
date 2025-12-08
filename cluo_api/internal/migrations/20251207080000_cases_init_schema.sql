-- +goose Up
-- SQL in this section is executed when the migration is applied.

-- Create cases schema if it doesn't exist
CREATE SCHEMA IF NOT EXISTS cases;

-- Cases table
CREATE TABLE IF NOT EXISTS cases.cases (
    id UUID PRIMARY KEY,
    client_id UUID NOT NULL,
    assigned_contact_id UUID,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    title_encrypted BYTEA NOT NULL,
    description_encrypted BYTEA NOT NULL,
    status_encrypted BYTEA NOT NULL,
    updated_at_encrypted BYTEA,
    dek_encrypted BYTEA NOT NULL,
    key_version INTEGER NOT NULL,
    metadata JSONB DEFAULT '{}',

    -- Foreign key constraints
    CONSTRAINT fk_cases_client_id
        FOREIGN KEY (client_id) REFERENCES clients.clients(id) ON DELETE RESTRICT,
    CONSTRAINT fk_cases_assigned_contact_id
        FOREIGN KEY (assigned_contact_id) REFERENCES clients.contacts(id) ON DELETE SET NULL
);

-- Cases table indexes
CREATE INDEX IF NOT EXISTS idx_cases_client_id ON cases.cases(client_id);
CREATE INDEX IF NOT EXISTS idx_cases_assigned_contact_id ON cases.cases(assigned_contact_id);
CREATE INDEX IF NOT EXISTS idx_cases_created_at ON cases.cases(created_at);
CREATE INDEX IF NOT EXISTS idx_cases_key_version ON cases.cases(key_version);
CREATE INDEX IF NOT EXISTS idx_cases_metadata ON cases.cases USING gin(metadata);

-- Cases table business constraints
ALTER TABLE cases.cases ADD CONSTRAINT chk_cases_key_version_positive
    CHECK (key_version >= 0);

ALTER TABLE cases.cases ADD CONSTRAINT chk_cases_client_id_not_empty
    CHECK (client_id IS NOT NULL);

-- Add comments for documentation
COMMENT ON SCHEMA cases IS 'Case management schema';
COMMENT ON TABLE cases.cases IS 'Stores encrypted case information';
COMMENT ON COLUMN cases.cases.id IS 'Unique identifier for the case';
COMMENT ON COLUMN cases.cases.client_id IS 'Foreign key reference to clients.clients.id';
COMMENT ON COLUMN cases.cases.assigned_contact_id IS 'Foreign key reference to clients.contacts.id (nullable)';
COMMENT ON COLUMN cases.cases.created_at IS 'Case creation timestamp';
COMMENT ON COLUMN cases.cases.title_encrypted IS 'Encrypted case title';
COMMENT ON COLUMN cases.cases.description_encrypted IS 'Encrypted case description';
COMMENT ON COLUMN cases.cases.status_encrypted IS 'Encrypted case status';
COMMENT ON COLUMN cases.cases.updated_at_encrypted IS 'Encrypted last update timestamp';
COMMENT ON COLUMN cases.cases.dek_encrypted IS 'Encrypted Data Encryption Key for field-level encryption';
COMMENT ON COLUMN cases.cases.key_version IS 'Version of the Key Encryption Key used';
COMMENT ON COLUMN cases.cases.metadata IS 'Additional metadata in JSON format';

-- +goose Down
-- SQL in this section is executed when the migration is rolled back.

DROP TABLE IF EXISTS cases.cases;
DROP SCHEMA IF EXISTS cases;
