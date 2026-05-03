package utils

import (
	"net/http"

	"github.com/omen77796/go-users-api/internal/common"
)

func GetRequestID(r *http.Request) string {
	if id, ok := r.Context().Value(common.RequestIDKey).(string); ok {
		return id
	}
	return "unknown"
}
