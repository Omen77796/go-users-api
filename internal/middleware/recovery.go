package middleware

import (
	"log"
	"net/http"

	"github.com/omen77796/go-users-api/internal/utils"
)

func Recoverer(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		defer func() {
			if err := recover(); err != nil {
				log.Println("PANIC RECOVERED:", err)
				utils.Error(w, http.StatusInternalServerError, "internal server error")
			}
		}()

		next.ServeHTTP(w, r)
	})
}
