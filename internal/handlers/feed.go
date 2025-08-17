package handlers

import (
	"encoding/xml"
	"net/http"
	"time"

	"github.com/victhorio/jambe-verte/internal"
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
	GUID        string `xml:"guid"`
}

func (h *Handler) RSSFeed(w http.ResponseWriter, r *http.Request) {
	// Get recent posts (max 20)
	posts := h.getCache().GetPosts()
	if len(posts) > 20 {
		posts = posts[:20]
	}

	// Determine the scheme from the request
	scheme := "http"
	if r.TLS != nil || r.Header.Get("X-Forwarded-Proto") == "https" {
		scheme = "https"
	}
	baseURL := scheme + "://" + r.Host

	// Build RSS items
	items := make([]Item, len(posts))
	for i, post := range posts {
		postURL := baseURL + "/blog/" + post.Slug
		items[i] = Item{
			Title:       post.Title,
			Link:        postURL,
			Description: post.Description,
			PubDate:     post.Date.Format(time.RFC1123Z),
			GUID:        postURL, // Use URL as unique identifier
		}
	}

	rss := RSS{
		Version: "2.0",
		Channel: Channel{
			Title:       "Jambe Verte",
			Link:        baseURL,
			Description: "A blog by Victhor Sartório",
			Items:       items,
		},
	}

	w.Header().Set("Content-Type", "application/rss+xml; charset=utf-8")
	if err := xml.NewEncoder(w).Encode(rss); err != nil {
		logger.Logger.ErrorContext(
			r.Context(),
			"Failed to encode RSS feed",
			"error", err,
			"posts_count", len(posts),
		)
		internal.WriteInternalError(w, "JVE-IHF-XE")
	}
}
