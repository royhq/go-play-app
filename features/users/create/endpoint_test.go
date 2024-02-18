package create_test

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"go-play-app/features/users/create"
)

func TestEndpointHandler_ServeHTTP(t *testing.T) {
	t.Parallel()

	t.Run("when handle command successfully should return 201", func(t *testing.T) {
		t.Parallel()

		// GIVEN
		reqBody := `
		{
			"name":"John Doe",
			"age": 32
		}`

		cmdHandlerFunc := create.CommandHandlerFunc(func(_ context.Context, cmd create.Command) (create.CommandOutput, error) {
			return create.CommandOutput{
				CreatedUser: create.User{
					ID:      "d4fd02d8-ec7e-4846-bd02-d8ec7e38462f",
					Name:    "John Doe",
					Age:     32,
					Created: time.Date(2023, time.March, 20, 14, 30, 12, 0, time.UTC),
				},
			}, nil
		})

		handler := create.NewEndpointHandler(cmdHandlerFunc, nil)

		// WHEN
		rec := httptest.NewRecorder()
		req := httptest.NewRequest(http.MethodPost, "/users", strings.NewReader(reqBody))
		handler.ServeHTTP(rec, req)

		// THEN
		assert.Equal(t, http.StatusCreated, rec.Code)

		expectedResponse := `
		{
			"id":"d4fd02d8-ec7e-4846-bd02-d8ec7e38462f",
			"name":"John Doe",
			"age":32,
			"created_at":"2023-03-20T14:30:12Z"
		}`
		assert.JSONEq(t, expectedResponse, rec.Body.String())
	})

	t.Run("bad request should be handled", func(t *testing.T) {
		t.Parallel()

		// GIVEN
		reqBody := `{xxx}` // bad json input
		handler := create.NewEndpointHandler(nil, nil)

		// WHEN
		rec := httptest.NewRecorder()
		req := httptest.NewRequest(http.MethodPost, "/users", strings.NewReader(reqBody))
		handler.ServeHTTP(rec, req)

		// THEN
		assert.Equal(t, http.StatusBadRequest, rec.Code)
		assert.Equal(t, "application/json", rec.Header().Get("Content-Type"))
		assert.JSONEq(t, `{"code":"bad_request","message":"bad request"}`, rec.Body.String())
	})

	t.Run("when handle command with error should be handled", func(t *testing.T) {
		t.Parallel()

		// GIVEN
		reqBody := `{"name":"John Doe","age":32}`

		cmdHandlerFunc := create.CommandHandlerFunc(func(_ context.Context, _ create.Command) (create.CommandOutput, error) {
			return create.CommandOutput{}, errors.New("something went wrong")
		})

		errHandler := errorHandlerFunc(func(_ context.Context, w http.ResponseWriter, e error) {
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusInternalServerError)

			resp := fmt.Sprintf(`{"message":"%s"}`, e)
			_, _ = w.Write([]byte(resp))
		})

		handler := create.NewEndpointHandler(cmdHandlerFunc, errHandler)

		// WHEN
		rec := httptest.NewRecorder()
		req := httptest.NewRequest(http.MethodPost, "/users", strings.NewReader(reqBody))
		handler.ServeHTTP(rec, req)

		// THEN
		assert.Equal(t, http.StatusInternalServerError, rec.Code)
		assert.Equal(t, "application/json", rec.Header().Get("Content-Type"))
		assert.JSONEq(t, `{"message":"something went wrong"}`, rec.Body.String())
	})
}

type errorHandlerFunc func(context.Context, http.ResponseWriter, error)

func (f errorHandlerFunc) HandleError(ctx context.Context, w http.ResponseWriter, e error) {
	f(ctx, w, e)
}
