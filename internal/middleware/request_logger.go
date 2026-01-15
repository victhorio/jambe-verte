package middleware

import (
	"net/http"
	"strings"
	"time"

	"github.com/go-chi/chi/v5/middleware"
	"github.com/victhorio/jambe-verte/internal/logger"
)

const slowRequestThreshold = time.Second

func RequestLogger(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Skip logging for static files
		if strings.HasPrefix(r.URL.Path, "/static/") {
			next.ServeHTTP(w, r)
			return
		}

		ww := middleware.NewWrapResponseWriter(w, r.ProtoMajor)
		start := time.Now()

		defer func() {
			duration := time.Since(start)
			log := logger.WithRequest(r.Context())
			attrs := []any{
				"method", r.Method,
				"path", r.URL.Path,
				"status", ww.Status(),
				"duration_ms", duration.Milliseconds(),
				"bytes", ww.BytesWritten(),
			}

			switch {
			case ww.Status() >= 500:
				log.Error("request", attrs...)
			case ww.Status() >= 400:
				log.Warn("request", attrs...)
			case duration >= slowRequestThreshold:
				log.Warn("slow request", attrs...)
			default:
				log.Info("request", attrs...)
			}
		}()

		next.ServeHTTP(ww, r)
	})
}
