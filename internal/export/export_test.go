package export_test

import (
	"bytes"
	"encoding/json"
	"strings"
	"testing"
	"time"

	"sshtrace/internal/export"
	"sshtrace/internal/session"
)

func makeSession() *session.Session {
	s := session.New("alice", "192.168.1.10")
	s.AddEvent(session.EventOutput, "hello world")
	s.AddEvent(session.EventInput, "ls -la")
	return s
}

func TestExportJSON(t *testing.T) {
	s := makeSession()
	var buf bytes.Buffer
	ex := export.New(&buf, export.FormatJSON)
	if err := ex.Export(s); err != nil {
		t.Fatalf("Export JSON error: %v", err)
	}
	var decoded session.Session
	if err := json.Unmarshal(buf.Bytes(), &decoded); err != nil {
		t.Fatalf("invalid JSON output: %v", err)
	}
	if decoded.User != "alice" {
		t.Errorf("expected user alice, got %s", decoded.User)
	}
	if len(decoded.Events) != 2 {
		t.Errorf("expected 2 events, got %d", len(decoded.Events))
	}
}

func TestExportText(t *testing.T) {
	s := makeSession()
	var buf bytes.Buffer
	ex := export.New(&buf, export.FormatText)
	if err := ex.Export(s); err != nil {
		t.Fatalf("Export text error: %v", err)
	}
	out := buf.String()
	if !strings.Contains(out, "alice") {
		t.Error("expected user alice in text output")
	}
	if !strings.Contains(out, "192.168.1.10") {
		t.Error("expected remote IP in text output")
	}
	if !strings.Contains(out, "ls -la") {
		t.Error("expected event data in text output")
	}
}

func TestExportTextEventCount(t *testing.T) {
	s := session.New("bob", "10.0.0.1")
	for i := 0; i < 5; i++ {
		s.AddEvent(session.EventOutput, "data")
	}
	var buf bytes.Buffer
	ex := export.New(&buf, export.FormatText)
	if err := ex.Export(s); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !strings.Contains(buf.String(), "Events     : 5") {
		t.Error("expected event count 5 in output")
	}
}

func TestExportUnsupportedFormat(t *testing.T) {
	s := makeSession()
	var buf bytes.Buffer
	ex := export.New(&buf, export.Format("xml"))
	if err := ex.Export(s); err == nil {
		t.Error("expected error for unsupported format")
	}
}

func TestExportJSONTimestamp(t *testing.T) {
	s := makeSession()
	var buf bytes.Buffer
	ex := export.New(&buf, export.FormatJSON)
	_ = ex.Export(s)
	var decoded session.Session
	_ = json.Unmarshal(buf.Bytes(), &decoded)
	if decoded.StartedAt.IsZero() {
		t.Error("expected non-zero StartedAt timestamp")
	}
	if decoded.StartedAt.After(time.Now().Add(time.Second)) {
		t.Error("StartedAt is in the future")
	}
}
