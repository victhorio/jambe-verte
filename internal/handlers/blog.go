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
	h.render(r.Context(), w, files, data)
}

func (h *Handler) ListPosts(w http.ResponseWriter, r *http.Request) {
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
	h.render(r.Context(), w, files, data)
}

func (h *Handler) ShowPost(w http.ResponseWriter, r *http.Request) {
	slug := chi.URLParam(r, "slug")
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
	h.render(r.Context(), w, files, data)
}

func (h *Handler) ShowPage(w http.ResponseWriter, r *http.Request) {
	slug := chi.URLParam(r, "page")
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
	h.render(r.Context(), w, files, data)
}

func (h *Handler) PostsByTag(w http.ResponseWriter, r *http.Request) {
	tag := chi.URLParam(r, "tag")
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
	h.render(r.Context(), w, files, data)
}

// Test endpoint to trigger template error
func (h *Handler) TestError(w http.ResponseWriter, r *http.Request) {
	// This will cause a template error because we're passing wrong data structure
	data := struct {
		WrongField string
	}{
		WrongField: "This will break the template",
	}

	// Try to execute template with wrong data - should trigger error
	files := []string{
		"templates/base.html",
		"templates/partials/nav.html",
		"templates/partials/footer.html",
		"templates/pages/post.html",
	}
	h.render(r.Context(), w, files, data)
}

func (h *Handler) render(ctx context.Context, w http.ResponseWriter, templateFiles []string, data interface{}) {
	// Check if context is cancelled
	select {
	case <-ctx.Done():
		logger.Logger.WarnContext(ctx, "Request cancelled", "error", ctx.Err())
		return
	default:
	}

	// Parse the template files for this specific request
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
	if _, err := buf.WriteTo(w); err != nil {
		logger.Logger.ErrorContext(
			ctx,
			"Failed to write response",
			"error", err,
		)
	}
}
