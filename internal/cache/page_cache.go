package cache

import (
	"sync"
)

// PageCache stores pre-rendered HTML pages
type PageCache struct {
	mu    sync.RWMutex
	pages map[string][]byte
}

func NewPageCache() *PageCache {
	return &PageCache{
		pages: make(map[string][]byte),
	}
}

// Get retrieves a cached page by its route path
func (pc *PageCache) Get(path string) ([]byte, bool) {
	pc.mu.RLock()
	defer pc.mu.RUnlock()

	page, ok := pc.pages[path]
	return page, ok
}

// Set stores a rendered page
func (pc *PageCache) Set(path string, content []byte) {
	pc.mu.Lock()
	defer pc.mu.Unlock()

	pc.pages[path] = content
}
