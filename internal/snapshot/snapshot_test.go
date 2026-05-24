package snapshot_test

import (
	"os"
	"path/filepath"
	"testing"
	"time"

	"github.com/sshtrace/sshtrace/internal/session"
	"github.com/sshtrace/sshtrace/internal/snapshot"
)

func makeSession(id, user, ip string) *session.Session {
	s := &session.Session{
		ID:        id,
		User:      user,
		ClientIP:  ip,
		StartedAt: time.Now().UTC(),
	}
	s.Events = append(s.Events, session.Event{
		Kind:      "output",
		Data:      "hello",
		Timestamp: time.Now().UTC(),
	})
	return s
}

func TestNewEmptyDirReturnsError(t *testing.T) {
	_, err := snapshot.New("")
	if err == nil {
		t.Fatal("expected error for empty dir")
	}
}

func TestNewCreatesDirectory(t *testing.T) {
	dir := filepath.Join(t.TempDir(), "snaps")
	_, err := snapshot.New(dir)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if _, err := os.Stat(dir); os.IsNotExist(err) {
		t.Fatal("expected directory to be created")
	}
}

func TestTakeNilSessionReturnsError(t *testing.T) {
	mgr, _ := snapshot.New(t.TempDir())
	_, err := mgr.Take(nil)
	if err == nil {
		t.Fatal("expected error for nil session")
	}
}

func TestTakeAndLoadRoundtrip(t *testing.T) {
	mgr, _ := snapshot.New(t.TempDir())
	sess := makeSession("abc123", "alice", "10.0.0.1")

	snap, err := mgr.Take(sess)
	if err != nil {
		t.Fatalf("Take: %v", err)
	}
	if snap.SessionID != sess.ID {
		t.Errorf("SessionID: got %q, want %q", snap.SessionID, sess.ID)
	}
	if snap.EventCount != 1 {
		t.Errorf("EventCount: got %d, want 1", snap.EventCount)
	}

	loaded, err := mgr.Load(sess.ID)
	if err != nil {
		t.Fatalf("Load: %v", err)
	}
	if loaded.Session.User != "alice" {
		t.Errorf("User: got %q, want %q", loaded.Session.User, "alice")
	}
	if loaded.EventCount != snap.EventCount {
		t.Errorf("EventCount mismatch after reload")
	}
}

func TestLoadMissingReturnsError(t *testing.T) {
	mgr, _ := snapshot.New(t.TempDir())
	_, err := mgr.Load("nonexistent")
	if err == nil {
		t.Fatal("expected error for missing snapshot")
	}
}

func TestLoadEmptyIDReturnsError(t *testing.T) {
	mgr, _ := snapshot.New(t.TempDir())
	_, err := mgr.Load("")
	if err == nil {
		t.Fatal("expected error for empty session ID")
	}
}

func TestDeleteRemovesFile(t *testing.T) {
	dir := t.TempDir()
	mgr, _ := snapshot.New(dir)
	sess := makeSession("del1", "bob", "192.168.1.1")

	if _, err := mgr.Take(sess); err != nil {
		t.Fatalf("Take: %v", err)
	}
	if err := mgr.Delete(sess.ID); err != nil {
		t.Fatalf("Delete: %v", err)
	}
	_, err := mgr.Load(sess.ID)
	if err == nil {
		t.Fatal("expected error after delete")
	}
}

func TestDeleteNonExistentIsNoOp(t *testing.T) {
	mgr, _ := snapshot.New(t.TempDir())
	if err := mgr.Delete("ghost"); err != nil {
		t.Fatalf("unexpected error deleting non-existent snapshot: %v", err)
	}
}
