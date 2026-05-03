package middleware

import (
	"context"
	"net/http"

	"github.com/google/uuid"
	"github.com/omen77796/go-users-api/internal/common"
)

func RequestID(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

		requestID := uuid.New().String()

		ctx := context.WithValue(r.Context(), common.RequestIDKey, requestID)

		w.Header().Set("X-Request-ID", requestID)

		next.ServeHTTP(w, r.WithContext(ctx))
	})
}
