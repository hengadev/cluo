-- +goose Up
-- Add encx support to document tables: add _encrypted columns and make original
-- columns nullable so that encx inserts can populate only the encrypted variants.

-- Estimates
ALTER TABLE estimates ALTER COLUMN case_id DROP NOT NULL;
ALTER TABLE estimates ALTER COLUMN client_id DROP NOT NULL;
ALTER TABLE estimates ALTER COLUMN estimate_number DROP NOT NULL;
ALTER TABLE estimates ALTER COLUMN line_items DROP NOT NULL;
ALTER TABLE estimates ALTER COLUMN estimated_total DROP NOT NULL;
ALTER TABLE estimates ALTER COLUMN notes DROP NOT NULL;

ALTER TABLE estimates ADD COLUMN IF NOT EXISTS caseid_encrypted BYTEA;
ALTER TABLE estimates ADD COLUMN IF NOT EXISTS clientid_encrypted BYTEA;
ALTER TABLE estimates ADD COLUMN IF NOT EXISTS estimatenumber_encrypted BYTEA;
ALTER TABLE estimates ADD COLUMN IF NOT EXISTS lineitems_encrypted BYTEA;
ALTER TABLE estimates ADD COLUMN IF NOT EXISTS estimatedtotal_encrypted BYTEA;
ALTER TABLE estimates ADD COLUMN IF NOT EXISTS notes_encrypted BYTEA;
ALTER TABLE estimates ADD COLUMN IF NOT EXISTS dek_encrypted BYTEA;
ALTER TABLE estimates ADD COLUMN IF NOT EXISTS key_version INT NOT NULL DEFAULT 1;
ALTER TABLE estimates ADD COLUMN IF NOT EXISTS metadata JSONB;

-- Mandates
ALTER TABLE mandates ALTER COLUMN case_id DROP NOT NULL;
ALTER TABLE mandates ALTER COLUMN client_id DROP NOT NULL;
ALTER TABLE mandates ALTER COLUMN mandate_number DROP NOT NULL;
ALTER TABLE mandates ALTER COLUMN scope_of_work DROP NOT NULL;
ALTER TABLE mandates ALTER COLUMN terms_conditions DROP NOT NULL;
ALTER TABLE mandates ALTER COLUMN client_signature DROP NOT NULL;
ALTER TABLE mandates ALTER COLUMN investigator_signature DROP NOT NULL;
ALTER TABLE mandates ALTER COLUMN special_instructions DROP NOT NULL;

ALTER TABLE mandates ADD COLUMN IF NOT EXISTS caseid_encrypted BYTEA;
ALTER TABLE mandates ADD COLUMN IF NOT EXISTS clientid_encrypted BYTEA;
ALTER TABLE mandates ADD COLUMN IF NOT EXISTS mandatenumber_encrypted BYTEA;
ALTER TABLE mandates ADD COLUMN IF NOT EXISTS scopeofwork_encrypted BYTEA;
ALTER TABLE mandates ADD COLUMN IF NOT EXISTS termsconditions_encrypted BYTEA;
ALTER TABLE mandates ADD COLUMN IF NOT EXISTS clientsignature_encrypted BYTEA;
ALTER TABLE mandates ADD COLUMN IF NOT EXISTS investigatorsignature_encrypted BYTEA;
ALTER TABLE mandates ADD COLUMN IF NOT EXISTS specialinstructions_encrypted BYTEA;
ALTER TABLE mandates ADD COLUMN IF NOT EXISTS dek_encrypted BYTEA;
ALTER TABLE mandates ADD COLUMN IF NOT EXISTS key_version INT NOT NULL DEFAULT 1;
ALTER TABLE mandates ADD COLUMN IF NOT EXISTS metadata JSONB;

-- Contracts
ALTER TABLE contracts ALTER COLUMN case_id DROP NOT NULL;
ALTER TABLE contracts ALTER COLUMN client_id DROP NOT NULL;
ALTER TABLE contracts ALTER COLUMN contract_number DROP NOT NULL;
ALTER TABLE contracts ALTER COLUMN scope_of_services DROP NOT NULL;
ALTER TABLE contracts ALTER COLUMN payment_terms DROP NOT NULL;
ALTER TABLE contracts ALTER COLUMN confidentiality DROP NOT NULL;
ALTER TABLE contracts ALTER COLUMN termination_clause DROP NOT NULL;
ALTER TABLE contracts ALTER COLUMN signatures DROP NOT NULL;
ALTER TABLE contracts ALTER COLUMN contract_value DROP NOT NULL;
ALTER TABLE contracts ALTER COLUMN renewal_terms DROP NOT NULL;

