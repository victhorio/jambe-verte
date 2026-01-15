package content

import (
	"bytes"
	"context"
	"fmt"
	"html/template"
	"os"
	"path/filepath"
	"regexp"
	"slices"
	"strings"
	"time"

	"github.com/goccy/go-yaml"
	"github.com/victhorio/jambe-verte/internal/logger"
	"github.com/yuin/goldmark"
	highlighting "github.com/yuin/goldmark-highlighting/v2"
	meta "github.com/yuin/goldmark-meta"
	"github.com/yuin/goldmark/parser"
)

var (
	// postFilenameRegex validates post filenames follow the YYYY-MM-DD-slug.md pattern
	postFilenameRegex = regexp.MustCompile(`^\d{4}-\d{2}-\d{2}-[a-z0-9-]+\.md$`)

	// mdParser is the shared goldmark instance for converting markdown to HTML
	mdParser = goldmark.New(
		goldmark.WithExtensions(
			meta.Meta,
			highlighting.NewHighlighting(highlighting.WithStyle("github")),
		),
	)
)

// LoadPosts reads all the files ending in Markdown in a given `dir`, returning a list
// of Post structs. If the `isPost` parameter is true, the naming convention for posts
// will be checked against the YYYY-MM-DD-slug.md pattern and results will be returned
// sorted from newest to oldest.
func LoadPosts(dir string, isPost bool) ([]*Post, error) {
	ctx := context.Background()
	start := time.Now()

	paths, err := filepath.Glob(filepath.Join(dir, "*.md"))
	if err != nil {
		return nil, err
	}

	// Since I pretty much always expect to have content on hand when calling this function
	// I'm just going to error out if nothing is found to assert my expectations instead of
	// being more general/friendly.
	if len(paths) == 0 {
		logger.Logger.ErrorContext(ctx, "No content found", "directory", dir)
		return nil, fmt.Errorf("LoadPosts: no content found in %s", dir)
	}

	var contentList []*Post
	for _, path := range paths {
		post, err := loadPost(path, isPost)
		if err != nil {
			logger.Logger.WarnContext(
				ctx,
				"Failed to load post",
				"path", path,
				"error", err,
			)
			continue
		}

		// Some content will return as nil indicating that we should skip it even though
		// there weren't any errors.
		if post != nil {
			contentList = append(contentList, post)
		}
	}

	if isPost {
		slices.SortFunc(contentList, func(a, b *Post) int {
			// Note that we do b.Date.Compare instead of a.Date.Compare since SortFunc will
			// sort in ascending order, but we want the highest dates (most recent) first.
			return b.Date.Compare(a.Date)
		})
	}

	duration := time.Since(start)
	logger.Logger.InfoContext(ctx, "Loaded content", "is_post", isPost, "count", len(contentList), "directory", dir, "duration", duration.String())
	return contentList, nil
}

// loadPost is a helper function that loads a post from a given path and returns a Post struct.
// If it's reading an actual isPost, it will assert the naming convention of YYYY-MM-DD-slug.md
// as well as clean up the date prefix when creating returning the slug.
func loadPost(path string, isPost bool) (*Post, error) {
	// First, if it's a post, let's make sure that the file follows the correct naming convention of YYYY-MM-DD-slug.md
	base := filepath.Base(path)
	if isPost {
		if !postFilenameRegex.MatchString(base) {
			return nil, fmt.Errorf("invalid filename for post `%s`: expected: YYYY-MM-DD-slug.md", path)
		}
	}

	content, err := os.ReadFile(path)
	if err != nil {
		return nil, fmt.Errorf("failed to read post `%s`: %w", path, err)
	}

	var htmlBuf bytes.Buffer
	context := parser.NewContext()
	if err := mdParser.Convert(content, &htmlBuf, parser.WithContext(context)); err != nil {
		return nil, fmt.Errorf("failed to convert post `%s`: %w", path, err)
	}

	// Get metadata
	var postMeta PostFrontmatter
	metaData := meta.Get(context)
	yamlBytes, err := yaml.Marshal(metaData)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal post `%s` metadata: %w", path, err)
	}
	if err := yaml.Unmarshal(yamlBytes, &postMeta); err != nil {
		return nil, fmt.Errorf("failed to unmarshal post `%s` metadata: %w", path, err)
	}

	// Skip drafts
	// TODO: Add a flag to the command line to include drafts
	if postMeta.Draft {
		return nil, nil
	}

	// Parse date
	date, err := time.Parse("2006-01-02", postMeta.Date)
	if err != nil {
		return nil, fmt.Errorf("invalid date format for post `%s`: %s", path, postMeta.Date)
	}

	// Generate slug from filename
	slug := strings.TrimSuffix(base, filepath.Ext(base))
	if isPost {
		slug = slug[11:] // Remove the date prefix, which we already asserted is present
	}

	return &Post{
		Slug:        slug,
		Title:       postMeta.Title,
		Date:        date,
		Tags:        postMeta.Tags,
		Description: postMeta.Description,
		HTML:        template.HTML(htmlBuf.String()),
	}, nil
}
