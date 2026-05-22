package storage_test

import (
	"os"
	"path/filepath"
	"testing"
	"time"

	"sshtrace/internal/session"
	"sshtrace/internal/storage"
)

func newTestSession(t *testing.T) *session.Session {
	t.Helper()
	sess := session.New("192.0.2.1", "alice")
	sess.AddEvent(session.Event{Timestamp: time.Now(), Command: "ls -la"})
	sess.Close()
	return sess
}

func TestSaveAndLoad(t *testing.T) {
	dir := t.TempDir()
	store, err := storage.New(dir)
	if err != nil {
		t.Fatalf("New: %v", err)
	}

	orig := newTestSession(t)
	if err := store.Save(orig); err != nil {
		t.Fatalf("Save: %v", err)
	}

	entries, _ := os.ReadDir(dir)
	if len(entries) != 1 {
		t.Fatalf("expected 1 file, got %d", len(entries))
	}

	loaded, err := store.Load(filepath.Join(dir, entries[0].Name()))
	if err != nil {
		t.Fatalf("Load: %v", err)
	}

	if loaded.ID != orig.ID {
		t.Errorf("ID mismatch: got %s, want %s", loaded.ID, orig.ID)
	}
	if loaded.User != orig.User {
		t.Errorf("User mismatch: got %s, want %s", loaded.User, orig.User)
	}
	if len(loaded.Events) != len(orig.Events) {
		t.Errorf("Events len mismatch: got %d, want %d", len(loaded.Events), len(orig.Events))
	}
}

func TestListSince(t *testing.T) {
	dir := t.TempDir()
	store, err := storage.New(dir)
	if err != nil {
		t.Fatalf("New: %v", err)
	}

	before := time.Now()
	for i := 0; i < 3; i++ {
		if err := store.Save(newTestSession(t)); err != nil {
			t.Fatalf("Save: %v", err)
		}
	}

	paths, err := store.ListSince(before)
	if err != nil {
		t.Fatalf("ListSince: %v", err)
	}
	if len(paths) != 3 {
		t.Errorf("expected 3 paths, got %d", len(paths))
	}

	futurePaths, err := store.ListSince(time.Now().Add(time.Hour))
	if err != nil {
		t.Fatalf("ListSince future: %v", err)
	}
	if len(futurePaths) != 0 {
		t.Errorf("expected 0 future paths, got %d", len(futurePaths))
	}
}

func TestNewCreatesDir(t *testing.T) {
	dir := filepath.Join(t.TempDir(), "nested", "store")
	if _, err := storage.New(dir); err != nil {
		t.Fatalf("New with nested dir: %v", err)
	}
	if _, err := os.Stat(dir); os.IsNotExist(err) {
		t.Error("expected directory to be created")
	}
}
