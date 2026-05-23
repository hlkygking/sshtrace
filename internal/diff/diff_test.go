package diff_test

import (
	"strings"
	"testing"
	"time"

	"sshtrace/internal/diff"
	"sshtrace/internal/session"
)

func makeSession(id string, cmds []string) *session.Session {
	s := session.New(id, "alice", "192.168.1.1")
	for _, c := range cmds {
		s.AddEvent(session.Event{Type: session.EventOutput, Data: c, Timestamp: time.Now()})
	}
	return s
}

func TestEqualSessions(t *testing.T) {
	a := makeSession("aaa", []string{"ls", "pwd"})
	b := makeSession("bbb", []string{"ls", "pwd"})
	res := diff.Compare(a, b)
	if !res.Equal() {
		t.Fatalf("expected equal sessions, got %d delta(s)", len(res.Deltas))
	}
}

func TestChangedEvent(t *testing.T) {
	a := makeSession("aaa", []string{"ls", "pwd"})
	b := makeSession("bbb", []string{"ls", "whoami"})
	res := diff.Compare(a, b)
	if res.Equal() {
		t.Fatal("expected differences")
	}
	if len(res.Deltas) != 1 {
		t.Fatalf("expected 1 delta, got %d", len(res.Deltas))
	}
	if res.Deltas[0].Type != diff.Changed {
		t.Errorf("expected Changed, got %s", res.Deltas[0].Type)
	}
}

func TestAddedEvents(t *testing.T) {
	a := makeSession("aaa", []string{"ls"})
	b := makeSession("bbb", []string{"ls", "pwd", "whoami"})
	res := diff.Compare(a, b)
	if len(res.Deltas) != 2 {
		t.Fatalf("expected 2 added deltas, got %d", len(res.Deltas))
	}
	for _, d := range res.Deltas {
		if d.Type != diff.Added {
			t.Errorf("expected Added, got %s", d.Type)
		}
	}
}

func TestRemovedEvents(t *testing.T) {
	a := makeSession("aaa", []string{"ls", "pwd", "whoami"})
	b := makeSession("bbb", []string{"ls"})
	res := diff.Compare(a, b)
	if len(res.Deltas) != 2 {
		t.Fatalf("expected 2 removed deltas, got %d", len(res.Deltas))
	}
	for _, d := range res.Deltas {
		if d.Type != diff.Removed {
			t.Errorf("expected Removed, got %s", d.Type)
		}
	}
}

func TestSummaryContainsIDs(t *testing.T) {
	a := makeSession("sess-1", []string{"ls"})
	b := makeSession("sess-2", []string{"pwd"})
	res := diff.Compare(a, b)
	sum := res.Summary()
	if !strings.Contains(sum, "sess-1") || !strings.Contains(sum, "sess-2") {
		t.Errorf("summary missing session IDs: %s", sum)
	}
}

func TestSummaryEqualSessions(t *testing.T) {
	a := makeSession("x", []string{"ls"})
	b := makeSession("y", []string{"ls"})
	res := diff.Compare(a, b)
	if !strings.Contains(res.Summary(), "identical") {
		t.Errorf("expected 'identical' in summary, got: %s", res.Summary())
	}
}
