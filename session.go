package main

import (
	"sync"
	"time"
)

type Store interface {
	Insert(token string, b []byte, expiresAt time.Time) (err error)
	Get(token string) (b []byte, exists bool, err error)
	Delete(token string) (err error)
}

type item struct {
	obj       []byte
	expiresAt int64
}

type MemoryStore struct {
	items map[string]item
	mu    sync.RWMutex
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
