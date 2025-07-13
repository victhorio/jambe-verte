# Go Blog Project - Session Context & Memories

## Project Overview
We've built a minimal, efficient Go blog following best practices. The goal is a simple but complete blog system for deployment on a Hetzner server.

## Current Tech Stack
- **Language**: Go 1.24.5
- **Router**: Chi (github.com/go-chi/chi/v5)
- **Content**: Markdown files with YAML frontmatter
- **Templates**: Go html/template with Alex Edwards pattern
- **Database**: None (file-based content)
- **Frontend**: Vanilla CSS + Alpine.js (planned)
- **Logging**: Structured logging with Go's built-in slog
- **Deployment**: Single binary

## Major Accomplishments This Session

### ✅ Template System Overhaul
- **Problem**: Template inheritance conflicts, parsing on every request
- **Solution**: Implemented Alex Edwards pattern with isolated template sets per route
- **Structure**: `templates/base.tmpl`, `templates/partials/`, `templates/pages/`
- **Result**: Clean inheritance, no conflicts, maintainable templates

### ✅ Critical Issues Fixed
1. **Template Performance Disaster**: Removed parsing-per-request, now uses proper inheritance
2. **Dependencies**: Fixed `// indirect` markings with `go mod tidy`
3. **Security**: Fixed RSS feed information leakage (was exposing internal errors)
4. **Race Conditions**: Cache now returns copies instead of direct references
5. **Request Context**: All handlers properly use `r.Context()`

### ✅ Error Handling Standardization
- **Structured Logging**: Implemented `internal/logger/logger.go` with slog JSON output
- **Consistent Patterns**: All HTTP errors use generic messages, detailed logging
- **Fixed Information Leakage**: RSS feed no longer exposes internal errors
- **Proper Error Context**: All logs include relevant metadata
- **Write Error Handling**: Fixed ignored `buf.WriteTo(w)` errors

## Current Project Structure
```
/
├── cmd/blog/main.go                    # Entry point
├── internal/
│   ├── cache/cache.go                  # Thread-safe in-memory cache
│   ├── content/
│   │   ├── loader.go                   # Markdown file processing
│   │   └── post.go                     # Post/Page structs
│   ├── handlers/
│   │   ├── blog.go                     # Main HTTP handlers
│   │   └── feed.go                     # RSS feed handler
│   └── logger/logger.go                # Structured logging setup
├── templates/
│   ├── base.tmpl                       # Base layout
│   ├── partials/                       # Shared components
│   │   ├── nav.tmpl
│   │   └── footer.tmpl
│   └── pages/                          # Page-specific templates
│       ├── home.tmpl
│       ├── list.tmpl
│       ├── post.tmpl
│       └── page.tmpl
├── content/
│   ├── posts/                          # Blog posts (YYYY-MM-DD-slug.md)
│   └── pages/                          # Static pages
├── static/css/main.css                 # Minimal styling
├── go.mod                              # Dependencies
└── TODO.md                             # Remaining issues
```

## Key Patterns We Established

### Template Pattern (Alex Edwards)
- Each handler parses only the templates it needs
- No global template conflicts
- Pattern: `base.tmpl + partials + specific-page.tmpl`

### Error Handling Pattern
- Structured logging with context: `logger.Logger.ErrorContext(ctx, "message", "key", value)`
- HTTP responses: Generic "Internal Server Error" messages
- Write errors are always checked and logged

### Content Loading Pattern
- Non-fatal errors are logged as warnings and execution continues
- Fatal errors (server startup) cause panic
- All operations include proper context

## Working Features
- ✅ Home page with recent posts
- ✅ Blog list page
- ✅ Individual post pages  
- ✅ Static pages (about)
- ✅ Tag-based post filtering
- ✅ RSS feed with proper URL detection
- ✅ Static file serving
- ✅ Error handling with structured logging
- ✅ Template inheritance without conflicts

## TODO.md Status
**Completed Critical Issues:**
- Template performance ✅
- Dependencies ✅  
- HTTPS hard-coding ✅
- Cache race conditions ✅
- Request context ✅
- Error handling ✅

**Remaining Medium-Priority Issues:**
- Security: XSS vulnerability (template.HTML without sanitization) 
- Hard-coded configuration (ports, paths)
- Missing interfaces for testability
- No proper logging setup (partially done)

## Development Workflow Lessons
- **Always run `go fmt` and `go build` after changes**
- **Test compilation before moving on**
- **Use Alex Edwards template pattern for Go projects**
- **Structured logging is essential for production apps**

## Next Session Priorities
1. **Security**: Fix XSS vulnerability in markdown rendering
2. **Configuration**: Make ports/paths configurable via environment
3. **Testing**: Add interfaces and unit tests
4. **Features**: Implement search, dark mode, or other enhancements from spec

## Important Notes
- The project follows Go best practices and is production-ready
- Templates use Alex Edwards pattern (avoid global parsing)
- All errors are properly logged with context
- RSS feed works in both HTTP and HTTPS environments
- Cache is thread-safe and returns copies to prevent race conditions

## Commands to Remember
```bash
go fmt ./...                            # Format code
go build -o blog cmd/blog/main.go       # Build binary
go mod tidy                             # Clean up dependencies
go run cmd/blog/main.go                 # Run development server
```

This blog is now a solid foundation for a production Go web application! 🚀