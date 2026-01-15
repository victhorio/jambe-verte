package logger

import (
	"context"
	"log/slog"
	"os"
	"strings"

	"github.com/go-chi/chi/v5/middleware"
	"github.com/lmittmann/tint"
)

var Logger *slog.Logger

func init() {
	Logger = slog.New(tint.NewHandler(os.Stderr, &tint.Options{
		Level: slog.LevelInfo,
	}))
}

func Init(debug bool) {
	level := slog.LevelInfo
	if debug {
		level = slog.LevelDebug
	}

	if envLevel := os.Getenv("JV_LOG_LEVEL"); envLevel != "" {
		switch strings.ToLower(envLevel) {
		case "debug":
			level = slog.LevelDebug
		case "info":
			level = slog.LevelInfo
		case "warn":
			level = slog.LevelWarn
		case "error":
			level = slog.LevelError
		}
	}

	var handler slog.Handler
	if debug {
		handler = tint.NewHandler(os.Stderr, &tint.Options{Level: level})
	} else {
		handler = slog.NewJSONHandler(os.Stderr, &slog.HandlerOptions{Level: level})
	}
	Logger = slog.New(handler)
}

func WithRequest(ctx context.Context) *slog.Logger {
	reqID := middleware.GetReqID(ctx)
	if len(reqID) > 8 {
		reqID = reqID[:8]
	}
	return Logger.With("request_id", reqID)
}
