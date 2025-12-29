-- +goose Up
-- SQL in this section is executed when the migration is applied.

-- Create cases schema if it doesn't exist
CREATE SCHEMA IF NOT EXISTS cases;

-- Cases table
CREATE TABLE IF NOT EXISTS cases.cases (
    id UUID PRIMARY KEY,
    client_id UUID NOT NULL,
    assigned_contact_id UUID,
    case_subject_id UUID,
    case_type VARCHAR(100) NOT NULL DEFAULT '',
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    title_encrypted BYTEA NOT NULL,
    description_encrypted BYTEA NOT NULL,
    external_reference_encrypted BYTEA,
    external_reference_hash VARCHAR(255),
    status_encrypted BYTEA NOT NULL,
    placename_encrypted BYTEA,
    placename_hash VARCHAR(255),
    address1_encrypted BYTEA,
    address1_hash VARCHAR(255),
    address2_encrypted BYTEA,
    address2_hash VARCHAR(255),
    city_encrypted BYTEA,
    city_hash VARCHAR(255),
    postal_code_encrypted BYTEA,
    postal_code_hash VARCHAR(255),
    country_encrypted BYTEA,
    country_hash VARCHAR(255),
    latitude_encrypted BYTEA,
    latitude_hash VARCHAR(255),
    longitude_encrypted BYTEA,
    longitude_hash VARCHAR(255),
    location_type_encrypted BYTEA,
    location_type_hash VARCHAR(255),
    location_notes_encrypted BYTEA,
    location_notes_hash VARCHAR(255),
    updated_at_encrypted BYTEA,
    dek_encrypted BYTEA NOT NULL,
    key_version INTEGER NOT NULL,
    metadata JSONB DEFAULT '{}',

    -- Foreign key constraints
    CONSTRAINT fk_cases_client_id
        FOREIGN KEY (client_id) REFERENCES clients.clients(id) ON DELETE RESTRICT,
    CONSTRAINT fk_cases_assigned_contact_id
        FOREIGN KEY (assigned_contact_id) REFERENCES clients.contacts(id) ON DELETE SET NULL,
    CONSTRAINT fk_cases_case_subject_id
        FOREIGN KEY (case_subject_id) REFERENCES cases.case_subjects(id) ON DELETE SET NULL
);

-- Cases table indexes
CREATE INDEX IF NOT EXISTS idx_cases_client_id ON cases.cases(client_id);
CREATE INDEX IF NOT EXISTS idx_cases_assigned_contact_id ON cases.cases(assigned_contact_id);
CREATE INDEX IF NOT EXISTS idx_cases_case_subject_id ON cases.cases(case_subject_id);
CREATE INDEX IF NOT EXISTS idx_cases_case_type ON cases.cases(case_type);
CREATE INDEX IF NOT EXISTS idx_cases_external_reference_hash ON cases.cases(external_reference_hash);
CREATE INDEX IF NOT EXISTS idx_cases_placename_hash ON cases.cases(placename_hash);
CREATE INDEX IF NOT EXISTS idx_cases_address1_hash ON cases.cases(address1_hash);
CREATE INDEX IF NOT EXISTS idx_cases_city_hash ON cases.cases(city_hash);
CREATE INDEX IF NOT EXISTS idx_cases_postal_code_hash ON cases.cases(postal_code_hash);
CREATE INDEX IF NOT EXISTS idx_cases_country_hash ON cases.cases(country_hash);
CREATE INDEX IF NOT EXISTS idx_cases_location_type_hash ON cases.cases(location_type_hash);
CREATE INDEX IF NOT EXISTS idx_cases_created_at ON cases.cases(created_at);
CREATE INDEX IF NOT EXISTS idx_cases_key_version ON cases.cases(key_version);
CREATE INDEX IF NOT EXISTS idx_cases_metadata ON cases.cases USING gin(metadata);

