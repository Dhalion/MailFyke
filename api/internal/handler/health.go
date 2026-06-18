package handler

import (
	"net/http"

	"github.com/Dhalion/MailFyke/internal/httputil"
)

func (h *Handler) Health(w http.ResponseWriter, r *http.Request) {
	httputil.WriteJSON(w, http.StatusOK, map[string]string{"status": "ok"})
}
