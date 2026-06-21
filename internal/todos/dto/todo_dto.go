package dto

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
	Title       string       `json:"title"`
	Description pgtype.Text       `json:"description"`
	DueDate     time.Time    `json:"due_date"`
	Status      TodoStatus   `json:"status"`
	Priority    TodoPriority `json:"priority"`
}

type UpdateTodoPriority struct {
	Status TodoStatus `json:"status"`
}