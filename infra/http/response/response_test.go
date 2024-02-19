package response_test

import (
	"io"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/royhq/go-play-app/infra/http/response"
)

func TestJSONResponse(t *testing.T) {
	t.Parallel()

	// WHEN
	rec := httptest.NewRecorder()
	response.JSONResponse(rec, http.StatusOK, map[string]any{
		"message": "test message",
	})

	// THEN
	result := rec.Result()
	assert.Equal(t, http.StatusOK, result.StatusCode)

	respBody, err := io.ReadAll(result.Body)
	require.NoError(t, err)
	assert.JSONEq(t, `{"message":"test message"}`, string(respBody))
}

func TestOk(t *testing.T) {
	t.Parallel()

	// GIVEN
	type resp struct {
		AttrString string `json:"attr_string"`
		AttrInt    int    `json:"attr_int"`
		AttrBool   bool   `json:"attr_bool"`
	}

	// WHEN
	rec := httptest.NewRecorder()
	response.Ok(rec, resp{
		AttrString: "a string",
		AttrInt:    123,
		AttrBool:   true,
	})

	// THEN
	result := rec.Result()
	assert.Equal(t, http.StatusOK, result.StatusCode)

	respBody, err := io.ReadAll(result.Body)
	require.NoError(t, err)
	assert.JSONEq(t, `{"attr_string":"a string","attr_int":123,"attr_bool":true}`, string(respBody))
}

func TestCreated(t *testing.T) {
	t.Parallel()

	// GIVEN
	type resp struct {
		AttrString string `json:"attr_string"`
		AttrInt    int    `json:"attr_int"`
		AttrBool   bool   `json:"attr_bool"`
	}

	// WHEN
	rec := httptest.NewRecorder()
	response.Created(rec, resp{
		AttrString: "a string",
		AttrInt:    123,
		AttrBool:   true,
	})

	// THEN
	result := rec.Result()
	assert.Equal(t, http.StatusCreated, result.StatusCode)

	respBody, err := io.ReadAll(result.Body)
	require.NoError(t, err)
	assert.JSONEq(t, `{"attr_string":"a string","attr_int":123,"attr_bool":true}`, string(respBody))
}

func TestBadRequest(t *testing.T) {
	t.Parallel()

	// WHEN
	resp := testErrorResponse{
		Code:    "bad_request",
		Message: "something went wrong",
	}

	rec := httptest.NewRecorder()
	response.BadRequest(rec, resp)

	// THEN
	result := rec.Result()
	assert.Equal(t, http.StatusBadRequest, result.StatusCode)

	respBody, err := io.ReadAll(result.Body)
	require.NoError(t, err)
	assert.JSONEq(t, `{"code":"bad_request","message":"something went wrong"}`, string(respBody))
}

func TestInternalError(t *testing.T) {
	t.Parallel()

	// WHEN
	resp := testErrorResponse{
		Code:    "internal_error",
		Message: "something went wrong",
	}

	rec := httptest.NewRecorder()
	response.InternalError(rec, resp)

	// THEN
	result := rec.Result()
	assert.Equal(t, http.StatusInternalServerError, result.StatusCode)

	respBody, err := io.ReadAll(result.Body)
	require.NoError(t, err)
	assert.JSONEq(t, `{"code":"internal_error","message":"something went wrong"}`, string(respBody))
}

type testErrorResponse struct {
	Code    string `json:"code"`
	Message string `json:"message"`
}
