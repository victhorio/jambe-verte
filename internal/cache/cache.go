package cache

import (
	"sync"

	"github.com/victhor/jv-blog/internal/content"
)

type Cache struct {
	mu    sync.RWMutex
	posts map[string]*content.Post
	pages map[string]*content.Post
	tags  map[string][]*content.Post
	all   []*content.Post
}

func New(posts []*content.Post, pages []*content.Post) *Cache {
	c := &Cache{
		posts: make(map[string]*content.Post),
		pages: make(map[string]*content.Post),
		tags:  make(map[string][]*content.Post),
		all:   posts,
	}

	// Index posts by slug
	for _, post := range posts {
		c.posts[post.Slug] = post

		// Index by tags
		for _, tag := range post.Tags {
			c.tags[tag] = append(c.tags[tag], post)
		}
	}

	// Index pages by slug
	for _, page := range pages {
		c.pages[page.Slug] = page
	}

	return c
}

func (c *Cache) GetPost(slug string) (*content.Post, bool) {
	c.mu.RLock()
	defer c.mu.RUnlock()
	post, ok := c.posts[slug]
	return post, ok
}

func (c *Cache) GetPage(slug string) (*content.Post, bool) {
	c.mu.RLock()
	defer c.mu.RUnlock()
	page, ok := c.pages[slug]
	return page, ok
}

func (c *Cache) GetAllPosts() []*content.Post {
	c.mu.RLock()
	defer c.mu.RUnlock()
	// Return a copy to prevent race conditions
	posts := make([]*content.Post, len(c.all))
	copy(posts, c.all)
	return posts
}

func (c *Cache) GetPostsByTag(tag string) []*content.Post {
	c.mu.RLock()
	defer c.mu.RUnlock()
	// Return a copy to prevent race conditions
	tagPosts := c.tags[tag]
	if tagPosts == nil {
		return nil
	}
	posts := make([]*content.Post, len(tagPosts))
	copy(posts, tagPosts)
	return posts
}

func (c *Cache) GetAllTags() []string {
	c.mu.RLock()
	defer c.mu.RUnlock()

	tags := make([]string, 0, len(c.tags))
	for tag := range c.tags {
		tags = append(tags, tag)
	}
	return tags
}
