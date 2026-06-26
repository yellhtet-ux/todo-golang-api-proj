package todos

import (
	"time"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgtype"
)

type TodoStatus string

const (
	StatusTodo      TodoStatus = "todo"
	StatusPending   TodoStatus = "pending"
	StatusCompleted TodoStatus = "completed"
)

type TodoPriority string

const (
	PriorityLow    TodoPriority = "low"
	PriorityMedium TodoPriority = "medium"
	PriorityHigh   TodoPriority = "high"
)

type Todo struct {
	ID          uuid.UUID
	Title       string
	Description string
	DueDate     time.Time
	Status      TodoStatus
	Priority    TodoPriority
}

type CreateTodoRequest struct {
	UserID      pgtype.UUID  `json:"user_id"`
	Title       string       `json:"title"`
	Description pgtype.Text  `json:"description"`
	DueDate     time.Time    `json:"due_date"`
	Status      TodoStatus   `json:"status"`
	Priority    TodoPriority `json:"priority"`
}

// CreateTodoRequestDoc documents the create todo request body for Swagger.
type CreateTodoRequestDoc struct {
	UserID      pgtype.UUID  `json:"user_id"`
	Title       string       `json:"title" example:"Buy groceries"`
	Description string       `json:"description" example:"Milk and eggs"`
	DueDate     time.Time    `json:"due_date"`
	Status      TodoStatus   `json:"status" example:"todo"`
	Priority    TodoPriority `json:"priority" example:"medium"`
}

type UpdateTodoStatus struct {
	Status TodoStatus `json:"status"`
}

type UpdateTodoPriority struct {
	Priority TodoPriority `json:"priority"`
}

// TodoResponse represents a todo item returned by the API.
type TodoResponse struct {
	ID          string       `json:"id" example:"550e8400-e29b-41d4-a716-446655440000"`
	Title       string       `json:"title" example:"Buy groceries"`
	Description string       `json:"description" example:"Milk and eggs"`
	Status      TodoStatus   `json:"status" example:"todo" enums:"todo,in_progress,completed"`
	Priority    TodoPriority `json:"priority" example:"medium" enums:"low,medium,high"`
	DueAt       time.Time    `json:"due_at"`
	CompletedAt *time.Time   `json:"completed_at,omitempty"`
	CreatedAt   time.Time    `json:"created_at"`
	UpdatedAt   time.Time    `json:"updated_at"`
}


// ListToDoByIDParam... 
type ListToDosByIDParam struct {
	ID     pgtype.UUID `json:"id"`
	UserID pgtype.UUID `json:"user_id"`
}

// ErrorResponse represents an error response body.
type ErrorResponse struct {
	Message string `json:"message" example:"Invalid request"`
	Error   string `json:"error" example:"invalid UUID"`
}
