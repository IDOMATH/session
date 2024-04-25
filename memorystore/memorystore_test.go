package memorystore

import (
	"testing"
	"time"
)

func TestMemoryStore_InsertGetDelete(t *testing.T) {
	token := "token"
	expected := "value"
	memStore := NewWithCustomCleanupInterval(0)
	memStore.Insert(token, []byte(expected), time.Now().Add(time.Minute))

	got, found, _ := memStore.Get(token)
	if !found {
		t.Errorf("Expected found: true, got: %v", found)
	}
	if string(got) != expected {
		t.Errorf("Expected %v, got: %v", expected, got)
	}

	memStore.Delete(token)
	_, found, _ = memStore.Get(token)
	if found {
		t.Errorf("Expected found: false, got: %v", found)
	}
}

func TestMemoryStoreCleanUp(t *testing.T) {
	token := "token"
	expected := "value"
	memStore := NewWithCustomCleanupInterval(time.Millisecond * 500)
	defer memStore.stopCleanup()

	memStore.Insert(token, []byte(expected), time.Now().Add(time.Millisecond*100))

	got, found, _ := memStore.Get(token)
	if !found {
		t.Errorf("Expected found: true, got: %v", found)
	}
	if string(got) != expected {
		t.Errorf("Expected %v, got: %v", expected, got)
	}

	time.Sleep(time.Millisecond * 600)
	got, found, _ = memStore.Get(token)
	if found {
		t.Errorf("Expected found: false, got: %v", found)
	}
}
