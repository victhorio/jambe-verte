package logger

import (
	"log/slog"
	"os"
)

var Logger *slog.Logger

func init() {
	// Create a structured logger with JSON output for production-like behavior
	Logger = slog.New(slog.NewJSONHandler(os.Stderr, &slog.HandlerOptions{
		Level: slog.LevelInfo,
	}))
}
