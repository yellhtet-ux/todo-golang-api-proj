-- +goose Up
-- 1. Enable UUID extension (PostgreSQL specific)
CREATE EXTENSION IF NOT EXISTS "uuid-ossp";

-- 2. Create custom ENUM types for strict data validation
CREATE TYPE todo_status AS ENUM ('todo', 'in_progress', 'completed');
CREATE TYPE todo_priority AS ENUM ('low', 'medium', 'high');

-- 3. Create the todos table
CREATE TABLE IF NOT EXISTS todos (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    title VARCHAR(255) NOT NULL,
    description TEXT,
    status todo_status NOT NULL DEFAULT 'todo',
    priority todo_priority NOT NULL DEFAULT 'medium',
    
    -- Timestamps
    due_at TIMESTAMP WITH TIME ZONE,
    completed_at TIMESTAMP WITH TIME ZONE,
    created_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT CURRENT_TIMESTAMP,
    deleted_at TIMESTAMP WITH TIME ZONE, 

    -- Enforce character limits
    CONSTRAINT chk_title_length CHECK (char_length(title) >= 1)
);

-- 4. Create an automatic update trigger function for `updated_at`
-- +goose StatementBegin
CREATE OR REPLACE FUNCTION update_updated_at_column()
RETURNS TRIGGER AS $$
BEGIN
    NEW.updated_at = CURRENT_TIMESTAMP;
    RETURN NEW;
END;
$$ LANGUAGE plpgsql;
-- +goose StatementEnd

-- 5. Attach the trigger to the todos table
DROP TRIGGER IF EXISTS set_todos_updated_at ON todos;
CREATE TRIGGER set_todos_updated_at
BEFORE UPDATE ON todos
FOR EACH ROW
EXECUTE FUNCTION update_updated_at_column();

-- 6. Production Performance Indexes
-- Updated to look up active tasks by status only (since user_id is removed)
CREATE INDEX IF NOT EXISTS idx_todos_active_status 
ON todos(status) 
WHERE deleted_at IS NULL;

CREATE INDEX IF NOT EXISTS idx_todos_due_at 
ON todos(due_at) 
WHERE deleted_at IS NULL AND status != 'completed';

-- +goose Down
DROP INDEX IF EXISTS idx_todos_due_at;
DROP INDEX IF EXISTS idx_todos_active_status;
DROP TRIGGER IF EXISTS set_todos_updated_at ON todos;
DROP FUNCTION IF EXISTS update_updated_at_column();
DROP TABLE IF EXISTS todos;
DROP TYPE IF EXISTS todo_priority;
DROP TYPE IF EXISTS todo_status;