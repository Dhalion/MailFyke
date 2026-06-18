package middleware

import (
	"context"
	"net/http"

	"github.com/go-chi/chi/v5"
)

const CtxOrgID ctxKey = "org_id"

func OrgIDFromCtx(ctx context.Context) string {
	v, _ := ctx.Value(CtxOrgID).(string)
	return v
}

func RequireOrgMembership(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		orgID := chi.URLParam(r, "orgId")
		if orgID == "" {
			http.Error(w, `{"error":"missing orgId"}`, http.StatusBadRequest)
			return
		}
		ctx := context.WithValue(r.Context(), CtxOrgID, orgID)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}
