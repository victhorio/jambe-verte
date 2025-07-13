---
title: "Building a Blog with Go"
date: "2024-01-16"
tags: ["golang", "web", "tutorial"]
description: "Learn how I built this blog using Go and minimal dependencies"
---

# Building a Blog with Go

In this post, I'll walk through how I built this blog using Go with minimal dependencies.

## The Stack

- **Go** - The programming language
- **Chi** - Lightweight HTTP router
- **Goldmark** - Markdown processor
- **Alpine.js** - For minimal interactivity

## Project Structure

```
/blog
├── cmd/blog/         # Application entry point
├── internal/         # Internal packages
│   ├── content/      # Content loading
│   ├── handlers/     # HTTP handlers
│   └── cache/        # In-memory cache
├── content/          # Markdown files
├── templates/        # HTML templates
└── static/           # CSS and JS
```

## Key Decisions

### 1. File-based Content

Instead of using a database, all content is stored as markdown files. This makes it easy to:

- Version control content
- Write posts in any text editor
- Deploy without database setup

### 2. In-memory Cache

Posts are loaded once at startup and cached in memory. This provides:

- Lightning-fast page loads
- No filesystem access per request
- Simple implementation

### 3. Single Binary Deployment

The entire blog compiles to a single binary, making deployment as simple as:

```bash
scp blog server:/var/www/
ssh server "systemctl restart blog"
```

## Performance

With this setup, pages load in under 10ms and the memory footprint is minimal. Perfect for a personal blog!

## Next Steps

- Add syntax highlighting
- Implement search
- Create an admin interface

The code is clean, fast, and easy to understand. Exactly what a blog should be!