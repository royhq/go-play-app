package ping_test

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/royhq/go-play-app/use_cases/ping"
)

func TestEndpointHandler_ServeHTTP(t *testing.T) {
	// GIVEN
	handler := ping.NewEndpointHandler()

	// WHEN
	rec := httptest.NewRecorder()
	handler.ServeHTTP(rec, httptest.NewRequest(http.MethodGet, "/ping", http.NoBody))

	// THEN
	assert.Equal(t, http.StatusOK, rec.Code)
	assert.Equal(t, "text/plain", rec.Header().Get("Content-Type"))
	assert.Equal(t, `pong`, rec.Body.String())
}
