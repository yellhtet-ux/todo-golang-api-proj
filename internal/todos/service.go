package todos

import (
	"context"
	"errors"
	"fmt"

	"github.com/jackc/pgx/v5/pgtype"
	repo "github.com/yellhtet-ux/todo-golang-api-proj/internal/adapters/postgresql/sqlc"
)

var (
	TodoNotFoundError = "todo not found"
	InvalidIDError = "Invalid IDs"
)

type Service interface {
	// GET
	ListToDos(ctx context.Context,userID pgtype.UUID) ([]repo.Todo,error) 
	ListToDosByID(ctx context.Context, params ListToDosByIDParam) (repo.Todo,error)

	// POST
	CreateTodo(ctx context.Context, todo CreateTodoRequest) (repo.Todo,error)

	// PUT
	UpdateTodoByStatus(ctx context.Context, todo *repo.UpdateToDoStatusParams) (repo.Todo,error)
	UpdateToDoByPriority(ctx context.Context,todo *repo.UpdateToDoPriorityParams) (repo.Todo,error)

	// DELETE
	DeleteTodoByID(ctx context.Context, params repo.DeleteTodoByIDParams) error
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

func (s *svc) ListToDos(ctx context.Context,userID pgtype.UUID) ([]repo.Todo,error) {
	// check user id is valid
	isUserIDValid := userID.Valid

	if isUserIDValid {
	// check user is already exist
		if _,err := s.repo.GetUserByID(ctx,userID); err != nil {
			return []repo.Todo{}, errors.New("user is not existed")
		}
	} else {
		return []repo.Todo{}, errors.New("user id should be valid")
	} 
	
	// check user is already exist
	if _,err := s.repo.GetUserByID(ctx,userID); err != nil {
			return []repo.Todo{}, errors.New("user is not existed")
	}
	
	return s.repo.ListToDos(ctx,userID)
}

func (s *svc) ListToDosByID(ctx context.Context, params ListToDosByIDParam) (repo.Todo,error) {
	// Check both User and Todo ID are valid
	isUserIDValid := params.UserID.Valid
	isToDoIDValid := params.ID.Valid
	
	// If yes, query the data from data with UserID and Todo ID
	if isUserIDValid && isToDoIDValid {
		listTodoParam := repo.ListToDosByIDParams {
			ID: params.ID,
			UserID: params.UserID,
		}
		todo, err := s.repo.ListToDosByID(ctx,listTodoParam)

		// If something went wrong, Not found error would be appeared.
		if err != nil {
			return repo.Todo{}, fmt.Errorf("error not found: %s", TodoNotFoundError)
		}else {
		// Unless Got the DATA Yayyyy
			return todo, nil
		}
	}else {
		// Unless Invalid ID error would be appeared 
		return repo.Todo{}, fmt.Errorf("invaild ids: %s", InvalidIDError)
	}
}

func (s *svc) CreateTodo(ctx context.Context, todo CreateTodoRequest) (repo.Todo,error) {
	isUserIDValid := todo.UserID.Valid

	if isUserIDValid {
		createTodoParams := repo.CreateToDoParams {
			UserID: todo.UserID,
			Title: todo.Title,
			Description: todo.Description,
			Status: repo.TodoStatus(todo.Status),
			Priority: repo.TodoPriority(todo.Priority),
			DueAt: pgtype.Timestamptz{Time: todo.DueDate,Valid: true},
		}

		todo, err := s.repo.CreateToDo(ctx,createTodoParams)
		if err != nil {
			return repo.Todo{},fmt.Errorf("not found error: %s", TodoNotFoundError)
		}else {
			return todo, nil 
		}
	}else {
		return repo.Todo{},fmt.Errorf("invalid id: %s", InvalidIDError)
	}
}

func (s *svc) UpdateTodoByStatus(ctx context.Context,todo *repo.UpdateToDoStatusParams) (repo.Todo,error) {
	return s.repo.UpdateToDoStatus(ctx,*todo)
} 

func (s *svc) UpdateToDoByPriority(ctx context.Context,todo *repo.UpdateToDoPriorityParams) (repo.Todo,error) {
	return s.repo.UpdateToDoPriority(ctx,*todo)
}

func (s *svc) DeleteTodoByID(ctx context.Context, params repo.DeleteTodoByIDParams) error {
	return s.repo.DeleteTodoByID(ctx,params)
}
