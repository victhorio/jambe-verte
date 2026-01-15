package handlers

import (
	"bytes"
	"context"
	"fmt"
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

// Handler manages HTTP request handling with hot-reloadable content caching.
//
// # Cache Consistency Model
//
// The cache can be swapped atomically at any time via AdminRefresh (e.g., when
// content is updated). This creates a subtle concurrency requirement: each request
// must use a single, consistent cache snapshot throughout its lifetime.
//
// The problem arises because each Cache instance contains its own PageCache for
// storing rendered HTML. If a handler fetches data from Cache A, then AdminRefresh
// swaps in Cache B, and then the handler writes rendered HTML to "the current cache",
// it would write stale content (rendered from A's data) into Cache B's PageCache.
// Future requests would then serve this stale cached HTML until the next refresh.
//
// To prevent this, handlers must:
//  1. Call getCache() exactly once at the start of the request
//  2. Extract both content data AND the PageCache from that same snapshot
//  3. Pass the PageCache explicitly to renderAndCache()
//
// This ensures that if a cache swap occurs mid-request, the rendered content is
// written to the old cache's PageCache (which will be garbage collected) rather
// than polluting the new cache with stale data.
type Handler struct {
	mu        sync.RWMutex
	cache     *cache.Cache
	debugMode bool

	// Pre-parsed templates (parsed once at startup)
	homeTmpl  *template.Template
	postsTmpl *template.Template
	postTmpl  *template.Template
	pageTmpl  *template.Template
}

type PostsPageData struct {
	Posts []*content.Post
	Tag   string
}

func New(cache *cache.Cache, debugMode bool) (*Handler, error) {
	homeTmpl, err := template.ParseFiles("templates/base.html", "templates/home.html")
	if err != nil {
		return nil, fmt.Errorf("parsing home template: %w", err)
	}

	postsTmpl, err := template.ParseFiles("templates/base.html", "templates/posts.html")
	if err != nil {
		return nil, fmt.Errorf("parsing posts template: %w", err)
	}

	postTmpl, err := template.ParseFiles("templates/base.html", "templates/post.html")
	if err != nil {
		return nil, fmt.Errorf("parsing post template: %w", err)
	}

	pageTmpl, err := template.ParseFiles("templates/base.html", "templates/page.html")
	if err != nil {
		return nil, fmt.Errorf("parsing page template: %w", err)
	}

	return &Handler{
		cache:     cache,
		debugMode: debugMode,
		homeTmpl:  homeTmpl,
		postsTmpl: postsTmpl,
		postTmpl:  postTmpl,
		pageTmpl:  pageTmpl,
	}, nil
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
	c := h.getCache()
	pageCache := c.GetPageCache()

	// Check page cache first, unless we're in debug mode
	if !h.debugMode {
		if cached, ok := pageCache.Get("/"); ok {
			w.Header().Set("Content-Type", "text/html; charset=utf-8")
			w.Write(cached)
			return
		}
	}

	h.renderAndCache(r.Context(), w, pageCache, "/", h.homeTmpl, nil)
}

func (h *Handler) ListPosts(w http.ResponseWriter, r *http.Request) {
	c := h.getCache()
	pageCache := c.GetPageCache()

	// Check page cache first, unless we're in debug mode
	if !h.debugMode {
		if cached, ok := pageCache.Get("/posts"); ok {
			w.Header().Set("Content-Type", "text/html; charset=utf-8")
			w.Write(cached)
			return
		}
	}

	data := PostsPageData{
		Posts: c.GetPosts(),
	}

	h.renderAndCache(r.Context(), w, pageCache, "/posts", h.postsTmpl, data)
}

func (h *Handler) ShowPost(w http.ResponseWriter, r *http.Request) {
	slug := chi.URLParam(r, "slug")
	c := h.getCache()
	pageCache := c.GetPageCache()

	post, ok := c.GetPost(slug)
	if !ok {
		http.NotFound(w, r)
		return
	}

	// Check page cache first, unless we're in debug mode
	route := "/blog/" + slug
	if !h.debugMode {
		if cached, ok := pageCache.Get(route); ok {
			w.Header().Set("Content-Type", "text/html; charset=utf-8")
			w.Write(cached)
			return
		}
	}

	h.renderAndCache(r.Context(), w, pageCache, route, h.postTmpl, post)
}

func (h *Handler) ShowPage(w http.ResponseWriter, r *http.Request) {
	slug := chi.URLParam(r, "page")
	c := h.getCache()
	pageCache := c.GetPageCache()

	page, ok := c.GetPage(slug)
	if !ok {
		http.NotFound(w, r)
		return
	}

	// Check page cache first, unless we're in debug mode
	route := "/" + slug
	if !h.debugMode {
		if cached, ok := pageCache.Get(route); ok {
			w.Header().Set("Content-Type", "text/html; charset=utf-8")
			w.Write(cached)
			return
		}
	}

	h.renderAndCache(r.Context(), w, pageCache, route, h.pageTmpl, page)
}

func (h *Handler) PostsByTag(w http.ResponseWriter, r *http.Request) {
	tag := chi.URLParam(r, "tag")
	c := h.getCache()
	pageCache := c.GetPageCache()

	posts := c.GetPostsByTag(tag)
	if len(posts) == 0 {
		http.NotFound(w, r)
		return
	}

	// Check page cache first, unless we're in debug mode
	route := "/tag/" + tag
	if !h.debugMode {
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

	h.renderAndCache(r.Context(), w, pageCache, route, h.postsTmpl, data)
}

// AdminRefresh is responsible for hot-reloading content by creating an entirely new cache
// and replacing it on the handler.
func (h *Handler) AdminRefresh(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	// Load posts
	posts, err := content.LoadContent("content/posts", true)
	if err != nil {
		logger.Logger.ErrorContext(ctx, "Error loading posts during refresh", "error", err)
		internal.WriteInternalError(w, "JVE-IHB-PO")
		return
	}

	// Load pages
	pages, err := content.LoadContent("content/pages", false)
	if err != nil {
		logger.Logger.ErrorContext(ctx, "Error loading pages during refresh", "error", err)
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

func (h *Handler) renderAndCache(ctx context.Context, w http.ResponseWriter, pageCache *cache.PageCache, route string, tmpl *template.Template, data any) {
	// Check if context is cancelled
	select {
	case <-ctx.Done():
		logger.Logger.WarnContext(ctx, "Request cancelled", "error", ctx.Err())
		return
	default:
	}

	startTime := time.Now()

	// Execute the template
	var buf bytes.Buffer
	if err := tmpl.ExecuteTemplate(&buf, "base", map[string]any{
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
