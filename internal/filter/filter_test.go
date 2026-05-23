package filter_test

import (
	"testing"
	"time"

	"sshtrace/internal/filter"
	"sshtrace/internal/session"
)

func makeSession(user, ip string, start time.Time, eventCount int) *session.Session {
	s := session.New(user, ip)
	s.StartedAt = start
	for i := 0; i < eventCount; i++ {
		s.AddEvent(session.Event{Kind: "output", Data: []byte("x"), At: start})
	}
	return s
}

var base = time.Date(2024, 1, 15, 12, 0, 0, 0, time.UTC)

var corpus = []*session.Session{
	makeSession("alice", "10.0.0.1", base, 3),
	makeSession("bob", "10.0.0.2", base.Add(time.Hour), 1),
	makeSession("alice", "10.0.0.3", base.Add(2*time.Hour), 5),
	makeSession("carol", "10.0.0.1", base.Add(3*time.Hour), 0),
}

func TestFilterByUser(t *testing.T) {
	got := filter.Filter(corpus, filter.Criteria{UserName: "alice"})
	if len(got) != 2 {
		t.Fatalf("expected 2 sessions, got %d", len(got))
	}
}

func TestFilterByIP(t *testing.T) {
	got := filter.Filter(corpus, filter.Criteria{RemoteIP: "10.0.0.1"})
	if len(got) != 2 {
		t.Fatalf("expected 2 sessions, got %d", len(got))
	}
}

func TestFilterBySince(t *testing.T) {
	got := filter.Filter(corpus, filter.Criteria{Since: base.Add(90 * time.Minute)})
	if len(got) != 2 {
		t.Fatalf("expected 2 sessions, got %d", len(got))
	}
}

func TestFilterByUntil(t *testing.T) {
	got := filter.Filter(corpus, filter.Criteria{Until: base.Add(30 * time.Minute)})
	if len(got) != 1 {
		t.Fatalf("expected 1 session, got %d", len(got))
	}
}

func TestFilterByMinEvents(t *testing.T) {
	got := filter.Filter(corpus, filter.Criteria{MinEvents: 3})
	if len(got) != 2 {
		t.Fatalf("expected 2 sessions, got %d", len(got))
	}
}

func TestFilterNoCriteria(t *testing.T) {
	got := filter.Filter(corpus, filter.Criteria{})
	if len(got) != len(corpus) {
		t.Fatalf("expected all %d sessions, got %d", len(corpus), len(got))
	}
}

func TestFilterCombined(t *testing.T) {
	got := filter.Filter(corpus, filter.Criteria{UserName: "alice", MinEvents: 4})
	if len(got) != 1 {
		t.Fatalf("expected 1 session, got %d", len(got))
	}
	if got[0].UserName != "alice" || len(got[0].Events) < 4 {
		t.Fatal("wrong session returned")
	}
}
