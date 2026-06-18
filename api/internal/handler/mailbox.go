package handler

import (
	"net/http"

	"github.com/Dhalion/MailFyke/internal/api"
	"github.com/Dhalion/MailFyke/internal/database/queries"
	"github.com/Dhalion/MailFyke/internal/httputil"
	"github.com/Dhalion/MailFyke/internal/middleware"
	"github.com/jackc/pgx/v5/pgtype"
)

func (h *Handler) ListEmails(w http.ResponseWriter, r *http.Request, orgId string, params api.ListEmailsParams) {
	page := 1
	perPage := 50
	if params.Page != nil && *params.Page > 0 {
		page = *params.Page
	}
	if params.PerPage != nil && *params.PerPage > 0 {
		perPage = *params.PerPage
	}
	offset := (page - 1) * perPage

	var orgUUID pgtype.UUID
	if err := orgUUID.Scan(orgId); err != nil {
		httputil.WriteError(w, http.StatusBadRequest, "invalid orgId")
		return
	}

	emails, err := h.queries.ListEmails(r.Context(), queries.ListEmailsParams{
		OrganizationID: orgUUID,
		Limit:          int32(perPage),
		Offset:         int32(offset),
	})
	if err != nil {
		httputil.WriteError(w, http.StatusInternalServerError, "failed to list emails")
		return
	}

	httputil.WriteJSON(w, http.StatusOK, emails)
}

func (h *Handler) GetEmail(w http.ResponseWriter, r *http.Request, orgId string, id string) {
	var orgUUID, emailUUID pgtype.UUID
	if err := orgUUID.Scan(orgId); err != nil {
		httputil.WriteError(w, http.StatusBadRequest, "invalid orgId")
		return
	}
	if err := emailUUID.Scan(id); err != nil {
		httputil.WriteError(w, http.StatusBadRequest, "invalid email id")
		return
	}

	email, err := h.queries.GetEmail(r.Context(), queries.GetEmailParams{
		ID:             emailUUID,
		OrganizationID: orgUUID,
	})
	if err != nil {
		httputil.WriteError(w, http.StatusNotFound, "email not found")
		return
	}

	_ = h.queries.MarkEmailRead(r.Context(), queries.MarkEmailReadParams{
		ID:             emailUUID,
		OrganizationID: orgUUID,
		Read:           true,
	})

	httputil.WriteJSON(w, http.StatusOK, email)
}

func (h *Handler) DeleteEmail(w http.ResponseWriter, r *http.Request, orgId string, id string) {
	var orgUUID, emailUUID pgtype.UUID
	if err := orgUUID.Scan(orgId); err != nil {
		httputil.WriteError(w, http.StatusBadRequest, "invalid orgId")
		return
	}
	if err := emailUUID.Scan(id); err != nil {
		httputil.WriteError(w, http.StatusBadRequest, "invalid email id")
		return
	}

	if err := h.queries.DeleteEmail(r.Context(), queries.DeleteEmailParams{
		ID:             emailUUID,
		OrganizationID: orgUUID,
	}); err != nil {
		httputil.WriteError(w, http.StatusNotFound, "email not found")
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

func (h *Handler) MarkRead(w http.ResponseWriter, r *http.Request, orgId string, id string) {
	var req api.MarkReadJSONBody
	if err := httputil.DecodeBody(r, &req); err != nil {
		httputil.WriteError(w, http.StatusBadRequest, "invalid request body")
		return
	}

	read := true
	if req.Read != nil {
		read = *req.Read
	}

	var orgUUID, emailUUID pgtype.UUID
	if err := orgUUID.Scan(orgId); err != nil {
		httputil.WriteError(w, http.StatusBadRequest, "invalid orgId")
		return
	}
	if err := emailUUID.Scan(id); err != nil {
		httputil.WriteError(w, http.StatusBadRequest, "invalid email id")
		return
	}

	if err := h.queries.MarkEmailRead(r.Context(), queries.MarkEmailReadParams{
		ID:             emailUUID,
		OrganizationID: orgUUID,
		Read:           read,
	}); err != nil {
		httputil.WriteError(w, http.StatusNotFound, "email not found")
		return
	}

	httputil.WriteJSON(w, http.StatusOK, map[string]bool{"read": read})
}

func (h *Handler) UnreadCount(w http.ResponseWriter, r *http.Request, orgId string) {
	var orgUUID pgtype.UUID
	if err := orgUUID.Scan(orgId); err != nil {
		httputil.WriteError(w, http.StatusBadRequest, "invalid orgId")
		return
	}

	count, err := h.queries.UnreadCount(r.Context(), orgUUID)
	if err != nil {
		httputil.WriteError(w, http.StatusInternalServerError, "failed to get unread count")
		return
	}

	httputil.WriteJSON(w, http.StatusOK, map[string]int64{"count": count})
}

func (h *Handler) ListOrganizations(w http.ResponseWriter, r *http.Request) {
	userID := middleware.UserIDFromCtx(r.Context())
	isAdmin := middleware.IsAdminFromCtx(r.Context())

	var rows []OrgRow

	if isAdmin {
		orgs, err := h.queries.ListOrganizationsByAdmin(r.Context())
		if err != nil {
			httputil.WriteError(w, http.StatusInternalServerError, "failed to list organizations")
			return
		}
		for _, o := range orgs {
			rows = append(rows, OrgRow{
				ID:   uuidStr(o.ID),
				Name: o.Name,
				Slug: o.Slug,
			})
		}
	} else {
		var uid pgtype.UUID
		if err := uid.Scan(userID); err != nil {
			httputil.WriteError(w, http.StatusBadRequest, "invalid user")
			return
		}
		orgs, err := h.queries.ListOrganizationsByUser(r.Context(), uid)
		if err != nil {
			httputil.WriteError(w, http.StatusInternalServerError, "failed to list organizations")
			return
		}
		for _, o := range orgs {
			rows = append(rows, OrgRow{
				ID:   uuidStr(o.ID),
				Name: o.Name,
				Slug: o.Slug,
				Role: o.Role,
			})
		}
	}

	httputil.WriteJSON(w, http.StatusOK, rows)
}
