-- +goose Up

CREATE TYPE todo_status AS ENUM (
    'todo',
    'in_progress',
    'completed'
);

CREATE TYPE todo_priority AS ENUM (
    'low',
    'medium',
    'high'
);

CREATE TABLE todos (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),

    user_id UUID NOT NULL,

    title VARCHAR(255) NOT NULL,
    description TEXT,

    status todo_status NOT NULL DEFAULT 'todo',
    priority todo_priority NOT NULL DEFAULT 'medium',

    due_at TIMESTAMPTZ,
    completed_at TIMESTAMPTZ,

    created_at TIMESTAMPTZ NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMPTZ NOT NULL DEFAULT CURRENT_TIMESTAMP,
    deleted_at TIMESTAMPTZ,

    CONSTRAINT chk_title_length
        CHECK (char_length(title) >= 1),

    CONSTRAINT fk_todos_user
        FOREIGN KEY (user_id)
        REFERENCES users(id)
        ON DELETE CASCADE
);

CREATE TRIGGER set_todos_updated_at
BEFORE UPDATE ON todos
FOR EACH ROW
EXECUTE FUNCTION update_updated_at_column();

CREATE INDEX idx_todos_user_status
ON todos(user_id, status)
WHERE deleted_at IS NULL;

CREATE INDEX idx_todos_user_due_at
ON todos(user_id, due_at)
WHERE deleted_at IS NULL
AND status != 'completed';


-- +goose Down

DROP INDEX IF EXISTS idx_todos_user_due_at;
DROP INDEX IF EXISTS idx_todos_user_status;

DROP TRIGGER IF EXISTS set_todos_updated_at ON todos;

DROP TABLE IF EXISTS todos;

DROP TYPE IF EXISTS todo_priority CASCADE;
DROP TYPE IF EXISTS todo_status CASCADE;
