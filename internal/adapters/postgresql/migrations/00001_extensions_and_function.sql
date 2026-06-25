-- +goose Up

-- Enable UUID generation
CREATE EXTENSION IF NOT EXISTS "uuid-ossp";

-- Shared trigger function for updated_at columns

-- +goose StatementBegin
CREATE OR REPLACE FUNCTION update_updated_at_column()
RETURNS TRIGGER AS $$
BEGIN
    NEW.updated_at = CURRENT_TIMESTAMP;
    RETURN NEW;
END;
$$ LANGUAGE plpgsql;
-- +goose StatementEnd


-- +goose Down

DROP FUNCTION IF EXISTS update_updated_at_column();
DROP EXTENSION IF EXISTS "uuid-ossp";
