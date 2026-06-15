package middleware

import (
	"context"
	"net/http"
)

type contextKey string

const (
	UserIDKey contextKey = "user_id"
	IsAdminKey contextKey = "is_admin"
)

func Auth(jwtSecret string) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			// TODO: validate JWT from Authorization header
			ctx := context.WithValue(r.Context(), UserIDKey, "")
			ctx = context.WithValue(ctx, IsAdminKey, false)
			next.ServeHTTP(w, r.WithContext(ctx))
		})
	}
}
