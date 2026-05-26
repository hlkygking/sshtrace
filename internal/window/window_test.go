package window_test

import (
	"testing"
	"time"

	"sshtrace/internal/session"
	"sshtrace/internal/window"
)

func makeSession(user string, events []session.Event) *session.Session {
	s := session.New(user, "192.168.1.1")
	s.Events = events
	return s
}

func evt(data string, ts time.Time) session.Event {
	return session.Event{Type: "output", Data: data, Timestamp: ts}
}

func TestNewInvalidWindowSize(t *testing.T) {
	_, err := window.New(0)
	if err == nil {
		t.Fatal("expected error for zero size")
	}
	_, err = window.New(-time.Second)
	if err == nil {
		t.Fatal("expected error for negative size")
	}
}

func TestNewValidWindowSize(t *testing.T) {
	agg, err := window.New(time.Minute)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if agg == nil {
		t.Fatal("expected non-nil aggregator")
	}
}

func TestAggregateNilSessionReturnsError(t *testing.T) {
	agg, _ := window.New(time.Minute)
	_, err := agg.Aggregate(nil)
	if err == nil {
		t.Fatal("expected error for nil session")
	}
}

func TestAggregateEmptySession(t *testing.T) {
	agg, _ := window.New(time.Minute)
	s := makeSession("alice", nil)
	results, err := agg.Aggregate(s)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(results) != 0 {
		t.Fatalf("expected 0 results, got %d", len(results))
	}
}

func TestAggregateSingleWindow(t *testing.T) {
	agg, _ := window.New(time.Minute)
	now := time.Now().Truncate(time.Minute)
	s := makeSession("bob", []session.Event{
		evt("ls", now.Add(5*time.Second)),
		evt("pwd", now.Add(10*time.Second)),
	})
	results, err := agg.Aggregate(s)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(results) != 1 {
		t.Fatalf("expected 1 window, got %d", len(results))
	}
	if results[0].EventCount != 2 {
		t.Errorf("expected 2 events, got %d", results[0].EventCount)
	}
	if len(results[0].Commands) != 2 {
		t.Errorf("expected 2 commands, got %d", len(results[0].Commands))
	}
}

func TestAggregateMultipleWindows(t *testing.T) {
	agg, _ := window.New(time.Minute)
	base := time.Now().Truncate(time.Minute)
	s := makeSession("carol", []session.Event{
		evt("ls", base.Add(10*time.Second)),
		evt("cd", base.Add(90*time.Second)), // second window
		evt("cat", base.Add(150*time.Second)), // third window
	})
	results, err := agg.Aggregate(s)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(results) != 3 {
		t.Fatalf("expected 3 windows, got %d", len(results))
	}
}

func TestPresetAggregatorKnown(t *testing.T) {
	for _, name := range []string{"second", "minute", "hour"} {
		agg := window.PresetAggregator(name)
		if agg == nil {
			t.Errorf("expected non-nil aggregator for preset %q", name)
		}
	}
}

func TestPresetAggregatorUnknownFallsBack(t *testing.T) {
	agg := window.PresetAggregator("unknown")
	if agg == nil {
		t.Fatal("expected non-nil fallback aggregator")
	}
}
