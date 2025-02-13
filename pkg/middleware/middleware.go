package middleware

import (
	"log"
	"net/http"
	"strings"
	"time"

	token "fr_lab_1/pkg/token"
)

func AuthMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if strings.Contains(r.URL.Path, "/static/") {
			next.ServeHTTP(w, r)
			return
		}

		user_token := r.Header.Get("Authorization")
		if user_token == "" {
			r.Header.Set("IsAuthorized", "false")
		}

		// Check if token is valid
		if token.CheckTokenExists(user_token) {
			r.Header.Set("IsAuthorized", "true")
		} else {
			r.Header.Set("IsAuthorized", "false")
		}

		authority := r.Header.Get("Authority")
		if authority == "" {
			r.Header.Set("Authority", "user")
		}

		next.ServeHTTP(w, r)
	})
}

func LoggingMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()
		log.Printf("Started %s %s", r.Method, r.URL.Path)

		next.ServeHTTP(w, r)

		log.Printf("Completed %s in %v", r.URL.Path, time.Since(start))
	})
}
