package user_middleware

import (
	"context"
	"errors"
	"net/http"
	"strings"

	"github.com/golang-jwt/jwt/v5"

	"github.com/yellhtet-ux/todo-golang-api-proj/env"
	"github.com/yellhtet-ux/todo-golang-api-proj/internal/json"
	"github.com/yellhtet-ux/todo-golang-api-proj/internal/user"
)

func JWTMiddleware (next http.Handler) http.Handler{
 	return http.HandlerFunc(func (w http.ResponseWriter,r *http.Request) {
		authHeader := r.Header.Get("Authorization")
		if authHeader == "" {
			json.NotAuthorizedError(w,errors.New("authorization header is required"),nil)
			return
		}


		parts := strings.Split(authHeader, " ")
		if len(parts) != 2 || parts[0] != "Bearer" {
			http.Error(w, "Invalid Token Format", http.StatusUnauthorized)
			return
		}
		tokenString := parts[1]

		claims := &user.CustomClaims{}
		token, err := jwt.ParseWithClaims(tokenString,claims,func (t *jwt.Token) (interface{},error) {
			return user.JWTSecretKey, nil
		})

		if err != nil || !token.Valid {
			json.Write(w,http.StatusUnauthorized,"unauthorized: invalid or expired token")
			return
		}

		ctx := context.WithValue(r.Context(),env.GetString("USER_CLAIMS","user_claims"),claims)
		next.ServeHTTP(w,r.WithContext(ctx))
	})
}
