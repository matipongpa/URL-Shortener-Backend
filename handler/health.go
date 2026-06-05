package handler

import (
	"net/http"
)

func (h *Handler) Health(w http.ResponseWriter, r *http.Request) {
	data := map[string]string{"status": "ok"}
	Response(w, http.StatusOK, data)
}
