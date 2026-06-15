package handler

import (
	"github.com/chris/MailFyke/internal/httputil"
	"net/http"
)

func (h *Handler) AuthLogin(w http.ResponseWriter, r *http.Request) {
	httputil.WriteJSON(w, http.StatusNotImplemented, map[string]string{"error": "not implemented"})
}
