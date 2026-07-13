package user

import (
	"github.com/golang-jwt/jwt/v5"
	"github.com/jackc/pgx/v5/pgtype"
)

// Create User Request
type CreateUserRequest struct {
  Email        string      `json:"email"`
	Password string      `json:"password"`
	DisplayName  pgtype.Text `json:"display_name"`
}

// Create User Response
type CreateUserResponse struct {
	ID pgtype.UUID `json:"id"`
	Email string  `json:"email"`
	DisplayName string `json:"display_name"`
}

type LoginUserRequest struct {
	Email string `json:"email"`
	Password string `json:"password"`
}

type LoginUserResponse struct {
	AccessToken string `json:"access_token"`
}

// Get User By ID Response
type GetUserByIDResponse struct {
	ID pgtype.UUID `json:"id"`
	Email string  `json:"email"`
	DisplayName string `json:"display_name"`
	CreatedAt string `json:"created_at"`
	UpdatedAt string `json:"updated_at"`
	DeletedAt string `json:"deleted_at"`
}

type CustomClaims struct {
	UserID pgtype.UUID `json:"user_id"`
	Email string `json:"email"`

	jwt.RegisteredClaims
}
