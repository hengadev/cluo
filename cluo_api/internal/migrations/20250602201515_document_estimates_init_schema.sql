-- +goose Up
-- SQL in this section is executed when the migration is applied.

-- Enable UUID extension if not already enabled
CREATE EXTENSION IF NOT EXISTS "uuid-ossp";

-- Estimates table
CREATE TABLE IF NOT EXISTS estimates (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    case_id UUID NOT NULL,
    client_id UUID NOT NULL,
    status VARCHAR(20) NOT NULL DEFAULT 'draft' CHECK (status IN ('draft', 'sent', 'signed', 'active', 'archived', 'cancelled', 'rejected', 'expired')),
    estimate_number VARCHAR(50) NOT NULL UNIQUE,
    issue_date TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    valid_until TIMESTAMPTZ,
    line_items JSONB NOT NULL DEFAULT '[]',
    estimated_total NUMERIC(12,2) NOT NULL DEFAULT 0.00,
    notes TEXT,
    accepted BOOLEAN NOT NULL DEFAULT false,
    accepted_at TIMESTAMPTZ,
    accepted_by UUID,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

-- Indexes for performance
CREATE INDEX IF NOT EXISTS idx_estimates_case_id ON estimates(case_id);
CREATE INDEX IF NOT EXISTS idx_estimates_client_id ON estimates(client_id);
CREATE INDEX IF NOT EXISTS idx_estimates_status ON estimates(status);
CREATE INDEX IF NOT EXISTS idx_estimates_estimate_number ON estimates(estimate_number);
CREATE INDEX IF NOT EXISTS idx_estimates_issue_date ON estimates(issue_date);
CREATE INDEX IF NOT EXISTS idx_estimates_valid_until ON estimates(valid_until) WHERE valid_until IS NOT NULL;

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
ALTER TABLE estimates ADD CONSTRAINT chk_estimates_positive_total
    CHECK (estimated_total >= 0);

ALTER TABLE estimates ADD CONSTRAINT chk_estimates_acceptance_logic
    CHECK (
        (accepted = false AND accepted_at IS NULL AND accepted_by IS NULL) OR
        (accepted = true AND accepted_at IS NOT NULL AND accepted_by IS NOT NULL)
    );

-- Add comments for documentation
COMMENT ON TABLE estimates IS 'Stores price quotations for investigative services';
COMMENT ON COLUMN estimates.estimate_number IS 'Unique identifier following format EST-YYYY-NNN';
COMMENT ON COLUMN estimates.line_items IS 'JSON array of estimate line items';
COMMENT ON COLUMN estimates.accepted IS 'Whether the estimate has been accepted by the client';

-- +goose Down
-- SQL in this section is executed when the migration is rolled back.

DROP TABLE IF EXISTS estimates;
