package create_test

import (
	"context"
	"errors"
	"fmt"
	"io"
	"log/slog"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
	"go-play-app/features/users/create"
)

func TestEndpointErrorHandler_HandleError(t *testing.T) {
	t.Parallel()

	testCases := map[string]struct {
		cmdErr             error
		expectedStatusCode int
		expectedResponse   string
	}{
		"validation error": {
			cmdErr:             create.ErrValidation,
			expectedStatusCode: http.StatusBadRequest,
			expectedResponse:   `{"message":"validation error"}`,
		},
		"validation error wrapped": {
			cmdErr:             fmt.Errorf("something went wrong: %w", create.ErrValidation),
			expectedStatusCode: http.StatusBadRequest,
			expectedResponse:   `{"message":"something went wrong: validation error"}`,
		},
		"another error": {
			cmdErr:             errors.New("something went wrong"),
			expectedStatusCode: http.StatusInternalServerError,
			expectedResponse:   `{"message":"something went wrong"}`,
		},
	}

	for name, tc := range testCases {
		t.Run(name, func(t *testing.T) {
			t.Parallel()
			// GIVEN
			handler := create.NewEndpointErrorHandler(noLogger())

			// WHEN
			rec := httptest.NewRecorder()
			handler.HandleError(context.Background(), rec, tc.cmdErr)

			// THEN
			assert.Equal(t, tc.expectedStatusCode, rec.Code)
			assert.Equal(t, "application/json", rec.Header().Get("Content-Type"))
			assert.JSONEq(t, tc.expectedResponse, rec.Body.String())
		})
	}
}

func noLogger() *slog.Logger {
	h := slog.NewTextHandler(io.Discard, nil)
	return slog.New(h)
}
