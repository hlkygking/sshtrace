package audit_test

import (
	"bytes"
	"encoding/json"
	"strings"
	"testing"
	"time"

	"github.com/sshtrace/sshtrace/internal/audit"
	"github.com/sshtrace/sshtrace/internal/session"
)

func makeSession() *session.Session {
	s := session.New("alice", "10.0.0.1")
	return s
}

func TestLogWritesJSON(t *testing.T) {
	var buf bytes.Buffer
	l := audit.New(&buf)
	s := makeSession()

	if err := l.LogSession(s, audit.EventSessionStarted); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	var r audit.Record
	if err := json.Unmarshal(buf.Bytes(), &r); err != nil {
		t.Fatalf("output is not valid JSON: %v", err)
	}

	if r.Kind != audit.EventSessionStarted {
		t.Errorf("expected kind %q, got %q", audit.EventSessionStarted, r.Kind)
	}
	if r.User != "alice" {
		t.Errorf("expected user alice, got %q", r.User)
	}
	if r.RemoteIP != "10.0.0.1" {
		t.Errorf("expected remote_ip 10.0.0.1, got %q", r.RemoteIP)
	}
	if r.SessionID == "" {
		t.Error("session_id should not be empty")
	}
}

func TestLogIncludesDetail(t *testing.T) {
	var buf bytes.Buffer
	l := audit.New(&buf)
	s := makeSession()

	if err := l.Log(s, audit.EventCommandRun, "ls -la"); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	var r audit.Record
	_ = json.Unmarshal(buf.Bytes(), &r)

	if r.Detail != "ls -la" {
		t.Errorf("expected detail %q, got %q", "ls -la", r.Detail)
	}
}

func TestLogTimestampIsRecent(t *testing.T) {
	var buf bytes.Buffer
	l := audit.New(&buf)
	s := makeSession()
	before := time.Now().UTC()
	_ = l.LogSession(s, audit.EventSessionClosed)
	after := time.Now().UTC()

	var r audit.Record
	_ = json.Unmarshal(buf.Bytes(), &r)

	if r.Timestamp.Before(before) || r.Timestamp.After(after) {
		t.Errorf("timestamp %v not in expected range [%v, %v]", r.Timestamp, before, after)
	}
}

func TestLogNewlineDelimited(t *testing.T) {
	var buf bytes.Buffer
	l := audit.New(&buf)
	s := makeSession()

	_ = l.LogSession(s, audit.EventSessionStarted)
	_ = l.LogSession(s, audit.EventSessionClosed)

	lines := strings.Split(strings.TrimSpace(buf.String()), "\n")
	if len(lines) != 2 {
		t.Errorf("expected 2 lines, got %d", len(lines))
	}
	for i, line := range lines {
		var r audit.Record
		if err := json.Unmarshal([]byte(line), &r); err != nil {
			t.Errorf("line %d is not valid JSON: %v", i, err)
		}
	}
}

func TestAllEventKinds(t *testing.T) {
	kinds := []audit.EventKind{
		audit.EventSessionStarted,
		audit.EventSessionClosed,
		audit.EventCommandRun,
		audit.EventAlertTriggered,
	}
	for _, k := range kinds {
		var buf bytes.Buffer
		l := audit.New(&buf)
		s := makeSession()
		if err := l.LogSession(s, k); err != nil {
			t.Errorf("kind %q: unexpected error: %v", k, err)
		}
	}
}
