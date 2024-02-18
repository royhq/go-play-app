package response

import (
	"encoding/json"
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
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCode)

	jsonResp, _ := json.Marshal(response)

	_, _ = w.Write(jsonResp)
}
