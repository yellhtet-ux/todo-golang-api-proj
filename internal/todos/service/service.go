package todos

import (
	"context"

	"github.com/jackc/pgx/v5/pgtype"
	repo "github.com/yellhtet-ux/todo-golang-api-proj/internal/adapters/postgresql/sqlc"
)

type Service interface {
	// GET
	ListToDos(ctx context.Context) ([]repo.Todo,error) 
	ListToDosByID(ctx context.Context, id pgtype.UUID) (repo.Todo,error)

	// POST
	CreateTodo(ctx context.Context, todo *repo.CreateToDoParams) (repo.Todo,error)

	// PUT
	UpdateTodoByStatus(ctx context.Context, todo *repo.UpdateToDoStatusParams) (repo.Todo,error)
	UpdateToDoByPriority(ctx context.Context,todo *repo.UpdateToDoPriorityParams) (repo.Todo,error)

	// DELETE
	DeleteTodoByID(ctx context.Context, id pgtype.UUID) error
}

type svc struct {
	// repositories
	repo repo.Querier
}

func NewService (repo repo.Querier) Service {
	return &svc{
		repo: repo,
	}
}

func (s *svc) ListToDos(ctx context.Context) ([]repo.Todo,error) {
	return s.repo.ListToDos(ctx)
}

func (s *svc) ListToDosByID(ctx context.Context, id pgtype.UUID) (repo.Todo,error) {
	return s.repo.ListToDosByID(ctx, id)
}

func (s *svc) CreateTodo(ctx context.Context, todo *repo.CreateToDoParams) (repo.Todo,error) {
	return s.repo.CreateToDo(ctx,repo.CreateToDoParams{
		Title: todo.Title,
		Description: todo.Description,
		Status: todo.Status,
		Priority: todo.Priority,
		DueAt: todo.DueAt,
	})
}

func (s *svc) UpdateTodoByStatus(ctx context.Context,todo *repo.UpdateToDoStatusParams) (repo.Todo,error) {
	return s.repo.UpdateToDoStatus(ctx,*todo)
} 

func (s *svc) UpdateToDoByPriority(ctx context.Context,todo *repo.UpdateToDoPriorityParams) (repo.Todo,error) {
	return s.repo.UpdateToDoPriority(ctx,*todo)
}

func (s *svc) DeleteTodoByID(ctx context.Context, id pgtype.UUID) error {
	return s.repo.DeleteTodoByID(ctx,id)
}