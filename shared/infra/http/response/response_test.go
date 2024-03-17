package response_test

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/royhq/go-play-app/shared/infra/http/response"
)

func TestJSONResponse(t *testing.T) {
	t.Parallel()

	t.Run("success marshal response", func(t *testing.T) {
		// WHEN
		rec := httptest.NewRecorder()
		response.JSONResponse(rec, http.StatusOK, map[string]any{
			"message": "test message",
		})

		// THEN
		assert.Equal(t, http.StatusOK, rec.Code)
		assert.JSONEq(t, `{"message":"test message"}`, rec.Body.String())
	})

	t.Run("marshal response error", func(t *testing.T) {
		// WHEN
		rec := httptest.NewRecorder()
		response.JSONResponse(rec, http.StatusInternalServerError, failedMarshal{})

		// THEN
		assert.Equal(t, http.StatusInternalServerError, rec.Code)
		assert.JSONEq(t, `{"message":"marshal response error"}`, rec.Body.String())
	})
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
	assert.Equal(t, http.StatusOK, rec.Code)
	assert.JSONEq(t, `{"attr_string":"a string","attr_int":123,"attr_bool":true}`, rec.Body.String())
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
	assert.Equal(t, http.StatusCreated, rec.Code)
	assert.JSONEq(t, `{"attr_string":"a string","attr_int":123,"attr_bool":true}`, rec.Body.String())
}

func TestAPIError(t *testing.T) {
	t.Parallel()

	// GIVEN
	apiErr := testAPIError{
		HTTPStatusCode: http.StatusTeapot,
		AttrString:     "a string",
		AttrInt:        123,
		AttrBool:       true,
	}

	// WHEN
	rec := httptest.NewRecorder()
	response.APIError(rec, apiErr)

	// THEN
	assert.Equal(t, http.StatusTeapot, rec.Code)
	assert.JSONEq(t, `{"attr_string":"a string","attr_int":123,"attr_bool":true}`, rec.Body.String())
}

type testAPIError struct {
	HTTPStatusCode int    `json:"-"`
	AttrString     string `json:"attr_string"`
	AttrInt        int    `json:"attr_int"`
	AttrBool       bool   `json:"attr_bool"`
}

func (e testAPIError) StatusCode() int {
	return e.HTTPStatusCode
}

func (e testAPIError) Response() ([]byte, error) {
	return json.Marshal(e)
}

type failedMarshal struct{}

func (f failedMarshal) MarshalJSON() ([]byte, error) {
	return nil, fmt.Errorf("marshal error")
}
