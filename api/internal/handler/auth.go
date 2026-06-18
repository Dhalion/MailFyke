package handler

import (
	"fmt"
	"net/http"
	"time"

	"github.com/Dhalion/MailFyke/internal/httputil"
	"github.com/golang-jwt/jwt/v5"
	"github.com/jackc/pgx/v5/pgtype"
	"golang.org/x/crypto/bcrypt"
)

type LoginRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

type LoginResponse struct {
	Token         string   `json:"token"`
	User          UserInfo `json:"user"`
	Organizations []OrgRow `json:"organizations"`
}

type UserInfo struct {
	ID      string `json:"id"`
	Email   string `json:"email"`
	IsAdmin bool   `json:"is_admin"`
}

type OrgRow struct {
	ID   string `json:"id"`
	Name string `json:"name"`
	Slug string `json:"slug"`
	Role string `json:"role"`
}

func uuidStr(id pgtype.UUID) string {
	return fmt.Sprintf("%08x-%04x-%04x-%04x-%012x", id.Bytes[0:4], id.Bytes[4:6], id.Bytes[6:8], id.Bytes[8:10], id.Bytes[10:16])
}

func (h *Handler) AuthLogin(w http.ResponseWriter, r *http.Request) {
	var req LoginRequest
	if err := httputil.DecodeBody(r, &req); err != nil {
		httputil.WriteError(w, http.StatusBadRequest, "invalid request body")
		return
	}

	if req.Email == "" || req.Password == "" {
		httputil.WriteError(w, http.StatusBadRequest, "email and password required")
		return
	}

	user, err := h.queries.GetUserByEmail(r.Context(), req.Email)
	if err != nil {
		httputil.WriteError(w, http.StatusUnauthorized, "invalid credentials")
		return
	}

	if err := bcrypt.CompareHashAndPassword([]byte(user.PasswordHash), []byte(req.Password)); err != nil {
		httputil.WriteError(w, http.StatusUnauthorized, "invalid credentials")
		return
	}

	orgRows, err := h.queries.ListOrganizationsByUser(r.Context(), user.ID)
	if err != nil {
		orgRows = nil
	}

	orgIDs := make([]string, len(orgRows))
	orgs := make([]OrgRow, len(orgRows))
	for i, o := range orgRows {
		orgIDs[i] = uuidStr(o.ID)
		orgs[i] = OrgRow{
			ID:   uuidStr(o.ID),
			Name: o.Name,
			Slug: o.Slug,
			Role: o.Role,
		}
	}

	now := time.Now()
	claims := jwt.MapClaims{
		"user_id":  uuidStr(user.ID),
		"email":    user.Email,
		"is_admin": user.IsAdmin,
		"org_ids":  orgIDs,
		"iat":      now.Unix(),
		"exp":      now.Add(24 * time.Hour).Unix(),
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenStr, err := token.SignedString([]byte(h.cfg.JWTSecret))
	if err != nil {
		httputil.WriteError(w, http.StatusInternalServerError, "failed to generate token")
		return
	}

	httputil.WriteJSON(w, http.StatusOK, LoginResponse{
		Token: tokenStr,
		User: UserInfo{
			ID:      uuidStr(user.ID),
			Email:   user.Email,
			IsAdmin: user.IsAdmin,
		},
		Organizations: orgs,
	})
}
