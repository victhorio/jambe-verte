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

// Template file paths for each template name
var templateFiles = map[string][]string{
	"home":  {"templates/base.html", "templates/home.html"},
	"posts": {"templates/base.html", "templates/posts.html"},
	"post":  {"templates/base.html", "templates/post.html"},
	"page":  {"templates/base.html", "templates/page.html"},
}

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

	// Pre-parsed templates (parsed once at startup, used in production)
	templates map[string]*template.Template
}

type PostsPageData struct {
	Posts []*content.Post
	Tag   string
}

func New(cache *cache.Cache, debugMode bool) (*Handler, error) {
	templates := make(map[string]*template.Template)

	for name, files := range templateFiles {
		tmpl, err := template.ParseFiles(files...)
		if err != nil {
			return nil, fmt.Errorf("parsing %s template: %w", name, err)
		}
		templates[name] = tmpl
	}

	return &Handler{
		cache:     cache,
		debugMode: debugMode,
		templates: templates,
	}, nil
}

// getTemplate returns a template by name. In debug mode, it re-parses from disk
// to pick up any changes. In production, it returns the cached template.
func (h *Handler) getTemplate(name string) (*template.Template, error) {
	files, ok := templateFiles[name]
	if !ok {
		return nil, fmt.Errorf("unknown template: %s", name)
	}
	if !h.debugMode {
		return h.templates[name], nil
	}
	return template.ParseFiles(files...)
}

func (h *Handler) getCache() (*cache.Cache, error) {
	// In debug mode, reload content from disk for hot-reload
	if h.debugMode {
		posts, err := content.LoadContent("content/posts", true)
		if err != nil {
			return nil, fmt.Errorf("loading posts: %w", err)
		}
		pages, err := content.LoadContent("content/pages", false)
		if err != nil {
			return nil, fmt.Errorf("loading pages: %w", err)
		}
		return cache.New(posts, pages), nil
	}
	h.mu.RLock()
	defer h.mu.RUnlock()
	return h.cache, nil
}

func (h *Handler) setCache(cache *cache.Cache) {
	h.mu.Lock()
	defer h.mu.Unlock()
	h.cache = cache
}

func (h *Handler) Home(w http.ResponseWriter, r *http.Request) {
	c, err := h.getCache()
	if err != nil {
		logger.WithRequest(r.Context()).Error("Failed to load content", "error", err)
		internal.WriteInternalError(w, "JVE-IHB-LC")
		return
	}
	pageCache := c.GetPageCache()

	// Check page cache first, unless we're in debug mode
	if !h.debugMode {
		if cached, ok := pageCache.Get("/"); ok {
			w.Header().Set("Content-Type", "text/html; charset=utf-8")
			w.Write(cached)
			return
		}
	}

	h.renderAndCache(r.Context(), w, pageCache, "/", "home", nil)
}

func (h *Handler) ListPosts(w http.ResponseWriter, r *http.Request) {
	c, err := h.getCache()
	if err != nil {
		logger.WithRequest(r.Context()).Error("Failed to load content", "error", err)
		internal.WriteInternalError(w, "JVE-IHB-LC")
		return
	}
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

	h.renderAndCache(r.Context(), w, pageCache, "/posts", "posts", data)
}

func (h *Handler) ShowPost(w http.ResponseWriter, r *http.Request) {
	slug := chi.URLParam(r, "slug")
	c, err := h.getCache()
	if err != nil {
		logger.WithRequest(r.Context()).Error("Failed to load content", "error", err)
		internal.WriteInternalError(w, "JVE-IHB-LC")
		return
	}
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

	h.renderAndCache(r.Context(), w, pageCache, route, "post", post)
}

func (h *Handler) ShowPage(w http.ResponseWriter, r *http.Request) {
	slug := chi.URLParam(r, "page")
	c, err := h.getCache()
	if err != nil {
		logger.WithRequest(r.Context()).Error("Failed to load content", "error", err)
		internal.WriteInternalError(w, "JVE-IHB-LC")
		return
	}
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

	h.renderAndCache(r.Context(), w, pageCache, route, "page", page)
}

func (h *Handler) PostsByTag(w http.ResponseWriter, r *http.Request) {
	tag := chi.URLParam(r, "tag")
	c, err := h.getCache()
	if err != nil {
		logger.WithRequest(r.Context()).Error("Failed to load content", "error", err)
		internal.WriteInternalError(w, "JVE-IHB-LC")
		return
	}
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

	h.renderAndCache(r.Context(), w, pageCache, route, "posts", data)
}

// AdminRefresh is responsible for hot-reloading content by creating an entirely new cache
// and replacing it on the handler.
func (h *Handler) AdminRefresh(w http.ResponseWriter, r *http.Request) {
	log := logger.WithRequest(r.Context())

	// Load posts
	posts, err := content.LoadContent("content/posts", true)
	if err != nil {
		log.Error("Error loading posts during refresh", "error", err)
		internal.WriteInternalError(w, "JVE-IHB-PO")
		return
	}

	// Load pages
	pages, err := content.LoadContent("content/pages", false)
	if err != nil {
		log.Error("Error loading pages during refresh", "error", err)
		internal.WriteInternalError(w, "JVE-IHB-PA")
		return
	}

	// Create new cache and replace the old one
	newCache := cache.New(posts, pages)
	h.setCache(newCache)

	log.Info("Cache refreshed successfully", "posts", len(posts), "pages", len(pages))

	// Also attempt to rebuild CSS
	content.RebuildCSS(r.Context())

	w.Header().Set("Content-Type", "text/plain; charset=utf-8")
	w.Write([]byte("OK"))
}

func (h *Handler) renderAndCache(ctx context.Context, w http.ResponseWriter, pageCache *cache.PageCache, route string, templateName string, data any) {
	log := logger.WithRequest(ctx)

	// Check if context is cancelled
	select {
	case <-ctx.Done():
		log.Warn("Request cancelled", "error", ctx.Err())
		return
	default:
	}

	// In debug mode, rebuild CSS to pick up any new Tailwind classes
	if h.debugMode {
		content.RebuildCSS(ctx)
	}

	// Get template (fresh parse in debug mode, cached in production)
	tmpl, err := h.getTemplate(templateName)
	if err != nil {
		log.Error("Template parsing failed", "error", err, "template", templateName)
		internal.WriteInternalError(w, "JVE-IHB-TP")
		return
	}

	startTime := time.Now()

	// Execute the template
	var buf bytes.Buffer
	if err := tmpl.ExecuteTemplate(&buf, "base", map[string]any{
		"DebugMode": h.debugMode,
		"Version":   internal.Version,
		"Data":      data,
	}); err != nil {
		log.Error("Template execution failed", "error", err, "template", templateName)
		internal.WriteInternalError(w, "JVE-IHB-TX")
		return
	}

	duration := time.Since(startTime)
	log.Info("Rendered route in cold path", "route", route, "duration", duration.String())

	// Cache the rendered content (skip caching in debug mode)
	rendered := buf.Bytes()
	if !h.debugMode {
		pageCache.Set(route, rendered)
	}

	// Check if context is cancelled before writing response
	select {
	case <-ctx.Done():
		log.Warn("Request cancelled before response", "error", ctx.Err())
		return
	default:
	}

	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	if _, err := w.Write(rendered); err != nil {
		log.Error("Failed to write cached response", "error", err)
	}
}
