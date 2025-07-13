## Architecture Overview

This is a Go-based blog application built with the chi router. The architecture follows clean separation of concerns:

### Key Components

1. **Entry Point** (`cmd/jv-server/main.go`): Initializes the server, sets up middleware, and defines routes.

   - There's also a cmd/jv-helper/main.go for better DX when creating new posts or pages.

2. **Content System** (`internal/content/`):

   - Loads markdown files from `content/posts/` and `content/pages/`
   - Parses YAML frontmatter for metadata (title, date, tags, description)
   - Posts follow naming: `YYYY-MM-DD-slug.md`
   - Content is cached in memory at startup via `internal/cache/`

3. **HTTP Handlers** (`internal/handlers/`):

   - `HomeHandler`: Shows 5 most recent posts
   - `BlogHandler`: Lists all blog posts
   - `PostHandler`: Renders individual posts
   - `TagHandler`: Filters posts by tag
   - `FeedHandler`: Generates RSS feed
   - `PageHandler`: Serves static pages
   - etc

4. **Templates** (`templates/`): Go HTML templates with partials for reusable components

5. **Middleware Stack** (order matters):
   - Request ID generation
   - Structured JSON logging to stderr
   - Panic recovery
   - 60-second timeout
   - gzip compression (level 5)

### Routing Structure

- `/` - Home page
- `/blog` - Blog listing
- `/blog/{slug}` - Individual post
- `/tag/{tag}` - Posts by tag
- `/feed.xml` - RSS feed
- `/{page}` - Static pages (e.g., /about)
- `/static/*` - Static assets
- etc

## Development Guidelines

1. **Adding New Posts**: Use `jv-helper post <slug>` or manually create markdown files in `content/posts/` with format `YYYY-MM-DD-title.md` and include required frontmatter.

2. **Logging**: Use the logger from `internal/logger/` which outputs structured JSON to stderr.

3. **Error Handling**: Return appropriate HTTP status codes. The middleware handles panic recovery.

4. **Templates**: When modifying templates, ensure partials are properly included using `{{template "partial-name" .}}`

5. **Static Assets**: Place CSS in `static/css/` and JavaScript in `static/js/`. They're served from `/static/`.

## Testing

Currently no test files exist. When adding tests:

```bash
# Run all tests
go test ./...

# Run tests with coverage
go test -cover ./...

# Run specific package tests
go test ./internal/handlers
```

## Code Quality Checks

Always run these after making changes:

```bash
# Ensure code compiles
go build ./...

# Format code
go fmt ./...

# Run static analysis
go vet ./...
```

## Dependencies

The project uses minimal dependencies:

- `chi/v5` for HTTP routing
- `goldmark` for markdown parsing
- `goldmark-meta` for YAML frontmatter
- `goccy/go-yaml` for YAML parsing

Frontend uses Alpine.js loaded from CDN.
