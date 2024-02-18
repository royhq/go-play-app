package response_test

import (
	"io"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go-play-app/infra/http/response"
)

func TestCreated(t *testing.T) {
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
	// WHEN
	rec := httptest.NewRecorder()
	response.BadRequest(rec, "something went wrong")

	// THEN
	result := rec.Result()
	assert.Equal(t, http.StatusBadRequest, result.StatusCode)

	respBody, err := io.ReadAll(result.Body)
	require.NoError(t, err)
	assert.JSONEq(t, `{"message":"something went wrong"}`, string(respBody))
}

func TestInternalError(t *testing.T) {
	// WHEN
	rec := httptest.NewRecorder()
	response.InternalError(rec, "something went wrong")

	// THEN
	result := rec.Result()
	assert.Equal(t, http.StatusInternalServerError, result.StatusCode)

	respBody, err := io.ReadAll(result.Body)
	require.NoError(t, err)
	assert.JSONEq(t, `{"message":"something went wrong"}`, string(respBody))
}
