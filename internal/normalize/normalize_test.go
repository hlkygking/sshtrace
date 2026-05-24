package normalize_test

import (
	"testing"
	"time"

	"github.com/sshtrace/sshtrace/internal/normalize"
	"github.com/sshtrace/sshtrace/internal/session"
)

func makeSession(user string) *session.Session {
	s := session.New(user, "192.168.1.1")
	return s
}

func addEvent(s *session.Session, data string) {
	s.AddEvent(session.Event{
		Timestamp: time.Now(),
		Kind:      "output",
		Data:      data,
	})
}

func TestDefaultOptions(t *testing.T) {
	opts := normalize.DefaultOptions()
	if !opts.TrimWhitespace {
		t.Error("expected TrimWhitespace to be true by default")
	}
	if !opts.NormalizeLineEndings {
		t.Error("expected NormalizeLineEndings to be true by default")
	}
}

func TestApplyNilSessionReturnsError(t *testing.T) {
	n := normalize.New(normalize.DefaultOptions())
	_, err := n.Apply(nil)
	if err == nil {
		t.Fatal("expected error for nil session")
	}
}

func TestTrimWhitespace(t *testing.T) {
	s := makeSession("alice")
	addEvent(s, "ls -la   ")

	n := normalize.New(normalize.Options{TrimWhitespace: true})
	out, err := n.Apply(s)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if out.Events[0].Data != "ls -la" {
		t.Errorf("expected trimmed data, got %q", out.Events[0].Data)
	}
}

func TestNormalizeLineEndings(t *testing.T) {
	s := makeSession("bob")
	addEvent(s, "echo hello\r\necho world\r")

	n := normalize.New(normalize.Options{NormalizeLineEndings: true})
	out, err := n.Apply(s)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	expected := "echo hello\necho world\n"
	if out.Events[0].Data != expected {
		t.Errorf("expected %q, got %q", expected, out.Events[0].Data)
	}
}

func TestLowercaseUser(t *testing.T) {
	s := makeSession("ALICE")

	n := normalize.New(normalize.Options{LowercaseUser: true})
	out, err := n.Apply(s)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if out.User != "alice" {
		t.Errorf("expected lowercase user, got %q", out.User)
	}
}

func TestNoOptionsIsNoop(t *testing.T) {
	s := makeSession("carol")
	addEvent(s, "pwd   ")

	n := normalize.New(normalize.Options{})
	out, err := n.Apply(s)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if out.Events[0].Data != "pwd   " {
		t.Errorf("expected unmodified data, got %q", out.Events[0].Data)
	}
}
