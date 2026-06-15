package handler

import (
	"github.com/chris/MailFyke/internal/httputil"
	"net/http"
)

func (h *Handler) Health(w http.ResponseWriter, r *http.Request) {
	httputil.WriteJSON(w, http.StatusOK, map[string]string{"status": "ok"})
}
