package memorystore

import (
	"sync"
	"time"
)

type MemoryStore struct {
	items           map[string]item
	mu              sync.RWMutex
	stopCleanupChan chan bool
}

type item struct {
	obj       []byte
	expiresAt int64
}

// New creates a new memory store with the default cleanup interval (1 minute)
func New() *MemoryStore {
	mem := NewWithCustomCleanupInterval(time.Minute)

	return mem
}

// NewWithCustomCleanupInterval creates a new memory store with a user defined cleanup interval
// sending 0 as the interval will not start the cleanup goroutine, so memorystore
// will persist for as long as server runs.
func NewWithCustomCleanupInterval(interval time.Duration) *MemoryStore {
	mem := &MemoryStore{
		items: make(map[string]item),
	}
	if interval != 0 {
		go mem.startCleanup(interval)
	}

	return mem
}

// Insert stores a token in the memory store.
func (s *MemoryStore) Insert(token string, b []byte, expiresAt time.Time) error {
	s.mu.Lock()
	s.items[token] = item{
		obj:       b,
		expiresAt: expiresAt.UnixNano(),
	}
	s.mu.Unlock()

	return nil
}

// Get retrieves a token from the memory store.
func (s *MemoryStore) Get(token string) (b []byte, found bool) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	item, found := s.items[token]
	if !found {
		return nil, false
	}

	if time.Now().UnixNano() > item.expiresAt {
		return nil, false
	}
	return item.obj, true
}

// Delete removes a given token from the memory store
func (s *MemoryStore) Delete(token string) error {
	s.mu.Lock()
	delete(s.items, token)
	s.mu.Unlock()

	return nil
}

// startCleanUp runs deleteExpiredTokens() at a given interval
func (s *MemoryStore) startCleanup(interval time.Duration) {
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

// deleteExpiredTokens deletes tokens that have expired.
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

// stopCleanup will tell the startCleanup() function to stop.
// This is really only used for tests.
func (s *MemoryStore) stopCleanup() {
	if s.stopCleanupChan != nil {
		s.stopCleanupChan <- true
	}
}
