package middleware

import (
	"crypto/subtle"
	"net/http"
	"os"
	"strings"

	"github.com/victhorio/jambe-verte/internal"
	"github.com/victhorio/jambe-verte/internal/logger"
)

var adminToken string

func init() {
	adminToken = os.Getenv("JV_ADMIN_TOKEN")
}

func AdminAuth(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		log := logger.WithRequest(r.Context())

		if adminToken == "" {
			log.Error("JV_ADMIN_TOKEN environment variable not set")
			internal.WriteInternalError(w, "JVE-IMA-MT")
			return
		}

		authHeader := r.Header.Get("Authorization")
		if authHeader == "" {
			log.Warn("Missing Authorization header", "path", r.URL.Path)
			http.Error(w, "Unauthorized", http.StatusUnauthorized)
			return
		}

		if !strings.HasPrefix(authHeader, "Bearer ") {
			log.Warn("Invalid Authorization header format", "path", r.URL.Path)
			http.Error(w, "Unauthorized", http.StatusUnauthorized)
			return
		}

		token := strings.TrimPrefix(authHeader, "Bearer ")
		if subtle.ConstantTimeCompare([]byte(token), []byte(adminToken)) != 1 {
			log.Warn("Invalid bearer token", "path", r.URL.Path)
			http.Error(w, "Unauthorized", http.StatusUnauthorized)
			return
		}

		log.Info("Admin access granted", "path", r.URL.Path)
		next.ServeHTTP(w, r)
	})
}
