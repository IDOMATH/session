package main

import (
	"time"
)

type Session struct {
	Store  Store
	Id     string
	Values map[interface{}]interface{}
}

type Store interface {
	Insert(token string, b []byte, expiresAt time.Time) (err error)
	Get(token string) (b []byte, exists bool, err error)
	Delete(token string) (err error)
}
