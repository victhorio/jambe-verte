package cache

import (
	"github.com/victhorio/jambe-verte/internal/content"
)

// Cache is immutable after initialization. Individual Cache instances never change
// their contents after creation. To update cached content, create a new Cache instance
// and replace the old one atomically.
type Cache struct {
	posts     map[string]*content.Post
	pages     map[string]*content.Post
	tags      map[string][]*content.Post
	postsFlat []*content.Post
}

func New(posts []*content.Post, pages []*content.Post) *Cache {
	c := &Cache{
		posts:     make(map[string]*content.Post),
		pages:     make(map[string]*content.Post),
		tags:      make(map[string][]*content.Post),
		postsFlat: posts,
	}

	// For each post, index it by slug on `c.posts` and index it
	// by its tag on `c.tags`
	for _, post := range posts {
		c.posts[post.Slug] = post
		for _, tag := range post.Tags {
			c.tags[tag] = append(c.tags[tag], post)
		}
	}

	// For each page, index it by slug on `c.pages`
	for _, page := range pages {
		c.pages[page.Slug] = page
	}
	return c
}

func (c *Cache) GetPost(slug string) (*content.Post, bool) {
	post, ok := c.posts[slug]
	return post, ok
}

func (c *Cache) GetPage(slug string) (*content.Post, bool) {
	page, ok := c.pages[slug]
	return page, ok
}

func (c *Cache) GetPostsByTag(tag string) []*content.Post {
	return c.tags[tag]
}

func (c *Cache) GetPosts() []*content.Post {
	return c.postsFlat
}
