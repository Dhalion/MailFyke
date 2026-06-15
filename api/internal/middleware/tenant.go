package middleware

import (
	"context"
	"github.com/go-chi/chi/v5"
	"net/http"
)

const OrgIDKey contextKey = "org_id"

func RequireOrgMembership(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		orgID := chi.URLParam(r, "orgId")
		if orgID == "" {
			http.Error(w, "missing orgId", http.StatusBadRequest)
			return
		}

		// TODO: verify user (from context) is member of this org or is admin
		ctx := context.WithValue(r.Context(), OrgIDKey, orgID)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}
