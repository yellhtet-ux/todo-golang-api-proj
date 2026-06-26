package user

import (
	"context"
	"errors"
	"fmt"

	"github.com/jackc/pgx/v5/pgtype"
	"github.com/jackc/pgx/v5/pgxpool"
	repo "github.com/yellhtet-ux/todo-golang-api-proj/internal/adapters/postgresql/sqlc"
)

var UserNotFoundError = "user is not found"


type Service interface {
  CreateUser(ctx context.Context, params createUserParam) (repo.User, error)
	GetUserByID (ctx context.Context,userID pgtype.UUID) (repo.User, error)
}

type svc struct {
	repo repo.Querier
	db *pgxpool.Pool
}

func NewService (repo repo.Querier,db *pgxpool.Pool) Service {
	return &svc{
		repo: repo,
		db: db,
	}
}

func (s *svc) CreateUser (ctx context.Context, params createUserParam) (repo.User,error) {
	if params.Email == "" {
		return repo.User{},fmt.Errorf("email should not be empty")
	}

	if params.PasswordHash == "" {
		return repo.User{},fmt.Errorf("password should not be empty")
	}

	userParams := repo.CreateUserParams {
		Email: params.Email,
		PasswordHash: params.PasswordHash,
		DisplayName: params.DisplayName,
	}

	user, err := s.repo.CreateUser(ctx,userParams)

	if err != nil {
		return repo.User{},nil
	}

	return user,nil
} 

func (s *svc) GetUserByID (ctx context.Context,userID pgtype.UUID) (repo.User, error) {
	isUserIDValid := userID.Valid

	if isUserIDValid {
		 user, err := s.repo.GetUserByID(ctx,userID);

		 if err != nil {
			 return repo.User{}, fmt.Errorf("error not found: %s",UserNotFoundError)
		 	}

			return user,nil
	}else {
			return repo.User{}, errors.New("user id should be valid")
	}
}
