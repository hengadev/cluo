-- +goose Up
-- SQL in this section is executed when the migration is applied.

-- Invoices table
CREATE TABLE IF NOT EXISTS invoices (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    case_id UUID NOT NULL,
    client_id UUID NOT NULL,
    status VARCHAR(20) NOT NULL DEFAULT 'draft' CHECK (status IN ('draft', 'sent', 'signed', 'active', 'archived', 'cancelled', 'rejected', 'expired')),
    invoice_number VARCHAR(50) NOT NULL UNIQUE,
    issue_date TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    due_date TIMESTAMPTZ NOT NULL,
    line_items JSONB NOT NULL DEFAULT '[]',
    total_amount NUMERIC(12,2) NOT NULL DEFAULT 0.00,
    tax_rate NUMERIC(5,2) NOT NULL DEFAULT 0.00,
    tax_amount NUMERIC(12,2) NOT NULL DEFAULT 0.00,
    notes TEXT,
    payment_status VARCHAR(20) NOT NULL DEFAULT 'unpaid' CHECK (payment_status IN ('unpaid', 'paid', 'partially_paid', 'overdue', 'refunded', 'void')),
    paid_at TIMESTAMPTZ,
    paid_amount NUMERIC(12,2),
    payment_method VARCHAR(50),
    linked_contract_id UUID REFERENCES contracts(id) ON DELETE SET NULL,
    currency VARCHAR(3) DEFAULT 'USD',
    payment_terms TEXT,
    late_fee NUMERIC(12,2),
    late_fee_rate NUMERIC(5,2),
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

-- Indexes for performance
CREATE INDEX IF NOT EXISTS idx_invoices_case_id ON invoices(case_id);
CREATE INDEX IF NOT EXISTS idx_invoices_client_id ON invoices(client_id);
CREATE INDEX IF NOT EXISTS idx_invoices_status ON invoices(status);
CREATE INDEX IF NOT EXISTS idx_invoices_invoice_number ON invoices(invoice_number);
CREATE INDEX IF NOT EXISTS idx_invoices_payment_status ON invoices(payment_status);
CREATE INDEX IF NOT EXISTS idx_invoices_linked_contract_id ON invoices(linked_contract_id);
CREATE INDEX IF NOT EXISTS idx_invoices_issue_date ON invoices(issue_date);
CREATE INDEX IF NOT EXISTS idx_invoices_due_date ON invoices(due_date);

-- Create trigger for invoices
DROP TRIGGER IF EXISTS update_invoices_updated_at ON invoices;
CREATE TRIGGER update_invoices_updated_at BEFORE UPDATE ON invoices
    FOR EACH ROW EXECUTE FUNCTION update_updated_at_column();

-- Business rule constraints
ALTER TABLE invoices ADD CONSTRAINT IF NOT EXISTS chk_invoices_date_logic
    CHECK (due_date >= issue_date);

ALTER TABLE invoices ADD CONSTRAINT IF NOT EXISTS chk_invoices_positive_amounts
    CHECK (
        total_amount >= 0 AND
        tax_amount >= 0 AND
        tax_rate >= 0
    );

ALTER TABLE invoices ADD CONSTRAINT IF NOT EXISTS chk_invoices_payment_logic
    CHECK (
        (paid_at IS NULL AND paid_amount IS NULL) OR
        (paid_at IS NOT NULL AND paid_amount IS NOT NULL AND paid_amount >= 0)
    );

ALTER TABLE invoices ADD CONSTRAINT IF NOT EXISTS chk_invoices_late_fee_logic
    CHECK (
        (late_fee IS NULL AND late_fee_rate IS NULL) OR
        (late_fee >= 0 AND late_fee_rate IS NULL) OR
        (late_fee IS NULL AND late_fee_rate >= 0) OR
        (late_fee >= 0 AND late_fee_rate >= 0)
    );

-- Add comments for documentation
COMMENT ON TABLE invoices IS 'Stores billing documents for services rendered';
COMMENT ON COLUMN invoices.invoice_number IS 'Unique identifier following format INV-YYYY-NNN';
COMMENT ON COLUMN invoices.line_items IS 'JSON array of invoice line items';
COMMENT ON COLUMN invoices.payment_status IS 'Current payment status of the invoice';
COMMENT ON COLUMN invoices.linked_contract_id IS 'Optional reference to the originating contract';

-- +goose Down
-- SQL in this section is executed when the migration is rolled back.

DROP TRIGGER IF EXISTS update_invoices_updated_at ON invoices;
DROP TABLE IF EXISTS invoices;