ALTER TABLE contracts ADD COLUMN IF NOT EXISTS caseid_encrypted BYTEA;
ALTER TABLE contracts ADD COLUMN IF NOT EXISTS clientid_encrypted BYTEA;
ALTER TABLE contracts ADD COLUMN IF NOT EXISTS contractnumber_encrypted BYTEA;
ALTER TABLE contracts ADD COLUMN IF NOT EXISTS scopeofservices_encrypted BYTEA;
ALTER TABLE contracts ADD COLUMN IF NOT EXISTS paymentterms_encrypted BYTEA;
ALTER TABLE contracts ADD COLUMN IF NOT EXISTS confidentiality_encrypted BYTEA;
ALTER TABLE contracts ADD COLUMN IF NOT EXISTS terminationclause_encrypted BYTEA;
ALTER TABLE contracts ADD COLUMN IF NOT EXISTS signatures_encrypted BYTEA;
ALTER TABLE contracts ADD COLUMN IF NOT EXISTS contractvalue_encrypted BYTEA;
ALTER TABLE contracts ADD COLUMN IF NOT EXISTS renewalterms_encrypted BYTEA;
ALTER TABLE contracts ADD COLUMN IF NOT EXISTS dek_encrypted BYTEA;
ALTER TABLE contracts ADD COLUMN IF NOT EXISTS key_version INT NOT NULL DEFAULT 1;
ALTER TABLE contracts ADD COLUMN IF NOT EXISTS metadata JSONB;

-- Invoices
ALTER TABLE invoices ALTER COLUMN case_id DROP NOT NULL;
ALTER TABLE invoices ALTER COLUMN client_id DROP NOT NULL;
ALTER TABLE invoices ALTER COLUMN invoice_number DROP NOT NULL;
ALTER TABLE invoices ALTER COLUMN line_items DROP NOT NULL;
ALTER TABLE invoices ALTER COLUMN total_amount DROP NOT NULL;
ALTER TABLE invoices ALTER COLUMN tax_amount DROP NOT NULL;
ALTER TABLE invoices ALTER COLUMN notes DROP NOT NULL;
ALTER TABLE invoices ALTER COLUMN paid_amount DROP NOT NULL;
ALTER TABLE invoices ALTER COLUMN payment_method DROP NOT NULL;
ALTER TABLE invoices ALTER COLUMN payment_terms DROP NOT NULL;
ALTER TABLE invoices ALTER COLUMN late_fee DROP NOT NULL;

ALTER TABLE invoices ADD COLUMN IF NOT EXISTS caseid_encrypted BYTEA;
ALTER TABLE invoices ADD COLUMN IF NOT EXISTS clientid_encrypted BYTEA;
ALTER TABLE invoices ADD COLUMN IF NOT EXISTS invoicenumber_encrypted BYTEA;
ALTER TABLE invoices ADD COLUMN IF NOT EXISTS lineitems_encrypted BYTEA;
ALTER TABLE invoices ADD COLUMN IF NOT EXISTS totalamount_encrypted BYTEA;
ALTER TABLE invoices ADD COLUMN IF NOT EXISTS taxamount_encrypted BYTEA;
ALTER TABLE invoices ADD COLUMN IF NOT EXISTS notes_encrypted BYTEA;
ALTER TABLE invoices ADD COLUMN IF NOT EXISTS paidamount_encrypted BYTEA;
ALTER TABLE invoices ADD COLUMN IF NOT EXISTS paymentmethod_encrypted BYTEA;
ALTER TABLE invoices ADD COLUMN IF NOT EXISTS paymentterms_encrypted BYTEA;
ALTER TABLE invoices ADD COLUMN IF NOT EXISTS latefee_encrypted BYTEA;
ALTER TABLE invoices ADD COLUMN IF NOT EXISTS dek_encrypted BYTEA;
ALTER TABLE invoices ADD COLUMN IF NOT EXISTS key_version INT NOT NULL DEFAULT 1;
ALTER TABLE invoices ADD COLUMN IF NOT EXISTS metadata JSONB;

