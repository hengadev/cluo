-- +goose Up
-- SQL in this section is executed when the migration is applied.

-- Enable UUID extension if not already enabled
CREATE EXTENSION IF NOT EXISTS "uuid-ossp";

-- Create auth schema if it doesn't exist
CREATE SCHEMA IF NOT EXISTS auth;

-- Users table for storing encrypted user authentication information
CREATE TABLE IF NOT EXISTS auth.users (
    id UUID PRIMARY KEY,
    email_hash VARCHAR(255) NOT NULL,
    email_encrypted BYTEA,
    password_hash_secure BYTEA NOT NULL,
    role_encrypted BYTEA NOT NULL,
    created_at_encrypted BYTEA,
    dek_encrypted BYTEA NOT NULL,
    key_version INTEGER NOT NULL,
    metadata JSONB DEFAULT '{}'
);

-- Indexes for performance and search
CREATE UNIQUE INDEX IF NOT EXISTS idx_users_email_hash ON auth.users(email_hash);
CREATE INDEX IF NOT EXISTS idx_users_key_version ON auth.users(key_version);
CREATE INDEX IF NOT EXISTS idx_users_id ON auth.users(id);

-- GIN index on metadata for efficient JSON queries
CREATE INDEX IF NOT EXISTS idx_users_metadata ON auth.users USING gin(metadata);

-- Business rule constraints
ALTER TABLE auth.users ADD CONSTRAINT chk_users_key_version_positive
    CHECK (key_version >= 0);

ALTER TABLE auth.users ADD CONSTRAINT chk_users_email_hash_not_empty
    CHECK (length(trim(email_hash)) > 0);

-- Add comments for documentation
COMMENT ON SCHEMA auth IS 'Authentication and authorization schema';
COMMENT ON TABLE auth.users IS 'Stores encrypted user authentication and authorization data';
COMMENT ON COLUMN auth.users.id IS 'Unique identifier for the user';
COMMENT ON COLUMN auth.users.email_hash IS 'Hashed email for indexing and authentication lookup';
COMMENT ON COLUMN auth.users.email_encrypted IS 'Encrypted email address for privacy';
COMMENT ON COLUMN auth.users.password_hash_secure IS 'Secure password hash using PBKDF2 or similar';
COMMENT ON COLUMN auth.users.role_encrypted IS 'Encrypted user role (Guest, Client, Administrator)';
COMMENT ON COLUMN auth.users.created_at_encrypted IS 'Encrypted account creation timestamp';
COMMENT ON COLUMN auth.users.dek_encrypted IS 'Encrypted Data Encryption Key for field-level encryption';
COMMENT ON COLUMN auth.users.key_version IS 'Version of the Key Encryption Key used';
COMMENT ON COLUMN auth.users.metadata IS 'Additional metadata in JSON format';

-- +goose Down
-- SQL in this section is executed when the migration is rolled back.

DROP TABLE IF EXISTS auth.users;
DROP SCHEMA IF EXISTS auth;