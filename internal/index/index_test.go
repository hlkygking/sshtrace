package index_test

import (
	"testing"
	"time"

	"sshtrace/internal/index"
	"sshtrace/internal/session"
)

func makeSession(user, ip string, start time.Time) *session.Session {
	s := session.New(user, ip)
	s.StartedAt = start
	return s
}

var base = time.Date(2024, 1, 1, 12, 0, 0, 0, time.UTC)

func buildIndex() *index.Index {
	sessions := []*session.Session{
		makeSession("alice", "10.0.0.1", base),
		makeSession("alice", "10.0.0.2", base.Add(1*time.Hour)),
		makeSession("bob", "10.0.0.1", base.Add(2*time.Hour)),
		makeSession("carol", "192.168.1.5", base.Add(3*time.Hour)),
	}
	return index.Build(sessions)
}

func TestSize(t *testing.T) {
	idx := buildIndex()
	if idx.Size() != 4 {
		t.Fatalf("expected 4 sessions, got %d", idx.Size())
	}
}

func TestLookupUserFound(t *testing.T) {
	idx := buildIndex()
	results := idx.LookupUser("alice")
	if len(results) != 2 {
		t.Fatalf("expected 2 sessions for alice, got %d", len(results))
	}
}

func TestLookupUserCaseInsensitive(t *testing.T) {
	idx := buildIndex()
	results := idx.LookupUser("ALICE")
	if len(results) != 2 {
		t.Fatalf("expected 2 sessions for ALICE, got %d", len(results))
	}
}

func TestLookupUserNotFound(t *testing.T) {
	idx := buildIndex()
	results := idx.LookupUser("unknown")
	if len(results) != 0 {
		t.Fatalf("expected 0 sessions, got %d", len(results))
	}
}

func TestLookupIPFound(t *testing.T) {
	idx := buildIndex()
	results := idx.LookupIP("10.0.0.1")
	if len(results) != 2 {
		t.Fatalf("expected 2 sessions for 10.0.0.1, got %d", len(results))
	}
}

func TestLookupIPNotFound(t *testing.T) {
	idx := buildIndex()
	results := idx.LookupIP("1.2.3.4")
	if len(results) != 0 {
		t.Fatalf("expected 0 sessions, got %d", len(results))
	}
}

func TestLookupRange(t *testing.T) {
	idx := buildIndex()
	from := base.Add(30 * time.Minute)
	to := base.Add(2*time.Hour + 30*time.Minute)
	results := idx.LookupRange(from, to)
	if len(results) != 2 {
		t.Fatalf("expected 2 sessions in range, got %d", len(results))
	}
}

func TestLookupRangeEmpty(t *testing.T) {
	idx := buildIndex()
	from := base.Add(10 * time.Hour)
	to := base.Add(11 * time.Hour)
	results := idx.LookupRange(from, to)
	if len(results) != 0 {
		t.Fatalf("expected 0 sessions, got %d", len(results))
	}
}
