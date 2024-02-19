package middleware_test

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/royhq/go-play-app/infra/http/middleware"
)

func TestWithRecover(t *testing.T) {
	t.Parallel()

	t.Run("when handler panics should return internal error", func(t *testing.T) {
		t.Parallel()

		// GIVEN
		panicHandler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			panic("something went wrong")
		})

		handler := middleware.WithRecover(panicHandler)

		// WHEN
		rec := httptest.NewRecorder()
		handler.ServeHTTP(rec, httptest.NewRequest(http.MethodGet, "/test", http.NoBody))

		// THEN
		assert.Equal(t, http.StatusInternalServerError, rec.Code)
		assert.Equal(t, "application/json", rec.Header().Get("Content-Type"))
		assert.JSONEq(t, `{"message":"internal error"}`, rec.Body.String())
	})

	t.Run("when handler does not panic return response", func(t *testing.T) {
		t.Parallel()

		// GIVEN
		var (
			statusCode = http.StatusOK
			response   = `{"message":"successfully response"}`
		)

		successHandler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(statusCode)
			_, _ = w.Write([]byte(response))
		})

		handler := middleware.WithRecover(successHandler)

		// WHEN
		rec := httptest.NewRecorder()
		handler.ServeHTTP(rec, httptest.NewRequest(http.MethodGet, "/test", http.NoBody))

		// THEN
		assert.Equal(t, statusCode, rec.Code)
		assert.Equal(t, "application/json", rec.Header().Get("Content-Type"))
		assert.JSONEq(t, response, rec.Body.String())
	})
}