-- +goose Down
-- Remove encx columns from document tables
ALTER TABLE estimates DROP COLUMN IF EXISTS caseid_encrypted;
ALTER TABLE estimates DROP COLUMN IF EXISTS clientid_encrypted;
ALTER TABLE estimates DROP COLUMN IF EXISTS estimatenumber_encrypted;
ALTER TABLE estimates DROP COLUMN IF EXISTS lineitems_encrypted;
ALTER TABLE estimates DROP COLUMN IF EXISTS estimatedtotal_encrypted;
ALTER TABLE estimates DROP COLUMN IF EXISTS notes_encrypted;
ALTER TABLE estimates DROP COLUMN IF EXISTS dek_encrypted;
ALTER TABLE estimates DROP COLUMN IF EXISTS key_version;
ALTER TABLE estimates DROP COLUMN IF EXISTS metadata;

ALTER TABLE mandates DROP COLUMN IF EXISTS caseid_encrypted;
ALTER TABLE mandates DROP COLUMN IF EXISTS clientid_encrypted;
ALTER TABLE mandates DROP COLUMN IF EXISTS mandatenumber_encrypted;
ALTER TABLE mandates DROP COLUMN IF EXISTS scopeofwork_encrypted;
ALTER TABLE mandates DROP COLUMN IF EXISTS termsconditions_encrypted;
ALTER TABLE mandates DROP COLUMN IF EXISTS clientsignature_encrypted;
ALTER TABLE mandates DROP COLUMN IF EXISTS investigatorsignature_encrypted;
ALTER TABLE mandates DROP COLUMN IF EXISTS specialinstructions_encrypted;
ALTER TABLE mandates DROP COLUMN IF EXISTS dek_encrypted;
ALTER TABLE mandates DROP COLUMN IF EXISTS key_version;
ALTER TABLE mandates DROP COLUMN IF EXISTS metadata;

ALTER TABLE contracts DROP COLUMN IF EXISTS caseid_encrypted;
ALTER TABLE contracts DROP COLUMN IF EXISTS clientid_encrypted;
ALTER TABLE contracts DROP COLUMN IF EXISTS contractnumber_encrypted;
ALTER TABLE contracts DROP COLUMN IF EXISTS scopeofservices_encrypted;
ALTER TABLE contracts DROP COLUMN IF EXISTS paymentterms_encrypted;
ALTER TABLE contracts DROP COLUMN IF EXISTS confidentiality_encrypted;
ALTER TABLE contracts DROP COLUMN IF EXISTS terminationclause_encrypted;
ALTER TABLE contracts DROP COLUMN IF EXISTS signatures_encrypted;
ALTER TABLE contracts DROP COLUMN IF EXISTS contractvalue_encrypted;
ALTER TABLE contracts DROP COLUMN IF EXISTS renewalterms_encrypted;
ALTER TABLE contracts DROP COLUMN IF EXISTS dek_encrypted;
ALTER TABLE contracts DROP COLUMN IF EXISTS key_version;
ALTER TABLE contracts DROP COLUMN IF EXISTS metadata;

ALTER TABLE invoices DROP COLUMN IF EXISTS caseid_encrypted;
ALTER TABLE invoices DROP COLUMN IF EXISTS clientid_encrypted;
ALTER TABLE invoices DROP COLUMN IF EXISTS invoicenumber_encrypted;
ALTER TABLE invoices DROP COLUMN IF EXISTS lineitems_encrypted;
ALTER TABLE invoices DROP COLUMN IF EXISTS totalamount_encrypted;
ALTER TABLE invoices DROP COLUMN IF EXISTS taxamount_encrypted;
ALTER TABLE invoices DROP COLUMN IF EXISTS notes_encrypted;
ALTER TABLE invoices DROP COLUMN IF EXISTS paidamount_encrypted;
ALTER TABLE invoices DROP COLUMN IF EXISTS paymentmethod_encrypted;
ALTER TABLE invoices DROP COLUMN IF EXISTS paymentterms_encrypted;
ALTER TABLE invoices DROP COLUMN IF EXISTS latefee_encrypted;
ALTER TABLE invoices DROP COLUMN IF EXISTS dek_encrypted;
ALTER TABLE invoices DROP COLUMN IF EXISTS key_version;
ALTER TABLE invoices DROP COLUMN IF EXISTS metadata;
