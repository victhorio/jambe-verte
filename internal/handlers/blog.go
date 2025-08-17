package handlers

import (
	"bytes"
	"context"
	"fmt"
	"html/template"
	"net/http"
	"sync"
	"time"

	// "github.com/go-chi/chi/v5"
	"github.com/victhorio/jambe-verte/internal/cache"
	"github.com/victhorio/jambe-verte/internal/content"
	"github.com/victhorio/jambe-verte/internal/logger"
)

type Handler struct {
	mu    sync.RWMutex
	cache *cache.Cache
}

func New(cache *cache.Cache) *Handler {
	return &Handler{
		cache: cache,
	}
}

func (h *Handler) getCache() *cache.Cache {
	h.mu.RLock()
	defer h.mu.RUnlock()
	return h.cache
}

func (h *Handler) setCache(cache *cache.Cache) {
	h.mu.Lock()
	defer h.mu.Unlock()
	h.cache = cache
}

func (h *Handler) Home(w http.ResponseWriter, r *http.Request) {
	// Check page cache first
	cache := h.getCache()
	pageCache := cache.GetPageCache()
	if cached, ok := pageCache.Get("/"); ok {
		w.Header().Set("Content-Type", "text/html; charset=utf-8")
		w.Write(cached)
		return
	}

	files := []string{
		"templates/base.html",
		"templates/home.html",
	}
	h.renderAndCache(r.Context(), w, "/", files, nil)
}

func (h *Handler) ListPosts(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/plain; charset=utf-8")
	w.Write([]byte("TODO: Posts listing"))
}

func (h *Handler) ShowPost(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/plain; charset=utf-8")
	w.Write([]byte("TODO: Individual post view"))
}

func (h *Handler) ShowPage(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/plain; charset=utf-8")
	w.Write([]byte("TODO: Static page view"))
}

func (h *Handler) PostsByTag(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/plain; charset=utf-8")
	w.Write([]byte("TODO: Posts by tag view"))
}

// AdminRefresh is responsible for hot-reloading content by creating an entirely new cache
// and replacing it on the handler.
func (h *Handler) AdminRefresh(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	// Load posts
	posts, err := content.LoadPosts("content/posts", true)
	if err != nil {
		logger.Logger.ErrorContext(ctx, "Error loading posts during refresh", "error", err)
		w.Header().Set("Content-Type", "text/plain; charset=utf-8")
		http.Error(w, "Failed to load posts: "+err.Error(), http.StatusInternalServerError)
		return
	}

	// Load pages
	pages, err := content.LoadPosts("content/pages", false)
	if err != nil {
		logger.Logger.ErrorContext(ctx, "Error loading pages during refresh", "error", err)
		w.Header().Set("Content-Type", "text/plain; charset=utf-8")
		http.Error(w, "Failed to load pages: "+err.Error(), http.StatusInternalServerError)
		return
	}

	// Create new cache and replace the old one
	newCache := cache.New(posts, pages)
	h.setCache(newCache)

	logger.Logger.InfoContext(ctx, "Cache refreshed successfully", "posts", len(posts), "pages", len(pages))

	// Also attempt to rebuild CSS
	content.RebuildCSS(ctx)

	w.Header().Set("Content-Type", "text/plain; charset=utf-8")
	w.Write([]byte("OK"))
}

func (h *Handler) renderAndCache(ctx context.Context, w http.ResponseWriter, route string, templateFiles []string, data any) {
	// Check if context is cancelled
	select {
	case <-ctx.Done():
		logger.Logger.WarnContext(ctx, "Request cancelled", "error", ctx.Err())
		return
	default:
	}

	startTime := time.Now()

	// Parse the template files
	ts, err := template.ParseFiles(templateFiles...)
	if err != nil {
		logger.Logger.ErrorContext(
			ctx,
			"Template parsing failed",
			"error", err,
			"files", templateFiles,
		)
		http.Error(w, fmt.Sprintf(content.InternalErrorTemplate, "JVE-IHB-TP"), http.StatusInternalServerError)
		return
	}

	// Execute the template
	var buf bytes.Buffer
	if err := ts.ExecuteTemplate(&buf, "base", data); err != nil {
		logger.Logger.ErrorContext(
			ctx,
			"Template execution failed",
			"error", err,
			"template", "base",
		)
		http.Error(w, fmt.Sprintf(content.InternalErrorTemplate, "JVE-IHB-TX"), http.StatusInternalServerError)
		return
	}

	duration := time.Since(startTime)
	logger.Logger.InfoContext(ctx, "Rendered route in cold path", "route", route, "duration", duration.String())

	// Cache the rendered content
	content := buf.Bytes()
	cache := h.getCache()
	pageCache := cache.GetPageCache()
	pageCache.Set(route, content)

	// Check if context is cancelled before writing response
	select {
	case <-ctx.Done():
		logger.Logger.WarnContext(
			ctx,
			"Request cancelled before response",
			"error", ctx.Err(),
		)
		return
	default:
	}

	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	if _, err := w.Write(content); err != nil {
		logger.Logger.ErrorContext(
			ctx,
			"Failed to write cached response",
			"error", err,
		)
	}
}
