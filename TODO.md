# TODO

## Implement Blog Handlers

The following handlers in `internal/handlers/blog.go` currently return "WIP" placeholders:

- [ ] `ListPosts` — Render list of all published posts at `/posts`
- [ ] `ShowPost` — Render individual post at `/blog/{slug}`
- [ ] `ShowPage` — Render static pages at `/{page}`
- [ ] `PostsByTag` — Filter and list posts by tag at `/tag/{tag}`

## Create Missing Templates

Only `base.html` and `home.html` exist. Need templates for the handlers above:

- [ ] `templates/posts.html` — Blog listing page
- [ ] `templates/post.html` — Individual post page
- [ ] `templates/page.html` — Static page template
- [ ] `templates/tag.html` — Posts filtered by tag (or reuse posts.html with conditional heading)

## Fix Navigation

- [ ] Update nav links in `templates/base.html` — currently hardcoded to `"/"` instead of actual routes

## Content

- [ ] Publish at least one non-draft post to test the full flow
- [ ] Consider adding a `--drafts` flag to `jv-helper` or server to include drafts in dev mode

## Nice to Have

- [ ] Add tests for handlers and content loading
- [ ] RSS feed currently hardcoded to 20 posts — make configurable if needed
