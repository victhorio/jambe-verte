package main

import (
	"context"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/victhorio/jambe-verte/internal/cache"
	"github.com/victhorio/jambe-verte/internal/content"
	"github.com/victhorio/jambe-verte/internal/handlers"
	"github.com/victhorio/jambe-verte/internal/logger"
	mymiddleware "github.com/victhorio/jambe-verte/internal/middleware"
)

func main() {
	ctx := context.Background()

	// Check if debug mode is enabled
	debugMode := os.Getenv("JV_DEBUG") == "1"
	if debugMode {
		logger.Logger.WarnContext(ctx, "===== Debug mode enabled =====")
	}

	// Load posts
	posts, err := content.LoadContent("content/posts", true)
	if err != nil {
		logger.Logger.ErrorContext(ctx, "Error loading posts", "error", err)
		os.Exit(1)
	}

	// Load pages
	pages, err := content.LoadContent("content/pages", false)
	if err != nil {
		logger.Logger.ErrorContext(ctx, "Error loading pages", "error", err)
		os.Exit(1)
	}

	// Create cache
	c := cache.New(posts, pages)

	// Create handlers
	h, err := handlers.New(c, debugMode)
	if err != nil {
		logger.Logger.ErrorContext(ctx, "Error parsing templates", "error", err)
		os.Exit(1)
	}

	// Setup routes
	r := chi.NewRouter()

	// Middleware
	r.Use(middleware.RealIP)
	r.Use(middleware.RequestID)
	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)
	r.Use(middleware.Timeout(60 * time.Second))
	r.Use(middleware.Compress(5))

	// Routes
	r.Get("/", h.Home)
	r.Get("/posts", h.ListPosts)
	r.Get("/blog/{slug}", h.ShowPost)
	r.Get("/tag/{tag}", h.PostsByTag)
	r.Get("/feed.xml", h.RSSFeed)
	r.Get("/{page}", h.ShowPage)

	// Protected admin routes
	r.Route("/admin", func(r chi.Router) {
		if !debugMode {
			r.Use(mymiddleware.AdminAuth)
		}
		r.Post("/refresh", h.AdminRefresh)
	})

	// Static files
	fileServer := http.FileServer(http.Dir("static"))
	r.Handle("/static/*", http.StripPrefix("/static/", fileServer))

	// Rebuild CSS on startup for better DX
	content.RebuildCSS(ctx)

	// Start server with timeouts
	addr := ":8080"
	srv := &http.Server{
		Addr:         addr,
		Handler:      r,
		ReadTimeout:  15 * time.Second,
		WriteTimeout: 15 * time.Second,
		IdleTimeout:  60 * time.Second,
	}

	// Start server in a goroutine
	go func() {
		logger.Logger.InfoContext(ctx, "Starting server", "address", addr)
		if err := srv.ListenAndServe(); err != http.ErrServerClosed {
			logger.Logger.ErrorContext(ctx, "Server error", "error", err)
			os.Exit(1)
		}
	}()

	// Wait for interrupt signal
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	logger.Logger.InfoContext(ctx, "Shutting down server...")

	// Give outstanding requests 10 seconds to complete
	shutdownCtx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	if err := srv.Shutdown(shutdownCtx); err != nil {
		logger.Logger.ErrorContext(ctx, "Server forced to shutdown", "error", err)
		os.Exit(1)
	}

	logger.Logger.InfoContext(ctx, "Server stopped gracefully")
}
