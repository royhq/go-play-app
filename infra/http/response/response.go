package response

import (
	"encoding/json"
	"net/http"
)

type errorResponse struct {
	Message string `json:"message"`
}

func Created(w http.ResponseWriter, response any) {
	writeOkResponse(w, http.StatusCreated, response)
}

func BadRequest(w http.ResponseWriter, msg string) {
	writeErrorResponse(w, http.StatusBadRequest, msg)
}

func InternalError(w http.ResponseWriter, msg string) {
	writeErrorResponse(w, http.StatusInternalServerError, msg)
}

func writeOkResponse(w http.ResponseWriter, statusCode int, response any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCode)

	jsonResp, _ := json.Marshal(response)

	_, _ = w.Write(jsonResp)
}

func writeErrorResponse(w http.ResponseWriter, statusCode int, msg string) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCode)

	jsonResp, _ := json.Marshal(errorResponse{Message: msg})

	_, _ = w.Write(jsonResp)
}
