package session_test

import (
	"testing"
	"time"

	"github.com/sshtrace/sshtrace/internal/session"
)

func TestNew(t *testing.T) {
	s := session.New("alice", "192.168.1.10")

	if s.ID == "" {
		t.Error("expected non-empty session ID")
	}
	if s.User != "alice" {
		t.Errorf("expected user 'alice', got '%s'", s.User)
	}
	if s.RemoteIP != "192.168.1.10" {
		t.Errorf("expected remote IP '192.168.1.10', got '%s'", s.RemoteIP)
	}
	if s.EndedAt != nil {
		t.Error("expected EndedAt to be nil for a new session")
	}
	if len(s.Events) != 0 {
		t.Errorf("expected 0 events, got %d", len(s.Events))
	}
}

func TestAddEvent(t *testing.T) {
	s := session.New("bob", "10.0.0.1")
	s.AddEvent(session.EventTypeConnect, "session started")
	s.AddEvent(session.EventTypeCommand, "ls -la")

	if len(s.Events) != 2 {
		t.Fatalf("expected 2 events, got %d", len(s.Events))
	}
	if s.Events[0].Type != session.EventTypeConnect {
		t.Errorf("expected first event type 'connect', got '%s'", s.Events[0].Type)
	}
	if s.Events[1].Data != "ls -la" {
		t.Errorf("expected second event data 'ls -la', got '%s'", s.Events[1].Data)
	}
}

func TestClose(t *testing.T) {
	s := session.New("carol", "172.16.0.5")
	if s.EndedAt != nil {
		t.Fatal("session should not be closed initially")
	}

	s.Close()

	if s.EndedAt == nil {
		t.Error("expected EndedAt to be set after Close()")
	}
}

func TestDuration(t *testing.T) {
	s := session.New("dave", "10.10.10.10")
	time.Sleep(10 * time.Millisecond)

	dur := s.Duration()
	if dur < 10*time.Millisecond {
		t.Errorf("expected duration >= 10ms, got %v", dur)
	}

	s.Close()
	closed := s.Duration()
	time.Sleep(20 * time.Millisecond)
	// Duration should not grow after session is closed
	if s.Duration() != closed {
		t.Error("duration should be fixed after session is closed")
	}
}

func TestUniqueIDs(t *testing.T) {
	s1 := session.New("user1", "1.1.1.1")
	s2 := session.New("user2", "2.2.2.2")
	if s1.ID == s2.ID {
		t.Error("expected unique session IDs")
	}
}
