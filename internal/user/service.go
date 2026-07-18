package user

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/jackc/pgx/v5/pgconn"
	"github.com/jackc/pgx/v5/pgtype"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/yellhtet-ux/todo-golang-api-proj/env"
	repo "github.com/yellhtet-ux/todo-golang-api-proj/internal/adapters/postgresql/sqlc"
	"golang.org/x/crypto/bcrypt"
)

var UserNotFoundError = "user is not found"
var JWTSecretKey = []byte(env.GetString("JWT_SECRET_KEY", "kit_ko-amyar_gyi-chit_tll"))

var (
	ErrInvalidSignup     = errors.New("invalid email or password too short")
	ErrUserAlreadyExists = errors.New("user with this email already exists")
)

func generateJWT(email string, UserID pgtype.UUID) (string, error) {
	expirationTime := time.Now().Add(1 * time.Hour)

	claims := &CustomClaims{
		Email:  email,
		UserID: UserID,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(expirationTime),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	tokenString, err := token.SignedString(JWTSecretKey)

	if err != nil {
		return "", err
	}

	return tokenString, nil
}

type Service interface {
	CreateUser(ctx context.Context, params CreateUserRequest) (CreateUserResponse, error)
	GetUserByID(ctx context.Context, userID pgtype.UUID) (GetUserByIDResponse, error)
	GetUserByEmail(ctx context.Context, params LoginUserRequest) (LoginUserResponse, error)
}

type svc struct {
	repo repo.Querier
	db   *pgxpool.Pool
}

func NewService(repo repo.Querier, db *pgxpool.Pool) Service {
	return &svc{
		repo: repo,
		db:   db,
	}
}

func (s *svc) CreateUser(ctx context.Context, params CreateUserRequest) (CreateUserResponse, error) {
	if params.Email == "" || len(params.Password) < 8 || params.DisplayName.String == "" {
		return CreateUserResponse{}, ErrInvalidSignup
	}

	// Hashed password and save it in database
	hashedBytes, err := bcrypt.GenerateFromPassword([]byte(params.Password), 12)
	if err != nil {
		return CreateUserResponse{}, errors.New("failed to process security credentials")
	}

	userParams := repo.CreateUserParams{
		Email:        params.Email,
		PasswordHash: string(hashedBytes),
		DisplayName:  params.DisplayName,
	}

	user, err := s.repo.CreateUser(ctx, userParams)

	if err != nil {
		var pgErr *pgconn.PgError
		if errors.As(err, &pgErr) && pgErr.Code == "23505" {
			return CreateUserResponse{}, ErrUserAlreadyExists
		}

		return CreateUserResponse{}, fmt.Errorf("create user: %w", err)
	}

	return CreateUserResponse{
		ID:          user.ID,
		Email:       user.Email,
		DisplayName: user.DisplayName.String,
	}, nil
}

func (s *svc) GetUserByID(ctx context.Context, userID pgtype.UUID) (GetUserByIDResponse, error) {
	isUserIDValid := userID.Valid

	if isUserIDValid {
		user, err := s.repo.GetUserByID(ctx, userID)

		if err != nil {
			return GetUserByIDResponse{}, fmt.Errorf("error not found: %s", UserNotFoundError)
		}

		return GetUserByIDResponse{
			ID:          user.ID,
			Email:       user.Email,
			DisplayName: user.DisplayName.String,
			CreatedAt:   user.CreatedAt.Time.GoString(),
			UpdatedAt:   user.UpdatedAt.Time.GoString(),
			DeletedAt:   user.DeletedAt.Time.GoString(),
		}, nil

	} else {
		return GetUserByIDResponse{}, errors.New("user id should be valid")
	}
}

func (s *svc) GetUserByEmail(ctx context.Context, params LoginUserRequest) (LoginUserResponse, error) {
	if params.Email == "" || len(params.Password) < 8 {
		return LoginUserResponse{}, errors.New("invalid email or password too short")
	}

	resultUser, err := s.repo.GetUserByEmail(ctx, params.Email)

	if err != nil {
		return LoginUserResponse{}, err
	} else {
		jwtToken, err := generateJWT(params.Email, resultUser.ID)

		if err != nil {
			return LoginUserResponse{}, err
		} else {
			return LoginUserResponse{
				AccessToken: jwtToken,
			}, nil
		}
	}
}
