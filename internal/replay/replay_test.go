package replay_test

import (
	"bytes"
	"strings"
	"testing"
	"time"

	"sshtrace/internal/replay"
	"sshtrace/internal/session"
)

func makeSession(commands []string) *session.Session {
	s := session.New("user", "127.0.0.1")
	base := time.Now()
	for i, cmd := range commands {
		s.AddEvent(session.Event{
			Timestamp: base.Add(time.Duration(i) * 100 * time.Millisecond),
			Data:      cmd,
		})
	}
	return s
}

func TestReplayInstant(t *testing.T) {
	s := makeSession([]string{"ls\n", "pwd\n", "exit\n"})
	var buf bytes.Buffer
	r := replay.New(replay.Options{Speed: 0, Writer: &buf})
	if err := r.Replay(s); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	got := buf.String()
	if !strings.Contains(got, "ls") || !strings.Contains(got, "exit") {
		t.Errorf("unexpected output: %q", got)
	}
}

func TestReplayEmptySession(t *testing.T) {
	s := session.New("user", "127.0.0.1")
	var buf bytes.Buffer
	r := replay.New(replay.Options{Speed: 0, Writer: &buf})
	if err := r.Replay(s); err != nil {
		t.Fatalf("unexpected error on empty session: %v", err)
	}
	if buf.Len() != 0 {
		t.Errorf("expected empty output, got %q", buf.String())
	}
}

func TestReplayOrder(t *testing.T) {
	s := makeSession([]string{"a", "b", "c"})
	var buf bytes.Buffer
	r := replay.New(replay.Options{Speed: 0, Writer: &buf})
	_ = r.Replay(s)
	if buf.String() != "abc" {
		t.Errorf("expected abc, got %q", buf.String())
	}
}

func TestReplayDefaultSpeed(t *testing.T) {
	// Speed 0 passed to New should default to 1.0 internally without panic.
	r := replay.New(replay.Options{})
	if r == nil {
		t.Fatal("expected non-nil replayer")
	}
}
