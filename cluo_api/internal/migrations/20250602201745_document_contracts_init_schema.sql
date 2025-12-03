-- +goose Up
-- +goose StatementBegin

-- Create uuid extension if not exists
CREATE EXTENSION IF NOT EXISTS "uuid-ossp";

-- Create clients schema
CREATE SCHEMA IF NOT EXISTS clients;

-- Clients table
CREATE TABLE IF NOT EXISTS clients.clients (
    id UUID PRIMARY KEY,
    name_encrypted BYTEA NOT NULL,
    name_hash VARCHAR(255) NOT NULL,
    type_encrypted BYTEA NOT NULL,
    type_hash VARCHAR(255) NOT NULL,
    contactids_encrypted BYTEA NOT NULL,
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

-- Create update trigger function
CREATE OR REPLACE FUNCTION update_updated_at_column()
RETURNS TRIGGER AS $$
BEGIN
    NEW.updated_at = NOW();
    RETURN NEW;
END;
$$ language 'plpgsql';

-- Create trigger for clients
DROP TRIGGER IF EXISTS update_clients_updated_at ON clients.clients;
CREATE TRIGGER update_clients_updated_at BEFORE UPDATE ON clients.clients
    FOR EACH ROW EXECUTE FUNCTION update_updated_at_column();

-- Add comments for clients documentation
COMMENT ON SCHEMA clients IS 'Client management schema';
COMMENT ON TABLE clients.clients IS 'Stores encrypted client information';
COMMENT ON COLUMN clients.clients.id IS 'Unique identifier for the client';
COMMENT ON COLUMN clients.clients.created_at IS 'Client creation timestamp';
COMMENT ON COLUMN clients.clients.updated_at IS 'Client last update timestamp';
COMMENT ON COLUMN clients.clients.name_encrypted IS 'Encrypted client name';
COMMENT ON COLUMN clients.clients.name_hash IS 'Hashed client name for indexing and search';
COMMENT ON COLUMN clients.clients.type_encrypted IS 'Encrypted client type';
COMMENT ON COLUMN clients.clients.type_hash IS 'Hashed client type for indexing and search';
COMMENT ON COLUMN clients.clients.contactids_encrypted IS 'Encrypted list of contact UUIDs';
COMMENT ON COLUMN clients.clients.dek_encrypted IS 'Encrypted Data Encryption Key for field-level encryption';
COMMENT ON COLUMN clients.clients.key_version IS 'Version of the Key Encryption Key used';
COMMENT ON COLUMN clients.clients.metadata IS 'Additional metadata in JSON format';

-- Contracts table
CREATE TABLE IF NOT EXISTS contracts (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    case_id UUID NOT NULL,
    client_id UUID NOT NULL,
    status VARCHAR(20) NOT NULL DEFAULT 'draft' CHECK (status IN ('draft', 'sent', 'signed', 'active', 'archived', 'cancelled', 'rejected', 'expired')),
    contract_number VARCHAR(50) NOT NULL UNIQUE,
    start_date TIMESTAMPTZ NOT NULL,
    end_date TIMESTAMPTZ,
    scope_of_services TEXT NOT NULL,
    payment_terms TEXT NOT NULL,
    confidentiality TEXT NOT NULL,
    termination_clause TEXT NOT NULL,
    signatures JSONB NOT NULL DEFAULT '[]',
    linked_mandate_id UUID REFERENCES mandates(id) ON DELETE SET NULL,
    contract_value NUMERIC(12,2),
    currency VARCHAR(3) DEFAULT 'USD',
    renewal_terms TEXT,
    governing_law VARCHAR(100),
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

-- Indexes for performance
CREATE INDEX IF NOT EXISTS idx_contracts_case_id ON contracts(case_id);
CREATE INDEX IF NOT EXISTS idx_contracts_client_id ON contracts(client_id);
CREATE INDEX IF NOT EXISTS idx_contracts_status ON contracts(status);
CREATE INDEX IF NOT EXISTS idx_contracts_contract_number ON contracts(contract_number);
CREATE INDEX IF NOT EXISTS idx_contracts_linked_mandate_id ON contracts(linked_mandate_id);
CREATE INDEX IF NOT EXISTS idx_contracts_start_date ON contracts(start_date);
CREATE INDEX IF NOT EXISTS idx_contracts_end_date ON contracts(end_date) WHERE end_date IS NOT NULL;

-- Add foreign key constraint for contracts.client_id
ALTER TABLE contracts ADD CONSTRAINT fk_contracts_client_id
    FOREIGN KEY (client_id) REFERENCES clients.clients(id) ON DELETE RESTRICT;

-- Create trigger for contracts
DROP TRIGGER IF EXISTS update_contracts_updated_at ON contracts;
CREATE TRIGGER update_contracts_updated_at BEFORE UPDATE ON contracts
    FOR EACH ROW EXECUTE FUNCTION update_updated_at_column();

-- Business rule constraints
ALTER TABLE contracts ADD CONSTRAINT chk_contracts_date_logic
    CHECK (
        end_date IS NULL OR
        end_date > start_date
    );

ALTER TABLE contracts ADD CONSTRAINT chk_contracts_positive_value
    CHECK (
        contract_value IS NULL OR
        contract_value >= 0
    );

-- Add comments for documentation
COMMENT ON TABLE contracts IS 'Stores formal agreements between parties';
COMMENT ON COLUMN contracts.contract_number IS 'Unique identifier following format CNT-YYYY-NNN';
COMMENT ON COLUMN contracts.signatures IS 'JSON array of all signatures on the contract';
COMMENT ON COLUMN contracts.linked_mandate_id IS 'Optional reference to the originating mandate';

-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin

DROP TRIGGER IF EXISTS update_contracts_updated_at ON contracts;
DROP TRIGGER IF EXISTS update_clients_updated_at ON clients.clients;
DROP TABLE IF EXISTS contracts;
DROP TABLE IF EXISTS clients.clients;
DROP SCHEMA IF EXISTS clients;

-- +goose StatementEnd
