package content

import (
	"context"
	"os/exec"

	"github.com/victhorio/jambe-verte/internal/logger"
)

const (
	InternalErrorTemplate = `Internal Server Error
=====================

Contact site administrator.

Error Code: %s`
)

func RebuildCSS(ctx context.Context) {
	logger.Logger.InfoContext(ctx, "Rebuilding CSS...")

	// Check if bun exists
	if _, err := exec.LookPath("bun"); err != nil {
		logger.Logger.WarnContext(ctx, "bun not found in PATH, skipping CSS rebuild")
		return
	}

	// Run bun build-css
	cmd := exec.Command("bun", "run", "build-css")
	if err := cmd.Run(); err != nil {
		logger.Logger.WarnContext(ctx, "CSS rebuild failed", "error", err)
	} else {
		logger.Logger.InfoContext(ctx, "CSS rebuild completed successfully")
	}
}
