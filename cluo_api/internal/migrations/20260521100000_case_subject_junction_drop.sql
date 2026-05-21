-- +goose Up
DROP TABLE IF EXISTS cases.case_subject_cases;

-- +goose Down
-- Recreate the junction table (for rollback safety)
CREATE TABLE IF NOT EXISTS cases.case_subject_cases (
    id UUID PRIMARY KEY,
    case_subject_id UUID NOT NULL,
    case_id UUID NOT NULL,
    roles TEXT[] NOT NULL DEFAULT '{}',
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    CONSTRAINT fk_case_subject_cases_subject_id
        FOREIGN KEY (case_subject_id) REFERENCES cases.case_subjects(id) ON DELETE CASCADE,
    CONSTRAINT fk_case_subject_cases_case_id
        FOREIGN KEY (case_id) REFERENCES cases.cases(id) ON DELETE CASCADE,
    CONSTRAINT uq_case_subject_cases_subject_case
        UNIQUE (case_subject_id, case_id)
);
