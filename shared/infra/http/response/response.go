package response

import (
	"encoding/json"
	"fmt"
	"net/http"
)

func Ok(w http.ResponseWriter, response any) {
	JSONResponse(w, http.StatusOK, response)
}

func Created(w http.ResponseWriter, response any) {
	JSONResponse(w, http.StatusCreated, response)
}

type apiError interface {
	StatusCode() int
	Response() ([]byte, error)
}

func APIError(w http.ResponseWriter, apiError apiError) {
	jsonResponse(w, apiError.StatusCode(), apiError.Response)
}

func JSONResponse(w http.ResponseWriter, statusCode int, response any) {
	jsonResponse(w, statusCode, responseMarshaler(response))
}

func jsonResponse(w http.ResponseWriter, statusCode int, response func() ([]byte, error)) {
	code := statusCode

	w.Header().Set("Content-Type", "application/json")

	jsonResp, err := response()
	if err != nil {
		code = http.StatusInternalServerError
		jsonResp = rawErrorResponse("marshal response error")
	}

	w.WriteHeader(code)
	_, _ = w.Write(jsonResp) //nolint: errcheck // no error here
}

func responseMarshaler(response any) func() ([]byte, error) {
	return func() ([]byte, error) {
		return json.Marshal(response)
	}
}

func rawErrorResponse(msg string) []byte {
	return []byte(fmt.Sprintf(`{"message":"%s"}`, msg))
}
