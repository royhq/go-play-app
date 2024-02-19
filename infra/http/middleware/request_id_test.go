package middleware_test

import (
	"context"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/royhq/go-play-app/infra/http/middleware"
)

func TestWithRequestID(t *testing.T) {
	t.Parallel()

	t.Run("should add request id header to response", func(t *testing.T) {
		t.Parallel()

		// GIVEN
		var reqID string

		handler := middleware.WithRequestID(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			reqID = middleware.RequestID(r.Context())

			w.WriteHeader(http.StatusOK)
			_, _ = w.Write([]byte(`response ok`))
		}))

		// WHEN
		rec := httptest.NewRecorder()
		handler.ServeHTTP(rec, httptest.NewRequest(http.MethodGet, "/test", http.NoBody))

		// THEN
		assert.Equal(t, http.StatusOK, rec.Code)
		assert.Equal(t, reqID, rec.Header().Get("X-Request-ID"))
		assert.Equal(t, `response ok`, rec.Body.String())
	})
}

func TestRequestID(t *testing.T) {
	t.Parallel()

	t.Run("when request id is not in context should return empty", func(t *testing.T) {
		reqID := middleware.RequestID(context.Background())
		assert.Empty(t, reqID)
	})
}
