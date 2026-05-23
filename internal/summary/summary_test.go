package summary_test

import (
	"testing"
	"time"

	"sshtrace/internal/session"
	"sshtrace/internal/summary"
)

func makeSession(user, ip string, events int, duration time.Duration) *session.Session {
	s := session.New(user, ip)
	for i := 0; i < events; i++ {
		s.AddEvent(session.Event{Type: "output", Data: "cmd"})
	}
	s.Close()
	// Shift EndedAt to simulate duration
	s.EndedAt = s.StartedAt.Add(duration)
	return s
}

func TestGenerateEmpty(t *testing.T) {
	r := summary.Generate(nil)
	if r.TotalSessions != 0 {
		t.Errorf("expected 0 sessions, got %d", r.TotalSessions)
	}
}

func TestGenerateTotalSessions(t *testing.T) {
	sessions := []*session.Session{
		makeSession("alice", "1.2.3.4", 2, 10*time.Second),
		makeSession("bob", "5.6.7.8", 3, 20*time.Second),
	}
	r := summary.Generate(sessions)
	if r.TotalSessions != 2 {
		t.Errorf("expected 2 sessions, got %d", r.TotalSessions)
	}
}

func TestGenerateUniqueUsers(t *testing.T) {
	sessions := []*session.Session{
		makeSession("alice", "1.1.1.1", 1, 5*time.Second),
		makeSession("alice", "2.2.2.2", 1, 5*time.Second),
		makeSession("bob", "3.3.3.3", 1, 5*time.Second),
	}
	r := summary.Generate(sessions)
	if len(r.UniqueUsers) != 2 {
		t.Errorf("expected 2 unique users, got %d", len(r.UniqueUsers))
	}
}

func TestGenerateTotalEvents(t *testing.T) {
	sessions := []*session.Session{
		makeSession("alice", "1.1.1.1", 4, 5*time.Second),
		makeSession("bob", "2.2.2.2", 6, 5*time.Second),
	}
	r := summary.Generate(sessions)
	if r.TotalEvents != 10 {
		t.Errorf("expected 10 total events, got %d", r.TotalEvents)
	}
}

func TestGenerateDurationStats(t *testing.T) {
	sessions := []*session.Session{
		makeSession("alice", "1.1.1.1", 1, 10*time.Second),
		makeSession("bob", "2.2.2.2", 1, 30*time.Second),
	}
	r := summary.Generate(sessions)
	if r.LongestSession != 30*time.Second {
		t.Errorf("expected longest 30s, got %v", r.LongestSession)
	}
	if r.ShortestSession != 10*time.Second {
		t.Errorf("expected shortest 10s, got %v", r.ShortestSession)
	}
	if r.AvgDuration != 20*time.Second {
		t.Errorf("expected avg 20s, got %v", r.AvgDuration)
	}
}

func TestGenerateSingleSession(t *testing.T) {
	sessions := []*session.Session{
		makeSession("alice", "1.1.1.1", 3, 15*time.Second),
	}
	r := summary.Generate(sessions)
	if r.LongestSession != 15*time.Second {
		t.Errorf("expected longest 15s, got %v", r.LongestSession)
	}
	if r.ShortestSession != 15*time.Second {
		t.Errorf("expected shortest 15s, got %v", r.ShortestSession)
	}
	if r.AvgDuration != 15*time.Second {
		t.Errorf("expected avg 15s, got %v", r.AvgDuration)
	}
}
