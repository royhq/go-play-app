package endpoints

import (
	"encoding/json"
	"log/slog"
)

type APIError struct {
	HTTPStatusCode int    `json:"-"`
	Code           string `json:"code"`
	Msg            string `json:"message"`
}

func (e *APIError) LogValue() slog.Value {
	v, _ := e.Response()
	return slog.StringValue(string(v))
}

func (e *APIError) StatusCode() int           { return e.HTTPStatusCode }
func (e *APIError) Response() ([]byte, error) { return json.Marshal(e) }
func (e *APIError) Error() string             { return e.Msg }
