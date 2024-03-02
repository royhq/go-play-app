package middleware

import (
	"log"
	"net/http"

	"github.com/royhq/go-play-app/shared/infra/http/response"
)

type recoverError string

func (e *recoverError) StatusCode() int { return http.StatusInternalServerError }

func (e *recoverError) Response() ([]byte, error) {
	return []byte(`{"message":"internal error"}`), nil
}

func (e *recoverError) Error() string { return string(*e) }

func WithRecover(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		defer func() {
			err := recover()
			if err != nil {
				log.Println("recovering from error:", err)
				response.APIError(w, new(recoverError))
			}
		}()

		next.ServeHTTP(w, r)
	})
}
