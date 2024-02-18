package ping

import (
	"net/http"
)

type Handler struct{}

func (h *Handler) ServeHTTP(w http.ResponseWriter, _ *http.Request) {
	w.WriteHeader(http.StatusOK)
	w.Header().Set("Content-Type", "application/json")
	_, _ = w.Write([]byte("pong"))
}

func NewHandler() http.Handler {
	h := &Handler{}
	return h
}
