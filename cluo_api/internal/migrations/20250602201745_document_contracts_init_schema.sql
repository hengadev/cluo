-- +goose Up
-- +goose StatementBegin

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
DROP TABLE IF EXISTS contracts;

-- +goose StatementEnd
