package checkpoint_test

import (
	"testing"
	"time"

	"github.com/sshtrace/sshtrace/internal/checkpoint"
	"github.com/sshtrace/sshtrace/internal/session"
)

func makeSession(id, user, ip string) *session.Session {
	s := session.New(user, ip)
	s.ID = id
	return s
}

func TestSaveAndLoad(t *testing.T) {
	tr := checkpoint.New()
	s := makeSession("sess-1", "alice", "10.0.0.1")

	if err := tr.Save(s, 5); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	cp, ok := tr.Load("sess-1")
	if !ok {
		t.Fatal("expected checkpoint to exist")
	}
	if cp.EventIndex != 5 {
		t.Errorf("expected index 5, got %d", cp.EventIndex)
	}
	if cp.SessionID != "sess-1" {
		t.Errorf("expected session id sess-1, got %s", cp.SessionID)
	}
}

func TestLoadMissing(t *testing.T) {
	tr := checkpoint.New()
	_, ok := tr.Load("nonexistent")
	if ok {
		t.Fatal("expected no checkpoint for unknown session")
	}
}

func TestSaveNilSessionReturnsError(t *testing.T) {
	tr := checkpoint.New()
	if err := tr.Save(nil, 0); err == nil {
		t.Fatal("expected error for nil session")
	}
}

func TestSaveNegativeIndexReturnsError(t *testing.T) {
	tr := checkpoint.New()
	s := makeSession("sess-2", "bob", "10.0.0.2")
	if err := tr.Save(s, -1); err == nil {
		t.Fatal("expected error for negative index")
	}
}

func TestDelete(t *testing.T) {
	tr := checkpoint.New()
	s := makeSession("sess-3", "carol", "10.0.0.3")
	_ = tr.Save(s, 2)
	tr.Delete("sess-3")
	_, ok := tr.Load("sess-3")
	if ok {
		t.Fatal("expected checkpoint to be deleted")
	}
}

func TestSavedAtIsRecent(t *testing.T) {
	tr := checkpoint.New()
	before := time.Now().UTC()
	s := makeSession("sess-4", "dave", "10.0.0.4")
	_ = tr.Save(s, 3)
	cp, _ := tr.Load("sess-4")
	if cp.SavedAt.Before(before) {
		t.Errorf("SavedAt %v is before test start %v", cp.SavedAt, before)
	}
}

func TestAll(t *testing.T) {
	tr := checkpoint.New()
	_ = tr.Save(makeSession("s1", "u1", "1.1.1.1"), 0)
	_ = tr.Save(makeSession("s2", "u2", "2.2.2.2"), 1)
	all := tr.All()
	if len(all) != 2 {
		t.Errorf("expected 2 checkpoints, got %d", len(all))
	}
}
