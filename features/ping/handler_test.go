package ping_test

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
	"go-play-app/features/ping"
)

func TestHandler_ServeHTTP(t *testing.T) {
	// GIVEN
	handler := ping.NewHandler()

	// WHEN
	rec := httptest.NewRecorder()
	handler.ServeHTTP(rec, httptest.NewRequest(http.MethodGet, "/ping", http.NoBody))

	// THEN
	assert.Equal(t, http.StatusOK, rec.Code)
	assert.Equal(t, "application/json", rec.Header().Get("Content-Type"))
	assert.Equal(t, `pong`, rec.Body.String())
}
