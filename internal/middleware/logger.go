package middleware

import (
	"net/http"
	"time"

	"github.com/omen77796/go-users-api/internal/logger"
	"github.com/omen77796/go-users-api/internal/utils"
	"go.uber.org/zap"
)

func Logger(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

		start := time.Now()

		rr := r

		next.ServeHTTP(w, rr)

		logger.Log.Info("http request",
			zap.String("request_id", utils.GetRequestID(rr)),
			zap.String("method", rr.Method),
			zap.String("url", rr.RequestURI),
			zap.Duration("duration", time.Since(start)),
		)
	})
}
