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

func BadRequest(w http.ResponseWriter, response any) {
	JSONResponse(w, http.StatusBadRequest, response)
}

func InternalError(w http.ResponseWriter, response any) {
	JSONResponse(w, http.StatusInternalServerError, response)
}

func JSONResponse(w http.ResponseWriter, statusCode int, response any) {
	code := statusCode

	w.Header().Set("Content-Type", "application/json")

	jsonResp, err := json.Marshal(response)
	if err != nil {
		code = http.StatusInternalServerError
		jsonResp = rawErrorResponse("unmarshal response error")
	}

	w.WriteHeader(code)
	_, _ = w.Write(jsonResp) //nolint: errcheck // no error here
}

func rawErrorResponse(msg string) []byte {
	return []byte(fmt.Sprintf(`{"message":"%s"}`, msg))
}
