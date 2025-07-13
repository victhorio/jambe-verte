package handlers

import (
	"bytes"
	"context"
	"html/template"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/victhorio/jambe-verte/internal/cache"
	"github.com/victhorio/jambe-verte/internal/content"
	"github.com/victhorio/jambe-verte/internal/logger"
)

type Handler struct {
	// TODO: Add a way to swap the cache during execution
	cache *cache.Cache
}

func New(cache *cache.Cache) *Handler {
	return &Handler{
		cache: cache,
	}
}

func (h *Handler) Home(w http.ResponseWriter, r *http.Request) {
	// Check page cache first
	pageCache := h.cache.GetPageCache()
	if cached, ok := pageCache.Get("/"); ok {
		w.Header().Set("Content-Type", "text/html; charset=utf-8")
		w.Write(cached)
		return
	}

	posts := h.cache.GetPosts()

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
	pageCache := h.cache.GetPageCache()
	if cached, ok := pageCache.Get("/blog"); ok {
		w.Header().Set("Content-Type", "text/html; charset=utf-8")
		w.Write(cached)
		return
	}

	posts := h.cache.GetPosts()

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
	pageCache := h.cache.GetPageCache()
	cachePath := "/blog/" + slug
	if cached, ok := pageCache.Get(cachePath); ok {
		w.Header().Set("Content-Type", "text/html; charset=utf-8")
		w.Write(cached)
		return
	}

	post, ok := h.cache.GetPost(slug)
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
	h.renderAndCache(r.Context(), w, cachePath, files, data)
}

func (h *Handler) ShowPage(w http.ResponseWriter, r *http.Request) {
	slug := chi.URLParam(r, "page")

	// Check page cache first
	pageCache := h.cache.GetPageCache()
	cachePath := "/" + slug
	if cached, ok := pageCache.Get(cachePath); ok {
		w.Header().Set("Content-Type", "text/html; charset=utf-8")
		w.Write(cached)
		return
	}

	page, ok := h.cache.GetPage(slug)
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
	h.renderAndCache(r.Context(), w, cachePath, files, data)
}

func (h *Handler) PostsByTag(w http.ResponseWriter, r *http.Request) {
	tag := chi.URLParam(r, "tag")

	// Check page cache first
	pageCache := h.cache.GetPageCache()
	cachePath := "/tag/" + tag
	if cached, ok := pageCache.Get(cachePath); ok {
		w.Header().Set("Content-Type", "text/html; charset=utf-8")
		w.Write(cached)
		return
	}

	posts := h.cache.GetPostsByTag(tag)

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
	h.renderAndCache(r.Context(), w, cachePath, files, data)
}

func (h *Handler) renderAndCache(ctx context.Context, w http.ResponseWriter, path string, templateFiles []string, data any) {
	// Check if context is cancelled
	select {
	case <-ctx.Done():
		logger.Logger.WarnContext(ctx, "Request cancelled", "error", ctx.Err())
		return
	default:
	}

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

	// Cache the rendered content
	content := buf.Bytes()
	pageCache := h.cache.GetPageCache()
	pageCache.Set(path, content)

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
