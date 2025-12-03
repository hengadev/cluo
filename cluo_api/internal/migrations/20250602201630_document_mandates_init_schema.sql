-- +goose Up
-- +goose StatementBegin

-- Mandates table
CREATE TABLE IF NOT EXISTS mandates (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    case_id UUID NOT NULL,
    client_id UUID NOT NULL,
    status VARCHAR(20) NOT NULL DEFAULT 'draft' CHECK (status IN ('draft', 'sent', 'signed', 'active', 'archived', 'cancelled', 'rejected', 'expired')),
    mandate_number VARCHAR(50) NOT NULL UNIQUE,
    issue_date TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    scope_of_work TEXT NOT NULL,
    valid_from TIMESTAMPTZ NOT NULL,
    valid_until TIMESTAMPTZ,
    terms_conditions TEXT NOT NULL,
    client_signature JSONB,
    investigator_signature JSONB,
    linked_estimate_id UUID REFERENCES estimates(id) ON DELETE SET NULL,
    special_instructions TEXT,
    jurisdiction VARCHAR(100),
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

-- Indexes for performance
CREATE INDEX IF NOT EXISTS idx_mandates_case_id ON mandates(case_id);
CREATE INDEX IF NOT EXISTS idx_mandates_client_id ON mandates(client_id);
CREATE INDEX IF NOT EXISTS idx_mandates_status ON mandates(status);
CREATE INDEX IF NOT EXISTS idx_mandates_mandate_number ON mandates(mandate_number);
CREATE INDEX IF NOT EXISTS idx_mandates_linked_estimate_id ON mandates(linked_estimate_id);
CREATE INDEX IF NOT EXISTS idx_mandates_valid_from ON mandates(valid_from);
CREATE INDEX IF NOT EXISTS idx_mandates_valid_until ON mandates(valid_until) WHERE valid_until IS NOT NULL;

-- Create trigger for mandates
DROP TRIGGER IF EXISTS update_mandates_updated_at ON mandates;
CREATE TRIGGER update_mandates_updated_at BEFORE UPDATE ON mandates
    FOR EACH ROW EXECUTE FUNCTION update_updated_at_column();

-- Business rule constraints
ALTER TABLE mandates ADD CONSTRAINT chk_mandates_date_logic
    CHECK (valid_from >= issue_date);

ALTER TABLE mandates ADD CONSTRAINT chk_mandates_valid_until_logic
    CHECK (
        valid_until IS NULL OR
        valid_until > valid_from
    );

-- Add comments for documentation
COMMENT ON TABLE mandates IS 'Stores legal authorization documents for investigations';
COMMENT ON COLUMN mandates.mandate_number IS 'Unique identifier following format MND-YYYY-NNN';
COMMENT ON COLUMN mandates.client_signature IS 'JSON object containing client signature details';
COMMENT ON COLUMN mandates.investigator_signature IS 'JSON object containing investigator signature details';
COMMENT ON COLUMN mandates.linked_estimate_id IS 'Optional reference to the originating estimate';

-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin

DROP TRIGGER IF EXISTS update_mandates_updated_at ON mandates;
DROP TABLE IF EXISTS mandates;
DROP FUNCTION IF EXISTS update_updated_at_column();

-- +goose StatementEnd
