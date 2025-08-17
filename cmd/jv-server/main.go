package main

import (
	"context"
	"net/http"
	"os"
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
	posts, err := content.LoadPosts("content/posts", true)
	if err != nil {
		logger.Logger.ErrorContext(ctx, "Error loading posts", "error", err)
		panic("No posts found")
	}

	// Load pages
	pages, err := content.LoadPosts("content/pages", false)
	if err != nil {
		logger.Logger.ErrorContext(ctx, "Error loading pages", "error", err)
		panic("No pages found")
	}

	// Create cache
	c := cache.New(posts, pages)

	// Create handlers
	h := handlers.New(c, debugMode)

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

	// Start server
	addr := ":8080"
	logger.Logger.InfoContext(ctx, "Starting server", "address", addr)
	if err := http.ListenAndServe(addr, r); err != nil {
		logger.Logger.ErrorContext(ctx, "Server failed to start", "error", err)
		panic(err)
	}
}
