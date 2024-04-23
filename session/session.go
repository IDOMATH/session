package session

import (
	"sync"
	"time"
)

type MemoryStore struct {
	items map[string]item
	mu    sync.RWMutex
}

func NewMemoryStore() *MemoryStore {
	return &MemoryStore{
		items: make(map[string]item),
	}
}

func (s *MemoryStore) Insert(token string, b []byte, expiresAt time.Time) error {
	s.mu.Lock()
	s.items[token] = item{
		obj:       b,
		expiresAt: expiresAt.UnixNano(),
	}
	s.mu.Unlock()

	return nil
}

func (s *MemoryStore) Get(token string) (b []byte, found bool, err error) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	item, found := s.items[token]
	if !found {
		return nil, false, nil
	}

	if time.Now().UnixNano() > item.expiresAt {
		return nil, false, nil
	}
	return item.obj, true, nil
}

func (s *MemoryStore) Delete(token string) error {
	s.mu.Lock()
	delete(s.items, token)
	s.mu.Unlock()

	return nil
}
