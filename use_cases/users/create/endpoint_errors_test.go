package create_test

import (
	"bytes"
	"context"
	"errors"
	"fmt"
	"io"
	"log/slog"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/royhq/go-play-app/use_cases/users/create"
)

func TestEndpointErrorHandler_HandleError(t *testing.T) {
	t.Parallel()

	t.Run("error handling", testErrorHandling)

	t.Run("should log error handled", func(t *testing.T) {
		// GIVEN
		logger, buf := buffLogger()
		handler := create.NewEndpointErrorHandler(logger)

		err1 := errors.New("error 1")
		err2 := fmt.Errorf("error 2: %w", err1)

		// WHEN
		rec := httptest.NewRecorder()
		handler.HandleError(context.Background(), rec, err2)

		// THEN
		expected :=
			`{"time":"((?:(\d{4}-\d{2}-\d{2})T(\d{2}:\d{2}:\d{2}(?:\.\d+)?))(Z|[\+-]\d{2}:\d{2})?)","level":"ERROR","msg":"error handled","error":"error 2: error 1","response":"{\\"code\\":\\"internal_error\\",\\"message\\":\\"unexpected error\\"}","status_code":500,"stack":\["error 2: error 1","error 1"\]}` // nolint:lll

		assert.Regexp(t, expected, buf.String())
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
		"command error should return internal error": {
			cmdErr: &create.CommandError{
				Msg:   "test error",
				Code:  "test_code",
				Cause: errors.New("cause error"),
			},
			expectedStatusCode: http.StatusInternalServerError,
			expectedResponse:   `{"code":"test_code","message":"test error"}`,
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
