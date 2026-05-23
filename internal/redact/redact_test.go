package redact

import (
	"strings"
	"testing"
	"time"

	"sshtrace/internal/session"
)

func makeEvent(data string) session.Event {
	return session.Event{Type: "output", Data: data, Timestamp: time.Now()}
}

func TestNewDefaultFields(t *testing.T) {
	r, err := New(nil)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(r.fields) == 0 {
		t.Error("expected default fields to be set")
	}
}

func TestNewInvalidPattern(t *testing.T) {
	// Passing a valid custom field should succeed.
	_, err := New([]string{"myfield"})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
}

func TestRedactPassword(t *testing.T) {
	r, _ := New(nil)
	e := makeEvent("login password=supersecret done")
	out := r.RedactEvent(e)
	if strings.Contains(out.Data, "supersecret") {
		t.Errorf("expected password to be redacted, got: %s", out.Data)
	}
	if !strings.Contains(out.Data, "[REDACTED]") {
		t.Errorf("expected [REDACTED] marker, got: %s", out.Data)
	}
}

func TestRedactToken(t *testing.T) {
	r, _ := New(nil)
	e := makeEvent("export token=abc123xyz")
	out := r.RedactEvent(e)
	if strings.Contains(out.Data, "abc123xyz") {
		t.Errorf("expected token to be redacted, got: %s", out.Data)
	}
}

func TestNoSensitiveDataUnchanged(t *testing.T) {
	r, _ := New(nil)
	original := "ls -la /home/user"
	e := makeEvent(original)
	out := r.RedactEvent(e)
	if out.Data != original {
		t.Errorf("expected data unchanged, got: %s", out.Data)
	}
}

func TestRedactSession(t *testing.T) {
	r, _ := New(nil)
	s := &session.Session{
		ID:       "sess-1",
		User:     "alice",
		SourceIP: "10.0.0.1",
		Events: []session.Event{
			makeEvent("password=hunter2"),
			makeEvent("echo hello"),
		},
	}
	out := r.RedactSession(s)
	if out.ID != s.ID || out.User != s.User {
		t.Error("session metadata should be preserved")
	}
	if strings.Contains(out.Events[0].Data, "hunter2") {
		t.Error("expected password redacted in session")
	}
	if out.Events[1].Data != "echo hello" {
		t.Error("non-sensitive event should be unchanged")
	}
}

func TestCaseInsensitiveRedaction(t *testing.T) {
	r, _ := New(nil)
	e := makeEvent("PASSWORD=MyPass123")
	out := r.RedactEvent(e)
	if strings.Contains(out.Data, "MyPass123") {
		t.Errorf("expected case-insensitive redaction, got: %s", out.Data)
	}
}
