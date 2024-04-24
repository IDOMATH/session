package session

import (
	"sync"
	"time"
)

type MemoryStore struct {
	items           map[string]item
	mu              sync.RWMutex
	stopCleanupChan chan bool
}

func NewMemoryStore() *MemoryStore {
	mem := NewMemoryStoreWithCustomCleanupInterval(time.Minute)

	return mem
}

func NewMemoryStoreWithCustomCleanupInterval(interval time.Duration) *MemoryStore {
	mem := &MemoryStore{
		items: make(map[string]item),
	}
	if interval != 0 {
		go mem.startCleanUp(interval)
	}

	return mem
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

func (s *MemoryStore) startCleanUp(interval time.Duration) {
	ticker := time.NewTicker(interval)
	for {
		select {
		case <-ticker.C:
			s.deleteExpiredTokens()
		case <-s.stopCleanupChan:
			return
		}
	}
}

func (s *MemoryStore) deleteExpiredTokens() {
	now := time.Now().UnixNano()
	s.mu.Lock()
	for token, val := range s.items {
		if now > val.expiresAt {
			delete(s.items, token)
		}
	}
	s.mu.Unlock()
}

func (s *MemoryStore) stopCleanup() {
	if s.stopCleanupChan != nil {
		s.stopCleanupChan <- true
	}
}
