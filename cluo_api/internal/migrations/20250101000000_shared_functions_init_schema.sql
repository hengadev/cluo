-- +goose Up
-- +goose StatementBegin

-- Create uuid extension if not exists
CREATE EXTENSION IF NOT EXISTS "uuid-ossp";

-- Create shared trigger function for updating updated_at columns
-- This function will be used by all tables that need automatic timestamp updates
CREATE OR REPLACE FUNCTION update_updated_at_column()
RETURNS TRIGGER AS $$
BEGIN
    NEW.updated_at = NOW();
    RETURN NEW;
END;
$$ language 'plpgsql';

-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin

-- Drop the shared trigger function
DROP FUNCTION IF EXISTS update_updated_at_column();

-- +goose StatementEnd