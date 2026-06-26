-- ============================================================================
-- USER QUERIES
-- ============================================================================

-- name: CreateUser :one
INSERT INTO users (
    email,
    password_hash,
    display_name
) VALUES (
    $1, $2, $3
)
RETURNING *;

-- name: GetUserByID :one
SELECT * FROM users 
WHERE id = $1 AND deleted_at IS NULL;


-- ============================================================================
-- TODO QUERIES (Updated with User Scope)
-- ============================================================================

-- name: ListToDos :many
-- Scoped to only list the current user's active todos
SELECT * FROM todos 
WHERE user_id = $1 AND deleted_at IS NULL;

-- name: ListToDosByID :one
-- Ensures users can only fetch their own todos
SELECT * FROM todos 
WHERE id = $1 AND user_id = $2 AND deleted_at IS NULL;

-- name: CreateToDo :one
INSERT INTO todos (
    user_id, -- 👈 Added link to the creator
    title,
    description,
    status,
    priority,
    due_at
) VALUES (
    $1, $2, $3, $4, $5, $6
)
RETURNING *;

-- name: UpdateToDoStatus :one
UPDATE todos
SET 
    status = $3,
    completed_at = CASE WHEN $3 = 'completed'::todo_status THEN CURRENT_TIMESTAMP ELSE NULL END,
    updated_at = CURRENT_TIMESTAMP
WHERE id = $1 AND user_id = $2 -- 👈 Security boundary check
RETURNING *;

-- name: UpdateToDoPriority :one
UPDATE todos
SET
    priority = $3,
    updated_at = CURRENT_TIMESTAMP
WHERE id = $1 AND user_id = $2 -- 👈 Security boundary check
RETURNING *;

-- name: DeleteTodoByID :exec
-- Soft delete or hard delete scoped to user
UPDATE todos 
SET deleted_at = CURRENT_TIMESTAMP 
WHERE id = $1 AND user_id = $2;
