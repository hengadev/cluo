-- +goose Up
-- SQL in this section is executed when the migration is applied.

-- Enable UUID extension if not already enabled
CREATE EXTENSION IF NOT EXISTS "uuid-ossp";

-- Create clients schema if it doesn't exist
CREATE SCHEMA IF NOT EXISTS clients;

-- Clients table
CREATE TABLE IF NOT EXISTS clients.clients (
    id UUID PRIMARY KEY,
    name_encrypted BYTEA NOT NULL,
    name_hash VARCHAR(255) NOT NULL,
    type_encrypted BYTEA NOT NULL,
    type_hash VARCHAR(255) NOT NULL,
    dek_encrypted BYTEA NOT NULL,
    key_version INTEGER NOT NULL,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    metadata JSONB DEFAULT '{}'
);

-- Clients table indexes
CREATE INDEX IF NOT EXISTS idx_clients_name_hash ON clients.clients(name_hash);
CREATE INDEX IF NOT EXISTS idx_clients_type_hash ON clients.clients(type_hash);
CREATE INDEX IF NOT EXISTS idx_clients_created_at ON clients.clients(created_at);
CREATE INDEX IF NOT EXISTS idx_clients_updated_at ON clients.clients(updated_at);
CREATE INDEX IF NOT EXISTS idx_clients_key_version ON clients.clients(key_version);
CREATE INDEX IF NOT EXISTS idx_clients_metadata ON clients.clients USING gin(metadata);

-- Clients table business constraints
ALTER TABLE clients.clients ADD CONSTRAINT chk_clients_key_version_positive
    CHECK (key_version >= 0);

ALTER TABLE clients.clients ADD CONSTRAINT chk_clients_name_hash_not_empty
    CHECK (length(trim(name_hash)) > 0);

ALTER TABLE clients.clients ADD CONSTRAINT chk_clients_type_hash_not_empty
    CHECK (length(trim(type_hash)) > 0);


-- Create trigger for clients
DROP TRIGGER IF EXISTS update_clients_updated_at ON clients.clients;
CREATE TRIGGER update_clients_updated_at BEFORE UPDATE ON clients.clients
    FOR EACH ROW EXECUTE FUNCTION update_updated_at_column();

-- Contacts table for storing encrypted client contact information
CREATE TABLE IF NOT EXISTS clients.contacts (
    id UUID PRIMARY KEY,
    client_id UUID NOT NULL REFERENCES clients.clients(id) ON DELETE CASCADE,
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
CREATE INDEX IF NOT EXISTS idx_contacts_client_id ON clients.contacts(client_id);
CREATE INDEX IF NOT EXISTS idx_contacts_email_hash ON clients.contacts(email_hash);
CREATE INDEX IF NOT EXISTS idx_contacts_created_at ON clients.contacts(created_at);
CREATE INDEX IF NOT EXISTS idx_contacts_updated_at ON clients.contacts(updated_at);
CREATE INDEX IF NOT EXISTS idx_contacts_key_version ON clients.contacts(key_version);

-- GIN index on metadata for efficient JSON queries
CREATE INDEX IF NOT EXISTS idx_contacts_metadata ON clients.contacts USING gin(metadata);

-- Create trigger for contacts
DROP TRIGGER IF EXISTS update_contacts_updated_at ON clients.contacts;
CREATE TRIGGER update_contacts_updated_at BEFORE UPDATE ON clients.contacts
    FOR EACH ROW EXECUTE FUNCTION update_updated_at_column();

-- Business rule constraints
ALTER TABLE clients.contacts ADD CONSTRAINT chk_contacts_key_version_positive
    CHECK (key_version >= 0);

-- Add comments for documentation
COMMENT ON SCHEMA clients IS 'Client management schema';
COMMENT ON TABLE clients.clients IS 'Stores encrypted client information';
COMMENT ON COLUMN clients.clients.id IS 'Unique identifier for the client';
COMMENT ON COLUMN clients.clients.created_at IS 'Client creation timestamp';
COMMENT ON COLUMN clients.clients.updated_at IS 'Client last update timestamp';
COMMENT ON COLUMN clients.clients.name_encrypted IS 'Encrypted client name';
COMMENT ON COLUMN clients.clients.name_hash IS 'Hashed client name for indexing and search';
COMMENT ON COLUMN clients.clients.type_encrypted IS 'Encrypted client type';
COMMENT ON COLUMN clients.clients.type_hash IS 'Hashed client type for indexing and search';
COMMENT ON COLUMN clients.clients.dek_encrypted IS 'Encrypted Data Encryption Key for field-level encryption';
COMMENT ON COLUMN clients.clients.key_version IS 'Version of the Key Encryption Key used';
COMMENT ON COLUMN clients.clients.metadata IS 'Additional metadata in JSON format';

COMMENT ON TABLE clients.contacts IS 'Stores encrypted contact information for clients';
COMMENT ON COLUMN clients.contacts.id IS 'Unique identifier for the contact record';
COMMENT ON COLUMN clients.contacts.client_id IS 'Foreign key reference to the clients table';
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

-- Add back contactids_encrypted column for rollback
ALTER TABLE clients.clients ADD COLUMN IF NOT EXISTS contactids_encrypted BYTEA NOT NULL DEFAULT '';

DROP TRIGGER IF EXISTS update_contacts_updated_at ON clients.contacts;
DROP TRIGGER IF EXISTS update_clients_updated_at ON clients.clients;
DROP TABLE IF EXISTS clients.contacts;
DROP TABLE IF EXISTS clients.clients;
DROP SCHEMA IF EXISTS clients;
