This file was last updated 2026-01-14.

## Overview

Personal blog/website for Victhor Sartório. Go server with chi router, markdown content with YAML frontmatter, Tailwind CSS and htmx.

## Architecture

- `cmd/jv-server/main.go` - Entry point. Loads content, sets up middleware stack (request ID, logging, recovery, timeout, gzip, auth), defines routes.
- `cmd/jv-helper/main.go` - CLI helper for creating new posts/pages with proper frontmatter.

- `internal/cache/cache.go` - In-memory content cache. Indexes posts by slug and tags.
- `internal/cache/page_cache.go` - Thread-safe rendered HTML page cache.
- `internal/content/loader.go` - Loads markdown files from disk, parses frontmatter. Skips drafts.
- `internal/content/post.go` - Post/Page structs and frontmatter definitions.
- `internal/content/utils.go` - CSS build utilities (runs `bun run build-css`).
- `internal/handlers/blog.go` - HTTP handlers: Home, ListPosts, ShowPost, ShowPage, PostsByTag, AdminRefresh.
- `internal/handlers/feed.go` - RSS feed generation.
- `internal/middleware/auth.go` - Bearer token authentication for admin routes.
- `internal/logger/logger.go` - Structured JSON logging to stderr.
- `internal/error.go` - Error response utilities with custom error codes.
- `internal/version.go` - Application version constant.

- `templates/base.html` - Master layout with nav and footer. Child templates define `title` and `main` blocks.
- `templates/home.html` - Home page template.

- `content/posts/` - Blog posts as `YYYY-MM-DD-slug.md` with frontmatter (title, date, tags, description, draft).
- `content/pages/` - Static pages as `slug.md` with frontmatter (title, description).

- `static/css/input.css` - Tailwind input with custom font declarations.
- `static/css/output.css` - Compiled CSS (run `bun run build-css` after changes).
- `static/fonts/` - Berkeley Mono font files (woff2).
- `static/favicons/` - Favicon variants.

- `deploy.sh` - Builds ARM64 Linux binary, creates tarball, deploys via rsync.

## Routes

- `/` - Home (displays About page content)
- `/posts` - Blog listing
- `/blog/{slug}` - Individual post
- `/tag/{tag}` - Posts filtered by tag
- `/feed.xml` - RSS feed
- `/{page}` - Static pages
- `/static/*` - Static assets
- `POST /admin/refresh` - Hot-reload content and rebuild CSS (requires `JV_ADMIN_TOKEN`)

## Environment Variables

- `JV_DEBUG` - Debug mode: skips auth, disables page caching, uses CDN CSS.
- `JV_ADMIN_TOKEN` - Bearer token for `/admin/refresh`.
- `JV_SERVER_IP` - Target server for `deploy.sh`.
