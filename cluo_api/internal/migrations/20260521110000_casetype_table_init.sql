-- +goose Up
CREATE TABLE IF NOT EXISTS cases.case_types (
    id UUID PRIMARY KEY,
    name TEXT NOT NULL,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    CONSTRAINT uq_case_types_name UNIQUE (name)
);

-- Seed common investigation types
INSERT INTO cases.case_types (id, name, created_at, updated_at) VALUES
    (gen_random_uuid(), 'surveillance', NOW(), NOW()),
    (gen_random_uuid(), 'insurance fraud', NOW(), NOW()),
    (gen_random_uuid(), 'arson investigation', NOW(), NOW()),
    (gen_random_uuid(), 'background check', NOW(), NOW()),
    (gen_random_uuid(), 'missing person', NOW(), NOW());

-- Add new FK column (nullable)
ALTER TABLE cases.cases ADD COLUMN case_type_id UUID NULL
    REFERENCES cases.case_types(id) ON DELETE RESTRICT;

-- Backfill: match existing case_type strings to seeded types
UPDATE cases.cases c
SET case_type_id = ct.id
FROM cases.case_types ct
WHERE lower(c.case_type) = lower(ct.name);

-- Drop old string column
ALTER TABLE cases.cases DROP COLUMN case_type;

CREATE INDEX IF NOT EXISTS idx_cases_case_type_id ON cases.cases(case_type_id);

-- +goose Down
ALTER TABLE cases.cases ADD COLUMN case_type VARCHAR(100) NOT NULL DEFAULT '';
UPDATE cases.cases c SET case_type = ct.name FROM cases.case_types ct WHERE c.case_type_id = ct.id;
ALTER TABLE cases.cases DROP COLUMN case_type_id;
DROP TABLE IF EXISTS cases.case_types;
