-- +goose Up
-- password_hash_secure was incorrectly defined as BYTEA; encx hash_secure
-- produces a bcrypt string, so the column must be text-compatible.
-- convert_from recovers the original string from the bytes written by pgx.
ALTER TABLE auth.users
    ALTER COLUMN password_hash_secure TYPE VARCHAR(255)
    USING convert_from(password_hash_secure, 'UTF8');

-- +goose Down
ALTER TABLE auth.users
    ALTER COLUMN password_hash_secure TYPE BYTEA
    USING convert_to(password_hash_secure, 'UTF8');
