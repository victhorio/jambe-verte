package main

import (
	"context"
	"fmt"
	"net/http"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/victhor/jv-blog/internal/cache"
	"github.com/victhor/jv-blog/internal/content"
	"github.com/victhor/jv-blog/internal/handlers"
	"github.com/victhor/jv-blog/internal/logger"
)

func main() {
	ctx := context.Background()

	// Load posts
	posts, err := content.LoadPosts("content/posts")
	if err != nil {
		logger.Logger.WarnContext(ctx, "Error loading posts", "error", err)
		posts = []*content.Post{} // Continue with empty posts
	}

	// Load pages
	pages, err := content.LoadPages("content/pages")
	if err != nil {
		logger.Logger.WarnContext(ctx, "Error loading pages", "error", err)
		pages = []*content.Post{} // Continue with empty pages
	}

	// Create cache
	c := cache.New(posts, pages)

	// Create handlers
	h := handlers.New(c)

	// Setup routes
	r := chi.NewRouter()

	// Middleware
	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)
	r.Use(middleware.Compress(5))
	r.Use(middleware.Timeout(60 * time.Second))

	// Routes
	r.Get("/", h.Home)
	r.Get("/blog", h.ListPosts)
	r.Get("/blog/{slug}", h.ShowPost)
	r.Get("/tag/{tag}", h.PostsByTag)
	r.Get("/feed.xml", h.RSSFeed)
	r.Get("/test-error", h.TestError) // Test endpoint for error handling
	r.Get("/{page}", h.ShowPage)

	// Static files
	fileServer := http.FileServer(http.Dir("static"))
	r.Handle("/static/*", http.StripPrefix("/static/", fileServer))

	// Start server
	addr := ":8080"
	logger.Logger.InfoContext(ctx, "Starting server", "address", addr)
	fmt.Printf("Starting server on http://localhost%s\n", addr)
	if err := http.ListenAndServe(addr, r); err != nil {
		logger.Logger.ErrorContext(ctx, "Server failed to start", "error", err)
		panic(err)
	}
}
