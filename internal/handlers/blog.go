package handlers

import (
	"bytes"
	"context"
	"html/template"
	"net/http"
	"sync"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/victhorio/jambe-verte/internal"
	"github.com/victhorio/jambe-verte/internal/cache"
	"github.com/victhorio/jambe-verte/internal/content"
	"github.com/victhorio/jambe-verte/internal/logger"
)

type Handler struct {
	mu        sync.RWMutex
	cache     *cache.Cache
	debugMode bool
}

type PostsPageData struct {
	Posts []*content.Post
	Tag   string
}

func New(cache *cache.Cache, debugMode bool) *Handler {
	return &Handler{
		cache:     cache,
		debugMode: debugMode,
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
	// Check page cache first, unless we're in debug mode
	if !h.debugMode {
		cache := h.getCache()
		pageCache := cache.GetPageCache()
		if cached, ok := pageCache.Get("/"); ok {
			w.Header().Set("Content-Type", "text/html; charset=utf-8")
			w.Write(cached)
			return
		}
	}

	files := []string{
		"templates/base.html",
		"templates/home.html",
	}
	h.renderAndCache(r.Context(), w, "/", files, nil)
}

func (h *Handler) ListPosts(w http.ResponseWriter, r *http.Request) {
	// Check page cache first, unless we're in debug mode
	if !h.debugMode {
		cache := h.getCache()
		pageCache := cache.GetPageCache()
		if cached, ok := pageCache.Get("/posts"); ok {
			w.Header().Set("Content-Type", "text/html; charset=utf-8")
			w.Write(cached)
			return
		}
	}

	data := PostsPageData{
		Posts: h.getCache().GetPosts(),
	}

	files := []string{
		"templates/base.html",
		"templates/posts.html",
	}
	h.renderAndCache(r.Context(), w, "/posts", files, data)
}

func (h *Handler) ShowPost(w http.ResponseWriter, r *http.Request) {
	slug := chi.URLParam(r, "slug")
	cache := h.getCache()

	post, ok := cache.GetPost(slug)
	if !ok {
		http.NotFound(w, r)
		return
	}

	// Check page cache first, unless we're in debug mode
	route := "/blog/" + slug
	if !h.debugMode {
		pageCache := cache.GetPageCache()
		if cached, ok := pageCache.Get(route); ok {
			w.Header().Set("Content-Type", "text/html; charset=utf-8")
			w.Write(cached)
			return
		}
	}

	files := []string{
		"templates/base.html",
		"templates/post.html",
	}
	h.renderAndCache(r.Context(), w, route, files, post)
}

func (h *Handler) ShowPage(w http.ResponseWriter, r *http.Request) {
	slug := chi.URLParam(r, "page")
	cache := h.getCache()

	page, ok := cache.GetPage(slug)
	if !ok {
		http.NotFound(w, r)
		return
	}

	// Check page cache first, unless we're in debug mode
	route := "/" + slug
	if !h.debugMode {
		pageCache := cache.GetPageCache()
		if cached, ok := pageCache.Get(route); ok {
			w.Header().Set("Content-Type", "text/html; charset=utf-8")
			w.Write(cached)
			return
		}
	}

	files := []string{
		"templates/base.html",
		"templates/page.html",
	}
	h.renderAndCache(r.Context(), w, route, files, page)
}

func (h *Handler) PostsByTag(w http.ResponseWriter, r *http.Request) {
	tag := chi.URLParam(r, "tag")
	cache := h.getCache()

	posts := cache.GetPostsByTag(tag)
	if len(posts) == 0 {
		http.NotFound(w, r)
		return
	}

	// Check page cache first, unless we're in debug mode
	route := "/tag/" + tag
	if !h.debugMode {
		pageCache := cache.GetPageCache()
		if cached, ok := pageCache.Get(route); ok {
			w.Header().Set("Content-Type", "text/html; charset=utf-8")
			w.Write(cached)
			return
		}
	}

	data := PostsPageData{
		Posts: posts,
		Tag:   tag,
	}

	files := []string{
		"templates/base.html",
		"templates/posts.html",
	}
	h.renderAndCache(r.Context(), w, route, files, data)
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
		internal.WriteInternalError(w, "JVE-IHB-PO")
		return
	}

	// Load pages
	pages, err := content.LoadPosts("content/pages", false)
	if err != nil {
		logger.Logger.ErrorContext(ctx, "Error loading pages during refresh", "error", err)
		w.Header().Set("Content-Type", "text/plain; charset=utf-8")
		internal.WriteInternalError(w, "JVE-IHB-PA")
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
		internal.WriteInternalError(w, "JVE-IHB-TP")
		return
	}

	// Execute the template
	var buf bytes.Buffer
	if err := ts.ExecuteTemplate(&buf, "base", map[string]any{
		"DebugMode": h.debugMode,
		"Version":   internal.Version,
		"Data":      data,
	}); err != nil {
		logger.Logger.ErrorContext(
			ctx,
			"Template execution failed",
			"error", err,
			"template", "base",
		)
		internal.WriteInternalError(w, "JVE-IHB-TX")
		return
	}

	duration := time.Since(startTime)
	logger.Logger.InfoContext(ctx, "Rendered route in cold path", "route", route, "duration", duration.String())

	// Cache the rendered content (skip caching in debug mode)
	content := buf.Bytes()
	if !h.debugMode {
		cache := h.getCache()
		pageCache := cache.GetPageCache()
		pageCache.Set(route, content)
	}

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
