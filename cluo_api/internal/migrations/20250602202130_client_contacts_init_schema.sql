-- +goose Up
-- SQL in this section is executed when the migration is applied.

-- Enable UUID extension if not already enabled
CREATE EXTENSION IF NOT EXISTS "uuid-ossp";

-- Create clients schema if it doesn't exist
CREATE SCHEMA IF NOT EXISTS clients;

-- Contacts table for storing encrypted client contact information
CREATE TABLE IF NOT EXISTS clients.contacts (
    id UUID PRIMARY KEY,
    client_id_hash VARCHAR(255) NOT NULL,
    client_id_encrypted BYTEA NOT NULL,
    lastname_encrypted BYTEA NOT NULL,
    firstname_encrypted BYTEA NOT NULL,
    email_hash VARCHAR(255) NOT NULL,
    email_encrypted BYTEA NOT NULL,
    phone_encrypted BYTEA,
    position_encrypted BYTEA,
    dek_encrypted BYTEA NOT NULL,
    key_version INTEGER NOT NULL,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    metadata JSONB DEFAULT '{}'
);

-- Indexes for performance and search
CREATE INDEX IF NOT EXISTS idx_contacts_client_id_hash ON clients.contacts(client_id_hash);
CREATE INDEX IF NOT EXISTS idx_contacts_email_hash ON clients.contacts(email_hash);
CREATE INDEX IF NOT EXISTS idx_contacts_created_at ON clients.contacts(created_at);
CREATE INDEX IF NOT EXISTS idx_contacts_updated_at ON clients.contacts(updated_at);
CREATE INDEX IF NOT EXISTS idx_contacts_key_version ON clients.contacts(key_version);

-- GIN index on metadata for efficient JSON queries
CREATE INDEX IF NOT EXISTS idx_contacts_metadata ON clients.contacts USING gin(metadata);

-- Note: Trigger functions for updated_at are commented out due to migration system limitations
-- To be implemented manually or with a different migration approach
-- CREATE OR REPLACE FUNCTION update_updated_at_column()
-- RETURNS TRIGGER AS $func$
-- BEGIN
--     NEW.updated_at = NOW();
--     RETURN NEW;
-- END;
-- $func$ LANGUAGE plpgsql;

-- Business rule constraints
ALTER TABLE clients.contacts ADD CONSTRAINT chk_contacts_key_version_positive
    CHECK (key_version >= 0);

-- Add comments for documentation
COMMENT ON SCHEMA clients IS 'Client management schema';
COMMENT ON TABLE clients.contacts IS 'Stores encrypted contact information for clients';
COMMENT ON COLUMN clients.contacts.id IS 'Unique identifier for the contact record';
COMMENT ON COLUMN clients.contacts.client_id_hash IS 'Hashed client identifier for indexing and searching';
COMMENT ON COLUMN clients.contacts.client_id_encrypted IS 'Encrypted client UUID reference';
COMMENT ON COLUMN clients.contacts.lastname_encrypted IS 'Encrypted contact last name';
COMMENT ON COLUMN clients.contacts.firstname_encrypted IS 'Encrypted contact first name';
COMMENT ON COLUMN clients.contacts.email_hash IS 'Hashed email for indexing and searching';
COMMENT ON COLUMN clients.contacts.email_encrypted IS 'Encrypted email address';
COMMENT ON COLUMN clients.contacts.phone_encrypted IS 'Encrypted phone number (optional)';
COMMENT ON COLUMN clients.contacts.position_encrypted IS 'Encrypted job position/title (optional)';
COMMENT ON COLUMN clients.contacts.dek_encrypted IS 'Encrypted Data Encryption Key for field-level encryption';
COMMENT ON COLUMN clients.contacts.key_version IS 'Version of the Key Encryption Key used';
COMMENT ON COLUMN clients.contacts.metadata IS 'Additional metadata in JSON format';

-- +goose Down
-- SQL in this section is executed when the migration is rolled back.

DROP TABLE IF EXISTS clients.contacts;
DROP SCHEMA IF EXISTS clients;
