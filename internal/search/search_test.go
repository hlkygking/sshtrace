package search_test

import (
	"testing"
	"time"

	"sshtrace/internal/search"
	"sshtrace/internal/session"
)

func makeSession(id, user string, events []session.Event) *session.Session {
	s := session.New(id, user, "192.168.1.1")
	for _, e := range events {
		s.AddEvent(e)
	}
	return s
}

func evt(typ, data string) session.Event {
	return session.Event{Type: typ, Data: data, Timestamp: time.Now()}
}

func TestSearchFindsMatch(t *testing.T) {
	s := makeSession("s1", "alice", []session.Event{
		evt("input", "ls -la"),
		evt("input", "sudo reboot"),
	})
	results := search.Search([]*session.Session{s}, search.Options{Query: "sudo"})
	if len(results) != 1 {
		t.Fatalf("expected 1 result, got %d", len(results))
	}
	if results[0].Event.Data != "sudo reboot" {
		t.Errorf("unexpected event data: %s", results[0].Event.Data)
	}
}

func TestSearchNoMatch(t *testing.T) {
	s := makeSession("s2", "bob", []session.Event{
		evt("input", "ls -la"),
	})
	results := search.Search([]*session.Session{s}, search.Options{Query: "sudo"})
	if len(results) != 0 {
		t.Fatalf("expected 0 results, got %d", len(results))
	}
}

func TestSearchCaseInsensitive(t *testing.T) {
	s := makeSession("s3", "carol", []session.Event{
		evt("input", "SUDO apt update"),
	})
	results := search.Search([]*session.Session{s}, search.Options{Query: "sudo", CaseSensitive: false})
	if len(results) != 1 {
		t.Fatalf("expected 1 result, got %d", len(results))
	}
}

func TestSearchCaseSensitiveMiss(t *testing.T) {
	s := makeSession("s4", "dave", []session.Event{
		evt("input", "SUDO apt update"),
	})
	results := search.Search([]*session.Session{s}, search.Options{Query: "sudo", CaseSensitive: true})
	if len(results) != 0 {
		t.Fatalf("expected 0 results, got %d", len(results))
	}
}

func TestSearchFilterByEventType(t *testing.T) {
	s := makeSession("s5", "eve", []session.Event{
		evt("input", "sudo ls"),
		evt("output", "sudo: command output"),
	})
	results := search.Search([]*session.Session{s}, search.Options{Query: "sudo", EventType: "input"})
	if len(results) != 1 {
		t.Fatalf("expected 1 result, got %d", len(results))
	}
	if results[0].Event.Type != "input" {
		t.Errorf("expected input event, got %s", results[0].Event.Type)
	}
}

func TestSearchMultipleSessions(t *testing.T) {
	s1 := makeSession("s6", "frank", []session.Event{evt("input", "sudo reboot")})
	s2 := makeSession("s7", "grace", []session.Event{evt("input", "sudo shutdown")})
	s3 := makeSession("s8", "henry", []session.Event{evt("input", "ls")})

	results := search.Search([]*session.Session{s1, s2, s3}, search.Options{Query: "sudo"})
	if len(results) != 2 {
		t.Fatalf("expected 2 results, got %d", len(results))
	}
}
