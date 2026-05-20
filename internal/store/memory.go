package store

import (
	"context"
	"sync"
	"errors"
)

type MemoryStore struct {
	mu sync.Mutex	
	ContentBySlug map[string]string
}

func NewMemoryStore() (*MemoryStore) {
	ms := &MemoryStore{
		ContentBySlug:  make(map[string]string, 0),
	}

	return ms	
}

func (m *MemoryStore)  Save(ctx context.Context, slug string, content string) error {
	m.mu.Lock()
	defer m.mu.Unlock()
	_, ok := m.ContentBySlug[slug]
	if (ok) {
		return errors.New("Attempted to overwrite stored value")
	}

	m.ContentBySlug[slug] = content
	return nil
}

func (m *MemoryStore) Get(ctx context.Context, slug string) (string, error) {
	m.mu.Lock()
	defer m.mu.Unlock()

	c, ok := m.ContentBySlug[slug]
	if (!ok) {
		return "", errors.New("Slug not found")
	}

	return c, nil
}

// Deletes a specified value from the store
// Returns an error in order to fulfill the interface
func (m *MemoryStore) Delete(ctx context.Context, slug string) error {
	m.mu.Lock()
	defer m.mu.Unlock()

	delete(m.ContentBySlug, slug)

	return nil 
}
