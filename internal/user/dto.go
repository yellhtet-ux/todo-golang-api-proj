package user

import "github.com/jackc/pgx/v5/pgtype"


type createUserParam struct {
  Email        string      `json:"email"`
	PasswordHash string      `json:"password"`
	DisplayName  pgtype.Text `json:"display_name"`
}
