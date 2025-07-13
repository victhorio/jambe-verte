package content

import (
	"bytes"
	"context"
	"fmt"
	"html/template"
	"os"
	"path/filepath"
	"sort"
	"strings"
	"time"

	"github.com/victhor/jv-blog/internal/logger"
	"github.com/yuin/goldmark"
	meta "github.com/yuin/goldmark-meta"
	"github.com/yuin/goldmark/parser"
	"gopkg.in/yaml.v3"
)

func LoadPosts(dir string) ([]*Post, error) {
	ctx := context.Background()

	files, err := filepath.Glob(filepath.Join(dir, "*.md"))
	if err != nil {
		return nil, err
	}

	var posts []*Post
	for _, file := range files {
		post, err := loadPost(file)
		if err != nil {
			logger.Logger.WarnContext(ctx, "Failed to load post",
				"file", file,
				"error", err)
			continue
		}
		if post != nil {
			posts = append(posts, post)
		}
	}

	// Sort posts by date, newest first
	sort.Slice(posts, func(i, j int) bool {
		return posts[i].Date.After(posts[j].Date)
	})

	logger.Logger.InfoContext(ctx, "Loaded posts", "count", len(posts), "directory", dir)
	return posts, nil
}

func loadPost(path string) (*Post, error) {
	content, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}

	// Parse markdown with frontmatter
	md := goldmark.New(
		goldmark.WithExtensions(meta.Meta),
	)

	var buf bytes.Buffer
	context := parser.NewContext()
	if err := md.Convert(content, &buf, parser.WithContext(context)); err != nil {
		return nil, err
	}

	// Get metadata
	metaData := meta.Get(context)

	// Parse metadata into PostMeta struct
	var postMeta PostMeta
	yamlBytes, err := yaml.Marshal(metaData)
	if err != nil {
		return nil, err
	}
	if err := yaml.Unmarshal(yamlBytes, &postMeta); err != nil {
		return nil, err
	}

	// Skip drafts
	if postMeta.Draft {
		return nil, nil
	}

	// Parse date
	date, err := time.Parse("2006-01-02", postMeta.Date)
	if err != nil {
		// Try parsing with time
		date, err = time.Parse("2006-01-02 15:04:05", postMeta.Date)
		if err != nil {
			return nil, fmt.Errorf("invalid date format: %s", postMeta.Date)
		}
	}

	// Generate slug from filename
	base := filepath.Base(path)
	slug := strings.TrimSuffix(base, filepath.Ext(base))
	// Remove date prefix if present (YYYY-MM-DD-slug.md format)
	if len(slug) > 11 && slug[4] == '-' && slug[7] == '-' && slug[10] == '-' {
		slug = slug[11:]
	}

	return &Post{
		Slug:        slug,
		Title:       postMeta.Title,
		Date:        date,
		Tags:        postMeta.Tags,
		Description: postMeta.Description,
		Content:     template.HTML(buf.String()),
		Raw:         string(content),
	}, nil
}

func LoadPage(path string) (*Post, error) {
	return loadPost(path)
}

func LoadPages(dir string) ([]*Post, error) {
	ctx := context.Background()

	files, err := filepath.Glob(filepath.Join(dir, "*.md"))
	if err != nil {
		return nil, err
	}

	var pages []*Post
	for _, file := range files {
		page, err := loadPost(file)
		if err != nil {
			logger.Logger.WarnContext(ctx, "Failed to load page",
				"file", file,
				"error", err)
			continue
		}
		if page != nil {
			// For pages, use filename as slug
			base := filepath.Base(file)
			page.Slug = strings.TrimSuffix(base, filepath.Ext(base))
			pages = append(pages, page)
		}
	}

	logger.Logger.InfoContext(ctx, "Loaded pages", "count", len(pages), "directory", dir)
	return pages, nil
}
