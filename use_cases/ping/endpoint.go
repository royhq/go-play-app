package ping

import (
	"net/http"
)

type EndpointHandler struct{}

func (h *EndpointHandler) ServeHTTP(w http.ResponseWriter, _ *http.Request) {
	w.WriteHeader(http.StatusOK)
	w.Header().Set("Content-Type", "text/plain")
	_, _ = w.Write([]byte("pong")) //nolint: errcheck // no error here
}

func NewEndpointHandler() http.Handler {
	return &EndpointHandler{}
}