-- Add foreign key constraint from case_subject_cases to cases
ALTER TABLE cases.case_subject_cases ADD CONSTRAINT fk_case_subject_cases_case_id
    FOREIGN KEY (case_id) REFERENCES cases.cases(id) ON DELETE CASCADE;

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
COMMENT ON COLUMN cases.cases.case_subject_id IS 'Foreign key reference to cases.case_subjects.id (nullable) - primary subject/person of interest';
COMMENT ON COLUMN cases.cases.case_type IS 'Type/category of the case (e.g., theft, accident, fraud)';
COMMENT ON COLUMN cases.cases.created_at IS 'Case creation timestamp';
COMMENT ON COLUMN cases.cases.title_encrypted IS 'Encrypted case title';
COMMENT ON COLUMN cases.cases.description_encrypted IS 'Encrypted case description';
COMMENT ON COLUMN cases.cases.external_reference_encrypted IS 'Encrypted external reference from other companies (e.g., insurance company client ID)';
COMMENT ON COLUMN cases.cases.external_reference_hash IS 'Hashed external reference for indexing and search';
COMMENT ON COLUMN cases.cases.status_encrypted IS 'Encrypted case status';
COMMENT ON COLUMN cases.cases.placename_encrypted IS 'Encrypted place name (e.g., Gare de Lyon)';
COMMENT ON COLUMN cases.cases.placename_hash IS 'Hashed place name for indexing and search';
COMMENT ON COLUMN cases.cases.address1_encrypted IS 'Encrypted address line 1';
COMMENT ON COLUMN cases.cases.address1_hash IS 'Hashed address line 1 for indexing and search';
COMMENT ON COLUMN cases.cases.address2_encrypted IS 'Encrypted address line 2';
COMMENT ON COLUMN cases.cases.address2_hash IS 'Hashed address line 2 for indexing and search';
COMMENT ON COLUMN cases.cases.city_encrypted IS 'Encrypted city name';
COMMENT ON COLUMN cases.cases.city_hash IS 'Hashed city name for indexing and search';
COMMENT ON COLUMN cases.cases.postal_code_encrypted IS 'Encrypted postal/zip code';
COMMENT ON COLUMN cases.cases.postal_code_hash IS 'Hashed postal code for indexing and search';
COMMENT ON COLUMN cases.cases.country_encrypted IS 'Encrypted country name';
COMMENT ON COLUMN cases.cases.country_hash IS 'Hashed country name for indexing and search';
COMMENT ON COLUMN cases.cases.latitude_encrypted IS 'Encrypted latitude coordinate';
COMMENT ON COLUMN cases.cases.latitude_hash IS 'Hashed latitude for indexing and search';
COMMENT ON COLUMN cases.cases.longitude_encrypted IS 'Encrypted longitude coordinate';
COMMENT ON COLUMN cases.cases.longitude_hash IS 'Hashed longitude for indexing and search';
COMMENT ON COLUMN cases.cases.location_type_encrypted IS 'Encrypted location type (e.g., residence, business, public space)';
COMMENT ON COLUMN cases.cases.location_type_hash IS 'Hashed location type for indexing and search';
COMMENT ON COLUMN cases.cases.location_notes_encrypted IS 'Encrypted additional notes about the location';
COMMENT ON COLUMN cases.cases.location_notes_hash IS 'Hashed location notes for indexing and search';
COMMENT ON COLUMN cases.cases.updated_at_encrypted IS 'Encrypted last update timestamp';
COMMENT ON COLUMN cases.cases.dek_encrypted IS 'Encrypted Data Encryption Key for field-level encryption';
COMMENT ON COLUMN cases.cases.key_version IS 'Version of the Key Encryption Key used';
COMMENT ON COLUMN cases.cases.metadata IS 'Additional metadata in JSON format';

-- +goose Down
-- SQL in this section is executed when the migration is rolled back.

DROP TABLE IF EXISTS cases.cases;
DROP SCHEMA IF EXISTS cases;
