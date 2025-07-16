package memorystore

import (
	"sync"
	"time"
)

type MemoryStore[T any] struct {
	items           map[string]item[T]
	mu              sync.RWMutex
	stopCleanupChan chan bool
}

type item[T any] struct {
	obj       T
	expiresAt int64
}

// New creates a new memory store with the default cleanup interval (1 minute)
func New[T any](dataType T) *MemoryStore[T] {
	mem := NewWithCustomCleanupInterval[T](time.Minute)

	return mem
}

// NewWithCustomCleanupInterval creates a new memory store with a user defined cleanup interval
// sending 0 as the interval will not start the cleanup goroutine, so memorystore
// will persist for as long as server runs.
func NewWithCustomCleanupInterval[T any](interval time.Duration) *MemoryStore[T] {
	mem := &MemoryStore[T]{
		items: make(map[string]item[T]),
	}
	if interval != 0 {
		go mem.startCleanup(interval)
	}

	return mem
}

// Insert stores a token in the memory store.
func (s *MemoryStore[T]) Insert(token string, data T, expiresAt time.Time) error {
	s.mu.Lock()
	s.items[token] = item[T]{
		obj:       data,
		expiresAt: expiresAt.UnixNano(),
	}
	s.mu.Unlock()

	return nil
}

// Get retrieves a token from the memory store.
func (s *MemoryStore[T]) Get(token string) (data T, found bool) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	var val T

	item, found := s.items[token]
	if !found {
		return val, false
	}

	if time.Now().UnixNano() > item.expiresAt {
		return val, false
	}
	return item.obj, true
}

// Delete removes a given token from the memory store
func (s *MemoryStore[T]) Delete(token string) error {
	s.mu.Lock()
	delete(s.items, token)
	s.mu.Unlock()

	return nil
}

// startCleanUp runs deleteExpiredTokens() at a given interval
func (s *MemoryStore[T]) startCleanup(interval time.Duration) {
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
func (s *MemoryStore[T]) deleteExpiredTokens() {
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
func (s *MemoryStore[T]) stopCleanup() {
	if s.stopCleanupChan != nil {
		s.stopCleanupChan <- true
	}
}
