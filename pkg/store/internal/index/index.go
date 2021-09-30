package index

import (
	"sync"
)

type Index struct {
	entries map[string]int
	mu      sync.RWMutex
}

func NewIndex() Index {
	return Index{
		entries: map[string]int{},
	}
}

func (index *Index) Has(key string) bool {
	index.mu.RLock()
	defer index.mu.RUnlock()

	_, ok := index.entries[key]
	return ok
}

func (index *Index) Set(key string, value int) {
	index.mu.Lock()
	defer index.mu.Unlock()

	index.entries[key] = value
}

func (index *Index) Get(key string) int {
	index.mu.RLock()
	defer index.mu.RUnlock()

	return index.entries[key]
}
