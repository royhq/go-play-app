package ping

import (
	"net/http"
)

type EndpointHandler struct{}

func (h *EndpointHandler) ServeHTTP(w http.ResponseWriter, _ *http.Request) {
	w.WriteHeader(http.StatusOK)
	w.Header().Set("Content-Type", "application/json")
	_, _ = w.Write([]byte("pong"))
}

func NewEndpointHandler() http.Handler {
	return &EndpointHandler{}
}
