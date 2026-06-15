package handler

import (
	"github.com/chris/MailFyke/internal/httputil"
	"net/http"
)

func (h *Handler) ListEmails(w http.ResponseWriter, r *http.Request) {
	httputil.WriteJSON(w, http.StatusNotImplemented, map[string]string{"error": "not implemented"})
}

func (h *Handler) GetEmail(w http.ResponseWriter, r *http.Request) {
	httputil.WriteJSON(w, http.StatusNotImplemented, map[string]string{"error": "not implemented"})
}

func (h *Handler) DeleteEmail(w http.ResponseWriter, r *http.Request) {
	httputil.WriteJSON(w, http.StatusNotImplemented, map[string]string{"error": "not implemented"})
}

func (h *Handler) MarkRead(w http.ResponseWriter, r *http.Request) {
	httputil.WriteJSON(w, http.StatusNotImplemented, map[string]string{"error": "not implemented"})
}

func (h *Handler) UnreadCount(w http.ResponseWriter, r *http.Request) {
	httputil.WriteJSON(w, http.StatusNotImplemented, map[string]string{"error": "not implemented"})
}

func (h *Handler) ListOrganizations(w http.ResponseWriter, r *http.Request) {
	httputil.WriteJSON(w, http.StatusNotImplemented, map[string]string{"error": "not implemented"})
}
