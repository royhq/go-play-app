package create_test

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"log/slog"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go-play-app/features/users/create"
)

func TestEndpointErrorHandler_HandleError(t *testing.T) {
	t.Parallel()

	t.Run("error handling", testErrorHandling)

	t.Run("should log error handled", func(t *testing.T) {
		// GIVEN
		logger, buff := buffLogger()
		handler := create.NewEndpointErrorHandler(logger)

		err1 := errors.New("error 1")
		err2 := fmt.Errorf("error 2: %w", err1)

		// WHEN
		rec := httptest.NewRecorder()
		handler.HandleError(context.Background(), rec, err2)

		// THEN
		var fields map[string]any

		err := json.Unmarshal(buff.Bytes(), &fields)
		require.NoError(t, err)

		assert.Equal(t, "ERROR", fields["level"])            // log level
		assert.Equal(t, "error handled", fields["msg"])      // log message
		assert.Equal(t, "error 2: error 1", fields["error"]) // log error

		// log response
		assert.Equal(t, "internal_error", fields["response"].(map[string]any)["code"])
		assert.Equal(t, "unexpected error", fields["response"].(map[string]any)["message"])
		assert.Equal(t, float64(http.StatusInternalServerError), fields["status_code"])

		// log error stack
		stack, ok := fields["stack"].([]any)
		require.True(t, ok)

		assert.Len(t, stack, 2)
		assert.Equal(t, "error 2: error 1", stack[0])
		assert.Equal(t, "error 1", stack[1])
	})
}

func testErrorHandling(t *testing.T) {
	t.Parallel()

	testCases := map[string]struct {
		cmdErr             error
		expectedStatusCode int
		expectedResponse   string
	}{
		"validation error should return bad request": {
			cmdErr:             create.NewValidationError("some data is not valid"),
			expectedStatusCode: http.StatusBadRequest,
			expectedResponse:   `{"code":"validation_error","message":"some data is not valid"}`,
		},
		"another error should return internal error": {
			cmdErr:             errors.New("something went wrong"),
			expectedStatusCode: http.StatusInternalServerError,
			expectedResponse:   `{"code":"internal_error","message":"unexpected error"}`,
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

func buffLogger() (*slog.Logger, *bytes.Buffer) {
	buff := new(bytes.Buffer)
	h := slog.NewJSONHandler(buff, nil) // use json handler to easily unmarshal record to map[string]any

	return slog.New(h), buff
}
