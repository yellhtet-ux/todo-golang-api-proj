package todos

import (
	"context"
	"fmt"

	"github.com/jackc/pgx/v5/pgtype"
	repo "github.com/yellhtet-ux/todo-golang-api-proj/internal/adapters/postgresql/sqlc"
)

var (
	TodoNotFoundError = "todo not found"
	InvalidToDoIDError = "todo id is not valid"
	InvalidUserIDError = "user id is not valid"
)

type Service interface {
	// GET
	ListToDos(ctx context.Context,userID pgtype.UUID) ([]repo.Todo,error) 
	ListToDosByID(ctx context.Context, todoID pgtype.UUID,userID pgtype.UUID) (repo.Todo,error)

	// POST
	CreateTodo(ctx context.Context,userID pgtype.UUID,todo CreateTodoRequest) (repo.Todo,error)

	// PUT
	UpdateTodoByStatus(ctx context.Context,todoID pgtype.UUID,userID pgtype.UUID,todo UpdateToDoStatusRequest) (repo.Todo,error)
	UpdateToDoByPriority(ctx context.Context,todoID pgtype.UUID,userID pgtype.UUID,todo UpdateToDoPriorityRequest) (repo.Todo,error)

	// DELETE
	DeleteTodoByID(ctx context.Context,todoID pgtype.UUID,userID pgtype.UUID) error
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

func (s *svc) ListToDos(ctx context.Context, userID pgtype.UUID) ([]repo.Todo,error) {
	return s.repo.ListToDos(ctx, userID)
}

func (s *svc) ListToDosByID(ctx context.Context,todoID pgtype.UUID,userID pgtype.UUID) (repo.Todo,error) {
	isUserIDValid := userID.Valid
	isToDoIDValid := todoID.Valid

	if isUserIDValid {
		if isToDoIDValid {
			params := repo.ListToDosByIDParams{
				ID: todoID,
				UserID: userID,
			}
		// If yes, query the data from data with UserID and Todo ID
			todo, err := s.repo.ListToDosByID(ctx,params)

		// If something went wrong, Not found error would be appeared.
		if err != nil {
			return repo.Todo{}, fmt.Errorf("error not found: %s", TodoNotFoundError)
		}else {
		// Unless Got the DATA Yayyyy
			return todo, nil
		}
		}else {
			return repo.Todo{}, fmt.Errorf("invalid todo id %s", InvalidToDoIDError)
		}
	}else {
			return repo.Todo{},fmt.Errorf("invalid user id %s",InvalidUserIDError)
	}
}

func (s *svc) CreateTodo(ctx context.Context,userID pgtype.UUID,todo CreateTodoRequest) (repo.Todo,error) {

	param := repo.CreateToDoParams{
		  UserID: userID,
			Title: todo.Title,
			Description: todo.Description,
			DueAt: pgtype.Timestamptz{Time: todo.DueDate,Valid: true},
			Priority: repo.TodoPriority(todo.Priority),
			Status: repo.TodoStatus(todo.Status),	
		}

		createdTodo, err := s.repo.CreateToDo(ctx,param)

		if err  != nil {
			return repo.Todo{},fmt.Errorf("error not found %s", TodoNotFoundError)
		}

		return createdTodo, nil

}

func (s *svc) UpdateTodoByStatus(ctx context.Context,todoID pgtype.UUID,userID pgtype.UUID,todo UpdateToDoStatusRequest) (repo.Todo,error) {
	isUserIDValid := userID.Valid
	isToDoIDValid := todoID.Valid

	if isUserIDValid {	
		if isToDoIDValid {
			updateTodoBySParams := repo.UpdateToDoStatusParams{
					ID: todoID,
					UserID: userID,
					Status: repo.TodoStatus(todo.Status),
			}
			todo, err := s.repo.UpdateToDoStatus(ctx,updateTodoBySParams)
			if err != nil {
				return repo.Todo{}, fmt.Errorf("error not found: %s",TodoNotFoundError)
			}else {
				return todo, nil
			}
		}else {
			return repo.Todo{}, fmt.Errorf("invalid todo id: %s", InvalidToDoIDError)
		}
	}else {
		return repo.Todo{},fmt.Errorf("invalid user id: %s", InvalidUserIDError)
	}
} 

func (s *svc) UpdateToDoByPriority(ctx context.Context,todoID pgtype.UUID,userID pgtype.UUID,todo UpdateToDoPriorityRequest) (repo.Todo,error) {

	isUserIDValid := userID.Valid
	isToDoIDValid := todoID.Valid

	if isUserIDValid {
		if isToDoIDValid {
						updateTodoPParams := repo.UpdateToDoPriorityParams{
							UserID: userID,
							ID: todoID,
							Priority: repo.TodoPriority(todo.Priority),
					}
					updatedTodo, err := s.repo.UpdateToDoPriority(ctx,updateTodoPParams)
						
					if err != nil {
						return repo.Todo{}, fmt.Errorf("not found error: %s",err)
					}else {
						return updatedTodo, nil
					}
			}else {
					return repo.Todo{}, fmt.Errorf("invalid todo id: %s", InvalidToDoIDError)
				}
			}else {
					return repo.Todo{},fmt.Errorf("invalid user id: %s",InvalidUserIDError)
			}
}

func (s *svc) DeleteTodoByID(ctx context.Context,todoID pgtype.UUID,userID pgtype.UUID) error {
	// Check IDs are valid 
	isUserIDValid := userID.Valid
	isToDoIDValid := todoID.Valid

	if isUserIDValid {
		if isToDoIDValid {
			deleteTodoByIDParam := repo.DeleteTodoByIDParams {
				UserID: userID,
				ID: todoID,
			}

			err := s.repo.DeleteTodoByID(ctx,deleteTodoByIDParam)
			if err != nil {
				return fmt.Errorf("error not found: %s",TodoNotFoundError)
			}else {
				return nil
			}
		}else {	
			return fmt.Errorf("invalid user id: %s",InvalidToDoIDError)
		}
	}else {
		return fmt.Errorf("invalid user id: %s",InvalidUserIDError)
	}
}
