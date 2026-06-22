-- name: ListToDos :many
SELECT * FROM todos;

-- name: ListToDosByID :one
SELECT * FROM todos WHERE id = $1;

-- name: CreateToDo :one
INSERT INTO todos (
    title,
    description,
    status,
    priority,
    due_at
) VALUES (
    $1, $2, $3, $4, $5
)
RETURNING *;

-- name: UpdateToDoStatus :one
UPDATE todos
SET 
    status = $2,
    completed_at = CASE WHEN $2 = 'completed'::todo_status THEN CURRENT_TIMESTAMP ELSE NULL END,
    updated_at = CURRENT_TIMESTAMP
WHERE id = $1
RETURNING *;

-- name: UpdateToDoPriority :one
UPDATE todos
SET
    priority = $2,
    updated_at = CURRENT_TIMESTAMP
WHERE id = $1
RETURNING *;

-- name: DeleteTodoByID :exec
DELETE FROM todos WHERE id = $1;