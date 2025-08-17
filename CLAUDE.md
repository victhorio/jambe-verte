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
   - Rendered pages are cached for performance

3. **HTTP Handlers** (`internal/handlers/`):

   - `Home`: Serves the About page content on the home route
   - `ListPosts`: Lists all blog posts
   - `ShowPost`: Renders individual posts
   - `PostsByTag`: Filters posts by tag
   - `RSSFeed`: Generates RSS feed (max 20 recent posts)
   - `ShowPage`: Serves static pages
   - `AdminRefresh`: Hot-reloads content and rebuilds CSS

4. **Templates** (`templates/`): Go HTML templates with a consolidated base layout

5. **Middleware Stack** (order matters):
   - Request ID generation
   - Structured JSON logging to stderr
   - Panic recovery
   - 60-second timeout
   - gzip compression (level 5)
   - Admin auth (for protected routes only)

### Routing Structure

- `/` - Home page (displays About page content)
- `/posts` - Blog listing
- `/blog/{slug}` - Individual post
- `/tag/{tag}` - Posts by tag
- `/feed.xml` - RSS feed
- `/{page}` - Static pages (e.g., /about)
- `/static/*` - Static assets
- `/admin/refresh` - Hot-reload content (protected)

## Development Guidelines

1. **Adding New Posts**: Use `jv-helper post <slug>` or manually create markdown files in `content/posts/` with format `YYYY-MM-DD-title.md` and include required frontmatter.

2. **Logging**: Use the logger from `internal/logger/` which outputs structured JSON to stderr.

3. **Hot Reloading**: Use `/admin/refresh` endpoint to reload content and rebuild CSS without restarting the server.

4. **Error Handling**: Return appropriate HTTP status codes. The middleware handles panic recovery.

5. **Templates**: The base template (`templates/base.html`) contains the complete page structure including nav and footer. Individual page templates define `title` and `main` blocks.

6. **Static Assets**: Place CSS in `static/css/` and JavaScript in `static/js/`. They're served from `/static/`.

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

## CSS Build Process

The application uses Tailwind CSS for styling. After making changes to templates or CSS:

```bash
# Install dependencies (first time only)
bun install

# Build CSS for production
bun run build-css

# Watch for changes during development
bun run watch-css
```

**CSS File Structure:**

- `static/css/input.css` - Main Tailwind input file with custom component styles
- `static/css/output.css` - Compiled output file (this is what the website loads)

**Important:** Always run `bun run build-css` after making changes to `input.css` to update the compiled `output.css` file.

## Dependencies

The project uses minimal dependencies:

- `chi/v5` for HTTP routing
- `goldmark` for markdown parsing
- `goldmark-meta` for YAML frontmatter
- `goldmark-highlighting` for syntax highlighting
- `goccy/go-yaml` for YAML parsing

Frontend uses Alpine.js loaded from CDN.
