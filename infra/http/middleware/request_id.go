package middleware

import (
	"context"
	"net/http"

	"go-play-app/infra/uuid"
)

type ctxKey string

const reqIDKey = ctxKey("request_id")

func WithRequestID(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		reqID := uuid.New()

		ctx := r.Context()
		ctx = context.WithValue(ctx, reqIDKey, reqID)
		newReq := r.WithContext(ctx)

		w.Header().Set("X-Request-ID", reqID)

		next.ServeHTTP(w, newReq)
	})
}

func RequestID(ctx context.Context) string {
	if value, ok := ctx.Value(reqIDKey).(string); ok {
		return value
	}

	return ""
}
