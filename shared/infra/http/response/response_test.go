package response_test

import (
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
		assert.JSONEq(t, `{"message":"unmarshal response error"}`, rec.Body.String())
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
	assert.Equal(t, http.StatusBadRequest, rec.Code)
	assert.JSONEq(t, `{"code":"bad_request","message":"something went wrong"}`, rec.Body.String())
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
	assert.Equal(t, http.StatusInternalServerError, rec.Code)
	assert.JSONEq(t, `{"code":"internal_error","message":"something went wrong"}`, rec.Body.String())
}

type testErrorResponse struct {
	Code    string `json:"code"`
	Message string `json:"message"`
}

type failedMarshal struct{}

func (f failedMarshal) MarshalJSON() ([]byte, error) {
	return nil, fmt.Errorf("marshal error")
}
