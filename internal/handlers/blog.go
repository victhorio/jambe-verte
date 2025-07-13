package handlers

import (
	"bytes"
	"context"
	"html/template"
	"net/http"
	"sync"
	"time"

	"github.com/go-chi/chi/v5"
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

	posts := cache.GetPosts()

	// Get recent posts (max 5)
	recentPosts := posts
	if len(posts) > 5 {
		recentPosts = posts[:5]
	}

	data := struct {
		Title string
		Posts []*content.Post
	}{
		Title: "Home",
		Posts: recentPosts,
	}

	files := []string{
		"templates/base.html",
		"templates/partials/nav.html",
		"templates/partials/footer.html",
		"templates/pages/home.html",
	}
	h.renderAndCache(r.Context(), w, "/", files, data)
}

func (h *Handler) ListPosts(w http.ResponseWriter, r *http.Request) {
	// Check page cache first
	cache := h.getCache()
	pageCache := cache.GetPageCache()
	if cached, ok := pageCache.Get("/blog"); ok {
		w.Header().Set("Content-Type", "text/html; charset=utf-8")
		w.Write(cached)
		return
	}

	posts := cache.GetPosts()

	data := struct {
		Title string
		Tag   string
		Posts []*content.Post
	}{
		Title: "Blog",
		Tag:   "", // Empty for general blog listing
		Posts: posts,
	}

	files := []string{
		"templates/base.html",
		"templates/partials/nav.html",
		"templates/partials/footer.html",
		"templates/pages/list.html",
	}
	h.renderAndCache(r.Context(), w, "/blog", files, data)
}

func (h *Handler) ShowPost(w http.ResponseWriter, r *http.Request) {
	slug := chi.URLParam(r, "slug")

	// Check page cache first
	cache := h.getCache()
	pageCache := cache.GetPageCache()
	cacheRoute := "/blog/" + slug
	if cached, ok := pageCache.Get(cacheRoute); ok {
		w.Header().Set("Content-Type", "text/html; charset=utf-8")
		w.Write(cached)
		return
	}

	post, ok := cache.GetPost(slug)
	if !ok {
		http.NotFound(w, r)
		return
	}

	data := struct {
		Title string
		Post  *content.Post
	}{
		Title: post.Title,
		Post:  post,
	}

	files := []string{
		"templates/base.html",
		"templates/partials/nav.html",
		"templates/partials/footer.html",
		"templates/pages/post.html",
	}
	h.renderAndCache(r.Context(), w, cacheRoute, files, data)
}

func (h *Handler) ShowPage(w http.ResponseWriter, r *http.Request) {
	slug := chi.URLParam(r, "page")

	// Check page cache first
	cache := h.getCache()
	pageCache := cache.GetPageCache()
	cacheRoute := "/" + slug
	if cached, ok := pageCache.Get(cacheRoute); ok {
		w.Header().Set("Content-Type", "text/html; charset=utf-8")
		w.Write(cached)
		return
	}

	page, ok := cache.GetPage(slug)
	if !ok {
		http.NotFound(w, r)
		return
	}

	data := struct {
		Title string
		Page  *content.Post
	}{
		Title: page.Title,
		Page:  page,
	}

	files := []string{
		"templates/base.html",
		"templates/partials/nav.html",
		"templates/partials/footer.html",
		"templates/pages/page.html",
	}
	h.renderAndCache(r.Context(), w, cacheRoute, files, data)
}

func (h *Handler) PostsByTag(w http.ResponseWriter, r *http.Request) {
	tag := chi.URLParam(r, "tag")

	// Check page cache first
	cache := h.getCache()
	pageCache := cache.GetPageCache()
	cacheRoute := "/tag/" + tag
	if cached, ok := pageCache.Get(cacheRoute); ok {
		w.Header().Set("Content-Type", "text/html; charset=utf-8")
		w.Write(cached)
		return
	}

	posts := cache.GetPostsByTag(tag)

	data := struct {
		Title string
		Tag   string
		Posts []*content.Post
	}{
		Title: "Posts tagged: " + tag,
		Tag:   tag,
		Posts: posts,
	}

	files := []string{
		"templates/base.html",
		"templates/partials/nav.html",
		"templates/partials/footer.html",
		"templates/pages/list.html",
	}
	h.renderAndCache(r.Context(), w, cacheRoute, files, data)
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
		http.Error(w, "Internal Server Error: ihr-tp", http.StatusInternalServerError)
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
		http.Error(w, "Internal Server Error: ihr-tx", http.StatusInternalServerError)
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
