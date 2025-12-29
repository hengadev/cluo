-- +goose Up
-- SQL in this section is executed when the migration is applied.

-- Create case_subjects schema if it doesn't exist
CREATE SCHEMA IF NOT EXISTS cases;

-- Case subjects table (persons of interest in cases)
CREATE TABLE IF NOT EXISTS cases.case_subjects (
    id UUID PRIMARY KEY,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    lastname_encrypted BYTEA NOT NULL,
    lastname_hash VARCHAR(255) NOT NULL,
    firstname_encrypted BYTEA NOT NULL,
    firstname_hash VARCHAR(255) NOT NULL,
    email_encrypted BYTEA,
    email_hash VARCHAR(255),
    phone_encrypted BYTEA,
    city_encrypted BYTEA,
    city_hash VARCHAR(255),
    postal_code_encrypted BYTEA,
    postal_code_hash VARCHAR(255),
    address1_encrypted BYTEA,
    address1_hash VARCHAR(255),
    address2_encrypted BYTEA,
    address2_hash VARCHAR(255),
    occupation_encrypted BYTEA,
    occupation_hash VARCHAR(255),
    notes_encrypted BYTEA,
    dek_encrypted BYTEA NOT NULL,
    key_version INTEGER NOT NULL,
    metadata JSONB DEFAULT '{}',

    -- Business constraints
    CONSTRAINT chk_case_subjects_key_version_positive
        CHECK (key_version >= 0),
    CONSTRAINT chk_case_subjects_lastname_hash_not_empty
        CHECK (length(trim(lastname_hash)) > 0),
    CONSTRAINT chk_case_subjects_firstname_hash_not_empty
        CHECK (length(trim(firstname_hash)) > 0)
);

-- Case subject to case relationship table (many-to-many with roles)
CREATE TABLE IF NOT EXISTS cases.case_subject_cases (
    id UUID PRIMARY KEY,
    case_subject_id UUID NOT NULL,
    case_id UUID NOT NULL,
    roles TEXT[] NOT NULL DEFAULT '{}',
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),

    -- Foreign key constraints
    CONSTRAINT fk_case_subject_cases_subject_id
        FOREIGN KEY (case_subject_id) REFERENCES cases.case_subjects(id) ON DELETE CASCADE,
    -- Note: FK to cases.cases will be added after cases table is created

    -- Prevent duplicate subject-case pairs
    CONSTRAINT uq_case_subject_cases_subject_case
        UNIQUE (case_subject_id, case_id)
);

-- Indexes for case_subjects table
CREATE INDEX IF NOT EXISTS idx_case_subjects_lastname_hash ON cases.case_subjects(lastname_hash);
CREATE INDEX IF NOT EXISTS idx_case_subjects_firstname_hash ON cases.case_subjects(firstname_hash);
CREATE INDEX IF NOT EXISTS idx_case_subjects_email_hash ON cases.case_subjects(email_hash);
CREATE INDEX IF NOT EXISTS idx_case_subjects_city_hash ON cases.case_subjects(city_hash);
CREATE INDEX IF NOT EXISTS idx_case_subjects_postal_code_hash ON cases.case_subjects(postal_code_hash);
CREATE INDEX IF NOT EXISTS idx_case_subjects_address1_hash ON cases.case_subjects(address1_hash);
CREATE INDEX IF NOT EXISTS idx_case_subjects_occupation_hash ON cases.case_subjects(occupation_hash);
CREATE INDEX IF NOT EXISTS idx_case_subjects_key_version ON cases.case_subjects(key_version);
CREATE INDEX IF NOT EXISTS idx_case_subjects_metadata ON cases.case_subjects USING gin(metadata);
CREATE INDEX IF NOT EXISTS idx_case_subjects_created_at ON cases.case_subjects(created_at);

-- Indexes for case_subject_cases table
CREATE INDEX IF NOT EXISTS idx_case_subject_cases_subject_id ON cases.case_subject_cases(case_subject_id);
CREATE INDEX IF NOT EXISTS idx_case_subject_cases_case_id ON cases.case_subject_cases(case_id);
CREATE INDEX IF NOT EXISTS idx_case_subject_cases_roles ON cases.case_subject_cases USING gin(roles);

-- Add comments for documentation
COMMENT ON TABLE cases.case_subjects IS 'Stores encrypted information about persons of interest in cases (victims, suspects, witnesses, etc.)';
COMMENT ON COLUMN cases.case_subjects.id IS 'Unique identifier for the case subject';
COMMENT ON COLUMN cases.case_subjects.lastname_encrypted IS 'Encrypted last name';
COMMENT ON COLUMN cases.case_subjects.lastname_hash IS 'Hashed last name for indexing and search';
COMMENT ON COLUMN cases.case_subjects.firstname_encrypted IS 'Encrypted first name';
COMMENT ON COLUMN cases.case_subjects.firstname_hash IS 'Hashed first name for indexing and search';
COMMENT ON COLUMN cases.case_subjects.email_encrypted IS 'Encrypted email address';
COMMENT ON COLUMN cases.case_subjects.email_hash IS 'Hashed email for indexing and search';
COMMENT ON COLUMN cases.case_subjects.phone_encrypted IS 'Encrypted phone number';
COMMENT ON COLUMN cases.case_subjects.city_encrypted IS 'Encrypted city';
COMMENT ON COLUMN cases.case_subjects.city_hash IS 'Hashed city for indexing and search';
COMMENT ON COLUMN cases.case_subjects.postal_code_encrypted IS 'Encrypted postal code';
COMMENT ON COLUMN cases.case_subjects.postal_code_hash IS 'Hashed postal code for indexing and search';
COMMENT ON COLUMN cases.case_subjects.address1_encrypted IS 'Encrypted address line 1';
COMMENT ON COLUMN cases.case_subjects.address1_hash IS 'Hashed address line 1 for indexing and search';
COMMENT ON COLUMN cases.case_subjects.address2_encrypted IS 'Encrypted address line 2';
COMMENT ON COLUMN cases.case_subjects.address2_hash IS 'Hashed address line 2 for indexing and search';
COMMENT ON COLUMN cases.case_subjects.occupation_encrypted IS 'Encrypted occupation';
COMMENT ON COLUMN cases.case_subjects.occupation_hash IS 'Hashed occupation for indexing and search';
COMMENT ON COLUMN cases.case_subjects.notes_encrypted IS 'Encrypted notes about the subject';
COMMENT ON COLUMN cases.case_subjects.dek_encrypted IS 'Encrypted Data Encryption Key for field-level encryption';
COMMENT ON COLUMN cases.case_subjects.key_version IS 'Version of the Key Encryption Key used';

COMMENT ON TABLE cases.case_subject_cases IS 'Many-to-many relationship between case subjects and cases with their roles';
COMMENT ON COLUMN cases.case_subject_cases.case_subject_id IS 'Foreign key to cases.case_subjects.id';
COMMENT ON COLUMN cases.case_subject_cases.case_id IS 'Foreign key to cases.cases.id';
COMMENT ON COLUMN cases.case_subject_cases.roles IS 'Array of roles the subject has in this case (victim, suspect, witness, etc.)';

-- +goose Down
-- SQL in this section is executed when the migration is rolled back.

DROP TABLE IF EXISTS cases.case_subject_cases;
DROP TABLE IF EXISTS cases.case_subjects;
-- Note: We don't drop the schema here as it may be used by other tables
