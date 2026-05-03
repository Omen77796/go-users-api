package middleware

import (
	"net/http"

	"github.com/omen77796/go-users-api/internal/logger"
	"github.com/omen77796/go-users-api/internal/utils"
	"go.uber.org/zap"
)

func Recovery(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		defer func() {
			if err := recover(); err != nil {
				logger.Log.Error("panic recovered", zap.Any("error", err))
				utils.JSONError(w, http.StatusInternalServerError, "internal server error")
			}
		}()

		next.ServeHTTP(w, r)
	})
}
