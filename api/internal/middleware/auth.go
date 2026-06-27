package middleware

import (
	"context"
	"fmt"
	"net/http"
	"strings"

	"github.com/Dhalion/MailFyke/internal/httputil"
	"github.com/golang-jwt/jwt/v5"
)

type ctxKey string

const (
	CtxUserID  ctxKey = "user_id"
	CtxIsAdmin ctxKey = "is_admin"
)

func UserIDFromCtx(ctx context.Context) string {
	v, _ := ctx.Value(CtxUserID).(string)
	return v
}

func IsAdminFromCtx(ctx context.Context) bool {
	v, _ := ctx.Value(CtxIsAdmin).(bool)
	return v
}

func Auth(jwtSecret string) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			auth := r.Header.Get("Authorization")

			tokenStr, found := strings.CutPrefix(auth, "Bearer ")
			if !found {
				httputil.WriteError(w, http.StatusUnauthorized, "missing or invalid authorization header")
				return
			}
			token, err := jwt.Parse(tokenStr, func(t *jwt.Token) (interface{}, error) {
				if _, ok := t.Method.(*jwt.SigningMethodHMAC); !ok {
					return nil, fmt.Errorf("unexpected signing method: %v", t.Header["alg"])
				}
				return []byte(jwtSecret), nil
			})
			if err != nil {
				httputil.WriteError(w, http.StatusUnauthorized, "invalid token")
				return
			}

			claims, ok := token.Claims.(jwt.MapClaims)
			if !ok || !token.Valid {
				httputil.WriteError(w, http.StatusUnauthorized, "invalid token")
				return
			}

			userID, _ := claims["user_id"].(string)
			isAdmin, _ := claims["is_admin"].(bool)

			ctx := context.WithValue(r.Context(), CtxUserID, userID)
			ctx = context.WithValue(ctx, CtxIsAdmin, isAdmin)
			next.ServeHTTP(w, r.WithContext(ctx))
		})
	}
}
