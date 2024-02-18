package middleware

import (
	"log"
	"net/http"

	"go-play-app/infra/http/response"
)

func WithRecover(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		defer func() {
			err := recover()
			if err != nil {
				log.Println("recovering from error:", err)
				response.InternalError(w, "internal error")
			}
		}()

		next.ServeHTTP(w, r)
	})
}
