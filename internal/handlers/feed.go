package handlers

import (
	"encoding/xml"
	"net/http"
	"time"

	"github.com/victhorio/jambe-verte/internal/logger"
)

type RSS struct {
	XMLName xml.Name `xml:"rss"`
	Version string   `xml:"version,attr"`
	Channel Channel  `xml:"channel"`
}

type Channel struct {
	Title       string `xml:"title"`
	Link        string `xml:"link"`
	Description string `xml:"description"`
	Items       []Item `xml:"item"`
}

type Item struct {
	Title       string `xml:"title"`
	Link        string `xml:"link"`
	Description string `xml:"description"`
	PubDate     string `xml:"pubDate"`
}

func (h *Handler) RSSFeed(w http.ResponseWriter, r *http.Request) {
	posts := h.cache.GetPosts()

	// Get recent posts (max 20)
	feedPosts := posts
	if len(posts) > 20 {
		feedPosts = posts[:20]
	}

	// Determine the scheme from the request
	scheme := "http"
	if r.TLS != nil || r.Header.Get("X-Forwarded-Proto") == "https" {
		scheme = "https"
	}
	baseURL := scheme + "://" + r.Host

	// Build RSS items
	items := make([]Item, len(feedPosts))
	for i, post := range feedPosts {
		items[i] = Item{
			Title:       post.Title,
			Link:        baseURL + "/blog/" + post.Slug,
			Description: post.Description,
			PubDate:     post.Date.Format(time.RFC1123Z),
		}
	}

	rss := RSS{
		Version: "2.0",
		Channel: Channel{
			Title:       "My Blog",
			Link:        baseURL,
			Description: "A minimal Go blog",
			Items:       items,
		},
	}

	w.Header().Set("Content-Type", "application/rss+xml")
	if err := xml.NewEncoder(w).Encode(rss); err != nil {
		logger.Logger.ErrorContext(r.Context(), "Failed to encode RSS feed",
			"error", err,
			"posts_count", len(feedPosts))
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
	}
}
