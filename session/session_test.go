package session

import (
	"testing"
	"time"
)

func TestMemoryStore_InsertGet(t *testing.T) {
	token := "token"
	expected := "value"
	memStore := NewMemoryStoreWithCustomCleanupInterval(0)
	memStore.Insert(token, []byte(expected), time.Now().Add(time.Minute))

	got, found, _ := memStore.Get(token)
	if !found {
		t.Errorf("Expected found: true, got: %v", found)
	}
	if string(got) != expected {
		t.Errorf("Expected %v, got: %v", expected, got)
	}
}
