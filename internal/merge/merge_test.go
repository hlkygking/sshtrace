package merge_test

import (
	"testing"
	"time"

	"sshtrace/internal/merge"
	"sshtrace/internal/session"
)

func makeSession(user, ip string, events []session.Event) *session.Session {
	s := session.New(user, ip)
	for _, e := range events {
		s.AddEvent(e)
	}
	return s
}

func evt(kind, data string, ts time.Time) session.Event {
	return session.Event{Kind: kind, Data: data, Timestamp: ts}
}

var t0 = time.Date(2024, 1, 1, 0, 0, 0, 0, time.UTC)

func TestMergeEmptySliceReturnsError(t *testing.T) {
	m, _ := merge.New(merge.DefaultOptions())
	_, err := m.Merge(nil)
	if err == nil {
		t.Fatal("expected error for empty slice")
	}
}

func TestMergeNilSessionReturnsError(t *testing.T) {
	m, _ := merge.New(merge.DefaultOptions())
	s := makeSession("alice", "1.2.3.4", nil)
	_, err := m.Merge([]*session.Session{s, nil})
	if err == nil {
		t.Fatal("expected error for nil session")
	}
}

func TestNewNegativeWindowReturnsError(t *testing.T) {
	_, err := merge.New(merge.Options{DeduplicateWindow: -time.Second})
	if err == nil {
		t.Fatal("expected error for negative window")
	}
}

func TestMergeSingleSession(t *testing.T) {
	events := []session.Event{
		evt("output", "hello", t0),
		evt("input", "ls", t0.Add(time.Second)),
	}
	s := makeSession("alice", "1.2.3.4", events)
	m, _ := merge.New(merge.DefaultOptions())
	out, err := m.Merge([]*session.Session{s})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(out.Events) != 2 {
		t.Fatalf("expected 2 events, got %d", len(out.Events))
	}
}

func TestMergeOrdersByTimestamp(t *testing.T) {
	s1 := makeSession("alice", "1.1.1.1", []session.Event{
		evt("input", "first", t0),
		evt("input", "third", t0.Add(2*time.Second)),
	})
	s2 := makeSession("bob", "2.2.2.2", []session.Event{
		evt("input", "second", t0.Add(time.Second)),
	})

	m, _ := merge.New(merge.DefaultOptions())
	out, err := m.Merge([]*session.Session{s1, s2})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(out.Events) != 3 {
		t.Fatalf("expected 3 events, got %d", len(out.Events))
	}
	expected := []string{"first", "second", "third"}
	for i, e := range out.Events {
		if e.Data != expected[i] {
			t.Errorf("event %d: got %q, want %q", i, e.Data, expected[i])
		}
	}
}

func TestMergePreservesOriginID(t *testing.T) {
	s1 := makeSession("alice", "1.1.1.1", []session.Event{evt("input", "cmd", t0)})
	opts := merge.Options{PreserveIDs: true}
	m, _ := merge.New(opts)
	out, err := m.Merge([]*session.Session{s1})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if out.Events[0].Meta["origin_session"] != s1.ID {
		t.Errorf("expected origin_session=%q", s1.ID)
	}
}

func TestMergeDeduplicatesWithinWindow(t *testing.T) {
	shared := evt("input", "ls", t0)
	s1 := makeSession("alice", "1.1.1.1", []session.Event{shared})
	s2 := makeSession("alice", "1.1.1.1", []session.Event{shared})

	opts := merge.Options{DeduplicateWindow: 5 * time.Second}
	m, _ := merge.New(opts)
	out, err := m.Merge([]*session.Session{s1, s2})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(out.Events) != 1 {
		t.Fatalf("expected 1 deduplicated event, got %d", len(out.Events))
	}
}

func TestDefaultOptions(t *testing.T) {
	opts := merge.DefaultOptions()
	if !opts.PreserveIDs {
		t.Error("expected PreserveIDs to be true by default")
	}
	if opts.DeduplicateWindow != 0 {
		t.Error("expected DeduplicateWindow to be zero by default")
	}
}
