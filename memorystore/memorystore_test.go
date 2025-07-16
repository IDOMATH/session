package memorystore

import (
	"testing"
	"time"
)

func TestMemoryStore_InsertGetDelete(t *testing.T) {
	token := "token"
	expected := "value"
	memStore := NewWithCustomCleanupInterval[string](0)
	memStore.Insert(token, expected, time.Now().Add(time.Minute))

	got, found := memStore.Get(token)
	if !found {
		t.Errorf("Expected found: true, got: %v", found)
	}
	if string(got) != expected {
		t.Errorf("Expected %v, got: %v", expected, got)
	}

	memStore.Delete(token)
	_, found = memStore.Get(token)
	if found {
		t.Errorf("Expected found: false, got: %v", found)
	}
}

func TestMemoryStoreCleanUp(t *testing.T) {
	token := "token"
	expected := "value"
	memStore := NewWithCustomCleanupInterval[string](time.Millisecond * 500)
	defer memStore.stopCleanup()

	memStore.Insert(token, expected, time.Now().Add(time.Millisecond*100))

	got, found := memStore.Get(token)
	if !found {
		t.Errorf("Expected found: true, got: %v", found)
	}
	if string(got) != expected {
		t.Errorf("Expected %v, got: %v", expected, got)
	}

	time.Sleep(time.Millisecond * 600)
	_, found = memStore.Get(token)
	if found {
		t.Errorf("Expected found: false, got: %v", found)
	}
}
