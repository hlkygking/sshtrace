package dedupe_test

import (
	"testing"
	"time"

	"sshtrace/internal/dedupe"
	"sshtrace/internal/session"
)

func makeSession(events ...session.Event) *session.Session {
	s := &session.Session{
		ID:        "test-id",
		User:      "alice",
		CreatedAt: time.Now(),
		Events:    events,
	}
	return s
}

func evt(typ, data string) session.Event {
	return session.Event{Type: typ, Data: data, At: time.Now()}
}

func TestNewInvalidWindowSize(t *testing.T) {
	_, err := dedupe.New(0)
	if err == nil {
		t.Fatal("expected error for windowSize=0")
	}
}

func TestNewValidWindowSize(t *testing.T) {
	d, err := dedupe.New(1)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if d.WindowSize != 1 {
		t.Errorf("expected WindowSize=1, got %d", d.WindowSize)
	}
}

func TestApplyNilSession(t *testing.T) {
	d, _ := dedupe.New(1)
	_, err := d.Apply(nil)
	if err == nil {
		t.Fatal("expected error for nil session")
	}
}

func TestApplyNoConsecutiveDuplicates(t *testing.T) {
	d, _ := dedupe.New(1)
	s := makeSession(evt("output", "ls"), evt("output", "pwd"), evt("output", "whoami"))
	out, err := d.Apply(s)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(out.Events) != 3 {
		t.Errorf("expected 3 events, got %d", len(out.Events))
	}
}

func TestApplyRemovesConsecutiveDuplicates(t *testing.T) {
	d, _ := dedupe.New(1)
	s := makeSession(
		evt("output", "ls"),
		evt("output", "ls"),
		evt("output", "pwd"),
	)
	out, err := d.Apply(s)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(out.Events) != 2 {
		t.Errorf("expected 2 events after dedup, got %d", len(out.Events))
	}
}

func TestApplyWindowSizeTwo(t *testing.T) {
	d, _ := dedupe.New(2)
	s := makeSession(
		evt("output", "ls"),
		evt("output", "pwd"),
		evt("output", "ls"), // within window of 2 — should be removed
	)
	out, err := d.Apply(s)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(out.Events) != 2 {
		t.Errorf("expected 2 events with window=2, got %d", len(out.Events))
	}
}

func TestApplyEmptySession(t *testing.T) {
	d, _ := dedupe.New(1)
	s := makeSession()
	out, err := d.Apply(s)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(out.Events) != 0 {
		t.Errorf("expected 0 events, got %d", len(out.Events))
	}
}
