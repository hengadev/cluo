-- +goose Up
-- SQL in this section is executed when the migration is applied.

ALTER TABLE auth.users ADD COLUMN name_encrypted BYTEA;

-- +goose Down
-- SQL in this section is executed when the migration is rolled back.

ALTER TABLE auth.users DROP COLUMN IF EXISTS name_encrypted;
