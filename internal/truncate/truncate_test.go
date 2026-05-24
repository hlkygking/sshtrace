package truncate_test

import (
	"strings"
	"testing"
	"time"

	"github.com/sshtrace/sshtrace/internal/session"
	"github.com/sshtrace/sshtrace/internal/truncate"
)

func makeSession(events []string) *session.Session {
	s := session.New("user1", "127.0.0.1")
	for _, data := range events {
		s.AddEvent(session.Event{
			Kind:      session.EventOutput,
			Data:      data,
			Timestamp: time.Now(),
		})
	}
	return s
}

func TestNewInvalidMaxBytes(t *testing.T) {
	_, err := truncate.New(0, 10)
	if err == nil {
		t.Fatal("expected error for maxBytes=0")
	}
}

func TestNewInvalidMaxEvents(t *testing.T) {
	_, err := truncate.New(100, 0)
	if err == nil {
		t.Fatal("expected error for maxEvents=0")
	}
}

func TestApplyNoTruncationNeeded(t *testing.T) {
	tr, _ := truncate.New(100, 50)
	s := makeSession([]string{"hello", "world"})
	tr.Apply(s)
	if len(s.Events) != 2 {
		t.Fatalf("expected 2 events, got %d", len(s.Events))
	}
	if s.Events[0].Data != "hello" {
		t.Errorf("unexpected data: %q", s.Events[0].Data)
	}
}

func TestApplyTruncatesLongData(t *testing.T) {
	tr, _ := truncate.New(20, 100)
	long := strings.Repeat("a", 100)
	s := makeSession([]string{long})
	tr.Apply(s)
	data := s.Events[0].Data
	if len(data) > 20 {
		t.Errorf("data length %d exceeds maxBytes 20", len(data))
	}
	if !strings.HasSuffix(data, "...[truncated]") {
		t.Errorf("expected truncation marker, got %q", data)
	}
}

func TestApplyDropsExcessEvents(t *testing.T) {
	tr, _ := truncate.New(1024, 3)
	events := []string{"a", "b", "c", "d", "e"}
	s := makeSession(events)
	tr.Apply(s)
	if len(s.Events) != 3 {
		t.Fatalf("expected 3 events, got %d", len(s.Events))
	}
	// First three should be retained.
	if s.Events[2].Data != "c" {
		t.Errorf("expected third event 'c', got %q", s.Events[2].Data)
	}
}

func TestApplyExactBoundary(t *testing.T) {
	tr, _ := truncate.New(5, 10)
	s := makeSession([]string{"hello"}) // exactly 5 bytes
	tr.Apply(s)
	if s.Events[0].Data != "hello" {
		t.Errorf("data should be unchanged at exact boundary, got %q", s.Events[0].Data)
	}
}
